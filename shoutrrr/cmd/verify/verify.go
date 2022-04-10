package verify

import (
	"fmt"
	"github.com/containrrr/shoutrrr/internal/util"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/fatih/color"
	"os"

	"github.com/spf13/cobra"
)

// Cmd verifies the validity of a service url
var Cmd = &cobra.Command{
	Use:    "verify",
	Short:  "Verify the validity of a notification service URL",
	PreRun: util.LoadFlagsFromAltSources,
	Run:    Run,
	Args:   cobra.MaximumNArgs(1),
}

var sr router.ServiceRouter

func init() {
	Cmd.Flags().StringP("url", "u", "", "The notification url")
	_ = Cmd.MarkFlagRequired("url")
}

// Run the verify command
func Run(cmd *cobra.Command, _ []string) {
	URL, _ := cmd.Flags().GetString("url")
	sr = router.ServiceRouter{}

	service, err := sr.Locate(URL)

	if err != nil {
		fmt.Printf("error verifying URL: %s\n", err)
		os.Exit(1)
	}

	config := format.GetServiceConfig(service)
	configNode := format.GetConfigFormat(config)

	_, _ = fmt.Fprintf(color.Output, format.ColorFormatTree(configNode, true))
}
