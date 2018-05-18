package evernote

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNotebook(t *testing.T) {
	assert := assert.New(t)
	t.Run("return notebook", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{&Notebook{Name: "Book"}}, nil }
		b, err := FindNotebook(c, "Book")
		assert.NoError(err, "Should not return an error")
		assert.Equal("Book", b.Name, "Wrong notebook name")
	})
	t.Run("return error if no notebook", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{&Notebook{Name: "Book"}}, nil }
		_, err := FindNotebook(c, "Missing")
		assert.Error(err, "Should return an error")
		assert.Equal(ErrNoNotebookFound, err, "Wrong error returned")
	})
	t.Run("return error from GetNoteStore", func(t *testing.T) {
		c := new(mockClient)
		expectedErr := errors.New("test error")
		c.getNotestore = func() (NotestoreClient, error) { return nil, expectedErr }
		_, err := FindNotebook(c, "Book")
		assert.Error(err, "Should return an error")
	})
}

func TestGetNotebooks(t *testing.T) {
	assert := assert.New(t)
	t.Run("multiple books", func(t *testing.T) {
		books := []*Notebook{&Notebook{}, &Notebook{}}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return books, nil }}
		c := &mockClient{getNotestore: func() (NotestoreClient, error) { return ns, nil }}
		bs, err := GetNotebooks(c)
		assert.NoError(err, "Should not return an error")
		assert.Len(bs, 2, "Incorrect number of notebooks returned")
	})
	t.Run("return error from GetNoteStore", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		c := &mockClient{getNotestore: func() (NotestoreClient, error) { return nil, expectedErr }}
		_, err := GetNotebooks(c)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedErr, err, "Wrong error returned")
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
			c := &mockClient{getNotestore: func() (NotestoreClient, error) { return ns, nil }}
			err := UpdateNotebook(c, oldName, test.Book)
			assert.NoError(err, "Should not return an error")
			assert.Equal(test.ExpectedBook, saved, "Saved notebook doesn't match")
		})
	}
	t.Run("return error from GetNoteStore", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		c := &mockClient{getNotestore: func() (NotestoreClient, error) { return nil, expectedErr }}
		err := UpdateNotebook(c, "", &Notebook{})
		assert.Error(err, "Should return an error")
		assert.Equal(expectedErr, err, "Wrong error returned")
	})
	t.Run("return error from UpdateNotebook", func(t *testing.T) {
		expectedErr := errors.New("expected error")
		oldBook := &Notebook{Name: oldName, Stack: oldStack}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return []*Notebook{oldBook}, nil },
			updateNotebook: func(book *Notebook) error { return expectedErr }}
		c := &mockClient{getNotestore: func() (NotestoreClient, error) { return ns, nil }}
		err := UpdateNotebook(c, oldName, &Notebook{})
		assert.Error(err, "Should return an error")
		assert.Equal(expectedErr, err, "Wrong error returned")
	})
	t.Run("return error when no notebook found", func(t *testing.T) {
		oldBook := &Notebook{Name: oldName, Stack: oldStack}
		ns := &mockNS{getAllNotebooks: func() ([]*Notebook, error) { return []*Notebook{oldBook}, nil },
			updateNotebook: func(book *Notebook) error { return nil }}
		c := &mockClient{getNotestore: func() (NotestoreClient, error) { return ns, nil }}
		err := UpdateNotebook(c, newName, &Notebook{})
		assert.Error(err, "Should return an error")
		assert.Equal(ErrNoNotebookFound, err, "Wrong error returned")
	})
}
