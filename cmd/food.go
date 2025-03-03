/*
Copyright © 2025 Lodewyk Roux lloderoux@protonmail.ch
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/LodewykRoux/fatsecret-cli/api"
	"github.com/LodewykRoux/fatsecret-cli/btea"
	"github.com/LodewykRoux/fatsecret-cli/favourites"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewFoodCmd(accessToken string, fileName string, apiF api.FoodItemGetter) *cobra.Command {
	var search_text string
	var selectedFood api.Food

	var foodCmd = &cobra.Command{
		Use:   "food",
		Short: "Search fod foods on the internet.",
		Long: `The food command is used to search for foods on the internet.
	The food command requires a term to work:

food --term="something": finds the top 10 related items, selecting chooses that item.`,
		Run: func(cmd *cobra.Command, args []string) {
			if search_text == "" {
				fmt.Println("A term to search for is required. use --term=<search>")
				return
			}
			var isQuitting bool
			l := tea.NewProgram(btea.InitialSpinnerModel(&isQuitting))
			done := make(chan struct{})

			go func() {
				_, err := l.Run()
				if err != nil {
					fmt.Println("Error running spinner:", err)
				}
			}()

			var result *api.FoodsResponse
			var err error

			go func() {
				result, err = apiF.GetFoodSuggestions(search_text, accessToken)
				close(done)
			}()

			<-done

			// Add way to quit while loading, also rename isQuitting to isLoading
			isQuitting = true
			l.Quit()

			if err != nil {
				fmt.Println("Error fetching food suggestions.")
				return
			}

			p := tea.NewProgram(btea.InitialFoodModel(result.Foods.Food, &selectedFood))
			_, err = p.Run()
			if err != nil {
				fmt.Printf("Error displaying food selection: %v", err)
			}

			if selectedFood.FoodDescription == "" {
				return
			}

			fmt.Println("Save to favourites?")
			var selectedAnswer string
			s := tea.NewProgram(btea.InitialConfirmationModel(&selectedAnswer))
			_, err = s.Run()
			if strings.ToLower(selectedAnswer) == "yes" {
				err = favourites.SaveFavourite(favourites.Favourite{
					Name:     selectedFood.FoodName,
					Uom:      selectedFood.ServingSize,
					Calories: selectedFood.Calories,
					Fat:      selectedFood.Fat,
					Carbs:    selectedFood.Carbs,
					Protein:  selectedFood.Protein,
				}, fileName)
				if err != nil {
					fmt.Printf("Error saving: %v", err)
				} else {
					fmt.Printf("%s has been saved to your favourites.", selectedFood.FoodName)
				}
			}
		},
	}
	foodCmd.Flags().StringVarP(&search_text, "term", "k", "", "Term to search for")

	return foodCmd
}
