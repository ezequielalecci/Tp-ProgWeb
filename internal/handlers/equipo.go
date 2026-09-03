package handlers

import (
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
	idStr := r.URL.Query().Get("id")
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

func (h *APIEquiposHandler) CreateEquipo(w http.ResponseWriter, r *http.Request) {
	var params database.CreateEquipoParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Payload JSON inválido", http.StatusBadRequest)
		return
	}

	equipo, err := h.queries.CreateEquipo(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear el equipo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(equipo)
}

func (h *APIEquiposHandler) UpdateEquipos(w http.ResponseWriter, r *http.Request) {
	var params database.UpdateEquipoParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Payload JSON inválido", http.StatusBadRequest)
		return
	}

	equipoActualizado, err := h.queries.UpdateEquipo(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar el equipo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipoActualizado)
}

func (h *APIEquiposHandler) DeleteEquipo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteEquipo(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar el equipo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
