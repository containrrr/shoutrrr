package format

import (
	"github.com/containrrr/shoutrrr/pkg/util"
	"strings"
)

func MarkdownFormattedTree(root *ContainerNode) string {

	sb := strings.Builder{}
	sb.WriteString("## Config fields\n\n")
	for _, node := range root.Items {
		field := node.Field()
		fieldKey := field.Name
		sb.WriteString("### ")
		sb.WriteString(fieldKey)
		sb.WriteString("  \n")

		//sb.WriteString("**")
		//typeName := field.Type.String()
		//sb.WriteString(typeName)
		//sb.WriteString("**  \n")

		sb.WriteString(field.Description)
		sb.WriteString("  \n")

		if field.Required {
			sb.WriteString(" **Required**  \n")
		} else {
			sb.WriteString(" **Default**: ")
			if field.DefaultValue == "" {
				sb.WriteString("*empty*")
			} else {
				sb.WriteRune('`')
				sb.WriteString(field.DefaultValue)
				sb.WriteRune('`')
			}
			sb.WriteString("  \n")
		}

		if len(field.Keys) > 1 {
			sb.WriteString(" **Aliases**: `")
			for i, key := range field.Keys {
				if i == 0 {
					// Skip primary alias (as it's the same as the field name)
					continue
				}
				if i > 1 {
					sb.WriteString("`, `")
				}
				sb.WriteString(key)
			}
			sb.WriteString("`  \n")
		}

		if field.EnumFormatter != nil {
			sb.WriteString(" **Possible values**: `")
			for i, name := range field.EnumFormatter.Names() {
				if i != 0 {
					sb.WriteString("`, `")
				}
				sb.WriteString(name)
			}

			sb.WriteString("`  \n")
		}

		sb.WriteRune('\n')
	}

	return sb.String()
}

func writeMarkdownNodeValue(sb *strings.Builder, node Node) int {
	if contNode, isContainer := node.(*ContainerNode); isContainer {
		return writeMarkdownContainer(sb, contNode)
	}

	if valNode, isValue := node.(*ValueNode); isValue {
		sb.WriteString(valNode.Value)
		return len(valNode.Value)
	}

	sb.WriteRune('?')
	return 1
}

func writeMarkdownContainer(sb *strings.Builder, node *ContainerNode) int {
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
		valLen := writeMarkdownNodeValue(sb, itemNode)
		totalLen += valLen
	}
	if hasKeys {
		sb.WriteString(" }")
	} else {
		sb.WriteString(" ]")
	}
	return totalLen
}
