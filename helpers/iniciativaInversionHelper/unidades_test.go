package iniciativaInversionHelper

import (
	"errors"
	"net/http"
	"testing"
)

// fuenteUnidadesFija imita a Terceros y Oikos en memoria.
type fuenteUnidadesFija struct {
	terceros      map[string]int
	vinculaciones map[int][]VinculacionTercero
	nombres       map[int]string
	err           error
}

func (f fuenteUnidadesFija) TerceroDe(documento string) (int, error) {
	return f.terceros[documento], f.err
}

func (f fuenteUnidadesFija) VinculacionesActivas(terceroID int) ([]VinculacionTercero, error) {
	return f.vinculaciones[terceroID], nil
}

func (f fuenteUnidadesFija) NombreDependencia(id int) (string, error) {
	return f.nombres[id], nil
}

func TestUnidadesSoloConLosCargosDeFormulacionYSinRepetir(t *testing.T) {
	fuente := fuenteUnidadesFija{
		terceros: map[string]int{"1020": 7},
		vinculaciones: map[int][]VinculacionTercero{7: {
			{CargoID: 320, DependenciaID: 41}, // asistente de dependencia
			{CargoID: 999, DependenciaID: 50}, // otro cargo: no da unidad
			{CargoID: 312, DependenciaID: 41}, // misma unidad con otro cargo
			{CargoID: 300, DependenciaID: 0},  // sin dependencia
			{CargoID: 312, DependenciaID: 60}, // jefe de oficina
		}},
		nombres: map[int]string{41: "BIENESTAR", 60: "PLANEACIÓN"},
	}

	unidades, err := NuevasUnidades(fuente).De(" 1020 ")
	if err != nil {
		t.Fatal(err)
	}
	if len(unidades) != 2 || unidades[0] != (Unidad{ID: 41, Nombre: "BIENESTAR"}) || unidades[1] != (Unidad{ID: 60, Nombre: "PLANEACIÓN"}) {
		t.Fatalf("unidades = %+v", unidades)
	}
}

func TestSinTerceroNoHayUnidades(t *testing.T) {
	unidades, err := NuevasUnidades(fuenteUnidadesFija{}).De("123")
	if err != nil || unidades == nil || len(unidades) != 0 {
		t.Fatalf("unidades = %+v, err = %v", unidades, err)
	}
}

func TestUnidadesExigeUnDocumentoValido(t *testing.T) {
	for _, documento := range []string{"  ", "123,Activo:false", "12 34", "123456789012345678901"} {
		_, err := NuevasUnidades(fuenteUnidadesFija{}).De(documento)
		exigirStatus(t, err, http.StatusBadRequest)
	}
}

func TestUnidadesPropagaElErrorDeTerceros(t *testing.T) {
	_, err := NuevasUnidades(fuenteUnidadesFija{err: errorServicio("Terceros no responde")}).De("123")
	var negocio *ErrorNegocio
	if !errors.As(err, &negocio) || negocio.Status != http.StatusBadGateway {
		t.Fatalf("err = %v", err)
	}
}
