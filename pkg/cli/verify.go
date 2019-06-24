package main

import (
	"flag"
	"fmt"
	"github.com/containrrr/shoutrrr"
)


func verify() action {
	return action{
		run: func(flags *flag.FlagSet) int {

			url := flags.Arg(0)
			if url == "" {
				return ExitCodeUsage
			}


			if err := shoutrrr.Verify(url); err != nil {
				fmt.Printf("error verifying URL: %s", err)
				return 1
			}


			return 0
		},
		FlagSet: *flag.NewFlagSet("verify", flag.ExitOnError),
		Usage: "%s send [OPTIONS] <URL> <Message [...]>\n",
	}
}