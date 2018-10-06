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

package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/TcM1911/clinote"
	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	assert := assert.New(t)
	db, tmpDir := setupTestDB(t)
	expected := &clinote.Settings{APIKey: "test session"}

	t.Run("Store", func(t *testing.T) {
		err := db.StoreSettings(expected)
		assert.NoError(err, "Should not fail when storing settings")
	})

	t.Run("Get", func(t *testing.T) {
		actual, err := db.GetSettings()
		assert.NoError(err, "Should not return an error")
		assert.Equal(expected, actual, "Wrong settings returned")
	})

	t.Run("Handle_no_bucket", func(t *testing.T) {
		db, tmpDir := setupTestDB(t)
		expected := new(clinote.Settings)
		actual, err := db.GetSettings()
		assert.NoError(err, "Should create a bucket without problems")
		assert.Equal(expected, actual, "Should return an empty settings")
		db.Close()
		os.RemoveAll(tmpDir)
	})
	// Cleanup
	db.Close()
	os.RemoveAll(tmpDir)
}

func TestNotebookCaching(t *testing.T) {
	assert := assert.New(t)
	db, tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)
	defer db.Close()
	books := make([]*clinote.Notebook, 3)
	expected := clinote.NewNotebookCacheList(books)

	t.Run("Store", func(t *testing.T) {
		err := db.StoreNotebookList(expected)
		assert.NoError(err, "Should not fail when storing notebook cache")
	})

	t.Run("Get", func(t *testing.T) {
		actual, err := db.GetNotebookCache()
		assert.NoError(err, "Should not return an error")
		compareCacheList(assert, expected, actual)
	})
}

func TestSearchCaching(t *testing.T) {
	assert := assert.New(t)
	db, tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)
	defer db.Close()
	expected := []*clinote.Note{
		&clinote.Note{Title: "Note 1"},
		&clinote.Note{Title: "Note 2"},
		&clinote.Note{Title: "Note 3"},
	}

	t.Run("Store", func(t *testing.T) {
		err := db.SaveSearch(expected)
		assert.NoError(err, "Should not fail when storing notebook cache")
	})

	t.Run("Get", func(t *testing.T) {
		actual, err := db.GetSearch()
		assert.NoError(err, "Should not return an error")
		assert.Equal(expected, actual, "Wrong data returned from store")
	})
}

func TestRecoveryPoint(t *testing.T) {
	assert := assert.New(t)
	db, tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)
	defer db.Close()
	expectedNote := &clinote.Note{Title: "Test note"}

	t.Run("Store", func(t *testing.T) {
		err := db.SaveNoteRecoveryPoint(expectedNote)
		assert.NoError(err, "Should not fail to save")
	})

	t.Run("Get", func(t *testing.T) {
		actual, err := db.GetNoteRecoveryPoint()
		assert.NoError(err, "Should not fail to return recovery point")
		assert.Equal(expectedNote, actual, "Wrong note returned")
	})
}

func TestCredentialStore(t *testing.T) {
	assert := assert.New(t)
	db, tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)
	defer db.Close()
	expectedCredentials := []*clinote.Credential{
		&clinote.Credential{Name: "Cred1", Secret: "Sec1", CredType: clinote.EvernoteCredential},
		&clinote.Credential{Name: "Cred2", Secret: "Sec2", CredType: clinote.EvernoteSandboxCredential},
		&clinote.Credential{Name: "Cred3", Secret: "Sec3", CredType: clinote.EvernoteCredential},
	}

	t.Run("Add and getall", func(t *testing.T) {
		for _, cred := range expectedCredentials {
			err := db.Add(cred)
			assert.NoError(err, "Should not fail to save")
		}
		creds, err := db.GetAll()
		assert.NoError(err, "Should get all without an error")
		for i, cred := range creds {
			assert.Equal(*expectedCredentials[i], *cred)
		}
	})

	t.Run("Get by index", func(t *testing.T) {
		for i, expected := range expectedCredentials {
			cred, err := db.GetByIndex(i)
			assert.NoError(err, "Failed to get credential by index")
			assert.Equal(*expected, *cred, "Credential returned doesn't match")
		}
	})

	t.Run("Out of range index checks", func(t *testing.T) {
		for _, index := range []int{-1, len(expectedCredentials), len(expectedCredentials) + 1} {
			cred, err := db.GetByIndex(index)
			assert.Error(err, "Should return error")
			assert.Nil(cred, "Nil should be returned")
			assert.Equal(ErrIndexOutOfRange, err, "Wrong error returned")
		}
	})

	t.Run("Remove a credential", func(t *testing.T) {
		err := db.Remove(expectedCredentials[len(expectedCredentials)-1])
		assert.NoError(err, "Should not fail removing")
		creds, err := db.GetAll()
		if assert.NoError(err, "Failed to get all creds") {
			for i, cred := range creds {
				assert.NotEqual(*expectedCredentials[len(expectedCredentials)-1], *cred, "Should not match removed credential")
				assert.Equal(*expectedCredentials[i], *cred, "Should maintain the same location in the array")
			}
		}
	})

	t.Run("Do not remove no existing", func(t *testing.T) {
		err := db.Remove(new(clinote.Credential))
		assert.Error(err, "Should not fail removing")
		assert.Equal(clinote.ErrNoMatchingCredentialFound, err)
	})
}

func compareCacheList(assert *assert.Assertions, expected *clinote.NotebookCacheList, actual *clinote.NotebookCacheList) {
	assert.Equal(expected.Limit, actual.Limit)
	assert.Equal(expected.Notebooks, actual.Notebooks)
	assert.True(expected.Timestamp.Equal(actual.Timestamp), "Wrong timestamp")
}

func setupTestDB(t *testing.T) (*Database, string) {
	tmpDir, err := ioutil.TempDir("", "clinote-test")
	if err != nil {
		t.Fatalf("Problem with creating temp folder: %s\n", err)
	}

	db, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("No db: %s\n", err)
	}
	return db, tmpDir
}
