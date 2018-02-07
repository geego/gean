package command

import (
	"github.com/gostores/goman"
	"github.com/gostores/notepad"
)

var autocompleteTarget string

// bash for now (zsh and others will come)
var autocompleteType string

var genautocompleteCmd = &goman.Command{
	Use:   "autocomplete",
	Short: "Generate shell autocompletion script for Gean",
	Long: `Generates a shell autocompletion script for Gean.

NOTE: The current version supports Bash only.
      This should work for *nix systems with Bash installed.

By default, the file is written directly to /etc/bash_completion.d
for convenience, and the command may need superuser rights, e.g.:

	$ sudo gean gen autocomplete

Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
file-path and name.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /etc/bash_completion`,

	RunE: func(cmd *goman.Command, args []string) error {
		if autocompleteType != "bash" {
			return newUserError("Only Bash is supported for now")
		}

		err := cmd.Root().GenBashCompletionFile(autocompleteTarget)

		if err != nil {
			return err
		}

		notepad.FEEDBACK.Println("Bash completion file for Gean saved to", autocompleteTarget)

		return nil
	},
}

func init() {
	genautocompleteCmd.PersistentFlags().StringVarP(&autocompleteTarget, "completionfile", "", "/etc/bash_completion.d/gean.sh", "autocompletion file")
	genautocompleteCmd.PersistentFlags().StringVarP(&autocompleteType, "type", "", "bash", "autocompletion type (currently only bash supported)")

	// For bash-completion
	genautocompleteCmd.PersistentFlags().SetAnnotation("completionfile", goman.BashCompFilenameExt, []string{})
}
