package pushbullet

import (
	"regexp"
)

// PushRequest ...
type PushRequest struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`

	Email      string `json:"email"`
	ChannelTag string `json:"channel_tag"`
	DeviceIden string `json:"device_iden"`
}

type PushResponse struct {
	Active                  bool    `json:"active"`
	Body                    string  `json:"body"`
	Created                 float64 `json:"created"`
	Direction               string  `json:"direction"`
	Dismissed               bool    `json:"dismissed"`
	Iden                    string  `json:"iden"`
	Modified                float64 `json:"modified"`
	ReceiverEmail           string  `json:"receiver_email"`
	ReceiverEmailNormalized string  `json:"receiver_email_normalized"`
	ReceiverIden            string  `json:"receiver_iden"`
	SenderEmail             string  `json:"sender_email"`
	SenderEmailNormalized   string  `json:"sender_email_normalized"`
	SenderIden              string  `json:"sender_iden"`
	SenderName              string  `json:"sender_name"`
	Title                   string  `json:"title"`
	Type                    string  `json:"type"`
}

type ErrorResponse struct {
	Error struct {
		Cat     string `json:"cat"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

var emailPattern = regexp.MustCompile(`.*@.*\..*`)

func (p *PushRequest) SetTarget(target string) {
	if emailPattern.MatchString(target) {
		p.Email = target
		return
	}

	if len(target) > 0 && string(target[0]) == "#" {
		p.ChannelTag = target[1:]
		return
	}

	p.DeviceIden = target
}

// NewNotePush creates a new push request
func NewNotePush(message, title string) *PushRequest {
	return &PushRequest{
		Type:  "note",
		Title: title,
		Body:  message,
	}
}
