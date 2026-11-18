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
	TipoCatalogo   = "catalogo"
	TipoCalculado  = "calculado"
	TipoRepetible  = "repetible"
	TipoReferencia = "referencia"
	TipoMatriz     = "matriz"
	TipoCasilla    = "casilla"
)

// Operaciones de un calculado numérico.
const (
	// OperacionResta: el primer operando menos los demás (p. ej. el déficit: oferta - demanda).
	// Si falta algún operando, no hay resultado.
	OperacionResta = "resta"
	// OperacionSuma: la suma de los operandos; los vacíos cuentan como 0 (p. ej. el costo total
	// de una actividad: la suma de sus periodos). Si todos están vacíos, no hay resultado.
	OperacionSuma = "suma"
	// OperacionMultiplicacion: el producto de los operandos (p. ej. el valor total de un año en
	// Ingresos y beneficios: cantidad × valor unitario). Si falta algún operando, no hay resultado.
	OperacionMultiplicacion = "multiplicacion"
)

// operacionValida indica si el calculado numérico está bien definido.
func operacionValida(dato Dato) bool {
	switch dato.Operacion {
	case OperacionResta, OperacionSuma, OperacionMultiplicacion:
		return len(dato.Operandos) >= 2
	}
	return false
}

// prefijoOrigen indica, en opciones_desde, desde y obligatorio_si, que el grupo o el campo
// es de la fila de origen del espejo más cercano (p. ej. "origen.periodo_0" en el Cronograma).
const prefijoOrigen = "origen."

// prefijoCatalogo indica, en espejo_de, que las filas salen de un catálogo: una por cada
// elemento (p. ej. "catalogo.resultados_esperados"). Cada elemento es una fila de origen
// con id = código y nombre; el cliente las recibe en relacionadas[claveRelacionadasCatalogo].
const prefijoCatalogo = "catalogo."

// claveRelacionadasCatalogo es la clave, en relacionadas, de los catálogos que espeja un
// submódulo: relacionadas["catalogo"]["<nombre>"] = filas de origen.
const claveRelacionadasCatalogo = "catalogo"

// catalogoDeRuta devuelve el nombre del catálogo de una ruta "catalogo.<nombre>".
func catalogoDeRuta(ruta string) (string, bool) {
	if !strings.HasPrefix(ruta, prefijoCatalogo) {
		return "", false
	}
	return strings.TrimPrefix(ruta, prefijoCatalogo), true
}

