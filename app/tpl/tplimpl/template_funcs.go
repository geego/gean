// Copyright 2017-present The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tplimpl

import (
	"html/template"

	"github.com/geego/gean/app/tpl/internal"

	// Init the namespaces
	_ "github.com/geego/gean/app/tpl/cast"
	_ "github.com/geego/gean/app/tpl/collections"
	_ "github.com/geego/gean/app/tpl/compare"
	_ "github.com/geego/gean/app/tpl/crypto"
	_ "github.com/geego/gean/app/tpl/data"
	_ "github.com/geego/gean/app/tpl/encoding"
	_ "github.com/geego/gean/app/tpl/fmt"
	_ "github.com/geego/gean/app/tpl/images"
	_ "github.com/geego/gean/app/tpl/inflect"
	_ "github.com/geego/gean/app/tpl/lang"
	_ "github.com/geego/gean/app/tpl/math"
	_ "github.com/geego/gean/app/tpl/os"
	_ "github.com/geego/gean/app/tpl/partials"
	_ "github.com/geego/gean/app/tpl/safe"
	_ "github.com/geego/gean/app/tpl/strings"
	_ "github.com/geego/gean/app/tpl/time"
	_ "github.com/geego/gean/app/tpl/transform"
	_ "github.com/geego/gean/app/tpl/urls"
)

func (t *templateFuncster) initFuncMap() {
	funcMap := template.FuncMap{}

	// Merge the namespace funcs
	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns := nsf(t.Deps)
		if _, exists := funcMap[ns.Name]; exists {
			panic(ns.Name + " is a duplicate template func")
		}
		funcMap[ns.Name] = ns.Context
		for _, mm := range ns.MethodMappings {
			for _, alias := range mm.Aliases {
				if _, exists := funcMap[alias]; exists {
					panic(alias + " is a duplicate template func")
				}
				funcMap[alias] = mm.Method
			}

		}
	}

	t.funcMap = funcMap
	t.Tmpl.(*templateHandler).setFuncs(funcMap)
}
