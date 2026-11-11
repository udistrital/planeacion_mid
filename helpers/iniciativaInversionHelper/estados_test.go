package iniciativaInversionHelper

import (
	"reflect"
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
		{"enviar con todo formulado", AccionEnviar, EstadoIniciativaEnFormulacion, formuladas, EstadoIniciativaEnRevision, true},
		{"enviar con un submódulo sin formular", AccionEnviar, EstadoIniciativaEnFormulacion, []string{EstadoSubmoduloFormulado, EstadoSubmoduloSinFormular}, "", false},
		{"reenviar tras corregir", AccionEnviar, EstadoIniciativaPorCorregir, []string{EstadoSubmoduloSinObservaciones, EstadoSubmoduloFormulado}, EstadoIniciativaEnRevision, true},
		{"reenviar con observaciones pendientes", AccionEnviar, EstadoIniciativaPorCorregir, []string{EstadoSubmoduloConObservaciones}, "", false},
		{"enviar estando en revisión", AccionEnviar, EstadoIniciativaEnRevision, formuladas, "", false},
		{"devolver con observaciones", AccionDevolver, EstadoIniciativaEnRevision, []string{EstadoSubmoduloSinObservaciones, EstadoSubmoduloConObservaciones}, EstadoIniciativaPorCorregir, true},
		{"devolver sin observaciones", AccionDevolver, EstadoIniciativaEnRevision, []string{EstadoSubmoduloSinObservaciones}, "", false},
		{"avalar todo sin observaciones", AccionAvalar, EstadoIniciativaEnRevision, []string{EstadoSubmoduloSinObservaciones, EstadoSubmoduloSinObservaciones}, EstadoIniciativaAvalado, true},
		{"avalar con submódulos sin revisar", AccionAvalar, EstadoIniciativaEnRevision, []string{EstadoSubmoduloSinObservaciones, EstadoSubmoduloFormulado}, "", false},
		{"no avalar en revisión", AccionNoAvalar, EstadoIniciativaEnRevision, formuladas, EstadoIniciativaNoAvalado, true},
		{"no avalar en formulación", AccionNoAvalar, EstadoIniciativaEnFormulacion, formuladas, "", false},
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

func TestAccionesPosibles(t *testing.T) {
	enRevision := []string{EstadoSubmoduloSinObservaciones, EstadoSubmoduloConObservaciones}
	if got := AccionesPosibles(EstadoIniciativaEnRevision, enRevision); !reflect.DeepEqual(got, []string{AccionDevolver, AccionNoAvalar}) {
		t.Fatalf("acciones = %v", got)
	}
	if got := AccionesPosibles(EstadoIniciativaAvalado, enRevision); len(got) != 0 {
		t.Fatalf("una iniciativa avalada no debe tener acciones: %v", got)
	}
}

func TestEstadosDeSubmodulo(t *testing.T) {
	if EstadoTrasGuardar(nil) != EstadoSubmoduloFormulado || EstadoTrasGuardar([]ErrorCampo{{}}) != EstadoSubmoduloSinFormular {
		t.Fatal("estado tras guardar incorrecto")
	}
	if !PermiteEditar(EstadoIniciativaPorCorregir) || PermiteEditar(EstadoIniciativaEnRevision) || !PermiteRevisar(EstadoIniciativaEnRevision) {
		t.Fatal("permisos por estado incorrectos")
	}
}
