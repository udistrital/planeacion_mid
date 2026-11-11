package iniciativaInversionHelper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// repoMemoria imita a planes_crud en memoria, incluido el bloqueo optimista de
// submodulo_iniciativa (responde 409 si la versión no es la vigente).
type repoMemoria struct {
	planes           map[string]*Plan
	submodulos       map[string]*SubmoduloIniciativa
	catalogos        map[string][]Estado
	contador         int
	fallarSubmodulos bool
	fallarListado    bool
}

func nuevoRepoMemoria() *repoMemoria {
	estado := func(id, codigo, nombre string) Estado { return Estado{ID: id, Codigo: codigo, Nombre: nombre} }
	return &repoMemoria{
		planes:     map[string]*Plan{},
		submodulos: map[string]*SubmoduloIniciativa{},
		catalogos: map[string][]Estado{
			RecursoTipoPlan: {estado(tipoPlanDePrueba, CodigoTipoPlanIniciativa, "Iniciativa de inversión")},
			RecursoEstadoPlan: {
				estado("ep-ef", EstadoIniciativaEnFormulacion, "En formulación"),
				estado("ep-er", EstadoIniciativaEnRevision, "Enviado a revisión"),
				estado("ep-pc", EstadoIniciativaPorCorregir, "Por corregir"),
				estado("ep-av", EstadoIniciativaAvalado, "Avalado"),
				estado("ep-na", EstadoIniciativaNoAvalado, "No avalado"),
			},
			RecursoEstadoSubmodulo: {
				estado("es-sf", EstadoSubmoduloSinFormular, "Sin formular"),
				estado("es-for", EstadoSubmoduloFormulado, "Formulado"),
				estado("es-co", EstadoSubmoduloConObservaciones, "Con observaciones"),
				estado("es-so", EstadoSubmoduloSinObservaciones, "Sin observaciones"),
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
		ID: r.nuevoID("obs"), Texto: nueva.Texto, Rol: nueva.Rol, Usuario: nueva.Usuario,
	})
	if nueva.EstadoSubmoduloID != "" {
		submodulo.EstadoSubmoduloID = nueva.EstadoSubmoduloID
	}
	copia := *submodulo
	copia.Observaciones = append([]Observacion(nil), submodulo.Observaciones...)
	return &copia, nil
}

type definicionFija struct{ definicion *Definicion }

func (d definicionFija) Definicion(string) (*Definicion, error) { return d.definicion, nil }

// ---------------------------------------------------------------------------

func nuevoServicioDePrueba(t *testing.T) (*Servicio, *repoMemoria) {
	t.Helper()
	repo := nuevoRepoMemoria()
	return NuevoServicio(repo, definicionFija{definicionDePrueba(t)}, catalogosDePrueba()), repo
}

func crearIniciativaDePrueba(t *testing.T, servicio *Servicio) *ResumenIniciativa {
	t.Helper()
	resumen, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{PlantillaID: "plantilla-1", DependenciaID: "dep-1", Vigencias: []byte(`[{"Id":1,"Nombre":"2026"}]`)})
	if err != nil {
		t.Fatal(err)
	}
	return resumen
}

// guardar guarda el submódulo con la versión vigente y devuelve la vista.
func guardar(t *testing.T, servicio *Servicio, iniciativaID, clave string, datos Datos) *SubmoduloVista {
	t.Helper()
	actual, err := servicio.ObtenerSubmodulo(iniciativaID, clave)
	if err != nil {
		t.Fatal(err)
	}
	vista, err := servicio.GuardarSubmodulo(iniciativaID, clave, SolicitudGuardarSubmodulo{Version: actual.Submodulo.Version, Datos: datos})
	if err != nil {
		t.Fatal(err)
	}
	return vista
}

// formularTodo guarda todos los submódulos del ejemplo y exige que queden formulados.
func formularTodo(t *testing.T, servicio *Servicio, iniciativaID string) {
	t.Helper()
	iniciativa := iniciativaDeEjemplo(t)
	for _, clave := range []string{"nombre", "articulacion_pnd", "problematica", "participantes", "poblacion", "objetivos", "alternativas"} {
		if vista := guardar(t, servicio, iniciativaID, clave, iniciativa[clave]); vista.Estado.Codigo != EstadoSubmoduloFormulado {
			t.Fatalf("%s quedó %s: %+v", clave, vista.Estado.Codigo, vista.Errores)
		}
	}
}

