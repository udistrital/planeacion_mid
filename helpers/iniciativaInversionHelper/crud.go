package iniciativaInversionHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RepositorioPlanes accede a los planes (iniciativas y plantillas).
type RepositorioPlanes interface {
	ObtenerPlan(id string) (*Plan, error)
	BuscarPlanes(filtros map[string]string) ([]Plan, error)
	CrearPlan(plan Plan) (*Plan, error)
	ActualizarPlan(plan Plan) error
	DesactivarPlan(id string) error
}

// RepositorioArbol lee la estructura de una plantilla.
type RepositorioArbol interface {
	HijosDe(id string) ([]Subgrupo, error)
	DetallesDe(subgrupoID string) ([]DetalleSubgrupo, error)
}

// RepositorioCatalogos lee los catálogos de planes_crud por código.
type RepositorioCatalogos interface {
	BuscarPorCodigo(recurso, codigo string) (*Estado, error)
	ListarActivos(recurso string) ([]Estado, error)
}

// RepositorioSubmodulos guarda los datos de los submódulos de una iniciativa.
type RepositorioSubmodulos interface {
	ListarSubmodulos(planID string) ([]SubmoduloIniciativa, error)
	CrearSubmodulos(submodulos []NuevoSubmodulo) ([]SubmoduloIniciativa, error)
	ActualizarSubmodulo(id string, cambios CambiosSubmodulo) (*SubmoduloIniciativa, error)
	AgregarObservacion(id string, nueva NuevaObservacion) (*SubmoduloIniciativa, error)
}

// Repositorio reúne todo lo que el servicio necesita de planes_crud.
type Repositorio interface {
	RepositorioPlanes
	RepositorioArbol
	RepositorioCatalogos
	RepositorioSubmodulos
}

// Recursos de catálogo de planes_crud.
const (
	RecursoTipoPlan        = "tipo-plan"
	RecursoEstadoPlan      = "estado-plan"
	RecursoEstadoSubmodulo = "estado-submodulo"
)

// ClienteCrud implementa Repositorio con HTTP. A diferencia de helpers.SendJson,
// trata como error cualquier respuesta 4xx/5xx o con Success=false (el crud
// antiguo responde HTTP 200 y pone el código real en Status).
type ClienteCrud struct {
	base string
	http *http.Client
}

