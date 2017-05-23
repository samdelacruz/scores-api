package main

import (
	"log"
	"net/http"
	"os"
)

type score struct {
	ID        int
	Game      string
	Player    string
	Score     int
	CreatedAt string
}

var createScores = `
	CREATE TABLE IF NOT EXISTS "scores" (
		id         serial      PRIMARY KEY,
		game       varchar(32) NOT NULL,
		player     varchar(32) NOT NULL,
		score      integer     NOT NULL DEFAULT 0,
		created_at datetime    NOT NULL DEFAULT NOW(),
		UNIQUE (game, player)
	);
	CREATE INDEX player_i ON scores (player);
	CREATE INDEX score_i ON scores (score);
	CREATE INDEX created_at_i ON scores (created_at);
`

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT environment variable must be set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	log.Println("Servicing requests at http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
