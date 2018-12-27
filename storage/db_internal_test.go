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

package storage

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testWaitTime = 1000 * time.Microsecond
)

func TestOpenAndCloseDB(t *testing.T) {
	assert := assert.New(t)
	oldTime := currentWaitTime
	currentWaitTime = testWaitTime
	db, tmpDir := setupTestDB(t)
	// db.waitTime = testWaitTime

	t.Run("multiopen_calls", func(t *testing.T) {
		newDB, err := db.open()
		assert.NoError(err)
		db.handlerMu.Lock()
		assert.Equal(db.bolt, newDB)
		db.handlerMu.Unlock()

		newDB, err = db.open()
		assert.NoError(err)
		db.handlerMu.Lock()
		assert.Equal(db.bolt, newDB)
		db.handlerMu.Unlock()
	})

	t.Run("multiclose_calls", func(t *testing.T) {
		err := db.closeDB()
		assert.NoError(err)

		err = db.closeDB()
		assert.NoError(err)
	})

	t.Run("open_for_gethandler_call", func(t *testing.T) {
		newDB, err := db.getDBHandler()
		assert.NoError(err)
		assert.NotNil(newDB)
		db.handlerMu.Unlock()
	})

	t.Run("close_db_after_wait_time", func(t *testing.T) {
		assert.Equal(testWaitTime, db.waitTime)
		newDB, err := db.getDBHandler()
		assert.NoError(err)
		assert.NotNil(newDB)
		db.handlerMu.Unlock()

		// Sleep to wait for time to expire
		runtime.Gosched()
		time.Sleep(testWaitTime * 2)

		db.handlerMu.Lock()
		assert.Nil(db.bolt)
		db.handlerMu.Unlock()
	})

	fmt.Println("Cleanup")

	// Cleanup
	db.Close()
	os.RemoveAll(tmpDir)
	currentWaitTime = oldTime
}
