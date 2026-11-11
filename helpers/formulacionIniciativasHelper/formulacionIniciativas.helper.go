package formulacionIniciativasHelper

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func DefinirFechasIniciativasFormulacion(body map[string]interface{}) []interface{} {
	body["planes_interes"] = "[]"
	var respuestasAcumuladas []interface{}
	var res map[string]interface{}
	var respuestaPost = make(map[string]interface{})
	var respuestaPost2 = make(map[string]interface{})

	if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas/buscar-unidad-planes/1", "POST", &respuestaPost, body); err == nil {
		data, ok := respuestaPost["Data"].([]interface{})
		if ok && len(data) > 0 { // Si sí existe se procede a verificar que las fechas sean iguales
			registro := data[0].(map[string]interface{})

			// Parsear las fechas
			fechaInicioBodyParsed, _ := time.Parse(time.RFC3339, body["fecha_inicio"].(string))
			fechaInicioRegistroParsed, _ := time.Parse(time.RFC3339, registro["fecha_inicio"].(string))
			fechaFinBodyParsed, _ := time.Parse(time.RFC3339, body["fecha_fin"].(string))
			fechaFinRegistroParsed, _ := time.Parse(time.RFC3339, registro["fecha_fin"].(string))

			if fechaInicioBodyParsed.Equal(fechaInicioRegistroParsed) && fechaFinBodyParsed.Equal(fechaFinRegistroParsed) {
				// Registro Fechas iguales
				registro = manejarUnidades(body, registro, 1)
				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas/"+registro["_id"].(string), "PUT", &res, registro); err != nil {
					panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
				}
				respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
			} else {
				// Registro con Fechas distintas, entonces se debe editar el registro existente quitando las unidades del body
				// Si unidades_interes del registro se queda vacio, se debe inactivar
				registro = manejarUnidades(body, registro, 2)
				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas/"+registro["_id"].(string), "PUT", &res, registro); err != nil {
					panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
				}
				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas/buscar-unidad-planes/2", "POST", &respuestaPost2, body); err == nil {
					data, ok := respuestaPost2["Data"].([]interface{})
					if ok && len(data) > 0 { // Existe el registro, entonces se le agregan las unidades
						registro2 := data[0].(map[string]interface{})
						registro2 = manejarUnidades(body, registro2, 1)
						if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas/"+registro2["_id"].(string), "PUT", &res, registro2); err != nil {
							panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
						}
						respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
					} else { // No existe el registro, se crea el registro con las unidades del body
						body["planes_interes"] = registro["planes_interes"]
						if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas", "POST", &res, body); err != nil {
							panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error agregaron periodo-seguimiento", "status": "400", "log": err})
						}
						respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
					}

					//? Actualiza registros de seguimiento en caso de que hayan planes avalados con el formato_id igual al id del plan_interes del body
					if body["tipo_seguimiento_id"] == "6af3a20c160752e5dee28c65" {
						periodoSeguimientoIdAntiguo := registro["_id"].(string)
						periodoSeguimientoIdNuevo := res["Data"].(map[string]interface{})["_id"].(string)
						unidadesBody, err := ObtenerArrayUnidadesInteres(body)
						if err != nil {
							panic(map[string]interface{}{"funcion": "manejarUnidades", "err": "Error en la decodificación JSON", "status": "400"})
						}
						CambiarFechasSeguimiento("", unidadesBody, periodoSeguimientoIdAntiguo, periodoSeguimientoIdNuevo)
					}
				}
			}
		} else { // No existe el registro con las unidades especificadas, se procede a validar por periodo, fecha_inicio y fecha_fin
			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas/buscar-unidad-planes/2", "POST", &respuestaPost2, body); err == nil {
				data, ok := respuestaPost2["Data"].([]interface{})
				if ok && len(data) > 0 { // Existe el registro, entonces se le agregan las unidades
					registro := data[0].(map[string]interface{})
					registro = manejarUnidades(body, registro, 1)
					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas/"+registro["_id"].(string), "PUT", &res, registro); err != nil {
						panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
					}
					respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
				} else { // No existe el registro, se crea
					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas", "POST", &res, body); err != nil {
						panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error versionando periodo-seguimiento", "status": "400", "log": err})
					}
					respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
				}
			}
		}
	} else {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error en la peticion", "status": "400", "log": err})
	}
	return respuestasAcumuladas
}

func manejarUnidades(body map[string]interface{}, registro map[string]interface{}, caso int) map[string]interface{} {
	unidadesBody, err := ObtenerArrayUnidadesInteres(body)
	if err != nil {
		panic(map[string]interface{}{"funcion": "manejarUnidades", "err": "Error en la decodificación JSON", "status": "400"})
	}

	unidadesRegistro, err := ObtenerArrayUnidadesInteres(registro)
	if err != nil {
		panic(map[string]interface{}{"funcion": "manejarUnidades", "err": "Error en la decodificación JSON", "status": "400"})
	}

	// Crear un mapa para almacenar los Ids ya presentes en "registro"
	idsPresentes := make(map[int]bool)
	for _, unidad := range unidadesRegistro {
		if id, ok := unidad["Id"].(int); ok {
			idsPresentes[id] = true
		} else {
			panic(map[string]interface{}{"funcion": "manejarUnidades", "err": "Error al obtener el Id de la unidad de interés del registro", "status": "400"})
		}
	}

	// Manejar las unidades de body que no existan en registro
	for _, unidadBody := range unidadesBody {
		if idBody, ok := unidadBody["Id"].(int); ok {
			switch caso {
			case 1:
				// Si el Id de la unidad no existe en registro, agregar la unidad
				if !idsPresentes[idBody] {
					unidadesRegistro = append(unidadesRegistro, unidadBody)
					idsPresentes[idBody] = true
				}
			case 2:
				// Eliminar la unidad del registro si también está en body
				if idsPresentes[idBody] {
					unidadesRegistro = eliminarUnidad(unidadesRegistro, idBody)
				}
			default:
				panic(map[string]interface{}{"funcion": "manejarUnidades", "err": "Caso no válido", "status": "400"})
			}
		} else {
			panic(map[string]interface{}{"funcion": "manejarUnidades", "err": "Error al obtener el Id de la unidad de interés del body", "status": "400"})
		}
	}
	if len(unidadesRegistro) == 0 {
		registro["activo"] = false
	}
	unidadesRegistroJSON, err := json.Marshal(unidadesRegistro)
	if err != nil {
		panic(map[string]interface{}{"funcion": "manejarUnidades", "err": "Error al convertir el array a JSON", "status": "400"})
	}
	registro["unidades_interes"] = string(unidadesRegistroJSON)
	return registro
}

