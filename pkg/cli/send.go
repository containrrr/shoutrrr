package main

import (
	"flag"
	"fmt"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/plugin"
	"log"
	"os"
	"strings"
)

var verbose bool

func Send() Action {


	action := Action {
		run: func(flags *flag.FlagSet) int {

			if flags.NArg() < 2 {
				fmt.Println(flags.NArg())
				fmt.Println(flags.Args())
				return ExitCodeUsage
			}

			fmt.Printf("Args: %d\n", flags.NArg())

			// Arg #0 is always the action verb
			url := flags.Arg(0)

			fmt.Printf("Url: %s\n", url)

			message := strings.Join(flags.Args()[1:], " ")

			fmt.Printf("Message: %s\n", message)


			logger := plugin.DiscardLogger
			if verbose {
				logger = log.New(os.Stderr, "SHOUTRRR ", log.LstdFlags)
			}

			opts := plugin.PluginOpts{
				Verbose: verbose,
				Logger: logger,
			}

			shoutrrr.Send(url, message, opts)

			return 1
		},
		FlagSet: *flag.NewFlagSet("send", flag.ExitOnError),
		Usage: "%s send [OPTIONS] <URL> <Message [...]>\n",
	}

	action.FlagSet.BoolVar(&verbose, "verbose", false, "display additional output")

	return action
}