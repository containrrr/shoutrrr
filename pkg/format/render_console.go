package format

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/nicholas-fedor/shoutrrr/pkg/util"
)

// Constants for console rendering.
const (
	DescriptionColumnWidth = 60 // Width of the description column in console output
	ItemSeparatorLength    = 2  // Length of the ", " separator between container items
	DefaultValueOffset     = 16 // Minimum offset before description when no values are shown
	ValueOffset            = 30 // Offset before description when values are shown
	ContainerBracketLength = 4  // Length of container delimiters (e.g., "{ }" or "[ ]")
	KeySeparatorLength     = 2  // Length of the ": " separator after a key in containers
)

// ConsoleTreeRenderer renders a ContainerNode tree into a ansi-colored console string.
type ConsoleTreeRenderer struct {
	WithValues bool
}

// RenderTree renders a ContainerNode tree into a ansi-colored console string.
func (r ConsoleTreeRenderer) RenderTree(root *ContainerNode, _ string) string {
	sb := strings.Builder{}

	for _, node := range root.Items {
		fieldKey := node.Field().Name
		sb.WriteString(fieldKey)

		for i := len(fieldKey); i <= root.MaxKeyLength; i++ {
			sb.WriteRune(' ')
		}

		valueLen := 0
		preLen := DefaultValueOffset // Default spacing before the description when no values are rendered

		field := node.Field()

		if r.WithValues {
			preLen = ValueOffset // Adjusts the spacing when values are included
			valueLen = r.writeNodeValue(&sb, node)
		} else {
			// Since no values was supplied, let's substitute the value with the type
			typeName := field.Type.String()

			// If the value is an enum type, providing the name is a bit pointless
			// Instead, use a common string "option" to signify the type
			if field.EnumFormatter != nil {
				typeName = "option"
			}

			valueLen = len(typeName)
			sb.WriteString(color.CyanString(typeName))
		}

		sb.WriteString(strings.Repeat(" ", util.Max(preLen-valueLen, 1)))
		sb.WriteString(ColorizeDesc(field.Description))
		sb.WriteString(strings.Repeat(" ", util.Max(DescriptionColumnWidth-len(field.Description), 1)))

		if len(field.URLParts) > 0 && field.URLParts[0] != URLQuery {
			sb.WriteString(" <URL: ")

			for i, part := range field.URLParts {
				if i > 0 {
					sb.WriteString(", ")
				}

				if part > URLPath {
					part = URLPath
				}

				sb.WriteString(ColorizeEnum(part))
			}

			sb.WriteString(">")
		}

		if len(field.Template) > 0 {
			sb.WriteString(fmt.Sprintf(" <Template: %s>", ColorizeString(field.Template)))
		}

		if len(field.DefaultValue) > 0 {
			sb.WriteString(fmt.Sprintf(" <Default: %s>", ColorizeValue(field.DefaultValue, field.EnumFormatter != nil)))
		}

		if field.Required {
			sb.WriteString(fmt.Sprintf(" <%s>", ColorizeFalse("Required")))
		}

		if len(field.Keys) > 1 {
			sb.WriteString(" <Aliases: ")

			for i, key := range field.Keys {
				if i == 0 {
					// Skip primary alias (as it's the same as the field name)
					continue
				}

				if i > 1 {
					sb.WriteString(", ")
				}

				sb.WriteString(ColorizeString(key))
			}

			sb.WriteString(">")
		}

		if field.EnumFormatter != nil {
			sb.WriteString(ColorizeContainer(" ["))

			for i, name := range field.EnumFormatter.Names() {
				if i != 0 {
					sb.WriteString(", ")
				}

				sb.WriteString(ColorizeEnum(name))
			}

			sb.WriteString(ColorizeContainer("]"))
		}

		sb.WriteRune('\n')
	}

	return sb.String()
}

func (r ConsoleTreeRenderer) writeNodeValue(sb *strings.Builder, node Node) int {
	if contNode, isContainer := node.(*ContainerNode); isContainer {
		return r.writeContainer(sb, contNode)
	}

	if valNode, isValue := node.(*ValueNode); isValue {
		sb.WriteString(ColorizeToken(valNode.Value, valNode.tokenType))

		return len(valNode.Value)
	}

	sb.WriteRune('?')

	return 1
}

func (r ConsoleTreeRenderer) writeContainer(sb *strings.Builder, node *ContainerNode) int {
	kind := node.Type.Kind()

	hasKeys := !util.IsCollection(kind)

	totalLen := ContainerBracketLength // Length of the opening and closing brackets ({ } or [ ])

	if hasKeys {
		sb.WriteString("{ ")
	} else {
		sb.WriteString("[ ")
	}

	for i, itemNode := range node.Items {
		if i != 0 {
			sb.WriteString(", ")

			totalLen += KeySeparatorLength // This accounts for the : separator between keys and values in containers
		}

		if hasKeys {
			itemKey := itemNode.Field().Name
			sb.WriteString(itemKey)
			sb.WriteString(": ")

			totalLen += len(itemKey) + ItemSeparatorLength
		}

		valLen := r.writeNodeValue(sb, itemNode)
		totalLen += valLen
	}

	if hasKeys {
		sb.WriteString(" }")
	} else {
		sb.WriteString(" ]")
	}

	return totalLen
}
