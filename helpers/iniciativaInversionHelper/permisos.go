package iniciativaInversionHelper

import (
	"strings"
)

// Quién puede qué sobre una iniciativa. Cada carril del diagrama (docs/Máquina de estados -
// Iniciativas de inversión) es un grupo de roles; en cada estado solo actúan los carriles que la
// tabla tareasPorEstado (o la transición de la acción) señala, y los demás consultan. Nada se
// asigna a una persona: actúa cualquier usuario con un rol del carril.
//
// Un usuario con varios roles puede todo lo que permita cualquiera de ellos (la suma de sus
// carriles). Cuando más de uno sirve, firma con el de mayor jerarquía (jerarquiaRoles).

// Roles del sistema que actúan sobre las iniciativas (los del token del usuario).
const (
	RolPlaneacion           = "PLANEACION"
	RolAsistentePlaneacion  = "ASISTENTE_PLANEACION"
	RolJefeDependencia      = "JEFE_DEPENDENCIA"
	RolAsistenteDependencia = "ASISTENTE_DEPENDENCIA"
)

// jerarquiaRoles va del rol de mayor jerarquía al de menor.
var jerarquiaRoles = []string{RolPlaneacion, RolAsistentePlaneacion, RolJefeDependencia, RolAsistenteDependencia}

// Carril es un grupo de roles que actúa sobre la iniciativa en ciertos estados.
type Carril struct {
	nombre string
	roles  []string
	// soloSusUnidades: el carril actúa solo sobre las iniciativas de las unidades del usuario.
	soloSusUnidades bool
}

var (
	// CarrilFormulador crea, diligencia, corrige y envía las iniciativas de sus unidades.
	CarrilFormulador = Carril{"formulador", []string{RolJefeDependencia, RolAsistenteDependencia, RolAsistentePlaneacion}, true}
	// CarrilJefe aprueba o devuelve las iniciativas de su unidad cuando están Por aprobar.
	CarrilJefe = Carril{"jefe de unidad", []string{RolJefeDependencia}, true}
	// CarrilOAP revisa las iniciativas de todas las unidades.
	CarrilOAP = Carril{"OAP", []string{RolPlaneacion, RolAsistentePlaneacion}, false}
)

// rolEn devuelve el rol con el que se actúa en el carril: el de mayor jerarquía entre los roles
// dados que son del carril, o "" si ninguno lo es.
func (c Carril) rolEn(roles []string) string {
	for _, rol := range jerarquiaRoles {
		if contiene(c.roles, rol) && contiene(roles, rol) {
			return rol
		}
	}
	return ""
}

// tareas dice qué carril hace cada cosa con los submódulos en un estado de la iniciativa; nil si
// en ese estado nadie la hace.
type tareas struct {
	// editar: guarda los submódulos no aprobados y sube anexos.
	editar *Carril
	// observar: publica observaciones que dejan el submódulo Con observaciones (también un aprobado).
	observar *Carril
	// aprobar: da el submódulo por aprobado.
	aprobar *Carril
}

var tareasPorEstado = map[string]tareas{
	EstadoIniciativaSinFormular:   {editar: &CarrilFormulador},
	EstadoIniciativaEnFormulacion: {editar: &CarrilFormulador},
	EstadoIniciativaPorCorregir:   {editar: &CarrilFormulador},
	EstadoIniciativaPorAprobar:    {observar: &CarrilJefe},
	EstadoIniciativaEnRevision:    {observar: &CarrilOAP, aprobar: &CarrilOAP},
}

// ---------------------------------------------------------------------------
// Quién hace la solicitud
// ---------------------------------------------------------------------------

// Actor es quien hace la solicitud: su documento y todos sus roles.
type Actor struct {
	Usuario string   `json:"usuario"`
	Roles   []string `json:"roles"`
}

// FuenteDependenciasUsuario da los ids de las unidades (dependencias de Oikos) de un usuario; la
// implementa Unidades.
type FuenteDependenciasUsuario interface {
	DependenciasDe(documento string) ([]string, error)
}

