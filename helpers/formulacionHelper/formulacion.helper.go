package formulacionhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

const (
	CodigoTipoPlan                  string = "PL_SP"
	CodigoTipoPlanAccionFormulacion string = "PAF_SP"
)

var validDataT = []string{}

func Limpia() {
	validDataT = []string{}
}

func ClonarHijos(hijos []map[string]interface{}, padre string) {

	clienteHttp := &http.Client{}
	url := "http://" + beego.AppConfig.String("PlanesService") + "/subgrupo/registrar_nodo/"

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

		cuerpoRespuesta, err := io.ReadAll(respuesta.Body)
		if err != nil {
			log.Fatalf("Error leyendo peticion: %v", err)
		}

		json.Unmarshal(cuerpoRespuesta, &resPost)
		resLimpia = resPost["Data"].(map[string]interface{})

		var respuestaHijos map[string]interface{}
		var respuestaHijosDetalle map[string]interface{}
		var subHijos []map[string]interface{}
		var subHijosDetalle []map[string]interface{}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijos[i]["_id"].(string), &respuestaHijosDetalle); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijosDetalle, &subHijosDetalle)
			ClonarHijosDetalle(subHijosDetalle, resLimpia["_id"].(string))
		}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+hijos[i]["_id"].(string), &respuestaHijos); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijos, &subHijos)
			ClonarHijos(subHijos, resLimpia["_id"].(string))
		}

	}
}

func ClonarHijosDetalle(subHijosDetalle []map[string]interface{}, subgrupo_id string) {
	clienteHttp := &http.Client{}
	url := "http://" + beego.AppConfig.String("PlanesService") + "/subgrupo-detalle/"

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

		cuerpoRespuesta, err := io.ReadAll(respuesta.Body)
		if err != nil {
			log.Fatalf("Error leyendo peticion: %v", err)
		}

		json.Unmarshal(cuerpoRespuesta, &resPost)

	}
}

func BuildTreeFa(hijos []map[string]interface{}, index string) [][]map[string]interface{} {

	var tree []map[string]interface{}
	var requeridos []map[string]interface{}
	armonizacion := make([]map[string]interface{}, 1)
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

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &res); err == nil {
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
	requeridos, armonizacion[0] = convert(validDataT, index)
	result = append(result, tree)
	result = append(result, requeridos)
	result = append(result, armonizacion)
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
		err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+childStr, &res)
		if err != nil {
			return
		}
		helpers.LimpiezaRespuestaRefactor(res, &nodo)
		if nodo["activo"] == true {
			forkData["id"] = nodo["_id"]
			forkData["nombre"] = nodo["nombre"]
			jsonString, _ := json.Marshal(nodo["_id"])
			json.Unmarshal(jsonString, &id)
			if err_ := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &resp); err_ == nil {
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
	if !contains(validDataT, id) && id != "" {
		validDataT = append(validDataT, id)
	}
}

