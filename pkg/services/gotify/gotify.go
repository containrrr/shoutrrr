package gotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

// Service providing Gotify as a notification service
type Service struct {
	standard.Standard
	config *Config
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger *log.Logger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{
		Title: "Shoutrrr notification",
	}
	err := service.config.SetURL(configURL)
	return err
}

const tokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_"

// The validation rules have been taken directly from the Gotify source code.
// These will have to be adapted in case of a change:
// https://github.com/gotify/server/blob/ad157a138b4985086c484a7aabfc2deada5a33dd/auth/token.go#L8
func isTokenValid(token string) bool {
	if len(token) != 15 {
		return false
	} else if token[0] != 'A' {
		return false
	}
	for _, c := range token {
		if !strings.ContainsRune(tokenChars, c) {
			return false
		}
	}
	return true
}

func buildURL(config *Config) (string, error) {
	token := config.Token
	if len(token) > 0 && token[0] == '/' {
		token = token[1:]
	}
	if !isTokenValid(token) {
		return "", fmt.Errorf("invalid gotify token \"%s\"", token)
	}
	return fmt.Sprintf("https://%s/message?token=%s", config.Host, token), nil
}

func getPriority(params map[string]string) int {
	priorityStr, ok := params["priority"]
	if !ok {
		priorityStr = "0"
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		priority = 0
	}
	return priority
}

func getTitle(params map[string]string, config *Config) string {
	title, ok := params["title"]
	if !ok {
		title = config.Title
	}
	return title
}

// Send a notification message to Gotify
func (service *Service) Send(message string, params *types.Params) error {
	if params == nil {
		params = &types.Params{}
	}
	config := service.config
	postURL, err := buildURL(config)
	if err != nil {
		return err
	}
	jsonBody, err := json.Marshal(JSON{
		Message:  message,
		Title:    getTitle(*params, config),
		Priority: getPriority(*params),
	})
	if err != nil {
		return err
	}
	jsonBuffer := bytes.NewBuffer(jsonBody)
	resp, err := http.Post(postURL, "application/json", jsonBuffer)
	if err != nil {
		return fmt.Errorf("failed to send notification to Gotify: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Gotify notification returned %d HTTP status code", resp.StatusCode)
	}

	return nil
}
