package btea

import (
	"bytes"
	"testing"
	"time"

	"github.com/LodewykRoux/fatsecret-cli/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

func TestFoodListSelection(t *testing.T) {
	var selectedFood api.Food

	favs := []api.Food{{
		BrandName: "name",
	}}

	tm := teatest.NewTestModel(t, InitialFoodModel(favs, &selectedFood),
		teatest.WithInitialTermSize(300, 100),
	)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Choose a food?"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*6))

	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))

	assert.Equal(t, selectedFood.BrandName, favs[len(favs)-1].BrandName)
}
