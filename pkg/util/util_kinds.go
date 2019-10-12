package util

import (
	"reflect"
)

// IsUnsignedDecimal is a check against the unsigned decimal types
func IsUnsignedDecimal(kind reflect.Kind) bool {
	unsignedDecimals := []reflect.Kind{
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uint,
	}
	return ContainsKind(unsignedDecimals, kind)
}

// IsSignedDecimal is a check against the signed decimal types
func IsSignedDecimal(kind reflect.Kind) bool {
	signedDecimals := []reflect.Kind{
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Int,
	}
	return ContainsKind(signedDecimals, kind)
}

// IsCollection is a check against slice and array
func IsCollection(kind reflect.Kind) bool {
	collections := []reflect.Kind{
		reflect.Slice,
		reflect.Array,
	}
	return ContainsKind(collections, kind)
}
