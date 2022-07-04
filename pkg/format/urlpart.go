package format

import (
	"log"
	"strings"
)

// URLPart is an indicator as to what part of an URL a field is serialized to
type URLPart string

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
	case URLPath1:
		fallthrough
	case URLPath2:
		fallthrough
	case URLPath3:
		fallthrough
	case URLPath4:
		fallthrough
	default:
		return '/'
	}
}

func (u URLPart) IsPath() bool {
	return strings.HasPrefix(string(u), "path")
}

// indicator as to what part of an URL a field is serialized to
const (
	URLScheme    URLPart = "scheme"
	URLQuery     URLPart = "query"
	URLUser      URLPart = "user"
	URLPassword  URLPart = "password"
	URLHost      URLPart = "host"
	URLPort      URLPart = "port"
	URLPath1     URLPart = "path1"
	URLPath2     URLPart = "path2"
	URLPath3     URLPart = "path3"
	URLPath4     URLPart = "path4"
	URLPath      URLPart = "path"
	URLPartCount         = 6
)

var URLPartOrder [9]URLPart = [9]URLPart{
	URLScheme,
	URLUser,
	URLPassword,
	URLHost,
	URLPort,
	URLPath1,
	URLPath2,
	URLPath3,
	URLPath4,
}

var URLPathParts [4]URLPart = [4]URLPart{
	URLPath1,
	URLPath2,
	URLPath3,
	URLPath4,
}

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
		return URLPath
	case "path1":
		return URLPath1
	case "path2":
		return URLPath2
	case "path3":
		return URLPath3
	case "path4":
		return URLPath4
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
