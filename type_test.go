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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentialType(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		expected string
		credType CredentialType
	}{
		{"evernote", "Evernote", EvernoteCredential},
		{"evernote_sandbox", "Evernote Sandbox", EvernoteSandboxCredential},
	}
	for _, test := range tests {
		t.Run("Credential type "+test.name, func(t *testing.T) {
			assert.Equal(test.expected, test.credType.String(), "String returns wrong string")
		})
	}
}
