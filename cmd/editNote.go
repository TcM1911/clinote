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

	"github.com/spf13/cobra"
	"github.com/tcm1911/clinote/config"
	"github.com/tcm1911/clinote/evernote"
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
			verify, err := cmd.Flags().GetBool("verify-content-length")
			if err != nil {
				fmt.Println("Error when parsing verify content length flag:", err)
				return
			}
			editNote(args[0], verify)
		}
	},
}

func init() {
	noteCmd.AddCommand(editNoteCmd)
	editNoteCmd.Flags().StringP("title", "t", "", "Change the note title to.")
	editNoteCmd.Flags().StringP("notebook", "b", "", "Move the note to notebook.")
	editNoteCmd.Flags().Bool("verify-content-length", false, "Verifies that the content length on the server matches content length saved.")
}

func editNote(title string, verifyLength bool) {
	n := evernote.GetNoteWithContent(title)
	n.MDHash = md5.Sum([]byte(n.Title + "\n\n" + n.MD))
	filename := string(n.GUID) + ".md"
	b, err := createTmpFileAndEdit(filename, n.Title, n.MD)
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
	err = parseFileChange(b, n)
	if err != nil {
		fmt.Println("Error parsing the changes:", err)
		return
	}
	err = evernote.SaveChanges(n)
	if err != nil {
		fmt.Println("Error when saving file:", err)
		saveRecoveredNote(n)
	}
	if verifyLength {
		verifyNoteContent(n)
	}
}

func verifyNoteContent(n *evernote.Note) {
	fmt.Println("Verifying that note content was saved.")
	saved, err := evernote.GetNoteByGUID(n.GUID)
	if err != nil {
		fmt.Println("Error when getting saved note from server:", err)
		fmt.Println("Saving note for just in case...")
		saveRecoveredNote(n)
		return
	}
	if len(n.Body) != len(saved.Body) {
		fmt.Println("Length of the note body on the server doesn't match the body length submitted. Saving recovered note...")
		saveRecoveredNote(n)
		return
	}
	fmt.Println("Content length okay.")
}

func saveRecoveredNote(n *evernote.Note) {
	tmpDir := config.GetCacheFolder()
	fp := filepath.Join(tmpDir, "recovered.md")
	f := createTempFile(fp)
	if f == nil {
		fmt.Println("Error when creating recovered file")
		return
	}
	defer f.Close()
	f.WriteString(n.Title + "\n\n" + n.MD)
	if err := f.Sync(); err != nil {
		fmt.Println("Error when syncing content to disk:", err)
		return
	}
}

func createTmpFileAndEdit(filename, title, content string) ([]byte, error) {
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
		panic("Error creating temp file: " + err.Error())
	}
	f.Close()
	f, err = os.OpenFile(fp, os.O_RDWR, 0600)
	if err != nil {
		panic("Error when opening temp file: " + err.Error())
	}
	return f
}

func parseFileChange(b []byte, n *evernote.Note) error {
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
	n.MD = string(bodyBytes)
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
