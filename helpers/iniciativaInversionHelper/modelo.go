package iniciativaInversionHelper

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Tipos de campo que entiende el motor (valor de dato.type en la plantilla).
const (
	TipoTexto      = "input"
	TipoNumerico   = "numeric"
	TipoSelect     = "select"
	TipoCatalogo   = "catalogo"
	TipoCalculado  = "calculado"
	TipoRepetible  = "repetible"
	TipoReferencia = "referencia"
	TipoMatriz     = "matriz"
)

// prefijoOrigen indica, en opciones_desde, que las opciones salen de la fila de
// origen del espejo más cercano (p. ej. "origen.causas_indirectas").
const prefijoOrigen = "origen."

// Datos es el contenido de un submódulo o de una fila: clave del campo -> valor.
type Datos = map[string]interface{}

// ---------------------------------------------------------------------------
// Definición (plantilla)
// ---------------------------------------------------------------------------

// Definicion es la construcción de una iniciativa: módulos (primer nivel, p. ej.
// Identificación) y, dentro de cada uno, sus submódulos u hojas (segundo nivel).
// Las claves de los submódulos son únicas en toda la plantilla.
type Definicion struct {
	PlantillaID string  `json:"plantilla_id"`
	Nombre      string  `json:"nombre"`
	TipoPlanID  string  `json:"tipo_plan_id"`
	Modulos     []*Nodo `json:"modulos"`
}

// Submodulos devuelve los submódulos de todos los módulos, en orden.
func (d *Definicion) Submodulos() []*Nodo {
	submodulos := []*Nodo{}
	for _, modulo := range d.Modulos {
		submodulos = append(submodulos, modulo.Hijos...)
	}
	return submodulos
}

// Submodulo devuelve el submódulo con la clave dada, o nil.
func (d *Definicion) Submodulo(clave string) *Nodo {
	return buscarNodo(d.Submodulos(), clave)
}

// Nodo es un subgrupo de la plantilla con su definición de campo.
type Nodo struct {
	ID           string  `json:"id"`
	Nombre       string  `json:"nombre"`
	Descripcion  string  `json:"descripcion"`
	BanderaTabla bool    `json:"bandera_tabla"`
	Dato         Dato    `json:"dato"`
	Hijos        []*Nodo `json:"hijos"`
}

// Clave técnica del nodo (dato.clave).
func (n *Nodo) Clave() string { return n.Dato.Clave }

// EsCampo indica si el nodo guarda un valor (tiene tipo). Si no, es un grupo.
func (n *Nodo) EsCampo() bool { return n.Dato.Tipo != "" }

// EsEspejo indica si es un grupo repetible cuyas filas salen de otro grupo.
func (n *Nodo) EsEspejo() bool { return n.Dato.Tipo == TipoRepetible && n.Dato.EspejoDe != "" }

// Hijo devuelve el hijo con la clave dada, o nil.
func (n *Nodo) Hijo(clave string) *Nodo {
	return buscarNodo(n.Hijos, clave)
}

func buscarNodo(nodos []*Nodo, clave string) *Nodo {
	for _, nodo := range nodos {
		if nodo.Clave() == clave {
			return nodo
		}
	}
	return nil
}

// Dato es la definición de un campo, tal como se guarda en subgrupo_detalle.dato.
type Dato struct {
	Clave          string          `json:"clave"`
	Modulo         bool            `json:"modulo,omitempty"`
	Tipo           string          `json:"type,omitempty"`
	Requerido      Booleano        `json:"required"`
	Multilinea     bool            `json:"multilinea,omitempty"`
	MaxLongitud    *int            `json:"maxlength,omitempty"`
	Min            *float64        `json:"min,omitempty"`
	Max            *float64        `json:"max,omitempty"`
	Decimales      *int            `json:"decimales,omitempty"`
	Opciones       []Opcion        `json:"options,omitempty"`
	Catalogo       string          `json:"catalogo,omitempty"`
	DependeDe      string          `json:"depende_de,omitempty"`
	Plantilla      string          `json:"plantilla,omitempty"`
	Transformar    string          `json:"transformar,omitempty"`
	Editable       bool            `json:"editable,omitempty"`
	NombrePlan     bool            `json:"nombre_plan,omitempty"`
	EtiquetaFila   string          `json:"etiqueta_fila,omitempty"`
	EspejoDe       string          `json:"espejo_de,omitempty"`
	Mostrar        string          `json:"mostrar,omitempty"`
	OpcionesDesde  string          `json:"opciones_desde,omitempty"`
	Multiple       bool            `json:"multiple,omitempty"`
	CubrirTodas    bool            `json:"cubrir_todas,omitempty"`
	TipoCelda      string          `json:"tipo_celda,omitempty"`
	Totales        *Totales        `json:"totales,omitempty"`
	Filas          []ElementoClave `json:"filas,omitempty"`
	GruposColumnas []GrupoColumnas `json:"grupos_columnas,omitempty"`
}

// Booleano acepta true/false y también "true"/"false", como guarda el cliente.
type Booleano bool

