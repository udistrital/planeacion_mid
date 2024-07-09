package formulacionhelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/helpers/formatoHelper"
	"github.com/udistrital/planeacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

const (
	CodigoTipoPlanAccionFormulacion string = "PAF_SP"
	CodigoPlanEnFormulacion         string = "EF_SP"
	Pregrado                        string = "PREGRADO"
	Posgrado                        string = "POSGRADO"
	RHVPosgrado                     string = "RHVPOS"
	MedioTiempo                     string = "MTO"
	TiempoCompleto                  string = "TCO"
	HCatedraHonorarios              string = "HCH"
	HCatedraPrestacional            string = "HCP"
	NoAplica                        string = "N/A"
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
		hijo["ref"] = hijos[i]["_id"]

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
			outputError := map[string]interface{}{
				"funcion": "getVigencias",
				"error":   err,
				"mensaje": "No se lograron obtener las vigencias",
			}
			panic(outputError)
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
			outputError := map[string]interface{}{
				"funcion": "getUnidades",
				"error":   err,
				"mensaje": "No se lograron obtener las unidades",
			}
			panic(outputError)
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
			outputError := map[string]interface{}{
				"funcion": "getEstados",
				"error":   err,
				"mensaje": "No se lograron obtener los estados",
			}
			panic(outputError)
		}
	}()
	var respuestaEstados map[string]interface{}
	var estadoFormulacion []map[string]interface{}
	estados := make(map[string]string)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-plan?query=activo:true", &respuestaEstados); err != nil {
		// Mal manejo de error en la funcion de request.GetJson()
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaEstados, &estadoFormulacion)
	for _, estado := range estadoFormulacion {
		estados[estado["_id"].(string)] = estado["nombre"].(string)
	}
	return estados
}

func obtenerNumeroVersion(planes []map[string]interface{}, planActual map[string]interface{}) int {
	versionesPlan := helpers.FiltrarArreglo(planes, func(plan map[string]interface{}) bool {
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
			outputError := map[string]interface{}{
				"funcion": "getPlanesPorTipoPlan",
				"error":   err,
				"mensaje": "No se lograron obtener los planes filtrados por el tipo de plan " + codigoDeAbreviacion,
			}
			panic(outputError)
		}
	}()
	var respuestaTipoPlan map[string]interface{}
	var tipoPlan []map[string]interface{}
	var respuestaPlanes map[string]interface{}
	var planes []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-plan?query=codigo_abreviacion:"+codigoDeAbreviacion, &respuestaTipoPlan); err != nil {
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaTipoPlan, &tipoPlan)
	// Obtener planes filtrados que sean formato y del tipo de plan especificado
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=formato:false"+",tipo_plan_id:"+tipoPlan[0]["_id"].(string), &respuestaPlanes); err != nil {
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaPlanes, &planes)

	return planes
}

func ObtenerPlanesFormulacion() (resumenPlanes []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		var funcionDelError string
		var codigoEstado string
		if err := recover(); err != nil {

			localError := err.(map[string]interface{})
			_, existeFuncion := localError["funcion"]
			if existeFuncion {
				funcionDelError = "/" + localError["funcion"].(string)
			} else {
				funcionDelError = ""
			}

			_, existeCodigo := localError["status"]
			if existeCodigo {
				codigoEstado = localError["status"].(string)
			} else {
				codigoEstado = "400"
			}
			outputError := map[string]interface{}{
				"funcion": "ObtenerPlanesFormulacion" + funcionDelError,
				"err":     err,
				"status":  codigoEstado,
			}
			panic(outputError)
		}
	}()
	tiposPlanes := []string{CodigoTipoPlanAccionFormulacion}

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
					planNuevo["fecha_creacion"] = plan["fecha_creacion"]
					planNuevo["activo"] = plan["activo"]
					planNuevo["aplicativo_id"] = plan["aplicativo_id"]
					planNuevo["tipo_plan_id"] = plan["tipo_plan_id"]
					planNuevo["descripcion"] = plan["descripcion"]
					resumenPlanes = append(resumenPlanes, planNuevo)
				}
			}
		}
	}
	sort.Slice(resumenPlanes, func(i, j int) bool {
		return resumenPlanes[i]["ultima_modificacion"].(string) > resumenPlanes[j]["ultima_modificacion"].(string)
	})
	return resumenPlanes, outputError
}

// Función para realizar la petición POST hacia Resoluciones Docentes
// Obteniendo los valores de desagregado planeación
func GetDesagregado(bodyResolucionesDocente []map[string]interface{}) (map[string]interface{}, error) {
	var respuestaPost map[string]interface{}
	err := request.SendJson("http://"+beego.AppConfig.String("ResolucionesDocentes")+"/services/desagregado_planeacion", "POST", &respuestaPost, bodyResolucionesDocente)
	if err != nil || !respuestaPost["Success"].(bool) {
		return nil, err
	}
	return respuestaPost, nil
}

