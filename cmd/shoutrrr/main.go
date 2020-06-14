package main

import (
	"github.com/containrrr/shoutrrr/cmd/generate"
	"github.com/containrrr/shoutrrr/cmd/send"
	"github.com/containrrr/shoutrrr/cmd/verify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "shoutrrr",
	Short: "Notification library for gophers and their furry friends",
}

func init() {
	viper.AutomaticEnv()
	rootCmd.AddCommand(verify.Cmd)
	rootCmd.AddCommand(generate.Cmd)
	rootCmd.AddCommand(send.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