// NuevoClienteCrud crea el cliente. host va sin esquema, como en conf/app.conf.
func NuevoClienteCrud(host string) *ClienteCrud {
	return &ClienteCrud{
		base: "http://" + strings.TrimSuffix(host, "/"),
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

type sobreCrud struct {
	Success *bool           `json:"Success"`
	Status  interface{}     `json:"Status"`
	Message string          `json:"Message"`
	Data    json.RawMessage `json:"Data"`
}

func (c *ClienteCrud) llamar(metodo, ruta string, cuerpo, destino interface{}) error {
	var lector io.Reader
	if cuerpo != nil {
		contenido, err := json.Marshal(cuerpo)
		if err != nil {
			return fmt.Errorf("no se pudo serializar la solicitud a %s: %w", ruta, err)
		}
		lector = bytes.NewReader(contenido)
	}
	solicitud, err := http.NewRequest(metodo, c.base+ruta, lector)
	if err != nil {
		return err
	}
	solicitud.Header.Set("Content-Type", "application/json")
	solicitud.Header.Set("Accept", "application/json")

	respuesta, err := c.http.Do(solicitud)
	if err != nil {
		return errorServicio("No se pudo conectar con planes_crud (%s %s): %v", metodo, ruta, err)
	}
	defer respuesta.Body.Close()

	var sobre sobreCrud
	if err := json.NewDecoder(respuesta.Body).Decode(&sobre); err != nil {
		return errorServicio("planes_crud respondió HTTP %d sin JSON válido (%s %s).", respuesta.StatusCode, metodo, ruta)
	}
	status := codigoReal(respuesta.StatusCode, sobre.Status)
	if status >= 400 || (sobre.Success != nil && !*sobre.Success) {
		return errorDesdeCrud(status, sobre.Message, metodo, ruta)
	}
	if destino == nil || len(sobre.Data) == 0 || string(sobre.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(sobre.Data, destino); err != nil {
		return errorServicio("Respuesta inesperada de planes_crud (%s %s): %v", metodo, ruta, err)
	}
	return nil
}

// codigoReal prefiere el Status del sobre cuando el HTTP es 200 (crud antiguo).
func codigoReal(http int, status interface{}) int {
	if http != 200 {
		return http
	}
	switch valor := status.(type) {
	case string:
		if n, err := strconv.Atoi(valor); err == nil {
			return n
		}
	case float64:
		return int(valor)
	}
	return http
}

// errorDesdeCrud traduce un error de planes_crud. Solo el 409 trae un mensaje
// pensado para el usuario (p. ej. «El submódulo cambió desde que se consultó»);
// los demás son técnicos (nombres de campos, errores de Mongo) y van al log.
func errorDesdeCrud(status int, mensaje, metodo, ruta string) *ErrorNegocio {
	detalle := fmt.Sprintf("planes_crud respondió %d en %s %s: %s", status, metodo, ruta, mensaje)
	switch status {
	case http.StatusConflict:
		if mensaje == "" {
			mensaje = "La información cambió mientras trabajaba. Recargue la página e intente de nuevo."
		}
		return errorConflicto("%s", mensaje)
	case http.StatusNotFound:
		return errorNoEncontrado("No se encontró la información solicitada.").conDetalle("%s", detalle)
	case http.StatusBadRequest:
		return errorFormatoSolicitud("%s", detalle)
	default:
		return errorServicio("%s", detalle)
	}
}

func consulta(filtros map[string]string) string {
	claves := make([]string, 0, len(filtros))
	for clave := range filtros {
		claves = append(claves, clave)
	}
	sort.Strings(claves)
	partes := make([]string, 0, len(claves))
	for _, clave := range claves {
		partes = append(partes, clave+":"+filtros[clave])
	}
	return "query=" + url.QueryEscape(strings.Join(partes, ","))
}

// --- RepositorioPlanes ------------------------------------------------------

// ObtenerPlan devuelve el plan o un error 404.
func (c *ClienteCrud) ObtenerPlan(id string) (*Plan, error) {
	var plan Plan
	if err := c.llamar(http.MethodGet, "/plan/"+url.PathEscape(id), nil, &plan); err != nil {
		return nil, err
	}
	if plan.ID == "" {
		return nil, errorNoEncontrado("No se encontró la información solicitada.").conDetalle("planes_crud no tiene el plan %s", id)
	}
	return &plan, nil
}

// BuscarPlanes filtra por igualdad de campos (query=campo:valor,...).
func (c *ClienteCrud) BuscarPlanes(filtros map[string]string) ([]Plan, error) {
	var planes []Plan
	err := c.llamar(http.MethodGet, "/plan?"+consulta(filtros)+"&sortby=fecha_creacion&order=desc", nil, &planes)
	return planes, err
}

// CrearPlan crea el plan y lo devuelve con su id.
func (c *ClienteCrud) CrearPlan(plan Plan) (*Plan, error) {
	var creado Plan
	if err := c.llamar(http.MethodPost, "/plan", plan, &creado); err != nil {
		return nil, err
	}
	return &creado, nil
}

// ActualizarPlan reenvía el plan completo.
func (c *ClienteCrud) ActualizarPlan(plan Plan) error {
	return c.llamar(http.MethodPut, "/plan/"+url.PathEscape(plan.ID), plan, nil)
}

// DesactivarPlan hace el borrado lógico (activo: false).
func (c *ClienteCrud) DesactivarPlan(id string) error {
	return c.llamar(http.MethodPut, "/plan/delete_plan/"+url.PathEscape(id), nil, nil)
}

// --- RepositorioArbol -------------------------------------------------------

// HijosDe devuelve los hijos ordenados (por hijos[] del padre, o por fecha si el padre es un plan).
func (c *ClienteCrud) HijosDe(id string) ([]Subgrupo, error) {
	var hijos []Subgrupo
	err := c.llamar(http.MethodGet, "/subgrupo/hijos/"+url.PathEscape(id), nil, &hijos)
	return hijos, err
}

// DetallesDe devuelve los subgrupo_detalle de un nodo.
func (c *ClienteCrud) DetallesDe(subgrupoID string) ([]DetalleSubgrupo, error) {
	var detalles []DetalleSubgrupo
	err := c.llamar(http.MethodGet, "/subgrupo-detalle/detalle/"+url.PathEscape(subgrupoID), nil, &detalles)
	return detalles, err
}

// --- RepositorioCatalogos ---------------------------------------------------

// BuscarPorCodigo devuelve el registro activo con ese codigo_abreviacion, o un 404.
func (c *ClienteCrud) BuscarPorCodigo(recurso, codigo string) (*Estado, error) {
	var registros []Estado
	ruta := "/" + recurso + "?" + consulta(map[string]string{"codigo_abreviacion": codigo, "activo": "true"})
	if err := c.llamar(http.MethodGet, ruta, nil, &registros); err != nil {
		return nil, err
	}
	if len(registros) == 0 {
		return nil, errorNoEncontrado("No se encontró la información solicitada.").
			conDetalle("no existe en %s un registro activo con código %s", recurso, codigo)
	}
	return &registros[0], nil
}

// ListarActivos devuelve los registros activos de un catálogo.
func (c *ClienteCrud) ListarActivos(recurso string) ([]Estado, error) {
	var registros []Estado
	err := c.llamar(http.MethodGet, "/"+recurso+"?"+consulta(map[string]string{"activo": "true"}), nil, &registros)
	return registros, err
}

// --- RepositorioSubmodulos ---------------------------------------------------

// ListarSubmodulos devuelve los submódulos activos de una iniciativa.
func (c *ClienteCrud) ListarSubmodulos(planID string) ([]SubmoduloIniciativa, error) {
	var submodulos []SubmoduloIniciativa
	ruta := "/submodulo-iniciativa?" + consulta(map[string]string{"plan_id": planID, "activo": "true"})
	err := c.llamar(http.MethodGet, ruta, nil, &submodulos)
	return submodulos, err
}

// CrearSubmodulos crea todos los submódulos en una sola llamada.
func (c *ClienteCrud) CrearSubmodulos(submodulos []NuevoSubmodulo) ([]SubmoduloIniciativa, error) {
	var creadas []SubmoduloIniciativa
	err := c.llamar(http.MethodPost, "/submodulo-iniciativa/lote", submodulos, &creadas)
	return creadas, err
}

// ActualizarSubmodulo responde 409 si la versión no es la vigente.
func (c *ClienteCrud) ActualizarSubmodulo(id string, cambios CambiosSubmodulo) (*SubmoduloIniciativa, error) {
	var actualizada SubmoduloIniciativa
	if err := c.llamar(http.MethodPut, "/submodulo-iniciativa/"+url.PathEscape(id), cambios, &actualizada); err != nil {
		return nil, err
	}
	return &actualizada, nil
}

// --- RepositorioElementosCatalogo -------------------------------------------

// ElementosDeCatalogo devuelve los elementos activos del catálogo, en su orden.
func (c *ClienteCrud) ElementosDeCatalogo(nombre string) ([]ItemCatalogo, error) {
	var items []ItemCatalogo
	ruta := "/elemento-catalogo?" + consulta(map[string]string{"catalogo": nombre, "activo": "true"}) + "&fields=codigo,nombre,padre&sortby=orden,codigo&order=asc"
	err := c.llamar(http.MethodGet, ruta, nil, &items)
	return items, err
}

// AgregarObservacion agrega la observación al historial sin cambiar la versión.
func (c *ClienteCrud) AgregarObservacion(id string, nueva NuevaObservacion) (*SubmoduloIniciativa, error) {
	var actualizada SubmoduloIniciativa
	if err := c.llamar(http.MethodPost, "/submodulo-iniciativa/"+url.PathEscape(id)+"/observaciones", nueva, &actualizada); err != nil {
		return nil, err
	}
	return &actualizada, nil
}
