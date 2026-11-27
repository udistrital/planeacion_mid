package iniciativaInversionHelper

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

// Con varios roles del mismo carril se actúa con el de mayor jerarquía; sin ninguno, no se actúa.
func TestElRolDelCarrilSigueLaJerarquia(t *testing.T) {
	casos := []struct {
		carril Carril
		roles  []string
		rol    string
	}{
		{CarrilFormulador, []string{RolAsistenteDependencia, RolJefeDependencia}, RolJefeDependencia},
		{CarrilFormulador, []string{RolJefeDependencia, RolAsistentePlaneacion}, RolAsistentePlaneacion},
		{CarrilFormulador, []string{RolPlaneacion}, ""},
		{CarrilJefe, []string{RolAsistenteDependencia, RolAsistentePlaneacion}, ""},
		{CarrilOAP, []string{RolAsistentePlaneacion, RolPlaneacion}, RolPlaneacion},
		{CarrilOAP, []string{"ADMINISTRADOR"}, ""},
	}
	for _, caso := range casos {
		if rol := caso.carril.rolEn(caso.roles); rol != caso.rol {
			t.Errorf("%s con %v: rol = %q, se esperaba %q", caso.carril.nombre, caso.roles, rol, caso.rol)
		}
	}
}

// La acción de cada submódulo sale de los permisos y del estado del submódulo.
func TestAccionDeSubmodulo(t *testing.T) {
	formula, revisa, nada := Permisos{Editar: true}, Permisos{Observar: true}, Permisos{}
	casos := []struct {
		permisos Permisos
		estado   string
		accion   string
	}{
		{formula, EstadoSubmoduloFormulando, AccionSubmoduloFormular},
		{formula, EstadoSubmoduloConObservaciones, AccionSubmoduloFormular},
		{formula, EstadoSubmoduloAprobado, AccionSubmoduloConsultar},
		{revisa, EstadoSubmoduloFormulado, AccionSubmoduloRevisar},
		{revisa, EstadoSubmoduloAprobado, AccionSubmoduloRevisar},
		{revisa, EstadoSubmoduloFormulando, AccionSubmoduloConsultar},
		{nada, EstadoSubmoduloFormulado, AccionSubmoduloConsultar},
	}
	for _, caso := range casos {
		if accion := accionDeSubmodulo(caso.permisos, caso.estado); accion != caso.accion {
			t.Errorf("%+v en %s: %s, se esperaba %s", caso.permisos, caso.estado, accion, caso.accion)
		}
	}
}

// La caché responde lo guardado hasta que vence, y no guarda las fallas.
func TestCacheDependencias(t *testing.T) {
	fuente := &unidadesFalsas{porUsuario: map[string][]string{"1020": {"41"}}}
	cache := NuevaCacheDependencias(fuente, time.Minute)
	ahora := time.Date(2026, 11, 26, 8, 0, 0, 0, time.UTC)
	cache.ahora = func() time.Time { return ahora }

	for i := 0; i < 2; i++ {
		if ids, err := cache.DependenciasDe("1020"); err != nil || !reflect.DeepEqual(ids, []string{"41"}) {
			t.Fatalf("ids = %v, err = %v", ids, err)
		}
	}
	if fuente.consultas != 1 {
		t.Fatalf("consultas = %d, se esperaba 1", fuente.consultas)
	}
	ahora = ahora.Add(time.Minute)
	fuente.err = errorServicio("Terceros no respondió")
	if _, err := cache.DependenciasDe("1020"); err == nil {
		t.Fatal("vencida, se consulta la fuente y responde su error")
	}
	fuente.err = nil
	if _, err := cache.DependenciasDe("1020"); err != nil || fuente.consultas != 3 {
		t.Fatalf("la falla no se guarda: err = %v, consultas = %d", err, fuente.consultas)
	}
}

// Sin fuente de unidades, lo que las necesita responde un error de configuración.
func TestSinFuenteDeUnidadesFallaComoConfiguracion(t *testing.T) {
	servicio, _ := nuevoServicioDePrueba(t)
	iniciativa := crearIniciativaDePrueba(t, servicio)
	servicio.unidades = nil

	_, err := servicio.Resumen(iniciativa.ID, formulador())
	exigirStatus(t, err, http.StatusInternalServerError)
}
