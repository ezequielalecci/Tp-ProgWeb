#!/usr/bin/env bash
set -e

echo "=== 1. Limpiando contenedores y volumenes previos ==="
docker compose down -v --remove-orphans

cleanup() {
    echo "=== 4. Tareas posteriores: Borrando contenedores y volumenes ==="
    docker compose down -v --remove-orphans
}
trap cleanup EXIT

echo "=== 2. Regenerando codigo con sqlc ==="
if command -v sqlc &> /dev/null; then
    sqlc generate
else
    echo "sqlc no esta instalado localmente, omitiendo ejecucion de sqlc generate..."
fi

echo "=== 3. Levantando base de datos PostgreSQL ==="
docker compose up -d postgres

echo "Esperando a que la base de datos este lista..."
until [ "$(docker inspect --format='{{json .State.Health.Status}}' gestion_futbol_db 2>/dev/null)" == "\"healthy\"" ]; do
    sleep 1
done

echo "Base de datos lista."

echo "=== Ejecutando Go Tests ==="
TEST_DB_URL="postgres://postgres:secret@localhost:5432/gestion_futbol?sslmode=disable" \
go test -v ./...