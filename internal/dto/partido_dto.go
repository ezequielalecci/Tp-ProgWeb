package dto

import (
	"errors"
	"strings"
	"time"
)

// DateOnly captura exclusivamente el formato YYYY-MM-DD del JSON
type DateOnly struct {
	time.Time
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

type CreatePartidoRequest struct {
	EquipoLocalID     int32    `json:"equipo_local_id"`
	EquipoVisitanteID int32    `json:"equipo_visitante_id"`
	Fecha             DateOnly `json:"fecha"` // Lee la cadena "YYYY-MM-DD"
	Estado            string   `json:"estado"`
}

func (req *CreatePartidoRequest) Validate() error {
	if req.EquipoLocalID <= 0 {
		return errors.New("el ID del equipo local debe ser un entero positivo")
	}
	if req.EquipoVisitanteID <= 0 {
		return errors.New("el ID del equipo visitante debe ser un entero positivo")
	}
	if req.EquipoLocalID == req.EquipoVisitanteID {
		return errors.New("el equipo local y el visitante no pueden ser el mismo")
	}
	if req.Fecha.IsZero() {
		return errors.New("la fecha del partido es obligatoria y debe tener formato YYYY-MM-DD")
	}
	if strings.TrimSpace(req.Estado) == "" {
		return errors.New("el estado del partido no puede estar vacío")
	}
	return nil
}

type UpdateResultadoPartidoRequest struct {
	GolesLocal     int16  `json:"goles_local"`
	GolesVisitante int16  `json:"goles_visitante"`
	Estado         string `json:"estado"` // PENDIENTE, JUGANDO, FINALIZADO
}

func (req *UpdateResultadoPartidoRequest) Validate() error {
	if req.GolesLocal < 0 {
		return errors.New("los goles del equipo local no pueden ser un valor negativo")
	}
	if req.GolesVisitante < 0 {
		return errors.New("los goles del equipo visitante no pueden ser un valor negativo")
	}
	if strings.TrimSpace(req.Estado) == "" {
		return errors.New("el estado del partido no puede estar vacío")
	}
	return nil
}
