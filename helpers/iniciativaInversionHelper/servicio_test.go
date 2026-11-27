package iniciativaInversionHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// repoMemoria imita a planes_crud en memoria, incluido el bloqueo optimista de
// submodulo_iniciativa (responde 409 si la versión no es la vigente). Empieza con las
// fechas de formulación habilitadas hoy para las dependencias de las pruebas (periodos).
type repoMemoria struct {
	planes            map[string]*Plan
	submodulos        map[string]*SubmoduloIniciativa
	catalogos         map[string][]Estado
	periodos          []map[string]interface{}
	consultasPeriodos int
	contador          int
	fallarSubmodulos  bool
	fallarListado     bool
	fallarPeriodos    bool
}

// tipoSeguimientoDePrueba es el id del tipo de seguimiento FIPI_SP en repoMemoria.
const tipoSeguimientoDePrueba = "ts-fipi"

func nuevoRepoMemoria() *repoMemoria {
	estado := func(id, codigo, nombre string) Estado { return Estado{ID: id, Codigo: codigo, Nombre: nombre} }
	return &repoMemoria{
		planes:     map[string]*Plan{},
		submodulos: map[string]*SubmoduloIniciativa{},
		periodos:   []map[string]interface{}{periodoDeFormulacion(-30, 30, "dep-1", "dep-2", "41", "42")},
		catalogos: map[string][]Estado{
			RecursoTipoSeguimiento: {estado(tipoSeguimientoDePrueba, CodigoTipoSeguimientoFormulacion, "Formulación de Iniciativas de Inversión")},
			RecursoTipoPlan:        {estado(tipoPlanDePrueba, CodigoTipoPlanIniciativa, "Iniciativa de inversión")},
			RecursoEstadoPlan: {
				estado("ep-sf", EstadoIniciativaSinFormular, "Sin formular"),
				estado("ep-ef", EstadoIniciativaEnFormulacion, "En formulación"),
				estado("ep-pa", EstadoIniciativaPorAprobar, "Por aprobar (Jefe de unidad)"),
				estado("ep-pr", EstadoIniciativaPorRevisar, "Por revisar OAP"),
				estado("ep-er", EstadoIniciativaEnRevision, "En revisión OAP"),
				estado("ep-pc", EstadoIniciativaPorCorregir, "Por corregir"),
				estado("ep-apr", EstadoIniciativaAprobada, "Aprobado"),
				estado("ep-rch", EstadoIniciativaRechazada, "Rechazada"),
			},
			RecursoEstadoSubmodulo: {
				estado("es-sf", EstadoSubmoduloSinFormular, "Sin formular"),
				estado("es-fn", EstadoSubmoduloFormulando, "Formulando"),
				estado("es-for", EstadoSubmoduloFormulado, "Formulado"),
				estado("es-co", EstadoSubmoduloConObservaciones, "Con observaciones"),
				estado("es-apr", EstadoSubmoduloAprobado, "Aprobado"),
			},
		},
	}
}

func (r *repoMemoria) nuevoID(prefijo string) string {
	r.contador++
	return fmt.Sprintf("%s-%d", prefijo, r.contador)
}

func (r *repoMemoria) ObtenerPlan(id string) (*Plan, error) {
	plan, existe := r.planes[id]
	if !existe {
		return nil, errorNoEncontrado("No existe el plan %s.", id)
	}
	copia := *plan
	return &copia, nil
}

func (r *repoMemoria) BuscarPlanes(filtros map[string]string) ([]Plan, error) {
	var planes []Plan
	for _, plan := range r.planes {
		if plan.TipoPlanID != filtros["tipo_plan_id"] || fmt.Sprint(plan.Formato) != filtros["formato"] {
			continue
		}
		if activo, filtra := filtros["activo"]; filtra && fmt.Sprint(plan.Activo) != activo {
			continue
		}
		planes = append(planes, *plan)
	}
	return planes, nil
}

// CrearPlan pone una fecha de creación creciente, como el crud.
func (r *repoMemoria) CrearPlan(plan Plan) (*Plan, error) {
	plan.ID = r.nuevoID("plan")
	plan.FechaCreacion = fmt.Sprintf("2026-11-10T08:%02d:00.000Z", r.contador)
	r.planes[plan.ID] = &plan
	copia := plan
	return &copia, nil
}

func (r *repoMemoria) ActualizarPlan(plan Plan) error {
	r.planes[plan.ID] = &plan
	return nil
}

func (r *repoMemoria) DesactivarPlan(id string) error {
	r.planes[id].Activo = false
	return nil
}

func (r *repoMemoria) HijosDe(string) ([]Subgrupo, error)           { return nil, nil }
func (r *repoMemoria) DetallesDe(string) ([]DetalleSubgrupo, error) { return nil, nil }

func (r *repoMemoria) BuscarPorCodigo(recurso, codigo string) (*Estado, error) {
	for _, registro := range r.catalogos[recurso] {
		if registro.Codigo == codigo {
			copia := registro
			return &copia, nil
		}
	}
	return nil, errorNoEncontrado("No existe %s %s.", recurso, codigo)
}

func (r *repoMemoria) ListarActivos(recurso string) ([]Estado, error) {
	return r.catalogos[recurso], nil
}

func (r *repoMemoria) ListarSubmodulos(planID string) ([]SubmoduloIniciativa, error) {
	if r.fallarListado {
		return nil, errorServicio("falla simulada")
	}
	var submodulos []SubmoduloIniciativa
	for _, submodulo := range r.submodulos {
		if submodulo.PlanID == planID {
			submodulos = append(submodulos, *submodulo)
		}
	}
	return submodulos, nil
}

func (r *repoMemoria) CrearSubmodulos(nuevos []NuevoSubmodulo) ([]SubmoduloIniciativa, error) {
	if r.fallarSubmodulos {
		return nil, errorServicio("falla simulada")
	}
	var creados []SubmoduloIniciativa
	for _, nuevo := range nuevos {
		submodulo := &SubmoduloIniciativa{ID: r.nuevoID("sec"), PlanID: nuevo.PlanID, Submodulo: nuevo.Submodulo, EstadoSubmoduloID: nuevo.EstadoSubmoduloID, Datos: nuevo.Datos, Activo: true}
		r.submodulos[submodulo.ID] = submodulo
		creados = append(creados, *submodulo)
	}
	return creados, nil
}

func (r *repoMemoria) ActualizarSubmodulo(id string, cambios CambiosSubmodulo) (*SubmoduloIniciativa, error) {
	submodulo, existe := r.submodulos[id]
	if !existe {
		return nil, errorNoEncontrado("No existe el submódulo %s.", id)
	}
	if submodulo.Version != cambios.Version {
		return nil, errorConflicto("El submódulo cambió desde que se consultó.")
	}
	if cambios.Datos != nil {
		submodulo.Datos = *cambios.Datos
	}
	if cambios.EstadoSubmoduloID != "" {
		submodulo.EstadoSubmoduloID = cambios.EstadoSubmoduloID
	}
	submodulo.Version++
	copia := *submodulo
	return &copia, nil
}

// AgregarObservacion imita el $push del crud: agrega al historial sin tocar la versión.
func (r *repoMemoria) AgregarObservacion(id string, nueva NuevaObservacion) (*SubmoduloIniciativa, error) {
	submodulo, existe := r.submodulos[id]
	if !existe {
		return nil, errorNoEncontrado("No existe el submódulo %s.", id)
	}
	submodulo.Observaciones = append(submodulo.Observaciones, Observacion{
		ID: r.nuevoID("obs"), Texto: nueva.Texto, Rol: nueva.Rol, Usuario: nueva.Usuario, UsuarioNombre: nueva.UsuarioNombre,
	})
	if nueva.EstadoSubmoduloID != "" {
		submodulo.EstadoSubmoduloID = nueva.EstadoSubmoduloID
	}
	copia := *submodulo
	copia.Observaciones = append([]Observacion(nil), submodulo.Observaciones...)
	return &copia, nil
}

