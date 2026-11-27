package iniciativaInversionHelper

import (
	"net/http"
	"sync"
)

// Códigos del tipo de plan y de los estados (catálogos tipo-plan, estado-plan y
// estado-submodulo de planes_crud). Son propios de las iniciativas de inversión
// para no mezclarse con los proyectos MACRO (PRI_SP) ni con Formulación.
//
// El flujo está en docs/Máquina de estados - Iniciativas de inversión.md:
//
//	01 Sin formular → 02 En formulación → 03 Por aprobar (Jefe) → 04 Por revisar OAP → 05 En revisión OAP
//	05 → 07 Aprobado (todos los submódulos aprobados) · 05 → 06 Por corregir (alguno con observaciones)
//	03 → 06 (el Jefe la devuelve) · 06 → 03 (se reenvía) · 05 → 08 Rechazada · 08 → 06 (la OAP se retracta)
const (
	CodigoTipoPlanIniciativa = "IPI_SP"

	EstadoIniciativaSinFormular   = "SF_IPI"
	EstadoIniciativaEnFormulacion = "EF_IPI"
	EstadoIniciativaPorAprobar    = "PA_IPI"
	EstadoIniciativaPorRevisar    = "PR_IPI"
	EstadoIniciativaEnRevision    = "ER_IPI"
	EstadoIniciativaPorCorregir   = "PC_IPI"
	EstadoIniciativaAprobada      = "APR_IPI"
	EstadoIniciativaRechazada     = "RCH_IPI"

	EstadoSubmoduloSinFormular      = "SF_SUB"
	EstadoSubmoduloFormulando       = "FN_SUB"
	EstadoSubmoduloFormulado        = "FOR_SUB"
	EstadoSubmoduloConObservaciones = "CO_SUB"
	EstadoSubmoduloAprobado         = "APR_SUB"

	// CodigoTipoSeguimientoFormulacion es el tipo de seguimiento de las fechas de formulación
	// de iniciativas (ver fechas.go).
	CodigoTipoSeguimientoFormulacion = "FIPI_SP"
)

var codigosEstadoIniciativa = []string{
	EstadoIniciativaSinFormular,
	EstadoIniciativaEnFormulacion,
	EstadoIniciativaPorAprobar,
	EstadoIniciativaPorRevisar,
	EstadoIniciativaEnRevision,
	EstadoIniciativaPorCorregir,
	EstadoIniciativaAprobada,
	EstadoIniciativaRechazada,
}

var codigosEstadoSubmodulo = []string{
	EstadoSubmoduloSinFormular,
	EstadoSubmoduloFormulando,
	EstadoSubmoduloFormulado,
	EstadoSubmoduloConObservaciones,
	EstadoSubmoduloAprobado,
}

// Acciones sobre el estado general de la iniciativa (PUT /iniciativas/:id/acciones/:accion).
const (
	AccionEnviar           = "enviar"
	AccionAprobar          = "aprobar"
	AccionDevolver         = "devolver"
	AccionIniciarRevision  = "iniciar-revision"
	AccionRechazar         = "rechazar"
	AccionRetractarRechazo = "retractar-rechazo"
)

// transicion describe una acción: desde qué estados se permite, a cuál lleva, qué carril la
// ejecuta y qué deben cumplir los submódulos.
type transicion struct {
	desde     []string
	hacia     string
	carril    Carril
	requisito func(estadosSubmodulos []string) string // "" si se cumple; si no, el motivo
}

var transiciones = map[string]transicion{
	AccionEnviar: {
		desde:     []string{EstadoIniciativaEnFormulacion, EstadoIniciativaPorCorregir},
		hacia:     EstadoIniciativaPorAprobar,
		carril:    CarrilFormulador,
		requisito: todasEn("Todos los submódulos deben estar formulados.", EstadoSubmoduloFormulado, EstadoSubmoduloAprobado),
	},
	AccionAprobar: {
		desde:     []string{EstadoIniciativaPorAprobar},
		hacia:     EstadoIniciativaPorRevisar,
		carril:    CarrilJefe,
		requisito: ningunaEn("Hay submódulos con observaciones: devuelva la iniciativa para que se corrijan.", EstadoSubmoduloConObservaciones),
	},
	AccionDevolver: {
		desde:     []string{EstadoIniciativaPorAprobar},
		hacia:     EstadoIniciativaPorCorregir,
		carril:    CarrilJefe,
		requisito: algunaEn("Para devolverla, deje al menos una observación en un submódulo.", EstadoSubmoduloConObservaciones),
	},
	AccionIniciarRevision: {
		desde:     []string{EstadoIniciativaPorRevisar},
		hacia:     EstadoIniciativaEnRevision,
		carril:    CarrilOAP,
		requisito: sinRequisito,
	},
	AccionRechazar: {
		desde:     []string{EstadoIniciativaEnRevision},
		hacia:     EstadoIniciativaRechazada,
		carril:    CarrilOAP,
		requisito: sinRequisito,
	},
	AccionRetractarRechazo: {
		desde:     []string{EstadoIniciativaRechazada},
		hacia:     EstadoIniciativaPorCorregir,
		carril:    CarrilOAP,
		requisito: sinRequisito,
	},
}

