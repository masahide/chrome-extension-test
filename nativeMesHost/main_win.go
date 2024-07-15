//go:build windows

package main

import (
	"fmt"
	"log"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

func register() {
	manifestPath, err := createManifest()
	if err != nil {
		log.Fatalf("Error writing manifest file: %v", err)
	}

	keyPath := `SOFTWARE\Google\Chrome\NativeMessagingHosts\` + Name
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		log.Fatalf("Error creating registry key: %v", err)
	}
	defer key.Close()

	if err = key.SetStringValue("", manifestPath); err != nil {
		log.Fatalf("Error setting registry value: %v", err)
	}

	fmt.Println("Native Messaging Host registered successfully.")
}

func openURLInBrowser(browser, profile, url string) {
	log.Printf("cmd /C start %s %s %s", browser, url, "--profile-directory="+profile)
	cmd := exec.Command("cmd", "/C", "start", browser, url, "--profile-directory="+profile)
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting browser: %v", err)
	}
}
