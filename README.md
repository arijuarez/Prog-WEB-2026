## Correr tests con:
**make test-hurl**

## Descripción de la Página Web
JamBox es una red social simil a Letterbox, pero de canciones y álbumes musicales. 
Funcionará através de una API gratis de Spotify, donde los usuarios podrán buscar sus álbumes, canciones u artistas favoritos, con el fin de calificarlos, "likearlos" y compartir su opinión con otros usuarios mediante reviews.
Las entidades definidas y más información sobre la Base de Datos está definida en el archivo "persistencia.md".
## Pruebas de Integración con Hurl

El proyecto utiliza **[Hurl](https://hurl.dev/)** para ejecutar pruebas de integración automatizadas de extremo a extremo sobre la API REST.

Los archivos `.hurl` definen secuencias de peticiones HTTP en texto plano que simulan el comportamiento real de un cliente interactuando con la API. En cada prueba se evalúa:
- **CRUD Completo:** Creación, lectura, actualización y eliminación de recursos (usuarios, artistas, álbumes, canciones y reseñas).
- **Aserciones en Respuestas:** Validación estricta de códigos de estado HTTP (`201 Created`, `200 OK`, `400 Bad Request`, `404 Not Found`), cabeceras (`Content-Type`) y estructuras del cuerpo JSON mediante expresiones `jsonpath`.
- **Encadenamiento de Peticiones:** Captura dinámica de variables (como los IDs generados tras un `POST`) para utilizarlas en consultas, ediciones (`PUT`) o borrados (`DELETE`) posteriores.
- **Manejo de Errores y Edge Cases:** Verificación de comportamiento ante inputs inválidos, campos faltantes o recursos inexistentes.
