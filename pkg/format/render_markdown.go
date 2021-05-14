package format

import (
	"sort"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/util"
)

// MarkdownTreeRenderer renders a ContainerNode tree into a markdown documentation string
type MarkdownTreeRenderer struct{}

// RenderTree renders a ContainerNode tree into a markdown documentation string
func (r MarkdownTreeRenderer) RenderTree(root *ContainerNode, scheme string) string {

	sb := strings.Builder{}

	queryFields := make([]*FieldInfo, 0, len(root.Items))
	urlFields := make([]*FieldInfo, URLPath+1)
	fieldsPrinted := make(map[string]bool)

	for _, node := range root.Items {
		field := node.Field()
		println(field.Name, len(field.URLParts))
		for _, urlPart := range field.URLParts {
			if urlPart == URLQuery {
				queryFields = append(queryFields, field)
			} else if urlPart > URLPath {
				urlFields = append(urlFields, field)
			} else {
				urlFields[urlPart] = field
			}
		}
		if len(field.URLParts) < 1 {
			queryFields = append(queryFields, field)
		}
	}

	sb.WriteString("## URL Fields\n\n")
	for _, field := range urlFields {
		if field == nil || fieldsPrinted[field.Name] {
			continue
		}
		r.writeFieldPrimary(&sb, field)

		sb.WriteString("  URL part: <code class=\"service-url\">")
		for i, uf := range urlFields {
			urlPart := URLPart(i)
			if urlPart > URLPath {
				// urlPart = URLPath
			}
			if urlPart == URLQuery {
				sb.WriteString(scheme)
				sb.WriteString("://")
				continue
			}
			if uf == nil {
				if urlPart == URLPath {
					sb.WriteRune(urlPart.Suffix())
				}
				continue
			} else if urlPart > URLUser {
				lastPart := urlPart - 1
				sb.WriteRune(lastPart.Suffix())
			}
			if field.IsURLPart(urlPart) {
				sb.WriteString("<strong>")
			}
			sb.WriteString(strings.ToLower(uf.Name))
			if field.IsURLPart(urlPart) {
				sb.WriteString("</strong>")
			}
		}
		sb.WriteString("</code>  \n")

		fieldsPrinted[field.Name] = true
	}

	sort.SliceStable(queryFields, func(i, j int) bool {
		return queryFields[i].Required && !queryFields[j].Required
	})

	sb.WriteString("## Query/Param Props\n\n")
	for _, field := range queryFields {
		r.writeFieldPrimary(&sb, field)
		r.writeFieldExtras(&sb, field)
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (MarkdownTreeRenderer) writeFieldExtras(sb *strings.Builder, field *FieldInfo) {
	if len(field.Keys) > 1 {
		sb.WriteString("  Aliases: `")
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
		sb.WriteString("  Possible values: `")
		for i, name := range field.EnumFormatter.Names() {
			if i != 0 {
				sb.WriteString("`, `")
			}
			sb.WriteString(name)
		}

		sb.WriteString("`  \n")
	}
}

func (MarkdownTreeRenderer) writeFieldPrimary(sb *strings.Builder, field *FieldInfo) {
	fieldKey := field.Name

	sb.WriteString("*  __")
	sb.WriteString(fieldKey)
	sb.WriteString("__")

	if field.Description != "" {
		sb.WriteString(" - ")
		sb.WriteString(field.Description)
	}

	if field.Required {
		sb.WriteString(" (**Required**)  \n")
	} else {
		sb.WriteString("  \n  Default: ")
		if field.DefaultValue == "" {
			sb.WriteString("*empty*")
		} else {
			sb.WriteRune('`')
			sb.WriteString(field.DefaultValue)
			sb.WriteRune('`')
		}
		sb.WriteString("  \n")
	}
}

func (r MarkdownTreeRenderer) writeNodeValue(sb *strings.Builder, node Node) int {
	if contNode, isContainer := node.(*ContainerNode); isContainer {
		return r.writeContainer(sb, contNode)
	}

	if valNode, isValue := node.(*ValueNode); isValue {
		sb.WriteString(valNode.Value)
		return len(valNode.Value)
	}

	sb.WriteRune('?')
	return 1
}

func (r MarkdownTreeRenderer) writeContainer(sb *strings.Builder, node *ContainerNode) int {
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
