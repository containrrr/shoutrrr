package generate

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/generators"
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
	Args:   cobra.MaximumNArgs(2),
}

func loadArgsFromAltSources(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		_ = cmd.Flags().Set("service", args[0])
	}
	if len(args) > 1 {
		_ = cmd.Flags().Set("generator", args[1])
	}
}

func init() {
	serviceRouter = router.ServiceRouter{}
	Cmd.Flags().StringP("service", "s", "", "The notification service to generate a URL for")

	Cmd.Flags().StringP("generator", "g", "basic", "The generator to use")

	Cmd.Flags().StringArrayP("property", "p", []string{}, "Configuration property in key=value format")
}

// Run the generate command
func Run(cmd *cobra.Command, _ []string) {

	var service types.Service
	var err error

	serviceSchema, _ := cmd.Flags().GetString("service")
	generatorName, _ := cmd.Flags().GetString("generator")
	propertyFlags, _ := cmd.Flags().GetStringArray("property")

	props := make(map[string]string, len(propertyFlags))
	for _, prop := range propertyFlags {
		parts := strings.Split(prop, "=")
		if len(parts) != 2 {
			_, _ = fmt.Fprintln(color.Output, "Invalid property key/value pair:", color.HiYellowString(prop))
			continue
		}
		props[parts[0]] = parts[1]
	}

	if len(propertyFlags) > 0 {
		fmt.Println()
	}

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

	var generator types.Generator

	var generatorFlag = cmd.Flags().Lookup("generator")

	if !generatorFlag.Changed {
		// try to use the service default generator if one exists
		generator, _ = generators.NewGenerator(serviceSchema)
	}

	if generator != nil {
		generatorName = serviceSchema
	} else {
		generator, err = generators.NewGenerator(generatorName)
	}

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	if generator == nil {
		generatorList := strings.Join(generators.ListGenerators(), ", ")

		cmd.SetUsageTemplate(cmd.UsageTemplate() + "\nAvailable generators: \n  " + generatorList + "\n")

		_ = cmd.Usage()
		os.Exit(1)
	}

	_, _ = fmt.Fprint(color.Output, "Generating URL for ", color.HiCyanString(serviceSchema))
	_, _ = fmt.Fprintln(color.Output, " using", color.HiMagentaString(generatorName), "generator")

	serviceConfig, err := generator.Generate(service, props, cmd.Flags().Args())

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("URL:", serviceConfig.GetURL().String())

}
