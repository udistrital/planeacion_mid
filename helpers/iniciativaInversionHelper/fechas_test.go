package iniciativaInversionHelper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// cerrarFechas deja a las dependencias de las pruebas sin fechas de formulación hoy: su único
// registro venció ayer.
func cerrarFechas(repo *repoMemoria) {
	repo.periodos = []map[string]interface{}{periodoDeFormulacion(-30, -1, "dep-1", "dep-2", "41", "42")}
}

func exigirFueraDeFechas(t *testing.T, err error) {
	t.Helper()
	exigirStatus(t, err, http.StatusConflict)
	if mensaje := ComoErrorNegocio(err).Mensaje; mensaje != mensajeFueraDeFechas {
		t.Fatalf("mensaje = %q, se esperaba %q", mensaje, mensajeFueraDeFechas)
	}
}

// guardarNombre guarda el Nombre del ejemplo como formulador, con la versión vigente, y devuelve el error.
func guardarNombre(t *testing.T, servicio *Servicio, iniciativaID string) error {
	t.Helper()
	_, err := intentarGuardar(t, servicio, iniciativaID, "nombre", iniciativaDeEjemplo(t)["nombre"])
	return err
}

func TestFueraDeFechasNoSeCreaLaIniciativa(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	cerrarFechas(repo)

	_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: formulador(), PlantillaID: "plantilla-1", DependenciaID: "dep-1"})
	exigirFueraDeFechas(t, err)
	if len(repo.planes) != 0 {
		t.Fatalf("no debió crearse ningún plan: %v", repo.planes)
	}
}

// Fuera de fechas la iniciativa solo se consulta: no se guarda, no se comenta ni se envía.
func TestFueraDeFechasSoloSeConsulta(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	formularTodo(t, servicio, iniciativa.ID)
	if resumen := resumenDe(t, servicio, iniciativa.ID, formulador()); resumen.FueraDeFechas || strings.Join(accionesDe(resumen.Permisos), ",") != AccionEnviar {
		t.Fatalf("dentro de fechas, con todo formulado, se envía: %+v", resumen)
	}
	cerrarFechas(repo)

	resumen := resumenDe(t, servicio, iniciativa.ID, formulador())
	if !resumen.FueraDeFechas || resumen.Permisos.Editar || resumen.Permisos.RolObservaciones != "" || strings.Join(accionesDe(resumen.Permisos), ",") != AccionEnviar+"!" {
		t.Fatalf("fuera de fechas no se edita ni se envía: %+v", resumen)
	}
	if resumen.Permisos.Acciones[0].Motivo != mensajeFueraDeFechas {
		t.Fatalf("enviar debe decir que está fuera de fechas: %+v", resumen.Permisos.Acciones[0])
	}
	exigirFueraDeFechas(t, guardarNombre(t, servicio, iniciativa.ID))
	_, err := servicio.EjecutarAccion(iniciativa.ID, AccionEnviar, SolicitudAccion{Actor: formulador()})
	exigirFueraDeFechas(t, err)
	_, err = servicio.AgregarObservacion(iniciativa.ID, "nombre", observacion(formulador(), "Se ajusta cuando se habiliten las fechas."))
	exigirFueraDeFechas(t, err)
	if estado := repo.planes[iniciativa.ID].EstadoPlanID; estado != "ep-ef" {
		t.Fatalf("la iniciativa no debió cambiar de estado y quedó %s", estado)
	}
}

// La revisión de la OAP también depende de las fechas de la unidad; quien no tiene nada que
// hacer en el estado no consulta las fechas.
func TestFueraDeFechasLaOAPNoRevisa(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	llevarARevision(t, servicio, iniciativa.ID)
	cerrarFechas(repo)

	repo.consultasPeriodos = 0
	if resumen := resumenDe(t, servicio, iniciativa.ID, formulador()); resumen.FueraDeFechas || repo.consultasPeriodos != 0 {
		t.Fatalf("el formulador no actúa en 05: fuera_de_fechas=%v consultas=%d", resumen.FueraDeFechas, repo.consultasPeriodos)
	}
	if resumen := resumenDe(t, servicio, iniciativa.ID, oap()); !resumen.FueraDeFechas || resumen.Permisos.AprobarSubmodulos {
		t.Fatalf("fuera de fechas la OAP no revisa: %+v", resumen)
	}
	actual, err := servicio.ObtenerSubmodulo(iniciativa.ID, "nombre")
	if err != nil {
		t.Fatal(err)
	}
	_, err = servicio.RevisarSubmodulo(iniciativa.ID, "nombre", SolicitudRevisarSubmodulo{Actor: oap(), Version: actual.Submodulo.Version})
	exigirFueraDeFechas(t, err)
	_, err = servicio.EjecutarAccion(iniciativa.ID, AccionRechazar, SolicitudAccion{Actor: oap(), Motivo: "Fuera de fechas."})
	exigirFueraDeFechas(t, err)
}

func TestFueraDeFechasNoSeSubenAnexos(t *testing.T) {
	servicio, repo, gestor, iniciativaID := servicioConAnexos(t)
	cerrarFechas(repo)

	_, err := servicio.SubirDocumento(iniciativaID, "anexos", SolicitudDocumento{Actor: formulador(), Clave: "carta_presentacion", Nombre: "carta.pdf", Archivo: archivoDePrueba(bytesMinimos)})
	exigirFueraDeFechas(t, err)
	if len(gestor.subidos) != 0 {
		t.Fatalf("no debió llamarse al gestor documental: %d archivos", len(gestor.subidos))
	}
}

