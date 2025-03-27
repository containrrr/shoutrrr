package generate

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nicholas-fedor/shoutrrr/pkg/generators"
	"github.com/nicholas-fedor/shoutrrr/pkg/router"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/spf13/cobra"
)

// MaximumNArgs defines the maximum number of positional arguments allowed.
const MaximumNArgs = 2

// serviceRouter manages the creation of notification services.
var serviceRouter router.ServiceRouter

// Cmd represents the Cobra command for generating notification service URLs.
var Cmd = &cobra.Command{
	Use:    "generate",
	Short:  "Generates a notification service URL from user input",
	Run:    Run,
	PreRun: loadArgsFromAltSources,
	Args:   cobra.MaximumNArgs(MaximumNArgs),
}

// loadArgsFromAltSources populates command flags from positional arguments if provided.
// It sets the "service" flag from the first argument and "generator" from the second.
func loadArgsFromAltSources(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		_ = cmd.Flags().Set("service", args[0])
	}

	if len(args) > 1 {
		_ = cmd.Flags().Set("generator", args[1])
	}
}

// init initializes the command flags for the generate command.
func init() {
	serviceRouter = router.ServiceRouter{}

	Cmd.Flags().StringP("service", "s", "", "Notification service to generate a URL for (e.g., discord, smtp)")
	Cmd.Flags().StringP("generator", "g", "basic", "Generator to use (e.g., basic, or service-specific)")
	Cmd.Flags().StringArrayP("property", "p", []string{}, "Configuration property in key=value format (e.g., token=abc123)")
	Cmd.Flags().BoolP("show-sensitive", "x", false, "Show sensitive data in the generated URL (default: masked)")
}

// maskSensitiveURL masks sensitive parts of a Shoutrrr URL based on the service schema.
// It redacts credentials like tokens, passwords, or user keys, tailoring the masking to specific services.
func maskSensitiveURL(serviceSchema, urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr // Return original URL if parsing fails
	}

	switch serviceSchema {
	case "discord", "slack", "teams":
		maskUser(u, "REDACTED")
	case "smtp":
		maskSMTPUser(u)
	case "pushover":
		maskPushoverQuery(u)
	case "gotify":
		maskGotifyQuery(u)
	default:
		maskGeneric(u)
	}

	return u.String()
}

// maskUser redacts the username in a URL, replacing it with a placeholder.
func maskUser(u *url.URL, placeholder string) {
	if u.User != nil {
		u.User = url.User(placeholder)
	}
}

// maskSMTPUser redacts the password in an SMTP URL, preserving the username.
func maskSMTPUser(u *url.URL) {
	if u.User != nil {
		u.User = url.UserPassword(u.User.Username(), "REDACTED")
	}
}

// maskPushoverQuery redacts token and user query parameters in a Pushover URL.
func maskPushoverQuery(u *url.URL) {
	q := u.Query()
	if q.Get("token") != "" {
		q.Set("token", "REDACTED")
	}

	if q.Get("user") != "" {
		q.Set("user", "REDACTED")
	}

	u.RawQuery = q.Encode()
}

// maskGotifyQuery redacts the token query parameter in a Gotify URL.
func maskGotifyQuery(u *url.URL) {
	q := u.Query()
	if q.Get("token") != "" {
		q.Set("token", "REDACTED")
	}

	u.RawQuery = q.Encode()
}

// maskGeneric redacts userinfo and all query parameters for unrecognized services.
func maskGeneric(u *url.URL) {
	maskUser(u, "REDACTED")

	q := u.Query()
	for key := range q {
		q.Set(key, "REDACTED")
	}

	u.RawQuery = q.Encode()
}

// Run executes the generate command, producing a notification service URL.
// It validates inputs, generates the URL using the specified service and generator,
// and outputs it, masking sensitive data unless --show-sensitive is set.
func Run(cmd *cobra.Command, _ []string) {
	var service types.Service

	var err error

	serviceSchema, _ := cmd.Flags().GetString("service")
	generatorName, _ := cmd.Flags().GetString("generator")
	propertyFlags, _ := cmd.Flags().GetStringArray("property")
	showSensitive, _ := cmd.Flags().GetBool("show-sensitive")

	// Parse properties into a key-value map.
	props := make(map[string]string, len(propertyFlags))

	for _, prop := range propertyFlags {
		parts := strings.Split(prop, "=")
		if len(parts) != MaximumNArgs {
			_, _ = fmt.Fprintln(color.Output, "Invalid property key/value pair:", color.HiYellowString(prop))

			continue
		}

		props[parts[0]] = parts[1]
	}

	if len(propertyFlags) > 0 {
		fmt.Println() // Add spacing after property warnings
	}

	// Validate and create the service.
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
		cmd.SetUsageTemplate(cmd.UsageTemplate() + "\nAvailable services:\n  " + serviceList + "\n")
		_ = cmd.Usage()

		os.Exit(1)
	}

	// Determine the generator to use.
	var generator types.Generator

	generatorFlag := cmd.Flags().Lookup("generator")
	if !generatorFlag.Changed {
		// Use the service-specific default generator if available and no explicit generator is set.
		generator, _ = generators.NewGenerator(serviceSchema)
	}

	if generator != nil {
		generatorName = serviceSchema
	} else {
		var genErr error

		generator, genErr = generators.NewGenerator(generatorName)
		if genErr != nil {
			fmt.Printf("Error: %s\n", genErr)
		}
	}

	if generator == nil {
		generatorList := strings.Join(generators.ListGenerators(), ", ")
		cmd.SetUsageTemplate(cmd.UsageTemplate() + "\nAvailable generators:\n  " + generatorList + "\n")
		_ = cmd.Usage()

		os.Exit(1)
	}

	// Generate and display the URL.
	_, _ = fmt.Fprint(color.Output, "Generating URL for ", color.HiCyanString(serviceSchema))
	_, _ = fmt.Fprintln(color.Output, " using ", color.HiMagentaString(generatorName), " generator")

	serviceConfig, err := generator.Generate(service, props, cmd.Flags().Args())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println()

	maskedURL := maskSensitiveURL(serviceSchema, serviceConfig.GetURL().String())

	if showSensitive {
		fmt.Println("URL:", serviceConfig.GetURL().String())
	} else {
		fmt.Println("URL:", maskedURL)
	}
}
