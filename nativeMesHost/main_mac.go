//go:build darwin

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var browserOSAppNameMap = map[string]string{
	"msedge":  "edge",
	"brave":   "Brave Browser",
	"firefox": "firefox",
}

func register() {
	manifestPath, err := createManifest()
	if err != nil {
		log.Fatalf("Error writing manifest file: %v", err)
	}

	chromeDir := filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Google", "Chrome", "NativeMessagingHosts")
	if err := os.MkdirAll(chromeDir, 0755); err != nil {
		log.Fatalf("Error creating Chrome directory: %v", err)
	}

	if err := os.Symlink(manifestPath, filepath.Join(chromeDir, Name+".json")); err != nil {
		log.Fatalf("Error creating symlink: %v", err)
	}

	fmt.Println("Native Messaging Host registered successfully.")
}

func openURLInBrowser(browser, profile, url string) {
	// Open the URL in the specified browser
	appName, ok := browserOSAppNameMap[browser]
	if !ok {
		log.Printf("Unsupported browser: %s", browser)
		return
	}

	// Open the URL in the specified browser
	cmd := exec.Command("open", "-a", appName, "--args", "--profile-directory="+profile, url)
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting browser: %v", err)
	}
}
