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

	"github.com/TcM1911/evernote-sdk-golang/types"
)

var (
	// ErrNoNotebookFound is returned if no matching notebook was found.
	ErrNoNotebookFound = errors.New("no notebook found")
	// ErrNoNotebookCached is returned when trying to update a notebook
	// that hasn't been pulled from the server.
	ErrNoNotebookCached = errors.New("no notebook found")
)

var notebookMu sync.Mutex
var cachedNotebooks map[types.GUID]*types.Notebook

func convertNotebooks(bs []*types.Notebook) []*Notebook {
	a := make([]*Notebook, len(bs), len(bs))
	for i, b := range bs {
		a[i] = &Notebook{GUID: string(b.GetGUID()), Name: b.GetName(), Stack: b.GetStack()}
	}
	return a
}

func transferNotebookData(src *Notebook, dst *types.Notebook) {
	dst.Name = &(src.Name)
	dst.Stack = &(src.Stack)
}

func cacheNotebook(nb *types.Notebook) {
	notebookMu.Lock()
	defer notebookMu.Unlock()
	if cachedNotebooks == nil {
		cachedNotebooks = make(map[types.GUID]*types.Notebook)
	}
	cachedNotebooks[*nb.GUID] = nb
}

func getCachedNotebook(guid types.GUID) (*types.Notebook, error) {
	notebookMu.Lock()
	defer notebookMu.Unlock()
	if cachedNotebooks == nil {
		return nil, ErrNoNotebookCached
	}
	nb, ok := cachedNotebooks[guid]
	if !ok {
		return nil, ErrNoNotebookFound
	}
	return nb, nil
}
