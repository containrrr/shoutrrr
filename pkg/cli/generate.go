package main

import (
	"flag"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
)

func Generate() Action {

	action := Action {
		run: func(flags *flag.FlagSet) int {

			if flags.NArg() < 1 {
				return ExitCodeUsage
			}

			serviceSchema := flags.Arg(0)

			fmt.Printf("Service: %s\n", serviceSchema)

			//logger := DiscardLogger
			//if verbose {
			//	logger = log.New(os.Stderr, "SHOUTRRR ", log.LstdFlags)
			//}
			//
			//opts := PluginOpts {
			//	Verbose: verbose,
			//	Logger: logger,
			//}

			serviceRouter := router.ServiceRouter{}

			service, err := serviceRouter.Locate(serviceSchema)
			if err != nil {
				fmt.Printf("invalid service schema '%s'\n", serviceSchema)
				return 2
			}

			config := service.GetConfig()
			configFormat := format.GetConfigMap(config) // TODO: GetConfigFormat
			for key, format := range configFormat {
				fmt.Printf("%s: %s", key, format)
			}

			return 1
		},
		FlagSet: *flag.NewFlagSet("generate", flag.ExitOnError),
		Usage: "%s send [OPTIONS] <URL> <Message [...]>\n",
	}

	// action.FlagSet.BoolVar(&verbose, "verbose", false, "display additional output")

	return action
}