func estadoDeSubmodulo(t *testing.T, servicio *Servicio, iniciativaID, clave string) string {
	t.Helper()
	resumen, err := servicio.Resumen(iniciativaID)
	if err != nil {
		t.Fatal(err)
	}
	for _, submodulo := range submodulosDe(resumen) {
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

// ---------------------------------------------------------------------------

func TestCrearIniciativaCreaUnSubmoduloSinFormularPorSubmoduloDeLaPlantilla(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	resumen := crearIniciativaDePrueba(t, servicio)

	if resumen.Estado.Codigo != EstadoIniciativaEnFormulacion || !resumen.PermiteEditar {
		t.Fatalf("estado inicial = %+v", resumen.Estado)
	}
	if len(submodulosDe(resumen)) != 7 {
		t.Fatalf("se esperaban 7 submódulos y hay %d", len(submodulosDe(resumen)))
	}
	for _, submodulo := range submodulosDe(resumen) {
		if submodulo.Estado.Codigo != EstadoSubmoduloSinFormular {
			t.Fatalf("%s quedó %s", submodulo.Clave, submodulo.Estado.Codigo)
		}
	}
	if len(resumen.AccionesPosibles) != 0 {
		t.Fatalf("no debería poder enviarse: %v", resumen.AccionesPosibles)
	}
}

func TestCrearIniciativaDesactivaElPlanSiFallanLosSubmodulos(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	repo.fallarSubmodulos = true

	if _, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{PlantillaID: "plantilla-1", DependenciaID: "dep-1"}); err == nil {
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
	_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{PlantillaID: "plantilla-1"})
	exigirStatus(t, err, http.StatusBadRequest)
}

func TestSinPlantillaSeCreaConLaMasReciente(t *testing.T) {
	repo := nuevoRepoMemoria()
	repo.planes["anterior"] = &Plan{ID: "anterior", TipoPlanID: tipoPlanDePrueba, Formato: true, Activo: true, FechaCreacion: "2026-10-01T00:00:00.000Z"}
	repo.planes["reciente"] = &Plan{ID: "reciente", TipoPlanID: tipoPlanDePrueba, Formato: true, Activo: true, FechaCreacion: "2026-11-01T00:00:00.000Z"}
	anterior, reciente := definicionDePrueba(t), definicionDePrueba(t)
	anterior.PlantillaID, reciente.PlantillaID = "anterior", "reciente"
	servicio := NuevoServicio(repo, definicionesPorPlantilla{"anterior": anterior, "reciente": reciente}, catalogosDePrueba())

	resumen, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{DependenciaID: "dep-1"})
	if err != nil {
		t.Fatal(err)
	}
	if resumen.PlantillaID != "reciente" || resumen.Nombre != nombreIniciativaInicial {
		t.Fatalf("plantilla = %s, nombre = %s", resumen.PlantillaID, resumen.Nombre)
	}
}

func TestSinPlantillasNoSeCreaLaIniciativa(t *testing.T) {
	servicio := NuevoServicio(nuevoRepoMemoria(), definicionesPorPlantilla{}, catalogosDePrueba())
	_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{DependenciaID: "dep-1"})
	exigirStatus(t, err, http.StatusConflict)
}

