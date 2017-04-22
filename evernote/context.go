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
 * Copyright (C) Joakim Kennedy, 2017
 */

package evernote

import (
	"context"
)

const (
	rawContentContextKey contextKey = iota
)

type contextKey int8

// AddUseRawContentToContext adds the value to the context.
func AddUseRawContentToContext(ctx context.Context, val bool) context.Context {
	return context.WithValue(ctx, rawContentContextKey, val)
}

// GetUseRawContentFromContext get's the useRawContent from the context.
func GetUseRawContentFromContext(ctx context.Context) bool {
	val, ok := ctx.Value(rawContentContextKey).(bool)
	if !ok {
		return false
	}
	return val
}
