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

type GamesEntry struct {
	Player1 string
	Player2 string
	Winner int
}

type Move struct {
	Column int `json:"column"`
}

type Message struct {
	Type string `json:"type"`
	Data any `json:"data"`
}

func (g *Game) clone() *Game {
	newG := &Game{
		Turn: g.Turn,
	}
	for r := range 6 {
		for c := range 7 {
			newG.Board[r][c] = g.Board[r][c]
		}
	}
	return newG
}

func (g *Game) simulateMove(col int, player int) int {
	for r := 5; r >= 0; r-- {
		if g.Board[r][col] == 0 {
			g.Board[r][col] = player
			return r
		}
	}
	return -1
}
