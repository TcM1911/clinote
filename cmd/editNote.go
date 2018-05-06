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
	"bufio"
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/TcM1911/clinote/config"
	"github.com/TcM1911/clinote/evernote"
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
		if title != "" {
			evernote.ChangeTitle(args[0], title)
		}
		if notebook != "" {
			evernote.MoveNote(args[0], notebook)
		}

		if title == "" && notebook == "" {
			editNote(args[0], raw)
		}
	},
}

func init() {
	noteCmd.AddCommand(editNoteCmd)
	editNoteCmd.Flags().StringP("title", "t", "", "Change the note title to.")
	editNoteCmd.Flags().StringP("notebook", "b", "", "Move the note to notebook.")
	editNoteCmd.Flags().Bool("raw", false, "Use raw content instead of markdown version.")
}

func editNote(title string, raw bool) {
	n := evernote.GetNoteWithContent(title)
	n.MDHash = md5.Sum([]byte(n.Title + "\n\n" + n.MD))
	if raw {
		n.MDHash = md5.Sum([]byte(n.Title + "\n\n" + n.Body))
	}
	filename := string(n.GUID) + ".md"
	var body string
	if raw {
		body = n.Body
	} else {
		body = n.MD
	}
	b, err := createTmpFileAndEdit(filename, n.Title, body)
	if err != nil {
		fmt.Println("Error when processing note:", err)
		return
	}
	hash := md5.Sum(b)
	if hash == n.MDHash {
		fmt.Println("No changes detected.")
		return
	}
	fmt.Println("Changes detected, saving note...")
	err = parseFileChange(b, n, raw)
	if err != nil {
		fmt.Println("Error parsing the changes:", err)
		return
	}
	evernote.SaveChanges(n, raw)
}

func createTmpFileAndEdit(filename, title, content string) ([]byte, error) {
	//tempDir := os.TempDir()
	tempDir := config.GetCacheFolder()
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

func parseFileChange(b []byte, n *evernote.Note, raw bool) error {
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
