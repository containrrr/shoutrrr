package smtp

import (
	"errors"
	"fmt"
	gosmtp "net/smtp"
	"strings"
)

type loginAuth struct {
	identity, username, password string
	host                         string
}

// Next implements smtp.Auth
func (a *loginAuth) Next(fromServer []byte, more bool) (toServer []byte, err error) {
	if !more {
		return nil, nil
	}
	msg := strings.Trim(strings.ToLower(string(fromServer)), ":")
	switch msg {
	case "user", "username":
		return []byte(a.username), nil
	case "pass", "password":
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("unknown challenge %q received", string(fromServer))
	}
}

// Start implements smtp.Auth
func (a *loginAuth) Start(server *gosmtp.ServerInfo) (proto string, toServer []byte, err error) {
	// Must have TLS, or else localhost server.
	// Note: If TLS is not true, then we can't trust ANYTHING in ServerInfo.
	// In particular, it doesn't matter if the server advertises LOGIN auth.
	// That might just be the attacker saying
	// "it's ok, you can trust me with your password."
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("unencrypted connection")
	}
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	return "LOGIN", nil, nil
}

// LoginAuth returns an Auth that implements the LOGIN authentication
func LoginAuth(identity, username, password, host string) gosmtp.Auth {
	return &loginAuth{identity, username, password, host}
}

func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}
