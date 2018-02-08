package command

import (
	"path/filepath"

	"github.com/geego/gean/app/geanlib"
	"github.com/gostores/goman"
	"github.com/gostores/notepad"
)

func init() {
	listCmd.AddCommand(listDraftsCmd)
	listCmd.AddCommand(listFutureCmd)
	listCmd.AddCommand(listExpiredCmd)
	listCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "filesystem path to read files relative from")
	listCmd.PersistentFlags().SetAnnotation("source", goman.BashCompSubdirsInDir, []string{})
}

var listCmd = &goman.Command{
	Use:   "list",
	Short: "Listing out various types of content",
	Long: `Listing out various types of content.

List requires a subcommand, e.g. ` + "`gean list drafts`.",
	RunE: nil,
}

var listDraftsCmd = &goman.Command{
	Use:   "drafts",
	Short: "List all drafts",
	Long:  `List all of the drafts in your content directory.`,
	RunE: func(cmd *goman.Command, args []string) error {

		cfg, err := InitializeConfig()
		if err != nil {
			return err
		}

		c, err := newCommandeer(cfg)
		if err != nil {
			return err
		}

		c.Set("buildDrafts", true)

		sites, err := geanlib.NewHugoSites(*cfg)

		if err != nil {
			return newSystemError("Error creating sites", err)
		}

		if err := sites.Build(geanlib.BuildCfg{SkipRender: true}); err != nil {
			return newSystemError("Error Processing Source Content", err)
		}

		for _, p := range sites.Pages() {
			if p.IsDraft() {
				notepad.FEEDBACK.Println(filepath.Join(p.File.Dir(), p.File.LogicalName()))
			}

		}

		return nil

	},
}

var listFutureCmd = &goman.Command{
	Use:   "future",
	Short: "List all posts dated in the future",
	Long: `List all of the posts in your content directory which will be
posted in the future.`,
	RunE: func(cmd *goman.Command, args []string) error {

		cfg, err := InitializeConfig()
		if err != nil {
			return err
		}

		c, err := newCommandeer(cfg)
		if err != nil {
			return err
		}

		c.Set("buildFuture", true)

		sites, err := geanlib.NewHugoSites(*cfg)

		if err != nil {
			return newSystemError("Error creating sites", err)
		}

		if err := sites.Build(geanlib.BuildCfg{SkipRender: true}); err != nil {
			return newSystemError("Error Processing Source Content", err)
		}

		for _, p := range sites.Pages() {
			if p.IsFuture() {
				notepad.FEEDBACK.Println(filepath.Join(p.File.Dir(), p.File.LogicalName()))
			}

		}

		return nil

	},
}

var listExpiredCmd = &goman.Command{
	Use:   "expired",
	Short: "List all posts already expired",
	Long: `List all of the posts in your content directory which has already
expired.`,
	RunE: func(cmd *goman.Command, args []string) error {

		cfg, err := InitializeConfig()
		if err != nil {
			return err
		}

		c, err := newCommandeer(cfg)
		if err != nil {
			return err
		}

		c.Set("buildExpired", true)

		sites, err := geanlib.NewHugoSites(*cfg)

		if err != nil {
			return newSystemError("Error creating sites", err)
		}

		if err := sites.Build(geanlib.BuildCfg{SkipRender: true}); err != nil {
			return newSystemError("Error Processing Source Content", err)
		}

		for _, p := range sites.Pages() {
			if p.IsExpired() {
				notepad.FEEDBACK.Println(filepath.Join(p.File.Dir(), p.File.LogicalName()))
			}

		}

		return nil

	},
}