func TestListarNumeraYFiltraPorUnidad(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	crear := func(dependencia string) string {
		t.Helper()
		resumen, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{PlantillaID: "plantilla-1", DependenciaID: dependencia})
		if err != nil {
			t.Fatal(err)
		}
		return resumen.ID
	}
	desactivada, activa, otraUnidad := crear("dep-1"), crear("dep-1"), crear("dep-2")
	repo.planes[desactivada].Activo = false

	todas, err := servicio.ListarIniciativas(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(todas) != 2 || todas[0].ID != otraUnidad || todas[0].Numero != 3 || todas[1].ID != activa || todas[1].Numero != 2 {
		t.Fatalf("listado = %+v", todas)
	}
	if todas[1].Estado.Codigo != EstadoIniciativaEnFormulacion || !todas[1].PermiteEditar {
		t.Fatalf("estado = %+v", todas[1])
	}

	deUnaUnidad, err := servicio.ListarIniciativas([]string{"dep-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(deUnaUnidad) != 1 || deUnaUnidad[0].ID != activa {
		t.Fatalf("listado de dep-1 = %+v", deUnaUnidad)
	}
}

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

func TestGuardarConErroresGuardaElBorradorSinFormular(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	objetivos := iniciativaDeEjemplo(t)["objetivos"]
	objetivos["objetivo_central"] = ""

	vista := guardar(t, servicio, iniciativa.ID, "objetivos", objetivos)

	if vista.Estado.Codigo != EstadoSubmoduloSinFormular || len(vista.Errores) != 1 {
		t.Fatalf("estado=%s errores=%+v", vista.Estado.Codigo, vista.Errores)
	}
	if len(filasEn(t, vista.Datos, "objetivos_especificos")) != 3 {
		t.Fatal("el borrador debió guardarse aunque tenga errores")
	}
}

func TestGuardarConUnaVersionViejaResponde409(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	guardar(t, servicio, iniciativa.ID, "nombre", iniciativaDeEjemplo(t)["nombre"]) // versión pasa de 0 a 1

	_, err := servicio.GuardarSubmodulo(iniciativa.ID, "nombre", SolicitudGuardarSubmodulo{Version: 0, Datos: Datos{}})
	exigirStatus(t, err, http.StatusConflict)
}

func TestAgregarUnObjetivoDevuelveAlternativasASinFormular(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)

	objetivos := iniciativaDeEjemplo(t)["objetivos"]
	objetivos["objetivos_especificos"] = append(filasEn(t, objetivos, "objetivos_especificos"), map[string]interface{}{
		"id": "oe-4", "objetivo_especifico": "Nuevo objetivo.", "linea_base": "Sin medición.",
		"medios": []interface{}{map[string]interface{}{"id": "me-9", "descripcion": "Nuevo medio."}},
	})
	guardar(t, servicio, iniciativa.ID, "objetivos", objetivos)

	if estado := estadoDeSubmodulo(t, servicio, iniciativa.ID, "alternativas"); estado != EstadoSubmoduloSinFormular {
		t.Fatalf("alternativas debió volver a sin formular (falta la fila del nuevo objetivo) y está %s", estado)
	}
	if estado := estadoDeSubmodulo(t, servicio, iniciativa.ID, "participantes"); estado != EstadoSubmoduloFormulado {
		t.Fatalf("participantes no depende de objetivos y quedó %s", estado)
	}
}

func TestCambiarLaProblematicaNoAfectaLosObjetivos(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)

	problematica := iniciativaDeEjemplo(t)["problematica"]
	problematica["causas_directas"] = append(filasEn(t, problematica, "causas_directas"), map[string]interface{}{
		"id": "cd-4", "descripcion": "Nueva causa directa.", "linea_base": "Sin medición.",
		"causas_indirectas": []interface{}{map[string]interface{}{"id": "ci-9", "descripcion": "Nueva causa indirecta."}},
	})
	guardar(t, servicio, iniciativa.ID, "problematica", problematica)

	if estado := estadoDeSubmodulo(t, servicio, iniciativa.ID, "objetivos"); estado != EstadoSubmoduloFormulado {
		t.Fatalf("objetivos no se relaciona con las causas y quedó %s", estado)
	}
}