// Función para realizar la petición GET hacia Parametros Service
// Obteniendo el valor del salario mínimo para la vigencia
func GetSalarioMinimo(vigenciaStr string) (map[string]interface{}, error) {
	var resPeriodo map[string]interface{}
	var periodo []map[string]interface{}
	var resParametro map[string]interface{}
	var parametro []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Nombre:`+vigenciaStr, &resPeriodo); err != nil {
		return nil, err
	}
	helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
	vigenciaId := resPeriodo["Data"].([]interface{})[0].(map[string]interface{})["Id"].(float64)
	vigenciaIdStr := strconv.FormatFloat(vigenciaId, 'f', 0, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Activo:true,ParametroId:1,PeriodoId:"+vigenciaIdStr+"&limit=1", &resParametro); err != nil {
		return nil, err
	}
	helpers.LimpiezaRespuestaRefactor(resParametro, &parametro)

	valorParametro, valorParametroOk := resParametro["Data"].([]interface{})
	primerElemento, primerElementoOk := valorParametro[0].(map[string]interface{})
	salarioMinimo, salarioMinimoOk := primerElemento["Valor"].(string)
	var valorSalarioMinimo map[string]interface{}
	err := json.Unmarshal([]byte(salarioMinimo), &valorSalarioMinimo)

	if !valorParametroOk || len(valorParametro) == 0 || !primerElementoOk || !salarioMinimoOk || err != nil {
		return nil, fmt.Errorf("no se pudo obtener el valor de salario mínimo")
	}

	return valorSalarioMinimo, nil
}

// Obtener calculos para la Identificación de Docentes
func GetCalculos(data map[string]interface{}) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("los datos de entrada están vacíos")
	}

	response := map[string]interface{}{
		"totalHoras":               DataFinal(GetTotalHoras(data, false)),
		"totalHorasIndividual":     DataFinal(GetTotalHoras(data, true)),
		"meses":                    GetMeses(data),
		"sueldoBasico":             DataFinal(GetSueldoBasico(data, false)),
		"sueldoBasicoIndividual":   DataFinal(GetSueldoBasico(data, true)),
		"sueldoMensual":            DataFinal(GetSueldoMensual(data, false)),
		"sueldoMensualIndividual":  DataFinal(GetSueldoMensual(data, true)),
		"primaServicios":           DataFinal(GetPrimaServicios(data)),
		"primaNavidad":             DataFinal(GetPrimaNavidad(data)),
		"primaVacaciones":          DataFinal(GetPrimaVacaciones(data)),
		"vacaciones":               DataFinal(GetVacacionesProyeccion(data)),
		"bonificacion":             DataFinal(GetBonificacionServicios(data)),
		"interesesCesantias":       DataFinal(GetInteresesCesantias(data)),
		"cesantias":                DataFinal(GetCesantias(data)),
		"totalCesantias":           DataFinal(GetTotalAportesCesantias(data, false)),
		"totalCesantiasIndividual": DataFinal(GetTotalAportesCesantias(data, true)),
		"totalSalud":               DataFinal(GetTotalAporteSalud(data, false)),
		"totalSaludIndividual":     DataFinal(GetTotalAporteSalud(data, true)),
		"totalPensiones":           DataFinal(GetTotalAportePension(data, false)),
		"totalPensionesIndividual": DataFinal(GetTotalAportePension(data, true)),
		"totalArl":                 DataFinal(GetTotalArl(data, false)),
		"totalArlIndividual":       DataFinal(GetTotalArl(data, true)),
		"caja":                     DataFinal(GetCajaCompensacion(data)),
		"icbf":                     DataFinal(GetIcbf(data)),
		"totalBasico":              DataFinal(GetTotalSueldoBasico(data, false)),
		"totalBasicoIndividual":    DataFinal(GetTotalSueldoBasico(data, true)),
		"totalAportes":             DataFinal(GetTotalAportes(data, false)),
		"totalAportesIndividual":   DataFinal(GetTotalAportes(data, true)),
		"total":                    DataFinal(GetTotalRecurso(data, false)),
		"totalIndividual":          DataFinal(GetTotalRecurso(data, true)),
	}

	if response["cesantias"] == NoAplica {
		response["cesantiasPrivado"] = NoAplica
		response["cesantiasPublico"] = NoAplica
	}
	if response["totalPensiones"] == NoAplica {
		response["pensionesPrivado"] = NoAplica
		response["pensionesPublico"] = NoAplica
	}
	return response, nil
}

// Contruir cuerpo para la petición hacia Resoluciones Docentes
func ConstruirCuerpoRD(data map[string]interface{}) []map[string]interface{} {
	var bodyResolucionesDocente []map[string]interface{}

	resolucionDocente := make(map[string]interface{})
	resolucionDocente["Vigencia"] = data["vigencia"].(float64)
	resolucionDocente["Categoria"] = data["categoria"].(string)
	resolucionDocente["NivelAcademico"] = Pregrado

	if data["tipoDocente"].(string) == RHVPosgrado {
		resolucionDocente["NivelAcademico"] = Posgrado
	}

	tipoDedicacion := map[string]string{
		"Medio Tiempo":            MedioTiempo,
		"Tiempo Completo":         TiempoCompleto,
		"H. Catedra Prestacional": HCatedraPrestacional,
		"H. Catedra Honorarios":   HCatedraHonorarios,
	}
	dedicacion, existeTipo := tipoDedicacion[data["tipo"].(string)]
	if existeTipo {
		resolucionDocente["Dedicacion"] = dedicacion
	}

	if strings.Contains(data["categoria"].(string), "UD") {
		resolucionDocente["Categoria"] = strings.Replace(data["categoria"].(string), " UD", "", -1)
		resolucionDocente["EsDePlanta"] = true
	}

	bodyResolucionesDocente = append(bodyResolucionesDocente, resolucionDocente)

	return bodyResolucionesDocente
}

// Redondear valores de los calculos
func DataFinal(numeroStr string) string {
	if numeroStr != NoAplica {
		numeroDecimal, err := strconv.ParseFloat(numeroStr, 64)
		if err != nil {
			return ""
		}
		numeroRedondeado := math.Round(numeroDecimal)
		numeroRedondeadoStr := strconv.FormatFloat(numeroRedondeado, 'f', -1, 64)
		return numeroRedondeadoStr
	}
	return numeroStr
}

// Calcular el total de horas
func GetTotalHoras(data map[string]interface{}, ind bool) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	cantidad, cantidadOk := data["cantidad"].(float64)

	var resultado float64

	if semanasOk && horasOk && cantidadOk {
		resultado = cantidad * semanas * horas
	}
	if ind {
		return strconv.FormatFloat(resultado/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(resultado, 'f', -1, 64)
}

// Calcular el total de meses
func GetMeses(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)

	var resultado float64

	if semanasOk {
		resultado = semanas / 4
	}
	return strconv.FormatFloat(resultado, 'f', 2, 64)
}

// Calcular sueldo básico
func GetSueldoBasico(data map[string]interface{}, ind bool) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	cantidad, cantidadOk := data["cantidad"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var sueldoBasico float64

	if semanasOk && horasOk && cantidadOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		salarioBasico := resolucionDocente["salarioBasico"].(float64)
		sueldoBasico = cantidad * (salarioBasico * horas) * semanas * (1 + incremento)
	}
	if ind {
		return strconv.FormatFloat(sueldoBasico/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(sueldoBasico, 'f', -1, 64)
}

// Calcular sueldo mensual
func GetSueldoMensual(data map[string]interface{}, ind bool) string {
	cantidad, cantidadOk := data["cantidad"].(float64)

	var sueldoMensual float64

	if cantidadOk {
		sueldoBasico, errSB := strconv.ParseFloat(DataFinal(GetSueldoBasico(data, true)), 64)
		meses, errM := strconv.ParseFloat(GetMeses(data), 64)
		if errSB != nil || errM != nil {
			return ""
		}
		sueldoMensual = sueldoBasico / meses * cantidad
	}
	if ind {
		return strconv.FormatFloat(sueldoMensual/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(sueldoMensual, 'f', -1, 64)
}

// Calcular prima de servicios
func GetPrimaServicios(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var primaServicios float64

	if semanasOk && horasOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}
		meses, err := strconv.ParseFloat(GetMeses(data), 64)
		if err != nil {
			return ""
		}
		if meses < 6 {
			return "0"
		}
		prima_servicios, primaServiciosOk := resolucionDocente["prima_servicios"].(float64)
		if !primaServiciosOk {
			return ""
		}
		primaServicios = (prima_servicios * horas) * semanas * (1 + incremento)
	}
	return strconv.FormatFloat(primaServicios, 'f', -1, 64)
}

// Calcular prima de navidad
func GetPrimaNavidad(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var primaNavidad float64

	if semanasOk && horasOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}
		prima_navidad, prima_navidadOk := resolucionDocente["primaNavidad"].(float64)
		if !prima_navidadOk {
			return ""
		}
		primaNavidad = (prima_navidad * horas) * semanas * (1 + incremento)
	}
	return strconv.FormatFloat(primaNavidad, 'f', -1, 64)
}

// Calcular prima de vacaciones
func GetPrimaVacaciones(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var primaVacaciones float64

	if semanasOk && horasOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}
		prima_vacaciones, primaVacacionesOk := resolucionDocente["primaVacaciones"].(float64)
		if !primaVacacionesOk {
			return ""
		}
		primaVacaciones = (prima_vacaciones * horas) * semanas * (1 + incremento)
	}
	return strconv.FormatFloat(primaVacaciones, 'f', -1, 64)
}

// Calcular vacaciones proyección
func GetVacacionesProyeccion(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var vacacionesProyeccion float64

	if semanasOk && horasOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}

		vacaciones, vacacionesOk := resolucionDocente["vacaciones"].(float64)
		if !vacacionesOk {
			return ""
		}
		vacacionesProyeccion = (vacaciones * horas) * semanas * (1 + incremento)
	}
	return strconv.FormatFloat(vacacionesProyeccion, 'f', -1, 64)
}

// Calcular bonificación por servicios
func GetBonificacionServicios(data map[string]interface{}) string {
	sueldoBasico, errSB := strconv.ParseFloat(DataFinal(GetSueldoBasico(data, true)), 64)
	meses, errM := strconv.ParseFloat(GetMeses(data), 64)

	var resultado float64

	if errSB != nil || errM != nil {
		return ""
	}
	if meses < 12 {
		return NoAplica
	}
	resultado = (sueldoBasico * 0.35) / meses
	return strconv.FormatFloat(resultado, 'f', -1, 64)
}

// Calcular intereses cesantias
func GetInteresesCesantias(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var interesesCesantias float64

	if semanasOk && horasOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}
		interes_cesantias, interes_cesantiasOk := resolucionDocente["interesCesantias"].(float64)
		if !interes_cesantiasOk {
			return ""
		}
		interesesCesantias = (float64(interes_cesantias) * horas) * semanas * (1 + incremento)
	}
	return strconv.FormatFloat(interesesCesantias, 'f', -1, 64)
}

// Calcular cesantias
func GetCesantias(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var cesantias float64

	if semanasOk && horasOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica // Cesantias, CesantiasPrivado y CesantiasPublico
		}

		cesantias_, cesantiasOk := resolucionDocente["cesantias"].(float64)
		if !cesantiasOk {
			return ""
		}
		cesantias = (cesantias_ * horas) * semanas * (1 + incremento)
	}
	return strconv.FormatFloat(cesantias, 'f', -1, 64)
}

// Calcular total aportes a cesantias
func GetTotalAportesCesantias(data map[string]interface{}, ind bool) string {
	cantidad, cantidadOk := data["cantidad"].(float64)

	var totalAportesCesantias float64

	if cantidadOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}

		intereses, errI := strconv.ParseFloat(DataFinal(GetInteresesCesantias(data)), 64)
		cesantias, errC := strconv.ParseFloat(DataFinal(GetCesantias(data)), 64)
		if errI != nil || errC != nil {
			return ""
		}
		totalAportesCesantias = cantidad * (intereses + cesantias)
	}
	if ind {
		return strconv.FormatFloat(totalAportesCesantias/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(totalAportesCesantias, 'f', -1, 64)
}

func evaluarSaludPrestacional(infoPrestacional map[string]interface{}, data map[string]interface{}) float64 {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	salarioMinimo, salarioMinimoOk := data["salarioMinimo"].(float64)
	salarioBasico, salarioBasicoOk := infoPrestacional["salarioBasico"].(float64)

	var resultado float64

	if semanasOk && horasOk && salarioMinimoOk && salarioBasicoOk {
		if (salarioBasico * horas * 4) >= salarioMinimo {
			resultado = salarioBasico * horas * semanas * 0.085
		} else {
			resultado = (salarioMinimo * (semanas / 4) * 0.125) - (salarioBasico * horas * semanas * 0.04)
		}
	}
	return resultado
}

// Calcular total aporte salud
func GetTotalAporteSalud(data map[string]interface{}, ind bool) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	cantidad, cantidadOk := data["cantidad"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var totalAporteSalud float64

	if semanasOk && horasOk && cantidadOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}
		if dedicacion == HCatedraPrestacional {
			totalAporteSalud = cantidad * evaluarSaludPrestacional(resolucionDocente, data) * (1 + incremento)
		} else {
			salarioBasico := resolucionDocente["salarioBasico"].(float64)
			totalAporteSalud = cantidad * (((salarioBasico * horas) * semanas) * 0.085) * (1 + incremento)
		}
	}
	if ind {
		return strconv.FormatFloat(totalAporteSalud/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(totalAporteSalud, 'f', -1, 64)
}

func evaluarPensionPrestacional(infoPrestacional map[string]interface{}, data map[string]interface{}) float64 {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	salarioMinimo, salarioMinimoOk := data["salarioMinimo"].(float64)
	salarioBasico, salarioBasicoOk := infoPrestacional["salarioBasico"].(float64)
	pension, pensionOk := infoPrestacional["pension"].(float64)

	var resultado float64

	if semanasOk && horasOk && salarioMinimoOk && salarioBasicoOk && pensionOk {
		if (salarioBasico * horas * 4) >= salarioMinimo {
			resultado = (pension * horas * semanas * 0.12) / 0.16
		} else {
			resultado = (salarioMinimo * (semanas / 4) * 0.16) - (salarioBasico * horas * semanas * 0.04)
		}
	}
	return resultado
}

// Calcular total aporte pensión
func GetTotalAportePension(data map[string]interface{}, ind bool) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	cantidad, cantidadOk := data["cantidad"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var totalAportePension float64

	if semanasOk && horasOk && cantidadOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica // Pensiones, PensionesPrivado y PensionesPublico
		}

		if dedicacion == HCatedraPrestacional {
			totalAportePension = cantidad * evaluarPensionPrestacional(resolucionDocente, data) * (1 + incremento)
		} else {
			pension := resolucionDocente["pension"].(float64)
			totalAportePension = cantidad * ((((pension * horas) * semanas) * 0.12) / 0.16) * (1 + incremento)
		}
	}
	if ind {
		return strconv.FormatFloat(totalAportePension/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(totalAportePension, 'f', -1, 64)
}

func evaluarArlPrestacional(infoPrestacional map[string]interface{}, data map[string]interface{}) float64 {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	salarioMinimo, salarioMinimoOk := data["salarioMinimo"].(float64)
	salarioBasico, salarioBasicoOk := infoPrestacional["salarioBasico"].(float64)

	var resultado float64

	if semanasOk && horasOk && salarioMinimoOk && salarioBasicoOk {
		if (salarioBasico * horas * 4) >= salarioMinimo {
			resultado = salarioBasico * horas * semanas * 0.00522
		} else {
			resultado = salarioMinimo * (semanas / 4) * 0.00522
		}
	}
	return resultado
}

// Calcular total ARL
func GetTotalArl(data map[string]interface{}, ind bool) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	cantidad, cantidadOk := data["cantidad"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var totalArl float64

	if semanasOk && horasOk && cantidadOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}
		if dedicacion == HCatedraPrestacional {
			totalArl = cantidad * evaluarArlPrestacional(resolucionDocente, data) * (1 + incremento)
		} else {
			salarioBasico := resolucionDocente["salarioBasico"].(float64)
			totalArl = cantidad * (salarioBasico * horas) * semanas * 0.00522 * (1 + incremento)
		}
	}
	if ind {
		return strconv.FormatFloat(totalArl/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(totalArl, 'f', -1, 64)
}

func evaluarCajaPrestacional(infoPrestacional map[string]interface{}, data map[string]interface{}) float64 {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	salarioMinimo, salarioMinimoOk := data["salarioMinimo"].(float64)
	salarioBasico, salarioBasicoOk := infoPrestacional["salarioBasico"].(float64)
	primaVacaciones, primaVacacionesOk := infoPrestacional["primaVacaciones"].(float64)

	var resultado float64

	if semanasOk && horasOk && salarioMinimoOk && salarioBasicoOk && primaVacacionesOk {
		if (salarioBasico * horas * 4) >= salarioMinimo {
			resultado = (salarioBasico + primaVacaciones) * horas * semanas * 0.04
		} else {
			resultado = (salarioMinimo * (semanas / 4) * 0.04) + (primaVacaciones * horas * semanas * 0.04)
		}
	}
	return resultado
}

// Calcular caja de compensación
func GetCajaCompensacion(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var cajaCompensacion float64

	if semanasOk && horasOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}
		salarioBasico := resolucionDocente["salarioBasico"].(float64)
		primaVacaciones, primaVacacionesOk := resolucionDocente["primaVacaciones"].(float64)
		if !primaVacacionesOk {
			return ""
		}
		if dedicacion == HCatedraPrestacional {
			cajaCompensacion = evaluarCajaPrestacional(resolucionDocente, data) * (1 + incremento)
		} else {
			cajaCompensacion = (((salarioBasico + primaVacaciones) * horas) * semanas) * 0.04 * (1 + incremento)
		}
	}
	return strconv.FormatFloat(cajaCompensacion, 'f', -1, 64)
}

func evaluarIcbfPrestacional(infoPrestacional map[string]interface{}, data map[string]interface{}) float64 {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	salarioMinimo, salarioMinimoOk := data["salarioMinimo"].(float64)
	salarioBasico, salarioBasicoOk := infoPrestacional["salarioBasico"].(float64)

	var resultado float64

	if semanasOk && horasOk && salarioMinimoOk && salarioBasicoOk {
		if (salarioBasico * horas * 4) >= salarioMinimo {
			primaVacaciones, primaVacacionesOk := infoPrestacional["primaVacaciones"].(float64)
			if !primaVacacionesOk {
				return -1
			}
			resultado = (salarioBasico + primaVacaciones) * horas * semanas * 0.03
		} else {
			resultado = salarioMinimo * (semanas / 4) * 0.03
		}
	}
	return resultado
}

// Calcular ICBF
func GetIcbf(data map[string]interface{}) string {
	semanas, semanasOk := data["semanas"].(float64)
	horas, horasOk := data["horas"].(float64)
	incremento, incrementoOk := data["incremento"].(float64)

	var icbf float64

	if semanasOk && horasOk && incrementoOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		if dedicacion == HCatedraHonorarios {
			return NoAplica
		}

		if dedicacion == HCatedraPrestacional {
			resultado := evaluarIcbfPrestacional(resolucionDocente, data)
			if resultado == -1 {
				return ""
			}
			icbf = resultado * (1 + incremento)
		} else {
			salarioBasico := resolucionDocente["salarioBasico"].(float64)
			primaVacaciones, primaVacacionesOk := resolucionDocente["primaVacaciones"].(float64)
			if !primaVacacionesOk {
				return ""
			}
			icbf = (((salarioBasico + primaVacaciones) * horas) * semanas) * 0.03 * (1 + incremento)
		}
	}
	return strconv.FormatFloat(icbf, 'f', -1, 64)
}

// Calcular total sueldo básico
func GetTotalSueldoBasico(data map[string]interface{}, ind bool) string {
	cantidad, cantidadOk := data["cantidad"].(float64)

	var resultado float64

	if cantidadOk {
		sueldoBasico, errSB := strconv.ParseFloat(DataFinal(GetSueldoBasico(data, true)), 64)
		if errSB != nil {
			return ""
		}

		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		var totalBasico float64
		if dedicacion == HCatedraHonorarios {
			totalBasico = sueldoBasico
		} else {
			var bonificacionF float64
			if GetBonificacionServicios(data) != "" && GetBonificacionServicios(data) != NoAplica {
				bonificacion, errB := strconv.ParseFloat(DataFinal(GetBonificacionServicios(data)), 64)
				if errB != nil {
					return ""
				}
				bonificacionF = bonificacion
			}

			primaServicios, errPS := strconv.ParseFloat(DataFinal(GetPrimaServicios(data)), 64)
			primaNavidad, errPN := strconv.ParseFloat(DataFinal(GetPrimaNavidad(data)), 64)
			primaVacaciones, errPV := strconv.ParseFloat(DataFinal(GetPrimaVacaciones(data)), 64)
			aportesCesantias, errAC := strconv.ParseFloat(DataFinal(GetTotalAportesCesantias(data, false)), 64)
			vacaciones, errV := strconv.ParseFloat(DataFinal(GetVacacionesProyeccion(data)), 64)
			if errPS != nil || errPN != nil || errPV != nil || errAC != nil || errV != nil {
				return ""
			}

			aportesCesantiasF := aportesCesantias / data["cantidad"].(float64)
			totalBasico = sueldoBasico + primaServicios + primaNavidad + primaVacaciones + bonificacionF + aportesCesantiasF + vacaciones
		}
		resultado = cantidad * math.Round(totalBasico)
	}
	if ind {
		return strconv.FormatFloat(resultado/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(resultado, 'f', -1, 64)
}

// Calcular total aportes
func GetTotalAportes(data map[string]interface{}, ind bool) string {
	cantidad, cantidadOk := data["cantidad"].(float64)

	var resultado float64

	if cantidadOk {
		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		var totalAportes float64
		if dedicacion == HCatedraHonorarios {
			totalAportes = 0
		} else {
			totalSalud, errTS := strconv.ParseFloat(DataFinal(GetTotalAporteSalud(data, false)), 64)
			totalPension, errTP := strconv.ParseFloat(DataFinal(GetTotalAportePension(data, false)), 64)
			totalArl, errTA := strconv.ParseFloat(DataFinal(GetTotalArl(data, false)), 64)
			caja, errC := strconv.ParseFloat(DataFinal(GetCajaCompensacion(data)), 64)
			icbf, errICBF := strconv.ParseFloat(DataFinal(GetIcbf(data)), 64)
			if errTS != nil || errTP != nil || errTA != nil || errC != nil || errICBF != nil {
				return ""
			}

			totalSaludF := totalSalud / data["cantidad"].(float64)
			totalPensionF := totalPension / data["cantidad"].(float64)
			totalArlF := totalArl / data["cantidad"].(float64)
			totalAportes = totalSaludF + totalPensionF + totalArlF + caja + icbf
		}
		resultado = cantidad * math.Round(totalAportes)
	}
	if ind {
		return strconv.FormatFloat(resultado/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(resultado, 'f', -1, 64)
}

// Calcular total recurso
func GetTotalRecurso(data map[string]interface{}, ind bool) string {
	cantidad, cantidadOk := data["cantidad"].(float64)

	var resultado float64

	if cantidadOk {
		var bonificacionF float64
		if GetBonificacionServicios(data) != "" && GetBonificacionServicios(data) != NoAplica {
			bonificacion, errB := strconv.ParseFloat(DataFinal(GetBonificacionServicios(data)), 64)
			if errB != nil {
				return ""
			}
			bonificacionF = bonificacion
		}

		sueldoBasico, errSB := strconv.ParseFloat(DataFinal(GetSueldoBasico(data, true)), 64)
		if errSB != nil {
			return ""
		}

		resolucionDocente := data["resolucionDocente"].(map[string]interface{})
		dedicacion := resolucionDocente["Dedicacion"].(string)

		var total float64
		if dedicacion == HCatedraHonorarios {
			total = sueldoBasico
		} else {
			aportesCesantias, errAC := strconv.ParseFloat(DataFinal(GetTotalAportesCesantias(data, false)), 64)
			primaServicios, errPS := strconv.ParseFloat(DataFinal(GetPrimaServicios(data)), 64)
			primaNavidad, errPN := strconv.ParseFloat(DataFinal(GetPrimaNavidad(data)), 64)
			primaVacaciones, errPV := strconv.ParseFloat(DataFinal(GetPrimaVacaciones(data)), 64)
			vacaciones, errV := strconv.ParseFloat(DataFinal(GetVacacionesProyeccion(data)), 64)
			totalSalud, errTS := strconv.ParseFloat(DataFinal(GetTotalAporteSalud(data, false)), 64)
			totalPension, errTP := strconv.ParseFloat(DataFinal(GetTotalAportePension(data, false)), 64)
			totalArl, errTA := strconv.ParseFloat(DataFinal(GetTotalArl(data, false)), 64)
			caja, errC := strconv.ParseFloat(DataFinal(GetCajaCompensacion(data)), 64)
			icbf, errICBF := strconv.ParseFloat(DataFinal(GetIcbf(data)), 64)

			if errPS != nil || errPN != nil || errPV != nil || errAC != nil || errV != nil || errTS != nil || errTP != nil || errTA != nil || errC != nil || errICBF != nil {
				return ""
			}

			aportesCesantiasF := aportesCesantias / data["cantidad"].(float64)
			totalSaludF := totalSalud / data["cantidad"].(float64)
			totalPensionF := totalPension / data["cantidad"].(float64)
			totalArlF := totalArl / data["cantidad"].(float64)

			total = sueldoBasico + primaServicios + primaNavidad + primaVacaciones + bonificacionF + aportesCesantiasF + totalSaludF + totalArlF + caja + icbf + vacaciones + totalPensionF
		}
		resultado = cantidad * math.Round(total)
	}
	if ind {
		return strconv.FormatFloat(resultado/cantidad, 'f', -1, 64)
	}
	return strconv.FormatFloat(resultado, 'f', -1, 64)
}

// Obtener una plantilla por id
func GetPlantilla(id string) (map[string]interface{}, error) {
	var resPlantilla map[string]interface{}
	var plantilla []map[string]interface{}

	err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=formato:true,_id:"+id, &resPlantilla)
	if err != nil {
		return nil, err
	}
	helpers.LimpiezaRespuestaRefactor(resPlantilla, &plantilla)
	return plantilla[0], nil
}

// Obtener el id de acuerdo al estado:"En formulación" por codigo de abreviación
func getIdEstadoEnFormulacion() (string, error) {
	var resEstado map[string]interface{}
	var estado []map[string]interface{}
	err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-plan?query=activo:true,codigo_abreviacion:"+CodigoPlanEnFormulacion, &resEstado)
	if err != nil {
		return "", err
	}
	helpers.LimpiezaRespuestaRefactor(resEstado, &estado)
	return estado[0]["_id"].(string), nil
}

// Obtener los planes en estado:"En formulación" por nombre de plantilla
func GetPlanesPorNombre(nombre string) ([]map[string]interface{}, error) {
	var resPlanes map[string]interface{}
	var planes []map[string]interface{}
	idEstado, err1 := getIdEstadoEnFormulacion()
	if err1 != nil {
		return nil, err1
	}
	err2 := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=formato:false,nombre:"+url.QueryEscape(nombre)+",estado_plan_id:"+idEstado, &resPlanes)
	if err2 != nil {
		return nil, err2
	}
	helpers.LimpiezaRespuestaRefactor(resPlanes, &planes)
	return planes, nil
}

// Obtener el formato de una plantilla o un plan
func GetFormato(id string) ([][]map[string]interface{}, error) {
	var res map[string]interface{}
	var hijos []models.Nodo
	var plan map[string]interface{}
	var hijosID []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		helpers.LimpiezaRespuestaRefactor(res, &hijosID)
		err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res)
		if err != nil {
			return nil, err
		}
		helpers.LimpiezaRespuestaRefactor(res, &plan)
		formatoHelper.Limpia(plan)
		tree := formatoHelper.BuildTreeFaActEst(hijos, hijosID)
		return tree, nil
	} else {
		return nil, err
	}
}

// Obtener un subgrupo por id
func getSubgrupo(id string) (map[string]interface{}, error) {
	var resSubgrupo map[string]interface{}
	var subgrupo map[string]interface{}
	err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+id, &resSubgrupo)
	if err != nil {
		return nil, err
	}
	helpers.LimpiezaRespuestaRefactor(resSubgrupo, &subgrupo)
	return subgrupo, nil
}

// Crear un subgrupo (registar_nodo)
func crearSubgrupo(bodySubgrupo map[string]interface{}) (map[string]interface{}, error) {
	var resSubgrupo map[string]interface{}
	err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/registrar_nodo", "POST", &resSubgrupo, bodySubgrupo)
	if err != nil || !resSubgrupo["Success"].(bool) {
		return nil, err
	}
	return resSubgrupo, nil
}

// Actualizar un subgrupo
func actualizarSubgrupo(bodySubgrupo map[string]interface{}, id string) (map[string]interface{}, error) {
	var resSubgrupo map[string]interface{}
	err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+id, "PUT", &resSubgrupo, bodySubgrupo)
	if err != nil {
		return nil, err
	}
	return resSubgrupo, nil
}

// Obtener el detalle de un subgrupo por id
func getSubgrupoDetalle(id string) (map[string]interface{}, error) {
	var resSubgrupoDetalle map[string]interface{}
	var subgrupoDetalle []map[string]interface{}
	err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+id, &resSubgrupoDetalle)
	if err != nil {
		return nil, err
	}
	helpers.LimpiezaRespuestaRefactor(resSubgrupoDetalle, &subgrupoDetalle)
	return subgrupoDetalle[0], nil
}

// Crear el detalle de un subgrupo
func crearSubgrupoDetalle(bodySubgrupoDetalle map[string]interface{}) (map[string]interface{}, error) {
	var resSubgrupoDetalle map[string]interface{}
	err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/", "POST", &resSubgrupoDetalle, bodySubgrupoDetalle)
	if err != nil || !resSubgrupoDetalle["Success"].(bool) {
		return nil, err
	}
	return resSubgrupoDetalle, nil
}

// Actualizar el detalle de un subgrupo
func actualizarSubgrupoDetalle(bodySubgrupoDetalle map[string]interface{}, id string) (map[string]interface{}, error) {
	var resSubgrupoDetalle map[string]interface{}
	err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+id, "PUT", &resSubgrupoDetalle, bodySubgrupoDetalle)
	if err != nil {
		return nil, err
	}
	return resSubgrupoDetalle, nil
}

// Convertir arbol a lista plana
func ConvArbolAListaPlana(lista []map[string]interface{}, id string, isFormato bool) ([]map[string]interface{}, error) {
	var listaPlana []map[string]interface{}
	for _, objeto := range lista {

		var referencia interface{}
		if !isFormato {
			ref, refOk := objeto["ref"].(string)
			if !refOk {
				return nil, fmt.Errorf("no se pudo obtener el valor de referencia")
			}
			referencia = ref
		} else {
			referencia = objeto["ref"]
		}

		nuevoMapa := map[string]interface{}{
			"id":       objeto["id"],
			"nombre":   objeto["nombre"],
			"required": objeto["required"],
			"type":     objeto["type"],
			"ref":      referencia,
			"padre":    id,
		}
		if objeto["type"] == "select" {
			nuevoMapa["options"] = objeto["options"]
		}
		listaPlana = append(listaPlana, nuevoMapa)

		if sub, ok := objeto["sub"].([]map[string]interface{}); ok && len(sub) > 0 {
			subListaPlana, err := ConvArbolAListaPlana(sub, objeto["id"].(string), isFormato)
			if err != nil {
				return nil, err
			}
			listaPlana = append(listaPlana, subListaPlana...)
		}
	}
	return listaPlana, nil
}

// Actualizar la estructura de los planes en base a una plantilla
// Comparando los formatos de la plantilla y el plan
func ActualizarEstructuraPlan(listaFormato, listaPlan []map[string]interface{}, idPlan string) error {
	for _, objeto1 := range listaFormato {
		var auxObjeto map[string]interface{}
		encontrado1 := false

		for _, objeto2 := range listaPlan {
			refObjeto2, refObjeto2Ok := objeto2["ref"].(string)
			if !refObjeto2Ok {
				return fmt.Errorf("no se pudo obtener el valor de referencia")
			}
			if objeto1["id"] == refObjeto2 {
				encontrado1 = true
				auxObjeto = objeto2
				break
			}
		}

		if encontrado1 {
			idNodoPlantilla := objeto1["id"].(string)
			idNodoPlan := auxObjeto["id"].(string)

			//Obtener subgrupo de la plantilla y el plan
			subgFormato, err := getSubgrupo(idNodoPlantilla)
			if err != nil {
				return err
			}
			subgPlan, err := getSubgrupo(idNodoPlan)
			if err != nil {
				return err
			}

			refSubgPlan, refSubgPlanOk := subgPlan["ref"].(string)
			if !refSubgPlanOk {
				return fmt.Errorf("no se pudo obtener el valor de referencia")
			}
			if !compararSubgrupos(subgFormato, subgPlan) {
				actSubgrupo := map[string]interface{}{
					"nombre":         subgFormato["nombre"],
					"descripcion":    subgFormato["descripcion"],
					"activo":         subgFormato["activo"],
					"bandera_tabla":  subgFormato["bandera_tabla"],
					"padre":          subgPlan["padre"],
					"ref":            refSubgPlan,
					"fecha_creacion": subgPlan["fecha_creacion"],
				}
				_, err := actualizarSubgrupo(actSubgrupo, idNodoPlan)
				if err != nil {
					return err
				}
			}

			//Obtener subgrupo-detalle de la plantilla y el plan
			subDetaFormato, err := getSubgrupoDetalle(idNodoPlantilla)
			if err != nil {
				return err
			}
			subDetaPlan, err := getSubgrupoDetalle(idNodoPlan)
			if err != nil {
				return err
			}
			if !compararSubgruposDetalle(subDetaFormato, subDetaPlan) {
				actSubgrupoDetalle := map[string]interface{}{
					"dato":           subDetaFormato["dato"],
					"fecha_creacion": subDetaPlan["fecha_creacion"],
					"subgrupo_id":    subDetaPlan["subgrupo_id"],
					"activo":         subDetaFormato["activo"],
					"descripcion":    subDetaFormato["descripcion"],
					"nombre":         subDetaFormato["nombre"],
					"dato_plan":      "",
				}
				_, err := actualizarSubgrupoDetalle(actSubgrupoDetalle, subDetaPlan["_id"].(string))
				if err != nil {
					return err
				}
			}
		} else {
			var padreID string
			encontrado2 := false

			for _, objeto := range listaPlan {
				refObjeto, refObjetoOk := objeto["ref"].(string)
				if !refObjetoOk {
					return fmt.Errorf("no se pudo obtener el valor de referencia")
				}
				if objeto1["padre"] == refObjeto {
					encontrado2 = true
					padreID = objeto["id"].(string)
					break
				}
			}

			if !encontrado2 {
				padreID = idPlan
			}

			// Realizar una copia del subgrupo del formato para el plan
			subFormato, err1 := getSubgrupo(objeto1["id"].(string))
			if err1 != nil {
				return err1
			}
			nuevoSubgrupo := map[string]interface{}{
				"nombre":        subFormato["nombre"],
				"descripcion":   subFormato["descripcion"],
				"padre":         padreID,
				"activo":        subFormato["activo"],
				"bandera_tabla": subFormato["bandera_tabla"],
				"ref":           subFormato["_id"],
			}
			resSubgrupo, err2 := crearSubgrupo(nuevoSubgrupo)
			if err2 != nil {
				return err2
			}
			resSubgrupoData := resSubgrupo["Data"].(map[string]interface{})

			// Realizar una copia del subgrupo-detalle del formato para el plan
			subDetFormato, err3 := getSubgrupoDetalle(objeto1["id"].(string))
			if err3 != nil {
				return err3
			}
			nuevoSubgrupoDetalle := map[string]interface{}{
				"nombre":      subDetFormato["nombre"],
				"descripcion": subDetFormato["descripcion"],
				"subgrupo_id": resSubgrupoData["_id"],
				"dato":        subDetFormato["dato"],
				"activo":      subDetFormato["activo"],
			}
			_, err4 := crearSubgrupoDetalle(nuevoSubgrupoDetalle)
			if err4 != nil {
				return err4
			}

			formatoPlanAct, err5 := GetFormato(idPlan)
			if err5 != nil {
				return err5
			}
			nuevaLista, err := ConvArbolAListaPlana(formatoPlanAct[0], idPlan, false)
			if err != nil {
				return err
			}
			listaPlan = nuevaLista
		}
	}
	return nil
}

// Comparar los subgrupos de una plantilla con un plan
func compararSubgrupos(objeto1, objeto2 map[string]interface{}) bool {
	if objeto1["nombre"] != objeto2["nombre"] {
		return false
	}
	if objeto1["descripcion"] != objeto2["descripcion"] {
		return false
	}
	if objeto1["activo"] != objeto2["activo"] {
		return false
	}
	if objeto1["bandera_tabla"] != objeto2["bandera_tabla"] {
		return false
	}
	return true
}

// Comparar el detalle de los subgrupos de una plantilla con un plan
func compararSubgruposDetalle(objeto1, objeto2 map[string]interface{}) bool {
	if objeto1["nombre"] != objeto2["nombre"] {
		return false
	}
	if objeto1["descripcion"] != objeto2["descripcion"] {
		return false
	}
	if objeto1["dato"] != objeto2["dato"] {
		return false
	}
	if objeto1["activo"] != objeto2["activo"] {
		return false
	}
	return true
}

func DefinirFechasFormulacionSeguimiento(body map[string]interface{}) []interface{} {
	var planesBody []map[string]interface{}

	planesBody, err := ObtenerArrayPlanesInteres(body)
	if err != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error en la decodificación JSON", "status": "400"})
	}

	// Llamada a la función para codificar cada elemento del array de planes de interés
	planesJSON, err := CodificarPlanesInteres(planesBody)
	if err != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error al codificar los planes de interés", "status": "400"})
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
						panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
					}
					respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
				} else {
					// Registro con Fechas distintas, entonces se debe editar el registro existente quitando las unidades del body
					// Si unidades_interes del registro se queda vacio, se debe inactivar
					registro = manejarUnidades(body, registro, 2)
					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+registro["_id"].(string), "PUT", &res, registro); err != nil {
						panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
					}
					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/buscar-unidad-planes/2", "POST", &respuestaPost2, body); err == nil {
						data, ok := respuestaPost2["Data"].([]interface{})
						if ok && len(data) > 0 { // Existe el registro, entonces se le agregan las unidades
							registro2 := data[0].(map[string]interface{})
							registro2 = manejarUnidades(body, registro2, 1)
							if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+registro2["_id"].(string), "PUT", &res, registro2); err != nil {
								panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
							}
							respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
						} else { // No existe el registro, se crea el registro con las unidades del body
							body["planes_interes"] = registro["planes_interes"]
							if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento", "POST", &res, body); err != nil {
								panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error agregaron periodo-seguimiento", "status": "400", "log": err})
							}
							respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
						}

						//? Actualiza registros de seguimiento en caso de que hayan planes avalados con el formato_id igual al id del plan_interes del body
						if body["tipo_seguimiento_id"] == "61f236f525e40c582a0840d0" || body["tipo_seguimiento_id"] == "6385fa136a0d19d7888837ed" {
							periodoSeguimientoIdAntiguo := registro["_id"].(string)
							periodoSeguimientoIdNuevo := res["Data"].(map[string]interface{})["_id"].(string)
							unidadesBody, err := ObtenerArrayUnidadesInteres(body)
							if err != nil {
								panic(map[string]interface{}{"funcion": "manejarUnidades", "err": "Error en la decodificación JSON", "status": "400"})
							}
							CambiarFechasSeguimiento(planJSON, unidadesBody, periodoSeguimientoIdAntiguo, periodoSeguimientoIdNuevo)
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
							panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error actualizando periodo-seguimiento \"registro[\"_id\"].(string)\"", "status": "400", "log": err})
						}
						respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
					} else { // No existe el registro, se crea
						if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento", "POST", &res, body); err != nil {
							panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error versionando periodo-seguimiento", "status": "400", "log": err})
						}
						respuestasAcumuladas = append(respuestasAcumuladas, res["Data"])
					}
				}
			}
		} else {
			panic(map[string]interface{}{"funcion": "DefinirFechasFormulacionSeguimiento", "err": "Error en la peticion", "status": "400", "log": err})
		}
	}
	return respuestasAcumuladas
}

// ? Función para cambiar el periodo_seguimiento_id de los registros de seguimiento de los planes avalados
func CambiarFechasSeguimiento(planInteresString string, unidadesInteres []map[string]interface{}, periodoSeguimientoIdAntiguo string, periodoSeguimientoIdNuevo string) {
	var resPlanesAvalados map[string]interface{}
	var resSeguimientos map[string]interface{}
	var planesAvalados []map[string]interface{}
	var seguimiento []map[string]interface{}
	var seguimientoActualizado map[string]interface{}
	var estadoAvaladoId = "6153355601c7a2365b2fb2a1"
	var planInteres models.PlanInteres

	err := json.Unmarshal([]byte(planInteresString), &planInteres)
	if err != nil {
		panic(map[string]interface{}{"funcion": "CambiarFechasSeguimiento", "err": "Error en la decodificación JSON", "status": "400"})
	}

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
		}
	}
	return unidades
}

func ValidarUnidadesPlanes(periodo_seguimiento map[string]interface{}, body_unidades map[string]interface{}) []map[string]interface{} {
	unidadesJSON1 := []interface{}{}
	if err := json.Unmarshal([]byte(periodo_seguimiento["unidades_interes"].(string)), &unidadesJSON1); err != nil {
		panic(map[string]interface{}{"funcion": "ValidarUnidadesPlanes", "err": "Error al deserializar unidadesJSON1", "status": "400", "log": err})
	}

	unidadesJSON2 := body_unidades["unidades_interes"].([]interface{})

	unidadesMap := make(map[int]bool)
	for _, unidad := range unidadesJSON1 {
		unidadMap := unidad.(map[string]interface{})
		id := int(unidadMap["Id"].(float64))
		unidadesMap[id] = true
	}

	var unidadesValidadas []map[string]interface{}

	for _, unidad := range unidadesJSON2 {
		unidadMap := unidad.(map[string]interface{})
		id := int(unidadMap["Id"].(float64))
		if !unidadesMap[id] {
			// Si hay una unidad en el body que no está en periodo_seguimiento, retornar null
			fmt.Println("Unidad en body no presente en periodo_seguimiento: ", unidadMap)
			return nil
		}
		unidadesValidadas = append(unidadesValidadas, unidadMap)
	}

	// Retorna la intersección de las unidades
	return unidadesValidadas
}

func ObtenerIdParametros() (float64, float64, float64, error) {
	var ParametroNoRegistra map[string]interface{}
	var ParametroJefeOficina map[string]interface{}
	var ParametroAsistenteDependencia map[string]interface{}
	baseURL := "http://" + beego.AppConfig.String("ParametrosService") + "/parametro?query="

	err := request.GetJson(baseURL+"CodigoAbreviacion:NR,TipoParametroId__CodigoAbreviacion:C,Activo:true", &ParametroNoRegistra)
	IdNoRegistra, ok := ParametroNoRegistra["Data"].([]interface{})[0].(map[string]interface{})["Id"].(float64)

	if err != nil || !ok {
		panic(map[string]interface{}{"funcion": "VinculacionTercero", "err": "Error get ParametroNoRegistra", "status": "400", "log": err})
	}

	err = request.GetJson(baseURL+"CodigoAbreviacion:JO,TipoParametroId__CodigoAbreviacion:C,Activo:true", &ParametroJefeOficina)
	IdJefeOficina, _ := ParametroJefeOficina["Data"].([]interface{})[0].(map[string]interface{})["Id"].(float64)
	if err != nil || IdJefeOficina == 0 {
		panic(map[string]interface{}{"funcion": "VinculacionTercero", "err": "Error get ParametroJefeOficina", "status": "400", "log": err})
	}

	err = request.GetJson(baseURL+"CodigoAbreviacion:AS_D,TipoParametroId__CodigoAbreviacion:C,Activo:true", &ParametroAsistenteDependencia)
	IdAsistenteDependencia, _ := ParametroAsistenteDependencia["Data"].([]interface{})[0].(map[string]interface{})["Id"].(float64)
	if err != nil || IdAsistenteDependencia == 0 {
		panic(map[string]interface{}{"funcion": "VinculacionTercero", "err": "Error get ParametroAsistenteDependencia", "status": "400", "log": err})
	}

	return IdNoRegistra, IdJefeOficina, IdAsistenteDependencia, nil
}

func obtenerCorreoPlaneacion() (string, error) {
	var respuestaPeticion map[string]interface{}
	var correoPlaneacion map[string]interface{}
	baseURL := "http://" + beego.AppConfig.String("ParametrosService") + "/parametro_periodo?query="
	err := request.GetJson(baseURL+"ParametroId.CodigoAbreviacion:CORREO_OAP,ParametroId.TipoParametroId.CodigoAbreviacion:P_SISGPLAN,Activo:true", &respuestaPeticion)
	jsonData, ok := respuestaPeticion["Data"].([]interface{})[0].(map[string]interface{})["Valor"].(string)

	if err != nil || !ok {
		return "", fmt.Errorf("no se pudo obtener el correo de la Oficina de Asesora de Planeación del módulo de parámetros, comuníquese con computo@udistrital.edu.co")
	}
	err1 := json.Unmarshal([]byte(jsonData), &correoPlaneacion)
	if err1 != nil {
		return "", fmt.Errorf("Error al deserializar el JSON de Correo Oficina Planeacion: ", err)
		}
	return correoPlaneacion["Valor"].(string), nil
}

func CambioCargoIdVinculacionTercero(idVinculacion string, body map[string]interface{}) (*models.Vinculacion, error) {
	const ROL_ASISTENTE_PLANEACION = "ASISTENTE_PLANEACION"
	var vinculacionPlaneacion, rolAsistentePlaneacion bool
	var vinculacion []models.Vinculacion

	correoPlaneacion, errorPeticion := obtenerCorreoPlaneacion()
	if errorPeticion != nil {
		panic(map[string]interface{}{"funcion": "CambioCargoIdVinculacionTercero", "err": errorPeticion.Error(), "status": "400", "log": errorPeticion})
	}

	idNoRegistra, idJefeOficina, idAsistenteDependencia, err := ObtenerIdParametros()
	if err != nil {
		panic(map[string]interface{}{"funcion": "VinculacionTercero", "err": "Error get parametros", "status": "400", "log": err})
	}

	err = request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/vinculacion?query=Activo:true,Id:"+idVinculacion, &vinculacion)
	if err != nil || vinculacion[0].CargoId == 0 {
		panic(map[string]interface{}{"funcion": "CambioCargoIdVinculacionTercero", "err": "Error get vinculacion", "status": "400", "log": err})
	}

	vinculaciones := body["user"].(map[string]interface{})["Vinculacion"].([]interface{})
	vinculacionSeleccionada := body["user"].(map[string]interface{})["VinculacionSeleccionadaId"]
	
	for _, vinculacion := range vinculaciones {
		if vinculacion.(map[string]interface{})["Id"].(float64) == vinculacionSeleccionada {
			if DependenciaCorreo, ok := vinculacion.(map[string]interface{})["DependenciaCorreo"].(string); ok {
        vinculacionPlaneacion = DependenciaCorreo == correoPlaneacion
				rolAsistentePlaneacion = (body["rol"].(string) == ROL_ASISTENTE_PLANEACION)
				if rolAsistentePlaneacion && !vinculacionPlaneacion && body["vincular"] == true {
					return nil, fmt.Errorf("El usuario no puede tener el rol de %s si no pertenece a la Oficina de Asesora de Planeación", ROL_ASISTENTE_PLANEACION)
				}
			} else {
				return nil, fmt.Errorf("no se pudo obtener la Dependencia asociada a la vinculación")
			}
		}
	}

	if vinculacion[0].CargoId == int(idJefeOficina) || vinculacion[0].CargoId == int(idAsistenteDependencia) || vinculacion[0].CargoId == int(idNoRegistra) {
		if body["vincular"] == true {
			vinculacion[0].CargoId = int(idAsistenteDependencia)
		} else {
			vinculacion[0].CargoId = int(idNoRegistra)
		}
		vinculacion[0].FechaCreacion = formatearFecha(vinculacion[0].FechaCreacion)
		vinculacion[0].FechaModificacion = time.Now().Format(time.RFC3339)
		if err := helpers.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/vinculacion/"+idVinculacion, "PUT", &vinculacion[0], vinculacion[0]); err != nil {
			panic(map[string]interface{}{"funcion": "CambioCargoIdVinculacionTercero", "err": "Error actualizando vinculacion", "status": "400", "log": err})
		}
		return &vinculacion[0], nil
	}

	return nil, errors.New("No se encontró la vinculación")
}

func formatearFecha(fecha string) string {
	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700 -0700", fecha)
	if err != nil {
		fmt.Println("Error parseando fecha:", err)
		return ""
	}
	// Formatear la fecha al nuevo formato
	return parsedTime.Format(time.RFC3339)
}
