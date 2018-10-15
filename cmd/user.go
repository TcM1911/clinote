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
 * Copyright (C) Joakim Kennedy, 2016-2018
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/TcM1911/clinote"
	"github.com/TcM1911/clinote/storage"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User functionality.",
	Long:  `User functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	RootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userRmCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userSetCmd)
	// Add flags
	userAddCmd.Flags().StringP("name", "n", "", "Username")
	userAddCmd.Flags().StringP("secret", "s", "", "Access token")
	userAddCmd.Flags().Bool("sandbox", false, "Use Evernote's Sandbox instance")
	// List flags
	userListCmd.Flags().Bool("show-secret", false, "Include credential secret in the output")
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all credentials",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := storage.Open((new(clinote.DefaultConfig)).GetConfigFolder())
		if err != nil {
			fmt.Println("Error when opening the database:", err.Error())
		}
		listCredentials(db, cmd)
	},
}

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new credential",
	Long:  "Add a new credential set for the user. Please follow the instructions on https://dev.evernote.com/doc/articles/dev_tokens.php to generate access tokens.",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := storage.Open((new(clinote.DefaultConfig)).GetConfigFolder())
		if err != nil {
			fmt.Println("Error when opening the database:", err.Error())
		}
		addCredential(db, cmd, args)
	},
}

var userRmCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a credential",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := storage.Open((new(clinote.DefaultConfig)).GetConfigFolder())
		if err != nil {
			fmt.Println("Error when opening the database:", err.Error())
		}
		rmCredential(db, args)
	},
}

var userSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a user configuration",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := storage.Open((new(clinote.DefaultConfig)).GetConfigFolder())
		if err != nil {
			fmt.Println("Error when opening the database:", err.Error())
		}
		setConfig(db, db, args)
	},
}

var setConfigOpts = []struct {
	val  string
	args string
	desc string
}{
	{"credential", "An index value.", "Set the active credential for the user."},
}

func setConfig(store clinote.UserCredentialStore, db clinote.Storager, args []string) {
	if len(args) != 2 {
		printConfigOptions()
		return
	}
	switch args[0] {
	case "credential":
		setCredential(store, db, args[1])
	default:
		printConfigOptions()
	}
}

func setCredential(store clinote.UserCredentialStore, db clinote.Storager, strIndex string) {
	index, err := strconv.Atoi(strIndex)
	if err != nil {
		fmt.Printf("%s is not a number\n", strIndex)
		return
	}
	creds, err := clinote.GetAllCredentials(store)
	if err != nil {
		fmt.Println("Error when getting credential list:", err)
		return
	}
	// Index is a 1 based index for the user.
	if index < 1 || index > len(creds) {
		fmt.Println("Error index out-of-range")
		return
	}
	settings, err := db.GetSettings()
	if err != nil {
		fmt.Println("Error when getting the settings:", err)
		return
	}
	settings.APIKey = creds[index-1].Secret
	settings.Credential = creds[index-1]
	err = db.StoreSettings(settings)
	if err != nil {
		fmt.Println("Error when saving the settings:", err)
	}
}

func printConfigOptions() {
	n := len(setConfigOpts)
	vals, args, descs := make([]string, n, n), make([]string, n, n), make([]string, n, n)
	for i, cfg := range setConfigOpts {
		vals[i] = cfg.val
		args[i] = cfg.args
		descs[i] = cfg.desc
	}
	clinote.WriteSettingsListing(os.Stdout, vals, args, descs)
}

func listCredentials(store clinote.UserCredentialStore, cmd *cobra.Command) {
	includeToken, err := cmd.Flags().GetBool("show-secret")
	if err != nil {
		fmt.Printf("Error when parsing arguments: %s\n", err.Error())
		return
	}
	list, err := clinote.GetAllCredentials(store)
	if err != nil {
		fmt.Println("Failed to get all credentials:", err)
		return
	}
	if includeToken {
		clinote.WriteCredentialListingWithSecret(os.Stdout, list)
		return
	}
	clinote.WriteCredentialListing(os.Stdout, list)
}

func rmCredential(store clinote.UserCredentialStore, args []string) {
	list := make([]int, 0, len(args))
	for _, arg := range args {
		index, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("%s is not a number, skipping.\n", arg)
			continue
		}
		list = append(list, index-1)
	}
	sort.Ints(list)
	// Track how many that's been removed so we can handle multi-delete.
	removed := 0
	for _, index := range list {
		if err := clinote.RemoveCredentialByIndex(store, index-removed); err != nil {
			fmt.Printf("Error when removing entry %d: %s\n", index+1, err.Error())
		}
		removed++
	}
}

func addCredential(store clinote.UserCredentialStore, cmd *cobra.Command, args []string) {
	name := parseStringFlag(cmd, "name", "Error when parsing the name:", "Please enter a name: ")
	secret := parseStringFlag(cmd, "secret", "Error when parsing the secret:", "Please enter the access token: ")
	sandbox, err := cmd.Flags().GetBool("sandbox")
	if err != nil {
		fmt.Println("Error when parsing the command flag:", err)
		return
	}
	credType := clinote.EvernoteCredential
	if sandbox {
		credType = clinote.EvernoteSandboxCredential
	}
	err = clinote.AddNewCredential(store, name, secret, credType)
	if err != nil {
		fmt.Println("Error when adding the new credentials:", err)
	}
}

func parseStringFlag(cmd *cobra.Command, flag, parseErr, scanLine string) string {
	var name string
	n, err := cmd.Flags().GetString(flag)
	if err != nil {
		fmt.Println(parseErr, err)
	}
	if n == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(scanLine)
		scanner.Scan()
		name = scanner.Text()
	} else {
		name = n
	}
	return name
}
