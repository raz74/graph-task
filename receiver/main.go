package main

import (
	"fmt"
	"github.com/rgamba/evtwebsocket"
	"io"
	"log"
	"net/http"
)

func main() {
	conn := initConnectionToBroker()
	handler := &handler{conn: *conn}

	setupRoutes(handler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initConnectionToBroker() *evtwebsocket.Conn {
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
	err := c.Dial("ws://localhost:8001/ws", "")
	if err != nil {
		log.Fatal(err)
	}
	return &c
}

func setupRoutes(ws *handler) {
	http.HandleFunc("/message", ws.receiveHandler)
}

type handler struct {
	conn evtwebsocket.Conn
}

func (h *handler) receiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		byteArr, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusBadRequest)
		}

		err = h.send(byteArr)
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
func (h *handler) send(byteArr []byte) error {
	msg := evtwebsocket.Msg{
		Body: byteArr,
	}

	return h.conn.Send(msg)
}
