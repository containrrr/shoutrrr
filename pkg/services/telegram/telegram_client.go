package telegram

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/common/webclient"
)

// Client for Telegram API
type Client struct {
	WebClient webclient.WebClient
	token     string
}

func (c *Client) apiURL(endpoint string) string {
	return fmt.Sprintf(apiFormat, c.token, endpoint)
}

// GetBotInfo returns the bot User info
func (c *Client) GetBotInfo() (*User, error) {
	response := &userResponse{}
	err := c.WebClient.Get(c.apiURL("getMe"), response)

	if !response.OK {
		return nil, c.getErrorResponse(err)
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
	err := c.WebClient.Post(c.apiURL("getUpdates"), request, response)

	if !response.OK {
		return nil, c.getErrorResponse(err)
	}

	return response.Result, nil
}

// SendMessage sends the specified Message
func (c *Client) SendMessage(message *SendMessagePayload) (*Message, error) {

	response := &messageResponse{}
	err := c.WebClient.Post(c.apiURL("sendMessage"), message, response)

	if !response.OK {
		return nil, c.getErrorResponse(err)
	}

	return response.Result, nil
}

// GetErrorResponse retrieves the error message from a failed request
func (c *Client) getErrorResponse(err error) error {
	var errResponse *errorResponse
	if c.WebClient.ErrorResponse(err, errResponse) {
		return errResponse
	} else {
		return err
	}
}
