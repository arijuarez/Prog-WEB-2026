package main

import (
	"os"
	"database/sql"
	"log"
	"net/http"
	"fmt"

	Handlers "github.com/arijuarez/Prog-WEB-2026/Handlers"
	bd "github.com/arijuarez/Prog-WEB-2026/bd/sqlc"

	_ "github.com/lib/pq" //importar drivers
)

func abrirDB() (*sql.DB, error) {
	// Lee las variables de entorno o usa los valores por defecto (para ejecución fuera de Docker)
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5431"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "1234"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	return db, nil
}

func main() {
	db, err := abrirDB()
	if err != nil {
		log.Fatalf("Error al conectar con base de datos %v\n", err)
	}
	defer db.Close()
	log.Println("Conexión a PostgreSQL establecida correctamente")

	queries := bd.New(db)

	userH := Handlers.NewUserHandler(queries)
	artistH := Handlers.NewArtistHandler(queries)
	songH := Handlers.NewSongHandler(queries)
	albumH := Handlers.NewAlbumHandler(queries)
	reviewH := Handlers.NewReviewHandler(queries)
	artistSongH := Handlers.NewArtistSongHandler(queries)

	http.HandleFunc("/users", userH.UsersHandler)
	http.HandleFunc("/users/", userH.UserHandler)

	http.HandleFunc("/artists", artistH.ArtistsHandler)
	http.HandleFunc("/artists/", artistH.ArtistHandler)

	http.HandleFunc("/songs", songH.SongsHandler)
	http.HandleFunc("/songs/", songH.SongHandler)

	http.HandleFunc("/albums", albumH.AlbumsHandler)
	http.HandleFunc("/albums/", albumH.AlbumHandler)

	http.HandleFunc("/reviews", reviewH.ReviewsHandler)
	http.HandleFunc("/reviews/", reviewH.ReviewHandler)

	http.HandleFunc("/artist-songs", artistSongH.ArtistSongsHandler)
	http.HandleFunc("/artist-songs/", artistSongH.ArtistSongHandler)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
