/*
Copyright © 2025 Lodewyk Roux lloderoux@protonmail.ch
*/
package cmd

import (
	"fmt"
	"log"
	"syscall"

	"golang.org/x/term"

	"github.com/LodewykRoux/fatsecret-cli/utils"
	"github.com/spf13/cobra"
)

type PasswordReader interface {
	ReadPassword() (string, error)
}

type StdInPasswordReader struct {
}

func (pr StdInPasswordReader) ReadPassword() (string, error) {
	pwd, error := term.ReadPassword(int(syscall.Stdin))
	return string(pwd), error
}

func NewConfigCmd(secretStorage string, idStorage string, encryptedKey []byte, pr PasswordReader) *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Sets the client secret and client id",
		Long:  `This command allows you to pass client secret and client id as a parameters.`,
		Run: func(cmd *cobra.Command, args []string) {
			var clientSecret string
			var clientId string
			if clientSecret == "" {
				var err error
				clientSecret, err = readClientSecret(pr)
				if err != nil {
					log.Fatalf("Error reading client secret: %v", err)
				}
			}

			if err := utils.StoreEncryptedClientSecret(clientSecret, secretStorage, encryptedKey); err != nil {
				log.Fatalf("Error saving client secret: %v", err)
			}

			fmt.Println("Client secret saved securely.")

			if clientId == "" {
				var err error
				clientId, err = readClientId(pr)
				if err != nil {
					log.Fatalf("Error reading client id: %v", err)
				}
			}

			if err := utils.StoreEncryptedClientId(clientId, idStorage, encryptedKey); err != nil {
				log.Fatalf("Error saving client id: %v", err)
			}

			fmt.Println("Client id saved securely.")
		},
	}
	return configCmd
}

func readClientSecret(pr PasswordReader) (string, error) {
	fmt.Print("Enter Client Secret: ")
	bytePassword, err := pr.ReadPassword()
	fmt.Println() // Newline after input
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

func readClientId(pr PasswordReader) (string, error) {
	fmt.Print("Enter Client Id: ")
	byteId, err := pr.ReadPassword()

	fmt.Println() // Newline after input
	if err != nil {
		return "", err
	}
	return string(byteId), nil
}
