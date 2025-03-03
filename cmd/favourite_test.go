package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/LodewykRoux/fatsecret-cli/favourites"
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

func TestFavouriteCommand(t *testing.T) {
	fName := SetupTempFile(t)

	var output bytes.Buffer
	cmd := NewFavouriteCmd(fName)
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestFavouriteListCommand(t *testing.T) {
	fName := SetupTempFile(t)

	// Capture output
	var out bytes.Buffer
	cmd := NewFavouriteListCmd(fName)
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{})

	// Execute
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestNewFavouriteDeleteListCmd(t *testing.T) {
	fName := SetupTempFile(t)

	// Capture output
	var out bytes.Buffer
	cmd := NewFavouriteDeleteListCmd(fName)
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{})

	// Execute
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestNewFavouriteSearchListCmd(t *testing.T) {
	fName := SetupTempFile(t)

	// Capture output
	var out bytes.Buffer
	cmd := NewFavouriteSearchListCmd(fName)
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{})

	// Execute
	err := cmd.Execute()
	assert.NoError(t, err)
}
