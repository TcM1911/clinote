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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCredential(t *testing.T) {
	assert := assert.New(t)
	store := new(mockCredentialStore)
	expected := &Credential{CredType: EvernoteSandboxCredential, Name: "cred name", Secret: "cred secret"}
	var saved *Credential
	store.add = func(c *Credential) error { saved = c; return nil }

	err := AddNewCredential(store, expected.Name, expected.Secret, expected.CredType)
	assert.NoError(err, "No error")
	assert.Equal(*expected, *saved, "Wrong data saved")
}

func TestRemoveCredential(t *testing.T) {
	assert := assert.New(t)
	expected := &Credential{Name: "Cred3"}
	list := []*Credential{
		&Credential{Name: "Cred1"},
		&Credential{Name: "Cred2"},
		expected,
	}
	store := new(mockCredentialStore)
	store.getAll = func() ([]*Credential, error) {
		return list, nil
	}
	var removed *Credential
	store.remove = func(c *Credential) error { removed = c; return nil }

	t.Run("remove match", func(t *testing.T) {
		err := RemoveCredential(store, "Cred3")
		assert.NoError(err, "Should not return an error")
		assert.Equal(*expected, *removed)
	})

	t.Run("error if no match", func(t *testing.T) {
		err := RemoveCredential(store, "Cred4")
		assert.Error(err, "Should return an error")
		assert.Equal(ErrNoMatchingCredentialFound, err)
	})

	t.Run("error from store", func(t *testing.T) {
		expectedErr := errors.New("error")
		store.getAll = func() ([]*Credential, error) { return nil, expectedErr }
		err := RemoveCredential(store, "")
		assert.Error(err, "Should return an error")
		assert.Equal(expectedErr, err)

		err = RemoveCredentialByIndex(store, 2)
		assert.Error(err, "Should return an error")
		assert.Equal(expectedErr, err)
	})

	store.remove = func(c *Credential) error { removed = c; return nil }
	store.getAll = func() ([]*Credential, error) { return list, nil }

	t.Run("remove by index", func(t *testing.T) {
		err := RemoveCredentialByIndex(store, 2)
		assert.NoError(err, "Should not return an error")
		assert.Equal(*expected, *removed)
	})

	t.Run("negative index", func(t *testing.T) {
		err := RemoveCredentialByIndex(store, -1)
		assert.Equal(ErrNegativeIndex, err)
	})

	t.Run("index to big", func(t *testing.T) {
		err := RemoveCredentialByIndex(store, 10)
		assert.Equal(ErrIndexToBig, err)
	})
}
