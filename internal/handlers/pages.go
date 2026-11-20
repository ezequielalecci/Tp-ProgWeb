package handlers

import (
	"errors"
	"log"
	"net/http"

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
	mux.HandleFunc("POST /jugadores", h.CrearJugador)
	mux.HandleFunc("POST /partidos", h.CrearPartido)
<<<<<<< Updated upstream
	mux.HandleFunc("DELETE /equipos/{id}", h.EliminarEquipo)
	mux.HandleFunc("DELETE /jugadores/{id}", h.EliminarJugador)
	mux.HandleFunc("DELETE /partidos/{id}", h.EliminarPartido)
=======
>>>>>>> Stashed changes
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
