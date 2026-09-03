package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "tpweb/db/sqlc"
)

type APIPartidosHandler struct {
	queries *db.Queries
}

func NewAPIPartidosHandler(q *db.Queries) *APIPartidosHandler {
	return &APIPartidosHandler{queries: q}
}

func (h *APIPartidosHandler) GetPartido(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	partido, err := h.queries.GetPartido(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partido)
}

func (h *APIPartidosHandler) ListPartidos(w http.ResponseWriter, r *http.Request) {
	partidos, err := h.queries.ListPartidos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener partidos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partidos)
}

func (h *APIPartidosHandler) CreatePartido(w http.ResponseWriter, r *http.Request) {
	var params db.CreatePartidoParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Payload JSON inválido", http.StatusBadRequest)
		return
	}

	partido, err := h.queries.CreatePartido(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear partido: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(partido)
}

func (h *APIPartidosHandler) UpdateResultadoPartido(w http.ResponseWriter, r *http.Request) {
	var params db.UpdateResultadoPartidoParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Payload JSON inválido", http.StatusBadRequest)
		return
	}

	partido, err := h.queries.UpdateResultadoPartido(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar resultado: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partido)
}

func (h *APIPartidosHandler) DeletePartido(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.DeletePartido(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar partido: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *APIPartidosHandler) GetUltimoPartidoByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.URL.Query().Get("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	partido, err := h.queries.GetUltimoPartidoByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partido)
}

func (h *APIPartidosHandler) GetUltimoPartidoFinalizadoByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.URL.Query().Get("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	partido, err := h.queries.GetUltimoPartidoFinalizadoByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partido)
}

func (h *APIPartidosHandler) GetUltimoPartidoConDetalleByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.URL.Query().Get("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	partidoDetalle, err := h.queries.GetUltimoPartidoConDetalleByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		http.Error(w, "Partido no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partidoDetalle)
}
