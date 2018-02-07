package compare

import (
	"yiqilai.tech/gean/app/deps"
	"yiqilai.tech/gean/app/tpl/internal"
)

const name = "compare"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Default,
			[]string{"default"},
			[][2]string{
				{`{{ "Hugo Rocks!" | default "Hugo Rules!" }}`, `Hugo Rocks!`},
				{`{{ "" | default "Hugo Rules!" }}`, `Hugo Rules!`},
			},
		)

		ns.AddMethodMapping(ctx.Eq,
			[]string{"eq"},
			[][2]string{
				{`{{ if eq .Section "blog" }}current{{ end }}`, `current`},
			},
		)

		ns.AddMethodMapping(ctx.Ge,
			[]string{"ge"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Gt,
			[]string{"gt"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Le,
			[]string{"le"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Lt,
			[]string{"lt"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Ne,
			[]string{"ne"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Conditional,
			[]string{"cond"},
			[][2]string{
				{`{{ cond (eq (add 2 2) 4) "2+2 is 4" "what?" | safeHTML }}`, `2+2 is 4`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
