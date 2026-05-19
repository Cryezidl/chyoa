package cli

import (
	"fmt"
	"strings"

	"github.com/Cryezidl/cyoa/cyoa"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	story      cyoa.Story
	cursor     int
	currentArc string
}

var (
	// Стиль для заголовка главы
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	// Стиль для основного текста истории
	storyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ABB2BF")).
			Italic(true).
			MaxWidth(130).
			PaddingLeft(2)

	// Стиль для активного (выбранного) пункта меню
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#EE6FF8")).
				Bold(true).
				PaddingLeft(2)

	// Стиль для обычного пункта меню
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#61AFEF")).
			PaddingLeft(4)
)

func InitialModel(s cyoa.Story, arcName string) model {
	return model{
		story:      s,
		currentArc: arcName,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	chapter, ok := m.story[m.currentArc]
	if !ok {
		return m, nil
	}

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if len(chapter.Options) > 0 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(chapter.Options) - 1 // прыгаем в конец
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if len(chapter.Options) > 0 {
				m.cursor++
				if m.cursor >= len(chapter.Options) {
					m.cursor = 0 // прыгаем в начало
				}
			}

		// The "enter" key and the space bar toggle the selected state
		// for the item that the cursor is pointing at.
		case "enter", "space":
			if len(chapter.Options) > 0 {
				m.currentArc = chapter.Options[m.cursor].Arc
				m.cursor = 0
			}

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	chapter, ok := m.story[m.currentArc]
	if !ok {
		return fmt.Sprintf("Arc not found: %s", m.currentArc)
	}

	// Пишем Тайтл и Параграфы
	var s strings.Builder
	s.WriteString(titleStyle.Render(chapter.Title) + "\n\n")

	for _, p := range chapter.Paragraphs {
		s.WriteString(storyStyle.Render(p) + "\n")
	}

	if len(chapter.Options) == 0 {
		s.WriteString("\n\nThe end. Thanks for playing!\n")
		s.WriteString("\nPress q to quit.\n")
		return s.String()
	}
	//// Пишем Опции
	s.WriteString(lipgloss.NewStyle().Underline(true).Render("What will you choose?") + "\n\n")
	for i, choice := range chapter.Options {
		// Is the cursor pointing at this choice?
		if m.cursor == i {
			s.WriteString(selectedItemStyle.Render("--> "+choice.Text) + "\n")
		} else {
			s.WriteString(itemStyle.Render(choice.Text) + "\n")
		}
	}

	// The footer
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")).MarginTop(1)
	s.WriteString("\n" + helpStyle.Render("Press q or enter to quit.") + "\n")
	// Send the UI for rendering
	return s.String()
}
