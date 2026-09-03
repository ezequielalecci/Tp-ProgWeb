package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	database "tpweb/db/sqlc"
)

type APIJugadoresHandler struct {
	queries *database.Queries
}

func NewAPIJugadoresHandler(q *database.Queries) *APIJugadoresHandler {
	return &APIJugadoresHandler{
		queries: q,
	}
}

func (h *APIJugadoresHandler) GetJugador(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	jugador, err := h.queries.GetJugador(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugador)
}

func (h *APIJugadoresHandler) ListJugadoresByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.URL.Query().Get("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	jugadores, err := h.queries.ListJugadoresByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		http.Error(w, "Error al obtener jugadores", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugadores)
}

func (h *APIJugadoresHandler) CreateJugador(w http.ResponseWriter, r *http.Request) {
	var params database.CreateJugadorParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Payload JSON inválido", http.StatusBadRequest)
		return
	}

	jugador, err := h.queries.CreateJugador(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear jugador: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jugador)
}

func (h *APIJugadoresHandler) UpdateJugadorStats(w http.ResponseWriter, r *http.Request) {
	var params database.UpdateJugadorStatsParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Payload JSON inválido", http.StatusBadRequest)
		return
	}

	jugador, err := h.queries.UpdateJugadorStats(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar estadísticas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugador)
}

func (h *APIJugadoresHandler) DeleteJugador(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteJugador(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar jugador: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *APIJugadoresHandler) GetMejorJugadorByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.URL.Query().Get("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	jugador, err := h.queries.GetMejorJugadorByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugador)
}
