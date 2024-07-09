package controllers

import (
	"net/http"
	"testing"
)

func TestGetObtenerArbol(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/arbol/6686e63640bc042575ae2b79"); err == nil {
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

// func TestActivarPlan(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/arbol/activar_plan/6686e63640bc042575ae2b79", bytes.NewBuffer(body)); err == nil {
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

// func TestActivarNodo(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/arbol/activar_nodo/6687230e40bc0402d2ae9b06", bytes.NewBuffer(body)); err == nil {
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

// func TestDeletePlan(t *testing.T) {
// 	body := []byte(`{}`)
// 	idPlan := "6686e63640bc042575ae2b79"

// 	if request, err := http.NewRequest(http.MethodDelete, "http://localhost:8081/v1/arbol/desactivar_plan/"+idPlan, bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestDeletePlan Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestDeletePlan Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestDeleteNodo(t *testing.T) {
// 	body := []byte(`{}`)
// 	id := "6687230e40bc0402d2ae9b06"

// 	if request, err := http.NewRequest(http.MethodDelete, "http://localhost:8081/v1/arbol/desactivar_nodo/"+id, bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestDeleteNodo Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestDeleteNodo Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }
