package web

import (
	"context"
	"net/http"

	"github.com/cmj0121/outage/status"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Web struct {
	closed <-chan struct{}

	*mux.Router    `kong:"-"`
	*status.Config `kong:"-"`

	Bind string `short:"b" default:"127.0.0.1:9999" help:"the address bind to"`
}

func New(closed <-chan struct{}, config *status.Config) (web *Web) {
	web = &Web{
		closed: closed,

		Router: mux.NewRouter(),
		Config: config,

		Bind: "127.0.0.1:9999",
	}

	web.Use(web.MiddlewareLog)
	web.HandleFunc("/", web.IndexPage)
	web.HandleFunc("/service", web.Response(web.Config.Services))
	web.HandleFunc("/summary", web.Response(web.Config.Summary))
	return
}

func (web *Web) ServeHTTP() (err error) {
	srv := http.Server{
		Addr:    web.Bind,
		Handler: web.Router,
	}

	go func() {
		<-web.closed

		log.Trace().Msg("start shutdown HTTP server ...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("cannot shutdown HTTP server")
			return
		}
		log.Info().Msg("shutdown HTTP server")
	}()

	log.Info().Str("bind", web.Bind).Msg("start HTTP server")
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatal().Err(err).Msg("cannot run HTTP server")
		return
	}

	err = nil
	return
}

func (web *Web) MiddlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("ip", r.RemoteAddr).Str("method", r.Method).Msg(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
