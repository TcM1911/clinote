// +build linux,!windows,!darwin

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
	"os"
	"path/filepath"
)

var (
	// configDir is the folder used to store the configuration files.
	configDir string
	// cacheDir is the folder used to store notes during edits.
	cacheDir string
)

func init() {
	xdgConf := os.Getenv("XDG_CONFIG_HOME")
	if xdgConf == "" {
		xdgConf = os.Getenv("HOME")
		configDir = filepath.Join(xdgConf, ".config", "clinote")
	} else {
		configDir = filepath.Join(xdgConf, "clinote")
	}
	if xdgConf == "" {
		panic("can't locate user's xdg config folder.")
	}

	xdgCache := os.Getenv("XDG_CACHE_HOME")
	if xdgCache == "" {
		xdgCache = os.Getenv("HOME")
		cacheDir = filepath.Join(xdgCache, ".cache", "clinote")
	} else {
		cacheDir = filepath.Join(xdgCache, "clinote")
	}
	if xdgCache == "" {
		panic("can't locate user's xdg cache folder.")
	}
}
