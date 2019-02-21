package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/klim0v/sequence-hashing/pkg/redis"
	"log"
	"net/http"
	"sync"
)

var (
	addr          = flag.String("addr", "localhost:8080", "http service address")
	upgrader      = websocket.Upgrader{}
	connections   = map[string]*websocket.Conn{}
	muConnections sync.RWMutex
)

func main() {
	flag.Parse()
	go worker()

	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	done := make(chan struct{})

	go func() {
		defer close(done)
		identification := conn.RemoteAddr().String()
		muConnections.RLock()
		connections[identification] = conn
		muConnections.RUnlock()

		defer func() {
			muConnections.RLock()
			delete(connections, identification)
			muConnections.RUnlock()
		}()

		for {
			if _, _, err := conn.NextReader(); err != nil {
				_ = conn.Close()
				break
			}
		}
	}()

	<-done
}

func worker() {
	for {
		result, err := redis.NewClient().Pop()
		if err != nil {
			continue
		}
		for _, c := range connections {
			if err := c.WriteMessage(websocket.TextMessage, result.Hash); err != nil {
				log.Println(err)
			}
		}
		//time.Sleep(3 * time.Second)
	}
}
