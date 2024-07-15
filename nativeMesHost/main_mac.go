//go:build darwin

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

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
	cmd := exec.Command("open", url)
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting browser: %v", err)
	}
}
