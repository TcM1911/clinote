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
	"encoding/json"

	"github.com/TcM1911/clinote"
	"github.com/boltdb/bolt"
)

func migrate(db *Database, currVersion uint64) error {
	if currVersion < uint64(1) {
		err := migrateOAuthCredential(db)
		if err != nil {
			return err
		}
	}
	return nil
}

func migrateOAuthCredential(db *Database) error {
	var data []byte
	d, err := db.getDBHandler()
	if err != nil {
		db.releaseDBHandler()
		return err
	}
	err = d.View(func(t *bolt.Tx) error {
		b := t.Bucket(settingsBucket)
		if b == nil {
			return errNoBucket
		}
		data = b.Get(settingsKey)
		return nil
	})
	if err == errNoBucket {
		db.releaseDBHandler()
		return nil
	}

	// Decode
	var settings struct {
		APIKey string
	}
	err = json.Unmarshal(data, &settings)
	if err != nil {
		db.releaseDBHandler()
		return err
	}
	db.releaseDBHandler()
	return db.Add(&clinote.Credential{Name: "OAuth", Secret: settings.APIKey, CredType: clinote.EvernoteCredential})
}