// PeriodosActivos devuelve los periodos solo para el tipo de seguimiento FIPI_SP, y cuenta las consultas.
func (r *repoMemoria) PeriodosActivos(tipoSeguimientoID string) ([]map[string]interface{}, error) {
	r.consultasPeriodos++
	if r.fallarPeriodos {
		return nil, errorServicio("falla simulada")
	}
	if tipoSeguimientoID != tipoSeguimientoDePrueba {
		return nil, errorServicio("se pidieron los periodos del tipo de seguimiento %q", tipoSeguimientoID)
	}
	return r.periodos, nil
}

// periodoDeFormulacion es un registro de fechas de formulación, como lo guarda el crud, que va de
// desde a hasta días contados a partir de hoy, con las dependencias indicadas.
func periodoDeFormulacion(desde, hasta int, dependencias ...string) map[string]interface{} {
	unidades := []map[string]string{}
	for _, dependencia := range dependencias {
		unidades = append(unidades, map[string]string{"Id": dependencia, "Nombre": "Unidad " + dependencia})
	}
	texto, err := json.Marshal(unidades)
	if err != nil {
		panic(err)
	}
	dia := func(dias int) string { return time.Now().AddDate(0, 0, dias).Format("2006-01-02") + "T00:00:00.000Z" }
	return map[string]interface{}{"_id": fmt.Sprintf("periodo%+d%+d", desde, hasta), "fecha_inicio": dia(desde), "fecha_fin": dia(hasta), "unidades_interes": string(texto)}
}

type definicionFija struct{ definicion *Definicion }

func (d definicionFija) Definicion(string) (*Definicion, error) { return d.definicion, nil }

// ---------------------------------------------------------------------------
// Usuarios de las pruebas
// ---------------------------------------------------------------------------

// Documentos de los usuarios de las pruebas.
const (
	responsableDePrueba = "52123456"   // formulador (ASISTENTE_DEPENDENCIA) de dep-1, dep-2, 41 y 42
	documentoJefe       = "79876543"   // JEFE_DEPENDENCIA de dep-1
	documentoOAP        = "1012345678" // PLANEACION, de la dependencia 11
	documentoJefeYOAP   = "80111222"   // JEFE_DEPENDENCIA y ASISTENTE_PLANEACION de dep-1
)

// personasDePrueba son los terceros que conocen las pruebas: documento -> nombre completo.
var personasDePrueba = map[string]string{
	responsableDePrueba: "LAURA GOMEZ RIOS",
	documentoJefe:       "MARTA LOPEZ DIAZ",
	documentoOAP:        "CARLOS RUIZ PENA",
	documentoJefeYOAP:   "ANDRES MORA VEGA",
}

// unidadesDePrueba son las unidades de cada usuario de las pruebas.
var unidadesDePrueba = map[string][]string{
	responsableDePrueba: {"dep-1", "dep-2", "41", "42"},
	documentoJefe:       {"dep-1"},
	documentoOAP:        {"11"},
	documentoJefeYOAP:   {"dep-1"},
}

func formulador() Actor {
	return Actor{Usuario: responsableDePrueba, Roles: []string{RolAsistenteDependencia}}
}

func jefe() Actor { return Actor{Usuario: documentoJefe, Roles: []string{RolJefeDependencia}} }

func oap() Actor { return Actor{Usuario: documentoOAP, Roles: []string{RolPlaneacion}} }

func jefeYOAP() Actor {
	return Actor{Usuario: documentoJefeYOAP, Roles: []string{RolJefeDependencia, RolAsistentePlaneacion}}
}

// tercerosFalsos imita a Terceros (nombre por documento) y cuenta las consultas. Un documento
// que no está en nombres responde con error, como si no estuviera o Terceros fallara.
type tercerosFalsos struct {
	nombres   map[string]string
	consultas int
}

func (f *tercerosFalsos) NombreTercero(documento string) (string, error) {
	f.consultas++
	nombre, existe := f.nombres[documento]
	if !existe {
		return "", errorServicio("Terceros no respondió")
	}
	return nombre, nil
}

// tercerosCaidos deja al servicio con un Terceros que no responde.
func tercerosCaidos(servicio *Servicio) *tercerosFalsos {
	caidos := &tercerosFalsos{}
	servicio.terceros = caidos
	return caidos
}

// unidadesFalsas imita las unidades de los usuarios (Terceros) y cuenta las consultas.
type unidadesFalsas struct {
	porUsuario map[string][]string
	consultas  int
	err        error
}

func (u *unidadesFalsas) DependenciasDe(documento string) ([]string, error) {
	u.consultas++
	return u.porUsuario[documento], u.err
}

// conPersonasDePrueba conecta el servicio con los terceros y las unidades de las pruebas.
func conPersonasDePrueba(servicio *Servicio) *Servicio {
	servicio.terceros = &tercerosFalsos{nombres: personasDePrueba}
	servicio.unidades = &unidadesFalsas{porUsuario: unidadesDePrueba}
	return servicio
}

// dependenciasFalsas imita a Oikos (nombre de una dependencia por id) y cuenta las
// consultas de cada una. Un id que no está en nombres responde con error, como si
// Oikos fallara.
type dependenciasFalsas struct {
	nombres  map[int]string
	llamadas map[int]int
}

func nuevasDependenciasFalsas(nombres map[int]string) *dependenciasFalsas {
	return &dependenciasFalsas{nombres: nombres, llamadas: map[int]int{}}
}

func (o *dependenciasFalsas) NombreDependencia(id int) (string, error) {
	o.llamadas[id]++
	nombre, existe := o.nombres[id]
	if !existe {
		return "", errors.New("Oikos no respondió")
	}
	return nombre, nil
}

// ---------------------------------------------------------------------------
// Ayudantes
// ---------------------------------------------------------------------------

func nuevoServicioDePrueba(t *testing.T) (*Servicio, *repoMemoria) {
	t.Helper()
	repo := nuevoRepoMemoria()
	servicio := conPersonasDePrueba(NuevoServicio(repo, definicionFija{definicionDePrueba(t)}, catalogosDePrueba(), nil))
	return servicio, repo
}

func crearEnDependencia(t *testing.T, servicio *Servicio, dependencia string) *ResumenIniciativa {
	t.Helper()
	resumen, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: formulador(), PlantillaID: "plantilla-1", DependenciaID: dependencia})
	if err != nil {
		t.Fatal(err)
	}
	return resumen
}

func crearIniciativaDePrueba(t *testing.T, servicio *Servicio) *ResumenIniciativa {
	t.Helper()
	return crearEnDependencia(t, servicio, "dep-1")
}

// guardar guarda el submódulo como formulador, con la versión vigente, y devuelve la vista.
func guardar(t *testing.T, servicio *Servicio, iniciativaID, clave string, datos Datos) *SubmoduloVista {
	t.Helper()
	vista, err := intentarGuardar(t, servicio, iniciativaID, clave, datos)
	if err != nil {
		t.Fatal(err)
	}
	return vista
}

func intentarGuardar(t *testing.T, servicio *Servicio, iniciativaID, clave string, datos Datos) (*SubmoduloVista, error) {
	t.Helper()
	actual, err := servicio.ObtenerSubmodulo(iniciativaID, clave)
	if err != nil {
		t.Fatal(err)
	}
	return servicio.GuardarSubmodulo(iniciativaID, clave, SolicitudGuardarSubmodulo{Actor: formulador(), Version: actual.Submodulo.Version, Datos: datos})
}

// formularTodo guarda, en el orden de la plantilla, todos los submódulos del ejemplo y exige
// que queden formulados.
func formularTodo(t *testing.T, servicio *Servicio, iniciativaID string) {
	t.Helper()
	iniciativa := iniciativaDeEjemplo(t)
	for _, submodulo := range definicionDePrueba(t).Submodulos() {
		clave := submodulo.Clave()
		if vista := guardar(t, servicio, iniciativaID, clave, iniciativa[clave]); vista.Estado.Codigo != EstadoSubmoduloFormulado {
			t.Fatalf("%s quedó %s: %+v", clave, vista.Estado.Codigo, vista.Errores)
		}
	}
}

// ejecutar ejecuta la acción como el actor y exige que pase.
func ejecutar(t *testing.T, servicio *Servicio, iniciativaID, accion string, actor Actor) *ResumenIniciativa {
	t.Helper()
	resumen, err := servicio.EjecutarAccion(iniciativaID, accion, SolicitudAccion{Actor: actor, Motivo: "Motivo de prueba."})
	if err != nil {
		t.Fatalf("%s: %v", accion, err)
	}
	return resumen
}

