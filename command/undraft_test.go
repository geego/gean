package command

// TODO Support Mac Encoding (\r)

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"yiqilai.tech/gean/app/parser"
)

var (
	jsonFM         = "{\n \"date\": \"12-04-06\",\n \"title\": \"test json\"\n}"
	jsonDraftFM    = "{\n \"draft\": true,\n \"date\": \"12-04-06\",\n \"title\":\"test json\"\n}"
	tomlFM         = "+++\n date= \"12-04-06\"\n title= \"test toml\"\n+++"
	tomlDraftFM    = "+++\n draft= true\n date= \"12-04-06\"\n title=\"test toml\"\n+++"
	yamlFM         = "---\n date: \"12-04-06\"\n title: \"test yaml\"\n---"
	yamlDraftFM    = "---\n draft: true\n date: \"12-04-06\"\n title: \"test yaml\"\n---"
	yamlYesDraftFM = "---\n draft: yes\n date: \"12-04-06\"\n title: \"test yaml\"\n---"
)

func TestUndraftContent(t *testing.T) {
	tests := []struct {
		fm          string
		expectedErr string
	}{
		{jsonFM, "not a Draft: nothing was done"},
		{jsonDraftFM, ""},
		{tomlFM, "not a Draft: nothing was done"},
		{tomlDraftFM, ""},
		{yamlFM, "not a Draft: nothing was done"},
		{yamlDraftFM, ""},
		{yamlYesDraftFM, ""},
	}

	for i, test := range tests {
		r := bytes.NewReader([]byte(test.fm))
		p, _ := parser.ReadFrom(r)
		res, err := undraftContent(p)
		if test.expectedErr != "" {
			if err == nil {
				t.Errorf("[%d] Expected error, got none", i)
				continue
			}
			if err.Error() != test.expectedErr {
				t.Errorf("[%d] Expected %q, got %q", i, test.expectedErr, err)
				continue
			}
		} else {
			r = bytes.NewReader(res.Bytes())
			p, _ = parser.ReadFrom(r)
			meta, err := p.Metadata()
			if err != nil {
				t.Errorf("[%d] unexpected error %q", i, err)
				continue
			}
			for k, v := range meta.(map[string]interface{}) {
				if k == "draft" {
					if v.(bool) {
						t.Errorf("[%d] Expected %q to be \"false\", got \"true\"", i, k)
						continue
					}
				}
				if k == "date" {
					if !strings.HasPrefix(v.(string), time.Now().Format("2006-01-02")) {
						t.Errorf("[%d] Expected %v to start with %v", i, v.(string), time.Now().Format("2006-01-02"))
					}
				}
			}
		}
	}
}
