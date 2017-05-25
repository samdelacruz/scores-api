package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type score struct {
	ID      *int64     `json:"id,omitempty"`
	Game    string     `json:"game"`
	Player  string     `json:"player"`
	Score   int64      `json:"score"`
	Created *time.Time `json:"created,omitempty"`
}

var (
	host     = os.Getenv("DB_HOST")
	port, _  = strconv.Atoi(os.Getenv("DB_PORT"))
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	dbname   = os.Getenv("DB_NAME")
	db       *sql.DB
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/scores", a.HandleNewScore).Methods("POST")
	a.Router.HandleFunc("/scores", a.HandleListScores).Methods("GET")
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT environment variable must be set")
	}

	app := App{DB: db, Router: mux.NewRouter()}
	app.initRoutes()

	http.Handle("/", app.Router)

	log.Println("Servicing requests at http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func (s *score) create(db *sql.DB) error {
	if s.Created != nil || s.ID != nil {
		return error(fmt.Errorf("Error trying to save an existing score"))
	}
	var (
		id      int64
		created time.Time
	)
	insertSQL := "INSERT INTO scores (game, player, score) values ($1, $2, $3) returning id, created"
	err := db.QueryRow(insertSQL, s.Game, s.Player, s.Score).Scan(&id, &created)

	if err != nil {
		return err
	}

	s.ID = &id
	s.Created = &created

	return nil
}

func getScores(db *sql.DB, lim int, off int) ([]score, error) {
	sql := "SELECT id, game, player, score, created FROM scores LIMIT $1 OFFSET $2"
	rows, err := db.Query(sql, lim, off)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	scores := []score{}

	for rows.Next() {
		var s score
		if err := rows.Scan(&s.ID, &s.Game, &s.Player, &s.Score, &s.Created); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}

	return scores, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
