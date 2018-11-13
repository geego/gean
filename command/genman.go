package command

import (
	"fmt"
	"strings"

	"github.com/geego/gean/app/geanfs"
	"github.com/geego/gean/app/helpers"
	"github.com/govenue/goman"
	"github.com/govenue/goman/doc"
	"github.com/govenue/notepad"
)

var genmandir string
var genmanCmd = &goman.Command{
	Use:   "man",
	Short: "Generate man pages for the Gean CLI",
	Long: `This command automatically generates up-to-date man pages of Gean's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,

	RunE: func(cmd *goman.Command, args []string) error {
		header := &doc.GenManHeader{
			Section: "1",
			Manual:  "Gean Manual",
			Source:  fmt.Sprintf("Gean %s", helpers.CurrentHugoVersion),
		}
		if !strings.HasSuffix(genmandir, helpers.FilePathSeparator) {
			genmandir += helpers.FilePathSeparator
		}
		if found, _ := helpers.Exists(genmandir, geanfs.Os); !found {
			notepad.FEEDBACK.Println("Directory", genmandir, "does not exist, creating...")
			if err := geanfs.Os.MkdirAll(genmandir, 0777); err != nil {
				return err
			}
		}
		cmd.Root().DisableAutoGenTag = true

		notepad.FEEDBACK.Println("Generating Gean man pages in", genmandir, "...")
		doc.GenManTree(cmd.Root(), header, genmandir)

		notepad.FEEDBACK.Println("Done.")

		return nil
	},
}

func init() {
	genmanCmd.PersistentFlags().StringVar(&genmandir, "dir", "man/", "the directory to write the man pages.")

	// For bash-completion
	genmanCmd.PersistentFlags().SetAnnotation("dir", goman.BashCompSubdirsInDir, []string{})
}
