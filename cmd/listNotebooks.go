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
	"os"

	"github.com/TcM1911/clinote/evernote"
	"github.com/spf13/cobra"
)

var listNotebooksCmd = &cobra.Command{
	Use:   "list",
	Short: "List notebooks.",
	Long: `
List notebooks returns all active notebooks.`,
	Run: func(cmd *cobra.Command, args []string) {
		listNotebooks()
	},
}

func init() {
	notebookCmd.AddCommand(listNotebooksCmd)
}

func listNotebooks() {
	bs, err := evernote.GetNotebooks()
	if err != nil {
		fmt.Println("Error when getting notebooks:", err)
		os.Exit(1)
	}
	var output []byte

	for i, b := range bs {
		output = append(output, []byte(fmt.Sprintf("%d : %s\n", i+1, b.Name))...)
	}
	fmt.Println(string(output))
}