func convert(valid []string, index string) ([]map[string]interface{}, map[string]interface{}) {
	var validadores []map[string]interface{}
	var res map[string]interface{}
	var dato_armonizacion map[string]interface{}
	armonizacion := make(map[string]interface{})
	forkData := make(map[string]interface{})
	for _, v := range valid {
		var subgrupo_detalle []map[string]interface{}
		var actividad map[string]interface{}
		var dato_plan map[string]interface{}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+v, &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &subgrupo_detalle)

			if len(subgrupo_detalle) > 0 {
				if subgrupo_detalle[0]["armonizacion_dato"] != nil {
					dato_armonizacion_str := subgrupo_detalle[0]["armonizacion_dato"].(string)
					json.Unmarshal([]byte(dato_armonizacion_str), &dato_armonizacion)
					aux := dato_armonizacion[index]
					if aux != nil {
						armonizacion["armo"] = dato_armonizacion[index].(map[string]interface{})["armonizacionPED"]
						armonizacion["armoPI"] = dato_armonizacion[index].(map[string]interface{})["armonizacionPI"]
						armonizacion["fuentesActividad"] = dato_armonizacion[index].(map[string]interface{})["fuentesActividad"]
						armonizacion["indexMetaSubProI"] = dato_armonizacion[index].(map[string]interface{})["indexMetaSubProI"]
						armonizacion["ponderacionH"] = dato_armonizacion[index].(map[string]interface{})["ponderacionH"]
						armonizacion["object"] = aux
						fmt.Println(armonizacion, "Meta consultada")
					}
				}
				if subgrupo_detalle[0]["dato_plan"] != nil {
					dato_plan_str := subgrupo_detalle[0]["dato_plan"].(string)
					json.Unmarshal([]byte(dato_plan_str), &dato_plan)
					if dato_plan[index] == nil {
						forkData[v] = ""
					} else {
						actividad = dato_plan[index].(map[string]interface{})
						if v != "" {
							forkData[v] = actividad["dato"]
							if actividad["observacion"] != nil {
								keyObservacion := v + "_o"
								forkData[keyObservacion] = getObservacion(actividad)
							} else {
								keyObservacion := v + "_o"
								forkData[keyObservacion] = "Sin observación"
							}
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
	return validadores, armonizacion
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
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijosSubgrupo[j].(string), &respuesta); err == nil {
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
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijos[i].(string), &respuesta); err == nil {
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

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo_id, &respuesta); err == nil {
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
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupoDetalle["_id"].(string), "PUT", &res, subgrupoDetalle); err != nil {
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
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijos[i].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &subgrupo)
			if subgrupo["bandera_tabla"] == true {
				displayed_columns = append(displayed_columns, subgrupo["nombre"].(string))
				getActividadTabla(subgrupo)
			}

			if len(subgrupo["hijos"].([]interface{})) != 0 {
				var respuestaHijos map[string]interface{}
				var subgrupoHijo map[string]interface{}
				for j := 0; j < len(subgrupo["hijos"].([]interface{})); j++ {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+subgrupo["hijos"].([]interface{})[j].(string), &respuestaHijos); err == nil {
						helpers.LimpiezaRespuestaRefactor(respuestaHijos, &subgrupoHijo)
						if subgrupoHijo["bandera_tabla"] == true {
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
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo["_id"].(string), &respuesta); err == nil {
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
					if dato_plan[data["index"].(string)] != nil {
						element := dato_plan[data["index"].(string)].(map[string]interface{})
						data[subgrupo["nombre"].(string)] = element["dato"]
					}

				}

			}
		}

	}
	sort.SliceStable(data_source, func(i, j int) bool {
		a, _ := strconv.Atoi(data_source[i]["index"].(string))
		b, _ := strconv.Atoi(data_source[j]["index"].(string))
		return a < b
	})
}

func GetArmonizacion(id string) map[string]interface{} {
	var armonizacion map[string]interface{}
	var respuesta map[string]interface{}
	var subgrupo map[string]interface{}
	var recorrido []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+id, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &subgrupo)
		recorrido = append(recorrido, subgrupo)
		for subgrupo != nil || subgrupo["padre"] != nil {
			var auxRes map[string]interface{}
			var auxSub map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+subgrupo["padre"].(string), &auxRes); err == nil {
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

func VersionarHijos(hijos []map[string]interface{}, padre string) {

	var respuestaPost map[string]interface{}
	var subgrupoVersionado map[string]interface{}
	for i := 0; i < len(hijos); i++ {
		hijo := make(map[string]interface{})
		hijo["nombre"] = hijos[i]["nombre"]
		hijo["descripcion"] = hijos[i]["descripcion"]
		hijo["activo"] = hijos[i]["activo"]
		hijo["padre"] = padre
		hijo["bandera_tabla"] = hijos[i]["bandera_tabla"]

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/registrar_nodo", "POST", &respuestaPost, hijo); err != nil {
			panic(map[string]interface{}{"funcion": "VersionarHijos", "err": "Error versionando subgrupo \"hijo[\"_id\"].(string)\"", "status": "400", "log": err})
		}
		subgrupoVersionado = respuestaPost["Data"].(map[string]interface{})

		var respuestaHijos map[string]interface{}
		var respuestaHijosDetalle map[string]interface{}
		var subHijos []map[string]interface{}
		var subHijosDetalle []map[string]interface{}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijos[i]["_id"].(string), &respuestaHijosDetalle); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijosDetalle, &subHijosDetalle)
			VersionarHijosDetalle(subHijosDetalle, subgrupoVersionado["_id"].(string))
		}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+hijos[i]["_id"].(string), &respuestaHijos); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijos, &subHijos)
			VersionarHijos(subHijos, subgrupoVersionado["_id"].(string))
		}

	}
}

