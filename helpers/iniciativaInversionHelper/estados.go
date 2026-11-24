package iniciativaInversionHelper

import (
	"net/http"
	"sync"
)

// Códigos del tipo de plan y de los estados (catálogos tipo-plan, estado-plan y
// estado-submodulo de planes_crud). Son propios de las iniciativas de inversión
// para no mezclarse con los proyectos MACRO (PRI_SP) ni con Formulación.
//
// Los estados son los de la «Relación de estados» del mockup (docs/Refinamiento mockup
// para Banco de proyectos OAP.pdf). Cada estado general va con uno detallado (submódulo):
//
//	General               Submódulo          Cuándo
//	Sin formular          Sin formular       nada diligenciado todavía
//	Formulación iniciada  Formulando         con información, sin enviar (el submódulo, incompleto)
//	Enviado a revisión    Formulado          todo formulado y enviado a Planeación (el submódulo, completo)
//	Avalado               Aprobado           todo revisado y aprobado por Planeación
//	Por corregir          Con observaciones  al menos un submódulo con observaciones por atender
const (
	CodigoTipoPlanIniciativa = "IPI_SP"

	EstadoIniciativaSinFormular         = "SF_IPI"
	EstadoIniciativaFormulacionIniciada = "EF_IPI"
	EstadoIniciativaEnRevision          = "ER_IPI"
	EstadoIniciativaPorCorregir         = "PC_IPI"
	EstadoIniciativaAvalado             = "AV_IPI"

	EstadoSubmoduloSinFormular      = "SF_SUB"
	EstadoSubmoduloFormulando       = "FN_SUB"
	EstadoSubmoduloFormulado        = "FOR_SUB"
	EstadoSubmoduloConObservaciones = "CO_SUB"
	EstadoSubmoduloAprobado         = "SO_SUB"

	// CodigoTipoSeguimientoFormulacion es el tipo de seguimiento de las fechas de formulación
	// de iniciativas (ver fechas.go).
	CodigoTipoSeguimientoFormulacion = "FIPI_SP"
)

var codigosEstadoIniciativa = []string{
	EstadoIniciativaSinFormular,
	EstadoIniciativaFormulacionIniciada,
	EstadoIniciativaEnRevision,
	EstadoIniciativaPorCorregir,
	EstadoIniciativaAvalado,
}

var codigosEstadoSubmodulo = []string{
	EstadoSubmoduloSinFormular,
	EstadoSubmoduloFormulando,
	EstadoSubmoduloFormulado,
	EstadoSubmoduloConObservaciones,
	EstadoSubmoduloAprobado,
}

// Acciones sobre el estado general de la iniciativa.
const (
	AccionEnviar   = "enviar"
	AccionDevolver = "devolver"
	AccionAvalar   = "avalar"
)

// transicion describe una acción: desde qué estados se permite, a cuál lleva y
// qué deben cumplir los submódulos.
type transicion struct {
	desde     []string
	hacia     string
	requisito func(estadosSubmodulos []string) string // "" si se cumple; si no, el motivo
}

var transiciones = map[string]transicion{
	AccionEnviar: {
		desde:     []string{EstadoIniciativaFormulacionIniciada, EstadoIniciativaPorCorregir},
		hacia:     EstadoIniciativaEnRevision,
		requisito: todasEn("Todos los submódulos deben estar formulados.", EstadoSubmoduloFormulado, EstadoSubmoduloAprobado),
	},
	AccionDevolver: {
		desde:     []string{EstadoIniciativaEnRevision},
		hacia:     EstadoIniciativaPorCorregir,
		requisito: algunaEn("Debe haber al menos un submódulo con observaciones.", EstadoSubmoduloConObservaciones),
	},
	AccionAvalar: {
		desde:     []string{EstadoIniciativaEnRevision},
		hacia:     EstadoIniciativaAvalado,
		requisito: todasEn("Todos los submódulos deben estar revisados y aprobados.", EstadoSubmoduloAprobado),
	},
}

// ordenAcciones fija el orden en que se listan las acciones posibles.
var ordenAcciones = []string{AccionEnviar, AccionDevolver, AccionAvalar}

// nombresAcciones son los nombres con que el usuario ve las acciones (los botones
// del cliente, ACCIONES_INICIATIVA); los mensajes no muestran la clave del API.
var nombresAcciones = map[string]string{
	AccionEnviar:   "Enviar a revisión",
	AccionDevolver: "Devolver para corregir",
	AccionAvalar:   "Avalar",
}

