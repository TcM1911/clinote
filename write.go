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
	"io"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

const timeFormat = "2006-01-02"

var (
	noteListingHeader     = []string{"#", "Title", "Notebook", "Modified", "Created"}
	notebookListingHeader = []string{"#", "Name"}
	credentialHeader      = append(notebookListingHeader, "Type")
	settingsHeader        = []string{"Setting", "Arguments", "Description"}
)

// WriteNoteListing creates and writes a note listing table using the writer.
func WriteNoteListing(w io.Writer, ns []*Note, nbs []*Notebook) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(noteListingHeader)

	for i, n := range ns {
		index := strconv.Itoa(i + 1)
		created := time.Unix(int64(n.Created)/1000, 0).Format(timeFormat)
		modified := time.Unix(int64(n.Updated)/1000, 0).Format(timeFormat)
		notebook := ""
		for _, nb := range nbs {
			if nb.GUID == n.Notebook.GUID {
				notebook = nb.Name
				break
			}
		}
		table.Append([]string{index, n.Title, notebook, modified, created})
	}
	table.Render()
}

// WriteNotebookListing creates and writes a notebook listing table using the writer.
func WriteNotebookListing(w io.Writer, nbs []*Notebook) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(notebookListingHeader)
	for i, nb := range nbs {
		index := strconv.Itoa(i + 1)
		table.Append([]string{index, nb.Name})
	}
	table.Render()
}

// WriteCredentialListing creates and writes a credential listing table using the writer.
func WriteCredentialListing(w io.Writer, creds []*Credential) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(credentialHeader)

	for i, cred := range creds {
		index := strconv.Itoa(i + 1)
		table.Append([]string{index, cred.Name, cred.CredType.String()})
	}
	table.Render()
}

// WriteSettingsListing writes the settings table to writer.
func WriteSettingsListing(w io.Writer, vals, args, desc []string) {
	if len(vals) != len(args) || len(vals) != len(desc) {
		return
	}
	table := tablewriter.NewWriter(w)
	table.SetHeader(settingsHeader)
	for i, val := range vals {
		table.Append([]string{val, args[i], desc[i]})
	}
	table.Render()
}
