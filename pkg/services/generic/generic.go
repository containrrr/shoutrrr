package generic

import (
	"encoding/json"
	"io/ioutil"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"

	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Service providing a generic notification service
type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

// Send a notification message to a generic webhook endpoint
func (service *Service) Send(message string, paramsPtr *types.Params) error {
	config := *service.config

	var params types.Params
	if paramsPtr == nil {
		params = types.Params{}
	} else {
		params = *paramsPtr
	}

	if err := service.pkr.UpdateConfigFromParams(&config, &params); err != nil {
		service.Logf("Failed to update params: %v", err)
	}

	updateParams(&config, params, message)

	if err := service.doSend(&config, params); err != nil {
		return fmt.Errorf("an error occurred while sending notification to generic webhook: %s", err.Error())
	}

	return nil
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	config, pkr := DefaultConfig()
	service.config = config
	service.pkr = pkr

	return service.config.setURL(&service.pkr, configURL)
}

// GetConfigURLFromCustom creates a regular service URL from one with a custom host
func (*Service) GetConfigURLFromCustom(customURL *url.URL) (serviceURL *url.URL, err error) {
	webhookURL := *customURL
	if strings.HasPrefix(webhookURL.Scheme, Scheme) {
		webhookURL.Scheme = webhookURL.Scheme[len(Scheme)+1:]
	}
	config, pkr, err := ConfigFromWebhookURL(webhookURL)
	if err != nil {
		return nil, err
	}
	return config.getURL(&pkr), nil
}

func (service *Service) doSend(config *Config, params types.Params) error {
	postURL := config.WebhookURL().String()
	payload, err := service.getPayload(config, params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(config.RequestMethod, postURL, payload)
	if err == nil {
		req.Header.Set("Content-Type", config.ContentType)
		req.Header.Set("Accept", config.ContentType)
		var res *http.Response
		res, err = http.DefaultClient.Do(req)
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			if body, errRead := ioutil.ReadAll(res.Body); errRead == nil {
				service.Log("Server response: ", string(body))
			}
		}
		if err == nil && res.StatusCode >= http.StatusMultipleChoices {
			err = fmt.Errorf("server returned response status code %s", res.Status)
		}
	}

	return err
}

func (service *Service) getPayload(config *Config, params types.Params) (io.Reader, error) {
	switch config.Template {
	case "":
		return bytes.NewBufferString(params[config.MessageKey]), nil
	case "json", "JSON":
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(jsonBytes), nil
	}
	tpl, found := service.GetTemplate(config.Template)
	if !found {
		return nil, fmt.Errorf("template %q has not been loaded", config.Template)
	}

	bb := &bytes.Buffer{}
	err := tpl.Execute(bb, params)
	return bb, err
}

func updateParams(config *Config, params types.Params, message string) {
	if title, found := params.Title(); found {
		if config.TitleKey != "title" {
			delete(params, "title")
			params[config.TitleKey] = title
		}
	}
	params[config.MessageKey] = message
}
