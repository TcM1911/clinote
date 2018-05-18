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

package evernote

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
func UpdateNotebook(client APIClient, name string, notebook *Notebook) error {
	ns, err := client.GetNoteStore()
	if err != nil {
		return err
	}
	b, err := findNotebook(ns, name)
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
func FindNotebook(client APIClient, name string) (*Notebook, error) {
	ns, err := client.GetNoteStore()
	if err != nil {
		return nil, err
	}
	return findNotebook(ns, name)
}

func findNotebook(ns NotestoreClient, name string) (*Notebook, error) {
	bs, err := ns.GetAllNotebooks()
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
func GetNotebooks(client APIClient) ([]*Notebook, error) {
	ns, err := client.GetNoteStore()
	if err != nil {
		return nil, err
	}
	return ns.GetAllNotebooks()
}

// GetNotebook returns a notebook from the user's notestore.
func GetNotebook(client APIClient, guid string) (*Notebook, error) {
	ns, err := client.GetNoteStore()
	if err != nil {
		return nil, err
	}
	return ns.GetNotebook(guid)
}

// CreateNotebook creates a new notebook.
func CreateNotebook(client APIClient, notebook *Notebook, defaultNotebook bool) error {
	ns, err := client.GetNoteStore()
	if err != nil {
		return err
	}
	return ns.CreateNotebook(notebook, defaultNotebook)
}
