package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gostores/configurator"
	"github.com/gostores/fsintra"
	"github.com/gostores/goman"
	"github.com/gostores/notepad"

	"yiqilai.tech/gean/app/create"
	"yiqilai.tech/gean/app/geanfs"
	"yiqilai.tech/gean/app/geanlib"
	"yiqilai.tech/gean/app/helpers"
	"yiqilai.tech/gean/app/parser"
)

var (
	configFormat  string
	contentEditor string
	contentType   string
)

func init() {
	newSiteCmd.Flags().StringVarP(&configFormat, "format", "f", "toml", "config & frontmatter format")
	newSiteCmd.Flags().Bool("force", false, "init inside non-empty directory")
	newCmd.Flags().StringVarP(&contentType, "kind", "k", "", "content type to create")
	newCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "filesystem path to read files relative from")
	newCmd.PersistentFlags().SetAnnotation("source", goman.BashCompSubdirsInDir, []string{})
	newCmd.Flags().StringVar(&contentEditor, "editor", "", "edit new content with this editor, if provided")

	newCmd.AddCommand(newSiteCmd)
	newCmd.AddCommand(newThemeCmd)

}

var newCmd = &goman.Command{
	Use:   "new [path]",
	Short: "Create new content for your site",
	Long: `Create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.

You can also specify the kind with ` + "`-k KIND`" + `.

If archetypes are provided in your theme or site, they will be used.`,

	RunE: NewContent,
}

var newSiteCmd = &goman.Command{
	Use:   "site [path]",
	Short: "Create a new site (skeleton)",
	Long: `Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use ` + "`gean new [contentPath]`" + ` to create new content.`,
	RunE: NewSite,
}

var newThemeCmd = &goman.Command{
	Use:   "theme [name]",
	Short: "Create a new theme",
	Long: `Create a new theme (skeleton) called [name] in the current directory.
New theme is a skeleton. Please add content to the touched files. Add your
name to the copyright line in the license and adjust the theme.toml file
as you see fit.`,
	RunE: NewTheme,
}

// NewContent adds new content to a gean site.
func NewContent(cmd *goman.Command, args []string) error {
	cfg, err := InitializeConfig()

	if err != nil {
		return err
	}

	c, err := newCommandeer(cfg)
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("editor") {
		c.Set("newContentEditor", contentEditor)
	}

	if len(args) < 1 {
		return newUserError("path needs to be provided")
	}

	createPath := args[0]

	var kind string

	createPath, kind = newContentPathSection(createPath)

	if contentType != "" {
		kind = contentType
	}

	ps, err := helpers.NewPathSpec(cfg.Fs, cfg.Cfg)
	if err != nil {
		return err
	}

	// If a site isn't in use in the archetype template, we can skip the build.
	siteFactory := func(filename string, siteUsed bool) (*geanlib.Site, error) {
		if !siteUsed {
			return geanlib.NewSite(*cfg)
		}
		var s *geanlib.Site
		if err := c.initSites(); err != nil {
			return nil, err
		}

		if err := Hugo.Build(geanlib.BuildCfg{SkipRender: true, PrintStats: false}); err != nil {
			return nil, err
		}

		s = Hugo.Sites[0]

		if len(Hugo.Sites) > 1 {
			// Find the best match.
			for _, ss := range Hugo.Sites {
				if strings.Contains(createPath, "."+ss.Language.Lang) {
					s = ss
					break
				}
			}
		}
		return s, nil
	}

	return create.NewContent(ps, siteFactory, kind, createPath)
}

func doNewSite(fs *geanfs.Fs, basepath string, force bool) error {
	archeTypePath := filepath.Join(basepath, "archetypes")
	dirs := []string{
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		archeTypePath,
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
		filepath.Join(basepath, "themes"),
	}

	if exists, _ := helpers.Exists(basepath, fs.Source); exists {
		if isDir, _ := helpers.IsDir(basepath, fs.Source); !isDir {
			return errors.New(basepath + " already exists but not a directory")
		}

		isEmpty, _ := helpers.IsEmpty(basepath, fs.Source)

		switch {
		case !isEmpty && !force:
			return errors.New(basepath + " already exists and is not empty")

		case !isEmpty && force:
			all := append(dirs, filepath.Join(basepath, "config."+configFormat))
			for _, path := range all {
				if exists, _ := helpers.Exists(path, fs.Source); exists {
					return errors.New(path + " already exists")
				}
			}
		}
	}

	for _, dir := range dirs {
		if err := fs.Source.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("Failed to create dir: %s", err)
		}
	}

	createConfig(fs, basepath, configFormat)

	// Create a defaul archetype file.
	helpers.SafeWriteToDisk(filepath.Join(archeTypePath, "default.md"),
		strings.NewReader(create.ArchetypeTemplateTemplate), fs.Source)

	notepad.FEEDBACK.Printf("Congratulations! Your new gean site is created in %s.\n\n", basepath)
	notepad.FEEDBACK.Println(nextStepsText())

	return nil
}

