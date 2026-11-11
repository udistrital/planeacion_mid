package iniciativaInversionHelper

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Unidad es una dependencia (Oikos) en la que el usuario formula iniciativas.
type Unidad struct {
	ID     int    `json:"Id"`
	Nombre string `json:"Nombre"`
}

// VinculacionTercero es la parte de una vinculación de Terceros que se usa aquí.
type VinculacionTercero struct {
	CargoID       int `json:"CargoId"`
	DependenciaID int `json:"DependenciaId"`
}

// FuenteUnidades consulta Terceros y Oikos.
type FuenteUnidades interface {
	// TerceroDe devuelve el id del tercero con ese documento, o 0 si no existe.
	TerceroDe(documento string) (int, error)
	VinculacionesActivas(terceroID int) ([]VinculacionTercero, error)
	NombreDependencia(dependenciaID int) (string, error)
}

// cargosConUnidad son los cargos con los que Formulación reconoce la unidad de un
// usuario (GET /v1/formulacion/vinculacion_tercero): vicerrector académico (319),
// jefe de oficina (312), asistente de dependencia (320) y administrativo de
// planta (300). Si Formulación cambia esa lista, hay que cambiarla aquí también.
var cargosConUnidad = map[int]bool{319: true, 312: true, 320: true, 300: true}

// patronDocumento acepta solo letras y números: el documento va dentro de la
// consulta de Terceros (query=Numero:<documento>,...) y una coma o dos puntos
// agregarían condiciones a esa consulta.
var patronDocumento = regexp.MustCompile(`^[0-9A-Za-z]{1,20}$`)

// Unidades resuelve las unidades de un usuario a partir de su documento. Se hace
// en el mid y no en el cliente porque Terceros no acepta peticiones del navegador
// desde cualquier origen (CORS), y así el navegador hace una sola llamada.
type Unidades struct {
	fuente FuenteUnidades
}

// NuevasUnidades arma el servicio con su fuente.
func NuevasUnidades(fuente FuenteUnidades) *Unidades {
	return &Unidades{fuente: fuente}
}

// NuevasUnidadesPorDefecto consulta Terceros y Oikos por HTTP (hosts sin esquema,
// como en conf/app.conf).
func NuevasUnidadesPorDefecto(hostTerceros, hostOikos string) *Unidades {
	return NuevasUnidades(NuevoClienteUnidades(hostTerceros, hostOikos))
}

// De devuelve las unidades del usuario, sin repetir y en el orden de sus
// vinculaciones. Si el documento no tiene tercero o vinculaciones con esos cargos,
// la lista va vacía (no es un error).
func (u *Unidades) De(documento string) ([]Unidad, error) {
	documento = strings.TrimSpace(documento)
	if !patronDocumento.MatchString(documento) {
		return nil, errorSolicitud("No se pudo identificar su documento. Vuelva a iniciar sesión e intente de nuevo.").
			conDetalle("el documento no cumple el formato: solo letras y números, máximo 20")
	}
	terceroID, err := u.fuente.TerceroDe(documento)
	if err != nil {
		return nil, err
	}
	unidades := []Unidad{}
	if terceroID == 0 {
		return unidades, nil
	}
	vinculaciones, err := u.fuente.VinculacionesActivas(terceroID)
	if err != nil {
		return nil, err
	}
	vistas := map[int]bool{}
	for _, vinculacion := range vinculaciones {
		dependencia := vinculacion.DependenciaID
		if !cargosConUnidad[vinculacion.CargoID] || dependencia == 0 || vistas[dependencia] {
			continue
		}
		vistas[dependencia] = true
		nombre, err := u.fuente.NombreDependencia(dependencia)
		if err != nil {
			return nil, err
		}
		unidades = append(unidades, Unidad{ID: dependencia, Nombre: nombre})
	}
	return unidades, nil
}

// ---------------------------------------------------------------------------
// Cliente HTTP de Terceros y Oikos
// ---------------------------------------------------------------------------

// ClienteUnidades implementa FuenteUnidades con HTTP. Trata como error cualquier
// respuesta distinta de 200.
type ClienteUnidades struct {
	terceros string
	oikos    string
	http     *http.Client
}

// NuevoClienteUnidades crea el cliente. Los hosts van sin esquema, como en conf/app.conf.
func NuevoClienteUnidades(hostTerceros, hostOikos string) *ClienteUnidades {
	return &ClienteUnidades{
		terceros: "http://" + strings.TrimSuffix(hostTerceros, "/"),
		oikos:    "http://" + strings.TrimSuffix(hostOikos, "/"),
		http:     &http.Client{Timeout: 30 * time.Second},
	}
}

// TerceroDe busca el tercero por número de documento. Terceros responde [{}]
// cuando no hay resultados.
func (c *ClienteUnidades) TerceroDe(documento string) (int, error) {
	var datos []struct {
		TerceroID *struct {
			ID int `json:"Id"`
		} `json:"TerceroId"`
	}
	consulta := url.QueryEscape("Numero:" + documento + ",Activo:true")
	if err := c.obtener("Terceros", c.terceros+"/datos_identificacion?query="+consulta+"&fields=TerceroId", &datos); err != nil {
		return 0, err
	}
	for _, dato := range datos {
		if dato.TerceroID != nil && dato.TerceroID.ID != 0 {
			return dato.TerceroID.ID, nil
		}
	}
	return 0, nil
}

// VinculacionesActivas devuelve todas las vinculaciones activas del tercero.
func (c *ClienteUnidades) VinculacionesActivas(terceroID int) ([]VinculacionTercero, error) {
	var vinculaciones []VinculacionTercero
	consulta := url.QueryEscape("Activo:true,TerceroPrincipalId:" + strconv.Itoa(terceroID))
	ruta := c.terceros + "/vinculacion?query=" + consulta + "&sortby=Id&order=asc&limit=0&fields=CargoId,DependenciaId"
	if err := c.obtener("Terceros", ruta, &vinculaciones); err != nil {
		return nil, err
	}
	return vinculaciones, nil
}

// NombreDependencia devuelve el nombre de la dependencia en Oikos.
func (c *ClienteUnidades) NombreDependencia(dependenciaID int) (string, error) {
	var dependencia struct {
		Nombre string `json:"Nombre"`
	}
	if err := c.obtener("Oikos", c.oikos+"/dependencia/"+strconv.Itoa(dependenciaID), &dependencia); err != nil {
		return "", err
	}
	if dependencia.Nombre == "" {
		return "", errorServicio("Oikos no devolvió el nombre de la dependencia %d", dependenciaID)
	}
	return dependencia.Nombre, nil
}

func (c *ClienteUnidades) obtener(servicio, ruta string, destino interface{}) error {
	respuesta, err := c.http.Get(ruta)
	if err != nil {
		return errorServicio("no se pudo conectar con %s: %v", servicio, err)
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK {
		return errorServicio("%s respondió HTTP %d en %s", servicio, respuesta.StatusCode, ruta)
	}
	if err := json.NewDecoder(respuesta.Body).Decode(destino); err != nil {
		return errorServicio("respuesta inesperada de %s en %s: %v", servicio, ruta, err)
	}
	return nil
}
