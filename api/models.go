package main

import (
	"sync"

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

	Mu sync.Mutex
}

type Move struct {
	Column int `json:"column"`
}

type Message struct {
	Type string `json:"type"`
	Data any `json:"data"`
}
