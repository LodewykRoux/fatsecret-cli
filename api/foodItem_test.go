package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFoodDescription(t *testing.T) {
	testFood := Food{
		BrandName:       "TestFood",
		FoodDescription: "Per 1 serving - Calories: 300kcal | Fat: 13.00g | Carbs: 32.00g | Protein: 15.00g",
		Calories:        "300",
		Fat:             "13.00",
		Carbs:           "32.00",
		Protein:         "15.00",
	}
	food := Food{
		BrandName:       "TestFood",
		FoodDescription: "Per 1 serving - Calories: 300kcal | Fat: 13.00g | Carbs: 32.00g | Protein: 15.00g",
	}

	ParseFoodDescription(&food)

	assert.Equal(t, testFood, food)

}
