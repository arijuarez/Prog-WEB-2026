package handlers
 
import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	bd "tp2/bd/sqlc"
	// reemplazar después por el link al github
)
 
// Contiene las dependencias
type SongHandler struct {
	queries *bd.Queries
}
 
// Constructor para crear el handler con las dependencias necesarias
func NewSongHandler(q *bd.Queries) *SongHandler {
	return &SongHandler{queries: q}
}
 
type Song struct {
	IDSong       int32  `json:"id_song"`
	SongName     string `json:"song_name"`
	SongDur      int32  `json:"song_dur"`
	Likes        int32  `json:"likes"`
	AlbumIDAlbum int32  `json:"album_id_album"`
}

 
// maneja peticiones a /songs
func (h *SongHandler) songsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getSongs(w, r)
	case http.MethodPost:
		h.createSong(w, r)
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}
 
// maneja peticiones a /songs/{id}
func (h *SongHandler) songHandler(w http.ResponseWriter, r *http.Request) {
 
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
 
	if err != nil {
		http.Error(w, "Invalid Song ID", http.StatusBadRequest)
		return
	}
 
	songID := int32(id)
 
	switch r.Method {
	case http.MethodGet:
		h.getSong(w, r, songID)
	case http.MethodPut:
		h.updateSong(w, r, songID)
	case http.MethodDelete:
		h.deleteSong(w, r, songID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
 
// GET /songs - Listar todas las canciones
func (h *SongHandler) getSongs(w http.ResponseWriter, r *http.Request) {
 
	songs, err := h.queries.ListSongs(r.Context())
	if err != nil {
		http.Error(w, "Error obtaining Songs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}
 
// POST /songs - Crear nueva canción
func (h *SongHandler) createSong(w http.ResponseWriter, r *http.Request) {
	var song Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateNewSong(song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := bd.CreateSongParams{
		SongName:     song.SongName,
		SongDur:      song.SongDur,
		AlbumIDAlbum: song.AlbumIDAlbum,
		Likes:        song.Likes,
	}
 
	nuevaCancion, err := h.queries.CreateSong(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear la cancion", http.StatusInternalServerError)
		return
	}
 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevaCancion)
}
 
// validaciones comunes a crear y actualizar
func validateSongFields(song Song) error {
	if strings.TrimSpace(song.SongName) == "" {
		return errors.New("Falta nombre de la cancion")
	}
	if song.SongDur <= 0 {
		return errors.New("Duracion de la cancion invalida")
	}
	if song.AlbumIDAlbum <= 0 {
		return errors.New("Falta el album de la cancion")
	}
	return nil
}
 
func validateNewSong(song Song) error {
	if err := validateSongFields(song); err != nil {
		return err
	}
	if song.Likes != 0 {
		return errors.New("Valor invalido de likes")
	}
	return nil
}
 
func validateUpdateSong(song Song) error {
	if err := validateSongFields(song); err != nil {
		return err
	}
	if song.Likes < 0 {
		return errors.New("Valor invalido de likes")
	}
	return nil
}
 
// GET /songs/{id} - Obtener una canción específica
func (h *SongHandler) getSong(w http.ResponseWriter, r *http.Request, id int32) {
	cancion, err := h.queries.GetSongById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
 
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cancion)
}
 
// PUT /songs/{id} - Actualizar una canción específica
func (h *SongHandler) updateSong(w http.ResponseWriter, r *http.Request, id int32) {
	var song Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateUpdateSong(song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := bd.UpdateSongParams{
		IDSong:       id,
		SongName:     song.SongName,
		SongDur:      song.SongDur,
		Likes:        song.Likes,
	}
 
	cancionActualizada, err := h.queries.UpdateSong(r.Context(), params)
	if err == sql.ErrNoRows {
		http.Error(w, "Cancion no encontrada para actualizar", http.StatusNotFound) // 404
		return
	}
	if err != nil {
		http.Error(w, "Error al actualizar la cancion", http.StatusInternalServerError) // 500
		return
	}
 
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cancionActualizada)
}
 
// DELETE /songs/{id} - Eliminar una canción específica
func (h *SongHandler) deleteSong(w http.ResponseWriter, r *http.Request, id int32) {
	err := h.queries.DeleteSong(r.Context(), id)
 
	if err == sql.ErrNoRows {
		http.Error(w, "Cancion no encontrada para eliminar", http.StatusNotFound) // 404
		return
	}
	if err != nil {
		http.Error(w, "Error al eliminar cancion", http.StatusInternalServerError) // 500
		return
	}
 
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
 
