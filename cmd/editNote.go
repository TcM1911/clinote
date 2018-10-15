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

package main

import (
	"fmt"
	"os"

	"github.com/TcM1911/clinote"
	"github.com/spf13/cobra"
)

var editNoteCmd = &cobra.Command{
	Use:   "edit \"note title\"",
	Short: "Edit note.",
	Long: `
Edit allows you to edit the note. If no flags are set, the note is opened
with the editor defined by the environment variable $EDITOR.

The first line will be used as the note title and the rest is encoded as
the note content.

To change to title, the title flag can be used.

The note can be moved to another notebook by defining the new notebook
with the notebook flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			fmt.Println("Error when paring raw flag:", err)
			return
		}
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			fmt.Println("Error parsing the title:", err)
			return
		}
		notebook, err := cmd.Flags().GetString("notebook")
		if err != nil {
			fmt.Println("Error parsing the notebook name:", err)
			return
		}
		recover, err := cmd.Flags().GetBool("recover")
		if err != nil {
			return
		}
		client := defaultClient()
		defer client.Close()
		ns, err := client.GetNoteStore()
		if err != nil {
			fmt.Println("Failed to get notestore:", err)
			return
		}
		opts := clinote.DefaultNoteOption
		if raw {
			opts = opts | clinote.RawNote
		}
		if recover {
			c := clinote.NewClient(client.Config, client.Config.Store(), ns, clinote.DefaultClientOptions)
			err := clinote.EditNote(c, "", opts|clinote.UseRecoveryPointNote)
			if err != nil {
				fmt.Println("Error when edit recovery note:", err)
				os.Exit(1)
			}
			return
		}
		if len(args) != 1 {
			fmt.Println("Error, a note has to be given.")
			return
		}
		if title != "" {
			clinote.ChangeTitle(client.Config.Store(), ns, args[0], title)
		}
		if notebook != "" {
			clinote.MoveNote(client.Config.Store(), ns, args[0], notebook)
		}

		if title == "" && notebook == "" {
			c := clinote.NewClient(client.Config, client.Config.Store(), ns, clinote.DefaultClientOptions)
			err := clinote.EditNote(c, args[0], opts)
			if err != nil {
				fmt.Println("Error when editing the note:", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	noteCmd.AddCommand(editNoteCmd)
	editNoteCmd.Flags().StringP("title", "t", "", "Change the note title to.")
	editNoteCmd.Flags().StringP("notebook", "b", "", "Move the note to notebook.")
	editNoteCmd.Flags().Bool("raw", false, "Use raw content instead of markdown version.")
	editNoteCmd.Flags().Bool("recover", false, "Recover previous note that failed to save.")
}
