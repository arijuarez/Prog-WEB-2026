# Interfaces dinámicas con HTMX (TP6)

En esta etapa la aplicación del TP5 pasó a usar **HTMX** para crear y borrar entidades **sin recargar la página**. El servidor sigue generando el HTML con `templ`, pero ahora devuelve solo el fragmento que cambió (una lista, un formulario, un contador) y HTMX lo inserta en el lugar correcto. No se escribió JavaScript propio.

Las entidades que se pueden crear y borrar son: usuarios, canciones, álbumes, reseñas y artistas.

---

## 1. Integración de HTMX

**Archivos:** `views/layout.templ`, `views/layout_entities.templ`

Se agregó el script de HTMX en el `<head>` de los dos layouts:

```html
<meta name="htmx-config" content='{"noSwap":[204,304,"4xx","5xx"]}'>
<script src="https://cdn.jsdelivr.net/npm/htmx.org@4.0.0/dist/htmx.min.js"></script>
```

Se usa **HTMX 4.0** (la versión vista en la teoría), no la 1.9.10 del enunciado. La línea `htmx-config` se explica en la sección 5.

---

## 2. Creación sin recarga de página

### 2.1 Listas con `id`

**Archivo:** `views/entity_list.templ`

Cada lista se envuelve en un `<section>` con un `id` único, que es el destino de la respuesta:

| Entidad  | Componente    | `id` del contenedor |
|----------|---------------|---------------------|
| Usuarios | `UserList`    | `#user-list`        |
| Canciones| `SongList`    | `#song-list`        |
| Álbumes  | `AlbumList`   | `#album-list`       |
| Reseñas  | `ReviewList`  | `#review-list`      |
| Artistas | `ArtistsList` | `#artist-list`      |

Los componentes se renombraron con **mayúscula inicial** (antes `userList`, etc.) para que sean **exportados** y el paquete `Handlers` pueda renderizarlos. En Go, un nombre en minúscula es privado del paquete.

### 2.2 Formularios con atributos HTMX

**Archivo:** `views/entity_form.templ`

Se reemplazó `method="POST" action="/users"` por:

```html
<form id="form-entidad" hx-post="/users" hx-target="#user-list" hx-swap="outerHTML">
```

- `hx-post`: envía el formulario por AJAX a la misma ruta que antes.
- `hx-target`: indica qué elemento se actualiza con la respuesta. Lleva `#` porque es un selector CSS de `id`.
- `hx-swap="outerHTML"`: reemplaza el `<section>` completo, no solo su contenido.

### 2.3 Handler `POST /<entidades>`

**Archivo:** `Handlers/Handlers.go`, función `HandlePOST`

Se **eliminó la redirección** (`http.Redirect`). Ahora, después de crear la entidad, el handler:

1. vuelve a consultar la lista completa (`queries.ListUsers(ctx)`, etc.);
2. renderiza **solo** el componente de la lista (`views.UserList(list)`);
3. renderiza además el formulario vacío como fragmento OOB (ver 2.4).

```go
list, err := queries.ListUsers(ctx)
if err != nil { ... }
views.UserList(list).Render(ctx, w)
views.FormUser(true).Render(ctx, w) // form vacío (OOB) → limpia el formulario
```

### 2.4 Limpieza del formulario (Out-of-Band swap)

Sin esto, después de crear una entidad los campos seguían con los datos cargados. Se resolvió con la técnica **OOB** del ejercicio 3:

- Los formularios se exportaron (`FormUser`, `FormSong`, ...) y reciben un parámetro `oob bool`:

  ```
  templ FormUser(oob bool) {
      <form id="form-entidad" hx-post="/users" hx-target="#user-list" hx-swap="outerHTML"
          if oob {
              hx-swap-oob="true"
          }
      >
  ```

- En la página se usan con `@FormUser(false)` (`views/index.templ`).
- En la respuesta del POST se usan con `FormUser(true)`. HTMX ve `hx-swap-oob="true"`, busca en la página el elemento con el mismo `id` (`form-entidad`) y lo reemplaza por el formulario vacío. El resto de la respuesta (la lista) va al `hx-target`.

Se eligió un parámetro en lugar de un componente envoltorio (`ProductFormOOB`) para no duplicar `id`s ni repetir el formulario.

---

## 3. Borrado sin recarga de página

### 3.1 Botón "Eliminar"

**Archivo:** `views/entity_list.templ`

El `<form method="POST" action="/users/{id}/delete">` de cada elemento se reemplazó por:

```html
<button type="button"
    hx-delete={ fmt.Sprintf("/users/%d", user.IDUser) }
    hx-target="closest li"
    hx-swap="outerHTML"
    hx-confirm="¿Seguro que querés eliminar este usuario?">
    Eliminar
</button>
```

- `hx-delete`: envía una petición HTTP `DELETE` a `/<entidades>/{id}`.
- `hx-target="closest li"`: el enunciado usa `closest tr` porque lista en una tabla; acá cada entidad está en un `<li>`.
- `hx-swap="outerHTML"`: el `<li>` se reemplaza por la respuesta. Como la respuesta no trae contenido para el target, el `<li>` desaparece.
- `hx-confirm`: pide confirmación antes de enviar.

