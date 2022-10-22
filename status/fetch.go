package status

import (
	"time"

	"github.com/rs/zerolog/log"
)

func (config *Config) Fetch(closed <-chan struct{}, worker int) (err error) {
	ticker := time.NewTicker(time.Duration(config.Interval))
	defer ticker.Stop()

	ch := make(chan *Service, worker)
	defer close(ch)
	for n := 0; n < worker; n++ {
		log.Debug().Int("worker", n).Msg("create service worker")
		go config.worker(closed, ch)
	}

	for {
		select {
		case <-closed:
			return
		case <-ticker.C:
			log.Trace().Msg("tick")

			now := Timestamp(time.Now())
			for _, service := range config.Services {
				switch {
				case service.UpdatedAt.IsZero():
					log.Debug().Str("service", service.Title).Msg("first blood")
					ch <- service
				case now.After(service.UpdatedAt.Add(config.Setting.Sampling)):
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
			service.Fetch()
		}
	}
}
