package Handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	bd "github.com/arijuarez/Prog-WEB-2026/bd/sqlc"
	views "github.com/arijuarez/Prog-WEB-2026/views"
)

// --------- helpers ---------

// intRequerido convierte un campo obligatorio a int32.
func intRequerido(s string) (int32, error) {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	return int32(n), err
}

// intOpcional convierte un campo opcional a int32 (vacío = 0).
func intOpcional(s string) (int32, error) {
	if strings.TrimSpace(s) == "" {
		return 0, nil
	}
	return intRequerido(s)
}

// nullInt32 convierte un campo opcional a sql.NullInt32 (vacío = NULL).
func nullInt32(s string) (sql.NullInt32, error) {
	if strings.TrimSpace(s) == "" {
		return sql.NullInt32{}, nil
	}
	n, err := intRequerido(s)
	return sql.NullInt32{Int32: n, Valid: err == nil}, err
}

// errorServidor registra el error real en el log y le responde al cliente un mensaje genérico.
func errorServidor(w http.ResponseWriter, msg string, err error, code int) {
	log.Printf("%s: %v", msg, err)
	http.Error(w, msg, code)
}

func badRequest(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

// --------- GET ---------

// HandleGET elige, según la ruta, qué página renderizar.
func HandleGET(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	switch r.URL.Path {
	case "/":
		views.IndexMain().Render(r.Context(), w)
	case "/users":
		getUsers(w, r, queries)
	case "/songs":
		getSongs(w, r, queries)
	case "/albums":
		getAlbums(w, r, queries)
	case "/reviews":
		getReviews(w, r, queries)
	case "/artists":
		getArtists(w, r, queries)
	default:
		http.NotFound(w, r)
	}
}

// getUsers atiende GET /users: página completa con la lista de usuarios.
func getUsers(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	list, err := queries.ListUsers(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener usuarios", err, http.StatusInternalServerError)
		return
	}
	views.IndexUsers(list).Render(ctx, w)
}

// getSongs atiende GET /songs: página completa con la lista de canciones.
func getSongs(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	list, err := queries.ListSongs(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener canciones", err, http.StatusInternalServerError)
		return
	}
	views.IndexSongs(list).Render(ctx, w)
}

// getAlbums atiende GET /albums: página completa con la lista de álbumes.
func getAlbums(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	list, err := queries.ListAlbums(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener álbumes", err, http.StatusInternalServerError)
		return
	}
	views.IndexAlbums(list).Render(ctx, w)
}

// getReviews atiende GET /reviews: página completa con la lista de reseñas.
func getReviews(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	list, err := queries.ListReviews(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener reseñas", err, http.StatusInternalServerError)
		return
	}
	views.IndexReviews(list).Render(ctx, w)
}

// getArtists atiende GET /artists: página completa con la lista de artistas.
func getArtists(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	list, err := queries.ListArtists(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener artistas", err, http.StatusInternalServerError)
		return
	}
	views.IndexArtists(list).Render(ctx, w)
}

// --------- POST ---------

// HandlePOST parsea el formulario y elige, según la ruta, qué entidad crear.
// Cada handler responde solo con la lista actualizada y el formulario vacío (OOB).
func HandlePOST(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	if err := r.ParseForm(); err != nil {
		badRequest(w, "Error al parsear el formulario")
		return
	}

	switch r.URL.Path {
	case "/users":
		postUsers(w, r, queries)
	case "/songs":
		postSongs(w, r, queries)
	case "/albums":
		postAlbums(w, r, queries)
	case "/reviews":
		postReviews(w, r, queries)
	case "/artists":
		postArtists(w, r, queries)
	default:
		http.NotFound(w, r)
	}
}

// postUsers atiende POST /users: valida, crea y devuelve la lista de usuarios.
func postUsers(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	userName := strings.TrimSpace(r.FormValue("user_name"))
	if userName == "" {
		badRequest(w, "El nombre de usuario es obligatorio")
		return
	}
	amigos, err := intOpcional(r.FormValue("amigos"))
	if err != nil || amigos < 0 {
		badRequest(w, "Cantidad de amigos inválida")
		return
	}
	photo, err := intOpcional(r.FormValue("user_ph"))
	if err != nil || photo < 0 {
		badRequest(w, "Valor de foto inválido")
		return
	}
	_, err = queries.CreateUser(ctx, bd.CreateUserParams{UserName: userName, Amigos: amigos, Photo: photo})
	if err != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}
	list, err := queries.ListUsers(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener usuarios", err, http.StatusInternalServerError)
		return
	}
	views.UserList(list).Render(ctx, w)
	views.FormUser(true).Render(ctx, w) // form vacío (OOB) → limpia el formulario
}

