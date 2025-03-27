package docs

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/nicholas-fedor/shoutrrr/pkg/router"
	"github.com/spf13/cobra"

	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/shoutrrr/cmd"
)

var (
	serviceRouter router.ServiceRouter
	services      = serviceRouter.ListServices()
)

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

func Run(cmd *cobra.Command, args []string) {
	format, _ := cmd.Flags().GetString("format")
	res := printDocs(format, args)

	if res.ExitCode != 0 {
		_, _ = fmt.Fprintf(os.Stderr, "%s", res.Message)
	}

	os.Exit(res.ExitCode)
}

func printDocs(docFormat string, services []string) cmd.Result {
	var renderer format.TreeRenderer

	switch docFormat {
	case "console":
		renderer = format.ConsoleTreeRenderer{WithValues: false}
	case "markdown":
		renderer = format.MarkdownTreeRenderer{
			HeaderPrefix:      "### ",
			PropsDescription:  "Props can be either supplied using the params argument, or through the URL using  \n`?key=value&key=value` etc.\n",
			PropsEmptyMessage: "*The services does not support any query/param props*",
		}
	default:
		return cmd.InvalidUsage("invalid format")
	}

	logger := log.New(os.Stderr, "", 0) // Concrete logger implementing types.StdLogger

	for _, scheme := range services {
		service, err := serviceRouter.NewService(scheme)
		if err != nil {
			return cmd.InvalidUsage("failed to init service: " + err.Error())
		}
		// Initialize the service to populate Config
		dummyURL, _ := url.Parse(fmt.Sprintf("%s://dummy@dummy.com", scheme))
		if err := service.Initialize(dummyURL, logger); err != nil {
			return cmd.InvalidUsage(fmt.Sprintf("failed to initialize service %q: %v", scheme, err))
		}

		config := format.GetServiceConfig(service)
		configNode := format.GetConfigFormat(config)
		fmt.Println(renderer.RenderTree(configNode, scheme))
	}

	return cmd.Success
}
