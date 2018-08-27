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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

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
		if len(args) != 1 {
			fmt.Println("Error, a note has to be given.")
			return
		}
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
		client := defaultClient()
		defer client.Close()
		ns, err := client.GetNoteStore()
		if err != nil {
			return
		}
		if title != "" {
			clinote.ChangeTitle(client.Config.Store(), ns, args[0], title)
		}
		if notebook != "" {
			clinote.MoveNote(client.Config.Store(), ns, args[0], notebook)
		}

		if title == "" && notebook == "" {
			opts := clinote.DefaultNoteOption
			if raw {
				opts = opts | clinote.RawNote
			}
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
}

// func editNote(client *evernote.Client, title string, raw bool) error {
// 	ns, err := client.GetNoteStore()
// 	if err != nil {
// 		return err
// 	}
// 	n, err := clinote.GetNoteWithContent(client.Config.Store(), ns, title)
// 	if err != nil {
// 		return err
// 	}
// 	n.MDHash = md5.Sum([]byte(n.Title + "\n\n" + n.MD))
// 	if raw {
// 		n.MDHash = md5.Sum([]byte(n.Title + "\n\n" + n.Body))
// 	}
// 	filename := string(n.GUID) + ".md"
// 	var body string
// 	if raw {
// 		body = n.Body
// 	} else {
// 		body = n.MD
// 	}
// 	b, err := createTmpFileAndEdit(filename, n.Title, body)
// 	if err != nil {
// 		return err
// 	}
// 	hash := md5.Sum(b)
// 	if hash == n.MDHash {
// 		return err
// 	}
// 	fmt.Println("Changes detected, saving note...")
// 	err = parseFileChange(b, n, raw)
// 	if err != nil {
// 		return err
// 	}
// 	return clinote.SaveChanges(ns, n, raw)
// }

func createTmpFileAndEdit(filename, title, content string) ([]byte, error) {
	cfg := new(clinote.DefaultConfig)
	tempDir := cfg.GetCacheFolder()
	if tempDir == "" {
		return nil, errors.New("no valid temp folder")
	}
	fp := filepath.Join(tempDir, filename)
	defer os.Remove(fp)
	f := createTempFile(fp)
	defer f.Close()
	if f == nil {
		return nil, errors.New("error when creating temp file")
	}
	f.WriteString(title + "\n\n" + content)
	err := f.Sync()
	if err != nil {
		return nil, errors.New("error when flushing data to temp file" + err.Error())
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	editInEditor(fp)
	f, err = os.Open(fp)
	if err != nil {
		return nil, err
	}
	// Go back to the beginning of the file.
	f.Seek(int64(0), 0)
	return ioutil.ReadAll(f)
}

func createTempFile(fp string) *os.File {
	f, err := os.OpenFile(fp, os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		f.Close()
		return nil
	}
	f.Close()
	f, err = os.OpenFile(fp, os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("Error when opening temp file:", err)
		return nil
	}
	return f
}

func parseFileChange(b []byte, n *clinote.Note, raw bool) error {
	r := bytes.NewReader(b)
	br := bufio.NewReader(r)
	// First line is the note title.
	line, _, err := br.ReadLine()
	if err != nil {
		return errors.New("error parsing title" + err.Error())
	}
	title := string(line)
	bodyBytes, err := ioutil.ReadAll(br)
	if err != nil {
		return errors.New("error parsing note body " + err.Error())
	}
	n.Title = title
	if raw {
		n.Body = string(bodyBytes)
	} else {
		n.MD = string(bodyBytes)
	}
	return nil
}

func editInEditor(file string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
