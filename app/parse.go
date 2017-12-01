package app

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/gostores/encoding/markdown"
	"github.com/gostores/encoding/yaml"
)

var globalConfig *GlobalConfig
var rootPath string

const (
	VERSION      = "0.1.0"
	DEFAULT_ROOT = ""
)

type SiteConfig struct {
	Root        string
	Title       string
	Subtitle    string
	Description string
	Logo        string
	Limit       int
	Theme       string
	Comment     string
	Lang        string
	Link        string
	Url         string
	Config      interface{}
}

type ItunesConfig struct {
	Category string
	Language string
	Author   string
	Email    string
	Logo     string
}

type AuthorConfig struct {
	Id     string
	Name   string
	Email  string
	Intro  string
	Avatar string
}

type BuildConfig struct {
	Port    string
	Watch   bool
	Copy    []string
	Publish string
}

type GlobalConfig struct {
	I18n    map[string]string
	Site    SiteConfig
	Itunes  ItunesConfig
	Authors map[string]AuthorConfig
	Build   BuildConfig
	Develop bool
}

type ArticleConfig struct {
	Title   string
	Create  string
	Update  string
	Author  string
	Poster  string
	Audio   string
	Topic   string
	Tags    []string
	Draft   bool
	Top     bool
	Summary template.HTML
	Config  interface{}
}

type Article struct {
	GlobalConfig
	ArticleConfig
	Create  int64
	Update  int64
	Author  AuthorConfig
	Summary template.HTML
	Content template.HTML
	Tags    []string
	Link    string
	Config  interface{}
}

type ThemeConfig struct {
	Copy []string
	Lang map[string]map[string]string
}

const (
	CONFIG_SPLIT = "+++"
	MORE_SPLIT   = "<!--more-->"
)

func Parse(markdownContent string) template.HTML {
	// html.UnescapeString
	return template.HTML(markdown.MarkdownCommon([]byte(markdownContent)))
}

func ReplaceRootFlag(content string) string {
	return strings.Replace(content, "-/", globalConfig.Site.Root+"/", -1)
}

func ParseGlobalConfig(configPath string, develop bool) *GlobalConfig {
	var config *GlobalConfig
	// Parse Global Config
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		Fatal(err.Error())
	}
	if config.Site.Config == nil {
		config.Site.Config = ""
	}
	config.Develop = develop
	if develop {
		config.Site.Root = ""
	}
	config.Site.Logo = strings.Replace(config.Site.Logo, "-/", config.Site.Root+"/", -1)
	// Parse Theme Config
	themeConfig := ParseThemeConfig(filepath.Join(rootPath, config.Site.Theme, "config.yml"))
	for _, copyItem := range themeConfig.Copy {
		config.Build.Copy = append(config.Build.Copy, filepath.Join(config.Site.Theme, copyItem))
	}
	config.I18n = make(map[string]string)
	for item, langItem := range themeConfig.Lang {
		config.I18n[item] = langItem[config.Site.Lang]
	}
	return config
}

func ParseThemeConfig(configPath string) *ThemeConfig {
	// Read data from file
	var themeConfig *ThemeConfig
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		Fatal(err.Error())
	}
	// Parse config content
	if err := yaml.Unmarshal(data, &themeConfig); err != nil {
		Fatal(err.Error())
	}
	return themeConfig
}

func ParseArticleConfig(markdownPath string) (config *ArticleConfig, content string) {
	var configStr string
	// Read data from file
	data, err := ioutil.ReadFile(markdownPath)
	if err != nil {
		Fatal(err.Error())
	}
	// Split config and markdown
	contentStr := string(data)
	contentStr = ReplaceRootFlag(contentStr)
	markdownStr := strings.SplitN(contentStr, CONFIG_SPLIT, 2)
	contentLen := len(markdownStr)
	if contentLen > 0 {
		configStr = markdownStr[0]
	}
	if contentLen > 1 {
		content = markdownStr[1]
	}
	// Parse config content
	if err := yaml.Unmarshal([]byte(configStr), &config); err != nil {
		Error(err.Error())
		return nil, ""
	}
	if config == nil {
		return nil, ""
	}
	// Parse Summary splited by MORE_SPLIT
	summaryAry := strings.SplitN(content, MORE_SPLIT, 2)
	if len(config.Summary) <= 0 && len(summaryAry) > 1 {
		config.Summary = Parse(summaryAry[0])
		content = strings.Replace(content, MORE_SPLIT, "", 1)
	}
	return config, content
}

func ParseArticle(markdownPath string) *Article {
	config, content := ParseArticleConfig(markdownPath)
	if config == nil {
		Error("Invalid format: " + markdownPath)
		return nil
	}
	if config.Config == nil {
		config.Config = ""
	}
	var article Article
	// Parse markdown content
	article.Summary = config.Summary
	article.Config = config.Config
	article.Content = Parse(content)
	article.Create = ParseDate(config.Create).Unix()
	article.Title = config.Title
	article.Tags = config.Tags
	article.Topic = config.Topic
	article.Draft = config.Draft
	article.Top = config.Top
	if config.Update != "" {
		article.Update = ParseDate(config.Update).Unix()
	}
	if author, ok := globalConfig.Authors[config.Author]; ok {
		author.Id = config.Author
		author.Avatar = ReplaceRootFlag(author.Avatar)
		article.Author = author
	}
	// Support topic and poster field
	if config.Poster != "" {
		article.Poster = config.Poster
	} else {
		article.Poster = config.Topic
	}

	// Support audio field
	if config.Audio != "" {
		article.Audio = config.Audio
	} else {
		article.Audio = config.Topic
	}

	return &article
}

func ParseGlobalConfigByCli(c *cli.Context, develop bool) {
	if len(c.Args()) > 0 {
		rootPath = c.Args()[0]
	} else {
		rootPath = "."
	}
	ParseGlobalConfigWrap(rootPath, develop)
	if globalConfig == nil {
		ParseGlobalConfigWrap(DEFAULT_ROOT, develop)
		if globalConfig == nil {
			Fatal("Parse config.yml failed, please specify a valid path")
		}
	}
}

func ParseGlobalConfigWrap(root string, develop bool) {
	rootPath = root
	globalConfig = ParseGlobalConfig(filepath.Join(rootPath, "config.yml"), develop)
	if globalConfig == nil {
		return
	}
}
