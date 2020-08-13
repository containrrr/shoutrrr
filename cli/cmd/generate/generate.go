package generate

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var serviceRouter router.ServiceRouter

// Cmd is used to generate a notification service URL from user input
var Cmd = &cobra.Command{
	Use:    "generate",
	Short:  "Generates a notification service URL from user input",
	Run:    Run,
	PreRun: loadArgsFromAltSources,
	Args:   cobra.MaximumNArgs(1),
}

func loadArgsFromAltSources(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		_ = cmd.Flags().Set("service", args[0])
	}
}

func init() {
	serviceRouter = router.ServiceRouter{}
	Cmd.Flags().StringP("service", "s", "", "The notification service to generate a URL for")

}

// Run the generate command
func Run(cmd *cobra.Command, _ []string) {

	var service types.Service
	var err error

	serviceSchema, _ := cmd.Flags().GetString("service")

	if serviceSchema == "" {
		err = errors.New("no service specified")
	} else {
		service, err = serviceRouter.NewService(serviceSchema)
	}

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	if service == nil {
		services := serviceRouter.ListServices()
		serviceList := strings.Join(services, ", ")

		cmd.SetUsageTemplate(cmd.UsageTemplate() + "\nAvailable services: \n  " + serviceList + "\n")

		_ = cmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("URL generation is not yet supported, the %q service configuration takes the following fields:\n", serviceSchema)

	configFormat, maxKeyLen := format.GetConfigFormat(service)
	for key, description := range configFormat {

		// TODO: Read user input

		pad := strings.Repeat(" ", maxKeyLen-len(key))
		_, _ = fmt.Fprintf(color.Output, "%s:  %s%s\n", key, pad, description)
	}
}
