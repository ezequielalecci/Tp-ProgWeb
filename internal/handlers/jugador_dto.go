package handlers

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

type UpdateJugadorStatsRequest struct {
	ID          int32 `json:"id"`
	Goles       int16 `json:"goles"`
	Asistencias int16 `json:"asistencias"`
}