// Simula que al guardar Objetivos falló devolver Alternativas a Sin formular: el
// nuevo objetivo quedó guardado y Alternativas sigue Formulado sin la fila de ese
// objetivo. Enviar debe detectarlo aunque el estado guardado diga Formulado.
func TestEnviarRevalidaLosSubmodulosAunqueSuEstadoDigaFormulado(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)

	objetivos := iniciativaDeEjemplo(t)["objetivos"]
	objetivos["objetivos_especificos"] = append(filasEn(t, objetivos, "objetivos_especificos"), map[string]interface{}{
		"id": "oe-4", "objetivo_especifico": "Nuevo objetivo.", "linea_base": "Sin medición.",
		"medios": []interface{}{map[string]interface{}{"id": "me-9", "descripcion": "Nuevo medio."}},
	})
	for _, documento := range repo.submodulos {
		if documento.PlanID == iniciativa.ID && documento.Submodulo == "objetivos" {
			documento.Datos = objetivos
		}
	}
	if estado := estadoDeSubmodulo(t, servicio, iniciativa.ID, "alternativas"); estado != EstadoSubmoduloFormulado {
		t.Fatalf("la prueba supone alternativas Formulado y está %s", estado)
	}

	_, err := servicio.EjecutarAccion(iniciativa.ID, AccionEnviar)
	exigirStatus(t, err, http.StatusConflict)
	if !strings.Contains(err.Error(), "«Alternativas de solución»") {
		t.Fatalf("el mensaje debe nombrar el submódulo con pendientes: %v", err)
	}
	if estado := repo.planes[iniciativa.ID].EstadoPlanID; estado != "ep-ef" {
		t.Fatalf("la iniciativa no debió cambiar de estado y quedó %s", estado)
	}
}

func TestObtenerSubmoduloIncluyeLosSubmodulosDeLasQueDepende(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)

	submodulo, err := servicio.ObtenerSubmodulo(iniciativa.ID, "alternativas")
	if err != nil {
		t.Fatal(err)
	}
	if _, existe := submodulo.Relacionadas["objetivos"]; !existe || len(submodulo.Relacionadas) != 1 {
		t.Fatalf("relacionadas = %v", submodulo.Relacionadas)
	}
	objetivos, err := servicio.ObtenerSubmodulo(iniciativa.ID, "objetivos")
	if err != nil {
		t.Fatal(err)
	}
	if len(objetivos.Relacionadas) != 0 {
		t.Fatalf("objetivos ya no depende de otros submódulos: %v", objetivos.Relacionadas)
	}
}

func TestFlujoDeRevisionConObservacionesYAval(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)

	resumen, err := servicio.EjecutarAccion(iniciativa.ID, AccionEnviar)
	if err != nil || resumen.Estado.Codigo != EstadoIniciativaEnRevision || resumen.PermiteEditar {
		t.Fatalf("enviar: %+v %v", resumen, err)
	}
	_, err = servicio.GuardarSubmodulo(iniciativa.ID, "nombre", SolicitudGuardarSubmodulo{Version: 1, Datos: Datos{}})
	exigirStatus(t, err, http.StatusConflict) // en revisión no se edita

	for _, submodulo := range submodulosDe(resumen) {
		sinObservaciones(t, servicio, iniciativa.ID, submodulo.Clave)
	}
	comentar(t, servicio, iniciativa.ID, "problematica", dePlaneacion("Precisar la fuente."))

	resumen, err = servicio.EjecutarAccion(iniciativa.ID, AccionDevolver)
	if err != nil || resumen.Estado.Codigo != EstadoIniciativaPorCorregir {
		t.Fatalf("devolver: %+v %v", resumen, err)
	}
	if _, err := servicio.EjecutarAccion(iniciativa.ID, AccionEnviar); err == nil {
		t.Fatal("no debe enviarse con un submódulo con observaciones sin corregir")
	}

	comentar(t, servicio, iniciativa.ID, "problematica", delGestor("Se agregó la fuente."))
	guardar(t, servicio, iniciativa.ID, "problematica", iniciativaDeEjemplo(t)["problematica"])
	if _, err := servicio.EjecutarAccion(iniciativa.ID, AccionEnviar); err != nil {
		t.Fatalf("reenviar: %v", err)
	}
	sinObservaciones(t, servicio, iniciativa.ID, "problematica")
	resumen, err = servicio.EjecutarAccion(iniciativa.ID, AccionAvalar)
	if err != nil || resumen.Estado.Codigo != EstadoIniciativaAvalado || len(resumen.AccionesPosibles) != 0 {
		t.Fatalf("avalar: %+v %v", resumen, err)
	}

	problematica, err := servicio.ObtenerSubmodulo(iniciativa.ID, "problematica")
	if err != nil {
		t.Fatal(err)
	}
	historial := problematica.Submodulo.Observaciones
	if len(historial) != 2 || historial[0].Rol != "PLANEACION" || historial[1].Rol != "JEFE_DEPENDENCIA" {
		t.Fatalf("«Sin observaciones» no debe borrar el historial: %+v", historial)
	}
}

