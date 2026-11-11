package Handlers

import (
	"database/sql"
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

func badRequest(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

// --------- GET ---------

func HandleGET(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	switch r.URL.Path {
	case "/":
		views.IndexMain().Render(ctx, w)
	case "/users":
		list, err := queries.ListUsers(ctx)
		if err != nil {
			http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexUsers(list).Render(ctx, w)
	case "/songs":
		list, err := queries.ListSongs(ctx)
		if err != nil {
			http.Error(w, "Error al obtener canciones: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexSongs(list).Render(ctx, w)
	case "/albums":
		list, err := queries.ListAlbums(ctx)
		if err != nil {
			http.Error(w, "Error al obtener álbumes: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexAlbums(list).Render(ctx, w)
	case "/reviews":
		list, err := queries.ListReviews(ctx)
		if err != nil {
			http.Error(w, "Error al obtener reseñas: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexReviews(list).Render(ctx, w)
	case "/artists":
		list, err := queries.ListArtists(ctx)
		if err != nil {
			http.Error(w, "Error al obtener artistas: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexArtists(list).Render(ctx, w)
	default:
		http.NotFound(w, r)
	}
}

// --------- POST ---------

func HandlePOST(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		badRequest(w, "Error al parsear el formulario")
		return
	}

	// Borrado: POST /<entidad>/<id>/delete
	if strings.HasSuffix(r.URL.Path, "/delete") {
		handleDelete(w, r, queries)
		return
	}

	switch r.URL.Path {
	case "/users":
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
		http.Redirect(w, r, "/users", http.StatusSeeOther)

	case "/songs":
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
		http.Redirect(w, r, "/songs", http.StatusSeeOther)

	case "/albums":
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
		http.Redirect(w, r, "/albums", http.StatusSeeOther)

	case "/reviews":
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
		http.Redirect(w, r, "/reviews", http.StatusSeeOther)

	case "/artists":
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
		http.Redirect(w, r, "/artists", http.StatusSeeOther)

	default:
		http.NotFound(w, r)
	}
}

// handleDelete atiende POST /<entidad>/<id>/delete
func handleDelete(w http.ResponseWriter, r *http.Request, queries *bd.Queries) {
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // [entidad, id, "delete"]
	if len(partes) != 3 {
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
		http.Error(w, "No se pudo eliminar (puede tener datos relacionados): "+err.Error(), http.StatusConflict)
		return
	}
	http.Redirect(w, r, "/"+partes[0], http.StatusSeeOther)
}
