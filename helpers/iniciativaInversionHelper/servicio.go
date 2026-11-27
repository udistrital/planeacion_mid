package iniciativaInversionHelper

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/astaxie/beego/logs"
)

const (
	aplicativoPlaneacion    = "idPlaneacion"
	nombreIniciativaInicial = "Iniciativa de inversión sin nombre"
)

// Servicio reúne los casos de uso de las iniciativas de inversión.
type Servicio struct {
	repo         Repositorio
	definiciones FuenteDefiniciones
	catalogos    ProveedorCatalogos
	gestor       GestorDocumental
	sistema      *catalogosSistema
	normalizador *Normalizador
	validador    *Validador

	// dependencias consulta el nombre de las dependencias (Oikos) para el listado.
	// Si es nil, el listado sale sin nombre de dependencia.
	dependencias FuenteNombreDependencia

	// terceros consulta en Terceros el nombre de quien crea la iniciativa, comenta o sube
	// un anexo. Si es nil, se guardan sin nombre.
	terceros FuenteNombreTercero

	// unidades da las unidades de cada usuario, con las que se decide qué iniciativas consulta y
	// sobre cuáles actúa como formulador o Jefe (ver permisos.go).
	unidades FuenteDependenciasUsuario
}

// duracionCacheUnidades es cuánto se guardan en memoria las unidades de un usuario.
const duracionCacheUnidades = 5 * time.Minute

// NuevoServicio arma el servicio con sus dependencias.
func NuevoServicio(repo Repositorio, definiciones FuenteDefiniciones, catalogos ProveedorCatalogos, gestor GestorDocumental) *Servicio {
	return &Servicio{
		repo:         repo,
		definiciones: definiciones,
		catalogos:    catalogos,
		gestor:       gestor,
		sistema:      nuevosCatalogosSistema(repo),
		normalizador: NuevoNormalizador(catalogos),
		validador:    NuevoValidador(catalogos),
	}
}

// NuevoServicioPorDefecto conecta con planes_crud, Terceros, Oikos y el gestor documental (hosts
// con esquema, como en conf/app.conf), con caché de definiciones y de catálogos: los de la colección
// elemento_catalogo; para el catálogo sedes, los espacios físicos de tipo SEDE de Oikos, y para
// los del Plan Distrital de Desarrollo, el Plan Estratégico de Desarrollo y el Plan Indicativo,
// el árbol de esos planes. Oikos también da el nombre de la dependencia en el listado de
// iniciativas, y Terceros, el de quien crea la iniciativa, comenta o sube un anexo, y las unidades
// de cada usuario.
func NuevoServicioPorDefecto(hostCrud, hostTerceros, hostOikos, hostGestor string) *Servicio {
	crud := NuevoClienteCrud(hostCrud)
	definiciones := NuevaCacheDefiniciones(NuevoCargadorDefiniciones(crud, crud, 8))
	catalogos := NuevoRepositorioConInstrumentos(NuevoRepositorioConSedes(crud, NuevoClienteSedesOikos(hostOikos)), crud, time.Minute, 8)
	servicio := NuevoServicio(crud, definiciones, NuevosCatalogosConCache(catalogos, 10*time.Minute), NuevoClienteGestorDocumental(hostGestor))
	personas := NuevoClienteUnidades(hostTerceros, hostOikos)
	servicio.dependencias = personas
	servicio.terceros = personas
	servicio.unidades = NuevaCacheDependencias(NuevasUnidades(personas), duracionCacheUnidades)
	return servicio
}

// ---------------------------------------------------------------------------
// Plantillas y iniciativas
// ---------------------------------------------------------------------------

// Plantillas lista las plantillas activas con las que se pueden crear
// iniciativas, de la más reciente a la más antigua.
func (s *Servicio) Plantillas() ([]PlantillaResumen, error) {
	tipo, err := s.sistema.TipoPlanIniciativa()
	if err != nil {
		return nil, err
	}
	planes, err := s.repo.BuscarPlanes(map[string]string{"tipo_plan_id": tipo.ID, "formato": "true", "activo": "true"})
	if err != nil {
		return nil, err
	}
	sort.SliceStable(planes, func(i, j int) bool { return planes[i].FechaCreacion > planes[j].FechaCreacion })
	plantillas := make([]PlantillaResumen, 0, len(planes))
	for _, plan := range planes {
		if motivo := s.motivoPlantillaNoUsable(plan.ID); motivo != "" {
			logs.Info("iniciativaInversion: la plantilla %s (%s) no se ofrece para crear iniciativas: %s", plan.ID, plan.Nombre, motivo)
			continue
		}
		plantillas = append(plantillas, PlantillaResumen{ID: plan.ID, Nombre: plan.Nombre, Descripcion: plan.Descripcion})
	}
	return plantillas, nil
}

// plantillaVigente es la plantilla usable más reciente: con ella se crean las
// iniciativas cuando no se indica otra.
func (s *Servicio) plantillaVigente() (string, error) {
	plantillas, err := s.Plantillas()
	if err != nil {
		return "", err
	}
	if len(plantillas) == 0 {
		return "", errorConflicto("No hay una plantilla de iniciativas de inversión disponible para crear la iniciativa. Comuníquese con el administrador del sistema.").
			conDetalle("no hay plantillas %s activas con submódulos. Se carga con database local/construccion-proyecto-inversion/cargar-construccion.js.", CodigoTipoPlanIniciativa)
	}
	return plantillas[0].ID, nil
}

