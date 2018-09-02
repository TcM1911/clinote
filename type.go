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

import "io"

// Storager is the interface for backend storage.
type Storager interface {
	io.Closer
	GetSettings() (*Settings, error)
	StoreSettings(*Settings) error
	// GetNotebookCache returns the stored NotebookCacheList.
	GetNotebookCache() (*NotebookCacheList, error)
	// StoreNotebookList saves the list to the database.
	StoreNotebookList(list *NotebookCacheList) error
	// SaveSearch stores a note search to the database.
	SaveSearch([]*Note) error
	// GetSearch returns a saved note search from the database.
	GetSearch() ([]*Note, error)
	// SaveNoteRecoveryPoint saves the note as a recovery point.
	SaveNoteRecoveryPoint(*Note) error
	// GetNoteREcoveryPoint returns the saved note.
	GetNoteRecoveryPoint() (*Note, error)
}

// Settings is a struct holding the user's settings for the application.
type Settings struct {
	// APIKey is the user's OAuth session key.
	APIKey string
}
