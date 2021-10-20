package formulacionhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

var validDataT = []string{}

func Limpia() {
	validDataT = []string{}
}

func ClonarHijos(hijos []map[string]interface{}, padre string) {

	clienteHttp := &http.Client{}
	url := beego.AppConfig.String("PlanesService") + "/subgrupo/registrar_nodo/"

	for i := 0; i < len(hijos); i++ {

		hijo := make(map[string]interface{})
		hijo["nombre"] = hijos[i]["nombre"]
		hijo["descripcion"] = hijos[i]["descripcion"]
		hijo["activo"] = hijos[i]["activo"]
		hijo["padre"] = padre
		hijo["bandera_tabla"] = hijos[i]["bandera_tabla"]

		var resPost map[string]interface{}
		var resLimpia map[string]interface{}

		aux, err := json.Marshal(hijo)
		if err != nil {
			log.Fatalf("Error codificado: %v", err)
		}

		peticion, err := http.NewRequest("POST", url, bytes.NewBuffer(aux))
		if err != nil {
			log.Fatalf("Error creando peticion: %v", err)
		}

		peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
		respuesta, err := clienteHttp.Do(peticion)
		if err != nil {
			log.Fatalf("Error haciendo peticion: %v", err)
		}

		defer respuesta.Body.Close()

		cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
		if err != nil {
			log.Fatalf("Error leyendo peticion: %v", err)
		}

		json.Unmarshal(cuerpoRespuesta, &resPost)
		resLimpia = resPost["Data"].(map[string]interface{})

		var respuestaHijos map[string]interface{}
		var respuestaHijosDetalle map[string]interface{}
		var subHijos []map[string]interface{}
		var subHijosDetalle []map[string]interface{}

		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijos[i]["_id"].(string), &respuestaHijosDetalle); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijosDetalle, &subHijosDetalle)
			ClonarHijosDetalle(subHijosDetalle, resLimpia["_id"].(string))
		}

		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+hijos[i]["_id"].(string), &respuestaHijos); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijos, &subHijos)
			ClonarHijos(subHijos, resLimpia["_id"].(string))
		}

	}
}

func ClonarHijosDetalle(subHijosDetalle []map[string]interface{}, subgrupo_id string) {
	clienteHttp := &http.Client{}
	url := beego.AppConfig.String("PlanesService") + "/subgrupo-detalle/"

	for i := 0; i < len(subHijosDetalle); i++ {
		hijoDetalle := make(map[string]interface{})
		hijoDetalle["nombre"] = subHijosDetalle[i]["nombre"]
		hijoDetalle["descripcion"] = subHijosDetalle[i]["descripcion"]
		hijoDetalle["subgrupo_id"] = subgrupo_id
		hijoDetalle["activo"] = subHijosDetalle[i]["activo"]
		hijoDetalle["dato"] = subHijosDetalle[i]["dato"]

		var resPost map[string]interface{}
		aux, err := json.Marshal(hijoDetalle)
		if err != nil {
			log.Fatalf("Error codificado: %v", err)
		}

		peticion, err := http.NewRequest("POST", url, bytes.NewBuffer(aux))
		if err != nil {
			log.Fatalf("Error creando peticion: %v", err)
		}

		peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
		respuesta, err := clienteHttp.Do(peticion)
		if err != nil {
			log.Fatalf("Error haciendo peticion: %v", err)
		}

		defer respuesta.Body.Close()

		cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
		if err != nil {
			log.Fatalf("Error leyendo peticion: %v", err)
		}

		json.Unmarshal(cuerpoRespuesta, &resPost)

	}
}

