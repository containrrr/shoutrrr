package format

import (
	r "reflect"
	"strconv"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// FieldInfo is the meta data about a config field
type FieldInfo struct {
	Name          string
	Type          r.Type
	EnumFormatter types.EnumFormatter
	Description   string
	DefaultValue  string
	Template      string
	Required      bool
	URLParts      []URLPart
	Title         bool
	Base          int
	Keys          []string
	ItemSeparator rune
}

// IsEnum returns whether a EnumFormatter has been assigned to the field and that it is of a suitable type
func (fi *FieldInfo) IsEnum() bool {
	return fi.EnumFormatter != nil && fi.Type.Kind() == r.Int
}

// IsURLPart returns whether the field is serialized as the specified part of an URL
func (fi *FieldInfo) IsURLPart(part URLPart) bool {
	for _, up := range fi.URLParts {
		if up == part {
			return true
		}
	}
	return false
}

func getStructFieldInfo(structType r.Type, enums map[string]types.EnumFormatter) []FieldInfo {

	numFields := structType.NumField()
	fields := make([]FieldInfo, 0, numFields)
	maxKeyLen := 0

	for i := 0; i < numFields; i++ {
		fieldDef := structType.Field(i)

		if isHiddenField(fieldDef) {
			// This is an embedded or private field, which should not be part of the Config output
			continue
		}

		info := FieldInfo{
			Name:          fieldDef.Name,
			Type:          fieldDef.Type,
			Required:      true,
			Title:         false,
			ItemSeparator: ',',
		}

		if util.IsNumeric(fieldDef.Type.Kind()) {
			info.Base = getFieldBase(fieldDef)
		}

		if tag, ok := fieldDef.Tag.Lookup("desc"); ok {
			info.Description = tag
		}

		if tag, ok := fieldDef.Tag.Lookup("tpl"); ok {
			info.Template = tag
		}

		if tag, ok := fieldDef.Tag.Lookup("default"); ok {
			info.Required = false
			info.DefaultValue = tag
		}

		if _, ok := fieldDef.Tag.Lookup("optional"); ok {
			info.Required = false
		}

		if _, ok := fieldDef.Tag.Lookup("title"); ok {
			info.Title = true
		}

		if tag, ok := fieldDef.Tag.Lookup("url"); ok {
			info.URLParts = ParseURLParts(tag)
		}

		if tag, ok := fieldDef.Tag.Lookup("key"); ok {
			tag := strings.ToLower(tag)
			info.Keys = strings.Split(tag, ",")
		}

		if tag, ok := fieldDef.Tag.Lookup("sep"); ok {
			info.ItemSeparator = rune(tag[0])
		}

		if ef, isEnum := enums[fieldDef.Name]; isEnum {
			info.EnumFormatter = ef
		}

		fields = append(fields, info)
		keyLen := len(fieldDef.Name)
		if keyLen > maxKeyLen {
			maxKeyLen = keyLen
		}
	}

	return fields
}

func isHiddenField(field r.StructField) bool {
	return field.Anonymous || strings.ToUpper(field.Name[0:1]) != field.Name[0:1]
}

func getFieldBase(field r.StructField) int {
	if tag, ok := field.Tag.Lookup("base"); ok {
		if base, err := strconv.ParseUint(tag, 10, 8); err == nil {
			return int(base)
		}
	}

	// Default to base 10 if not tagged
	return 10
}
