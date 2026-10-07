package reporteshelper

import "testing"

func TestFiltrarPlanesPorDependencias(t *testing.T) {
	planes := []map[string]interface{}{
		{"_id": "plan-1", "dependencia_id": "8"},
		{"_id": "plan-2", "dependencia_id": "14"},
		{"_id": "plan-3", "dependencia_id": "23"},
		{"_id": "plan-sin-dependencia"},
	}

	filtrados := FiltrarPlanesPorDependencias(planes, []string{"14", "23", "999"})

	if len(filtrados) != 2 {
		t.Fatalf("se esperaban 2 planes filtrados y se obtuvieron %d", len(filtrados))
	}
	if filtrados[0]["_id"] != "plan-2" || filtrados[1]["_id"] != "plan-3" {
		t.Fatalf("el filtro no conservó las dependencias seleccionadas: %#v", filtrados)
	}
}

func TestFiltrarPlanesPorDependenciasVacias(t *testing.T) {
	planes := []map[string]interface{}{{"_id": "plan-1", "dependencia_id": "8"}}

	filtrados := FiltrarPlanesPorDependencias(planes, []string{})

	if len(filtrados) != 0 {
		t.Fatalf("se esperaba un resultado vacío y se obtuvieron %d planes", len(filtrados))
	}
}
