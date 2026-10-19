#album requests.hurl

# CREACION DE ARTISTA AUXILIAR
POST http://localhost:8080/artists
Content-Type: application/json
{
	"artist_name": "Pink Floyd",
    "photo": 1,
    "artist_bio": "influyente banda britanica"
}
HTTP 201
[Asserts]
jsonpath "$.id_artist" isInteger
jsonpath "$.artist_name" == "Pink Floyd"
jsonpath "$.photo" == 1
jsonpath "$.artist_name" == "influyente banda britanica"
[Captures]
auxArtist_id: jsonpath "$.id_artist"

# CREACION DE ALBUM1
POST http://localhost:8080/albums
Content-Type: application/json
{
	"album_name": "The Dark Side of the Moon",
	"album_dur": 100 ,
	"album_ph": 1,
	"release_date": "1973-03-01",
	"artist_id_artist": {{auxArtist_id}},
	"likes": 0
}
HTTP 201
[Asserts]
jsonpath "$.id_album" isInteger
jsonpath "$.album_name" == "The Dark Side of the Moon"
jsonpath "$.album_dur" == 100
jsonpath "$.album_ph" == 1
jsonpath "$.release_date" == "1973-03-01"
jsonpath "$.artist_id_artist" == auxArtist_id
jsonpath "$.likes" == 0
[Captures]
albumcreado_id: jsonpath "$.id_album"

# LISTAR TODOS LOS ARTISTAS
GET http://localhost:8080/albums
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.length()" == 1
jsonpath "$[0].album_name" == "The Dark Side of the Moon"

#LISTAR ALBUM CREADO
GET http://localhost:8080/albums/{{albumcreado_id}}
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$.id_album" = {{albumcreado_id}}
jsonpath "$.album_name" == "The Dark Side of the Moon"
jsonpath "$.album_dur" == 100
jsonpath "$.album_ph" == 0
jsonpath "$.release_date" == "1973-03-01"
jsonpath "$.artist_id_artist" == auxArtist_id
jsonpath "$.likes" == 0

# UPDATE SOBRE ALBUM CREADO
PUT http://localhost:8080/albums/{{albumcreado_id}}
Content-Type: application/json
{
	"album_name": "The Dark Side of the Moon (Remastered)",
	"album_dur": 150 ,
	"album_ph": 2,
	"release_date": "2023-03-01",
	"artist_id_artist": {{auxArtist_id}},
	"likes": 100
}
HTTP 200
[Asserts]
jsonpath "$.id_album"= {{albumcreado_id}}
jsonpath "$.album_name" == "The Dark Side of the Moon (Remastered)"
jsonpath "$.album_dur" == 150
jsonpath "$.album_ph" == 2
jsonpath "$.release_date" == "2023-03-01"
jsonpath "$.artist_id_artist" == auxArtist_id
jsonpath "$.likes" == 100

# DELETE SOBRE ARTISTA CREADO
DELETE http://localhost:8080/albums/{{albumcreado_id}}
HTTP 204

#==========

#ERRORES 

# LISTAR ALBUM INEXISTENTE 

GET http://localhost:8080/albums/{{albumcreado_id}}
HTTP 404

# LISTAR ALBUM CON ID NO NUMERICO
GET http://localhost:8080/albums/abc
HTTP 400

# ACTUALIZACION DE ALBUM CON ID INEXISTENTE 
PUT http://localhost:8080/albums/1
HTTP 404

# ELIMINAR ALBUM INEXISTENTE
DELETE http://localhost:8080/albums/1
HTTP 404

# CREACION DE ALBUM CON PARAMETROS ERRONEOS

POST http://localhost:8080/albums
Content-Type: application/json
{
	"album_name": " ",
	"album_dur": 100 ,
	"album_ph": -1,
	"release_date": " ",
	"artist_id_artist": {{auxArtist_id}},
	"likes": -1
}
HTTP 400