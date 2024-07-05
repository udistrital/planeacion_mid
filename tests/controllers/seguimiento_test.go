package controllers

import (
	"net/http"
	"testing"
)

func TestTrimestres(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/trimestres/35"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestTrimestres Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestTrimestres Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestTrimestres:", err.Error())
		t.Fail()
	}
}
func TestPeriodos(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_periodos/35"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPeriodos Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPeriodos Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPeriodos:", err.Error())
		t.Fail()
	}
}
func TestActividades(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_actividades/61f60e4525e40c6f5d084185"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestActividades Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestActividades Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestActividades:", err.Error())
		t.Fail()
	}
}
func TestSeguimiento(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_seguimiento/61f08edc25e40c91b0083e4f/1/635b1f995073f2675157dc7f"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestSeguimiento Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestSeguimiento:", err.Error())
		t.Fail()
	}
}

func TestIndicadores(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_indicadores/6201d43f25e40c205608b459"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestIndicadores Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestIndicadores Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestIndicadores:", err.Error())
		t.Fail()
	}
}

func TestConsultarEstadoTrimestre(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/get_estado_trimestre/666a3b79252f5d633cc097de/T2"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarEstadoTrimestre Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarEstadoTrimestre Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarEstadoTrimestre:", err.Error())
		t.Fail()
	}
}

func TestConsultarEstadoTrimestres(t *testing.T) {
	if response, err := http.Get("http://localhost:8081/v1/seguimiento/estado_trimestres/666a3b79252f5d633cc097de"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestConsultarEstadoTrimestres Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestConsultarEstadoTrimestres Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestConsultarEstadoTrimestres:", err.Error())
		t.Fail()
	}
}

// TODO: Faltan datos de prueba
// func TestAvalarPlan(t *testing.T) {
// 	body := []byte(`{{}}`)

// 	if response, err := http.Post("http://localhost:8081/v1/seguimiento/avalar/:id", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestAvalarPlan Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestAvalarPlan Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestAvalarPlan:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestCrearReportes(t *testing.T) {
// 	body := []byte(`{{}}`)

// 	if response, err := http.Post("http://localhost:8081/v1/seguimiento/crear_reportes/61f08edc25e40c91b0083e4f/61f236f525e40c582a0840d0", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestCrearReportes Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestCrearReportes Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestCrearReportes:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestObtenerAvanceIndicador(t *testing.T) {
// 	body := []byte(`{{
// 		"plan_id": "",
// 		"periodo_seguimiento_id": "",
// 		"index": "",
// 		"Nombre_del_indicador": ""
// 	}}`)

// 	if response, err := http.Post("http://localhost:8081/v1/seguimiento/get_avance", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestObtenerAvanceIndicador Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestObtenerAvanceIndicador Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestObtenerAvanceIndicador:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestMigrarInformacion(t *testing.T) {
// 	body := []byte(`{}`)

// 	if response, err := http.Post("http://localhost:8081/v1/seguimiento/migrar_seguimiento/61f08edc25e40c91b0083e4f/635b1f995073f2675157dc7f", "application/json", bytes.NewBuffer(body)); err == nil {
// 		if response.StatusCode != 200 {
// 			t.Error("Error TestMigrarInformacion Se esperaba 200 y se obtuvo", response.StatusCode)
// 			t.Fail()
// 		} else {
// 			t.Log("TestMigrarInformacion Finalizado Correctamente (OK)")
// 		}
// 	} else {
// 		t.Error("Error TestMigrarInformacion:", err.Error())
// 		t.Fail()
// 	}
// }

