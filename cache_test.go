package clinote

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotebookListExpiration(t *testing.T) {
	assert := assert.New(t)
	t.Run("should_return_false_if_not_expired", func(t *testing.T) {
		list := &NotebookCacheList{Limit: DefaultNotebookCacheTime, Timestamp: time.Now()}
		assert.False(list.IsOutdated(), "Should return false if not expired.")
	})
	t.Run("should_return_true_if_expired", func(t *testing.T) {
		list := &NotebookCacheList{Limit: 1 * time.Nanosecond, Timestamp: time.Now()}
		time.Sleep(10 * time.Microsecond)
		assert.True(list.IsOutdated(), "Should return true if expired.")
	})
}

func TestNewNotebookCacheList(t *testing.T) {
	assert := assert.New(t)
	count := 3
	books := make([]*Notebook, count)
	t.Run("default_expiration", func(t *testing.T) {
		list := NewNotebookCacheList(books)
		assert.Len(list.Notebooks, count, "Incorrect length.")
		assert.Equal(DefaultNotebookCacheTime, list.Limit, "Incorrect limit.")
	})
	t.Run("set_expiration", func(t *testing.T) {
		expectedLimit := 6 * time.Hour
		list := NewNotebookCacheListWithLimit(books, expectedLimit)
		assert.Len(list.Notebooks, count, "Incorrect length.")
		assert.Equal(expectedLimit, list.Limit, "Incorrect limit.")
	})
}