// motivoPlantillaNoUsable devuelve por qué una plantilla no sirve para crear
// iniciativas (p. ej. una en construcción sin módulos ni submódulos), o "".
func (s *Servicio) motivoPlantillaNoUsable(plantillaID string) string {
	definicion, err := s.definiciones.Definicion(plantillaID)
	if err != nil {
		return err.Error()
	}
	if len(definicion.Submodulos()) == 0 {
		return "no tiene submódulos"
	}
	return ""
}

// ListarIniciativas lista las iniciativas activas que el usuario consulta (la OAP, todas; los
// demás, las de sus unidades), de la más reciente a la más antigua, con su número, su estado
// general, lo que el usuario hace con cada una y si está fuera de las fechas de su unidad.
func (s *Servicio) ListarIniciativas(actor Actor) ([]IniciativaListado, error) {
	q, err := s.solicitante(actor)
	if err != nil {
		return nil, err
	}
	tipo, err := s.sistema.TipoPlanIniciativa()
	if err != nil {
		return nil, err
	}
	// Sin filtrar por activo: el número se cuenta sobre todas las iniciativas.
	planes, err := s.repo.BuscarPlanes(map[string]string{"tipo_plan_id": tipo.ID, "formato": "false"})
	if err != nil {
		return nil, err
	}
	numeros := numerarPorCreacion(planes)

	nombresDeDependencias := nombresDependencias{fuente: s.dependencias, nombres: map[string]string{}}
	fechas := fechasDeConsulta{servicio: s}
	iniciativas := make([]IniciativaListado, 0, len(planes))
	for i := range planes {
		plan := &planes[i]
		if !plan.Activo {
			continue
		}
		visible, err := q.puedeConsultar(plan.DependenciaID)
		if err != nil {
			return nil, err
		}
		if !visible {
			continue
		}
		estado, err := s.sistema.EstadoIniciativaPorID(plan.EstadoPlanID)
		if err != nil {
			return nil, err
		}
		permisos, fueraDeFechas, err := s.permisosDe(q, plan, estado.Codigo, nil, &fechas)
		if err != nil {
			return nil, err
		}
		iniciativas = append(iniciativas, IniciativaListado{
			ID:                plan.ID,
			Numero:            numeros[plan.ID],
			Nombre:            plan.Nombre,
			DependenciaID:     plan.DependenciaID,
			DependenciaNombre: nombresDeDependencias.de(plan.ID, plan.DependenciaID),
			Responsable:       plan.ResponsableNombre,
			Estado:            estado,
			Accion:            accionEnListado(permisos),
			FueraDeFechas:     fueraDeFechas,
			FechaCreacion:     plan.FechaCreacion,
		})
	}
	sort.SliceStable(iniciativas, func(i, j int) bool { return iniciativas[i].Numero > iniciativas[j].Numero })
	return iniciativas, nil
}

// FuenteNombreTercero devuelve el nombre completo de quien tiene un documento en Terceros (la
// implementa ClienteUnidades).
type FuenteNombreTercero interface {
	NombreTercero(documento string) (string, error)
}

// nombreTercero devuelve el nombre en Terceros de quien tiene el documento. Si falla, quien
// llama guarda sin nombre y deja el motivo en el log.
func (s *Servicio) nombreTercero(documento string) (string, error) {
	if s.terceros == nil {
		return "", errors.New("el servicio no tiene cliente de Terceros")
	}
	return s.terceros.NombreTercero(documento)
}

// FuenteNombreDependencia devuelve el nombre de una dependencia (la implementa ClienteUnidades).
type FuenteNombreDependencia interface {
	NombreDependencia(dependenciaID int) (string, error)
}

// nombresDependencias guarda, durante un listado, el nombre de cada dependencia según
// Oikos, para consultar cada una una sola vez.
type nombresDependencias struct {
	fuente  FuenteNombreDependencia
	nombres map[string]string
}

// de devuelve el nombre de la dependencia de la iniciativa planID. Queda vacío si no hay
// fuente, si el id no es un número o si Oikos falla; en esos casos el motivo va al log
// y el listado sale igual.
func (n nombresDependencias) de(planID, dependenciaID string) string {
	dependenciaID = strings.TrimSpace(dependenciaID)
	if n.fuente == nil || dependenciaID == "" {
		return ""
	}
	if nombre, consultada := n.nombres[dependenciaID]; consultada {
		return nombre
	}
	nombre := ""
	id, err := strconv.Atoi(dependenciaID) // NombreDependencia recibe un int
	if err == nil {
		nombre, err = n.fuente.NombreDependencia(id)
	}
	if err != nil {
		logs.Warn("iniciativaInversion: sin nombre para la dependencia %q de la iniciativa %s: %v", dependenciaID, planID, err)
		nombre = ""
	}
	n.nombres[dependenciaID] = nombre
	return nombre
}

// numerarPorCreacion da a cada plan su consecutivo por fecha de creación (y por
// _id si dos tienen la misma fecha), empezando en 1.
func numerarPorCreacion(planes []Plan) map[string]int {
	ordenados := append([]Plan(nil), planes...)
	sort.SliceStable(ordenados, func(i, j int) bool {
		if ordenados[i].FechaCreacion != ordenados[j].FechaCreacion {
			return ordenados[i].FechaCreacion < ordenados[j].FechaCreacion
		}
		return ordenados[i].ID < ordenados[j].ID
	})
	numeros := make(map[string]int, len(ordenados))
	for i, plan := range ordenados {
		numeros[plan.ID] = i + 1
	}
	return numeros
}

