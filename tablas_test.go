package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	sqlc "tp2/bd/sqlc" // generado por sqlc

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	queries *sqlc.Queries
	ctx     context.Context
)

func TestMain(m *testing.M) {
	/* Fue comentado para que pueda accederse sin el .env
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable", user, password, db_name, host)
	*/
	connStr := "user=postgres password=12345 dbname=DB host=database port=5432 sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	queries = sqlc.New(db)
	ctx = context.Background()

	code := m.Run()

	os.Exit(code)

}

// Pruebas de Album
func TestAlbum(t *testing.T) {
	// Creación de artiste temporal para cumplir foreign key
	artist_temp, err := queries.CreateArtist(ctx, sqlc.CreateArtistParams{
		ArtistName: "Pink Floyd",
		Photo:      1,
		ArtistBio:  "Pink Floyd es una legendaria banda británica de rock formada en Londres en 1965. Es uno de los grupos más influyentes, exitosos e importantes de la historia de la música popular.",
	})

	createdAlbum, err := queries.CreateAlbum(ctx, sqlc.CreateAlbumParams{
		AlbumName:      "The Dark Side of the moon",
		AlbumDur:       230,
		AlbumPh:        1,
		ReleaseDate:    time.Now(),
		ArtistIDArtist: artist_temp.IDArtist,
		Likes:          0,
	})
	if err != nil {
		log.Fatalf("Failed to create album: %v", err)
	}
	fmt.Printf("Created user: %+v\n", createdAlbum)

	album, err := queries.BuscarAlbumById(ctx, createdAlbum.IDAlbum)
	if err != nil {
		log.Fatalf("Failed to get Album id %v\n", err)
	}
	fmt.Printf("Retrieved Album: %+v\n", album)

	albumsByName, err := queries.BuscarAlbumByName(ctx, createdAlbum.AlbumName)
	if err != nil {
		log.Fatalf("Failed to get Album id %v\n", err)
	}
	fmt.Printf("Retrieved Album: %+v\n", albumsByName)

	albums, err := queries.ListAlbums(ctx)
	if err != nil {
		log.Fatalf("Failed to get Album id %v\n", err)
	}
	fmt.Printf("Retrieved Album: %+v\n", albums)

	albumsByArtist, err := queries.GetAlbumByArtist(ctx, createdAlbum.ArtistIDArtist)
	if err != nil {
		log.Fatalf("Failed to get Album id %v\n", err)
	}
	fmt.Printf("Retrieved Album: %+v\n", albumsByArtist)

	albumLikes, err := queries.GetAlbumLikes(ctx, createdAlbum.IDAlbum)
	if err != nil {
		log.Fatalf("Failed to get Album id %v\n", err)
	}
	fmt.Printf("Retrieved Album: %+v\n", albumLikes)

	err = queries.UpdateAlbum(ctx, sqlc.UpdateAlbumParams{
		IDAlbum:        createdAlbum.IDAlbum,
		AlbumName:      "THE DARK SIDE OF THE MOON",
		AlbumDur:       300,
		AlbumPh:        1,
		ReleaseDate:    time.Now(),
		ArtistIDArtist: artist_temp.IDArtist,
		Likes:          3,
	})
	if err != nil {
		log.Fatalf("Falló al actualizar album %v\n", err)
	}
	fmt.Printf("Album actualizado")

	err = queries.DeleteAlbum(ctx, createdAlbum.IDAlbum)
	if err != nil {
		log.Fatalf("Falló al eliminar album %v\n", err)
	}
	fmt.Printf("Album eliminado:\n")

	// Eliminar entidades temporales
	_ = queries.DeleteArtist(ctx, artist_temp.IDArtist)

}

