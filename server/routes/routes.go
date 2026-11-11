package routes

import (
	"net/http"

	db "tpweb/db/sqlc"
	"tpweb/internal/handlers"
	"tpweb/internal/service"
)

func RegisterRoutes(queries *db.Queries) http.Handler {
	mux := http.NewServeMux()

	// solo queda la hoja de estilos.
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	equipoService := service.NewEquipoService(queries)
	jugadorService := service.NewJugadorService(queries)
	partidosService := service.NewPartidoService(queries)

	pages := handlers.NewPageHandler(equipoService, jugadorService, partidosService)
	pages.RegisterRoutes(mux)

	return mux
}
