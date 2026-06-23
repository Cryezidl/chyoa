package web

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Cryezidl/cyoa/cyoa"
	"github.com/go-chi/chi/v5"
)

func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testTemplates(t *testing.T) *template.Template {
	t.Helper()
	// Tests run from the package directory, so templates live under ./templates.
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}
	return tmpl
}

func testStory() cyoa.Story {
	return cyoa.Story{
		"intro": {
			Title:      "The Beginning",
			Paragraphs: []string{"First paragraph."},
			Options: []cyoa.Options{
				{Text: "Go on", Arc: "end"},
			},
		},
		"end": {
			Title:      "The End",
			Paragraphs: []string{"It is over."},
			Options:    nil,
		},
	}
}

// doRequest routes a GET request for the given chapter param through chi so
// chi.URLParam resolves correctly, then returns the recorded response.
func doRequest(t *testing.T, h *WebHandler, chapter string) *httptest.ResponseRecorder {
	t.Helper()
	r := chi.NewRouter()
	r.Get("/", h.GetChapter)
	r.Get("/{chapter}", h.GetChapter)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/"+chapter, nil)
	r.ServeHTTP(rec, req)
	return rec
}

func TestGetChapter_Existing(t *testing.T) {
	h := NewWebHandler(testStory(), quietLogger(), testTemplates(t))
	rec := doRequest(t, h, "intro")

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "The Beginning") {
		t.Errorf("expected chapter title in body, got:\n%s", body)
	}
	if !strings.Contains(body, "Go on") {
		t.Errorf("expected option text in body, got:\n%s", body)
	}
}

func TestGetChapter_EmptyParamDefaultsToIntro(t *testing.T) {
	h := NewWebHandler(testStory(), quietLogger(), testTemplates(t))
	rec := doRequest(t, h, "")

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "The Beginning") {
		t.Errorf("empty param should default to intro chapter")
	}
}

func TestGetChapter_Ending(t *testing.T) {
	h := NewWebHandler(testStory(), quietLogger(), testTemplates(t))
	rec := doRequest(t, h, "end")

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "The End") {
		t.Errorf("expected ending title in body")
	}
}

func TestGetChapter_NotFound(t *testing.T) {
	h := NewWebHandler(testStory(), quietLogger(), testTemplates(t))
	rec := doRequest(t, h, "missing")

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
	// Falls back to the error chapter, which offers a route home.
	if !strings.Contains(rec.Body.String(), "Return to the beginning") {
		t.Errorf("expected error chapter fallback in body, got:\n%s", rec.Body.String())
	}
}
