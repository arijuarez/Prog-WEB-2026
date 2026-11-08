/* ==========================================================
   JamBox: main.js
   Go dibuja las páginas (templates). Este archivo solo:
     1. abre y cierra los modales
     2. envía los formularios a los endpoints de Go
     3. maneja el selector de estrellas

   Endpoints: POST /users (registro), POST /login, POST /logout,
              POST /reviews
   ========================================================== */

"use strict";

const $ = (id) => document.getElementById(id);

// Envía JSON a Go y devuelve { ok, status }
async function enviar(url, cuerpo) {
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "same-origin",
    body: cuerpo ? JSON.stringify(cuerpo) : undefined,
  });
  return { ok: res.ok, status: res.status };
}

function mostrarError(id, mensaje) {
  const caja = $(id);
  caja.textContent = mensaje;
  caja.hidden = false;
}

function ocultarError(id) {
  $(id).hidden = true;
}

/* ---------- Modales ---------- */
const overlay = $("modal-overlay");
const modales = [$("modal-login"), $("modal-signup"), $("modal-review")];

function abrirModal(modal) {
  modales.forEach((m) => (m.hidden = true));
  modal.hidden = false;
  overlay.hidden = false;
  document.body.style.overflow = "hidden";
}

function cerrarModales() {
  overlay.hidden = true;
  modales.forEach((m) => (m.hidden = true));
  document.body.style.overflow = "";
}

// Los botones de login y registro solo existen si no hay sesión
const abrirLogin = () => abrirModal($("modal-login"));
const abrirRegistro = () => abrirModal($("modal-signup"));

["btn-login"].forEach((id) => $(id) && $(id).addEventListener("click", abrirLogin));
["btn-signup", "hero-cta", "bottom-cta-btn"].forEach((id) => $(id) && $(id).addEventListener("click", abrirRegistro));

["modal-login-close", "modal-signup-close", "modal-review-close"].forEach((id) =>
  $(id).addEventListener("click", cerrarModales)
);

overlay.addEventListener("click", (e) => {
  if (e.target === overlay) cerrarModales();
});

document.addEventListener("keydown", (e) => {
  if (e.key === "Escape") cerrarModales();
});

/* ---------- Registro: POST /users ---------- */
$("signup-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  ocultarError("signup-error");

  const password = $("signup-password").value;
  if (password.length < 8) {
    return mostrarError("signup-error", "La contraseña debe tener al menos 8 caracteres.");
  }

  try {
    const { ok, status } = await enviar("/users", {
      username: $("signup-username").value.trim(),
      email: $("signup-email").value.trim(),
      password,
    });

    if (!ok) {
      return mostrarError(
        "signup-error",
        status === 409 ? "Ese usuario o email ya está registrado." : "No se pudo crear la cuenta."
      );
    }
    location.reload(); // Go vuelve a dibujar la página con la sesión iniciada
  } catch (_) {
    mostrarError("signup-error", "No se pudo conectar con el servidor.");
  }
});

/* ---------- Login: POST /login ---------- */
$("login-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  ocultarError("login-error");

  try {
    const { ok } = await enviar("/login", {
      username: $("login-username").value.trim(),
      password: $("login-password").value,
    });

    if (!ok) return mostrarError("login-error", "Usuario o contraseña incorrectos.");
    location.reload();
  } catch (_) {
    mostrarError("login-error", "No se pudo conectar con el servidor.");
  }
});

/* ---------- Menú de usuario y logout: POST /logout ---------- */
if ($("user-menu-toggle")) {
  $("user-menu-toggle").addEventListener("click", (e) => {
    e.stopPropagation();
    $("user-menu").hidden = !$("user-menu").hidden;
  });

  document.addEventListener("click", (e) => {
    if (!$("auth-user").contains(e.target)) $("user-menu").hidden = true;
  });

  $("btn-logout").addEventListener("click", async () => {
    await enviar("/logout");
    location.reload();
  });
}

/* ---------- Nueva reseña: POST /reviews ---------- */
// Go imprime el id del usuario en <body data-user-id="...">; vacío = sin sesión
const userId = document.body.dataset.userId;

$("btn-new-review").addEventListener("click", () => {
  abrirModal(userId ? $("modal-review") : $("modal-login"));
});

// Estrellas
const estrellas = Array.from($("review-rating").querySelectorAll("button"));
const valorPuntaje = $("review-rating-value");

function pintarEstrellas(hasta) {
  estrellas.forEach((s) => s.classList.toggle("active", Number(s.dataset.value) <= hasta));
}

estrellas.forEach((s) => {
  s.addEventListener("click", () => {
    valorPuntaje.value = s.dataset.value;
    pintarEstrellas(Number(s.dataset.value));
  });
  s.addEventListener("mouseenter", () => pintarEstrellas(Number(s.dataset.value)));
});

$("review-rating").addEventListener("mouseleave", () => pintarEstrellas(Number(valorPuntaje.value)));

// El campo de búsqueda usa un <datalist> que dibuja Go con las opciones
function buscarObjetivo(texto) {
  return Array.from($("review-options").options).find((o) => o.value === texto);
}

$("review-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const opcion = buscarObjetivo($("review-search").value);
  if (!opcion) return alert("Elegí una canción, álbum o artista de la lista.");
  if (Number(valorPuntaje.value) === 0) return alert("Elegí una calificación de 1 a 5 estrellas.");

  // data-type es "song", "album" o "artist" -> song_id, album_id, artist_id
  try {
    const { ok } = await enviar("/reviews", {
      user_id: Number(userId),
      [`${opcion.dataset.type}_id`]: Number(opcion.dataset.id),
      rating: Number(valorPuntaje.value),
      content: $("review-content").value.trim(),
      spoiler: $("review-spoiler").checked,
    });

    if (!ok) return alert("No se pudo publicar la reseña. Probá de nuevo.");
    location.reload();
  } catch (_) {
    alert("No se pudo conectar con el servidor.");
  }
});