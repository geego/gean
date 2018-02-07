package transform

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestLiveReloadInject(t *testing.T) {
	doTestLiveReloadInject(t, "</body>")
	doTestLiveReloadInject(t, "</BODY>")
}

func doTestLiveReloadInject(t *testing.T, bodyEndTag string) {
	out := new(bytes.Buffer)
	in := strings.NewReader(bodyEndTag)

	tr := NewChain(LiveReloadInject(1313))
	tr.Apply(out, in, []byte("path"))

	expected := fmt.Sprintf(`<script data-no-instant>document.write('<script src="/livereload.js?port=1313&mindelay=10"></' + 'script>')</script>%s`, bodyEndTag)
	if string(out.Bytes()) != expected {
		t.Errorf("Expected %s got %s", expected, string(out.Bytes()))
	}
}
