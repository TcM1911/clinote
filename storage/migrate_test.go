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
	"os"
	"testing"

	"github.com/TcM1911/clinote"
	"github.com/stretchr/testify/assert"
)

func TestCredentialMigration(t *testing.T) {
	assert := assert.New(t)
	db, tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)
	defer db.Close()

	t.Run("no migration on first run", func(t *testing.T) {
		err := migrate(db, uint64(0))
		assert.NoError(err)
		list, err := db.GetAll()
		assert.NoError(err, "Error when getting all credentials")
		assert.Empty(list, "Credential list should be empty")
	})

	t.Run("migrate API key to credential store", func(t *testing.T) {
		secret := "test token"
		settings := &clinote.Settings{APIKey: secret}
		err := db.StoreSettings(settings)
		assert.NoError(err, "Failed to setup database")

		err = migrate(db, uint64(0))
		assert.NoError(err)
		list, err := db.GetAll()
		assert.NoError(err, "Error when getting all credentials")
		assert.Len(list, 1, "Credential list should have one entry")
		assert.Equal("OAuth", list[0].Name)
		assert.Equal(secret, list[0].Secret)
		assert.Equal(clinote.EvernoteCredential, list[0].CredType)
	})
}
