package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geego/gean/core"
	"github.com/gorilla/websocket"
	"github.com/gostores/fsnotify"
	"github.com/gostores/symwalk"
)

var watcher *fsnotify.Watcher
var conn *websocket.Conn

func Watch() {
	// Listen watched file change event
	if watcher != nil {
		watcher.Close()
	}
	watcher, _ = fsnotify.NewWatcher()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Write {
					// Handle when file change
					fmt.Println(event.Name)
					Build()
					if conn != nil {
						if err := conn.WriteMessage(websocket.TextMessage, []byte("change")); err != nil {
							Warn(err.Error())
						}
					}
				}
			case err := <-watcher.Errors:
				Warn(err.Error())
			}
		}
	}()
	var dirs = []string{"source"}
	for _, source := range dirs {
		dirPath := filepath.Join(rootPath, source)
		symwalk.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				if err := watcher.Add(path); err != nil {
					Warn(err.Error())
				}
			}
			return nil
		})
	}
}

func Websocket(ctx *core.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	if c, err := upgrader.Upgrade(ctx.Res, ctx.Req, nil); err != nil {
		Fatal(err)
	} else {
		conn = c
	}
	ctx.Stop()
}

func Serve() {
	editorWeb := core.New()

	editorWeb.Get("/articles", ApiListArticle)
	editorWeb.Get("/articles/:id", ApiGetArticle)
	editorWeb.Post("/articles", ApiCreateArticle)
	editorWeb.Put("/articles/:id", ApiSaveArticle)
	editorWeb.Delete("/articles/:id", ApiRemoveArticle)
	editorWeb.Get("/config", ApiGetConfig)
	editorWeb.Put("/config", ApiSaveConfig)
	editorWeb.Post("/upload", ApiUploadFile)
	editorWeb.Use(core.Cors)
	editorWeb.Get("*", core.Static(filepath.Join("editor/assets")))

	Log("Access http://localhost:" + globalConfig.Build.Port + "/ to open editor")
	go editorWeb.Listen(":" + globalConfig.Build.Port)

	previewWeb := core.New()
	previewWeb.Get("/live", Websocket)
	previewWeb.Get("*", core.Static(filepath.Join(rootPath, "public")))

	Log("Access http://localhost:9090" + "/ to open preview")
	previewWeb.Listen(":9090")
}
