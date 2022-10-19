package web

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func (web *Web) IndexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (web Web) Response(raw interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		format := r.URL.Query().Get("format")
		switch format {
		case "yaml":
			data, err := yaml.Marshal(raw)
			if err != nil {
				web.ResponseError(w, r, err)
				return
			}

			w.Header().Set("Content-Type", "text/yaml")
			if _, err := w.Write([]byte(data)); err != nil {
				log.Error().Err(err).Msg("cannot write HTTP response")
				return
			}
		case "json":
			fallthrough
		default:
			data, err := json.Marshal(raw)
			if err != nil {
				web.ResponseError(w, r, err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(data)); err != nil {
				log.Error().Err(err).Msg("cannot write HTTP response")
				return
			}
		}
	}
}

func (web Web) ResponseError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(err.Error())); err != nil {
		log.Error().Err(err).Msg("cannot write HTTP response")
		return
	}
}
