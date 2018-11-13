package cast

import (
	"github.com/geego/gean/app/deps"
	"github.com/geego/gean/app/docshelper"
	"github.com/geego/gean/app/tpl/internal"
	"github.com/govenue/configurator"
)

// This file provides documentation support and is randomly put into this package.
func init() {
	docsProvider := func() map[string]interface{} {
		docs := make(map[string]interface{})
		d := &deps.Deps{Cfg: configurator.New()}

		var namespaces internal.TemplateFuncsNamespaces

		for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
			nf := nsf(d)
			namespaces = append(namespaces, nf)

		}

		docs["funcs"] = namespaces
		return docs
	}

	docshelper.AddDocProvider("tpl", docsProvider)
}