// llevarARevision formula toda la iniciativa y la lleva hasta En revisión OAP: el formulador la
// envía, el Jefe la aprueba y la OAP inicia la revisión.
func llevarARevision(t *testing.T, servicio *Servicio, iniciativaID string) *ResumenIniciativa {
	t.Helper()
	formularTodo(t, servicio, iniciativaID)
	ejecutar(t, servicio, iniciativaID, AccionEnviar, formulador())
	ejecutar(t, servicio, iniciativaID, AccionAprobar, jefe())
	return ejecutar(t, servicio, iniciativaID, AccionIniciarRevision, oap())
}

func resumenDe(t *testing.T, servicio *Servicio, iniciativaID string, actor Actor) *ResumenIniciativa {
	t.Helper()
	resumen, err := servicio.Resumen(iniciativaID, actor)
	if err != nil {
		t.Fatal(err)
	}
	return resumen
}

func estadoDeSubmodulo(t *testing.T, servicio *Servicio, iniciativaID, clave string) string {
	t.Helper()
	for _, submodulo := range submodulosDe(resumenDe(t, servicio, iniciativaID, oap())) {
		if submodulo.Clave == clave {
			return submodulo.Estado.Codigo
		}
	}
	t.Fatalf("no existe el submódulo %s", clave)
	return ""
}

func exigirStatus(t *testing.T, err error, status int) {
	t.Helper()
	var negocio *ErrorNegocio
	if !errors.As(err, &negocio) || negocio.Status != status {
		t.Fatalf("se esperaba un error %d y llegó %v", status, err)
	}
}

// accionesDe resume las acciones de los permisos como "accion" (habilitada) o "accion!" (no).
func accionesDe(permisos Permisos) []string {
	acciones := []string{}
	for _, accion := range permisos.Acciones {
		if accion.Habilitada {
			acciones = append(acciones, accion.Accion)
		} else {
			acciones = append(acciones, accion.Accion+"!")
		}
	}
	return acciones
}

// ---------------------------------------------------------------------------
// Crear
// ---------------------------------------------------------------------------

func TestCrearIniciativaCreaUnSubmoduloSinFormularPorSubmoduloDeLaPlantilla(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	resumen := crearIniciativaDePrueba(t, servicio)

	if resumen.Estado.Codigo != EstadoIniciativaSinFormular || !resumen.Permisos.Editar || resumen.DependenciaID != "dep-1" {
		t.Fatalf("resumen inicial = %+v", resumen)
	}
	if esperados := len(definicionDePrueba(t).Submodulos()); len(submodulosDe(resumen)) != esperados {
		t.Fatalf("se esperaban %d submódulos y hay %d", esperados, len(submodulosDe(resumen)))
	}
	for _, submodulo := range submodulosDe(resumen) {
		if submodulo.Estado.Codigo != EstadoSubmoduloSinFormular || submodulo.Accion != AccionSubmoduloFormular {
			t.Fatalf("%s quedó %s (%s)", submodulo.Clave, submodulo.Estado.Codigo, submodulo.Accion)
		}
	}
	if len(resumen.Permisos.Acciones) != 0 {
		t.Fatalf("Sin formular no tiene acciones: %v", accionesDe(resumen.Permisos))
	}
}

func TestCrearIniciativaDesactivaElPlanSiFallanLosSubmodulos(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	repo.fallarSubmodulos = true

	if _, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: formulador(), PlantillaID: "plantilla-1", DependenciaID: "dep-1"}); err == nil {
		t.Fatal("se esperaba un error")
	}
	for _, plan := range repo.planes {
		if plan.Activo {
			t.Fatal("el plan debió quedar desactivado")
		}
	}
}

func TestCrearIniciativaExigeDependencia(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: formulador(), PlantillaID: "plantilla-1"})
	exigirStatus(t, err, http.StatusBadRequest)
}

// Sin un documento válido o sin roles no se crea nada ni se consulta Terceros.
func TestCrearIniciativaExigeElUsuarioYSusRoles(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	terceros := servicio.terceros.(*tercerosFalsos)
	actores := []Actor{
		{Usuario: "", Roles: []string{RolAsistenteDependencia}},
		{Usuario: "   ", Roles: []string{RolAsistenteDependencia}},
		{Usuario: "123,Activo:false", Roles: []string{RolAsistenteDependencia}},
		{Usuario: "12 34", Roles: []string{RolAsistenteDependencia}},
		{Usuario: "123456789012345678901", Roles: []string{RolAsistenteDependencia}},
		{Usuario: responsableDePrueba},
		{Usuario: responsableDePrueba, Roles: []string{" "}},
	}
	for _, actor := range actores {
		_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: actor, PlantillaID: "plantilla-1", DependenciaID: "dep-1"})
		exigirStatus(t, err, http.StatusBadRequest)
		exigirMensajeParaUsuario(t, err, "roles")
	}
	if len(repo.planes) != 0 || terceros.consultas != 0 {
		t.Fatalf("planes = %d, consultas a Terceros = %d; se esperaba 0 y 0", len(repo.planes), terceros.consultas)
	}
}

// Crea el carril del formulador (también el Jefe, que formula) y solo para sus unidades; la OAP
// sin un rol de formulador no crea.
func TestSoloElFormuladorCreaYSoloEnSusUnidades(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	for nombre, caso := range map[string]struct {
		actor       Actor
		dependencia string
	}{
		"la OAP":             {oap(), "dep-1"},
		"en otra unidad":     {formulador(), "dep-9"},
		"el Jefe en otra":    {jefe(), "dep-2"},
		"un rol desconocido": {Actor{Usuario: responsableDePrueba, Roles: []string{"ADMINISTRADOR"}}, "dep-1"},
	} {
		_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: caso.actor, PlantillaID: "plantilla-1", DependenciaID: caso.dependencia})
		if negocio := ComoErrorNegocio(err); err == nil || negocio.Status != http.StatusForbidden {
			t.Errorf("%s: se esperaba 403 y llegó %v", nombre, err)
		}
	}
	if len(repo.planes) != 0 {
		t.Fatalf("no debió crearse ningún plan: %d", len(repo.planes))
	}
	creada, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: jefe(), PlantillaID: "plantilla-1", DependenciaID: "dep-1"})
	if err != nil || repo.planes[creada.ID].Responsable != documentoJefe {
		t.Fatalf("el Jefe formula en su unidad: %+v %v", creada, err)
	}
}

// Al crear se guardan el documento (sin espacios) y el nombre del responsable; el listado los
// muestra sin consultar Terceros.
func TestElNombreDelResponsableSeGuardaAlCrearYElListadoNoConsultaTerceros(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	terceros := servicio.terceros.(*tercerosFalsos)
	actor := formulador()
	actor.Usuario = " " + actor.Usuario + " "
	resumen, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: actor, PlantillaID: "plantilla-1", DependenciaID: "dep-1"})
	if err != nil {
		t.Fatal(err)
	}
	if plan := repo.planes[resumen.ID]; plan.Responsable != responsableDePrueba || plan.ResponsableNombre != personasDePrueba[responsableDePrueba] {
		t.Fatalf("responsable = %q, nombre = %q", plan.Responsable, plan.ResponsableNombre)
	}

	terceros.consultas = 0
	listado, err := servicio.ListarIniciativas(oap())
	if err != nil {
		t.Fatal(err)
	}
	if len(listado) != 1 || listado[0].Responsable != personasDePrueba[responsableDePrueba] || terceros.consultas != 0 {
		t.Fatalf("listado = %+v, consultas a Terceros = %d", listado, terceros.consultas)
	}
}

// Si Terceros no da el nombre, la iniciativa se crea igual y el listado la muestra sin él.
func TestSiTercerosNoRespondeLaIniciativaSeCreaSinNombre(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	tercerosCaidos(servicio)
	resumen := crearIniciativaDePrueba(t, servicio)

	if plan := repo.planes[resumen.ID]; plan.Responsable != responsableDePrueba || plan.ResponsableNombre != "" {
		t.Fatalf("responsable = %q, nombre = %q", plan.Responsable, plan.ResponsableNombre)
	}
	listado, err := servicio.ListarIniciativas(oap())
	if err != nil {
		t.Fatal(err)
	}
	if len(listado) != 1 || listado[0].Responsable != "" {
		t.Fatalf("listado = %+v", listado)
	}
}

