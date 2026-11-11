package handlers

import (
	"log"
	"net/http"

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
	mux.HandleFunc("POST /equipos", h.CrearEquipo)
	mux.HandleFunc("POST /jugadores", h.CrearJugador)
	mux.HandleFunc("POST /partidos", h.CrearPartido)
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

func errorInterno(w http.ResponseWriter, accion string, err error) {
	log.Printf("error al %s: %v", accion, err)
	http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
}
