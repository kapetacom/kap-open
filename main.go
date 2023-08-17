package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

type Authentication struct {
	AccessToken string `json:"access_token"`
}

func main() {
	// Get the port number from the command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: kap-open <port>")
		return
	}
	port := os.Args[1]
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

	// Construct the URL and open it in the default web browser
	url := fmt.Sprintf("http://localhost:"+port+"/kapeta/?token=%s", accessToken)
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
