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
	err := sendMessage("Testing", "More testing")
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
	chatIDs, err := getAllChatIds()
	if err != nil {
		log.Fatal(err)
		return err
	}
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

func getAllChatIds() ([]string, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return nil, err
	}

	// Retrieve environment variables
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN not set in environment")
	}

	// URL to get updates
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", botToken)

	// Send the GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get updates: status code %d", resp.StatusCode)
	}

	// Parse the response body
	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		log.Printf("Failed to decode response body: %v", err)
		return nil, err
	}

	// Extract chat IDs
	chatIDSet := make(map[string]struct{})
	result, ok := respBody["result"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	for _, v := range result {
		message, ok := v.(map[string]interface{})["message"].(map[string]interface{})
		if !ok {
			continue
		}
		chat, ok := message["chat"].(map[string]interface{})
		if !ok {
			continue
		}
		chatID, ok := chat["id"].(float64)
		if !ok {
			continue
		}
		// Add chat ID to the set to remove duplicates
		chatIDSet[fmt.Sprintf("%.0f", chatID)] = struct{}{}

	}
	var chatIds []string
	for id := range chatIDSet {
		chatIds = append(chatIds, id)
	}

	return chatIds, nil
}
