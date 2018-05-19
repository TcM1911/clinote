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
 * Copyright (C) Joakim Kennedy, 2016-2017
 */

package evernote

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/TcM1911/clinote/markdown"
)

const (
	// XMLHeader is the header that needs to added to the note content.
	XMLHeader string = `<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">`
)

var (
	// NoteFilterOrderCreated sorts the notes by create time.
	NoteFilterOrderCreated = int32(1)
	// NoteFilterOrderUpdated sorts the notes by update time.
	NoteFilterOrderUpdated = int32(2)
	// NoteFilterOrderRelevance sorts the notes by relevance.
	NoteFilterOrderRelevance = int32(3)
	// NoteFilterOrderSequenceNumber sorts the notes by sequence number.
	NoteFilterOrderSequenceNumber = int32(4)
	// NoteFilterOrderTitle sorts the notes by title.
	NoteFilterOrderTitle = int32(5)
)

var (
	// ErrNoNoteFound is returned if search resulted in no notes found.
	ErrNoNoteFound = errors.New("no note found")
)

// Note is the structure of an Evernote note.
type Note struct {
	// Title is the note tile.
	Title string
	// GUID is the unique identifier.
	GUID string
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
	// Created
	Created int64
	// Updated
	Updated int64
}

// NoteFilter is the search filter for notes.
type NoteFilter struct {
	// NotebookGUID is the GUID for the notebook to limit the search to.
	NotebookGUID string
	// Words can be a search string or note title.
	Words string
	// Order
	Order int32
}

// FindNotes searches for notes.
func FindNotes(client APIClient, filter *NoteFilter, offset int, count int) ([]*Note, error) {
	ns, err := client.GetNoteStore()
	if err != nil {
		return nil, err
	}
	return ns.FindNotes(filter, offset, count)
}

// GetNote gets the note metadata in the notebook from the server.
// If the notebook is an empty string, the first matching note will
// be returned.
func GetNote(client APIClient, title, notebook string) (*Note, error) {
	ns, err := client.GetNoteStore()
	if err != nil {
		return nil, err
	}
	filter := new(NoteFilter)
	if notebook != "" {
		nb, err := findNotebook(ns, notebook)
		if err != nil {
			return nil, err
		}
		filter.NotebookGUID = nb.GUID
	}
	filter.Words = title
	notes, err := ns.FindNotes(filter, 0, 20)
	if err != nil {
		return nil, err
	}
	var note *Note
	for _, n := range notes {
		if n.Title == title {
			note = n
			break
		}
	}
	if note == nil {
		return nil, ErrNoNoteFound
	}
	return note, nil
}

// GetNoteWithContent returns the note with content from the user's notestore.
func GetNoteWithContent(client APIClient, title string) (*Note, error) {
	n, err := GetNote(client, title, "")
	ns, err := client.GetNoteStore()
	if err != nil {
		return nil, err
	}
	content, err := ns.GetNoteContent(n.GUID)
	if err != nil {
		return nil, err
	}
	err = decodeXML(content, n)
	if err != nil {
		return nil, err
	}
	n.MD, err = markdown.FromHTML(n.Body)
	if err != nil {
		return nil, err
	}
	return n, nil
}

// SaveChanges updates the changes to the note on the server.
func SaveChanges(client APIClient, n *Note, useRawContent bool) error {
	// useRawContent := GetUseRawContentFromContext(ctx)
	return saveChanges(client, n, true, useRawContent)
}

// ChangeTitle changes the note's title.
func ChangeTitle(client APIClient, old, new string) error {
	n, err := GetNote(client, old, "")
	if err != nil {
		return err
	}
	n.Title = new
	return saveChanges(client, n, false, false)
}

// MoveNote moves the note to a new notebook.
func MoveNote(client APIClient, noteTitle, notebookName string) error {
	n, err := GetNote(client, noteTitle, "")
	if err != nil {
		return err
	}
	b, err := FindNotebook(client, notebookName)
	if err != nil {
		return err
	}
	n.Notebook = b
	return saveChanges(client, n, false, false)
}

// DeleteNote moves a note from the notebook to the trash can.
func DeleteNote(client APIClient, title, notebook string) error {
	n, err := GetNote(client, title, notebook)
	if err != nil {
		return err
	}
	ns, err := client.GetNoteStore()
	if err != nil {
		return err
	}
	err = ns.DeleteNote(n.GUID)
	if err != nil {
		return err
	}
	return nil
}

func saveChanges(client APIClient, n *Note, updateContent, useRawContent bool) error {
	ns, err := client.GetNoteStore()
	if err != nil {
		return err
	}
	if updateContent {
		body := toXML(n.MD)
		if useRawContent {
			body = fmt.Sprintf("%s<en-note>%s</en-note>", XMLHeader, n.Body)
		}
		n.Body = body
	}
	err = ns.UpdateNote(n)
	if err != nil {
		return err
	}
	return nil
}

// SaveNewNote pushes the new note to the server.
func SaveNewNote(client APIClient, n *Note, raw bool) error {
	var body string
	if !raw && n.MD != "" {
		body = toXML(n.MD)
	} else if raw {
		body = fmt.Sprintf("%s<en-note>%s</en-note>", XMLHeader, n.Body)
	} else {
		body = XMLHeader + "<en-note></en-note>"
	}
	n.Body = body
	ns, err := client.GetNoteStore()
	if err != nil {
		return err
	}
	if err = ns.CreateNote(n); err != nil {
		return err
	}
	return nil
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

func decodeXML(content string, v interface{}) error {
	d := xml.NewDecoder(strings.NewReader(content))
	d.Strict = false
	d.Entity = xml.HTMLEntity
	d.AutoClose = xml.HTMLAutoClose
	return d.Decode(&v)
}
