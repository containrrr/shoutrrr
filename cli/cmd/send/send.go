package send

import (
	"fmt"
	"log"
	"os"

	cli "github.com/containrrr/shoutrrr/cli/cmd"
	u "github.com/containrrr/shoutrrr/internal/util"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/util"
	"github.com/spf13/cobra"
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

	Cmd.Flags().StringSliceP("url", "u", []string{}, "The notification url")
	_ = Cmd.MarkFlagRequired("url")

	Cmd.Flags().StringP("message", "m", "", "The message to send to the notification url")
	_ = Cmd.MarkFlagRequired("message")
}

// Run the send command
func Run(cmd *cobra.Command, _ []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")

	urls, _ := cmd.Flags().GetStringSlice("url")
	message, _ := cmd.Flags().GetString("message")

	var logger *log.Logger

	if verbose {
		fmt.Println("URLs:")
		for _, url := range urls {
			fmt.Printf("  %s\n", url)
		}
		fmt.Printf("Message: %s\n", message)
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
		errs := sr.SendAsync(message, nil)
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
