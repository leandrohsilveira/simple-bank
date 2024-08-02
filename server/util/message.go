package util

import "encoding/json"

type MessageType string

const (
	MessageSuccess MessageType = "success"
	MessageError   MessageType = "error"
)

type Message struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

func (m Message) Json() (result string, err error) {
	bytes, err := json.Marshal(map[string]Message{
		"showMessage": m,
	})

	result = string(bytes)

	return
}
