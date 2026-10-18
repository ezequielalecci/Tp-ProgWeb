#!/usr/bin/env bash
# Espera a que la base este healthy y a que la API responda.
DB_CONTAINER=${DB_CONTAINER:-gestion_futbol_db}
BASE_URL=${BASE_URL:-http://localhost:8080}
TIMEOUT=${TIMEOUT:-90}

echo ">> Esperando a la base de datos ($DB_CONTAINER)..."
for ((i = 0; i < TIMEOUT; i++)); do
  health=$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$DB_CONTAINER" 2>/dev/null || echo "missing")
  case "$health" in
    healthy) echo "   base lista"; break ;;
    none) docker exec "$DB_CONTAINER" pg_isready >/dev/null 2>&1 && { echo "   base lista"; break; } ;;
    unhealthy) echo "!! La base quedo unhealthy"; docker logs --tail=40 "$DB_CONTAINER"; exit 1 ;;
  esac
  sleep 1
done
[[ $i -ge $TIMEOUT ]] && { echo "!! Timeout esperando la base"; exit 1; }

echo ">> Esperando a la aplicacion ($BASE_URL)..."
for ((i = 0; i < TIMEOUT; i++)); do
  code=$(curl -s -o /dev/null -w '%{http_code}' "$BASE_URL/api/equipos" 2>/dev/null)
  [[ "$code" == "200" ]] && { echo "   aplicacion lista"; exit 0; }
  sleep 1
done
echo "!! Timeout esperando la aplicacion (ultimo codigo: $code)"
docker compose logs --tail=40
exit 1