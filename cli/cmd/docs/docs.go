package docs

import (
	"fmt"
	"os"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/spf13/cobra"

	cli "github.com/containrrr/shoutrrr/cli/cmd"
	f "github.com/containrrr/shoutrrr/pkg/format"
)

var serviceRouter router.ServiceRouter
var services = serviceRouter.ListServices()

// Cmd prints documentation for services
var Cmd = &cobra.Command{
	Use:   "docs",
	Short: "Print documentation for services",
	Run:   Run,
	Args: func(cmd *cobra.Command, args []string) error {
		serviceList := strings.Join(services, ", ")
		cmd.SetUsageTemplate(cmd.UsageTemplate() + "\nAvailable services: \n  " + serviceList + "\n")
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	ValidArgs: services,
}

func init() {
	Cmd.Flags().StringP("format", "f", "console", "Output format")
}

// Run the docs command
func Run(cmd *cobra.Command, args []string) {
	format, _ := cmd.Flags().GetString("format")

	res := printDocs(format, args)
	if res.ExitCode != 0 {
		_, _ = fmt.Fprintf(os.Stderr, res.Message)
	}
	os.Exit(res.ExitCode)
}

func printDocs(format string, services []string) cli.Result {
	var renderer f.TreeRenderer

	switch format {
	case "console":
		renderer = f.ConsoleTreeRenderer{WithValues: false}
	case "markdown":
		renderer = f.MarkdownTreeRenderer{}
	default:
		return cli.InvalidUsage("invalid format")
	}

	for _, scheme := range services {
		service, err := serviceRouter.NewService(scheme)
		if err != nil {
			return cli.InvalidUsage("failed to init service: " + err.Error())
		}
		config := f.GetServiceConfig(service)
		configNode := f.GetConfigFormat(config)
		fmt.Println(renderer.RenderTree(configNode, scheme))
	}

	return cli.Success
}
