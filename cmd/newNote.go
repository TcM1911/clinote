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
 * Copyright (C) Joakim Kennedy, 2016
 */

package cmd

import (
	"fmt"

	"github.com/TcM1911/clinote/evernote"
	"github.com/spf13/cobra"
)

var newNoteCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new note.",
	Long: `
New creates a new note. A title needs to be given for the
note.

If no notebook is given, the default notebook will be used.

The new note can be open in the $EDITOR by using the edit
flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			fmt.Println("Error when parsing note title:", err)
			return
		}
		edit, err := cmd.Flags().GetBool("edit")
		if err != nil {
			fmt.Println("Error when parsing edit flag:", err)
			return
		}
		if title == "" && !edit {
			fmt.Println("Note title has to be given")
			return
		}
		notebook, err := cmd.Flags().GetString("notebook")
		if err != nil {
			fmt.Println("Error when parsing notebook name:", err)
			return
		}
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			fmt.Println("Error when parsing raw parameter:", err)
			return
		}
		createNote(title, notebook, edit, raw)
	},
}

func init() {
	noteCmd.AddCommand(newNoteCmd)
	newNoteCmd.Flags().StringP("title", "t", "", "Note title.")
	newNoteCmd.Flags().StringP("notebook", "b", "", "The notebook to save note to, if not set the default notebook will be used.")
	newNoteCmd.Flags().BoolP("edit", "e", false, "Open note in the editor.")
	newNoteCmd.Flags().Bool("raw", false, "Edit the content in raw mode.")
}

func createNote(title, notebook string, edit, raw bool) {
	note := new(evernote.Note)
	if edit {
		var t string
		if title == "" {
			t = "<TITLE>"
		} else {
			t = title
		}
		var b []byte
		var err error
		if raw {
			b, err = createTmpFileAndEdit("new-note.xml", t, "")
		} else {
			b, err = createTmpFileAndEdit("new-note.md", t, "<CONTENT>")
		}
		if err != nil {
			fmt.Println("Error when processing note:", err)
			return
		}
		if err = parseFileChange(b, note, raw); err != nil {
			fmt.Println("Error processing changes:", err)
			return
		}
	} else {
		note.Title = title
	}

	client := defaultClient()
	if notebook != "" {
		nb, err := evernote.FindNotebook(client, notebook)
		if err != nil {
			fmt.Println("Error when searching for notebook:", err)
			return
		}
		note.Notebook = nb
	}
	evernote.SaveNewNote(client, note, raw)
}
