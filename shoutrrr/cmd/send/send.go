package send

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	intutil "github.com/containrrr/shoutrrr/internal/util"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/containrrr/shoutrrr/pkg/util"
	cli "github.com/containrrr/shoutrrr/shoutrrr/cmd"
)

// Cmd sends a notification using a service url
var Cmd = &cobra.Command{
	Use:    "send",
	Short:  "Send a notification using a service url",
	Args:   cobra.MaximumNArgs(2),
	PreRun: intutil.LoadFlagsFromAltSources,
	RunE:   Run,
}

func init() {
	Cmd.Flags().BoolP("verbose", "v", false, "")

	Cmd.Flags().StringArrayP("url", "u", []string{}, "The notification url")
	_ = Cmd.MarkFlagRequired("url")

	Cmd.Flags().StringP("message", "m", "", "The message to send to the notification url")
	_ = Cmd.MarkFlagRequired("message")

	Cmd.Flags().StringP("title", "t", "", "The title used for services that support it")
}

func logf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func run(cmd *cobra.Command) error {
	flags := cmd.Flags()
	verbose, _ := flags.GetBool("verbose")

	urls, _ := flags.GetStringArray("url")
	message, _ := flags.GetString("message")
	title, _ := flags.GetString("title")

	if message == "-" {
		logf("Reading from STDIN...")
		msgBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read message from stdin: %v", err)
		}
		message = string(msgBytes)
	}

	var logger *log.Logger

	if verbose {
		urlsPrefix := "URLs:"
		for i, url := range urls {
			logf("%s %s", urlsPrefix, url)
			if i == 0 {
				// Only display "URLs:" prefix for first line, replace with indentation for the the subsequent
				urlsPrefix = strings.Repeat(" ", len(urlsPrefix))
			}
		}
		logf("Message: %s", util.Ellipsis(message, 100))
		if title != "" {
			logf("Title: %v", title)
		}
		logger = log.New(os.Stderr, "SHOUTRRR ", log.LstdFlags)
	} else {
		logger = util.DiscardLogger
	}

	sr, err := router.New(logger, urls...)
	if err != nil {
		return cli.ConfigurationError(fmt.Sprintf("error invoking send: %s", err))
	} else {
		params := make(types.Params)
		if title != "" {
			params["title"] = title
		}
		errs := sr.SendAsync(message, &params)
		for err := range errs {
			if err != nil {
				return cli.TaskUnavailable(err.Error())
			}
			logf("Notification sent")
		}
	}

	return nil
}

// Run the send command
func Run(cmd *cobra.Command, _ []string) error {
	err := run(cmd)
	if err != nil {
		if result, ok := err.(cli.Result); ok && result.ExitCode != cli.ExUsage {
			// If the error is not related to the CLI usage, report error and exit to not invoke cobra error output
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(result.ExitCode)
		}
	}
	return err
}
