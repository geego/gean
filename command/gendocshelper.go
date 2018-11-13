package command

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/geego/gean/app/docshelper"
	"github.com/govenue/goman"
)

type genDocsHelper struct {
	target string
	cmd    *goman.Command
}

func createGenDocsHelper() *genDocsHelper {
	g := &genDocsHelper{
		cmd: &goman.Command{
			Use:    "docshelper",
			Short:  "Generate some data files for the Gean docs.",
			Hidden: true,
		},
	}

	g.cmd.RunE = func(cmd *goman.Command, args []string) error {
		return g.generate()
	}

	g.cmd.PersistentFlags().StringVarP(&g.target, "dir", "", "docs/data", "data dir")

	return g
}

func (g *genDocsHelper) generate() error {
	fmt.Println("Generate docs data to", g.target)

	targetFile := filepath.Join(g.target, "docs.json")

	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	if err := enc.Encode(docshelper.DocProviders); err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil

}
