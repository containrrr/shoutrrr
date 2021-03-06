package format

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/util"
	"github.com/fatih/color"
	"strings"
)

func getColorFormattedTree(root *ContainerNode, withValues bool) string {

	sb := strings.Builder{}
	packageName := root.Type.String()
	packageName = packageName[:strings.LastIndexByte(packageName, '.')+1]

	for _, node := range root.Items {
		fieldKey := node.Field().Name
		sb.WriteString(fieldKey)
		for i := len(fieldKey); i <= root.MaxKeyLength; i++ {
			sb.WriteRune(' ')
		}

		valueLen := 0
		preLen := 16

		field := node.Field()

		if withValues {
			preLen = 30
			valueLen = writeColoredNodeValue(&sb, node)
		} else {
			// Since no values was supplied, let's substitute the value with the type
			typeName := strings.TrimPrefix(field.Type.String(), packageName)
			valueLen = len(typeName)
			sb.WriteString(color.CyanString(typeName))
		}

		//if len(field.Description) > 0 {
		sb.WriteString(strings.Repeat(" ", util.Max(preLen-valueLen, 1)))
		sb.WriteString(ColorizeDesc(field.Description))
		sb.WriteString(strings.Repeat(" ", util.Max(60-len(field.Description), 1)))
		//}

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

func writeColoredNodeValue(sb *strings.Builder, node Node) int {
	if contNode, isContainer := node.(*ContainerNode); isContainer {
		return writeColoredContainer(sb, contNode)
	}

	if valNode, isValue := node.(*ValueNode); isValue {
		sb.WriteString(ColorizeToken(valNode.Value, valNode.tokenType))
		return len(valNode.Value)
	}

	sb.WriteRune('?')
	return 1
}

func writeColoredContainer(sb *strings.Builder, node *ContainerNode) int {
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
		valLen := writeColoredNodeValue(sb, itemNode)
		totalLen += valLen
	}
	if hasKeys {
		sb.WriteString(" }")
	} else {
		sb.WriteString(" ]")
	}
	return totalLen
}
