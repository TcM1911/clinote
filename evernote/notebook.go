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

import (
	"errors"
	"fmt"

	"github.com/tcm1911/clinote/config"
	"github.com/tcm1911/evernote-sdk-golang/types"
)

// Notebook is a struct for the notebook.
type Notebook struct {
	// Name is the notebook's name
	Name string
	// GUID is the notebook's GUID.
	GUID types.GUID
	// Stack is the stack that the notebook belongs too.
	Stack string
}

// UpdateNotebook updates the notebook.
func UpdateNotebook(cfg config.Configuration, name string, notebook *Notebook) {
	b, err := findNotebook(cfg, name)
	if err != nil {
		fmt.Println("Error when looking for", name, ":", err)
		return
	}
	if notebook.Name != "" {
		fmt.Println("Changing notebook name to", notebook.Name)
		b.Name = &notebook.Name
	}
	if notebook.Stack != "" {
		fmt.Println("Changing notebook stack to", notebook.Stack)
		b.Stack = &notebook.Stack
	}
	ns := GetNoteStore(cfg)
	if _, err := ns.UpdateNotebook(AuthToken, b); err != nil {
		fmt.Println("Error when updating the notebook:", err)
		return
	}
	fmt.Println("Notebook updated.")
}

// FindNotebook gets the notebook matching with the name.
// If no notebook is found, nil is returned.
func FindNotebook(cfg config.Configuration, name string) (*Notebook, error) {
	b, err := findNotebook(cfg, name)
	if err != nil {
		return nil, err
	}
	book := new(Notebook)
	book.Name = b.GetName()
	book.GUID = b.GetGUID()
	book.Stack = b.GetStack()
	return book, nil
}

func findNotebook(cfg config.Configuration, name string) (*types.Notebook, error) {
	bs, err := getNotebooks(cfg)
	if err != nil {
		return nil, err
	}
	for _, b := range bs {
		if b.GetName() == name {
			return b, nil
		}
	}
	return nil, errors.New("no notebook found")
}

// GetNotebooks returns all the user's notebooks.
func GetNotebooks(cfg config.Configuration) ([]*Notebook, error) {
	books, err := getNotebooks(cfg)
	if err != nil {
		return nil, err
	}
	bs := make([]*Notebook, len(books))
	for i, book := range books {
		p := new(Notebook)
		p.Name = book.GetName()
		p.GUID = book.GetGUID()
		p.Stack = book.GetStack()
		bs[i] = p
	}
	return bs, nil
}

func getNotebooks(cfg config.Configuration) ([]*types.Notebook, error) {
	ns := GetNoteStore(cfg)
	return ns.ListNotebooks(AuthToken)
}
