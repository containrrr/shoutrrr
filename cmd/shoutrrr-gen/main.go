package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"

	"github.com/containrrr/shoutrrr/cmd/shoutrrr-gen/confwriter"
	"github.com/containrrr/shoutrrr/internal/logging"
	"github.com/containrrr/shoutrrr/internal/util"
	"github.com/containrrr/shoutrrr/pkg/cli"
	"github.com/containrrr/shoutrrr/pkg/conf"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	cmd = &cobra.Command{
		Use:  "shoutrrr-gen",
		RunE: run,
		// Args: cobra.MinimumNArgs(1),
	}
	opts struct {
		lang   string
		source string
		pkg    string
		root   string
		debug  bool
	}
)

func init() {
	cmd.Flags().StringVarP(&opts.lang, "lang", "l", "", "Output language")
	cmd.Flags().StringVarP(&opts.source, "source", "s", "", "Source file")
	cmd.Flags().StringVarP(&opts.root, "root", "r", "", "Repository root")
	cmd.Flags().BoolVarP(&opts.debug, "debug", "d", false, "Debug")
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(cli.ExUsage)
	}
}

func run(c *cobra.Command, args []string) error {

	logging.InitLogging(opts.debug)
	logger := logging.GetLogger("shoutrrr-gen")

	if val, found := os.LookupEnv("GOFILE"); found {
		opts.source = val
	}
	if val, found := os.LookupEnv("GOPACKAGE"); found {
		opts.pkg = val
	}

	if len(args) > 0 {
		opts.pkg = args[0]
	}
	if opts.pkg == "" {
		return fmt.Errorf("no package specified")
	}

	logger.WithField("service", opts.pkg).Info("Generating service config")

	if opts.root == "" {
		root, err := util.FindGitRootFromCwd()
		if err != nil {
			return err
		}
		logger.WithField("root", root).Debug("Found git root")
		opts.root = root
	}

	serviceRoot := fmt.Sprintf("%v/pkg/services/%v", opts.root, opts.pkg)

	if opts.source == "" {
		logger.Infof("Running outside go generate")
		opts.source = fmt.Sprintf("%v/%v_config.go", serviceRoot, opts.pkg)
	}

	specFile := fmt.Sprintf("%v/spec/%v.yml", opts.root, opts.pkg)
	genFile := fmt.Sprintf("%v/%v_config.gen.go", serviceRoot, opts.pkg)

	// logger.WithFields(log.Fields{
	// 	"context": "shoutrrr-gen",
	// 	"Root":    opts.root,
	// 	"Service": opts.pkg,
	// 	"Source":  opts.source,
	// 	// "Lang":    opts.lang,
	// 	"Spec":   specFile,
	// 	"Target": genFile,
	// }).Debug("")

	logger.WithField("file", specFile).Debug("Loading spec")
	specBytes, err := os.ReadFile(specFile)
	if err != nil {
		return err
	}

	spec := conf.Spec{}
	if err := yaml.Unmarshal(specBytes, &spec); err != nil {
		return err
	}

	cw := confwriter.New(&spec)

	logger.WithField("file", genFile).Debug("Creating output")
	src, err := os.Create(genFile)
	if err != nil {
		return err
	}
	defer src.Close()

	buf := bytes.Buffer{}

	if err := cw.WriteConfig(&buf, os.Args); err != nil {
		return err
	}

	logger.WithField("bytes", buf.Len()).Info("Formatting result")

	if formatted, err := format.Source(buf.Bytes()); err != nil {
		return err
	} else {
		wrote, err := src.Write(formatted)
		if err != nil {
			return err
		} else {
			logger.WithField("bytes", wrote).Info("Wrote result to output")
		}
	}

	return nil
}
