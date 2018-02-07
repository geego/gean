package collections

import (
	"fmt"
	"testing"

	"github.com/gostores/require"
	"yiqilai.tech/gean/app/deps"
	"yiqilai.tech/gean/app/tpl"
)

type templateFinder int

func (templateFinder) Lookup(name string) *tpl.TemplateAdapter {
	return nil
}

func (templateFinder) GetFuncs() map[string]interface{} {
	return map[string]interface{}{
		"print": fmt.Sprint,
	}
}

func TestApply(t *testing.T) {
	t.Parallel()

	ns := New(&deps.Deps{Tmpl: new(templateFinder)})

	strings := []interface{}{"a\n", "b\n"}

	result, err := ns.Apply(strings, "print", "a", "b", "c")
	require.NoError(t, err)
	require.Equal(t, []interface{}{"abc", "abc"}, result, "testing variadic")

	_, err = ns.Apply(strings, "apply", ".")
	require.Error(t, err)

	var nilErr *error
	_, err = ns.Apply(nilErr, "chomp", ".")
	require.Error(t, err)

	_, err = ns.Apply(strings, "dobedobedo", ".")
	require.Error(t, err)

	_, err = ns.Apply(strings, "foo.Chomp", "c\n")
	if err == nil {
		t.Errorf("apply with unknown func should fail")
	}

}
