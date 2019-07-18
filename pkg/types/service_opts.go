package types

import "log"

// ServiceOpts is the interface describing the service options
type ServiceOpts interface {
	Verbose() bool
	Logger() *log.Logger
	Props() map[string]string
}
