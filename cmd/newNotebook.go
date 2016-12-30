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

	"github.com/spf13/cobra"
	"github.com/tcm1911/clinote/user"
	"github.com/tcm1911/evernote-sdk-golang/types"
)

var newBookCmd = &cobra.Command{
	Use:   "new \"notebook name\"",
	Short: "Create a new notebook.",
	Long: `
New creates a new notebook.`,
	Run: func(cmd *cobra.Command, args []string) {
		createNotebook(cmd, args)
	},
}

func init() {
	notebookCmd.AddCommand(newBookCmd)
	newBookCmd.Flags().StringP("stack", "s", "", "Add notebook to stack.")
	newBookCmd.Flags().BoolP("default", "d", false, "If notebook should be set to the default notebook.")
}

func createNotebook(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("No notebook name given")
		os.Exit(1)
	}
	nb := types.NewNotebook()
	nb.Name = &args[0]

	stack, err := cmd.Flags().GetString("stack")
	if err != nil {
		fmt.Println("Error when parsing stack name:", err)
		os.Exit(1)
	}
	if stack != "" {
		nb.Stack = &stack
	}

	d, err := cmd.Flags().GetBool("default")
	if err != nil {
		fmt.Println("Error when parsing default value:", err)
		os.Exit(1)
	}
	nb.DefaultNotebook = &d

	ns := user.GetNoteStore()

	_, err = ns.CreateNotebook(user.AuthToken, nb)
	if err != nil {
		fmt.Println("Error when creating the notebook:", err)
		os.Exit(1)
	}
}
