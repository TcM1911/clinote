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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWritingNoteAndNotebookTables(t *testing.T) {
	assert := assert.New(t)
	nbs := []*Notebook{
		&Notebook{GUID: "GUID1", Name: "Notebook1"},
		&Notebook{GUID: "GUID2", Name: "Notebook2"},
		&Notebook{GUID: "GUID3", Name: "Notebook3"},
	}
	notes := []*Note{
		&Note{Title: "Note1", Notebook: &Notebook{GUID: "GUID1"}, Created: int64(0), Updated: int64(0)},
		&Note{Title: "Note2", Notebook: &Notebook{GUID: "GUID2"}, Created: int64(0), Updated: int64(0)},
		&Note{Title: "Note3", Notebook: &Notebook{GUID: "GUID3"}, Created: int64(0), Updated: int64(0)},
	}

	t.Run("NotebookList", func(t *testing.T) {
		buf := new(bytes.Buffer)
		WriteNotebookListing(buf, nbs)
		assert.Equal(expectedNotebooklist, string(buf.Bytes()), "Notebook list table doesn't match")
	})

	t.Run("NoteList", func(t *testing.T) {
		buf := new(bytes.Buffer)
		WriteNoteListing(buf, notes, nbs)
		assert.Equal(expectedNotelist, string(buf.Bytes()), "Note list table doesn't match")
	})
}

func TestCredentialTable(t *testing.T) {
	assert := assert.New(t)
	creds := []*Credential{
		&Credential{Name: "Cred1", Secret: "test12", CredType: EvernoteCredential},
		&Credential{Name: "Cred2", Secret: "test23", CredType: EvernoteSandboxCredential},
	}
	t.Run("print without secret", func(t *testing.T) {
		buf := new(bytes.Buffer)
		WriteCredentialListing(buf, creds)
		assert.Equal(expectedCredentialList, string(buf.Bytes()))
	})

	t.Run("print with secret", func(t *testing.T) {
		buf := new(bytes.Buffer)
		WriteCredentialListingWithSecret(buf, creds)
		assert.Equal(expectedCredentialListWithSecret, string(buf.Bytes()))
	})
}

func TestSettingsTable(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)

	WriteSettingsListing(buf, []string{"credential"}, []string{"An index value."}, []string{"Set the active credential for the user."})

	assert.Equal(expectedSettingList, string(buf.Bytes()))
}

const expectedNotebooklist = `+---+-----------+
| # |   NAME    |
+---+-----------+
| 1 | Notebook1 |
| 2 | Notebook2 |
| 3 | Notebook3 |
+---+-----------+
`
const expectedNotelist = `+---+-------+-----------+------------+------------+
| # | TITLE | NOTEBOOK  |  MODIFIED  |  CREATED   |
+---+-------+-----------+------------+------------+
| 1 | Note1 | Notebook1 | 1970-01-01 | 1970-01-01 |
| 2 | Note2 | Notebook2 | 1970-01-01 | 1970-01-01 |
| 3 | Note3 | Notebook3 | 1970-01-01 | 1970-01-01 |
+---+-------+-----------+------------+------------+
`
const expectedCredentialList = `+---+-------+------------------+
| # | NAME  |       TYPE       |
+---+-------+------------------+
| 1 | Cred1 | Evernote         |
| 2 | Cred2 | Evernote Sandbox |
+---+-------+------------------+
`

const expectedCredentialListWithSecret = `+---+-------+------------------+--------+
| # | NAME  |       TYPE       | SECRET |
+---+-------+------------------+--------+
| 1 | Cred1 | Evernote         | test12 |
| 2 | Cred2 | Evernote Sandbox | test23 |
+---+-------+------------------+--------+
`

const expectedSettingList = `+------------+-----------------+--------------------------------+
|  SETTING   |    ARGUMENTS    |          DESCRIPTION           |
+------------+-----------------+--------------------------------+
| credential | An index value. | Set the active credential for  |
|            |                 | the user.                      |
+------------+-----------------+--------------------------------+
`
