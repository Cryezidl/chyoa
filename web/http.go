package web

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Cryezidl/cyoa/cyoa"
	"github.com/go-chi/chi/v5"
)

type WebHandler struct {
	story  cyoa.Story
	logger *slog.Logger
	tmp    *template.Template
}

func NewWebHandler(s cyoa.Story, l *slog.Logger, t *template.Template) *WebHandler {
	return &WebHandler{story: s, logger: l, tmp: t}
}

func (h *WebHandler) GetChapter(w http.ResponseWriter, r *http.Request) {
	chapterParam := chi.URLParam(r, "chapter")
	if strings.TrimSpace(chapterParam) == "" {
		chapterParam = "intro"
	}

	chapter, ok := cyoa.GetChapter(h.story, chapterParam)

	status := http.StatusOK
	tmplName := "story"

	switch {
	case !ok:
		h.logger.Warn("Chapter not found", "param", chapterParam)
		status = http.StatusNotFound
		chapter = cyoa.Chapter{
			Title: "A tale such as this remains untold......",
			Paragraphs: []string{
				"You have wandered into a dark corner of the library, where books have not yet been written.",
				"Perhaps it's time to return to the beginning?",
			},
			Options: []cyoa.Options{
				{Text: "Return to the beginning", Arc: "intro"},
			},
		}
	case len(chapter.Options) == 0:
		tmplName = "the end"
	}

	// Render into a buffer first: if the template fails we can still send a
	// clean 500 instead of a partial body under an already-committed status.
	var buf bytes.Buffer
	if err := h.tmp.ExecuteTemplate(&buf, tmplName, chapter); err != nil {
		h.logger.Error("Failed to execute template", "template", tmplName, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := buf.WriteTo(w); err != nil {
		h.logger.Error("Failed to write response", "err", err)
	}
}
