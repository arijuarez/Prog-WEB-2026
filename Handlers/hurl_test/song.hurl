# ==========================================
# CREACIÓN DE ENTIDADES AUXILIARES 

# Auxiliar 1: Crear Artista
POST http://localhost:8080/artists
Content-Type: application/json
{
  "artist_name": "Pink Floyd",
  "photo": 1,
  "artist_bio": "Banda británica de rock"
}
HTTP 201
[Captures]
aux_artist_id: jsonpath "$.id_artist"

# Auxiliar 2: Crear Álbum
POST http://localhost:8080/albums
Content-Type: application/json
{
  "album_name": "The Dark Side of the Moon",
  "album_dur": 2580,
  "album_ph": 1,
  "release_date": "1973-03-01",
  "artist_id_artist": {{aux_artist_id}},
  "likes": 0
}
HTTP 201
[Captures]
aux_album_id: jsonpath "$.id_album"

# ==========================================

# 1. Crear nueva canción
POST http://localhost:8080/songs
Content-Type: application/json
{
  "song_name": "Time",
  "song_dur": 413,
  "likes": 0,
  "album_id_album": {{aux_album_id}}
}
HTTP 201
[Captures]
song_id: jsonpath "$.id_song"
[Asserts]
jsonpath "$.song_name" == "Time"
jsonpath "$.song_dur" == 413
jsonpath "$.likes" == 0
jsonpath "$.album_id_album" == {{aux_album_id}}

# 2. Obtener la canción recién creada
GET http://localhost:8080/songs/{{song_id}}
HTTP 200
[Asserts]
jsonpath "$.id_song" == {{song_id}}
jsonpath "$.song_name" == "Time"
jsonpath "$.song_dur" == 413
jsonpath "$.album_id_album" == {{aux_album_id}}

# 3. Listar canciones
GET http://localhost:8080/songs
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"

# 4. Actualizar canción
PUT http://localhost:8080/songs/{{song_id}}
Content-Type: application/json
{
  "song_name": "Time (2011 Remaster)",
  "song_dur": 413,
  "likes": 15,
  "album_id_album": {{aux_album_id}}
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.song_name" == "Time (2011 Remaster)"
jsonpath "$.likes" == 15

# 5. Verificar que se actualizó la canción
GET http://localhost:8080/songs/{{song_id}}
HTTP 200
[Asserts]
jsonpath "$.id_song" == {{song_id}}
jsonpath "$.song_name" == "Time (2011 Remaster)"
jsonpath "$.likes" == 15

# 6. Eliminar canción
DELETE http://localhost:8080/songs/{{song_id}}
HTTP 204

# 7. Verificar que se eliminó
GET http://localhost:8080/songs/{{song_id}}
HTTP 404

# ==========================================
# TESTS DE ERROR DE SONG

# 8. Intentar obtener canción no existente
GET http://localhost:8080/songs/999999
HTTP 404

# 9. Intento de creación con campos vacíos o inválidos
POST http://localhost:8080/songs
Content-Type: application/json
{
  "song_name": "   ",
  "song_dur": -10,
  "likes": 0,
  "album_id_album": 0
}
HTTP 400

# 10. Búsqueda con ID no numérico
GET http://localhost:8080/songs/abc
HTTP 400

# 11. Intentar actualizar una canción que no existe
PUT http://localhost:8080/songs/999999
Content-Type: application/json
{
  "song_name": "Inexistente",
  "song_dur": 200,
  "likes": 0,
  "album_id_album": {{aux_album_id}}
}
HTTP 404

# 12. Intentar eliminar una canción que no existe
DELETE http://localhost:8080/songs/999999
HTTP 404