func TestSinPlantillaSeCreaConLaMasReciente(t *testing.T) {
	repo := nuevoRepoMemoria()
	repo.planes["anterior"] = &Plan{ID: "anterior", TipoPlanID: tipoPlanDePrueba, Formato: true, Activo: true, FechaCreacion: "2026-10-01T00:00:00.000Z"}
	repo.planes["reciente"] = &Plan{ID: "reciente", TipoPlanID: tipoPlanDePrueba, Formato: true, Activo: true, FechaCreacion: "2026-11-01T00:00:00.000Z"}
	anterior, reciente := definicionDePrueba(t), definicionDePrueba(t)
	anterior.PlantillaID, reciente.PlantillaID = "anterior", "reciente"
	servicio := conPersonasDePrueba(NuevoServicio(repo, definicionesPorPlantilla{"anterior": anterior, "reciente": reciente}, catalogosDePrueba(), nil))

	resumen, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: formulador(), DependenciaID: "dep-1"})
	if err != nil {
		t.Fatal(err)
	}
	if plantilla := repo.planes[resumen.ID].FormatoID; plantilla != "reciente" || resumen.Nombre != nombreIniciativaInicial {
		t.Fatalf("plantilla = %s, nombre = %s", plantilla, resumen.Nombre)
	}
}

func TestSinPlantillasNoSeCreaLaIniciativa(t *testing.T) {
	servicio := conPersonasDePrueba(NuevoServicio(nuevoRepoMemoria(), definicionesPorPlantilla{}, catalogosDePrueba(), nil))
	_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: formulador(), DependenciaID: "dep-1"})
	exigirStatus(t, err, http.StatusConflict)
}

// ---------------------------------------------------------------------------
// Listado
// ---------------------------------------------------------------------------

// La OAP ve todas las iniciativas activas; los demás, las de sus unidades.
func TestListarNumeraYMuestraACadaUnoLasQueConsulta(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	desactivada := crearEnDependencia(t, servicio, "dep-1").ID
	activa := crearEnDependencia(t, servicio, "dep-1").ID
	otraUnidad := crearEnDependencia(t, servicio, "dep-2").ID
	repo.planes[desactivada].Activo = false

	todas, err := servicio.ListarIniciativas(oap())
	if err != nil {
		t.Fatal(err)
	}
	if len(todas) != 2 || todas[0].ID != otraUnidad || todas[0].Numero != 3 || todas[1].ID != activa || todas[1].Numero != 2 {
		t.Fatalf("listado = %+v", todas)
	}
	if todas[1].Estado.Codigo != EstadoIniciativaSinFormular || todas[1].Accion != AccionSubmoduloConsultar || todas[1].Responsable != personasDePrueba[responsableDePrueba] {
		t.Fatalf("la OAP solo consulta una iniciativa Sin formular: %+v", todas[1])
	}

	deSuUnidad, err := servicio.ListarIniciativas(jefe())
	if err != nil {
		t.Fatal(err)
	}
	if len(deSuUnidad) != 1 || deSuUnidad[0].ID != activa || deSuUnidad[0].Accion != AccionSubmoduloFormular {
		t.Fatalf("el Jefe ve solo las de dep-1, y las formula: %+v", deSuUnidad)
	}
}

// Oikos responde: la fila trae el nombre. Cada dependencia se consulta una sola vez, y si
// Oikos falla con una, esa fila sale sin nombre y las demás no se afectan.
func TestListarTraeElNombreDeLaDependenciaYLaConsultaUnaSolaVez(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	oikos := nuevasDependenciasFalsas(map[int]string{41: "OFICINA ASESORA DE PLANEACION"})
	servicio.dependencias = oikos
	crearEnDependencia(t, servicio, "41")
	crearEnDependencia(t, servicio, "41")
	crearEnDependencia(t, servicio, "42") // Oikos no la conoce

	listado, err := servicio.ListarIniciativas(oap())
	if err != nil {
		t.Fatal(err)
	}
	if len(listado) != 3 {
		t.Fatalf("listado = %+v", listado)
	}
	for _, fila := range listado {
		esperado := ""
		if fila.DependenciaID == "41" {
			esperado = "OFICINA ASESORA DE PLANEACION"
		}
		if fila.DependenciaNombre != esperado {
			t.Fatalf("la dependencia %s quedó con nombre %q, se esperaba %q", fila.DependenciaID, fila.DependenciaNombre, esperado)
		}
	}
	if oikos.llamadas[41] != 1 || oikos.llamadas[42] != 1 {
		t.Fatalf("consultas a Oikos = %v, se esperaba una por dependencia (también la que falló)", oikos.llamadas)
	}
}

// Si la dependencia_id no es un número no se consulta Oikos (NombreDependencia recibe un int)
// y el listado sale igual, sin nombre.
func TestListarNoConsultaOikosConUnaDependenciaQueNoEsNumero(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	oikos := nuevasDependenciasFalsas(map[int]string{})
	servicio.dependencias = oikos
	crearEnDependencia(t, servicio, "dep-1")

	listado, err := servicio.ListarIniciativas(oap())
	if err != nil {
		t.Fatal(err)
	}
	if len(listado) != 1 || listado[0].DependenciaNombre != "" || len(oikos.llamadas) != 0 {
		t.Fatalf("listado = %+v, consultas a Oikos = %v", listado, oikos.llamadas)
	}
}

// Sin fuente de dependencias (servicio.dependencias == nil) el listado sale sin nombre, sin pánico.
func TestListarSaleSinNombreDeDependenciaSiNoHayFuente(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	crearEnDependencia(t, servicio, "41")

	listado, err := servicio.ListarIniciativas(oap())
	if err != nil {
		t.Fatal(err)
	}
	if len(listado) != 1 || listado[0].DependenciaNombre != "" {
		t.Fatalf("listado = %+v", listado)
	}
}

// Las unidades del usuario se consultan una sola vez por listado, y la OAP sin otro rol no las
// necesita.
func TestListarConsultaLasUnidadesUnaSolaVez(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	unidades := servicio.unidades.(*unidadesFalsas)
	crearEnDependencia(t, servicio, "dep-1")
	crearEnDependencia(t, servicio, "dep-2")

	unidades.consultas = 0
	if _, err := servicio.ListarIniciativas(formulador()); err != nil || unidades.consultas != 1 {
		t.Fatalf("formulador: err = %v, consultas = %d", err, unidades.consultas)
	}
	unidades.consultas = 0
	if _, err := servicio.ListarIniciativas(oap()); err != nil || unidades.consultas != 0 {
		t.Fatalf("OAP: err = %v, consultas = %d", err, unidades.consultas)
	}
}

// ---------------------------------------------------------------------------
// Guardar
// ---------------------------------------------------------------------------

func TestGuardarNombreLoFormulaYActualizaElNombreDelPlan(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)

	vista := guardar(t, servicio, iniciativa.ID, "nombre", iniciativaDeEjemplo(t)["nombre"])

	if vista.Estado.Codigo != EstadoSubmoduloFormulado || len(vista.Errores) != 0 {
		t.Fatalf("estado=%s errores=%+v", vista.Estado.Codigo, vista.Errores)
	}
	if nombre := repo.planes[iniciativa.ID].Nombre; nombre != vista.Datos["nombre_proyecto"] {
		t.Fatalf("el plan quedó con nombre %q", nombre)
	}
}

func TestGuardarConErroresGuardaElBorradorFormulando(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	participantes := iniciativaDeEjemplo(t)["participantes"]
	participantes["analisis"] = ""

	vista := guardar(t, servicio, iniciativa.ID, "participantes", participantes)

	if vista.Estado.Codigo != EstadoSubmoduloFormulando || len(vista.Errores) != 1 {
		t.Fatalf("estado=%s errores=%+v", vista.Estado.Codigo, vista.Errores)
	}
	if len(filasEn(t, vista.Datos, "actores")) != 4 {
		t.Fatal("el borrador debió guardarse aunque tenga errores")
	}
	if estado := repo.planes[iniciativa.ID].EstadoPlanID; estado != "ep-ef" {
		t.Fatalf("con la primera información la iniciativa debió pasar a En formulación y está %s", estado)
	}
}

