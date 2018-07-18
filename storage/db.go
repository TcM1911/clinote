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
	"errors"
	"path/filepath"

	"github.com/TcM1911/clinote"
	"github.com/boltdb/bolt"
)

// List of buckets
var (
	settingsBucket = []byte("settings")
)

// List of keys
var (
	settingsKey = []byte("user_settings")
)

var (
	errNoBucket = errors.New("no bucket")
)

// Open returns an instance of the database.
func Open(cfgFolder string) (*Database, error) {
	b, err := bolt.Open(filepath.Join(cfgFolder, "clinote.db"), 0600, nil)
	if err != nil {
		return nil, err
	}
	return &Database{bolt: b}, nil
}

// Database is a representation of the backend storage.
type Database struct {
	bolt *bolt.DB
}

// GetSettings returns the settings from the database.
func (d *Database) GetSettings() (*clinote.Settings, error) {
	var data []byte
	err := d.bolt.View(func(t *bolt.Tx) error {
		b := t.Bucket(settingsBucket)
		if b == nil {
			return errNoBucket
		}
		data = b.Get(settingsKey)
		return nil
	})
	if err == errNoBucket {
		s := new(clinote.Settings)
		err := d.bolt.Update(func(t *bolt.Tx) error {
			b, err := t.CreateBucket(settingsBucket)
			if err != nil {
				return nil
			}
			data, err := json.Marshal(s)
			if err != nil {
				return err
			}
			return b.Put(settingsKey, data)
		})
		return s, err
	}
	if err != nil {
		return nil, err
	}
	var settings clinote.Settings
	err = json.Unmarshal(data, &settings)
	return &settings, err
}

// StoreSettings saves the settings to the database.
func (d *Database) StoreSettings(settings *clinote.Settings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return d.bolt.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists(settingsBucket)
		if err != nil {
			return err
		}
		return b.Put(settingsKey, data)
	})
}

// Close shuts down the connection to the database.
func (d *Database) Close() error {
	return d.bolt.Close()
}
