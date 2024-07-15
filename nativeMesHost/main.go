package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	extensionID = "hkmknjbjeeljdgpfkakimhcdfaghlofc"
	Name        = "com.my_company.my_application"
	LogFilename = "testlog.txt"
)

var (
	manifest = Manifest{
		Name:           Name,
		Description:    "OpenAltBrowser allows you to effortlessly open the current tab in another browser. Enhance your browsing experience by seamlessly switching between browsers with just a click. Perfect for developers, testers, and anyone who uses multiple web browsers",
		Path:           exePath,
		Type:           "stdio",
		AllowedOrigins: []string{"chrome-extension://" + extensionID + "/"},
	}
)

type Message struct {
	Browser string `json:"browser"`
	Profile string `json:"profile"`
	URL     string `json:"url"`
}

type Manifest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Path           string   `json:"path"`
	Type           string   `json:"type"`
	AllowedOrigins []string `json:"allowed_origins"`
}

func getNativeByteOrder() binary.ByteOrder {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(0x01020304))
	return map[bool]binary.ByteOrder{
		true:  binary.BigEndian,
		false: binary.LittleEndian,
	}[b[0] == 0x01]
}

func init() {
	var err error
	nativeByteOrder = getNativeByteOrder()
	exePath, err = os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}
	fname := filepath.Join(filepath.Dir(exePath), LogFilename)
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v", err)
		os.Exit(1)
	}
	w := io.MultiWriter(f, os.Stderr)
	// log to file
	log.SetOutput(w)

}

func readMessage() (Message, error) {
	var length int32

	err := binary.Read(os.Stdin, nativeByteOrder, &length)
	if err != nil {
		if err == io.EOF {
			return Message{}, err
		}
		return Message{}, fmt.Errorf("failed to read length: %w", err)
	}

	messageBytes := make([]byte, length)
	_, err = io.ReadFull(os.Stdin, messageBytes)
	if err != nil {
		return Message{}, fmt.Errorf("failed to read message: %w", err)
	}
	log.Printf("read length: %v messsage:%s", length, messageBytes)

	var message Message
	err = json.Unmarshal(messageBytes, &message)
	if err != nil {
		return Message{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return message, nil
}

func sendMessage(message Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	stdout := bufio.NewWriter(os.Stdout)
	length := int32(len(messageBytes))
	err = binary.Write(stdout, nativeByteOrder, length)
	if err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	n, err := stdout.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	log.Printf("write length: %v messsage:%s", length, messageBytes)
	if n != int(length) {
		return io.ErrShortWrite
	}
	if err := stdout.Flush(); err != nil {
		log.Fatalf("failed to flush buffer: %s", err)
	}
	log.Printf("sent message: %s", string(messageBytes))

	return nil
}

var (
	exePath         string
	nativeByteOrder binary.ByteOrder
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "register" {
		register()
		return
	}

	for {
		message, err := readMessage()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("Error reading message: %v", err)
			return
		}
		log.Printf("Received message: %v", message)
		// Open the URL in the browser
		openURLInBrowser(message.Browser, message.Profile, message.URL)
		if err := sendMessage(Message{URL: "Browser opened successfully."}); err != nil {
			log.Fatalf("Error sending message: %v", err)
			return
		}
	}
}

func createManifest() (string, error) {
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling manifest: %w", err)
	}

	manifestPath := filepath.Join(filepath.Dir(exePath), "manifest.json")

	err = os.WriteFile(manifestPath, manifestBytes, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing manifest file: %w", err)
	}
	return manifestPath, nil
}
