package collections

import (
	"fmt"
	"testing"

	"github.com/geego/gean/app/deps"
	"github.com/govenue/assert"
	"github.com/govenue/require"
)

func TestIndex(t *testing.T) {
	t.Parallel()

	ns := New(&deps.Deps{})

	for i, test := range []struct {
		item    interface{}
		indices []interface{}
		expect  interface{}
		isErr   bool
	}{
		{[]int{0, 1}, []interface{}{0}, 0, false},
		{[]int{0, 1}, []interface{}{9}, nil, false}, // index out of range
		{[]uint{0, 1}, nil, []uint{0, 1}, false},
		{[][]int{{1, 2}, {3, 4}}, []interface{}{0, 0}, 1, false},
		{map[int]int{1: 10, 2: 20}, []interface{}{1}, 10, false},
		{map[int]int{1: 10, 2: 20}, []interface{}{0}, 0, false},
		// errors
		{nil, nil, nil, true},
		{[]int{0, 1}, []interface{}{"1"}, nil, true},
		{[]int{0, 1}, []interface{}{nil}, nil, true},
		{tstNoStringer{}, []interface{}{0}, nil, true},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test)

		result, err := ns.Index(test.item, test.indices...)

		if test.isErr {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
