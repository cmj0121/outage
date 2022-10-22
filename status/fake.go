package status

import (
	"time"
)

func Fake() (config *Config) {
	config = &Config{
		Interval: DEFAULT_INTERVAL,
		Setting: Setting{
			Sampling: Interval(time.Duration(time.Second * 3600)),
			Timeout:  Interval(time.Duration(time.Second * 60)),
		},

		Services: []*Service{
			{
				Mode: ON,
				Meta: Meta{
					Title:   "example",
					Subject: "The example service",
					Link:    "https://example.com",
					Tags: []string{
						"Group 1",
						"Group 2",
					},
				},
			},
			{
				Mode: OFF,
				Meta: Meta{
					Title:   "example",
					Subject: "The example service",
					Link:    "https://example.com/off",
					Tags: []string{
						"Group 1",
						"Group 2",
					},
				},
			},
			{
				Mode: INCIDENT,
				Meta: Meta{
					Title:   "example",
					Subject: "The example service",
					Link:    "https://example.com/incident",
					Tags: []string{
						"Group 1",
					},
				},
			},
		},
		Summary: map[string][]*Service{},
	}

	config.epologue()
	return
}
