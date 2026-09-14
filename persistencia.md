# Capa de Persistencia (Base de Datos)

El proyecto utiliza **PostgreSQL** como motor de base de datos relacional, **sqlc** para la generación automática de código en Go a partir de las queries SQL seguras y tipadas, y Docker para crear la imagen compuesta por la imagen sql y la imagen de los tests hechos en Go.

##  Las Tecnologías:
* **Motor de BD:** PostgreSQL
* **Generador de Código SQL:** [sqlc](https://sqlc.dev/)
* **Driver de Go:** `github.com/jackc/pgx`

---

## Base De Datos:

La base de datos de JamBox cuenta con 6 entidades(tablas) y 7 relaciones entre ellas:
ESTRUCTURA DE LA BASE DE DATOS MUSICAL

1. ARTIST TABLE (Artista)
- id_artist: serial (Primary Key)
- artist_name: varchar
- photo: int (momentaneamente hasta que resolvamos poner las imagenes en base64)
- artist_bio: text

2. ALBUM TABLE(Álbum)
- id_album: serial (Primary Key)
- album_name: int
- album_dur: int
- album_ph: int (momentaneamente hasta que resolvamos poner las imagenes en base64)
- release_date: date
- Artist_id_artist: int (Foreign Key -> Artist)
- Likes: int

3. SONG TABLE (Canción)
- id_song: serial (Primary Key)
- song_name: varchar
- song_dur: int
- likes: int
- Album_id_album: int (Foreign Key -> Album)

4. USER TABLE (Usuario)
- id_user: serial (Primary Key)
- user_name: varchar
- amigos: int
- photo: int (momentaneamente hasta que resolvamos poner las imagenes en base64)

5. ARTIST_SONG TABLE (Relación Muchos a Muchos: Artista y Canción) -> 
- Artist_id_artist: int (Primary Key, Foreign Key -> Artist)
- Song_id_song: int (Primary Key, Foreign Key -> Song)
- Roll: bool

6. REVIEW TABLE (Reseña)
- id_review: serial (Primary Key)
- review_text: text
- review_rate: float
- review_date: date
- User_id_user: int (Foreign Key -> User)
- Song_id_song: int (Nullable Foreign Key -> Song)
- Album_id_album: int (Nullable Foreign Key -> Album)
- Likes: int

RESUMEN DE RELACIONES:
- Un artista tiene muchos álbumes (1:N). (relacion artist->album)
- Un álbum tiene muchas canciones (1:N). (relacion album->song)
- Un artista puede participar de muchas canciones como creador o acompañante (1:N) (relacion artist->artist_song)
- Una cancion puede tener a muchos artistas involucrados en ella (1:N) (relacion song->artist_song)
- Un usuario puede hacer muchas reseñas o no tener ninguna reseña (0:N)
- Una album puede no tener o tener muchas reviews (0:N)(relacion album->review) 
- Una cancion puede no tener o tener muchas reviews (0:N)(relacion song->review) 

![SQL model](./bd/DataModelerSQL.png)
## 2. Estructura de Directorio
DataModelerSQL.png

```
tp2/
├── bd/
│   ├── sqlc/               # Código Go generado por sqlc
│   │   ├── db.go          
│   │   ├── models.go       # Estructuras de las entidades 
│   │   └── queries.sql.go  # Métodos para ejecutar consultas SQL
│   ├── queries.sql         # Sentencias SQL con anotaciones para sqlc (-- name: ...)
│   └── schema.sql          # Definición de las tablas y relaciones
├── docker-compose.yml      # Contenedor PostgreSQL
├── Dockerfile              # Imagen para pruebas e integración continua
├── Makefile                # Archivo para ejecutar los tests y activar el docker etc
├── sqlc.yaml               # Configuración del generador sqlc
├── persistencia.md         # Documentación de la capa de datos
└── tablas_test.go          # Unit tests para las tablas
