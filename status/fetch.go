package status

import (
	"time"

	"github.com/rs/zerolog/log"
)

func (config *Config) Fetch(closed <-chan struct{}) (err error) {
	ticker := time.NewTicker(time.Duration(config.Interval))
	defer ticker.Stop()

	ch := make(chan *Service, 1)
	defer close(ch)
	go config.worker(closed, ch)

	for {
		select {
		case <-closed:
			return
		case <-ticker.C:
			log.Debug().Msg("tick")

			now := time.Now().UTC()
			for _, service := range config.Services {
				switch {
				case service.UpdatedAt.IsZero():
					log.Debug().Str("service", service.Title).Msg("first blood")
					ch <- service
				case now.After(service.UpdatedAt.Add(time.Duration(10 * time.Second))):
					log.Debug().Str("service", service.Title).Msg("need update")
					ch <- service
				default:
					log.Trace().Str("service", service.Title).Str("UpdatedAt", service.UpdatedAt.String()).Str("now", now.String()).Msg("check")
				}
			}
		}
	}
}

func (config *Config) worker(closed <-chan struct{}, c <-chan *Service) {
	for {
		select {
		case <-closed:
			return
		case service := <-c:
			log.Info().Str("service", service.Title).Msg("process service")

			now := time.Now().UTC()
			service.UpdatedAt = now
		}
	}
}