// CrearIniciativa crea el plan de la iniciativa y un submódulo vacío por cada
// submódulo de la plantilla (la indicada o, si no se indica, la más reciente). Si
// falla la creación de los submódulos, desactiva el plan. La crea el carril del
// formulador, para una de sus unidades (403 si no) y dentro de las fechas de
// formulación de esa unidad (409 si no).
//
// El plan guarda el documento del responsable (quien la crea) y su nombre en Terceros,
// que el listado muestra sin volver a consultarlo. Si no se obtiene el nombre, la
// iniciativa se crea igual, sin él, y el motivo va al log.
func (s *Servicio) CrearIniciativa(solicitud SolicitudCrearIniciativa) (*ResumenIniciativa, error) {
	dependenciaID := strings.TrimSpace(solicitud.DependenciaID)
	if dependenciaID == "" {
		return nil, errorSolicitud("Seleccione la unidad para la que se crea la iniciativa.")
	}
	q, err := s.solicitante(solicitud.Actor)
	if err != nil {
		return nil, err
	}
	vigencias, err := textoVigencias(solicitud.Vigencias)
	if err != nil {
		return nil, err
	}
	rol, err := q.rolEn(&CarrilFormulador, dependenciaID)
	if err != nil {
		return nil, err
	}
	if rol == "" {
		return nil, errorProhibido("Solo puede crear iniciativas para sus unidades y con un rol de formulador.").
			conDetalle("el usuario %s con roles %v no formula en la dependencia %s", q.Usuario, q.Roles, dependenciaID)
	}
	if err := s.exigirFechasHabilitadas(dependenciaID); err != nil {
		return nil, err
	}
	responsable := q.Usuario
	tipo, err := s.sistema.TipoPlanIniciativa()
	if err != nil {
		return nil, err
	}
	plantillaID := strings.TrimSpace(solicitud.PlantillaID)
	if plantillaID == "" {
		if plantillaID, err = s.plantillaVigente(); err != nil {
			return nil, err
		}
	}
	definicion, err := s.definiciones.Definicion(plantillaID)
	if err != nil {
		return nil, err
	}
	if definicion.TipoPlanID != tipo.ID {
		return nil, errorSolicitud("La plantilla indicada no es de iniciativas de inversión.").
			conDetalle("la plantilla %s tiene tipo_plan_id %s y se esperaba %s", plantillaID, definicion.TipoPlanID, tipo.ID)
	}
	submodulos := definicion.Submodulos()
	if len(submodulos) == 0 {
		return nil, errorSolicitud("La plantilla indicada no tiene submódulos.").conDetalle("la plantilla %s no tiene submódulos", plantillaID)
	}
	iniciativaSinFormular, err := s.sistema.EstadoIniciativa(EstadoIniciativaSinFormular)
	if err != nil {
		return nil, err
	}
	sinFormular, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloSinFormular)
	if err != nil {
		return nil, err
	}

	nombre := strings.TrimSpace(solicitud.Nombre)
	if nombre == "" {
		nombre = nombreIniciativaInicial
	}
	responsableNombre, err := s.nombreTercero(responsable)
	if err != nil {
		logs.Warn("iniciativaInversion: la iniciativa de la dependencia %s se crea sin el nombre de su responsable: %v", dependenciaID, err)
	}
	creado, err := s.repo.CrearPlan(Plan{
		Nombre:            nombre,
		Descripcion:       "Iniciativa de inversión formulada con la plantilla «" + definicion.Nombre + "»",
		TipoPlanID:        tipo.ID,
		AplicativoID:      aplicativoPlaneacion,
		Responsable:       responsable,
		ResponsableNombre: responsableNombre,
		Activo:            true,
		Formato:           false,
		FormatoID:         definicion.PlantillaID,
		DependenciaID:     dependenciaID,
		EstadoPlanID:      iniciativaSinFormular.ID,
		VigenciaAplica:    vigencias,
	})
	if err != nil {
		return nil, err
	}

	nuevos := make([]NuevoSubmodulo, 0, len(submodulos))
	for _, submodulo := range submodulos {
		nuevos = append(nuevos, NuevoSubmodulo{PlanID: creado.ID, Submodulo: submodulo.Clave(), EstadoSubmoduloID: sinFormular.ID, Datos: Datos{}})
	}
	if _, err := s.repo.CrearSubmodulos(nuevos); err != nil {
		if errDesactivar := s.repo.DesactivarPlan(creado.ID); errDesactivar != nil {
			logs.Error("iniciativaInversion: no se pudo desactivar el plan %s tras fallar sus submódulos: %v", creado.ID, errDesactivar)
		}
		return nil, err
	}
	return s.Resumen(creado.ID, q.Actor)
}

// textoVigencias valida que las vigencias sean una lista JSON y la devuelve como
// texto, igual que vigencia_aplica de PED y PI.
func textoVigencias(vigencias json.RawMessage) (string, error) {
	if len(vigencias) == 0 || string(vigencias) == "null" {
		return "", nil
	}
	var lista []interface{}
	if err := json.Unmarshal(vigencias, &lista); err != nil {
		return "", errorFormatoSolicitud("vigencias debe ser una lista JSON: %v", err)
	}
	return string(vigencias), nil
}

