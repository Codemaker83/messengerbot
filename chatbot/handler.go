package chatbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"net/http"
	
	"gopkg.in/yaml.v3"
)

func deferClose(obj io.Closer, err *error) {
	closeErr := obj.Close()

	if *err == nil && closeErr != nil && !errors.Is(closeErr, os.ErrClosed) {
		*err = closeErr
		return
	}

	if closeErr != nil && !errors.Is(closeErr, os.ErrClosed) {
		*err = fmt.Errorf("%s.\n\nWhile closing file after previous error: %w", (*err).Error(), closeErr)
	}
}

func checkConfig(cfg Config) (Config, error) {
	if cfg.VerifyToken == "" {
		return cfg, fmt.Errorf("verify required")
	}
	if cfg.AccessToken == "" {
		return cfg, fmt.Errorf("access_token required")
	}
	if cfg.GraphQlURL == "" {
		return cfg, fmt.Errorf("graphql_url required")
	}
	
	if cfg.ServerIP == "" {
		cfg.ServerIP = "localhost"
	}
	return cfg, nil
}

// GetConfig gets config values from a yaml file
func GetConfig(configFile string) (cfg Config, err error) {
	var f *os.File
	f, err = os.Open(configFile)
	if err != nil {
		return
	}
	defer deferClose(f, &err)
	
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)

	err = decoder.Decode(&cfg)
	cfg, err = checkConfig(cfg)
	return
}

// chatHandler handles Webhook server
func (c *Config) chatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// log.Printf("%+v\n", r.Header)

	if r.Method == http.MethodGet {
		verifyToken := r.URL.Query().Get("hub.verify_token")

		// verifying incoming request
		if verifyToken != c.VerifyToken {
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

		var message MessageEvent
		if err := json.Unmarshal(body, &message); err != nil {
			log.Printf("error while unmarshaling request body: %v", err)
			return
		}
		log.Printf("%+v", message)

		// send response
		err = c.sendMessage(message.Entry[0].Messaging[0].Sender.ID, "Thanks for communicating to u! We will be in touch soon.")
		if err != nil {
			log.Printf("error sending message: %v", err)
		}

		return
	}
	
	log.Printf("invalid http method: should be get or post")
	return
}

// sendMessage sends a message to end-user
func (c *Config) sendMessage(senderId, message string) error {
	var request Message
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
	url := fmt.Sprintf("%s/me/messages?access_token=%s", c.GraphQlURL, c.AccessToken)
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