func BuildTreeFa(hijos []map[string]interface{}, index string) [][]map[string]interface{} {

	var tree []map[string]interface{}
	var requeridos []map[string]interface{}
	var res map[string]interface{}
	var resLimpia []map[string]interface{}
	var result [][]map[string]interface{}
	var nodo map[string]interface{}

	for i := 0; i < len(hijos); i++ {
		if hijos[i]["activo"] == true {

			forkData := make(map[string]interface{})
			var id string
			forkData["id"] = hijos[i]["_id"]
			forkData["nombre"] = hijos[i]["nombre"]
			jsonString, _ := json.Marshal(hijos[i]["_id"])
			json.Unmarshal(jsonString, &id)

			if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &resLimpia)
				nodo = resLimpia[0]
				if len(nodo) == 0 {
					//forkData["type"] = ""
					//forkData["required"] = ""
				} else {

					var deta map[string]interface{}
					dato_str := nodo["dato"].(string)
					json.Unmarshal([]byte(dato_str), &deta)
					if (deta["type"] != nil) && (deta["required"] != nil) && (deta["options"] == nil) {
						forkData["type"] = deta["type"]
						forkData["required"] = deta["required"]
					} else if (deta["type"] != nil) && (deta["required"] != nil) && (deta["options"] != nil) {
						forkData["type"] = deta["type"]
						forkData["required"] = deta["required"]
						forkData["options"] = deta["options"]
					} else {
						forkData["type"] = " "
						forkData["required"] = " "
					}

				}
			}

			if len(hijos[i]["hijos"].([]interface{})) > 0 {

				forkData["sub"] = make([]map[string]interface{}, len(getChildren(hijos[i]["hijos"].([]interface{}))))
				forkData["sub"] = getChildren(hijos[i]["hijos"].([]interface{}))
			} else {
				forkData["sub"] = ""
			}

			tree = append(tree, forkData)
			add(id)
		}
	}
	requeridos = convert(validDataT, index)
	result = append(result, tree)
	result = append(result, requeridos)
	return result
}

func getChildren(children []interface{}) (childrenTree []map[string]interface{}) {
	var res map[string]interface{}
	var resp map[string]interface{}
	var nodo map[string]interface{}
	var detalle []map[string]interface{}

	for _, child := range children {
		childStr := fmt.Sprintf("%v", child)
		forkData := make(map[string]interface{})
		var id string
		err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+childStr, &res)
		if err != nil {
			return
		}
		helpers.LimpiezaRespuestaRefactor(res, &nodo)
		if nodo["activo"] == true {
			forkData["id"] = nodo["_id"]
			forkData["nombre"] = nodo["nombre"]
			jsonString, _ := json.Marshal(nodo["_id"])
			json.Unmarshal(jsonString, &id)
			if err_ := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &resp); err_ == nil {
				helpers.LimpiezaRespuestaRefactor(resp, &detalle)
				if len(detalle) == 0 {
					// forkData["type"] = ""
					// forkData["required"] = ""
				} else {

					var deta map[string]interface{}
					dato_str := fmt.Sprintf("%v", detalle[0]["dato"])
					json.Unmarshal([]byte(dato_str), &deta)
					// forkData["type"] = deta["type"]
					// forkData["required"] = deta["required"]
					if (deta["type"] != nil) && (deta["required"] != nil) && (deta["options"] == nil) {
						forkData["type"] = deta["type"]
						forkData["required"] = deta["required"]
					} else if (deta["type"] != nil) && (deta["required"] != nil) && (deta["options"] != nil) {
						forkData["type"] = deta["type"]
						forkData["required"] = deta["required"]
						forkData["options"] = deta["options"]
					} else {
						forkData["type"] = " "
						forkData["required"] = " "
					}
				}
			}
			if len(nodo["hijos"].([]interface{})) > 0 {
				if len(getChildren(nodo["hijos"].([]interface{}))) == 0 {
					forkData["sub"] = ""
				} else {
					forkData["sub"] = getChildren(nodo["hijos"].([]interface{}))
				}
			}

			childrenTree = append(childrenTree, forkData)
		}
		add(id)
	}
	return
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func add(id string) {
	if !contains(validDataT, id) {
		validDataT = append(validDataT, id)
	}
}