// func TestGuardarSeguimiento(t *testing.T) {
// 	body := []byte(`{"_id":"667a6f4f252f5da293ccc9d9","cualitativo":{"dificultades":"asdfsfa","observaciones":"","productos":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.","reporte":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old."},"cuantitativo":{"indicadores":[{"denominador":"","detalleReporte":"","formula":"investigadores/200","meta":90,"nombre":"Investigadores involucrados","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"1","reporteNumerador":"90","tendencia":"Creciente","unidad":"Unidad"},{"denominador":"","detalleReporte":"","formula":"conocen el proyecto/total estudiantes pregrado","meta":100,"nombre":"Porcentaje estudiantes pregrado que conozcan el proyecto","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"100","reporteNumerador":"30","tendencia":"Creciente","unidad":"Porcentaje"},{"denominador":"","detalleReporte":"","formula":"gente que le gusta el proyecto/total gente","meta":70,"nombre":"Tasa de aprobación","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"2","reporteNumerador":"30","tendencia":"Creciente","unidad":"Tasa"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":90,"avanceAcumulado":1,"brechaExistente":0,"divisionCero":false,"indicador":90,"indicadorAcumulado":90,"meta":90,"nombre":"Investigadores involucrados","unidad":"Unidad"},{"acumuladoDenominador":100,"acumuladoNumerador":30,"avanceAcumulado":0.3,"brechaExistente":0.7,"divisionCero":false,"indicador":0.3,"indicadorAcumulado":0.3,"meta":100,"nombre":"Porcentaje estudiantes pregrado que conozcan el proyecto","unidad":"Porcentaje"},{"acumuladoDenominador":2,"acumuladoNumerador":30,"avanceAcumulado":0.214,"brechaExistente":55,"divisionCero":false,"indicador":15,"indicadorAcumulado":15,"meta":70,"nombre":"Tasa de aprobación","unidad":"Tasa"}]},"estado":{"id":"635fb7205795c6891895e89b","nombre":"Con observaciones"},"estadoSeguimiento":"Revisión Verificada con Observaciones","evidencia":[],"id":"6687133e40bc04262eae6ef8","id_actividad":"AdxMCCIOI0CUXVvEX","informacion":{"descripcion":"Actividad 1","index":"1","nombre":"Pruebas AFM","periodo":"Toda la vigencia","ponderacion":100,"producto":"Investigación sobre IA","tarea":"","trimestre":"T1","unidad":"FACULTAD DE INGENIERIA"}}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/guardar_seguimiento/667a6d98252f5d44a5ccbb2c/1/66857cbe40bc04f718ad368f", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestGuardarSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestGuardarSeguimiento Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestGuardarDocumentosSeguimiento(t *testing.T) {
// 	body := []byte(`{
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
// 		"evidencia": []
// 	}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/guardar_documentos/prueba/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestGuardarDocumentosSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestGuardarDocumentosSeguimiento Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestGuardarCualitativo(t *testing.T) {
// 	body := []byte(`{
// 		"_id": "",
// 		"informacion": {
// 			"descripcion": "Generar y desarrollar proyectos institucionales e interinstitucionales para apoyar desarrollo de competencia didáctica de los profesores.",
// 			"index": "1",
// 			"nombre": "Plan de acción 2023 Prod Seguimiento",
// 			"periodo": "Toda la vigencia",
// 			"ponderacion": 25,
// 			"producto": "• Prototipo de talleres para formación de estudiantes para profesores validados\n• Informe",
// 			"tarea": " Busqueda de convocatoria para presentar el Proyecto (0,05)\n • Desarrollo de actividades mediante el proyecto transmigrarts (0,3)\n • Convenio interinstitucional UD-IBERO, Vértice SAS e Inclusive Movimiento (0,3)\n • Consolidar equipos de investigadores para el diseño de Ambientes de Aprendizaje Accesibles y Afectivos en el marco del proyecto SGRColciencias (0,05) \n •  Desarrollar estrategias de articulación con unidades institucionales para apoyar al profesorado de la UDFJC (0,3)",
// 			"trimestre": "T2",
// 			"unidad": "FACULTAD DE INGENIERIA"
// 		},
// 		"evidencias": [
// 			{
// 				"Activo": true,
// 				"Enlace": "c275df0d-b27a-4446-b434-69ee6b71d2c1",
// 				"Id": 148376,
// 				"Observacion": "",
// 				"TipoDocumento": {
// 					"codigoAbreviacion": "DSPA",
// 					"id": 60
// 				},
// 				"nombre": "UNAL - 230210 - SOSTENIBILIDAD ECONÓMICA (CC).docx"
// 			},
// 			{
// 				"Activo": true,
// 				"Enlace": "927b1653-7c23-40f1-bd04-e35a98e9e72a",
// 				"Id": 148377,
// 				"Observacion": "",
// 				"TipoDocumento": {
// 					"codigoAbreviacion": "DSPA",
// 					"id": 60
// 				},
// 				"nombre": "Acreditación de programas.xlsx"
// 			}
// 		],
// 		"cualitativo": {
// 			"dificultades": "grdfvbfdesafghndawSHGFRDFRGTHYJKULIHGJUYTRFEHJVNGBFEW3R4T5Y67UTJYHFGRCUTYRFESRGTHYJUGHGjhkjgfrdgh",
// 			"observaciones": "",
// 			"productos": "gbvnmbhgfvdsafghjm",
// 			"reporte": "Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:"
// 		},
// 		"cuantitativo": {
// 			"indicadores": [
// 				{
// 					"denominador": "Denominador fijo",
// 					"detalleReporte": "",
// 					"formula": "∑ % avance en la tarea *ponderación de la tarea",
// 					"meta": "80",
// 					"nombre": "Avance en el Proyecto SGR Colciencias",
// 					"observaciones": "",
// 					"reporteDenominador": "1",
// 					"reporteNumerador": "0.3",
// 					"tendencia": "Creciente",
// 					"unidad": "Porcentaje"
// 				},
// 				{
// 					"denominador": "Denominador fijo",
// 					"detalleReporte": "",
// 					"formula": "∑ cursos impartidos durante la vigencia ",
// 					"meta": 7,
// 					"nombre": "Cursos Cultiva articulados a la Red Acacia",
// 					"observaciones": "",
// 					"reporteDenominador": "1",
// 					"reporteNumerador": "2",
// 					"tendencia": "Creciente",
// 					"unidad": "Unidad"
// 				}
// 			],
// 			"resultados": [
// 				{
// 					"acumuladoDenominador": 1,
// 					"acumuladoNumerador": 0.55,
// 					"avanceAcumulado": 0.6875,
// 					"brechaExistente": 0.25,
// 					"divisionCero": false,
// 					"indicador": 0.3,
// 					"indicadorAcumulado": 0.55,
// 					"meta": 80,
// 					"nombre": "Avance en el Proyecto SGR Colciencias",
// 					"unidad": "Porcentaje"
// 				},
// 				{
// 					"acumuladoDenominador": 1,
// 					"acumuladoNumerador": 3,
// 					"avanceAcumulado": 0.429,
// 					"brechaExistente": 4,
// 					"divisionCero": false,
// 					"indicador": 2,
// 					"indicadorAcumulado": 3,
// 					"meta": 7,
// 					"nombre": "Cursos Cultiva articulados a la Red Acacia",
// 					"unidad": "Unidad"
// 				}
// 			]
// 		},
// 		"dependencia": false
// 	}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/guardar_cualitativo/63b5fecb1598303a848fe7b8/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestGuardarCualitativo Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestGuardarCualitativo Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestGuardarCuantitativo(t *testing.T) {
// 	body := []byte(`{"_id":"6687133e40bc04262eae6ef8","informacion":{"descripcion":"Actividad 1","index":"1","nombre":"Pruebas AFM","periodo":"Toda la vigencia","ponderacion":100,"producto":"Investigación sobre IA","tarea":"","trimestre":"T1","unidad":"FACULTAD DE INGENIERIA"},"evidencias":[],"cualitativo":{"dificultades":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.","observaciones":"","productos":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.","reporte":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old."},"cuantitativo":{"indicadores":[{"denominador":"","detalleReporte":"","formula":"investigadores/200","meta":90,"nombre":"Investigadores involucrados","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"1","reporteNumerador":"90","tendencia":"Creciente","unidad":"Unidad"},{"denominador":"","detalleReporte":"","formula":"conocen el proyecto/total estudiantes pregrado","meta":100,"nombre":"Porcentaje estudiantes pregrado que conozcan el proyecto","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"100","reporteNumerador":"30","tendencia":"Creciente","unidad":"Porcentaje"},{"denominador":"","detalleReporte":"","formula":"gente que le gusta el proyecto/total gente","meta":70,"nombre":"Tasa de aprobación","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"2","reporteNumerador":"30","tendencia":"Creciente","unidad":"Tasa"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":90,"avanceAcumulado":1,"brechaExistente":0,"divisionCero":false,"indicador":90,"indicadorAcumulado":90,"meta":90,"nombre":"Investigadores involucrados","unidad":"Unidad"},{"acumuladoDenominador":100,"acumuladoNumerador":30,"avanceAcumulado":0.3,"brechaExistente":0.7,"divisionCero":false,"indicador":0.3,"indicadorAcumulado":0.3,"meta":100,"nombre":"Porcentaje estudiantes pregrado que conozcan el proyecto","unidad":"Porcentaje"},{"acumuladoDenominador":2,"acumuladoNumerador":30,"avanceAcumulado":0.214,"brechaExistente":55,"divisionCero":false,"indicador":15,"indicadorAcumulado":15,"meta":70,"nombre":"Tasa de aprobación","unidad":"Tasa"}]},"dependencia":false}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/guardar_cuantitativo/667a6d98252f5d44a5ccbb2c/1/66857cbe40bc04f718ad368f", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestGuardarCuantitativo Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestGuardarCuantitativo Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestReportarActividad(t *testing.T) {
// 	body := []byte(`{"SeguimientoId":"667a6f4f252f5da293ccc9d9"}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/reportar_actividad/1", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestReportarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestReportarActividad Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestReportarSeguimiento(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/reportar_seguimiento/61f238db25e40ccb450840db", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestReportarSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestReportarSeguimiento Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestRevisarActividad(t *testing.T) {
// 	body := []byte(`{"_id":"6687133e40bc04262eae6ef8","informacion":{"descripcion":"Actividad 1","index":"1","nombre":"Pruebas AFM","periodo":"Toda la vigencia","ponderacion":100,"producto":"Investigación sobre IA","tarea":"","trimestre":"T1","unidad":"FACULTAD DE INGENIERIA"},"evidencias":[],"cualitativo":{"dificultades":"asdfsfa","observaciones":"","productos":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.","reporte":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old."},"cuantitativo":{"indicadores":[{"denominador":"","detalleReporte":"","formula":"investigadores/200","meta":90,"nombre":"Investigadores involucrados","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"1","reporteNumerador":"90","tendencia":"Creciente","unidad":"Unidad"},{"denominador":"","detalleReporte":"","formula":"conocen el proyecto/total estudiantes pregrado","meta":100,"nombre":"Porcentaje estudiantes pregrado que conozcan el proyecto","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"100","reporteNumerador":"30","tendencia":"Creciente","unidad":"Porcentaje"},{"denominador":"","detalleReporte":"","formula":"gente que le gusta el proyecto/total gente","meta":70,"nombre":"Tasa de aprobación","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"2","reporteNumerador":"30","tendencia":"Creciente","unidad":"Tasa"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":90,"avanceAcumulado":1,"brechaExistente":0,"divisionCero":false,"indicador":90,"indicadorAcumulado":90,"meta":90,"nombre":"Investigadores involucrados","unidad":"Unidad"},{"acumuladoDenominador":100,"acumuladoNumerador":30,"avanceAcumulado":0.3,"brechaExistente":0.7,"divisionCero":false,"indicador":0.3,"indicadorAcumulado":0.3,"meta":100,"nombre":"Porcentaje estudiantes pregrado que conozcan el proyecto","unidad":"Porcentaje"},{"acumuladoDenominador":2,"acumuladoNumerador":30,"avanceAcumulado":0.214,"brechaExistente":55,"divisionCero":false,"indicador":15,"indicadorAcumulado":15,"meta":70,"nombre":"Tasa de aprobación","unidad":"Tasa"}]},"dependencia":false}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/revision_actividad/667a6d98252f5d44a5ccbb2c/1/66857cbe40bc04f718ad368f", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestRevisarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestRevisarActividad Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestRevisarActividadJefeDependencia(t *testing.T) {
// 	body := []byte(`{"_id":"667a6f4f252f5da293ccc9d9","cualitativo":{"dificultades":"asdfsfa","observaciones":"","productos":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.","reporte":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old."},"cuantitativo":{"indicadores":[{"denominador":"","detalleReporte":"","formula":"investigadores/200","meta":90,"nombre":"Investigadores involucrados","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"1","reporteNumerador":"90","tendencia":"Creciente","unidad":"Unidad"},{"denominador":"","detalleReporte":"","formula":"conocen el proyecto/total estudiantes pregrado","meta":100,"nombre":"Porcentaje estudiantes pregrado que conozcan el proyecto","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"100","reporteNumerador":"30","tendencia":"Creciente","unidad":"Porcentaje"},{"denominador":"","detalleReporte":"","formula":"gente que le gusta el proyecto/total gente","meta":70,"nombre":"Tasa de aprobación","observaciones":"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. ","reporteDenominador":"2","reporteNumerador":"30","tendencia":"Creciente","unidad":"Tasa"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":90,"avanceAcumulado":1,"brechaExistente":0,"divisionCero":false,"indicador":90,"indicadorAcumulado":90,"meta":90,"nombre":"Investigadores involucrados","unidad":"Unidad"},{"acumuladoDenominador":100,"acumuladoNumerador":30,"avanceAcumulado":0.3,"brechaExistente":0.7,"divisionCero":false,"indicador":0.3,"indicadorAcumulado":0.3,"meta":100,"nombre":"Porcentaje estudiantes pregrado que conozcan el proyecto","unidad":"Porcentaje"},{"acumuladoDenominador":2,"acumuladoNumerador":30,"avanceAcumulado":0.214,"brechaExistente":55,"divisionCero":false,"indicador":15,"indicadorAcumulado":15,"meta":70,"nombre":"Tasa de aprobación","unidad":"Tasa"}]},"estado":{"id":"635fb6eb5795c6891895e899","nombre":"Actividad reportada"},"estadoSeguimiento":"En revisión JU","evidencia":[],"id":"6687133e40bc04262eae6ef8","id_actividad":"AdxMCCIOI0CUXVvEX","informacion":{"descripcion":"Actividad 1","index":"1","nombre":"Pruebas AFM","periodo":"Toda la vigencia","ponderacion":100,"producto":"Investigación sobre IA","tarea":"","trimestre":"T1","unidad":"FACULTAD DE INGENIERIA"}}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/revision_actividad_jefe_dependencia/667a6d98252f5d44a5ccbb2c/1/66857cbe40bc04f718ad368f", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestRevisarActividadJefeDependencia Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestRevisarActividadJefeDependencia Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }
// func TestRevisarSeguimiento(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/revision_seguimiento/639a42e954a3d2399c3bb6ff", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestRevisarSeguimiento Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestRevisarSeguimiento Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// TODO: Faltan datos de prueba
// func TestRevisarSeguimientoJefeDependencia(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/revision_seguimiento_jefe_dependencia/:id", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestRevisarSeguimientoJefeDependencia Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestRevisarSeguimientoJefeDependencia Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// func TestRetornarActividad(t *testing.T) {
// 	body := []byte(`{"_id":"63b6ec611598306b5c90095b","cualitativo":{"dificultades":"grdfvbfdesafghndawSHGFRDFRGTHYJKULIHGJUYTRFEHJVNGBFEW3R4T5Y67UTJYHFGRCUTYRFESRGTHYJUGHGjhkjgfrdgh","observaciones":"","productos":"gbvnmbhgfvdsafghjm","reporte":"Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:Adicionalmente, con el propósito de consolidar el número de programas Acreditados en Alta Calidad, se tramitaron diferentes procesos ante el Consejo Nacional de Acreditación, CNA, tanto para obtención como para renovación, de los cuales se recoge información en la siguiente tabla:"},"cuantitativo":{"indicadores":[{"denominador":"Denominador fijo","detalleReporte":"","formula":"∑ % avance en la tarea *ponderación de la tarea","meta":"80","nombre":"Avance en el Proyecto SGR Colciencias","observaciones":"","reporteDenominador":"1","reporteNumerador":"0.3","tendencia":"Creciente","unidad":"Porcentaje"},{"denominador":"Denominador fijo","detalleReporte":"","formula":"∑ cursos impartidos durante la vigencia ","meta":7,"nombre":"Cursos Cultiva articulados a la Red Acacia","observaciones":"","reporteDenominador":"1","reporteNumerador":"2","tendencia":"Creciente","unidad":"Unidad"}],"resultados":[{"acumuladoDenominador":1,"acumuladoNumerador":0.55,"avanceAcumulado":0.6875,"brechaExistente":0.25,"divisionCero":false,"indicador":0.3,"indicadorAcumulado":0.55,"meta":80,"nombre":"Avance en el Proyecto SGR Colciencias","unidad":"Porcentaje"},{"acumuladoDenominador":1,"acumuladoNumerador":3,"avanceAcumulado":0.429,"brechaExistente":4,"divisionCero":false,"indicador":2,"indicadorAcumulado":3,"meta":7,"nombre":"Cursos Cultiva articulados a la Red Acacia","unidad":"Unidad"}]},"estado":{"id":"63793207242b813898e9856b","nombre":"Actividad avalada"},"estadoSeguimiento":"En revisión OAPC","evidencia":[{"Activo":true,"Enlace":"c275df0d-b27a-4446-b434-69ee6b71d2c1","Id":148376,"Observacion":"","TipoDocumento":{"codigoAbreviacion":"DSPA","id":60},"nombre":"UNAL - 230210 - SOSTENIBILIDAD ECONÓMICA (CC).docx"},{"Activo":true,"Enlace":"927b1653-7c23-40f1-bd04-e35a98e9e72a","Id":148377,"Observacion":"","TipoDocumento":{"codigoAbreviacion":"DSPA","id":60},"nombre":"Acreditación de programas.xlsx"}],"id":"","informacion":{"descripcion":"Generar y desarrollar proyectos institucionales e interinstitucionales para apoyar desarrollo de competencia didáctica de los profesores.","index":"1","nombre":"Plan de acción 2023 Prod Seguimiento","periodo":"Toda la vigencia","ponderacion":25,"producto":"• Prototipo de talleres para formación de estudiantes para profesores validados\n• Informe","tarea":" Busqueda de convocatoria para presentar el Proyecto (0,05)\n • Desarrollo de actividades mediante el proyecto transmigrarts (0,3)\n • Convenio interinstitucional UD-IBERO, Vértice SAS e Inclusive Movimiento (0,3)\n • Consolidar equipos de investigadores para el diseño de Ambientes de Aprendizaje Accesibles y Afectivos en el marco del proyecto SGRColciencias (0,05) \n •  Desarrollar estrategias de articulación con unidades institucionales para apoyar al profesorado de la UDFJC (0,3)","trimestre":"T2","unidad":"FACULTAD DE INGENIERIA"}}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8080/v1/seguimiento/retornar_actividad/63b5fecb1598303a848fe7b8/1/635b1f795073f2675157dc7d", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestRetornarActividad Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestRetornarActividad Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }

