package btea

import (
	"fmt"

	"github.com/LodewykRoux/fatsecret-cli/api"
	tea "github.com/charmbracelet/bubbletea"
)

type foodSelection struct {
	choices      []api.Food
	cursor       int
	choice       api.Food
	selectedFood *api.Food
	quitting     bool
}

func InitialFoodModel(foods []api.Food, selectedFood *api.Food) foodSelection {
	return foodSelection{
		choices:      foods,
		selectedFood: selectedFood,
	}
}

func (m foodSelection) Init() tea.Cmd {
	return nil
}

func (m foodSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}

		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case tea.KeyEnter, tea.KeySpace:
			m.choice = m.choices[m.cursor]
			if m.selectedFood != nil {
				*m.selectedFood = m.choice
			}
			return m, tea.Quit

		case tea.KeyRunes:
			switch msg.String() {
			case "q":
				return m, tea.Quit

			case "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			}
		}
	}
	return m, nil
}

func (m foodSelection) View() string {
	if m.choice.FoodDescription != "" {
		// Get the Food choice here and use its selection
		return QuitTextStyle.Render(api.FormatFoodDetails(m.choice))
	}
	if m.quitting {
		return QuitTextStyle.Render("Quiting.")
	}
	s := "Choose a food?\n\n"
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s - Per %s\n", cursor, choice.FoodName, choice.ServingSize)
	}
	s += "\nPress q to quit.\n"
	return s
}
