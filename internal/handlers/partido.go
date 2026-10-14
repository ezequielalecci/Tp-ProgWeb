package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"tpweb/internal/dto"
	"tpweb/internal/service"
)

type APIPartidosHandler struct {
	service *service.PartidoService
}

func NewAPIPartidosHandler(s *service.PartidoService) *APIPartidosHandler {
	return &APIPartidosHandler{service: s}
}

// GET /api/partidos/{id}
func (h *APIPartidosHandler) GetPartido(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "ID de partido inválido")
		return
	}

	partido, err := h.service.GetPartido(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrPartidoNotFound) {
			respondError(w, http.StatusNotFound, "Partido no encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al obtener el partido")
		return
	}

	respondJSON(w, http.StatusOK, partido)
}

// GET /api/partidos
func (h *APIPartidosHandler) ListPartidos(w http.ResponseWriter, r *http.Request) {
	partidos, err := h.service.ListPartidos(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error al obtener la lista de partidos")
		return
	}

	respondJSON(w, http.StatusOK, partidos)
}

// POST /api/partidos
func (h *APIPartidosHandler) CrearPartido(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePartidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload JSON o formato de fecha inválido: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	partido, err := h.service.CrearPartido(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEquipoLocalNotFound) || errors.Is(err, service.ErrEquipoVisitNotFound) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al crear el partido en la base de datos")
		return
	}

	respondJSON(w, http.StatusCreated, partido)
}

// PUT /api/partidos/{id}/resultado
func (h *APIPartidosHandler) UpdateResultadoPartido(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "ID de partido inválido")
		return
	}

	var req dto.UpdateResultadoPartidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	partido, err := h.service.UpdateResultadoPartido(r.Context(), int32(id), req)
	if err != nil {
		if errors.Is(err, service.ErrPartidoNotFound) {
			respondError(w, http.StatusNotFound, "El partido a actualizar no existe")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al actualizar el resultado del partido")
		return
	}

	respondJSON(w, http.StatusOK, partido)
}

// DELETE /api/partidos/{id}
func (h *APIPartidosHandler) DeletePartido(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "ID de partido inválido")
		return
	}

	err = h.service.DeletePartido(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrPartidoNotFound) {
			respondError(w, http.StatusNotFound, "El partido a eliminar no existe")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al eliminar el partido")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/equipos/{equipo_id}/partidos/ultimo
func (h *APIPartidosHandler) GetUltimoPartidoByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.PathValue("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil || equipoID <= 0 {
		respondError(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	partido, err := h.service.GetUltimoPartidoByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		if errors.Is(err, service.ErrPartidoNotFound) {
			respondError(w, http.StatusNotFound, "No se encontró ningún partido para este equipo")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al consultar el último partido del equipo")
		return
	}

	respondJSON(w, http.StatusOK, partido)
}

// GET /api/equipos/{equipo_id}/partidos/ultimo-finalizado
func (h *APIPartidosHandler) GetUltimoPartidoFinalizadoByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.PathValue("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil || equipoID <= 0 {
		respondError(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	partido, err := h.service.GetUltimoPartidoFinalizadoByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		if errors.Is(err, service.ErrPartidoNotFound) {
			respondError(w, http.StatusNotFound, "No se encontró ningún partido finalizado para este equipo")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al consultar el último partido finalizado")
		return
	}

	respondJSON(w, http.StatusOK, partido)
}

// GET /api/equipos/{equipo_id}/partidos/ultimo-detalle
func (h *APIPartidosHandler) GetUltimoPartidoConDetalleByEquipo(w http.ResponseWriter, r *http.Request) {
	equipoIDStr := r.PathValue("equipo_id")
	equipoID, err := strconv.Atoi(equipoIDStr)
	if err != nil || equipoID <= 0 {
		respondError(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	partidoDetalle, err := h.service.GetUltimoPartidoConDetalleByEquipo(r.Context(), int32(equipoID))
	if err != nil {
		if errors.Is(err, service.ErrPartidoNotFound) {
			respondError(w, http.StatusNotFound, "No se encontraron detalles del último partido para este equipo")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al consultar el detalle del último partido")
		return
	}

	respondJSON(w, http.StatusOK, partidoDetalle)
}
