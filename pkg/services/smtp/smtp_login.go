package smtp

import (
	"errors"
	"net/smtp"
	"strings"
)

type loginAuth struct {
	username, password, host string
}

// PlainAuth returns an Auth that implements the PLAIN authentication
// mechanism as defined in RFC 4616. The returned Auth uses the given
// username and password to authenticate to host and act as identity.
// Usually identity should be the empty string, to act as username.
//
// PlainAuth will only send the credentials if the connection is using TLS
// or is connected to localhost. Otherwise authentication will fail with an
// error, without sending the credentials.
func LoginAuth(username, password, host string) smtp.Auth {
	return &loginAuth{username, password, host}
}

func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// Must have TLS, or else localhost server.
	// Note: If TLS is not true, then we can't trust ANYTHING in ServerInfo.
	// In particular, it doesn't matter if the server advertises PLAIN auth.
	// That might just be the attacker saying
	// "it's ok, you can trust me with your password."
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("unencrypted connection")
	}
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.username)
	return "LOGIN", resp, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		challenge := string(fromServer)
		if strings.ToLower(challenge) != "password" {
			return nil, errors.New(`unexpected server challenge "` + challenge + `"`)
		}
		return []byte(a.password), nil

	}
	return nil, nil
}
