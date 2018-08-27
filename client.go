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
	"bytes"
	"os"
	"path/filepath"
	"syscall"
)

// NewClient creates a new Client struct.
func NewClient(config Configuration, store Storager, ns NotestoreClient, opts ClientOption) *Client {
	c := &Client{
		Config:     config,
		Store:      store,
		NoteStore:  ns,
		clientOpts: opts,
	}
	if opts&MemoryBasedCacheFile != 0 {
		c.newCacheFile = newMemoryCacheFile
	} else {
		c.newCacheFile = newFileCacheFile
	}
	if opts&VimEditer != 0 {
		c.Editor = new(VimEditor)
	} else {
		c.Editor = new(EnvEditor)
	}
	return c
}

// ClientOption is an bit mask of options for the client.
type ClientOption int32

const (
	// DefaultClientOptions creates a client with default options.
	DefaultClientOptions ClientOption = 0
	// EvernoteSandbox tells the client to use Evernote's sandbox instead of production.
	EvernoteSandbox = 1 << iota
	// MemoryBasedCacheFile for using a ram based file for editing.
	MemoryBasedCacheFile
	// VimEditer for using Vim as the editor.
	VimEditer
)

// Client is a client for all note operations.
type Client struct {
	// Config is the systems configuration
	Config Configuration
	// Store is the storage client
	Store Storager
	// Notestore is a client to interact with the note store.
	NoteStore NotestoreClient
	// Editor is the editor.
	Editor       Editer
	newCacheFile func(c *Client, filename string) (CacheFile, error)
	clientOpts   ClientOption
}

// NewCacheFile creates a new cache file for editing.
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

func newMemoryCacheFile(c *Client, filename string) (CacheFile, error) {
	fp := filepath.Join(c.Config.GetCacheFolder(), filename)
	err := syscall.Mkfifo(fp, 0600)
	if err != nil {
		return nil, err
	}
	cachefile := &MemoryCacheFile{buf: new(bytes.Buffer), pipePath: fp}
	return cachefile, nil
}

// Edit edits the cache file using the client's editor.
func (c *Client) Edit(file CacheFile) error {
	return c.Editor.Edit(file)
}
