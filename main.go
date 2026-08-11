package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	var err error
	connstr := "postgres://user:password@localhost:5433/videodb?sslmode=disable"

	db, err = pgxpool.New(context.Background(), connstr)
	if err != nil {
		log.Fatalf("error connecting to db: %v\n", err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("db not answering: %v\n", err)
	}
	fmt.Println("is working!")

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/stream", handleStream)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./templates/index.html")
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("id")
	if videoID == "" {
		http.Error(w, "write id", http.StatusBadRequest)
		return
	}

	var filepath string

	err := db.QueryRow(context.Background(), "SELECT filepath FROM videos WHERE id = $1", videoID).Scan(&filepath)
	if err != nil {
		http.Error(w, "id video not found", http.StatusNotFound)
		return
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filepath)
}
