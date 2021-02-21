package format

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	r "reflect"
	"sort"
	"strconv"
	"strings"
)

type NodeTokenType int

const (
	UnknownToken NodeTokenType = iota
	NumberToken
	StringToken
	EnumToken
	TrueToken
	FalseToken
	PropToken
	ErrorToken
	ContainerToken
)

type Node interface {
	Field() *FieldInfo
	TokenType() NodeTokenType
}

type ValueNode struct {
	*FieldInfo
	Value     string
	tokenType NodeTokenType
}

func (n *ValueNode) Field() *FieldInfo {
	return n.FieldInfo
}

func (n *ValueNode) TokenType() NodeTokenType {
	return n.tokenType
}

type ContainerNode struct {
	*FieldInfo
	Items        []Node
	MaxKeyLength int
}

func (n *ContainerNode) Field() *FieldInfo {
	return n.FieldInfo
}

func (n *ContainerNode) TokenType() NodeTokenType {
	return ContainerToken
}

func getArrayNode(arrayValue r.Value, fieldInfo *FieldInfo) (node *ContainerNode) {
	itemCount := arrayValue.Len()
	nodeItems := make([]Node, 0, itemCount)

	for i := 0; i < itemCount; i++ {
		key := strconv.Itoa(i)
		val := arrayValue.Index(i)
		nodeItems = append(nodeItems, getValueNode(val, &FieldInfo{
			Name: key,
			Type: val.Type(),
		}))
	}

	return &ContainerNode{
		FieldInfo:    fieldInfo,
		Items:        nodeItems,
		MaxKeyLength: 0,
	}
}

func getMapNode(mapValue r.Value, fieldInfo *FieldInfo) *ContainerNode {
	if mapValue.Kind() == r.Ptr {
		mapValue = mapValue.Elem()
	}

	mapKeys := mapValue.MapKeys()
	nodeItems := make([]Node, len(mapKeys))
	base := fieldInfo.Base
	if base == 0 {
		base = 10
	}

	maxKeyLength := 0
	for i, keyVal := range mapKeys {
		// The keys will always be strings
		key := keyVal.String()
		val := mapValue.MapIndex(keyVal)
		nodeItems[i] = getValueNode(val, &FieldInfo{
			Name: key,
			Type: val.Type(),
			Base: base,
		})
		maxKeyLength = util.Max(len(key), maxKeyLength)
	}

	return sortedContainerNode(nodeItems, fieldInfo, maxKeyLength)
}

func getNode(fieldVal r.Value, fieldInfo *FieldInfo) Node {
	switch fieldInfo.Type.Kind() {
	case r.Array, r.Slice:
		return getArrayNode(fieldVal, fieldInfo)
	case r.Map:
		return getMapNode(fieldVal, fieldInfo)
	default:
		return getValueNode(fieldVal, fieldInfo)
	}
}

func getRootNode(config types.ServiceConfig) *ContainerNode {
	structValue := r.ValueOf(config)
	if structValue.Kind() == r.Ptr {
		structValue = structValue.Elem()
	}

	fieldInfo := &FieldInfo{
		Type: structValue.Type(),
	}

	infoFields := getStructFieldInfo(fieldInfo.Type, config.Enums())

	numFields := len(infoFields)
	nodeItems := make([]Node, numFields)
	maxKeyLength := 0
	fieldOffset := 0
	for i := range infoFields {
		field := infoFields[i]
		if fieldInfo.Type.Field(fieldOffset + i).Anonymous {
			// The current field is Anonymous and not present in the FieldInfo slice
			fieldOffset += 1
		}
		fieldValue := structValue.Field(fieldOffset + i)

		nodeItems[i] = getNode(fieldValue, &field)
		maxKeyLength = util.Max(len(field.Name), maxKeyLength)
	}

	return sortedContainerNode(nodeItems, fieldInfo, maxKeyLength)
}

func sortedContainerNode(nodeItems []Node, fieldInfo *FieldInfo, maxKeyLength int) *ContainerNode {
	sort.Slice(nodeItems, func(i, j int) bool {
		return nodeItems[i].Field().Name < nodeItems[j].Field().Name
	})

	return &ContainerNode{
		FieldInfo:    fieldInfo,
		Items:        nodeItems,
		MaxKeyLength: maxKeyLength,
	}
}

func getValueNode(fieldVal r.Value, fieldInfo *FieldInfo) (node *ValueNode) {
	value, tokenType := getValueNodeValue(fieldVal, fieldInfo)
	return &ValueNode{
		FieldInfo: fieldInfo,
		Value:     value,
		tokenType: tokenType,
	}
}

func getValueNodeValue(fieldValue r.Value, fieldInfo *FieldInfo) (string, NodeTokenType) {
	kind := fieldValue.Kind()
	base := fieldInfo.Base
	if base == 0 {
		base = 10
	}

	if fieldInfo.IsEnum() {
		return fieldInfo.EnumFormatter.Print(int(fieldValue.Int())), EnumToken
	}
	switch kind {
	case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
		val := strconv.FormatUint(fieldValue.Uint(), base)
		if base == 16 {
			val = "0x" + val
		}
		return val, NumberToken
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		return strconv.FormatInt(fieldValue.Int(), base), NumberToken
	case r.String:
		return fieldValue.String(), StringToken
	case r.Bool:
		val := fieldValue.Bool()
		if val {
			return PrintBool(val), TrueToken
		}
		return PrintBool(val), FalseToken
	case r.Array, r.Slice, r.Map:
		return getContainerValueString(fieldValue, fieldInfo), UnknownToken
	case r.Struct, r.Ptr:
		if val, err := GetConfigPropString(fieldValue); err != nil {
			return val, PropToken
		}
		return "<ERR>", ErrorToken
	}

	// Unsupported value
	return fmt.Sprintf("<?%s>", kind.String()), UnknownToken
}

func getContainerValueString(fieldValue r.Value, fieldInfo *FieldInfo) string {
	sliceLen := fieldValue.Len()
	var mapKeys []r.Value
	if fieldInfo.Type.Kind() == r.Map {
		mapKeys = fieldValue.MapKeys()
	}

	sb := strings.Builder{}
	var itemFieldInfo *FieldInfo
	for i := 0; i < sliceLen; i++ {
		var itemValue r.Value
		if mapKeys != nil {
			mapKey := mapKeys[i]
			sb.WriteString(mapKey.String())
			sb.WriteRune(':')
			itemValue = fieldValue.MapIndex(mapKey)
		} else {
			itemValue = fieldValue.Index(i)
		}

		if i == 0 {
			itemFieldInfo = &FieldInfo{
				Type: itemValue.Type(),
				// Inherit the base from the container
				Base: fieldInfo.Base,
			}

			if itemFieldInfo.Base == 0 {
				itemFieldInfo.Base = 10
			}
		} else {
			sb.WriteRune(',')
		}
		strVal, _ := getValueNodeValue(itemValue, itemFieldInfo)
		sb.WriteString(strVal)
	}
	return sb.String()
}