func TestGuardarVacioDejaSinFormularElSubmoduloYLaIniciativa(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)

	vista := guardar(t, servicio, iniciativa.ID, "objetivos", Datos{})

	if vista.Estado.Codigo != EstadoSubmoduloSinFormular || len(vista.Errores) == 0 {
		t.Fatalf("estado=%s errores=%d", vista.Estado.Codigo, len(vista.Errores))
	}
	if estado := repo.planes[iniciativa.ID].EstadoPlanID; estado != "ep-sf" {
		t.Fatalf("sin información la iniciativa sigue Sin formular y está %s", estado)
	}
}

func TestPasarAEnFormulacionConservaElNombreQueSePropone(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)

	vista := guardar(t, servicio, iniciativa.ID, "nombre", iniciativaDeEjemplo(t)["nombre"])

	plan := repo.planes[iniciativa.ID]
	if plan.EstadoPlanID != "ep-ef" || plan.Nombre != vista.Datos["nombre_proyecto"] {
		t.Fatalf("el plan debió quedar En formulación y con el nombre propuesto: estado=%s nombre=%q", plan.EstadoPlanID, plan.Nombre)
	}
}

func TestGuardarConUnaVersionViejaResponde409(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	guardar(t, servicio, iniciativa.ID, "nombre", iniciativaDeEjemplo(t)["nombre"]) // versión pasa de 0 a 1

	_, err := servicio.GuardarSubmodulo(iniciativa.ID, "nombre", SolicitudGuardarSubmodulo{Actor: formulador(), Version: 0, Datos: Datos{}})
	exigirStatus(t, err, http.StatusConflict)
}

// Guarda solo el formulador de la unidad: ni la OAP ni el Jefe de otra unidad.
func TestSoloGuardaElFormuladorDeLaUnidad(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	deOtraUnidad := crearEnDependencia(t, servicio, "dep-2")
	for nombre, actor := range map[string]Actor{"la OAP": oap(), "el Jefe de otra unidad": jefe()} {
		_, err := servicio.GuardarSubmodulo(deOtraUnidad.ID, "nombre", SolicitudGuardarSubmodulo{Actor: actor, Datos: iniciativaDeEjemplo(t)["nombre"]})
		if negocio := ComoErrorNegocio(err); err == nil || negocio.Status != http.StatusForbidden {
			t.Errorf("%s: se esperaba 403 y llegó %v", nombre, err)
		}
	}
}

func TestAgregarUnaCausaDirectaDevuelveObjetivosAFormulando(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)

	guardar(t, servicio, iniciativa.ID, "problematica", problematicaConOtraCausa(t))

	if estado := estadoDeSubmodulo(t, servicio, iniciativa.ID, "objetivos"); estado != EstadoSubmoduloFormulando {
		t.Fatalf("objetivos debió volver a formulando (falta el objetivo de la nueva causa) y está %s", estado)
	}
	if estado := estadoDeSubmodulo(t, servicio, iniciativa.ID, "participantes"); estado != EstadoSubmoduloFormulado {
		t.Fatalf("participantes no depende de la problemática y quedó %s", estado)
	}
}

func TestCambiarUnTextoDeLaProblematicaNoAfectaLosObjetivos(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)

	problematica := iniciativaDeEjemplo(t)["problematica"]
	problematica["problema_central"] = "Baja capacidad tecnológica para la transformación digital de la Universidad."
	guardar(t, servicio, iniciativa.ID, "problematica", problematica)

	if estado := estadoDeSubmodulo(t, servicio, iniciativa.ID, "objetivos"); estado != EstadoSubmoduloFormulado {
		t.Fatalf("objetivos sigue con un objetivo por causa directa y quedó %s", estado)
	}
}

// Simula que al guardar Problemática falló devolver Objetivos a Formulando: la nueva
// causa quedó guardada y Objetivos sigue Formulado sin el objetivo de esa causa. Enviar
// debe detectarlo aunque el estado guardado diga Formulado.
func TestEnviarRevalidaLosSubmodulosAunqueSuEstadoDigaFormulado(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)

	cambiarDatos(repo, iniciativa.ID, "problematica", problematicaConOtraCausa(t))
	if estado := estadoDeSubmodulo(t, servicio, iniciativa.ID, "objetivos"); estado != EstadoSubmoduloFormulado {
		t.Fatalf("la prueba supone objetivos Formulado y está %s", estado)
	}

	_, err := servicio.EjecutarAccion(iniciativa.ID, AccionEnviar, SolicitudAccion{Actor: formulador()})
	exigirStatus(t, err, http.StatusConflict)
	if !strings.Contains(err.Error(), "«"+definicionDePrueba(t).Submodulo("objetivos").Nombre+"»") {
		t.Fatalf("el mensaje debe nombrar el submódulo con pendientes: %v", err)
	}
	if estado := repo.planes[iniciativa.ID].EstadoPlanID; estado != "ep-ef" {
		t.Fatalf("la iniciativa no debió cambiar de estado y quedó %s", estado)
	}
}

// cambiarDatos cambia los datos guardados de un submódulo sin pasar por el servicio.
func cambiarDatos(repo *repoMemoria, iniciativaID, clave string, datos Datos) {
	for _, documento := range repo.submodulos {
		if documento.PlanID == iniciativaID && documento.Submodulo == clave {
			documento.Datos = datos
		}
	}
}

func TestObtenerSubmoduloIncluyeLosSubmodulosDeLasQueDepende(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)

	cronograma, err := servicio.ObtenerSubmodulo(iniciativa.ID, "cronograma")
	if err != nil {
		t.Fatal(err)
	}
	if _, existe := cronograma.Relacionadas["cadena_valor"]; !existe || len(cronograma.Relacionadas) != 1 {
		t.Fatalf("relacionadas = %v", cronograma.Relacionadas)
	}
	participantes, err := servicio.ObtenerSubmodulo(iniciativa.ID, "participantes")
	if err != nil {
		t.Fatal(err)
	}
	if len(participantes.Relacionadas) != 0 {
		t.Fatalf("participantes no depende de otros submódulos: %v", participantes.Relacionadas)
	}
}

// El resumen dice qué submódulos leen datos de cada uno (Cronograma lee de Cadena de valor).
func TestElResumenTraeLosDependientesDeCadaSubmodulo(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)

	for _, submodulo := range submodulosDe(iniciativa) {
		if submodulo.Dependientes == nil {
			t.Fatalf("%s: dependientes nil", submodulo.Clave)
		}
		if submodulo.Clave == "cadena_valor" && !contiene(submodulo.Dependientes, "cronograma") {
			t.Fatalf("cronograma lee de cadena_valor: %v", submodulo.Dependientes)
		}
		if submodulo.Clave == "participantes" && len(submodulo.Dependientes) != 0 {
			t.Fatalf("nadie lee de participantes: %v", submodulo.Dependientes)
		}
	}
}

// ---------------------------------------------------------------------------
// Máquina de estados
// ---------------------------------------------------------------------------

