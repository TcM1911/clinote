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
	"bytes"
)

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

type mockStore struct {
	getNotebookCache      func() (*NotebookCacheList, error)
	storeNotebookList     func(list *NotebookCacheList) error
	getSearch             func() ([]*Note, error)
	saveNoteRecoveryPoint func(*Note) error
	getNoteRecoveryPoint  func() (*Note, error)
}

func (m *mockStore) SaveNoteRecoveryPoint(n *Note) error {
	return m.saveNoteRecoveryPoint(n)
}

func (m *mockStore) GetNoteRecoveryPoint() (*Note, error) {
	return m.getNoteRecoveryPoint()
}

func (m *mockStore) SaveSearch([]*Note) error {
	panic("not implemented")
}

func (m *mockStore) GetSearch() ([]*Note, error) {
	return m.getSearch()
}

func (m *mockStore) Close() error {
	panic("not implemented")
}

func (m *mockStore) GetSettings() (*Settings, error) {
	panic("not implemented")
}

func (m *mockStore) StoreSettings(*Settings) error {
	panic("not implemented")
}

func (m *mockStore) GetNotebookCache() (*NotebookCacheList, error) {
	return m.getNotebookCache()
}

func (m *mockStore) StoreNotebookList(list *NotebookCacheList) error {
	return m.storeNotebookList(list)
}

type mockEditor struct {
	edit func(CacheFile) error
}

func (m *mockEditor) Edit(file CacheFile) error {
	return m.edit(file)
}

type mockCacheFile struct {
	buffer *bytes.Buffer
	write  func([]byte) (int, error)
	read   func([]byte) (int, error)
	close  func() error
	reopen func() error
}

func (m *mockCacheFile) Read(p []byte) (n int, err error) {
	if m.read != nil {
		return m.read(p)
	}
	return m.buffer.Read(p)
}

func (m *mockCacheFile) Write(p []byte) (n int, err error) {
	if m.write != nil {
		return m.write(p)
	}
	return m.buffer.Write(p)
}

func (m *mockCacheFile) Close() error {
	if m.close != nil {
		return m.close()
	}
	return nil
}

func (m *mockCacheFile) FilePath() string {
	return ""
}

func (m *mockCacheFile) ReOpen() error {
	if m.reopen != nil {
		return m.reopen()
	}
	return nil
}

func (m *mockCacheFile) CloseAndRemove() error {
	return nil
}

type mockCredentialStore struct {
	add        func(*Credential) error
	remove     func(*Credential) error
	getAll     func() ([]*Credential, error)
	getByIndex func(int) (*Credential, error)
}

func (m *mockCredentialStore) Add(c *Credential) error {
	return m.add(c)
}

func (m *mockCredentialStore) Remove(c *Credential) error {
	return m.remove(c)
}

func (m *mockCredentialStore) GetAll() ([]*Credential, error) {
	return m.getAll()
}

func (m *mockCredentialStore) GetByIndex(index int) (*Credential, error) {
	return m.getByIndex(index)
}
