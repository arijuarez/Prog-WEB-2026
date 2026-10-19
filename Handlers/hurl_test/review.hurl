# ==========================================
# CREACIÓN DE ENTIDADES AUXILIARES 

# Auxiliar 1: Crear Usuario
POST http://localhost:8080/users
Content-Type: application/json
{
  "user_name": "Crítico Musical",
  "amigos": 5,
  "photo": 1
}
HTTP 201
[Captures]
aux_user_id: jsonpath "$.id_user"

# Auxiliar 2: Crear Artista
POST http://localhost:8080/artists
Content-Type: application/json
{
  "artist_name": "Daft Punk",
  "photo": 1,
  "artist_bio": "Dúo de música electrónica"
}
HTTP 201
[Captures]
aux_artist_id: jsonpath "$.id_artist"

# Auxiliar 3: Crear Álbum
POST http://localhost:8080/albums
Content-Type: application/json
{
  "album_name": "Random Access Memories",
  "album_dur": 4438,
  "album_ph": 1,
  "release_date": "2013-05-17",
  "artist_id_artist": {{aux_artist_id}},
  "likes": 0
}
HTTP 201
[Captures]
aux_album_id: jsonpath "$.id_album"

# Auxiliar 4: Crear Canción
POST http://localhost:8080/songs
Content-Type: application/json
{
  "song_name": "Get Lucky",
  "song_dur": 248,
  "likes": 0,
  "album_id_album": {{aux_album_id}}
}
HTTP 201
[Captures]
aux_song_id: jsonpath "$.id_song"

# ==========================================

# 1. Crear nueva review
POST http://localhost:8080/reviews
Content-Type: application/json
{
  "review_text": "Excelente álbum y producción impecable.",
  "review_rate": 4.8,
  "review_date": "2024-05-10",
  "user_id_user": {{aux_user_id}},
  "song_id_song": {{aux_song_id}},
  "album_id_album": {{aux_album_id}},
  "likes": 0
}
HTTP 201
[Captures]
review_id: jsonpath "$.id_review"
[Asserts]
jsonpath "$.review_text" == "Excelente álbum y producción impecable."
jsonpath "$.review_rate" == 4.8
jsonpath "$.user_id_user" == {{aux_user_id}}
jsonpath "$.song_id_song" == {{aux_song_id}}
jsonpath "$.album_id_album" == {{aux_album_id}}

# 2. Obtener la review recién creada
GET http://localhost:8080/reviews/{{review_id}}
HTTP 200
[Asserts]
jsonpath "$.id_review" == {{review_id}}
jsonpath "$.review_text" == "Excelente álbum y producción impecable."
jsonpath "$.user_id_user" == {{aux_user_id}}

# 3. Listar reviews
GET http://localhost:8080/reviews
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$" isArray
jsonpath "$.length()" == 1
jsonpath "$[0].id_review" == {{review_id}}
jsonpath "$[0].review_text" == "Excelente álbum y producción impecable."
jsonpath "$[0].user_id_user" == {{aux_user_id}}

# 4. Actualizar review
PUT http://localhost:8080/reviews/{{review_id}}
Content-Type: application/json
{
  "review_text": "Una obra maestra absoluta del siglo XXI.",
  "review_rate": 5.0,
  "review_date": "2024-05-11",
  "user_id_user": {{aux_user_id}},
  "song_id_song": {{aux_song_id}},
  "album_id_album": {{aux_album_id}},
  "likes": 12
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.review_text" == "Una obra maestra absoluta del siglo XXI."
jsonpath "$.review_rate" == 5.0
jsonpath "$.likes" == 12

# 5. Verificar que se actualizó la review
GET http://localhost:8080/reviews/{{review_id}}
HTTP 200
[Asserts]
jsonpath "$.id_review" == {{review_id}}
jsonpath "$.review_text" == "Una obra maestra absoluta del siglo XXI."
jsonpath "$.review_rate" == 5.0
jsonpath "$.likes" == 12

# 6. Eliminar review
DELETE http://localhost:8080/reviews/{{review_id}}
HTTP 204

# 7. Verificar que se eliminó
GET http://localhost:8080/reviews/{{review_id}}
HTTP 404

# ==========================================
# TESTS DE ERROR DE REVIEW

# 8. Intentar obtener review no existente
GET http://localhost:8080/reviews/999999
HTTP 404

# 9. Intento de creación con campos vacíos o inválidos
POST http://localhost:8080/reviews
Content-Type: application/json
{
  "review_text": "   ",
  "review_rate": -1.0,
  "review_date": "invalid-date",
  "user_id_user": 0,
  "song_id_song": null,
  "album_id_album": null,
  "likes": 0
}
HTTP 400

# 10. Búsqueda con ID no numérico
GET http://localhost:8080/reviews/abc
HTTP 400

# 11. Intentar actualizar una review que no existe
PUT http://localhost:8080/reviews/999999
Content-Type: application/json
{
  "review_text": "Texto válido",
  "review_rate": 4.0,
  "review_date": "2024-05-10",
  "user_id_user": {{aux_user_id}},
  "song_id_song": {{aux_song_id}},
  "album_id_album": {{aux_album_id}},
  "likes": 0
}
HTTP 404

# 12. Intentar eliminar una review que no existe
DELETE http://localhost:8080/reviews/999999
HTTP 404