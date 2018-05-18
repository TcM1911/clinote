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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNote(t *testing.T) {
	assert := assert.New(t)
	t.Run("handle notestore error", func(t *testing.T) {
		testError := errors.New("expected error")
		client := new(mockClient)
		client.getNotestore = func() (NotestoreClient, error) { return nil, testError }
		_, err := GetNote(client, "note", "notebook")
		assert.EqualError(err, testError.Error())
	})
	t.Run("find note by title", func(t *testing.T) {
		title := "Expected Note"
		expectedNote := &Note{Title: title}
		c, _ := nsWithNote(expectedNote)
		note, err := GetNote(c, title, "")
		assert.NoError(err)
		assert.Equal(expectedNote, note)
	})
	t.Run("return error from FindNotes", func(t *testing.T) {
		expectedError := errors.New("Expected error")
		ns := new(mockNS)
		c := new(mockClient)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return nil, expectedError }
		_, err := GetNote(c, "title", "")
		assert.EqualError(err, expectedError.Error())
	})
	t.Run("error when note not found", func(t *testing.T) {
		title := "Note Title"
		otherNote1 := &Note{Title: "Other note"}
		otherNote2 := &Note{Title: "Other note2"}
		notes := []*Note{otherNote1, otherNote2}

		ns := new(mockNS)
		c := new(mockClient)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return notes, nil }
		_, err := GetNote(c, title, "")
		assert.EqualError(err, ErrNoNoteFound.Error())
	})
	t.Run("restrict notes by notebook", func(t *testing.T) {
		title := "Expected Note"
		notebook := "Expected Notebook"
		otherNote := &Note{Title: "Other note"}
		expectedNote := &Note{Title: title}
		notes := []*Note{otherNote, expectedNote}
		otherBook := &Notebook{Name: "Other Notebook"}
		expectedNotebook := &Notebook{Name: notebook, GUID: "GUID"}
		books := []*Notebook{otherBook, expectedNotebook}

		ns := new(mockNS)
		c := new(mockClient)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return notes, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return books, nil }
		note, err := GetNote(c, title, notebook)
		assert.NoError(err)
		assert.Equal(expectedNote, note)
	})
	t.Run("should return error from findNotebook", func(t *testing.T) {
		title := "Note"
		expectedError := errors.New("Expected error")

		ns := new(mockNS)
		c := new(mockClient)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return nil, expectedError }
		_, err := GetNote(c, title, "Notebook")
		assert.EqualError(err, expectedError.Error())
	})
}

func TestGetNoteContent(t *testing.T) {
	assert := assert.New(t)
	t.Run("return note with content", func(t *testing.T) {
		title := "Note title"
		expectedContent := "<p>Note content</p>"
		expectedNote := &Note{Title: title}
		c, ns := nsWithNote(expectedNote)
		ns.getNoteContent = func(guid string) (string, error) { return "<en-note>" + expectedContent + "</en-note>", nil }
		n, err := GetNoteWithContent(c, title)
		assert.NoError(err, "Should not return an error")
		assert.Equal(expectedNote, n, "Note doesn't match")
		assert.Equal(expectedContent, n.Body)
	})
	t.Run("return error from GetNoteContent", func(t *testing.T) {
		title := "Note title"
		expectedError := errors.New("Expected error")
		expectedNote := &Note{Title: title}
		c, ns := nsWithNote(expectedNote)
		ns.getNoteContent = func(guid string) (string, error) { return "", expectedError }
		_, err := GetNoteWithContent(c, title)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})
	t.Run("return error from GetNoteClient", func(t *testing.T) {
		title := "Note title"
		expectedError := errors.New("Expected error")
		c := new(mockClient)
		c.getNotestore = func() (NotestoreClient, error) { return nil, expectedError }
		_, err := GetNoteWithContent(c, title)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})
	t.Run("return error from decoder", func(t *testing.T) {
		title := "Note title"
		note := &Note{Title: title}
		c, ns := nsWithNote(note)
		ns.getNoteContent = func(string) (string, error) { return "", nil }
		_, err := GetNoteWithContent(c, title)
		assert.Error(err, "Expected an error")
	})
}

