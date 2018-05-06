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
 * Copyright (C) Joakim Kennedy, 2016-2017
 */

package evernote

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/TcM1911/clinote/config"
	"github.com/mrjones/oauth"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	assert := assert.New(t)
	t.Run("Error when not logged in", func(t *testing.T) {
		cfg := new(cfgMock)
		tmpdir := os.TempDir()
		folder := filepath.Join(tmpdir, "logout_test_dir1")
		if err := os.MkdirAll(folder, os.ModeDir|0700); err != nil {
			t.Fatal(err)
		}
		cfg.getCacheFolder = func() string { return folder }
		err := Logout(cfg)
		assert.Equal(ErrNotLoggedIn, err, "Wrong error message")
	})
	t.Run("Should logout", func(t *testing.T) {
		tmpdir := os.TempDir()
		folder := filepath.Join(tmpdir, "logout_test_dir2")
		if err := os.MkdirAll(folder, os.ModeDir|0700); err != nil {
			t.Fatal(err)
		}
		fp := filepath.Join(folder, "session")
		f, err := os.OpenFile(fp, os.O_CREATE, 0600)
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
		cfg := new(cfgMock)
		cfg.getCacheFolder = func() string { return folder }
		err = Logout(cfg)
		assert.Nil(err, "Should not return an error. Returned:", err)
		// Clean up
		err = os.RemoveAll(folder)
		if err != nil {
			t.Log("Error when removing test folder:", err.Error())
		}
	})
}

type cfgMock struct {
	getCacheFolder func() string
	getConfFolder  func() string
}

func (c *cfgMock) GetConfigFolder() string {
	return c.getConfFolder()
}

func (c *cfgMock) GetCacheFolder() string {
	return c.getCacheFolder()
}

func TestCallbackHandler(t *testing.T) {
	assert := assert.New(t)
	tempToken := "internal-dev.14CD91FCE1F.687474703A2F2F6C6F63616C686F7374.6E287AD298969B6F8C0B4B1D67BCAB1D"
	verifier := "40793F8BAE15D4E3B6DD5CA8AB4BF62F"
	sandbox := "false"

	c := make(chan *callbackValues)
	url := fmt.Sprintf("http://www.sample.com/?oauth_token=%s&&oauth_verifier=%s&&sandbox_lnb=%s", tempToken, verifier, sandbox)
	r := httptest.NewRequest(http.MethodGet, url, nil)
	w := new(httptest.ResponseRecorder)
	go oathCallbackHandler(c).ServeHTTP(w, r)
	vals := <-c

	assert.Equal(verifier, vals.Verifier)
	assert.Equal(tempToken, vals.TempToken)
	assert.False(vals.SandboxLnb)
}

func TestLogin(t *testing.T) {
	os.Unsetenv("BROWSER")
	assert := assert.New(t)
	t.Run("should login", func(t *testing.T) {
		testDir := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().Unix()))
		err := os.MkdirAll(testDir, os.ModeDir|0777)
		if err != nil {
			assert.FailNow("Error creating test folder", err)
		}
		err = loginHelperFunction(t, testDir, false, false)
		assert.NoError(err, "Should not return a login error")
		err = os.RemoveAll(testDir)
		if err != nil {
			assert.FailNow("Error when removing test folder", err)
		}
	})
	t.Run("error when access revoked", func(t *testing.T) {
		testDir := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().Unix()))
		err := os.MkdirAll(testDir, os.ModeDir|0777)
		if err != nil {
			assert.FailNow("Error creating test folder", err)
		}
		err = loginHelperFunction(t, testDir, true, false)
		assert.EqualError(err, ErrAccessRevoked.Error(), "Expected access denied")
		err = os.RemoveAll(testDir)
		if err != nil {
			assert.FailNow("Error when removing test folder", err)
		}
	})
	t.Run("error when temporary token is incorrect", func(t *testing.T) {
		testDir := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().Unix()))
		err := os.MkdirAll(testDir, os.ModeDir|0777)
		if err != nil {
			assert.FailNow("Error creating test folder", err)
		}
		err = loginHelperFunction(t, testDir, false, true)
		assert.EqualError(err, ErrTempTokenMismatch.Error(), "Expected token error")
		err = os.RemoveAll(testDir)
		if err != nil {
			assert.FailNow("Error when removing test folder", err)
		}
	})
	t.Run("error when already logged in", func(t *testing.T) {
		testDir := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().Unix()))
		err := os.MkdirAll(testDir, os.ModeDir|0777)
		if err != nil {
			assert.FailNow("Error creating test folder", err)
		}
		file := filepath.Join(testDir, "session")
		os.OpenFile(file, os.O_CREATE, 0777)
		err = loginHelperFunction(t, testDir, false, true)
		assert.EqualError(err, ErrAlreadyLoggedIn.Error(), "Expected already logged in error")
		err = os.RemoveAll(testDir)
		if err != nil {
			assert.FailNow("Error when removing test folder", err)
		}
	})
}

func loginHelperFunction(t *testing.T, testFolder string, verify, tokenMismatch bool) error {
	tmpToken := &oauth.RequestToken{Token: "testToken"}
	client := new(mockClient)
	cfg := new(cfgMock)
	cfg.getCacheFolder = func() string { return testFolder }
	client.getConfig = func() config.Configuration { return cfg }
	client.getAuthorizedToken = func(tmpToken *oauth.RequestToken, verify string) (string, error) {
		return "oauth_token", nil
	}
	ch := make(chan string)
	client.getRequestToken = func(callback string) (*oauth.RequestToken, string, error) {
		time.Sleep(1 * time.Second)
		ch <- callback
		return tmpToken, "http://test", nil
	}
	go func() {
		callback := <-ch
		verifier := "verifier"
		if verify {
			verifier = ""
		}
		url := fmt.Sprintf("%s?oauth_token=%s&&oauth_verifier=%s&&sandbox_lnb=false", callback, tmpToken.Token, verifier)
		if tokenMismatch {
			url = fmt.Sprintf("%s?oauth_token=%s&&oauth_verifier=%s&&sandbox_lnb=false", callback, "mismatch", verifier)
		}
		fmt.Println("Sending request to", url)
		r, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		r.Body.Close()
	}()
	err := Login(client)
	return err
}

type mockClient struct {
	getAuthorizedToken func(*oauth.RequestToken, string) (string, error)
	getRequestToken    func(string) (*oauth.RequestToken, string, error)
	getConfig          func() config.Configuration
	apiToken           string
	getNotestore       func() (NotestoreClient, error)
}

func (c *mockClient) GetNoteStore() (NotestoreClient, error) {
	return c.getNotestore()
}

func (c *mockClient) GetAuthorizedToken(tmpToken *oauth.RequestToken, verifier string) (token string, err error) {
	return c.getAuthorizedToken(tmpToken, verifier)
}

func (c *mockClient) GetRequestToken(callbackURL string) (token *oauth.RequestToken, url string, err error) {
	return c.getRequestToken(callbackURL)
}

func (c *mockClient) GetConfig() config.Configuration {
	return c.getConfig()
}

func (c *mockClient) GetAPIToken() string {
	return c.apiToken
}
