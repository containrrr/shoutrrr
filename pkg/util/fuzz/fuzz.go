//go:build gofuzz || cgo
// +build gofuzz cgo
// Note that the `cgo` above is a hack to prevent the file from being built in releases, but still
// included when checking for doc comments (lint)

package fuzz

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	t "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// FuzzPartitionMessage fuzzes the util.PartitionMessage function
func FuzzPartitionMessage(data []byte) int {
	f := fuzz.NewConsumer(data)

	input, err := f.GetString()
	if err != nil {
		return 0
	}

	limits := t.MessageLimit{}
	err = f.GenerateStruct(&limits)
	if err != nil {
		return 0
	}

	distance, err := f.GetInt()
	if err != nil {
		return 0
	}
	_, _ = util.PartitionMessage(input, limits, distance)
	return 1
}
