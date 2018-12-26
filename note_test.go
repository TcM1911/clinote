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

package clinote

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNote(t *testing.T) {
	assert := assert.New(t)
	store := &mockStore{
		getNotebookCache:  func() (*NotebookCacheList, error) { return &NotebookCacheList{Notebooks: []*Notebook{}}, nil },
		storeNotebookList: func(list *NotebookCacheList) error { return nil },
	}
	t.Run("find note by title", func(t *testing.T) {
		title := "Expected Note"
		expectedNote := &Note{Title: title}
		ns := nsWithNote(expectedNote)
		note, err := GetNote(store, ns, title, "")
		assert.NoError(err)
		assert.Equal(expectedNote, note)
	})
	t.Run("get note from search", func(t *testing.T) {
		expectedNote := new(Note)
		store.getSearch = func() ([]*Note, error) {
			return []*Note{new(Note), expectedNote, new(Note)}, nil
		}
		ns := nsWithNote(expectedNote)
		note, err := GetNote(store, ns, "2", "")
		assert.NoError(err)
		assert.Equal(expectedNote, note)
	})
	t.Run("handle cache note index overflow", func(t *testing.T) {
		store.getSearch = func() ([]*Note, error) {
			return []*Note{new(Note), new(Note), new(Note)}, nil
		}
		notes := []*Note{new(Note), new(Note)}
		ns := new(mockNS)
		ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return notes, nil }
		_, err := GetNote(store, ns, "4", "")
		assert.Error(err)
		assert.EqualError(err, ErrNoNoteFound.Error())
	})
	t.Run("return error from FindNotes", func(t *testing.T) {
		expectedError := errors.New("Expected error")
		ns := new(mockNS)
		ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return nil, expectedError }
		_, err := GetNote(store, ns, "title", "")
		assert.EqualError(err, expectedError.Error())
	})
	t.Run("error when note not found", func(t *testing.T) {
		title := "Note Title"
		otherNote1 := &Note{Title: "Other note"}
		otherNote2 := &Note{Title: "Other note2"}
		notes := []*Note{otherNote1, otherNote2}

		ns := new(mockNS)
		ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return notes, nil }
		_, err := GetNote(store, ns, title, "")
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
		ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return notes, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return books, nil }
		note, err := GetNote(store, ns, title, notebook)
		assert.NoError(err)
		assert.Equal(expectedNote, note)
	})
	t.Run("should return error from findNotebook", func(t *testing.T) {
		title := "Note"
		expectedError := errors.New("Expected error")

		ns := new(mockNS)
		ns.getAllNotebooks = func() ([]*Notebook, error) { return nil, expectedError }
		_, err := GetNote(store, ns, title, "Notebook")
		assert.EqualError(err, expectedError.Error())
	})
}

func TestGetNoteContent(t *testing.T) {
	assert := assert.New(t)
	store := &mockStore{
		getNotebookCache:  func() (*NotebookCacheList, error) { return &NotebookCacheList{Notebooks: []*Notebook{}}, nil },
		storeNotebookList: func(list *NotebookCacheList) error { return nil },
	}
	t.Run("return note with content", func(t *testing.T) {
		title := "Note title"
		expectedContent := "<p>Note content</p>\n"
		expectedNote := &Note{Title: title}
		ns := nsWithNote(expectedNote)
		ns.getNoteContent = func(guid string) (string, error) { return "<en-note>" + expectedContent + "</en-note>", nil }
		n, err := GetNoteWithContent(store, ns, title)
		assert.NoError(err, "Should not return an error")
		assert.Equal(expectedNote, n, "Note doesn't match")
		assert.Equal(expectedContent, n.Body)
	})
	t.Run("return error from GetNoteContent", func(t *testing.T) {
		title := "Note title"
		expectedError := errors.New("Expected error")
		expectedNote := &Note{Title: title}
		ns := nsWithNote(expectedNote)
		ns.getNoteContent = func(guid string) (string, error) { return "", expectedError }
		_, err := GetNoteWithContent(store, ns, title)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})
	t.Run("return error from decoder", func(t *testing.T) {
		title := "Note title"
		note := &Note{Title: title}
		ns := nsWithNote(note)
		ns.getNoteContent = func(string) (string, error) { return "", nil }
		_, err := GetNoteWithContent(store, ns, title)
		assert.Error(err, "Expected an error")
	})
}

