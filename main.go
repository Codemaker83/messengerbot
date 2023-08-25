package main

import (
	"log"

	"messengerbot/chatbot"
)

func main() {
	cfg, err := chatbot.GetConfig("config.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv := chatbot.New(cfg)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

