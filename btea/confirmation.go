package btea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type conFirmationModel struct {
	choices        []string
	cursor         int
	selectedAnswer *string
	quitting       bool
}

func InitialConfirmationModel(selectedAnswer *string) conFirmationModel {
	return conFirmationModel{
		choices:        []string{"No", "Yes"},
		selectedAnswer: selectedAnswer,
	}
}

func (m conFirmationModel) Init() tea.Cmd {
	return nil
}

func (m conFirmationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			*m.selectedAnswer = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m conFirmationModel) View() string {
	if *m.selectedAnswer == "Yes" {
		// Get the Food choice here and use its selection
		return QuitTextStyle.Render("Saving to favourites")
	} else if *m.selectedAnswer == "No" {
		// Get the Food choice here and use its selection
		return QuitTextStyle.Render("")
	}
	if m.quitting {
		return QuitTextStyle.Render("Quiting.")
	}
	s := ""
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s \n", cursor, choice)
	}
	s += "\nPress q to quit.\n"
	return s
}
