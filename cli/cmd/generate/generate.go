package generate

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"reflect"
	"strconv"
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

	_, _ = fmt.Fprint(color.Output, "Generating URL for ", color.HiCyanString(serviceSchema))
	_, _ = fmt.Fprintln(color.Output, " using", color.HiMagentaString(generatorName), "generator")
	fmt.Println("Enter the configuration values as prompted")
	fmt.Println()

	// fmt.Printf("URL generation is not yet supported, the %q service configuration takes the following fields:\n", serviceSchema)

	configType, fields := format.GetConfigFormat(service)
	configPtr := reflect.New(configType)
	config := configPtr.Elem()

	scanner := bufio.NewScanner(os.Stdin)

	for _, field := range fields {

		// pad := strings.Repeat(" ", maxKeyLen-len(key))

		var inputValue string
		valueValid := false

		for !valueValid {
			err = nil

			if propValue, ok := props[field.Name]; ok && len(propValue) > 0 {
				inputValue = propValue
				_, _ = fmt.Fprint(color.Output, "Using property ", color.HiCyanString(propValue), " for ", color.HiMagentaString(field.Name), " field\n")
				// Clear the property value to skip it next iteration in case of errors
				props[field.Name] = ""
			} else {
				if len(field.DefaultValue) > 0 {
					_, _ = fmt.Fprint(color.Output, color.HiWhiteString(field.Name), "[", field.DefaultValue, "]: ")
				} else {
					_, _ = fmt.Fprint(color.Output, color.HiWhiteString(field.Name), ": ")
				}

				if scanner.Scan() {
					inputValue = scanner.Text()
				} else {
					if err := scanner.Err(); err != nil {
						log.Fatal(err)
					}
				}
			}

			// Handle empty lines
			if len(inputValue) == 0 {
				if field.Required {
					_, _ = fmt.Fprint(color.Output, "Field ", color.HiCyanString(field.Name), " is required!\n\n")
					continue
				} else {
					// Use the default value instead of the user input
					inputValue = field.DefaultValue
					if len(inputValue) == 0 {
						// No default value is specified, leave it uninitialized
						valueValid = true
						continue
					}
				}
			}

			configField := config.FieldByName(field.Name)
			fieldKind := field.Type.Kind()

			if fieldKind == reflect.String {
				configField.SetString(inputValue)
				valueValid = true
			} else if field.EnumFormatter != nil {
				value := field.EnumFormatter.Parse(inputValue)
				if value == format.EnumInvalid {
					enumNames := strings.Join(field.EnumFormatter.Names(), ", ")
					err = fmt.Errorf("not a one of %v", enumNames)
				} else {
					configField.SetInt(int64(value))
					valueValid = true
				}
			} else if fieldKind >= reflect.Uint && fieldKind <= reflect.Uint64 {
				var value uint64
				value, err = strconv.ParseUint(inputValue, 10, field.Type.Bits())
				if err == nil {
					configField.SetUint(value)
					valueValid = true
				}
			} else if fieldKind >= reflect.Int && fieldKind <= reflect.Int64 {
				var value int64
				value, err = strconv.ParseInt(inputValue, 10, field.Type.Bits())
				if err == nil {
					configField.SetInt(value)
					valueValid = true
				}
			} else if fieldKind == reflect.Bool {
				if value, ok := format.ParseBool(inputValue, false); !ok {
					err = errors.New("accepted values are 1, true, yes or 0, false, no")
				} else {
					configField.SetBool(value)
					valueValid = true
				}
			} else if fieldKind >= reflect.Slice {
				elemKind := field.Type.Elem().Kind()
				if elemKind != reflect.String {
					err = errors.New("field format is not supported")
				} else {
					values := strings.Split(inputValue, ",")
					configField.Set(reflect.ValueOf(values))
					valueValid = true
				}
			} else {
				_, _ = fmt.Fprint(color.Output, "Invalid type ", color.HiYellowString(fieldKind.String()))
				_, _ = fmt.Fprint(color.Output, "for field ", color.HiCyanString(field.Name), "\n\n")
			}

			if err != nil {
				_, _ = fmt.Fprint(color.Output, "Invalid format for field ", color.HiCyanString(field.Name), ": ", err.Error(), "\n\n")
			}

		}

		//
		// _, _ = fmt.Fprintf(color.Output, "%s:  %s%s\n", key, pad, )
	}

	serviceConfig := configPtr.Interface().(types.ServiceConfig)

	fmt.Println()
	fmt.Println("URL:", serviceConfig.GetURL().String())

}
