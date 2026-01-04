package main

func botMove(g *Game) int {
	bestCol := -1
	maxScore := -100000 
	alpha := -100000
	beta := 100000

	for c := range 7 {
		if g.isValidMove(c) {
			next := g.clone()
			row := next.simulateMove(c, 2)

			if next.checkWinner(row, c, 2) {
				return c
			}
			score := minimax(next, 5, alpha, beta, false)
			
			if score > maxScore {
				maxScore = score
				bestCol = c
			}
			
			if score > alpha {
				alpha = score
			}
		}
	}

	if bestCol == -1 {
		for c := range 7 {
			if g.isValidMove(c) {
				return c
			}
		}
	}
	return bestCol
}

func minimax(g *Game, depth, alpha, beta int, maximizing bool) int {
	if depth == 0 {
		return 0
	}

	if maximizing {
		maxEval := -100000
		for c := 0; c < 7; c++ {
			if g.isValidMove(c) {
				next := g.clone()
				row := next.simulateMove(c, 2)

				if next.checkWinner(row, c, 2) {
					return 1000 + depth
				}

				eval := minimax(next, depth-1, alpha, beta, false)
				if eval > maxEval {
					maxEval = eval
				}
				if eval > alpha {
					alpha = eval
				}
				if beta <= alpha {
					break
				}
			}
		}
		return maxEval
	} else {
		minEval := 100000
		for c := range 7 {
			if g.isValidMove(c) {
				next := g.clone()
				row := next.simulateMove(c, 1)

				if next.checkWinner(row, c, 1) {
					return -1000 - depth
				}

				eval := minimax(next, depth-1, alpha, beta, true)
				if eval < minEval {
					minEval = eval
				}
				if eval < beta {
					beta = eval
				}
				if beta <= alpha {
					break
				}
			}
		}
		return minEval
	}
}
