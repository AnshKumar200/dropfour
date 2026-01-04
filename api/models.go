package main

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID string
	Name string
	Conn *websocket.Conn
	IsBot bool

	WriteMu sync.Mutex
}

type Game struct {
	Board [6][7]int
	Turn int
	Players [2]*Player
	Winner int
	Over bool
	LastMoveTime int64

	Mu sync.Mutex
}

type GameResult struct {
	ID string
	Player1 string
	Player2 string
	Winner int
	EndedAt time.Time
}

type LeaderboardEntry struct {
	Name string
	Wins int
}

type Move struct {
	Column int `json:"column"`
}

type Message struct {
	Type string `json:"type"`
	Data any `json:"data"`
}