// El recorrido completo: el Jefe devuelve, la OAP observa (la iniciativa pasa sola a Por
// corregir), el Jefe desaprueba un aprobado y, al final, todo aprobado (pasa sola a Aprobado).
func TestFlujoCompletoHastaAprobado(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	id := crearIniciativaDePrueba(t, servicio).ID
	formularTodo(t, servicio, id)

	// 02 → 03: desde que la envía, el formulador ya no edita.
	if resumen := ejecutar(t, servicio, id, AccionEnviar, formulador()); resumen.Estado.Codigo != EstadoIniciativaPorAprobar || resumen.Permisos.Editar {
		t.Fatalf("enviar: %+v", resumen)
	}
	_, err := intentarGuardar(t, servicio, id, "nombre", iniciativaDeEjemplo(t)["nombre"])
	exigirStatus(t, err, http.StatusConflict)

	// 03: el Jefe observa; con observaciones no aprueba, devuelve.
	comentar(t, servicio, id, "problematica", observacion(jefe(), "Precisar la fuente."))
	if acciones := accionesDe(resumenDe(t, servicio, id, jefe()).Permisos); strings.Join(acciones, ",") != "devolver,aprobar!" {
		t.Fatalf("acciones del Jefe con observaciones = %v", acciones)
	}
	_, err = servicio.EjecutarAccion(id, AccionAprobar, SolicitudAccion{Actor: jefe()})
	exigirStatus(t, err, http.StatusConflict)
	ejecutar(t, servicio, id, AccionDevolver, jefe())

	// 06 → 03 → 04 → 05: corrige, reenvía, el Jefe aprueba y la OAP inicia la revisión.
	_, err = servicio.EjecutarAccion(id, AccionEnviar, SolicitudAccion{Actor: formulador()})
	exigirStatus(t, err, http.StatusConflict) // problemática sigue Con observaciones
	guardar(t, servicio, id, "problematica", iniciativaDeEjemplo(t)["problematica"])
	ejecutar(t, servicio, id, AccionEnviar, formulador())
	ejecutar(t, servicio, id, AccionAprobar, jefe())
	if resumen := ejecutar(t, servicio, id, AccionIniciarRevision, oap()); resumen.Estado.Codigo != EstadoIniciativaEnRevision {
		t.Fatalf("iniciar revisión: %+v", resumen.Estado)
	}

	// 05: aprueba todo menos problemática; el último veredicto (Con observaciones) la pasa a 06.
	for _, submodulo := range definicionDePrueba(t).Submodulos() {
		if submodulo.Clave() != "problematica" {
			if vista := aprobarSubmodulo(t, servicio, id, submodulo.Clave()); vista.EstadoIniciativa != nil {
				t.Fatalf("la iniciativa cambió antes del último veredicto: %+v", vista.EstadoIniciativa)
			}
		}
	}
	vista := comentar(t, servicio, id, "problematica", observacion(oap(), "Falta la fuente del dato."))
	if vista.EstadoIniciativa == nil || vista.EstadoIniciativa.Codigo != EstadoIniciativaPorCorregir {
		t.Fatalf("el último veredicto debió pasarla a Por corregir: %+v", vista.EstadoIniciativa)
	}

	// 06: solo se edita lo no aprobado.
	_, err = intentarGuardar(t, servicio, id, "nombre", iniciativaDeEjemplo(t)["nombre"])
	exigirStatus(t, err, http.StatusConflict)
	guardar(t, servicio, id, "problematica", iniciativaDeEjemplo(t)["problematica"])
	ejecutar(t, servicio, id, AccionEnviar, formulador())

	// 03: el Jefe observa un aprobado (lo desaprueba) y la devuelve; se corrige y vuelve a subir.
	comentar(t, servicio, id, "nombre", observacion(jefe(), "Ajustar el nombre."))
	ejecutar(t, servicio, id, AccionDevolver, jefe())
	guardar(t, servicio, id, "nombre", iniciativaDeEjemplo(t)["nombre"])
	ejecutar(t, servicio, id, AccionEnviar, formulador())
	ejecutar(t, servicio, id, AccionAprobar, jefe())
	ejecutar(t, servicio, id, AccionIniciarRevision, oap())

	// 05: con los dos que faltaban aprobados, pasa sola a Aprobado.
	aprobarSubmodulo(t, servicio, id, "nombre")
	vista = aprobarSubmodulo(t, servicio, id, "problematica")
	if vista.EstadoIniciativa == nil || vista.EstadoIniciativa.Codigo != EstadoIniciativaAprobada {
		t.Fatalf("con todo aprobado debió pasar a Aprobado: %+v", vista.EstadoIniciativa)
	}
	for _, actor := range []Actor{formulador(), jefe(), oap()} {
		if permisos := resumenDe(t, servicio, id, actor).Permisos; permisos.Editar || permisos.Observar || len(permisos.Acciones) != 0 {
			t.Fatalf("en Aprobado nadie hace nada: %+v", permisos)
		}
	}
	_, err = servicio.AgregarObservacion(id, "nombre", observacion(oap(), "Después de aprobada."))
	exigirStatus(t, err, http.StatusConflict)
}

// La OAP rechaza con motivo; la rechazada no se inactiva, nadie la edita y solo la OAP se
// retracta, lo que la devuelve a Por corregir.
func TestRechazoYRetraccion(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	id := crearIniciativaDePrueba(t, servicio).ID
	llevarARevision(t, servicio, id)

	for _, motivo := range []string{"", "   ", strings.Repeat("a", MaximoTextoObservacion+1)} {
		_, err := servicio.EjecutarAccion(id, AccionRechazar, SolicitudAccion{Actor: oap(), Motivo: motivo})
		exigirStatus(t, err, http.StatusBadRequest)
	}
	if resumen := ejecutar(t, servicio, id, AccionRechazar, oap()); resumen.Estado.Codigo != EstadoIniciativaRechazada {
		t.Fatalf("rechazar: %+v", resumen.Estado)
	}
	if !repo.planes[id].Activo {
		t.Fatal("la rechazada no se inactiva")
	}
	_, err := intentarGuardar(t, servicio, id, "nombre", iniciativaDeEjemplo(t)["nombre"])
	exigirStatus(t, err, http.StatusConflict)
	listado, err := servicio.ListarIniciativas(formulador())
	if err != nil || len(listado) != 1 || listado[0].Accion != AccionSubmoduloConsultar {
		t.Fatalf("la rechazada sigue en el listado, para consultarla: %+v %v", listado, err)
	}

	_, err = servicio.EjecutarAccion(id, AccionRetractarRechazo, SolicitudAccion{Actor: jefe()})
	exigirStatus(t, err, http.StatusForbidden)
	if acciones := accionesDe(resumenDe(t, servicio, id, oap()).Permisos); strings.Join(acciones, ",") != AccionRetractarRechazo {
		t.Fatalf("en 08 la OAP solo se retracta: %v", acciones)
	}
	if resumen := ejecutar(t, servicio, id, AccionRetractarRechazo, oap()); resumen.Estado.Codigo != EstadoIniciativaPorCorregir {
		t.Fatalf("retractar: %+v", resumen.Estado)
	}
	if acciones := accionesDe(resumenDe(t, servicio, id, formulador()).Permisos); strings.Join(acciones, ",") != AccionEnviar {
		t.Fatalf("tras la retracción el formulador la reenvía: %v", acciones)
	}
}

// En cada estado solo actúa su carril: los demás consultan.
func TestCadaEstadoLoTrabajaSuCarril(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	id := crearIniciativaDePrueba(t, servicio).ID
	formularTodo(t, servicio, id)

	esperado := func(estado string, actor Actor, editar, observar, aprobar bool, acciones string) {
		t.Helper()
		resumen := resumenDe(t, servicio, id, actor)
		if resumen.Estado.Codigo != estado {
			t.Fatalf("la prueba supone %s y está %s", estado, resumen.Estado.Codigo)
		}
		p := resumen.Permisos
		if p.Editar != editar || p.Observar != observar || p.AprobarSubmodulos != aprobar || strings.Join(accionesDe(p), ",") != acciones {
			t.Fatalf("%s, %s: editar=%v observar=%v aprobar=%v acciones=%v", estado, actor.Usuario, p.Editar, p.Observar, p.AprobarSubmodulos, accionesDe(p))
		}
	}
	esperado(EstadoIniciativaEnFormulacion, formulador(), true, false, false, AccionEnviar)
	esperado(EstadoIniciativaEnFormulacion, jefe(), true, false, false, AccionEnviar)
	esperado(EstadoIniciativaEnFormulacion, oap(), false, false, false, "")

	ejecutar(t, servicio, id, AccionEnviar, formulador())
	esperado(EstadoIniciativaPorAprobar, formulador(), false, false, false, "")
	esperado(EstadoIniciativaPorAprobar, jefe(), false, true, false, "devolver!,aprobar")
	esperado(EstadoIniciativaPorAprobar, oap(), false, false, false, "")

	ejecutar(t, servicio, id, AccionAprobar, jefe())
	esperado(EstadoIniciativaPorRevisar, jefe(), false, false, false, "")
	esperado(EstadoIniciativaPorRevisar, oap(), false, false, false, AccionIniciarRevision)

	ejecutar(t, servicio, id, AccionIniciarRevision, oap())
	esperado(EstadoIniciativaEnRevision, formulador(), false, false, false, "")
	esperado(EstadoIniciativaEnRevision, jefe(), false, false, false, "")
	esperado(EstadoIniciativaEnRevision, oap(), false, true, true, AccionRechazar)
}

