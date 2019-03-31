package main

import (
	"flag"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/klim0v/sequence-hashing/pkg/store"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	upgrader      = websocket.Upgrader{}
	connections   = map[string]*websocket.Conn{}
	muConnections sync.RWMutex
)

func main() {
	addr := flag.String("addr", "localhost:8080", "http service address")
	flag.Parse()
	go worker(store.NewClient(redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})))

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

func worker(storeClient store.Client) {
	for {
		result, err := storeClient.Pop()
		if err != nil {
			continue
		}
		for _, c := range connections {
			if err := c.WriteMessage(websocket.TextMessage, result.Hash); err != nil {
				log.Println(err)
			}
		}
		time.Sleep(3 * time.Second)
	}
}