// ---------------------------------------------------------------------------
// Observaciones
// ---------------------------------------------------------------------------

func dePlaneacion(texto string) SolicitudObservacion {
	return SolicitudObservacion{Texto: texto, Rol: "PLANEACION", Usuario: "1012345678"}
}

func delGestor(texto string) SolicitudObservacion {
	return SolicitudObservacion{Texto: texto, Rol: "JEFE_DEPENDENCIA", Usuario: "79876543"}
}

func comentar(t *testing.T, servicio *Servicio, iniciativaID, clave string, solicitud SolicitudObservacion) *SubmoduloVista {
	t.Helper()
	vista, err := servicio.AgregarObservacion(iniciativaID, clave, solicitud)
	if err != nil {
		t.Fatal(err)
	}
	return vista
}

func sinObservaciones(t *testing.T, servicio *Servicio, iniciativaID, clave string) *SubmoduloVista {
	t.Helper()
	actual, err := servicio.ObtenerSubmodulo(iniciativaID, clave)
	if err != nil {
		t.Fatal(err)
	}
	vista, err := servicio.RevisarSubmodulo(iniciativaID, clave, SolicitudRevisarSubmodulo{Version: actual.Submodulo.Version})
	if err != nil {
		t.Fatal(err)
	}
	return vista
}

func TestPlaneacionComentaYElSubmoduloQuedaConObservacionesSinCambiarLaVersion(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)
	if _, err := servicio.EjecutarAccion(iniciativa.ID, AccionEnviar); err != nil {
		t.Fatal(err)
	}
	antes, _ := servicio.ObtenerSubmodulo(iniciativa.ID, "poblacion")

	vista := comentar(t, servicio, iniciativa.ID, "poblacion", dePlaneacion("  Revisar el total de juventud.  "))

	if vista.Estado.Codigo != EstadoSubmoduloConObservaciones || vista.Version != antes.Submodulo.Version {
		t.Fatalf("estado %s, versión %d (antes %d)", vista.Estado.Codigo, vista.Version, antes.Submodulo.Version)
	}
	esperada := Observacion{ID: vista.Observaciones[0].ID, Texto: "Revisar el total de juventud.", Rol: "PLANEACION", Usuario: "1012345678"}
	if len(vista.Observaciones) != 1 || vista.Observaciones[0] != esperada {
		t.Fatalf("observaciones = %+v", vista.Observaciones)
	}
}

func TestElGestorComentaMientrasFormulaSinChocarConSuEdicion(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	abierta, _ := servicio.ObtenerSubmodulo(iniciativa.ID, "nombre")

	vista := comentar(t, servicio, iniciativa.ID, "nombre", delGestor("El nombre se ajustó con la OAPC."))
	if vista.Estado.Codigo != EstadoSubmoduloSinFormular || vista.Version != abierta.Submodulo.Version {
		t.Fatalf("el comentario del gestor no cambia estado ni versión: %s v%d", vista.Estado.Codigo, vista.Version)
	}

	// Guarda con la versión que tenía abierta antes de comentar: no hay 409.
	if _, err := servicio.GuardarSubmodulo(iniciativa.ID, "nombre", SolicitudGuardarSubmodulo{Version: abierta.Submodulo.Version, Datos: iniciativaDeEjemplo(t)["nombre"]}); err != nil {
		t.Fatalf("guardar tras comentar: %v", err)
	}
}

