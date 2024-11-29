package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	ChannelID string    `json:"channel_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func (msg Message) ToBytes() ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type WriteMessage struct {
	RequestID string `json:"request_id"`
	//GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	Message   string `json:"message"`
}

func ToWriteMessage(msg []byte) (WriteMessage, error) {
	var wMsg WriteMessage
	err := json.Unmarshal(msg, &wMsg)
	if err != nil {
		return WriteMessage{}, fmt.Errorf("error deserializing message: %w", err)
	}

	return wMsg, nil
}
