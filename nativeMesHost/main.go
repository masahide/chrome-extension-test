package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

const (
	extensionID = "hkmknjbjeeljdgpfkakimhcdfaghlofc"
	Name        = "com.my_company.my_application"
)

type Message struct {
	URL string `json:"url"`
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
	nativeByteOder = getNativeByteOrder()
}

func readMessage() (Message, error) {
	var length int32

	err := binary.Read(os.Stdin, nativeByteOder, &length)
	if err != nil {
		if err == io.EOF {
			runtime.Goexit()
		}
		return Message{}, fmt.Errorf("failed to read length: %w", err)
	}

	messageBytes := make([]byte, length)
	_, err = io.ReadFull(os.Stdin, messageBytes)
	if err != nil {
		return Message{}, fmt.Errorf("failed to read message: %w", err)
	}
	logPrintf(fmt.Sprintf("read length: %v messsage:%s\n", length, messageBytes))

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
	err = binary.Write(stdout, nativeByteOder, length)
	if err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	n, err := stdout.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	logPrintf("write length: %v messsage:%s\n", length, messageBytes)
	if n != int(length) {
		return io.ErrShortWrite
	}
	if err := stdout.Flush(); err != nil {
		logPrintf("failed to flush buffer: %s", err)
		return fmt.Errorf("failed to flush buffer: %w", err)
	}
	logPrintf("sent message: %s\n", string(messageBytes))

	return nil
}

var (
	exePath        string
	nativeByteOder binary.ByteOrder
	fname          string
)

func init() {
	var err error
	exePath, err = os.Executable()
	if err != nil {
		errPrintf("Error getting executable path: %v\n", err)
		return
	}
	fname = filepath.Join(filepath.Dir(exePath), "test.txt")
}

func register() {

	manifest := Manifest{
		Name:           Name,
		Description:    "Example Native Messaging Host",
		Path:           exePath,
		Type:           "stdio",
		AllowedOrigins: []string{"chrome-extension://" + extensionID + "/"},
	}

	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		errPrintf("Error marshaling manifest: %v\n", err)
		return
	}

	manifestPath := filepath.Join(filepath.Dir(exePath), "manifest.json")

	err = os.WriteFile(manifestPath, manifestBytes, 0644)
	if err != nil {
		errPrintf("Error writing manifest file: %v\n", err)
		return
	}

	keyPath := `SOFTWARE\Google\Chrome\NativeMessagingHosts\` + Name
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		errPrintf("Error creating registry key: %v\n", err)
		return
	}
	defer key.Close()

	err = key.SetStringValue("", manifestPath)
	if err != nil {
		errPrintf("Error setting registry value: %v\n", err)
		return
	}

	fmt.Println("Native Messaging Host registered successfully.")
}

func errPrintf(format string, a ...any) {
	logPrintf(format, a...)
	fmt.Fprintf(os.Stderr, format, a...)
}

func logPrintf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		runtime.Goexit()
	}
	defer f.Close()
	_, err = f.Write([]byte(s))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		runtime.Goexit()
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "register" {
		register()
		return
	}

	for {
		message, err := readMessage()
		if err != nil {
			errPrintf("Error reading message: %v\n", err)
			return
		}
		logPrintf("Received message: %v\n", message)
		// Open the URL in Firefox
		cmd := exec.Command("cmd", "/C", "start", "brave", message.URL)
		err = cmd.Start()
		if err != nil {
			errPrintf("Error starting browser: %v\n", err)
			return
		}
		if err := sendMessage(Message{URL: "Browser opened successfully."}); err != nil {
			errPrintf("Error sending message: %v\n", err)
			return
		}
	}
}
