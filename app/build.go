package app

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gostores/feeds"
	"github.com/gostores/symwalk"
)

// Parse config
var articleTpl, pageTpl, archiveTpl, tagTpl template.Template
var themePath, publicPath, sourcePath string

// For concurrency
var wg sync.WaitGroup

// Data struct
type ArticleInfo struct {
	DetailDate int64
	Create     string
	Title      string
	Link       string
	Top        bool
}

type Archive struct {
	Year     string
	Articles Collections
}

type Tag struct {
	Name     string
	Count    int
	Articles Collections
}

// For sort
type Collections []interface{}

func (v Collections) Len() int      { return len(v) }
func (v Collections) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v Collections) Less(i, j int) bool {
	switch v[i].(type) {
	case ArticleInfo:
		return v[i].(ArticleInfo).DetailDate > v[j].(ArticleInfo).DetailDate
	case Article:
		article1 := v[i].(Article)
		article2 := v[j].(Article)
		if article1.Top && !article2.Top {
			return true
		} else if !article1.Top && article2.Top {
			return false
		} else {
			return article1.Create > article2.Create
		}
	case Archive:
		return v[i].(Archive).Year > v[j].(Archive).Year
	case Tag:
		if v[i].(Tag).Count == v[j].(Tag).Count {
			return v[i].(Tag).Name > v[j].(Tag).Name
		}
		return v[i].(Tag).Count > v[j].(Tag).Count
	}
	return false
}

func Build() {
	startTime := time.Now()
	var articles = make(Collections, 0)
	var tagMap = make(map[string]Collections)
	var archiveMap = make(map[string]Collections)
	// Parse config
	themePath = filepath.Join(rootPath, globalConfig.Site.Theme)
	publicPath = filepath.Join(rootPath, "public")
	sourcePath = filepath.Join(rootPath, "source")
	// Append all partial html
	var partialTpl string
	files, _ := filepath.Glob(filepath.Join(themePath, "*.html"))
	for _, path := range files {
		fileExt := strings.ToLower(filepath.Ext(path))
		baseName := strings.ToLower(filepath.Base(path))
		if fileExt == ".html" && strings.HasPrefix(baseName, "_") {
			html, err := ioutil.ReadFile(path)
			if err != nil {
				Fatal(err.Error())
			}
			tplName := strings.TrimPrefix(baseName, "_")
			tplName = strings.TrimSuffix(tplName, ".html")
			htmlStr := "{{define \"" + tplName + "\"}}" + string(html) + "{{end}}"
			partialTpl += htmlStr
		}
	}
	// Compile template
	articleTpl = CompileTpl(filepath.Join(themePath, "article.html"), partialTpl, "article")
	pageTpl = CompileTpl(filepath.Join(themePath, "page.html"), partialTpl, "page")
	archiveTpl = CompileTpl(filepath.Join(themePath, "archive.html"), partialTpl, "archive")
	tagTpl = CompileTpl(filepath.Join(themePath, "tag.html"), partialTpl, "tag")
	// Clean public folder
	cleanPatterns := []string{"post", "tag", "image", "js", "css", "*.html", "favicon.ico", "robots.txt"}
	for _, pattern := range cleanPatterns {
		files, _ := filepath.Glob(filepath.Join(publicPath, pattern))
		for _, path := range files {
			os.RemoveAll(path)
		}
	}
	// Find all .md to generate article
	symwalk.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		fileExt := strings.ToLower(filepath.Ext(path))
		if fileExt == ".md" {
			// Parse markdown data
			article := ParseArticle(path)
			if article == nil || article.Draft {
				return nil
			}
			// Generate page name
			fileName := strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ".md")
			Log("Building " + fileName)
			// Genetate custom link
			unixTime := time.Unix(article.Create, 0)
			linkMap := map[string]string{
				"{year}":  unixTime.Format("2006"),
				"{month}": unixTime.Format("01"),
				"{day}":   unixTime.Format("02"),
				"{title}": fileName,
			}
			var link string
			if globalConfig.Site.Link == "" {
				link = fileName + ".html"
			} else {
				link = globalConfig.Site.Link
				for key, val := range linkMap {
					link = strings.Replace(link, key, val, -1)
				}
			}
			directory := filepath.Dir(link)
			err := os.MkdirAll(filepath.Join(publicPath, directory), 0777)
			if err != nil {
				Fatal(err.Error())
			}
			// Generate file path
			article.Link = link
			article.GlobalConfig = *globalConfig
			articles = append(articles, *article)
			// Get tags info
			for _, tag := range article.Tags {
				if _, ok := tagMap[tag]; !ok {
					tagMap[tag] = make(Collections, 0)
				}
				tagMap[tag] = append(tagMap[tag], *article)
			}
			// Get archive info
			createYear := unixTime.Format("2006")
			if _, ok := archiveMap[createYear]; !ok {
				archiveMap[createYear] = make(Collections, 0)
			}
			articleInfo := ArticleInfo{
				DetailDate: article.Create,
				Create:     unixTime.Format("2006-01-02"),
				Title:      article.Title,
				Link:       article.Link,
				Top:        article.Top,
			}
			archiveMap[createYear] = append(archiveMap[createYear], articleInfo)
		}
		return nil
	})
	if len(articles) == 0 {
		Fatal("Must be have at least one article")
	}
	// Sort by create
	sort.Sort(articles)
	// Generate rss page
	wg.Add(1)
	go GeneratePodcast(articles)
	// Render article
	wg.Add(1)
	go RenderArticles(articleTpl, articles)
	// Generate article list pages
	wg.Add(1)
	go RenderArticleList("", articles, "")
	// Generate article list pages by tag
	for tagName, articles := range tagMap {
		wg.Add(1)
		go RenderArticleList(filepath.Join("tag", tagName), articles, tagName)
	}
	// Generate archive page
	archives := make(Collections, 0)
	for year, articleInfos := range archiveMap {
		// Sort by create
		sort.Sort(articleInfos)
		archives = append(archives, Archive{
			Year:     year,
			Articles: articleInfos,
		})
	}
	// Sort by year
	sort.Sort(archives)
	wg.Add(1)
	go RenderPage(archiveTpl, map[string]interface{}{
		"Total":   len(articles),
		"Archive": archives,
		"Site":    globalConfig.Site,
		"I18n":    globalConfig.I18n,
	}, filepath.Join(publicPath, "archive.html"))
	// Generate tag page
	tags := make(Collections, 0)
	for tagName, tagArticles := range tagMap {
		articleInfos := make(Collections, 0)
		for _, article := range tagArticles {
			articleValue := article.(Article)
			articleInfos = append(articleInfos, ArticleInfo{
				DetailDate: articleValue.Create,
				Create:     time.Unix(articleValue.Create, 0).Format("2006-01-02"),
				Title:      articleValue.Title,
				Link:       articleValue.Link,
				Top:        articleValue.Top,
			})
		}
		// Sort by create
		sort.Sort(articleInfos)
		tags = append(tags, Tag{
			Name:     tagName,
			Count:    len(tagArticles),
			Articles: articleInfos,
		})
	}
	// Sort by count
	sort.Sort(Collections(tags))
	wg.Add(1)
	go RenderPage(tagTpl, map[string]interface{}{
		"Total": len(articles),
		"Tag":   tags,
		"Site":  globalConfig.Site,
		"I18n":  globalConfig.I18n,
	}, filepath.Join(publicPath, "tag.html"))
	// Generate other pages
	files, _ = filepath.Glob(filepath.Join(sourcePath, "*.html"))
	for _, path := range files {
		fileExt := strings.ToLower(filepath.Ext(path))
		baseName := filepath.Base(path)
		if fileExt == ".html" && !strings.HasPrefix(baseName, "_") {
			htmlTpl := CompileTpl(path, partialTpl, baseName)
			relPath, _ := filepath.Rel(sourcePath, path)
			wg.Add(1)
			go RenderPage(htmlTpl, globalConfig, filepath.Join(publicPath, relPath))
		}
	}
	// Copy static files
	Copy()
	wg.Wait()
	endTime := time.Now()
	usedTime := endTime.Sub(startTime)
	fmt.Printf("\nFinished to build in public folder (%v)\n", usedTime)
}

