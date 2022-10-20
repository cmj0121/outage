package status

import (
	"time"

	"github.com/rs/zerolog/log"
)

func (config *Config) Fetch(closed <-chan struct{}) (err error) {
	ticker := time.NewTicker(time.Duration(config.Interval))
	defer ticker.Stop()

	for {
		select {
		case <-closed:
			return
		case <-ticker.C:
			log.Trace().Msg("tick")
		}
	}
}
