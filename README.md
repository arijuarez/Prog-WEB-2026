## Correr aplicación con:
### **make run-app**  
### Al finalizar, **no olvidar** de ejecutar **make delete** para borrar el contenedor. 

### **IMPORTANTE! PARA USAR LA PÁGINA CON SUS POST Y GET:**
Los botones funcionales son las columnas de arriba: artistas, canciones, usuarios, reseñas y álbumes. El resto, por el momento, es sólo ilustrativo, incluidos los botones de login y logout.
En las columnas declaradas se podrán crear y eliminar entidades.
No se podrán borrar entidades que estén siendo referenciadas por otra.

## Descripción de la Página Web
JamBox es una red social simil a Letterbox, pero de canciones y álbumes musicales. 
Funcionará através de una API gratis de Spotify, donde los usuarios podrán buscar sus álbumes, canciones u artistas favoritos, con el fin de calificarlos, "likearlos" y compartir su opinión con otros usuarios mediante reviews.
Las entidades definidas y más información sobre la Base de Datos está definida en el archivo "persistencia.md".
La creación y el borrado de entidades sin recargar la página (HTMX, TP6) están documentados en el archivo "htmx.md".
Al correr **make run-app** se hará lo siguiente: se dará de baja cualquier contenedor activo, se generará el sqlc, se generará el código go de templ, se construirá la imagen, se activará el docker, y se correrá el servidor go. Para desactivar el docker, se debe ejecutar **make delete**, ya que si se ponía todo en un mismo comando del makefile no se iba a poder correr la página.

## Documentación extra del trabajo
Se puede ver más sobre las decisiones de diseño y particularidades del pasaje al HTMX en el markdown llamado **'htmx.md'**.
