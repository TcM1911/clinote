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
	"html"

	"github.com/lunny/html2md"
)

func ToHTML(body string) string {
	html2md.AddRule("code", code())
	html2md.AddRule("pre", pre())
	return html2md.Convert(html.UnescapeString(body))
}

func code() *html2md.Rule {
	return &html2md.Rule{
		Patterns: []string{"code"},
		Replacement: func(innerHTML string, attrs []string) string {
			if len(attrs) > 1 {
				return "```\n" + attrs[1] + "```\n"
			}
			return ""
		},
	}
}

func pre() *html2md.Rule {
	return &html2md.Rule{
		Patterns: []string{"pre"},
		Replacement: func(innerHTML string, attrs []string) string {
			if len(attrs) > 1 {
				return "" + attrs[1] + ""
			}
			return ""
		},
	}
}
