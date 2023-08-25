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
