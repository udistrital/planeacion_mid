package controllers

import (
	"net/http"
	"testing"
)

func TestTrimestres(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/trimestres/35"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestTrimestres Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestTrimestres Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestTrimestres:", err.Error())
		t.Fail()
	}
}
func TestPeriodos(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_periodos/35"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPeriodos Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPeriodos Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPeriodos:", err.Error())
		t.Fail()
	}
}
func TestActividades(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_actividades/61f60e4525e40c6f5d084185"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestActividades Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestActividades Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestActividades:", err.Error())
		t.Fail()
	}
}
func TestSeguimiento(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_seguimiento/61f08edc25e40c91b0083e4f/1/635b1f995073f2675157dc7f"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestSeguimiento Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestSeguimiento:", err.Error())
		t.Fail()
	}
}

func TestIndicadores(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_indicadores/6201d43f25e40c205608b459"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestIndicadores Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestIndicadores Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestIndicadores:", err.Error())
		t.Fail()
	}
}

func TestConsultarEstadoTrimestre(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_estado_trimestre/666a3b79252f5d633cc097de/T2"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarEstadoTrimestre Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarEstadoTrimestre Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarEstadoTrimestre:", err.Error())
		t.Fail()
	}
}

func TestConsultarEstadoTrimestres(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/estado_trimestres/666a3b79252f5d633cc097de"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarEstadoTrimestres Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarEstadoTrimestres Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarEstadoTrimestres:", err.Error())
		t.Fail()
	}
}
