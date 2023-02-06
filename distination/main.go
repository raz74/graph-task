package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8002", nil))
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi  second Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	listen(ws)
}

// Listening on the connection continuously
func listen(conn *websocket.Conn) {
	var count int
	var volume int
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		count += 1
		volume += len(p)
		log.Printf("count: %v volume: %v\n", count, volume)
	}

}
