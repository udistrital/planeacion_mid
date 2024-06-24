package controllers

import (
	"net/http"
	"testing"
)

func TestProyecto(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/proyecto/617c3171f6fc97294f27a041"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestProyecto Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestProyecto Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestProyecto:", err.Error())
		t.Fail()
	}
}
func TestConsultarMetasProyecto(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/metaspro/62266cfa16511ea35c5be15f"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarMetasProyecto Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarMetasProyecto Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarMetasProyecto:", err.Error())
		t.Fail()
	}
}
func TestProyectos(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/proyectos/string"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestProyectos Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestProyectos Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestProyectos:", err.Error())
		t.Fail()
	}
}
func TestConsultarPlanIdentificador(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/get_plan/62150f8525e40c9662093e58"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarPlanIdentificador Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarPlanIdentificador Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarPlanIdentificador:", err.Error())
		t.Fail()
	}
}
func TestConsultarPlanInversion(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/get_infoPlan/63f467aeccee49e93a859ace"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarPlanInversion Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarPlanInversion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarPlanInversion:", err.Error())
		t.Fail()
	}
}
func TestConsultarTodasMetasPlan(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/all_metas/63f467aeccee49e93a859ace"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarTodasMetasPlan Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarTodasMetasPlan Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarTodasMetasPlan:", err.Error())
		t.Fail()
	}
}
func TestConsultarMagnitudesProgramadas(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/magnitudes/64367379a280497994a41e46/1"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarMagnitudesProgramadas Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarMagnitudesProgramadas Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarMagnitudesProgramadas:", err.Error())
		t.Fail()
	}
}
func TestVerificarMagnitudesProgramadas(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/inversion/verificar_magnitudes/64367379a280497994a41e46"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestVerificarMagnitudesProgramadas Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestVerificarMagnitudesProgramadas Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestVerificarMagnitudesProgramadas:", err.Error())
		t.Fail()
	}
}

// func TestAgregarProyecto(t *testing.T) {
// 	body := []byte(`{
// 		"codigo_proyecto": "prueba",
// 		"fecha_creacion": "2023-10-29T17:37:53.934Z",
// 		"nombre_proyecto": "Plan Prueba Clonacion",
// 		"fuentes": []
// 	}`)

// 	if response, err := http.Post("http://localhost:8081/v1/inversion/proyecto/", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestAgregarProyecto Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestAgregarProyecto Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestAgregarProyecto:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestGuardarDocumentosInversion(t *testing.T) {
// 	body := []byte(`{{
// 		"file": "si",
// 		"documento": [{
//         "IdTipoDocumento": 66,
//           "nombre": "PRUEBA",
//           "metadatos": {
//             "dato_a": "Soportes planeacion"
//           },
//           "descripcion": "Documento de soporte para proyectos de plan de acción de inversión",
//           "file": "DATA"
//     	}]
// 	}}`)

// 	if response, err := http.Post("http://localhost:8081/v1/inversion/guardar_documentos", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestGuardarDocumentosInversion Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestGuardarDocumentosInversion Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestGuardarDocumentosInversion:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestCrearPlan(t *testing.T) {
// 	body := []byte(`{
// 		"nombre": "Plan de Acción Proyecto de Inversión",
// 		"descripcion": "Formato plan de acción proyecto de inversión",
// 		"tipo_plan_id": "611af8364a34b3b2df3799a0",
// 		"vigencia": "20",
// 		"dependencia_id": "abc",
// 		"aplicativo_id": "plan",
// 		"activo": true,
//         "id": "61398379df020f786256e5a7"
// 	  }`)

// 	if response, err := http.Post("http://localhost:8081/v1/inversion/crearplan", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestCrearPlan Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestCrearPlan Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestCrearPlan:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestCrearGrupoMeta(t *testing.T) {
// 	body := []byte(`{
// 		"id": "618de204f6fc97904a27d902",
// 		"activo": true,
// 		"aplicativo_id": "idPlaneacion",
// 		"dependencia_id": "10",
// 		"descripcion": "Plan de Acción de Funcionamiento 2022 - 01",
// 		"estado_plan_id": "614d3aeb01c7a245952fabff",
// 		"formato": false,
// 		"nombre": "Plan de Acción de Funcionamiento 2022",
// 		"padre_plan_id": "616700961634ad4d31ed6bd6",
// 		"tipo_plan_id": "61639b8c1634adf976ed4b4c",
// 		"vigencia": "3",
// 		"indexMeta": "documento"
// 	  }`)

