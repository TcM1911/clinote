package clinote

// NotestoreClient is the interface for the notestore.
type NotestoreClient interface {
	// FindNotes searches for the notes based on the filter.
	FindNotes(filter *NoteFilter, offset, count int) ([]*Note, error)
	// GetAllNotebooks returns all the of users notebooks.
	GetAllNotebooks() ([]*Notebook, error)
	// GetNotebook
	GetNotebook(guid string) (*Notebook, error)
	// CreateNotebook
	CreateNotebook(b *Notebook, defaultNotebook bool) error
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
