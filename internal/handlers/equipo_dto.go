package handlers

type CreateEquipoRequest struct {
	Nombre          string  `json:"nombre"`
	FormacionActual string  `json:"formacion_actual"`
	EscudoUrl       *string `json:"escudo_url"`
	Valoracion      int16   `json:"valoracion"`
	PosicionTabla   int16   `json:"posicion_tabla"`
}

type UpdateEquipoRequest struct {
	Nombre          string  `json:"nombre"`
	FormacionActual string  `json:"formacion_actual"`
	EscudoUrl       *string `json:"escudo_url"`
	MejorJugador    *int32  `json:"mejor_jugador"`
	Valoracion      int16   `json:"valoracion"`
	PosicionTabla   int16   `json:"posicion_tabla"`
}
