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
 * Copyright (C) Joakim Kennedy, 2016, 2018
 */

package evernote

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TcM1911/clinote/config"
)

func migrateOldSession(cfg config.Configuration) string {
	cacheDir := cfg.GetCacheFolder()
	fp := filepath.Join(cacheDir, "session")
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return ""
	}
	f, err := os.OpenFile(fp, os.O_RDONLY, 0600)
	if err != nil {
		panic(err.Error())
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}

	f.Close()
	err = os.Remove(fp)
	if err != nil {
		panic(err.Error())
	}
	apiKey := string(b)
	settings, err := cfg.Store().GetSettings()
	if err != nil {
		panic(err.Error())
	}
	settings.APIKey = apiKey
	if err = cfg.Store().StoreSettings(settings); err != nil {
		panic(err.Error())
	}
	return apiKey
}
