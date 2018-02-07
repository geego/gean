package command

import (
	"reflect"
	"sort"

	"github.com/gostores/configurator"
	"github.com/gostores/goman"
	"github.com/gostores/notepad"
)

var configCmd = &goman.Command{
	Use:   "config",
	Short: "Print the site configuration",
	Long:  `Print the site configuration, both default and custom settings.`,
}

func init() {
	configCmd.RunE = printConfig
}

func printConfig(cmd *goman.Command, args []string) error {
	cfg, err := InitializeConfig(configCmd)

	if err != nil {
		return err
	}

	allSettings := cfg.Cfg.(*configurator.Configurator).AllSettings()

	var separator string
	if allSettings["metadataformat"] == "toml" {
		separator = " = "
	} else {
		separator = ": "
	}

	var keys []string
	for k := range allSettings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		kv := reflect.ValueOf(allSettings[k])
		if kv.Kind() == reflect.String {
			notepad.FEEDBACK.Printf("%s%s\"%+v\"\n", k, separator, allSettings[k])
		} else {
			notepad.FEEDBACK.Printf("%s%s%+v\n", k, separator, allSettings[k])
		}
	}

	return nil
}
