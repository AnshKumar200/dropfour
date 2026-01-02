package main

import "github.com/gorilla/websocket"

type Player struct {
	ID string
	Name string
	Conn *websocket.Conn
	IsBot bool
}

type Game struct {
	Board [6][7]int
	Turn int
	Players [2]*Player
	Over bool
}

type Move struct {
	Column int `json:"column"`
}

type Message struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}