// Pruebas de Song
func TestSong(t *testing.T) {
	// Creacion de Artista temporal para cumplir Foreign Key
	artist_temp, err := queries.CreateArtist(ctx, sqlc.CreateArtistParams{
		ArtistName: "Pink Floyd",
		Photo:      1,
		ArtistBio:  "Pink Floyd es una legendaria banda británica de rock formada en Londres en 1965. Es uno de los grupos más influyentes, exitosos e importantes de la historia de la música popular.",
	})
	// Creacion de Album temporal para cumplir Foreign Key

	album_temp, err := queries.CreateAlbum(ctx, sqlc.CreateAlbumParams{
		AlbumName:      "The Dark Side of the moon",
		AlbumDur:       230,
		AlbumPh:        1,
		ReleaseDate:    time.Now(),
		ArtistIDArtist: artist_temp.IDArtist,
		Likes:          0,
	})

	//CreateSong
	createdSong, err := queries.CreateSong(ctx, sqlc.CreateSongParams{
		SongName:     "Hit me baby one more time",
		SongDur:      320,
		AlbumIDAlbum: album_temp.IDAlbum,
	})
	if err != nil {
		log.Fatalf("Falló al añadir la canción: %v", err)
	}
	fmt.Printf("Se añadió la canción: %v\n", createdSong)

	//GetSongById
	songById, err := queries.GetSongById(ctx, createdSong.IDSong)
	if err != nil {
		log.Fatalf("Falló al querer retornar la canción: %v", err)
	}
	fmt.Printf("Se retornó la canción: %v\n", songById)

	//GetSongByName
	songByName, err := queries.GetSongByName(ctx, createdSong.SongName)
	if err != nil {
		log.Fatalf("Falló al querer retornar la canción: %v", err)
	}
	fmt.Printf("Se retornó la canción: %v\n", songByName)

	//ListSongs
	songs, err := queries.ListSongs(ctx)
	if err != nil {
		log.Fatalf("Falló al listar las canciones: %v", err)
	}
	fmt.Printf("Todas las canciones: %v\n", songs)

	//GetSongLikes
	songLikes, err := queries.GetSongLikes(ctx, createdSong.IDSong)
	if err != nil {
		log.Fatalf("Falló al retornar los likes de la canción: %v", err)
	}
	fmt.Printf("Likes de la canción: %v\n", songLikes)

	//GetAlbumSongById
	songAlbumById, err := queries.GetAlbumSongById(ctx, createdSong.IDSong)
	if err != nil {
		log.Fatalf("Falló al retornar el album de la canción: %v", err)
	}
	fmt.Printf("Album de la canción: %v\n", songAlbumById)

	//GetAlbumSongByName
	songAlbumByName, err := queries.GetAlbumSongByName(ctx, createdSong.SongName)
	if err != nil {
		log.Fatalf("Falló al retornar el album de la canción: %v", err)
	}
	fmt.Printf("Album de la canción: %v\n", songAlbumByName)
	//UpdateSong
	updatedSong, err := queries.UpdateSong(ctx, sqlc.UpdateSongParams{
		IDSong:   createdSong.IDSong,
		SongName: "HIT ME BABY ONE MORE TIME",
		SongDur:  315,
	})
	if err != nil {
		log.Fatalf("Falló al actualizar la canción: %v", err)
	}
	fmt.Printf("La canción fue actualizada: %v\n", updatedSong)

	//DeleteSong
	err = queries.DeleteSong(ctx, createdSong.IDSong)
	if err != nil {
		log.Fatalf("Falló al eliminar la canción: %v", err)
	}
	fmt.Printf("La canción fue eliminada:\n")

	// Borrar registros temporales
	_ = queries.DeleteArtist(ctx, artist_temp.IDArtist)
	_ = queries.DeleteAlbum(ctx, album_temp.IDAlbum)
}

func TestArtist(t *testing.T) {
	// CreateArtist
	createdArtist, err := queries.CreateArtist(ctx, sqlc.CreateArtistParams{
		ArtistName: "Daft Punk",
		Photo:      1,
		ArtistBio:  "Dúo francés de música electrónica formado en 1993 en París.",
	})
	if err != nil {
		log.Fatalf("Falló al crear el artista: %v", err)
	}
	fmt.Printf("Artista creado: %+v\n", createdArtist)

	// GetArtistById
	artistById, err := queries.GetArtistById(ctx, createdArtist.IDArtist)
	if err != nil {
		log.Fatalf("Falló al obtener el artista por ID: %v", err)
	}
	fmt.Printf("Artista por ID: %+v\n", artistById)

	// GetArtistByName
	artistsByName, err := queries.GetArtistByName(ctx, createdArtist.ArtistName)
	if err != nil {
		log.Fatalf("Falló al obtener el artista por nombre: %v", err)
	}
	fmt.Printf("Artistas por nombre: %+v\n", artistsByName)

	// ListArtists
	artists, err := queries.ListArtists(ctx)
	if err != nil {
		log.Fatalf("Falló al listar artistas: %v", err)
	}
	fmt.Printf("Lista de artistas: %+v\n", artists)

	// UpdateArtist
	updatedArtist, err := queries.UpdateArtist(ctx, sqlc.UpdateArtistParams{
		IDArtist:   createdArtist.IDArtist,
		ArtistName: "Daft Punk (Oficial)",
		Photo:      2,
		ArtistBio:  "Dúo legendario de música electrónica.",
	})
	if err != nil {
		log.Fatalf("Falló al actualizar el artista: %v", err)
	}
	fmt.Printf("Artista actualizado: %+v\n", updatedArtist)

	// DeleteArtist
	err = queries.DeleteArtist(ctx, createdArtist.IDArtist)
	if err != nil {
		log.Fatalf("Falló al eliminar el artista: %v", err)
	}
	fmt.Println("Artista eliminado con éxito")
}

