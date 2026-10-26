package formulacionhelper

import (
	"math"
	"strconv"
	"testing"
)

func TestGetCajaCompensacionHCPUsaVacaciones(t *testing.T) {
	const (
		salarioBasico   = 64594.79582992248
		vacaciones      = 5382.680454567342
		primaVacaciones = 3588.8850098714306
		semanas         = 19.0
		horas           = 8.0
		incremento      = 0.093
	)

	data := map[string]interface{}{
		"semanas":       semanas,
		"horas":         horas,
		"incremento":    incremento,
		"salarioMinimo": 1_300_000.0,
		"resolucionDocente": map[string]interface{}{
			"Dedicacion":      HCatedraPrestacional,
			"salarioBasico":   salarioBasico,
			"vacaciones":      vacaciones,
			"primaVacaciones": primaVacaciones,
		},
	}

	obtenidoTexto := GetCajaCompensacion(data)
	obtenido, err := strconv.ParseFloat(obtenidoTexto, 64)
	if err != nil {
		t.Fatalf("GetCajaCompensacion retornó un valor inválido: %q", obtenidoTexto)
	}

	esperado := (salarioBasico + vacaciones) * horas * semanas * 0.04 * (1 + incremento)
	if math.Abs(obtenido-esperado) > 0.01 {
		t.Fatalf("CCF HCP incorrecta: se obtuvo %.2f y se esperaba %.2f", obtenido, esperado)
	}

	if math.Round(obtenido) != 465031 {
		t.Fatalf("CCF HCP inesperada para 19 semanas y 8 horas: se obtuvo %.2f", obtenido)
	}
}

func TestGetIcbfHCPUsaVacaciones(t *testing.T) {
	const (
		salarioBasico   = 64594.79582992248
		vacaciones      = 5382.680454567342
		primaVacaciones = 3588.8850098714306
		semanas         = 19.0
		horas           = 8.0
		incremento      = 0.093
	)

	data := map[string]interface{}{
		"semanas":       semanas,
		"horas":         horas,
		"incremento":    incremento,
		"salarioMinimo": 1_300_000.0,
		"resolucionDocente": map[string]interface{}{
			"Dedicacion":      HCatedraPrestacional,
			"salarioBasico":   salarioBasico,
			"vacaciones":      vacaciones,
			"primaVacaciones": primaVacaciones,
		},
	}

	obtenidoTexto := GetIcbf(data)
	obtenido, err := strconv.ParseFloat(obtenidoTexto, 64)
	if err != nil {
		t.Fatalf("GetIcbf retornó un valor inválido: %q", obtenidoTexto)
	}

	esperado := (salarioBasico + vacaciones) * horas * semanas * 0.03 * (1 + incremento)
	if math.Abs(obtenido-esperado) > 0.01 {
		t.Fatalf("ICBF HCP incorrecto: se obtuvo %.2f y se esperaba %.2f", obtenido, esperado)
	}

	if math.Round(obtenido) != 348773 {
		t.Fatalf("ICBF HCP inesperado para 19 semanas y 8 horas: se obtuvo %.2f", obtenido)
	}
}
