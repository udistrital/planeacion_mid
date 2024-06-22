package controllers

import (
	"net/http"
	"testing"
)

func TestGetFormato(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formato/666791d5252f5d2f0dbf9f6e"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetFormato Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetFormato Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetFormato:", err.Error())
		t.Fail()
	}
}
