//go:build !tinygo
// +build !tinygo

package generic

import (
	"fmt"
	"net/http"

	"github.com/containrrr/shoutrrr/pkg/types"
)

func (service *Service) doSend(config *Config, message string, params *types.Params) error {
	postURL := config.WebhookURL().String()
	payload, err := service.getPayload(config.Template, message, params)
	if err != nil {
		return err
	}

	res, err := http.Post(postURL, config.ContentType, payload)
	if err == nil && res.StatusCode != http.StatusOK {
		err = fmt.Errorf("server returned response status code %s", res.Status)
	}

	return err
}