func TestSaveChanges(t *testing.T) {
	assert := assert.New(t)
	expectedError := errors.New("Expected error")
	body := "This is the note content"
	expectedMDContent := XMLHeader + "<en-note><p>" + body + "</p>\n</en-note>"
	expectedRawContent := XMLHeader + "<en-note><p>" + body + "</p></en-note>"
	var opts NoteOption
	t.Run("return error from UpdateNote", func(t *testing.T) {
		ns := new(mockNS)
		ns.updateNote = func(n *Note) error { return expectedError }
		err := SaveChanges(ns, &Note{}, opts)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})
	t.Run("UpdateNote without content change", func(t *testing.T) {
		ns := new(mockNS)
		ns.updateNote = func(n *Note) error { return nil }
		err := SaveChanges(ns, &Note{}, opts)
		assert.NoError(err, "Should not return an error")
	})
	t.Run("UpdateNote with MD content", func(t *testing.T) {
		ns := new(mockNS)
		note := new(Note)
		note.MD = body
		ns.updateNote = func(n *Note) error { return nil }
		err := SaveChanges(ns, note, opts)
		assert.NoError(err, "Should not return an error")
		assert.Equal(expectedMDContent, note.Body, "Note content doesn't match")
	})
	t.Run("UpdateNote with raw content", func(t *testing.T) {
		ns := new(mockNS)
		note := new(Note)
		note.MD = body
		note.Body = "<p>" + body + "</p>"
		ns.updateNote = func(n *Note) error { return nil }
		err := SaveChanges(ns, note, RawNote)
		assert.NoError(err, "Should not return an error")
		assert.Equal(expectedRawContent, note.Body, "Note content doesn't match")
	})
}

