package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	bd "tp2/bd/sqlc"

	// reemplazar después por el link al github
	//"Prog-WEB-2026/bd"
)

// Contiene las dependencias
type ArtistHandler struct {
	queries *bd.Queries
}

// Constructor para crear el handler con las dependencias necesarias
func NewArtistHandler(q *bd.Queries) *ArtistHandler {
	return &ArtistHandler{queries: q}
}

// Estructura auxiliar para implementar algunos métodos
type Artista struct {
	IDArtist   int32  `json:"id_artist"`
	ArtistName string `json:"artist_name"`
	Photo      int32  `json:"photo"`
	ArtistBio  string `json:"artist_bio"`
}

// Manejar peticiones a /artists
func (h *ArtistHandler) artistsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getArtists(w, r)
	case http.MethodPost:
		h.createArtist(w, r)
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}

// Maneja peticiones a /artists/{id}
func (h *ArtistHandler) artistHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer id del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])

	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artistID := int32(id)

	switch r.Method {
	case http.MethodGet:
		h.getArtist(w, r, artistID)
	case http.MethodPut:
		h.updateArtist(w, r, artistID)
	case http.MethodDelete:
		h.deleteArtist(w, r, artistID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /artists - listar todos los artistas
func (h *ArtistHandler) getArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := h.queries.ListArtists(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener artista", http.StatusInternalServerError) // Error 500 para problemas inesperados
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

// crear artista maneja POST /artists
func (h *ArtistHandler) createArtist(w http.ResponseWriter, r *http.Request) {
	var artista Artista

	err := json.NewDecoder(r.Body).Decode(&artista)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validar que se hayan ingresado datos
	if strings.TrimSpace(artista.ArtistName) == "" ||
		artista.Photo <= 0 ||
		strings.TrimSpace(artista.ArtistBio) == "" {
		http.Error(w, "Faltan campos obligatorios al crear artista", http.StatusBadRequest)
		return
	}

	params := bd.CreateArtistParams{
		ArtistName: artista.ArtistName,
		Photo:      artista.Photo,
		ArtistBio:  artista.ArtistBio,
	}
	artist, err := h.queries.CreateArtist(r.Context(), params)

	if err != nil {
		http.Error(w, "Error al crear el artista", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artist)
}

// GET /artists/{id}
func (h *ArtistHandler) getArtist(w http.ResponseWriter, r *http.Request, id int32) {
	artista, err := h.queries.GetArtistById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artista)
}

// PUT /artists/{id} - Actualizar artista
func (h *ArtistHandler) updateArtist(w http.ResponseWriter, r *http.Request, id int32) {
	var artista Artista

	err := json.NewDecoder(r.Body).Decode(&artista)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(artista.ArtistName) == "" ||
		artista.Photo <= 0 ||
		strings.TrimSpace(artista.ArtistBio) == "" {
		http.Error(w, "Faltan campos obligatorios al actualizar artista", http.StatusBadRequest)
		return
	}

	params := bd.UpdateArtistParams{
		ArtistName: artista.ArtistName,
		Photo:      artista.Photo,
		ArtistBio:  artista.ArtistBio,
	}
	artist, err := h.queries.UpdateArtist(r.Context(), params)

	if err == sql.ErrNoRows {
		http.Error(w, "Artista no encontrado para actualizar", http.StatusNotFound) // 404
		return
	}

	if err != nil {
		http.Error(w, "Error al actualizar el artista", http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(artist)
}

func (h *ArtistHandler) deleteArtist(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.queries.DeleteArtist(r.Context(), id)

	if err == sql.ErrNoRows {
		http.Error(w, "Artista no encontrado para eliminar", http.StatusNotFound) // 404
		return
	}
	if err != nil {
		http.Error(w, "Error al eliminar artista", http.StatusInternalServerError) // 500
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
