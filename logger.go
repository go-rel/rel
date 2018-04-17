package grimoire

import (
	"log"
	"time"
)

// Logger defines function signature for custom logger.
type Logger func(string, time.Duration, error)

// DefaultLogger log query suing standard log library.
func DefaultLogger(query string, duration time.Duration, err error) {
	if err != nil {
		log.Print("[", duration, "] ", err, " - ", query)
	} else {
		log.Print("[", duration, "] OK - ", query)
	}
}

// Log using multiple logger.
// This function intended to be used within adapter.
func Log(logger []Logger, statement string, duration time.Duration, err error) {
	for _, l := range logger {
		l(statement, duration, err)
	}
}
