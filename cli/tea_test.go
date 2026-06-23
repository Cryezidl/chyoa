package cli

import (
	"strings"
	"testing"

	"github.com/Cryezidl/cyoa/cyoa"
	tea "github.com/charmbracelet/bubbletea"
)

func testStory() cyoa.Story {
	return cyoa.Story{
		"intro": {
			Title:      "The Beginning",
			Paragraphs: []string{"First paragraph.", "Second paragraph."},
			Options: []cyoa.Options{
				{Text: "Go left", Arc: "left"},
				{Text: "Go right", Arc: "right"},
			},
		},
		"left": {
			Title:      "Left Path",
			Paragraphs: []string{"You went left."},
			Options:    nil,
		},
	}
}

// key builds a tea.KeyMsg from a single key string like "down" or "enter".
func key(s string) tea.KeyMsg {
	switch s {
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

// update applies a key and returns the concrete model for assertions.
func update(t *testing.T, m model, k tea.KeyMsg) model {
	t.Helper()
	next, _ := m.Update(k)
	out, ok := next.(model)
	if !ok {
		t.Fatalf("Update returned unexpected type %T", next)
	}
	return out
}

func TestInitialModel(t *testing.T) {
	m := InitialModel(testStory(), "intro")
	if m.currentArc != "intro" {
		t.Errorf("expected currentArc 'intro', got %q", m.currentArc)
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor 0, got %d", m.cursor)
	}
}

func TestCursorMovement(t *testing.T) {
	m := InitialModel(testStory(), "intro")

	m = update(t, m, key("down"))
	if m.cursor != 1 {
		t.Errorf("after down: expected cursor 1, got %d", m.cursor)
	}

	// Wrap from the last option back to the first.
	m = update(t, m, key("down"))
	if m.cursor != 0 {
		t.Errorf("after wrap down: expected cursor 0, got %d", m.cursor)
	}

	// Wrap from the first option to the last when going up.
	m = update(t, m, key("up"))
	if m.cursor != 1 {
		t.Errorf("after wrap up: expected cursor 1, got %d", m.cursor)
	}
}

func TestSelectOptionFollowsArc(t *testing.T) {
	m := InitialModel(testStory(), "intro")

	// Move to "Go right" (index 1) and select it.
	m = update(t, m, key("down"))
	m = update(t, m, key("enter"))

	if m.currentArc != "right" {
		t.Errorf("expected currentArc 'right', got %q", m.currentArc)
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor reset to 0 after selection, got %d", m.cursor)
	}
}

func TestQuitReturnsQuitCommand(t *testing.T) {
	m := InitialModel(testStory(), "intro")
	_, cmd := m.Update(key("q"))
	if cmd == nil {
		t.Fatal("expected a quit command, got nil")
	}
	// tea.Quit returns a tea.QuitMsg when invoked.
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Error("expected command to produce tea.QuitMsg")
	}
}

func TestViewRendersChapter(t *testing.T) {
	m := InitialModel(testStory(), "intro")
	out := m.View()

	for _, want := range []string{"The Beginning", "First paragraph.", "Go left", "Go right"} {
		if !strings.Contains(out, want) {
			t.Errorf("view missing %q\n---\n%s", want, out)
		}
	}
}

func TestViewRendersEnding(t *testing.T) {
	m := InitialModel(testStory(), "left")
	out := m.View()

	if !strings.Contains(out, "The end") {
		t.Errorf("ending view should show end message, got:\n%s", out)
	}
}

func TestViewUnknownArc(t *testing.T) {
	m := InitialModel(testStory(), "ghost")
	out := m.View()
	if !strings.Contains(out, "Arc not found") {
		t.Errorf("expected 'Arc not found' for unknown arc, got:\n%s", out)
	}
}
