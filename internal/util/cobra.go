package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadArgsFromAltSources is a WORKAROUND to make cobra count env vars and positional arguments when checking required flags
func LoadFlagsFromAltSources(cmd *cobra.Command, args []string) {

	if len(args) > 0 {
		_ = cmd.Flags().Set("url", args[0])

		if len(args) > 1 {
			_ = cmd.Flags().Set("message", args[1])
		}

		return
	}

	if hasURLInEnvButNotFlag(cmd) {
		_ = cmd.Flags().Set("url", viper.GetViper().GetString("SHOUTRRR_URL"))
	}
}

func hasURLInEnvButNotFlag(cmd *cobra.Command) bool {
	s, _ := cmd.Flags().GetString("url")
	return s == "" && viper.GetViper().GetString("SHOUTRRR_URL") != ""
}
