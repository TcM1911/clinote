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
 * Copyright (C) Joakim Kennedy, 2017
 */

// Package api provides interface for the Evernote API.
package api

import "github.com/TcM1911/evernote-sdk-golang/notestore"
import "github.com/TcM1911/evernote-sdk-golang/types"

// Notestore is the API interface for the notestore client.
type Notestore interface {
	// ListNotebooks returns a list of all the user's notebooks.
	ListNotebooks(apiKey string) (r []*types.Notebook, err error)
	// CreateNotebook creates a new notebook for the user.
	CreateNotebook(apiKey string, notebook *types.Notebook) (r *types.Notebook, err error)
	// UpdateNotebook sends an updated notebook to the server.
	UpdateNotebook(apiKey string, notebook *types.Notebook) (r int32, err error)
	// CreateNote creates a new note on the server.
	CreateNote(apiKey string, note *types.Note) (r *types.Note, err error)
	// DeleteNote moves a note to the trash can.
	DeleteNote(apiKey string, guid types.GUID) (int32, error)
	// FindNotes searches the server and returns notes matching the filter.
	FindNotes(apiKey string, filter *notestore.NoteFilter, offset int32, maxNumNotes int32) (r *notestore.NoteList, err error)
}
