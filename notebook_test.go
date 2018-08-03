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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindNotebook(t *testing.T) {
	assert := assert.New(t)
	store := &mockStore{
		getNotebookCache:  func() (*NotebookCacheList, error) { return &NotebookCacheList{Notebooks: []*Notebook{}}, nil },
		storeNotebookList: func(list *NotebookCacheList) error { return nil },
	}
	t.Run("return notebook", func(t *testing.T) {
		ns := new(mockNS)
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{&Notebook{Name: "Book"}}, nil }
		b, err := FindNotebook(store, ns, "Book")
		assert.NoError(err, "Should not return an error")
		assert.Equal("Book", b.Name, "Wrong notebook name")
	})
	t.Run("return error if no notebook", func(t *testing.T) {
		ns := new(mockNS)
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{&Notebook{Name: "Book"}}, nil }
		_, err := FindNotebook(store, ns, "Missing")
		assert.Error(err, "Should return an error")
		assert.Equal(ErrNoNotebookFound, err, "Wrong error returned")
	})
}

func TestGetNotebooks(t *testing.T) {
	assert := assert.New(t)
	var storedList *NotebookCacheList
	createMocks := func(empty, ex bool) (*mockNS, *mockStore, []*Notebook, *NotebookCacheList) {
		var cachedBooks *NotebookCacheList
		if empty {
			cachedBooks = &NotebookCacheList{Notebooks: []*Notebook{}}
		} else if ex {
			cachedBooks = NewNotebookCacheListWithLimit([]*Notebook{&Notebook{}, &Notebook{}}, 1*time.Nanosecond)
		} else {
			cachedBooks = NewNotebookCacheList([]*Notebook{&Notebook{}, &Notebook{}})
		}
		books := []*Notebook{&Notebook{}, &Notebook{}}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return books, nil }}
		store := &mockStore{
			getNotebookCache:  func() (*NotebookCacheList, error) { return cachedBooks, nil },
			storeNotebookList: func(list *NotebookCacheList) error { storedList = list; return nil },
		}
		return ns, store, books, cachedBooks
	}
	t.Run("multiple books from notestore", func(t *testing.T) {
		ns, db, expectedBooks, _ := createMocks(true, false)
		bs, err := GetNotebooks(db, ns, false)
		assert.NoError(err, "Should not return an error")
		assert.Len(bs, 2, "Incorrect number of notebooks returned")
		assert.Equal(expectedBooks, bs, "Wrong books returned")
		assert.Equal(expectedBooks, storedList.Notebooks, "Wrong books cached")
	})
	t.Run("refresh if expired", func(t *testing.T) {
		ns, db, expectedBooks, _ := createMocks(false, true)
		time.Sleep(10 * time.Microsecond)
		bs, err := GetNotebooks(db, ns, false)
		assert.NoError(err, "Should not return an error")
		assert.Len(bs, 2, "Incorrect number of notebooks returned")
		assert.Equal(expectedBooks, bs, "Wrong books returned")
		assert.Equal(expectedBooks, storedList.Notebooks, "Wrong books cached")
	})
	t.Run("multiple books from cache", func(t *testing.T) {
		ns, db, _, expectedCache := createMocks(false, false)
		bs, err := GetNotebooks(db, ns, false)
		assert.NoError(err, "Should not return an error")
		assert.Len(bs, 2, "Incorrect number of notebooks returned")
		assert.Equal(expectedCache.Notebooks, bs, "Wrong books returned")
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
	store := &mockStore{
		getNotebookCache:  func() (*NotebookCacheList, error) { return &NotebookCacheList{Notebooks: []*Notebook{}}, nil },
		storeNotebookList: func(list *NotebookCacheList) error { return nil },
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			oldBook := &Notebook{Name: oldName, Stack: oldStack}
			var saved *Notebook
			ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return []*Notebook{oldBook}, nil },
				updateNotebook: func(book *Notebook) error { saved = book; return nil }}
			err := UpdateNotebook(store, ns, oldName, test.Book)
			assert.NoError(err, "Should not return an error")
			assert.Equal(test.ExpectedBook, saved, "Saved notebook doesn't match")
		})
	}
	t.Run("return error from UpdateNotebook", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		oldBook := &Notebook{Name: oldName, Stack: oldStack}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return []*Notebook{oldBook}, nil },
			updateNotebook: func(book *Notebook) error { return expectedErr }}
		err := UpdateNotebook(store, ns, oldName, &Notebook{})
		assert.Error(err, "Should return an error")
		assert.Equal(expectedErr, err, "Wrong error returned")
	})
	t.Run("return error when no notebook found", func(t *testing.T) {
		oldBook := &Notebook{Name: oldName, Stack: oldStack}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return []*Notebook{oldBook}, nil },
			updateNotebook: func(book *Notebook) error { return nil }}
		err := UpdateNotebook(store, ns, newName, &Notebook{})
		assert.Error(err, "Should return an error")
		assert.Equal(ErrNoNotebookFound, err, "Wrong error returned")
	})
}
