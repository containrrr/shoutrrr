package main

import (
	cli "github.com/containrrr/shoutrrr/cli/cmd"
	"github.com/containrrr/shoutrrr/cli/cmd/docs"
	"github.com/containrrr/shoutrrr/cli/cmd/generate"
	"github.com/containrrr/shoutrrr/cli/cmd/send"
	"github.com/containrrr/shoutrrr/cli/cmd/verify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cmd = &cobra.Command{
	Use:   "shoutrrr",
	Short: "Notification library for gophers and their furry friends",
}

func init() {
	viper.AutomaticEnv()
	cmd.AddCommand(verify.Cmd)
	cmd.AddCommand(generate.Cmd)
	cmd.AddCommand(send.Cmd)
	cmd.AddCommand(docs.Cmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(cli.ExUsage)
	}
}
