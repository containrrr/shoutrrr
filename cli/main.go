package main

import (
	"github.com/containrrr/shoutrrr/cli/cmd/root"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := root.Cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
