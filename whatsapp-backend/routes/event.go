package routes

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	//Message type sent
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// This function is used to affect messages on the socket and trigger depending on the type
type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
)

// Payload sent in the send_message event
type SendMessageEvent struct {
	Message  string `json:"message"`
	From     int    `json:"from"`
	To       int    `json:"to"`
	Chatroom string `json:"chatroom"`
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

// SendMessageHandler will send out a message to all other participants in the chat
func SendMessageHandler(event Event, c *Client) error {
	//Marshal the payload
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// Prepare an Outgoing message to others
	var broadMessage NewMessageEvent
	broadMessage.Sent = time.Now()
	broadMessage.Message = chatevent.Message
	broadMessage.From = chatevent.From
	broadMessage.To = chatevent.To
	broadMessage.Chatroom = chatevent.Chatroom

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	//Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage

	//Broadcast to all other clients
	for client := range c.manager.clients {
		//Only send to clients inside the same chatroom
		if client.chatroom == c.chatroom {
			client.egress <- outgoingEvent
		}
	}

	return nil
}
