package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type FoodItemGetter interface {
	GetFoodSuggestions(searchText string, accessToken string) (*FoodsResponse, error)
}

type ApiFoodItemGetter struct{}
type Food struct {
	BrandName       string `json:"brand_name,omitempty"`
	FoodDescription string `json:"food_description"`
	FoodID          string `json:"food_id"`
	FoodName        string `json:"food_name"`
	FoodType        string `json:"food_type"`
	FoodURL         string `json:"food_url"`
	ServingSize     string `json:"-"`
	Calories        string `json:"-"`
	Fat             string `json:"-"`
	Carbs           string `json:"-"`
	Protein         string `json:"-"`
}

type FoodsResponse struct {
	Foods struct {
		Food         []Food `json:"food"`
		MaxResults   string `json:"max_results"`
		PageNumber   string `json:"page_number"`
		TotalResults string `json:"total_results"`
	} `json:"foods"`
}

func FormatFoodDetails(food Food) string {
	return fmt.Sprintf(
		"    🍽  Food: %s\n"+
			"    📏 Serving Size: %s\n"+
			"    🔥 Calories: %s kcal\n"+
			"    🥑 Fat: %s g\n"+
			"    🍞 Carbs: %s g\n"+
			"    🍗 Protein: %s g\n",
		food.FoodName,
		food.ServingSize,
		food.Calories,
		food.Fat,
		food.Carbs,
		food.Protein,
	)
}

func ParseFoodDescription(food *Food) {
	if food.FoodDescription == "" {
		return
	}

	sizeRegex := regexp.MustCompile(`Per (\d+g)`)
	matches := sizeRegex.FindStringSubmatch(food.FoodDescription)
	if len(matches) > 1 {
		food.ServingSize = matches[1]
	}

	nutritionRegex := regexp.MustCompile(`Calories:\s*([\d.]+)kcal\s*\|\s*Fat:\s*([\d.]+)g\s*\|\s*Carbs:\s*([\d.]+)g\s*\|\s*Protein:\s*([\d.]+)g`)
	matches = nutritionRegex.FindStringSubmatch(food.FoodDescription)

	if len(matches) == 5 {
		food.Calories = matches[1]
		food.Fat = matches[2]
		food.Carbs = matches[3]
		food.Protein = matches[4]
	}
}

func (apiF ApiFoodItemGetter) GetFoodSuggestions(searchText string, accessToken string) (*FoodsResponse, error) {

	baseUrl := "https://platform.fatsecret.com/rest/foods/search/v1"

	params := url.Values{}
	params.Add("search_expression", searchText)
	params.Add("format", "json")
	params.Add("max_results", "10")

	fullURL := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var FoodsResponse FoodsResponse
	err = json.Unmarshal(body, &FoodsResponse)
	if err != nil {
		return nil, err
	}

	for i := range FoodsResponse.Foods.Food {
		ParseFoodDescription(&FoodsResponse.Foods.Food[i])
	}

	return &FoodsResponse, nil
}