// PermiteEditar indica si en ese estado de la iniciativa se pueden guardar submódulos.
func PermiteEditar(estadoIniciativa string) bool {
	return contiene([]string{EstadoIniciativaSinFormular, EstadoIniciativaFormulacionIniciada, EstadoIniciativaPorCorregir}, estadoIniciativa)
}

// PermiteRevisar indica si en ese estado de la iniciativa se pueden revisar submódulos.
func PermiteRevisar(estadoIniciativa string) bool {
	return estadoIniciativa == EstadoIniciativaEnRevision
}

// EstadoTrasGuardar: un submódulo sin errores queda Formulado; con errores, Formulando si ya
// tiene información y Sin formular si no tiene nada diligenciado.
func EstadoTrasGuardar(errores []ErrorCampo, conInformacion bool) string {
	switch {
	case len(errores) == 0:
		return EstadoSubmoduloFormulado
	case conInformacion:
		return EstadoSubmoduloFormulando
	default:
		return EstadoSubmoduloSinFormular
	}
}

// EstadoIniciativaTrasGuardar: la iniciativa Sin formular pasa a Formulación iniciada en cuanto
// algún submódulo tiene información. En los demás estados no cambia ("" = sin cambio).
func EstadoIniciativaTrasGuardar(estadoIniciativa string, conInformacion bool) string {
	if estadoIniciativa == EstadoIniciativaSinFormular && conInformacion {
		return EstadoIniciativaFormulacionIniciada
	}
	return ""
}

// EvaluarAccion devuelve el estado destino o el motivo por el que no se puede.
func EvaluarAccion(accion, estadoIniciativa string, estadosSubmodulos []string) (string, string) {
	t, existe := transiciones[accion]
	if !existe {
		return "", "La acción solicitada no existe."
	}
	if !contiene(t.desde, estadoIniciativa) {
		return "", "La acción «" + nombresAcciones[accion] + "» no se permite en el estado actual de la iniciativa."
	}
	if motivo := t.requisito(estadosSubmodulos); motivo != "" {
		return "", motivo
	}
	return t.hacia, ""
}

// AccionesPosibles lista las acciones que hoy se pueden ejecutar.
func AccionesPosibles(estadoIniciativa string, estadosSubmodulos []string) []string {
	posibles := []string{}
	for _, accion := range ordenAcciones {
		if _, motivo := EvaluarAccion(accion, estadoIniciativa, estadosSubmodulos); motivo == "" {
			posibles = append(posibles, accion)
		}
	}
	return posibles
}

// SubmoduloRevisable indica si un submódulo en ese estado se puede revisar (ya está formulado).
func SubmoduloRevisable(estadoSubmodulo string) bool {
	return estadoSubmodulo != EstadoSubmoduloSinFormular && estadoSubmodulo != EstadoSubmoduloFormulando
}

func todasEn(motivo string, permitidos ...string) func([]string) string {
	return func(estados []string) string {
		if len(estados) == 0 {
			return motivo
		}
		for _, estado := range estados {
			if !contiene(permitidos, estado) {
				return motivo
			}
		}
		return ""
	}
}

func algunaEn(motivo string, buscados ...string) func([]string) string {
	return func(estados []string) string {
		for _, estado := range estados {
			if contiene(buscados, estado) {
				return ""
			}
		}
		return motivo
	}
}

