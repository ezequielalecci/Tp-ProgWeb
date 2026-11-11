/* ---------- Configuración de endpoints ---------- */
const ESTADOS = ["PENDIENTE","EN CURSO","FINALIZADO"];
const ENT = {
  equipos: {
    label:"Equipos", singular:"equipo", path:"/equipos", del:id=>`/equipos/${id}`,
    columns:[
      {k:"escudo_url",t:"Escudo",type:"img"},{k:"id",t:"ID",num:1},{k:"nombre",t:"Nombre"},
      {k:"formacion_actual",t:"Formación"},{k:"valoracion",t:"Valoración",num:1},
      {k:"posicion_tabla",t:"Posición",num:1},{k:"mejor_jugador",t:"Mejor jugador",type:"jugador"}],
    create:[
      {k:"nombre",l:"Nombre",req:1},{k:"formacion_actual",l:"Formación actual",ph:"4-3-3"},
      {k:"escudo_url",l:"URL del escudo",type:"url",opt:1},
      {k:"valoracion",l:"Valoración",type:"number",min:0},{k:"posicion_tabla",l:"Posición en la tabla",type:"number",min:1}],
    edit:{method:"PUT",path:id=>`/equipos/${id}`,fields:[
      {k:"nombre",l:"Nombre",req:1},{k:"formacion_actual",l:"Formación actual"},
      {k:"escudo_url",l:"URL del escudo",type:"url",opt:1},
      {k:"mejor_jugador",l:"ID del mejor jugador",type:"number",min:1,opt:1},
      {k:"valoracion",l:"Valoración",type:"number",min:0},{k:"posicion_tabla",l:"Posición en la tabla",type:"number",min:1}]}
  },
  jugadores: {
    label:"Jugadores", singular:"jugador", path:"/jugadores", del:id=>`/jugadores/${id}`,
    columns:[
      {k:"foto_url",t:"Foto",type:"img"},{k:"id",t:"ID",num:1},{k:"nombre",t:"Nombre"},
      {k:"equipo_id",t:"Equipo",type:"equipo"},{k:"posicion",t:"Pos."},{k:"fecha_nacimiento",t:"Nacimiento",type:"fecha"},
      {k:"altura",t:"Altura"},{k:"media_general",t:"Media",num:1},{k:"ritmo",t:"RIT",num:1},{k:"tiro",t:"TIR",num:1},
      {k:"pase",t:"PAS",num:1},{k:"regate",t:"REG",num:1},{k:"defensa",t:"DEF",num:1},{k:"fisico",t:"FIS",num:1},
      {k:"goles",t:"Goles",num:1},{k:"asistencias",t:"Asist.",num:1}],
    create:[
      {k:"equipo_id",l:"Equipo",type:"equipo",req:1},{k:"nombre",l:"Nombre",req:1},{k:"posicion",l:"Posición",ph:"DEL",req:1},
      {k:"fecha_nacimiento",l:"Fecha de nacimiento",type:"date",req:1},{k:"altura",l:"Altura",ph:"1.78 m"},
      {k:"media_general",l:"Media general (0–99)",type:"number",min:0,max:99},
      {k:"ritmo",l:"Ritmo",type:"number",min:0,max:99},{k:"tiro",l:"Tiro",type:"number",min:0,max:99},
      {k:"pase",l:"Pase",type:"number",min:0,max:99},{k:"regate",l:"Regate",type:"number",min:0,max:99},
      {k:"defensa",l:"Defensa",type:"number",min:0,max:99},{k:"fisico",l:"Físico",type:"number",min:0,max:99},
      {k:"foto_url",l:"URL de la foto",type:"url",opt:1}],
    edit:{method:"PUT",path:id=>`/jugadores/${id}/stats`,title:"Actualizar goles y asistencias",fields:[
      {k:"goles",l:"Goles",type:"number",min:0},{k:"asistencias",l:"Asistencias",type:"number",min:0}]}
  },
  partidos: {
    label:"Partidos", singular:"partido", path:"/partidos", del:id=>`/partidos/${id}`,
    columns:[
      {k:"id",t:"ID",num:1},{k:"fecha",t:"Fecha",type:"fecha"},{k:"equipo_local_id",t:"Local",type:"equipo"},
      {k:"_marcador",t:"Resultado",type:"marcador"},{k:"equipo_visitante_id",t:"Visitante",type:"equipo"},
      {k:"estado",t:"Estado",type:"estado"}],
    create:[
      {k:"equipo_local_id",l:"Equipo local",type:"equipo",req:1},{k:"equipo_visitante_id",l:"Equipo visitante",type:"equipo",req:1},
      {k:"fecha",l:"Fecha",type:"date",req:1},{k:"estado",l:"Estado",type:"estado",req:1}],
    edit:{method:"PUT",path:id=>`/partidos/${id}/resultado`,title:"Cargar resultado",fields:[
      {k:"goles_local",l:"Goles local",type:"number",min:0},{k:"goles_visitante",l:"Goles visitante",type:"number",min:0},
      {k:"estado",l:"Estado",type:"estado",req:1}]}
  }
};

