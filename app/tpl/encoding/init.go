package encoding

import (
	"yiqilai.tech/gean/app/deps"
	"yiqilai.tech/gean/app/tpl/internal"
)

const name = "encoding"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Base64Decode,
			[]string{"base64Decode"},
			[][2]string{
				{`{{ "SGVsbG8gd29ybGQ=" | base64Decode }}`, `Hello world`},
				{`{{ 42 | base64Encode | base64Decode }}`, `42`},
			},
		)

		ns.AddMethodMapping(ctx.Base64Encode,
			[]string{"base64Encode"},
			[][2]string{
				{`{{ "Hello world" | base64Encode }}`, `SGVsbG8gd29ybGQ=`},
			},
		)

		ns.AddMethodMapping(ctx.Jsonify,
			[]string{"jsonify"},
			[][2]string{
				{`{{ (slice "A" "B" "C") | jsonify }}`, `["A","B","C"]`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
