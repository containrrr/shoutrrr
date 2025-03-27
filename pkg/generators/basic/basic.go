package basic

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/nicholas-fedor/shoutrrr/pkg/format"
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
)

// Generator is the Basic Generator implementation.
type Generator struct{}

// Generate generates a service URL from a set of user questions/answers.
func (g *Generator) Generate(service types.Service, props map[string]string, _ []string) (types.ServiceConfig, error) {
	fmt.Println("Enter the configuration values as prompted")
	fmt.Println()

	var err error

	configNode := format.GetServiceConfigFormat(service)
	configPtr := reflect.New(configNode.Type)
	config := configPtr.Elem()

	scanner := bufio.NewScanner(os.Stdin)

	for _, item := range configNode.Items {
		field := item.Field()

		var inputValue string

		valueValid := false

		for !valueValid {
			err = nil

			propKey := strings.ToLower(field.Name)
			if propValue, ok := props[propKey]; ok && len(propValue) > 0 {
				inputValue = propValue
				_, _ = fmt.Fprint(color.Output, "Using property ", color.HiCyanString(propValue), " for ", color.HiMagentaString(field.Name), " field\n")
				// Clear the property value to skip it next iteration in case of errors
				props[propKey] = ""
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
						return nil, err
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

			if valueValid, err = format.SetConfigField(config, *field, inputValue); !valueValid && err == nil {
				_, _ = fmt.Fprint(color.Output, "Invalid type ", color.HiYellowString(field.Type.Kind().String()))
				_, _ = fmt.Fprint(color.Output, " for field ", color.HiCyanString(field.Name), "\n\n")
			}

			if err != nil {
				_, _ = fmt.Fprint(color.Output, "Invalid format for field ", color.HiCyanString(field.Name), ": ", err.Error(), "\n\n")
			}
		}
	}

	return configPtr.Interface().(types.ServiceConfig), nil
}
