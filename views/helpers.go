package views

import (
	"database/sql"
	"net/url"
	"strconv"
	"strings"
	"time"

	database "tpweb/db/sqlc"
)

// Pestañas de la página principal. El valor viaja en la URL: /?tab=jugadores
const (
	TabEquipos   = "equipos"
	TabJugadores = "jugadores"
	TabPartidos  = "partidos"
)

// opcion es un par valor/texto para armar pestañas, <select> y grupos de inputs.
type opcion struct {
	Valor string
	Texto string
}

var pestanias = []opcion{
	{TabEquipos, "Equipos"},
	{TabJugadores, "Jugadores"},
	{TabPartidos, "Partidos"},
}

var posiciones = []opcion{
	{"POR", "Arquero"},
	{"DEF", "Defensor"},
	{"MED", "Mediocampista"},
	{"DEL", "Delantero"},
}

var estadosPartido = []opcion{
	{"PENDIENTE", "Pendiente"},
	{"EN CURSO", "En curso"},
	{"FINALIZADO", "Finalizado"},
}

// Atributos del jugador que van de 1 a 99. Valor es el name del input.
var statsJugador = []opcion{
	{"media_general", "Media general"},
	{"ritmo", "Ritmo"},
	{"tiro", "Tiro"},
	{"pase", "Pase"},
	{"regate", "Regate"},
	{"defensa", "Defensa"},
	{"fisico", "Físico"},
}

// PageData es todo lo que necesita IndexPage para dibujarse.
// El handler lo arma con los datos de la base y se lo pasa al componente.
type PageData struct {
	Tab       string
	Equipos   []database.Equipo
	Jugadores []database.Jugadore
	Partidos  []database.Partido
	Form      FormState
}

// FormState guarda lo que el usuario envió y el error de validación, para
// volver a mostrar el formulario sin que pierda lo que ya había escrito.
// En un GET normal está vacío.
type FormState struct {
	Values url.Values
	Error  string
}

func (f FormState) Get(campo string) string {
	return f.Values.Get(campo)
}

func (f FormState) GetOr(campo, porDefecto string) string {
	if v := f.Values.Get(campo); v != "" {
		return v
	}
	return porDefecto
}

// NormalizarTab devuelve una pestaña válida; cualquier otro valor cae en equipos.
func NormalizarTab(tab string) string {
	switch tab {
	case TabJugadores, TabPartidos:
		return tab
	}
	return TabEquipos
}

// URLTab es la dirección de una pestaña. Equipos es la principal y vive en "/".
func URLTab(tab string) string {
	if tab == TabEquipos {
		return "/"
	}
	return "/?tab=" + tab
}

func tituloTab(tab string) string {
	for _, p := range pestanias {
		if p.Valor == tab {
			return p.Texto
		}
	}
	return "Equipos"
}

func idStr(id int32) string {
	return strconv.Itoa(int(id))
}

func fecha(t time.Time) string {
	return t.Format("02/01/2006")
}

func nombreEquipo(equipos []database.Equipo, id int32) string {
	for _, e := range equipos {
		if e.ID == id {
			return e.Nombre
		}
	}
	return "—"
}

// nombreEquipoPartido contempla que el equipo haya sido eliminado:
// la FK de partidos es ON DELETE SET NULL, así que la columna puede ser NULL.
func nombreEquipoPartido(equipos []database.Equipo, id sql.NullInt32) string {
	if !id.Valid {
		return "(equipo eliminado)"
	}
	return nombreEquipo(equipos, id.Int32)
}

func nombreJugador(jugadores []database.Jugadore, id sql.NullInt32) string {
	if !id.Valid {
		return "—"
	}
	for _, j := range jugadores {
		if j.ID == id.Int32 {
			return j.Nombre
		}
	}
	return "—"
}

// claseEstado convierte "EN CURSO" en "estado-en-curso" para usarlo como clase CSS.
func claseEstado(estado string) string {
	return "estado-" + strings.ReplaceAll(strings.ToLower(estado), " ", "-")
}
