package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"tpweb/internal/dto"
	"tpweb/internal/service"
	"tpweb/views"
)

type PageHandler struct {
	equipos   *service.EquipoService
	jugadores *service.JugadorService
	partidos  *service.PartidoService
}

func NewPageHandler(e *service.EquipoService, j *service.JugadorService, p *service.PartidoService) *PageHandler {
	return &PageHandler{equipos: e, jugadores: j, partidos: p}
}

func (h *PageHandler) RegisterRoutes(mux *http.ServeMux) {
	// {$} hace que el patrón coincida solo con "/" y no con cualquier ruta.
	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("GET /equipos/lista", h.ListarEquiposPartial)
	mux.HandleFunc("POST /equipos", h.CrearEquipo)
	mux.HandleFunc("GET /equipos/{id}/edit", h.EditarEquipoForm)
	mux.HandleFunc("PUT /equipos/{id}", h.ActualizarEquipo)
	mux.HandleFunc("DELETE /equipos/{id}", h.EliminarEquipo)
	mux.HandleFunc("GET /jugadores/lista", h.ListarJugadoresPartial)
	mux.HandleFunc("POST /jugadores", h.CrearJugador)
	mux.HandleFunc("GET /jugadores/{id}/edit", h.EditarJugadorForm)
	mux.HandleFunc("PUT /jugadores/{id}", h.ActualizarJugador)
	mux.HandleFunc("DELETE /jugadores/{id}", h.EliminarJugador)
	mux.HandleFunc("GET /partidos/lista", h.ListarPartidosPartial)
	mux.HandleFunc("POST /partidos", h.CrearPartido)
	mux.HandleFunc("GET /partidos/{id}/edit", h.EditarPartidoForm)
	mux.HandleFunc("PUT /partidos/{id}", h.ActualizarPartido)
	mux.HandleFunc("DELETE /partidos/{id}", h.EliminarPartido)
}

// GET /  (opcional: ?tab=jugadores o ?tab=partidos)
func (h *PageHandler) Index(w http.ResponseWriter, r *http.Request) {
	tab := views.NormalizarTab(r.URL.Query().Get("tab"))
	h.renderIndex(w, r, http.StatusOK, tab, views.FormState{})
}

