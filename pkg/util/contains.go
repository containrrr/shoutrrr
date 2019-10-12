package util

import (
	"reflect"
)

// ContainsKind returns true if the candidate is present in the kind slice
func ContainsKind(kinds []reflect.Kind, candidate reflect.Kind) bool {
	for _, item := range kinds {
		if item == candidate {
			return true
		}
	}
	return false
}
