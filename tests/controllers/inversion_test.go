package controllers

import (
	"net/http"
	"testing"
)

func TestProyecto(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/proyecto/617c3171f6fc97294f27a041"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestProyecto Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestProyecto Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestProyecto:", err.Error())
		t.Fail()
	}
}
func TestConsultarMetasProyecto(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/metaspro/62266cfa16511ea35c5be15f"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarMetasProyecto Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarMetasProyecto Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarMetasProyecto:", err.Error())
		t.Fail()
	}
}
func TestProyectos(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/proyectos/string"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestProyectos Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestProyectos Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestProyectos:", err.Error())
		t.Fail()
	}
}
func TestConsultarPlanIdentificador(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/get_plan/62150f8525e40c9662093e58"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarPlanIdentificador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarPlanIdentificador Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarPlanIdentificador:", err.Error())
		t.Fail()
	}
}
func TestConsultarPlanInversion(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/get_infoPlan/63f467aeccee49e93a859ace"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarPlanInversion Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarPlanInversion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarPlanInversion:", err.Error())
		t.Fail()
	}
}
func TestConsultarTodasMetasPlan(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/all_metas/63f467aeccee49e93a859ace"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarTodasMetasPlan Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarTodasMetasPlan Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarTodasMetasPlan:", err.Error())
		t.Fail()
	}
}
func TestConsultarMagnitudesProgramadas(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/magnitudes/64367379a280497994a41e46/1"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarMagnitudesProgramadas Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarMagnitudesProgramadas Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarMagnitudesProgramadas:", err.Error())
		t.Fail()
	}
}
func TestVerificarMagnitudesProgramadas(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/verificar_magnitudes/64367379a280497994a41e46"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestVerificarMagnitudesProgramadas Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestVerificarMagnitudesProgramadas Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestVerificarMagnitudesProgramadas:", err.Error())
		t.Fail()
	}
}
