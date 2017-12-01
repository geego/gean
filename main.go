package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/geego/gean/app"
)

func main() {
	cmd := cli.NewApp()
	cmd.Name = "gean"
	cmd.Usage = "a simple podcast generator"
	cmd.Version = "0.1.0"
	cmd.Commands = []cli.Command{
		{
			Name:  "build",
			Usage: "Generate source to public folder",
			Action: func(c *cli.Context) {
				app.ParseGlobalConfigByCli(c, false)
				app.Build()
			},
		},
		{
			Name:  "preview",
			Usage: "Run in server mode to preview",
			Action: func(c *cli.Context) {
				app.ParseGlobalConfigByCli(c, true)
				app.Build()
				app.Watch()
				app.Serve()
			},
		},
		{
			Name:  "publish",
			Usage: "Generate source to public folder and publish",
			Action: func(c *cli.Context) {
				app.ParseGlobalConfigByCli(c, false)
				app.Build()
				app.Publish()
			},
		},
		{
			Name:  "serve",
			Usage: "Run in server mode",
			Action: func(c *cli.Context) {
				app.ParseGlobalConfigByCli(c, true)
				app.Build()
				app.Watch()
				app.Serve()
			},
		},
	}
	cmd.Run(os.Args)
}
