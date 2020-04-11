package generate

import (
	"fmt"
	"os"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serviceRouter router.ServiceRouter

func hasURLInEnvButNotFlag(cmd *cobra.Command) bool {
	s, _ := cmd.Flags().GetString("url")
	return s == "" && viper.GetViper().GetString("url") != ""
}

// Cmd used to generate and display a config from a notification service URL
var Cmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates and displays a config from a notification service URL.",
	PreRun: func(cmd *cobra.Command, args []string) {
		// WORKAROUND: make cobra count env vars when checking required flags
		if hasURLInEnvButNotFlag(cmd) {
			cmd.Flags().Set("url", viper.GetViper().GetString("url"))
		}
	},
	Run:  Run,
	Args: cobra.NoArgs,
}

func init() {
	serviceRouter = router.ServiceRouter{}
}

// Run the generate command
func Run(cmd *cobra.Command, args []string) {
	URL, _ := cmd.Flags().GetString("url")

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
