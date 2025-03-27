package standard

import (
	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/nicholas-fedor/shoutrrr/pkg/util"
)

// Logger provides the utility methods Log* that maps to Logger.Print*.
type Logger struct {
	logger types.StdLogger
}

// Logf maps to the service loggers Logger.Printf function.
func (sl *Logger) Logf(format string, v ...any) {
	sl.logger.Printf(format, v...)
}

// Log maps to the service loggers Logger.Print function.
func (sl *Logger) Log(v ...any) {
	sl.logger.Print(v...)
}

// SetLogger maps the specified logger to the Log* helper methods.
func (sl *Logger) SetLogger(logger types.StdLogger) {
	if logger == nil {
		sl.logger = util.DiscardLogger
	} else {
		sl.logger = logger
	}
}
