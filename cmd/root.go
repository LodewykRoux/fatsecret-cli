/*
Copyright © 2025 Lodewyk Roux lloderoux@protonmail.ch
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/LodewykRoux/fatsecret-cli/api"
	"github.com/LodewykRoux/fatsecret-cli/favourites"
	"github.com/LodewykRoux/fatsecret-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "fatsecret-cli",
	Short: "Fat Secret Cli is a non-official cli app to search for foods in the fat secret database.",
	Long:  `Use the command fatsecret-cli config to set up your initial client secret and client id.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fatsecret-cli.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".fatsecret-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func SetupCmd() {
	secretStorage := utils.GetSecretStoragePath()
	idStorage := utils.GetClientIdStoragePath()
	encryptionKey := utils.LoadEncryptionKey()
	favouritePath := favourites.GetFavouriteFile()

	// config
	configCmd := NewConfigCmd(secretStorage, idStorage, encryptionKey, StdInPasswordReader{})
	rootCmd.AddCommand(configCmd)

	// favourite
	favouriteCmd := NewFavouriteCmd(favouritePath)
	rootCmd.AddCommand(favouriteCmd)

	favouriteListCmd := NewFavouriteListCmd(favouritePath)
	favouriteCmd.AddCommand(favouriteListCmd)

	deleteListCmd := NewFavouriteDeleteListCmd(favouritePath)
	favouriteCmd.AddCommand(deleteListCmd)

	searchListCmd := NewFavouriteSearchListCmd(favouritePath)
	favouriteCmd.AddCommand(searchListCmd)

	accessToken := api.GetAccessToken(secretStorage, idStorage, encryptionKey)

	foodCmd := NewFoodCmd(accessToken, favouritePath, api.ApiFoodItemGetter{})
	rootCmd.AddCommand(foodCmd)
}
