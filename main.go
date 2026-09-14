/*
package main

import (

	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" //importar drivers

)

	type User struct {
		id         int
		name       string
		email      string
		created_at time.Time
	}

	type userRepository struct {
		db *sql.DB
	}

	func abrirDB() (*sql.DB, error) {
		/*
			//Abre base de datos
			db, err := sql.Open("mysql", "usuario:clave@tcp(localhost:3306)/mi_app")

		// Alternativa: db, err := sql.Open("postgres", "postgres://postgres:1234@localhost:5431/postgres?sslmode=disable")

		//no confundir nombre del docker con nombre de la db
		db, err := sql.Open("postgres", "host=localhost port=5431 user=postgres password=1234 dbname=postgres sslmode=disable")

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

	func (r *userRepository) createUser(user *User) {
		_, err := r.db.Exec("INSERT INTO users(name, email) VALUES($1,$2)", user.name, user.email)
		if err != nil {
			log.Fatalf("Error al crear usuario %v\n", err)
		}
	}

	func main() {
		db, err := abrirDB()
		if err != nil {
			log.Fatalf("Error al conectar con base de datos %v\n", err)
		}
		defer db.Close()

		fmt.Println("Se conectó exitosamente")

		usuario := User{
			name:  "Nacho",
			email: "nacho@gmail.com",
		}

		//Instanciamos repositorio
		repo := userRepository{db: db}
		repo.createUser(&usuario)

}
*/
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Main.go andando")
}
