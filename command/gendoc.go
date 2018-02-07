package command

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gostores/goman"
	"github.com/gostores/goman/doc"
	"github.com/gostores/notepad"
	"yiqilai.tech/gean/app/geanfs"
	"yiqilai.tech/gean/app/helpers"
)

const gendocFrontmatterTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`

var gendocdir string
var gendocCmd = &goman.Command{
	Use:   "doc",
	Short: "Generate Markdown documentation for the Gean CLI.",
	Long: `Generate Markdown documentation for the Gean CLI.

This command is, mostly, used to create up-to-date documentation
of Gean's command-line interface.

It creates one Markdown file per command with front matter suitable
for rendering in Gean.`,

	RunE: func(cmd *goman.Command, args []string) error {
		if !strings.HasSuffix(gendocdir, helpers.FilePathSeparator) {
			gendocdir += helpers.FilePathSeparator
		}
		if found, _ := helpers.Exists(gendocdir, geanfs.Os); !found {
			notepad.FEEDBACK.Println("Directory", gendocdir, "does not exist, creating...")
			if err := geanfs.Os.MkdirAll(gendocdir, 0777); err != nil {
				return err
			}
		}
		now := time.Now().Format(time.RFC3339)
		prepender := func(filename string) string {
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			url := "/commands/" + strings.ToLower(base) + "/"
			return fmt.Sprintf(gendocFrontmatterTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
		}

		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "/commands/" + strings.ToLower(base) + "/"
		}

		notepad.FEEDBACK.Println("Generating Gean command-line documentation in", gendocdir, "...")
		doc.GenMarkdownTreeCustom(cmd.Root(), gendocdir, prepender, linkHandler)
		notepad.FEEDBACK.Println("Done.")

		return nil
	},
}

func init() {
	gendocCmd.PersistentFlags().StringVar(&gendocdir, "dir", "/tmp/geandoc/", "the directory to write the doc.")

	// For bash-completion
	gendocCmd.PersistentFlags().SetAnnotation("dir", goman.BashCompSubdirsInDir, []string{})
}
