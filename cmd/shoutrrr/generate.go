package main

import (
	"flag"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
)

func generate() action {

	action := action{
		run: func(flags *flag.FlagSet) int {

			if flags.NArg() < 1 {
				return ExitCodeUsage
			}

			serviceSchema := flags.Arg(0)

			fmt.Printf("Service: %s\n", serviceSchema)

			serviceRouter := router.ServiceRouter{}

			service, err := serviceRouter.Locate(serviceSchema)
			if err != nil {
				fmt.Printf("invalid service schema '%s'\n", serviceSchema)
				return 2
			}

			configFormat, _ := format.GetConfigMap(service) // TODO: GetConfigFormat
			for key, format := range configFormat {
				fmt.Printf("%s: %s", key, format)
			}

			return 1
		},
		FlagSet: *flag.NewFlagSet("generate", flag.ExitOnError),
		Usage:   "%s generate [OPTIONS] <service>\n",
	}

	// action.FlagSet.BoolVar(&verbose, "verbose", false, "display additional output")

	return action
}
