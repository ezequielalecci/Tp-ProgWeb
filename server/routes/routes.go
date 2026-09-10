package routes

import (
	"net/http"

	db "tpweb/db/sqlc"
	"tpweb/internal/handlers"
)

func RegisterRoutes(queries *db.Queries) http.Handler {
	mux := http.NewServeMux()

	// archivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("GET /resultados-api.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/resultados-api.html")
	})

	// raiz
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// handlers
	equiposH := handlers.NewAPIEquiposHandler(queries)
	jugadoresH := handlers.NewAPIJugadoresHandler(queries)
	partidosH := handlers.NewAPIPartidosHandler(queries)

	// Equipos
	mux.HandleFunc("GET /api/equipos", equiposH.GetEquipos)
	mux.HandleFunc("GET /api/equipos/{id}", equiposH.GetEquipoByID)
	mux.HandleFunc("POST /api/equipos", equiposH.CrearEquipo)
	mux.HandleFunc("PUT /api/equipos/{id}", equiposH.UpdateEquipo)
	mux.HandleFunc("DELETE /api/equipos/{id}", equiposH.DeleteEquipo)

	// Jugadores
	mux.HandleFunc("GET /api/jugadores", jugadoresH.GetJugadores)
	mux.HandleFunc("GET /api/jugador/{id}", jugadoresH.GetJugador)
	mux.HandleFunc("GET /api/jugadores/equipo", jugadoresH.ListJugadoresByEquipo)
	mux.HandleFunc("GET /api/equipos/{equipo_id}/mejor-jugador", jugadoresH.GetMejorJugadorByEquipo)
	mux.HandleFunc("POST /api/jugador", jugadoresH.CreateJugador)
	mux.HandleFunc("PUT /api/jugador/{id}/stats", jugadoresH.UpdateJugadorStats)
	mux.HandleFunc("DELETE /api/jugador/{id}", jugadoresH.DeleteJugador)

	// Partidos
	mux.HandleFunc("GET /api/partidos", partidosH.ListPartidos)
	mux.HandleFunc("GET /api/partidos/{id}", partidosH.GetPartido)
	mux.HandleFunc("POST /api/partidos", partidosH.CrearPartido)
	mux.HandleFunc("PUT /api/partidos/{id}/resultado", partidosH.UpdateResultadoPartido)
	mux.HandleFunc("DELETE /api/partidos/{id}", partidosH.DeletePartido)

	mux.HandleFunc("GET /api/equipos/{equipo_id}/partidos/ultimo", partidosH.GetUltimoPartidoByEquipo)
	mux.HandleFunc("GET /api/equipos/{equipo_id}/partidos/ultimo-finalizado", partidosH.GetUltimoPartidoFinalizadoByEquipo)
	mux.HandleFunc("GET /api/equipos/{equipo_id}/partidos/ultimo-detalle", partidosH.GetUltimoPartidoConDetalleByEquipo)

	return mux
}
