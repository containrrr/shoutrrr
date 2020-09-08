package basic

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/fatih/color"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Generator struct{}

func (g *Generator) Generate(service types.Service, props map[string]string, _ []string) (types.ServiceConfig, error) {
	fmt.Println("Enter the configuration values as prompted")
	fmt.Println()

	var err error

	configType, fields := format.GetConfigFormat(service)
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

	}

	return configPtr.Interface().(types.ServiceConfig), nil
}