func contiene(lista []string, valor string) bool {
	for _, elemento := range lista {
		if elemento == valor {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Catálogos del sistema (ids de tipo de plan y estados)
// ---------------------------------------------------------------------------

// catalogosSistema resuelve códigos a registros de planes_crud y los guarda en
// memoria. Si falta un registro, el log indica qué cargar.
type catalogosSistema struct {
	repo RepositorioCatalogos

	mu               sync.Mutex
	cargado          bool
	tipoPlan         Estado
	iniciativaPorCod map[string]Estado
	iniciativaPorID  map[string]Estado
	submoduloPorCod  map[string]Estado
	submoduloPorID   map[string]Estado
	tipoSeguimiento  *Estado
}

func nuevosCatalogosSistema(repo RepositorioCatalogos) *catalogosSistema {
	return &catalogosSistema{repo: repo}
}

func (c *catalogosSistema) asegurar() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cargado {
		return nil
	}
	tipo, err := c.repo.BuscarPorCodigo(RecursoTipoPlan, CodigoTipoPlanIniciativa)
	if err != nil {
		return faltaCatalogo(err)
	}
	iniciativaPorCod := map[string]Estado{}
	iniciativaPorID := map[string]Estado{}
	for _, codigo := range codigosEstadoIniciativa {
		estado, err := c.repo.BuscarPorCodigo(RecursoEstadoPlan, codigo)
		if err != nil {
			return faltaCatalogo(err)
		}
		iniciativaPorCod[codigo] = *estado
		iniciativaPorID[estado.ID] = *estado
	}
	estadosSubmodulo, err := c.repo.ListarActivos(RecursoEstadoSubmodulo)
	if err != nil {
		return err
	}
	submoduloPorCod := map[string]Estado{}
	submoduloPorID := map[string]Estado{}
	for _, estado := range estadosSubmodulo {
		submoduloPorCod[estado.Codigo] = estado
		submoduloPorID[estado.ID] = estado
	}
	for _, codigo := range codigosEstadoSubmodulo {
		if _, existe := submoduloPorCod[codigo]; !existe {
			return errorConfiguracion("no existe en %s un registro activo con código %s. %s", RecursoEstadoSubmodulo, codigo, comoCargarEstados)
		}
	}
	c.tipoPlan, c.iniciativaPorCod, c.iniciativaPorID = *tipo, iniciativaPorCod, iniciativaPorID
	c.submoduloPorCod, c.submoduloPorID = submoduloPorCod, submoduloPorID
	c.cargado = true
	return nil
}

// comoCargarEstados va en el log (nunca en la respuesta) cuando faltan el tipo de
// plan, los estados o el tipo de seguimiento.
const comoCargarEstados = "El tipo de plan, los estados y el tipo de seguimiento se crean con database local/construccion-proyecto-inversion/cargar-estados.js."

// faltaCatalogo convierte la falta de un registro base (tipo de plan o estado) en
// un error de configuración. Las fallas del servicio se devuelven tal cual.
func faltaCatalogo(err error) error {
	negocio := ComoErrorNegocio(err)
	if negocio.Status != http.StatusNotFound {
		return err
	}
	causa := negocio.Detalle
	if causa == "" {
		causa = negocio.Mensaje
	}
	return errorConfiguracion("%s. %s", causa, comoCargarEstados)
}

// TipoSeguimientoFormulacion devuelve el tipo de seguimiento de las fechas de formulación. No se
// carga en asegurar: si falta, las iniciativas se consultan igual y solo queda cerrada la
// formulación. Se guarda en memoria la primera vez que se encuentra.
func (c *catalogosSistema) TipoSeguimientoFormulacion() (Estado, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.tipoSeguimiento != nil {
		return *c.tipoSeguimiento, nil
	}
	tipo, err := c.repo.BuscarPorCodigo(RecursoTipoSeguimiento, CodigoTipoSeguimientoFormulacion)
	if err != nil {
		return Estado{}, faltaCatalogo(err)
	}
	c.tipoSeguimiento = tipo
	return *tipo, nil
}

func (c *catalogosSistema) TipoPlanIniciativa() (Estado, error) {
	if err := c.asegurar(); err != nil {
		return Estado{}, err
	}
	return c.tipoPlan, nil
}

func (c *catalogosSistema) EstadoIniciativa(codigo string) (Estado, error) {
	if err := c.asegurar(); err != nil {
		return Estado{}, err
	}
	return c.iniciativaPorCod[codigo], nil
}

func (c *catalogosSistema) EstadoIniciativaPorID(id string) (Estado, error) {
	if err := c.asegurar(); err != nil {
		return Estado{}, err
	}
	estado, existe := c.iniciativaPorID[id]
	if !existe {
		return Estado{ID: id, Nombre: "Desconocido"}, nil
	}
	return estado, nil
}

func (c *catalogosSistema) EstadoSubmodulo(codigo string) (Estado, error) {
	if err := c.asegurar(); err != nil {
		return Estado{}, err
	}
	return c.submoduloPorCod[codigo], nil
}

func (c *catalogosSistema) EstadoSubmoduloPorID(id string) (Estado, error) {
	if err := c.asegurar(); err != nil {
		return Estado{}, err
	}
	estado, existe := c.submoduloPorID[id]
	if !existe {
		return Estado{ID: id, Nombre: "Desconocido"}, nil
	}
	return estado, nil
}