// solicitante es el actor de una solicitud. Sus unidades se consultan solo si un carril las
// necesita, y una sola vez por solicitud.
type solicitante struct {
	Actor
	fuente       FuenteDependenciasUsuario
	consultadas  bool
	dependencias map[string]bool
	err          error
}

// solicitante valida el actor: el documento va dentro de consultas a Terceros y los roles deciden
// todo lo demás, así que sin ellos la solicitud no se procesa.
func (s *Servicio) solicitante(actor Actor) (*solicitante, error) {
	actor.Usuario = strings.TrimSpace(actor.Usuario)
	roles := make([]string, 0, len(actor.Roles))
	for _, rol := range actor.Roles {
		if rol = strings.TrimSpace(rol); rol != "" {
			roles = append(roles, rol)
		}
	}
	actor.Roles = roles
	if !patronDocumento.MatchString(actor.Usuario) || len(actor.Roles) == 0 {
		return nil, errorSolicitud("No se pudo identificar su usuario o sus roles. Vuelva a iniciar sesión e intente de nuevo.").
			conDetalle("la solicitud necesita usuario (documento: solo letras y números, máximo 20) y roles (usuario %q, %d roles)", actor.Usuario, len(actor.Roles))
	}
	return &solicitante{Actor: actor, fuente: s.unidades}, nil
}

// esDeSusUnidades indica si la dependencia es una de las unidades del usuario.
func (q *solicitante) esDeSusUnidades(dependenciaID string) (bool, error) {
	if !q.consultadas {
		q.consultadas = true
		q.dependencias = map[string]bool{}
		if q.fuente == nil {
			q.err = errorConfiguracion("el servicio de iniciativas no tiene fuente de unidades de usuario")
		} else {
			var ids []string
			if ids, q.err = q.fuente.DependenciasDe(q.Usuario); q.err == nil {
				for _, id := range ids {
					q.dependencias[id] = true
				}
			}
		}
	}
	return q.dependencias[strings.TrimSpace(dependenciaID)], q.err
}

// rolEn devuelve el rol con el que el usuario actúa en el carril sobre una iniciativa de esa
// dependencia, o "" si no actúa (ningún rol del carril o, si el carril lo exige, no es su unidad).
func (q *solicitante) rolEn(carril *Carril, dependenciaID string) (string, error) {
	if carril == nil {
		return "", nil
	}
	rol := carril.rolEn(q.Roles)
	if rol == "" || !carril.soloSusUnidades {
		return rol, nil
	}
	suya, err := q.esDeSusUnidades(dependenciaID)
	if err != nil || !suya {
		return "", err
	}
	return rol, nil
}

// veTodas indica si el usuario consulta las iniciativas de todas las unidades (la OAP); los demás
// consultan solo las de sus unidades.
func (q *solicitante) veTodas() bool {
	return CarrilOAP.rolEn(q.Roles) != ""
}

// puedeConsultar indica si el usuario puede abrir una iniciativa de esa dependencia.
func (q *solicitante) puedeConsultar(dependenciaID string) (bool, error) {
	if q.veTodas() {
		return true, nil
	}
	return q.esDeSusUnidades(dependenciaID)
}

// autorizar exige que el usuario actúe en el carril sobre la iniciativa y que hoy la unidad de la
// iniciativa esté dentro de sus fechas de formulación. Devuelve el rol con el que actúa.
func (s *Servicio) autorizar(q *solicitante, ctx *contextoIniciativa, carril *Carril) (string, error) {
	rol, err := q.rolEn(carril, ctx.plan.DependenciaID)
	if err != nil {
		return "", err
	}
	if rol == "" {
		return "", errorProhibido("Su rol o su unidad no le permiten hacer esto mientras la iniciativa está «%s».", ctx.estado.Nombre).
			conDetalle("el usuario %s con roles %v no actúa en el carril %s de la iniciativa %s (dependencia %s)",
				q.Usuario, q.Roles, carrilNombre(carril), ctx.plan.ID, ctx.plan.DependenciaID)
	}
	if err := s.exigirFechasHabilitadas(ctx.plan.DependenciaID); err != nil {
		return "", err
	}
	return rol, nil
}

func carrilNombre(carril *Carril) string {
	if carril == nil {
		return "(ninguno)"
	}
	return carril.nombre
}

