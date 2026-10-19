package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	bd "tp2/bd/sqlc"

)

// Contiene las dependencias
type UserHandler struct {
	queries *bd.Queries
}

// Constructor para crear el handler con las dependencias necesarias
func NewUserHandler(q *bd.Queries) *UserHandler {
	return &UserHandler{queries: q}
}

type User struct {
	IDUser   int32  `json:"id_user"`
	UserName string `json:"user_name"`
	Amigos   int32  `json:"amigos"`
	Photo    int32  `json:"photo"`
}

// maneja peticiones a /users
func (h *UserHandler) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// maneja peticiones a /users/{id}
func (h *UserHandler) userHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	userID := int32(id)

	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r, userID)
	case http.MethodPut:
		h.updateUser(w, r, userID)
	case http.MethodDelete:
		h.deleteUser(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /users - Listar todos los usuarios
func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "Error obtaining Users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// POST /users - Crear nuevo usuario
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateNewUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := bd.CreateUserParams{
		UserName: user.UserName,
		Amigos:   user.Amigos,
		Photo:    user.Photo,
	}

	usario, err := h.queries.CreateUser(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usario)
}

func validateNewUser(usuario User) error {
	if strings.TrimSpace(usuario.UserName) == "" {
		return errors.New("Falta nombre de usuario")
	}
	if usuario.Amigos != 0 {
		return errors.New("Valor invalido de amigos")
	}
	return nil
}

// GET /users/{id} - Obtener un usuario específico
func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request, id int32) {
	usuario, err := h.queries.GetUserById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// PUT /users/{id} - Actualizar un usuario específico
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id int32) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateNewUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := bd.UpdateUserParams{
		IDUser:   id,
		UserName: user.UserName,
		Amigos:   user.Amigos,
		Photo:    user.Photo,
	}

	usario, err := h.queries.UpdateUser(r.Context(), params)
	if err == sql.ErrNoRows {
		http.Error(w, "Usuario no encontrado para actualizar", http.StatusNotFound) // 404
		return
	}
	if err != nil {
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usario)
}
// DELETE /users/{id} - Eliminar un usuario específico
func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.queries.DeleteUser(r.Context(), id)

	if err == sql.ErrNoRows {
		http.Error(w, "Usuario no encontrado para eliminar", http.StatusNotFound) // 404
		return
	}
	if err != nil {
		http.Error(w, "Error al eliminar usuario", http.StatusInternalServerError) // 500
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
