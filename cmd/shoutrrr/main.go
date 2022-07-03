package main

import (
	"os"

	"github.com/containrrr/shoutrrr/cmd/shoutrrr/docs"
	"github.com/containrrr/shoutrrr/cmd/shoutrrr/generate"
	"github.com/containrrr/shoutrrr/cmd/shoutrrr/send"
	"github.com/containrrr/shoutrrr/cmd/shoutrrr/verify"
	"github.com/containrrr/shoutrrr/internal/meta"
	"github.com/containrrr/shoutrrr/pkg/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:     "shoutrrr",
	Version: meta.Version,
	Short:   "Shoutrrr CLI",
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
