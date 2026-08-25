package smtp

import (
	"errors"
	"net/smtp"
	"strings"
)

type loginAuth struct {
	username, password string
}

// LoginAuth returns an Auth that implements the LOGIN authentication mechanism.
// It is not part of the net/smtp package since LOGIN is not a formal standard,
// but it is still expected by some servers such as Microsoft Exchange and
// Office 365.
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	// The server sends the username and password prompts as separate
	// challenges. The wording is not standardized, so match on the prefix.
	switch prompt := strings.ToLower(strings.TrimSpace(string(fromServer))); {
	case strings.HasPrefix(prompt, "user"):
		return []byte(a.username), nil
	case strings.HasPrefix(prompt, "pass"):
		return []byte(a.password), nil
	default:
		return nil, errors.New("unexpected server challenge during LOGIN authentication")
	}
}
