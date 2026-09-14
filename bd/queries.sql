/* QUERIES DE ALBUM */
-- name: CreateAlbum :one
INSERT INTO Album (album_name, album_dur, album_ph, release_date, Artist_id_artist, Likes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: BuscarAlbumById :one
SELECT * FROM Album
WHERE id_album = $1;

    /* El user no sabrá el id del album, debe poder buscarlo por nombre
    Puede devolver más de uno si es por nombre por si dos albumes se llaman igual */
-- name: BuscarAlbumByName :many
SELECT * FROM Album
WHERE album_name = $1;

-- name: ListAlbums :many
SELECT * FROM Album
ORDER BY album_name;

-- name: UpdateAlbum :exec
UPDATE Album
SET album_name = $2, album_dur = $3, album_ph = $4, release_date = $5, Artist_id_artist = $6, Likes = $7
WHERE id_album = $1;

-- name: DeleteAlbum :exec
DELETE FROM Album
WHERE id_album = $1;

-- name: GetAlbumByArtist :many
SELECT * FROM Album
WHERE Artist_id_artist = $1;

-- name: GetAlbumLikes :many
SELECT SUM(Likes) AS Total_likes
FROM Album
WHERE id_album = $1;

/*  QUERIES DE ARTIST */
-- name: CreateArtist :one
INSERT INTO Artist (artist_name, photo, artist_bio)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetArtistById :one
SELECT * FROM Artist
WHERE id_artist = $1;

-- name: GetArtistByName :many
SELECT * FROM Artist
WHERE artist_name = $1;

-- name: ListArtists :many
SELECT * FROM Artist
ORDER BY id_artist;

-- name: UpdateArtist :one
UPDATE Artist
SET artist_name = $2, photo = $3, artist_bio = $4
WHERE id_artist = $1
RETURNING *;

-- name: DeleteArtist :exec
DELETE FROM Artist
WHERE id_artist = $1;

/*QUERIES DE SONG*/

-- name: CreateSong :one
INSERT INTO Song (song_name, song_dur, Album_id_album, likes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSongById :one
SELECT * FROM Song
WHERE id_song = $1;

-- name: GetSongByName :many
SELECT * FROM Song
WHERE song_name = $1;

-- name: GetSongByAlbum :many
SELECT * FROM Song
WHERE Album_id_album = $1;

-- name: ListSongs :many
SELECT * FROM Song
ORDER BY id_song;

-- name: GetSongLikes :one
SELECT SUM(Likes) AS Total_likes
FROM Song
WHERE id_song = $1;

-- name: GetAlbumSongById :one
SELECT  Album_id_album 
FROM Song
WHERE id_song = $1;

-- name: GetAlbumSongByName :many
SELECT  Album_id_album 
FROM Song
WHERE song_name = $1;

-- name: UpdateSong :one
UPDATE Song
SET song_name = $2, song_dur = $3, likes = $4
WHERE id_song = $1
RETURNING *;

-- name: DeleteSong :exec
DELETE FROM Song
WHERE id_song = $1;

/* QUERIES DE USER */

-- name: CreateUser :one
INSERT INTO "User" (user_name, amigos, photo)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM "User"
WHERE id_user = $1;

-- name: GetUserByName :many
SELECT * FROM "User"
WHERE user_name = $1;

-- name: ListUsers :many
SELECT * FROM "User"
ORDER BY id_user;

-- name: UpdateUser :one
UPDATE "User"
SET user_name = $2, amigos = $3, photo = $4
WHERE id_user = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "User"
WHERE id_user = $1;

/* REVIEW QUERIES */

-- name: CreateReview :one
INSERT INTO Review (review_text, review_rate, review_date, User_id_user, Song_id_song, Album_id_album, Likes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetReview :one
SELECT * FROM Review
WHERE id_review = $1;

-- name: ListReviews :many
SELECT * FROM Review
ORDER BY id_review;

-- name: UpdateReview :one
UPDATE Review
SET review_text = $2, review_rate = $3, review_date = $4, User_id_user = $5, Song_id_song = $6, Album_id_album = $7, Likes = $8
WHERE id_review = $1
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM Review
WHERE id_review = $1;

-- name: GetReviewLikes :one
SELECT sum(likes) as review_likes
FROM Review
WHERE id_review = $1;

/* QUERIES DE ARTIST_SONG (MANY TO MANY)*/

-- name: CreateArtistSong :one
INSERT INTO Artist_Song (Artist_id_artist, Song_id_song, Roll)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetArtistSong :one
SELECT * FROM Artist_Song
WHERE Artist_id_artist = $1 AND Song_id_song = $2;

-- name: ListArtistSongs :many
SELECT * FROM Artist_Song;

-- name: GetSongsByArtist :many
SELECT * FROM Artist_Song
WHERE Artist_id_artist = $1;

-- name: GetArtistsBySong :many
SELECT * FROM Artist_Song
WHERE Song_id_song = $1;

-- name: UpdateArtistSong :one
UPDATE Artist_Song
SET Roll = $3
WHERE Artist_id_artist = $1 AND Song_id_song = $2
RETURNING *;

-- name: DeleteArtistSong :exec
DELETE FROM Artist_Song
WHERE Artist_id_artist = $1 AND Song_id_song = $2;