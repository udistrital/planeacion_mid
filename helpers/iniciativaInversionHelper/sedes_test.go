package iniciativaInversionHelper

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// oikosDePrueba responde los espacios físicos de tipo SEDE como Oikos (incluye una inactiva y
// dos sedes con el mismo Codigo).
func oikosDePrueba(t *testing.T, estado int) *httptest.Server {
	t.Helper()
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/espacio_fisico" || r.URL.Query().Get("query") != "tipo_espacio_fisico_id:38" {
			t.Errorf("consulta inesperada a Oikos: %s", r.URL.String())
		}
		w.WriteHeader(estado)
		w.Write([]byte(`[
			{"Id":13,"Estado":"Activo","TipoEspacio":{"Id":38,"Nombre":"SEDE"},"Nombre":"TEATRO CAMARIN DEL CARMEN","Codigo":"FATC"},
			{"Id":12,"Estado":"Inactivo","TipoEspacio":{"Id":38,"Nombre":"SEDE"},"Nombre":"CENTRO","Codigo":"FC"},
			{"Id":10,"Estado":"Activo","TipoEspacio":{"Id":38,"Nombre":"SEDE"},"Nombre":"ACADEMICA  LUIS A. CALVO","Codigo":"FALC"},
			{"Id":18,"Estado":"Activo","TipoEspacio":{"Id":38,"Nombre":"SEDE"},"Nombre":"CALLE 58","Codigo":"FC"}
		]`))
	}))
	t.Cleanup(servidor.Close)
	return servidor
}

func hostDe(servidor *httptest.Server) string {
	return servidor.URL + "/v1"
}

func TestLasSedesSalenDeOikosActivasPorNombreYConOtraAlFinal(t *testing.T) {
	repo := NuevoRepositorioConSedes(&repoCatalogosContado{}, NuevoClienteSedesOikos(hostDe(oikosDePrueba(t, http.StatusOK))))

	sedes, err := repo.ElementosDeCatalogo(CatalogoSedes)
	if err != nil {
		t.Fatal(err)
	}
	esperadas := []ItemCatalogo{
		{Codigo: "10", Nombre: "ACADEMICA LUIS A. CALVO"},
		{Codigo: "18", Nombre: "CALLE 58"},
		{Codigo: "13", Nombre: "TEATRO CAMARIN DEL CARMEN"},
		{Codigo: CodigoOtraSede, Nombre: CodigoOtraSede},
	}
	if len(sedes) != len(esperadas) {
		t.Fatalf("se esperaban %v y llegaron %v", esperadas, sedes)
	}
	for i := range esperadas {
		if sedes[i] != esperadas[i] {
			t.Fatalf("sede %d: se esperaba %v y llegó %v", i, esperadas[i], sedes[i])
		}
	}
}

func TestLosDemasCatalogosSiguenSaliendoDelCrud(t *testing.T) {
	crud := &repoCatalogosContado{items: []ItemCatalogo{{Codigo: "2202001", Nombre: "Documentos de planeación"}}}
	repo := NuevoRepositorioConSedes(crud, NuevoClienteSedesOikos("host-que-no-se-usa"))

	items, err := repo.ElementosDeCatalogo("productos")
	if err != nil || len(items) != 1 || crud.consultas != 1 {
		t.Fatalf("debía consultar el crud: items %v, err %v, consultas %d", items, err, crud.consultas)
	}
}

func TestSiOikosFallaLasSedesNoSeInventan(t *testing.T) {
	repo := NuevoRepositorioConSedes(&repoCatalogosContado{}, NuevoClienteSedesOikos(hostDe(oikosDePrueba(t, http.StatusInternalServerError))))

	if sedes, err := repo.ElementosDeCatalogo(CatalogoSedes); err == nil {
		t.Fatalf("se esperaba un error y llegaron %v", sedes)
	}
}

type sedesFijas struct {
	items []ItemCatalogo
	err   error
}

func (s sedesFijas) SedesActivas() ([]ItemCatalogo, error) { return s.items, s.err }

func TestLaValidacionUsaLasSedesDeOikos(t *testing.T) {
	repo := NuevoRepositorioConSedes(&repoCatalogosContado{}, sedesFijas{items: []ItemCatalogo{{Codigo: "10", Nombre: "ACADEMICA LUIS A. CALVO"}}})
	validador := NuevoValidador(NuevosCatalogosConCache(repo, time.Hour))
	submodulo := &Nodo{Nombre: "Localización", Dato: Dato{Clave: "localizacion"}, Hijos: []*Nodo{
		{Nombre: "Sedes", Dato: Dato{Clave: "sedes"}, Hijos: []*Nodo{
			campoDePrueba("seleccionadas", TipoCatalogo, Dato{Requerido: true, Catalogo: CatalogoSedes, Multiple: true}),
			campoDePrueba("otra", TipoTexto, Dato{Requerido: true, SoloSi: &Condicion{Campo: "seleccionadas", Valor: CodigoOtraSede}}),
		}},
	}}
	validar := func(sedes ...interface{}) []ErrorCampo {
		return validador.ValidarSubmodulo(submodulo, map[string]Datos{"localizacion": {"sedes": Datos{"seleccionadas": sedes}}})
	}

	exigirSinErrores(t, "localizacion", validar("10"))
	exigirError(t, validar("10", CodigoOtraSede), "sedes.otra", "es obligatorio")
	exigirError(t, validar("99"), "sedes.seleccionadas", "ya no está en la lista")

	conError := NuevoRepositorioConSedes(&repoCatalogosContado{}, sedesFijas{err: errors.New("Oikos caído")})
	errores := NuevoValidador(NuevosCatalogosConCache(conError, time.Hour)).ValidarSubmodulo(submodulo,
		map[string]Datos{"localizacion": {"sedes": Datos{"seleccionadas": []interface{}{"10"}}}})
	exigirError(t, errores, "sedes.seleccionadas", "no se pudieron consultar")
}
