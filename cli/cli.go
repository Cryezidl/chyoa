package cli

import (
	"fmt"
	"log/slog"

	"github.com/Cryezidl/cyoa/cyoa"
	tea "github.com/charmbracelet/bubbletea"
)

func StartGame(s cyoa.Story, arcName string, logger *slog.Logger) error {
	m := InitialModel(s, arcName)
	m.currentArc = arcName
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("Alas, there's been an error: %v", err)
	}

	return nil
}
