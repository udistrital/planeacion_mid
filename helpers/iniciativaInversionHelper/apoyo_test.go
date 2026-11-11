package iniciativaInversionHelper

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// Datos de prueba:
//   - testdata/identificacion.plantilla.json: copia de la construcción de
//     database local/construccion-proyecto-inversion/identificacion.plantilla.json
//   - testdata/iniciativa_ejemplo.json: la iniciativa de ejemplo del Excel (OATI), sin
//     errores. Las matrices de población van en 0 porque el Excel no trae cifras.

const tipoPlanDePrueba = "tipo-iniciativa"

func definicionDePrueba(t *testing.T) *Definicion {
	t.Helper()
	contenido, err := os.ReadFile("testdata/identificacion.plantilla.json")
	if err != nil {
		t.Fatal(err)
	}
	var plantilla struct {
		Plan struct {
			Nombre string `json:"nombre"`
		} `json:"plan"`
		Nodos []*Nodo `json:"nodos"`
	}
	if err := json.Unmarshal(contenido, &plantilla); err != nil {
		t.Fatal(err)
	}
	return &Definicion{PlantillaID: "plantilla-1", Nombre: plantilla.Plan.Nombre, TipoPlanID: tipoPlanDePrueba, Modulos: plantilla.Nodos}
}

// iniciativaDeEjemplo devuelve una copia nueva de la iniciativa de ejemplo en cada llamada.
func iniciativaDeEjemplo(t *testing.T) map[string]Datos {
	t.Helper()
	contenido, err := os.ReadFile("testdata/iniciativa_ejemplo.json")
	if err != nil {
		t.Fatal(err)
	}
	var iniciativa map[string]Datos
	if err := json.Unmarshal(contenido, &iniciativa); err != nil {
		t.Fatal(err)
	}
	return iniciativa
}

// validarIniciativa normaliza y valida todos los submódulos, como al guardarlos.
func validarIniciativa(t *testing.T, definicion *Definicion, iniciativa map[string]Datos) map[string][]ErrorCampo {
	t.Helper()
	catalogos := catalogosDePrueba()
	normalizador := NuevoNormalizador(catalogos)
	normalizado := map[string]Datos{}
	for _, submodulo := range definicion.Submodulos() {
		normalizado[submodulo.Clave()] = normalizador.Normalizar(submodulo, iniciativa[submodulo.Clave()])
	}
	validador := NuevoValidador(catalogos)
	errores := map[string][]ErrorCampo{}
	for _, submodulo := range definicion.Submodulos() {
		errores[submodulo.Clave()] = validador.ValidarSubmodulo(submodulo, normalizado)
	}
	return errores
}

func filasEn(t *testing.T, datos Datos, clave string) []interface{} {
	t.Helper()
	filas, ok := datos[clave].([]interface{})
	if !ok {
		t.Fatalf("%s no es una lista", clave)
	}
	return filas
}

func filaConID(t *testing.T, filas []interface{}, id string) Datos {
	t.Helper()
	for _, elemento := range filas {
		if fila := elemento.(map[string]interface{}); fila["id"] == id {
			return fila
		}
	}
	t.Fatalf("no existe la fila %s", id)
	return nil
}

// submodulosDe aplana los submódulos de todos los módulos de la vista general.
func submodulosDe(resumen *ResumenIniciativa) []ResumenSubmodulo {
	var submodulos []ResumenSubmodulo
	for _, modulo := range resumen.Modulos {
		submodulos = append(submodulos, modulo.Submodulos...)
	}
	return submodulos
}

func exigirSinErrores(t *testing.T, submodulo string, errores []ErrorCampo) {
	t.Helper()
	if len(errores) > 0 {
		t.Fatalf("%s: se esperaban 0 errores y hay %d: %+v", submodulo, len(errores), errores)
	}
}

func exigirError(t *testing.T, errores []ErrorCampo, ruta, fragmento string) {
	t.Helper()
	for _, e := range errores {
		if e.Ruta == ruta && strings.Contains(e.Mensaje, fragmento) {
			return
		}
	}
	t.Fatalf("no se encontró un error en %q que contenga %q; errores: %+v", ruta, fragmento, errores)
}
