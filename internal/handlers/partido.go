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
	idStr := r.PathValue("id")
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

func (h *APIPartidosHandler) CrearPartido(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la petición JSON entrante
	var req CreatePartidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Fecha o JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Mapear los datos al struct generado por sqlc
	params := db.CreatePartidoParams{
		EquipoLocalID:     req.EquipoLocalID,
		EquipoVisitanteID: req.EquipoVisitanteID,
		Fecha:             req.Fecha.Time, // Pasa el time.Time compatible con sqlc
		Estado:            req.Estado,
	}

	// 3. Ejecutar la consulta en la base de datos
	partido, err := h.queries.CreatePartido(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(partido)
}

func (h *APIPartidosHandler) UpdateResultadoPartido(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de partido inválido", http.StatusBadRequest)
		return
	}

	var req UpdateResultadoPartidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	params := db.UpdateResultadoPartidoParams{
		ID:             int32(id),
		GolesLocal:     req.GolesLocal,
		GolesVisitante: req.GolesVisitante,
		Estado:         req.Estado,
	}

	partido, err := h.queries.UpdateResultadoPartido(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar el resultado del partido", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partido)
}

func (h *APIPartidosHandler) DeletePartido(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.DeletePartido(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar el partido", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *APIPartidosHandler) GetUltimoPartidoByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.PathValue("equipo_id")
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
	equipoIDStr := r.PathValue("equipo_id")
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
	equipoIDStr := r.PathValue("equipo_id")
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
