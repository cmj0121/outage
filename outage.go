package outage

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Agent struct {
	Verbose int  `group:"log" xor:"log" short:"v" type:"counter" help:"the log verbose level"`
	Quiet   bool `group:"log" xor:"log" short:"q" help:"disabled all log"`
	Pretty  bool `group:"log" short:"p" help:"show the pretty log"`
	Color   bool `group:"log" default:"true" negatable:"" help:"show color"`

	Version Version `short:"V" help:"Show version info"`
}

func New() *Agent {
	return &Agent{}
}

func (agent *Agent) Run() (err error) {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	kong.Parse(agent)
	agent.prologue()
	defer agent.epologue()

	return
}

func (agent *Agent) prologue() {
	log.Info().Msg("starting prologue")
	agent.setup_logger()
	log.Info().Msg("finished prologue")
}

func (agent *Agent) epologue() {
	log.Info().Msg("starting epologue")
	log.Info().Msg("finished epologue")
}

func (agent *Agent) setup_logger() {
	if agent.Pretty {
		console := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !agent.Color}
		log.Logger = log.Output(console)
	}

	switch agent.Quiet {
	case true:
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case false:
		switch agent.Verbose {
		case 0:
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case 1:
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case 2:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case 3:
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		}
	}
}
