package evernote

import (
	"time"

	"github.com/TcM1911/clinote"
	"github.com/TcM1911/clinote/evernote/api"
	"github.com/TcM1911/evernote-sdk-golang/notestore"
	"github.com/TcM1911/evernote-sdk-golang/types"
)

// Notestore is an implementation of the NotestoreClient.
type Notestore struct {
	evernoteNS api.Notestore
	apiToken   string
}

// GetAllNotebooks returns all the of users notebooks.
func (s *Notestore) GetAllNotebooks() ([]*clinote.Notebook, error) {
	bs, err := s.evernoteNS.ListNotebooks(s.apiToken)
	if err != nil {
		return nil, err
	}
	return convertNotebooks(bs), nil
}

// UpdateNotebook updates the notebook on the server.
func (s *Notestore) UpdateNotebook(b *clinote.Notebook) error {
	nb, err := getCachedNotebook(types.GUID(b.GUID))
	if err != nil {
		return err
	}
	transferNotebookData(b, nb)
	_, err = s.evernoteNS.UpdateNotebook(s.apiToken, nb)
	return err
}

//CreateNotebook creates a new notebook for the user.
func (s *Notestore) CreateNotebook(b *clinote.Notebook, defaultNotebook bool) error {
	nb := types.NewNotebook()
	nb.DefaultNotebook = &defaultNotebook
	transferNotebookData(b, nb)
	_, err := s.evernoteNS.CreateNotebook(s.apiToken, nb)
	return err
}

// GetNotebook returns the notebook with the specific GUID.
func (s *Notestore) GetNotebook(guid string) (*clinote.Notebook, error) {
	nb, err := s.evernoteNS.GetNotebook(s.apiToken, types.GUID(guid))
	if err != nil {
		return nil, err
	}
	return convertNotebooks([]*types.Notebook{nb})[0], nil
}

// CreateNote creates a new note and saves it to the server.
func (s *Notestore) CreateNote(n *clinote.Note) error {
	note := types.NewNote()
	now := types.Timestamp(time.Now().Unix() * 1000)
	note.Created = &now
	note.Title = &n.Title
	if n.Body != "" {
		note.Content = &n.Body
	}
	if n.Notebook != nil && n.Notebook.Name != "" {
		guid := string(n.Notebook.GUID)
		note.NotebookGuid = &guid
	}
	_, err := s.evernoteNS.CreateNote(s.apiToken, note)
	return err
}

// DeleteNote removes a note from the user's notebook.
func (s *Notestore) DeleteNote(guid string) error {
	_, err := s.evernoteNS.DeleteNote(s.apiToken, types.GUID(guid))
	return err
}

// UpdateNote update's the note.
func (s *Notestore) UpdateNote(note *clinote.Note) error {
	if note.GUID == "" {
		return ErrNoGUIDSet
	}
	if note.Title == "" {
		return ErrNoTitleSet
	}
	n := types.NewNote()
	n.Title = &note.Title
	guid := types.GUID(note.GUID)
	n.GUID = &guid
	if note.Body != "" {
		n.Content = &note.Body
	}
	_, err := s.evernoteNS.UpdateNote(s.apiToken, n)
	return err
}

// FindNotes searches for the notes based on the filter.
func (s *Notestore) FindNotes(filter *clinote.NoteFilter, offset, count int) ([]*clinote.Note, error) {
	r, err := s.evernoteNS.FindNotes(s.apiToken, createFilter(filter), int32(offset), int32(count))
	if err != nil {
		return nil, err
	}
	return convertNotes(r.GetNotes()), nil
}

// GetNoteContent gets the note's content from the notestore.
func (s *Notestore) GetNoteContent(guid string) (string, error) {
	return s.evernoteNS.GetNoteContent(s.apiToken, types.GUID(guid))
}

func createFilter(filter *clinote.NoteFilter) *notestore.NoteFilter {
	searchFilter := notestore.NewNoteFilter()
	if filter.NotebookGUID != "" {
		guid := types.GUID(filter.NotebookGUID)
		searchFilter.NotebookGuid = &guid
	}
	if filter.Words != "" {
		searchFilter.Words = &(filter.Words)
	}
	return searchFilter
}
