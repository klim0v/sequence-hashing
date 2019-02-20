package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/klim0v/sequence-hashing/pkg/entity"
	"github.com/klim0v/sequence-hashing/pkg/redis"
	"log"
	"net/http"
	"sync"
	"time"
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
	pubsub := redis.NewClient().Subscribe()
	if _, err := pubsub.Receive(); err != nil {
		log.Println(err)
	}
	for msg := range pubsub.Channel() {
		fmt.Println(msg.Payload)
		var data entity.Result
		if err := data.UnmarshalBinary([]byte(msg.Payload)); err != nil {
			continue
		}
		for _, c := range connections {
			if err := c.WriteMessage(websocket.TextMessage, data.Hash); err != nil {
				log.Println(err)
			}
		}
		time.Sleep(3 * time.Second)
	}
}