// Resumen arma la vista general: estado de la iniciativa y de cada submódulo, por módulo, y lo que
// el usuario puede hacer con ella y con cada submódulo. Quien no es de la OAP solo consulta las
// iniciativas de sus unidades (403 si no).
func (s *Servicio) Resumen(iniciativaID string, actor Actor) (*ResumenIniciativa, error) {
	q, err := s.solicitante(actor)
	if err != nil {
		return nil, err
	}
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	visible, err := q.puedeConsultar(ctx.plan.DependenciaID)
	if err != nil {
		return nil, err
	}
	if !visible {
		return nil, errorProhibido("Solo puede consultar las iniciativas de sus unidades.").
			conDetalle("el usuario %s con roles %v no ve la dependencia %s de la iniciativa %s", q.Usuario, q.Roles, ctx.plan.DependenciaID, ctx.plan.ID)
	}
	codigos, err := ctx.codigosEstadoSubmodulos(s.sistema)
	if err != nil {
		return nil, err
	}
	permisos, fueraDeFechas, err := s.permisosDe(q, ctx.plan, ctx.estado.Codigo, codigos, &fechasDeConsulta{servicio: s})
	if err != nil {
		return nil, err
	}
	dependientes := dependientesPorSubmodulo(ctx.definicion)
	modulos := make([]ResumenModulo, 0, len(ctx.definicion.Modulos))
	for _, modulo := range ctx.definicion.Modulos {
		resumenModulo := ResumenModulo{Clave: modulo.Clave(), Nombre: modulo.Nombre, Submodulos: []ResumenSubmodulo{}}
		for _, nodo := range modulo.Hijos {
			documento, err := ctx.documento(nodo)
			if err != nil {
				return nil, err
			}
			estado, err := s.sistema.EstadoSubmoduloPorID(documento.EstadoSubmoduloID)
			if err != nil {
				return nil, err
			}
			resumenModulo.Submodulos = append(resumenModulo.Submodulos, ResumenSubmodulo{
				Clave:             nodo.Clave(),
				Nombre:            nodo.Nombre,
				Estado:            estado,
				Accion:            accionDeSubmodulo(permisos, estado.Codigo),
				Observaciones:     len(documento.Observaciones),
				UltimaObservacion: ultimaObservacionPorAtender(estado.Codigo, documento.Observaciones),
				Dependientes:      dependientes[nodo.Clave()],
			})
		}
		modulos = append(modulos, resumenModulo)
	}
	return &ResumenIniciativa{
		ID:            ctx.plan.ID,
		Nombre:        ctx.plan.Nombre,
		DependenciaID: ctx.plan.DependenciaID,
		Estado:        ctx.estado,
		FueraDeFechas: fueraDeFechas,
		Permisos:      permisos,
		Modulos:       modulos,
	}, nil
}

// ultimaObservacionPorAtender es la observación más reciente de un submódulo Con observaciones;
// nil en los demás estados.
func ultimaObservacionPorAtender(estadoSubmodulo string, observaciones []Observacion) *Observacion {
	if estadoSubmodulo != EstadoSubmoduloConObservaciones || len(observaciones) == 0 {
		return nil
	}
	ultima := observaciones[len(observaciones)-1]
	return &ultima
}

// dependientesPorSubmodulo da, por clave de submódulo, las claves de los submódulos que leen datos
// de él, en el orden de la plantilla (nunca nil).
func dependientesPorSubmodulo(definicion *Definicion) map[string][]string {
	submodulos := definicion.Submodulos()
	dependientes := make(map[string][]string, len(submodulos))
	for _, submodulo := range submodulos {
		dependientes[submodulo.Clave()] = []string{}
	}
	for _, submodulo := range submodulos {
		for origen := range dependenciasDe(submodulo) {
			if lista, existe := dependientes[origen]; existe {
				dependientes[origen] = append(lista, submodulo.Clave())
			}
		}
	}
	return dependientes
}

// Definicion devuelve la construcción de la plantilla de la iniciativa.
func (s *Servicio) Definicion(iniciativaID string) (*Definicion, error) {
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	return ctx.definicion, nil
}

// ---------------------------------------------------------------------------
// Submodulos
// ---------------------------------------------------------------------------

// ObtenerSubmodulo devuelve el submódulo, sus errores vigentes y los datos de los
// submódulos de los que depende (para armar espejos y referencias).
func (s *Servicio) ObtenerSubmodulo(iniciativaID, clave string) (*SubmoduloConRelacionadas, error) {
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	nodo, documento, err := ctx.submodulo(clave)
	if err != nil {
		return nil, err
	}
	iniciativa := ctx.datosIniciativa()
	vista, err := s.vista(*documento, s.validador.ValidarSubmodulo(nodo, iniciativa))
	if err != nil {
		return nil, err
	}
	relacionadas := map[string]Datos{}
	for dependencia := range dependenciasDe(nodo) {
		relacionadas[dependencia] = iniciativa[dependencia]
	}
	if catalogos := catalogosEspejadosDe(nodo); len(catalogos) > 0 {
		filas := Datos{}
		for nombre := range catalogos {
			catalogo, err := s.catalogos.Catalogo(nombre)
			if err != nil {
				// El formulario se muestra igual; la tabla queda sin filas y el validador avisa.
				logs.Warn("iniciativaInversion: no se pudo consultar el catálogo %s para el submódulo %s: %v", nombre, clave, err)
				continue
			}
			filas[nombre] = filasDeCatalogo(catalogo)
		}
		relacionadas[claveRelacionadasCatalogo] = filas
	}
	return &SubmoduloConRelacionadas{Submodulo: *vista, Relacionadas: relacionadas}, nil
}

