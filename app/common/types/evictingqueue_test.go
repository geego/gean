package types

import (
	"sync"
	"testing"

	"github.com/gostores/require"
)

func TestEvictingStringQueue(t *testing.T) {
	assert := require.New(t)

	queue := NewEvictingStringQueue(3)

	assert.Equal("", queue.Peek())
	queue.Add("a")
	queue.Add("b")
	queue.Add("a")
	assert.Equal("b", queue.Peek())
	queue.Add("b")
	assert.Equal("b", queue.Peek())

	queue.Add("a")
	queue.Add("b")

	assert.Equal([]string{"b", "a"}, queue.PeekAll())
	assert.Equal("b", queue.Peek())
	queue.Add("c")
	queue.Add("d")
	// Overflowed, a should now be removed.
	assert.Equal([]string{"d", "c", "b"}, queue.PeekAll())
	assert.Len(queue.PeekAllSet(), 3)
	assert.True(queue.PeekAllSet()["c"])
}

func TestEvictingStringQueueConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	val := "someval"

	queue := NewEvictingStringQueue(3)

	for j := 0; j < 100; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			queue.Add(val)
			v := queue.Peek()
			if v != val {
				t.Error("wrong val")
			}
			vals := queue.PeekAll()
			if len(vals) != 1 || vals[0] != val {
				t.Error("wrong val")
			}
		}()
	}
	wg.Wait()
}
