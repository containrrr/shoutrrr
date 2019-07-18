package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var actionWords = [...]string {"send", "verify", "generate"}

// ExitCodeUsage is used to signify that the command was not properly invoked
const ExitCodeUsage = 64
const mainUsage = "%s <ActionVerb> [...]\n"

func usage(syntax string) {
	fmt.Println("Usage:")
	fmt.Printf(syntax, os.Args[0])
}

func main() {

	if len(os.Args) < 2 {
		showMainUsage()
		return
	}

	actionResult := 1025
	actionWord := os.Args[1]
	var action action

	switch actionWord {
	case "send":
		action = send()
	case "verify":
		action = verify()
	case "generate":
		action = generate()
	default:
		showMainUsage()
		return
	}

	actionResult = action.Run()

	if parseErr := action.FlagSet.Parse(os.Args[2:]); parseErr == nil {
		actionResult = action.Run()
	} else {
		fmt.Print(parseErr)
	}

	if actionResult == ExitCodeUsage {
		usage(action.Usage)
		fmt.Println("\nOPTIONS:")
		action.FlagSet.PrintDefaults()
	}

	os.Exit(actionResult)
}

func showMainUsage() {
	usage(mainUsage)
	fmt.Printf("Possible actions: %s\n", strings.Join(actionWords[:	], ", "))
	os.Exit(ExitCodeUsage)
}


type action struct {
	run func(flags *flag.FlagSet) int
	FlagSet flag.FlagSet
	Usage string
}

func (a *action) Run() int {
	return a.run(&a.FlagSet)
}