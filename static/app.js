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

/* ==========================================================
   NUEVO: lista de reseñas (GET /reviews) y borrado (DELETE)
   ========================================================== */
const lista = $("popular-reviews-list");

// sqlc serializa los campos opcionales como {Int32, Valid}; esto lo tolera
function valor(x) {
  if (x && typeof x === "object") return x.Valid === false ? null : x.Int32 ?? x.String ?? null;
  return x;
}

function nombreObjetivo(r) {
  const tipos = [["song", "Canción"], ["album", "Álbum"], ["artist", "Artista"]];
  for (const [tipo, nombre] of tipos) {
    const id = valor(r[`${tipo}_id`]);
    if (id) return `${nombre} #${id}`;
  }
  return "Reseña";
}

// Construye una tarjeta con textContent (nunca innerHTML) para evitar inyección de HTML
function crearTarjeta(r) {
  const rating = Number(valor(r.rating)) || 0;

  const card = document.createElement("article");
  card.className = "review-card";
  card.dataset.reviewId = r.id;

  const body = document.createElement("div");
  body.className = "review-body";

  const titulo = document.createElement("h3");
  titulo.className = "review-target";
  titulo.textContent = nombreObjetivo(r);

  const autor = document.createElement("div");
  autor.className = "review-author";
  const usuario = document.createElement("span");
  usuario.textContent = `Usuario #${valor(r.user_id)}`;
  const estrellasTxt = document.createElement("span");
  estrellasTxt.className = "stars";
  estrellasTxt.textContent = "★".repeat(rating) + "☆".repeat(5 - rating);
  autor.append(usuario, estrellasTxt);

  const texto = document.createElement("p");
  texto.className = "review-text";
  texto.textContent = valor(r.content) || "";

  const acciones = document.createElement("div");
  acciones.className = "review-actions";
  const btn = document.createElement("button");
  btn.type = "button";
  btn.className = "btn-delete";
  btn.textContent = "Eliminar";
  acciones.appendChild(btn);

  body.append(titulo, autor, texto, acciones);
  card.appendChild(body);
  return card;
}

async function cargarReviews() {
  try {
    const res = await fetch("/reviews");
    if (!res.ok) throw new Error(res.status);
    const reviews = (await res.json()) || []; // Go devuelve null si el slice está vacío

    lista.replaceChildren(); // vacía el contenedor
    if (reviews.length === 0) {
      const vacio = document.createElement("p");
      vacio.textContent = "Todavía no hay reseñas.";
      lista.appendChild(vacio);
      return;
    }
    reviews.forEach((r) => lista.appendChild(crearTarjeta(r)));
  } catch (_) {
    lista.textContent = "No se pudieron cargar las reseñas.";
  }
}

// Delegación de eventos: UN listener en el contenedor para todos los botones
lista.addEventListener("click", async (e) => {
  const btn = e.target.closest(".btn-delete");
  if (!btn) return;

  const id = btn.closest(".review-card").dataset.reviewId;
  try {
    const res = await fetch(`/reviews/${id}`, { method: "DELETE" });
    if (res.ok) cargarReviews();
    else alert("No se pudo eliminar la reseña.");
  } catch (_) {
    alert("No se pudo conectar con el servidor.");
  }
});

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
    $("signup-form").reset();
    cerrarModales();
  } catch (_) {
    mostrarError("signup-error", "No se pudo conectar con el servidor.");
  }
});

/* ---------- Login / menú / logout ----------
   Siguen igual, pero tu API todavía no tiene /login ni /logout.
   No los pide la consigna; sin esos endpoints mostrarán error. */
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
// Ya no hay sesión desde templates: el botón abre siempre el modal de reseña
$("btn-new-review").addEventListener("click", () => abrirModal($("modal-review")));

// Estrellas (sin cambios)
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

$("review-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  if (Number(valorPuntaje.value) === 0) return alert("Elegí una calificación de 1 a 5 estrellas.");

  try {
    const { ok } = await enviar("/reviews", {
      user_id: Number($("review-user-id").value),
      [`${$("review-type").value}_id`]: Number($("review-target-id").value), // song_id, album_id o artist_id
      rating: Number(valorPuntaje.value),
      content: $("review-content").value.trim(),
      spoiler: $("review-spoiler").checked,
    });

    if (!ok) return alert("No se pudo publicar la reseña. Probá de nuevo.");

    // Consigna: tras crear, volver a pedir la lista y refrescar la vista
    $("review-form").reset();
    valorPuntaje.value = 0;
    pintarEstrellas(0);
    cerrarModales();
    cargarReviews();
  } catch (_) {
    alert("No se pudo conectar con el servidor.");
  }
});

// Al cargar la página
cargarReviews();