/* ---------- Estado y utilidades ---------- */
const $ = s => document.querySelector(s);
const state = {tab:"equipos", data:{equipos:[],jugadores:[],partidos:[]}, sort:{}, loaded:{}};
const esc = v => String(v ?? "").replace(/[&<>"']/g, c => ({"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"}[c]));
const base = () => "/api";
const equipoNombre = id => (state.data.equipos.find(e => e.id === id) || {}).nombre || (id ? `#${id}` : "—");
const jugadorNombre = id => (state.data.jugadores.find(j => j.id === id) || {}).nombre || (id ? `#${id}` : "—");

function toast(msg, ok=true){
  const t = $("#toast"); t.textContent = msg; t.className = (ok ? "ok" : "err") + " show"; t.style.display = "block";
  clearTimeout(toast.h); toast.h = setTimeout(() => t.style.display = "none", 3500);
}

/* Convierte sql.NullInt32 / NullString / NullTime… ({Int32|String|Time…, Valid}) en valor simple o null */
function unwrap(v){
  if (Array.isArray(v)) return v.map(unwrap);
  if (v && typeof v === "object"){
    if ("Valid" in v){
      if (!v.Valid) return null;
      const k = ["Int16","Int32","Int64","Float64","String","Time","Bool","Byte"].find(k => k in v);
      return k ? v[k] : null;
    }
    return Object.fromEntries(Object.entries(v).map(([k,x]) => [k, unwrap(x)]));
  }
  return v;
}

async function api(method, path, body){
  const res = await fetch(base() + path, {
    method, headers: body ? {"Content-Type":"application/json"} : {},
    body: body ? JSON.stringify(body) : undefined
  });
  const text = await res.text();
  let json = null; try { json = text ? JSON.parse(text) : null; } catch {}
  if (!res.ok) throw new Error((json && (json.error || json.message)) || text || `Error ${res.status}`);
  return unwrap(json);
}

async function load(name){
  try{
    const json = await api("GET", ENT[name].path);
    state.data[name] = Array.isArray(json) ? json : (json && (json.data || json.items)) || [];
    state.loaded[name] = true;
  }catch(e){
    state.data[name] = [];
    toast(`No se pudo cargar ${ENT[name].label.toLowerCase()}: ${e.message}`, false);
  }
}
const loadAll = async () => { await Promise.all(Object.keys(ENT).map(load)); render(); };

/* ---------- Render de la tabla ---------- */
function cell(col, row){
  const v = row[col.k];
  switch(col.type){
    case "img": return v ? `<img src="${esc(v)}" alt="" loading="lazy" onerror="this.style.visibility='hidden'">` : "—";
    case "equipo": return esc(equipoNombre(v));
    case "jugador": return v ? esc(jugadorNombre(v)) : "—";
    case "fecha": return v ? esc(String(v).slice(0,10)) : "—";
    case "estado": return v ? `<span class="badge ${esc(v)}">${esc(v)}</span>` : "—";
    case "marcador": {
      const hay = row.goles_local != null && row.goles_visitante != null;
      return `<span class="score">${hay ? `${esc(row.goles_local)} – ${esc(row.goles_visitante)}` : "– –"}</span>`;
    }
    default: return v === null || v === undefined || v === "" ? "—" : esc(v);
  }
}
const sortVal = (col,row) => col.type==="equipo" ? equipoNombre(row[col.k]) : col.type==="marcador" ? (row.goles_local ?? -1) : row[col.k];

function render(){
  const ent = ENT[state.tab];
  let rows = [...state.data[state.tab]];
  const tf = $("#teamFilter").value;
  if (state.tab === "jugadores" && tf) rows = rows.filter(r => String(r.equipo_id) === tf);
  const s = state.sort[state.tab];
  if (s){
    const col = ent.columns.find(c => c.k === s.k);
    rows.sort((a,b) => {
      const x = sortVal(col,a), y = sortVal(col,b);
      const r = typeof x === "number" && typeof y === "number" ? x - y : String(x ?? "").localeCompare(String(y ?? ""), "es", {numeric:true});
      return s.dir * r;
    });
  }
  $("thead").innerHTML = "<tr>" + ent.columns.map(c =>
    `<th data-k="${c.k}" class="${c.num?"num":""}" title="Ordenar">${esc(c.t)}${s && s.k===c.k ? (s.dir>0?" ▲":" ▼") : ""}</th>`).join("") + "<th></th></tr>";
  $("tbody").innerHTML = rows.length
    ? rows.map(r => "<tr>" + ent.columns.map(c => `<td class="${c.num?"num":""}">${cell(c,r)}</td>`).join("") +
        `<td class="acts"><button class="btn small" data-edit="${esc(r.id)}">${esc(ent.edit.title ? ent.edit.title.split(" ")[0] : "Editar")}</button>` +
        (state.tab === "equipos" ? `<button class="btn small" data-info="${esc(r.id)}">Consultas</button>` : "") +
        `<button class="btn small" data-del="${esc(r.id)}">Eliminar</button></td></tr>`).join("")
    : `<tr><td class="empty" colspan="${ent.columns.length+1}">No hay ${ent.label.toLowerCase()} para mostrar. Cargá el primero con el botón de arriba.</td></tr>`;
  $("#count").textContent = `${rows.length} de ${state.data[state.tab].length}`;
  const sel = $("#teamFilter"), prev = sel.value;
  sel.style.display = state.tab === "jugadores" ? "" : "none";
  sel.innerHTML = '<option value="">Todos los equipos</option>' + state.data.equipos.map(e => `<option value="${esc(e.id)}">${esc(e.nombre)}</option>`).join("");
  sel.value = prev;
  $("#add").textContent = `Agregar ${ent.singular}`;
  document.querySelectorAll("#tabs button").forEach(b => b.setAttribute("aria-selected", b.dataset.tab === state.tab));
}

/* ---------- Formularios ---------- */
function inputFor(f, value){
  const v = value ?? "";
  const common = `name="${f.k}" id="f_${f.k}" ${f.req ? "required" : ""}`;
  if (f.type === "equipo")
    return `<select ${common}><option value="">Elegí un equipo…</option>${state.data.equipos.map(e =>
      `<option value="${esc(e.id)}" ${String(v)===String(e.id)?"selected":""}>${esc(e.nombre)}</option>`).join("")}</select>`;
  if (f.type === "estado")
    return `<select ${common}>${ESTADOS.map(s => `<option ${v===s?"selected":""}>${s}</option>`).join("")}</select>`;
  const t = f.type || "text";
  const extra = (f.min!==undefined?` min="${f.min}"`:"") + (f.max!==undefined?` max="${f.max}"`:"") + (f.ph?` placeholder="${esc(f.ph)}"`:"");
  const val = t === "date" ? String(v).slice(0,10) : v;
  return `<input type="${t}" ${common}${extra} value="${esc(val)}">`;
}

function openForm(mode, row){
  const ent = ENT[state.tab];
  const isEdit = mode === "edit";
  const fields = isEdit ? ent.edit.fields : ent.create;
  const title = isEdit ? (ent.edit.title || `Editar ${ent.singular}`) + ` #${row.id}` : `Nuevo ${ent.singular}`;
  $("#form").innerHTML = `<h2>${esc(title)}</h2><div class="grid">` +
    fields.map(f => `<div class="field"><label for="f_${f.k}">${esc(f.l)}${f.opt?" (opcional)":""}</label>${inputFor(f, isEdit ? row[f.k] : "")}</div>`).join("") +
    `</div><div class="actions"><button type="button" class="btn" id="cancel">Cancelar</button><button class="btn primary" type="submit">${isEdit?"Guardar cambios":"Crear"}</button></div>`;
  $("#cancel").onclick = () => $("#dlg").close();
  $("#form").onsubmit = async ev => {
    ev.preventDefault();
    const body = {};
    for (const f of fields){
      const raw = $("#form").elements[f.k].value.trim();
      if (f.type === "number" || f.type === "equipo") body[f.k] = raw === "" ? (f.opt ? null : 0) : Number(raw);
      else body[f.k] = raw === "" && f.opt ? null : raw;
    }
    // Validaciones equivalentes a las del backend
    if (state.tab === "partidos" && !isEdit && body.equipo_local_id === body.equipo_visitante_id)
      return toast("El equipo local y el visitante no pueden ser el mismo.", false);
    try{
      if (isEdit) await api(ent.edit.method, ent.edit.path(row.id), body);
      else await api("POST", ent.path, body);
      $("#dlg").close();
      toast(isEdit ? "Cambios guardados" : `${ent.singular[0].toUpperCase()+ent.singular.slice(1)} creado`);
      await loadAll();
    }catch(e){ toast(e.message, false); }
  };
  $("#dlg").showModal();
}

/* ---------- Eventos ---------- */
$("#tabs").innerHTML = Object.entries(ENT).map(([k,e]) => `<button role="tab" data-tab="${k}">${e.label}</button>`).join("");
$("#tabs").onclick = e => { const b = e.target.closest("button"); if (b){ state.tab = b.dataset.tab; $("#teamFilter").value=""; render(); } };
$("#add").onclick = () => openForm("create");
$("thead").onclick = e => {
  const th = e.target.closest("th[data-k]"); if (!th) return;
  const cur = state.sort[state.tab];
  state.sort[state.tab] = cur && cur.k === th.dataset.k ? {k:cur.k, dir:-cur.dir} : {k:th.dataset.k, dir:1};
  render();
};
$("#teamFilter").onchange = render;
$("tbody").onclick = async e => {
  const b = e.target.closest("button"); if (!b) return;
  const ent = ENT[state.tab];
  const find = id => state.data[state.tab].find(r => String(r.id) === id);
  if (b.dataset.edit){ const row = find(b.dataset.edit); if (row) openForm("edit", row); }
  else if (b.dataset.info){ const row = find(b.dataset.info); if (row) openInfo(row); }
  else if (b.dataset.del){
    if (!confirm(`¿Eliminar ${ent.singular} #${b.dataset.del}? Esta acción no se puede deshacer.`)) return;
    try{ await api("DELETE", ent.del(b.dataset.del)); toast(`${ent.singular[0].toUpperCase()+ent.singular.slice(1)} eliminado`); await loadAll(); }
    catch(err){ toast(err.message, false); }
  }
};

/* Consultas por equipo: mejor jugador y últimos partidos */
const CONSULTAS = [
  ["Mejor jugador", id => `/equipos/${id}/mejor-jugador`],
  ["Último partido", id => `/equipos/${id}/partidos/ultimo`],
  ["Último finalizado", id => `/equipos/${id}/partidos/ultimo-finalizado`],
  ["Último con detalle", id => `/equipos/${id}/partidos/ultimo-detalle`]
];
function kv(o){
  if (o === null || typeof o !== "object") return `<p>${esc(o)}</p>`;
  return `<dl class="kv">` + Object.entries(o).map(([k,v]) => {
    const val = v && typeof v === "object" ? `<code>${esc(JSON.stringify(v))}</code>` : esc(v ?? "—");
    return `<dt>${esc(k.replace(/_/g," "))}</dt><dd>${val}</dd>`;
  }).join("") + `</dl>`;
}
function openInfo(eq){
  $("#form").innerHTML = `<h2>${esc(eq.nombre)}</h2><div class="toolbar">` +
    CONSULTAS.map((c,i) => `<button type="button" class="btn small" data-q="${i}">${c[0]}</button>`).join("") +
    `</div><div id="res" class="count">Elegí una consulta.</div><div class="actions"><button type="button" class="btn" id="cancel">Cerrar</button></div>`;
  $("#cancel").onclick = () => $("#dlg").close();
  $("#form").onsubmit = ev => ev.preventDefault();
  $("#form").onclick = async ev => {
    const q = ev.target.closest("[data-q]"); if (!q) return;
    $("#res").textContent = "Cargando…";
    try{ $("#res").innerHTML = kv(await api("GET", CONSULTAS[q.dataset.q][1](eq.id))); }
    catch(err){ $("#res").textContent = err.message; }
  };
  $("#dlg").showModal();
}

render();
loadAll();