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

	"github.com/spf13/cobra"
	"github.com/tcm1911/clinote/evernote"
)

var noteCmd = &cobra.Command{
	Use:   "note \"note title\"",
	Short: "View, edit and create a note.",
	Long:  `Displays the content of a note.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			return
		}
		getNote(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(noteCmd)
	noteCmd.Flags().Bool("raw", false, "Display raw content instead of markdown encoded.")
}

func getNote(cmd *cobra.Command, args []string) {
	name := args[0]
	raw, err := cmd.Flags().GetBool("raw")
	if err != nil {
		fmt.Println("Error when paring raw flag:", err)
		return
	}
	n := evernote.GetNoteWithContent(name)
	if raw {
		fmt.Println(n.Title, "\n\n", n.Body)
	} else {
		fmt.Println(n.Title, "\n\n", n.MD)
	}
}
