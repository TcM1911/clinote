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

package markdown

import (
	"errors"

	"github.com/russross/blackfriday"
)

// ErrEmptyConvertedBody is returned if the decoder returns an empty body when content was passed in
// to decode. This means the decoding failed due to a issue in the content.
var ErrEmptyConvertedBody = errors.New("markdown decoding failed and returned an empty body")

// ToXML converts the markdown body to Evernote's xml body style.
func ToXML(mdBody string) ([]byte, error) {
	b := blackfriday.MarkdownCommon([]byte(mdBody))
	if len(b) == 0 && mdBody != "" {
		return b, ErrEmptyConvertedBody
	}
	return b, nil
}
