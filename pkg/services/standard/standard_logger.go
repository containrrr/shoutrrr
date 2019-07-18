package standard

import (
	"log"

	"github.com/containrrr/shoutrrr/pkg/util"
)

// Logger provides the utility methods Log* that maps to Logger.Print*
type Logger struct {
	logger *log.Logger
}

// Logf maps to the service loggers Logger.Printf function
func (sl *Logger) Logf(format string, v ...interface{}) {
	sl.logger.Printf(format, v...)
}

// Log maps to the service loggers Logger.Print function
func (sl *Logger) Log(v ...interface{}) {
	sl.logger.Print(v...)
}

// SetLogger maps the specified logger to the Log* helper methods
func (sl *Logger) SetLogger(logger *log.Logger) {
	if logger == nil {
		sl.logger = util.DiscardLogger
	} else {
		sl.logger = logger
	}
}
