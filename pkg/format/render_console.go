package format

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/containrrr/shoutrrr/pkg/util"
)

// ConsoleTreeRenderer renders a ContainerNode tree into a ansi-colored console string
type ConsoleTreeRenderer struct {
	WithValues bool
}

// RenderTree renders a ContainerNode tree into a ansi-colored console string
func (r ConsoleTreeRenderer) RenderTree(root *ContainerNode, _ string) string {

	sb := strings.Builder{}

	for _, node := range root.Items {
		fieldKey := node.Field().Name
		sb.WriteString(fieldKey)
		for i := len(fieldKey); i <= root.MaxKeyLength; i++ {
			sb.WriteRune(' ')
		}

		valueLen := 0
		preLen := 16

		field := node.Field()

		if r.WithValues {
			preLen = 30
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
		sb.WriteString(strings.Repeat(" ", util.Max(60-len(field.Description), 1)))

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

	totalLen := 4
	if hasKeys {
		sb.WriteString("{ ")
	} else {
		sb.WriteString("[ ")
	}
	for i, itemNode := range node.Items {
		if i != 0 {
			sb.WriteString(", ")
			totalLen += 2
		}
		if hasKeys {
			itemKey := itemNode.Field().Name
			sb.WriteString(itemKey)
			sb.WriteString(": ")
			totalLen += len(itemKey) + 2
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
