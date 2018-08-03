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
	cacheBucket    = []byte("cache")
)

// List of keys
var (
	settingsKey      = []byte("user_settings")
	notebookCacheKey = []byte("notebook_cache")
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

func (d *Database) getData(bucket, key []byte) ([]byte, error) {
	var data []byte
	err := d.bolt.View(func(t *bolt.Tx) error {
		b := t.Bucket(bucket)
		if b == nil {
			return errNoBucket
		}
		data = b.Get(key)
		return nil
	})
	if err == errNoBucket {
		err := d.bolt.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucket(bucket)
			if err != nil {
				return err
			}
			return nil
		})
		return data, err
	}
	return data, err
}

func (d *Database) storeData(bucket, key, data []byte) error {
	return d.bolt.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return b.Put(key, data)
	})
}

// GetSettings returns the settings from the database.
func (d *Database) GetSettings() (*clinote.Settings, error) {
	var settings clinote.Settings
	data, err := d.getData(settingsBucket, settingsKey)
	if err == nil && data != nil {
		err = json.Unmarshal(data, &settings)
	}
	return &settings, err
}

// StoreSettings saves the settings to the database.
func (d *Database) StoreSettings(settings *clinote.Settings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return d.storeData(settingsBucket, settingsKey, data)
}

// GetNotebookCache returns the stored NotebookCacheList.
func (d *Database) GetNotebookCache() (*clinote.NotebookCacheList, error) {
	var list clinote.NotebookCacheList
	data, err := d.getData(cacheBucket, notebookCacheKey)
	if err == nil && data != nil {
		err = json.Unmarshal(data, &list)
	}
	return &list, err
}

// StoreNotebookList saves the list to the database.
func (d *Database) StoreNotebookList(list *clinote.NotebookCacheList) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return d.storeData(cacheBucket, notebookCacheKey, data)
}

// Close shuts down the connection to the database.
func (d *Database) Close() error {
	return d.bolt.Close()
}
