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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNotebook(t *testing.T) {
	assert := assert.New(t)
	t.Run("return notebook", func(t *testing.T) {
		ns := new(mockNS)
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{&Notebook{Name: "Book"}}, nil }
		b, err := FindNotebook(ns, "Book")
		assert.NoError(err, "Should not return an error")
		assert.Equal("Book", b.Name, "Wrong notebook name")
	})
	t.Run("return error if no notebook", func(t *testing.T) {
		ns := new(mockNS)
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{&Notebook{Name: "Book"}}, nil }
		_, err := FindNotebook(ns, "Missing")
		assert.Error(err, "Should return an error")
		assert.Equal(ErrNoNotebookFound, err, "Wrong error returned")
	})
}

func TestGetNotebooks(t *testing.T) {
	assert := assert.New(t)
	t.Run("multiple books", func(t *testing.T) {
		books := []*Notebook{&Notebook{}, &Notebook{}}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return books, nil }}
		bs, err := GetNotebooks(ns)
		assert.NoError(err, "Should not return an error")
		assert.Len(bs, 2, "Incorrect number of notebooks returned")
	})
}

func TestUpdateNotebook(t *testing.T) {
	assert := assert.New(t)
	newName, oldName, newStack, oldStack := "New Name", "Old Name", "New Stack", "Old Stack"
	tests := []struct {
		Name         string
		Book         *Notebook
		ExpectedBook *Notebook
	}{
		{"Change name", &Notebook{Name: newName}, &Notebook{Name: newName, Stack: oldStack}},
		{"Change stack", &Notebook{Stack: newStack}, &Notebook{Name: oldName, Stack: newStack}},
		{"Change name and stack", &Notebook{Name: newName, Stack: newStack}, &Notebook{Name: newName, Stack: newStack}},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			oldBook := &Notebook{Name: oldName, Stack: oldStack}
			var saved *Notebook
			ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return []*Notebook{oldBook}, nil },
				updateNotebook: func(book *Notebook) error { saved = book; return nil }}
			err := UpdateNotebook(ns, oldName, test.Book)
			assert.NoError(err, "Should not return an error")
			assert.Equal(test.ExpectedBook, saved, "Saved notebook doesn't match")
		})
	}
	t.Run("return error from UpdateNotebook", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		oldBook := &Notebook{Name: oldName, Stack: oldStack}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return []*Notebook{oldBook}, nil },
			updateNotebook: func(book *Notebook) error { return expectedErr }}
		err := UpdateNotebook(ns, oldName, &Notebook{})
		assert.Error(err, "Should return an error")
		assert.Equal(expectedErr, err, "Wrong error returned")
	})
	t.Run("return error when no notebook found", func(t *testing.T) {
		oldBook := &Notebook{Name: oldName, Stack: oldStack}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return []*Notebook{oldBook}, nil },
			updateNotebook: func(book *Notebook) error { return nil }}
		err := UpdateNotebook(ns, newName, &Notebook{})
		assert.Error(err, "Should return an error")
		assert.Equal(ErrNoNotebookFound, err, "Wrong error returned")
	})
}