func TestUnaIniciativaAvaladaYaNoRecibeObservaciones(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)
	resumen, err := servicio.EjecutarAccion(iniciativa.ID, AccionEnviar)
	if err != nil {
		t.Fatal(err)
	}
	for _, submodulo := range submodulosDe(resumen) {
		sinObservaciones(t, servicio, iniciativa.ID, submodulo.Clave)
	}
	if _, err := servicio.EjecutarAccion(iniciativa.ID, AccionAvalar); err != nil {
		t.Fatal(err)
	}

	_, err = servicio.AgregarObservacion(iniciativa.ID, "nombre", dePlaneacion("Después del aval."))
	exigirStatus(t, err, http.StatusConflict)
}

func TestPlaneacionNoComentaUnSubmoduloSinFormular(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	repo.planes[iniciativa.ID].EstadoPlanID = "ep-er" // en revisión, con todo sin formular

	_, err := servicio.AgregarObservacion(iniciativa.ID, "nombre", dePlaneacion("x"))
	exigirStatus(t, err, http.StatusConflict)
}

func TestObservacionInvalidaResponde400(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	largo := make([]rune, MaximoTextoObservacion+1)
	for i := range largo {
		largo[i] = 'á'
	}
	casos := map[string]SolicitudObservacion{
		"sin texto":       {Texto: "   ", Rol: "R", Usuario: "U"},
		"texto muy largo": {Texto: string(largo), Rol: "R", Usuario: "U"},
		"sin rol":         {Texto: "x", Usuario: "U"},
		"sin usuario":     {Texto: "x", Rol: "R"},
	}
	for nombre, solicitud := range casos {
		if _, err := servicio.AgregarObservacion(iniciativa.ID, "nombre", solicitud); err == nil {
			t.Errorf("%s: se esperaba error", nombre)
		} else {
			exigirStatus(t, err, http.StatusBadRequest)
		}
	}
}

func TestRevisarYaNoRecibeObservaciones(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)

	_, err := servicio.RevisarSubmodulo(iniciativa.ID, "nombre", SolicitudRevisarSubmodulo{Observaciones: []Observacion{{Texto: "x"}}})
	exigirStatus(t, err, http.StatusBadRequest)
}

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
	servicio := NuevoServicio(repo, definiciones, catalogosDePrueba())

	plantillas, err := servicio.Plantillas()
	if err != nil {
		t.Fatal(err)
	}
	if len(plantillas) != 1 || plantillas[0].ID != "completa" {
		t.Fatalf("plantillas ofrecidas = %+v", plantillas)
	}
}

func TestAccionInexistenteResponde400(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	_, err := servicio.EjecutarAccion(iniciativa.ID, "borrar")
	exigirStatus(t, err, http.StatusBadRequest)
}

func TestUnaPlantillaNoSeConsultaComoIniciativa(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	repo.planes["plantilla-1"] = &Plan{ID: "plantilla-1", TipoPlanID: tipoPlanDePrueba, Formato: true, Activo: true}

	_, err := servicio.Resumen("plantilla-1")
	exigirStatus(t, err, http.StatusNotFound)
}

// Un plan de otro tipo (p. ej. un proyecto MACRO, PRI_SP) no es una iniciativa,
// aunque apunte a una plantilla.
func TestUnPlanDeOtroTipoNoEsIniciativa(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	repo.planes["macro-1"] = &Plan{ID: "macro-1", Nombre: "Proyecto MACRO", TipoPlanID: "tipo-macro", FormatoID: "plantilla-1", Activo: true}

	iniciativas, err := servicio.ListarIniciativas(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(iniciativas) != 1 || iniciativas[0].ID != iniciativa.ID {
		t.Fatalf("iniciativas listadas = %+v", iniciativas)
	}
	_, err = servicio.Resumen("macro-1")
	exigirStatus(t, err, http.StatusNotFound)
}

// El plan y sus submódulos se piden a la vez: si el plan no existe, eso es lo que
// se responde, aunque también falle el listado de submódulos.
func TestIniciativaInexistenteResponde404AunqueFalleElListado(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	repo.fallarListado = true

	_, err := servicio.Resumen("no-existe")
	exigirStatus(t, err, http.StatusNotFound)
}

func TestSiFallaElListadoDeSubmodulosSeReporta(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	repo.fallarListado = true

	_, err := servicio.Resumen(iniciativa.ID)
	exigirStatus(t, err, http.StatusBadGateway)
}
