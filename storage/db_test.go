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
		assert.Equal(expected, actual, "Wrong notebook cache returned")
	})
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
