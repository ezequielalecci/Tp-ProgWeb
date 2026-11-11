#!/usr/bin/env bash
# Prueba las paginas renderizadas en el servidor y el alta por formularios (Post/Redirect/Get).
# Requiere los datos de test-data.sql (equipos 1 y 2, jugadores 1 y 2, partidos 1 y 2).
BASE_URL=${BASE_URL:-http://localhost:8080}
TMP=$(mktemp); HDR=$(mktemp); trap 'rm -f "$TMP" "$HDR"' EXIT
PASS=0; FAIL=0; STATUS=""; BODY=""
G=$'\e[32m'; R=$'\e[31m'; B=$'\e[1m'; N=$'\e[0m'

# req "nombre" "codigos esperados (regex)" METODO /ruta [campo=valor ...]
# Los campos se envian como un <form>: application/x-www-form-urlencoded.
# No se siguen los redirects, asi se puede verificar el 303 y su Location.
req() {
  local name=$1 expected=$2 method=$3 path=$4; shift 4
  local args=(-s -o "$TMP" -D "$HDR" -w '%{http_code}' -X "$method" "$BASE_URL$path")
  local kv; for kv in "$@"; do args+=(--data-urlencode "$kv"); done
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
# redirige "nombre" "/destino" : verifica el header Location de la ultima respuesta
redirige() {
  local loc; loc=$(grep -i '^location:' "$HDR" | tr -d '\r' | sed 's/^[^:]*: *//')
  if [[ $loc == "$2" ]]; then PASS=$((PASS + 1)); echo "  ${G}PASS${N} $1 (Location: $loc)"
  else FAIL=$((FAIL + 1)); echo "  ${R}FAIL${N} $1 (Location: '${loc}', esperado '$2')"; fi
}
section() { echo; echo "${B}== $1 ==${N}"; }

HOY=$(date +%F)

section "PAGINAS"
req "Pagina principal" 200 GET /;                          has "Es HTML completo" "</html>"
has "Lista a Real Madrid" "Real Madrid"
has "Tiene el formulario de alta" 'action="/equipos"'
req "Pestania jugadores" 200 GET "/?tab=jugadores";        has "Lista a Vinicius" "Vinicius"
req "Pestania partidos" 200 GET "/?tab=partidos";          has "Muestra el estado" "FINALIZADO"
req "Hoja de estilos" 200 GET /static/styles.css
req "Ruta inexistente" 404 GET /no-existe
req "La API JSON ya no existe" 404 GET /api/equipos

section "EQUIPOS"
req "Crear equipo sin nombre (validacion)" 400 POST /equipos "nombre=  " "formacion_actual=3-5-2" "valoracion=80" "posicion_tabla=3"
has "Muestra el error" "no puede estar vac"
has "Conserva lo que se escribio" 'value="3-5-2"'
req "Crear equipo con valoracion no numerica" 400 POST /equipos "nombre=X" "formacion_actual=4-3-3" "valoracion=abc" "posicion_tabla=3"
req "Crear AC Milan" 303 POST /equipos "nombre=AC Milan" "formacion_actual=4-2-3-1" "escudo_url=https://example.com/milan.svg" "valoracion=84" "posicion_tabla=3"
redirige "Vuelve a la pagina principal" "/"
req "Lista actualizada" 200 GET /;                         has "Incluye AC Milan" "AC Milan"
req "Crear AC Milan repetido (nombre unico)" 400 POST /equipos "nombre=AC Milan" "formacion_actual=4-2-3-1" "valoracion=84" "posicion_tabla=3"

# El id del Milan sale del <select> de equipos del formulario de jugadores.
req "Formulario de jugadores" 200 GET "/?tab=jugadores"
MILAN=$(grep -oE '<option value="[0-9]+"[^>]*>AC Milan' <<<"$BODY" | head -1 | grep -oE '[0-9]+' | head -1)
echo "     id AC Milan = ${MILAN:-?}"
[[ -z $MILAN ]] && { echo "${R}No se pudo obtener el id del Milan; abortando${N}"; exit 1; }

section "JUGADORES"
req "Crear jugador con fecha invalida (validacion)" 400 POST /jugadores "equipo_id=1" "nombre=X" "posicion=DEL" "fecha_nacimiento=12/07/2000" "altura=1.80" "media_general=80" "ritmo=1" "tiro=1" "pase=1" "regate=1" "defensa=1" "fisico=1"
req "Crear jugador sin equipo (validacion)" 400 POST /jugadores "equipo_id=" "nombre=X" "posicion=DEL" "fecha_nacimiento=2000-07-12" "altura=1.80" "media_general=80" "ritmo=1" "tiro=1" "pase=1" "regate=1" "defensa=1" "fisico=1"
req "Crear Luka Modric" 303 POST /jugadores "equipo_id=$MILAN" "nombre=Luka Modric" "posicion=MED" "fecha_nacimiento=1985-09-09" "altura=1.72" "media_general=85" "ritmo=70" "tiro=76" "pase=88" "regate=86" "defensa=72" "fisico=66" "foto_url=https://example.com/fotos/luka-modric.jpg"
redirige "Vuelve a la pestania de jugadores" "/?tab=jugadores"
req "Lista actualizada" 200 GET "/?tab=jugadores";         has "Incluye a Modric" "Luka Modric"

section "PARTIDOS"
req "Crear partido con mismo equipo (validacion)" 400 POST /partidos "equipo_local_id=2" "equipo_visitante_id=2" "fecha=$HOY" "estado=PENDIENTE"
req "Crear partido Barcelona vs Milan (en curso)" 303 POST /partidos "equipo_local_id=2" "equipo_visitante_id=$MILAN" "fecha=$HOY" "estado=EN CURSO"
redirige "Vuelve a la pestania de partidos" "/?tab=partidos"
req "Lista actualizada" 200 GET "/?tab=partidos";          has "Incluye el partido en curso" "EN CURSO"

echo; echo "${B}Resultado: $PASS correctas, $FAIL fallidas${N}"
[[ $FAIL -eq 0 ]]
