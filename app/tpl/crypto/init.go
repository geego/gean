package crypto

import (
	"yiqilai.tech/gean/app/deps"
	"yiqilai.tech/gean/app/tpl/internal"
)

const name = "crypto"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New()

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.MD5,
			[]string{"md5"},
			[][2]string{
				{`{{ md5 "Hello world, gophers!" }}`, `b3029f756f98f79e7f1b7f1d1f0dd53b`},
				{`{{ crypto.MD5 "Hello world, gophers!" }}`, `b3029f756f98f79e7f1b7f1d1f0dd53b`},
			},
		)

		ns.AddMethodMapping(ctx.SHA1,
			[]string{"sha1"},
			[][2]string{
				{`{{ sha1 "Hello world, gophers!" }}`, `c8b5b0e33d408246e30f53e32b8f7627a7a649d4`},
			},
		)

		ns.AddMethodMapping(ctx.SHA256,
			[]string{"sha256"},
			[][2]string{
				{`{{ sha256 "Hello world, gophers!" }}`, `6ec43b78da9669f50e4e422575c54bf87536954ccd58280219c393f2ce352b46`},
			},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
