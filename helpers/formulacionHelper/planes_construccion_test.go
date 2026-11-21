package formulacionhelper

import (
	"reflect"
	"testing"
)

func TestSinTipoPlanQuitaLosPlanesDeEseTipo(t *testing.T) {
	planes := []map[string]interface{}{
		{"_id": "p1", "tipo_plan_id": "paf"},
		{"_id": "p2", "tipo_plan_id": "ipi"},
		{"_id": "p3"},
		{"_id": "p4", "tipo_plan_id": "pi"},
	}

	obtenidos := SinTipoPlan(planes, "ipi")

	esperados := []map[string]interface{}{planes[0], planes[2], planes[3]}
	if !reflect.DeepEqual(obtenidos, esperados) {
		t.Fatalf("SinTipoPlan = %v; se esperaba %v", obtenidos, esperados)
	}
}

// Si el tipo no existe en el ambiente (id vacío), la lista sale completa en vez de fallar.
func TestSinTipoPlanSinTipoDevuelveTodos(t *testing.T) {
	planes := []map[string]interface{}{{"_id": "p1", "tipo_plan_id": "paf"}, {"_id": "p2", "tipo_plan_id": ""}}

	if obtenidos := SinTipoPlan(planes, ""); !reflect.DeepEqual(obtenidos, planes) {
		t.Fatalf("SinTipoPlan = %v; se esperaba la lista completa", obtenidos)
	}
}

func TestIDPrimerTipoPlan(t *testing.T) {
	casos := []struct {
		nombre   string
		tipos    []map[string]interface{}
		esperado string
	}{
		{"con resultado", []map[string]interface{}{{"_id": "ipi", "codigo_abreviacion": "IPI_SP"}}, "ipi"},
		{"sin resultados", nil, ""},
		{"lista vacía", []map[string]interface{}{}, ""},
		{"sin _id", []map[string]interface{}{{"codigo_abreviacion": "IPI_SP"}}, ""},
	}
	for _, caso := range casos {
		if obtenido := IDPrimerTipoPlan(caso.tipos); obtenido != caso.esperado {
			t.Errorf("%s: IDPrimerTipoPlan = %q; se esperaba %q", caso.nombre, obtenido, caso.esperado)
		}
	}
}
