"use strict";

const $ = (id) => document.getElementById(id);

// Cada página indica su endpoint en <body data-endpoint="/artists">
const endpoint = document.body.dataset.endpoint;
const lista = $("lista");
const form = $("form-entidad");

// sqlc serializa los opcionales como {Int32, Valid} o {String, Valid}
function valor(x) {
  if (x && typeof x === "object") {
    if (x.Valid === false) return null;
    return x.Int32 ?? x.Int64 ?? x.String ?? x.Time ?? JSON.stringify(x);
  }
  return x;
}

// Dibuja un <li> con todas las propiedades del objeto, excepto el id
function crearItem(obj, index) {
  const li = document.createElement("li");
  li.className = "item";
  // Busca el campo id sin importar prefijos (id, ID, id_review, id_song, etc.)
  const idKey = Object.keys(obj).find(k => k.toLowerCase() === "id" || k.toLowerCase().startsWith("id_"));
  li.dataset.id = idKey ? obj[idKey] : undefined;

  // Número de índice
  const num = document.createElement("span");
  num.className = "item-index";
  num.textContent = index + 1;

  const datos = document.createElement("div");
  datos.className = "item-datos";

  const keys = Object.entries(obj).filter(([k]) => k.toLowerCase() !== "id");
  keys.forEach(([clave, v], i) => {
    const texto = valor(v);
    if (texto === null || texto === "") return;

    const fila = document.createElement("span");
    const etiqueta = document.createElement("strong");
    etiqueta.textContent = `${clave}: `;

    if (i === 0) {
      // primera clave como "nombre" destacado
      fila.className = "item-name";
      fila.append(String(texto));
    } else {
      fila.append(etiqueta, String(texto));
    }
    datos.appendChild(fila);
  });

  const btn = document.createElement("button");
  btn.type = "button";
  btn.className = "btn-delete";
  btn.textContent = "Eliminar";

  li.append(num, datos, btn);
  return li;
}

// GET: pedir la lista y dibujarla
async function cargar() {
  try {
    const res = await fetch(endpoint);
    if (!res.ok) throw new Error(res.status);
    const items = (await res.json()) || []; // Go devuelve null si está vacía

    lista.replaceChildren();

    // actualizar contador si existe el elemento
    const contEl = document.getElementById("lista-count");
    if (contEl) contEl.textContent = items.length === 0 ? "" : `${items.length} elemento${items.length !== 1 ? "s" : ""}`;

    if (items.length === 0) {
      const vacio = document.createElement("li");
      vacio.textContent = "No hay elementos todavía.";
      lista.appendChild(vacio);
      return;
    }
    items.forEach((o, i) => lista.appendChild(crearItem(o, i)));
  } catch (_) {
    lista.textContent = "No se pudo cargar la lista.";
  }
}

// Arma el JSON a partir de los inputs del formulario (según su "name")
function leerFormulario() {
  const cuerpo = {};
  for (const el of form.elements) {
    if (!el.name) continue;
    if (el.type === "checkbox") cuerpo[el.name] = el.checked;
    else if (el.value.trim() === "") {
      // Los number opcionales se envían como null; el resto se omite
      if (el.type === "number" && !el.required) cuerpo[el.name] = null;
      // campos de texto vacíos no se envían (el backend los rechazaría igual)
    }
    else if (el.type === "number") cuerpo[el.name] = Number(el.value);
    else cuerpo[el.name] = el.value.trim();
  }
  return cuerpo;
}

// POST y luego volver a pedir la lista
form.addEventListener("submit", async (e) => {
  e.preventDefault();
  try {
    const res = await fetch(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(leerFormulario()),
    });
    if (!res.ok) return alert(`No se pudo crear (error ${res.status}).`);
    form.reset();
    cargar();
  } catch (_) {
    alert("No se pudo conectar con el servidor.");
  }
});

// DELETE con delegación de eventos: un solo listener en la lista
lista.addEventListener("click", async (e) => {
  const btn = e.target.closest(".btn-delete");
  if (!btn) return;
  const id = btn.closest(".item").dataset.id;
  try {
    const res = await fetch(`${endpoint}/${id}`, { method: "DELETE" });
    if (res.ok) cargar();
    else alert(`No se pudo eliminar (error ${res.status}).`);
  } catch (_) {
    alert("No se pudo conectar con el servidor.");
  }
});

cargar(); // al cargar la página