/*
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
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
 * Copyright (C) Joakim Kennedy, 2018
 */

package evernote

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/TcM1911/clinote"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)
	expectedSession := "test session string"

	// Setup cache folder
	cacheDir, err := ioutil.TempDir("", "clinote-test")
	if err != nil {
		assert.FailNow("Error when creating cache folder:" + err.Error())
	}
	defer os.Remove(cacheDir)
	fp := filepath.Join(cacheDir, "session")
	_, err = os.OpenFile(fp, os.O_CREATE, 0600)
	if err != nil {
		assert.FailNow("Error when creating session file:" + err.Error())
	}
	f, err := os.OpenFile(fp, os.O_WRONLY, 0600)
	if err != nil {
		assert.FailNow("Error when opening session file:" + err.Error())
	}
	_, err = f.Write([]byte(expectedSession))
	if err != nil {
		assert.FailNow("Error when writing to session file:" + err.Error())
	}
	f.Close()

	// Setup config folder
	configDir, err := ioutil.TempDir("", "clinote-test")
	if err != nil {
		assert.FailNow("Error when creating config folder:" + err.Error())
	}
	defer os.Remove(configDir)

	// Setup store mock
	settings := &clinote.Settings{}
	store := &mockStore{settings: settings}

	// Setup config mock
	cfg := &cfgMock{
		getConfFolder:  func() string { return configDir },
		getCacheFolder: func() string { return cacheDir },
		getStore:       func() clinote.Storager { return store },
	}

	// Tests
	t.Run("migration", func(t *testing.T) {
		client := NewClient(cfg)
		assert.NotNil(client, "Should return a client")
		assert.Equal(expectedSession, client.apiToken)
		// Check that the session was migrated
		assert.Equal(expectedSession, store.settings.APIKey)
		_, err = os.Stat(fp)
		assert.True(os.IsNotExist(err), "Session file not removed")
	})
	t.Run("session_from_storage", func(t *testing.T) {
		// Ensure session file doesn't exist.
		_, err = os.Stat(fp)
		assert.True(os.IsNotExist(err), "Session file not removed")
		s, _ := cfg.Store().GetSettings()
		assert.Equal(expectedSession, s.APIKey)

		client := NewClient(cfg)
		assert.NotNil(client)
		assert.Equal(expectedSession, client.apiToken)
	})
}

type mockStore struct {
	settings *clinote.Settings
}

func (m *mockStore) GetNotebookCache() (*clinote.NotebookCacheList, error) {
	panic("not implemented")
}

func (m *mockStore) StoreNotebookList(list *clinote.NotebookCacheList) error {
	panic("not implemented")
}

func (m *mockStore) Close() error {
	return nil
}

func (m *mockStore) GetSettings() (*clinote.Settings, error) {
	return m.settings, nil
}

func (m *mockStore) StoreSettings(s *clinote.Settings) error {
	m.settings.APIKey = s.APIKey
	return nil
}
