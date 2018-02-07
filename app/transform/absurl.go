package transform

var ar = newAbsURLReplacer()

// AbsURL replaces relative URLs with absolute ones
// in HTML files, using the baseURL setting.
var AbsURL = func(ct contentTransformer) {
	ar.replaceInHTML(ct)
}

// AbsURLInXML replaces relative URLs with absolute ones
// in XML files, using the baseURL setting.
var AbsURLInXML = func(ct contentTransformer) {
	ar.replaceInXML(ct)
}
