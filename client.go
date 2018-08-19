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

package clinote

import (
	"os"
	"path/filepath"
)

func NewClient(config Configuration, store Storager, opts ClientOption) *Client {
	c := &Client{
		Config: config,
		Store:  store,
	}
	return c
}

type ClientOption int32

const (
	DefaultClientOptions ClientOption = 0
	EvernoteSandbox                   = 1 << iota
	MemoryBasedCacheFile
)

type Client struct {
	Config       Configuration
	Store        Storager
	NoteStore    NotestoreClient
	Editor       Editer
	newCacheFile func(c *Client, filename string) (CacheFile, error)
}

func (c *Client) NewCacheFile(filename string) (CacheFile, error) {
	return c.newCacheFile(c, filename)
}

func newFileCacheFile(c *Client, filename string) (CacheFile, error) {
	fp := filepath.Join(c.Config.GetCacheFolder(), filename)
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	cachefile := &FileCacheFile{file: f, fp: fp}
	return cachefile, nil
}

func (c *Client) Edit(file CacheFile) error {
	return c.Editor.Edit(file)
}
