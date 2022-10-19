package outage

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/cmj0121/outage/status"
	"github.com/cmj0121/outage/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Agent struct {
	closed chan struct{}

	Verbose int  `group:"log" xor:"log" short:"v" type:"counter" help:"the log verbose level"`
	Quiet   bool `group:"log" xor:"log" short:"q" help:"disabled all log"`
	Pretty  bool `group:"log" short:"p" help:"show the pretty log"`
	Color   bool `group:"log" default:"true" negatable:"" help:"show color"`

	Version Version `short:"V" help:"Show version info"`

	*web.Web

	Config *status.Config `group:"config" xor:"config" short:"c" help:"The task config settings"`
	Fake   bool           `group:"config" xor:"config" short:"F" help:"load the fake config"`

	Action string `arg:"" default:"server" enum:"server,dump" help:"run as"`
}

func New() (agent *Agent) {
	agent = &Agent{
		closed: make(chan struct{}),
	}

	agent.Web = web.New(agent.closed)
	return
}

func (agent *Agent) Run() (err error) {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	kong.Parse(agent)
	agent.prologue()
	defer agent.epologue()

	err = agent.run()
	return
}

func (agent *Agent) run() (err error) {
	switch agent.Action {
	case "dump":
		fmt.Println(agent.Config)
	case "server":
		err = agent.ServeHTTP()
	}

	log.Trace().Err(err).Msg("agent stoped")
	return
}

func (agent *Agent) prologue() {
	log.Info().Msg("starting prologue")

	agent.setup_logger()
	agent.setup_graceful_shutdown()

	if agent.Fake {
		log.Warn().Msg("load the fake config")
		agent.Config = status.Fake()
	}

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

func (agent *Agent) setup_graceful_shutdown() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)

		<-sigint

		log.Info().Msg("starting graceful shutdown")
		close(agent.closed)
	}()
}
