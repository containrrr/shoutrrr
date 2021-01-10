package util

import (
	"reflect"
)

// IsUnsignedInt is a check against the unsigned integer types
func IsUnsignedInt(kind reflect.Kind) bool {
	return kind >= reflect.Uint && kind <= reflect.Uint64
}

// IsSignedInt is a check against the signed decimal types
func IsSignedInt(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Int64
}

// IsCollection is a check against slice and array
func IsCollection(kind reflect.Kind) bool {
	return kind == reflect.Slice || kind == reflect.Array
}

// IsNumeric returns whether the Kind is one of the numeric ones
func IsNumeric(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Complex128
}
