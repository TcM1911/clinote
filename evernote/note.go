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
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tcm1911/clinote/markdown"
	"github.com/tcm1911/clinote/user"
	"github.com/tcm1911/evernote-sdk-golang/notestore"
	"github.com/tcm1911/evernote-sdk-golang/types"
)

const (
	// XMLHeader is the header that needs to added to the note content.
	XMLHeader string = `<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">`
)

var (
	NoteFilterOrderCreated        = int32(1)
	NoteFilterOrderUpdated        = int32(2)
	NoteFilterOrderRelevance      = int32(3)
	NoteFilterOrderSequenceNumber = int32(4)
	NoteFilterOrderTitle          = int32(5)
)

var cacheMu sync.Mutex
var cache map[types.GUID]*types.Note

// Note is the structure of an Evernote note.
type Note struct {
	// Title is the note tile.
	Title string
	// GUID is the unique identifier.
	GUID types.GUID
	// Body contains the body of the note.
	Body string `xml:",innerxml"`
	// MD is a Markdown representation of the note body.
	MD string
	// MDHash is the MD5 hash of the MD body.
	MDHash [16]byte
	// Deleted is set true if the note is marked for deletion.
	Deleted bool
	// Notebook the note belongs to.
	Notebook *Notebook
}

// GetNote gets the note metadata in the notebook from the server.
// If the notebook is an empty string, the first matching note will
// be returned.
func GetNote(title, notebook string) *Note {
	ns := user.GetNoteStore()
	filter := notestore.NewNoteFilter()
	if notebook != "" {
		nb, err := findNotebook(notebook)
		if err != nil {
			fmt.Println("Error when getting the notebook:", err)
			os.Exit(1)
		}
		nbGUID := nb.GetGUID()
		filter.NotebookGuid = &nbGUID
	}
	filter.Words = &title
	notes, err := ns.FindNotes(user.AuthToken, filter, 0, 20)
	if err != nil {
		fmt.Println("Error when search for the note:", err)
		os.Exit(1)
	}
	var note *types.Note
	for _, n := range notes.GetNotes() {
		if n.GetTitle() == title {
			note = n
			cacheMu.Lock()
			cache[*n.GUID] = n
			cacheMu.Unlock()
			break
		}
	}
	if note == nil {
		fmt.Println("Could not find a note with title", title)
		os.Exit(1)
	}

	return convert(note)
}

func convert(note *types.Note) *Note {
	n := new(Note)
	n.Title = note.GetTitle()
	n.GUID = note.GetGUID()
	notebook := new(Notebook)
	notebookGUID := types.GUID(note.GetNotebookGuid())
	n.Notebook = notebook
	n.Notebook.GUID = notebookGUID
	return n
}

// GetNoteWithContent returns the note with content from the user's notestore.
func GetNoteWithContent(title string) *Note {
	n := GetNote(title, "")
	ns := user.GetNoteStore()
	content, err := ns.GetNoteContent(user.AuthToken, n.GUID)
	if err != nil {
		fmt.Println("Error when downloading note content:", err)
		os.Exit(1)
	}
	decodeXML(content, n)
	n.MD = markdown.ToHTML(n.Body)
	return n
}

// SaveChanges updates the changes to the note on the server.
func SaveChanges(n *Note) {
	saveChanges(n, true)
}

func ChangeTitle(old, new string) {
	n := GetNote(old, "")
	n.Title = new
	saveChanges(n, false)
}

func MoveNote(noteTitle, notebookName string) {
	n := GetNote(noteTitle, "")
	b, err := FindNotebook(notebookName)
	if err != nil {
		fmt.Println("Error when trying to retrieve notebook:", err)
		return
	}
	n.Notebook.GUID = b.GUID
	saveChanges(n, false)
}

func DeleteNote(title, notebook string) {
	n := GetNote(title, notebook)
	ns := user.GetNoteStore()
	_, err := ns.DeleteNote(user.AuthToken, n.GUID)
	if err != nil {
		fmt.Println("Error when removing the note:", err)
		return
	}
}

func saveChanges(n *Note, updateContent bool) {
	cacheMu.Lock()
	note, ok := cache[n.GUID]
	if !ok {
		// No cached note, so we can't update.
		fmt.Println("Failed to update the changes.")
		return
	}
	// Remove cached note.
	delete(cache, n.GUID)
	cacheMu.Unlock()
	note.Title = &n.Title
	if updateContent {
		xmlBody := toXML(n.MD)
		note.Content = &xmlBody
	}

	notebookGUID := string(n.Notebook.GUID)
	note.NotebookGuid = &notebookGUID

	now := types.Timestamp(time.Now().Unix() * 1000)
	note.Updated = &now
	ns := user.GetNoteStore()
	_, err := ns.UpdateNote(user.AuthToken, note)
	if err != nil {
		fmt.Println("Error when saving the note to server:", err)
		return
	}
}

// SaveNewNote pushes the new note to the server.
func SaveNewNote(n *Note) {
	note := types.NewNote()
	now := types.Timestamp(time.Now().Unix() * 1000)
	note.Created = &now
	note.Title = &n.Title
	if n.MD != "" {
		body := toXML(n.MD)
		note.Content = &body
	} else {
		body := XMLHeader + "<en-note></en-note>"
		note.Content = &body
	}
	if n.Notebook != nil && n.Notebook.Name != "" {
		guid := string(n.Notebook.GUID)
		note.NotebookGuid = &guid
	}
	ns := user.GetNoteStore()
	if _, err := ns.CreateNote(user.AuthToken, note); err != nil {
		fmt.Println("Error when creating the note:", err)
		return
	}
}

func toXML(mdBody string) string {
	b := []byte("")
	content := bytes.NewBuffer(b)
	content.WriteString(XMLHeader)
	content.WriteString("<en-note>")
	content.Write(markdown.ToXML(mdBody))
	content.WriteString("</en-note>")
	return content.String()
}

func decodeXML(content string, v interface{}) {
	d := xml.NewDecoder(strings.NewReader(content))
	d.Strict = false
	d.Entity = xml.HTMLEntity
	d.AutoClose = xml.HTMLAutoClose
	err := d.Decode(&v)
	if err != nil {
		fmt.Println("Error when decoding note content:", err)
		os.Exit(1)
	}
}

func init() {
	cache = make(map[types.GUID]*types.Note)
}
