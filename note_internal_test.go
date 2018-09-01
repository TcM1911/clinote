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

func TestNoteParsing(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		content []byte
	}{
		{"simple", []byte(testContent)},
		{"compact", []byte(compactContent)},
		{"with_white_space", []byte(contentWithWhiteSpace)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n := new(Note)
			r := bytes.NewReader(test.content)
			err := parseNote(r, n, DefaultNoteOption)
			assert.NoError(err, "Should not return an error")
			assert.Equal(noteTitle, n.Title, "Wrong title parsed")
			assert.Equal(noteContent, n.MD, "Wrong content parsed")
		})
	}
}

func TestNoteWriting(t *testing.T) {
	assert := assert.New(t)
	n := &Note{Title: noteTitle, MD: noteContent}
	w := new(bytes.Buffer)

	err := WriteNote(w, n, DefaultNoteOption)
	assert.NoError(err, "Should not fail")
	assert.Equal(testContent, string(w.Bytes()), "Wrong content written")
}

const (
	noteTitle   = "Note title"
	noteContent = "Body\nof\nthe\nnote"
)

const testContent = `---
title: Note title
---
Body
of
the
note`

const compactContent = `---
title:Note title
---
Body
of
the
note`

const contentWithWhiteSpace = `---
title: Note title
---

Body
of
the
note


`
