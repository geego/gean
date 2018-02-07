package collections

import (
	"yiqilai.tech/gean/app/deps"
	"yiqilai.tech/gean/app/tpl/internal"
)

const name = "collections"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.After,
			[]string{"after"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Apply,
			[]string{"apply"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Delimit,
			[]string{"delimit"},
			[][2]string{
				{`{{ delimit (slice "A" "B" "C") ", " " and " }}`, `A, B and C`},
			},
		)

		ns.AddMethodMapping(ctx.Dictionary,
			[]string{"dict"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.EchoParam,
			[]string{"echoParam"},
			[][2]string{
				{`{{ echoParam .Params "langCode" }}`, `en`},
			},
		)

		ns.AddMethodMapping(ctx.First,
			[]string{"first"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.KeyVals,
			[]string{"keyVals"},
			[][2]string{
				{`{{ keyVals "key" "a" "b" }}`, `key: [a b]`},
			},
		)

		ns.AddMethodMapping(ctx.In,
			[]string{"in"},
			[][2]string{
				{`{{ if in "this string contains a substring" "substring" }}Substring found!{{ end }}`, `Substring found!`},
			},
		)

		ns.AddMethodMapping(ctx.Index,
			[]string{"index"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Intersect,
			[]string{"intersect"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.IsSet,
			[]string{"isSet", "isset"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Last,
			[]string{"last"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Querify,
			[]string{"querify"},
			[][2]string{
				{
					`{{ (querify "foo" 1 "bar" 2 "baz" "with spaces" "qux" "this&that=those") | safeHTML }}`,
					`bar=2&baz=with+spaces&foo=1&qux=this%26that%3Dthose`},
				{
					`<a href="https://www.google.com?{{ (querify "q" "test" "page" 3) | safeURL }}">Search</a>`,
					`<a href="https://www.google.com?page=3&amp;q=test">Search</a>`},
			},
		)

		ns.AddMethodMapping(ctx.Shuffle,
			[]string{"shuffle"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Slice,
			[]string{"slice"},
			[][2]string{
				{`{{ slice "B" "C" "A" | sort }}`, `[A B C]`},
			},
		)

		ns.AddMethodMapping(ctx.Sort,
			[]string{"sort"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Union,
			[]string{"union"},
			[][2]string{
				{`{{ union (slice 1 2 3) (slice 3 4 5) }}`, `[1 2 3 4 5]`},
			},
		)

		ns.AddMethodMapping(ctx.Where,
			[]string{"where"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.Seq,
			[]string{"seq"},
			[][2]string{
				{`{{ seq 3 }}`, `[1 2 3]`},
			},
		)
		ns.AddMethodMapping(ctx.Uniq,
			[]string{"uniq"},
			[][2]string{
				{`{{ slice 1 2 3 2 | uniq }}`, `[1 2 3]`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
