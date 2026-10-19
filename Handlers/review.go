package handlers
 
import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	bd "github.com/arijuarez/Prog-WEB-2026/bd/sqlc"
)
 
// Contiene las dependencias
type ReviewHandler struct {
	queries *bd.Queries
}
 
// Constructor para crear el handler con las dependencias necesarias
func NewReviewHandler(q *bd.Queries) *ReviewHandler {
	return &ReviewHandler{queries: q}
}
 
// SongID y AlbumID son punteros porque en la base son FK nullable:
// si el campo no viene en el JSON (o viene null) queda en nil.
type Review struct {
	IDReview    int32   `json:"id_review"`
	ReviewText  string  `json:"review_text"`
	ReviewRate  float64 `json:"review_rate"`
	ReviewDate  string  `json:"review_date"` // formato YYYY-MM-DD
	UserID      int32   `json:"user_id_user"`
	SongID      *int32  `json:"song_id_song"`
	AlbumID     *int32  `json:"album_id_album"`
	Likes       int32   `json:"likes"`
}
 
// maneja peticiones a /reviews
func (h *ReviewHandler) ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getReviews(w, r)
	case http.MethodPost:
		h.createReview(w, r)
	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}
 
// maneja peticiones a /reviews/{id}
func (h *ReviewHandler) ReviewHandler(w http.ResponseWriter, r *http.Request) {
 
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
 
	if err != nil {
		http.Error(w, "Invalid Review ID", http.StatusBadRequest)
		return
	}
 
	reviewID := int32(id)
 
	switch r.Method {
	case http.MethodGet:
		h.getReview(w, r, reviewID)
	case http.MethodPut:
		h.updateReview(w, r, reviewID)
	case http.MethodDelete:
		h.deleteReview(w, r, reviewID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
 
// convierte un *int32 del JSON en el tipo nullable de la base
func toNullInt32(v *int32) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *v, Valid: true}
}
 
// GET /reviews - Listar todas las reseñas
func (h *ReviewHandler) getReviews(w http.ResponseWriter, r *http.Request) {
 
	reviews, err := h.queries.ListReviews(r.Context())
	if err != nil {
		http.Error(w, "Error obtaining Reviews", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}
 
// POST /reviews - Crear nueva reseña
func (h *ReviewHandler) createReview(w http.ResponseWriter, r *http.Request) {
	var review Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateNewReview(review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reviewDate, _ := time.Parse("2006-01-02", review.ReviewDate) // ya validada
 
	params := bd.CreateReviewParams{
		ReviewText:   review.ReviewText,
		ReviewRate:   review.ReviewRate,
		ReviewDate:   reviewDate,
		UserIDUser:   review.UserID,
		SongIDSong:   toNullInt32(review.SongID),
		AlbumIDAlbum: toNullInt32(review.AlbumID),
		Likes:        review.Likes,
	}
 
	nuevaReview, err := h.queries.CreateReview(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear la review", http.StatusInternalServerError)
		return
	}
 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevaReview)
}
 
// validaciones comunes a crear y actualizar
func validateReviewFields(review Review) error {
	if strings.TrimSpace(review.ReviewText) == "" {
		return errors.New("Falta texto de la review")
	}
	if review.ReviewRate < 0 || review.ReviewRate > 5 {
		return errors.New("Puntaje invalido (debe estar entre 0 y 5)")
	}
	if _, err := time.Parse("2006-01-02", review.ReviewDate); err != nil {
		return errors.New("Fecha de review invalida (formato YYYY-MM-DD)")
	}
	if review.UserID <= 0 {
		return errors.New("Falta el usuario de la review")
	}
	// la reseña es de una cancion o de un album, pero no de ambos ni de ninguno
	if (review.SongID == nil) == (review.AlbumID == nil) {
		return errors.New("La review debe tener song_id_song o album_id_album (solo uno)")
	}
	return nil
}
 
func validateNewReview(review Review) error {
	if err := validateReviewFields(review); err != nil {
		return err
	}
	if review.Likes != 0 {
		return errors.New("Valor invalido de likes")
	}
	return nil
}
 
func validateUpdateReview(review Review) error {
	if err := validateReviewFields(review); err != nil {
		return err
	}
	if review.Likes < 0 {
		return errors.New("Valor invalido de likes")
	}
	return nil
}
 
// GET /reviews/{id} - Obtener una reseña específica
func (h *ReviewHandler) getReview(w http.ResponseWriter, r *http.Request, id int32) {
	review, err := h.queries.GetReview(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
 
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}
 
// PUT /reviews/{id} - Actualizar una reseña específica
func (h *ReviewHandler) updateReview(w http.ResponseWriter, r *http.Request, id int32) {
	var review Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateUpdateReview(review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reviewDate, _ := time.Parse("2006-01-02", review.ReviewDate) // ya validada
 
	params := bd.UpdateReviewParams{
		IDReview:     id,
		ReviewText:   review.ReviewText,
		ReviewRate:   review.ReviewRate,
		ReviewDate:   reviewDate,
		UserIDUser:   review.UserID,
		SongIDSong:   toNullInt32(review.SongID),
		AlbumIDAlbum: toNullInt32(review.AlbumID),
		Likes:        review.Likes,
	}
 
	reviewActualizada, err := h.queries.UpdateReview(r.Context(), params)
	if err == sql.ErrNoRows {
		http.Error(w, "Review no encontrada para actualizar", http.StatusNotFound) // 404
		return
	}
	if err != nil {
		http.Error(w, "Error al actualizar la review", http.StatusInternalServerError) // 500
		return
	}
 
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviewActualizada)
}
 
// DELETE /reviews/{id} - Eliminar una reseña específica
func (h *ReviewHandler) deleteReview(w http.ResponseWriter, r *http.Request, id int32) {
	// Verificar si la review existe antes de borrar
	_, err := h.queries.GetReview(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Review no encontrado para eliminar", http.StatusNotFound) // Estado 404
			return
		}
		http.Error(w, "Error al buscar review", http.StatusInternalServerError)
		return
	}	
	
	err = h.queries.DeleteReview(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al eliminar review", http.StatusInternalServerError) // 500
		return
	}
 
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
