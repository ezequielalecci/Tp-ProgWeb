package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // postgres

	"tp2/internal/database"
	"tp2/internal/handlers"
)

func main() {
	// inicio de base de datos
	connStr := "host=postgres port=5432 user=postgres password=secret dbname=gestion_futbol sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a PostgreSQL: %v", err)
	}

	queries := database.New(db)
	apiHandler := handlers.NewAPIHandler(queries)

	// rutas
	http.HandleFunc("/api/equipos", apiHandler.GetEquipos)
	http.HandleFunc("/", handlers.ServeViews)

	// inicio de servidor
	port := ":8080"
	fmt.Printf("Servidor corriendo en http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
