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
 * Copyright (C) Joakim Kennedy
 , 2016
*/

package evernote

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/tcm1911/clinote/config"
	"github.com/tcm1911/evernote-sdk-golang/client"
	"github.com/tcm1911/evernote-sdk-golang/notestore"
)

var evernote *client.EvernoteClient
var apiConsumer = "clinote"
var apiSecret = "e9a3234ceefed62b"
var setup sync.Once
var devBuild = false

// GetNoteStore returns a notestore client for the user.
func GetNoteStore(cfg config.Configuration) *notestore.NoteStoreClient {
	setupClient(cfg)
	if AuthToken == "" {
		fmt.Println("No valid token. Please login again.")
		os.Exit(1)
	}
	nb, err := evernote.GetNoteStore(AuthToken)
	if err != nil {
		fmt.Println("Error when getting user's notestore:", err)
		os.Exit(1)
	}
	return nb
}

// GetClient returns the evernote client.
func GetClient(cfg config.Configuration) *client.EvernoteClient {
	setupClient(cfg)
	return evernote
}

func setupClient(cfg config.Configuration) {
	setup.Do(func() {
		env := client.PRODUCTION
		if devBuild {
			env = client.SANDBOX
		}
		evernote = client.NewClient(apiConsumer, apiSecret, env)
		devToken := os.Getenv("EVERNOTE_DEV_TOKEN")
		if devToken != "" {
			AuthToken = devToken
		} else {
			cacheDir := cfg.GetCacheFolder()
			fp := filepath.Join(cacheDir, "session")
			if _, err := os.Stat(fp); os.IsNotExist(err) {
				return
			}
			f, err := os.OpenFile(fp, os.O_RDONLY, 0600)
			if err != nil {
				panic(err.Error())
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err.Error())
			}
			AuthToken = string(b)
		}
	})
}
