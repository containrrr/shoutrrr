//go:build xmpp
// +build xmpp

package router

import t "github.com/nicholas-fedor/shoutrrr/pkg/types"

func init() {
	serviceMap["xmpp"] = func() t.Service { return &xmpp.Service{} }
}