func TestSaveChanges(t *testing.T) {
	assert := assert.New(t)
	expectedError := errors.New("Expected error")
	body := "This is the note content"
	expectedMDContent := XMLHeader + "<en-note><p>" + body + "</p>\n</en-note>"
	expectedRawContent := XMLHeader + "<en-note><p>" + body + "</p></en-note>"
	t.Run("return error from GetNoteStore", func(t *testing.T) {
		client := new(mockClient)
		client.getNotestore = func() (NotestoreClient, error) { return nil, expectedError }
		// ctx := AddUseRawContentToContext(context.Background(), true)
		err := SaveChanges(client, &Note{}, true)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})
	t.Run("return error from UpdateNote", func(t *testing.T) {
		client := new(mockClient)
		ns := new(mockNS)
		ns.updateNote = func(n *Note) error { return expectedError }
		client.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		// ctx := AddUseRawContentToContext(context.Background(), false)
		err := SaveChanges(client, &Note{}, false)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})
	t.Run("UpdateNote without content change", func(t *testing.T) {
		client := new(mockClient)
		ns := new(mockNS)
		ns.updateNote = func(n *Note) error { return nil }
		client.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		// ctx := AddUseRawContentToContext(context.Background(), false)
		err := SaveChanges(client, &Note{}, false)
		assert.NoError(err, "Should not return an error")
	})
	t.Run("UpdateNote with MD content", func(t *testing.T) {
		client := new(mockClient)
		ns := new(mockNS)
		note := new(Note)
		note.MD = body
		ns.updateNote = func(n *Note) error { return nil }
		client.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		// ctx := AddUseRawContentToContext(context.Background(), false)
		err := SaveChanges(client, note, false)
		assert.NoError(err, "Should not return an error")
		assert.Equal(expectedMDContent, note.Body, "Note content doesn't match")
	})
	t.Run("UpdateNote with raw content", func(t *testing.T) {
		client := new(mockClient)
		ns := new(mockNS)
		note := new(Note)
		note.MD = body
		note.Body = "<p>" + body + "</p>"
		ns.updateNote = func(n *Note) error { return nil }
		client.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		// ctx := AddUseRawContentToContext(context.Background(), true)
		err := SaveChanges(client, note, true)
		assert.NoError(err, "Should not return an error")
		assert.Equal(expectedRawContent, note.Body, "Note content doesn't match")
	})
}

func TestChangeTitle(t *testing.T) {
	assert := assert.New(t)
	expectedError := errors.New("expected error")
	t.Run("should change title", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		note := &Note{Title: "Old"}
		var savedNote *Note
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.updateNote = func(n *Note) error { savedNote = n; return nil }

		err := ChangeTitle(c, "Old", "New")
		assert.NoError(err, "Should not return an error")
		assert.Equal(note, savedNote, "Same note should be saved")
		assert.Equal("New", savedNote.Title, "Title should be New")
	})
	t.Run("should handle error from saveChanges", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		note := &Note{Title: "Old"}
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.updateNote = func(*Note) error { return expectedError }

		err := ChangeTitle(c, "Old", "New")
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
	t.Run("should handle error from GetNote", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return nil, expectedError }

		err := ChangeTitle(c, "Old", "New")
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
}

func TestMoveNote(t *testing.T) {
	assert := assert.New(t)
	expectedError := errors.New("expected error")
	noteName := "Expected Note"
	notebookGUID := "Notebook GUID"
	notebookName := "New notebook"

	t.Run("should move note", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)

		notebook := &Notebook{Name: notebookName, GUID: notebookGUID}
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{notebook}, nil }

		note := &Note{Title: noteName, Notebook: &Notebook{Name: "Old", GUID: "Old GUID"}}
		var savedNote *Note
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.updateNote = func(n *Note) error { savedNote = n; return nil }

		err := MoveNote(c, noteName, notebookName)
		assert.NoError(err, "Should not return an error")
		assert.Equal(note, savedNote, "Same note should be saved")
		assert.Equal(notebook, savedNote.Notebook, "Incorrect notebook set")
		assert.Equal(notebookGUID, savedNote.Notebook.GUID, "The notebook should be New")
	})
	t.Run("should handle error from saveChanges", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		notebook := &Notebook{Name: notebookName, GUID: notebookGUID}
		note := &Note{Title: noteName, Notebook: notebook}
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{notebook}, nil }
		ns.updateNote = func(*Note) error { return expectedError }

		err := MoveNote(c, noteName, notebookName)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
	t.Run("should handle error from GetNote", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return nil, expectedError }

		err := MoveNote(c, noteName, notebookName)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
	t.Run("should handle error from FindNote", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		note := &Note{Title: noteName}
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return nil, expectedError }
		ns.updateNote = func(*Note) error { return expectedError }

		err := MoveNote(c, noteName, notebookName)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
}