// UnmarshalJSON implementa json.Unmarshaler.
func (b *Booleano) UnmarshalJSON(data []byte) error {
	switch strings.Trim(strings.TrimSpace(string(data)), `"`) {
	case "true":
		*b = true
	case "false", "", "null", " ":
		*b = false
	default:
		return fmt.Errorf("valor booleano inválido: %s", data)
	}
	return nil
}

// Opcion de un select.
type Opcion struct {
	Valor string `json:"valor"`
}

// ElementoClave es una fila o columna de una matriz.
type ElementoClave struct {
	Clave  string `json:"clave"`
	Nombre string `json:"nombre"`
}

// GrupoColumnas agrupa columnas de una matriz bajo un encabezado.
type GrupoColumnas struct {
	Nombre   string          `json:"nombre"`
	Columnas []ElementoClave `json:"columnas"`
}

// Totales indica qué totales muestra una matriz (se calculan, no se guardan).
type Totales struct {
	Filas    bool `json:"filas"`
	Columnas bool `json:"columnas"`
}

// ---------------------------------------------------------------------------
// Documentos de planes_crud
// ---------------------------------------------------------------------------

// Plan es la vista del documento plan que usa este paquete. Al actualizar se
// reenvía completo porque PUT /plan/:id del crud valida el documento entero.
type Plan struct {
	ID                string `json:"_id,omitempty"`
	Nombre            string `json:"nombre"`
	Descripcion       string `json:"descripcion"`
	TipoPlanID        string `json:"tipo_plan_id"`
	AplicativoID      string `json:"aplicativo_id"`
	Activo            bool   `json:"activo"`
	Formato           bool   `json:"formato"`
	FormatoID         string `json:"formato_id,omitempty"`
	DependenciaID     string `json:"dependencia_id,omitempty"`
	EstadoPlanID      string `json:"estado_plan_id,omitempty"`
	VigenciaAplica    string `json:"vigencia_aplica,omitempty"`
	FechaCreacion     string `json:"fecha_creacion,omitempty"`
	FechaModificacion string `json:"fecha_modificacion,omitempty"`
}

// Subgrupo es un nodo del árbol tal como lo devuelve el crud.
type Subgrupo struct {
	ID           string   `json:"_id"`
	Nombre       string   `json:"nombre"`
	Descripcion  string   `json:"descripcion"`
	Activo       bool     `json:"activo"`
	BanderaTabla bool     `json:"bandera_tabla"`
	Hijos        []string `json:"hijos"`
}

// DetalleSubgrupo es el subgrupo_detalle de un nodo; Dato es texto JSON.
type DetalleSubgrupo struct {
	ID     string `json:"_id"`
	Dato   string `json:"dato"`
	Activo bool   `json:"activo"`
}

// Estado es un registro de los catálogos estado-plan o estado-submodulo.
type Estado struct {
	ID     string `json:"_id"`
	Nombre string `json:"nombre"`
	Codigo string `json:"codigo_abreviacion"`
	Color  string `json:"color,omitempty"`
}

// Observacion es una entrada del historial de comentarios de un submódulo (siempre
// sobre el submódulo completo). Se agregan y nunca se reemplazan.
type Observacion struct {
	ID      string `json:"_id,omitempty"`
	Texto   string `json:"texto"`
	Rol     string `json:"rol"`
	Usuario string `json:"usuario"`
	Fecha   string `json:"fecha,omitempty"`
}

// NuevaObservacion es el cuerpo para agregar una observación en planes_crud. Si
// EstadoSubmoduloID no es vacío, el estado del submódulo cambia en la misma operación.
type NuevaObservacion struct {
	Texto             string `json:"texto"`
	Rol               string `json:"rol"`
	Usuario           string `json:"usuario"`
	EstadoSubmoduloID string `json:"estado_submodulo_id,omitempty"`
}

// SubmoduloIniciativa es un documento de la colección submodulo_iniciativa.
type SubmoduloIniciativa struct {
	ID                string        `json:"_id"`
	PlanID            string        `json:"plan_id"`
	Submodulo         string        `json:"submodulo"`
	EstadoSubmoduloID string        `json:"estado_submodulo_id"`
	Datos             Datos         `json:"dato"`
	Observaciones     []Observacion `json:"observaciones"`
	Version           int           `json:"version"`
	Activo            bool          `json:"activo"`
	FechaModificacion string        `json:"fecha_modificacion,omitempty"`
}

// NuevoSubmodulo es el cuerpo para crear un submódulo.
type NuevoSubmodulo struct {
	PlanID            string `json:"plan_id"`
	Submodulo         string `json:"submodulo"`
	EstadoSubmoduloID string `json:"estado_submodulo_id"`
	Datos             Datos  `json:"dato"`
}

// CambiosSubmodulo es el cuerpo para actualizar un submódulo (bloqueo optimista por
// versión). El puntero distingue "no cambiar" (nil) de "dejar vacío". Las
// observaciones no van aquí: se agregan con NuevaObservacion.
type CambiosSubmodulo struct {
	Version           int    `json:"version"`
	Datos             *Datos `json:"dato,omitempty"`
	EstadoSubmoduloID string `json:"estado_submodulo_id,omitempty"`
}

