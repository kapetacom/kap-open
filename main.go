package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"runtime"

	"gopkg.in/square/go-jose.v2/jwt"
)

type Authentication struct {
	AccessToken string `json:"access_token"`
	BaseURL     string `json:"base_url"`
}

func main() {
	// Get the port number from the command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: kap-open <port> <base path>")
		return
	}
	port := os.Args[1]
	path := os.Args[2]

	// Run 'kap registry ls' and discard its output
	registryCmd := exec.Command("kap", "registry", "ls")
	//registryCmd.Stdout = os.Stdout
	registryCmd.Run()

	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting user's home directory:", err)
		return
	}
	homeDir := usr.HomeDir

	// Construct the path to authentication.json
	authFilePath := homeDir + "/.kapeta/authentication.json"

	// Read authentication.json and extract access token
	authFile, err := os.Open(authFilePath)
	if err != nil {
		fmt.Println("Error opening authentication.json:", err)
		return
	}
	defer authFile.Close()

	var authData Authentication
	err = json.NewDecoder(authFile).Decode(&authData)
	if err != nil {
		fmt.Println("Error decoding authentication.json:", err)
		return
	}

	accessToken := authData.AccessToken
	contextString := getContext(accessToken)
	// url encode the context string
	// Construct the URL and open it in the default web browser
	if path == "" {
		path = "/kapeta"
	}
	url := fmt.Sprintf("http://localhost:"+port+"/%s/?token=%s&context=%s", path, accessToken, url.QueryEscape(contextString))
	err = openBrowser(url) // Use 'xdg-open' for Linux, use 'open' for macOS
	if err != nil {
		fmt.Println("Error opening browser:", err)
		return
	}

	fmt.Println("Browser opened to:", url)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

func getContext(tokenString string) string {

	var claims map[string]interface{} // generic map to store parsed token

	// decode JWT token without verifying the signature
	token, _ := jwt.ParseSigned(tokenString)
	_ = token.UnsafeClaimsWithoutVerification(&claims)
	// Replace with your actual JWT token

	// Print the claims as json
	claimsJSON, _ := json.MarshalIndent(claims, "", "  ")

	// convert to TokenInfo struct
	var tokenInfo TokenInfo
	err := json.Unmarshal(claimsJSON, &tokenInfo)
	if err != nil {
		log.Fatal("Error unmarshalling claims:", err)
	}

	// encode the context as a URL parameter
	context, _ := json.Marshal(tokenInfo.Contexts[0])
	return string(context)

}

type TokenInfo struct {
	AuthID   string    `json:"auth_id"`
	AuthType string    `json:"auth_type"`
	Contexts []Context `json:"contexts"`
	Exp      int64     `json:"exp"`
	Iat      int64     `json:"iat"`
	Iss      string    `json:"iss"`
	Purpose  string    `json:"purpose"`
	Scopes   []string  `json:"scopes"`
	Sub      string    `json:"sub"`
	Type     string    `json:"type"`
}

type Context struct {
	Handle string   `json:"handle"`
	ID     string   `json:"id"`
	Type   string   `json:"type"`
	Scopes []string `json:"scopes"`
}
