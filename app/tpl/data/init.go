package data

import (
	"yiqilai.tech/gean/app/deps"
	"yiqilai.tech/gean/app/tpl/internal"
)

const name = "data"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.GetCSV,
			[]string{"getCSV"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.GetJSON,
			[]string{"getJSON"},
			[][2]string{},
		)
		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
