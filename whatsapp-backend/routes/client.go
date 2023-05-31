package routes

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (

	//pongWait is how long we will await a pong response from the client
	pongWait = 10 * time.Second

	pingInterval = (pongWait * 9) / 10
)

// ClientList is just map for looking up client
type ClientList map[*Client]bool

// Client is a websocket client,
type Client struct {
	connection *websocket.Conn

	//manager used to manage the Client
	manager *Manager

	//egress is used to avoid concurrent writeon the websocket
	egress chan Event

	//Client room
	chatroom string
}

// Constructor for the client
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		chatroom:   "",
	}
}

// function for reading the messages and handling them

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	//setting max size of messages in bytes
	c.connection.SetReadLimit(1024)

	//Configure wait time for Pong response, using current time+ pongwait
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	//Confguring how to handle Pong response
	c.connection.SetPongHandler(c.pongHandler)

	for {
		//reading the message from queue
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		//Marshal the incoming messages into an Event struct
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		//Route the event
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handeling message: ", err)
		}

		// log.Println("MessageType: ", messageType)
		// log.Println("Payload: ", string(payload))

		// for wsclient := range c.manager.clients {
		// 	wsclient.egress <- payload
		// }
	}
}

func (c *Client) pongHandler(pongmsg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {
	//create a ticker
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}

				// close the go routine
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
			log.Println("sent message")

		case <-ticker.C:
			log.Println("ping")
			//Send the Ping
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return //return to break this goroutine triggering cleanup
			}

		}
	}
}
