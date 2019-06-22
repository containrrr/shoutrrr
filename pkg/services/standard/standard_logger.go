package standard

import (
	"log"
)

// Logger is the standard implementation of SetLogger and provides the utility methods Log* that maps to Logger.Print*
type Logger struct {
	logger *log.Logger
}



// Logf maps to the service loggers Logger.Printf function
func (sl *Logger) Logf(format string, v ...interface{}) {
	sl.logger.Printf(format, v...)
}

// Logln maps to the service loggers Logger.Println function
func (sl *Logger) Logln(v ...interface{}) {
	sl.logger.Println(v...)
}

// Log maps to the service loggers Logger.Print function
func (sl *Logger) Log(v ...interface{}) {
	sl.logger.Print(v...)
}