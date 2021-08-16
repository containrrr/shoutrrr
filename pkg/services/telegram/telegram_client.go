package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
)

// Client for Telegram API
type Client struct {
	token string
}

func (c *Client) apiURL(endpoint string) string {
	return fmt.Sprintf(apiFormat, c.token, endpoint)
}

// GetBotInfo returns the bot User info
func (c *Client) GetBotInfo() (*User, error) {
	response := &userResponse{}
	err := jsonclient.Get(c.apiURL("getMe"), response)

	if !response.OK {
		return nil, GetErrorResponse(jsonclient.ErrorBody(err))
	}

	return &response.Result, nil
}

// GetUpdates retrieves the latest updates
func (c *Client) GetUpdates(offset int, limit int, timeout int, allowedUpdates []string) ([]Update, error) {

	request := &updatesRequest{
		Offset:         offset,
		Limit:          limit,
		Timeout:        timeout,
		AllowedUpdates: allowedUpdates,
	}
	response := &updatesResponse{}
	err := jsonclient.Post(c.apiURL("getUpdates"), request, response)

	if !response.OK {
		return nil, GetErrorResponse(jsonclient.ErrorBody(err))
	}

	return response.Result, nil
}

// SendMessage sends the specified Message
func (c *Client) SendMessage(message *SendMessagePayload) (*Message, error) {

	response := &messageResponse{}
	err := jsonclient.Post(c.apiURL("sendMessage"), message, response)

	if !response.OK {
		return nil, GetErrorResponse(jsonclient.ErrorBody(err))
	}

	return response.Result, nil
}

// GetErrorResponse retrieves the error message from a failed request
func GetErrorResponse(body string) error {
	response := &errorResponse{}
	if err := json.Unmarshal([]byte(body), response); err == nil {
		return response
	}
	return nil
}
