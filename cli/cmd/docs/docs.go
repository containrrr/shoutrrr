package docs

import (
	"fmt"
	"github.com/containrrr/shoutrrr/internal/meta"
	"github.com/mattn/go-isatty"
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
		if showVersion, _ := cmd.Flags().GetBool("version"); showVersion {
			return nil
		}
		serviceList := strings.Join(services, ", ")
		cmd.SetUsageTemplate(cmd.UsageTemplate() + "\nAvailable services: \n  " + serviceList + "\n")
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	ValidArgs: services,
}

func init() {
	Cmd.Flags().StringP("format", "f", "console", "Output format")
	Cmd.Flags().BoolP("version", "V", false, "Show docs version")
}

// Run the docs command
func Run(cmd *cobra.Command, args []string) {
	format, _ := cmd.Flags().GetString("format")

	if showVersion, _ := cmd.Flags().GetBool("version"); showVersion {
		_, _ = os.Stdout.WriteString(meta.DocsVersion)
		if isatty.IsTerminal(os.Stdout.Fd()) {
			// write a newline if the output is not being redirected
			_, _ = os.Stdout.WriteString("\n")
		}
		os.Exit(cli.ExSuccess)
	}

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
		renderer = f.MarkdownTreeRenderer{
			HeaderPrefix:      "### ",
			PropsDescription:  "Props can be either supplied using the params argument, or through the URL using  \n`?key=value&key=value` etc.\n",
			PropsEmptyMessage: "*The services does not support any query/param props*",
		}
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
