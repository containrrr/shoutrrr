package basic

import (
	"bufio"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/fatih/color"
	"os"
	"reflect"
)

type Generator struct{}

func (g *Generator) Generate(service types.Service, props map[string]string, _ []string) (types.ServiceConfig, error) {
	fmt.Println("Enter the configuration values as prompted")
	fmt.Println()

	var err error

	configType, fields := format.GetServiceConfigFormat(service)
	configPtr := reflect.New(configType)
	config := configPtr.Elem()

	scanner := bufio.NewScanner(os.Stdin)

	for _, field := range fields {

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

			if valueValid, err = format.SetConfigField(config, field, inputValue); !valueValid {
				_, _ = fmt.Fprint(color.Output, "Invalid type ", color.HiYellowString(field.Type.Kind().String()))
				_, _ = fmt.Fprint(color.Output, "for field ", color.HiCyanString(field.Name), "\n\n")
			}

			if err != nil {
				_, _ = fmt.Fprint(color.Output, "Invalid format for field ", color.HiCyanString(field.Name), ": ", err.Error(), "\n\n")
			}

		}

	}

	return configPtr.Interface().(types.ServiceConfig), nil
}