// TODO: Faltan datos de prueba
// func TestRetornarActividadJefeDependencia(t *testing.T) {
// 	body := []byte(`{}`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/retornar_actividad_jefe_dependencia/:plan_id/:index/:trimestre", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestRetornarActividadJefeDependencia Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestRetornarActividadJefeDependencia Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }
// func TestHabilitarReportes(t *testing.T) {
// 	body := []byte(`{
// 		"_id": "635b1f995073f2675157dc7f",
// 		"fecha_inicio": "2024-01-17T00:00:00.000Z",
// 		"fecha_fin": "2024-02-23T23:59:59.000Z",
// 		"periodo_id": "314",
// 		"activo": true,
// 		"fecha_creacion": "2022-10-28T00:17:29.116Z",
// 		"fecha_modificacion": "2024-03-04T14:42:36.388Z",
// 		"__v": 0,
// 		"tipo_seguimiento_id": "61f236f525e40c582a0840d0",
// 		"unidades_interes": "[{\"Id\":8,\"Nombre\":\"VICERRECTORIA ACADEMICA\"}]",
// 		"planes_interes": "[{\"_id\":\"628ce817ebe1e6512a74b32e\",\"nombre\":\"prueba nueva\"}]"
// 	  }`)

// 	if request, err := http.NewRequest(http.MethodPut, "http://localhost:8081/v1/seguimiento/habilitar_reportes", bytes.NewBuffer(body)); err == nil {
// 		client := &http.Client{}
// 		if response, err := client.Do(request); err == nil {
// 			if response.StatusCode != 200 {
// 				t.Error("Error TestHabilitarReportes Se esperaba 200 y se obtuvo", response.StatusCode)
// 				t.Fail()
// 			} else {
// 				t.Log("TestHabilitarReportes Finalizado Correctamente (OK)")
// 			}
// 		}
// 	} else {
// 		t.Error("Error al crear la solicitud PUT: ", err.Error())
// 		t.Fail()
// 	}
// }
