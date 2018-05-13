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
	"log"
	"os"
	"time"

	"github.com/TcM1911/clinote/evernote"
	"github.com/spf13/cobra"
)

const timeFormat = "2006-01-02"

var listNoteCmd = &cobra.Command{
	Use:   "list",
	Short: "List note based on a search filter.",
	Long: `
List returns a list of notes based on a search filter.
The search term flag can be used to define a search term
to be used. The search can be restricted to a notebook
by using the notebook flag.

Count can be used to restrict the maximum number of notes
returned.

If no search term is given, a wild card search will be used.
The notes will be sorted by the modified time.`,
	Run: func(cmd *cobra.Command, args []string) {
		findNotes(cmd, args)
	},
}

func init() {
	noteCmd.AddCommand(listNoteCmd)
	listNoteCmd.Flags().IntP("count", "c", 20, "How many notes to show in the result.")
	listNoteCmd.Flags().StringP("search", "s", "", "Search term.")
	listNoteCmd.Flags().StringP("notebook", "b", "", "Restrict search to notebook.")
}

func findNotes(cmd *cobra.Command, args []string) {
	client := defaultClient()

	// Create filter
	filter := &evernote.NoteFilter{}
	filter.Order = evernote.NoteFilterOrderUpdated
	c, err := cmd.Flags().GetInt("count")
	if err != nil {
		fmt.Println("Error when parsing count value, using default:", err)
		c = 20
	}
	searchBook, err := cmd.Flags().GetString("notebook")
	if err != nil {
		fmt.Println("Error when parsing notebook:", err)
		return
	}
	search, err := cmd.Flags().GetString("search")
	if err != nil {
		fmt.Println("Error when parsing search term", err)
		return
	}

	if search != "" {
		filter.Words = search
	}

	if searchBook != "" {
		book, err := evernote.FindNotebook(client, searchBook)
		if err != nil {
			fmt.Println("Error when trying to filter by notebook: ", err)
			os.Exit(1)
		}
		filter.NotebookGUID = book.GUID
	}

	list, err := evernote.FindNotes(client, filter, 0, c)
	if err != nil {
		log.Fatal(err)
	}

	outputStr := []byte("Search request:.\n")
	outputStr = append(outputStr, []byte(fmt.Sprintf("Found %d items\n", len(list)))...)
	outputStr = append(outputStr, []byte(fmt.Sprintf("%3s : %10s | %10s | %-25s | %-25s\n",
		"#",
		"Created",
		"Updated",
		"Notebook",
		"Title"))...)
	for i, n := range list {
		book, err := evernote.GetNotebook(client, n.Notebook.GUID)
		bookName := ""
		if err != nil {
			log.Println("Error when getting notebook name:", err)

		} else {
			bookName = book.Name
		}
		outputStr = append(outputStr, []byte(fmt.Sprintf("%3d : %10s | %10s | %-25s | %s\n", i+1,
			time.Unix(int64(n.Created)/1000, 0).Format(timeFormat),
			time.Unix(int64(n.Updated)/1000, 0).Format(timeFormat),
			bookName, n.Title))...)
	}

	fmt.Println(string(outputStr))
}
