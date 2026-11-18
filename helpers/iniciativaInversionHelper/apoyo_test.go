package iniciativaInversionHelper

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// Datos de prueba:
//   - testdata/iniciativa.plantilla.json: la plantilla de database local/construccion-proyecto-inversion/
//     iniciativa.plantilla.json con los módulos Identificación y Preparación. Anexos no va: sus
//     campos son documentos, un tipo que el validador no tiene, y ninguna iniciativa de prueba
//     quedaría completa para enviarla a revisión.
//   - testdata/iniciativa_ejemplo.json: database local/construccion-proyecto-inversion/
//     iniciativa.ejemplo.json con la sede de prueba «1» en Localización. Todos sus submódulos
//     quedan formulados.
//   - testdata/catalogos/: copias de database local/catalogos con los catálogos que usa la
//     plantilla, y sedes.json, de prueba (en producción las sedes salen de Oikos).

const tipoPlanDePrueba = "tipo-iniciativa"

func definicionDePrueba(t *testing.T) *Definicion {
	t.Helper()
	contenido, err := os.ReadFile("testdata/iniciativa.plantilla.json")
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

// problematicaConOtraCausa es la Problemática del ejemplo con una cuarta causa directa.
func problematicaConOtraCausa(t *testing.T) Datos {
	t.Helper()
	problematica := iniciativaDeEjemplo(t)["problematica"]
	problematica["causas_directas"] = append(filasEn(t, problematica, "causas_directas"), map[string]interface{}{
		"id": "cd-4", "descripcion": "Nueva causa directa.",
		"causas_indirectas": []interface{}{map[string]interface{}{"id": "ci-9", "descripcion": "Nueva causa indirecta."}},
	})
	return problematica
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
