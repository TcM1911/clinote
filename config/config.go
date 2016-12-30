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
 * Copyright (C) Joakim Kennedy, 2016
 */

package config

import (
	"fmt"
	"os"
)

// GetConfigFolder returns the folder used to store configurations.
func GetConfigFolder() string {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// Create folder
		if err = os.MkdirAll(configDir, os.ModeDir|0700); err != nil {
			fmt.Println("Error when creating config folder:", err)
			return ""
		}
	}
	return configDir
}

// GetCacheFolder returns the folder used to cache.
func GetCacheFolder() string {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		// Create cache folder.
		if err = os.MkdirAll(cacheDir, os.ModeDir|0700); err != nil {
			fmt.Println("Error when creating cache folder:", err)
			return ""
		}
	}
	return cacheDir
}
