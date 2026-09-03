package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	database "tpweb/db/sqlc"
	"tpweb/server/routes"

	_ "github.com/lib/pq" // postgres
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

	// rutas
	router := routes.RegisterRoutes(queries)

	// inicio de servidor
	port := ":8080"
	fmt.Printf("Servidor corriendo en http://localhost%s\n", port)
	//con router definimos nuestro propio MUX para que maneje las peticiones
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
