package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	_ "embed"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

//go:embed template/index.htm
var TMPL_HTML_INDEX string

var (
	TemplateFuncMap = template.FuncMap{
		"safe_tag": func(raw string) (str string) {
			str = strings.ReplaceAll(raw, " ", "_")
			str = strings.ToLower(str)
			return
		},
	}
)

func (web *Web) IndexPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	tmpl, err := template.New("index").Funcs(TemplateFuncMap).Parse(TMPL_HTML_INDEX)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot generate template")
		return
	}

	if err = tmpl.Execute(w, web.Config); err != nil {
		log.Fatal().Err(err).Msg("cannot render template")
		return
	}
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
