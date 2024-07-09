package controllers

import (
	"net/http"
	"testing"
)

func TestPlanesDeAccion(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/planes_accion/"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPlanesDeAccion Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPlanesDeAccion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPlanesDeAccion:", err.Error())
		t.Fail()
	}
}
