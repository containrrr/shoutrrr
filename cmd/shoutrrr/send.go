package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/util"
)

var verbose bool
var usageExample = "%s send [OPTIONS] <URL> <Message [...]>\n"

func send() action {
	action := action{
		run: func(flags *flag.FlagSet) int {

			if flags.NArg() < 2 {
				return ExitCodeUsage
			}

			url := flags.Arg(0)
			fmt.Printf("URL: %s\n", url)

			message := strings.Join(flags.Args()[1:], " ")
			fmt.Printf("Message: %s\n", message)

			var logger *log.Logger

			if verbose {
				logger = log.New(os.Stderr, "SHOUTRRR ", log.LstdFlags)
			} else {
				logger = util.DiscardLogger
			}

			shoutrrr.SetLogger(logger)
			err := shoutrrr.Send(url, message)

			if err != nil {
				fmt.Printf("error invoking send: %s", err)
				return 1
			}

			return 0
		},
		FlagSet: *flag.NewFlagSet("send", flag.ExitOnError),
		Usage:   usageExample,
	}

	action.FlagSet.BoolVar(&verbose, "verbose", false, "display additional output")

	return action
}
