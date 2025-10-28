package openapi

import (
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed openapi.yaml
var openAPISpec []byte

//go:embed swagger.html
var swaggerHTMLTemplate string

var swaggerTemplate = template.Must(template.New("swagger").Parse(swaggerHTMLTemplate))

type Config struct {
	Title   string
	SpecURL string
}

func NewHandler(basePath, title string) http.Handler {
	mux := http.NewServeMux()

	config := Config{
		Title:   title,
		SpecURL: basePath + "/openapi.yaml",
	}

	mux.HandleFunc(basePath+"/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != basePath+"/" && r.URL.Path != basePath {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		swaggerTemplate.Execute(w, config)
	})

	mux.HandleFunc(basePath+"/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write(openAPISpec)
	})

	return mux
}
