package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var queue []*Player
var mu sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new conn")
	conn, _ := upgrader.Upgrade(w, r, nil)

	player := &Player{
		ID: uuid.NewString(),
		Name: r.URL.Query().Get("name"),
		Conn: conn,
	}

	addToQueue(player)
}

func addToQueue(p *Player) {
	mu.Lock()
	queue = append(queue, p)
	mu.Unlock()
	
	fmt.Println("player is added: ", p.Name)

	go waitForMatch(p)

//	p2 := &Player{
//		ID: "TestUser",
//		Name: "test12",
//		IsBot: true,
//	} 
//	startGame(p, p2)
}

func waitForMatch(p *Player) {
	timer := time.After(10000000 * time.Second) // for testing

	for {
		select {
		case <-timer:
			startGameWithBot(p)
			return
		default:
			mu.Lock()
			if len(queue) >= 2 {
				p1 := queue[0]
				p2 := queue[1]
				queue = queue[2:]
				mu.Unlock()

				startGame(p1, p2)
				return
			}
			mu.Unlock()
			time.Sleep(200 * time.Millisecond)
		}

	}
}

func removePlayer(p *Player) {
	p.Conn.Close()
}
