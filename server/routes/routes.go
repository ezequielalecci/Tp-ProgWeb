package routes

import (
	"net/http"

	db "tpweb/db/sqlc"
	"tpweb/internal/handlers"
	"tpweb/internal/service"
)

func RegisterRoutes(queries *db.Queries) http.Handler {
	mux := http.NewServeMux()

	// Archivos estáticos: solo queda la hoja de estilos propia.
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Instanciación de Servicios
	equipoService := service.NewEquipoService(queries)
	jugadorService := service.NewJugadorService(queries)
	partidosService := service.NewPartidoService(queries)

	// Un único handler de páginas: renderiza las vistas y procesa los formularios.
	pages := handlers.NewPageHandler(equipoService, jugadorService, partidosService)
	pages.RegisterRoutes(mux)

	return mux
}
