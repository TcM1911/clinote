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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	assert := assert.New(t)
	t.Run("add use raw content", func(t *testing.T) {
		ctx := AddUseRawContentToContext(context.Background(), true)
		assert.True(ctx.Value(rawContentContextKey).(bool), "Should return true")
	})
	t.Run("false if value not set", func(t *testing.T) {
		assert.False(GetUseRawContentFromContext(context.Background()), "Should return false")
	})
	t.Run("return set value", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), rawContentContextKey, true)
		val := GetUseRawContentFromContext(ctx)
		assert.True(val, "Should return set true value")
	})
}