func TestChangeTitle(t *testing.T) {
	assert := assert.New(t)
	expectedError := errors.New("expected error")
	store := &mockStore{
		getNotebookCache:  func() (*NotebookCacheList, error) { return &NotebookCacheList{Notebooks: []*Notebook{}}, nil },
		storeNotebookList: func(list *NotebookCacheList) error { return nil },
	}
	t.Run("should change title", func(t *testing.T) {
		ns := new(mockNS)
		note := &Note{Title: "Old"}
		var savedNote *Note
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.updateNote = func(n *Note) error { savedNote = n; return nil }

		err := ChangeTitle(store, ns, "Old", "New")
		assert.NoError(err, "Should not return an error")
		assert.Equal(note, savedNote, "Same note should be saved")
		assert.Equal("New", savedNote.Title, "Title should be New")
	})
	t.Run("should handle error from saveChanges", func(t *testing.T) {
		ns := new(mockNS)
		note := &Note{Title: "Old"}
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.updateNote = func(*Note) error { return expectedError }

		err := ChangeTitle(store, ns, "Old", "New")
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
	t.Run("should handle error from GetNote", func(t *testing.T) {
		ns := new(mockNS)
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return nil, expectedError }

		err := ChangeTitle(store, ns, "Old", "New")
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
	store := &mockStore{
		getNotebookCache:  func() (*NotebookCacheList, error) { return &NotebookCacheList{Notebooks: []*Notebook{}}, nil },
		storeNotebookList: func(list *NotebookCacheList) error { return nil },
	}

	t.Run("should move note", func(t *testing.T) {
		ns := new(mockNS)

		notebook := &Notebook{Name: notebookName, GUID: notebookGUID}
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{notebook}, nil }

		note := &Note{Title: noteName, Notebook: &Notebook{Name: "Old", GUID: "Old GUID"}}
		var savedNote *Note
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.updateNote = func(n *Note) error { savedNote = n; return nil }

		err := MoveNote(store, ns, noteName, notebookName)
		assert.NoError(err, "Should not return an error")
		assert.Equal(note, savedNote, "Same note should be saved")
		assert.Equal(notebook, savedNote.Notebook, "Incorrect notebook set")
		assert.Equal(notebookGUID, savedNote.Notebook.GUID, "The notebook should be New")
	})
	t.Run("should handle error from saveChanges", func(t *testing.T) {
		ns := new(mockNS)
		notebook := &Notebook{Name: notebookName, GUID: notebookGUID}
		note := &Note{Title: noteName, Notebook: notebook}
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return []*Notebook{notebook}, nil }
		ns.updateNote = func(*Note) error { return expectedError }

		err := MoveNote(store, ns, noteName, notebookName)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
	t.Run("should handle error from GetNote", func(t *testing.T) {
		ns := new(mockNS)
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return nil, expectedError }

		err := MoveNote(store, ns, noteName, notebookName)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
	t.Run("should handle error from FindNote", func(t *testing.T) {
		ns := new(mockNS)
		note := &Note{Title: noteName}
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.getAllNotebooks = func() ([]*Notebook, error) { return nil, expectedError }
		ns.updateNote = func(*Note) error { return expectedError }

		err := MoveNote(store, ns, noteName, notebookName)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Not the correct error")
	})
}

func TestDeleteNote(t *testing.T) {
	assert := assert.New(t)
	noteGUID := "Note GUID"
	noteTitle := "Note title"
	expectedError := errors.New("expected error")
	store := &mockStore{
		getNotebookCache:  func() (*NotebookCacheList, error) { return &NotebookCacheList{Notebooks: []*Notebook{}}, nil },
		storeNotebookList: func(list *NotebookCacheList) error { return nil },
	}
	t.Run("should delete note", func(t *testing.T) {
		note := &Note{Title: noteTitle, GUID: noteGUID}
		ns := nsWithNote(note)
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.deleteNote = func(g string) error {
			if g == noteGUID {
				return nil
			}
			return errors.New("wrong GUID")
		}
		err := DeleteNote(store, ns, noteTitle, "")
		assert.NoError(err, "Should note return an error")
	})
	t.Run("should return error from DeleteNote", func(t *testing.T) {
		note := &Note{Title: noteTitle, GUID: noteGUID}
		ns := nsWithNote(note)
		ns.findNotes = func(*NoteFilter, int, int) ([]*Note, error) { return []*Note{note}, nil }
		ns.deleteNote = func(g string) error { return expectedError }
		err := DeleteNote(store, ns, noteTitle, "")
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
			ns := new(mockNS)
			var createdNote *Note
			ns.createNote = func(n *Note) error { createdNote = n; return nil }
			err := SaveNewNote(ns, test.N, test.Raw)
			assert.NoError(err, "Should not return an error")
			assert.Equal(test.N, createdNote, "Should save the correct note")
		})
	}
	t.Run("return error from CreateNote", func(t *testing.T) {
		ns := new(mockNS)
		ns.createNote = func(*Note) error { return expectedError }
		err := SaveNewNote(ns, &Note{}, false)
		assert.Error(err, "should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})
}

func TestEditNote(t *testing.T) {
	assert := assert.New(t)

	// Setup test fixtures.
	var setupClientAndStore = func(addToNote string) (*Client, *mockNS, *[]byte, *Note, string, *mockStore) {
		// Setup store
		store := new(mockStore)

		// Setup notestore
		noteTitle := "Note Title"
		originalContent := "Body content"
		expectedNote := &Note{
			Title: noteTitle,
			Body:  "<en-note><p>" + originalContent + "</p></en-note>",
			MD:    originalContent,
			GUID:  "NOTEGUID",
		}
		ns := nsWithNote(expectedNote)
		ns.getNoteContent = func(guid string) (string, error) { return expectedNote.Body, nil }

		// Setup editer
		var writtenData []byte
		editor := &mockEditor{
			edit: func(file CacheFile) error {
				cache, ok := file.(*mockCacheFile)
				if !ok {
					t.Fatalf("Wrong CacheFile type\n")
				}
				data := cache.buffer.Bytes()
				d := make([]byte, len(data))
				copy(d, data)
				writtenData = d
				if addToNote != "" {
					_, err := cache.buffer.WriteString("\n\n" + addToNote + "\n")
					return err
				}
				return nil
			},
		}

		// Setup client
		c := &Client{
			Store:     store,
			Config:    new(DefaultConfig),
			NoteStore: ns,
			Editor:    editor,
		}
		c.newCacheFile = func(c *Client, filename string) (CacheFile, error) {
			buf := new(bytes.Buffer)
			return &mockCacheFile{buffer: buf}, nil
		}
		return c, ns, &writtenData, expectedNote, originalContent, store
	}
	var setupClient = func(addToNote string) (*Client, *mockNS, *[]byte, *Note, string) {
		a, b, c, d, e, _ := setupClientAndStore(addToNote)
		return a, b, c, d, e
	}

	// No edit
	t.Run("no_change_md", func(t *testing.T) {
		c, ns, writtenData, expectedNote, originalContent := setupClient("")
		saveNoteCalled := false
		ns.updateNote = func(*Note) error {
			saveNoteCalled = true
			return nil
		}
		err := EditNote(c, expectedNote.Title, DefaultNoteOption)
		assert.NoError(err, "Should not return an error")
		assert.NotNil(writtenData, "Should record the data")
		assert.Contains(string(*writtenData), expectedNote.Title, "Should write title to file")
		assert.Contains(string(*writtenData), originalContent, "Should write content to file")
		assert.NotContains(string(*writtenData), "<p>", "Should not include HTML")
		assert.False(saveNoteCalled, "Should not call SaveNote")
	})

	t.Run("no_change_raw", func(t *testing.T) {
		c, ns, writtenData, expectedNote, originalContent := setupClient("")
		saveNoteCalled := false
		ns.updateNote = func(*Note) error {
			saveNoteCalled = true
			return nil
		}
		err := EditNote(c, expectedNote.Title, RawNote)
		assert.NoError(err, "Should not return an error")
		assert.NotNil(writtenData, "Should record the data")
		assert.Contains(string(*writtenData), expectedNote.Title, "Should write title to file")
		assert.Contains(string(*writtenData), originalContent, "Should write content to file")
		assert.Contains(string(*writtenData), "<p>", "Should include HTML")
		assert.False(saveNoteCalled, "Should not call SaveNote")
	})

	// Detect edit
	t.Run("change_md", func(t *testing.T) {
		addedToNote := "New content added"
		c, ns, writtenData, expectedNote, originalContent := setupClient(addedToNote)
		saveNoteCalled := false
		var savedNote *Note
		ns.updateNote = func(n *Note) error {
			saveNoteCalled = true
			savedNote = n
			return nil
		}
		err := EditNote(c, expectedNote.Title, DefaultNoteOption)
		assert.NoError(err, "Should not return an error")
		assert.NotNil(writtenData, "Should record the data")
		assert.Contains(string(*writtenData), expectedNote.Title, "Should write title to file")
		assert.Contains(string(*writtenData), originalContent, "Should write content to file")
		assert.NotContains(string(*writtenData), "<p>", "Should not include HTML")
		assert.True(saveNoteCalled, "Should call SaveNote")
		assert.NotNil(savedNote, "Saved note should not be nil")
		assert.Contains(savedNote.Body, addedToNote, "Saved note should include added data")
	})

	t.Run("change_raw", func(t *testing.T) {
		addedToNote := "<p>New content added</p>"
		c, ns, writtenData, expectedNote, originalContent := setupClient(addedToNote)
		saveNoteCalled := false
		var savedNote *Note
		ns.updateNote = func(n *Note) error {
			saveNoteCalled = true
			savedNote = n
			return nil
		}
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|RawNote)
		assert.NoError(err, "Should not return an error")
		assert.NotNil(writtenData, "Should record the data")
		assert.Contains(string(*writtenData), expectedNote.Title, "Should write title to file")
		assert.Contains(string(*writtenData), originalContent, "Should write content to file")
		assert.Contains(string(*writtenData), "<p>", "Should include HTML")
		assert.True(saveNoteCalled, "Should call SaveNote")
		assert.NotNil(savedNote, "Saved note should not be nil")
		assert.Contains(savedNote.Body, addedToNote, "Saved note should include added data")
	})

	// Error tests
	expectedError := errors.New("test error")

	t.Run("error_from_ns", func(t *testing.T) {
		c, ns, _, expectedNote, _ := setupClient("")
		ns.getNoteContent = func(string) (string, error) { return "", expectedError }
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|RawNote)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})

	t.Run("error_from_new_cachefile", func(t *testing.T) {
		c, _, _, expectedNote, _ := setupClient("")
		c.newCacheFile = func(*Client, string) (CacheFile, error) { return nil, expectedError }
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|RawNote)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})

	t.Run("error_from_cachefile_write", func(t *testing.T) {
		c, _, _, expectedNote, _ := setupClient("")
		cache := &mockCacheFile{
			write: func([]byte) (int, error) { return 0, expectedError },
		}
		c.newCacheFile = func(*Client, string) (CacheFile, error) { return cache, nil }
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|RawNote)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})

	t.Run("error_from_cachefile_read", func(t *testing.T) {
		c, _, _, expectedNote, _ := setupClient("")
		cache := &mockCacheFile{
			read:   func([]byte) (int, error) { return 0, expectedError },
			buffer: new(bytes.Buffer),
		}
		c.newCacheFile = func(*Client, string) (CacheFile, error) { return cache, nil }
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|RawNote)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})

	t.Run("error_from_cachefile_close", func(t *testing.T) {
		addedToNote := "<p>New content added</p>"
		c, _, _, expectedNote, _ := setupClient(addedToNote)
		cache := &mockCacheFile{
			close:  func() error { return expectedError },
			buffer: new(bytes.Buffer),
		}
		c.newCacheFile = func(*Client, string) (CacheFile, error) { return cache, nil }
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|RawNote)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})

	t.Run("error_from_cachefile_reopen", func(t *testing.T) {
		addedToNote := "<p>New content added</p>"
		c, _, _, expectedNote, _ := setupClient(addedToNote)
		cache := &mockCacheFile{
			reopen: func() error { return expectedError },
			buffer: new(bytes.Buffer),
		}
		c.newCacheFile = func(*Client, string) (CacheFile, error) { return cache, nil }
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|RawNote)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})

	t.Run("error_from_cachefile_edit", func(t *testing.T) {
		addedToNote := "<p>New content added</p>"
		c, _, _, expectedNote, _ := setupClient(addedToNote)
		c.Editor = &mockEditor{
			edit: func(file CacheFile) error { return expectedError },
		}
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|RawNote)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")
	})

	t.Run("save_recovery_point_if_saves_fails", func(t *testing.T) {
		c, ns, _, expectedNote, _, store := setupClientAndStore("added text")
		ns.updateNote = func(*Note) error { return expectedError }
		var savedNote *Note
		store.saveNoteRecoveryPoint = func(n *Note) error {
			savedNote = n
			return nil
		}

		err := EditNote(c, expectedNote.Title, DefaultNoteOption)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedError, err, "Wrong error returned")

		assert.Equal(expectedNote, savedNote, "Note not saved")
	})

	t.Run("warn_if_recovery_fails", func(t *testing.T) {
		c, ns, _, expectedNote, _, store := setupClientAndStore("added text")
		ns.updateNote = func(*Note) error { return expectedError }
		expectedSaveError := errors.New("recovery error")
		store.saveNoteRecoveryPoint = func(n *Note) error {
			return expectedSaveError
		}

		err := EditNote(c, expectedNote.Title, DefaultNoteOption)
		assert.Error(err, "Should return an error")
		assert.Contains(err.Error(), expectedError.Error(), "Should include notestore error")
		assert.Contains(err.Error(), expectedSaveError.Error(), "Should include recovery point error")
	})

	t.Run("recover_note", func(t *testing.T) {
		c, ns, _, expectedNote, _, store := setupClientAndStore("added text")
		store.saveNoteRecoveryPoint = func(n *Note) error {
			return nil
		}
		store.getNoteRecoveryPoint = func() (*Note, error) {
			return expectedNote, nil
		}
		ns.getNoteContent = func(string) (string, error) { return "", errors.New("should not be called") }

		saveNoteCalled := false
		var savedNote *Note
		ns.updateNote = func(n *Note) error {
			saveNoteCalled = true
			savedNote = n
			return nil
		}
		err := EditNote(c, expectedNote.Title, DefaultNoteOption|UseRecoveryPointNote)
		assert.NoError(err, "Should not return an error")
		assert.True(saveNoteCalled)
		assert.Equal(expectedNote, savedNote, "Wrong note saved")
	})

	t.Run("error_recover_note_if_empty", func(t *testing.T) {
		c, ns, _, expectedNote, _, store := setupClientAndStore("added text")
		expectedNote.GUID = ""
		store.saveNoteRecoveryPoint = func(n *Note) error {
			return nil
		}
		store.getNoteRecoveryPoint = func() (*Note, error) {
			return expectedNote, nil
		}
		ns.getNoteContent = func(string) (string, error) { return "", errors.New("should not be called") }

		err := EditNote(c, expectedNote.Title, DefaultNoteOption|UseRecoveryPointNote)
		assert.Error(err, "Should return an error")
		assert.Equal(ErrNoNoteFound, err, "Wrong error returned")
	})
}

