package main

import (
	"fmt"
	"os"
 	"github.com/containrrr/shoutrrr"
	"strings"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Printf("%s <Message>\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
		return
	}

	message := strings.Join(os.Args[1:], " ")

	shoutrrr.SendEnv(message)
}