// 	if response, err := http.Post("http://localhost:8081/v1/inversion/crear_grupo_meta", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestCrearGrupoMeta Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestCrearGrupoMeta Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestCrearGrupoMeta:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestVersionarPlanInversion(t *testing.T) {
// 	body := []byte(`{}`)

// 	if response, err := http.Post("http://localhost:8081/v1/inversion/versionar_plan/618de204f6fc97904a27d902", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestVersionarPlanInversion Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestVersionarPlanInversion Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestVersionarPlanInversion:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestEditarProyecto(t *testing.T) {
// 	body := []byte(`{{
// 		"codigo_proyecto": "prueba",
// 		"fecha_creacion": "2023-10-29T17:37:53.934Z",
// 		"nombre_proyecto": "Plan Prueba Clonacion",
// 		"fuentes": []
// 	}}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/proyecto/65d2c11092a08b3774a7d737", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestEditarProyecto Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestEditarProyecto Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestActualizarSubgrupoDetalle(t *testing.T) {
// 	body := []byte(`{
// 		"nombre": "subgrupo detalle Ejemplo 4 de segundo nivel",
// 		"descripcion": "Ejemplo 4 de segundo nivel",
// 		"subgrupo_id": "6127b50bd40348063138c94a",
// 		"dato": "{\"type\":\"numeric\",\"required\":\"true\"}",
// 		"activo": true,
// 		"fecha_creacion": "2021-09-01T19:54:48.668+00:00"
// 	  }`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/actualiza_sub_detalle/612fda88df020f752656b310", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestActualizarSubgrupoDetalle Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestActualizarSubgrupoDetalle Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestActualizarProyectoGeneral(t *testing.T) {
// 	body := []byte(`{
// 		"nombre": "Formato plan estratégico 2",
// 		"descripcion": "Descripción modificada 2",
// 		"tipo_plan_id": "611af8464a34b3599e3799a2",
// 		"aplicativo_id": "idPlaneacion",
// 		"activo": false,
// 		"fecha_creacion": "2021-08-19T12:10:21.911Z"
// 	  }`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/actualizar_proyecto/611e4a2dd403481fb638b6e9", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestActualizarProyectoGeneral Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestActualizarProyectoGeneral Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// // TODO: TestGuardarMeta
// // SE NECESITA EL JSON
// func TestGuardarMeta(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/guardar_meta", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestGuardarMeta Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestGuardarMeta Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// // TODO: ArmonizarInversion

// // TODO: ActualizarMetaPlan
// // SE NECESITA EL JSON
// func TestActualizarMetaPlan(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/actualizar_meta/622676b216511e8ea55be6de/1", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error ActualizarMetaPlan Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("ActualizarMetaPlan Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// // TODO: InactivarMeta
// // BUSCAR DATO QUE COINCIDA
// func TestInactivarMeta(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/inactivar_meta/622676b216511e8ea55be6de/1", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestInactivarMeta Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestInactivarMeta Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// // TODO: ProgMagnitudesPlan
// // SE NECESITA EL JSON
// func TestProgramarMagnitudesPlan(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/magnitudes/", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestProgramarMagnitudesPlan Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestProgramarMagnitudesPlan Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// // TODO: ActualizarActividad
// // SE NECESITA EL JSON
// func TestActualizarActividadInversion(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/actualizar_actividad/", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestActualizarActividadInversion Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestActualizarActividadInversion Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// // TODO: ActualizarTablaActividad
// // SE NECESITA EL JSON
// func TestActualizarTablaActividad(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/actividad/id/index/tabla", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestActualizarTablaActividad Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestActualizarTablaActividad Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// // TODO: ActualizarPresupuestoMeta
// // SE NECESITA EL JSON
// func TestActualizarPresupuestoMeta(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/inversion/metas/presupuestos/id/index", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestActualizarPresupuestoMeta Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestActualizarPresupuestoMeta Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }
