package util

import (
	fuzz "github.com/AdaLogics/go-fuzz-headers"
	t "github.com/containrrr/shoutrrr/pkg/types"
)

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
	_, _ = PartitionMessage(input, limits, distance)
	return 1
}
