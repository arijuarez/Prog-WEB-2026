## Correr tests con:
**make test**

## Descripción de la Página Web
JamBox es una red social simil a Letterbox, pero de canciones y álbumes musicales. 
Funcionará através de una API gratis de Spotify, donde los usuarios podrán buscar sus álbumes, canciones u artistas favoritos, con el fin de calificarlos, "likearlos" y compartir su opinión con otros usuarios mediante reviews.
Las entidades definidas y más información sobre la Base de Datos está definida en el archivo "persistencia.md".
Al correr **make test** se hará lo siguiente: se dará de baja cualquier contenedor activo, se generará el sqlc, se construirá la imagen, se activará el docker, se correrán los tests dentro de "tablas_test.go", y finalmente se desactivará el docker.