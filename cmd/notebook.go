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
	"errors"

	"github.com/spf13/cobra"
	"github.com/tcm1911/evernote-sdk-golang/notestore"
	"github.com/tcm1911/evernote-sdk-golang/types"
)

var notebookCmd = &cobra.Command{
	Use:   "notebook",
	Short: "View, create and edit notebooks.",
	Long:  `View, create and edit notebooks.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	RootCmd.AddCommand(notebookCmd)
}

func findNoteBook(ns *notestore.NoteStoreClient, token, bookName string) (*types.Notebook, error) {
	books, err := ns.ListNotebooks(token)
	if err != nil {
		return nil, err
	}
	for _, b := range books {
		if *b.Name == bookName {
			return b, nil
		}
	}
	return nil, errors.New("no matching notebook")
}
