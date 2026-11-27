package iniciativaInversionHelper

import (
	"errors"
	"net/http"
	"net/http/httptest"
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

// DependenciasDe da los mismos ids que De, como texto (igual que la dependencia_id de las iniciativas).
func TestDependenciasDeDaLosIdsComoTexto(t *testing.T) {
	fuente := fuenteUnidadesFija{
		terceros:      map[string]int{"1020": 7},
		vinculaciones: map[int][]VinculacionTercero{7: {{CargoID: 320, DependenciaID: 41}, {CargoID: 312, DependenciaID: 60}}},
	}
	ids, err := NuevasUnidades(fuente).DependenciasDe("1020")
	if err != nil || len(ids) != 2 || ids[0] != "41" || ids[1] != "60" {
		t.Fatalf("ids = %v, err = %v", ids, err)
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

// ClienteUnidades busca el tercero activo por documento: NombreTercero da su nombre y TerceroDe,
// su id. Un documento sin tercero ([{}]) no tiene nombre, y lo que no es un documento no se
// consulta.
func TestClienteUnidadesBuscaElTerceroActivoPorDocumento(t *testing.T) {
	consultas := 0
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consultas++
		if r.URL.Path != "/datos_identificacion" {
			t.Errorf("ruta inesperada: %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") == "Numero:52123456,Activo:true" {
			w.Write([]byte(`[{"Numero":"52123456","TerceroId":{"Id":7,"NombreCompleto":"LAURA GOMEZ RIOS"}}]`))
			return
		}
		w.Write([]byte(`[{}]`))
	}))
	defer servidor.Close()
	cliente := NuevoClienteUnidades(servidor.URL, "")

	if nombre, err := cliente.NombreTercero("52123456"); err != nil || nombre != "LAURA GOMEZ RIOS" {
		t.Fatalf("nombre = %q, err = %v", nombre, err)
	}
	if id, err := cliente.TerceroDe("52123456"); err != nil || id != 7 {
		t.Fatalf("id = %d, err = %v", id, err)
	}
	if _, err := cliente.NombreTercero("999"); err == nil {
		t.Fatal("un documento sin tercero no tiene nombre: se esperaba un error")
	}
	if id, err := cliente.TerceroDe("999"); err != nil || id != 0 {
		t.Fatalf("id = %d, err = %v", id, err)
	}
	antes := consultas
	if _, err := cliente.NombreTercero("123,Activo:false"); err == nil || consultas != antes {
		t.Fatalf("lo que no es un documento no se consulta: err = %v, consultas = %d", err, consultas-antes)
	}
}

func TestUnidadesPropagaElErrorDeTerceros(t *testing.T) {
	_, err := NuevasUnidades(fuenteUnidadesFija{err: errorServicio("Terceros no responde")}).De("123")
	var negocio *ErrorNegocio
	if !errors.As(err, &negocio) || negocio.Status != http.StatusBadGateway {
		t.Fatalf("err = %v", err)
	}
}
