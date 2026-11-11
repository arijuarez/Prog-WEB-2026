## Correr tests con:
**make run-app**

### Forma alternativa para que los datos perduren:
**make exec**
Hace que no se borren los datos del docker, manteniendo los datos ingresados al reiniciar la página. Al terminar se probarla se deberá correr: **docker compose down -v** para eliminar todo.

## Descripción de la Página Web
JamBox es una red social simil a Letterbox, pero de canciones y álbumes musicales. 
Funcionará através de una API gratis de Spotify, donde los usuarios podrán buscar sus álbumes, canciones u artistas favoritos, con el fin de calificarlos, "likearlos" y compartir su opinión con otros usuarios mediante reviews.
Las entidades definidas y más información sobre la Base de Datos está definida en el archivo "persistencia.md".
Al correr **make run-app** se hará lo siguiente: se dará de baja cualquier contenedor activo, se generará el sqlc, se generará el código go de templ, se construirá la imagen, se activará el docker, se correrá el servidor go, y finalmente se desactivará el docker, borrando todos los datos del contenedor.