### 3.2 Ruta `DELETE` en el router

**Archivo:** `main.go`

El router solo aceptaba GET y POST (cualquier otro método devolvía 405). Se agregó:

```go
case http.MethodDelete:
    Handlers.HandleDELETE(w, r, queries)
```

### 3.3 Handler `DELETE /<entidades>/{id}`

**Archivo:** `Handlers/Handlers.go`, función `HandleDELETE`

1. Toma la entidad y el ID de la URL (`/users/5` → `["users", "5"]`).
2. Llama al método de sqlc correspondiente (`DeleteUser`, `DeleteSong`, `DeleteAlbum`, `DeleteReview`, `DeleteArtist`). Las queries ya existían, no hubo que regenerar sqlc.
3. Responde **200 OK**. Se usa 200 y no 204 porque HTMX **no hace el swap** con 204, y el `<li>` quedaría visible.

Si la entidad está referenciada por otra (por ejemplo, un álbum con canciones), la base rechaza el borrado y el handler responde **409 Conflict**.

### 3.4 Contador de elementos (extra)

Al borrar solo se reemplaza el `<li>`, así que el contador de la cabecera de la lista no cambiaba. Se resolvió también con OOB:

- El contador pasó a ser un componente `ListCount(n int, oob bool)` en `views/entity_list.templ`, usado en las cinco listas con `oob = false`.
- Al final de `HandleDELETE` se cuentan los elementos restantes (`contarEntidades`) y se devuelve `ListCount(n, true)`.

La respuesta del DELETE queda así:

```html
<span class="entity-count" id="lista-count" hx-swap-oob="true">4</span>
```

HTMX extrae el `<span>` OOB y actualiza el número. Lo que queda para el `<li>` es vacío, así que la fila desaparece igual.

---

## 4. Verificación

1. `make run-app` (genera sqlc y templ, levanta la base y el servidor).
2. Entrar a `/users`, `/songs`, `/albums`, `/reviews` o `/artists`.
3. **Crear:** completar el formulario y enviarlo. La lista se actualiza, el contador sube y el formulario queda vacío, sin recargar la página.
4. **Borrar:** apretar "Eliminar" y confirmar. El elemento desaparece y el contador baja.
5. En DevTools → Network se ven las peticiones `POST` y `DELETE`, y la respuesta es solo el fragmento HTML, no la página completa.

---

## 5. Correcciones adicionales

### 5.1 Errores del servidor que rompían la página

En **HTMX 4**, a diferencia de HTMX 2, las respuestas de error (4xx y 5xx) **se insertan en la página** por defecto. Por ejemplo, si un formulario tenía un dato inválido, el texto "Cantidad de amigos inválida" reemplazaba toda la lista. Como el `#user-list` desaparecía, los envíos siguientes ya no tenían dónde insertarse.

Solución, en los dos layouts:

- `<meta name="htmx-config" content='{"noSwap":[204,304,"4xx","5xx"]}'>`: las respuestas de error **no** se insertan, como en HTMX 2.
- `hx-on:htmx:response:error="alert(...)"` en el `<body>`: ante cualquier error de una petición HTMX se muestra un aviso al usuario. El evento sube (*bubbling*) desde el formulario o botón hasta el `<body>`.

### 5.2 Otras

- **Mensajes de error copiados:** en el POST de canciones, álbumes, reseñas y artistas, el error de listado decía "Error al obtener usuarios".
- **Errores internos expuestos:** varias respuestas incluían el mensaje crudo de PostgreSQL. Se agregó `errorServidor`, que registra el error real con `log.Printf` y responde un mensaje genérico.
- **CSS inexistente:** `layout_entities.templ` cargaba `static/entidades.css`, que no existe (daba 404). Se quitó el enlace.
- **Rutas relativas:** `href="static/styles.css"` pasó a `href="/static/styles.css"` para que funcione desde cualquier ruta.
- **Fechas:** la fecha de las reseñas y de lanzamiento de los álbumes se mostraba como `2026-11-20 00:00:00 +0000 UTC`. Ahora se muestra como `20/11/2026` (`.Format("02/01/2006")`).

---

## Resumen de archivos modificados

| Archivo | Cambio |
|---|---|
| `views/layout.templ`, `views/layout_entities.templ` | Script de HTMX, `htmx-config`, aviso de errores, rutas de CSS |
| `views/entity_form.templ` | `hx-post`/`hx-target`/`hx-swap`, formularios exportados con versión OOB |
| `views/entity_list.templ` | `id` en las listas, componentes exportados, botón `hx-delete`, `ListCount`, formato de fechas |
| `views/index.templ` | Llamadas a los componentes renombrados |
| `main.go` | Ruta para el método `DELETE` |
| `Handlers/Handlers.go` | POST sin redirección (lista + form OOB), `HandleDELETE` (200 + contador OOB), `errorServidor`, mensajes corregidos |
