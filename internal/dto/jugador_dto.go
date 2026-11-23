package dto

import (
	"errors"
	"strings"
	"time"
)

type CreateJugadorRequest struct {
	EquipoID        int32   `json:"equipo_id"`
	Nombre          string  `json:"nombre"`
	Posicion        string  `json:"posicion"`
	FechaNacimiento string  `json:"fecha_nacimiento"` // Formato AAAA-MM-DD
	MediaGeneral    int16   `json:"media_general"`
	Altura          string  `json:"altura"`
	Ritmo           int16   `json:"ritmo"`
	Tiro            int16   `json:"tiro"`
	Pase            int16   `json:"pase"`
	Regate          int16   `json:"regate"`
	Defensa         int16   `json:"defensa"`
	Fisico          int16   `json:"fisico"`
	FotoUrl         *string `json:"foto_url"`
}

func (req *CreateJugadorRequest) Validate() error {
	if req.EquipoID <= 0 {
		return errors.New("el ID del equipo debe ser un entero positivo")
	}
	if strings.TrimSpace(req.Nombre) == "" {
		return errors.New("el nombre del jugador no puede estar vacío")
	}
	if strings.TrimSpace(req.Posicion) == "" {
		return errors.New("la posición del jugador no puede estar vacía")
	}
	if _, err := time.Parse("2006-01-02", req.FechaNacimiento); err != nil {
		return errors.New("formato de fecha de nacimiento inválido (debe ser AAAA-MM-DD)")
	}
	if req.MediaGeneral < 0 || req.MediaGeneral > 99 {
		return errors.New("la media general debe estar entre 0 y 99")
	}
	return nil
}

type UpdateJugadorRequest struct {
	EquipoID        int32   `json:"equipo_id"`
	Nombre          string  `json:"nombre"`
	Posicion        string  `json:"posicion"`
	FechaNacimiento string  `json:"fecha_nacimiento"` // AAAA-MM-DD
	MediaGeneral    int16   `json:"media_general"`
	Altura          string  `json:"altura"`
	Ritmo           int16   `json:"ritmo"`
	Tiro            int16   `json:"tiro"`
	Pase            int16   `json:"pase"`
	Regate          int16   `json:"regate"`
	Defensa         int16   `json:"defensa"`
	Fisico          int16   `json:"fisico"`
	Goles           int16   `json:"goles"`
	Asistencias     int16   `json:"asistencias"`
	FotoUrl         *string `json:"foto_url"`
}

func (req *UpdateJugadorRequest) Validate() error {
	if req.EquipoID <= 0 {
		return errors.New("el ID del equipo debe ser un entero positivo")
	}
	if strings.TrimSpace(req.Nombre) == "" {
		return errors.New("el nombre del jugador no puede estar vacío")
	}
	if strings.TrimSpace(req.Posicion) == "" {
		return errors.New("la posición del jugador no puede estar vacía")
	}
	if _, err := time.Parse("2006-01-02", req.FechaNacimiento); err != nil {
		return errors.New("formato de fecha de nacimiento inválido (debe ser AAAA-MM-DD)")
	}
	stats := []struct {
		nombre string
		valor  int16
	}{
		{"la media general", req.MediaGeneral},
		{"el ritmo", req.Ritmo},
		{"el tiro", req.Tiro},
		{"el pase", req.Pase},
		{"el regate", req.Regate},
		{"la defensa", req.Defensa},
		{"el físico", req.Fisico},
	}
	for _, s := range stats {
		if s.valor < 1 || s.valor > 99 {
			return errors.New(s.nombre + " debe estar entre 1 y 99")
		}
	}
	if req.Goles < 0 {
		return errors.New("los goles no pueden ser un valor negativo")
	}
	if req.Asistencias < 0 {
		return errors.New("las asistencias no pueden ser un valor negativo")
	}
	return nil
}
