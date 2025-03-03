/*
Copyright © 2025 Lodewyk Roux lloderoux@protonmail.ch
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/LodewykRoux/fatsecret-cli/btea"
	"github.com/LodewykRoux/fatsecret-cli/favourites"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewFavouriteCmd(fileName string) *cobra.Command {
	var favouriteCmd = &cobra.Command{
		Use:   "favourite",
		Short: "Retrieve, search and delete local favourites",
		Long: `A list of commands that deal with the locally saved favourites:

favourite list: Lists all favourites, selecting views the details.
favourite delete: Lists all favourites, selecting prompts you to delete it.
favourite search --term"something": Lists all favourites related to your search.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("favourite called")
		},
	}
	return favouriteCmd
}

func NewFavouriteListCmd(fileName string) *cobra.Command {
	var selectedFavourite favourites.Favourite
	var favouriteListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all favourite foods",
		Run: func(cmd *cobra.Command, args []string) {
			var isQuitting bool
			l := tea.NewProgram(btea.InitialSpinnerModel(&isQuitting))
			done := make(chan struct{})

			go func() {
				_, err := l.Run()
				if err != nil {
					fmt.Println("Error running spinner:", err)
				}
			}()

			var favouritesList []favourites.Favourite
			var err error

			go func() {
				favouritesList, err = favourites.RetrieveFavourites(fileName)
				close(done)
			}()

			<-done

			isQuitting = true
			l.Quit()

			if err != nil {
				fmt.Println("Error while retrieving favourites")
			}

			if len(favouritesList) == 0 {
				fmt.Println("No favourites found.")
				return
			}
			p := tea.NewProgram(btea.InitialFavouriteModel(favouritesList, &selectedFavourite))
			_, err = p.Run()
			if err != nil {
				fmt.Printf("Error displaying food selection: %v", err)
			}

			if selectedFavourite.Name == "" {
				return
			}

			favourites.FormatFavourite(selectedFavourite)
		},
	}
	return favouriteListCmd
}

func NewFavouriteDeleteListCmd(fileName string) *cobra.Command {
	var selectedFavourite favourites.Favourite
	var deleteListCmd = &cobra.Command{
		Use:   "delete",
		Short: "List all favourite foods, select do delete",
		Run: func(cmd *cobra.Command, args []string) {
			var isQuitting bool
			l := tea.NewProgram(btea.InitialSpinnerModel(&isQuitting))
			done := make(chan struct{})

			go func() {
				_, err := l.Run()
				if err != nil {
					fmt.Println("Error running spinner:", err)
				}
			}()

			var favouritesList []favourites.Favourite
			var err error

			go func() {
				favouritesList, err = favourites.RetrieveFavourites(fileName)
				close(done)
			}()

			<-done

			isQuitting = true
			l.Quit()

			if err != nil {
				fmt.Println("Error while retrieving favourites")
			}

			if len(favouritesList) == 0 {
				fmt.Println("No favourites found.")
				return
			}

			p := tea.NewProgram(btea.InitialFavouriteModel(favouritesList, &selectedFavourite))
			_, err = p.Run()
			if err != nil {
				fmt.Printf("Error displaying food selection: %v", err)
			}

			if selectedFavourite.Name == "" {
				return
			}

			fmt.Printf("Delete %s from favourites?", selectedFavourite.Name)
			var selectedAnswer string
			s := tea.NewProgram(btea.InitialConfirmationModel(&selectedAnswer))
			_, err = s.Run()
			if strings.ToLower(selectedAnswer) == "yes" {
				err = favourites.DeleteFavourite(selectedFavourite, fileName)

				if err != nil {
					fmt.Println("Unable to delete favourite")
				} else {
					fmt.Printf("%s has been deleted", selectedFavourite.Name)
				}
			}
		},
	}
	return deleteListCmd
}

func NewFavouriteSearchListCmd(fileName string) *cobra.Command {
	var search string
	var selectedFavourite favourites.Favourite

	var searchListCmd = &cobra.Command{
		Use:   "search",
		Short: "List all favourite foods for the search term",
		Run: func(cmd *cobra.Command, args []string) {
			var isQuitting bool
			l := tea.NewProgram(btea.InitialSpinnerModel(&isQuitting))
			done := make(chan struct{})

			go func() {
				_, err := l.Run()
				if err != nil {
					fmt.Println("Error running spinner:", err)
				}
			}()

			var favouritesList []favourites.Favourite
			var err error

			go func() {
				favouritesList, err = favourites.RetrieveFavourites(fileName)
				close(done)
			}()

			<-done

			isQuitting = true
			l.Quit()

			if err != nil {
				fmt.Println("Error while retrieving favourites")
			}

			favouritesList, err = favourites.FilterFavourites(favouritesList, search)

			if len(favouritesList) == 0 {
				fmt.Println("No favourites found.")
				return
			}

			if err != nil {
				fmt.Println("Error while retrieving favourites")
			}

			p := tea.NewProgram(btea.InitialFavouriteModel(favouritesList, &selectedFavourite))
			_, err = p.Run()
			if err != nil {
				fmt.Printf("Error displaying food selection: %v", err)
			}

			if selectedFavourite.Name == "" {
				return
			}

			favourites.FormatFavourite(selectedFavourite)
		},
	}

	searchListCmd.Flags().StringVarP(&search, "term", "k", "", "Search for favourite item")

	return searchListCmd
}
