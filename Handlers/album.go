package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	bd "tp2/bd/sqlc"

	// Reemplazar esta ruta por la de GitHub
	//bd "Prog-WEB-2026/bd"
)

// AlbumHandler contiene la dependencia de la base de datos (sqlc)
type AlbumHandler struct {
	queries *bd.Queries
}

// NewAlbumHandler es el constructor para crear la instancia del handler
func NewAlbumHandler(q *bd.Queries) *AlbumHandler {
	return &AlbumHandler{queries: q}
}

// Estructura auxiliar
type Album struct {
	IDAlbum        int32     `json:"id_album"`
	AlbumName      string    `json:"album_name"`
	AlbumDur       int32     `json:"album_dur"`
	AlbumPh        int32     `json:"album_ph"`
	ReleaseDate    time.Time `json:"release_date"`
	ArtistIDArtist int32     `json:"artist_id_artist"`
	Likes          int32     `json:"likes"`
}

// AlbumsHandler maneja las peticiones dirigidas a /albums (GET y POST)
func (h *AlbumHandler) AlbumsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAlbums(w, r)
	case http.MethodPost:
		h.createAlbum(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AlbumHandler maneja las peticiones dirigidas a /albums/{id} (GET, PUT, DELETE)
func (h *AlbumHandler) AlbumHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAlbum(w, r, int32(id))
	case http.MethodPut:
		h.updateAlbum(w, r, int32(id))
	case http.MethodDelete:
		h.deleteAlbum(w, r, int32(id))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// FUNCIONES CRUD

// GET /albums - Listar todos los álbumes
func (h *AlbumHandler) getAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := h.queries.ListAlbums(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener albumes", http.StatusInternalServerError)
		return
	}

	// Si no hay álbumes cargados, devolvemos un array vacío [] en vez de null
	if albums == nil {
		albums = []bd.Album{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(albums)
}

// POST /albums - Crear un nuevo álbum
func (h *AlbumHandler) createAlbum(w http.ResponseWriter, r *http.Request) {
	var album Album

	err := json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar que los campos obligatorios no vengan vacíos
	if strings.TrimSpace(album.AlbumName) == "" ||
		strings.TrimSpace(album.ReleaseDate.String()) == "" ||
		album.AlbumDur < 0 ||
		album.AlbumPh < 0 ||
		album.ArtistIDArtist < 0 ||
		album.Likes < 0 {
		http.Error(w, "Faltan campos obligatorios o hay campos erroneos al crear el álbum", http.StatusBadRequest)
		return
	}

	params := bd.CreateAlbumParams{
		AlbumName:      album.AlbumName,
		AlbumDur:       album.AlbumDur,
		AlbumPh:        album.AlbumPh,
		ReleaseDate:    album.ReleaseDate,
		ArtistIDArtist: album.ArtistIDArtist,
		Likes:          album.Likes,
	}

	createdAlbum, err := h.queries.CreateAlbum(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear el álbum", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Estado 201 Created
	json.NewEncoder(w).Encode(createdAlbum)
}

// GET /albums/{id} - Obtener un álbum por ID
func (h *AlbumHandler) getAlbum(w http.ResponseWriter, r *http.Request, id int32) {
	album, err := h.queries.BuscarAlbumById(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Álbum no encontrado", http.StatusNotFound) // Estado 404
			return
		}
		http.Error(w, "Error al obtener el álbum", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(album)
}

// PUT /albums/{id} - Actualizar un álbum existente
func (h *AlbumHandler) updateAlbum(w http.ResponseWriter, r *http.Request, id int32) {
	// Verificar si el álbum existe antes de actualizar
	_, err := h.queries.BuscarAlbumById(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Álbum no encontrado para actualizar", http.StatusNotFound) // Estado 404
			return
		}
		http.Error(w, "Error al buscar el álbum", http.StatusInternalServerError)
		return
	}

	var album Album
	err = json.NewDecoder(r.Body).Decode(&album)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar datos 
	if strings.TrimSpace(album.AlbumName) == "" ||
		strings.TrimSpace(album.ReleaseDate.String()) == "" ||
		album.AlbumDur < 0 ||
		album.AlbumPh < 0 ||
		album.ArtistIDArtist < 0 ||
		album.Likes < 0 {
		http.Error(w, "Faltan campos obligatorios al actualizar el álbum", http.StatusBadRequest)
		return
	}

	params := bd.UpdateAlbumParams{
		IDAlbum:        id,
		AlbumName:      album.AlbumName,
		AlbumDur:       album.AlbumDur,
		AlbumPh:        album.AlbumPh,
		ReleaseDate:    album.ReleaseDate,
		ArtistIDArtist: album.ArtistIDArtist,
		Likes:          album.Likes,
	}

	err = h.queries.UpdateAlbum(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar el álbum", http.StatusInternalServerError)
		return
	}

	updatedAlbum, err := h.queries.BuscarAlbumById(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al obtener el álbum actualizado", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Estado 200 OK
	json.NewEncoder(w).Encode(updatedAlbum)
}

// DELETE /albums/{id} - Eliminar álbum
func (h *AlbumHandler) deleteAlbum(w http.ResponseWriter, r *http.Request, id int32) {
	// Verificar si el álbum existe antes de borrar
	_, err := h.queries.BuscarAlbumById(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Álbum no encontrado para eliminar", http.StatusNotFound) // Estado 404
			return
		}
		http.Error(w, "Error al buscar el álbum", http.StatusInternalServerError)
		return
	}

	// Eliminar de la base de datos
	err = h.queries.DeleteAlbum(r.Context(), id)
	if err != nil {
		http.Error(w, "Error al eliminar el álbum", http.StatusInternalServerError)
		return
	}

	// 204 No Content
	w.WriteHeader(http.StatusNoContent)
}