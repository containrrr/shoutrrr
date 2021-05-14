//go:generate stringer -type=URLPart -trimprefix URL

package format

import (
	"log"
	"strings"
)

// URLPart is an indicator as to what part of an URL a field is serialized to
type URLPart int

// Suffix returns the separator between the URLPart and it's subsequent part
func (u URLPart) Suffix() rune {
	switch u {
	case URLUser:
		return ':'
	case URLPassword:
		return '@'
	case URLHost:
		return ':'
	case URLPort:
		fallthrough
	case URLPath:
		fallthrough
	default:
		return '/'
	}
}

// indicator as to what part of an URL a field is serialized to
const (
	URLQuery URLPart = iota
	URLUser
	URLPassword
	URLHost
	URLPort
	URLPath
)

// ParseURLPart returns the URLPart that matches the supplied string
func ParseURLPart(s string) URLPart {
	switch strings.ToLower(s) {
	case "user":
		return URLUser
	case "pass":
		fallthrough
	case "password":
		return URLPassword
	case "host":
		return URLHost
	case "port":
		return URLPort
	case "path":
		fallthrough
	case "path1":
		return URLPath
	case "path2":
		return URLPath + 1
	case "path3":
		return URLPath + 2
	case "path4":
		return URLPath + 3
	case "query":
		fallthrough
	case "":
		return URLQuery
	default:
		log.Fatal("invalid URLPart")
		return URLQuery
	}
}

// ParseURLParts returns the URLParts that matches the supplied string
func ParseURLParts(s string) []URLPart {
	rawParts := strings.Split(s, ",")
	urlParts := make([]URLPart, len(rawParts))
	for i, raw := range rawParts {
		urlParts[i] = ParseURLPart(raw)
	}
	return urlParts
}
