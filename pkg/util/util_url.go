package util

import "net/url"

// UrlUserPassword is a replacement/wrapper around url.UserPassword that treats empty string arguments as not specified
// If no user or password is specified, it returns nil (which serializes in url.URL to "")
func UrlUserPassword(user, password string) *url.Userinfo {
	if len(password) > 0 {
		return url.UserPassword(user, password)
	} else if len(user) > 0 {
		return url.User(user)
	}
	return nil
}
