package btea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Spinner struct {
	spinner    spinner.Model
	isQuitting *bool
	err        error
}

func InitialSpinnerModel(isQuitting *bool) Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Spinner{spinner: s, isQuitting: isQuitting}
}

func (m Spinner) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			*m.isQuitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Spinner) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Loading food...press q to quit\n\n", m.spinner.View())
	if *m.isQuitting {
		return ""
	}
	return str
}
