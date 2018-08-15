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
	"io"
	"os"
	"time"
)

var (
	// DefaultNotebookCacheTime is the default time limit for when the
	// list is considered outdated.
	DefaultNotebookCacheTime = 24 * time.Hour
)

// NewNotebookCacheListWithLimit creates a new cache list with the given expiration limit.
func NewNotebookCacheListWithLimit(notebooks []*Notebook, limit time.Duration) *NotebookCacheList {
	return &NotebookCacheList{
		Notebooks: notebooks,
		Limit:     limit,
		Timestamp: time.Now(),
	}
}

// NewNotebookCacheList creates a cache list with the default expiration limit.
func NewNotebookCacheList(notebooks []*Notebook) *NotebookCacheList {
	return NewNotebookCacheListWithLimit(notebooks, DefaultNotebookCacheTime)
}

// NotebookCacheList is a list of cached notebooks.
type NotebookCacheList struct {
	// Notebooks is the list of notebooks.
	Notebooks []*Notebook
	// Timestamp of when the list was created.
	Timestamp time.Time
	// Limit is the until the list outdated.
	Limit time.Duration
}

// IsOutdated returns true if the list has expired.
func (n *NotebookCacheList) IsOutdated() bool {
	return time.Since(n.Timestamp) > n.Limit
}

// CacheFile has the note content written and the user
// edits the content in the CacheFile to update the note's
// content.
type CacheFile interface {
	io.ReadWriteCloser
	FilePath() string
	ReOpen() error
	CloseAndRemove() error
}

// FileCacheFile implements the CacheFile interface and uses
// a temporary file for storing the data on disk.
type FileCacheFile struct {
	file *os.File
	fp   string
}

// Read returns content from the file.
func (f *FileCacheFile) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

// Write adds content to the file.
func (f *FileCacheFile) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

// Close closes the file.
// This should be called before the file is edited
// by the editor.
func (f *FileCacheFile) Close() error {
	return f.file.Close()
}

// CloseAndRemove closes the file and removes it.
func (f *FileCacheFile) CloseAndRemove() error {
	err := f.Close()
	if err != nil {
		return err
	}
	return os.Remove(f.fp)
}

// ReOpen opens the file again after it's been closed.
// This should be called after the file has been edited.
func (f *FileCacheFile) ReOpen() error {
	file, err := os.OpenFile(f.fp, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	f.file = file
	return nil
}

// FilePath returns the absolute path to the temporary file.
func (f *FileCacheFile) FilePath() string {
	return f.fp
}
