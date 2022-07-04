package main

import (
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/containrrr/shoutrrr/cmd/shoutrrr-gen/confwriter"
	"github.com/containrrr/shoutrrr/internal/logging"
	"github.com/containrrr/shoutrrr/pkg/cli"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	cmd = &cobra.Command{
		Use:  "shoutrrr-gen",
		RunE: run,
		Args: cobra.MinimumNArgs(1),
	}
	opts struct {
		lang   string
		source string
	}
)

func init() {
	cmd.Flags().StringVarP(&opts.lang, "lang", "l", "", "Output language")
	cmd.Flags().StringVarP(&opts.source, "source", "s", "", "Source file")
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(cli.ExUsage)
	}
}

func run(c *cobra.Command, args []string) error {

	logging.InitLogging(true)
	logger := logging.GetLogger("shoutrrr-gen")

	// logConfig := zap.NewDevelopmentConfig()
	// logConfig.EncoderConfig.EncodeCaller
	// logger, err := logConfig.Build()
	// if err != nil {
	// 	panic(err)
	// }
	// zap.ReplaceGlobals(logger)
	// // zap.Development()
	// log := zap.S()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	gofile, _ := os.LookupEnv("GOFILE")
	gopackage, _ := os.LookupEnv("GOPACKAGE")

	specFile := args[0]

	logger.WithFields(log.Fields{
		"context": "shoutrrr-gen",
		"CWD":     cwd,
		"Source":  opts.source,
		"Lang":    opts.lang,
		"Config":  gofile,
		"Service": gopackage,
	}).Debug("Started")

	specBytes, err := os.ReadFile(specFile)
	if err != nil {
		return err
	}

	spec := format.ConfigSpec{}
	if err := yaml.Unmarshal(specBytes, &spec); err != nil {
		return err
	}

	genfile := strings.Replace(gofile, "config.go", "config.gen.go", 1)
	src, err := os.Create(genfile)
	if err != nil {
		return err
	}
	defer src.Close()

	cw := confwriter.New(&spec, src)

	cw.WriteHeader(gopackage, os.Args)
	if err := cw.WriteProps(); err != nil {
		return err
	}
	cw.WriteGetURL()
	cw.WriteSetURL()
	cw.WriteEnums()
	cw.WriteUpdate()
	cw.WriteHelpers()

	return nil
}
