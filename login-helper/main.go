package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type loginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type loginResponse struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresIn        int    `json:"expiresIn"`
	RefreshExpiresIn int    `json:"refreshExpiresIn"`
	TfaKey           string `json:"tfaKey"`
	AccessMethod     string `json:"accessMethod"`
	LoginType        string `json:"loginType"`
}

func main() {
	// Define command-line flags
	account := flag.String("account", "", "Account name (required)")
	password := flag.String("password", "", "Password (required)")
	apiURL := flag.String("api-url", "https://api.bambulab.com", "API URL")
	flag.Parse()

	// Validate required flags
	if *account == "" || *password == "" {
		fmt.Println("Error: --account and --password flags are required")
		flag.Usage()
		os.Exit(1)
	}

	httpClient := &http.Client{}
	l := loginRequest{
		Account:  *account,
		Password: *password,
		Code:     "",
	}

	respLogin, err := doLogin(l, *apiURL, httpClient)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		os.Exit(1)
	}

	// If two-factor authentication (TFA) is required
	if respLogin.LoginType != "" {
		fmt.Print("Enter access code: ")
		var accessCode string
		_, err := fmt.Scanln(&accessCode)
		if err != nil {
			fmt.Printf("TFA login failed: %v\n", err)
			os.Exit(1)
		}

		l = loginRequest{
			Account:  *account,
			Password: "",
			Code:     accessCode,
		}

		respLogin, err = doLogin(l, *apiURL, httpClient)
		if err != nil {
			fmt.Printf("TFA login failed: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Login successful!")
	fmt.Println("Access Token:", respLogin.AccessToken)
}

func doLogin(l loginRequest, url string, httpClient *http.Client) (*loginResponse, error) {
	marshal, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(marshal)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/user-service/user/login", url), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	respLogin := loginResponse{}
	err = json.Unmarshal(responseBody, &respLogin)
	if err != nil {
		return nil, err
	}

	return &respLogin, nil
}
