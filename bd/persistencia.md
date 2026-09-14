# Capa de Persistencia (Base de Datos)

El proyecto utiliza **PostgreSQL** como motor de base de datos relacional, **sqlc** para la generación automática de código en Go a partir de consultas SQL seguras y tipadas, y Docker para crear la imagen compuesta por los archivos sql y los tests hechos en Go.

##  Las Tecnologías:
* **Motor de BD:** PostgreSQL
* **Generador de Código SQL:** [sqlc](https://sqlc.dev/)
* **Driver de Go:** `github.com/jackc/pgx`

---

## Base De Datos:

La Base 

##  Estructura de Directorios

El código relacionado con la base de datos se encuentra aislado en la carpeta `db/`:

```text
tu_proyecto/
├── db/                    
│   ├── schema.sql         # Definición de tablas (DDL: CREATE, ALTER)
│   ├── queries.sql        # Consultas a la base de datos (DML: SELECT, INSERT, UPDATE)
│   └── sqlc/              # ⚠️ CÓDIGO AUTOGENERADO. No editar manualmente.
│       ├── db.go          # Manejador de la conexión a la base de datos
│       ├── models.go      # Structs de Go (Mapeo de las tablas)
│       └── queries.sql.go # Métodos para ejecutar las consultas SQL
├── sqlc.yaml              # Archivo de configuración de sqlc
└── main.go