// ordenAcciones fija el orden en que se listan las acciones.
var ordenAcciones = []string{AccionEnviar, AccionDevolver, AccionAprobar, AccionIniciarRevision, AccionRechazar, AccionRetractarRechazo}

// nombresAcciones son los nombres con que el usuario ve las acciones (los botones del cliente,
// ACCIONES_INICIATIVA); los mensajes no muestran la clave del API.
var nombresAcciones = map[string]string{
	AccionEnviar:           "Enviar a revisión",
	AccionAprobar:          "Aprobar",
	AccionDevolver:         "Devolver",
	AccionIniciarRevision:  "Iniciar revisión",
	AccionRechazar:         "Rechazar",
	AccionRetractarRechazo: "Retractar rechazo",
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

// EstadoIniciativaTrasGuardar: la iniciativa Sin formular pasa a En formulación en cuanto algún
// submódulo tiene información. En los demás estados no cambia ("" = sin cambio).
func EstadoIniciativaTrasGuardar(estadoIniciativa string, conInformacion bool) string {
	if estadoIniciativa == EstadoIniciativaSinFormular && conInformacion {
		return EstadoIniciativaEnFormulacion
	}
	return ""
}

// EstadoTrasRevision es el cambio automático de una iniciativa En revisión OAP después de cada
// veredicto: con todos los submódulos aprobados pasa a Aprobado; con todos con veredicto y alguno
// con observaciones, a Por corregir. Mientras falte algún veredicto no cambia ("").
func EstadoTrasRevision(estadosSubmodulos []string) string {
	switch {
	case todosEn(estadosSubmodulos, EstadoSubmoduloAprobado):
		return EstadoIniciativaAprobada
	case todosEn(estadosSubmodulos, EstadoSubmoduloAprobado, EstadoSubmoduloConObservaciones):
		return EstadoIniciativaPorCorregir
	default:
		return ""
	}
}

// EvaluarAccion devuelve el estado destino o el motivo por el que no se puede, según el estado de
// la iniciativa y de sus submódulos (no mira quién la ejecuta: eso es del carril, ver permisos.go).
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

// SubmoduloRevisable indica si un submódulo en ese estado se puede observar o aprobar: ya está
// formulado (Formulado, Con observaciones o Aprobado).
func SubmoduloRevisable(estadoSubmodulo string) bool {
	return estadoSubmodulo != EstadoSubmoduloSinFormular && estadoSubmodulo != EstadoSubmoduloFormulando
}

// Requisitos de las transiciones sobre los estados de los submódulos: devuelven "" si se cumplen
// y, si no, el motivo.

func sinRequisito([]string) string { return "" }

func todasEn(motivo string, permitidos ...string) func([]string) string {
	return requisito(motivo, func(estados []string) bool { return todosEn(estados, permitidos...) })
}

func algunaEn(motivo string, buscados ...string) func([]string) string {
	return requisito(motivo, func(estados []string) bool { return algunoEn(estados, buscados...) })
}

func ningunaEn(motivo string, prohibidos ...string) func([]string) string {
	return requisito(motivo, func(estados []string) bool { return !algunoEn(estados, prohibidos...) })
}

func requisito(motivo string, cumple func([]string) bool) func([]string) string {
	return func(estados []string) string {
		if cumple(estados) {
			return ""
		}
		return motivo
	}
}

// todosEn indica si hay estados y todos son de los permitidos.
func todosEn(estados []string, permitidos ...string) bool {
	if len(estados) == 0 {
		return false
	}
	for _, estado := range estados {
		if !contiene(permitidos, estado) {
			return false
		}
	}
	return true
}

// algunoEn indica si alguno de los estados es de los buscados.
func algunoEn(estados []string, buscados ...string) bool {
	for _, estado := range estados {
		if contiene(buscados, estado) {
			return true
		}
	}
	return false
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
