package main

import (
	"log"
	"time"
)

func newGame(p1, p2 *Player) *Game {
	g := &Game{
		Turn:         1,
		Players:      [2]*Player{p1, p2},
		LastMoveTime: time.Now().UnixMilli(),
	}

	for r := range 6 {
		for c := range 7 {
			g.Board[r][c] = 0
		}
	}
	return g
}

func startGame(p1, p2 *Player) {
	game := newGame(p1, p2)

	activeGamesMu.Lock()
	activeGames[p1.ID] = game
	activeGames[p2.ID] = game
	activeGamesMu.Unlock()

	if !p1.IsBot {
		p1.Conn.WriteJSON(Message{
			Type: "start",
		})
	}

	if !p2.IsBot {
		p2.Conn.WriteJSON(Message{
			Type: "start",
		})
	}

	game.sendGameState(p1, p2)

	go monitorTimeout(game)

	go listenMove(game, p1, 1)

	if p2.IsBot {
		go botLoop(game)
	} else {
		go listenMove(game, p2, 2)
	}

}

func monitorTimeout(g *Game) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().UnixMilli()

		g.Mu.Lock()
		if g.Over {
			g.Mu.Unlock()
			return
		}

		if now-g.LastMoveTime >= 30_000 {
			if g.Turn == 1 {
				g.Winner = 2
			} else {
				g.Winner = 1
			}
			endGame(g)

			g.sendGameState(g.Players[0], g.Players[1])
			g.Mu.Unlock()
			return
		}
		g.Mu.Unlock()
	}
}

func startGameWithBot(p *Player) {
	log.Println("bot was here ---- ")
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

		switch msg.Type {
		case "move":
			if g.Over {
				continue
			}
			data := msg.Data.(map[string]any)
			column := int(data["column"].(float64))

			if g.Turn == playerNum && g.isValidMove(column) {
				g.makeMove(playerNum, column)
			}
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
		case "game_queue":
			leaveGame(g, p)
			addToQueue(p)
			return
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

	var row int
	for i := 5; i >= 0; i-- {
		if g.Board[i][column] == 0 {
			g.Board[i][column] = playerNum
			row = i
			break
		}
	}

	g.LastMoveTime = time.Now().UnixMilli()

	if g.checkWinner(row, column, playerNum) {
		g.Winner = playerNum
		endGame(g)
	} else if g.checkDraw() {
		g.Winner = 0
		endGame(g)
	} else {
		if g.Turn == 1 {
			g.Turn = 2
		} else {
			g.Turn = 1
		}
	}

	g.sendGameState(g.Players[0], g.Players[1])
}

func (g *Game) checkWinner(row, col, player int) bool {
	directions := [][2]int{
		{0, 1},
		{1, 0},
		{1, 1},
		{1, -1},
	}

	for _, d := range directions {
		count := 1
		count += g.countDir(row, col, d[0], d[1], player)
		count += g.countDir(row, col, -d[0], -d[1], player)

		if count >= 4 {
			return true
		}
	}

	return false
}

func (g *Game) countDir(r, c, dr, dc, player int) int {
	cnt := 0
	r += dr
	c += dc

	for r >= 0 && r < 6 && c >= 0 && c < 7 && g.Board[r][c] == player {
		cnt++
		r += dr
		c += dc
	}

	return cnt
}

func (g *Game) isValidMove(column int) bool {
	if g.Board[0][column] == 0 {
		return true
	} else {
		return false
	}
}

func (g *Game) checkDraw() bool {
	draw := true
	for i := range 7 {
		if g.Board[0][i] == 0 {
			draw = false
			break
		}
	}
	return draw
}

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

func leaveGame(g *Game, p *Player) {
	g.Mu.Lock()
	defer g.Mu.Unlock()
}

func endGame(g *Game) {
	if g.Over {
		return
	}
	g.Over = true

	activeGamesMu.Lock()
	delete(activeGames, g.Players[0].ID)
	delete(activeGames, g.Players[1].ID)
	activeGamesMu.Unlock()

	go storeResult(g)
}
