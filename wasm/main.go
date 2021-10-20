package main

import (
	"fmt"

	"github.com/containrrr/shoutrrr/internal/meta"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
)

func main() {
	println("Shoutrrr WASM", meta.Version)
}

func verify(url string) (string, bool) {
	sr := router.ServiceRouter{}

	service, err := sr.Locate(url)

	if err != nil {
		return fmt.Sprintf("error verifying URL: %s\n", err), false
	}

	config := format.GetServiceConfig(service)
	configNode := format.GetConfigFormat(config)

	return format.ColorFormatTree(configNode, true), true
}