// ---------------------------------------------------------------------------
// Solicitudes y respuestas del API
// ---------------------------------------------------------------------------

// ErrorCampo es un error de validación ubicado por la ruta del campo dentro de la
// submódulo, p. ej. "causas_directas[cd-1].causas_indirectas[ci-2].descripcion".
type ErrorCampo struct {
	Ruta    string `json:"ruta"`
	Mensaje string `json:"mensaje"`
}

// SolicitudCrearIniciativa es el cuerpo de POST /iniciativas. Sin PlantillaID se usa
// la plantilla más reciente; Nombre y Vigencias también son opcionales.
type SolicitudCrearIniciativa struct {
	PlantillaID   string          `json:"plantilla_id"`
	DependenciaID string          `json:"dependencia_id"`
	Nombre        string          `json:"nombre"`
	Vigencias     json.RawMessage `json:"vigencias"`
}

// SolicitudGuardarSubmodulo es el cuerpo de PUT /iniciativas/:id/submodulos/:clave.
type SolicitudGuardarSubmodulo struct {
	Version int   `json:"version"`
	Datos   Datos `json:"datos"`
}

// SolicitudRevisarSubmodulo es el cuerpo de PUT /iniciativas/:id/submodulos/:clave/revision:
// Planeación da el submódulo por revisado sin observaciones. Observaciones ya no
// se acepta aquí; solo se lee para responder 400 a quien la siga enviando.
type SolicitudRevisarSubmodulo struct {
	Version       int           `json:"version"`
	Observaciones []Observacion `json:"observaciones,omitempty"`
}

// SolicitudObservacion es el cuerpo de POST /iniciativas/:id/submodulos/:clave/observaciones.
type SolicitudObservacion struct {
	Texto   string `json:"texto"`
	Rol     string `json:"rol"`
	Usuario string `json:"usuario"`
}

// PlantillaResumen describe una plantilla disponible para crear iniciativas.
type PlantillaResumen struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// IniciativaListado es una fila del listado de iniciativas.
//   - Numero: consecutivo por orden de creación (1 = la primera). Se cuenta sobre
//     todas las iniciativas, activas o no, para que no cambie al desactivar otra.
//   - Estado: estado general de la iniciativa.
type IniciativaListado struct {
	ID            string `json:"id"`
	Numero        int    `json:"numero"`
	Nombre        string `json:"nombre"`
	DependenciaID string `json:"dependencia_id"`
	PlantillaID   string `json:"plantilla_id"`
	Estado        Estado `json:"estado"`
	PermiteEditar bool   `json:"permite_editar"`
	FechaCreacion string `json:"fecha_creacion"`
}

// ResumenSubmodulo es una fila de la vista general de la iniciativa.
type ResumenSubmodulo struct {
	Clave         string `json:"clave"`
	Nombre        string `json:"nombre"`
	Descripcion   string `json:"descripcion"`
	Orden         int    `json:"orden"`
	Estado        Estado `json:"estado"`
	Version       int    `json:"version"`
	Observaciones int    `json:"observaciones"`
}

// ResumenModulo agrupa los submódulos de un módulo en la vista general.
type ResumenModulo struct {
	Clave       string             `json:"clave"`
	Nombre      string             `json:"nombre"`
	Descripcion string             `json:"descripcion"`
	Submodulos  []ResumenSubmodulo `json:"submodulos"`
}

// ResumenIniciativa es la vista general: estado de la iniciativa y de cada submódulo, por módulo.
type ResumenIniciativa struct {
	ID               string          `json:"id"`
	Nombre           string          `json:"nombre"`
	DependenciaID    string          `json:"dependencia_id"`
	PlantillaID      string          `json:"plantilla_id"`
	VigenciaAplica   string          `json:"vigencia_aplica"`
	Estado           Estado          `json:"estado"`
	PermiteEditar    bool            `json:"permite_editar"`
	PermiteRevisar   bool            `json:"permite_revisar"`
	AccionesPosibles []string        `json:"acciones_posibles"`
	Modulos          []ResumenModulo `json:"modulos"`
}

// SubmoduloVista es un submódulo con su estado, observaciones y errores vigentes.
type SubmoduloVista struct {
	ID            string        `json:"id"`
	Clave         string        `json:"clave"`
	Estado        Estado        `json:"estado"`
	Datos         Datos         `json:"datos"`
	Observaciones []Observacion `json:"observaciones"`
	Version       int           `json:"version"`
	Errores       []ErrorCampo  `json:"errores"`
}

// SubmoduloConRelacionadas incluye los datos de los submódulos de los que depende
// (espejos y referencias), para que el cliente pueda armar el formulario.
type SubmoduloConRelacionadas struct {
	Submodulo    SubmoduloVista   `json:"submodulo"`
	Relacionadas map[string]Datos `json:"relacionadas"`
}
