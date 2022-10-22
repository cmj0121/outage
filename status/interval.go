package status

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Interval time.Duration

const (
	// the default interval = 1s
	DEFAULT_INTERVAL Interval = Interval(time.Second)
)

func (interval Interval) String() (str string) {
	str = time.Duration(interval).String()
	return
}

func (interval *Interval) UnmarshalText(data []byte) (err error) {
	var str string
	var duration time.Duration

	if err = yaml.Unmarshal(data, &str); err != nil {
		log.Error().Err(err).Msg("cannot unmarshal YAML")
		return
	}

	if duration, err = time.ParseDuration(str); err != nil {
		log.Error().Err(err).Msg("cannot unmarshal Sampling")
		return
	}

	*interval = Interval(duration)
	return
}

func (interval Interval) MarshalJSON() ([]byte, error) {
	duration := time.Duration(interval)
	return json.Marshal(duration.String())
}

func (interval Interval) MarshalYAML() (interface{}, error) {
	duration := time.Duration(interval)
	return duration.String(), nil
}

func (interval Interval) Milliseconds() int64 {
	return time.Duration(interval).Milliseconds()
}
