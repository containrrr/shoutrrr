package format

import (
	"reflect"
	"sort"
	"strings"
)

// MarkdownTreeRenderer renders a ContainerNode tree into a markdown documentation string
type MarkdownTreeRenderer struct {
	HeaderPrefix      string
	PropsDescription  string
	PropsEmptyMessage string
}

// RenderTree renders a ContainerNode tree into a markdown documentation string
func (r MarkdownTreeRenderer) RenderTree(root *ContainerNode, scheme string) string {

	sb := strings.Builder{}

	queryFields := make([]*FieldInfo, 0, len(root.Items))
	urlFields := make([]*FieldInfo, URLPath+1)

	for _, node := range root.Items {
		field := node.Field()
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

	r.writeURLFields(&sb, urlFields, scheme)

	sort.SliceStable(queryFields, func(i, j int) bool {
		return queryFields[i].Required && !queryFields[j].Required
	})

	r.writeHeader(&sb, "Query/Param Props")
	if len(queryFields) > 0 {
		sb.WriteString(r.PropsDescription)
	} else {
		sb.WriteString(r.PropsEmptyMessage)
	}
	sb.WriteRune('\n')
	for _, field := range queryFields {
		r.writeFieldPrimary(&sb, field)
		r.writeFieldExtras(&sb, field)
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (r MarkdownTreeRenderer) writeURLFields(sb *strings.Builder, urlFields []*FieldInfo, scheme string) {
	fieldsPrinted := make(map[string]bool)

	sort.SliceStable(urlFields, func(i, j int) bool {
		if urlFields[i] == nil || urlFields[j] == nil {
			return false
		}

		urlPartA := URLQuery
		if len(urlFields[i].URLParts) > 0 {
			urlPartA = urlFields[i].URLParts[0]
		}

		urlPartB := URLQuery
		if len(urlFields[j].URLParts) > 0 {
			urlPartB = urlFields[j].URLParts[0]
		}

		return urlPartA < urlPartB
	})

	r.writeHeader(sb, "URL Fields")
	for _, field := range urlFields {
		if field == nil || fieldsPrinted[field.Name] {
			continue
		}
		r.writeFieldPrimary(sb, field)

		sb.WriteString("  URL part: <code class=\"service-url\">")

		for i, uf := range urlFields {
			urlPart := URLPart(i)
			if urlPart == URLQuery {
				sb.WriteString(scheme)
				sb.WriteString("://")
				continue
			}
			if uf == nil {
				if urlPart == URLPath {
					sb.WriteRune(urlPart.Suffix())
				} else if urlPart == URLHost {
					// Host cannot be empty
					if urlFields[URLPassword] != nil || urlFields[URLUser] != nil {
						sb.WriteRune(URLPassword.Suffix())
					}
					sb.WriteString(scheme)
				}
				continue
			} else if urlPart == URLHost && urlFields[URLUser] == nil && urlFields[URLPassword] == nil {
			} else if urlPart > URLUser {
				lastPart := urlPart - 1
				sb.WriteRune(lastPart.Suffix())
			}
			if field.IsURLPart(urlPart) {
				sb.WriteString("<strong>")
			}

			slug := strings.ToLower(uf.Name)

			// Hard coded override for host:port 😓
			if slug == "host" && urlPart == URLPort {
				slug = "port"
			}
			sb.WriteString(slug)

			if field.IsURLPart(urlPart) {
				sb.WriteString("</strong>")
			}

		}
		sb.WriteString("</code>  \n")

		fieldsPrinted[field.Name] = true
	}
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
			if field.Type.Kind() == reflect.Bool {
				defaultValue, _ := ParseBool(field.DefaultValue, false)
				if defaultValue {
					sb.WriteString("✔ ")
				} else {
					sb.WriteString("❌ ")
				}
			}
			sb.WriteRune('`')
			sb.WriteString(field.DefaultValue)
			sb.WriteRune('`')
		}
		sb.WriteString("  \n")
	}
}

func (r MarkdownTreeRenderer) writeHeader(sb *strings.Builder, text string) {
	sb.WriteString(r.HeaderPrefix)
	sb.WriteString(text)
	sb.WriteString("\n\n")
}
