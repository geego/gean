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

	"yiqilai.tech/gean/app/tpl/internal"

	// Init the namespaces
	_ "yiqilai.tech/gean/app/tpl/cast"
	_ "yiqilai.tech/gean/app/tpl/collections"
	_ "yiqilai.tech/gean/app/tpl/compare"
	_ "yiqilai.tech/gean/app/tpl/crypto"
	_ "yiqilai.tech/gean/app/tpl/data"
	_ "yiqilai.tech/gean/app/tpl/encoding"
	_ "yiqilai.tech/gean/app/tpl/fmt"
	_ "yiqilai.tech/gean/app/tpl/images"
	_ "yiqilai.tech/gean/app/tpl/inflect"
	_ "yiqilai.tech/gean/app/tpl/lang"
	_ "yiqilai.tech/gean/app/tpl/math"
	_ "yiqilai.tech/gean/app/tpl/os"
	_ "yiqilai.tech/gean/app/tpl/partials"
	_ "yiqilai.tech/gean/app/tpl/safe"
	_ "yiqilai.tech/gean/app/tpl/strings"
	_ "yiqilai.tech/gean/app/tpl/time"
	_ "yiqilai.tech/gean/app/tpl/transform"
	_ "yiqilai.tech/gean/app/tpl/urls"
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
