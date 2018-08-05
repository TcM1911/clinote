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

package clinote

import "errors"

var (
	// ErrNoNotebookFound is returned if no matching notebook was found.
	ErrNoNotebookFound = errors.New("no notebook found")
	// ErrNoNotebookCached is returned when trying to update a notebook
	// that hasn't been pulled from the server.
	ErrNoNotebookCached = errors.New("no notebook found")
)

// Notebook is a struct for the notebook.
type Notebook struct {
	// Name is the notebook's name
	Name string
	// GUID is the notebook's GUID.
	GUID string
	// Stack is the stack that the notebook belongs too.
	Stack string
}

// UpdateNotebook updates the notebook.
func UpdateNotebook(db Storager, ns NotestoreClient, name string, notebook *Notebook) error {
	b, err := findNotebook(db, ns, name)
	if err != nil {
		return err
	}
	if notebook.Name != "" {
		b.Name = notebook.Name
	}
	if notebook.Stack != "" {
		b.Stack = notebook.Stack
	}
	return ns.UpdateNotebook(b)
}

// FindNotebook gets the notebook matching with the name.
// If no notebook is found, nil is returned.
func FindNotebook(db Storager, ns NotestoreClient, name string) (*Notebook, error) {
	return findNotebook(db, ns, name)
}

func findNotebook(db Storager, ns NotestoreClient, name string) (*Notebook, error) {
	bs, err := GetNotebooks(db, ns, false)
	if err != nil {
		return nil, err
	}
	for _, b := range bs {
		if b.Name == name {
			return b, nil
		}
	}
	return nil, ErrNoNotebookFound
}

// GetNotebooks returns all the user's notebooks.
func GetNotebooks(db Storager, ns NotestoreClient, forceSync bool) ([]*Notebook, error) {
	list, err := db.GetNotebookCache()
	if err != nil {
		return nil, err
	}
	if !list.IsOutdated() && len(list.Notebooks) > 0 && !forceSync {
		return list.Notebooks, nil
	}
	bs, err := ns.GetAllNotebooks()
	if err != nil {
		return nil, err
	}
	list = NewNotebookCacheList(bs)
	err = db.StoreNotebookList(list)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// GetNotebook returns a notebook from the user's notestore.
func GetNotebook(ns NotestoreClient, guid string) (*Notebook, error) {
	return ns.GetNotebook(guid)
}

// CreateNotebook creates a new notebook.
func CreateNotebook(ns NotestoreClient, notebook *Notebook, defaultNotebook bool) error {
	return ns.CreateNotebook(notebook, defaultNotebook)
}
