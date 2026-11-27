package iniciativaInversionHelper

import (
	"testing"
)

func TestEvaluarAccion(t *testing.T) {
	formuladas := []string{EstadoSubmoduloFormulado, EstadoSubmoduloFormulado}
	casos := []struct {
		nombre     string
		accion     string
		iniciativa string
		submodulos []string
		destino    string
		permitida  bool
	}{
		{"enviar con todo formulado", AccionEnviar, EstadoIniciativaEnFormulacion, formuladas, EstadoIniciativaPorAprobar, true},
		{"enviar con un submódulo sin formular", AccionEnviar, EstadoIniciativaEnFormulacion, []string{EstadoSubmoduloFormulado, EstadoSubmoduloSinFormular}, "", false},
		{"enviar con un submódulo formulando", AccionEnviar, EstadoIniciativaEnFormulacion, []string{EstadoSubmoduloFormulado, EstadoSubmoduloFormulando}, "", false},
		{"enviar sin haber empezado", AccionEnviar, EstadoIniciativaSinFormular, formuladas, "", false},
		{"reenviar tras corregir", AccionEnviar, EstadoIniciativaPorCorregir, []string{EstadoSubmoduloAprobado, EstadoSubmoduloFormulado}, EstadoIniciativaPorAprobar, true},
		{"reenviar con observaciones pendientes", AccionEnviar, EstadoIniciativaPorCorregir, []string{EstadoSubmoduloConObservaciones}, "", false},
		{"enviar estando por aprobar", AccionEnviar, EstadoIniciativaPorAprobar, formuladas, "", false},
		{"el Jefe aprueba sin observaciones", AccionAprobar, EstadoIniciativaPorAprobar, []string{EstadoSubmoduloAprobado, EstadoSubmoduloFormulado}, EstadoIniciativaPorRevisar, true},
		{"el Jefe no aprueba con observaciones", AccionAprobar, EstadoIniciativaPorAprobar, []string{EstadoSubmoduloConObservaciones, EstadoSubmoduloFormulado}, "", false},
		{"el Jefe devuelve con observaciones", AccionDevolver, EstadoIniciativaPorAprobar, []string{EstadoSubmoduloAprobado, EstadoSubmoduloConObservaciones}, EstadoIniciativaPorCorregir, true},
		{"el Jefe no devuelve sin observaciones", AccionDevolver, EstadoIniciativaPorAprobar, formuladas, "", false},
		{"devolver en revisión OAP", AccionDevolver, EstadoIniciativaEnRevision, []string{EstadoSubmoduloConObservaciones}, "", false},
		{"iniciar la revisión", AccionIniciarRevision, EstadoIniciativaPorRevisar, formuladas, EstadoIniciativaEnRevision, true},
		{"iniciar la revisión dos veces", AccionIniciarRevision, EstadoIniciativaEnRevision, formuladas, "", false},
		{"rechazar en revisión", AccionRechazar, EstadoIniciativaEnRevision, formuladas, EstadoIniciativaRechazada, true},
		{"rechazar antes de revisar", AccionRechazar, EstadoIniciativaPorRevisar, formuladas, "", false},
		{"retractar el rechazo", AccionRetractarRechazo, EstadoIniciativaRechazada, formuladas, EstadoIniciativaPorCorregir, true},
		{"retractar sin rechazo", AccionRetractarRechazo, EstadoIniciativaEnRevision, formuladas, "", false},
		{"nada sale de Aprobado", AccionRetractarRechazo, EstadoIniciativaAprobada, formuladas, "", false},
		{"acción inexistente", "borrar", EstadoIniciativaEnRevision, formuladas, "", false},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			destino, motivo := EvaluarAccion(caso.accion, caso.iniciativa, caso.submodulos)
			if (motivo == "") != caso.permitida || destino != caso.destino {
				t.Fatalf("destino=%q motivo=%q; se esperaba destino=%q permitida=%v", destino, motivo, caso.destino, caso.permitida)
			}
		})
	}
}

