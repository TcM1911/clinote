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
	tmpDir, err := ioutil.TempDir("", "clinote-test")
	if err != nil {
		t.Fatalf("Problem with creating temp folder: %s\n", err)
	}

	expected := &clinote.Settings{APIKey: "test session"}
	db, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("No db: %s\n", err)
	}

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
		tmpDir, err := ioutil.TempDir("", "clinote-test")
		if err != nil {
			t.Fatalf("Problem with creating temp folder: %s\n", err)
		}
		db, err := Open(tmpDir)
		if err != nil {
			t.Fatalf("No db: %s\n", err)
		}
		expected := new(clinote.Settings)
		actual, err := db.GetSettings()
		assert.NoError(err, "Should create a bucket without problems")
		assert.Equal(expected, actual, "Should return an empty settings")
		db.Close()
		os.Remove(tmpDir)
	})
	// Cleanup
	db.Close()
	os.Remove(tmpDir)
}
