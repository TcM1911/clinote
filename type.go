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

import "io"

// Storager is the interface for backend storage.
type Storager interface {
	io.Closer
	GetSettings() (*Settings, error)
	StoreSettings(*Settings) error
	// GetNotebookCache returns the stored NotebookCacheList.
	GetNotebookCache() (*NotebookCacheList, error)
	// StoreNotebookList saves the list to the database.
	StoreNotebookList(list *NotebookCacheList) error
	// SaveSearch stores a note search to the database.
	SaveSearch([]*Note) error
	// GetSearch returns a saved note search from the database.
	GetSearch() ([]*Note, error)
	// SaveNoteRecoveryPoint saves the note as a recovery point.
	SaveNoteRecoveryPoint(*Note) error
	// GetNoteREcoveryPoint returns the saved note.
	GetNoteRecoveryPoint() (*Note, error)
}

// UserCredentialStore provides an interface to a backend that stores
// a user's credentials.
type UserCredentialStore interface {
	// Add saves a new credential to the store.
	Add(*Credential) error
	// Remove deletes a credential from the store.
	Remove(*Credential) error
	// GetAll returns all credentials in the store.
	GetAll() ([]*Credential, error)
	// GetByIndex returns a user credential by its index.
	GetByIndex(index int) (*Credential, error)
}

// Settings is a struct holding the user's settings for the application.
type Settings struct {
	// APIKey is the user's session key.
	APIKey string
}

// Credential is a struct that holds credential information.
type Credential struct {
	// Name is the user given name for the credential.
	Name string
	// Secret is used to authenticate to the note store.
	Secret string
	// CredType is used to identify credential type.
	CredType CredentialType
}

// CredentialType is a type of credential. Used to identify which backend to use
type CredentialType uint8

func (c CredentialType) String() string {
	return credtypeStringMapper[c]
}

const (
	// EvernoteCredential is used for credentials that can authenticate with Evernote.
	EvernoteCredential CredentialType = iota
	// EvernoteSandboxCredential is used for credentials that can authenticate with
	// Evernote's sandbox server.
	EvernoteSandboxCredential
)

var credtypeStringMapper = []string{"Evernote", "Evernote Sandbox"}
