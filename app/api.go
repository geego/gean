package app

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geego/gean/core"
	"github.com/gostores/symwalk"
)

type NewArticle struct {
	Name    string
	Content string
}

type OldArticle struct {
	Content string
}

type CacheArticleInfo struct {
	Name    string
	Path    string
	Create  time.Time
	Article *ArticleConfig
}

var articleCache map[string]CacheArticleInfo

func hashPath(path string) string {
	md5Hex := md5.Sum([]byte(path))
	return hex.EncodeToString(md5Hex[:])
}

func replyJSON(ctx *core.Context, status int, data interface{}) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		http.Error(ctx.Res, err.Error(), http.StatusInternalServerError)
		ctx.Stop()
		return
	}
	if status == http.StatusOK {
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Res.Write(jsonStr)
	} else {
		Warn(data)
		http.Error(ctx.Res, data.(string), status)
	}
	ctx.Stop()
}

func UpdateArticleCache() {
	articleCache = make(map[string]CacheArticleInfo, 0)
	symwalk.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		fileExt := strings.ToLower(filepath.Ext(path))
		if fileExt == ".md" {
			fileName := strings.TrimPrefix(strings.TrimSuffix(strings.ToLower(path), ".md"), "template/source/article")
			config, _ := ParseArticleConfig(path)
			id := hashPath(path)
			articleCache[string(id)] = CacheArticleInfo{
				Name:    fileName,
				Path:    path,
				Create:  ParseDate(config.Create),
				Article: config,
			}
		}
		return nil
	})
}

func ApiListArticle(ctx *core.Context) {
	UpdateArticleCache()
	replyJSON(ctx, http.StatusOK, articleCache)
}

func ApiGetArticle(ctx *core.Context) {
	UpdateArticleCache()
	article, ok := articleCache[ctx.Param["id"]]
	if !ok {
		replyJSON(ctx, http.StatusNotFound, "Not Found")
		return
	}
	filePath := article.Path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	replyJSON(ctx, http.StatusOK, string(data))
}

func ApiRemoveArticle(ctx *core.Context) {
	UpdateArticleCache()
	article, ok := articleCache[ctx.Param["id"]]
	if !ok {
		replyJSON(ctx, http.StatusNotFound, "Not Found")
		return
	}
	filePath := article.Path
	err := os.Remove(filePath)
	if err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	replyJSON(ctx, http.StatusOK, nil)
}

func ApiCreateArticle(ctx *core.Context) {
	decoder := json.NewDecoder(ctx.Req.Body)
	var article NewArticle
	err := decoder.Decode(&article)
	if err != nil {
		replyJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}
	filePath := filepath.Join(sourcePath, article.Name+".md")
	err = ioutil.WriteFile(filePath, []byte(article.Content), 0644)
	if err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	replyJSON(ctx, http.StatusOK, map[string]string{
		"id": hashPath(filePath),
	})
}

func ApiSaveArticle(ctx *core.Context) {
	UpdateArticleCache()
	decoder := json.NewDecoder(ctx.Req.Body)
	var article OldArticle
	err := decoder.Decode(&article)
	if err != nil {
		replyJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}
	cacheArticle, ok := articleCache[ctx.Param["id"]]
	if !ok {
		replyJSON(ctx, http.StatusNotFound, "Not Found")
		return
	}
	// Write
	path := cacheArticle.Path
	err = ioutil.WriteFile(path, []byte(article.Content), 0644)
	if err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	replyJSON(ctx, http.StatusOK, nil)
}

func getFormFile(ctx *core.Context, field string) (data []byte, handler *multipart.FileHeader, err error) {
	file, handler, err := ctx.Req.FormFile(field)
	if err != nil {
		replyJSON(ctx, http.StatusBadRequest, err.Error())
		return nil, handler, err
	}
	data, err = ioutil.ReadAll(file)
	if err != nil {
		replyJSON(ctx, http.StatusBadRequest, err.Error())
		return data, handler, err
	}
	return data, handler, err
}

func ApiUploadFile(ctx *core.Context) {
	UpdateArticleCache()
	fileData, handler, err := getFormFile(ctx, "file")
	if err != nil {
		replyJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}
	articleId := ctx.Req.FormValue("article_id")
	article, ok := articleCache[articleId]
	if !ok {
		replyJSON(ctx, http.StatusNotFound, "Not Found")
		return
	}
	fileDirPath := filepath.Join(sourcePath, "image", article.Name)
	err = os.MkdirAll(fileDirPath, 0777)
	if err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if err = ioutil.WriteFile(filepath.Join(fileDirPath, handler.Filename), fileData, 0777); err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	replyJSON(ctx, http.StatusOK, map[string]string{
		"path": "-/" + filepath.Join("image", article.Name, handler.Filename),
	})
}

func ApiGetConfig(ctx *core.Context) {
	filePath := filepath.Join(rootPath, "config.yml")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	replyJSON(ctx, http.StatusOK, string(data))
}

func ApiSaveConfig(ctx *core.Context) {
	content, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	filePath := filepath.Join(rootPath, "config.yml")
	err = ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		replyJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	replyJSON(ctx, http.StatusOK, nil)
}