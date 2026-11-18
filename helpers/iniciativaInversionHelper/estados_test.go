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
		{"enviar con todo formulado", AccionEnviar, EstadoIniciativaFormulacionIniciada, formuladas, EstadoIniciativaEnRevision, true},
		{"enviar con un submódulo sin formular", AccionEnviar, EstadoIniciativaFormulacionIniciada, []string{EstadoSubmoduloFormulado, EstadoSubmoduloSinFormular}, "", false},
		{"enviar con un submódulo formulando", AccionEnviar, EstadoIniciativaFormulacionIniciada, []string{EstadoSubmoduloFormulado, EstadoSubmoduloFormulando}, "", false},
		{"enviar sin haber empezado", AccionEnviar, EstadoIniciativaSinFormular, formuladas, "", false},
		{"reenviar tras corregir", AccionEnviar, EstadoIniciativaPorCorregir, []string{EstadoSubmoduloAprobado, EstadoSubmoduloFormulado}, EstadoIniciativaEnRevision, true},
		{"reenviar con observaciones pendientes", AccionEnviar, EstadoIniciativaPorCorregir, []string{EstadoSubmoduloConObservaciones}, "", false},
		{"enviar estando en revisión", AccionEnviar, EstadoIniciativaEnRevision, formuladas, "", false},
		{"devolver con observaciones", AccionDevolver, EstadoIniciativaEnRevision, []string{EstadoSubmoduloAprobado, EstadoSubmoduloConObservaciones}, EstadoIniciativaPorCorregir, true},
		{"devolver sin observaciones", AccionDevolver, EstadoIniciativaEnRevision, []string{EstadoSubmoduloAprobado}, "", false},
		{"avalar todo aprobado", AccionAvalar, EstadoIniciativaEnRevision, []string{EstadoSubmoduloAprobado, EstadoSubmoduloAprobado}, EstadoIniciativaAvalado, true},
		{"avalar con submódulos sin revisar", AccionAvalar, EstadoIniciativaEnRevision, []string{EstadoSubmoduloAprobado, EstadoSubmoduloFormulado}, "", false},
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
	enRevision := []string{EstadoSubmoduloAprobado, EstadoSubmoduloConObservaciones}
	if got := AccionesPosibles(EstadoIniciativaEnRevision, enRevision); !reflect.DeepEqual(got, []string{AccionDevolver}) {
		t.Fatalf("acciones = %v", got)
	}
	if got := AccionesPosibles(EstadoIniciativaAvalado, enRevision); len(got) != 0 {
		t.Fatalf("una iniciativa avalada no debe tener acciones: %v", got)
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
	if SubmoduloRevisable(EstadoSubmoduloFormulando) || SubmoduloRevisable(EstadoSubmoduloSinFormular) || !SubmoduloRevisable(EstadoSubmoduloFormulado) {
		t.Fatal("solo se revisa lo formulado")
	}
	if !PermiteEditar(EstadoIniciativaSinFormular) || !PermiteEditar(EstadoIniciativaFormulacionIniciada) || !PermiteEditar(EstadoIniciativaPorCorregir) ||
		PermiteEditar(EstadoIniciativaEnRevision) || !PermiteRevisar(EstadoIniciativaEnRevision) {
		t.Fatal("permisos por estado incorrectos")
	}
}

func TestLaIniciativaPasaAFormulacionIniciadaConLaPrimeraInformacion(t *testing.T) {
	if EstadoIniciativaTrasGuardar(EstadoIniciativaSinFormular, true) != EstadoIniciativaFormulacionIniciada {
		t.Fatal("con información, de Sin formular a Formulación iniciada")
	}
	for _, sinCambio := range []struct {
		estado         string
		conInformacion bool
	}{
		{EstadoIniciativaSinFormular, false},
		{EstadoIniciativaFormulacionIniciada, true},
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
