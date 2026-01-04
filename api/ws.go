package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var queue []*Player
var mu sync.Mutex

var players = make(map[string]*Player)
var playersMU sync.Mutex

var activeGames = make(map[string]*Game)
var activeGamesMu sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new conn")
	conn, _ := upgrader.Upgrade(w, r, nil)
	token := r.URL.Query().Get("token")
	name := r.URL.Query().Get("name")

	playersMU.Lock()
	p, exists := players[token]
	playersMU.Unlock()

	if !exists {
		if token == "" {
			token = uuid.NewString()
			conn.WriteJSON(Message{
				Type: "token",
				Data: token,
			})
		}
		player := &Player{
			ID:    token,
			Name:  name,
			Conn:  conn,
			IsBot: false,
		}
		playersMU.Lock()
		players[token] = player
		playersMU.Unlock()

		go listenLobby(player)
	} else {
		p.WriteMu.Lock()
		p.Name = name
		p.Conn = conn
		p.WriteMu.Unlock()

		activeGamesMu.Lock()
		game, inGame := activeGames[token]
		activeGamesMu.Unlock()

		if inGame {
			playerNum := 1
			if game.Players[1].ID == p.ID {
				playerNum = 2
			}

			p.Conn.WriteJSON(Message{
				Type: "start",
			})

			p.WriteMu.Lock()
			p.Conn.WriteJSON(Message{
				Type: "state",
				Data: game,
			})
			p.WriteMu.Unlock()

			go listenMove(game, p, playerNum)
		} else {
			go listenLobby(p)
		}
	}
}

func addToQueue(p *Player) {
	mu.Lock()
	queue = append(queue, p)
	mu.Unlock()

	go waitForMatch(p)
}

func removeFromQueue(p *Player) bool {
	for i, qp := range queue {
		if qp == p {
			queue = append(queue[:i], queue[i+1:]...)
			return true
		}
	}
	return false
}

func waitForMatch(p *Player) {
	timer := time.After(10 * time.Second)

	for {
		select {
		case <-timer:
			mu.Lock()
			if removeFromQueue(p) {
				mu.Unlock()
				startGameWithBot(p)
			} else {
				mu.Unlock()
			}
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

func listenLobby(p *Player) {
	for {
		var msg Message
		err := p.Conn.ReadJSON(&msg)
		if err != nil {
			removePlayer(p)
			return
		}

		switch msg.Type {
		case "queue":
			addToQueue(p)
			return
		case "leaderboard":
			data, err := leaderboardData()
			if err != nil {
				log.Println("leaderboard failed: ", err)
				continue
			}
			p.WriteMu.Lock()
			p.Conn.WriteJSON(Message{
				Type: "leaderboard",
				Data: data,
			})
			p.WriteMu.Unlock()
		case "games":
			data, err := gamesData()
			if err != nil {
				log.Println("games failed: ", err)
				continue
			}
			p.WriteMu.Lock()
			p.Conn.WriteJSON(Message{
				Type: "games",
				Data: data,
			})
			p.WriteMu.Unlock()
		}
	}
}

func removePlayer(p *Player) {
	p.Conn.Close()
}
