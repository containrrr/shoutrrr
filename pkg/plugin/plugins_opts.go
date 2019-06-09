package plugin

import (
	"io/ioutil"
	"log"
)

type PluginOpts struct {
	Verbose bool
	Logger *log.Logger
	Props map[string]string
}

func PluginOptsDefault(logger log.Logger) PluginOpts {
	return PluginOpts{
		Verbose: false,
		Logger: DiscardLogger,
	}
}

var DiscardLogger = log.New(ioutil.Discard, "", 0)