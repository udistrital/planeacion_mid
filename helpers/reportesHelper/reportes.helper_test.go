package reporteshelper

import "testing"

func docenteParaTotales(cantidad float64) map[string]interface{} {
	return map[string]interface{}{
		"cantidad":         cantidad,
		"sueldoBasico":     "$1,000.00",
		"primaServicios":   "$100.00",
		"primaNavidad":     "$200.00",
		"primaVacaciones":  "$300.00",
		"bonificacion":     "$400.00",
		"intereses":        "$50.00",
		"cesantiasPublico": "$600.00",
		"cesantiasPrivado": "$700.00",
		"totalSalud":       "$800.00",
		"pensionesPublico": "$900.00",
		"pensionesPrivado": "$1,000.00",
		"totalArl":         "$1,100.00",
		"caja":             "$120.00",
		"icbf":             "$130.00",
	}
}

func TestTotalDocentesDistingueValoresTotalesEIndividuales(t *testing.T) {
	esperados := map[string]int{
		"sueldoBasico":      1000,
		"primaServicios":    200,
		"primaNavidad":      400,
		"primaVacaciones":   600,
		"bonificacion":      800,
		"cesantias":         100,
		"cesantiasPublicas": 600,
		"cesantiasPrivadas": 700,
		"salud":             800,
		"pensionesPublicas": 900,
		"pensionesPrivadas": 1000,
		"arl":               1100,
		"caja":              240,
		"icbf":              260,
	}

	for _, grupo := range []string{"rhf", "rhv_pre", "rhv_pos"} {
		t.Run(grupo, func(t *testing.T) {
			activo := docenteParaTotales(2)
			activo["activo"] = true
			inactivo := docenteParaTotales(100)
			inactivo["activo"] = false
			docentes := map[string]interface{}{
				"rhf":     "{}",
				"rhv_pre": "{}",
				"rhv_pos": "{}",
			}
			docentes[grupo] = []map[string]interface{}{activo, inactivo}

			resultado := TotalDocentes(docentes)
			for campo, esperado := range esperados {
				if obtenido := resultado[campo].(int); obtenido != esperado {
					t.Errorf("%s: se esperaba %d y se obtuvo %d", campo, esperado, obtenido)
				}
			}

			totales := resultado["TotalesPorTipo"].(TotalesDocentes)
			var totalGrupo TotalDocentVal
			switch grupo {
			case "rhf":
				totalGrupo = totales.Rhf
			case "rhv_pre":
				totalGrupo = totales.Rhv_pre
			case "rhv_pos":
				totalGrupo = totales.Rhv_pos
			}
			if totalGrupo.SalarioBasico != 1000 || totalGrupo.PrimaServicios != 200 || totalGrupo.InteresesCesantias != 100 {
				t.Errorf("los totales de %s no respetan la misma regla: %#v", grupo, totalGrupo)
			}
		})
	}
}

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
