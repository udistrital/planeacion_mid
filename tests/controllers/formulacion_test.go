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

func TestVinculacionTerceroByEmail(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/vinculacion_tercero_email/mpcastroc@udistrital.edu.co"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestVinculacionTerceroByEmail Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestVinculacionTerceroByEmail Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestVinculacionTerceroByEmail:", err.Error())
		t.Fail()
	}
}

func TestVinculacionTerceroByIdentificacion(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/formulacion/vinculacion_tercero_identificacion/1022435418"); err == nil {
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

// func TestClonarFormato(t *testing.T) {
// 	body := []byte(`{
// 		"dependencia_id": "PRB",
// 		"vigencia": "2024"
// 	}`)

// 	if response, err := http.Post("http://localhost:8081/v1/formulacion/clonar_formato/611e4a2dd403481fb638b6e9", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestClonarFormato Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestClonarFormato Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestClonarFormato:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestConsultarArbolArmonizacion(t *testing.T) {
// 	body := []byte(`{}`)

// 	if response, err := http.Post("http://localhost:8081/v1/formulacion/get_arbol_armonizacion/611db9b4d403482fec38b637", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestConsultarArbolArmonizacion Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestConsultarArbolArmonizacion Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestConsultarArbolArmonizacion:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestVersionarPlan(t *testing.T) {
// 	body := []byte(`{}`)

// 	if response, err := http.Post("http://localhost:8081/v1/formulacion/versionar_plan/61398379df020f786256e5a7", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestVersionarPlan Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestVersionarPlan Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestVersionarPlan:", err.Error())
// 		t.Fail()
// 	}
// }
// TODO: Faltan datos de prueba
// func TestCalculosDocentes(t *testing.T) {
// 	body := []byte(`{}`)

// 	if response, err := http.Post("http://localhost:8081/v1/formulacion/calculos_docentes", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestCalculosDocentes Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestCalculosDocentes Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestCalculosDocentes:", err.Error())
// 		t.Fail()
// 	}
// }
// func TestDefinirFechasFuncionamiento(t *testing.T) {

// 	body := []byte(`{
// 	"activo": true,
// 	"fecha_fin": "2024-07-14T00:00:00.000Z",
// 	"fecha_inicio": "2024-07-01T00:00:00.000Z",
// 	"periodo_id": "39",
// 	"planes_interes": "[{\"_id\":\"667a5074252f5d1000cc7f14\",\"nombre\":\"Pruebas AFM\"}]",
// 	"tipo_seguimiento_id": "6260e975ebe1e6498f7404ee",
// 	"unidades_interes": "[{\"Id\":8,\"Nombre\":\"VICERRECTORIA ACADEMICA\"}]",
// 	"usuario_modificacion": "52505205"}`)

// 	if response, err := http.Post("http://localhost:8081/v1/formulacion/habilitar_fechas_funcionamiento", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestDefinirFechasFuncionamiento Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestDefinirFechasFuncionamiento Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestDefinirFechasFuncionamiento:", err.Error())
// 		t.Fail()
// 	}

// }

// GetPlanesUnidadesComun no se consume desde el cliente

// func TestGuardarActividad(t *testing.T) {
// 	body := []byte(`{"armo":"616647531634ade426ed535a","armoPI":"652ca721ce026b594962570e","entrada":{"66872d8540bc041bddaea70e":"","66872d8540bc041bddaea70e_o":"","66872d8540bc044cd4aea720":"investigadores","66872d8540bc044cd4aea720_o":"","66872d8540bc04b2deaea715":"","66872d8540bc04b2deaea715_o":"","66872d8640bc043c1faea72a":"Σ investigadores","66872d8640bc043c1faea72a_o":"","66872d8640bc04e044aea734":200,"66872d8640bc04e044aea734_o":"","66872d8740bc0407d3aea759":"Denominador fijo","66872d8740bc0407d3aea759_o":"","66872d8740bc046293aea741":"Unidad","66872d8740bc046293aea741_o":"","66872d8740bc04b147aea74d":"Creciente","66872d8740bc04b147aea74d_o":"","66872d8840bc0431e5aea765":100,"66872d8840bc0431e5aea765_o":"","66872d8840bc044ec8aea76d":"CERI","66872d8840bc044ec8aea76d_o":""}}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/guardar_actividad/66872d8340bc047b72aea6ef", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestGuardarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestGuardarActividad Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestActualizarActividad(t *testing.T) {
// 	body := []byte(`{"armo":"616647531634ade426ed535a","armoPI":"652ca721ce026b594962570e","entrada":{"66872d8340bc0464f3aea6f2":"Actividad 1","66872d8340bc0470cfaea6f9":"1.1 actividad 1\n1.2 act 2\n1.3 act","66872d8440bc040c37aea707":"Toda la vigencia","66872d8440bc041501aea700":"Investigación","66872d8540bc041bddaea70e":"","66872d8540bc041bddaea70e_o":"Sin observación","66872d8540bc044cd4aea720":"investigadores","66872d8540bc044cd4aea720_o":"Sin observación","66872d8540bc04b2deaea715":"","66872d8540bc04b2deaea715_o":"Sin observación","66872d8640bc043c1faea72a":"Σ investigadores","66872d8640bc043c1faea72a_o":"Sin observación","66872d8640bc04e044aea734":200,"66872d8640bc04e044aea734_o":"Sin observación","66872d8740bc0407d3aea759":"Denominador fijo","66872d8740bc0407d3aea759_o":"Sin observación","66872d8740bc046293aea741":"Unidad","66872d8740bc046293aea741_o":"Sin observación","66872d8740bc04b147aea74d":"Creciente","66872d8740bc04b147aea74d_o":"Sin observación","66872d8840bc0431e5aea765":100,"66872d8840bc0431e5aea765_o":"Sin observación","66872d8840bc044ec8aea76d":"CERI","66872d8840bc044ec8aea76d_o":"Sin observación"}}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/actualizar_actividad/66872d8340bc047b72aea6ef/1", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestActualizarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestActualizarActividad Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestDesactivarActividad(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/delete_actividad/613991d2df020fd81056e5c8/1", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestDesactivarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestDesactivarActividad Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestGuardarIdentificacion(t *testing.T) {
// 	body := []byte(`{"0":{"descripcionNecesidad":"Investigadores iniciales","perfil":487,"cantidad":20,"meses":6,"dias":0,"valorUnitario":"$6,382,169.00","valorUnitarioInc":"$7,020,385.90","valorTotal":"$765,860,280.00","valorTotalInc":"$842,446,308.00","actividades":["1"],"rubro":"2.1.1.01.01.001.01.1","rubroNombre":"Sueldo Básico Administrativos","total":"","totalInc":"842446308.00","activo":true,"index":"1"}}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/guardar_identificacion/66872d8340bc047b72aea6ef/6184b3e6f6fc97850127bb68", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestGuardarIdentificacion Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestGuardarIdentificacion Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestDesactivarIdentificacion(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/delete_identificacion/616f6911a985e921bca12e96/6184b3e6f6fc97850127bb68/0", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestDesactivarIdentificacion Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestDesactivarIdentificacion Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// TODO: Faltan datos de prueba
// func TestCambioCargoIdVinculacionTercero(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/cargo_vinculacion/7230282", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestCambioCargoIdVinculacionTercero Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestCambioCargoIdVinculacionTercero Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestEstructuraPlanes(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/formulacion/estructura_planes/6686e63640bc042575ae2b79", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestEstructuraPlanes Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestEstructuraPlanes Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }
