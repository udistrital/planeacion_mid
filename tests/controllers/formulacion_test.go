package controllers

import (
	"bytes"
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

func TestClonarFormato(t *testing.T) {
	body := []byte(`{
		"dependencia_id": "PRB",
		"vigencia": "2024"
	}`)

	if response, err := http.Post("http://localhost:8081/v1/formulacion/clonar_formato/611e4a2dd403481fb638b6e9", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestClonarFormato Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestClonarFormato Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestClonarFormato:", err.Error())
		t.Fail()
	}
}

func TestConsultarArbolArmonizacion(t *testing.T) {
	body := []byte(`{}`)

	if response, err := http.Post("http://localhost:8081/v1/formulacion/get_arbol_armonizacion/611db9b4d403482fec38b637", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarArbolArmonizacion Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarArbolArmonizacion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarArbolArmonizacion:", err.Error())
		t.Fail()
	}
}

func TestVersionarPlan(t *testing.T) {
	body := []byte(`{}`)

	if response, err := http.Post("http://localhost:8081/v1/formulacion/versionar_plan/61398379df020f786256e5a7", "application/json", bytes.NewBuffer(body)); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestVersionarPlan Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestVersionarPlan Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestVersionarPlan:", err.Error())
		t.Fail()
	}
}

// TODO: CalculosDocentes
// TODO: DefinirFechasFuncionamiento
// TODO: GetPlanesUnidadesComun

func TestGuardarActividad(t *testing.T) {
	body := []byte(`{{
		"armo": "613991d2df020fd81056e5c8",
		"armoPI": "613991d2df020fd81056e5c8",
		"entrada": {"1":{"dato":"Prueba Meta segplan","index":1},"10":{"activo":false,"dato":"Prueba Meta segplan","index":"10"},"11":{"activo":true,"dato":"Prueba Meta segplan","index":11},"12":{"activo":true,"dato":"Prueba Meta segplan","index":12},"13":{"activo":true,"dato":"Prueba Meta segplan","index":13},"14":{"activo":false,"dato":"Prueba Meta segplan","index":"14"},"2":{"dato":"Prueba Meta segplan","index":"2"},"3":{"dato":"Prueba Meta segplan","index":"3"},"4":{"dato":"Prueba Meta segplan","index":4},"5":{"dato":"Prueba Meta segplan","index":5},"6":{"dato":"Prueba Meta segplan","index":6},"7":{"dato":"Prueba Meta segplan","index":7},"8":{"dato":"Prueba Meta segplan","index":8},"9":{"activo":false,"dato":"Prueba Meta segplan","index":"9"}}
	}}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/guardar_actividad/613991d2df020fd81056e5c8", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestGuardarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestGuardarActividad Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

func TestActualizarActividad(t *testing.T) {
	body := []byte(`{{
		"armo": "613991d2df020fd81056e5c8",
		"armoPI": "613991d2df020fd81056e5c8",
		"entrada": {"1":{"dato":"Prueba Meta segplan","index":1},"10":{"activo":false,"dato":"Prueba Meta segplan","index":"10"},"11":{"activo":true,"dato":"Prueba Meta segplan","index":11},"12":{"activo":true,"dato":"Prueba Meta segplan","index":12},"13":{"activo":true,"dato":"Prueba Meta segplan","index":13},"14":{"activo":false,"dato":"Prueba Meta segplan","index":"14"},"2":{"dato":"Prueba Meta segplan","index":"2"},"3":{"dato":"Prueba Meta segplan","index":"3"},"4":{"dato":"Prueba Meta segplan","index":4},"5":{"dato":"Prueba Meta segplan","index":5},"6":{"dato":"Prueba Meta segplan","index":6},"7":{"dato":"Prueba Meta segplan","index":7},"8":{"dato":"Prueba Meta segplan","index":8},"9":{"activo":false,"dato":"Prueba Meta segplan","index":"9"}}
	}}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/actualizar_actividad/613991d2df020fd81056e5c8/1", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestActualizarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestActualizarActividad Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

func TestDesactivarActividad(t *testing.T) {
	body := []byte(`{}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/delete_actividad/613991d2df020fd81056e5c8/1", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestDesactivarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestDesactivarActividad Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

func TestGuardarIdentificacion(t *testing.T) {
	body := []byte(`{
		"nombre": "Identificación de Contratistas Plan de Acción de Funcionamiento 2022",
		"descripcion": "Identificación de Contratistas Plan de Acción de Funcionamiento 2022 OFICINA ASESORA DE ASUNTOS DISCIPLINARIOS",
		"plan_id": "616f6911a985e921bca12e96",
		"dato": "{}",
		"tipo_identificacion_id": "6184b3e6f6fc97850127bb68",
		"activo": true
	  }`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/guardar_identificacion/616f6911a985e921bca12e96/6184b3e6f6fc97850127bb68", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestGuardarIdentificacion Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestGuardarIdentificacion Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

func TestDesactivarIdentificacion(t *testing.T) {
	body := []byte(`{}`)

	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/delete_identificacion/616f6911a985e921bca12e96/6184b3e6f6fc97850127bb68/0", bytes.NewBuffer(body)); err == nil {
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != 200 {
				t.Error("Error TestDesactivarIdentificacion Se esperaba 200 y se obtuvo", response.StatusCode)
				t.Fail()
			} else {
				t.Log("TestDesactivarIdentificacion Finalizado Correctamente (OK)")
			}
		}
	} else {
		t.Error("Error al crear la solicitud PUT: ", err.Error())
		t.Fail()
	}
}

// TODO: CambioCargoIdVinculacionTercero
// TODO: EstructuraPlanes