func ObtenerArrayUnidadesInteres(body map[string]interface{}) ([]map[string]interface{}, error) {
	// Acceder a la lista de unidades de interés
	unidadesInteresString, ok := body["unidades_interes"].(string)
	if !ok {
		panic(map[string]interface{}{"funcion": "ObtenerArrayUnidadesInteres", "err": "No se pudo obtener la lista de unidades de interés como cadena JSON", "status": "400"})
	}

	// Decodificar la cadena JSON a una lista de interfaces
	var unidadesInteres []interface{}
	err := json.Unmarshal([]byte(unidadesInteresString), &unidadesInteres)
	if err != nil {
		panic(map[string]interface{}{"funcion": "ObtenerArrayUnidadesInteres", "err": "Error al decodificar la lista de unidades de interés", "status": "400"})
	}

	// Almacenar las unidades de interés en un array
	var unidades []map[string]interface{}
	for _, unidad := range unidadesInteres {
		if unidadMap, ok := unidad.(map[string]interface{}); ok {
			// Verificar si la clave "Id" está presente y tiene el tipo correcto
			if id, idOk := unidadMap["Id"].(float64); idOk {
				// Convertir el valor de "Id" a entero
				idInt := int(id)
				unidades = append(unidades, map[string]interface{}{
					"Id":     idInt,
					"Nombre": unidadMap["Nombre"],
					// Otras propiedades de interés si es necesario
				})
			} else {
				panic(map[string]interface{}{"funcion": "ObtenerArrayUnidadesInteres", "err": "El valor de Id no es un número", "status": "400"})
			}
		}
	}
	return unidades, nil
}

func eliminarUnidad(unidades []map[string]interface{}, id int) []map[string]interface{} {
	for i, unidad := range unidades {
		if unidad["Id"].(int) == id {
			// Eliminar la unidad del slice
			unidades = append(unidades[:i], unidades[i+1:]...)
		}
	}
	return unidades
}

func CambiarFechasSeguimiento(planInteresString string, unidadesInteres []map[string]interface{}, periodoSeguimientoIdAntiguo string, periodoSeguimientoIdNuevo string) {
	var resPlanesAvalados map[string]interface{}
	var resSeguimientos map[string]interface{}
	var planesAvalados []map[string]interface{}
	var seguimiento []map[string]interface{}
	var seguimientoActualizado map[string]interface{}
	var estadoAvaladoId = "6153355601c7a2365b2fb2a1"
	var planInteres models.PlanInteres

	for _, unidad := range unidadesInteres {
		idUnidad := strconv.Itoa(unidad["Id"].(int))

		// 1. Se buscan planes avalados con el _id de la plantilla
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=formato_id:"+planInteres.Id+",estado_plan_id:"+estadoAvaladoId+",dependencia_id:"+idUnidad+",activo:true", &resPlanesAvalados); err != nil {
			panic(map[string]interface{}{"funcion": "CambiarFechasSeguimiento", "err": "Error en la peticion", "status": "400", "log": err})
		}
		helpers.LimpiezaRespuestaRefactor(resPlanesAvalados, &planesAvalados)

		if len(planesAvalados) < 1 { //? No se encontraron planes avalados con el formato_id del id del plan_interes
			continue
		}

		// 2. Se buscan los registro de seguimiento en caso de que la respuesta anterior contenga registros

		for _, planAvalado := range planesAvalados {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=plan_id:"+planAvalado["_id"].(string)+",periodo_seguimiento_id:"+periodoSeguimientoIdAntiguo+",activo:true", &resSeguimientos); err != nil {
				panic(map[string]interface{}{"funcion": "CambiarFechasSeguimiento", "err": "Error en la peticion", "status": "400", "log": err})
			}
			helpers.LimpiezaRespuestaRefactor(resSeguimientos, &seguimiento)

			if len(seguimiento) < 1 { //? No se encontraron registros de seguimiento
				continue
			}

			// 3. Se actualizan los registros de seguimiento
			seguimientoRegistro := seguimiento[0]
			seguimientoRegistro["periodo_seguimiento_id"] = periodoSeguimientoIdNuevo
			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento[0]["_id"].(string), "PUT", &seguimientoActualizado, seguimientoRegistro); err != nil {
				panic(map[string]interface{}{"funcion": "CambiarFechasSeguimiento", "err": "Error actualizando seguimiento", "status": "400", "log": err})
			}
		}
	}
}