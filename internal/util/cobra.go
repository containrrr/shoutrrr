package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MoveEnvVarToFlag is a WORKAROUND to make cobra count env vars when checking required flags
func MoveEnvVarToFlag(cmd *cobra.Command, args []string) {

	if len(args) == 1 {
		cmd.Flags().Set("url", args[0])
		return
	}

	if hasURLInEnvButNotFlag(cmd) {
		cmd.Flags().Set("url", viper.GetViper().GetString("SHOUTRRR_URL"))
		return
	}
}

func hasURLInEnvButNotFlag(cmd *cobra.Command) bool {
	s, _ := cmd.Flags().GetString("url")
	return s == "" && viper.GetViper().GetString("SHOUTRRR_URL") != ""
}
