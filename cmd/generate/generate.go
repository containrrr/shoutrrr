package generate

import (
	"fmt"
	"os"

	"github.com/containrrr/shoutrrr/internal/util"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/spf13/cobra"
)

var serviceRouter router.ServiceRouter

// Cmd used to generate and display a config from a notification service URL
var Cmd = &cobra.Command{
	Use:    "generate",
	Short:  "Generates and displays a config from a notification service URL.",
	PreRun: util.MoveEnvVarToFlag,
	Run:    Run,
	Args:   cobra.MaximumNArgs(1),
}

func init() {
	serviceRouter = router.ServiceRouter{}

}

// Run the generate command
func Run(cmd *cobra.Command, args []string) {
	URL := args[0]

	if _, err := serviceRouter.Locate(URL); err != nil {
		fmt.Printf("invalid service schema '%s', %s", URL, err)
		os.Exit(1)
	}
	fmt.Printf("Service: %s\n", URL)

	serviceSchema := URL
	service, _ := serviceRouter.Locate(serviceSchema)

	configFormat, _ := format.GetConfigMap(service) // TODO: GetConfigFormat
	for key, format := range configFormat {
		fmt.Printf("%s: %s\n", key, format)
	}
}
