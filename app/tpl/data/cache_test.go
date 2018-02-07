package data

import (
	"fmt"
	"testing"

	"github.com/gostores/assert"
	"github.com/gostores/configurator"
	"github.com/gostores/fsintra"
)

func TestCache(t *testing.T) {
	t.Parallel()

	fs := new(fsintra.MemMapFs)

	for i, test := range []struct {
		path    string
		content []byte
		ignore  bool
	}{
		{"http://Foo.Bar/foo_Bar-Foo", []byte(`T€st Content 123`), false},
		{"fOO,bar:foo%bAR", []byte(`T€st Content 123 fOO,bar:foo%bAR`), false},
		{"FOo/BaR.html", []byte(`FOo/BaR.html T€st Content 123`), false},
		{"трям/трям", []byte(`T€st трям/трям Content 123`), false},
		{"은행", []byte(`T€st C은행ontent 123`), false},
		{"Банковский кассир", []byte(`Банковский кассир T€st Content 123`), false},
		{"Банковский кассир", []byte(`Банковский кассир T€st Content 456`), true},
	} {
		msg := fmt.Sprintf("Test #%d: %v", i, test)

		cfg := configurator.New()

		c, err := getCache(test.path, fs, cfg, test.ignore)
		assert.NoError(t, err, msg)
		assert.Nil(t, c, msg)

		err = writeCache(test.path, test.content, fs, cfg, test.ignore)
		assert.NoError(t, err, msg)

		c, err = getCache(test.path, fs, cfg, test.ignore)
		assert.NoError(t, err, msg)

		if test.ignore {
			assert.Nil(t, c, msg)
		} else {
			assert.Equal(t, string(test.content), string(c))
		}
	}
}