func (h *PageHandler) renderIndex(w http.ResponseWriter, r *http.Request, status int, tab string, form views.FormState) {
	ctx := r.Context()

	equipos, err := h.equipos.GetEquipos(ctx)
	if err != nil {
		errorInterno(w, "listar equipos", err)
		return
	}
	jugadores, err := h.jugadores.GetJugadores(ctx)
	if err != nil {
		errorInterno(w, "listar jugadores", err)
		return
	}
	partidos, err := h.partidos.ListPartidos(ctx)
	if err != nil {
		errorInterno(w, "listar partidos", err)
		return
	}

	page := views.IndexPage(views.PageData{
		Tab:       tab,
		Equipos:   equipos,
		Jugadores: jugadores,
		Partidos:  partidos,
		Form:      form,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := page.Render(ctx, w); err != nil {
		log.Printf("error al renderizar la página: %v", err)
	}
}

// ---------- EQUIPOS ----------

// GET /equipos/lista
func (h *PageHandler) ListarEquiposPartial(w http.ResponseWriter, r *http.Request) {
	h.renderEquipoList(w, r)
}

// GET /equipos/{id}/edit
func (h *PageHandler) EditarEquipoForm(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	equipo, err := h.equipos.GetEquipoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrEquipoNotFound) {
			http.Error(w, "Equipo no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("error al obtener equipo %d: %v", id, err)
		http.Error(w, "No se pudo cargar el equipo", http.StatusInternalServerError)
		return
	}

	jugadores, err := h.jugadores.GetJugadores(r.Context())
	if err != nil {
		log.Printf("error al listar jugadores: %v", err)
		http.Error(w, "No se pudieron cargar los jugadores", http.StatusInternalServerError)
		return
	}

	if err := views.EquipoEditRow(equipo, jugadores, views.FormState{}).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar fila de edición: %v", err)
		http.Error(w, "No se pudo renderizar el formulario", http.StatusInternalServerError)
	}
}

// PUT /equipos/{id}
func (h *PageHandler) ActualizarEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	if !parsearFormulario(w, r) {
		return
	}

	valoracion, err := formInt16(r, "valoracion", "La valoración")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posicion, err := formInt16(r, "posicion_tabla", "La posición en la tabla")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var mejorJugador *int32
	if mj, err := formInt32(r, "mejor_jugador"); err == nil && mj > 0 {
		mejorJugador = &mj
	}

	req := dto.UpdateEquipoRequest{
		Nombre:          formTexto(r, "nombre"),
		FormacionActual: formTexto(r, "formacion_actual"),
		EscudoUrl:       formOpcional(r, "escudo_url"),
		MejorJugador:    mejorJugador,
		Valoracion:      valoracion,
		PosicionTabla:   posicion,
	}

	if err := req.Validate(); err != nil {
		http.Error(w, mayuscula(err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := h.equipos.UpdateEquipo(r.Context(), id, req); err != nil {
		if errors.Is(err, service.ErrEquipoNotFound) {
			http.Error(w, "Equipo no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("error al actualizar equipo %d: %v", id, err)
		http.Error(w, "No se pudo actualizar el equipo", http.StatusInternalServerError)
		return
	}

	h.renderEquipoList(w, r)
}

func (h *PageHandler) renderEquipoList(w http.ResponseWriter, r *http.Request) {
	equipos, err := h.equipos.GetEquipos(r.Context())
	if err != nil {
		log.Printf("error al listar equipos: %v", err)
		http.Error(w, "No se pudieron cargar los equipos", http.StatusInternalServerError)
		return
	}

	jugadores, err := h.jugadores.GetJugadores(r.Context())
	if err != nil {
		log.Printf("error al listar jugadores: %v", err)
		http.Error(w, "No se pudieron cargar los jugadores", http.StatusInternalServerError)
		return
	}

	if err := views.EquipoList(equipos, jugadores).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar lista de equipos: %v", err)
		http.Error(w, "No se pudo renderizar la lista", http.StatusInternalServerError)
	}
}

func errorInterno(w http.ResponseWriter, accion string, err error) {
	log.Printf("error al %s: %v", accion, err)
	http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
}

// ---------- JUGADORES ----------

// GET /jugadores/lista
func (h *PageHandler) ListarJugadoresPartial(w http.ResponseWriter, r *http.Request) {
	h.renderJugadorList(w, r)
}

// GET /jugadores/{id}/edit
func (h *PageHandler) EditarJugadorForm(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de jugador inválido", http.StatusBadRequest)
		return
	}
	jugador, err := h.jugadores.GetJugador(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrJugadorNotFound) {
			http.Error(w, "Jugador no encontrado", http.StatusNotFound)
			return
		}
		errorInterno(w, "obtener jugador", err)
		return
	}
	equipos, err := h.equipos.GetEquipos(r.Context())
	if err != nil {
		errorInterno(w, "listar equipos", err)
		return
	}
	if err := views.JugadorEditRow(jugador, equipos, views.FormState{}).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar fila de edición: %v", err)
	}
}

// PUT /jugadores/{id}
func (h *PageHandler) ActualizarJugador(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de jugador inválido", http.StatusBadRequest)
		return
	}
	if !parsearFormulario(w, r) {
		return
	}

	equipoID, err := formInt32(r, "equipo_id")
	if err != nil {
		http.Error(w, "Elegí un equipo para el jugador.", http.StatusBadRequest)
		return
	}
	altura := formTexto(r, "altura")
	if _, err := strconv.ParseFloat(altura, 64); err != nil {
		http.Error(w, "La altura debe ser un número, por ejemplo 1.78.", http.StatusBadRequest)
		return
	}

	req := dto.UpdateJugadorRequest{
		EquipoID:        equipoID,
		Nombre:          formTexto(r, "nombre"),
		Posicion:        formTexto(r, "posicion"),
		FechaNacimiento: formTexto(r, "fecha_nacimiento"),
		Altura:          altura,
		FotoUrl:         formOpcional(r, "foto_url"),
	}

	// Mismo truco que CrearJugador: una tabla campo → destino para no repetir 9 veces el parseo.
	numeros := []struct {
		campo, etiqueta string
		destino         *int16
	}{
		{"media_general", "La media general", &req.MediaGeneral},
		{"ritmo", "El ritmo", &req.Ritmo},
		{"tiro", "El tiro", &req.Tiro},
		{"pase", "El pase", &req.Pase},
		{"regate", "El regate", &req.Regate},
		{"defensa", "La defensa", &req.Defensa},
		{"fisico", "El físico", &req.Fisico},
		{"goles", "Los goles", &req.Goles},
		{"asistencias", "Las asistencias", &req.Asistencias},
	}
	for _, n := range numeros {
		v, err := formInt16(r, n.campo, n.etiqueta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		*n.destino = v
	}

	if err := req.Validate(); err != nil {
		http.Error(w, mayuscula(err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := h.jugadores.UpdateJugador(r.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, service.ErrJugadorNotFound):
			http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		case errors.Is(err, service.ErrEquipoNotFound):
			http.Error(w, "El equipo elegido ya no existe.", http.StatusBadRequest)
		case errors.Is(err, service.ErrFechaInvalida):
			http.Error(w, mayuscula(err.Error()), http.StatusBadRequest)
		default:
			errorInterno(w, "actualizar jugador", err)
		}
		return
	}

	h.renderJugadorList(w, r)
}

func (h *PageHandler) renderJugadorList(w http.ResponseWriter, r *http.Request) {
	jugadores, err := h.jugadores.GetJugadores(r.Context())
	if err != nil {
		errorInterno(w, "listar jugadores", err)
		return
	}
	equipos, err := h.equipos.GetEquipos(r.Context())
	if err != nil {
		errorInterno(w, "listar equipos", err)
		return
	}
	if err := views.JugadorList(jugadores, equipos).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar lista de jugadores: %v", err)
	}
}

// ---------- PARTIDOS ----------

// GET /partidos/lista
func (h *PageHandler) ListarPartidosPartial(w http.ResponseWriter, r *http.Request) {
	h.renderPartidoList(w, r)
}

// GET /partidos/{id}/edit
func (h *PageHandler) EditarPartidoForm(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de partido inválido", http.StatusBadRequest)
		return
	}
	partido, err := h.partidos.GetPartido(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPartidoNotFound) {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}
		errorInterno(w, "obtener partido", err)
		return
	}
	equipos, err := h.equipos.GetEquipos(r.Context())
	if err != nil {
		errorInterno(w, "listar equipos", err)
		return
	}
	if err := views.PartidoEditRow(partido, equipos, views.FormState{}).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar fila de edición: %v", err)
	}
}