func convert(valid []string, index string) []map[string]interface{} {
	var validadores []map[string]interface{}
	var res map[string]interface{}
	var subgrupo_detalle []map[string]interface{}
	var dato_plan map[string]interface{}
	var actividad map[string]interface{}

	forkData := make(map[string]interface{})
	for _, v := range valid {
		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+v, &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &subgrupo_detalle)

			if len(subgrupo_detalle) > 0 {
				if subgrupo_detalle[0]["dato_plan"] != nil {
					dato_plan_str := subgrupo_detalle[0]["dato_plan"].(string)
					json.Unmarshal([]byte(dato_plan_str), &dato_plan)

					if dato_plan[index] == nil {

					} else {
						actividad = dato_plan[index].(map[string]interface{})
						forkData[v] = actividad["dato"]
						if actividad["observacion"] != nil {
							keyObservacion := v + "_o"
							forkData[keyObservacion] = getObservacion(actividad)

						}

					}

				} else {
					forkData[v] = ""
				}
			} else {
				forkData[v] = ""

			}
		}

	}
	validadores = append(validadores, forkData)
	return validadores
}

func getObservacion(actividad map[string]interface{}) string {
	if actividad["observacion"] == nil {
		return ""
	} else {
		str := fmt.Sprintf("%v", actividad["observacion"])
		return str
	}
}

func RecorrerHijos(hijos []map[string]interface{}, index string) {
	var subgrupo map[string]interface{}
	var respuesta map[string]interface{}

	for i := 0; i < len(hijos); i++ {
		if len(hijos[i]["hijos"].([]interface{})) != 0 {
			hijosSubgrupo := hijos[i]["hijos"].([]interface{})
			for j := 0; j < len(hijosSubgrupo); j++ {
				if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijosSubgrupo[j].(string), &respuesta); err == nil {
					helpers.LimpiezaRespuestaRefactor(respuesta, &subgrupo)
					if len(subgrupo["hijos"].([]interface{})) != 0 {
						recorrerSubgrupos(subgrupo["hijos"].([]interface{}), index)
					} else {
						desactivarActividad(subgrupo["_id"].(string), index)
					}
				} else {
					panic(map[string]interface{}{"funcion": "DeleteActividad", "err": "Error obteniendo subgrupo \"subgrupo[\"_id\"].(string)\"", "status": "400", "log": err})
				}
			}
		} else {
			desactivarActividad(hijos[i]["_id"].(string), index)
		}
	}
}

func recorrerSubgrupos(hijos []interface{}, index string) {
	var respuesta map[string]interface{}
	var subgrupo map[string]interface{}

	for i := 0; i < len(hijos); i++ {
		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijos[i].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &subgrupo)
			if len(subgrupo["hijos"].([]interface{})) != 0 {
				recorrerSubgrupos(subgrupo["hijos"].([]interface{}), index)
			} else {
				desactivarActividad(subgrupo["_id"].(string), index)
			}
		} else {
			panic(map[string]interface{}{"funcion": "DeleteActividad", "err": "Error obteniendo subgrupo \"subgrupo[\"_id\"].(string)\"", "status": "400", "log": err})
		}
	}
}

func desactivarActividad(subgrupo_id string, index string) {
	var respuesta map[string]interface{}
	var res map[string]interface{}
	var respuestaLimpia []map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var dato_plan map[string]interface{}

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo_id, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
		subgrupoDetalle = respuestaLimpia[0]
		if subgrupoDetalle["dato_plan"] != nil {
			actividad := make(map[string]interface{})
			dato_plan_str := subgrupoDetalle["dato_plan"].(string)
			json.Unmarshal([]byte(dato_plan_str), &dato_plan)
			for index_actividad := range dato_plan {
				if index_actividad == index {
					aux_actividad := dato_plan[index_actividad].(map[string]interface{})
					actividad["index"] = index_actividad
					actividad["dato"] = aux_actividad["dato"]
					if aux_actividad["observacion"] != nil {
						actividad["observacion"] = aux_actividad["observacion"]
					}
					actividad["activo"] = false
					dato_plan[index_actividad] = actividad
				}
			}
			b, _ := json.Marshal(dato_plan)
			str := string(b)
			subgrupoDetalle["dato_plan"] = str
		}
		if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupoDetalle["_id"].(string), "PUT", &res, subgrupoDetalle); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteActividad", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
		}
	} else {
		panic(map[string]interface{}{"funcion": "DeleteActividad", "err": "Error obteniendo subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
	}
}

var data_source []map[string]interface{}
var displayed_columns []string

func LimpiaTabla() {
	data_source = nil
	displayed_columns = nil
}

func GetTabla(hijos []interface{}) map[string]interface{} {
	tabla := make(map[string]interface{})
	var respuesta map[string]interface{}
	var subgrupo map[string]interface{}

	for i := 0; i < len(hijos); i++ {
		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijos[i].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &subgrupo)
			if subgrupo["bandera_tabla"] == true {
				displayed_columns = append(displayed_columns, subgrupo["nombre"].(string))
				getActividadTabla(subgrupo)
			}

			if len(subgrupo["hijos"].([]interface{})) != 0 {
				var respuestaHijos map[string]interface{}
				var subgrupoHijo map[string]interface{}
				for j := 0; j < len(subgrupo["hijos"].([]interface{})); j++ {
					if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijos[i].(string), &respuestaHijos); err == nil {
						helpers.LimpiezaRespuestaRefactor(respuesta, &subgrupoHijo)
						if subgrupo["bandera_tabla"] == true {
							displayed_columns = append(displayed_columns, subgrupoHijo["nombre"].(string))
							getActividadTabla(subgrupoHijo)
						}
					}
				}
			}
		}
	}
	tabla["displayed_columns"] = displayed_columns
	tabla["data_source"] = data_source
	return tabla
}

