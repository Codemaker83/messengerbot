package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"messengerbot/chatbot"
)

const (
	// VerifyToken use to verify the incoming request
	VerifyToken = "PcvLctz8wOqGb0D3W9D5ksevJkdhTmACKeDUs62p"
	// AccessToken use to access Messenger API
	AccessToken = "EAAzJAA5FJkoBOxuLauRedPycxZC3vhgmU2U0140HzonUNyL1eLDJcJH2Oq37hOxej5AzCALLYwLPutAFNgpRTvsHzmu2opjbIuS6HQOOXndIJwbuNwXXdHgNGHfbtDmZBMhXKuCgYwg63wkAwcjNsgL1DEPA4YvczrZBTvFB7zfFLQQPCIpV19e8RRisMkZCw0ntngD9MEd4aTUZD"
	// GraphQlURL is a base URL v17.0 for Messenger API
	GraphQlURL = "https://graph.facebook.com/v17.0"
)

// chatHandler handles Webhook server
func chatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	log.Printf("%+v\n", r.Header)

	if r.Method == http.MethodGet {
		verifyToken := r.URL.Query().Get("hub.verify_token")

		// verifying incoming request
		if verifyToken != VerifyToken {
			log.Printf("invalid verification token: %s", verifyToken)
			return
		}

		// returning challenge
		if _, err := w.Write([]byte(r.URL.Query().Get("hub.challenge"))); err != nil {
			log.Printf("error writing response body: %v", err)
		}

		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading request body: %v", err)
			return
		}

		var message chatbot.MessageEvent
		if err := json.Unmarshal(body, &message); err != nil {
			log.Printf("error while unmarshaling request body: %v", err)
			return
		}
		log.Printf("%+v", message)

		// send response
		err = sendMessage(message.Entry[0].Messaging[0].Sender.ID, "Thanks for communicating to u! We will be in touch soon.")
		if err != nil {
			log.Printf("error sending message: %v", err)
		}

		return
	}
	
	log.Printf("invalid http method: should be get or post")
	return
}

// sendMessage sends a message to end-user
func sendMessage(senderId, message string) error {
	var request chatbot.Message
	request.Recipient.ID = senderId
	request.Message.Text = message

	if len(message) == 0 {
		return errors.New("empty message")
	}

	// marshal request data
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("while marshaling request: %w", err)
	}

	// create http request
	url := fmt.Sprintf("%s/me/messages?access_token=%s", GraphQlURL, AccessToken)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed wrap request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	
	// send request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed send request: %w", err)
	}
	defer res.Body.Close()

	// print response
	log.Printf("Response:\n%#v", res)
	
	return nil
}

func main() {
	// create the handler
	handler := http.NewServeMux()
	handler.HandleFunc("/", chatHandler)

	// configure http server
	srv := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("localhost:%d", 3030),
	}

	// start http server
	log.Printf("http server listening at %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