func VersionarHijosDetalle(subHijosDetalle []map[string]interface{}, subgrupo_id string) {
	for i := 0; i < len(subHijosDetalle); i++ {
		hijoDetalle := make(map[string]interface{})
		hijoDetalle["nombre"] = subHijosDetalle[i]["nombre"]
		hijoDetalle["descripcion"] = subHijosDetalle[i]["descripcion"]
		hijoDetalle["subgrupo_id"] = subgrupo_id
		hijoDetalle["activo"] = subHijosDetalle[i]["activo"]
		hijoDetalle["dato"] = subHijosDetalle[i]["dato"]
		hijoDetalle["dato_plan"] = subHijosDetalle[i]["dato_plan"]
		hijoDetalle["armonizacion_dato"] = subHijosDetalle[i]["armonizacion_dato"]

		var respuestaPost map[string]interface{}

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle", "POST", &respuestaPost, hijoDetalle); err != nil {
			panic(map[string]interface{}{"funcion": "VersionarHijosDetalle", "err": "Error versionando subgrupo_detalle ", "status": "400", "log": err})
		}

	}
}

func VersionarIdentificaciones(identificaciones []map[string]interface{}, id string) {
	for i := 0; i < len(identificaciones); i++ {
		var aux map[string]interface{} = identificaciones[i]
		identificacion := make(map[string]interface{})
		var respuestaPost map[string]interface{}

		identificacion["nombre"] = aux["nombre"]
		identificacion["descripcion"] = aux["descripcion"]
		identificacion["plan_id"] = id
		identificacion["dato"] = aux["dato"]
		identificacion["tipo_identificacion_id"] = aux["tipo_identificacion_id"]
		identificacion["activo"] = aux["activo"]

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion", "POST", &respuestaPost, identificacion); err != nil {
			panic(map[string]interface{}{"funcion": "VersionaIdentificaciones", "err": "Error versionando identificaciones", "status": "400", "log": err})
		}
	}
}

func OrdenarVersiones(versiones []map[string]interface{}) []map[string]interface{} {
	var versionesOrdenadas []map[string]interface{}

	for i := range versiones {
		if versiones[i]["padre_plan_id"] == nil || versiones[i]["padre_plan_id"] == versiones[i]["_id"] {
			versionesOrdenadas = append(versionesOrdenadas, versiones[i])
		}
	}

	for len(versionesOrdenadas) < len(versiones) {
		versionesOrdenadas = append(versionesOrdenadas, getVersionHija(versionesOrdenadas[len(versionesOrdenadas)-1]["_id"], versiones))
	}

	return versionesOrdenadas
}

func getVersionHija(id interface{}, versiones []map[string]interface{}) map[string]interface{} {
	for i := range versiones {
		if versiones[i]["padre_plan_id"] == id {
			return versiones[i]
		}
	}
	return nil
}

func GetHijosRubro(entrada []interface{}) []map[string]interface{} {
	var hojas []map[string]interface{}

	var respuesta map[string]interface{}
	var resLimpia interface{}
	for i := 0; i < len(entrada); i++ {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanCuentasService")+"/arbol_rubro/"+entrada[i].(string), &respuesta); err == nil {
			if respuesta["Body"] != nil {
				resLimpia = respuesta["Body"].(map[string]interface{})
				if len(resLimpia.(map[string]interface{})["Hijos"].([]interface{})) != 0 {
					var aux = resLimpia.(map[string]interface{})["Hijos"]
					hojas = append(hojas, GetHijosRubro(aux.([]interface{}))...)
				} else {
					hojas = append(hojas, resLimpia.(map[string]interface{}))
				}
			}
		} else {
			panic(map[string]interface{}{"funcion": "GetHijosRubros", "err": "Error arbol_rubros", "status": "400", "log": err})
		}
	}
	return hojas
}

func VerificarDataIdentificaciones(identificaciones []map[string]interface{}, tipoUnidad string) bool {
	var bandera bool

	if tipoUnidad == "facultad" {
		for i := 0; i < len(identificaciones); i++ {
			identificacion := identificaciones[i]
			if identificacion["tipo_identificacion_id"] == "61897518f6fc97091727c3c3" {
				if identificacion["dato"] == "{}" {
					bandera = false
					break
				} else {
					bandera = true
				}
			}
			if identificacion["tipo_identificacion_id"] == "6184b3e6f6fc97850127bb68" {
				if identificacion["dato"] == "{}" {
					bandera = false
					break
				} else {
					bandera = true
				}
			}
		}
	} else if tipoUnidad == "unidad" {
		for i := 0; i < len(identificaciones); i++ {
			identificacion := identificaciones[i]
			if identificacion["tipo_identificacion_id"] == "6184b3e6f6fc97850127bb68" {
				if identificacion["dato"] == "{}" {
					bandera = false
					break
				} else {
					bandera = true
				}
			}
		}
	}
	return bandera
}

