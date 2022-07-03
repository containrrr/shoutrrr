package main

import (
	"fmt"
	"os"
	"path"

	"github.com/containrrr/shoutrrr/pkg/cli"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/spf13/cobra"
)

// cmd is used to migrate service configs between Shoutrrr versions
var cmd = &cobra.Command{
	Use:   "shoutrrr-mkspec",
	Short: "",
	Run:   Run,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	flags := cmd.Flags()

	// Cmd.Flags().StringP("service", "s", "", "The notification service to generate a URL for")

	// Cmd.Flags().StringP("generator", "g", "basic", "The generator to use")

	flags.StringP("output", "o", "spec", "Output dir")
	flags.BoolP("file", "f", false, "Output file")
	// flags.
}

func Run(cmd *cobra.Command, _ []string) {
	services := cmd.Flags().Args()
	serviceRouter := router.ServiceRouter{}

	if services[0] == "all" {
		services = serviceRouter.ListServices()
	}

	var err error
	target := os.Stdout

	for _, serviceSchema := range services {
		logf("Migrating service %q...\n", serviceSchema)
		if f, _ := cmd.Flags().GetBool("file"); f {
			output, _ := cmd.Flags().GetString("output")
			fileName := path.Join(output, serviceSchema+".yml")
			target, err = os.Create(fileName)
			if err != nil {
				exit(fmt.Sprintf("Error opening output file %q: %v\n", fileName, err))
				return
			}
		}
		service, err := serviceRouter.NewService(serviceSchema)
		if err != nil {
			exit(fmt.Sprintf("Error resolving service: %v\n", err))
			return
		}
		if err = Export(service, serviceSchema, target); err != nil {
			exit(fmt.Sprintf("Error exporting service: %v\n", err))
			return
		}
	}
}

func logf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func exit(message string) {
	logf("%v\nExiting!", message)
	os.Exit(1)
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(cli.ExUsage)
	}
}
