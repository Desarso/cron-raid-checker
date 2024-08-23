package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := sendMessage("Testing", "Testing send message function.")
	if err != nil {
		log.Fatal(err)
	}
}

func sendMessage(subject, body string) error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return err
	}

	// Retrieve environment variables
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN not set in environment")
		return fmt.Errorf("BOT_TOKEN not set")
	}

	// Define chat IDs
	chatIDs := []string{"6611371097", "6995936214"}
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// URL to send the message
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Iterate over chat IDs and send messages
	for _, chatID := range chatIDs {
		// Create the payload
		payload := map[string]string{
			"chat_id": chatID,
			"text":    message,
		}

		// Convert payload to JSON
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal payload for chat ID %s: %v", chatID, err)
			continue
		}

		// Send the POST request
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("Failed to send request to chat ID %s: %v", chatID, err)
			continue
		}
		defer resp.Body.Close()

		// Check if the message was sent successfully
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Message sent successfully to chat ID %s!\n", chatID)
		} else {
			var respBody map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				log.Printf("Failed to decode response body for chat ID %s: %v", chatID, err)
			}
			//fmt.Printf("Failed to send message to chat ID %s. Status code: %d\n", chatID, resp.StatusCode)
			//fmt.Println("Response:", respBody)
		}
	}
	return nil
}
