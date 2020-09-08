package root

import (
	"github.com/containrrr/shoutrrr/cli/cmd/generate"
	"github.com/containrrr/shoutrrr/cli/cmd/send"
	"github.com/containrrr/shoutrrr/cli/cmd/verify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd is the root command for the shoutrrr CLI
var Cmd = &cobra.Command{
	Use:   "shoutrrr",
	Short: "Notification library for gophers and their furry friends",
}

func init() {
	viper.AutomaticEnv()
	Cmd.AddCommand(verify.Cmd)
	Cmd.AddCommand(generate.Cmd)
	Cmd.AddCommand(send.Cmd)
}
