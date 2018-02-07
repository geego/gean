package fmt

import (
	"yiqilai.tech/gean/app/deps"
	"yiqilai.tech/gean/app/tpl/internal"
)

const name = "fmt"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Print,
			[]string{"print"},
			[][2]string{
				{`{{ print "works!" }}`, `works!`},
			},
		)

		ns.AddMethodMapping(ctx.Println,
			[]string{"println"},
			[][2]string{
				{`{{ println "works!" }}`, "works!\n"},
			},
		)

		ns.AddMethodMapping(ctx.Printf,
			[]string{"printf"},
			[][2]string{
				{`{{ printf "%s!" "works" }}`, `works!`},
			},
		)

		ns.AddMethodMapping(ctx.Errorf,
			[]string{"errorf"},
			[][2]string{
				{`{{ errorf "%s." "failed" }}`, `failed.`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
