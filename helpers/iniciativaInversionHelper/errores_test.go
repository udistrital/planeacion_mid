package iniciativaInversionHelper

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

// Los mensajes llegan tal cual a la pantalla del usuario; la causa técnica
// (scripts de carga, servicios, rutas, códigos) solo puede ir en Detalle.

// fragmentosTecnicos no deben aparecer en un mensaje para el usuario.
var fragmentosTecnicos = []string{".js", "database local", "planes_crud", "IPI_SP", "_SUB", "estado-submodulo", "tipo-plan", "HTTP", "JSON", "dial tcp", "Cast to"}

func exigirMensajeParaUsuario(t *testing.T, err error, detalleEsperado string) *ErrorNegocio {
	t.Helper()
	var negocio *ErrorNegocio
	if !errors.As(err, &negocio) {
		t.Fatalf("se esperaba un ErrorNegocio y llegó %v", err)
	}
	for _, fragmento := range fragmentosTecnicos {
		if strings.Contains(negocio.Mensaje, fragmento) {
			t.Errorf("el mensaje para el usuario contiene %q: %s", fragmento, negocio.Mensaje)
		}
	}
	if !strings.Contains(negocio.Detalle, detalleEsperado) {
		t.Errorf("el detalle para el log debería contener %q: %s", detalleEsperado, negocio.Detalle)
	}
	return negocio
}

func TestSinTipoDePlanElUsuarioNoVeElScriptDeCarga(t *testing.T) {
	repo := nuevoRepoMemoria()
	delete(repo.catalogos, RecursoTipoPlan)
	servicio := NuevoServicio(repo, definicionFija{definicionDePrueba(t)}, catalogosDePrueba())

	_, err := servicio.ListarIniciativas(nil)
	negocio := exigirMensajeParaUsuario(t, err, "cargar-estados.js")
	if negocio.Mensaje != mensajeConfiguracion {
		t.Fatalf("mensaje = %q", negocio.Mensaje)
	}
}

func TestSinEstadoDeSubmoduloElUsuarioNoVeElScriptDeCarga(t *testing.T) {
	repo := nuevoRepoMemoria()
	repo.catalogos[RecursoEstadoSubmodulo] = repo.catalogos[RecursoEstadoSubmodulo][1:]
	servicio := NuevoServicio(repo, definicionFija{definicionDePrueba(t)}, catalogosDePrueba())

	_, err := servicio.ListarIniciativas(nil)
	exigirMensajeParaUsuario(t, err, "cargar-estados.js")
}

func TestSinPlantillaElUsuarioNoVeElScriptDeCarga(t *testing.T) {
	servicio := NuevoServicio(nuevoRepoMemoria(), definicionesPorPlantilla{}, catalogosDePrueba())
	_, err := servicio.CrearIniciativa(SolicitudCrearIniciativa{DependenciaID: "dep-1"})
	exigirStatus(t, err, http.StatusConflict)
	exigirMensajeParaUsuario(t, err, "cargar-construccion.js")
}

func TestCatalogoVacioElUsuarioNoVeElScriptDeCarga(t *testing.T) {
	_, err := NuevosCatalogosConCache(&repoCatalogosContado{}, time.Hour).Catalogo("sectores_pnd")
	exigirStatus(t, err, http.StatusNotFound)
	exigirMensajeParaUsuario(t, err, "cargar-catalogos.js")
}

func TestLosMensajesTecnicosDelCrudVanSoloAlLog(t *testing.T) {
	casos := map[int]string{
		http.StatusBadRequest:          mensajeSolicitudInvalida,
		http.StatusNotFound:            "No se encontró la información solicitada.",
		http.StatusInternalServerError: mensajeServicio,
	}
	for status, mensaje := range casos {
		err := errorDesdeCrud(status, "Cast to ObjectId failed for value \"x\"", "GET", "/plan/x")
		negocio := exigirMensajeParaUsuario(t, err, "Cast to ObjectId")
		if negocio.Mensaje != mensaje {
			t.Errorf("%d: mensaje = %q", status, negocio.Mensaje)
		}
	}
}

func TestElConflictoDeVersionDelCrudSeMuestraTalCual(t *testing.T) {
	mensaje := "El submódulo cambió desde que se consultó. Vuelva a cargarlo antes de guardar."
	negocio := errorDesdeCrud(http.StatusConflict, mensaje, "PUT", "/submodulo-iniciativa/x")
	if negocio.Status != http.StatusConflict || negocio.Mensaje != mensaje {
		t.Fatalf("error = %+v", negocio)
	}
}

func TestUnErrorDesconocidoNoMuestraSuTexto(t *testing.T) {
	err := ComoErrorNegocio(errors.New("dial tcp 10.0.0.1:8080: connection refused"))
	negocio := exigirMensajeParaUsuario(t, err, "dial tcp")
	if negocio.Status != http.StatusInternalServerError || negocio.Mensaje != mensajeErrorInterno {
		t.Fatalf("error = %+v", negocio)
	}
}

func TestLaAccionNoPermitidaSeNombraComoEnPantalla(t *testing.T) {
	_, motivo := EvaluarAccion(AccionNoAvalar, EstadoIniciativaEnFormulacion, nil)
	if !strings.Contains(motivo, "«No avalar»") {
		t.Fatalf("motivo = %q", motivo)
	}
}

func TestElCatalogoDependienteNombraElCampoYNoSuClave(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	iniciativa["articulacion_pnd"]["sector"] = ""

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["articulacion_pnd"], "programa", "primero seleccione «Sector»")
}
