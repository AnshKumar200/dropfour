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

		switch msg.Type {
		case "move":
			if g.Over { continue }
			data := msg.Data.(map[string]interface{})
			column := int(data["column"].(float64))

			if g.Turn == playerNum && g.isValidMove(column) {
				g.makeMove(playerNum, column)
			}
		case "queue":
			fmt.Println("putting back in queue: ", p.Name)
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

	if g.checkWinner(row, column, playerNum) {
		g.Over = true
		g.Winner = playerNum
	}
	
	if !g.Over {
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
			fmt.Println("winner is :", player)
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
		cnt++;
		r += dr
		c += dc
	}

	return cnt
}

func (g *Game) isValidMove(column int) bool {
	if(g.Board[0][column] == 0) {
		return true
	} else { return false }
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
