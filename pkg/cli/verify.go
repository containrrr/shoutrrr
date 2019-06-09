package main

import (
	"flag"
	"fmt"
	"github.com/containrrr/shoutrrr"
)


func verify() action {

	action := action{
		run: func(flags *flag.FlagSet) int {

			if flags.NArg() < 1 {
				return ExitCodeUsage
			}

			//fmt.Printf("Args: %d\n", flags.NArg())

			// Arg #0 is always the action verb
			url := flags.Arg(0)

			//fmt.Printf("Url: %s\n", url)

			//logger := DiscardLogger
			//if verbose {
			//	logger = log.New(os.Stderr, "SHOUTRRR ", log.LstdFlags)
			//}
			//
			//opts := PluginOpts {
			//	Verbose: verbose,
			//	Logger: logger,
			//}

			if err := shoutrrr.Verify(url); err != nil {
				fmt.Printf("error verifying URL: %s", err)
				return 1
			}


			return 0
		},
		FlagSet: *flag.NewFlagSet("verify", flag.ExitOnError),
		Usage: "%s send [OPTIONS] <URL> <Message [...]>\n",
	}

	// action.FlagSet.BoolVar(&verbose, "verbose", false, "display additional output")

	return action
}