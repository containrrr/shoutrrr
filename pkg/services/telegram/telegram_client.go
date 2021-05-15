package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
	"strconv"
	"strings"
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

// SetCommands sets what commands are available in the bot
func (c *Client) SetCommands(commands map[string]string) error {

	var cmds []command
	for cmd, desc := range commands {
		cmds = append(cmds, command{
			Command:     cmd,
			Description: desc,
		})
	}

	request := &commandsRequest{Commands: cmds}
	response := &errorResponse{}
	if err := jsonclient.Post(c.apiURL("setMyCommands"), request, response); err != nil {
		return err
	}

	if !response.OK {
		return response
	}

	return nil
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

func (c *Client) parseCommand(message *Message, botName string, private bool) (command string, params []string, err error) {
	parts := strings.Split(message.Text, " ")
	if len(parts) < 1 {
		return "", nil, fmt.Errorf("Message is empty")
	}

	atBot := "@" + botName

	if parts[0] == atBot {
		if len(parts) > 1 {
			params = parts[2:]
			command = parts[1]
		}
	} else {
		cmdParts := strings.Split(parts[0], "@")

		if len(cmdParts) < 2 || cmdParts[1] != botName {
			if private {
				command = parts[0]
			} else {
				return "", nil, fmt.Errorf("no bot tag")
			}
		} else {
			command = cmdParts[0]
		}

		params = parts[1:]
	}

	if len(command) < 1 || command[0] != '/' {
		return "", nil, fmt.Errorf("Message is not a command")
	}

	return command[1:], params, nil
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

// UpdateMessage updates an already sent Message
func (c *Client) UpdateMessage(update *updateMessagePayload) error {
	response := &errorResponse{}
	if err := jsonclient.Post(c.apiURL("editMessageText"), update, response); err != nil {
		return err
	}

	if !response.OK {
		return response
	}

	return nil
}

// Reply sends a Message containing another Message as a reply
func (c *Client) Reply(original *Message, text string) (*Message, error) {
	return c.SendMessage(&SendMessagePayload{
		Text:      text,
		ID:        strconv.FormatInt(original.Chat.ID, 10),
		ParseMode: ParseModes.Markdown.String(),
		ReplyTo:   original.MessageID,
	})
}

func (c *Client) answerCallbackQuery(answer *callbackQueryAnswer) error {
	response := &errorResponse{}
	if err := jsonclient.Post(c.apiURL("answerCallbackQuery"), answer, response); err != nil {
		return err
	}

	if !response.OK {
		return response
	}

	return nil
}

// GetErrorResponse retrieves the error message from a failed request
func GetErrorResponse(body string) error {
	response := &errorResponse{}
	if err := json.Unmarshal([]byte(body), response); err == nil {
		return response
	}
	return nil
}