// Un Jefe no actúa ni consulta las iniciativas de otras unidades.
func TestElJefeDeOtraUnidadNiActuaNiConsulta(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	id := crearEnDependencia(t, servicio, "dep-2").ID
	formularTodo(t, servicio, id)
	ejecutar(t, servicio, id, AccionEnviar, formulador())

	_, err := servicio.Resumen(id, jefe())
	exigirStatus(t, err, http.StatusForbidden)
	_, err = servicio.EjecutarAccion(id, AccionAprobar, SolicitudAccion{Actor: jefe()})
	exigirStatus(t, err, http.StatusForbidden)
	_, err = servicio.AgregarObservacion(id, "nombre", observacion(jefe(), "No es de su unidad."))
	exigirStatus(t, err, http.StatusForbidden)
}

// Con varios roles se suman los carriles: el mismo usuario observa como Jefe en 03 y como OAP en
// 05, y cada observación queda firmada con el rol de ese carril. Como tiene un rol de la OAP,
// consulta las iniciativas de todas las unidades.
func TestConVariosRolesSeSumanLosCarriles(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	id := crearIniciativaDePrueba(t, servicio).ID
	formularTodo(t, servicio, id)
	ejecutar(t, servicio, id, AccionEnviar, formulador())

	comoJefe := comentar(t, servicio, id, "nombre", observacion(jefeYOAP(), "Como Jefe."))
	if rol := comoJefe.Observaciones[len(comoJefe.Observaciones)-1].Rol; rol != RolJefeDependencia {
		t.Fatalf("en 03 firma como Jefe y firmó %s", rol)
	}
	ejecutar(t, servicio, id, AccionDevolver, jefeYOAP())
	guardar(t, servicio, id, "nombre", iniciativaDeEjemplo(t)["nombre"])
	ejecutar(t, servicio, id, AccionEnviar, formulador())
	ejecutar(t, servicio, id, AccionAprobar, jefeYOAP())
	ejecutar(t, servicio, id, AccionIniciarRevision, jefeYOAP())

	comoOAP := comentar(t, servicio, id, "nombre", observacion(jefeYOAP(), "Como OAP."))
	if rol := comoOAP.Observaciones[len(comoOAP.Observaciones)-1].Rol; rol != RolAsistentePlaneacion {
		t.Fatalf("en 05 firma como OAP y firmó %s", rol)
	}
	otra := crearEnDependencia(t, servicio, "dep-2")
	if _, err := servicio.Resumen(otra.ID, jefeYOAP()); err != nil {
		t.Fatalf("con un rol de la OAP consulta las de todas las unidades: %v", err)
	}
}

// Fuera de las fechas de la unidad nadie hace nada, ni el Jefe ni la OAP: todo queda congelado.
func TestFueraDeFechasNadieActua(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	id := crearIniciativaDePrueba(t, servicio).ID
	formularTodo(t, servicio, id)
	ejecutar(t, servicio, id, AccionEnviar, formulador())
	cerrarFechas(repo)

	resumen := resumenDe(t, servicio, id, jefe())
	if !resumen.FueraDeFechas || resumen.Permisos.Observar || strings.Join(accionesDe(resumen.Permisos), ",") != "devolver!,aprobar!" {
		t.Fatalf("fuera de fechas el Jefe solo consulta: %+v", resumen)
	}
	if resumen := resumenDe(t, servicio, id, oap()); resumen.FueraDeFechas {
		t.Fatal("la OAP no tiene nada que hacer en 03: no se le avisa de las fechas")
	}
	_, err := servicio.EjecutarAccion(id, AccionAprobar, SolicitudAccion{Actor: jefe()})
	exigirFueraDeFechas(t, err)
	_, err = servicio.AgregarObservacion(id, "nombre", observacion(jefe(), "Fuera de fechas."))
	exigirFueraDeFechas(t, err)
}

func TestAccionInexistenteResponde400(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	_, err := servicio.EjecutarAccion(iniciativa.ID, "borrar", SolicitudAccion{Actor: formulador()})
	exigirStatus(t, err, http.StatusBadRequest)
}

// ---------------------------------------------------------------------------
// Submódulos aprobados
// ---------------------------------------------------------------------------

// llevarAPorCorregirConObjetivosAprobados deja la iniciativa en Por corregir con Problemática
// observada y los demás submódulos aprobados.
func llevarAPorCorregirConObjetivosAprobados(t *testing.T, servicio *Servicio) string {
	t.Helper()
	id := crearIniciativaDePrueba(t, servicio).ID
	llevarARevision(t, servicio, id)
	for _, submodulo := range definicionDePrueba(t).Submodulos() {
		if submodulo.Clave() != "problematica" {
			aprobarSubmodulo(t, servicio, id, submodulo.Clave())
		}
	}
	comentar(t, servicio, id, "problematica", observacion(oap(), "Agregar la causa que falta."))
	return id
}

// Un aprobado no se revalida al guardar otro submódulo ni al enviar: quien revisa decide si lo
// vuelve a observar.
func TestUnAprobadoNoSeRevalidaNiSeEdita(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	id := llevarAPorCorregirConObjetivosAprobados(t, servicio)

	guardar(t, servicio, id, "problematica", problematicaConOtraCausa(t))
	if estado := estadoDeSubmodulo(t, servicio, id, "objetivos"); estado != EstadoSubmoduloAprobado {
		t.Fatalf("objetivos sigue aprobado aunque le falte la nueva causa y quedó %s", estado)
	}
	if resumen := ejecutar(t, servicio, id, AccionEnviar, formulador()); resumen.Estado.Codigo != EstadoIniciativaPorAprobar {
		t.Fatalf("se envía sin revalidar los aprobados: %+v", resumen.Estado)
	}
	_, err := intentarGuardar(t, servicio, id, "objetivos", iniciativaDeEjemplo(t)["objetivos"])
	exigirStatus(t, err, http.StatusConflict)
}

// En Por corregir el resumen trae la última observación de cada submódulo Con observaciones.
func TestEnPorCorregirElResumenTraeLaUltimaObservacion(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	id := llevarAPorCorregirConObjetivosAprobados(t, servicio)

	for _, submodulo := range submodulosDe(resumenDe(t, servicio, id, formulador())) {
		switch {
		case submodulo.Clave == "problematica":
			if submodulo.UltimaObservacion == nil || submodulo.UltimaObservacion.Texto != "Agregar la causa que falta." || submodulo.Accion != AccionSubmoduloFormular {
				t.Fatalf("problemática: %+v", submodulo)
			}
		case submodulo.UltimaObservacion != nil || submodulo.Accion != AccionSubmoduloConsultar:
			t.Fatalf("%s está aprobado: se consulta y no trae observación por atender: %+v", submodulo.Clave, submodulo)
		}
	}
}

// ---------------------------------------------------------------------------
// Observaciones
// ---------------------------------------------------------------------------

func observacion(actor Actor, texto string) SolicitudObservacion {
	return SolicitudObservacion{Actor: actor, Texto: texto}
}

func comentar(t *testing.T, servicio *Servicio, iniciativaID, clave string, solicitud SolicitudObservacion) *SubmoduloVista {
	t.Helper()
	vista, err := servicio.AgregarObservacion(iniciativaID, clave, solicitud)
	if err != nil {
		t.Fatal(err)
	}
	return vista
}

// aprobarSubmodulo da el veredicto Aprobado como la OAP, con la versión vigente.
func aprobarSubmodulo(t *testing.T, servicio *Servicio, iniciativaID, clave string) *SubmoduloVista {
	t.Helper()
	actual, err := servicio.ObtenerSubmodulo(iniciativaID, clave)
	if err != nil {
		t.Fatal(err)
	}
	vista, err := servicio.RevisarSubmodulo(iniciativaID, clave, SolicitudRevisarSubmodulo{Actor: oap(), Version: actual.Submodulo.Version})
	if err != nil {
		t.Fatal(err)
	}
	return vista
}

