package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/LodewykRoux/fatsecret-cli/utils"
)

var fileName = "token.json"

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenFile struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     string `json:"expires"`
}

func GetAccessToken(secretStorage string, idStorage string, encryptionKey []byte) string {
	access := retrieveAccessToken(encryptionKey)
	if access != "" {
		return access
	}

	clientID, err := utils.GetDecryptedClientId(idStorage, encryptionKey)
	if err != nil {
		log.Fatalf("Error retrieving client ID.")
	}
	if clientID == "" {
		return ""
	}

	clientSecret, err := utils.GetDecryptedClientSecret(secretStorage, encryptionKey)
	if err != nil {
		log.Fatalf("Error retrieving client secret.")
	}
	if clientSecret == "" {
		return ""
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("scope", "basic")

	req, err := http.NewRequest("POST", "https://oauth.fatsecret.com/connect/token", strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		errorText := string(body)
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			errorText = errorResp.Error
		}
		log.Fatalf("Error: %s. Failed to get access token: %s", http.StatusText(resp.StatusCode), errorText)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		panic(err)
	}

	encryptedToken, err := utils.Encrypt(tokenResp.AccessToken, encryptionKey)

	if err != nil {
		panic(err)
	}

	saveAccessToken(tokenResp, encryptedToken)

	return tokenResp.AccessToken
}

func saveAccessToken(tokenResp TokenResponse, encryptedToken string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	expirationTime := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	tokenFile := TokenFile{AccessToken: encryptedToken, TokenType: tokenResp.TokenType, Expires: expirationTime.Format(time.RFC3339)}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tokenFile); err != nil {
		panic(err)
	}
}

func retrieveAccessToken(encryptionKey []byte) string {
	file, err := os.Open(fileName)
	if err != nil {
		return ""
	}
	defer file.Close()

	var tokenResp TokenFile
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tokenResp); err != nil {
		return ""
	}

	expirationTimeParsed, err := time.Parse(time.RFC3339, tokenResp.Expires)

	if err != nil {
		panic(err)
	}

	if time.Now().After(expirationTimeParsed) {
		return ""
	}

	decryptedToken, err := utils.Decrypt(tokenResp.AccessToken, encryptionKey)
	if err != nil {
		return ""
	}

	return decryptedToken
}
