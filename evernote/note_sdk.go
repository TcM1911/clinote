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

package evernote

import (
	"errors"
	"sync"
	"time"

	"github.com/tcm1911/evernote-sdk-golang/types"
)

// ErrNoCachedNote is return if the note wasn't cached and can't be
// updated.
var ErrNoCachedNote = errors.New("no cache note found")

var noteMu sync.Mutex
var cache map[types.GUID]*types.Note

func init() {
	cache = make(map[types.GUID]*types.Note)
}

func createNoteType(n *Note) (*types.Note, error) {
	noteMu.Lock()
	note, ok := cache[types.GUID(n.GUID)]
	if !ok {
		// No cached note, so we can't update.
		return nil, ErrNoCachedNote
	}
	// Remove cached note.
	delete(cache, note.GetGUID())
	noteMu.Unlock()
	note.Title = &n.Title

	notebookGUID := string(n.Notebook.GUID)
	note.NotebookGuid = &notebookGUID

	now := types.Timestamp(time.Now().Unix() * 1000)
	note.Updated = &now
	return note, nil
}

func convert(note *types.Note) *Note {
	n := new(Note)
	n.Title = note.GetTitle()
	n.GUID = string(note.GetGUID())
	notebook := new(Notebook)
	notebookGUID := note.GetNotebookGuid()
	n.Notebook = notebook
	n.Notebook.GUID = notebookGUID
	return n
}

func convertNotes(notes []*types.Note) []*Note {
	a := make([]*Note, len(notes))
	for i, n := range notes {
		// Cache notes for later.
		noteMu.Lock()
		cache[*n.GUID] = n
		noteMu.Unlock()
		a[i] = convert(n)
	}
	return a
}
