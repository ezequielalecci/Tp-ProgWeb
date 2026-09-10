package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	database "tpweb/db/sqlc"
)

// APIHandler guarda la referencia a las consultas de sqlc
type APIEquiposHandler struct {
	queries *database.Queries
}

// constructor para instanciar el handler desde el main
func NewAPIEquiposHandler(q *database.Queries) *APIEquiposHandler {
	return &APIEquiposHandler{
		queries: q,
	}
}

func (h *APIEquiposHandler) GetEquipos(w http.ResponseWriter, r *http.Request) {
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

	if equipos == nil {
		equipos = []database.Equipo{}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(equipos)
}

func (h *APIEquiposHandler) GetEquipoByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	equipo, err := h.queries.GetEquipo(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Equipo no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipo)
}

func (h *APIEquiposHandler) CrearEquipo(w http.ResponseWriter, r *http.Request) {
	var req CreateEquipoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	params := database.CreateEquipoParams{
		Nombre:          req.Nombre,
		FormacionActual: req.FormacionActual,
		Valoracion:      req.Valoracion,
		PosicionTabla:   req.PosicionTabla,
	}

	if req.EscudoUrl != nil {
		params.EscudoUrl = sql.NullString{String: *req.EscudoUrl, Valid: true}
	}

	equipo, err := h.queries.CreateEquipo(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear equipo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(equipo)
}

func (h *APIEquiposHandler) UpdateEquipo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	var req UpdateEquipoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convertir DTO a los tipos de sql.Null* que requiere sqlc
	params := database.UpdateEquipoParams{
		ID:              int32(id),
		Nombre:          req.Nombre,
		FormacionActual: req.FormacionActual,
		Valoracion:      req.Valoracion,
		PosicionTabla:   req.PosicionTabla,
	}

	if req.EscudoUrl != nil {
		params.EscudoUrl = sql.NullString{String: *req.EscudoUrl, Valid: true}
	}

	if req.MejorJugador != nil {
		params.MejorJugador = sql.NullInt32{Int32: *req.MejorJugador, Valid: true}
	}

	equipoActualizado, err := h.queries.UpdateEquipo(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar el equipo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipoActualizado)
}

func (h *APIEquiposHandler) DeleteEquipo(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteEquipo(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar el equipo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