// catalogosEspejadosDe devuelve los catálogos de los que salen las filas de algún espejo del
// submódulo (espejo_de "catalogo.<nombre>").
func catalogosEspejadosDe(submodulo *Nodo) map[string]bool {
	catalogos := map[string]bool{}
	var recorrer func(nodos []*Nodo)
	recorrer = func(nodos []*Nodo) {
		for _, nodo := range nodos {
			if nombre, esCatalogo := catalogoDeRuta(nodo.Dato.EspejoDe); esCatalogo {
				catalogos[nombre] = true
			}
			recorrer(nodo.Hijos)
		}
	}
	recorrer(submodulo.Hijos)
	return catalogos
}

// GuardarSubmodulo guarda los datos aunque tengan errores (borrador). El submódulo
// queda Formulado solo si no tiene errores; si no, Formulando (con información) o Sin
// formular (vacío). La iniciativa pasa de Sin formular a En formulación con la primera
// información. Responde 409 si el estado no deja editarla, si el submódulo está aprobado, si
// otra persona lo cambió desde que se consultó o si está fuera de las fechas de su unidad, y 403
// si el usuario no la formula.
func (s *Servicio) GuardarSubmodulo(iniciativaID, clave string, solicitud SolicitudGuardarSubmodulo) (*SubmoduloVista, error) {
	ctx, nodo, documento, err := s.submoduloParaEditar(iniciativaID, clave, solicitud.Actor)
	if err != nil {
		return nil, err
	}

	datos := s.normalizador.Normalizar(nodo, solicitud.Datos)
	iniciativa := ctx.datosIniciativa()
	iniciativa[clave] = datos
	errores := s.validador.ValidarSubmodulo(nodo, iniciativa)
	conInformacion := tieneInformacion(nodo.Hijos, datos)
	estado, err := s.sistema.EstadoSubmodulo(EstadoTrasGuardar(errores, conInformacion))
	if err != nil {
		return nil, err
	}
	guardada, err := s.repo.ActualizarSubmodulo(documento.ID, CambiosSubmodulo{
		Version:           solicitud.Version,
		Datos:             &datos,
		EstadoSubmoduloID: estado.ID,
	})
	if err != nil {
		return nil, err
	}

	s.revalidarDependientes(ctx, clave, iniciativa)
	s.actualizarEstadoIniciativa(ctx, conInformacion)
	s.actualizarNombrePlan(ctx, nodo, datos)
	return s.vista(*guardada, errores)
}

// submoduloParaEditar carga la iniciativa y el submódulo que el usuario quiere guardar o al que
// quiere subir un anexo: el estado debe dejar editarla, el usuario debe formularla (su carril y su
// unidad), la unidad debe estar dentro de sus fechas y el submódulo no puede estar aprobado.
func (s *Servicio) submoduloParaEditar(iniciativaID, clave string, actor Actor) (*contextoIniciativa, *Nodo, *SubmoduloIniciativa, error) {
	q, err := s.solicitante(actor)
	if err != nil {
		return nil, nil, nil, err
	}
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, nil, nil, err
	}
	carril := tareasPorEstado[ctx.estado.Codigo].editar
	if carril == nil {
		return nil, nil, nil, errorConflicto("La iniciativa está «%s»; no se puede editar.", ctx.estado.Nombre)
	}
	if _, err := s.autorizar(q, ctx, carril); err != nil {
		return nil, nil, nil, err
	}
	nodo, documento, err := ctx.submodulo(clave)
	if err != nil {
		return nil, nil, nil, err
	}
	estado, err := s.sistema.EstadoSubmoduloPorID(documento.EstadoSubmoduloID)
	if err != nil {
		return nil, nil, nil, err
	}
	if estado.Codigo == EstadoSubmoduloAprobado {
		return nil, nil, nil, errorConflicto("El submódulo está aprobado; no se puede editar.")
	}
	return ctx, nodo, documento, nil
}

// RevisarSubmodulo es el veredicto «Aprobado» de la OAP sobre un submódulo formulado, mientras
// revisa la iniciativa (05). El veredicto «Con observaciones» se da con AgregarObservacion, y el
// historial de observaciones no se borra. Después de cada veredicto la iniciativa puede cambiar de
// estado sola (ver cerrarRevision).
func (s *Servicio) RevisarSubmodulo(iniciativaID, clave string, solicitud SolicitudRevisarSubmodulo) (*SubmoduloVista, error) {
	q, err := s.solicitante(solicitud.Actor)
	if err != nil {
		return nil, err
	}
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	carril := tareasPorEstado[ctx.estado.Codigo].aprobar
	if carril == nil {
		return nil, errorConflicto("Los submódulos se aprueban mientras la OAP revisa la iniciativa; esta está «%s».", ctx.estado.Nombre)
	}
	if _, err := s.autorizar(q, ctx, carril); err != nil {
		return nil, err
	}
	nodo, documento, err := ctx.submodulo(clave)
	if err != nil {
		return nil, err
	}
	if err := s.exigirRevisable(nodo, documento); err != nil {
		return nil, err
	}

	aprobado, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloAprobado)
	if err != nil {
		return nil, err
	}
	guardada, err := s.repo.ActualizarSubmodulo(documento.ID, CambiosSubmodulo{
		Version:           solicitud.Version,
		EstadoSubmoduloID: aprobado.ID,
	})
	if err != nil {
		return nil, err
	}
	vista, err := s.vista(*guardada, s.validador.ValidarSubmodulo(nodo, ctx.datosIniciativa()))
	if err != nil {
		return nil, err
	}
	vista.EstadoIniciativa = s.cerrarRevision(iniciativaID)
	return vista, nil
}

