-- name: GetEquipo :one
SELECT * FROM "equipos"
WHERE "id" = $1 LIMIT 1;

-- name: ListEquipos :many
SELECT * FROM "equipos"
ORDER BY "posicion_tabla" ASC;

-- name: CreateEquipo :one
INSERT INTO "equipos" (
    "nombre",
    "formacion_actual",
    "escudo_url",
    "valoracion",
    "posicion_tabla"
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateEquipo :one
UPDATE "equipos"
SET 
    "nombre" = COALESCE($2, "nombre"),
    "formacion_actual" = COALESCE($3, "formacion_actual"),
    "escudo_url" = COALESCE($4, "escudo_url"),
    "mejor_jugador" = COALESCE($5, "mejor_jugador"),
    "valoracion" = COALESCE($6, "valoracion"),
    "posicion_tabla" = COALESCE($7, "posicion_tabla")
WHERE "id" = $1
RETURNING *;

-- name: DeleteEquipo :exec
DELETE FROM "equipos"
WHERE "id" = $1;


-- name: GetUltimoPartidoByEquipo :one
SELECT * FROM "partidos"
WHERE "equipo_local_id" = $1 OR "equipo_visitante_id" = $1
ORDER BY "fecha" DESC, "id" DESC
LIMIT 1;

-- name: GetUltimoPartidoFinalizadoByEquipo :one
SELECT * FROM "partidos"
WHERE ("equipo_local_id" = $1 OR "equipo_visitante_id" = $1)
  AND "estado" = 'FINALIZADO'
ORDER BY "fecha" DESC, "id" DESC
LIMIT 1;

-- name: GetUltimoPartidoConDetalleByEquipo :one
SELECT 
    p.id,
    p.fecha,
    p.estado,
    p.goles_local,
    p.goles_visitante,
    el.id AS local_id,
    el.nombre AS local_nombre,
    el.escudo_url AS local_escudo,
    ev.id AS visitante_id,
    ev.nombre AS visitante_nombre,
    ev.escudo_url AS visitante_escudo
FROM "partidos" p
JOIN "equipos" el ON p."equipo_local_id" = el."id"
JOIN "equipos" ev ON p."equipo_visitante_id" = ev."id"
WHERE p."equipo_local_id" = $1 OR p."equipo_visitante_id" = $1
ORDER BY p."fecha" DESC, p."id" DESC
LIMIT 1;

-- name: GetJugador :one
SELECT * FROM "jugadores"
WHERE "id" = $1 LIMIT 1;

-- name: ListJugadoresByEquipo :many
SELECT * FROM "jugadores"
WHERE "equipo_id" = $1
ORDER BY "media_general" DESC;

-- name: CreateJugador :one
INSERT INTO "jugadores" (
    "equipo_id",
    "nombre",
    "posicion",
    "fecha_nacimiento",
    "media_general",
    "altura",
    "ritmo",
    "tiro",
    "pase",
    "regate",
    "defensa",
    "fisico",
    "foto_url"
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)

RETURNING *;

-- name: UpdateJugadorStats :one
UPDATE "jugadores"
SET 
    "goles" = "goles" + $2,
    "asistencias" = "asistencias" + $3
WHERE "id" = $1
RETURNING *;

-- name: DeleteJugador :exec
DELETE FROM "jugadores"
WHERE "id" = $1;

-- name: GetMejorJugadorByEquipo :one
SELECT * FROM "jugadores"
WHERE "equipo_id" = $1
ORDER BY "media_general" DESC, "id" ASC
LIMIT 1;

-- name: GetPartido :one
SELECT * FROM "partidos"
WHERE "id" = $1 LIMIT 1;

-- name: ListPartidos :many
SELECT * FROM "partidos"
ORDER BY "fecha" DESC;

-- name: CreatePartido :one
INSERT INTO "partidos" (
    "equipo_local_id",
    "equipo_visitante_id",
    "fecha",
    "estado"
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateResultadoPartido :one
UPDATE "partidos"
SET 
    "goles_local" = $2,
    "goles_visitante" = $3,
    "estado" = $4
WHERE "id" = $1
RETURNING *;

-- name: DeletePartido :exec
DELETE FROM "partidos"
WHERE "id" = $1;