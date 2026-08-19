package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed templates/* static/*
var embeddedFiles embed.FS

var db *pgxpool.Pool

func main() {
	var err error

	connstr := os.Getenv("DATABASE_URL")
	if connstr == "" {
		connstr = "postgres://user:password@localhost:5433/videodb?sslmode=disable"
	}

	log.Println("Connecting to database...")
	for i := 0; i < 15; i++ {
		db, err = pgxpool.New(context.Background(), connstr)
		if err == nil {
			err = db.Ping(context.Background())
			if err == nil {
				break
			}
		}
		log.Printf("Database not ready (attempt %d/15), waiting 2 seconds...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Critical error: failed to connect to db: %v\n", err)
	}
	defer db.Close()
	fmt.Println("Database connection established!")

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/stream", handleStream)

	staticFS := http.FS(embeddedFiles)
	http.Handle("/static/", http.FileServer(staticFS))

	fmt.Println("Server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data, err := embeddedFiles.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error: template missing", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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
