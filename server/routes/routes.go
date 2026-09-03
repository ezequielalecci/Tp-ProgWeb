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
	mux.HandleFunc("POST /api/equipos", equiposH.CreateEquipo)
	mux.HandleFunc("PUT /api/equipos", equiposH.UpdateEquipos)
	mux.HandleFunc("DELETE /api/equipos", equiposH.DeleteEquipo)

	// Jugadores
	mux.HandleFunc("GET /api/jugadores", jugadoresH.GetJugador)
	mux.HandleFunc("GET /api/jugadores/equipo", jugadoresH.ListJugadoresByEquipo)
	mux.HandleFunc("GET /api/jugadores/mejor", jugadoresH.GetMejorJugadorByEquipo)
	mux.HandleFunc("POST /api/jugadores", jugadoresH.CreateJugador)
	mux.HandleFunc("PUT /api/jugadores/stats", jugadoresH.UpdateJugadorStats)
	mux.HandleFunc("DELETE /api/jugadores", jugadoresH.DeleteJugador)

	// Partidos
	mux.HandleFunc("GET /api/partidos", partidosH.ListPartidos)
	mux.HandleFunc("GET /api/partidos/detalle", partidosH.GetPartido)
	mux.HandleFunc("GET /api/partidos/ultimo", partidosH.GetUltimoPartidoByEquipo)
	mux.HandleFunc("GET /api/partidos/ultimo-finalizado", partidosH.GetUltimoPartidoFinalizadoByEquipo)
	mux.HandleFunc("GET /api/partidos/ultimo-detalle", partidosH.GetUltimoPartidoConDetalleByEquipo)
	mux.HandleFunc("POST /api/partidos", partidosH.CreatePartido)
	mux.HandleFunc("PUT /api/partidos/resultado", partidosH.UpdateResultadoPartido)
	mux.HandleFunc("DELETE /api/partidos", partidosH.DeletePartido)

	return mux
}