func nextStepsText() string {
	var nextStepsText bytes.Buffer

	nextStepsText.WriteString(`Just a few more steps and you're ready to go:

1. Download a theme into the same-named folder.
   create your own with the "gean new theme <THEMENAME>" command.
2. Perhaps you want to add some content. You can add single files
   with "gean new `)

	nextStepsText.WriteString(filepath.Join("<SECTIONNAME>", "<FILENAME>.<FORMAT>"))

	nextStepsText.WriteString(`".
3. Start the built-in live server via "gean server".

`)

	return nextStepsText.String()
}

// NewSite creates a new Hugo site and initializes a structured Hugo directory.
func NewSite(cmd *goman.Command, args []string) error {
	if len(args) < 1 {
		return newUserError("path needs to be provided")
	}

	createpath, err := filepath.Abs(filepath.Clean(args[0]))
	if err != nil {
		return newUserError(err)
	}

	forceNew, _ := cmd.Flags().GetBool("force")

	return doNewSite(geanfs.NewDefault(configurator.New()), createpath, forceNew)
}

// NewTheme creates a new Hugo theme.
func NewTheme(cmd *goman.Command, args []string) error {
	cfg, err := InitializeConfig()

	if err != nil {
		return err
	}

	if len(args) < 1 {
		return newUserError("theme name needs to be provided")
	}

	c, err := newCommandeer(cfg)
	if err != nil {
		return err
	}

	createpath := c.PathSpec().AbsPathify(filepath.Join(c.Cfg.GetString("themesDir"), args[0]))
	notepad.INFO.Println("creating theme at", createpath)

	if x, _ := helpers.Exists(createpath, cfg.Fs.Source); x {
		return newUserError(createpath, "already exists")
	}

	mkdir(createpath, "layouts", "_default")
	mkdir(createpath, "layouts", "partials")

	touchFile(cfg.Fs.Source, createpath, "layouts", "index.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "404.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "_default", "list.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "_default", "single.html")

	touchFile(cfg.Fs.Source, createpath, "layouts", "partials", "header.html")
	touchFile(cfg.Fs.Source, createpath, "layouts", "partials", "footer.html")

	mkdir(createpath, "archetypes")

	archDefault := []byte("+++\n+++\n")

	err = helpers.WriteToDisk(filepath.Join(createpath, "archetypes", "default.md"), bytes.NewReader(archDefault), cfg.Fs.Source)
	if err != nil {
		return err
	}

	mkdir(createpath, "static", "js")
	mkdir(createpath, "static", "css")

	by := []byte(`The MIT License (MIT)

Copyright (c) ` + time.Now().Format("2006") + ` YOUR_NAME_HERE

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`)

	err = helpers.WriteToDisk(filepath.Join(createpath, "LICENSE.md"), bytes.NewReader(by), cfg.Fs.Source)
	if err != nil {
		return err
	}

	createThemeMD(cfg.Fs, createpath)

	return nil
}

func mkdir(x ...string) {
	p := filepath.Join(x...)

	err := os.MkdirAll(p, 0777) // before umask
	if err != nil {
		notepad.FATAL.Fatalln(err)
	}
}

func touchFile(fs fsintra.Fs, x ...string) {
	inpath := filepath.Join(x...)
	mkdir(filepath.Dir(inpath))
	err := helpers.WriteToDisk(inpath, bytes.NewReader([]byte{}), fs)
	if err != nil {
		notepad.FATAL.Fatalln(err)
	}
}

func createThemeMD(fs *geanfs.Fs, inpath string) (err error) {

	by := []byte(`# theme.toml template for a gean theme
# See https://yiqilai.tech/gean/appThemes#themetoml for an example

name = "` + strings.Title(helpers.MakeTitle(filepath.Base(inpath))) + `"
license = "MIT"
licenselink = "https://github.com/yourname/yourtheme/blob/master/LICENSE.md"
description = ""
homepage = "http://example.com/"
tags = []
features = []
min_version = "0.31"

[author]
  name = ""
  homepage = ""

# If porting an existing theme
[original]
  name = ""
  homepage = ""
  repo = ""
`)

	err = helpers.WriteToDisk(filepath.Join(inpath, "theme.toml"), bytes.NewReader(by), fs.Source)
	if err != nil {
		return
	}

	return nil
}

func newContentPathSection(path string) (string, string) {
	// Forward slashes is used in all examples. Convert if needed.
	// Issue #1133
	createpath := filepath.FromSlash(path)
	var section string
	// assume the first directory is the section (kind)
	if strings.Contains(createpath[1:], helpers.FilePathSeparator) {
		section = helpers.GuessSection(createpath)
	}

	return createpath, section
}

func createConfig(fs *geanfs.Fs, inpath string, kind string) (err error) {
	in := map[string]string{
		"baseURL":      "http://example.org/",
		"title":        "my new gean site",
		"languageCode": "en-us",
	}
	kind = parser.FormatSanitize(kind)

	var buf bytes.Buffer
	err = parser.InterfaceToConfig(in, parser.FormatToLeadRune(kind), &buf)
	if err != nil {
		return err
	}

	return helpers.WriteToDisk(filepath.Join(inpath, "config."+kind), &buf, fs.Source)
}