// Pruebas de User

func TestUser(t *testing.T) {
	// CreateUser
	createdUser, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		UserName: "jambox_user",
		Amigos:   5,
		Photo:    1,
	})
	if err != nil {
		log.Fatalf("Falló al crear el usuario: %v", err)
	}
	fmt.Printf("Usuario creado: %+v\n", createdUser)

	// GetUserById
	userById, err := queries.GetUserById(ctx, createdUser.IDUser)
	if err != nil {
		log.Fatalf("Falló al obtener el usuario por ID: %v", err)
	}
	fmt.Printf("Usuario por ID: %+v\n", userById)

	// GetUserByName
	usersByName, err := queries.GetUserByName(ctx, createdUser.UserName)
	if err != nil {
		log.Fatalf("Falló al obtener el usuario por nombre: %v", err)
	}
	fmt.Printf("Usuarios por nombre: %+v\n", usersByName)

	// ListUsers
	users, err := queries.ListUsers(ctx)
	if err != nil {
		log.Fatalf("Falló al listar usuarios: %v", err)
	}
	fmt.Printf("Lista de usuarios: %+v\n", users)

	// UpdateUser
	updatedUser, err := queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		IDUser:   createdUser.IDUser,
		UserName: "jambox_user_updated",
		Amigos:   10,
		Photo:    2,
	})
	if err != nil {
		log.Fatalf("Falló al actualizar el usuario: %v", err)
	}
	fmt.Printf("Usuario actualizado: %+v\n", updatedUser)

	// DeleteUser
	err = queries.DeleteUser(ctx, createdUser.IDUser)
	if err != nil {
		log.Fatalf("Falló al eliminar el usuario: %v", err)
	}
	fmt.Println("Usuario eliminado con éxito")
}

// Pruebas Review

