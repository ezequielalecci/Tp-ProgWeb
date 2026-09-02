-- name: CreateEquipo :one
INSERT INTO equipos (
    nombre,
    formacion_actual,
    escudo_url
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetEquipoByID :one
SELECT 
    id,
    nombre,
    formacion_actual,
    escudo_url
FROM equipos
WHERE id = $1 LIMIT 1;

-- name: ListEquipos :many
SELECT 
    id,
    nombre,
    formacion_actual,
    escudo_url
FROM equipos
ORDER BY nombre ASC;

-- name: UpdateEquipo :one
UPDATE equipos
SET 
    nombre = $1,
    formacion_actual = $2,
    escudo_url = $3
WHERE id = $4
RETURNING *;

-- name: DeleteEquipo :exec
DELETE FROM equipos
WHERE id = $1;