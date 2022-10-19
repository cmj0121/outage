package status

import (
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Meta struct {
	// the title of the service
	Title string

	// the brief of service
	Subject string

	// the service link
	Link string

	// the service tags
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// / The service meta and status
type Service struct {
	Meta `json:",inline" yaml:",inline"`

	Mode Mode `yaml:",omitempty"`

	UpdatedAt *time.Time `yaml:",omitempty"`
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
