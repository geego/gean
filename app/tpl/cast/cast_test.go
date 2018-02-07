package cast

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/gostores/assert"
	"github.com/gostores/require"
)

func TestToInt(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{"1", 1},
		{template.HTML("2"), 2},
		{template.CSS("3"), 3},
		{template.HTMLAttr("4"), 4},
		{template.JS("5"), 5},
		{template.JSStr("6"), 6},
		{"a", false},
		{t, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.v)

		result, err := ns.ToInt(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestToString(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{1, "1"},
		{template.HTML("2"), "2"},
		{"a", "a"},
		{t, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.v)

		result, err := ns.ToString(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestToFloat(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		v      interface{}
		expect interface{}
	}{
		{"1", 1.0},
		{template.HTML("2"), 2.0},
		{template.CSS("3"), 3.0},
		{template.HTMLAttr("4"), 4.0},
		{template.JS("-5.67"), -5.67},
		{template.JSStr("6"), 6.0},
		{"1.23", 1.23},
		{"-1.23", -1.23},
		{"0", 0.0},
		{float64(2.12), 2.12},
		{int64(123), 123.0},
		{2, 2.0},
		{t, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.v)

		result, err := ns.ToFloat(test.v)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
