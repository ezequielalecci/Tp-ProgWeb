package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"tpweb/internal/dto"
	"tpweb/internal/service"
)

// APIEquiposHandler delega la lógica de negocio a EquipoService
type APIEquiposHandler struct {
	service *service.EquipoService
}

// NewAPIEquiposHandler es el constructor para instanciar el handler inyectándole el servicio
func NewAPIEquiposHandler(s *service.EquipoService) *APIEquiposHandler {
	return &APIEquiposHandler{
		service: s,
	}
}

func (h *APIEquiposHandler) RegisterEquiposRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/equipos", h.GetEquipos)
	mux.HandleFunc("GET /api/equipos/{id}", h.GetEquipoByID)
	mux.HandleFunc("POST /api/equipos", h.CrearEquipo)
	mux.HandleFunc("PUT /api/equipos/{id}", h.UpdateEquipo)
	mux.HandleFunc("DELETE /api/equipos/{id}", h.DeleteEquipo)
}

// GET /api/equipos - Obtener todos los equipos
func (h *APIEquiposHandler) GetEquipos(w http.ResponseWriter, r *http.Request) {
	equipos, err := h.service.GetEquipos(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error al obtener los equipos")
		return
	}

	respondJSON(w, http.StatusOK, equipos)
}

// GET /api/equipos/{id} - Obtener equipo por ID
func (h *APIEquiposHandler) GetEquipoByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	equipo, err := h.service.GetEquipoByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrEquipoNotFound) {
			respondError(w, http.StatusNotFound, "Equipo no encontrado")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al buscar el equipo")
		return
	}

	respondJSON(w, http.StatusOK, equipo)
}

// POST /api/equipos - Crear un nuevo equipo
func (h *APIEquiposHandler) CrearEquipo(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEquipoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	equipo, err := h.service.CrearEquipo(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error al crear el equipo")
		return
	}

	respondJSON(w, http.StatusCreated, equipo)
}

// PUT /api/equipos/{id} - Actualizar equipo existente
func (h *APIEquiposHandler) UpdateEquipo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	var req dto.UpdateEquipoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload JSON inválido: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	equipoActualizado, err := h.service.UpdateEquipo(r.Context(), int32(id), req)
	if err != nil {
		if errors.Is(err, service.ErrEquipoNotFound) {
			respondError(w, http.StatusNotFound, "El equipo no existe")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al actualizar el equipo")
		return
	}

	respondJSON(w, http.StatusOK, equipoActualizado)
}

// DELETE /api/equipos/{id} - Eliminar equipo
func (h *APIEquiposHandler) DeleteEquipo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	err = h.service.DeleteEquipo(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrEquipoNotFound) {
			respondError(w, http.StatusNotFound, "El equipo a eliminar no existe")
			return
		}
		respondError(w, http.StatusInternalServerError, "Error al eliminar el equipo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
