/*
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
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

package clinote

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)
	store := new(mockStore)
	cfg := new(DefaultConfig)
	ns := new(mockNS)

	t.Run("default client", func(t *testing.T) {
		c := NewClient(cfg, store, ns, DefaultClientOptions)
		assert.NotNil(c, "None nil client")
		expectedFunc := reflect.ValueOf(newFileCacheFile).Pointer()
		actualFunc := reflect.ValueOf(c.newCacheFile).Pointer()
		assert.Equal(expectedFunc, actualFunc, "Wrong cache creating function")
		assert.IsType(new(EnvEditor), c.Editor, "Wrong default editor")
	})

	t.Run("client using Vim", func(t *testing.T) {
		c := NewClient(cfg, store, ns, DefaultClientOptions|VimEditer)
		assert.NotNil(c, "None nil client")
		assert.IsType(new(VimEditor), c.Editor, "Wrong editer type")
	})
}