// Cuentan las fechas de la unidad de la iniciativa, y el listado las lee una sola vez.
func TestCuentanLasFechasDeLaUnidadDeLaIniciativa(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	deUno := crearIniciativaDePrueba(t, servicio)
	deDos := crearEnDependencia(t, servicio, "dep-2")
	repo.periodos = []map[string]interface{}{periodoDeFormulacion(-30, -1, "dep-1"), periodoDeFormulacion(-1, 1, "dep-2")}

	exigirFueraDeFechas(t, guardarNombre(t, servicio, deUno.ID))
	if err := guardarNombre(t, servicio, deDos.ID); err != nil {
		t.Fatalf("dep-2 está dentro de sus fechas: %v", err)
	}

	repo.consultasPeriodos = 0
	listado, err := servicio.ListarIniciativas(formulador())
	if err != nil {
		t.Fatal(err)
	}
	fuera := map[string]bool{}
	for _, fila := range listado {
		fuera[fila.ID] = fila.FueraDeFechas
	}
	if len(listado) != 2 || !fuera[deUno.ID] || fuera[deDos.ID] {
		t.Fatalf("listado = %+v", listado)
	}
	if repo.consultasPeriodos != 1 {
		t.Fatalf("el listado leyó las fechas %d veces; se esperaba una", repo.consultasPeriodos)
	}
}

// Si no se pueden revisar las fechas, la consulta sale igual con la iniciativa cerrada, y guardar
// responde el error.
func TestSiNoSePuedenRevisarLasFechasSoloSeConsulta(t *testing.T) {
	casos := []struct {
		nombre  string
		romper  func(repo *repoMemoria)
		status  int
		detalle string
	}{
		{"el crud falla", func(repo *repoMemoria) { repo.fallarPeriodos = true }, http.StatusBadGateway, "falla simulada"},
		{"un registro está mal formado", func(repo *repoMemoria) {
			repo.periodos = append([]map[string]interface{}{{"_id": "roto", "unidades_interes": "no es json"}}, repo.periodos...)
		}, http.StatusInternalServerError, "roto"},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			servicio, repo := nuevoServicioDePrueba(t)
			iniciativa := crearIniciativaDePrueba(t, servicio)
			caso.romper(repo)

			resumen, err := servicio.Resumen(iniciativa.ID, formulador())
			if err != nil || !resumen.FueraDeFechas {
				t.Fatalf("resumen = %+v, err = %v", resumen, err)
			}
			listado, err := servicio.ListarIniciativas(formulador())
			if err != nil || len(listado) != 1 || !listado[0].FueraDeFechas {
				t.Fatalf("listado = %+v, err = %v", listado, err)
			}
			err = guardarNombre(t, servicio, iniciativa.ID)
			exigirStatus(t, err, caso.status)
			exigirMensajeParaUsuario(t, err, caso.detalle)
		})
	}
}

// Sin el tipo de seguimiento FIPI_SP no se formula; en cuanto se carga, funciona sin reiniciar.
func TestSinTipoDeSeguimientoNoSeFormulaHastaQueSeCargue(t *testing.T) {
	servicio, repo := nuevoServicioDePrueba(t)
	tipos := repo.catalogos[RecursoTipoSeguimiento]
	delete(repo.catalogos, RecursoTipoSeguimiento)

	_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{Actor: formulador(), PlantillaID: "plantilla-1", DependenciaID: "dep-1"})
	exigirStatus(t, err, http.StatusInternalServerError)
	exigirMensajeParaUsuario(t, err, "cargar-estados.js")

	repo.catalogos[RecursoTipoSeguimiento] = tipos
	crearIniciativaDePrueba(t, servicio)
}

// ClienteCrud pide los registros activos del tipo de seguimiento con los campos que usa la
// validación, el más reciente primero, y trata como error la respuesta con Success: false.
func TestClienteCrudPeriodosActivos(t *testing.T) {
	periodo := map[string]interface{}{"_id": "p1", "fecha_inicio": "2026-11-20T00:00:00.000Z", "fecha_fin": "2026-12-10T00:00:00.000Z", "unidades_interes": `[{"Id":14,"Nombre":"FACULTAD DE INGENIERIA"}]`}
	exito := true
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consulta := r.URL.Query()
		if r.URL.Path != "/periodo-seguimiento" || consulta.Get("query") != "activo:true,tipo_seguimiento_id:ts-1" || consulta.Get("fields") != "fecha_inicio,fecha_fin,unidades_interes" ||
			consulta.Get("sortby") != "fecha_modificacion" || consulta.Get("order") != "desc" {
			t.Errorf("consulta inesperada: %s", r.URL.String())
		}
		if !exito {
			json.NewEncoder(w).Encode(map[string]interface{}{"Success": false, "Status": "400", "Message": "falla del crud", "Data": nil})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"Success": true, "Status": "200", "Message": "Request succesfull", "Data": []interface{}{periodo}})
	}))
	defer servidor.Close()
	cliente := NuevoClienteCrud(servidor.URL)

	periodos, err := cliente.PeriodosActivos("ts-1")
	if err != nil || len(periodos) != 1 || periodos[0]["unidades_interes"] != periodo["unidades_interes"] {
		t.Fatalf("periodos = %v, err = %v", periodos, err)
	}
	exito = false
	if _, err := cliente.PeriodosActivos("ts-1"); err == nil {
		t.Fatal("con Success: false debe responder error")
	}
}
