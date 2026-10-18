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

	// Instanciación de Servicios
	equipoService := service.NewEquipoService(queries)
	jugadorService := service.NewJugadorService(queries)
	partidosService := service.NewPartidoService(queries)

	// Instanciación de Handlers
	equipoHandler := handlers.NewAPIEquiposHandler(equipoService)
	jugadorHandler := handlers.NewAPIJugadoresHandler(jugadorService)
	partidoHandler := handlers.NewAPIPartidosHandler(partidosService)

	// Delegación del registro de rutas a cada handler
	equipoHandler.RegisterEquiposRoutes(mux)
	jugadorHandler.RegisterJugadoresRoutes(mux)
	partidoHandler.RegisterPartidosRoutes(mux)

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
