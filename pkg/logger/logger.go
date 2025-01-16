// logger.go — post-passport.
// Author: d28035203

package logger

import (
	"log"
	"os"
)

// Init initializes the logger with the specified log level
func Init(level string) {
	switch level {
	case "DEBUG":
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.SetOutput(os.Stdout)
	case "INFO":
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stdout)
	case "WARN":
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stdout)
	case "ERROR":
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)
	default:
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stdout)
	}
}
