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

import "errors"

var (
	// ErrNoMatchingCredentialFound is returned if no matching credential is found.
	ErrNoMatchingCredentialFound = errors.New("no matching credential found")
	// ErrNegativeIndex is returned if index is negative
	ErrNegativeIndex = errors.New("index negative")
	// ErrIndexToBig is returned if index is greater than the list.
	ErrIndexToBig = errors.New("index greater than the array")
)

// AddNewCredential creates and add a new credential to the store.
func AddNewCredential(store UserCredentialStore, name, secret string, credType CredentialType) error {
	cred := &Credential{
		Name:     name,
		Secret:   secret,
		CredType: credType,
	}
	return store.Add(cred)
}

// RemoveCredential removes the first credential with the matching name.
func RemoveCredential(store UserCredentialStore, name string) error {
	creds, err := GetAllCredentials(store)
	if err != nil {
		return err
	}
	for _, cred := range creds {
		if cred.Name == name {
			return store.Remove(cred)
		}
	}
	return ErrNoMatchingCredentialFound
}

// RemoveCredentialByIndex removes the credential at the index provided.
func RemoveCredentialByIndex(store UserCredentialStore, index int) error {
	cred, err := GetCredential(store, index)
	if err != nil {
		return err
	}
	return store.Remove(cred)
}

// GetAllCredentials returns all the credentials in the store.
func GetAllCredentials(store UserCredentialStore) ([]*Credential, error) {
	return store.GetAll()
}

// GetCredential returns the credential with the matching index.
func GetCredential(store UserCredentialStore, index int) (*Credential, error) {
	if index < 0 {
		return nil, ErrNegativeIndex
	}
	creds, err := GetAllCredentials(store)
	if err != nil {
		return nil, err
	}
	if index > len(creds)-1 {
		return nil, ErrIndexToBig
	}
	return creds[index], nil
}
