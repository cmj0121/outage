package status

import (
	"encoding/json"
)

type Mode int

const (
	UNKNOWN Mode = iota
	ON
	OFF
	INCIDENT
)

func (mode Mode) String() (str string) {
	switch mode {
	case ON:
		str = "on"
	case OFF:
		str = "off"
	case INCIDENT:
		str = "incident"
	default:
		str = "unknown"
	}
	return
}

func (mode Mode) MarshalJSON() ([]byte, error) {
	var str = mode.String()
	return json.Marshal(str)
}

func (mode Mode) MarshalYAML() (interface{}, error) {
	return mode.String(), nil
}
