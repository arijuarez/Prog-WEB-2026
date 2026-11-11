package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	Handlers "github.com/arijuarez/Prog-WEB-2026/Handlers"
	bd "github.com/arijuarez/Prog-WEB-2026/bd/sqlc"

	_ "github.com/lib/pq" // importar drivers
)

func abrirDB() (*sql.DB, error) {
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

	// Archivos estáticos (styles.css, imágenes, js) desde ./static
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			Handlers.HandleGET(w, r, queries)
		case http.MethodPost:
			Handlers.HandlePOST(w, r, queries)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}