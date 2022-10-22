package status

import (
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Meta struct {
	// the title of the service
	Title string `json:"title"`

	// the brief of service
	Subject string `json:"subject"`

	// the service link
	Link string `json:"link"`

	// the service tags
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// / The service meta and status
type Service struct {
	Meta `json:",inline" yaml:",inline"`

	Mode Mode `json:"mode" yaml:",omitempty"`

	UpdatedAt time.Time `json:"updated_at" yaml:",omitempty"`
}

func (svc Service) String() (raw string) {
	data, err := yaml.Marshal(svc)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot marshal Service to YAML")
		return
	}

	raw = string(data)
	return
}

func (svc *Service) Fetch() {
	svc.UpdatedAt = time.Now().UTC()

	parse_url, err := url.Parse(svc.Link)
	if err != nil {
		log.Error().Err(err).Msg("invalid link")
		return
	}

	switch scheme := parse_url.Scheme; scheme {
	case "http", "https":
		client := &http.Client{
			// not follow the redirect
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Get(svc.Link)
		switch {
		case err != nil:
			log.Warn().Err(err).Msg("HTTP GET failure")
			svc.Mode = OFF
		case resp.StatusCode >= http.StatusInternalServerError:
			log.Warn().Str("link", svc.Link).Str("status", resp.Status).Msg("HTTP GET fail")
			svc.Mode = OFF
		case resp.StatusCode >= http.StatusBadRequest:
			log.Warn().Str("link", svc.Link).Str("status", resp.Status).Msg("HTTP GET denied")
			svc.Mode = INCIDENT
		case resp.StatusCode >= http.StatusBadRequest:
			log.Info().Str("link", svc.Link).Str("status", resp.Status).Msg("HTTP GET redirect")
			svc.Mode = INCIDENT
		case resp.StatusCode >= http.StatusOK:
			log.Info().Str("link", svc.Link).Str("status", resp.Status).Msg("HTTP GET success")
			svc.Mode = ON
		default:
			log.Warn().Str("link", svc.Link).Str("status", resp.Status).Msg("HTTP GET unknown")
			svc.Mode = UNKNOWN
		}
	default:
		log.Warn().Str("scheme", scheme).Msg("not implemented scheme")
		return
	}
}
