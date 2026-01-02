package main

func botMove(g *Game) int {
	for c := 0; c < 7; c++ {
		if g.isValidMove(c) { return c }
	}

	return 0
}

func evaluateMove(g *Game, column int) int { return 0 }
