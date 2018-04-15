package grimoire

import (
	"log"
	"time"
)

// Logger defines function signature for custom logger.
type Logger func(string, time.Duration, error)

// logger is default logger used to log query. It's exported for adapter testing purpose.
func logger(query string, duration time.Duration, err error) {
	if err != nil {
		log.Print("[", duration, "] ", err, " - ", query)
	} else {
		log.Print("[", duration, "] OK - ", query)
	}
}
