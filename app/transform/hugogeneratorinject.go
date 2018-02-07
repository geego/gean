package transform

import (
	"bytes"
	"fmt"
	"regexp"

	"yiqilai.tech/gean/app/helpers"
)

var metaTagsCheck = regexp.MustCompile(`(?i)<meta\s+name=['|"]?generator['|"]?`)
var hugoGeneratorTag = fmt.Sprintf(`<meta name="generator" content="Hugo %s" />`, helpers.CurrentHugoVersion)

// HugoGeneratorInject injects a meta generator tag for Hugo if none present.
func HugoGeneratorInject(ct contentTransformer) {
	if metaTagsCheck.Match(ct.Content()) {
		if _, err := ct.Write(ct.Content()); err != nil {
			helpers.DistinctWarnLog.Println("Failed to inject Hugo generator tag:", err)
		}
		return
	}

	head := "<head>"
	replace := []byte(fmt.Sprintf("%s\n\t%s", head, hugoGeneratorTag))
	newcontent := bytes.Replace(ct.Content(), []byte(head), replace, 1)

	if len(newcontent) == len(ct.Content()) {
		head := "<HEAD>"
		replace := []byte(fmt.Sprintf("%s\n\t%s", head, hugoGeneratorTag))
		newcontent = bytes.Replace(ct.Content(), []byte(head), replace, 1)
	}

	if _, err := ct.Write(newcontent); err != nil {
		helpers.DistinctWarnLog.Println("Failed to inject Hugo generator tag:", err)
	}

}
