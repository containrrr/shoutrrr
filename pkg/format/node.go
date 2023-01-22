package format

import (
	"fmt"
	r "reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// NodeTokenType is used to represent the type of value that a node has for syntax highlighting
type NodeTokenType int

const (
	// UnknownToken represents all unknown/unspecified tokens
	UnknownToken NodeTokenType = iota
	// NumberToken represents all numbers
	NumberToken
	// StringToken represents strings and keys
	StringToken
	// EnumToken represents enum values
	EnumToken
	// TrueToken represent boolean true
	TrueToken
	// FalseToken represent boolean false
	FalseToken
	// PropToken represent a serializable struct prop
	PropToken
	// ErrorToken represent a value that was not serializable or otherwise invalid
	ErrorToken
	// ContainerToken is used for Array/Slice and Map tokens
	ContainerToken
)

// Node is the generic config tree item
type Node interface {
	Field() *FieldInfo
	TokenType() NodeTokenType
	Update(tv r.Value)
}

// ValueNode is a Node without any child items
type ValueNode struct {
	*FieldInfo
	Value     string
	tokenType NodeTokenType
}

// Field returns the inner FieldInfo
func (n *ValueNode) Field() *FieldInfo {
	return n.FieldInfo
}

// TokenType returns a NodeTokenType that matches the value
func (n *ValueNode) TokenType() NodeTokenType {
	return n.tokenType

}

// Update updates the value string from the provided value
func (n *ValueNode) Update(tv r.Value) {
	value, token := getValueNodeValue(tv, n.FieldInfo)
	n.Value = value
	n.tokenType = token
}

// ContainerNode is a Node with child items
type ContainerNode struct {
	*FieldInfo
	Items        []Node
	MaxKeyLength int
}

// Field returns the inner FieldInfo
func (n *ContainerNode) Field() *FieldInfo {
	return n.FieldInfo
}

// TokenType always returns ContainerToken for ContainerNode
func (n *ContainerNode) TokenType() NodeTokenType {
	return ContainerToken
}

// Update updates the items to match the provided value
func (n *ContainerNode) Update(tv r.Value) {
	switch n.FieldInfo.Type.Kind() {
	case r.Array, r.Slice:
		n.updateArrayNode(tv)
	case r.Map:
		n.updateMapNode(tv)
	}
}

func (n *ContainerNode) updateArrayNode(arrayValue r.Value) {
	itemCount := arrayValue.Len()
	n.Items = make([]Node, 0, itemCount)

	elemType := arrayValue.Type().Elem()
	for i := 0; i < itemCount; i++ {
		key := strconv.Itoa(i)
		val := arrayValue.Index(i)
		n.Items = append(n.Items, getValueNode(val, &FieldInfo{
			Name: key,
			Type: elemType,
		}))
	}
}

func getArrayNode(arrayValue r.Value, fieldInfo *FieldInfo) (node *ContainerNode) {
	node = &ContainerNode{
		FieldInfo:    fieldInfo,
		MaxKeyLength: 0,
	}

	node.updateArrayNode(arrayValue)

	return node
}

func sortNodeItems(nodeItems []Node) {
	sort.Slice(nodeItems, func(i, j int) bool {
		return nodeItems[i].Field().Name < nodeItems[j].Field().Name
	})
}

func (n *ContainerNode) updateMapNode(mapValue r.Value) {
	base := n.FieldInfo.Base
	if base == 0 {
		base = 10
	}
	elemType := mapValue.Type().Elem()
	mapKeys := mapValue.MapKeys()
	nodeItems := make([]Node, len(mapKeys))
	maxKeyLength := 0
	for i, keyVal := range mapKeys {
		// The keys will always be strings
		key := keyVal.String()
		val := mapValue.MapIndex(keyVal)
		nodeItems[i] = getValueNode(val, &FieldInfo{
			Name: key,
			Type: elemType,
			Base: base,
		})
		maxKeyLength = util.Max(len(key), maxKeyLength)
	}
	sortNodeItems(nodeItems)

	n.Items = nodeItems
	n.MaxKeyLength = maxKeyLength
}

func getMapNode(mapValue r.Value, fieldInfo *FieldInfo) (node *ContainerNode) {
	if mapValue.Kind() == r.Ptr {
		mapValue = mapValue.Elem()
	}

	node = &ContainerNode{
		FieldInfo: fieldInfo,
	}

	node.updateMapNode(mapValue)

	return
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

func getRootNode(v interface{}) *ContainerNode {
	structValue := r.ValueOf(v)
	if structValue.Kind() == r.Ptr {
		structValue = structValue.Elem()
	}

	fieldInfo := &FieldInfo{
		Type: structValue.Type(),
	}

	enums := map[string]types.EnumFormatter{}
	if enummer, isEnummer := v.(types.Enummer); isEnummer {
		enums = enummer.Enums()
	}

	infoFields := getStructFieldInfo(fieldInfo.Type, enums)

	numFields := len(infoFields)
	nodeItems := make([]Node, numFields)
	maxKeyLength := 0
	fieldOffset := 0
	for i := range infoFields {
		field := infoFields[i]
		for isHiddenField(fieldInfo.Type.Field(fieldOffset + i)) {
			// The current field is Anonymous and not present in the FieldInfo slice
			fieldOffset++
		}
		fieldValue := structValue.Field(fieldOffset + i)

		nodeItems[i] = getNode(fieldValue, &field)
		maxKeyLength = util.Max(len(field.Name), maxKeyLength)
	}

	sortNodeItems(nodeItems)

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
	case r.Ptr, r.Struct:
		if val, err := GetConfigPropString(fieldValue); err == nil {
			return val, PropToken
		}
		return "<ERR>", ErrorToken
	}

	// Unsupported value
	return fmt.Sprintf("<?%s>", kind.String()), UnknownToken
}

func getContainerValueString(fieldValue r.Value, fieldInfo *FieldInfo) string {
	itemSep := fieldInfo.ItemSeparator
	sliceLen := fieldValue.Len()
	var mapKeys []r.Value
	if fieldInfo.Type.Kind() == r.Map {
		mapKeys = fieldValue.MapKeys()
		sort.Slice(mapKeys, func(a, b int) bool {
			return mapKeys[a].String() < mapKeys[b].String()
		})
	}

	sb := strings.Builder{}
	var itemFieldInfo *FieldInfo
	for i := 0; i < sliceLen; i++ {
		var itemValue r.Value
		if i > 0 {
			sb.WriteRune(itemSep)
		}

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
		}
		strVal, _ := getValueNodeValue(itemValue, itemFieldInfo)
		sb.WriteString(strVal)
	}
	return sb.String()
}
