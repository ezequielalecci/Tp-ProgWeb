package dto

import (
	"errors"
	"strings"
)

type CreateEquipoRequest struct {
	Nombre          string  `json:"nombre"`
	FormacionActual string  `json:"formacion_actual"`
	EscudoUrl       *string `json:"escudo_url"`
	Valoracion      int16   `json:"valoracion"`
	PosicionTabla   int16   `json:"posicion_tabla"`
}

func (req *CreateEquipoRequest) Validate() error {
	if strings.TrimSpace(req.Nombre) == "" {
		return errors.New("el nombre del equipo no puede estar vacío")
	}
	if req.Valoracion < 0 {
		return errors.New("la valoración no puede ser un número negativo")
	}
	if req.PosicionTabla <= 0 {
		return errors.New("la posición en la tabla debe ser mayor a 0")
	}
	return nil
}

type UpdateEquipoRequest struct {
	Nombre          string  `json:"nombre"`
	FormacionActual string  `json:"formacion_actual"`
	EscudoUrl       *string `json:"escudo_url"`
	MejorJugador    *int32  `json:"mejor_jugador"`
	Valoracion      int16   `json:"valoracion"`
	PosicionTabla   int16   `json:"posicion_tabla"`
}

func (req *UpdateEquipoRequest) Validate() error {
	if strings.TrimSpace(req.Nombre) == "" {
		return errors.New("el nombre del equipo no puede estar vacío")
	}
	if req.Valoracion < 0 {
		return errors.New("la valoración no puede ser un número negativo")
	}
	if req.PosicionTabla <= 0 {
		return errors.New("la posición en la tabla debe ser mayor a 0")
	}
	if req.MejorJugador != nil && *req.MejorJugador <= 0 {
		return errors.New("el ID del mejor jugador debe ser un entero positivo")
	}
	return nil
}