func GetIndexActividad(entrada map[string]interface{}) int {
	var respuesta map[string]interface{}

	var respuestaLimpia []map[string]interface{}
	var subgrupo_detalle map[string]interface{}
	dato_plan := make(map[string]interface{})
	var maxIndex = 0

	for key, element := range entrada {
		if element != "" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+key, &respuesta); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
			}
			helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
			subgrupo_detalle = respuestaLimpia[0]

			if subgrupo_detalle["dato_plan"] == nil {
				maxIndex = 0
			} else {
				dato_plan_str := subgrupo_detalle["dato_plan"].(string)
				json.Unmarshal([]byte(dato_plan_str), &dato_plan)
				for key2 := range dato_plan {
					index, err := strconv.Atoi(key2)
					if err != nil {
						log.Panic(err)
					}
					if index > maxIndex {
						maxIndex = index
					}
				}
			}
		}
	}

	return maxIndex
}

func filtrarPlanes(data []map[string]interface{}, f func(map[string]interface{}) bool) []map[string]interface{} {
	fltd := make([]map[string]interface{}, 0)
	for _, v := range data {
		if f(v) {
			fltd = append(fltd, v)
		}
	}
	return fltd
}

func getNumVersion(data []map[string]interface{}, f func(map[string]interface{}) bool) int {
	for i, e := range data {
		if f(e) {
			return i + 1
		}
	}
	return -1
}

func getVigencias() map[string]float64 {
	defer func() {
		if err := recover(); err != nil {
			panic(map[string]interface{}{"funcion": "getVigencias", "err": "Error obteniendo las vigencias", "status": "400", "log": err})
		}
	}()
	var respuestaVigencias map[string]interface{}
	var respuesta []map[string]interface{}
	vigencias := make(map[string]float64)

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/periodo?query=CodigoAbreviacion:VG,activo:true", &respuestaVigencias); err != nil {
		panic(err)
	}

	helpers.LimpiezaRespuestaRefactor(respuestaVigencias, &respuesta)
	for _, vigencia := range respuesta {
		idVigencia := strconv.FormatFloat(vigencia["Id"].(float64), 'f', -1, 64)
		vigencias[idVigencia] = vigencia["Year"].(float64)
	}

	return vigencias
}

func getUnidades() map[string]string {
	defer func() {
		if err := recover(); err != nil {
			panic(map[string]interface{}{"funcion": "getUnidades", "err": "Error obteniendo las unidades", "status": "400", "log": err})
		}
	}()
	unidades := make(map[string]string)
	var respuesta []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia?limit=0", &respuesta); err != nil {
		panic(err)
	}
	for _, u := range respuesta {
		idDependencia := strconv.FormatFloat(u["Id"].(float64), 'f', -1, 64)
		if u["Nombre"] == nil {
			unidades[idDependencia] = ""
		} else {
			unidades[idDependencia] = u["Nombre"].(string)
		}
	}
	return unidades
}

func getEstados() map[string]string {
	defer func() {
		if err := recover(); err != nil {
			panic(map[string]interface{}{"funcion": "getEstados", "err": "Error obteniendo los estados", "status": "400", "log": err})
		}
	}()
	estados := make(map[string]string)
	var respuestaEstados map[string]interface{}
	var estadoFormulacion []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-plan?query=activo:true", &respuestaEstados); err != nil {
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaEstados, &estadoFormulacion)
	for _, estado := range estadoFormulacion {
		estados[estado["_id"].(string)] = estado["nombre"].(string)
	}
	return estados
}

func obtenerNumeroVersion(planes []map[string]interface{}, planActual map[string]interface{}) int {
	versionesPlan := filtrarPlanes(planes, func(plan map[string]interface{}) bool {
		return plan["dependencia_id"] == planActual["dependencia_id"] && plan["vigencia"] == planActual["vigencia"] && plan["nombre"] == planActual["nombre"]
	})
	versionesOrdenadas := OrdenarVersiones(versionesPlan)

	return getNumVersion(versionesOrdenadas, func(plan map[string]interface{}) bool {
		return plan["_id"] == planActual["_id"]
	})

}

