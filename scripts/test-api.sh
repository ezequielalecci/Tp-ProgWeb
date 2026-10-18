#!/usr/bin/env bash
# Prueba todos los endpoints de equipos, jugadores y partidos.
# Requiere los datos de test-data.sql (equipos 1 y 2, jugadores 1 y 2, partidos 1 y 2).
BASE_URL=${BASE_URL:-http://localhost:8080}
TMP=$(mktemp); trap 'rm -f "$TMP"' EXIT
PASS=0; FAIL=0; STATUS=""; BODY=""
G=$'\e[32m'; R=$'\e[31m'; B=$'\e[1m'; N=$'\e[0m'

# req "nombre" "codigos esperados (regex)" METODO /ruta ['cuerpo json']
req() {
  local name=$1 expected=$2 method=$3 path=$4 body=${5:-}
  local args=(-s -o "$TMP" -w '%{http_code}' -X "$method" "$BASE_URL$path")
  [[ -n $body ]] && args+=(-H 'Content-Type: application/json' -d "$body")
  STATUS=$(curl "${args[@]}"); BODY=$(cat "$TMP")
  if [[ $STATUS =~ ^($expected)$ ]]; then
    PASS=$((PASS + 1)); echo "  ${G}PASS${N} $name ($method $path -> $STATUS)"
  else
    FAIL=$((FAIL + 1)); echo "  ${R}FAIL${N} $name ($method $path -> $STATUS, esperado $expected)"
    echo "       respuesta: ${BODY:0:300}"
  fi
}
# has "nombre" "texto" : verifica que la ultima respuesta contenga el texto
has() {
  if [[ $BODY == *"$2"* ]]; then PASS=$((PASS + 1)); echo "  ${G}PASS${N} $1 (contiene '$2')"
  else FAIL=$((FAIL + 1)); echo "  ${R}FAIL${N} $1 (no contiene '$2')"; echo "       respuesta: ${BODY:0:300}"; fi
}
json_id() { grep -oiE '"id" *: *[0-9]+' <<<"$BODY" | head -1 | grep -oE '[0-9]+$'; }
section() { echo; echo "${B}== $1 ==${N}"; }

HOY=$(date +%F)
FOTO_LARGA="https://example.com/fotos/luka-modric.jpg?token=$(printf 'a%.0s' {1..300})"   # >255 caracteres

section "EQUIPOS"
req "Listar equipos" 200 GET /api/equipos;                has "Incluye Real Madrid" "Real Madrid"
req "Obtener equipo 1" 200 GET /api/equipos/1;            has "Es Real Madrid" "Real Madrid"
req "Equipo inexistente" 404 GET /api/equipos/999999
req "Crear equipo sin nombre (validacion)" 400 POST /api/equipos '{"nombre":"  ","formacion_actual":"4-3-3","escudo_url":null,"valoracion":80,"posicion_tabla":3}'
req "Crear AC Milan" '200|201' POST /api/equipos '{"nombre":"AC Milan","formacion_actual":"4-2-3-1","escudo_url":"https://example.com/milan.svg","valoracion":84,"posicion_tabla":3}'
MILAN=$(json_id); echo "     id AC Milan = ${MILAN:-?}"
[[ -z $MILAN ]] && { echo "${R}No se pudo obtener el id del Milan; abortando${N}"; exit 1; }
req "Obtener AC Milan" 200 GET "/api/equipos/$MILAN";     has "Es AC Milan" "AC Milan"
req "Actualizar AC Milan" '200|204' PUT "/api/equipos/$MILAN" '{"nombre":"AC Milan","formacion_actual":"4-4-2","escudo_url":"https://example.com/milan.svg","mejor_jugador":null,"valoracion":85,"posicion_tabla":3}'
req "Actualizar equipo invalido (posicion 0)" 400 PUT "/api/equipos/$MILAN" '{"nombre":"AC Milan","formacion_actual":"4-4-2","escudo_url":null,"mejor_jugador":null,"valoracion":85,"posicion_tabla":0}'

section "JUGADORES"
req "Listar jugadores" 200 GET /api/jugadores;            has "Incluye Vinicius" "Vinicius"
req "Obtener jugador 1" 200 GET /api/jugadores/1;         has "Es Vinicius" "Vinicius"
req "Jugador inexistente" 404 GET /api/jugadores/999999
req "Crear jugador con fecha invalida (validacion)" 400 POST /api/jugadores '{"equipo_id":1,"nombre":"X","posicion":"DEL","fecha_nacimiento":"12/07/2000","media_general":80,"altura":"1.80","ritmo":1,"tiro":1,"pase":1,"regate":1,"defensa":1,"fisico":1,"foto_url":null}'
req "Crear Luka Modric (foto >255 caracteres)" '200|201' POST /api/jugadores "{\"equipo_id\":$MILAN,\"nombre\":\"Luka Modric\",\"posicion\":\"MED\",\"fecha_nacimiento\":\"1985-09-09\",\"media_general\":85,\"altura\":\"1.72\",\"ritmo\":70,\"tiro\":76,\"pase\":88,\"regate\":86,\"defensa\":72,\"fisico\":66,\"foto_url\":\"$FOTO_LARGA\"}"
MODRIC=$(json_id); echo "     id Modric = ${MODRIC:-?}"
[[ -z $MODRIC ]] && { echo "${R}No se pudo obtener el id de Modric; abortando${N}"; exit 1; }
req "Obtener Modric" 200 GET "/api/jugadores/$MODRIC";    has "Es Luka Modric" "Luka Modric"
req "Jugadores por equipo (Milan)" 200 GET "/api/jugadores/equipo?equipo_id=$MILAN"; has "Incluye Modric" "Luka Modric"
req "Actualizar stats de Modric" '200|204' PUT "/api/jugadores/$MODRIC/stats" '{"goles":3,"asistencias":5}'
req "Stats negativas (validacion)" 400 PUT "/api/jugadores/$MODRIC/stats" '{"goles":-1,"asistencias":0}'

section "MEJOR JUGADOR"
req "Asignar a Modric como mejor jugador del Milan" '200|204' PUT "/api/equipos/$MILAN" "{\"nombre\":\"AC Milan\",\"formacion_actual\":\"4-4-2\",\"escudo_url\":\"https://example.com/milan.svg\",\"mejor_jugador\":$MODRIC,\"valoracion\":85,\"posicion_tabla\":3}"
req "Mejor jugador del Milan" 200 GET "/api/equipos/$MILAN/mejor-jugador"; has "Es Luka Modric" "Luka Modric"
req "Mejor jugador del Real Madrid" 200 GET /api/equipos/1/mejor-jugador;  has "Es Vinicius" "Vinicius"

section "PARTIDOS"
req "Listar partidos" 200 GET /api/partidos
req "Obtener partido 1" 200 GET /api/partidos/1
req "Partido inexistente" 404 GET /api/partidos/999999
req "Crear partido con mismo equipo (validacion)" 400 POST /api/partidos "{\"equipo_local_id\":2,\"equipo_visitante_id\":2,\"fecha\":\"$HOY\",\"estado\":\"PENDIENTE\"}"
req "Crear partido Barcelona vs Milan (en curso)" '200|201' POST /api/partidos "{\"equipo_local_id\":2,\"equipo_visitante_id\":$MILAN,\"fecha\":\"$HOY\",\"estado\":\"EN CURSO\"}"
PARTIDO=$(json_id); echo "     id partido = ${PARTIDO:-?}"
[[ -z $PARTIDO ]] && { echo "${R}No se pudo obtener el id del partido; abortando${N}"; exit 1; }
req "Obtener partido nuevo" 200 GET "/api/partidos/$PARTIDO"; has "Estado EN CURSO" "EN CURSO"
req "Actualizar resultado 1-1 (en curso)" '200|204' PUT "/api/partidos/$PARTIDO/resultado" '{"goles_local":1,"goles_visitante":1,"estado":"EN CURSO"}'
req "Resultado negativo (validacion)" 400 PUT "/api/partidos/$PARTIDO/resultado" '{"goles_local":-1,"goles_visitante":0,"estado":"EN CURSO"}'
req "Ultimo partido del Milan" 200 GET "/api/equipos/$MILAN/partidos/ultimo"
req "Ultimo partido del Barcelona" 200 GET /api/equipos/2/partidos/ultimo
req "Ultimo finalizado del Barcelona" 200 GET /api/equipos/2/partidos/ultimo-finalizado
req "Ultimo con detalle del Barcelona" 200 GET /api/equipos/2/partidos/ultimo-detalle

section "ELIMINACION"
req "Quitar mejor jugador del Milan" '200|204' PUT "/api/equipos/$MILAN" '{"nombre":"AC Milan","formacion_actual":"4-4-2","escudo_url":null,"mejor_jugador":null,"valoracion":85,"posicion_tabla":3}'
req "Eliminar partido" '200|204' DELETE "/api/partidos/$PARTIDO"
req "Partido eliminado ya no existe" 404 GET "/api/partidos/$PARTIDO"
req "Eliminar jugador Modric" '200|204' DELETE "/api/jugadores/$MODRIC"
req "Jugador eliminado ya no existe" 404 GET "/api/jugadores/$MODRIC"
req "Eliminar equipo AC Milan" '200|204' DELETE "/api/equipos/$MILAN"
req "Equipo eliminado ya no existe" 404 GET "/api/equipos/$MILAN"

echo; echo "${B}Resultado: $PASS correctas, $FAIL fallidas${N}"
[[ $FAIL -eq 0 ]]