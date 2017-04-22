package evernote

import (
	"time"

	"github.com/tcm1911/clinote/evernote/api"
	"github.com/tcm1911/evernote-sdk-golang/notestore"
	"github.com/tcm1911/evernote-sdk-golang/types"
)

// NotestoreClient is the interface for the notestore.
type NotestoreClient interface {
	// GetClient returns the client for the notestore.
	GetClient() APIClient
	// FindNotes searches for the notes based on the filter.
	FindNotes(filter *NoteFilter, offset, count int) ([]*Note, error)
	// GetAllNotebooks returns all the of users notebooks.
	GetAllNotebooks() ([]*Notebook, error)
	// GetNoteContent gets the note's content from the notestore.
	GetNoteContent(guid string) (string, error)
	// UpdateNote update's the note.
	UpdateNote(note *Note) error
	// DeleteNote removes a note from the user's notebook.
	DeleteNote(guid string) error
	// CreateNote creates a new note on the server.
	CreateNote(note *Note) error
	// UpdateNotebook updates the notebook on the server.
	UpdateNotebook(book *Notebook) error
}

// Notestore is an implementation of the NotestoreClient.
type Notestore struct {
	evernoteNS api.Notestore
	client     APIClient
}

// GetAllNotebooks returns all the of users notebooks.
func (s *Notestore) GetAllNotebooks() ([]*Notebook, error) {
	bs, err := s.evernoteNS.ListNotebooks(s.GetClient().GetAPIToken())
	if err != nil {
		return nil, err
	}
	return convertNotebooks(bs), nil
}

// UpdateNotebook updates the notebook on the server.
func (s *Notestore) UpdateNotebook(b *Notebook) error {
	nb, err := getCachedNotebook(types.GUID(b.GUID))
	if err != nil {
		return err
	}
	transferNotebookData(b, nb)
	_, err = s.evernoteNS.UpdateNotebook(s.GetClient().GetAPIToken(), nb)
	return err
}

//CreateNotebook creates a new notebook for the user.
func (s *Notestore) CreateNotebook(b *Notebook, defaultNotebook bool) error {
	nb := types.NewNotebook()
	nb.DefaultNotebook = &defaultNotebook
	transferNotebookData(b, nb)
	_, err := s.evernoteNS.CreateNotebook(s.client.GetAPIToken(), nb)
	return err
}

// CreateNote creates a new note and saves it to the server.
func (s *Notestore) CreateNote(n *Note) error {
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
	_, err := s.evernoteNS.CreateNote(s.client.GetAPIToken(), note)
	return err
}

// DeleteNote removes a note from the user's notebook.
func (s *Notestore) DeleteNote(guid string) error {
	_, err := s.evernoteNS.DeleteNote(s.client.GetAPIToken(), types.GUID(guid))
	return err
}

// UpdateNote update's the note.
func (s *Notestore) UpdateNote(note *Note) error {
	panic("not implemented")
}

// GetClient returns the client for the notestore.
func (s *Notestore) GetClient() APIClient {
	return s.client
}

// FindNotes searches for the notes based on the filter.
func (s *Notestore) FindNotes(filter *NoteFilter, offset, count int) ([]*Note, error) {
	r, err := s.evernoteNS.FindNotes(s.client.GetAPIToken(), createFilter(filter), int32(offset), int32(count))
	if err != nil {
		return nil, err
	}
	return convertNotes(r.GetNotes()), nil
}

// GetNoteContent gets the note's content from the notestore.
func (s *Notestore) GetNoteContent(guid string) (string, error) {
	panic("not implemented")
}

func createFilter(filter *NoteFilter) *notestore.NoteFilter {
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
