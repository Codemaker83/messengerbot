package chatbot

// node contains information about a node present in the message event
type node struct {
	ID string `json:"id"`
}

type message struct {
	Mid  string `json:"mid"`
	Text string `json:"text"`
}

type messaging struct {
	Sender    node    `json:"sender"`
	Recipient node    `json:"recipient"`
	Timestamp int64   `json:"timestamp"`
	Message   message `json:"message"`
}

type entry struct {
	ID        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []messaging `json:"messaging"`
}

// MessageEvent contains message event data
type MessageEvent struct {
	Object string   `json:"object"`
	Entry  []entry  `json:"entry"`
}

// Message contains data for message to be sent
type Message struct {
	Recipient node `json:"recipient"`
	Message   struct {
		Text string `json:"text"`
	} `json:"message"`
}

// Config contains all configuration for a chatbot
type Config struct {
	ServerIP    string `yaml:"server_ip,omitempty"`    // ServerIP is the IP address or url of the server running the chatbot (defaults to localhost)
	Port        int    `yaml:"port,omitempty"`         // Port used to listen and serve
	VerifyToken string `yaml:"verify,omitempty"`       // VerifyToken is the token used to verify the incoming request
	AccessToken string `yaml:"access_token,omitempty"` // AccessToken is Messenger API's token
	GraphQlURL  string `yaml:"graphql_url,omitempty"`  // GraphQlURL is Messenger API's base URL 
}