// filasDeCatalogo convierte los elementos de un catálogo en filas de origen de un espejo, con
// su descripción si la tienen (el cliente la muestra como ayuda bajo el nombre).
func filasDeCatalogo(catalogo *Catalogo) []Datos {
	filas := make([]Datos, 0, len(catalogo.Items))
	for _, item := range catalogo.Items {
		fila := Datos{"id": item.Codigo, "nombre": item.Nombre}
		if item.Descripcion != "" {
			fila["descripcion"] = item.Descripcion
		}
		filas = append(filas, fila)
	}
	return filas
}

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
	ID          string  `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Dato        Dato    `json:"dato"`
	Hijos       []*Nodo `json:"hijos"`
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
	Clave             string          `json:"clave"`
	Modulo            bool            `json:"modulo,omitempty"`
	Tipo              string          `json:"type,omitempty"`
	Requerido         Booleano        `json:"required"`
	Multilinea        bool            `json:"multilinea,omitempty"`
	MaxLongitud       *int            `json:"maxlength,omitempty"`
	Min               *float64        `json:"min,omitempty"`
	Max               *float64        `json:"max,omitempty"`
	Decimales         *int            `json:"decimales,omitempty"`
	Catalogo          string          `json:"catalogo,omitempty"` // en un calculado con desde, el catálogo del que se muestra el nombre (solo cliente)
	DependeDe         string          `json:"depende_de,omitempty"`
	Plantilla         string          `json:"plantilla,omitempty"`
	Operacion         string          `json:"operacion,omitempty"` // calculado numérico (OperacionResta, OperacionSuma, OperacionMultiplicacion); sin ella, el calculado es texto (plantilla)
	Desde             string          `json:"desde,omitempty"`     // calculado que solo muestra "<submodulo>.<campo>" de otro submódulo, "origen.<campo>" de la fila de origen del espejo o "<ruta de un grupo>.<campo>" de la fila de ese grupo que contiene la de origen; no se guarda
	Operandos         []string        `json:"operandos,omitempty"` // claves de los campos hermanos, en orden
	Transformar       string          `json:"transformar,omitempty"`
	Editable          bool            `json:"editable,omitempty"`
	NombrePlan        bool            `json:"nombre_plan,omitempty"`
	EtiquetaFila      string          `json:"etiqueta_fila,omitempty"`
	EspejoDe          string          `json:"espejo_de,omitempty"`          // "<submodulo>.<grupo>" o "catalogo.<nombre>"
	FilasOpcionales   bool            `json:"filas_opcionales,omitempty"`   // espejo: la fila vacía no se exige; min = cuántas se llenan como mínimo
	FilasObligatorias []string        `json:"filas_obligatorias,omitempty"` // con filas_opcionales: ids de origen (códigos del catálogo) cuya fila sí se exige
	Mostrar           string          `json:"mostrar,omitempty"`
	MostrarDetalle    string          `json:"mostrar_detalle,omitempty"` // referencia: segundo texto de la opción (p. ej. la descripción bajo el nombre de la alternativa); solo lo usa el cliente
	OpcionesDesde     string          `json:"opciones_desde,omitempty"`
	MarcadaEn         string          `json:"marcada_en,omitempty"`     // referencia a la fila marcada en esa casilla única (se llena sola)
	Multiple          bool            `json:"multiple,omitempty"`       // referencia o catálogo con varias opciones (lista de ids o códigos)
	SoloSi            *Condicion      `json:"solo_si,omitempty"`        // el campo se muestra, se exige y se guarda solo si se cumple
	ObligatorioSi     string          `json:"obligatorio_si,omitempty"` // "origen.<campo>": se exige solo si ese campo de la fila de origen del espejo tiene valor (un número distinto de 0)
	NumerarPor        string          `json:"numerar_por,omitempty"`    // lista: sus filas se numeran n.m, con n la fila de este grupo (de su ruta o de espejo_de) que las contiene; solo lo usa el cliente
	Unica             bool            `json:"unica,omitempty"`          // casilla que se marca en una sola fila de su lista
	Totales           *Totales        `json:"totales,omitempty"`        // matriz; en un espejo, columnas = fila de totales (solo cliente)
	Filas             []ElementoClave `json:"filas,omitempty"`
	GruposColumnas    []GrupoColumnas `json:"grupos_columnas,omitempty"`
}

// Condicion de solo_si: el campo hermano Campo vale Valor o, si es una lista (catálogo
// multiple), lo incluye. P. ej. «¿Cuál otra sede?» solo si en las sedes se marcó «Otra».
// Con Distinto, al revés: el campo aplica si no vale (ni incluye) Valor (p. ej. el
// discriminado de todos los trazadores presupuestales menos el de grupos étnicos).
type Condicion struct {
	Campo    string `json:"campo"`
	Valor    string `json:"valor"`
	Distinto bool   `json:"distinto,omitempty"`
}

// Cumple indica si la condición se cumple con los valores del mismo nivel.
func (c *Condicion) Cumple(hermanos Datos) bool {
	return c.coincide(hermanos[c.Campo]) != c.Distinto
}

func (c *Condicion) coincide(valor interface{}) bool {
	if texto, esTexto := comoTexto(valor); esTexto {
		return texto == c.Valor
	}
	lista, _ := comoLista(valor)
	for _, elemento := range lista {
		if texto, esTexto := comoTexto(elemento); esTexto && texto == c.Valor {
			return true
		}
	}
	return false
}

// visible indica si el nodo aplica con los valores de su nivel (sin solo_si, siempre).
func (n *Nodo) visible(hermanos Datos) bool {
	return n.Dato.SoloSi == nil || n.Dato.SoloSi.Cumple(hermanos)
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
	ID          string   `json:"_id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Activo      bool     `json:"activo"`
	Hijos       []string `json:"hijos"`
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

// ErrorCampo es un error de validación ubicado por la ruta del campo dentro del
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
// Planeación da el submódulo por revisado sin observaciones. Si trae Observaciones
// responde 400: las observaciones se registran con POST .../observaciones.
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
	Estado        Estado `json:"estado"`
	PermiteEditar bool   `json:"permite_editar"`
	FechaCreacion string `json:"fecha_creacion"`
}

// ResumenSubmodulo es una fila de la vista general de la iniciativa.
type ResumenSubmodulo struct {
	Clave         string `json:"clave"`
	Nombre        string `json:"nombre"`
	Estado        Estado `json:"estado"`
	Observaciones int    `json:"observaciones"`
}

// ResumenModulo agrupa los submódulos de un módulo en la vista general.
type ResumenModulo struct {
	Clave      string             `json:"clave"`
	Nombre     string             `json:"nombre"`
	Submodulos []ResumenSubmodulo `json:"submodulos"`
}

// ResumenIniciativa es la vista general: estado de la iniciativa y de cada submódulo, por módulo.
type ResumenIniciativa struct {
	ID               string          `json:"id"`
	Nombre           string          `json:"nombre"`
	Estado           Estado          `json:"estado"`
	PermiteEditar    bool            `json:"permite_editar"`
	PermiteRevisar   bool            `json:"permite_revisar"`
	AccionesPosibles []string        `json:"acciones_posibles"`
	Modulos          []ResumenModulo `json:"modulos"`
}

// SubmoduloVista es un submódulo con su estado, observaciones y errores vigentes.
type SubmoduloVista struct {
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
