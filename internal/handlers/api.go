package handlers

import (
	"encoding/json"
	"net/http"

	"tp2/internal/database"
)

// APIHandler guarda la referencia a las consultas de sqlc
type APIHandler struct {
	queries *database.Queries
}

// NewAPIHandler es el constructor para instanciar el handler desde el main
func NewAPIHandler(q *database.Queries) *APIHandler {
	return &APIHandler{
		queries: q,
	}
}

// GetEquipos lista todos los equipos (Solo GET)
func (h *APIHandler) GetEquipos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	equipos, err := h.queries.ListEquipos(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error de base de datos"})
		return
	}

	// Si la lista de la BD es nil, enviamos un slice vacío [] en lugar de null
	if equipos == nil {
		equipos = []database.Equipo{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(equipos)
}
