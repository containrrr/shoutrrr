package lark

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/services/standard"
	"github.com/containrrr/shoutrrr/pkg/types"
)

const (
	apiFormat   = "https://%s/open-apis/bot/v2/hook/%s"
	maxLength   = 4096
	defaultTime = 30 * time.Second
)

type Service struct {
	standard.Standard
	config *Config
	pkr    format.PropKeyResolver
}

var (
	ErrInvalidHost  = errors.New("invalid host, use 'open.larksuite.com' or 'open.feishu.cn'")
	ErrNoPath       = errors.New("no path, path like 'xxx' in 'https://open.larksuite.com/open-apis/bot/v2/hook/xxx'")
	ErrLargeMessage = errors.New("message exceeds the max length")

	httpClient = &http.Client{Timeout: defaultTime}
)

const (
	larkHost   = "open.larksuite.com"
	feishuHost = "open.feishu.com"
)

// Send notification to Lark
func (service *Service) Send(message string, params *types.Params) error {
	if len(message) > maxLength {
		return ErrLargeMessage
	}

	config := *service.config
	if err := service.pkr.UpdateConfigFromParams(&config, params); err != nil {
		return err
	}

	if config.Host != larkHost && config.Host != feishuHost {
		return ErrInvalidHost
	}

	if config.Path == "" {
		return ErrNoPath
	}

	return service.sendMessage(message, config)
}

// Initialize loads ServiceConfig from configURL and sets logger for this Service
func (service *Service) Initialize(configURL *url.URL, logger types.StdLogger) error {
	service.Logger.SetLogger(logger)
	service.config = &Config{}
	service.pkr = format.NewPropKeyResolver(service.config)
	if err := service.config.setURL(&service.pkr, configURL); err != nil {
		return err
	}
	return nil
}

func (service *Service) genSign(secret string, timestamp int64) string {
	//timestamp + key calculate sha256, then base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	h.Write(data)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature
}

func (service *Service) sendMessage(message string, cfg Config) error {
	url := fmt.Sprintf(apiFormat, cfg.Host, cfg.Path)
	body := service.getRequestBody(message, cfg.Title, cfg.Secret)
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	service.Logf("Lark Request Body: %s", string(data))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := Response{}
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	if response.Code == 0 {
		return nil
	}
	return fmt.Errorf("code: %d, msg: %s", response.Code, response.Msg)
}

func (service *Service) getRequestBody(message, title, secret string) *RequestBody {
	body := &RequestBody{}
	if secret != "" {
		ts := time.Now().Unix()
		body.Timestamp = strconv.FormatInt(ts, 10)
		body.Sign = service.genSign(secret, ts)
	}
	if title == "" {
		body.MsgType = MsgTypeText
		body.Content.Text = message
		return body
	}
	body.MsgType = MsgTypePost
	body.Content.Post = &Post{
		En: &Message{
			Title: title,
			Content: [][]Item{{
				{Tag: TagValueText, Text: message},
			}},
		},
	}
	return body
}
