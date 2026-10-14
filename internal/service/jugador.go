package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	database "tpweb/db/sqlc"
	"tpweb/internal/dto"
)

var (
	ErrJugadorNotFound = errors.New("jugador no encontrado")
	ErrFechaInvalida   = errors.New("formato de fecha inválido (debe ser AAAA-MM-DD)")
)

type JugadorService struct {
	queries *database.Queries
}

func NewJugadorService(q *database.Queries) *JugadorService {
	return &JugadorService{
		queries: q,
	}
}

// GetJugador obtiene un jugador por su ID
func (s *JugadorService) GetJugador(ctx context.Context, id int32) (database.Jugadore, error) {
	jugador, err := s.queries.GetJugador(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Jugadore{}, ErrJugadorNotFound
		}
		return database.Jugadore{}, err
	}
	return jugador, nil
}

// GetJugadores retorna todos los jugadores
func (s *JugadorService) GetJugadores(ctx context.Context) ([]database.Jugadore, error) {
	jugadores, err := s.queries.GetJugadores(ctx)
	if err != nil {
		return nil, err
	}
	if jugadores == nil {
		return []database.Jugadore{}, nil
	}
	return jugadores, nil
}

// ListJugadoresByEquipo lista todos los jugadores pertenecientes a un equipo
func (s *JugadorService) ListJugadoresByEquipo(ctx context.Context, equipoID int32) ([]database.Jugadore, error) {
	jugadores, err := s.queries.ListJugadoresByEquipo(ctx, equipoID)
	if err != nil {
		return nil, err
	}
	if jugadores == nil {
		return []database.Jugadore{}, nil
	}
	return jugadores, nil
}

// CreateJugador valida existencia de equipo, parsea la fecha y crea el jugador
func (s *JugadorService) CreateJugador(ctx context.Context, req dto.CreateJugadorRequest) (database.Jugadore, error) {
	// Parsear fecha
	fechaNac, err := time.Parse("2006-01-02", req.FechaNacimiento)
	if err != nil {
		return database.Jugadore{}, ErrFechaInvalida
	}

	// Verificar si el equipo asignado existe
	_, err = s.queries.GetEquipo(ctx, req.EquipoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Jugadore{}, ErrEquipoNotFound
		}
		return database.Jugadore{}, err
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

	return s.queries.CreateJugador(ctx, params)
}

// UpdateJugadorStats actualiza goles y asistencias del jugador verificando existencia previa
func (s *JugadorService) UpdateJugadorStats(ctx context.Context, id int32, req dto.UpdateJugadorStatsRequest) (database.Jugadore, error) {
	// Verificar existencia previa
	_, err := s.queries.GetJugador(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Jugadore{}, ErrJugadorNotFound
		}
		return database.Jugadore{}, err
	}

	params := database.UpdateJugadorStatsParams{
		ID:          id,
		Goles:       req.Goles,
		Asistencias: req.Asistencias,
	}

	return s.queries.UpdateJugadorStats(ctx, params)
}

// DeleteJugador elimina un jugador tras validar su existencia
func (s *JugadorService) DeleteJugador(ctx context.Context, id int32) error {
	_, err := s.queries.GetJugador(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrJugadorNotFound
		}
		return err
	}

	return s.queries.DeleteJugador(ctx, id)
}

// GetMejorJugadorByEquipo consulta el jugador con mejor valoración de un equipo
func (s *JugadorService) GetMejorJugadorByEquipo(ctx context.Context, equipoID int32) (database.Jugadore, error) {
	jugador, err := s.queries.GetMejorJugadorByEquipo(ctx, equipoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Jugadore{}, ErrJugadorNotFound
		}
		return database.Jugadore{}, err
	}
	return jugador, nil
}
