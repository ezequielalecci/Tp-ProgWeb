package handlers_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	database "tpweb/db/sqlc"
	"tpweb/server/routes"

	_ "github.com/lib/pq"
)

// Estas pruebas usan la base real (la de docker compose) con los datos de
// db/test-data.sql: equipos 1 y 2 ya cargados.
func setupTest(t *testing.T) (http.Handler, *sql.DB) {
	t.Helper()
	dbURL := os.Getenv("TEST_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:secret@localhost:5432/gestion_futbol?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Error de conexión a la base de datos de test: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("No se pudo conectar a PostgreSQL en test: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	return routes.RegisterRoutes(database.New(db)), db
}

func get(t *testing.T, router http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
	return rr
}

// postForm envía los datos igual que un <form method="post"> del navegador.
func postForm(t *testing.T, router http.Handler, path string, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func esperarRedirect(t *testing.T, rr *httptest.ResponseRecorder, destino string) {
	t.Helper()
	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status %d, esperaba 303 (body: %s)", rr.Code, rr.Body.String())
	}
	if loc := rr.Header().Get("Location"); loc != destino {
		t.Errorf("Location %q, esperaba %q", loc, destino)
	}
}

func TestIndexRenderizaHTML(t *testing.T) {
	router, _ := setupTest(t)

	for _, path := range []string{"/", "/?tab=jugadores", "/?tab=partidos"} {
		rr := get(t, router, path)
		if rr.Code != http.StatusOK {
			t.Errorf("GET %s devolvió status %d, esperaba 200", path, rr.Code)
		}
		if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
			t.Errorf("GET %s devolvió Content-Type %q, esperaba text/html", path, ct)
		}
		if body := rr.Body.String(); !strings.Contains(body, "<form") || !strings.Contains(body, "</html>") {
			t.Errorf("GET %s no devolvió la página completa con el formulario", path)
		}
	}
}

func TestIndexListaEquipos(t *testing.T) {
	router, _ := setupTest(t)

	rr := get(t, router, "/")
	if !strings.Contains(rr.Body.String(), "Real Madrid") {
		t.Errorf("GET / no muestra los equipos de la base")
	}
}

func TestCrearEquipoRedirigeYAparece(t *testing.T) {
	router, db := setupTest(t)

	nombre := fmt.Sprintf("Equipo de prueba %d", time.Now().UnixNano())
	t.Cleanup(func() { db.Exec(`DELETE FROM equipos WHERE nombre = $1`, nombre) })

	rr := postForm(t, router, "/equipos", url.Values{
		"nombre":           {nombre},
		"formacion_actual": {"4-4-2"},
		"valoracion":       {"84"},
		"posicion_tabla":   {"3"},
		"escudo_url":       {""},
	})
	esperarRedirect(t, rr, "/")

	// El navegador seguiría el redirect con un GET: ahí tiene que estar el equipo nuevo.
	if body := get(t, router, "/").Body.String(); !strings.Contains(body, nombre) {
		t.Errorf("el equipo creado no aparece en la lista")
	}
}

func TestCrearEquipoInvalidoNoRedirige(t *testing.T) {
	router, _ := setupTest(t)

	rr := postForm(t, router, "/equipos", url.Values{
		"nombre":           {""},
		"formacion_actual": {"3-5-2"},
		"valoracion":       {"80"},
		"posicion_tabla":   {"5"},
	})
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status %d, esperaba 400", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "El nombre del equipo no puede estar vacío") {
		t.Errorf("no se muestra el mensaje de validación")
	}
	if !strings.Contains(body, `value="3-5-2"`) {
		t.Errorf("el formulario no conserva lo que se había escrito")
	}
}

func TestCrearJugadorRedirige(t *testing.T) {
	router, db := setupTest(t)

	nombre := fmt.Sprintf("Jugador de prueba %d", time.Now().UnixNano())
	t.Cleanup(func() { db.Exec(`DELETE FROM jugadores WHERE nombre = $1`, nombre) })

	rr := postForm(t, router, "/jugadores", url.Values{
		"equipo_id":        {"1"},
		"nombre":           {nombre},
		"posicion":         {"MED"},
		"fecha_nacimiento": {"1985-09-09"},
		"altura":           {"1.72"},
		"media_general":    {"87"},
		"ritmo":            {"72"},
		"tiro":             {"76"},
		"pase":             {"89"},
		"regate":           {"88"},
		"defensa":          {"72"},
		"fisico":           {"65"},
	})
	esperarRedirect(t, rr, "/?tab=jugadores")

	if body := get(t, router, "/?tab=jugadores").Body.String(); !strings.Contains(body, nombre) {
		t.Errorf("el jugador creado no aparece en la lista")
	}
}

func TestCrearPartidoRedirige(t *testing.T) {
	router, db := setupTest(t)

	const fecha = "2031-03-15"
	t.Cleanup(func() { db.Exec(`DELETE FROM partidos WHERE fecha = $1`, fecha) })

	rr := postForm(t, router, "/partidos", url.Values{
		"equipo_local_id":     {"1"},
		"equipo_visitante_id": {"2"},
		"fecha":               {fecha},
		"estado":              {"PENDIENTE"},
	})
	esperarRedirect(t, rr, "/?tab=partidos")
}

func TestCrearPartidoMismoEquipo(t *testing.T) {
	router, _ := setupTest(t)

	rr := postForm(t, router, "/partidos", url.Values{
		"equipo_local_id":     {"1"},
		"equipo_visitante_id": {"1"},
		"fecha":               {"2031-03-15"},
		"estado":              {"PENDIENTE"},
	})
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status %d, esperaba 400", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "no pueden ser el mismo") {
		t.Errorf("no se muestra el mensaje de validación")
	}
}