func getPlanesPorTipoPlan(codigoDeAbreviacion string) []map[string]interface{} {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			panic(map[string]interface{}{"funcion": "getPlanesPorTipoPlan", "err": localError["err"], "status": "400", "log": localError["log"]})
		}
	}()
	var respuestaTipoPlan map[string]interface{}
	var tipoPlan []map[string]interface{}
	var respuestaPlanes map[string]interface{}
	var planes []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-plan?query=codigo_abreviacion:"+codigoDeAbreviacion, &respuestaTipoPlan); err != nil {
		panic(map[string]interface{}{"log": "Error obteniendo el tipo de plan " + codigoDeAbreviacion, "err": err})
	}
	helpers.LimpiezaRespuestaRefactor(respuestaTipoPlan, &tipoPlan)
	// Obtener planes filtrados que sean formato y del tipo de plan especificado
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=formato:false"+",tipo_plan_id:"+tipoPlan[0]["_id"].(string) /*+"&limit=100"*/, &respuestaPlanes); err != nil {
		panic(map[string]interface{}{"log": "Error obteniendo los planes filtrados por el tipo de plan " + codigoDeAbreviacion, "err": err})
	}
	helpers.LimpiezaRespuestaRefactor(respuestaPlanes, &planes)

	return planes
}

func ObtenerPlanesFormulacion() []map[string]interface{} {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			beego.Debug(localError["err"])
			panic(map[string]interface{}{
				"funcion": "ObtenerPlanesFormulacion/" + localError["funcion"].(string),
				"err":     localError["err"],
				"status":  localError["status"],
			})
		}
	}()
	tiposPlanes := []string{CodigoTipoPlan, CodigoTipoPlanAccionFormulacion}
	var resumenPlanes []map[string]interface{}

	estados := getEstados()
	vigencias := getVigencias()
	unidades := getUnidades()
	// Obtener planes filtrados por el tipo de plan
	for _, tipoPlan := range tiposPlanes {
		planes := getPlanesPorTipoPlan(tipoPlan)
		for _, plan := range planes {
			if plan["dependencia_id"] != nil && plan["vigencia"] != nil {
				_, errD := strconv.Atoi(plan["dependencia_id"].(string))
				_, errV := strconv.Atoi(plan["vigencia"].(string))
				if errD == nil && errV == nil {
					planNuevo := make(map[string]interface{})
					planNuevo["id"] = plan["_id"]
					planNuevo["dependencia_id"] = plan["dependencia_id"]
					planNuevo["dependencia_nombre"] = unidades[plan["dependencia_id"].(string)]
					planNuevo["vigencia_id"] = plan["vigencia"]
					planNuevo["vigencia"] = vigencias[plan["vigencia"].(string)]
					planNuevo["nombre"] = plan["nombre"]
					planNuevo["version"] = obtenerNumeroVersion(planes, plan)
					planNuevo["estado_id"] = plan["estado_plan_id"]
					planNuevo["estado"] = estados[plan["estado_plan_id"].(string)]
					planNuevo["ultima_modificacion"] = plan["fecha_modificacion"]
					resumenPlanes = append(resumenPlanes, planNuevo)
				}
			}
		}
	}
	sort.Slice(resumenPlanes, func(i, j int) bool {
		return resumenPlanes[i]["ultima_modificacion"].(string) > resumenPlanes[j]["ultima_modificacion"].(string)
	})
	return resumenPlanes
}

