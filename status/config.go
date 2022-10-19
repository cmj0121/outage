package status

import (
	"bufio"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Setting struct {
	Sampling Interval
	Timeout  Interval
}

type Config struct {
	Interval Interval
	Setting  Setting

	Services []*Service            `json:"services,omitempty" yaml:"services,omitempty"`
	Summary  map[string][]*Service `json:"summary,omitempty" yaml:"summary,omitempty"`
}

func (config Config) String() (raw string) {
	data, err := yaml.Marshal(config)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot marshal Config to YAML")
		return
	}

	raw = string(data)
	return
}

func Load[T []byte | string](data T) (config *Config, err error) {
	config = &Config{}
	err = config.UnmarshalText([]byte(data))
	return
}

func (config *Config) UnmarshalText(text []byte) (err error) {
	var data []byte

	switch path := string(text); path {
	case "-":
		log.Info().Msg("load config from STDIN")
		reader := bufio.NewReader(os.Stdin)
		if data, err = io.ReadAll(reader); err != nil {
			log.Error().Err(err).Msg("cannot open config file")
			return
		}
	default:
		log.Info().Str("file", path).Msg("load config from file")
		if data, err = os.ReadFile(path); err != nil {
			log.Error().Err(err).Msg("cannot open config file")
			return
		}
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Error().Err(err).Msg("cannot load config")
		return
	}

	config.epologue()
	return
}

func (config *Config) epologue() {
	config.Summary = map[string][]*Service{}

	for _, service := range config.Services {
		for _, tag := range service.Meta.Tags {
			switch _, ok := config.Summary[tag]; ok {
			case true:
				config.Summary[tag] = append(config.Summary[tag], service)
			case false:
				config.Summary[tag] = []*Service{service}
			}
		}
	}
}
