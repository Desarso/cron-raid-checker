package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	failure, err := checkRaidStatus()
	if err != nil {
		log.Fatal(err)
	}

	if failure {
		err := sendMessage("RAID Failure Alert", "A failure was detected in the RAID array.")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("RAID status is normal.")
	}
}

func sendMessage(subject, body string) error {
	// Replace with your bot's token and chat ID
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Retrieve environment variables
	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// URL to send the message
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Create the payload
	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
		return err
	}

	// Send the POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Check if the message was sent successfully
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully!")
	} else {
		var respBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			log.Fatalf("Failed to decode response body: %v", err)
			return err
		}
		fmt.Printf("Failed to send message. Status code: %d\n", resp.StatusCode)
		fmt.Println("Response:", respBody)
	}
	return nil
}

func checkRaidStatus() (bool, error) {
	cmd := exec.Command("mdadm", "--detail", "/dev/md0")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return false, err
	}

	// Check for RAID failure status in the output
	if strings.Contains(string(output), "failed") {
		fmt.Println("RAID failure detected!")
		return true, nil
	}

	return false, nil
}
