package logging

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/multi"
	"github.com/apex/log/handlers/text"
)

var logLevel = log.InfoLevel
var start = time.Now()
var timestamp = false

// InitLogging sets up the logging according to the passed config
func InitLogging(debug bool) {

	var invalidLoggers []string

	handler := &multi.Handler{
		Handlers: []log.Handler{
			&LogHandler{
				Writer: os.Stderr,
			},
		},
	}

	log.SetHandler(handler)

	if len(handler.Handlers) < 1 {
		handler.Handlers = append(handler.Handlers, &LogHandler{
			Writer: os.Stderr,
		})
		log.Warn("no valid loggers configured, using default")
	}

	for _, loggerId := range invalidLoggers {
		log.WithField("logger", loggerId).Warn("invalid logger in configuration")
	}

	if debug {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

	if err := initPlatform(); err != nil {
		log.Errorf("failed to initialize platform: %v", err)
	}
}

// GetLogger returns a log.Interface with the context set to the specified value
func GetLogger(context string) log.Interface {
	return log.WithField("context", context)
}

// LogHandler implements log.Handler
type LogHandler struct {
	mu     sync.Mutex
	Writer io.Writer
}

type Fields log.Fields

// HandleLog implements log.Handler.HandleLog
func (lh *LogHandler) HandleLog(e *log.Entry) error {
	color := text.Colors[e.Level]
	if e.Level == log.DebugLevel {
		color = 36
	}
	level := text.Strings[e.Level][:3]
	names := e.Fields.Names()

	lh.mu.Lock()
	defer lh.mu.Unlock()

	ts := time.Since(start) / time.Second

	message := e.Message
	if context, ok := e.Fields.Fields()["context"]; ok {
		message = context.(string) + ": " + message
	}

	time := ""
	if timestamp {
		time = fmt.Sprintf("%04d", ts)
	}
	fmt.Fprintf(lh.Writer, "%s[\033[%dm%3s\033[0m] %-40s", time, color, level, message)
	// fmt.Fprintf(lh.Writer, "\033[%dm%1s \033[0m[%04d] %-35s", color, level, ts, e.Message)

	for _, name := range names {
		if name != "context" {
			fmt.Fprintf(lh.Writer, " \033[%dm%s\033[0m=%v", color, name, e.Fields.Get(name))
		}
	}

	fmt.Fprintln(lh.Writer)

	return nil
}
