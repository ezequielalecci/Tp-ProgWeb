package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"tpweb/internal/dto"
	"tpweb/internal/service"
)

type APIJugadoresHandler struct {
	service *service.JugadorService
}

func NewAPIJugadoresHandler(s *service.JugadorService) *APIJugadoresHandler {
	return &APIJugadoresHandler{
		service: s,
	}
}

// GET /api/jugador/{id}
func (h *APIJugadoresHandler) GetJugador(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de jugador inválido")
		return
	}

	jugador, err := h.service.GetJugador(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrJugadorNotFound) {
			respondError(w, http.StatusNotFound, "Jugador no encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al buscar el jugador")
		return
	}

	respondJSON(w, http.StatusOK, jugador)
}

// GET /api/jugadores
func (h *APIJugadoresHandler) GetJugadores(w http.ResponseWriter, r *http.Request) {
	jugadores, err := h.service.GetJugadores(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error al obtener la lista de jugadores")
		return
	}

	respondJSON(w, http.StatusOK, jugadores)
}

// GET /api/jugadores/equipo?equipo_id=X
func (h *APIJugadoresHandler) ListJugadoresByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.URL.Query().Get("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de equipo query param inválido o faltante")
		return
	}

	jugadores, err := h.service.ListJugadoresByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error al obtener jugadores del equipo")
		return
	}

	respondJSON(w, http.StatusOK, jugadores)
}

// POST /api/jugadores
func (h *APIJugadoresHandler) CreateJugador(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJugadorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	jugador, err := h.service.CreateJugador(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrFechaInvalida) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrEquipoNotFound) {
			respondError(w, http.StatusBadRequest, "El equipo asignado no existe")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al crear jugador en la base de datos")
		return
	}

	respondJSON(w, http.StatusCreated, jugador)
}

// PUT /api/jugadores/{id}/stats
func (h *APIJugadoresHandler) UpdateJugadorStats(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de jugador inválido")
		return
	}

	var req dto.UpdateJugadorStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	jugadorActualizado, err := h.service.UpdateJugadorStats(r.Context(), int32(id), req)
	if err != nil {
		if errors.Is(err, service.ErrJugadorNotFound) {
			respondError(w, http.StatusNotFound, "El jugador no existe")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al actualizar estadísticas del jugador")
		return
	}

	respondJSON(w, http.StatusOK, jugadorActualizado)
}

// DELETE /api/jugadores/{id}
func (h *APIJugadoresHandler) DeleteJugador(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de jugador inválido")
		return
	}

	err = h.service.DeleteJugador(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrJugadorNotFound) {
			respondError(w, http.StatusNotFound, "El jugador no existe")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al eliminar el jugador")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/equipos/{equipo_id}/mejor-jugador
func (h *APIJugadoresHandler) GetMejorJugadorByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.PathValue("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	jugador, err := h.service.GetMejorJugadorByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		if errors.Is(err, service.ErrJugadorNotFound) {
			respondError(w, http.StatusNotFound, "No se encontró el jugador")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al consultar el mejor jugador del equipo")
		return
	}

	respondJSON(w, http.StatusOK, jugador)
}
