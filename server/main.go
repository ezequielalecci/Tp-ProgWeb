package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	database "tpweb/db/sqlc"
	"tpweb/server/routes"

	_ "github.com/lib/pq" // postgres
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		log.Fatal("Faltan variables de entorno para la base de datos (DB_HOST, etc.)")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto si no está en el .env
	}

	fmt.Printf("Servidor corriendo en http://localhost%s\n", port)

	//con router definimos nuestro propio MUX para que maneje las peticiones
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
