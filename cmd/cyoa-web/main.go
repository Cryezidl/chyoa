package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/Cryezidl/cyoa/cyoa"
	"github.com/Cryezidl/cyoa/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	filepath := flag.String("filepath", "gopher.json", "path to the file with story")
	flag.Parse()

	logger := slog.Default()
	story, err := cyoa.LoadStory(*filepath, logger)
	if err != nil {
		logger.Error("Couldn't load story file", "Error", err)
		os.Exit(1)
	}
	t, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		logger.Error("Couldn't parse html templates", "Error", err)
		os.Exit(1)
	}

	h := web.NewWebHandler(story, logger, t)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/api/v1/cyoa", func(r chi.Router) {
		r.Get("/", h.GetChapter)
		r.Get("/{chapter}", h.GetChapter)
	})

	logger.Info("Starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error("Server failed", "Error", err)
	}
}
