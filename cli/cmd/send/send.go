package send

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	cli "github.com/containrrr/shoutrrr/cli/cmd"
	u "github.com/containrrr/shoutrrr/internal/util"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
)

// Cmd sends a notification using a service url
var Cmd = &cobra.Command{
	Use:    "send",
	Short:  "Send a notification using a service url",
	Args:   cobra.MaximumNArgs(2),
	PreRun: u.LoadFlagsFromAltSources,
	Run:    Run,
}

func init() {
	Cmd.Flags().BoolP("verbose", "v", false, "")

	Cmd.Flags().StringArrayP("url", "u", []string{}, "The notification url")
	_ = Cmd.MarkFlagRequired("url")

	Cmd.Flags().StringP("message", "m", "", "The message to send to the notification url")
	_ = Cmd.MarkFlagRequired("message")

	Cmd.Flags().StringP("title", "t", "", "The title used for services that support it")
}

// Run the send command
func Run(cmd *cobra.Command, _ []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")

	urls, _ := cmd.Flags().GetStringArray("url")
	message, _ := cmd.Flags().GetString("message")
	title, _ := cmd.Flags().GetString("title")

	var logger *log.Logger

	if verbose {
		fmt.Println("URLs:")
		for _, url := range urls {
			fmt.Printf("  %s\n", url)
		}
		fmt.Printf("Message: %s\n", message)
		if title != "" {
			fmt.Printf("Title: %v\n", title)
		}
		logger = log.New(os.Stderr, "SHOUTRRR ", log.LstdFlags)
	} else {
		logger = util.DiscardLogger
	}

	exitCode := cli.ExSuccess

	sr, err := router.New(logger, urls...)
	if err != nil {
		fmt.Printf("error invoking send: %s\n", err)
		exitCode = cli.ExConfig
	} else {
		params := make(types.Params)
		if title != "" {
			params["title"] = title
		}
		errs := sr.SendAsync(message, &params)
		for err := range errs {
			if err == nil {
				fmt.Println("Notification sent")
			} else {
				fmt.Printf("Error: %v\n", err)
				exitCode = cli.ExUnavailable
			}
		}
	}

	os.Exit(exitCode)

}
