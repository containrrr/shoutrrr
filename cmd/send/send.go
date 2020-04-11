package send

import (
	"fmt"
	"log"
	"os"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func hasURLInEnvButNotFlag(cmd *cobra.Command) bool {
	s, _ := cmd.Flags().GetString("url")
	return s == "" && viper.GetViper().GetString("url") != ""
}

// Cmd sends a notification using a service url
var Cmd = &cobra.Command{
	Use:   "send",
	Short: "Send a notification using a service url",
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		// WORKAROUND: make cobra count env vars when checking required flags
		if hasURLInEnvButNotFlag(cmd) {
			cmd.Flags().Set("url", viper.GetViper().GetString("url"))
		}
	},
	Run: Run,
}

func init() {
	Cmd.Flags().StringP("message", "m", "", "The message to send to the notification url")
	Cmd.MarkFlagRequired("message")
}

// Run the send command
func Run(cmd *cobra.Command, args []string) {
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
