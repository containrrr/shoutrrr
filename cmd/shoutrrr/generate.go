package main

import (
	"flag"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/services/smtp"
	"log"
)

func generate() action {

	action := action{
		run: func(flags *flag.FlagSet) int {

			if flags.NArg() < 1 {
				return ExitCodeUsage
			}

			serviceSchema := flags.Arg(0)

			fmt.Printf("Service: %s\n", serviceSchema)

			if serviceSchema == "smtp" {

				guide := flags.Arg(1)
				var url string
				var err error

				credFile := flags.Arg(2)

				if guide == "oauth2" {
					if len(credFile) > 0 {
						url, err = smtp.OAuth2GeneratorFile(credFile)
					} else {
						url, err = smtp.OAuth2Generator()
					}
				} else if guide == "gmail" {
					url, err = smtp.OAuth2GeneratorGmail(credFile)
				} else {
					err = fmt.Errorf("unknown guide %q", guide)
				}

				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("\nService URL:\n%q", url)

				return 0
			}

			serviceRouter := router.ServiceRouter{}

			service, err := serviceRouter.Locate(serviceSchema)
			if err != nil {
				fmt.Printf("invalid service schema '%s': %s\n", serviceSchema, err.Error())
				return 2
			}

			configFormat, _ := format.GetConfigMap(service) // TODO: GetConfigFormat
			for key, cf := range configFormat {
				fmt.Printf("%s: %s", key, cf)
			}

			return 1
		},
		FlagSet: *flag.NewFlagSet("generate", flag.ExitOnError),
		Usage:   "%s generate [OPTIONS] <service>\n",
	}

	// action.FlagSet.BoolVar(&verbose, "verbose", false, "display additional output")

	return action
}
