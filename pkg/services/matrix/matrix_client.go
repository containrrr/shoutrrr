package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	apiURL      url.URL
	accessToken string
	logger      *log.Logger
}

func NewClient(host string, disableTLS bool, logger *log.Logger) (c *Client) {
	c = &Client{
		logger: logger,
		apiURL: url.URL{
			Host:   host,
			Scheme: "https",
		},
	}

	if disableTLS {
		c.apiURL.Scheme = c.apiURL.Scheme[:4]
	}

	logger.Printf("Using server: %v\n", c.apiURL.String())

	return c
}

func (c *Client) UseToken(token string) {
	c.accessToken = token
	c.updateAccessToken()
}

func (c *Client) Login(user string, password string) error {
	c.apiURL.RawQuery = ""
	defer c.updateAccessToken()

	resLogin := APIResLoginFlows{}
	if err := c.apiGet(APILogin, &resLogin); err != nil {
		return fmt.Errorf("failed to get login flows: %v", err)
	}

	var flows []string
	for _, flow := range resLogin.Flows {
		flows = append(flows, string(flow.Type))
		if flow.Type == FlowLoginPassword {
			c.logf("Using login flow '%v'", flow.Type)
			return c.LoginPassword(user, password)
		}
	}

	return fmt.Errorf("none of the server login flows are supported: %v", strings.Join(flows, ", "))
}

func (c *Client) LoginPassword(user string, password string) error {

	response := APIResLogin{}
	if err := c.apiPost(APILogin, APIReqLogin{
		Type:       FlowLoginPassword,
		Password:   password,
		Identifier: NewUserIdentifier(user),
	}, &response); err != nil {
		return fmt.Errorf("failed to log in: %v", err)
	}

	c.accessToken = response.AccessToken
	tokenHint := ""
	if len(response.AccessToken) > 3 {
		tokenHint = response.AccessToken[:3]
	}

	c.logf("AccessToken: %v...\n", tokenHint)
	c.logf("HomeServer: %v\n", response.HomeServer)
	c.logf("User: %v\n", response.UserID)

	return nil
}

func (c *Client) SendMessage(message string, rooms []string) (errors []error) {
	if len(rooms) > 0 {
		return c.sendToExplicitRooms(rooms, message)
	} else {
		return c.sendToJoinedRooms(message)
	}
}

func (c *Client) sendToExplicitRooms(rooms []string, message string) (errors []error) {
	var err error

	for _, room := range rooms {
		c.logf("Sending message to '%v'...\n", room)

		var roomID string
		if roomID, err = c.JoinRoom(room); err != nil {
			errors = append(errors, fmt.Errorf("error joining room %v: %v", roomID, err))
			continue
		}

		if room != roomID {
			c.logf("Resolved room alias '%v' to ID '%v'", room, roomID)
		}

		if err := c.SendMessageToRoom(message, roomID); err != nil {
			errors = append(errors, fmt.Errorf("failed to send message to room '%v': %v", roomID, err))
		}
	}

	return errors
}

func (c *Client) sendToJoinedRooms(message string) (errors []error) {
	joinedRooms, err := c.GetJoinedRooms()
	if err != nil {
		return append(errors, fmt.Errorf("failed to get joined rooms: %v", err))
	}

	// Send to all rooms that are joined
	for _, roomID := range joinedRooms {
		c.logf("Sending message to '%v'...\n", roomID)
		if err := c.SendMessageToRoom(message, roomID); err != nil {
			errors = append(errors, fmt.Errorf("failed to send message to room '%v': %v", roomID, err))
		}
	}

	return errors
}

func (c *Client) JoinRoom(room string) (roomID string, err error) {
	resRoom := APIResRoom{}
	if err = c.apiPost(fmt.Sprintf(APIRoomJoin, room), nil, &resRoom); err != nil {
		return "", err
	}
	return resRoom.RoomID, nil
}

func (c *Client) SendMessageToRoom(message string, roomID string) error {
	resEvent := APIResEvent{}
	if err := c.apiPost(fmt.Sprintf(APISendMessage, roomID), APIReqSend{
		MsgType: MsgTypeText,
		Body:    message,
	}, &resEvent); err != nil {
		return err
	}

	return nil
}

func (c *Client) apiGet(path string, response interface{}) error {
	c.apiURL.Path = path

	var err error
	var res *http.Response
	res, err = http.Get(c.apiURL.String())
	if err != nil {
		return err
	}

	var body []byte
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		resError := &APIResError{}
		if err == nil {
			if err = json.Unmarshal(body, resError); err == nil {
				return resError
			}
		}

		return fmt.Errorf("got HTTP %v", res.Status)
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

func (c *Client) apiPost(path string, request interface{}, response interface{}) error {
	c.apiURL.Path = path

	var err error
	var body []byte

	body, err = json.Marshal(request)
	if err != nil {
		return err
	}

	var res *http.Response
	res, err = http.Post(c.apiURL.String(), ContentType, bytes.NewReader(body))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		resError := &APIResError{}
		if err == nil {
			if err = json.Unmarshal(body, resError); err == nil {
				return resError
			}
		}

		return fmt.Errorf("got HTTP %v", res.Status)
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

func (c *Client) updateAccessToken() {
	query := c.apiURL.Query()
	query.Set(AccessTokenKey, c.accessToken)
	c.apiURL.RawQuery = query.Encode()
}

func (c *Client) logf(format string, v ...interface{}) {
	c.logger.Printf(format, v...)
}

func (c *Client) GetJoinedRooms() ([]string, error) {
	response := APIResJoinedRooms{}
	if err := c.apiGet(APIJoinedRooms, &response); err != nil {
		return []string{}, err
	}
	return response.Rooms, nil
}
