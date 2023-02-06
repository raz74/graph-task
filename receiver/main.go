package main

import (
	"fmt"
	"github.com/rgamba/evtwebsocket"
	"io"
	"log"
	"net/http"
)

func main() {
	//init websocket connection
	conn := initConnectionToBroker()
	handler := &client{conn: *conn}

	setupRoutes(handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initConnectionToBroker() *evtwebsocket.Conn {
	c := evtwebsocket.Conn{
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected")
		},
		OnMessage: func(msg []byte, conn *evtwebsocket.Conn) {
			fmt.Printf("Received producer: %s\n", msg)
		},
		OnError: func(err error) {
			fmt.Printf("** ERROR **\n%s\n", err.Error())
		},
	}
	// Connect
	// Dial sets up the connection with the remote
	// host provided in the url parameter.
	err := c.Dial("ws://localhost:8001/ws", "")
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

func setupRoutes(ws *client) {
	http.HandleFunc("/message", ws.produce)
}

type client struct {
	conn evtwebsocket.Conn
}

func (ws *client) produce(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		byteArr, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusBadRequest)
		}

		err = ws.send(byteArr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Send sends a message through the connection.
func (ws *client) send(byteArr []byte) error {
	msg := evtwebsocket.Msg{
		Body: byteArr,
		Callback: func(resp []byte, conn *evtwebsocket.Conn) {
			fmt.Printf("Got response: %s\n", resp)
		},
	}

	return ws.conn.Send(msg)
}
