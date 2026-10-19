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

// ArtistSongHandler contiene las dependencias de la base de datos
type ArtistSongHandler struct {
	queries *bd.Queries
}

// NewArtistSongHandler es el constructor del handler
func NewArtistSongHandler(q *bd.Queries) *ArtistSongHandler {
	return &ArtistSongHandler{queries: q}
}

// Estructura auxiliar para decodificar JSON entrante
type ArtistSongRequest struct {
	ArtistIDArtist int32 `json:"artist_id_artist"`
	SongIDSong     int32 `json:"song_id_song"`
	Roll           *bool `json:"roll"`
}

// Manejador principal para /artist-songs
func (h *ArtistSongHandler) ArtistSongsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getArtistSongs(w, r)
	case http.MethodPost:
		h.createArtistSong(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Manejador principal para /artist-songs/{artist_id}/{song_id}
func (h *ArtistSongHandler) ArtistSongHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	artistID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	songID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getArtistSong(w, r, int32(artistID), int32(songID))
	case http.MethodPut:
		h.updateArtistSong(w, r, int32(artistID), int32(songID))
	case http.MethodDelete:
		h.deleteArtistSong(w, r, int32(artistID), int32(songID))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /artist-songs - Listar todas las relaciones
func (h *ArtistSongHandler) getArtistSongs(w http.ResponseWriter, r *http.Request) {
	artistSongs, err := h.queries.ListArtistSongs(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las relaciones artista-canción", http.StatusInternalServerError)
		return
	}

	if artistSongs == nil {
		artistSongs = []bd.ArtistSong{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(artistSongs)
}

// POST /artist-songs - Crear nueva relación
func (h *ArtistSongHandler) createArtistSong(w http.ResponseWriter, r *http.Request) {
	var req ArtistSongRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ArtistIDArtist <= 0 || req.SongIDSong <= 0 || req.Roll == nil {
		http.Error(w, "Faltan campos obligatorios al crear la relación", http.StatusBadRequest)
		return
	}

	params := bd.CreateArtistSongParams{
		ArtistIDArtist: req.ArtistIDArtist,
		SongIDSong:     req.SongIDSong,
		Roll:           *req.Roll,
	}

	created, err := h.queries.CreateArtistSong(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear la relación artista-canción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GET /artist-songs/{artist_id}/{song_id} - Obtener relación por clave compuesta
func (h *ArtistSongHandler) getArtistSong(w http.ResponseWriter, r *http.Request, artistID, songID int32) {
	params := bd.GetArtistSongParams{
		ArtistIDArtist: artistID,
		SongIDSong:     songID,
	}

	relation, err := h.queries.GetArtistSong(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Relación artista-canción no encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al obtener la relación artista-canción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(relation)
}

// PUT /artist-songs/{artist_id}/{song_id} - Actualizar el rol de la relación
func (h *ArtistSongHandler) updateArtistSong(w http.ResponseWriter, r *http.Request, artistID, songID int32) {
	getParams := bd.GetArtistSongParams{
		ArtistIDArtist: artistID,
		SongIDSong:     songID,
	}

	_, err := h.queries.GetArtistSong(r.Context(), getParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Relación artista-canción no encontrada para actualizar", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al buscar la relación artista-canción", http.StatusInternalServerError)
		return
	}

	var req struct {
		Roll *bool `json:"roll"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Roll == nil {
		http.Error(w, "El campo roll es obligatorio al actualizar", http.StatusBadRequest)
		return
	}

	updateParams := bd.UpdateArtistSongParams{
		ArtistIDArtist: artistID,
		SongIDSong:     songID,
		Roll:           *req.Roll,
	}

	updated, err := h.queries.UpdateArtistSong(r.Context(), updateParams)
	if err != nil {
		http.Error(w, "Error al actualizar la relación artista-canción", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// DELETE /artist-songs/{artist_id}/{song_id} - Eliminar relación
func (h *ArtistSongHandler) deleteArtistSong(w http.ResponseWriter, r *http.Request, artistID, songID int32) {
	getParams := bd.GetArtistSongParams{
		ArtistIDArtist: artistID,
		SongIDSong:     songID,
	}

	_, err := h.queries.GetArtistSong(r.Context(), getParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Relación artista-canción no encontrada para eliminar", http.StatusNotFound)
			return
		}
		http.Error(w, "Error al buscar la relación artista-canción", http.StatusInternalServerError)
		return
	}

	deleteParams := bd.DeleteArtistSongParams{
		ArtistIDArtist: artistID,
		SongIDSong:     songID,
	}

	err = h.queries.DeleteArtistSong(r.Context(), deleteParams)
	if err != nil {
		http.Error(w, "Error al eliminar la relación artista-canción", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}