func DefinirFechasFuncionamiento(body map[string]interface{}) []interface{} {
	var planesBody []map[string]interface{}

	planesBody, err := ObtenerArrayPlanesInteres(body)
	if err != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasFuncionamientoFormulacion", "err": "Error en la decodificación JSON", "status": "400"})
	}

	// Llamada a la función para codificar cada elemento del array de planes de interés
	planesJSON, err := CodificarPlanesInteres(planesBody)
	if err != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasFuncionamientoFormulacion", "err": "Error al codificar los planes de interés", "status": "400"})
	}

	body["planes_interes"] = ""
	var respuestasAcumuladas []interface{}
	for _, planJSON := range planesJSON { // Busca registros por cada _id de planes de interés y unidades
		var res map[string]interface{}
		var respuestaPost = make(map[string]interface{})
		var respuestaPost2 = make(map[string]interface{})
		body["planes_interes"] = "[" + planJSON + "]"
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/buscar-unidad-planes/1", "POST", &respuestaPost, body); err == nil {
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
					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+registro["_id"].(string), "PUT", &res, registro); err != nil {
						panic(map[string]interface{}{"funcion": "DefinirFechasFuncionamientoFormulacion", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
					}
					respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
				} else {
					// Registro con Fechas distintas, entonces se debe editar el registro existente quitando las unidades del body
					// Si unidades_interes del registro se queda vacio, se debe inactivar
					registro = manejarUnidades(body, registro, 2)
					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+registro["_id"].(string), "PUT", &res, registro); err != nil {
						panic(map[string]interface{}{"funcion": "DefinirFechasFuncionamientoFormulacion", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
					}
					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/buscar-unidad-planes/2", "POST", &respuestaPost2, body); err == nil {
						data, ok := respuestaPost2["Data"].([]interface{})
						if ok && len(data) > 0 { // Existe el registro, entonces se le agregan las unidades
							registro2 := data[0].(map[string]interface{})
							registro2 = manejarUnidades(body, registro2, 1)
							if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+registro2["_id"].(string), "PUT", &res, registro2); err != nil {
								panic(map[string]interface{}{"funcion": "DefinirFechasFuncionamientoFormulacion", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
							}
							respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
						} else { // No existe el registro, se crea el registro con las unidades del body
							body["planes_interes"] = registro["planes_interes"]
							if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento", "POST", &res, body); err != nil {
								panic(map[string]interface{}{"funcion": "DefinirFechasFuncionamientoFormulacion", "err": "Error agregaron periodo-seguimiento", "status": "400", "log": err})
							}
							respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
						}
					}

				}
			} else { // No existe el registro con las unidades especificadas, se procede a validar por periodo, fecha_inicio y fecha_fin
				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/buscar-unidad-planes/2", "POST", &respuestaPost2, body); err == nil {
					data, ok := respuestaPost2["Data"].([]interface{})
					if ok && len(data) > 0 { // Existe el registro, entonces se le agregan las unidades
						registro := data[0].(map[string]interface{})
						registro = manejarUnidades(body, registro, 1)
						if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+registro["_id"].(string), "PUT", &res, registro); err != nil {
							panic(map[string]interface{}{"funcion": "DefinirFechasFuncionamientoFormulacion", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
						}
						respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
					} else { // No existe el registro, se crea
						if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento", "POST", &res, body); err != nil {
							panic(map[string]interface{}{"funcion": "DefinirFechasFuncionamientoFormulacion", "err": "Error versionando periodo-seguimiento", "status": "400", "log": err})
						}
						respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
					}
				}
			}
		} else {
			panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionUnidadesPlanes", "err": "Error en la peticion", "status": "400", "log": err})
		}
	}
	return respuestasAcumuladas
}

func ObtenerArrayPlanesInteres(body map[string]interface{}) ([]map[string]interface{}, error) {
	// Acceder a la lista de planes de interés
	planesInteresString, ok := body["planes_interes"].(string)
	if !ok {
		return nil, fmt.Errorf("No se pudo obtener la lista de planes de interés como cadena JSON")
	}

	// Decodificar la cadena JSON a una lista de interfaces
	var planesInteres []interface{}
	err := json.Unmarshal([]byte(planesInteresString), &planesInteres)
	if err != nil {
		return nil, fmt.Errorf("Error al decodificar la lista de planes de interés: %v", err)
	}

	// Almacenar los planes de interés en un array
	var planes []map[string]interface{}
	for _, plan := range planesInteres {
		if planMap, ok := plan.(map[string]interface{}); ok {
			planes = append(planes, map[string]interface{}{
				"_id":    planMap["_id"],
				"nombre": planMap["nombre"],
			})
		}
	}

	return planes, nil
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

func CodificarPlanesInteres(planes []map[string]interface{}) ([]string, error) {
	// Codificar cada elemento del array de planes de interés a formato JSON
	var planesJSON []string
	for _, plan := range planes {
		planJSON, err := json.Marshal(plan)
		if err != nil {
			return nil, fmt.Errorf("Error al codificar un elemento del array de planes de interés: %v", err)
		}
		planesJSON = append(planesJSON, string(planJSON))
	}

	return planesJSON, nil
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

func eliminarUnidad(unidades []map[string]interface{}, id int) []map[string]interface{} {
	for i, unidad := range unidades {
		if unidad["Id"].(int) == id {
			// Eliminar la unidad del slice
			unidades = append(unidades[:i], unidades[i+1:]...)
			break
		}
	}
	return unidades
}