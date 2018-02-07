package transform

import (
	"bytes"
	"fmt"
)

// LiveReloadInject returns a function that can be used
// to inject a script tag for the livereload JavaScript in a HTML document.
func LiveReloadInject(port int) func(ct contentTransformer) {
	return func(ct contentTransformer) {
		endBodyTag := "</body>"
		match := []byte(endBodyTag)
		replaceTemplate := `<script data-no-instant>document.write('<script src="/livereload.js?port=%d&mindelay=10"></' + 'script>')</script>%s`
		replace := []byte(fmt.Sprintf(replaceTemplate, port, endBodyTag))

		newcontent := bytes.Replace(ct.Content(), match, replace, 1)
		if len(newcontent) == len(ct.Content()) {
			endBodyTag = "</BODY>"
			replace := []byte(fmt.Sprintf(replaceTemplate, port, endBodyTag))
			match := []byte(endBodyTag)
			newcontent = bytes.Replace(ct.Content(), match, replace, 1)
		}

		ct.Write(newcontent)
	}
}
