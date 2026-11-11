package Handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	bd "github.com/arijuarez/Prog-WEB-2026/bd/sqlc"
	views "github.com/arijuarez/Prog-WEB-2026/views"
)

func handleGET(w http.ResponseWriter, r *http.Request, queries *bd.Queries){
	ctx := r.Context()
	switch r.URL.Path {
	case "/":
		views.IndexMain().Render(ctx, w)
	case "/users":
		usersList, err := queries.ListUsers(ctx)
		if err != nil {
			http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexUsers(usersList).Render(ctx, w)
	case "/songs":
		songsList, err := queries.ListSongs(ctx)
		if err != nil {
			http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexSongs(songsList).Render(ctx, w)
	case "/albums":
		albumsList, err := queries.ListAlbums(ctx)
		if err != nil {
			http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexAlbums(albumsList).Render(ctx, w)
	case "/reviews":
		reviewsList, err := queries.ListReviews(ctx)
		if err != nil {
			http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexReviews(userReview).Render(ctx, w)
	case "/artists":
		artistsList, err := queries.ListArtists(ctx)
		if err != nil {
			http.Error(w, "Error al obtener usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}
		views.IndexArtists(artistsList).Render(ctx, w)
	default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
}

func handlePOST(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	//Parsea los datos del formulario.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}
	switch r.URL.Path {
		case "/": {}
		case "/users":{
			//Obtiene los valores del formulario.
			UserName := r.FormValue("user_name")
			Amigos:= r.FormValue("amigos")
			Photo:= r.FormValue("user_ph")

			if UserName == "" || Amigos == "" || Photo == "" {
				http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
				return
			}

			// formatea los amigos a int
			Amigos, err := strconv.Atoi(Amigos)
			if err != nil {
				http.Error(w, "Cantidad de amigos inválida", http.StatusBadRequest)
				return
			}

			Photo, err := strconv.Atoi(Photo)
			if err != nil {
				http.Error(w, "Valor de foto inválido", http.StatusBadRequest)
				return
			}

			// chequea amigos validos
			if Amigos < 0 {
				http.Error(w, "La cantidad de amigos no puede ser negativa", http.StatusBadRequest)
				return
			}

			params := bd.CreateUserParams{
				UserName: UserName,
				Amigos:   Amigos,
				Photo:    Photo,
			}

			_, err := h.queries.CreateUser(r.Context(), params)
			if err != nil {
				http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
				return
			}
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		}
		case "/songs":{
			//Obtiene los valores del formulario.
			SongName := r.FormValue("song_name")
			SongDur := r.FormValue("song_dur")
			AlbumIDAlbum := r.FormValue("album_id_album")
			Likes := r.FormValue("likes")
			

			if SongName == "" || SongDur == "" || AlbumIDAlbum == "" || Likes == "" {
				http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
				return
			}

			// formatea likes a int
			Likes, err := strconv.Atoi(Likes)
			if err != nil {
				http.Error(w, "Likes inválidos", http.StatusBadRequest)
				return
			}

			AlbumIDAlbum, err := strconv.Atoi(AlbumIDAlbum)
			if err != nil {
				http.Error(w, "Id album inválido", http.StatusBadRequest)
				return
			}
			
			SongDur, err := strconv.Atoi(SongDur)
			if err != nil {
				http.Error(w, "SongDur inválido", http.StatusBadRequest)
				return
			}

			// chequea likes validos
			if Likes < 0 {
				http.Error(w, "La cantidad de likes no puede ser negativa", http.StatusBadRequest)
				return
			}
			if AlbumIDAlbum < 0 {
				http.Error(w, "La cantidad de likes no puede ser negativa", http.StatusBadRequest)
				return
			}
			if SongDur < 0 {
				http.Error(w, "La cantidad de likes no puede ser negativa", http.StatusBadRequest)
				return
			}

			params := bd.CreateSongParams{
				SongName:     SongName,
				SongDur:      SongDur,
				AlbumIDAlbum: AlbumIDAlbum,
				Likes:        Likes,
			}
		
			_, err := h.queries.CreateSong(r.Context(), params)
			if err != nil {
				http.Error(w, "Error al crear la cancion", http.StatusInternalServerError)
				return
			}
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
		}
		case "/albums":{
						//Obtiene los valores del formulario.
			AlbumName := r.FormValue("album_name")
			AlbumDur:= r.FormValue("artist_id_artist")
			AlbumPh:= r.FormValue("release_date")
			ReleaseDate:= r.FormValue("album_ph")
			ArtistIDArtist:= r.FormValue("album_dur")
			Likes:= r.FormValue("likes")

			if AlbumName == "" || AlbumDur == "" || AlbumPh == "" || ReleaseDate == "" || ArtistIDArtist == "" || Likes == "" {
				http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
				return
			}

			// formatea los amia int
			AlbumDur, err := strconv.Atoi(AlbumDur)
			if err != nil {
				http.Error(w, "Duracion de album inválido", http.StatusBadRequest)
				return
			}

			ArtistIDArtist, err := strconv.Atoi(ArtistIDArtist)
			if err != nil {
				http.Error(w, "Id artist invalido", http.StatusBadRequest)
				return
			}

			AlbumPh, err := strconv.Atoi(AlbumPh)
			if err != nil {
				http.Error(w, "Foto inválida", http.StatusBadRequest)
				return
			}

			Likes, err := strconv.Atoi(Likes)
			if err != nil {
				http.Error(w, "Likes inválidos", http.StatusBadRequest)
				return
			}

			// chequea likes validos
			if ArtistIDArtist < 0 {
				http.Error(w, "El id artista no puede ser negativo", http.StatusBadRequest)
				return
			}
			if AlbumDur < 0 {
				http.Error(w, "La duracion no puede ser negativa", http.StatusBadRequest)
				return
			}
			if AlbumPh < 0 {
				http.Error(w, "Valor de foto invalido", http.StatusBadRequest)
				return
			}
			if Likes < 0 {
				http.Error(w, "La cantidad de likes no puede ser negativa", http.StatusBadRequest)
				return
			}

			params := bd.CreateUserParams{
   			AlbumName := AlbumName,
			AlbumDur:= AlbumDur,
			AlbumPh:= AlbumPh,
			ReleaseDate:= ReleaseDate,
			ArtistIDArtist:= ArtistIDArtist,
			Likes:= Likes,
			}

			_, err := h.queries.CreateUser(r.Context(), params)
			if err != nil {
				http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/albums", http.StatusSeeOther)
		}
		case "/reviews":{
			ReviewText := r.FormValue("review_text")
			ReviewRate:= r.FormValue("review_rate")
			ReviewDate:= r.FormValue("review_date")
			UserIDUser:= r.FormValue("user_id_user")
			SongIDSong:= r.FormValue("song_id_song")
			AlbumIDAlbum:= r.FormValue("album_id_album")
			Likes:= r.FormValue("likes")

			if ReviewText == "" || ReviewRate == "" || 
				ReviewDate == "" || UserIDUser == "" || 
				(SongIDSong =="" & AlbumIDAlbum == "") ||
				Likes =="" {
				http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
				return
			}

			// formatea los likes a int
			Likes, err := strconv.Atoi(Likes)
			if err != nil {
				http.Error(w, "Likes inválidos", http.StatusBadRequest)
				return
			}

			// chequea likes validos
			if Likes < 0 {
				http.Error(w, "La cantidad de likes no puede ser negativa", http.StatusBadRequest)
				return
			}

			// formatea los id a int
			UserIDUser, err := strconv.Atoi(UserIDUser)
			if err != nil {
				http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
				return
			}
			SongIDSong, err := strconv.Atoi(SongIDSong)
			if err != nil {
				http.Error(w, "ID de canción inválido", http.StatusBadRequest)
				return
			}
			AlbumIDAlbum, err := strconv.Atoi(AlbumIDAlbum)
			if err != nil {
				http.Error(w, "ID de álbum inválido", http.StatusBadRequest)
				return
			}

			params := bd.CreateReviewParams{
				ReviewText:   ReviewText,
				ReviewRate:   ReviewRate,
				ReviewDate:   reviewDate,
				UserIDUser:   UserID,
				SongIDSong:   toNullInt32(SongIDSong),
				AlbumIDAlbum: toNullInt32(AlbumIDAlbum),
				Likes:        Likes,
			}
		
			_, err := h.queries.CreateReview(r.Context(), params)
			if err != nil {
				http.Error(w, "Error al crear la review", http.StatusInternalServerError)
				return
			}
		}
		case "/artists": {
			ArtistName := r.FormValue("artist_name")
			Photo:= r.FormValue("artist_ph")
			ArtistBio:= r.FormValue("likes")

			if ArtistName == "" || Photo == "" || 
				ArtistBio == "" {
				http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
				return
			}

			// formatea la foto a int
			Photo, err := strconv.Atoi(Likes)
			if err != nil {
				http.Error(w, "Likes inválidos", http.StatusBadRequest)
				return
			}

			params := bd.CreateArtistParams{
				ArtistName: ArtistName,
				Photo:      Photo,
				ArtistBio:  ArtistBio,
			}
			_, err := h.queries.CreateArtist(r.Context(), params)

			if err != nil {
				http.Error(w, "Error al crear el artista", http.StatusInternalServerError)
				return
			}
		}
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
	// Redirije al usuario de vuelta a la página principal.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
