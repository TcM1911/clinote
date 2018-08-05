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

var listNotebooksCmd = &cobra.Command{
	Use:   "list",
	Short: "List notebooks.",
	Long: `
List notebooks returns all active notebooks.`,
	Run: func(cmd *cobra.Command, args []string) {
		sync, err := cmd.Flags().GetBool("sync")
		if err != nil {
			fmt.Println(err)
			return
		}
		listNotebooks(sync)
	},
}

func init() {
	notebookCmd.AddCommand(listNotebooksCmd)
	listNotebooksCmd.Flags().BoolP("sync", "s", false, "Force a resync of notebooks from the server.")
}

func listNotebooks(sync bool) {
	client := defaultClient()
	defer client.Close()
	ns, err := client.GetNoteStore()
	if err != nil {
		return
	}
	bs, err := clinote.GetNotebooks(client.Config.Store(), ns, sync)
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
