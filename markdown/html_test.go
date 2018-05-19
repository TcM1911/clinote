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
 * Copyright (C) Joakim Kennedy, 2018
 */

package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromHTML(t *testing.T) {
	assert := assert.New(t)

	para := "Test paragraph"
	doc := "<p>" + para + "</p>"
	expected := para + "\n\n\n"

	actual, err := FromHTML(doc)
	assert.NoError(err, "Should parse the doc without an error")
	assert.Equal(expected, actual, "Not converted")
}
