package shoutrrr

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"net/url"
)

// Send lets you send shoutrrr notifications using a supplied url and message
func Send(rawURL string, message string, opts types.ServiceOpts) error {
	routing := router.ServiceRouter{}
	if plugin, err := routing.Locate(rawURL); err != nil {
		return err
	} else if serviceURL, err := url.Parse(rawURL); err != nil {
		return err
	} else {
		return plugin.Send(serviceURL, message, opts)
	}
}

func Verify(rawURL string) error {

	routing := router.ServiceRouter{}

	svc, url, err := routing.ExtractServiceName(rawURL)
	if err != nil {
		return err
	}

	plugin, err := routing.Locate(svc)
	if err != nil {
		return err
	}

	config := plugin.GetConfig()
	if err := config.SetURL(url); err != nil {
		return err
	}

	configMap := format.GetConfigMap(config)
	for key, value := range configMap {
		fmt.Printf("%s: %s\n", key, value)
	}

	return nil
}