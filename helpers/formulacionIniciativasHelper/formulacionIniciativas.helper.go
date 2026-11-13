package formulacionIniciativasHelper

import (
	"encoding/json"
	"time"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
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

func ValidarFechasIniciativasFormulacion(body map[string]interface{}) map[string]interface{} {
	var resFechas map[string]interface{}
	var resTipoSeguimiento map[string]interface{}
	var fechas []map[string]interface{}
	var tipoSeguimiento []map[string]interface{}
	respuesta := map[string]interface{}{"unidadEncontrada": false, "habilitado": false}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-seguimiento/?query=codigo_abreviacion:FIPI_SP", &resTipoSeguimiento); err != nil {
		panic(map[string]interface{}{"funcion": "ValidarFechasIniciativasFormulacion", "err": "Error al obtener el tipo de seguimiento", "status": "400", "log": err})
	}
	helpers.LimpiezaRespuestaRefactor(resTipoSeguimiento, &tipoSeguimiento)
	if len(tipoSeguimiento) == 0 {
		panic(map[string]interface{}{"funcion": "ValidarFechasIniciativasFormulacion", "err": "No se encontró el tipo de seguimiento", "status": "404"})
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/?query=tipo_seguimiento_id:"+tipoSeguimiento[0]["_id"].(string)+",activo:true", &resFechas); err != nil {
		panic(map[string]interface{}{"funcion": "ValidarFechasIniciativasFormulacion", "err": "Error al obtener las fechas", "status": "400", "log": err})
	}
	helpers.LimpiezaRespuestaRefactor(resFechas, &fechas)
	if len(fechas) == 0 {
		panic(map[string]interface{}{"funcion": "ValidarFechasIniciativasFormulacion", "err": "No hay periodos activos", "status": "404"})
	}

	unidadesFechas, err := ObtenerArrayUnidadesInteres(fechas[0])
	if err != nil {
		panic(map[string]interface{}{"funcion": "ValidarFechasIniciativasFormulacion", "err": "Error en la decodificación JSON", "status": "400"})
	}

	// 1. Validar que la unidad esté; si no, se descarta todo el proceso
	unidadEncontrada := false
	for _, unidad := range unidadesFechas {
		if fmt.Sprint(body["dependenciaId"]) == fmt.Sprint(unidad["Id"]) {
			unidadEncontrada = true
			break
		}
	}
	if !unidadEncontrada {
		return respuesta
	}
	respuesta["unidadEncontrada"] = true

	// 2. Comparar la fecha de hoy con el rango
	bogota := time.FixedZone("COT", -5*60*60)
	ahora := time.Now().In(bogota)
	hoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, bogota)

	inicio, err := parseFecha(fmt.Sprint(fechas[0]["fecha_inicio"]), bogota)
	if err != nil {
		panic(map[string]interface{}{"funcion": "ValidarFechasIniciativasFormulacion", "err": "fecha_inicio inválida", "status": "400", "log": err})
	}
	fin, err := parseFecha(fmt.Sprint(fechas[0]["fecha_fin"]), bogota)
	if err != nil {
		panic(map[string]interface{}{"funcion": "ValidarFechasIniciativasFormulacion", "err": "fecha_fin inválida", "status": "400", "log": err})
	}

	// Inclusivo: hoy >= inicio y hoy <= fin
	respuesta["habilitado"] = !hoy.Before(inicio) && !hoy.After(fin)

	return respuesta
}

// Toma solo la parte YYYY-MM-DD para comparar por día y evitar desfases de zona horaria
func parseFecha(valor string, loc *time.Location) (time.Time, error) {
	if len(valor) < 10 {
		return time.Time{}, fmt.Errorf("fecha inválida: %q", valor)
	}
	return time.ParseInLocation("2006-01-02", valor[:10], loc)
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