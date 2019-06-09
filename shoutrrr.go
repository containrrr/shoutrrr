package shoutrrr

import (
	"errors"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/plugin"
	"github.com/containrrr/shoutrrr/pkg/router"
	"net/url"
	"os"
)

// Send lets you send shoutrrr notifications using a supplied url and message
func Send(rawUrl string, message string, opts plugin.PluginOpts) error {
	routing := router.ServiceRouter{}
	if plugin, err := routing.Locate(rawUrl); err != nil {
		return err
	} else if serviceUrl, err := url.Parse(rawUrl); err != nil {
		return err
	} else {
		return plugin.Send(*serviceUrl, message, opts)
	}
}

// SendEnv lets you send shoutrrr notifications using an url stored in your env variables and a supplied message
func SendEnv(message string) error {
	envURL := os.Getenv("SHOUTRRR_URL")
	if envURL == "" {
		return errors.New("trying to use SendEnv but SHOUTRRR_URL is not set")
	}
	return Send(envURL, message, plugin.PluginOpts{})
}

func Verify(rawUrl string) error {

	routing := router.ServiceRouter{}

	svc, url, err := routing.ExtractServiceName(rawUrl)
	if err != nil {
		return err
	}

	if plugin, err := routing.Locate(svc); err != nil {
		return err
	} else {
		config := plugin.GetConfig()
		if err := config.SetURL(url); err != nil {
			return err
		}
		configMap := format.GetConfigMap(config)
		for key, value := range configMap {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	return nil
}