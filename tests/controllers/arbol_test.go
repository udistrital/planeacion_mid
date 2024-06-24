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

// TODO: falta armar petición con :id y JSON
// func TestActivarPlan(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/arbol/activar_plan/:id", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestActivarPlan Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestActivarPlan Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// // TODO: falta armar petición con :id y JSON
// func TestActivarNodo(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/arbol/activar_nodo", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestActivarNodo Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestActivarNodo Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }
