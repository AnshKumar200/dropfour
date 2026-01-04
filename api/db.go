package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func conncetDB() {
	log.Println("started the conn to db")
	dsn := os.Getenv("DATABASE_URL")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("failed to parse the db config: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	if err := DB.Ping(ctx); err != nil {
		log.Fatal("db not recheable: ", err)
	}

	log.Println("connceted to the db")
}

func storeResult(g *Game) {
	if DB == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Println("tx begin failed: ", err)
		return
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `insert into games (player1_id, player2_id, winner) values ($1, $2, $3)`, g.Players[0].ID, g.Players[1].ID, g.Winner)
	if err != nil {
		log.Println("inset game failed: ", err)
		return
	}

	for _, p := range g.Players {
		if p.IsBot { continue }

		_, err = tx.Exec(ctx, `insert into players (id, name, wins) values ($1, $2, 0) on conflict (id) do update set name = excluded.name`, p.ID, p.Name)
		if err != nil {
			log.Println("upsert player failed: ", err)
		}
	}

	if g.Winner > 0 {
		winner := g.Players[g.Winner-1]
		_, err = tx.Exec(ctx, `update players set wins = wins + 1 where id = $1`, winner.ID)
		if err != nil {
			log.Println("update wins failed: ", err)
			return
		}
	}
	if err := tx.Commit(ctx); err != nil {
		log.Println("commit failed: ", err)
	}
}

func leaderboardData() ([]LeaderboardEntry, error) {
	if DB == nil {
		return nil, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	rows, err := DB.Query(ctx, `select name, wins from players order by wins desc limit 10`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []LeaderboardEntry
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.Name, &e.Wins); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, e)
	}

	return leaderboard, nil
}
