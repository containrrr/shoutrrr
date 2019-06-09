package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var actionWords = [...]string {"send", "verify", "generate"}

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
	var action Action

	switch actionWord {
	case "send":
		action = Send()
	case "verify":
		action = Verify()
	case "generate":
		action = Generate()
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


type Action struct {
	run func(flags *flag.FlagSet) int
	FlagSet flag.FlagSet
	Usage string
}

func (a *Action) Run() int {
	return a.run(&a.FlagSet)
}