// postSongs atiende POST /songs: valida, crea y devuelve la lista de canciones.
func postSongs(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	songName := strings.TrimSpace(r.FormValue("song_name"))
	if songName == "" {
		badRequest(w, "El título es obligatorio")
		return
	}
	albumID, err := intRequerido(r.FormValue("album_id_album"))
	if err != nil || albumID < 1 {
		badRequest(w, "ID de álbum inválido")
		return
	}
	dur, err := intOpcional(r.FormValue("song_dur"))
	if err != nil || dur < 0 {
		badRequest(w, "Duración inválida")
		return
	}
	likes, err := intOpcional(r.FormValue("likes"))
	if err != nil || likes < 0 {
		badRequest(w, "Likes inválidos")
		return
	}
	_, err = queries.CreateSong(ctx, bd.CreateSongParams{
		SongName: songName, SongDur: dur, AlbumIDAlbum: albumID, Likes: likes,
	})
	if err != nil {
		http.Error(w, "Error al crear la canción", http.StatusInternalServerError)
		return
	}
	list, err := queries.ListSongs(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener canciones", err, http.StatusInternalServerError)
		return
	}
	views.SongList(list).Render(ctx, w)
	views.FormSong(true).Render(ctx, w) // form vacío (OOB) → limpia el formulario
}

// postAlbums atiende POST /albums: valida, crea y devuelve la lista de álbumes.
func postAlbums(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	albumName := strings.TrimSpace(r.FormValue("album_name"))
	if albumName == "" {
		badRequest(w, "El título es obligatorio")
		return
	}
	artistID, err := intRequerido(r.FormValue("artist_id_artist"))
	if err != nil || artistID < 1 {
		badRequest(w, "ID de artista inválido")
		return
	}
	releaseDate, err := time.Parse("2006-01-02", r.FormValue("release_date"))
	if err != nil {
		badRequest(w, "Fecha de lanzamiento inválida")
		return
	}
	albumPh, err := intOpcional(r.FormValue("album_ph"))
	if err != nil || albumPh < 0 {
		badRequest(w, "Valor de portada inválido")
		return
	}
	albumDur, err := intOpcional(r.FormValue("album_dur"))
	if err != nil || albumDur < 0 {
		badRequest(w, "Duración inválida")
		return
	}
	likes, err := intOpcional(r.FormValue("likes"))
	if err != nil || likes < 0 {
		badRequest(w, "Likes inválidos")
		return
	}
	_, err = queries.CreateAlbum(ctx, bd.CreateAlbumParams{
		AlbumName: albumName, AlbumDur: albumDur, AlbumPh: albumPh,
		ReleaseDate: releaseDate, ArtistIDArtist: artistID, Likes: likes,
	})
	if err != nil {
		http.Error(w, "Error al crear el álbum", http.StatusInternalServerError)
		return
	}
	list, err := queries.ListAlbums(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener álbumes", err, http.StatusInternalServerError)
		return
	}
	views.AlbumList(list).Render(ctx, w)
	views.FormAlbum(true).Render(ctx, w) // form vacío (OOB) → limpia el formulario
}

