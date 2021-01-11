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
	userState   UserState
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

	resLogin := APIResLogin{}
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

	response := APIResLoginPassword{}
	if err := c.apiPost(APILogin, APIReqLoginPassword{
		Type:     FlowLoginPassword,
		User:     user,
		Password: password,
	}, &response); err != nil {
		return err
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
	// Update State
	if err := c.SyncState(); err != nil {
		errors = append(errors, fmt.Errorf("failed to update state: %v", err))
	}

	if len(rooms) > 0 {
		return c.sendToExplicitRooms(rooms, message, errors)
	} else {
		// Send to all rooms that are joined
		for roomID := range c.userState.Rooms.Join {
			c.logf("Sending message to '%v'...\n", roomID)
			if err := c.SendMessageToRoom(message, roomID); err != nil {
				errors = append(errors, fmt.Errorf("failed to send message to room '%v': %v", roomID, err))
			}
		}
	}
	return errors
}

func (c *Client) sendToExplicitRooms(rooms []string, message string, errors []error) []error {
	var err error
	for _, room := range rooms {
		c.logf("Sending message to '%v'...\n", room)

		var roomID string
		if roomID, err = c.GetRoomID(room); err != nil {
			errors = append(errors, fmt.Errorf("failed to look up alias '%v': %v", room, err))
			continue
		}

		c.logf("Resolved room alias '%v' to ID '%v'", room, roomID)

		if _, found := c.userState.Rooms.Join[roomID]; !found {

			c.logf("User has not joined room '%v'", roomID)

			if _, found := c.userState.Rooms.Invite[roomID]; !found {
				if err := c.JoinRoom(roomID); err != nil {
					errors = append(errors, fmt.Errorf("error joining room %v: %v", roomID, err))
					continue
				}
				c.logf("Joined room '%v'", roomID)
			} else {
				errors = append(errors, fmt.Errorf("user was not invited to '%v'", roomID))
				continue
			}
		}

		if err := c.SendMessageToRoom(message, roomID); err != nil {
			errors = append(errors, fmt.Errorf("failed to send message to room '%v': %v", roomID, err))
		}
	}
	return errors
}

func (c *Client) SyncState() error {
	c.logf("Syncing state...")

	userState := UserState{}
	if err := c.apiGet(APISync, &userState); err != nil {
		return err
	}
	c.userState = userState

	return nil
}

func (c *Client) JoinRoom(roomID string) error {
	resRoom := APIResRoom{}
	if err := c.apiPost(fmt.Sprintf(APIJoinInvite, roomID), nil, resRoom); err != nil {
		return err
	}

	return nil
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

func (c *Client) GetRoomID(room string) (roomID string, err error) {
	// If room begins with an '!', its already a room ID
	if room[0] == '!' {
		return room, nil
	}

	// If room does not begin with a '#' let's prepend it
	if room[0] != '#' {
		room = "#" + room
	}

	resRoom := APIResRoom{}
	if err = c.apiGet(fmt.Sprintf(APILookupRoom, room), &resRoom); err == nil {
		roomID = resRoom.RoomID
	}
	return roomID, err
}

func (c *Client) apiGet(path string, response interface{}) error {
	c.apiURL.Path = path

	var err error
	var res *http.Response
	res, err = http.Get(c.apiURL.String())
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("got HTTP %v", res.Status)
	}

	defer res.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(res.Body)
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

	if res.StatusCode >= 400 {
		return fmt.Errorf("got HTTP %v", res.Status)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
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
