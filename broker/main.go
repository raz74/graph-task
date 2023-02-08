package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rgamba/evtwebsocket"
	"log"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	conn := initConnectionToDestination()
	handler := &handler{conn: *conn}

	setupRoutes(handler)
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func initConnectionToDestination() *evtwebsocket.Conn {
	c := evtwebsocket.Conn{
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected")
		},
		OnMessage: func(msg []byte, conn *evtwebsocket.Conn) {
			fmt.Printf("Received sender: %s\n", msg)
		},
		OnError: func(err error) {
			fmt.Printf("** ERROR **\n%s\n", err.Error())
		},
	}
	// Connect
	// Dial sets up the connection with the remote
	// host provided in the url parameter.
	err := c.Dial("ws://localhost:8002/ws", "")
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

func setupRoutes(h *handler) {
	http.HandleFunc("/ws", h.wsEndpoint)
}

type handler struct {
	conn evtwebsocket.Conn
}

func (h *handler) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	h.listen(ws)
}

// Listening on the connection continuously
func (h *handler) listen(conn *websocket.Conn) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		log.Println(string(p))
		err = h.send(p)

	}
}

// Send sends a message through the connection.
func (h *handler) send(byteArr []byte) error {
	msg := evtwebsocket.Msg{
		Body: byteArr,
		Callback: func(resp []byte, conn *evtwebsocket.Conn) {
			fmt.Printf("Got response: %s\n", resp)
		},
	}

	return h.conn.Send(msg)
}
