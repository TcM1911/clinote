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

var editNotebookCmd = &cobra.Command{
	Use:   "edit \"notebook name\"",
	Short: "Edit a notebook.",
	Long: `
Edit a notebook. The notebook's name can be changed using the
name flag.

To move the notebook to another stack, use the stack flag to
define the new stack.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Error, a notebook has to be given.")
			return
		}
		change := false
		notebook := new(evernote.Notebook)
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println("Error when parsing new notebook name:", err)
			return
		}
		if name != "" {
			notebook.Name = name
			change = true
		}

		stack, err := cmd.Flags().GetString("stack")
		if err != nil {
			fmt.Println("Error when parsing the new stack:", err)
			return
		}
		if stack != "" {
			notebook.Stack = stack
			change = true
		}

		if !change {
			fmt.Println("No changes detected, aborting.")
			return
		}
		evernote.UpdateNotebook(args[0], notebook)
	},
}

func init() {
	notebookCmd.AddCommand(editNotebookCmd)
	editNotebookCmd.Flags().StringP("name", "n", "", "Change notebook name to.")
	editNotebookCmd.Flags().StringP("stack", "s", "", "Change notebook stack to.")
}
