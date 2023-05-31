package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}

	ErrEventNotSupported = errors.New("this event type is not supported")
)

//checkOrigin will the origin of the request

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:3000":
		return true
	default:
		return false
	}
}

// Hold references to all the clients
type Manager struct {
	clients ClientList

	//Using this to be able to lock state before editing clients
	sync.RWMutex

	// handlers are functions that are used to handle events
	handlers map[string]EventHandler
}

// Constructor for the Manager
func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
}

// routeEvent is used to make sure the correct venet goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	//Checking if handler is present
	if handler, ok := m.handlers[event.Type]; ok {
		//executing the handler
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

// serveWS is a HTTP handler that has manager to allow connection
func (m *Manager) ServeWS(c *gin.Context) {

	jwt := c.Query("jwt")
	if jwt == "" {
		log.Printf("Error invalid token!: %v", jwt)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
		return
	}

	//Verify the jwt token
	if err := validateToken(jwt); err != nil {
		log.Println("Error invalid token!")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	log.Println("New Connection...!")

	//Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//Creating a new Client
	client := NewClient(conn, m)

	//Adding the newly created client to the manager
	m.addClient(client)

	//Start the read/write processes

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	//Add client
	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	//Check if client exists, then delete it
	if _, ok := m.clients[client]; ok {
		//close connection
		client.connection.Close()

		//remove
		delete(m.clients, client)
	}
}

func validateToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte("my_secret"), nil
	})

	if err != nil || !token.Valid {
		return err
	}

	return nil
}
