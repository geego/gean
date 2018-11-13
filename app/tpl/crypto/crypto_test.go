package crypto

import (
	"fmt"
	"testing"

	"github.com/govenue/assert"
	"github.com/govenue/require"
)

func TestMD5(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		in     interface{}
		expect interface{}
	}{
		{"Hello world, gophers!", "b3029f756f98f79e7f1b7f1d1f0dd53b"},
		{"Lorem ipsum dolor", "06ce65ac476fc656bea3fca5d02cfd81"},
		{t, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.in)

		result, err := ns.MD5(test.in)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestSHA1(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		in     interface{}
		expect interface{}
	}{
		{"Hello world, gophers!", "c8b5b0e33d408246e30f53e32b8f7627a7a649d4"},
		{"Lorem ipsum dolor", "45f75b844be4d17b3394c6701768daf39419c99b"},
		{t, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.in)

		result, err := ns.SHA1(test.in)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}

func TestSHA256(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		in     interface{}
		expect interface{}
	}{
		{"Hello world, gophers!", "6ec43b78da9669f50e4e422575c54bf87536954ccd58280219c393f2ce352b46"},
		{"Lorem ipsum dolor", "9b3e1beb7053e0f900a674dd1c99aca3355e1275e1b03d3cb1bc977f5154e196"},
		{t, false},
	} {
		errMsg := fmt.Sprintf("[%d] %v", i, test.in)

		result, err := ns.SHA256(test.in)

		if b, ok := test.expect.(bool); ok && !b {
			require.Error(t, err, errMsg)
			continue
		}

		require.NoError(t, err, errMsg)
		assert.Equal(t, test.expect, result, errMsg)
	}
}