func TestCreateAndEditNewNote(t *testing.T) {
	assert := assert.New(t)
	store := &mockStore{}

	note := &Note{
		Title: "Untitled note",
	}
	ns := nsWithNote(note)
	ns.createNote = func(n *Note) error { return nil }

	var actualFilename string

	client := &Client{
		Config:    &DefaultConfig{},
		Store:     store,
		NoteStore: ns,
		newCacheFile: func(c *Client, filename string) (CacheFile, error) {
			buf := new(bytes.Buffer)
			actualFilename = filename
			return &mockCacheFile{buffer: buf}, nil
		},
		Editor: &mockEditor{
			edit: func(file CacheFile) error {
				return nil
			},
		},
	}
	expectedError := errors.New("expected error")

	t.Run("create_random_file_for_new_note", func(t *testing.T) {
		err := CreateAndEditNewNote(client, note, DefaultNoteOption)
		assert.NoError(err)
		assert.Contains(actualFilename, newNotePrependString)
		// Length of UUID string + length of the prepended string + file extension.
		assert.Len(actualFilename, 36+len(newNotePrependString)+3)
	})

	t.Run("handle_error_from_parsing", func(t *testing.T) {
		client.newCacheFile = func(_ *Client, _ string) (CacheFile, error) {
			return &mockCacheFile{
				read:   func([]byte) (int, error) { return 0, expectedError },
				buffer: new(bytes.Buffer),
			}, nil
		}
		err := CreateAndEditNewNote(client, note, DefaultNoteOption)
		assert.Error(err)
		assert.Equal(expectedError, err)
	})

	t.Run("handle_error_from_edit", func(t *testing.T) {
		client.newCacheFile = func(_ *Client, _ string) (CacheFile, error) { return nil, expectedError }
		err := CreateAndEditNewNote(client, note, DefaultNoteOption)
		assert.Error(err)
		assert.Equal(expectedError, err)
	})
}

func nsWithNote(note *Note) *mockNS {
	notes := []*Note{&Note{Title: "Other note"}, note}
	ns := new(mockNS)
	ns.findNotes = func(filter *NoteFilter, o, max int) ([]*Note, error) { return notes, nil }
	return ns
}
