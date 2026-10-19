package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"Prog-WEB-2026/Handlers"
	_ "github.com/lib/pq" //importar drivers
)

func abrirDB() (*sql.DB, error) {

	//Abre base de datos
	db, err := sql.Open("mysql", "usuario:clave@tcp(localhost:3306)/mi_app")

	// Alternativa: db, err := sql.Open("postgres", "postgres://postgres:1234@localhost:5431/postgres?sslmode=disable")

	//no confundir nombre del docker con nombre de la db
	db, err = sql.Open("postgres", "host=localhost port=5431 user=postgres password=1234 dbname=postgres sslmode=disable")

	//Controla errores
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		// Verifica que la BD responde
		return nil, err
	}

	db.SetMaxOpenConns(25) // Configura el tamaño del pool
	return db, nil
}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	http.HandleFunc("/artists", artistsHandler)
	http.HandleFunc("/artists/", artistHandler)
	http.HandleFunc("/songs", songsHandler)
	http.HandleFunc("/songs/", songHandler)
	http.HandleFunc("/albums", albumsHandler)
	http.HandleFunc("/albums/", albumHandler)
	http.HandleFunc("/reviews", reviewsHandler)
	http.HandleFunc("/reviews/", reviewHandler)

	db, err := abrirDB()
	if err != nil {
		log.Fatalf("Error al conectar con base de datos %v\n", err)
	}
	repo := userRepository{db: db}

	
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
	defer db.Close()

	fmt.Println("Se conectó exitosamente")

	//Instanciamos repositorio
}
