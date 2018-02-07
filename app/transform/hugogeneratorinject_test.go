package transform

import (
	"bytes"
	"strings"
	"testing"
)

func TestHugoGeneratorInject(t *testing.T) {
	hugoGeneratorTag = "META"
	for i, this := range []struct {
		in     string
		expect string
	}{
		{`<head>
	<foo />
</head>`, `<head>
	META
	<foo />
</head>`},
		{`<HEAD>
	<foo />
</HEAD>`, `<HEAD>
	META
	<foo />
</HEAD>`},
		{`<head><meta name="generator" content="Jekyll" /></head>`, `<head><meta name="generator" content="Jekyll" /></head>`},
		{`<head><meta name='generator' content='Jekyll' /></head>`, `<head><meta name='generator' content='Jekyll' /></head>`},
		{`<head><meta name=generator content=Jekyll /></head>`, `<head><meta name=generator content=Jekyll /></head>`},
		{`<head><META     NAME="GENERATOR" content="Jekyll" /></head>`, `<head><META     NAME="GENERATOR" content="Jekyll" /></head>`},
		{"", ""},
		{"</head>", "</head>"},
		{"<head>", "<head>\n\tMETA"},
	} {
		in := strings.NewReader(this.in)
		out := new(bytes.Buffer)

		tr := NewChain(HugoGeneratorInject)
		tr.Apply(out, in, []byte(""))

		if out.String() != this.expect {
			t.Errorf("[%d] Expected \n%q got \n%q", i, this.expect, out.String())
		}
	}

}
