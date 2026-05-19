package web

import (
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
	if !ok {
		h.logger.Warn("Chapter not found", "param", chapterParam)

		w.WriteHeader(http.StatusNotFound)
		errorChapter := cyoa.Chapter{
			Title: "Такая история нам незнакома...",
			Paragraphs: []string{
				"Вы забрели в темный угол библиотеки, где книги еще не написаны.",
				"Возможно, стоит вернуться в начало?",
			},
			Options: []cyoa.Options{
				{Text: "Вернуться в начало", Arc: "intro"},
			},
		}
		if err := h.tmp.ExecuteTemplate(w, "story", errorChapter); err != nil {
			h.logger.Error("Failed to execute error template", "err", err)
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

	if len(chapter.Options) == 0 {
		if err := h.tmp.ExecuteTemplate(w, "the end", chapter); err != nil {
			h.logger.Error("Failed to execute 'the end' template", "err", err)
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	if err := h.tmp.ExecuteTemplate(w, "story", chapter); err != nil {
		h.logger.Error("Failed to execute story template", "err", err)
		http.Error(w, "Internal Server Error", 500)
	}
}