func TestDeleteNote(t *testing.T) {
	assert := assert.New(t)
	noteGUID := "Note GUID"
	noteTitle := "Note title"
	expectedError := errors.New("expected error")
	t.Run("should delete note", func(t *testing.T) {
		note := &Note{Title: noteTitle, GUID: noteGUID}
		c, ns := nsWithNote(note)
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.deleteNote = func(g string) error {
			if g == noteGUID {
				return nil
			}
			return errors.New("wrong GUID")
		}
		err := DeleteNote(c, noteTitle, "")
		assert.NoError(err, "Should note return an error")
	})
	t.Run("should return error from GetNoteStore", func(t *testing.T) {
		note := &Note{Title: noteTitle, GUID: noteGUID}
		c, ns := nsWithNote(note)
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		count := 0
		c.getNotestore = func() (NotestoreClient, error) {
			if count == 0 {
				count++
				return ns, nil
			}
			return nil, expectedError
		}
		err := DeleteNote(c, noteTitle, "")
		assert.Error(err, "Should note return an error")
		assert.Equal(err, expectedError, "Wrong error returned")
	})
	t.Run("should return error from GetNote", func(t *testing.T) {
		c := new(mockClient)
		c.getNotestore = func() (NotestoreClient, error) { return nil, expectedError }
		err := DeleteNote(c, noteTitle, "")
		assert.Error(err, "Should note return an error")
		assert.Equal(err, expectedError, "Wrong error returned")
	})
	t.Run("should return error from DeleteNote", func(t *testing.T) {
		note := &Note{Title: noteTitle, GUID: noteGUID}
		c, ns := nsWithNote(note)
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.deleteNote = func(g string) error { return expectedError }
		err := DeleteNote(c, noteTitle, "")
		assert.Error(err, "Should note return an error")
		assert.Equal(err, expectedError, "Wrong error returned")
	})

}

func TestSaveNewNote(t *testing.T) {
	assert := assert.New(t)
	expectedError := errors.New("expected error")
	cases := []struct {
		Name string
		N    *Note
		Raw  bool
	}{
		{"empty note", &Note{}, false},
		{"with MD", &Note{MD: "content"}, false},
		{"raw content", &Note{Body: "<p>content</p>"}, true},
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			c := new(mockClient)
			ns := new(mockNS)
			var createdNote *Note
			c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
			ns.createNote = func(n *Note) error { createdNote = n; return nil }
			err := SaveNewNote(c, test.N, test.Raw)
			assert.NoError(err, "Should not return an error")
			assert.Equal(test.N, createdNote, "Should save the correct note")
		})
	}
	t.Run("return error from GetNoteStore", func(t *testing.T) {
		c := new(mockClient)
		c.getNotestore = func() (NotestoreClient, error) { return nil, expectedError }
		err := SaveNewNote(c, &Note{}, false)
		assert.Error(err, "should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})
	t.Run("return error from CreateNote", func(t *testing.T) {
		c := new(mockClient)
		ns := new(mockNS)
		c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
		ns.createNote = func(*Note) error { return expectedError }
		err := SaveNewNote(c, &Note{}, false)
		assert.Error(err, "should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})

}

func nsWithNote(note *Note) (*mockClient, *mockNS) {
	notes := []*Note{&Note{Title: "Other note"}, note}
	ns := new(mockNS)
	c := new(mockClient)
	c.getNotestore = func() (NotestoreClient, error) { return ns, nil }
	ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return notes, nil }
	return c, ns
}

type mockNS struct {
	findNotes       func(*NoteFilter, int, int) ([]*Note, error)
	getAllNotebooks func() ([]*Notebook, error)
	getNoteContent  func(guid string) (string, error)
	updateNote      func(n *Note) error
	deleteNote      func(guid string) error
	saveNewNote     func(n *Note) error
	createNote      func(n *Note) error
	updateNotebook  func(b *Notebook) error
}

func (s *mockNS) UpdateNotebook(b *Notebook) error {
	return s.updateNotebook(b)
}

func (s *mockNS) CreateNote(n *Note) error {
	return s.createNote(n)
}

func (s *mockNS) GetClient() APIClient {
	panic("not implemented")
}

func (s *mockNS) SaveNewNote(n *Note) error {
	return s.saveNewNote(n)
}

func (s *mockNS) DeleteNote(guid string) error {
	return s.deleteNote(guid)
}

func (s *mockNS) UpdateNote(n *Note) error {
	return s.updateNote(n)
}

func (s *mockNS) GetNoteContent(guid string) (string, error) {
	return s.getNoteContent(guid)
}

func (s *mockNS) FindNotes(filter *NoteFilter, offset int, count int) ([]*Note, error) {
	return s.findNotes(filter, offset, count)
}

func (s *mockNS) GetAllNotebooks() ([]*Notebook, error) {
	return s.getAllNotebooks()
}

func (s *mockNS) CreateNotebook(b *Notebook, defaultNotebook bool) error {
	panic("not implemented")
}

func (s *mockNS) GetNotebook(guid string) (*Notebook, error) {
	panic("not implemented")
}
