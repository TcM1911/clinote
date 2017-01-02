/*
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (C) Joakim Kennedy, 2016
 */

package user

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/mrjones/oauth"
	"github.com/tcm1911/clinote/config"
)

// AuthToken is the user's authentication token.
var AuthToken string

type callbackValues struct {
	TempToken  string
	Verifier   string
	SandboxLnb bool
}

// Logout removes the session stored.
func Logout() {
	fp := filepath.Join(config.GetCacheFolder(), "session")
	if !checkLogin(fp) {
		fmt.Println("You are not logged in.")
		return
	}
	if err := os.Remove(fp); err != nil {
		fmt.Println("Error when trying to remove the session file:", err)
		return
	}
	fmt.Println("Successfully logged out.")
}

func checkLogin(fp string) bool {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return false
	}
	return true
}

// Login logs the user in to the server.
func Login() bool {
	fp := filepath.Join(config.GetCacheFolder(), "session")
	if checkLogin(fp) {
		fmt.Println("You are already logged in.")
		return false
	}
	c := make(chan *callbackValues)
	http.HandleFunc("/", oathCallbackHandler(c))
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	defer listener.Close()
	fmt.Println("Starting callback listener on", listener.Addr().String())
	if err != nil {
		fmt.Println("Error when starting callback listener:", err)
		return false
	}
	go func() {
		err = http.Serve(listener, nil)
		if err != nil {
			fmt.Println("Error when starting listener web server:", err)
			os.Exit(1)
		}
	}()
	callbackURL := fmt.Sprintf("http://%s/", listener.Addr().String())
	tmpToken, loginURL, err := getTempToken(callbackURL)
	if err != nil {
		fmt.Println("Error when getting request temporary token:", err)
		return false
	}
	go tryOpenLoginInBrowser(loginURL)
	fmt.Println("Waiting for access...")
	callback := <-c
	if callback.TempToken != tmpToken.Token {
		fmt.Println("Temporary token doesn't match, aborting...")
		return false
	}
	if callback.Verifier == "" {
		fmt.Println("Access revoked.")
		return false
	}
	token, err := getAuthToken(tmpToken, callback.Verifier)
	if err != nil {
		fmt.Println("Error when retrieving auth token:", err)
		return false
	}
	if err = saveToken(token); err != nil {
		fmt.Println("Saving to file failed:", err)
		return false
	}
	return true
}

func saveToken(token string) error {
	dir := config.GetCacheFolder()
	fp := filepath.Join(dir, "session")
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		f, err := os.OpenFile(fp, os.O_CREATE, 0600)
		if err != nil {
			return errors.New("error when creating session file: " + err.Error())
		}
		f.Close()
	}
	f, err := os.OpenFile(fp, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(token); err != nil {
		return errors.New("error when saving token to file:" + err.Error())
	}
	return nil
}

func getAuthToken(tmpToken *oauth.RequestToken, verifier string) (string, error) {
	client := GetClient()
	token, err := client.GetAuthorizedToken(tmpToken, verifier)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func tryOpenLoginInBrowser(url string) {
	browser := os.Getenv("BROWSER")
	if browser == "" {
		fmt.Printf("Open %s in your browser to give CLInote access to Evernote.\n", url)
		return
	}
	cmd := exec.Command(browser, url)
	fmt.Printf("Opening %s in %s\n", url, browser)
	cmd.Run()
}

func getTempToken(callback string) (*oauth.RequestToken, string, error) {
	client := GetClient()
	token, url, err := client.GetRequestToken(callback)
	return token, url, err
}

func oathCallbackHandler(returnChan chan *callbackValues) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vals := new(callbackValues)
		requestVals := r.URL.Query()
		vals.TempToken = requestVals.Get("oauth_token")
		vals.Verifier = requestVals.Get("oauth_verifier")
		sandboxBool := requestVals.Get("sandbox_lnb")
		if sandboxBool != "" {
			sandbox, err := strconv.ParseBool(sandboxBool)
			if err != nil {
				fmt.Println("Error when parsing OAth callback request:", err)
			}
			vals.SandboxLnb = sandbox
		}
		w.Write([]byte("You can now close this tab"))
		returnChan <- vals
	}
}
