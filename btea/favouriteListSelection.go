package btea

import (
	"fmt"

	"github.com/LodewykRoux/fatsecret-cli/favourites"
	tea "github.com/charmbracelet/bubbletea"
)

type favouriteSelection struct {
	choices           []favourites.Favourite
	cursor            int
	choice            favourites.Favourite
	selectedFavourite *favourites.Favourite
	quitting          bool
}

func InitialFavouriteModel(favouriteList []favourites.Favourite, selectedFavourite *favourites.Favourite) favouriteSelection {
	return favouriteSelection{
		choices:           favouriteList,
		selectedFavourite: selectedFavourite,
	}
}

func (m favouriteSelection) Init() tea.Cmd {
	return nil
}

func (m favouriteSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.selectedFavourite != nil {
				*m.selectedFavourite = m.choice
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

func (m favouriteSelection) View() string {
	if m.choice.Name != "" {
		return QuitTextStyle.Render(favourites.FormatFavourite(m.choice))
	}
	if m.quitting {
		return QuitTextStyle.Render("Quiting.")
	}
	s := "Choose a favourite?\n\n"
	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s - Per %s\n", cursor, choice.Name, choice.Uom)
	}
	s += "\nPress q to quit.\n"
	return s
}