// PUT /partidos/{id}
func (h *PageHandler) ActualizarPartido(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de partido inválido", http.StatusBadRequest)
		return
	}
	if !parsearFormulario(w, r) {
		return
	}

	// Vacío = el partido no tiene ese equipo (fue eliminado); queda en NULL.
	local, err := formInt32Opcional(r, "equipo_local_id")
	if err != nil {
		http.Error(w, "Equipo local inválido.", http.StatusBadRequest)
		return
	}
	visitante, err := formInt32Opcional(r, "equipo_visitante_id")
	if err != nil {
		http.Error(w, "Equipo visitante inválido.", http.StatusBadRequest)
		return
	}
	golesLocal, err := formInt16(r, "goles_local", "Los goles del local")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	golesVisitante, err := formInt16(r, "goles_visitante", "Los goles del visitante")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fecha, _ := time.Parse("2006-01-02", formTexto(r, "fecha")) // si falla queda en cero y Validate lo rechaza

	req := dto.UpdatePartidoRequest{
		EquipoLocalID:     local,
		EquipoVisitanteID: visitante,
		GolesLocal:        golesLocal,
		GolesVisitante:    golesVisitante,
		Fecha:             dto.DateOnly{Time: fecha},
		Estado:            formTexto(r, "estado"),
	}
	if err := req.Validate(); err != nil {
		http.Error(w, mayuscula(err.Error()), http.StatusBadRequest)
		return
	}

	if _, err := h.partidos.UpdatePartido(r.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, service.ErrPartidoNotFound):
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
		case errors.Is(err, service.ErrEquipoLocalNotFound), errors.Is(err, service.ErrEquipoVisitNotFound):
			http.Error(w, mayuscula(err.Error()), http.StatusBadRequest)
		default:
			errorInterno(w, "actualizar partido", err)
		}
		return
	}

	h.renderPartidoList(w, r)
}

func (h *PageHandler) renderPartidoList(w http.ResponseWriter, r *http.Request) {
	partidos, err := h.partidos.ListPartidos(r.Context())
	if err != nil {
		errorInterno(w, "listar partidos", err)
		return
	}
	equipos, err := h.equipos.GetEquipos(r.Context())
	if err != nil {
		errorInterno(w, "listar equipos", err)
		return
	}
	if err := views.PartidoList(partidos, equipos).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar lista de partidos: %v", err)
	}
}
