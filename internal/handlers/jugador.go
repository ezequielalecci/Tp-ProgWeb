package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

	idStr := r.PathValue("id")
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

func (h *APIJugadoresHandler) GetJugadores(w http.ResponseWriter, r *http.Request) {
	jugadores, err := h.queries.GetJugadores(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los jugadores", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugadores)
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
	var req CreateJugadorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parsear fecha flexible (AAAA-MM-DD)
	fechaNac, err := time.Parse("2006-01-02", req.FechaNacimiento)
	if err != nil {
		http.Error(w, "Formato de fecha inválido (debe ser AAAA-MM-DD)", http.StatusBadRequest)
		return
	}

	params := database.CreateJugadorParams{
		EquipoID:        req.EquipoID,
		Nombre:          req.Nombre,
		Posicion:        req.Posicion,
		FechaNacimiento: fechaNac,
		MediaGeneral:    req.MediaGeneral,
		Altura:          req.Altura,
		Ritmo:           req.Ritmo,
		Tiro:            req.Tiro,
		Pase:            req.Pase,
		Regate:          req.Regate,
		Defensa:         req.Defensa,
		Fisico:          req.Fisico,
	}

	if req.FotoUrl != nil {
		params.FotoUrl = sql.NullString{String: *req.FotoUrl, Valid: true}
	}

	jugador, err := h.queries.CreateJugador(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al crear jugador", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(jugador)
}

func (h *APIJugadoresHandler) UpdateJugadorStats(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de jugador inválido", http.StatusBadRequest)
		return
	}

	var req UpdateJugadorStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	params := database.UpdateJugadorStatsParams{
		ID:          int32(id),
		Goles:       req.Goles,
		Asistencias: req.Asistencias,
	}

	jugador, err := h.queries.UpdateJugadorStats(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al actualizar estadísticas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugador)
}

func (h *APIJugadoresHandler) DeleteJugador(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
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

	equipoIDStr := r.PathValue("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	jugador, err := h.queries.GetMejorJugadorByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		http.Error(w, "No se encontró el mejor jugador para este equipo", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jugador)
}
