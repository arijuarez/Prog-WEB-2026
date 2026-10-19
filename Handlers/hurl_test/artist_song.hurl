# ==========================================
# CREACIÓN DE ENTIDADES AUXILIARES

# Auxiliar 1: Crear Artista
POST http://localhost:8080/artists
Content-Type: application/json
{
  "artist_name": "Gorillaz",
  "photo": 1,
  "artist_bio": "Banda virtual británica"
}
HTTP 201
[Captures]
aux_artist_id: jsonpath "$.id_artist"

# Auxiliar 2: Crear Álbum
POST http://localhost:8080/albums
Content-Type: application/json
{
  "album_name": "Demon Days",
  "album_dur": 3043,
  "album_ph": 1,
  "release_date": "2005-05-11",
  "artist_id_artist": {{aux_artist_id}},
  "likes": 0
}
HTTP 201
[Captures]
aux_album_id: jsonpath "$.id_album"

# Auxiliar 3: Crear Canción
POST http://localhost:8080/songs
Content-Type: application/json
{
  "song_name": "Feel Good Inc.",
  "song_dur": 222,
  "likes": 0,
  "album_id_album": {{aux_album_id}}
}
HTTP 201
[Captures]
aux_song_id: jsonpath "$.id_song"

# ==========================================

# 1. Crear nueva relación artista-canción
POST http://localhost:8080/artist-songs
Content-Type: application/json
{
  "artist_id_artist": {{aux_artist_id}},
  "song_id_song": {{aux_song_id}},
  "roll": true
}
HTTP 201
[Asserts]
jsonpath "$.artist_id_artist" == {{aux_artist_id}}
jsonpath "$.song_id_song" == {{aux_song_id}}
jsonpath "$.roll" == true

# 2. Obtener relación por clave compuesta /artist-songs/{artist_id}/{song_id}
GET http://localhost:8080/artist-songs/{{aux_artist_id}}/{{aux_song_id}}
HTTP 200
[Asserts]
jsonpath "$.artist_id_artist" == {{aux_artist_id}}
jsonpath "$.song_id_song" == {{aux_song_id}}
jsonpath "$.roll" == true

# 3. Listar todas las relaciones
GET http://localhost:8080/artist-songs
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$" isArray
jsonpath "$.length()" == 1
jsonpath "$[0].artist_id_artist" == {{aux_artist_id}}
jsonpath "$[0].song_id_song" == {{aux_song_id}}
jsonpath "$[0].roll" == true

# 4. Actualizar el rol de la relación
PUT http://localhost:8080/artist-songs/{{aux_artist_id}}/{{aux_song_id}}
Content-Type: application/json
{
  "roll": false
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.artist_id_artist" == {{aux_artist_id}}
jsonpath "$.song_id_song" == {{aux_song_id}}
jsonpath "$.roll" == false

# 5. Verificar que se actualizó la relación
GET http://localhost:8080/artist-songs/{{aux_artist_id}}/{{aux_song_id}}
HTTP 200
[Asserts]
jsonpath "$.artist_id_artist" == {{aux_artist_id}}
jsonpath "$.song_id_song" == {{aux_song_id}}
jsonpath "$.roll" == false

# 6. Eliminar la relación
DELETE http://localhost:8080/artist-songs/{{aux_artist_id}}/{{aux_song_id}}
HTTP 204

# 7. Verificar que se eliminó
GET http://localhost:8080/artist-songs/{{aux_artist_id}}/{{aux_song_id}}
HTTP 404

# ==========================================
# TESTS DE ERROR DE ARTIST_SONG

# 8. Intentar obtener relación no existente
GET http://localhost:8080/artist-songs/999999/999999
HTTP 404

# 9. Intento de creación con campos inválidos
POST http://localhost:8080/artist-songs
Content-Type: application/json
{
  "artist_id_artist": 0,
  "song_id_song": -1,
  "roll": null
}
HTTP 400

# 10. Búsqueda con IDs no numéricos
GET http://localhost:8080/artist-songs/abc/xyz
HTTP 400

# 11. Intentar actualizar una relación que no existe
PUT http://localhost:8080/artist-songs/999999/999999
Content-Type: application/json
{
  "roll": true
}
HTTP 404

# 12. Intentar eliminar una relación que no existe
DELETE http://localhost:8080/artist-songs/999999/999999
HTTP 404