package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadFlagsFromAltSources is a WORKAROUND to make cobra count env vars and
// positional arguments when checking required flags.
func LoadFlagsFromAltSources(cmd *cobra.Command, args []string) {
	flags := cmd.Flags()

	if len(args) > 0 {
		_ = flags.Set("url", args[0])

		if len(args) > 1 {
			_ = flags.Set("message", args[1])
		}

		return
	}

	if hasURLInEnvButNotFlag(cmd) {
		_ = flags.Set("url", viper.GetViper().GetString("SHOUTRRR_URL"))

		// If the URL has been set in ENV, default the message to read from stdin.
		if msg, _ := flags.GetString("message"); msg == "" {
			flags.Set("message", "-")
		}
	}
}

func hasURLInEnvButNotFlag(cmd *cobra.Command) bool {
	s, _ := cmd.Flags().GetString("url")

	return s == "" && viper.GetViper().GetString("SHOUTRRR_URL") != ""
}