func TestLaOAPObservaYElSubmoduloQuedaConObservacionesSinCambiarLaVersion(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	id := crearIniciativaDePrueba(t, servicio).ID
	llevarARevision(t, servicio, id)
	antes, _ := servicio.ObtenerSubmodulo(id, "poblacion")

	vista := comentar(t, servicio, id, "poblacion", observacion(oap(), "  Revisar el total de juventud.  "))

	if vista.Estado.Codigo != EstadoSubmoduloConObservaciones || vista.Version != antes.Submodulo.Version {
		t.Fatalf("estado %s, versión %d (antes %d)", vista.Estado.Codigo, vista.Version, antes.Submodulo.Version)
	}
	esperada := Observacion{ID: vista.Observaciones[0].ID, Texto: "Revisar el total de juventud.", Rol: RolPlaneacion, Usuario: documentoOAP, UsuarioNombre: personasDePrueba[documentoOAP]}
	if len(vista.Observaciones) != 1 || vista.Observaciones[0] != esperada {
		t.Fatalf("observaciones = %+v", vista.Observaciones)
	}
}

func TestElFormuladorComentaMientrasFormulaSinChocarConSuEdicion(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	abierta, _ := servicio.ObtenerSubmodulo(iniciativa.ID, "nombre")

	vista := comentar(t, servicio, iniciativa.ID, "nombre", observacion(formulador(), "El nombre se ajustó con la OAP."))
	if vista.Estado.Codigo != EstadoSubmoduloSinFormular || vista.Version != abierta.Submodulo.Version || vista.Observaciones[0].Rol != RolAsistenteDependencia {
		t.Fatalf("el comentario del formulador no cambia estado ni versión: %s v%d %+v", vista.Estado.Codigo, vista.Version, vista.Observaciones)
	}

	// Guarda con la versión que tenía abierta antes de comentar: no hay 409.
	if _, err := servicio.GuardarSubmodulo(iniciativa.ID, "nombre", SolicitudGuardarSubmodulo{Actor: formulador(), Version: abierta.Submodulo.Version, Datos: iniciativaDeEjemplo(t)["nombre"]}); err != nil {
		t.Fatalf("guardar tras comentar: %v", err)
	}
}

func TestLaOAPNoObservaUnSubmoduloSinFormular(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	repo.planes[iniciativa.ID].EstadoPlanID = "ep-er" // en revisión, con todo sin formular

	_, err := servicio.AgregarObservacion(iniciativa.ID, "nombre", observacion(oap(), "x"))
	exigirStatus(t, err, http.StatusConflict)
}

// Fuera de su carril nadie publica: la OAP no comenta mientras se formula y el formulador no
// observa en revisión.
func TestSoloPublicaElCarrilDelEstado(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	id := crearIniciativaDePrueba(t, servicio).ID
	_, err := servicio.AgregarObservacion(id, "nombre", observacion(oap(), "Mientras se formula."))
	exigirStatus(t, err, http.StatusForbidden)

	llevarARevision(t, servicio, id)
	_, err = servicio.AgregarObservacion(id, "nombre", observacion(formulador(), "En revisión."))
	exigirStatus(t, err, http.StatusForbidden)
}

func TestObservacionInvalidaResponde400(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	largo := make([]rune, MaximoTextoObservacion+1)
	for i := range largo {
		largo[i] = 'á'
	}
	casos := map[string]SolicitudObservacion{
		"sin texto":                   observacion(formulador(), "   "),
		"texto muy largo":             observacion(formulador(), string(largo)),
		"sin roles":                   observacion(Actor{Usuario: responsableDePrueba}, "x"),
		"sin usuario":                 observacion(Actor{Roles: []string{RolAsistenteDependencia}}, "x"),
		"usuario que no es documento": observacion(Actor{Usuario: "123,Activo:false", Roles: []string{RolAsistenteDependencia}}, "x"),
	}
	for nombre, solicitud := range casos {
		if _, err := servicio.AgregarObservacion(iniciativa.ID, "nombre", solicitud); err == nil {
			t.Errorf("%s: se esperaba error", nombre)
		} else {
			exigirStatus(t, err, http.StatusBadRequest)
		}
	}
}

func TestLaObservacionSePublicaSinNombreSiTercerosNoResponde(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	tercerosCaidos(servicio)

	vista := comentar(t, servicio, iniciativa.ID, "nombre", observacion(formulador(), "El nombre se ajustó con la OAP."))
	if len(vista.Observaciones) != 1 || vista.Observaciones[0].Usuario != responsableDePrueba || vista.Observaciones[0].UsuarioNombre != "" {
		t.Fatalf("observaciones = %+v", vista.Observaciones)
	}
}

// ---------------------------------------------------------------------------
// Plantillas y consultas
// ---------------------------------------------------------------------------

// definicionesPorPlantilla sirve una definición distinta (o un error) por plantilla.
type definicionesPorPlantilla map[string]*Definicion

func (d definicionesPorPlantilla) Definicion(id string) (*Definicion, error) {
	if definicion, existe := d[id]; existe {
		return definicion, nil
	}
	return nil, errorServicio("La plantilla %s no está organizada por módulos.", id)
}

func TestSoloSeOfrecenPlantillasUsables(t *testing.T) {
	repo := nuevoRepoMemoria()
	for _, id := range []string{"completa", "vacia", "sin-modulos"} {
		repo.planes[id] = &Plan{ID: id, Nombre: id, TipoPlanID: tipoPlanDePrueba, Formato: true, Activo: true}
	}
	definiciones := definicionesPorPlantilla{
		"completa": definicionDePrueba(t),
		"vacia":    {PlantillaID: "vacia", TipoPlanID: tipoPlanDePrueba},
	}
	servicio := NuevoServicio(repo, definiciones, catalogosDePrueba(), nil)

	plantillas, err := servicio.Plantillas()
	if err != nil {
		t.Fatal(err)
	}
	if len(plantillas) != 1 || plantillas[0].ID != "completa" {
		t.Fatalf("plantillas ofrecidas = %+v", plantillas)
	}
}

func TestUnaPlantillaNoSeConsultaComoIniciativa(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	repo.planes["plantilla-1"] = &Plan{ID: "plantilla-1", TipoPlanID: tipoPlanDePrueba, Formato: true, Activo: true}

	_, err := servicio.Resumen("plantilla-1", oap())
	exigirStatus(t, err, http.StatusNotFound)
}

// Un plan de otro tipo (p. ej. un proyecto MACRO, PRI_SP) no es una iniciativa,
// aunque apunte a una plantilla.
func TestUnPlanDeOtroTipoNoEsIniciativa(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	repo.planes["macro-1"] = &Plan{ID: "macro-1", Nombre: "Proyecto MACRO", TipoPlanID: "tipo-macro", FormatoID: "plantilla-1", Activo: true}

	iniciativas, err := servicio.ListarIniciativas(oap())
	if err != nil {
		t.Fatal(err)
	}
	if len(iniciativas) != 1 || iniciativas[0].ID != iniciativa.ID {
		t.Fatalf("iniciativas listadas = %+v", iniciativas)
	}
	_, err = servicio.Resumen("macro-1", oap())
	exigirStatus(t, err, http.StatusNotFound)
}

// El plan y sus submódulos se piden a la vez: si el plan no existe, eso es lo que
// se responde, aunque también falle el listado de submódulos.
func TestIniciativaInexistenteResponde404AunqueFalleElListado(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	repo.fallarListado = true

	_, err := servicio.Resumen("no-existe", oap())
	exigirStatus(t, err, http.StatusNotFound)
}

func TestSiFallaElListadoDeSubmodulosSeReporta(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	repo.fallarListado = true

	_, err := servicio.Resumen(iniciativa.ID, oap())
	exigirStatus(t, err, http.StatusBadGateway)
}

// Si no se pueden consultar las unidades del usuario, lo que las necesita responde el error.
func TestSiFallanLasUnidadesSeReporta(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	servicio.unidades = &unidadesFalsas{err: errorServicio("Terceros no respondió")}

	_, err := servicio.Resumen(iniciativa.ID, formulador())
	exigirStatus(t, err, http.StatusBadGateway)
}