// Generate podcast page
func GeneratePodcast(articles Collections) {
	defer wg.Done()
	var feedArticles Collections
	if len(articles) < globalConfig.Site.Limit {
		feedArticles = articles
	} else {
		feedArticles = articles[0:globalConfig.Site.Limit]
	}
	if globalConfig.Site.Url != "" {
		feed := &feeds.Feed{
			Title:       globalConfig.Site.Title,
			Subtitle:    globalConfig.Site.Subtitle,
			Description: globalConfig.Site.Description,
			Link:        &feeds.Link{Href: globalConfig.Site.Url, Rel: "self", Type: "application/rss+xml"},
			Author:      &feeds.Author{Name: globalConfig.Itunes.Author, Email: globalConfig.Itunes.Email},
			Itunes: &feeds.Itunes{
				Email:    globalConfig.Itunes.Email,
				Author:   globalConfig.Itunes.Author,
				Category: globalConfig.Itunes.Category,
				Language: globalConfig.Itunes.Language,
				Logo:     globalConfig.Itunes.Logo,
			},
			Created: time.Now(),
		}
		feed.Items = make([]*feeds.Item, 0)
		for _, item := range feedArticles {
			article := item.(Article)

			var audioSize int64
			if article.Audio != "" {
				fileinfo, err := os.Stat(filepath.Join(sourcePath, article.Audio))
				if err == nil {
					audioSize = fileinfo.Size()
				}
			} else {
				audioSize = 0
			}

			feed.Items = append(feed.Items, &feeds.Item{
				Title:       article.Title,
				Description: string(article.Content),
				Link:        &feeds.Link{Href: globalConfig.Site.Url + article.Link},
				Author:      &feeds.Author{article.Author.Name, ""},
				Itunes: &feeds.ItunesItem{
					Subtitle:  string(article.Summary),
					Image:     article.Poster,
					AudioHref: article.Audio,
					AudioType: "mp3",
					AudioSize: audioSize,
				},
				Created: time.Unix(article.Create, 0),
				Updated: time.Unix(article.Update, 0),
			})
		}
		if rss, err := feed.ToPodcast(); err == nil {
			err := ioutil.WriteFile(filepath.Join(publicPath, "index.xml"), []byte(rss), 0644)
			if err != nil {
				Fatal(err.Error())
			}
		} else {
			Fatal(err.Error())
		}
	}
}

// Copy static files
func Copy() {
	srcList := globalConfig.Build.Copy
	for _, source := range srcList {
		if matches, err := filepath.Glob(filepath.Join(rootPath, source)); err == nil {
			for _, srcPath := range matches {
				Log("Copying " + srcPath)
				file, err := os.Stat(srcPath)
				if err != nil {
					Fatal("Not exist: " + srcPath)
				}
				fileName := file.Name()
				desPath := filepath.Join(publicPath, fileName)
				wg.Add(1)
				if file.IsDir() {
					go CopyDir(srcPath, desPath)
				} else {
					go CopyFile(srcPath, desPath)
				}
			}
		} else {
			Fatal(err.Error())
		}
	}
}
