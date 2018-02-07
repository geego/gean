package helpers

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/gostores/configurator"
	"github.com/gostores/require"
)

// Renders a codeblock using markdown
func (c ContentSpec) render(input string) string {
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewMarkdown()}
	render := c.getHTMLRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html")
	return buf.String()
}

// Renders a codeblock using Mmark
func (c ContentSpec) renderWithMmark(input string) string {
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewMarkdown()}
	render := c.getMmarkHTMLRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html", []byte(""), false, false)
	return buf.String()
}

func TestCodeFence(t *testing.T) {
	assert := require.New(t)

	type test struct {
		enabled         bool
		input, expected string
	}

	// Pygments 2.0 and 2.1 have slightly different outputs so only do partial matching
	data := []test{
		{true, "<html></html>", `(?s)^<div class="highlight">\n?<pre.*><code class="language-html" data-lang="html">.*?</code></pre>\n?</div>\n?$`},
		{false, "<html></html>", `(?s)^<pre.*><code class="language-html">.*?</code></pre>\n$`},
	}

	for _, useClassic := range []bool{false, true} {
		for i, d := range data {
			v := configurator.New()
			v.Set("pygmentsStyle", "monokai")
			v.Set("pygmentsUseClasses", true)
			v.Set("pygmentsCodeFences", d.enabled)
			v.Set("pygmentsUseClassic", useClassic)

			c, err := NewContentSpec(v)
			assert.NoError(err)

			result := c.render(d.input)

			expectedRe, err := regexp.Compile(d.expected)

			if err != nil {
				t.Fatal("Invalid regexp", err)
			}
			matched := expectedRe.MatchString(result)

			if !matched {
				t.Errorf("Test %d failed. Markdown enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
			}

			result = c.renderWithMmark(d.input)
			matched = expectedRe.MatchString(result)
			if !matched {
				t.Errorf("Test %d failed. Mmark enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
			}
		}
	}
}

func TestMarkdownTaskList(t *testing.T) {
	c := newTestContentSpec()

	for i, this := range []struct {
		markdown        string
		taskListEnabled bool
		expect          string
	}{
		{`
TODO:

- [x] On1
- [X] On2
- [ ] Off

END
`, true, `<p>TODO:</p>

<ul class="task-list">
<li><label><input type="checkbox" checked disabled class="task-list-item"> On1</label></li>
<li><label><input type="checkbox" checked disabled class="task-list-item"> On2</label></li>
<li><label><input type="checkbox" disabled class="task-list-item"> Off</label></li>
</ul>

<p>END</p>
`},
		{`- [x] On1`, false, `<ul>
<li>[x] On1</li>
</ul>
`},
		{`* [ ] Off

END`, true, `<ul class="task-list">
<li><label><input type="checkbox" disabled class="task-list-item"> Off</label></li>
</ul>

<p>END</p>
`},
	} {
		markdownConfig := c.NewMarkdown()
		markdownConfig.TaskLists = this.taskListEnabled
		ctx := &RenderingContext{Content: []byte(this.markdown), PageFmt: "markdown", Config: markdownConfig}

		result := string(c.RenderBytes(ctx))

		if result != this.expect {
			t.Errorf("[%d] got \n%v but expected \n%v", i, result, this.expect)
		}
	}
}
