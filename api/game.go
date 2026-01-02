package main

import "fmt"

func newGame(p1, p2 *Player) *Game {
	g := &Game{
		Turn:    1,
		Players: [2]*Player{p1, p2},
	}

	for r := 0; r < 6; r++ {
		for c := 0; c < 7; c++ {
			g.Board[r][c] = 0
		}
	}
	return g
}

func startGame(p1, p2 *Player) {
	fmt.Println("starting the game ------ ")
	game := newGame(p1, p2)

	if !p1.IsBot {
		p1.Conn.WriteJSON(Message{
			Type: "start",
			Data: map[string]any{
				"player": 1,
			},
		})
	}

	if !p2.IsBot {
		p2.Conn.WriteJSON(Message{
			Type: "start",
			Data: map[string]any{
				"player": 2,
			},
		})
	}

	game.sendGameState(p1, p2)

	go listenMove(game, p1, 1)

	if p2.IsBot {
		go botLoop(game)
	} else {
		go listenMove(game, p2, 2)
	}

}

func startGameWithBot(p *Player) {
	bot := &Player{
		ID:    "BOT",
		Name:  "BOT",
		IsBot: true,
	}

	startGame(p, bot)
}

func listenMove(g *Game, p *Player, playerNum int) {
	for {
		var msg Message
		err := p.Conn.ReadJSON(&msg)
		if err != nil {
			removePlayer(p)
			return
		}

		if msg.Type == "move" {
			data := msg.Data.(map[string]interface{})
			column := int(data["column"].(float64))

			if g.Turn == playerNum && g.isValidMove(column) {
				g.makeMove(playerNum, column)
			}
		}
	}
}

func botLoop(g *Game) {
	for !g.Over {
		if g.Turn == 2 {
			col := botMove(g)
			g.makeMove(2, col)
		}
	}
}

func (g *Game) makeMove(playerNum int, column int) {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	for i := 5; i >= 0; i-- {
		if g.Board[i][column] == 0 {
			g.Board[i][column] = playerNum
			break
		}
	}
	
	if g.Turn == 1 {
		g.Turn = 2
	} else {
		g.Turn = 1
	}

	if g.checkWinner() {
		g.Over = true
	}
	
	g.sendGameState(g.Players[0], g.Players[1])
}

func (g *Game) checkWinner() bool { return false }

func (g *Game) isValidMove(column int) bool { return true }

func (g *Game) sendGameState(p1, p2 *Player) {
	if !p1.IsBot {
		p1.WriteMu.Lock()
		p1.Conn.WriteJSON(Message{
			Type: "state",
			Data: g,
		})
		p1.WriteMu.Unlock()
	}

	if !p2.IsBot {
		p2.WriteMu.Lock()
		p2.Conn.WriteJSON(Message{
			Type: "state",
			Data: g,
		})
		p2.WriteMu.Unlock()
	}
}
