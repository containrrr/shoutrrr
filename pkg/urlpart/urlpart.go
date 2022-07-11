package urlpart

import (
	"log"
	"strings"
)

// URLPart is an indicator as to what part of an URL a field is serialized to
type URLPart string

// Suffix returns the separator between the URLPart and it's subsequent part
func (u URLPart) Suffix() rune {
	switch u {
	case User:
		return ':'
	case Password:
		return '@'
	case Host:
		return ':'
	case Port:
		fallthrough
	case Path1:
		fallthrough
	case Path2:
		fallthrough
	case Path3:
		fallthrough
	case Path4:
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
	Scheme   URLPart = "scheme"
	Query    URLPart = "query"
	User     URLPart = "user"
	Password URLPart = "password"
	Host     URLPart = "host"
	Port     URLPart = "port"
	Path1    URLPart = "path1"
	Path2    URLPart = "path2"
	Path3    URLPart = "path3"
	Path4    URLPart = "path4"
	Path     URLPart = "path"
	Count            = 6
)

var Order [9]URLPart = [9]URLPart{
	Scheme,
	User,
	Password,
	Host,
	Port,
	Path1,
	Path2,
	Path3,
	Path4,
}

var PathParts [4]URLPart = [4]URLPart{
	Path1,
	Path2,
	Path3,
	Path4,
}

// ParseOne returns the URLPart that matches the supplied string
func ParseOne(s string) URLPart {
	switch strings.ToLower(s) {
	case "user":
		return User
	case "pass":
		fallthrough
	case "password":
		return Password
	case "host":
		return Host
	case "port":
		return Port
	case "path":
		return Path
	case "path1":
		return Path1
	case "path2":
		return Path2
	case "path3":
		return Path3
	case "path4":
		return Path4
	case "query":
		fallthrough
	case "":
		return Query
	default:
		log.Fatal("invalid URLPart")
		return Query
	}
}

// ParseAll returns the URLParts that matches the supplied string
func ParseAll(s string) []URLPart {
	rawParts := strings.Split(s, ",")
	urlParts := make([]URLPart, len(rawParts))
	for i, raw := range rawParts {
		urlParts[i] = ParseOne(raw)
	}
	return urlParts
}
