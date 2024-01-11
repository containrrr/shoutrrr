package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

type client struct {
	apiURL      url.URL
	accessToken string
	logger      types.StdLogger
	counter     uint64
}

func newClient(host string, disableTLS bool, logger types.StdLogger) (c *client) {
	c = &client{
		logger: logger,
		apiURL: url.URL{
			Host:   host,
			Scheme: "https",
		},
	}

	if c.logger == nil {
		c.logger = util.DiscardLogger
	}

	if disableTLS {
		c.apiURL.Scheme = c.apiURL.Scheme[:4]
	}

	c.logger.Printf("Using server: %v\n", c.apiURL.String())

	return c
}

func (c *client) txId() uint64 {
	return atomic.AddUint64(&c.counter, 1)
}

func (c *client) useToken(token string) {
	c.accessToken = token
	c.updateAccessToken()
}

func (c *client) login(user string, password string) error {
	c.apiURL.RawQuery = ""
	defer c.updateAccessToken()

	resLogin := apiResLoginFlows{}
	if err := c.apiGet(apiLogin, &resLogin); err != nil {
		return fmt.Errorf("failed to get login flows: %w", err)
	}

	var flows []string
	for _, flow := range resLogin.Flows {
		flows = append(flows, string(flow.Type))
		if flow.Type == flowLoginPassword {
			c.logf("Using login flow '%v'", flow.Type)
			return c.loginPassword(user, password)
		}
	}

	return fmt.Errorf("none of the server login flows are supported: %v", strings.Join(flows, ", "))
}

func (c *client) loginPassword(user string, password string) error {
	response := apiResLogin{}
	if err := c.apiPost(apiLogin, apiReqLogin{
		Type:       flowLoginPassword,
		Password:   password,
		Identifier: newUserIdentifier(user),
	}, &response); err != nil {
		return fmt.Errorf("failed to log in: %w", err)
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

func (c *client) sendMessage(message string, rooms []string) (errors []error) {
	if len(rooms) > 0 {
		return c.sendToExplicitRooms(rooms, message)
	}
	return c.sendToJoinedRooms(message)
}

func (c *client) sendToExplicitRooms(rooms []string, message string) (errors []error) {
	var err error

	for _, room := range rooms {
		c.logf("Sending message to '%v'...\n", room)

		var roomID string
		if roomID, err = c.joinRoom(room); err != nil {
			errors = append(errors, fmt.Errorf("error joining room %v: %w", roomID, err))
			continue
		}

		if room != roomID {
			c.logf("Resolved room alias '%v' to ID '%v'", room, roomID)
		}

		if err := c.sendMessageToRoom(message, roomID); err != nil {
			errors = append(errors, fmt.Errorf("failed to send message to room '%v': %w", roomID, err))
		}
	}

	return errors
}

func (c *client) sendToJoinedRooms(message string) (errors []error) {
	joinedRooms, err := c.getJoinedRooms()
	if err != nil {
		return append(errors, fmt.Errorf("failed to get joined rooms: %w", err))
	}

	if len(joinedRooms) == 0 {
		return append(errors, fmt.Errorf("no rooms has been joined"))
	}

	// Send to all rooms that are joined
	for _, roomID := range joinedRooms {
		c.logf("Sending message to '%v'...\n", roomID)
		if err := c.sendMessageToRoom(message, roomID); err != nil {
			errors = append(errors, fmt.Errorf("failed to send message to room '%v': %w", roomID, err))
		}
	}

	return errors
}

func (c *client) joinRoom(room string) (roomID string, err error) {
	resRoom := apiResRoom{}
	if err = c.apiPost(fmt.Sprintf(apiRoomJoin, room), nil, &resRoom); err != nil {
		return "", err
	}
	return resRoom.RoomID, nil
}

func (c *client) sendMessageToRoom(message string, roomID string) error {
	resEvent := apiResEvent{}
	return c.apiPut(fmt.Sprintf(apiSendMessage, roomID, c.txId()), apiReqSend{
		MsgType: msgTypeText,
		Body:    message,
	}, &resEvent)
}

func (c *client) apiGet(path string, response interface{}) error {
	return c.apiReq(path, "GET", nil, response)
}

func (c *client) apiPost(path string, request interface{}, response interface{}) error {
	return c.apiReq(path, "POST", request, response)
}

func (c *client) apiPut(path string, request interface{}, response interface{}) error {
	return c.apiReq(path, "PUT", request, response)
}

func (c *client) apiReq(path string, method string, request interface{}, response interface{}) error {
	c.apiURL.Path = path

	var payload io.Reader = nil
	if request != nil {
		body, err := json.Marshal(request)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, c.apiURL.String(), payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		resError := &apiResError{}
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

func (c *client) updateAccessToken() {
	query := c.apiURL.Query()
	query.Set(accessTokenKey, c.accessToken)
	c.apiURL.RawQuery = query.Encode()
}

func (c *client) logf(format string, v ...interface{}) {
	c.logger.Printf(format, v...)
}

func (c *client) getJoinedRooms() ([]string, error) {
	response := apiResJoinedRooms{}
	if err := c.apiGet(apiJoinedRooms, &response); err != nil {
		return []string{}, err
	}
	return response.Rooms, nil
}
