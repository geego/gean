package command

import (
	"bytes"
	"errors"
	"os"
	"time"

	"github.com/govenue/goman"

	"github.com/geego/gean/app/parser"
)

var undraftCmd = &goman.Command{
	Use:   "undraft path/to/content",
	Short: "Undraft resets the content's draft status",
	Long: `Undraft resets the content's draft status
and updates the date to the current date and time.
If the content's draft status is 'False', nothing is done.`,
	RunE: Undraft,
}

// Undraft publishes the specified content by setting its draft status
// to false and setting its publish date to now. If the specified content is
// not a draft, it will log an error.
func Undraft(cmd *goman.Command, args []string) error {
	cfg, err := InitializeConfig()

	if err != nil {
		return err
	}

	if len(args) < 1 {
		return newUserError("a piece of content needs to be specified")
	}

	location := args[0]
	// open the file
	f, err := cfg.Fs.Source.Open(location)
	if err != nil {
		return err
	}

	// get the page from file
	p, err := parser.ReadFrom(f)
	f.Close()
	if err != nil {
		return err
	}

	w, err := undraftContent(p)
	if err != nil {
		return newSystemErrorF("an error occurred while undrafting %q: %s", location, err)
	}

	f, err = cfg.Fs.Source.OpenFile(location, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return newSystemErrorF("%q not be undrafted due to error opening file to save changes: %q\n", location, err)
	}
	defer f.Close()
	_, err = w.WriteTo(f)
	if err != nil {
		return newSystemErrorF("%q not be undrafted due to save error: %q\n", location, err)
	}
	return nil
}

// undraftContent: if the content is a draft, change its draft status to
// 'false' and set the date to time.Now(). If the draft status is already
// 'false', don't do anything.
func undraftContent(p parser.Page) (bytes.Buffer, error) {
	var buff bytes.Buffer
	// get the metadata; easiest way to see if it's a draft
	meta, err := p.Metadata()
	if err != nil {
		return buff, err
	}
	// since the metadata was obtainable, we can also get the key/value separator for
	// Front Matter
	fm := p.FrontMatter()
	if fm == nil {
		return buff, errors.New("Front Matter was found, nothing was finalized")
	}

	var isDraft, gotDate bool
	var date string
L:
	for k, v := range meta.(map[string]interface{}) {
		switch k {
		case "draft":
			if !v.(bool) {
				return buff, errors.New("not a Draft: nothing was done")
			}
			isDraft = true
			if gotDate {
				break L
			}
		case "date":
			date = v.(string) // capture the value to make replacement easier
			gotDate = true
			if isDraft {
				break L
			}
		}
	}

	// if draft wasn't found in FrontMatter, it isn't a draft.
	if !isDraft {
		return buff, errors.New("not a Draft: nothing was done")
	}

	// get the front matter as bytes and split it into lines
	var lineEnding []byte
	fmLines := bytes.Split(fm, []byte("\n"))
	if len(fmLines) == 1 { // if the result is only 1 element, try to split on dos line endings
		fmLines = bytes.Split(fm, []byte("\r\n"))
		if len(fmLines) == 1 {
			return buff, errors.New("unable to split FrontMatter into lines")
		}
		lineEnding = append(lineEnding, []byte("\r\n")...)
	} else {
		lineEnding = append(lineEnding, []byte("\n")...)
	}

	// Write the front matter lines to the buffer, replacing as necessary
	for _, v := range fmLines {
		pos := bytes.Index(v, []byte("draft"))
		if pos != -1 {
			continue
		}
		pos = bytes.Index(v, []byte("date"))
		if pos != -1 { // if date field wasn't found, add it
			v = bytes.Replace(v, []byte(date), []byte(time.Now().Format(time.RFC3339)), 1)
		}
		buff.Write(v)
		buff.Write(lineEnding)
	}

	// append the actual content
	buff.Write(p.Content())

	return buff, nil
}