func TestReview(t *testing.T) {
	// Crear usuario, artista, album y canción temporales para la Clave Foránea (User_id_user, song_id_song, Album_id_album)
	tempUser, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		UserName: "reviewer_test",
		Amigos:   0,
		Photo:    1,
	})
	if err != nil {
		log.Fatalf("Falló al crear usuario para la reseña: %v", err)
	}
	artist_temp, err := queries.CreateArtist(ctx, sqlc.CreateArtistParams{
		ArtistName: "Pink Floyd",
		Photo:      1,
		ArtistBio:  "Pink Floyd es una legendaria banda británica de rock formada en Londres en 1965. Es uno de los grupos más influyentes, exitosos e importantes de la historia de la música popular.",
	})
	if err != nil {
		log.Fatalf("Falló al crear artista para la reseña: %v", err)
	}
	album_temp, err := queries.CreateAlbum(ctx, sqlc.CreateAlbumParams{
		AlbumName:      "The Dark Side of the moon",
		AlbumDur:       230,
		AlbumPh:        1,
		ReleaseDate:    time.Now(),
		ArtistIDArtist: artist_temp.IDArtist,
		Likes:          100,
	})
	if err != nil {
		log.Fatalf("Falló al crear album para la reseña: %v", err)
	}
	song_temp, err := queries.CreateSong(ctx, sqlc.CreateSongParams{
		SongName:     "Hit me baby one more time",
		SongDur:      320,
		AlbumIDAlbum: album_temp.IDAlbum,
	})
	if err != nil {
		log.Fatalf("Falló al crear canción para la reseña: %v", err)
	}

	// CreateReview
	createdReview, err := queries.CreateReview(ctx, sqlc.CreateReviewParams{
		ReviewText:   "Un álbum excelente de principio a fin.",
		ReviewRate:   4.8,
		ReviewDate:   time.Now(),
		UserIDUser:   tempUser.IDUser,
		SongIDSong:   sql.NullInt32{Int32: song_temp.IDSong, Valid: true},
		AlbumIDAlbum: sql.NullInt32{Int32: album_temp.IDAlbum, Valid: true},
	})
	if err != nil {
		log.Fatalf("Falló al crear la reseña: %v", err)
	}
	fmt.Printf("Reseña creada: %+v\n", createdReview)

	// GetReview
	reviewById, err := queries.GetReview(ctx, createdReview.IDReview)
	if err != nil {
		log.Fatalf("Falló al obtener la reseña por ID: %v", err)
	}
	fmt.Printf("Reseña obtenida: %+v\n", reviewById)

	// GetReviewLikes
	reviewLikes, err := queries.GetReviewLikes(ctx, createdReview.IDReview)
	if err != nil {
		log.Fatalf("Falló al obtener los likes de la reseña: %v", err)
	}
	fmt.Printf("Likes de la reseña: %v\n", reviewLikes)

	// ListReviews
	reviews, err := queries.ListReviews(ctx)
	if err != nil {
		log.Fatalf("Falló al listar reseñas: %v", err)
	}
	fmt.Printf("Lista de reseñas: %+v\n", reviews)

	// UpdateReview
	updatedReview, err := queries.UpdateReview(ctx, sqlc.UpdateReviewParams{
		IDReview:     createdReview.IDReview,
		ReviewText:   "Edición: Definitivamente una obra maestra.",
		ReviewRate:   5.0,
		ReviewDate:   time.Now(),
		UserIDUser:   tempUser.IDUser,
		SongIDSong:   sql.NullInt32{Int32: song_temp.IDSong, Valid: true},
		AlbumIDAlbum: sql.NullInt32{Int32: album_temp.IDAlbum, Valid: true},
	})
	if err != nil {
		log.Fatalf("Falló al actualizar la reseña: %v", err)
	}
	fmt.Printf("Reseña actualizada: %+v\n", updatedReview)

	// DeleteReview
	err = queries.DeleteReview(ctx, createdReview.IDReview)
	if err != nil {
		log.Fatalf("Falló al eliminar la reseña: %v", err)
	}
	fmt.Println("Reseña eliminada con éxito")

	// Limpieza de usuario temporal
	_ = queries.DeleteUser(ctx, tempUser.IDUser)
	_ = queries.DeleteUser(ctx, song_temp.IDSong)
	_ = queries.DeleteAlbum(ctx, album_temp.IDAlbum)
	_ = queries.DeleteArtist(ctx, artist_temp.IDArtist)

}

// Pruebas de ArtistSong