// cerrarRevision aplica el cambio automático de una iniciativa En revisión OAP tras un veredicto
// (EstadoTrasRevision) y devuelve el estado al que pasó, o nil si no cambió. Relee la iniciativa
// para ver los veredictos de todos los revisores: si dos dan a la vez los dos últimos, el segundo
// la encuentra ya fuera de En revisión OAP y no hace nada. Si falla, el veredicto ya quedó: solo
// se avisa en el log, y el cambio se aplica con el siguiente veredicto.
func (s *Servicio) cerrarRevision(iniciativaID string) *Estado {
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		logs.Warn("iniciativaInversion: no se pudo releer la iniciativa %s para cerrar la revisión: %v", iniciativaID, err)
		return nil
	}
	if ctx.estado.Codigo != EstadoIniciativaEnRevision {
		return nil
	}
	codigos, err := ctx.codigosEstadoSubmodulos(s.sistema)
	if err != nil {
		logs.Warn("iniciativaInversion: no se pudo cerrar la revisión de la iniciativa %s: %v", iniciativaID, err)
		return nil
	}
	destino := EstadoTrasRevision(codigos)
	if destino == "" {
		return nil
	}
	estado, err := s.cambiarEstado(ctx, destino)
	if err != nil {
		logs.Warn("iniciativaInversion: no se pudo cerrar la revisión de la iniciativa %s: %v", iniciativaID, err)
		return nil
	}
	return &estado
}

// EjecutarAccion cambia el estado general de la iniciativa con una acción (ver transiciones). La
// ejecuta el carril de la acción, dentro de las fechas de la unidad. Enviar revalida los
// submódulos; rechazar exige el motivo, que no se guarda (va en la notificación).
func (s *Servicio) EjecutarAccion(iniciativaID, accion string, solicitud SolicitudAccion) (*ResumenIniciativa, error) {
	t, existe := transiciones[accion]
	if !existe {
		return nil, errorFormatoSolicitud("la acción %q no existe; las válidas son %s", accion, strings.Join(ordenAcciones, ", "))
	}
	if err := motivoValido(accion, solicitud.Motivo); err != nil {
		return nil, err
	}
	q, err := s.solicitante(solicitud.Actor)
	if err != nil {
		return nil, err
	}
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	codigos, err := ctx.codigosEstadoSubmodulos(s.sistema)
	if err != nil {
		return nil, err
	}
	destino, motivo := EvaluarAccion(accion, ctx.estado.Codigo, codigos)
	if motivo != "" {
		return nil, errorConflicto("%s", motivo)
	}
	if _, err := s.autorizar(q, ctx, &t.carril); err != nil {
		return nil, err
	}
	if accion == AccionEnviar {
		if err := s.exigirSubmodulosValidos(ctx); err != nil {
			return nil, err
		}
	}
	if _, err := s.cambiarEstado(ctx, destino); err != nil {
		return nil, err
	}
	return s.Resumen(iniciativaID, q.Actor)
}

// motivoValido exige el motivo al rechazar, con el mismo largo máximo que una observación.
func motivoValido(accion, motivo string) error {
	if accion != AccionRechazar {
		return nil
	}
	motivo = strings.TrimSpace(motivo)
	if motivo == "" {
		return errorSolicitud("Escriba el motivo del rechazo.")
	}
	if utf8.RuneCountInString(motivo) > MaximoTextoObservacion {
		return errorSolicitud("El motivo no puede superar %d caracteres.", MaximoTextoObservacion)
	}
	return nil
}

// cambiarEstado deja la iniciativa en el estado con ese código y devuelve el registro del estado.
func (s *Servicio) cambiarEstado(ctx *contextoIniciativa, codigo string) (Estado, error) {
	estado, err := s.sistema.EstadoIniciativa(codigo)
	if err != nil {
		return Estado{}, err
	}
	plan := *ctx.plan
	plan.EstadoPlanID = estado.ID
	if err := s.repo.ActualizarPlan(plan); err != nil {
		return Estado{}, err
	}
	ctx.plan.EstadoPlanID = estado.ID
	ctx.estado = estado
	return estado, nil
}

