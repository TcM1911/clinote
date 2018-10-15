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
	"os"
	"os/exec"
)

var (
	// ErrNoEditorFound is returned if no editor was found.
	ErrNoEditorFound = errors.New("no editor found")
)

// Editer is an object that can edit notes.
type Editer interface {
	// Edit allows the user to edit the note.
	Edit(CacheFile) error
}

// VimEditor opens the note in VIM and lets the user edit
// the note.
type VimEditor struct{}

// Edit opens the CacheFile with VIM.
func (e *VimEditor) Edit(file CacheFile) error {
	return executeEditorViaCommand("vim", file.FilePath())
}

func memoryReadLoop(file MemoryCacheFile) chan<- struct{} {
	doneChan := make(chan struct{})
	go func(closeChan <-chan struct{}) {
		for {
			select {
			case <-closeChan:
				return
			}
		}
	}(doneChan)
	return doneChan
}

// EnvEditor opens the note the note using the program defined
// in the environment variable $EDITOR.
type EnvEditor struct{}

// Edit opens the CacheFile with the editor defined in $EDITOR.
func (e *EnvEditor) Edit(file CacheFile) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return ErrNoEditorFound
	}
	return executeEditorViaCommand(editor, file.FilePath())
}

func executeEditorViaCommand(editor, filepath string) error {
	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
