package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/util/jsonclient"
	"strconv"
	"strings"
)

type Client struct {
	token string
}

func (c *Client) ApiURL(endpoint string) string {
	return fmt.Sprintf(apiFormat, c.token, endpoint)
}

func (c *Client) GetBotInfo() (*User, error) {
	response := &UserResponse{}
	err := jsonclient.Get(c.ApiURL("getMe"), response)

	if !response.OK {
		return nil, GetErrorResponse(err.Body)
	}

	return &response.Result, nil
}

func (c *Client) SetCommands(commands map[string]string) error {

	var cmds []Command
	for cmd, desc := range commands {
		cmds = append(cmds, Command{
			Command:     cmd,
			Description: desc,
		})
	}

	request := &CommandsRequest{Commands: cmds}
	response := &ErrorResponse{}
	if err := jsonclient.Post(c.ApiURL("setMyCommands"), request, response); err != nil {
		return err
	}

	if !response.OK {
		return response
	}

	return nil
}

func (c *Client) GetUpdates(offset int, limit int, timeout int, allowedUpdates []string) ([]Update, error) {

	request := &UpdatesRequest{
		Offset:         offset,
		Limit:          limit,
		Timeout:        timeout,
		AllowedUpdates: allowedUpdates,
	}
	response := &UpdatesResponse{}
	err := jsonclient.Post(c.ApiURL("getUpdates"), request, response)

	if !response.OK {
		return nil, GetErrorResponse(err.Body)
	}

	return response.Result, nil
}

func (c *Client) ParseCommand(message *Message, botName string, private bool) (command string, params []string, err error) {
	parts := strings.Split(message.Text, " ")
	if len(parts) < 1 {
		return "", nil, fmt.Errorf("message is empty")
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
		return "", nil, fmt.Errorf("message is not a command")
	}

	return command[1:], params, nil
}

func (c *Client) SendMessage(message *SendMessagePayload) (*Message, error) {

	response := &MessageResponse{}
	err := jsonclient.Post(c.ApiURL("sendMessage"), message, response)

	if !response.OK {
		return nil, GetErrorResponse(err.Body)
	}

	return response.Result, nil
}

func (c *Client) UpdateMessage(update *UpdateMessagePayload) error {
	response := &ErrorResponse{}
	if err := jsonclient.Post(c.ApiURL("editMessageText"), update, response); err != nil {
		return err
	}

	if !response.OK {
		return response
	}

	return nil
}

func (c *Client) Reply(original *Message, text string) (*Message, error) {
	return c.SendMessage(&SendMessagePayload{
		Text:      text,
		ID:        strconv.FormatInt(original.Chat.ID, 10),
		ParseMode: ParseModes.Markdown.String(),
		ReplyTo:   original.MessageID,
	})
}

func (c *Client) AnswerCallbackQuery(answer *CallbackQueryAnswer) error {
	response := &ErrorResponse{}
	if err := jsonclient.Post(c.ApiURL("answerCallbackQuery"), answer, response); err != nil {
		return err
	}

	if !response.OK {
		return response
	}

	return nil
}

func GetErrorResponse(body string) error {
	response := &ErrorResponse{}
	err := json.Unmarshal([]byte(body), response)
	if err == nil {
		err = fmt.Errorf("telegram API error: %v", response)
	}
	return err
}