func TestArtistSong(t *testing.T) {
	// Crear artista y canción temporales para las Claves Foráneas
	artist_temp, err := queries.CreateArtist(ctx, sqlc.CreateArtistParams{
		ArtistName: "Pink Floyd",
		Photo:      1,
		ArtistBio:  "Pink Floyd es una legendaria banda británica de rock formada en Londres en 1965. Es uno de los grupos más influyentes, exitosos e importantes de la historia de la música popular.",
	})
	// Creacion de Album temporal para cumplir Foreign Key

	album_temp, err := queries.CreateAlbum(ctx, sqlc.CreateAlbumParams{
		AlbumName:      "The Dark Side of the moon",
		AlbumDur:       230,
		AlbumPh:        1,
		ReleaseDate:    time.Now(),
		ArtistIDArtist: artist_temp.IDArtist,
		Likes:          0,
	})

	//CreateSong
	createdSong, err := queries.CreateSong(ctx, sqlc.CreateSongParams{
		SongName:     "Hit me baby one more time",
		SongDur:      320,
		AlbumIDAlbum: album_temp.IDAlbum,
	})
	if err != nil {
		log.Fatalf("Falló al añadir la canción: %v", err)
	}

	tempArtist2, err := queries.CreateArtist(ctx, sqlc.CreateArtistParams{
		ArtistName: "Artista Colaborador",
		Photo:      1,
		ArtistBio:  "Bio temporal",
	})
	if err != nil {
		log.Fatalf("Falló al crear artista temporal: %v", err)
	}

	// CreateArtistSong
	createdArtistSong1, err := queries.CreateArtistSong(ctx, sqlc.CreateArtistSongParams{
		ArtistIDArtist: artist_temp.IDArtist,
		SongIDSong:     createdSong.IDSong,
		Roll:           true,
	})
	if err != nil {
		log.Fatalf("Falló al crear la relación Artist_Song: %v", err)
	}
	fmt.Printf("Relación Artist_Song creada: %+v\n", createdArtistSong1)

	createdArtistSong2, err := queries.CreateArtistSong(ctx, sqlc.CreateArtistSongParams{
		ArtistIDArtist: tempArtist2.IDArtist,
		SongIDSong:     createdSong.IDSong,
		Roll:           false,
	})
	if err != nil {
		log.Fatalf("Falló al crear la relación Artist_Song: %v", err)
	}
	fmt.Printf("Relación Artist_Song creada: %+v\n", createdArtistSong2)

	// GetArtistSong
	artistSong1, err := queries.GetArtistSong(ctx, sqlc.GetArtistSongParams{
		ArtistIDArtist: artist_temp.IDArtist,
		SongIDSong:     createdSong.IDSong,
	})
	if err != nil {
		log.Fatalf("Falló al obtener la relación Artist_Song: %v", err)
	}
	fmt.Printf("Relación Artist_Song obtenida: %+v\n", artistSong1)

	artistSong2, err := queries.GetArtistSong(ctx, sqlc.GetArtistSongParams{
		ArtistIDArtist: artist_temp.IDArtist,
		SongIDSong:     createdSong.IDSong,
	})
	if err != nil {
		log.Fatalf("Falló al obtener la relación Artist_Song: %v", err)
	}
	fmt.Printf("Relación Artist_Song obtenida: %+v\n", artistSong2)

	// ListArtistSongs
	artistSongs, err := queries.ListArtistSongs(ctx)
	if err != nil {
		log.Fatalf("Falló al listar relaciones Artist_Song: %v", err)
	}
	fmt.Printf("Lista de relaciones Artist_Song: %+v\n", artistSongs)

	SongsByArtist, err := queries.GetSongsByArtist(ctx, artist_temp.IDArtist)
	if err != nil {
		log.Fatalf("Falló al listar las canciones de el artista: %v", err)
	}
	fmt.Printf("Lista de relaciones Artist_Song: %+v\n", SongsByArtist)

	ArtistsBySong, err := queries.GetArtistsBySong(ctx, createdSong.IDSong)
	if err != nil {
		log.Fatalf("Falló al listar los artistas de la cancion: %v", err)
	}
	fmt.Printf("Lista de relaciones Artist_Song: %+v\n", ArtistsBySong)

	// UpdateArtistSong
	updatedArtistSong, err := queries.UpdateArtistSong(ctx, sqlc.UpdateArtistSongParams{
		ArtistIDArtist: tempArtist2.IDArtist,
		SongIDSong:     createdSong.IDSong,
		Roll:           true,
	})
	if err != nil {
		log.Fatalf("Falló al actualizar la relación Artist_Song: %v", err)
	}
	fmt.Printf("Relación Artist_Song actualizada: %+v\n", updatedArtistSong)

	// DeleteArtistSong
	err = queries.DeleteArtistSong(ctx, sqlc.DeleteArtistSongParams{
		ArtistIDArtist: tempArtist2.IDArtist,
		SongIDSong:     createdSong.IDSong,
	})
	if err != nil {
		log.Fatalf("Falló al eliminar la relación Artist_Song: %v", err)
	}
	fmt.Println("Relación Artist_Song eliminada con éxito")

	// Limpieza de entidades temporales
	_ = queries.DeleteSong(ctx, createdSong.IDSong)
	_ = queries.DeleteAlbum(ctx, album_temp.IDAlbum)
	_ = queries.DeleteArtist(ctx, artist_temp.IDArtist)
	_ = queries.DeleteArtist(ctx, tempArtist2.IDArtist)
}
