package cast

import (
	"github.com/geego/gean/app/deps"
	"github.com/geego/gean/app/tpl/internal"
)

const name = "cast"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.ToInt,
			[]string{"int"},
			[][2]string{
				{`{{ "1234" | int | printf "%T" }}`, `int`},
			},
		)

		ns.AddMethodMapping(ctx.ToString,
			[]string{"string"},
			[][2]string{
				{`{{ 1234 | string | printf "%T" }}`, `string`},
			},
		)

		ns.AddMethodMapping(ctx.ToFloat,
			[]string{"float"},
			[][2]string{
				{`{{ "1234" | float | printf "%T" }}`, `float64`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
