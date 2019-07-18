package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/format"
)


func verify() action {
	return action{
		run: func(flags *flag.FlagSet) int {

			url := flags.Arg(0)
			if url == "" {
				return ExitCodeUsage
			}

			service, err := shoutrrr.CreateSender(url)
			if err != nil {
				fmt.Printf("error verifying URL: %s", err)
				return 1
			}

			configMap, maxKeyLen := format.GetConfigMap(service)
			for key, value := range configMap {
				pad := strings.Repeat(" ", maxKeyLen -len(key))
				fmt.Printf("%s%s: %s\n", pad, key, value)
			}

			return 0
		},
		FlagSet: *flag.NewFlagSet("verify", flag.ExitOnError),
		Usage:   "%s send [OPTIONS] <URL> <Message [...]>\n",
	}
}