/*/*
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

	"github.com/TcM1911/clinote/evernote"
	"github.com/spf13/cobra"
)

var deleteNoteCmd = &cobra.Command{
	Use:   "delete \"note title\"",
	Short: "Delete note.",
	Long: `Moves the note into the trash. The note may still be undeleted, unless it is expunged.
To expunge the note you need to use the official client or the web client.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Error, a note title has to be given")
			return
		}
		nb, err := cmd.Flags().GetString("notebook")
		if err != nil {
			fmt.Println("Error when parsing the notebook name:", err)
			return
		}
		client := defaultClient()
		defer client.Close()
		err = evernote.DeleteNote(client, args[0], nb)
		if err != nil {
			fmt.Println("Error when deleting the note:", err)
			os.Exit(1)
		}
	},
}

func init() {
	noteCmd.AddCommand(deleteNoteCmd)
	deleteNoteCmd.Flags().StringP("notebook", "b", "", "The notebook of the note.")
}
