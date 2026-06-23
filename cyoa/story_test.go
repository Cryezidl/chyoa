package cyoa

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTempStory writes content to a temp file and returns its path.
func writeTempStory(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "story.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp story: %v", err)
	}
	return path
}

const validStory = `{
	"intro": {
		"title": "The Beginning",
		"story": ["First paragraph.", "Second paragraph."],
		"options": [
			{"text": "Go left", "arc": "left"},
			{"text": "Go right", "arc": "right"}
		]
	},
	"end": {
		"title": "The End",
		"story": ["It is over."],
		"options": []
	}
}`

func TestLoadStory(t *testing.T) {
	t.Run("valid file parses all chapters", func(t *testing.T) {
		path := writeTempStory(t, validStory)

		s, err := LoadStory(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(s) != 2 {
			t.Fatalf("expected 2 chapters, got %d", len(s))
		}

		intro, ok := s["intro"]
		if !ok {
			t.Fatal("expected 'intro' chapter to exist")
		}
		if intro.Title != "The Beginning" {
			t.Errorf("unexpected title: %q", intro.Title)
		}
		if len(intro.Paragraphs) != 2 {
			t.Errorf("expected 2 paragraphs, got %d", len(intro.Paragraphs))
		}
		if len(intro.Options) != 2 {
			t.Fatalf("expected 2 options, got %d", len(intro.Options))
		}
		if intro.Options[0].Text != "Go left" || intro.Options[0].Arc != "left" {
			t.Errorf("unexpected first option: %+v", intro.Options[0])
		}
	})

	t.Run("missing file returns error", func(t *testing.T) {
		_, err := LoadStory(filepath.Join(t.TempDir(), "does-not-exist.json"))
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})

	t.Run("invalid json returns error", func(t *testing.T) {
		path := writeTempStory(t, `{"intro": { not valid json `)

		_, err := LoadStory(path)
		if err == nil {
			t.Fatal("expected error for invalid json, got nil")
		}
	})

	t.Run("empty json object yields empty story", func(t *testing.T) {
		path := writeTempStory(t, `{}`)

		s, err := LoadStory(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(s) != 0 {
			t.Errorf("expected empty story, got %d chapters", len(s))
		}
	})
}

func TestGetChapter(t *testing.T) {
	s := Story{
		"intro": {Title: "The Beginning"},
		"end":   {Title: "The End"},
	}

	t.Run("existing chapter is returned", func(t *testing.T) {
		chapter, ok := GetChapter(s, "intro")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if chapter.Title != "The Beginning" {
			t.Errorf("unexpected title: %q", chapter.Title)
		}
	})

	t.Run("missing chapter reports not found", func(t *testing.T) {
		_, ok := GetChapter(s, "nope")
		if ok {
			t.Error("expected ok to be false for missing chapter")
		}
	})
}
