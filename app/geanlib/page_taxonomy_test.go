package geanlib

import (
	"reflect"
	"strings"
	"testing"
)

var pageYamlWithTaxonomiesA = `---
tags: ['a', 'B', 'c']
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var pageYamlWithTaxonomiesB = `---
tags:
 - "a"
 - "B"
 - "c"
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var pageYamlWithTaxonomiesC = `---
tags: 'E'
categories: 'd'
---
YAML frontmatter with tags and categories taxonomy.`

var pageJSONWithTaxonomies = `{
  "categories": "D",
  "tags": [
    "a",
    "b",
    "c"
  ]
}
JSON Front Matter with tags and categories`

var pageTomlWithTaxonomies = `+++
tags = [ "a", "B", "c" ]
categories = "d"
+++
TOML Front Matter with tags and categories`

func TestParseTaxonomies(t *testing.T) {
	t.Parallel()
	for _, test := range []string{pageTomlWithTaxonomies,
		pageJSONWithTaxonomies,
		pageYamlWithTaxonomiesA,
		pageYamlWithTaxonomiesB,
		pageYamlWithTaxonomiesC,
	} {

		s := newTestSite(t)
		p, _ := s.NewPage("page/with/taxonomy")
		_, err := p.ReadFrom(strings.NewReader(test))
		if err != nil {
			t.Fatalf("Failed parsing %q: %s", test, err)
		}

		param := p.GetParam("tags")

		if params, ok := param.([]string); ok {
			expected := []string{"a", "b", "c"}
			if !reflect.DeepEqual(params, expected) {
				t.Errorf("Expected %s: got: %s", expected, params)
			}
		} else if params, ok := param.(string); ok {
			expected := "e"
			if params != expected {
				t.Errorf("Expected %s: got: %s", expected, params)
			}
		}

		param = p.GetParam("categories")
		singleparam := param.(string)

		if singleparam != "d" {
			t.Fatalf("Expected: d, got: %s", singleparam)
		}
	}
}
