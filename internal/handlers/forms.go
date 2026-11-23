package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"tpweb/internal/dto"
	"tpweb/internal/service"
	"tpweb/views"
)

// POST /equipos
func (h *PageHandler) CrearEquipo(w http.ResponseWriter, r *http.Request) {
	if !parsearFormulario(w, r) {
		return
	}
	fallar := func(msg string) { h.formularioInvalido(w, r, views.TabEquipos, msg) }

	valoracion, err := formInt16(r, "valoracion", "La valoración")
	if err != nil {
		fallar(err.Error())
		return
	}
	posicion, err := formInt16(r, "posicion_tabla", "La posición en la tabla")
	if err != nil {
		fallar(err.Error())
		return
	}

	req := dto.CreateEquipoRequest{
		Nombre:          formTexto(r, "nombre"),
		FormacionActual: formTexto(r, "formacion_actual"),
		EscudoUrl:       formOpcional(r, "escudo_url"),
		Valoracion:      valoracion,
		PosicionTabla:   posicion,
	}
	if err := req.Validate(); err != nil {
		fallar(mayuscula(err.Error()))
		return
	}

	if _, err := h.equipos.CrearEquipo(r.Context(), req); err != nil {
		log.Printf("error al crear equipo: %v", err)
		fallar("No se pudo crear el equipo. Revisá que no exista otro con el mismo nombre y que la valoración esté entre 1 y 99.")
		return
	}

	equipos, err := h.equipos.GetEquipos(r.Context())
	if err != nil {
		log.Printf("error al listar equipos: %v", err)
		http.Error(w, "No se pudieron cargar los equipos", http.StatusInternalServerError)
		return
	}

	jugadores, err := h.jugadores.GetJugadores(r.Context())
	if err != nil {
		log.Printf("error al listar jugadores: %v", err)
		http.Error(w, "No se pudieron cargar los jugadores", http.StatusInternalServerError)
		return
	}

	if err := views.EquipoList(equipos, jugadores).Render(r.Context(), w); err != nil {
		log.Printf("error al renderizar lista de equipos: %v", err)
		http.Error(w, "No se pudo renderizar la lista", http.StatusInternalServerError)
	}
}

// POST /jugadores
func (h *PageHandler) CrearJugador(w http.ResponseWriter, r *http.Request) {
	if !parsearFormulario(w, r) {
		return
	}
	fallar := func(msg string) { h.formularioInvalido(w, r, views.TabJugadores, msg) }

	equipoID, err := formInt32(r, "equipo_id")
	if err != nil {
		fallar("Elegí un equipo para el jugador.")
		return
	}

	altura := formTexto(r, "altura")
	if _, err := strconv.ParseFloat(altura, 64); err != nil {
		fallar("La altura debe ser un número, por ejemplo 1.78.")
		return
	}

	req := dto.CreateJugadorRequest{
		EquipoID:        equipoID,
		Nombre:          formTexto(r, "nombre"),
		Posicion:        formTexto(r, "posicion"),
		FechaNacimiento: formTexto(r, "fecha_nacimiento"),
		Altura:          altura,
		FotoUrl:         formOpcional(r, "foto_url"),
	}

	stats := []struct {
		campo, etiqueta string
		destino         *int16
	}{
		{"media_general", "La media general", &req.MediaGeneral},
		{"ritmo", "El ritmo", &req.Ritmo},
		{"tiro", "El tiro", &req.Tiro},
		{"pase", "El pase", &req.Pase},
		{"regate", "El regate", &req.Regate},
		{"defensa", "La defensa", &req.Defensa},
		{"fisico", "El físico", &req.Fisico},
	}
	for _, s := range stats {
		v, err := formInt16(r, s.campo, s.etiqueta)
		if err != nil {
			fallar(err.Error())
			return
		}
		*s.destino = v
	}

	if err := req.Validate(); err != nil {
		fallar(mayuscula(err.Error()))
		return
	}

	if _, err := h.jugadores.CreateJugador(r.Context(), req); err != nil {
		switch {
		case errors.Is(err, service.ErrFechaInvalida):
			fallar(mayuscula(err.Error()))
		case errors.Is(err, service.ErrEquipoNotFound):
			fallar("El equipo elegido ya no existe.")
		default:
			log.Printf("error al crear jugador: %v", err)
			fallar("No se pudo crear el jugador. Revisá que los atributos estén entre 1 y 99 y que la altura sea válida.")
		}
		return
	}

	http.Redirect(w, r, views.URLTab(views.TabJugadores), http.StatusSeeOther)
}