func getActividadTabla(subgrupo map[string]interface{}) {
	var respuesta map[string]interface{}
	var respuestaLimpia []map[string]interface{}
	var subgrupo_detalle map[string]interface{}
	var dato_plan map[string]interface{}
	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo["_id"].(string), &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
		subgrupo_detalle = respuestaLimpia[0]
		if data_source == nil {
			if subgrupo_detalle["dato_plan"] != nil {
				dato_plan_str := subgrupo_detalle["dato_plan"].(string)
				json.Unmarshal([]byte(dato_plan_str), &dato_plan)
				for key := range dato_plan {
					actividad := make(map[string]interface{})
					element := dato_plan[key].(map[string]interface{})
					actividad["index"] = key
					actividad[subgrupo["nombre"].(string)] = element["dato"]
					actividad["activo"] = element["activo"]
					data_source = append(data_source, actividad)
				}
			}
		} else {
			for i := 0; i < len(data_source); i++ {
				if subgrupo_detalle["dato_plan"] != nil {
					var data = data_source[i]
					dato_plan_str := subgrupo_detalle["dato_plan"].(string)
					json.Unmarshal([]byte(dato_plan_str), &dato_plan)
					element := dato_plan[data["index"].(string)].(map[string]interface{})
					data[subgrupo["nombre"].(string)] = element["dato"]

				}

			}
		}

	}
}

func GetArmonizacion(id string) map[string]interface{} {
	var armonizacion map[string]interface{}
	var respuesta map[string]interface{}
	var subgrupo map[string]interface{}
	var recorrido []map[string]interface{}
	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+id, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &subgrupo)
		recorrido = append(recorrido, subgrupo)
		for subgrupo != nil || subgrupo["padre"] != nil {
			var auxRes map[string]interface{}
			var auxSub map[string]interface{}
			if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+subgrupo["padre"].(string), &auxRes); err == nil {
				helpers.LimpiezaRespuestaRefactor(auxRes, &auxSub)
				recorrido = append(recorrido, auxSub)
			}
			subgrupo = auxSub
		}
	}
	recorrido = recorrido[:len(recorrido)-1]
	armonizacion = getRama(recorrido)

	return armonizacion
}

func getRama(recorrido []map[string]interface{}) map[string]interface{} {

	var armonizacion map[string]interface{}

	for i := 0; i < len(recorrido); i++ {
		forkData := make(map[string]interface{})
		forkData["_id"] = recorrido[i]["_id"]
		forkData["nombre"] = recorrido[i]["nombre"]
		forkData["descripcion"] = recorrido[i]["descripcion"]
		forkData["activo"] = recorrido[i]["activo"]
		if armonizacion != nil {
			forkData["children"] = armonizacion
		}

		armonizacion = forkData

	}

	return armonizacion
}
