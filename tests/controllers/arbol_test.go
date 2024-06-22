package controllers

import (
	"net/http"
	"testing"
)

func TestGetObtenerArbol(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/arbol/616610841634adc541ed5333"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetObtenerArbol Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetObtenerArbol Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetObtenerArbol:", err.Error())
		t.Fail()
	}
}