// POST /partidos
func (h *PageHandler) CrearPartido(w http.ResponseWriter, r *http.Request) {
	if !parsearFormulario(w, r) {
		return
	}
	fallar := func(msg string) { h.formularioInvalido(w, r, views.TabPartidos, msg) }

	localID, err := formInt32(r, "equipo_local_id")
	if err != nil {
		fallar("Elegí el equipo local.")
		return
	}
	visitanteID, err := formInt32(r, "equipo_visitante_id")
	if err != nil {
		fallar("Elegí el equipo visitante.")
		return
	}

	fecha, _ := time.Parse("2006-01-02", formTexto(r, "fecha"))

	req := dto.CreatePartidoRequest{
		EquipoLocalID:     localID,
		EquipoVisitanteID: visitanteID,
		Fecha:             dto.DateOnly{Time: fecha},
		Estado:            formTexto(r, "estado"),
	}
	if err := req.Validate(); err != nil {
		fallar(mayuscula(err.Error()))
		return
	}

	if _, err := h.partidos.CrearPartido(r.Context(), req); err != nil {
		if errors.Is(err, service.ErrEquipoLocalNotFound) || errors.Is(err, service.ErrEquipoVisitNotFound) {
			fallar(mayuscula(err.Error()))
			return
		}
		log.Printf("error al crear partido: %v", err)
		fallar("No se pudo crear el partido. Revisá el estado elegido.")
		return
	}

	http.Redirect(w, r, views.URLTab(views.TabPartidos), http.StatusSeeOther)
}

func (h *PageHandler) formularioInvalido(w http.ResponseWriter, r *http.Request, tab, msg string) {
	h.renderIndex(w, r, http.StatusBadRequest, tab, views.FormState{Values: r.PostForm, Error: msg})
}

func parsearFormulario(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
		return false
	}
	return true
}

func formTexto(r *http.Request, campo string) string {
	return strings.TrimSpace(r.PostFormValue(campo))
}

func formOpcional(r *http.Request, campo string) *string {
	v := formTexto(r, campo)
	if v == "" {
		return nil
	}
	return &v
}

func formInt16(r *http.Request, campo, etiqueta string) (int16, error) {
	n, err := strconv.ParseInt(formTexto(r, campo), 10, 16)
	if err != nil {
		return 0, fmt.Errorf("%s debe ser un número entero.", etiqueta)
	}
	return int16(n), nil
}

func formInt32(r *http.Request, campo string) (int32, error) {
	n, err := strconv.ParseInt(formTexto(r, campo), 10, 32)
	return int32(n), err
}

func mayuscula(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:] + "."
}

// DELETE /equipos/{id}
func (h *PageHandler) EliminarEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de equipo inválido", http.StatusBadRequest)
		return
	}

	if err := h.equipos.DeleteEquipo(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrEquipoNotFound) {
			http.Error(w, "Equipo no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("error al eliminar equipo %d: %v", id, err)
		http.Error(w, "No se pudo eliminar el equipo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func formInt32Path(r *http.Request, campo string) (int32, error) {
	n, err := strconv.ParseInt(r.PathValue(campo), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

// DELETE /jugadores/{id}
func (h *PageHandler) EliminarJugador(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de jugador inválido", http.StatusBadRequest)
		return
	}

	if err := h.jugadores.DeleteJugador(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrJugadorNotFound) {
			http.Error(w, "Jugador no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("error al eliminar jugador %d: %v", id, err)
		http.Error(w, "No se pudo eliminar el jugador", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DELETE /partidos/{id}
func (h *PageHandler) EliminarPartido(w http.ResponseWriter, r *http.Request) {
	id, err := formInt32Path(r, "id")
	if err != nil || id <= 0 {
		http.Error(w, "ID de partido inválido", http.StatusBadRequest)
		return
	}

	if err := h.partidos.DeletePartido(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrPartidoNotFound) {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}
		log.Printf("error al eliminar partido %d: %v", id, err)
		http.Error(w, "No se pudo eliminar el partido", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// formInt32Opcional devuelve nil si el campo vino vacío. Ejemplo, si se actualiza un partido y el equipo fue borrado
func formInt32Opcional(r *http.Request, campo string) (*int32, error) {
	if formTexto(r, campo) == "" {
		return nil, nil
	}
	n, err := formInt32(r, campo)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
