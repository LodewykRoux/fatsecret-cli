package cmd

import (
	"bytes"
	"testing"

	"github.com/LodewykRoux/fatsecret-cli/api"
	"github.com/stretchr/testify/assert"
)

type stubFoodItemGetter struct {
	foodItem    api.FoodsResponse
	ReturnError bool
}

func (fi stubFoodItemGetter) GetFoodSuggestions(searchText string, accessToken string) (*api.FoodsResponse, error) {
	return &fi.foodItem, nil
}

func TestFoodCommand(t *testing.T) {
	fName := SetupTempFile(t)

	var output bytes.Buffer
	cmd := NewFoodCmd("", fName, stubFoodItemGetter{})
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	err := cmd.Execute()
	assert.NoError(t, err)
}
