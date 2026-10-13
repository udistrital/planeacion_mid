package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGuardarIdentificacionNoDocentePersisteDato(t *testing.T) {
	var almacenado map[string]interface{}
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			if err := json.NewDecoder(r.Body).Decode(&almacenado); err != nil {
				t.Fatalf("no se pudo decodificar el PUT: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"Success": true,
				"Status":  "200",
				"Message": "Successful",
				"Data":    almacenado,
			})
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(almacenado)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer servidor.Close()

	entrada := map[string]interface{}{
		"0": map[string]interface{}{
			"activo":               true,
			"descripcionNecesidad": "Prueba de persistencia",
			"index":                "1",
		},
	}
	identificacion := map[string]interface{}{
		"_id":                    "identificacion-prueba",
		"plan_id":                "plan-prueba",
		"tipo_identificacion_id": "tipo-prueba",
	}

	err := guardarIdentificacionNoDocente(
		context.Background(),
		strings.TrimPrefix(servidor.URL, "http://"),
		"identificacion-prueba",
		identificacion,
		entrada,
	)
	if err != nil {
		t.Fatalf("se esperaba un guardado exitoso: %v", err)
	}

	respuesta, err := http.Get(servidor.URL + "/identificacion/identificacion-prueba")
	if err != nil {
		t.Fatalf("no se pudo consultar el dato guardado: %v", err)
	}
	defer respuesta.Body.Close()

	var consultado map[string]interface{}
	if err := json.NewDecoder(respuesta.Body).Decode(&consultado); err != nil {
		t.Fatalf("no se pudo decodificar el dato consultado: %v", err)
	}
	dato, ok := consultado["dato"].(string)
	if !ok {
		t.Fatalf("el registro persistido no contiene dato como string: %#v", consultado["dato"])
	}
	var persistido map[string]interface{}
	if err := json.Unmarshal([]byte(dato), &persistido); err != nil {
		t.Fatalf("el dato persistido no contiene JSON válido: %v", err)
	}
	registro := persistido["0"].(map[string]interface{})
	if registro["descripcionNecesidad"] != "Prueba de persistencia" || registro["activo"] != true {
		t.Fatalf("el valor consultado no coincide con el guardado: %#v", registro)
	}
}

func TestGuardarIdentificacionNoDocenteReportaRechazoDelCRUD(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"Success": false,
			"Status":  "400",
			"Message": "dato inválido",
			"Data":    nil,
		})
	}))
	defer servidor.Close()

	err := guardarIdentificacionNoDocente(
		context.Background(),
		strings.TrimPrefix(servidor.URL, "http://"),
		"identificacion-prueba",
		map[string]interface{}{"_id": "identificacion-prueba"},
		map[string]interface{}{"0": map[string]interface{}{"activo": true}},
	)
	if err == nil || err.Error() != "dato inválido" {
		t.Fatalf("se esperaba propagar el rechazo del CRUD y se obtuvo: %v", err)
	}
}
