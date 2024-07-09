package controllers

import (
	"net/http"
	"testing"
)

func TestGetPlanesPeriodo(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/evaluacion/planes_periodo/25/8"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetPlanesPeriodo Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetPlanesPeriodo Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetPlanesPeriodo:", err.Error())
		t.Fail()
	}
}
func TestGetEvaluacion(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/evaluacion/25/662f0f8717e615569749a748/663a248407051740146fed88"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetEvaluacion Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetEvaluacion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetEvaluacion:", err.Error())
		t.Fail()
	}
}
func TestPlanesAEvaluar(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/evaluacion/planes"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPlanesAEvaluar Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPlanesAEvaluar Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPlanesAEvaluar:", err.Error())
		t.Fail()
	}
}
func TestGetUnidades(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/evaluacion/unidades/Seguimiento%20PruebaSure/25"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetUnidades Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetUnidades Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetUnidades:", err.Error())
		t.Fail()
	}
}
func TestGetAvances(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/evaluacion/avance/Seguimiento%20PruebaSure/25/8"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetAvances Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetAvances Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetAvances:", err.Error())
		t.Fail()
	}
}
