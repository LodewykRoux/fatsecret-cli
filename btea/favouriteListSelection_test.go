package btea

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/LodewykRoux/fatsecret-cli/favourites"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

var tmpFileName = "favourites_test_*.csv"

func SeedData(fileName string) {
	favourites.SaveFavourite(favourites.Favourite{Name: "Apple", Uom: "1 piece", Calories: "95", Fat: "0.3", Carbs: "25", Protein: "0.5"}, fileName)
	favourites.SaveFavourite(favourites.Favourite{Name: "Banana", Uom: "1 piece", Calories: "105", Fat: "0.3", Carbs: "27", Protein: "1.3"}, fileName)
}

func SetupTempFile(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", tmpFileName)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	return tmpFile.Name()
}

func TestFavouriteListSelection(t *testing.T) {
	var selectedFavourite favourites.Favourite
	fName := SetupTempFile(t)

	SeedData(fName)
	favs, err := favourites.RetrieveFavourites(fName)
	if err != nil {
		t.Fatalf("Failed to retrieve favourites: %v", err)
	}

	tm := teatest.NewTestModel(t, InitialFavouriteModel(favs, &selectedFavourite),
		teatest.WithInitialTermSize(300, 100),
	)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Choose a favourite?"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*3))

	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))

	assert.Equal(t, selectedFavourite.Name, favs[len(favs)-1].Name)
}
