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
	"testing"
	"time"

	"github.com/TcM1911/clinote"
	"github.com/mrjones/oauth"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	assert := assert.New(t)
	t.Run("Error when not logged in", func(t *testing.T) {
		cfg := new(cfgMock)
		store := new(mockStore)
		settings := new(clinote.Settings)
		store.settings = settings
		cfg.getStore = func() clinote.Storager { return store }
		err := Logout(cfg)
		assert.Equal(ErrNotLoggedIn, err, "Wrong error message")
	})
	t.Run("Should logout", func(t *testing.T) {
		cfg := new(cfgMock)
		store := new(mockStore)
		settings := new(clinote.Settings)
		settings.APIKey = "test session"
		store.settings = settings
		cfg.getStore = func() clinote.Storager { return store }
		err := Logout(cfg)
		assert.Nil(err, "Should not return an error. Returned:", err)
		assert.Equal("", settings.APIKey, "Session key should be empty")
	})
}

type cfgMock struct {
	getCacheFolder func() string
	getConfFolder  func() string
	getStore       func() clinote.Storager
}

func (c *cfgMock) GetConfigFolder() string {
	return c.getConfFolder()
}

func (c *cfgMock) GetCacheFolder() string {
	return c.getCacheFolder()
}

func (c *cfgMock) Store() clinote.Storager {
	return c.getStore()
}

func (c *cfgMock) Close() error {
	return nil
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
		settings := new(clinote.Settings)
		err := loginHelperFunction(t, settings, false, false)
		assert.NoError(err, "Should not return a login error")
		assert.Equal("oauth_token", settings.APIKey, "Session key not set")
	})
	t.Run("error when access revoked", func(t *testing.T) {
		settings := new(clinote.Settings)
		err := loginHelperFunction(t, settings, true, false)
		assert.EqualError(err, ErrAccessRevoked.Error(), "Expected access denied")
	})
	t.Run("error when temporary token is incorrect", func(t *testing.T) {
		settings := new(clinote.Settings)
		err := loginHelperFunction(t, settings, false, true)
		assert.EqualError(err, ErrTempTokenMismatch.Error(), "Expected token error")
	})
	t.Run("error when already logged in", func(t *testing.T) {
		settings := new(clinote.Settings)
		settings.APIKey = "test token"
		err := loginHelperFunction(t, settings, false, true)
		assert.EqualError(err, ErrAlreadyLoggedIn.Error(), "Expected already logged in error")
	})
}

func loginHelperFunction(t *testing.T, settings *clinote.Settings, verify, tokenMismatch bool) error {
	tmpToken := &oauth.RequestToken{Token: "testToken"}
	client := new(mockClient)
	cfg := new(cfgMock)
	store := new(mockStore)
	store.settings = settings
	cfg.getStore = func() clinote.Storager { return store }
	client.getConfig = func() clinote.Configuration { return cfg }
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
	getConfig          func() clinote.Configuration
	apiToken           string
	getNotestore       func() (clinote.NotestoreClient, error)
}

func (c *mockClient) GetNoteStore() (clinote.NotestoreClient, error) {
	return c.getNotestore()
}

func (c *mockClient) GetAuthorizedToken(tmpToken *oauth.RequestToken, verifier string) (token string, err error) {
	return c.getAuthorizedToken(tmpToken, verifier)
}

func (c *mockClient) GetRequestToken(callbackURL string) (token *oauth.RequestToken, url string, err error) {
	return c.getRequestToken(callbackURL)
}

func (c *mockClient) GetConfig() clinote.Configuration {
	return c.getConfig()
}

func (c *mockClient) GetAPIToken() string {
	return c.apiToken
}
