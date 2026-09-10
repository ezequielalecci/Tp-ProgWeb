package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	database "tpweb/db/sqlc"
	"tpweb/server/routes"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) (*database.Queries, func()) {
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

	queries := database.New(db)
	cleanup := func() {
		db.Close()
	}

	return queries, cleanup
}

func TestListEquipos(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	router := routes.RegisterRoutes(queries)

	req, err := http.NewRequest("GET", "/api/equipos", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /api/equipos devolvió status %v, esperaba %v", status, http.StatusOK)
	}
}

func TestCrearEquipo(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	router := routes.RegisterRoutes(queries)

	body := map[string]interface{}{
		"nombre":           "Atlético de Madrid",
		"formacion_actual": "4-4-2",
		"escudo_url":       "https://example.com/atleti.png",
		"valoracion":       84,
		"posicion_tabla":   3,
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "/api/equipos", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated && status != http.StatusOK {
		t.Errorf("POST /api/equipos devolvió status %v, esperaba 201/200", status)
	}
}

func TestListJugadoresByEquipo(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	router := routes.RegisterRoutes(queries)

	req, err := http.NewRequest("GET", "/api/jugadores/equipo?equipo_id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /api/jugadores/equipo?equipo_id=1 devolvió status %v, esperaba %v", status, http.StatusOK)
	}
}

func TestCrearJugador(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	router := routes.RegisterRoutes(queries)

	// Se envía "altura" como string ("1.72") acorde a la definición del DTO
	jsonPayload := []byte(`{
		"equipo_id": 1,
		"nombre": "Luka Modric",
		"posicion": "MED",
		"fecha_nacimiento": "1985-09-09",
		"media_general": 87,
		"altura": "1.72",
		"ritmo": 72,
		"tiro": 76,
		"pase": 89,
		"regate": 88,
		"defensa": 72,
		"fisico": 65,
		"foto_url": "https://example.com/modric.png"
	}`)

	req, err := http.NewRequest("POST", "/api/jugador", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated && status != http.StatusOK {
		t.Errorf("POST /api/jugador devolvió status %v (body: %s), esperaba 201/200", status, rr.Body.String())
	}
}

func TestListPartidos(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	router := routes.RegisterRoutes(queries)

	req, err := http.NewRequest("GET", "/api/partidos", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET /api/partidos devolvió status %v, esperaba %v", status, http.StatusOK)
	}
}

func TestDeleteEquipoCascade(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	router := routes.RegisterRoutes(queries)

	req, err := http.NewRequest("DELETE", "/api/equipos/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent && status != http.StatusOK {
		t.Errorf("DELETE /api/equipos/1 devolvió status %v, esperaba 204/200", status)
	}
}
