package iniciativaInversionHelper

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"
)

// repoCatalogosArchivos lee los catálogos de testdata/catalogos/ (copias de los que
// genera database local/catalogos/extraer-catalogos.js).
type repoCatalogosArchivos struct{}

func (repoCatalogosArchivos) ElementosDeCatalogo(nombre string) ([]ItemCatalogo, error) {
	contenido, err := os.ReadFile("testdata/catalogos/" + nombre + ".json")
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var archivo struct {
		Items []ItemCatalogo `json:"items"`
	}
	if err := json.Unmarshal(contenido, &archivo); err != nil {
		return nil, err
	}
	return archivo.Items, nil
}

// catalogosDePrueba es el mismo proveedor de producción, sobre los archivos de prueba.
func catalogosDePrueba() ProveedorCatalogos {
	return NuevosCatalogosConCache(repoCatalogosArchivos{}, time.Hour)
}

// repoCatalogosContado cuenta las consultas y puede fallar.
type repoCatalogosContado struct {
	consultas int
	items     []ItemCatalogo
	err       error
}

func (r *repoCatalogosContado) ElementosDeCatalogo(string) ([]ItemCatalogo, error) {
	r.consultas++
	return r.items, r.err
}

func TestElCatalogoSeGuardaEnMemoriaHastaQueVence(t *testing.T) {
	repo := &repoCatalogosContado{items: []ItemCatalogo{{Codigo: "45", Nombre: "Educación"}}}
	proveedor := NuevosCatalogosConCache(repo, 10*time.Minute)
	reloj := time.Date(2026, 11, 4, 8, 0, 0, 0, time.UTC)
	proveedor.ahora = func() time.Time { return reloj }

	for i := 0; i < 3; i++ {
		if _, err := proveedor.Catalogo("sectores_pnd"); err != nil {
			t.Fatal(err)
		}
	}
	if repo.consultas != 1 {
		t.Fatalf("dentro de la vigencia se consulta una sola vez; hubo %d", repo.consultas)
	}

	reloj = reloj.Add(11 * time.Minute)
	catalogo, err := proveedor.Catalogo("sectores_pnd")
	if err != nil || repo.consultas != 2 {
		t.Fatalf("al vencer se vuelve a consultar: %d consultas, %v", repo.consultas, err)
	}
	if item, existe := catalogo.Buscar("45"); !existe || item.Nombre != "Educación" {
		t.Fatalf("Buscar(45) = %+v, %v", item, existe)
	}
}

func TestUnCatalogoVacioRespondeNoEncontradoYNoSeGuarda(t *testing.T) {
	repo := &repoCatalogosContado{}
	proveedor := NuevosCatalogosConCache(repo, time.Hour)

	_, err := proveedor.Catalogo("regiones")
	exigirStatus(t, err, http.StatusNotFound)

	repo.items = []ItemCatalogo{{Codigo: "Caribe", Nombre: "Caribe"}}
	if _, err := proveedor.Catalogo("regiones"); err != nil {
		t.Fatalf("en cuanto se carga el catálogo debe encontrarse: %v", err)
	}
}

func TestUnErrorDelCrudNoSeGuarda(t *testing.T) {
	repo := &repoCatalogosContado{err: errorServicio("crud caído")}
	proveedor := NuevosCatalogosConCache(repo, time.Hour)

	if _, err := proveedor.Catalogo("regiones"); err == nil {
		t.Fatal("se esperaba el error del crud")
	}
	repo.err, repo.items = nil, []ItemCatalogo{{Codigo: "Caribe", Nombre: "Caribe"}}
	if _, err := proveedor.Catalogo("regiones"); err != nil || repo.consultas != 2 {
		t.Fatalf("tras un error se vuelve a consultar: %d consultas, %v", repo.consultas, err)
	}
}

func TestNombreDeCatalogoInvalidoResponde400(t *testing.T) {
	repo := &repoCatalogosContado{}
	_, err := NuevosCatalogosConCache(repo, time.Hour).Catalogo("../secretos")
	exigirStatus(t, err, http.StatusBadRequest)
	if repo.consultas != 0 {
		t.Fatal("un nombre inválido no llega al crud")
	}
}

func TestFiltrarPorPadre(t *testing.T) {
	catalogo := NuevoCatalogo("programas_pnd", []ItemCatalogo{
		{Codigo: "4501", Nombre: "a", Padre: "45"}, {Codigo: "4502", Nombre: "b", Padre: "45"}, {Codigo: "2301", Nombre: "c", Padre: "23"},
	})
	if len(catalogo.Filtrar("45")) != 2 || len(catalogo.Filtrar("")) != 3 || len(catalogo.Filtrar("99")) != 0 {
		t.Fatalf("filtro por padre incorrecto")
	}
}