// ---------------------------------------------------------------------------
// Permisos de quien consulta
// ---------------------------------------------------------------------------

// Lo que se hace con un submódulo desde la vista de la iniciativa y el listado.
const (
	AccionSubmoduloFormular  = "formular"
	AccionSubmoduloRevisar   = "revisar"
	AccionSubmoduloConsultar = "consultar"
)

// permisosDe calcula lo que el usuario puede hacer hoy con la iniciativa. Las fechas se consultan
// solo si el usuario tiene algo que hacer en ese estado; si están cerradas, nada queda habilitado
// y el segundo valor es true. estadosSubmodulos sirve para saber si cada acción está habilitada
// (con nil, las acciones salen sin evaluar sus requisitos, como en el listado).
func (s *Servicio) permisosDe(q *solicitante, plan *Plan, estado string, estadosSubmodulos []string, fechas *fechasDeConsulta) (Permisos, bool, error) {
	permisos := Permisos{Acciones: []AccionDisponible{}}
	tareasDelEstado := tareasPorEstado[estado]

	rolEditar, err := q.rolEn(tareasDelEstado.editar, plan.DependenciaID)
	if err != nil {
		return permisos, false, err
	}
	rolObservar, err := q.rolEn(tareasDelEstado.observar, plan.DependenciaID)
	if err != nil {
		return permisos, false, err
	}
	rolAprobar, err := q.rolEn(tareasDelEstado.aprobar, plan.DependenciaID)
	if err != nil {
		return permisos, false, err
	}
	for _, accion := range ordenAcciones {
		t := transiciones[accion]
		if !contiene(t.desde, estado) {
			continue
		}
		rol, err := q.rolEn(&t.carril, plan.DependenciaID)
		if err != nil {
			return permisos, false, err
		}
		if rol == "" {
			continue
		}
		disponible := AccionDisponible{Accion: accion, Habilitada: true}
		if estadosSubmodulos != nil {
			if _, motivo := EvaluarAccion(accion, estado, estadosSubmodulos); motivo != "" {
				disponible.Habilitada, disponible.Motivo = false, motivo
			}
		}
		permisos.Acciones = append(permisos.Acciones, disponible)
	}

	tieneTareas := rolEditar != "" || rolObservar != "" || rolAprobar != "" || len(permisos.Acciones) > 0
	if tieneTareas && fechas.fueraDeFechas(plan) {
		for i := range permisos.Acciones {
			permisos.Acciones[i].Habilitada, permisos.Acciones[i].Motivo = false, mensajeFueraDeFechas
		}
		return permisos, true, nil
	}
	permisos.Editar = rolEditar != ""
	permisos.Observar = rolObservar != ""
	permisos.AprobarSubmodulos = rolAprobar != ""
	return permisos, false, nil
}

// accionDeSubmodulo es lo que el usuario hace con un submódulo: lo formula si edita la iniciativa y
// el submódulo no está aprobado; lo revisa si observa o aprueba y el submódulo ya está formulado; si
// no, lo consulta.
func accionDeSubmodulo(permisos Permisos, estadoSubmodulo string) string {
	switch {
	case permisos.Editar && estadoSubmodulo != EstadoSubmoduloAprobado:
		return AccionSubmoduloFormular
	case (permisos.Observar || permisos.AprobarSubmodulos) && SubmoduloRevisable(estadoSubmodulo):
		return AccionSubmoduloRevisar
	default:
		return AccionSubmoduloConsultar
	}
}

// accionEnListado es lo que el usuario hace con la iniciativa desde el listado: formularla, revisarla
// (observar, aprobar submódulos o ejecutar una acción de su carril) o consultarla.
func accionEnListado(permisos Permisos) string {
	switch {
	case permisos.Editar:
		return AccionSubmoduloFormular
	case permisos.Observar || permisos.AprobarSubmodulos || algunaHabilitada(permisos.Acciones):
		return AccionSubmoduloRevisar
	default:
		return AccionSubmoduloConsultar
	}
}

func algunaHabilitada(acciones []AccionDisponible) bool {
	for _, accion := range acciones {
		if accion.Habilitada {
			return true
		}
	}
	return false
}
