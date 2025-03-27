//go:generate stringer -type=URLPart -trimprefix URL

package format

import (
	"log"
	"strconv"
	"strings"
)

// URLPart is an indicator as to what part of an URL a field is serialized to.
type URLPart int

// Suffix returns the separator between the URLPart and it's subsequent part.
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

// indicator as to what part of an URL a field is serialized to.
const (
	URLQuery URLPart = iota
	URLUser
	URLPassword
	URLHost
	URLPort
	URLPath // Base path; additional paths are URLPath + N
)

// ParseURLPart returns the URLPart that matches the supplied string.
func ParseURLPart(s string) URLPart {
	s = strings.ToLower(s)
	switch s {
	case "user":
		return URLUser
	case "pass", "password":
		return URLPassword
	case "host":
		return URLHost
	case "port":
		return URLPort
	case "path", "path1":
		return URLPath
	case "query", "":
		return URLQuery
	}

	// Handle dynamic path segments (e.g., "path2", "path3", etc.).
	if strings.HasPrefix(s, "path") && len(s) > 4 {
		if num, err := strconv.Atoi(s[4:]); err == nil && num >= 2 {
			return URLPath + URLPart(num-1) // Offset from URLPath; "path2" -> URLPath+1
		}
	}

	log.Printf("invalid URLPart: %s, defaulting to URLQuery", s)

	return URLQuery
}

// ParseURLParts returns the URLParts that matches the supplied string.
func ParseURLParts(s string) []URLPart {
	rawParts := strings.Split(s, ",")
	urlParts := make([]URLPart, len(rawParts))

	for i, raw := range rawParts {
		urlParts[i] = ParseURLPart(raw)
	}

	return urlParts
}