// exigirSubmodulosValidos vuelve a validar los submódulos no aprobados antes de enviar la
// iniciativa, porque el estado guardado puede no estar al día: si al guardar un
// submódulo falló devolver a Formulando uno que depende de él (revalidarDependientes),
// ese sigue Formulado con datos que ya no son válidos. Los aprobados no se revalidan: si una
// corrección los dejó desactualizados, quien revisa decide si los vuelve a observar.
func (s *Servicio) exigirSubmodulosValidos(ctx *contextoIniciativa) error {
	aprobado, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloAprobado)
	if err != nil {
		return err
	}
	iniciativa := ctx.datosIniciativa()
	conPendientes := []string{}
	for _, nodo := range ctx.definicion.Submodulos() {
		if ctx.submodulos[nodo.Clave()].EstadoSubmoduloID == aprobado.ID {
			continue
		}
		if len(s.validador.ValidarSubmodulo(nodo, iniciativa)) > 0 {
			conPendientes = append(conPendientes, "«"+nodo.Nombre+"»")
		}
	}
	switch len(conPendientes) {
	case 0:
		return nil
	case 1:
		return errorConflicto("El submódulo %s tiene pendientes. Ábralo, corrija lo señalado y guárdelo de nuevo.", conPendientes[0])
	default:
		ultimo := len(conPendientes) - 1
		return errorConflicto("Los submódulos %s y %s tienen pendientes. Ábralos, corrija lo señalado y guárdelos de nuevo.",
			strings.Join(conPendientes[:ultimo], ", "), conPendientes[ultimo])
	}
}

// Catalogo devuelve los elementos de un catálogo; con padre, solo sus hijos.
func (s *Servicio) Catalogo(nombre, padre string) ([]ItemCatalogo, error) {
	catalogo, err := s.catalogos.Catalogo(nombre)
	if err != nil {
		return nil, err
	}
	return catalogo.Filtrar(padre), nil
}

// ---------------------------------------------------------------------------
// Apoyo
// ---------------------------------------------------------------------------

// revalidarDependientes devuelve a Formulando (o a Sin formular, si están vacíos) los
// submódulos que dependen del que cambió y dejaron de ser válidos (p. ej. Objetivos si se
// agregó una causa). Los que ya estaban sin formular o formulando no cambian, y los aprobados
// tampoco: nadie los edita, y quien revisa decide si los vuelve a observar.
func (s *Servicio) revalidarDependientes(ctx *contextoIniciativa, cambiada string, iniciativa map[string]Datos) {
	sinFormular, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloSinFormular)
	if err != nil {
		return
	}
	formulando, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloFormulando)
	if err != nil {
		return
	}
	aprobado, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloAprobado)
	if err != nil {
		return
	}
	for _, nodo := range ctx.definicion.Submodulos() {
		clave := nodo.Clave()
		documento, existe := ctx.submodulos[clave]
		if clave == cambiada || !existe || !dependenciasDe(nodo)[cambiada] {
			continue
		}
		if actual := documento.EstadoSubmoduloID; actual == sinFormular.ID || actual == formulando.ID || actual == aprobado.ID {
			continue
		}
		if len(s.validador.ValidarSubmodulo(nodo, iniciativa)) == 0 {
			continue
		}
		nuevo := sinFormular
		if tieneInformacion(nodo.Hijos, iniciativa[clave]) {
			nuevo = formulando
		}
		cambios := CambiosSubmodulo{Version: documento.Version, EstadoSubmoduloID: nuevo.ID}
		if _, err := s.repo.ActualizarSubmodulo(documento.ID, cambios); err != nil {
			logs.Warn("iniciativaInversion: no se pudo marcar el submódulo %s como %s: %v", clave, nuevo.Nombre, err)
		}
	}
}

// actualizarEstadoIniciativa pasa la iniciativa de Sin formular a En formulación cuando un
// submódulo se guarda con información. Si falla, el guardado del submódulo ya quedó: solo se avisa
// en el log, y se vuelve a intentar en el siguiente guardado. cambiarEstado deja al día el plan
// del contexto: actualizarNombrePlan parte de él y no debe devolver el estado.
func (s *Servicio) actualizarEstadoIniciativa(ctx *contextoIniciativa, conInformacion bool) {
	destino := EstadoIniciativaTrasGuardar(ctx.estado.Codigo, conInformacion)
	if destino == "" {
		return
	}
	if _, err := s.cambiarEstado(ctx, destino); err != nil {
		logs.Warn("iniciativaInversion: no se pudo pasar la iniciativa %s a %s: %v", ctx.plan.ID, destino, err)
	}
}

// actualizarNombrePlan copia al plan el valor del campo marcado con nombre_plan.
func (s *Servicio) actualizarNombrePlan(ctx *contextoIniciativa, submodulo *Nodo, datos Datos) {
	for _, hijo := range submodulo.Hijos {
		if !hijo.Dato.NombrePlan {
			continue
		}
		nombre := strings.TrimSpace(textoPlano(datos[hijo.Clave()]))
		if nombre == "" || nombre == ctx.plan.Nombre {
			return
		}
		plan := *ctx.plan
		plan.Nombre = nombre
		if err := s.repo.ActualizarPlan(plan); err != nil {
			logs.Warn("iniciativaInversion: no se pudo actualizar el nombre del plan %s: %v", plan.ID, err)
		}
		return
	}
}

func (s *Servicio) vista(documento SubmoduloIniciativa, errores []ErrorCampo) (*SubmoduloVista, error) {
	estado, err := s.sistema.EstadoSubmoduloPorID(documento.EstadoSubmoduloID)
	if err != nil {
		return nil, err
	}
	datos := documento.Datos
	if datos == nil {
		datos = Datos{}
	}
	observaciones := documento.Observaciones
	if observaciones == nil {
		observaciones = []Observacion{}
	}
	return &SubmoduloVista{
		Estado:        estado,
		Datos:         datos,
		Observaciones: observaciones,
		Version:       documento.Version,
		Errores:       errores,
	}, nil
}

