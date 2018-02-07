package bufferpool

import (
	"testing"

	"github.com/gostores/assert"
)

func TestBufferPool(t *testing.T) {
	buff := GetBuffer()
	buff.WriteString("do be do be do")
	assert.Equal(t, "do be do be do", buff.String())
	PutBuffer(buff)
	assert.Equal(t, 0, buff.Len())
}
