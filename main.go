package main

import (
	"fmt"
	"net/http"
)

func main() {
	static := "./static"

	fileserver := http.FileServer(http.Dir(static))

	http.Handle("/", fileserver)

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", static)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
