package routes

import (
	"net/http"

	db "tpweb/db/sqlc"
	"tpweb/internal/handlers"
	"tpweb/internal/service"
)

func RegisterRoutes(queries *db.Queries) http.Handler {
	mux := http.NewServeMux()

	registerStaticRoutes(mux)

	equipoService := service.NewEquipoService(queries)
	jugadorService := service.NewJugadorService(queries)
	partidosService := service.NewPartidoService(queries)

	equipoHandler := handlers.NewAPIEquiposHandler(equipoService)
	jugadorHandler := handlers.NewAPIJugadoresHandler(jugadorService)
	partidoHandler := handlers.NewAPIPartidosHandler(partidosService)

	registerEquiposRoutes(mux, equipoHandler)
	registerJugadoresRoutes(mux, jugadorHandler)
	registerPartidosRoutes(mux, partidoHandler)

	return mux
}

func registerStaticRoutes(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /resultados-api.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/resultados-api.html")
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

}

func registerEquiposRoutes(mux *http.ServeMux, h *handlers.APIEquiposHandler) {
	mux.HandleFunc("GET /api/equipos", h.GetEquipos)
	mux.HandleFunc("GET /api/equipos/{id}", h.GetEquipoByID)
	mux.HandleFunc("POST /api/equipos", h.CrearEquipo)
	mux.HandleFunc("PUT /api/equipos/{id}", h.UpdateEquipo)
	mux.HandleFunc("DELETE /api/equipos/{id}", h.DeleteEquipo)
}

func registerJugadoresRoutes(mux *http.ServeMux, h *handlers.APIJugadoresHandler) {
	mux.HandleFunc("GET /api/jugadores", h.GetJugadores)
	mux.HandleFunc("GET /api/jugadores/{id}", h.GetJugador)
	mux.HandleFunc("GET /api/jugadores/equipo", h.ListJugadoresByEquipo)
	mux.HandleFunc("GET /api/equipos/{equipo_id}/mejor-jugador", h.GetMejorJugadorByEquipo)
	mux.HandleFunc("POST /api/jugadores", h.CreateJugador)
	mux.HandleFunc("PUT /api/jugadores/{id}/stats", h.UpdateJugadorStats)
	mux.HandleFunc("DELETE /api/jugadoress/{id}", h.DeleteJugador)

}

func registerPartidosRoutes(mux *http.ServeMux, h *handlers.APIPartidosHandler) {
	mux.HandleFunc("GET /api/partidos", h.ListPartidos)
	mux.HandleFunc("GET /api/partidos/{id}", h.GetPartido)
	mux.HandleFunc("POST /api/partidos", h.CrearPartido)
	mux.HandleFunc("PUT /api/partidos/{id}/resultado", h.UpdateResultadoPartido)
	mux.HandleFunc("DELETE /api/partidos/{id}", h.DeletePartido)
	mux.HandleFunc("GET /api/equipos/{equipo_id}/partidos/ultimo", h.GetUltimoPartidoByEquipo)
	mux.HandleFunc("GET /api/equipos/{equipo_id}/partidos/ultimo-finalizado", h.GetUltimoPartidoFinalizadoByEquipo)
	mux.HandleFunc("GET /api/equipos/{equipo_id}/partidos/ultimo-detalle", h.GetUltimoPartidoConDetalleByEquipo)
}
