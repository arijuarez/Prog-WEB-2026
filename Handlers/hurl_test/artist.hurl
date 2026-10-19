# 1. Crear un artista (POST /artists)
POST http://localhost:8080/artists
Content-Type: application/json
{
  "artist_name": "Charly García",
  "photo": 1,
  "artist_bio": "Músico, compositor y productor argentino de rock."
}
HTTP 201
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.id_artist" exists
jsonpath "$.artist_name" == "Charly García"
jsonpath "$.photo" == 1
jsonpath "$.artist_bio" == "Músico, compositor y productor argentino de rock."
[Captures]
artist_id: jsonpath "$.id_artist"

# 2. Listar todos los artistas (GET /artists)
GET http://localhost:8080/artists
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$" isArray

# 3. Obtener el artista recién creado (GET /artists/{id})
GET http://localhost:8080/artists/{{artist_id}}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.id_artist" == {{artist_id}}
jsonpath "$.artist_name" == "Charly García"

# 4. Actualizar el artista (PUT /artists/{id})
PUT http://localhost:8080/artists/{{artist_id}}
Content-Type: application/json
{
  "artist_name": "Charly García (Leyenda)",
  "photo": 2,
  "artist_bio": "Icono fundamental del rock nacional argentino."
}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.id_artist" == {{artist_id}}
jsonpath "$.artist_name" == "Charly García (Leyenda)"

# 5. Eliminar el artista (DELETE /artists/{id})
DELETE http://localhost:8080/artists/{{artist_id}}
HTTP 204

# 6. Verificar que ya no existe (GET /artists/{id})
GET http://localhost:8080/artists/{{artist_id}}
HTTP 404

# ======================
# 7. Pruebas de errores

# Intentar obtener un artista que no existe
GET http://localhost:8080/artists/999999
HTTP 404

# Intento de creación con campos vacíos
POST http://localhost:8080/artists
Content-Type: application/json
{
  "artist_name": "   ",
  "photo": 0,
  "artist_bio": ""
}
HTTP 400

# Búsqueda con un ID que no es numérico
GET http://localhost:8080/artists/abc
HTTP 400

# Intentar actualizar un artista que no existe
PUT http://localhost:8080/artists/999999
HTTP 404

# Intentar eliminar un artista que no existe
DELETE http://localhost:8080/artists/999999
HTTP 404