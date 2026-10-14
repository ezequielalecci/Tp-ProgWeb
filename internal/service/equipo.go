package service

import (
	"context"
	"database/sql"
	"errors"

	database "tpweb/db/sqlc"
	"tpweb/internal/dto"
)

// ErrEquipoNotFound es un error de dominio reutilizable
var ErrEquipoNotFound = errors.New("equipo no encontrado")

// EquipoService encapsula la interacción con la base de datos y la lógica de negocio
type EquipoService struct {
	queries *database.Queries
}

// NewEquipoService constructor del servicio
func NewEquipoService(q *database.Queries) *EquipoService {
	return &EquipoService{
		queries: q,
	}
}

// GetEquipos lista todos los equipos
func (s *EquipoService) GetEquipos(ctx context.Context) ([]database.Equipo, error) {
	equipos, err := s.queries.ListEquipos(ctx)
	if err != nil {
		return nil, err
	}
	if equipos == nil {
		return []database.Equipo{}, nil
	}
	return equipos, nil
}

// GetEquipoByID obtiene un equipo por su ID
func (s *EquipoService) GetEquipoByID(ctx context.Context, id int32) (database.Equipo, error) {
	equipo, err := s.queries.GetEquipo(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Equipo{}, ErrEquipoNotFound
		}
		return database.Equipo{}, err
	}
	return equipo, nil
}

// CrearEquipo procesa la creación de un nuevo equipo mapeando el DTO a los tipos de sqlc
func (s *EquipoService) CrearEquipo(ctx context.Context, req dto.CreateEquipoRequest) (database.Equipo, error) {
	params := database.CreateEquipoParams{
		Nombre:          req.Nombre,
		FormacionActual: req.FormacionActual,
		Valoracion:      req.Valoracion,
		PosicionTabla:   req.PosicionTabla,
	}

	if req.EscudoUrl != nil {
		params.EscudoUrl = sql.NullString{String: *req.EscudoUrl, Valid: true}
	}

	return s.queries.CreateEquipo(ctx, params)
}

// UpdateEquipo verifica existencia previa y actualiza los datos
func (s *EquipoService) UpdateEquipo(ctx context.Context, id int32, req dto.UpdateEquipoRequest) (database.Equipo, error) {
	// 1. Verificar existencia previa
	_, err := s.queries.GetEquipo(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Equipo{}, ErrEquipoNotFound
		}
		return database.Equipo{}, err
	}

	// 2. Mapear DTO a parámetros requeridos por sqlc
	params := database.UpdateEquipoParams{
		ID:              id,
		Nombre:          req.Nombre,
		FormacionActual: req.FormacionActual,
		Valoracion:      req.Valoracion,
		PosicionTabla:   req.PosicionTabla,
	}

	if req.EscudoUrl != nil {
		params.EscudoUrl = sql.NullString{String: *req.EscudoUrl, Valid: true}
	}

	if req.MejorJugador != nil {
		params.MejorJugador = sql.NullInt32{Int32: *req.MejorJugador, Valid: true}
	}

	return s.queries.UpdateEquipo(ctx, params)
}

func (s *EquipoService) DeleteEquipo(ctx context.Context, id int32) error {
	// Verificar existencia del equipo
	_, err := s.queries.GetEquipo(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEquipoNotFound
		}
		return err
	}

	// 2. Convertir id a sql.NullInt32 para la consulta en partidos
	equipoIDNull := sql.NullInt32{
		Int32: id,
		Valid: true,
	}

	// 3. Eliminar ÚNICAMENTE los partidos PENDIENTES (no finalizados)
	// Query SQL: DELETE FROM partidos WHERE (equipo_local_id = $1 OR equipo_visitante_id = $1) AND estado != 'FINALIZADO';
	err = s.queries.DeletePartidosPendientesByEquipo(ctx, equipoIDNull)
	if err != nil {
		return err
	}

	//  Eliminar el equipo.
	// Los partidos FINALIZADOS permanecen en la BD; PostgreSQL pondrá equipo_local_id / equipo_visitante_id en NULL.
	return s.queries.DeleteEquipo(ctx, id)
}
