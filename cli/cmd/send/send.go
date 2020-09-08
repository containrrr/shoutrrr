package send

import (
	"fmt"
	"log"
	"os"

	"github.com/containrrr/shoutrrr"
	u "github.com/containrrr/shoutrrr/internal/util"
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
	Cmd.Flags().StringP("url", "u", "", "The notification url")
	_ = Cmd.MarkFlagRequired("url")

	Cmd.Flags().StringP("message", "m", "", "The message to send to the notification url")
	_ = Cmd.MarkFlagRequired("message")
}

// Run the send command
func Run(cmd *cobra.Command, _ []string) {
	debug, _ := cmd.Flags().GetBool("debug")

	url, _ := cmd.Flags().GetString("url")
	message, _ := cmd.Flags().GetString("message")

	var logger *log.Logger

	if debug {
		fmt.Printf("URL: %s\n", url)
		fmt.Printf("Message: %s\n", message)
		logger = log.New(os.Stderr, "SHOUTRRR ", log.LstdFlags)
	} else {
		logger = util.DiscardLogger
	}

	shoutrrr.SetLogger(logger)
	err := shoutrrr.Send(url, message)

	if err != nil {
		fmt.Printf("error invoking send: %s", err)
		os.Exit(1)
	}
}
