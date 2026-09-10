package handlers

import (
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

// Reemplaza el struct HTTP por este para procesar partidos
type CreatePartidoRequest struct {
	EquipoLocalID     int32    `json:"equipo_local_id"`
	EquipoVisitanteID int32    `json:"equipo_visitante_id"`
	Fecha             DateOnly `json:"fecha"` // Lee la cadena "YYYY-MM-DD"
	Estado            string   `json:"estado"`
}

type UpdateResultadoPartidoRequest struct {
	GolesLocal     int16  `json:"goles_local"`
	GolesVisitante int16  `json:"goles_visitante"`
	Estado         string `json:"estado"` // PENDIENTE, JUGANDO, FINALIZADO
}