// postReviews atiende POST /reviews: valida, crea y devuelve la lista de reseñas.
func postReviews(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	text := strings.TrimSpace(r.FormValue("review_text"))
	if text == "" {
		badRequest(w, "El texto de la reseña es obligatorio")
		return
	}
	rate, err := strconv.ParseFloat(r.FormValue("review_rate"), 64)
	if err != nil || rate < 0 || rate > 5 {
		badRequest(w, "Calificación inválida (0 a 5)")
		return
	}
	date, err := time.Parse("2006-01-02", r.FormValue("review_date"))
	if err != nil {
		badRequest(w, "Fecha inválida")
		return
	}
	userID, err := intRequerido(r.FormValue("user_id_user"))
	if err != nil || userID < 1 {
		badRequest(w, "ID de usuario inválido")
		return
	}
	songID, err := nullInt32(r.FormValue("song_id_song"))
	if err != nil {
		badRequest(w, "ID de canción inválido")
		return
	}
	albumID, err := nullInt32(r.FormValue("album_id_album"))
	if err != nil {
		badRequest(w, "ID de álbum inválido")
		return
	}
	if !songID.Valid && !albumID.Valid {
		badRequest(w, "Indicá un ID de canción o de álbum")
		return
	}
	likes, err := intOpcional(r.FormValue("likes"))
	if err != nil || likes < 0 {
		badRequest(w, "Likes inválidos")
		return
	}
	_, err = queries.CreateReview(ctx, bd.CreateReviewParams{
		ReviewText: text, ReviewRate: rate, ReviewDate: date, UserIDUser: userID,
		SongIDSong: songID, AlbumIDAlbum: albumID, Likes: likes,
	})
	if err != nil {
		http.Error(w, "Error al crear la reseña", http.StatusInternalServerError)
		return
	}
	list, err := queries.ListReviews(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener reseñas", err, http.StatusInternalServerError)
		return
	}
	views.ReviewList(list).Render(ctx, w)
	views.FormReview(true).Render(ctx, w) // form vacío (OOB) → limpia el formulario
}

// postArtists atiende POST /artists: valida, crea y devuelve la lista de artistas.
func postArtists(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	artistName := strings.TrimSpace(r.FormValue("artist_name"))
	if artistName == "" {
		badRequest(w, "El nombre del artista es obligatorio")
		return
	}
	photo, err := intOpcional(r.FormValue("artist_ph"))
	if err != nil || photo < 0 {
		badRequest(w, "Valor de foto inválido")
		return
	}
	_, err = queries.CreateArtist(ctx, bd.CreateArtistParams{
		ArtistName: artistName, Photo: photo, ArtistBio: strings.TrimSpace(r.FormValue("artist_bio")),
	})
	if err != nil {
		http.Error(w, "Error al crear el artista", http.StatusInternalServerError)
		return
	}
	list, err := queries.ListArtists(ctx)
	if err != nil {
		errorServidor(w, "Error al obtener artistas", err, http.StatusInternalServerError)
		return
	}
	views.ArtistsList(list).Render(ctx, w)
	views.FormArtist(true).Render(ctx, w) // form vacío (OOB) → limpia el formulario
}

// HandleDELETE atiende DELETE /<entidad>/<id>
func HandleDELETE(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // [entidad, id]
	if len(partes) != 2 {
		http.NotFound(w, r)
		return
	}
	id, err := intRequerido(partes[1])
	if err != nil {
		badRequest(w, "ID inválido")
		return
	}

	ctx := r.Context()
	switch partes[0] {
	case "users":
		err = queries.DeleteUser(ctx, id)
	case "songs":
		err = queries.DeleteSong(ctx, id)
	case "albums":
		err = queries.DeleteAlbum(ctx, id)
	case "reviews":
		err = queries.DeleteReview(ctx, id)
	case "artists":
		err = queries.DeleteArtist(ctx, id)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		errorServidor(w, "No se pudo eliminar (puede tener datos relacionados)", err, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)

	// La respuesta solo contiene el contador con hx-swap-oob="true":
	// HTMX lo saca de la respuesta y reemplaza #lista-count en la página.
	// Lo que queda para el hx-target (el <li>) es vacío, así que el <li> desaparece.
	n, err := contarEntidades(ctx, queries, partes[0])
	if err != nil {
		return // el borrado ya se hizo; si falla el conteo, solo no se actualiza el número
	}
	views.ListCount(n, true).Render(ctx, w)
}

// contarEntidades devuelve cuántos elementos quedan de una entidad.
func contarEntidades(ctx context.Context, queries *bd.Queries, entidad string) (int, error) {
	switch entidad {
	case "users":
		l, err := queries.ListUsers(ctx)
		return len(l), err
	case "songs":
		l, err := queries.ListSongs(ctx)
		return len(l), err
	case "albums":
		l, err := queries.ListAlbums(ctx)
		return len(l), err
	case "reviews":
		l, err := queries.ListReviews(ctx)
		return len(l), err
	case "artists":
		l, err := queries.ListArtists(ctx)
		return len(l), err
	}
	return 0, nil
}