func TestEstadoTrasRevision(t *testing.T) {
	casos := []struct {
		nombre     string
		submodulos []string
		destino    string
	}{
		{"todos aprobados", []string{EstadoSubmoduloAprobado, EstadoSubmoduloAprobado}, EstadoIniciativaAprobada},
		{"todos con veredicto y uno con observaciones", []string{EstadoSubmoduloAprobado, EstadoSubmoduloConObservaciones}, EstadoIniciativaPorCorregir},
		{"falta un veredicto", []string{EstadoSubmoduloAprobado, EstadoSubmoduloFormulado}, ""},
		{"falta un veredicto aunque haya observaciones", []string{EstadoSubmoduloConObservaciones, EstadoSubmoduloFormulado}, ""},
		{"sin submódulos", nil, ""},
	}
	for _, caso := range casos {
		if destino := EstadoTrasRevision(caso.submodulos); destino != caso.destino {
			t.Errorf("%s: destino = %q, se esperaba %q", caso.nombre, destino, caso.destino)
		}
	}
}

func TestEstadosDeSubmodulo(t *testing.T) {
	conErrores := []ErrorCampo{{}}
	if EstadoTrasGuardar(nil, true) != EstadoSubmoduloFormulado || EstadoTrasGuardar(nil, false) != EstadoSubmoduloFormulado {
		t.Fatal("sin errores, el submódulo queda Formulado")
	}
	if EstadoTrasGuardar(conErrores, true) != EstadoSubmoduloFormulando || EstadoTrasGuardar(conErrores, false) != EstadoSubmoduloSinFormular {
		t.Fatal("con errores, Formulando si tiene información y Sin formular si está vacío")
	}
	if SubmoduloRevisable(EstadoSubmoduloFormulando) || SubmoduloRevisable(EstadoSubmoduloSinFormular) ||
		!SubmoduloRevisable(EstadoSubmoduloFormulado) || !SubmoduloRevisable(EstadoSubmoduloAprobado) || !SubmoduloRevisable(EstadoSubmoduloConObservaciones) {
		t.Fatal("se revisa lo formulado: Formulado, Con observaciones o Aprobado")
	}
}

func TestLaIniciativaPasaAEnFormulacionConLaPrimeraInformacion(t *testing.T) {
	if EstadoIniciativaTrasGuardar(EstadoIniciativaSinFormular, true) != EstadoIniciativaEnFormulacion {
		t.Fatal("con información, de Sin formular a En formulación")
	}
	for _, sinCambio := range []struct {
		estado         string
		conInformacion bool
	}{
		{EstadoIniciativaSinFormular, false},
		{EstadoIniciativaEnFormulacion, true},
		{EstadoIniciativaPorCorregir, true},
	} {
		if destino := EstadoIniciativaTrasGuardar(sinCambio.estado, sinCambio.conInformacion); destino != "" {
			t.Fatalf("%s con información=%v no debe cambiar y cambió a %s", sinCambio.estado, sinCambio.conInformacion, destino)
		}
	}
}

func TestTieneInformacionNoCuentaLoQuePoneElSistema(t *testing.T) {
	nodos := []*Nodo{
		campoDePrueba("nombre", TipoCalculado, Dato{Plantilla: "{x}"}),
		campoDePrueba("filas", TipoRepetible, Dato{}, campoDePrueba("texto", TipoTexto, Dato{}), campoDePrueba("marcada", TipoCasilla, Dato{})),
		{Nombre: "grupo", Dato: Dato{Clave: "grupo"}, Hijos: []*Nodo{campoDePrueba("numero", TipoNumerico, Dato{})}},
	}
	vacios := []Datos{
		{},
		{"nombre": "PROPUESTO", "filas": []interface{}{map[string]interface{}{"id": "f-1", "origen": "cd-1", "texto": "  ", "marcada": false}}},
		{"grupo": map[string]interface{}{"numero": nil}},
	}
	for _, datos := range vacios {
		if tieneInformacion(nodos, datos) {
			t.Fatalf("no debía contar como información: %v", datos)
		}
	}
	llenos := []Datos{
		{"filas": []interface{}{map[string]interface{}{"id": "f-1", "texto": "algo"}}},
		{"filas": []interface{}{map[string]interface{}{"id": "f-1", "marcada": true}}},
		{"grupo": map[string]interface{}{"numero": 0.0}},
	}
	for _, datos := range llenos {
		if !tieneInformacion(nodos, datos) {
			t.Fatalf("debía contar como información: %v", datos)
		}
	}
}
