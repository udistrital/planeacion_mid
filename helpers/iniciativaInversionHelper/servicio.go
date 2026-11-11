package iniciativaInversionHelper

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

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
	sistema      *catalogosSistema
	normalizador *Normalizador
	validador    *Validador
}

// NuevoServicio arma el servicio con sus dependencias.
func NuevoServicio(repo Repositorio, definiciones FuenteDefiniciones, catalogos ProveedorCatalogos) *Servicio {
	return &Servicio{
		repo:         repo,
		definiciones: definiciones,
		catalogos:    catalogos,
		sistema:      nuevosCatalogosSistema(repo),
		normalizador: NuevoNormalizador(catalogos),
		validador:    NuevoValidador(catalogos),
	}
}

// NuevoServicioPorDefecto conecta con planes_crud (host sin esquema, como en
// conf/app.conf), con caché de definiciones y de catálogos (colección elemento_catalogo).
func NuevoServicioPorDefecto(hostCrud string) *Servicio {
	crud := NuevoClienteCrud(hostCrud)
	definiciones := NuevaCacheDefiniciones(NuevoCargadorDefiniciones(crud, crud, 8), 10*time.Minute)
	return NuevoServicio(crud, definiciones, NuevosCatalogosConCache(crud, 10*time.Minute))
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

// ListarIniciativas lista las iniciativas activas, de la más reciente a la más
// antigua, con su número y su estado general. Si se indican dependencias, solo
// las de esas dependencias.
func (s *Servicio) ListarIniciativas(dependencias []string) ([]IniciativaListado, error) {
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
	deInteres := map[string]bool{}
	for _, dependencia := range dependencias {
		deInteres[dependencia] = true
	}

	iniciativas := make([]IniciativaListado, 0, len(planes))
	for _, plan := range planes {
		if !plan.Activo || (len(deInteres) > 0 && !deInteres[plan.DependenciaID]) {
			continue
		}
		estado, err := s.sistema.EstadoIniciativaPorID(plan.EstadoPlanID)
		if err != nil {
			return nil, err
		}
		iniciativas = append(iniciativas, IniciativaListado{
			ID:            plan.ID,
			Numero:        numeros[plan.ID],
			Nombre:        plan.Nombre,
			DependenciaID: plan.DependenciaID,
			PlantillaID:   plan.FormatoID,
			Estado:        estado,
			PermiteEditar: PermiteEditar(estado.Codigo),
			FechaCreacion: plan.FechaCreacion,
		})
	}
	sort.SliceStable(iniciativas, func(i, j int) bool { return iniciativas[i].Numero > iniciativas[j].Numero })
	return iniciativas, nil
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
// falla la creación de los submódulos, desactiva el plan.
func (s *Servicio) CrearIniciativa(solicitud SolicitudCrearIniciativa) (*ResumenIniciativa, error) {
	if strings.TrimSpace(solicitud.DependenciaID) == "" {
		return nil, errorSolicitud("Seleccione la unidad para la que se crea la iniciativa.")
	}
	vigencias, err := textoVigencias(solicitud.Vigencias)
	if err != nil {
		return nil, err
	}
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
	enFormulacion, err := s.sistema.EstadoIniciativa(EstadoIniciativaEnFormulacion)
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
	creado, err := s.repo.CrearPlan(Plan{
		Nombre:         nombre,
		Descripcion:    "Iniciativa de inversión formulada con la plantilla «" + definicion.Nombre + "»",
		TipoPlanID:     tipo.ID,
		AplicativoID:   aplicativoPlaneacion,
		Activo:         true,
		Formato:        false,
		FormatoID:      definicion.PlantillaID,
		DependenciaID:  solicitud.DependenciaID,
		EstadoPlanID:   enFormulacion.ID,
		VigenciaAplica: vigencias,
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
	return s.Resumen(creado.ID)
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

// Resumen arma la vista general: estado de la iniciativa y de cada submódulo, por módulo.
func (s *Servicio) Resumen(iniciativaID string) (*ResumenIniciativa, error) {
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	modulos := make([]ResumenModulo, 0, len(ctx.definicion.Modulos))
	codigos := []string{}
	orden := 0
	for _, modulo := range ctx.definicion.Modulos {
		resumenModulo := ResumenModulo{Clave: modulo.Clave(), Nombre: modulo.Nombre, Descripcion: modulo.Descripcion, Submodulos: []ResumenSubmodulo{}}
		for _, nodo := range modulo.Hijos {
			documento, err := ctx.documento(nodo)
			if err != nil {
				return nil, err
			}
			estado, err := s.sistema.EstadoSubmoduloPorID(documento.EstadoSubmoduloID)
			if err != nil {
				return nil, err
			}
			codigos = append(codigos, estado.Codigo)
			resumenModulo.Submodulos = append(resumenModulo.Submodulos, ResumenSubmodulo{
				Clave:         nodo.Clave(),
				Nombre:        nodo.Nombre,
				Descripcion:   nodo.Descripcion,
				Orden:         orden,
				Estado:        estado,
				Version:       documento.Version,
				Observaciones: len(documento.Observaciones),
			})
			orden++
		}
		modulos = append(modulos, resumenModulo)
	}
	return &ResumenIniciativa{
		ID:               ctx.plan.ID,
		Nombre:           ctx.plan.Nombre,
		DependenciaID:    ctx.plan.DependenciaID,
		PlantillaID:      ctx.plan.FormatoID,
		VigenciaAplica:   ctx.plan.VigenciaAplica,
		Estado:           ctx.estado,
		PermiteEditar:    PermiteEditar(ctx.estado.Codigo),
		PermiteRevisar:   PermiteRevisar(ctx.estado.Codigo),
		AccionesPosibles: AccionesPosibles(ctx.estado.Codigo, codigos),
		Modulos:          modulos,
	}, nil
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

// ObtenerSubmodulo devuelve el submódulo, sus errores vigentes y los datos de las
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
	return &SubmoduloConRelacionadas{Submodulo: *vista, Relacionadas: relacionadas}, nil
}

// GuardarSubmodulo guarda los datos aunque tengan errores (borrador). El submódulo
// queda Formulada solo si no tiene errores. Si otra persona la cambió desde que
// se consultó, responde 409.
func (s *Servicio) GuardarSubmodulo(iniciativaID, clave string, solicitud SolicitudGuardarSubmodulo) (*SubmoduloVista, error) {
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	if !PermiteEditar(ctx.estado.Codigo) {
		return nil, errorConflicto("La iniciativa está «%s»; no se puede editar.", ctx.estado.Nombre)
	}
	nodo, documento, err := ctx.submodulo(clave)
	if err != nil {
		return nil, err
	}

	datos := s.normalizador.Normalizar(nodo, solicitud.Datos)
	iniciativa := ctx.datosIniciativa()
	iniciativa[clave] = datos
	errores := s.validador.ValidarSubmodulo(nodo, iniciativa)
	estado, err := s.sistema.EstadoSubmodulo(EstadoTrasGuardar(errores))
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
	s.actualizarNombrePlan(ctx, nodo, datos)
	return s.vista(*guardada, errores)
}

// RevisarSubmodulo da el submódulo por revisado sin observaciones. Las observaciones
// se registran con AgregarObservacion, y su historial no se borra.
func (s *Servicio) RevisarSubmodulo(iniciativaID, clave string, solicitud SolicitudRevisarSubmodulo) (*SubmoduloVista, error) {
	if tieneTexto(solicitud.Observaciones) {
		return nil, errorFormatoSolicitud("la revisión no recibe observaciones; se registran con POST /iniciativas/%s/submodulos/%s/observaciones", iniciativaID, clave)
	}
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	nodo, documento, err := ctx.submodulo(clave)
	if err != nil {
		return nil, err
	}
	if err := s.exigirRevisable(ctx, nodo, documento); err != nil {
		return nil, err
	}

	sinObservaciones, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloSinObservaciones)
	if err != nil {
		return nil, err
	}
	guardada, err := s.repo.ActualizarSubmodulo(documento.ID, CambiosSubmodulo{
		Version:           solicitud.Version,
		EstadoSubmoduloID: sinObservaciones.ID,
	})
	if err != nil {
		return nil, err
	}
	return s.vista(*guardada, s.validador.ValidarSubmodulo(nodo, ctx.datosIniciativa()))
}

// EjecutarAccion cambia el estado general de la iniciativa (enviar, devolver, avalar, no-avalar).
func (s *Servicio) EjecutarAccion(iniciativaID, accion string) (*ResumenIniciativa, error) {
	if _, existe := transiciones[accion]; !existe {
		return nil, errorFormatoSolicitud("la acción %q no existe; las válidas son %s", accion, strings.Join(ordenAcciones, ", "))
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
	if accion == AccionEnviar {
		if err := s.exigirSubmodulosValidos(ctx); err != nil {
			return nil, err
		}
	}
	estado, err := s.sistema.EstadoIniciativa(destino)
	if err != nil {
		return nil, err
	}
	plan := *ctx.plan
	plan.EstadoPlanID = estado.ID
	if err := s.repo.ActualizarPlan(plan); err != nil {
		return nil, err
	}
	return s.Resumen(iniciativaID)
}

// exigirSubmodulosValidos vuelve a validar todos los submódulos antes de enviar la
// iniciativa, porque el estado guardado puede no estar al día: si al guardar un
// submódulo falló devolver a Sin formular uno que depende de él (revalidarDependientes),
// ese sigue Formulado con datos que ya no son válidos.
func (s *Servicio) exigirSubmodulosValidos(ctx *contextoIniciativa) error {
	iniciativa := ctx.datosIniciativa()
	conPendientes := []string{}
	for _, nodo := range ctx.definicion.Submodulos() {
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

// revalidarDependientes devuelve a "Sin formular" los submódulos que dependen de
// la que cambió y dejaron de ser válidas (p. ej. Objetivos si se agregó una causa).
func (s *Servicio) revalidarDependientes(ctx *contextoIniciativa, cambiada string, iniciativa map[string]Datos) {
	sinFormular, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloSinFormular)
	if err != nil {
		return
	}
	for _, nodo := range ctx.definicion.Submodulos() {
		clave := nodo.Clave()
		documento, existe := ctx.submodulos[clave]
		if clave == cambiada || !existe || documento.EstadoSubmoduloID == sinFormular.ID || !dependenciasDe(nodo)[cambiada] {
			continue
		}
		if len(s.validador.ValidarSubmodulo(nodo, iniciativa)) == 0 {
			continue
		}
		cambios := CambiosSubmodulo{Version: documento.Version, EstadoSubmoduloID: sinFormular.ID}
		if _, err := s.repo.ActualizarSubmodulo(documento.ID, cambios); err != nil {
			logs.Warn("iniciativaInversion: no se pudo marcar el submódulo %s como sin formular: %v", clave, err)
		}
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
		ID:            documento.ID,
		Clave:         documento.Submodulo,
		Estado:        estado,
		Datos:         datos,
		Observaciones: observaciones,
		Version:       documento.Version,
		Errores:       errores,
	}, nil
}

// dependenciasDe devuelve los submódulos de los que lee datos un submódulo
// (por espejo_de u opciones_desde absolutos).
func dependenciasDe(submodulo *Nodo) map[string]bool {
	dependencias := map[string]bool{}
	var recorrer func(nodos []*Nodo)
	recorrer = func(nodos []*Nodo) {
		for _, nodo := range nodos {
			for _, ruta := range []string{nodo.Dato.EspejoDe, nodo.Dato.OpcionesDesde} {
				if ruta == "" || strings.HasPrefix(ruta, prefijoOrigen) {
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
