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
	"os"
	"os/exec"
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
	cmd := exec.Command("vim", file.FilePath())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return nil
}