// dependenciasDe devuelve los submódulos de los que lee datos un submódulo
// (por espejo_de, opciones_desde absolutos o un calculado con desde).
func dependenciasDe(submodulo *Nodo) map[string]bool {
	dependencias := map[string]bool{}
	var recorrer func(nodos []*Nodo)
	recorrer = func(nodos []*Nodo) {
		for _, nodo := range nodos {
			for _, ruta := range []string{nodo.Dato.EspejoDe, nodo.Dato.OpcionesDesde, nodo.Dato.Desde} {
				if ruta == "" || strings.HasPrefix(ruta, prefijoOrigen) || strings.HasPrefix(ruta, prefijoCatalogo) {
					continue
				}
				if otra := strings.SplitN(ruta, ".", 2)[0]; otra != submodulo.Clave() {
					dependencias[otra] = true
				}
			}
			recorrer(nodo.Hijos)
		}
	}
	recorrer(submodulo.Hijos)
	return dependencias
}

// ---------------------------------------------------------------------------
// Contexto de una iniciativa
// ---------------------------------------------------------------------------

type contextoIniciativa struct {
	plan       *Plan
	definicion *Definicion
	estado     Estado
	submodulos map[string]SubmoduloIniciativa
}

// cargarIniciativa carga la iniciativa: un plan activo de tipo IPI_SP que no es
// plantilla y apunta a la suya (formato_id). El plan y los submódulos se piden a la
// vez porque los dos solo necesitan el id; si el plan falla, manda su error.
func (s *Servicio) cargarIniciativa(iniciativaID string) (*contextoIniciativa, error) {
	var lista []SubmoduloIniciativa
	var errLista error
	listados := make(chan struct{})
	go func() {
		defer close(listados)
		lista, errLista = s.repo.ListarSubmodulos(iniciativaID)
	}()
	plan, err := s.repo.ObtenerPlan(iniciativaID)
	<-listados
	if err != nil {
		if ComoErrorNegocio(err).Status == http.StatusNotFound {
			return nil, iniciativaNoExiste()
		}
		return nil, err
	}
	tipo, err := s.sistema.TipoPlanIniciativa()
	if err != nil {
		return nil, err
	}
	if plan.TipoPlanID != tipo.ID || plan.Formato || !plan.Activo || strings.TrimSpace(plan.FormatoID) == "" {
		return nil, iniciativaNoExiste()
	}
	definicion, err := s.definiciones.Definicion(plan.FormatoID)
	if err != nil {
		return nil, err
	}
	estado, err := s.sistema.EstadoIniciativaPorID(plan.EstadoPlanID)
	if err != nil {
		return nil, err
	}
	if errLista != nil {
		return nil, errLista
	}
	submodulos := make(map[string]SubmoduloIniciativa, len(lista))
	for _, submodulo := range lista {
		submodulos[submodulo.Submodulo] = submodulo
	}
	return &contextoIniciativa{plan: plan, definicion: definicion, estado: estado, submodulos: submodulos}, nil
}

func iniciativaNoExiste() *ErrorNegocio {
	return errorNoEncontrado("La iniciativa no existe o fue eliminada.")
}

// documento devuelve los datos guardados del submódulo. Si faltan, la iniciativa
// quedó incompleta (p. ej. la plantilla ganó un submódulo después de crearla).
func (c *contextoIniciativa) documento(nodo *Nodo) (*SubmoduloIniciativa, error) {
	documento, existe := c.submodulos[nodo.Clave()]
	if !existe {
		return nil, nuevoError(http.StatusInternalServerError,
			"La iniciativa está incompleta: le falta el submódulo «%s». Comuníquese con el administrador del sistema.", nodo.Nombre).
			conDetalle("la iniciativa %s no tiene documento en submodulo_iniciativa para %s", c.plan.ID, nodo.Clave())
	}
	return &documento, nil
}

func (c *contextoIniciativa) submodulo(clave string) (*Nodo, *SubmoduloIniciativa, error) {
	nodo := c.definicion.Submodulo(clave)
	if nodo == nil {
		return nil, nil, errorNoEncontrado("El submódulo solicitado no existe en esta iniciativa.")
	}
	documento, err := c.documento(nodo)
	return nodo, documento, err
}

// datosIniciativa devuelve una copia del mapa submódulo -> datos.
func (c *contextoIniciativa) datosIniciativa() map[string]Datos {
	datos := make(map[string]Datos, len(c.submodulos))
	for clave, submodulo := range c.submodulos {
		datos[clave] = submodulo.Datos
	}
	return datos
}

func (c *contextoIniciativa) codigosEstadoSubmodulos(sistema *catalogosSistema) ([]string, error) {
	submodulos := c.definicion.Submodulos()
	codigos := make([]string, 0, len(submodulos))
	for _, nodo := range submodulos {
		documento, err := c.documento(nodo)
		if err != nil {
			return nil, err
		}
		estado, err := sistema.EstadoSubmoduloPorID(documento.EstadoSubmoduloID)
		if err != nil {
			return nil, err
		}
		codigos = append(codigos, estado.Codigo)
	}
	return codigos, nil
}
