package service

import (
	"context"
	"database/sql"
	"errors"

	database "tpweb/db/sqlc"
	"tpweb/internal/dto"
)

var (
	ErrPartidoNotFound     = errors.New("partido no encontrado")
	ErrEquipoLocalNotFound = errors.New("el equipo local especificado no existe")
	ErrEquipoVisitNotFound = errors.New("el equipo visitante especificado no existe")
)

type PartidoService struct {
	queries *database.Queries
}

func NewPartidoService(q *database.Queries) *PartidoService {
	return &PartidoService{
		queries: q,
	}
}

// GetPartido obtiene un partido por su ID
func (s *PartidoService) GetPartido(ctx context.Context, id int32) (database.Partido, error) {
	partido, err := s.queries.GetPartido(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Partido{}, ErrPartidoNotFound
		}
		return database.Partido{}, err
	}
	return partido, nil
}

// ListPartidos obtiene todos los partidos
func (s *PartidoService) ListPartidos(ctx context.Context) ([]database.Partido, error) {
	partidos, err := s.queries.ListPartidos(ctx)
	if err != nil {
		return nil, err
	}
	if partidos == nil {
		return []database.Partido{}, nil
	}
	return partidos, nil
}

func (s *PartidoService) CrearPartido(ctx context.Context, req dto.CreatePartidoRequest) (database.Partido, error) {
	// 1. Verificar existencia del equipo local
	_, err := s.queries.GetEquipo(ctx, req.EquipoLocalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Partido{}, ErrEquipoLocalNotFound
		}
		return database.Partido{}, err
	}

	// 2. Verificar existencia del equipo visitante
	_, err = s.queries.GetEquipo(ctx, req.EquipoVisitanteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Partido{}, ErrEquipoVisitNotFound
		}
		return database.Partido{}, err
	}

	// 3. Mapear a los parámetros requeridos por sqlc (usando sql.NullInt32)
	params := database.CreatePartidoParams{
		EquipoLocalID: sql.NullInt32{
			Int32: req.EquipoLocalID,
			Valid: true,
		},
		EquipoVisitanteID: sql.NullInt32{
			Int32: req.EquipoVisitanteID,
			Valid: true,
		},
		Fecha:  req.Fecha.Time,
		Estado: req.Estado,
	}

	return s.queries.CreatePartido(ctx, params)
}

// UpdateResultadoPartido actualiza marcadores y estado previa verificación
func (s *PartidoService) UpdateResultadoPartido(ctx context.Context, id int32, req dto.UpdateResultadoPartidoRequest) (database.Partido, error) {
	_, err := s.queries.GetPartido(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Partido{}, ErrPartidoNotFound
		}
		return database.Partido{}, err
	}

	params := database.UpdateResultadoPartidoParams{
		ID:             id,
		GolesLocal:     req.GolesLocal,
		GolesVisitante: req.GolesVisitante,
		Estado:         req.Estado,
	}

	return s.queries.UpdateResultadoPartido(ctx, params)
}

// DeletePartido elimina un partido verificando su existencia
func (s *PartidoService) DeletePartido(ctx context.Context, id int32) error {
	_, err := s.queries.GetPartido(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPartidoNotFound
		}
		return err
	}

	return s.queries.DeletePartido(ctx, id)
}

// GetUltimoPartidoByEquipo
func (s *PartidoService) GetUltimoPartidoByEquipo(ctx context.Context, equipoID int32) (database.Partido, error) {
	equipoIDNull := sql.NullInt32{Int32: equipoID, Valid: true}

	partido, err := s.queries.GetUltimoPartidoByEquipo(ctx, equipoIDNull)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Partido{}, ErrPartidoNotFound
		}
		return database.Partido{}, err
	}
	return partido, nil
}

// GetUltimoPartidoFinalizadoByEquipo
func (s *PartidoService) GetUltimoPartidoFinalizadoByEquipo(ctx context.Context, equipoID int32) (database.Partido, error) {
	equipoIDNull := sql.NullInt32{Int32: equipoID, Valid: true}

	partido, err := s.queries.GetUltimoPartidoFinalizadoByEquipo(ctx, equipoIDNull)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Partido{}, ErrPartidoNotFound
		}
		return database.Partido{}, err
	}
	return partido, nil
}

// GetUltimoPartidoConDetalleByEquipo obtiene el detalle completo del último partido
func (s *PartidoService) GetUltimoPartidoConDetalleByEquipo(ctx context.Context, equipoID int32) (database.GetUltimoPartidoConDetalleByEquipoRow, error) {
	// Convertir int32 a sql.NullInt32
	equipoIDNull := sql.NullInt32{
		Int32: equipoID,
		Valid: true,
	}

	partidoDetalle, err := s.queries.GetUltimoPartidoConDetalleByEquipo(ctx, equipoIDNull)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.GetUltimoPartidoConDetalleByEquipoRow{}, ErrPartidoNotFound
		}
		return database.GetUltimoPartidoConDetalleByEquipoRow{}, err
	}
	return partidoDetalle, nil
}
