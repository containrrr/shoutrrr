//+build xmpp

package router

func init() {
	serviceMap["xmpp"] = func() t.Service { return &xmpp.Service{} }
}
