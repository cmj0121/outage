package web

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Web struct {
	closed      <-chan struct{}
	*mux.Router `kong:"-"`

	Bind string `short:"b" default:"127.0.0.1:9999" help:"the address bind to"`
}

func New(closed <-chan struct{}) (web *Web) {
	web = &Web{
		closed: closed,

		Router: mux.NewRouter(),
		Bind:   "127.0.0.1:9999",
	}

	web.Use(web.MiddlewareLog)
	web.HandleFunc("/", web.IndexPage)
	return
}

func (web *Web) IndexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
