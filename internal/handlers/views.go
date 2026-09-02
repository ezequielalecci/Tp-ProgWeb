package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ServeViews(w http.ResponseWriter, r *http.Request) {
	// bloquear los metodos que no sean "GET"
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		// Si la petición viene de la API, responder con JSON
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Método no permitido. Solo se acepta GET",
			})
			return
		}

		// Si es una petición HTML web normal, se reponde con texto/plano
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("405 Método no permitido"))
		return
	}

	// endpoint inexistente
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Endpoint de API no encontrado",
		})
		return
	}

	//servir HTMLS validos
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch r.URL.Path {
	case "/", "/index.html":
		http.ServeFile(w, r, "static/index.html")
	case "/equipos", "/equipos.html":
		http.ServeFile(w, r, "static/equipos.html")
	default:
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "static/no-encontrado.html")
	}
}
