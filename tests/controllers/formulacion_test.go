package controllers

import (
	"net/http"
	"testing"
)

func TestGetPlan(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/get_plan/666a3b79252f5d633cc097de/4"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetPlan Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetPlan Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetPlan:", err.Error())
		t.Fail()
	}
}
func TestGetAllActividades(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/get_all_actividades/666a3b79252f5d633cc097de"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetAllActividades Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetAllActividades Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetAllActividades:", err.Error())
		t.Fail()
	}
}

func TestConsultarIdentificaciones(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/get_all_identificacion/6675d1b1252f5d60c7c6f56c/6184b3e6f6fc97850127bb68"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarIdentificaciones Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarIdentificaciones Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarIdentificaciones:", err.Error())
		t.Fail()
	}
}

func TestConsultarPlanVersiones(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/get_plan_versiones/pruebas/2020/Plan de acción Inversión 20202"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarPlanVersiones Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarPlanVersiones Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarPlanVersiones:", err.Error())
		t.Fail()
	}
}

func TestPonderacionActividades(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/ponderacion_actividades/64368f5aa280496519a44efc"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPonderacionActividades Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestCPonderacionActividades Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPonderacionActividades:", err.Error())
		t.Fail()
	}
}

func TestConsultarUnidades(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/get_unidades"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarUnidades Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarUnidades Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarUnidades:", err.Error())
		t.Fail()
	}
}

func TestVinculacionTercero(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/vinculacion_tercero/59769"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestVinculacionTercero Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestVinculacionTercero Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestVinculacionTercero:", err.Error())
		t.Fail()
	}
}

// TODO: VinculacionTerceroByEmail
// TODO: VinculacionTerceroByIdentificacion
func TestPlanes(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/planes"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPlanes Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPlanes Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPlanes:", err.Error())
		t.Fail()
	}
}

func TestVerificarIdentificaciones(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/verificar_identificaciones/59769"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestVerificarIdentificaciones Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestVerificarIdentificaciones Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestVerificarIdentificaciones:", err.Error())
		t.Fail()
	}
}
func TestPlanesEnFormulacion(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/planes_formulacion/"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPlanesEnFormulacion Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPlanesEnFormulacion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPlanesEnFormulacion:", err.Error())
		t.Fail()
	}
}
