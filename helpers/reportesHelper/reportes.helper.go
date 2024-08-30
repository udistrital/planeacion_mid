package reporteshelper

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"log"
	// "strings"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/xuri/excelize/v2"
)

var validDataT = []string{}
var hijos_key []interface{}
var hijos_data [][]map[string]interface{}
var detalles []map[string]interface{}
var detalles_armonizacion map[string]interface{}
var ids [][]string
var id_arr []string
var detallesLlenados bool

func LimpiarDetalles() {
	detalles = []map[string]interface{}{}
	detalles_armonizacion = map[string]interface{}{}
	detallesLlenados = false
}

func Limpia() {
	//validDataT = []string{}
	//ids = [][]string{}
	//hijos_data = nil
	//hijos_key = nil
}

func LimpiaIds() {
	id_arr = []string{}
}

func Limp() {
	validDataT = []string{}
	ids = [][]string{}
	hijos_data = nil
	hijos_key = nil
}

func GetActividades(subgrupo_id string) []map[string]interface{} {
	var res map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var datoPlan map[string]interface{}
	var actividades []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+subgrupo_id+"&fields=dato_plan", &res); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(res, &aux)
		subgrupoDetalle = aux[0]
		if subgrupoDetalle["dato_plan"] != nil {
			dato_plan_str := subgrupoDetalle["dato_plan"].(string)
			json.Unmarshal([]byte(dato_plan_str), &datoPlan)
			for indexActividad, element := range datoPlan {
				_ = indexActividad
				if err != nil {
					log.Panic(err)
				}
				if element.(map[string]interface{})["activo"] == true {
					actividades = append(actividades, element.(map[string]interface{}))
				}
			}

		}
	} else {
		panic(map[string]interface{}{"Code": "400", "Body": err, "Type": "error"})

	}
	return actividades
}

func BuildTreeFa(hijos []map[string]interface{}, index string) [][]map[string]interface{} {
	var tree []map[string]interface{}
	var requeridos []map[string]interface{}
	armonizacion := make([]map[string]interface{}, 1)
	var result [][]map[string]interface{}
	for i := 0; i < len(hijos); i++ {
		if hijos[i]["activo"] == true {
			forkData := make(map[string]interface{})
			var id string
			forkData["id"] = hijos[i]["_id"]
			forkData["nombre"] = hijos[i]["nombre"]
			id = hijos[i]["_id"].(string)

			if len(hijos[i]["hijos"].([]interface{})) > 0 {
				var aux []map[string]interface{}
				if len(hijos_key) == 0 {
					hijos_key = append(hijos_key, hijos[i]["hijos"])
					hijos_data = append(hijos_data, getChildren(hijos[i]["hijos"].([]interface{}), true))
					aux = hijos_data[len(hijos_data)-1]
				} else {
					flag := false
					var posicion int
					for j := 0; j < len(hijos_key); j++ {
						if reflect.DeepEqual(hijos[i]["hijos"], hijos_key[j]) {
							flag = true
							posicion = j
						}
					}
					if !flag {
						hijos_key = append(hijos_key, hijos[i]["hijos"])
						hijos_data = append(hijos_data, getChildren(hijos[i]["hijos"].([]interface{}), true))
						aux = hijos_data[len(hijos_data)-1]
					} else {
						aux = hijos_data[posicion]
						for k := 0; k < len(ids[posicion]); k++ {
							add(ids[posicion][k])
						}
					}
				}

				forkData["sub"] = make([]map[string]interface{}, len(aux))
				forkData["sub"] = aux
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

	LimpiaIds()

	return result
}

func add(id string) {
	if !contains(validDataT, id) {
		validDataT = append(validDataT, id)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func convert(valid []string, index string) ([]map[string]interface{}, map[string]interface{}) {
	var validadores []map[string]interface{}
	var actividad map[string]interface{}
	var dato_armonizacion map[string]interface{}
	armonizacion := make(map[string]interface{})
	forkData := make(map[string]interface{})
	for i, v := range valid {
		var res map[string]interface{}
		var subgrupo_detalle []map[string]interface{}
		var dato_plan map[string]interface{}

		if !detallesLlenados {
			detalles = append(detalles, map[string]interface{}{})
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+v+"&fields=dato_plan,armonizacion_dato", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &subgrupo_detalle)

				if len(subgrupo_detalle) > 0 {
					if subgrupo_detalle[0]["armonizacion_dato"] != nil {
						dato_armonizacion_str := subgrupo_detalle[0]["armonizacion_dato"].(string)
						json.Unmarshal([]byte(dato_armonizacion_str), &dato_armonizacion)
						detalles_armonizacion = dato_armonizacion
						armonizacion["armo"] = dato_armonizacion[index]
					}
					if subgrupo_detalle[0]["dato_plan"] != nil {
						dato_plan_str := subgrupo_detalle[0]["dato_plan"].(string)
						json.Unmarshal([]byte(dato_plan_str), &dato_plan)

						if dato_plan[index] != nil {
							actividad = dato_plan[index].(map[string]interface{})
							detalles[i] = dato_plan
							if v != "" {
								forkData[v] = actividad["dato"]
							}
						} else {
							detalles = append(detalles, map[string]interface{}{})
						}
					}
				}
			}
		} else {
			if detalles[i][index] != nil {
				forkData[v] = detalles[i][index].(map[string]interface{})["dato"]
			}
			if detalles_armonizacion[index] != nil {
				armonizacion["armo"] = detalles_armonizacion[index]
			}
		}
	}
	if !detallesLlenados {
		detallesLlenados = true
	}

	if detalles_armonizacion[index] == nil {
		armonizacion["armo"] = map[string]interface{}{
			"armonizacionPED": "",
			"armonizacionPI":  "",
		}
	}

	validadores = append(validadores, forkData)
	return validadores, armonizacion
}

func getChildren(children []interface{}, exist bool) (childrenTree []map[string]interface{}) {
	var res map[string]interface{}
	var nodo []map[string]interface{}

	for _, child := range children {
		childStr := child.(string)
		forkData := make(map[string]interface{})
		var id string
		err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=_id:"+childStr+"&fields=nombre,_id,hijos,activo", &res)
		if err != nil {
			return
		}
		helpers.LimpiezaRespuestaRefactor(res, &nodo)
		if nodo[0]["activo"] == true {
			forkData["id"] = nodo[0]["_id"]
			forkData["nombre"] = nodo[0]["nombre"]
			id = nodo[0]["_id"].(string)

			if len(nodo[0]["hijos"].([]interface{})) > 0 {
				aux := getChildren(nodo[0]["hijos"].([]interface{}), true)
				if len(aux) == 0 {
					forkData["sub"] = ""
				} else {
					forkData["sub"] = aux
				}
			}

			childrenTree = append(childrenTree, forkData)
		}
		id_arr = append(id_arr, id)
		add(id)
	}
	ids = append(ids, id_arr)
	return
}

func ArbolArmonizacionV2(armonizacion string) []map[string]interface{} {

	var estrategias []map[string]interface{}
	var metas []map[string]interface{}
	var lineamientos []map[string]interface{}
	var arregloArmo []map[string]interface{}

	if armonizacion != "" {
		armonizacionPED := strings.Split(armonizacion, ",")
		for i := 0; i < len(armonizacionPED); i++ {
			var respuesta map[string]interface{}
			var respuestaSubgrupo map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+armonizacionPED[i], &respuesta); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaSubgrupo)

				if len(respuestaSubgrupo) > 0 {
					nombre := strings.ToLower(respuestaSubgrupo["nombre"].(string))
					if strings.Contains(nombre, "lineamiento") {
						lineamientos = append(lineamientos, respuestaSubgrupo)
					} else if strings.Contains(nombre, "meta") {
						metas = append(metas, respuestaSubgrupo)
					} else if strings.Contains(nombre, "estrategia") {
						estrategias = append(estrategias, respuestaSubgrupo)
					}
				}
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}
		}

		for i := 0; i < len(lineamientos); i++ {
			arregloArmo = append(arregloArmo, map[string]interface{}{
				"_id":                  lineamientos[i]["_id"],
				"nombreLineamiento":    lineamientos[i]["nombre"],
				"meta":                 []map[string]interface{}{},
				"nombrePlanDesarrollo": "Plan Estrategico de Desarrollo",
				"hijos":                lineamientos[i]["hijos"],
			})
		}

		for i := 0; i < len(metas); i++ {
			foundPadreMeta := false
			for j := 0; j < len(arregloArmo); j++ {
				if arregloArmo[j]["_id"] == metas[i]["padre"] {
					arregloArmo[j]["meta"] = append(arregloArmo[j]["meta"].([]map[string]interface{}), map[string]interface{}{
						"_id":         metas[i]["_id"],
						"nombreMeta":  metas[i]["nombre"],
						"estrategias": []map[string]interface{}{},
					})
					foundPadreMeta = true
					break
				}
			}
			if !foundPadreMeta {
				arregloArmo = append(arregloArmo, map[string]interface{}{
					"_id":               metas[i]["padre"],
					"nombreLineamiento": "No seleccionado",
					"meta": []map[string]interface{}{
						{
							"_id":         metas[i]["_id"],
							"nombreMeta":  metas[i]["nombre"],
							"estrategias": []map[string]interface{}{},
						},
					},
					"nombrePlanDesarrollo": "Plan Estrategico de Desarrollo",
					"hijos":                []interface{}{},
				})
			}
		}

		for i := 0; i < len(estrategias); i++ {
			foundPadreEstrategia := false
			for j := 0; j < len(arregloArmo); j++ {
				for k := 0; k < len(arregloArmo[j]["meta"].([]map[string]interface{})); k++ {
					if arregloArmo[j]["meta"].([]map[string]interface{})[k]["_id"] == estrategias[i]["padre"] {
						arregloArmo[j]["meta"].([]map[string]interface{})[k]["estrategias"] = append(arregloArmo[j]["meta"].([]map[string]interface{})[k]["estrategias"].([]map[string]interface{}), map[string]interface{}{
							"_id":                   estrategias[i]["_id"],
							"nombreEstrategia":      estrategias[i]["nombre"],
							"descripcionEstrategia": estrategias[i]["descripcion"],
						})
						foundPadreEstrategia = true
						break
					}
				}
			}
			if !foundPadreEstrategia {
				for j := 0; j < len(arregloArmo); j++ {
					for k := 0; k < len(arregloArmo[j]["hijos"].([]interface{})); k++ {
						if arregloArmo[j]["hijos"].([]interface{})[k] == estrategias[i]["padre"] {
							arregloArmo[j]["meta"] = append(arregloArmo[j]["meta"].([]map[string]interface{}), map[string]interface{}{
								"_id":        estrategias[i]["padre"],
								"nombreMeta": "No seleccionado",
								"estrategias": []map[string]interface{}{
									{
										"_id":                   estrategias[i]["_id"],
										"nombreEstrategia":      estrategias[i]["nombre"],
										"descripcionEstrategia": estrategias[i]["descripcion"],
									},
								},
							})
							foundPadreEstrategia = true
							break
						}
					}
					if foundPadreEstrategia {
						break
					}
				}
			}
			if !foundPadreEstrategia {
				arregloArmo = append(arregloArmo, map[string]interface{}{
					"_id":               "",
					"nombreLineamiento": "No seleccionado",
					"meta": []map[string]interface{}{
						{
							"_id":        estrategias[i]["padre"],
							"nombreMeta": "No seleccionado",
							"estrategias": []map[string]interface{}{
								{
									"_id":                   estrategias[i]["_id"],
									"nombreEstrategia":      estrategias[i]["nombre"],
									"descripcionEstrategia": estrategias[i]["descripcion"],
								},
							},
						},
					},
					"nombrePlanDesarrollo": "Plan Estrategico de Desarrollo",
					"hijos": []interface{}{
						estrategias[i]["padre"],
					},
				})
			}
		}

		if len(arregloArmo) > 0 {
			for i := 0; i < len(arregloArmo); i++ {
				if len(arregloArmo[i]["meta"].([]map[string]interface{})) == 0 {
					arregloArmo[i]["meta"] = append(arregloArmo[i]["meta"].([]map[string]interface{}), map[string]interface{}{
						"_id":        "",
						"nombreMeta": "No seleccionado",
						"estrategias": []map[string]interface{}{
							{
								"_id":                   "",
								"nombreEstrategia":      "No seleccionado",
								"descripcionEstrategia": "No seleccionado",
							},
						},
					})
				} else {
					for j := 0; j < len(arregloArmo[i]["meta"].([]map[string]interface{})); j++ {
						if len(arregloArmo[i]["meta"].([]map[string]interface{})[j]["estrategias"].([]map[string]interface{})) == 0 {
							arregloArmo[i]["meta"].([]map[string]interface{})[j]["estrategias"] = append(arregloArmo[i]["meta"].([]map[string]interface{})[j]["estrategias"].([]map[string]interface{}), map[string]interface{}{
								"_id":                   "",
								"nombreEstrategia":      "No seleccionado",
								"descripcionEstrategia": "No seleccionado",
							})
						}
					}
				}
				delete(arregloArmo[i], "hijos")
			}
		} else {
			arregloArmo = append(arregloArmo, map[string]interface{}{
				"_id":               "",
				"nombreLineamiento": "No seleccionado",
				"meta": []map[string]interface{}{
					{
						"_id":        "",
						"nombreMeta": "No seleccionado",
						"estrategias": []map[string]interface{}{
							{
								"_id":                   "",
								"nombreEstrategia":      "No seleccionado",
								"descripcionEstrategia": "No seleccionado",
							},
						},
					},
				},
				"nombrePlanDesarrollo": "Plan Estrategico de Desarrollo",
			})
		}
	} else {
		arregloArmo = append(arregloArmo, map[string]interface{}{
			"_id":               "",
			"nombreLineamiento": "No seleccionado",
			"meta": []map[string]interface{}{
				{
					"_id":        "",
					"nombreMeta": "No seleccionado",
					"estrategias": []map[string]interface{}{
						{
							"_id":                   "",
							"nombreEstrategia":      "No seleccionado",
							"descripcionEstrategia": "No seleccionado",
						},
					},
				},
			},
			"nombrePlanDesarrollo": "Plan Estrategico de Desarrollo",
		})
	}

	return arregloArmo
}

func ArbolArmonizacion(armonizacion string) []map[string]interface{} {

	var respuesta map[string]interface{}
	var lineamientos []map[string]interface{}
	var metas []map[string]interface{}
	var estrategias []map[string]interface{}
	var arreglo []map[string]interface{}
	armonizacionPED := strings.Split(armonizacion, ",")
	for i := 0; i < len(armonizacionPED); i++ {
		var respuestaSubgrupo map[string]interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+armonizacionPED[i], &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaSubgrupo)
			if len(respuestaSubgrupo) > 0 {
				nombre := strings.ToLower(respuestaSubgrupo["nombre"].(string))
				if strings.Contains(nombre, "lineamiento") {
					lineamientos = append(lineamientos, respuestaSubgrupo)
				}
				if strings.Contains(nombre, "meta") {
					metas = append(metas, respuestaSubgrupo)
				}
				if strings.Contains(nombre, "estrategia") {
					estrategias = append(estrategias, respuestaSubgrupo)
				}
			}
		} else {
			panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
		}
	}

	for i := 0; i < len(lineamientos); i++ {
		meta := make(map[string]interface{})
		lineamiento := make(map[string]interface{})
		estrategia := make(map[string]interface{})
		var arregloEstrategias []map[string]interface{}
		var arregloMetas []map[string]interface{}

		estrategia["_id"] = ""
		estrategia["nombreEstrategia"] = ""
		estrategia["descripcionEstrategia"] = ""

		arregloEstrategias = append(arregloEstrategias, estrategia)

		meta["_id"] = ""
		meta["nombreMeta"] = ""
		meta["estrategias"] = arregloEstrategias

		arregloMetas = append(arregloMetas, meta)

		lineamiento["_id"] = lineamientos[i]["_id"]
		lineamiento["nombreLineamiento"] = lineamientos[i]["nombre"]
		lineamiento["nombrePlanDesarrollo"] = "Plan Estrategico de Desarrollo"
		lineamiento["meta"] = arregloMetas

		arreglo = append(arreglo, lineamiento)
	}

	for i := 0; i < len(metas); i++ {

		meta := make(map[string]interface{})
		estrategia := make(map[string]interface{})
		var auxMeta = metas[i]
		bandera := false
		var arregloEstrategias []map[string]interface{}
		estrategia["_id"] = ""
		estrategia["nombreEstrategia"] = ""
		estrategia["descripcionEstrategia"] = ""

		arregloEstrategias = append(arregloEstrategias, estrategia)

		meta["_id"] = auxMeta["_id"]
		meta["nombreMeta"] = auxMeta["nombre"]
		meta["estrategias"] = arregloEstrategias

		for j := 0; j < len(arreglo); j++ {
			if arreglo[j]["_id"] == auxMeta["padre"] {

				bandera = true
				aux := arreglo[j]["meta"].([]map[string]interface{})
				if aux[0]["_id"] == "" {
					aux = append(aux[:0], aux[1:]...)
				}
				aux = append(aux, meta)
				arreglo[j]["meta"] = aux
				break
			}
		}
		if bandera == false {
			lineamiento := make(map[string]interface{})
			var arregloMetas []map[string]interface{}
			arregloMetas = append(arregloMetas, meta)
			lineamiento["_id"] = ""
			lineamiento["nombreLineamiento"] = ""
			lineamiento["nombrePlanDesarrollo"] = "Plan Estrategico de Desarrollo"
			lineamiento["meta"] = arregloMetas
			arreglo = append(arreglo, lineamiento)
		}
	}

	for i := 0; i < len(estrategias); i++ {
		var auxEstrategia = estrategias[i]
		estrategia := make(map[string]interface{})
		bandera := false

		estrategia["_id"] = auxEstrategia["_id"]
		estrategia["nombreEstrategia"] = auxEstrategia["nombre"]
		estrategia["descripcionEstrategia"] = auxEstrategia["descripcion"]

		for j := 0; j < len(metas); j++ {
			if metas[j]["_id"] == auxEstrategia["padre"] {
				for n := 0; n < len(arreglo); n++ {
					if arreglo[n]["_id"] == metas[j]["padre"] {
						bandera = true
						auxMetas := arreglo[n]["meta"].([]map[string]interface{})

						for k := 0; k < len(auxMetas); k++ {
							if auxMetas[k]["_id"] == auxEstrategia["padre"] {
								aux2 := auxMetas[k]["estrategias"].([]map[string]interface{})
								if aux2[0]["_id"] == "" {
									aux2 = append(aux2[:0], aux2[1:]...)
								}
								aux2 = append(aux2, estrategia)
								auxMetas[k]["estrategias"] = aux2

								arreglo[n]["meta"] = auxMetas
								break
							}
						}
						break
					}
				}
				break
			}
		}
		if !bandera {
			meta := make(map[string]interface{})
			lineamiento := make(map[string]interface{})
			var arregloEstrategias []map[string]interface{}
			var arregloMetas []map[string]interface{}
			arregloEstrategias = append(arregloEstrategias, estrategia)

			meta["_id"] = ""
			meta["nombreMeta"] = ""
			meta["estrategias"] = arregloEstrategias

			arregloMetas = append(arregloMetas, meta)

			lineamiento["_id"] = ""
			lineamiento["nombreLineamiento"] = ""
			lineamiento["nombrePlanDesarrollo"] = "Plan Estrategico de Desarrollo"
			lineamiento["meta"] = arregloMetas
			arreglo = append(arreglo, lineamiento)
		}
	}

	return arreglo
}

func ArbolArmonizacionPIV2(armonizacion string) []map[string]interface{} {

	var estrategias []map[string]interface{}
	var lineamientos []map[string]interface{}
	var factores []map[string]interface{}
	var arregloArmo []map[string]interface{}

	if armonizacion != "" {
		armonizacionPI := strings.Split(armonizacion, ",")
		for i := 0; i < len(armonizacionPI); i++ {
			var respuesta map[string]interface{}
			var respuestaSubgrupo map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+armonizacionPI[i], &respuesta); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaSubgrupo)
				if len(respuestaSubgrupo) > 0 {
					nombre := strings.ToLower(respuestaSubgrupo["nombre"].(string))
					if (strings.Contains(nombre, "eje") || strings.Contains(nombre, "transformador")) || strings.Contains(nombre, "nivel 1") {
						factores = append(factores, respuestaSubgrupo)
					} else if strings.Contains(nombre, "lineamientos") || strings.Contains(nombre, "lineamiento") || strings.Contains(nombre, "nivel 2") {
						lineamientos = append(lineamientos, respuestaSubgrupo)
					} else if strings.Contains(nombre, "estrategia") || strings.Contains(nombre, "proyecto") || strings.Contains(nombre, "nivel 3") {
						estrategias = append(estrategias, respuestaSubgrupo)
					}
				}
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}
		}

		for i := 0; i < len(factores); i++ {
			arregloArmo = append(arregloArmo, map[string]interface{}{
				"_id":                  factores[i]["_id"],
				"nombreFactor":         factores[i]["nombre"],
				"lineamientos":         []map[string]interface{}{},
				"nombrePlanDesarrollo": "Plan Indicativo",
				"hijos":                factores[i]["hijos"],
			})
		}

		for i := 0; i < len(lineamientos); i++ {
			foundPadreMeta := false
			for j := 0; j < len(arregloArmo); j++ {
				if arregloArmo[j]["_id"] == lineamientos[i]["padre"] {
					arregloArmo[j]["lineamientos"] = append(arregloArmo[j]["lineamientos"].([]map[string]interface{}), map[string]interface{}{
						"_id":               lineamientos[i]["_id"],
						"nombreLineamiento": lineamientos[i]["nombre"],
						"estrategias":       []map[string]interface{}{},
					})
					foundPadreMeta = true
					break
				}
			}
			if !foundPadreMeta {
				arregloArmo = append(arregloArmo, map[string]interface{}{
					"_id":          lineamientos[i]["padre"],
					"nombreFactor": "No seleccionado",
					"lineamientos": []map[string]interface{}{
						{
							"_id":               lineamientos[i]["_id"],
							"nombreLineamiento": lineamientos[i]["nombre"],
							"estrategias":       []map[string]interface{}{},
						},
					},
					"nombrePlanDesarrollo": "Plan Indicativo",
					"hijos":                []interface{}{},
				})
			}
		}

		for i := 0; i < len(estrategias); i++ {
			foundPadreEstrategia := false
			for j := 0; j < len(arregloArmo); j++ {
				for k := 0; k < len(arregloArmo[j]["lineamientos"].([]map[string]interface{})); k++ {
					if arregloArmo[j]["lineamientos"].([]map[string]interface{})[k]["_id"] == estrategias[i]["padre"] {
						arregloArmo[j]["lineamientos"].([]map[string]interface{})[k]["estrategias"] = append(arregloArmo[j]["lineamientos"].([]map[string]interface{})[k]["estrategias"].([]map[string]interface{}), map[string]interface{}{
							"_id":                   estrategias[i]["_id"],
							"nombreEstrategia":      estrategias[i]["nombre"],
							"descripcionEstrategia": estrategias[i]["descripcion"],
						})
						foundPadreEstrategia = true
						break
					}
				}
			}
			if !foundPadreEstrategia {
				for j := 0; j < len(arregloArmo); j++ {
					for k := 0; k < len(arregloArmo[j]["hijos"].([]interface{})); k++ {
						if arregloArmo[j]["hijos"].([]interface{})[k] == estrategias[i]["padre"] {
							arregloArmo[j]["lineamientos"] = append(arregloArmo[j]["lineamientos"].([]map[string]interface{}), map[string]interface{}{
								"_id":               estrategias[i]["padre"],
								"nombreLineamiento": "No seleccionado",
								"estrategias": []map[string]interface{}{
									{
										"_id":                   estrategias[i]["_id"],
										"nombreEstrategia":      estrategias[i]["nombre"],
										"descripcionEstrategia": estrategias[i]["descripcion"],
									},
								},
							})
							foundPadreEstrategia = true
							break
						}
					}
					if foundPadreEstrategia {
						break
					}
				}
			}
			if !foundPadreEstrategia {
				arregloArmo = append(arregloArmo, map[string]interface{}{
					"_id":          "",
					"nombreFactor": "No seleccionado",
					"lineamientos": []map[string]interface{}{
						{
							"_id":               estrategias[i]["padre"],
							"nombreLineamiento": "No seleccionado",
							"estrategias": []map[string]interface{}{
								{
									"_id":                   estrategias[i]["_id"],
									"nombreEstrategia":      estrategias[i]["nombre"],
									"descripcionEstrategia": estrategias[i]["descripcion"],
								},
							},
						},
					},
					"nombrePlanDesarrollo": "Plan Indicativo",
					"hijos": []interface{}{
						estrategias[i]["padre"],
					},
				})
			}
		}

		if len(arregloArmo) > 0 {
			for i := 0; i < len(arregloArmo); i++ {
				if len(arregloArmo[i]["lineamientos"].([]map[string]interface{})) == 0 {
					arregloArmo[i]["lineamientos"] = append(arregloArmo[i]["lineamientos"].([]map[string]interface{}), map[string]interface{}{
						"_id":               "",
						"nombreLineamiento": "No seleccionado",
						"estrategias": []map[string]interface{}{
							{
								"_id":                   "",
								"nombreEstrategia":      "No seleccionado",
								"descripcionEstrategia": "No seleccionado",
							},
						},
					})
				} else {
					for j := 0; j < len(arregloArmo[i]["lineamientos"].([]map[string]interface{})); j++ {
						if len(arregloArmo[i]["lineamientos"].([]map[string]interface{})[j]["estrategias"].([]map[string]interface{})) == 0 {
							arregloArmo[i]["lineamientos"].([]map[string]interface{})[j]["estrategias"] = append(arregloArmo[i]["lineamientos"].([]map[string]interface{})[j]["estrategias"].([]map[string]interface{}), map[string]interface{}{
								"_id":                   "",
								"nombreEstrategia":      "No seleccionado",
								"descripcionEstrategia": "No seleccionado",
							})
						}
					}
				}
				delete(arregloArmo[i], "hijos")
			}
		} else {
			arregloArmo = append(arregloArmo, map[string]interface{}{
				"_id":          "",
				"nombreFactor": "No seleccionado",
				"lineamientos": []map[string]interface{}{
					{
						"_id":               "",
						"nombreLineamiento": "No seleccionado",
						"estrategias": []map[string]interface{}{
							{
								"_id":                   "",
								"nombreEstrategia":      "No seleccionado",
								"descripcionEstrategia": "No seleccionado",
							},
						},
					},
				},
				"nombrePlanDesarrollo": "Plan Indicativo",
			})
		}
	} else {
		arregloArmo = append(arregloArmo, map[string]interface{}{
			"_id":          "",
			"nombreFactor": "No seleccionado",
			"lineamientos": []map[string]interface{}{
				{
					"_id":               "",
					"nombreLineamiento": "No seleccionado",
					"estrategias": []map[string]interface{}{
						{
							"_id":                   "",
							"nombreEstrategia":      "No seleccionado",
							"descripcionEstrategia": "No seleccionado",
						},
					},
				},
			},
			"nombrePlanDesarrollo": "Plan Indicativo",
		})
	}

	return arregloArmo
}

func ArbolArmonizacionPI(armonizacion interface{}) []map[string]interface{} {

	var respuesta map[string]interface{}
	var lineamientos []map[string]interface{}
	var factores []map[string]interface{}
	var estrategias []map[string]interface{}
	var arreglo []map[string]interface{}
	if armonizacion != "" {

		armonizacionPI := strings.Split(armonizacion.(string), ",")

		for i := 0; i < len(armonizacionPI); i++ {
			var respuestaSubgrupo map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+armonizacionPI[i], &respuesta); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaSubgrupo)
				if len(respuestaSubgrupo) > 0 {
					nombre := strings.ToLower(respuestaSubgrupo["nombre"].(string))
					if (strings.Contains(nombre, "eje") && strings.Contains(nombre, "transformador")) || strings.Contains(nombre, "nivel 1") {
						factores = append(factores, respuestaSubgrupo)
					}
					if strings.Contains(nombre, "lineamientos") || strings.Contains(nombre, "lineamiento") || strings.Contains(nombre, "nivel 2") {
						lineamientos = append(lineamientos, respuestaSubgrupo)
					}
					if strings.Contains(nombre, "estrategia") || strings.Contains(nombre, "proyecto") || strings.Contains(nombre, "nivel 3") {
						estrategias = append(estrategias, respuestaSubgrupo)
					}
				}
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}
		}

		for i := 0; i < len(factores); i++ {
			factor := make(map[string]interface{})
			lineamiento := make(map[string]interface{})
			estrategia := make(map[string]interface{})
			var arregloEstrategias []map[string]interface{}
			var arregloLineamientos []map[string]interface{}

			estrategia["_id"] = ""
			estrategia["nombreEstrategia"] = ""
			estrategia["descripcionEstrategia"] = ""

			arregloEstrategias = append(arregloEstrategias, estrategia)

			lineamiento["_id"] = ""
			lineamiento["nombreLineamiento"] = ""
			lineamiento["estrategias"] = arregloEstrategias

			arregloLineamientos = append(arregloLineamientos, lineamiento)

			factor["_id"] = factores[i]["_id"]
			factor["nombreFactor"] = factores[i]["nombre"]
			factor["nombrePlanDesarrollo"] = "Plan Indicativo"
			factor["lineamientos"] = arregloLineamientos

			arreglo = append(arreglo, factor)
		}

		for i := 0; i < len(lineamientos); i++ {

			lineamiento := make(map[string]interface{})
			estrategia := make(map[string]interface{})
			var auxLineamiento = lineamientos[i]
			bandera := false
			var arregloEstrategias []map[string]interface{}
			estrategia["_id"] = ""
			estrategia["nombreEstrategia"] = ""
			estrategia["descripcionEstrategia"] = ""

			arregloEstrategias = append(arregloEstrategias, estrategia)

			lineamiento["_id"] = auxLineamiento["_id"]
			lineamiento["nombreLineamiento"] = auxLineamiento["nombre"]
			lineamiento["estrategias"] = arregloEstrategias

			for j := 0; j < len(arreglo); j++ {

				if arreglo[j]["_id"] == auxLineamiento["padre"] {

					bandera = true
					aux := arreglo[j]["lineamientos"].([]map[string]interface{})
					if aux[0]["_id"] == "" {
						aux = append(aux[:0], aux[1:]...)
					}
					aux = append(aux, lineamiento)
					arreglo[j]["lineamientos"] = aux
					break
				}
			}
			if bandera == false {
				factor := make(map[string]interface{})
				var arregloLineamientos []map[string]interface{}
				arregloLineamientos = append(arregloLineamientos, lineamiento)
				factor["_id"] = ""
				factor["nombreFactor"] = ""
				factor["nombrePlanDesarrollo"] = "Plan Indicativo"
				factor["lineamientos"] = arregloLineamientos
				arreglo = append(arreglo, factor)
			}
		}

		for i := 0; i < len(estrategias); i++ {
			var auxEstrategia = estrategias[i]
			estrategia := make(map[string]interface{})
			bandera := false

			estrategia["_id"] = auxEstrategia["_id"]
			estrategia["nombreEstrategia"] = auxEstrategia["nombre"]
			estrategia["descripcionEstrategia"] = auxEstrategia["descripcion"]

			for j := 0; j < len(lineamientos); j++ {
				if lineamientos[j]["_id"] == auxEstrategia["padre"] {
					for n := 0; n < len(arreglo); n++ {
						if arreglo[n]["_id"] == lineamientos[j]["padre"] {
							bandera = true
							auxLineamientos := arreglo[n]["lineamientos"].([]map[string]interface{})

							for k := 0; k < len(auxLineamientos); k++ {
								if auxLineamientos[k]["_id"] == auxEstrategia["padre"] {
									aux2 := auxLineamientos[k]["estrategias"].([]map[string]interface{})
									if aux2[0]["_id"] == "" {
										aux2 = append(aux2[:0], aux2[1:]...)
									}
									aux2 = append(aux2, estrategia)
									auxLineamientos[k]["estrategias"] = aux2

									arreglo[n]["lineamientos"] = auxLineamientos
									break
								}
							}
							break
						}
					}
					break
				}
			}
			if !bandera {
				lineamiento := make(map[string]interface{})
				factor := make(map[string]interface{})
				var arregloEstrategias []map[string]interface{}
				var arregloLineamientos []map[string]interface{}
				arregloEstrategias = append(arregloEstrategias, estrategia)

				lineamiento["_id"] = ""
				lineamiento["nombreLineamiento"] = ""
				lineamiento["estrategias"] = arregloEstrategias

				arregloLineamientos = append(arregloLineamientos, lineamiento)

				factor["_id"] = ""
				factor["nombreFactor"] = ""
				factor["nombrePlanDesarrollo"] = "Plan Indicativo"
				factor["lineamientos"] = arregloLineamientos
				arreglo = append(arreglo, factor)
			}
		}
	} else {
		lineamiento := make(map[string]interface{})
		factor := make(map[string]interface{})
		estrategia := make(map[string]interface{})

		estrategia["_id"] = ""
		estrategia["nombreEstrategia"] = ""
		estrategia["descripcionEstrategia"] = ""

		var arregloEstrategias []map[string]interface{}
		var arregloLineamientos []map[string]interface{}
		arregloEstrategias = append(arregloEstrategias, estrategia)

		lineamiento["_id"] = ""
		lineamiento["nombreLineamiento"] = ""
		lineamiento["estrategias"] = arregloEstrategias

		arregloLineamientos = append(arregloLineamientos, lineamiento)

		factor["_id"] = ""
		factor["nombreFactor"] = ""
		factor["nombrePlanDesarrollo"] = "Plan Indicativo"
		factor["lineamientos"] = arregloLineamientos
		arreglo = append(arreglo, factor)
	}

	return arreglo
}

type nodo struct {
	valor     int
	divisible bool
	hijos     []*nodo
}

func MinComMul_Armonization(armoPED, armoPI []map[string]interface{}, lenIndicadores int) int {
	sizePED := &nodo{valor: len(armoPED)}
	for _, n2 := range armoPED {
		h1 := &nodo{valor: len(n2["meta"].([]map[string]interface{}))}
		sizePED.hijos = append(sizePED.hijos, h1)
		for _, n3 := range n2["meta"].([]map[string]interface{}) {
			h2 := &nodo{valor: len(n3["estrategias"].([]map[string]interface{}))}
			h1.hijos = append(h1.hijos, h2)
		}
	}

	sizePI := &nodo{valor: len(armoPI)}
	for _, n2 := range armoPI {
		h1 := &nodo{valor: len(n2["lineamientos"].([]map[string]interface{}))}
		sizePI.hijos = append(sizePI.hijos, h1)
		for _, n3 := range n2["lineamientos"].([]map[string]interface{}) {
			h2 := &nodo{valor: len(n3["estrategias"].([]map[string]interface{}))}
			h1.hijos = append(h1.hijos, h2)
		}
	}

	fitSize1 := false
	fitSize2 := false
	fitSize3 := false
	rowMax := lenIndicadores
	for !(fitSize1 && fitSize2 && fitSize3) {
		fitSize1 = calcMinCol(sizePED, rowMax)
		fitSize2 = calcMinCol(sizePI, rowMax)
		fitSize3 = (rowMax % lenIndicadores) == 0
		rowMax++
	}

	return rowMax - 1
}

func calcMinCol(node *nodo, size int) bool {
	if node.valor == 0 {
		node.valor = 1
	}
	if (size % node.valor) == 0 {
		node.divisible = true
		div := size / node.valor
		for _, hijo := range node.hijos {
			if !calcMinCol(hijo, div) {
				return false
			}
		}
	} else {
		node.divisible = false
	}
	return node.divisible
}

func TablaIdentificaciones(consolidadoExcelPlanAnual *excelize.File, planId string) *excelize.File {
	var res map[string]interface{}
	var identificaciones []map[string]interface{}
	var recursos []map[string]interface{}
	var contratistas []map[string]interface{}
	var docentes map[string]interface{}
	var data_identi []map[string]interface{}
	var rubro string
	var nombreRubro string

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+planId, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &identificaciones)
	}

	for i := 0; i < len(identificaciones); i++ {
		identificacion := identificaciones[i]
		nombre := strings.ToLower(identificacion["nombre"].(string))
		if strings.Contains(nombre, "recurso") {
			if identificacion["dato"] != nil {
				var dato map[string]interface{}
				dato_str := identificacion["dato"].(string)
				json.Unmarshal([]byte(dato_str), &dato)
				for key := range dato {
					element := dato[key].(map[string]interface{})
					if element["activo"] == true {
						data_identi = append(data_identi, element)
					}
				}
				recursos = data_identi
			}
		} else if strings.Contains(nombre, "contratista") {
			if identificacion["dato"] != nil {
				var dato map[string]interface{}
				var data_identi []map[string]interface{}
				dato_str := identificacion["dato"].(string)
				json.Unmarshal([]byte(dato_str), &dato)
				for key := range dato {
					element := dato[key].(map[string]interface{})
					if element["rubro"] == nil {
						rubro = "Información no suministrada"
					} else {
						rubro = element["rubro"].(string)
					}
					if element["rubroNombre"] == nil {
						nombreRubro = "Información no suministrada"
					} else {
						nombreRubro = element["rubroNombre"].(string)
					}
					if element["activo"] == true {
						data_identi = append(data_identi, element)
					}
				}
				contratistas = data_identi
			}
		} else if strings.Contains(nombre, "docente") {
			dato := map[string]interface{}{}
			var data_identi []map[string]interface{}
			if identificacion["dato"] != nil && identificacion["dato"] != "{}" {
				result := make(map[string]interface{})
				dato_str := identificacion["dato"].(string)

				// ? Se tiene en cuenta la nueva estructura la info ahora está en identificacion-detalle, pero tambien tiene en cuenta la estructura de indentificaciones viejas (else)
				if strings.Contains(dato_str, "ids_detalle") {
					json.Unmarshal([]byte(dato_str), &dato)

					var identi map[string]interface{}
					iddetail := ""
					identificacionDetalle := map[string]interface{}{}
					errIdentificacionDetalle := error(nil)

					identi = nil
					data_identi = nil
					identificacionDetalle = map[string]interface{}{}
					iddetail = dato["ids_detalle"].(map[string]interface{})["rhf"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						if dato_aux == "{}" {
							result["rhf"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rhf"] = data_identi
						}
					} else {
						result["rhf"] = "{}"
					}

					identi = nil
					data_identi = nil
					identificacionDetalle = map[string]interface{}{}
					iddetail = dato["ids_detalle"].(map[string]interface{})["rhv_pre"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						if dato_aux == "{}" {
							result["rhv_pre"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rhv_pre"] = data_identi
						}
					} else {
						result["rhv_pre"] = "{}"
					}

					identi = nil
					data_identi = nil
					identificacionDetalle = map[string]interface{}{}
					iddetail = dato["ids_detalle"].(map[string]interface{})["rhv_pos"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						if dato_aux == "{}" {
							result["rhv_pos"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rhv_pos"] = data_identi
						}
					} else {
						result["rhv_pos"] = "{}"
					}

					identi = nil
					data_identi = nil
					identificacionDetalle = map[string]interface{}{}
					iddetail = dato["ids_detalle"].(map[string]interface{})["rubros"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						if dato_aux == "{}" {
							result["rubros"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rubros"] = data_identi
						}
					} else {
						result["rubros"] = "{}"
					}

					identi = nil
					data_identi = nil
					identificacionDetalle = map[string]interface{}{}
					iddetail = dato["ids_detalle"].(map[string]interface{})["rubros_pos"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						if dato_aux == "{}" {
							result["rubros_pos"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rubros_pos"] = data_identi
						}
					} else {
						result["rubros_pos"] = "{}"
					}

				} else {
					json.Unmarshal([]byte(dato_str), &dato)

					var identi map[string]interface{}
					identi = nil
					data_identi = nil
					dato_aux := dato["rhf"].(string)
					if dato_aux == "{}" {
						result["rhf"] = "{}"
					} else {
						json.Unmarshal([]byte(dato_aux), &identi)
						for key := range identi {
							element := identi[key].(map[string]interface{})
							if element["activo"] == true {
								data_identi = append(data_identi, element)
							}
						}
						result["rhf"] = data_identi
					}

					identi = nil
					data_identi = nil
					dato_aux = dato["rhv_pre"].(string)
					if dato_aux == "{}" {
						result["rhv_pre"] = "{}"
					} else {
						json.Unmarshal([]byte(dato_aux), &identi)
						for key := range identi {
							element := identi[key].(map[string]interface{})
							if element["activo"] == true {
								data_identi = append(data_identi, element)
							}
						}
						result["rhv_pre"] = data_identi
					}

					identi = nil
					data_identi = nil
					dato_aux = dato["rhv_pos"].(string)
					if dato_aux == "{}" {
						result["rhv_pos"] = "{}"
					} else {
						json.Unmarshal([]byte(dato_aux), &identi)
						for key := range identi {
							element := identi[key].(map[string]interface{})
							if element["activo"] == true {
								data_identi = append(data_identi, element)
							}
						}
						result["rhv_pos"] = data_identi
					}

					identi = nil
					data_identi = nil
					if dato["rubros"] != nil {
						dato_aux = dato["rubros"].(string)
						if dato_aux == "{}" {
							result["rubros"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rubros"] = data_identi
						}
					}

					identi = nil
					data_identi = nil
					if dato["rubros_pos"] != nil {
						dato_aux = dato["rubros_pos"].(string)
						if dato_aux == "{}" {
							result["rubros_pos"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rubros_pos"] = data_identi
						}
					}
				}

				docentes = result
			}
		}
	}
	return construirTablas(consolidadoExcelPlanAnual, recursos, contratistas, docentes, rubro, nombreRubro)
}

func construirTablas(consolidadoExcelPlanAnual *excelize.File, recursos []map[string]interface{}, contratistas []map[string]interface{}, docentes map[string]interface{}, rubro string, nombreRubro string) *excelize.File {
	stylecontent, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1}}})
	stylecontentMR, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		NumFmt:    183,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1}}})
	stylecontentC, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1}}})
	stylecontentS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1}}})
	stylecontentMRS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		NumFmt:    183,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1}}})
	stylecontentCS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1}}})
	styletitles, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Family: "Arial", Size: 12, Color: "000000"},
		Border:    []excelize.Border{{Type: "bottom", Color: "000000", Style: 1}}})
	stylehead, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Color: "ffffff"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1}}})

	sheetName := "Identificaciones"

	consolidadoExcelPlanAnual.NewSheet(sheetName)
	consolidadoExcelPlanAnual.InsertCols(sheetName, "A", 1)
	disable := false
	if err := consolidadoExcelPlanAnual.SetSheetView(sheetName, -1, &excelize.ViewOptions{
		ShowGridLines: &disable,
	}); err != nil {
		fmt.Println(err)
	}
	consolidadoExcelPlanAnual.MergeCell(sheetName, "B1", "F1")

	consolidadoExcelPlanAnual.SetColWidth(sheetName, "A", "A", 2)
	consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "D", 26)
	consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "C", 30)
	consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "E", 35)
	consolidadoExcelPlanAnual.SetColWidth(sheetName, "F", "G", 20)
	consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "I", 20)

	consolidadoExcelPlanAnual.SetCellValue(sheetName, "B1", "Identificación de recursos")
	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B1", "F1", styletitles)

	consolidadoExcelPlanAnual.SetCellValue(sheetName, "B3", "Código del rubro")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "C3", "Nombre del rubro")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "D3", "Valor")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "E3", "Descripción del bien y/o servicio")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "F3", "Actividades")
	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B3", "F3", stylehead)
	consolidadoExcelPlanAnual.SetRowHeight(sheetName, 2, 7)

	contador := 4
	for i := 0; i < len(recursos); i++ {
		aux := recursos[i]
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), aux["codigo"])
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), aux["Nombre"])
		strValor := strings.TrimLeft(fmt.Sprintf("%v", aux["valor"]), "$")
		strValor = strings.ReplaceAll(strValor, ",", "")
		arrValor := strings.Split(strValor, ".")
		auxValor, err := strconv.Atoi(arrValor[0])
		if err == nil {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), auxValor)
		}
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), aux["descripcion"])
		auxStrString := aux["actividades"].([]interface{})
		var strActividades string
		for j := 0; j < len(auxStrString); j++ {
			if j != len(auxStrString)-1 {
				strActividades = strActividades + " " + auxStrString[j].(string) + ","
			} else {
				strActividades = strActividades + " " + auxStrString[j].(string)
			}
		}
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contador), strActividades)
		SombrearCeldas(consolidadoExcelPlanAnual, i, sheetName, "B"+fmt.Sprint(contador), "F"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, i, sheetName, "D"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
	}
	contador++
	consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "H"+fmt.Sprint(contador))
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Identificación de contratistas")
	consolidadoExcelPlanAnual.SetRowHeight(sheetName, contador+1, 7)
	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "H"+fmt.Sprint(contador), styletitles)

	contador++
	contador++

	consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador))
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Descripción de la necesidad")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), "Perfil")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "Cantidad")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contador), "Valor Total")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contador), "Valor Total Incremeto")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contador), "Actividades")
	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "H"+fmt.Sprint(contador), stylehead)

	contador++
	var total float64 = 0
	var valorTotal int = 0
	var valorTotalInc int = 0
	for i := 0; i < len(contratistas); i++ {
		var respuestaParametro map[string]interface{}
		var perfil map[string]interface{}

		aux := contratistas[i]

		total = total + aux["cantidad"].(float64)
		aux1 := fmt.Sprintf("%v", aux["valorTotal"])
		strValorTotal := strings.TrimLeft(aux1, "$")
		strValorTotal = strings.ReplaceAll(strValorTotal, ",", "")
		arrValorTotal := strings.Split(strValorTotal, ".")
		auxValorTotal, err := strconv.Atoi(arrValorTotal[0])

		if err == nil {
			valorTotal = valorTotal + auxValorTotal
		}
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), aux["descripcionNecesidad"])
		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro/"+fmt.Sprint(aux["perfil"]), &respuestaParametro); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaParametro, &perfil)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), perfil["Nombre"])
		}
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), aux["cantidad"])
		aux2 := fmt.Sprintf("%v", aux["valorTotal"])
		strValor := strings.TrimLeft(aux2, "$")
		strValor = strings.ReplaceAll(strValor, ",", "")
		arrValor := strings.Split(strValor, ".")
		auxValor, err := strconv.Atoi(arrValor[0])
		if err == nil {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contador), auxValor)
		}
		aux3 := fmt.Sprintf("%v", aux["valorTotalInc"])
		strValorInc := strings.TrimLeft(aux3, "$")
		strValorInc = strings.ReplaceAll(strValorInc, ",", "")
		arrValorInc := strings.Split(strValorInc, ".")
		auxValorInc, err := strconv.Atoi(arrValorInc[0])
		if err == nil {
			valorTotalInc = valorTotalInc + auxValorInc
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contador), auxValorInc)
		}
		auxStrString := aux["actividades"].([]interface{})
		var strActividades string
		for j := 0; j < len(auxStrString); j++ {
			if j != len(auxStrString)-1 {
				strActividades = strActividades + " " + auxStrString[j].(string) + ","
			} else {
				strActividades = strActividades + " " + auxStrString[j].(string)
			}
		}
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contador), strActividades)
		SombrearCeldas(consolidadoExcelPlanAnual, i, sheetName, "B"+fmt.Sprint(contador), "H"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, i, sheetName, "F"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		SombrearCeldas(consolidadoExcelPlanAnual, i, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentC, stylecontentCS)
		contador++
	}
	consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador))
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Total")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), total)
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contador), valorTotal)
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contador), valorTotalInc)

	stylecontentTotal, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"D9D9D9"}},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 6},
			{Type: "bottom", Color: "000000", Style: 1}}})
	stylecontentTotalCant, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"D9D9D9"}},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 6},
			{Type: "bottom", Color: "000000", Style: 1}}})
	stylecontentTotalM, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		NumFmt:    183,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"D9D9D9"}},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 6},
			{Type: "bottom", Color: "000000", Style: 1}}})

	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentTotal)
	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentTotalCant)
	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "F"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentTotalM)

	contador++
	consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador))
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Rubro")
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), rubro)
	consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contador), "G"+fmt.Sprint(contador))
	consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), nombreRubro)

	stylecontentRubro, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Font:      &excelize.Font{Bold: true},
		Border: []excelize.Border{{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1}}})
	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador), stylecontentRubro)
	consolidadoExcelPlanAnual.SetCellStyle(sheetName, "D"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentC)

	contador++
	contador++

	if docentes != nil {
		infoDocentes := TotalDocentes(docentes)["TotalesPorTipo"].(TotalesDocentes)
		rubros := docentes["rubros"].([]map[string]interface{})
		rubros_pos := docentes["rubros_pos"].([]map[string]interface{})
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Identificación docente")
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contador+1, 7)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), styletitles)

		contador++
		contador++

		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Categoría")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), "Código del rubro")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), "Nombre del rubro")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "Valor")
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylehead)

		contador++

		//Cuerpo Tabla
		content, _ := ioutil.ReadFile("static/json/rubros.json")
		rubrosJson := []map[string]interface{}{}
		_ = json.Unmarshal(content, &rubrosJson)

		code := ""
		nombre := ""

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Salario básico")
		code = codigoRubrosDocentes(rubros, "Salario básico")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.SalarioBasico+infoDocentes.Rhv_pre.SalarioBasico)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Salario básico")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.SalarioBasico <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.SalarioBasico)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Prima de Servicios")
		code = codigoRubrosDocentes(rubros, "Prima de Servicios")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.PrimaServicios+infoDocentes.Rhv_pre.PrimaServicios)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Prima de Servicios")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.PrimaServicios <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.PrimaServicios)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Prima de navidad")
		code = codigoRubrosDocentes(rubros, "Prima de navidad")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.PrimaNavidad+infoDocentes.Rhv_pre.PrimaNavidad)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Prima de navidad")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.PrimaNavidad <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.PrimaNavidad)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Prima de vacaciones")
		code = codigoRubrosDocentes(rubros, "Prima de vacaciones")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.PrimaVacaciones+infoDocentes.Rhv_pre.PrimaVacaciones)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Prima de vacaciones")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.PrimaVacaciones <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.PrimaVacaciones)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Fondo pensiones público")
		code = codigoRubrosDocentes(rubros, "Fondo pensiones público")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.PensionesPublicas+infoDocentes.Rhv_pre.PensionesPublicas)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Fondo pensiones público")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.PensionesPublicas <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.PensionesPublicas)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Fondo pensiones privado")
		code = codigoRubrosDocentes(rubros, "Fondo pensiones privado")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.PensionesPrivadas+infoDocentes.Rhv_pre.PensionesPrivadas+infoDocentes.Rhv_pos.PensionesPrivadas)
		//consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), infoDocentes.Rhf.PensionesPrivadas+infoDocentes.Rhv_pre.PensionesPrivadas)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Fondo pensiones privado")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		//consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), infoDocentes.Rhv_pos.PensionesPrivadas)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Aporte salud")
		code = codigoRubrosDocentes(rubros, "Aporte salud")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.Salud+infoDocentes.Rhv_pre.Salud)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Aporte salud")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.Salud <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.Salud)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Aporte cesantías público")
		code = codigoRubrosDocentes(rubros, "Aporte cesantías público")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.CesantiasPublicas+infoDocentes.Rhv_pre.CesantiasPublicas)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Aporte cesantías público")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.CesantiasPublicas <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.CesantiasPublicas)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Aporte cesantías privado")
		code = codigoRubrosDocentes(rubros, "Aporte cesantías privado")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.CesantiasPrivadas+infoDocentes.Rhv_pre.CesantiasPrivadas)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Aporte cesantías privado")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.CesantiasPrivadas <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.CesantiasPrivadas)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Aporte CCF")
		code = codigoRubrosDocentes(rubros, "Aporte CCF")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.Caja+infoDocentes.Rhv_pre.Caja)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Aporte CCF")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.Caja <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.Caja)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Aporte ARL")
		code = codigoRubrosDocentes(rubros, "Aporte ARL")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.Arl+infoDocentes.Rhv_pre.Arl)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Aporte ARL")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.Arl <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.Arl)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Aporte ICBF")
		code = codigoRubrosDocentes(rubros, "Aporte ICBF")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhf.Icbf+infoDocentes.Rhv_pre.Icbf)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++
		code = codigoRubrosDocentes(rubros_pos, "Aporte ICBF")
		nombre = NombreRubroByCodigo(rubrosJson, code)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), code)
		if code == "No definido" && infoDocentes.Rhv_pos.Icbf <= 0 {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre+" Posgrado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "N/A")
		} else {
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), nombre)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), infoDocentes.Rhv_pos.Icbf)
		}
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

	}

	consolidadoExcelPlanAnual.InsertRows(sheetName, 1, 7)
	consolidadoExcelPlanAnual.MergeCell(sheetName, "C2", "G6")

	return consolidadoExcelPlanAnual
}

func NombreRubroByCodigo(rubros []map[string]interface{}, codigo string) string {
	nombre := "No definido"
	if codigo == "No definido" {
		return nombre
	}
	for i := 0; i < len(rubros); i++ {
		if rubros[i]["Codigo"] == codigo {
			nombre = rubros[i]["Nombre"].(string)
			break
		}
	}
	return nombre
}

func codigoRubrosDocentes(rubros []map[string]interface{}, categoria string) string {
	var codigo string
	for i := 0; i < len(rubros); i++ {
		rubro := rubros[i]
		if rubro["categoria"] == categoria {
			if _, exist := rubro["rubro"]; exist {
				codigo = rubro["rubro"].(string)
			} else {
				codigo = ""
			}
			break
		}
	}
	if codigo == "" {
		codigo = "No definido"
	}
	return codigo
}

type TotalDocentVal struct {
	SalarioBasico      int
	PrimaServicios     int
	PrimaNavidad       int
	PrimaVacaciones    int
	Bonificacion       int
	PensionesPublicas  int
	PensionesPrivadas  int
	Salud              int
	InteresesCesantias int
	CesantiasPublicas  int
	CesantiasPrivadas  int
	Caja               int
	Arl                int
	Icbf               int
}

type TotalesDocentes struct {
	Rhf     TotalDocentVal
	Rhv_pre TotalDocentVal
	Rhv_pos TotalDocentVal
}

func TotalDocentes(docentes map[string]interface{}) map[string]interface{} {

	var rhf []map[string]interface{}
	var rhvPre []map[string]interface{}
	var rhvPos []map[string]interface{}

	if docentes["rhf"] != "{}" {
		rhf = docentes["rhf"].([]map[string]interface{})
	}
	if docentes["rhv_pre"] != "{}" {
		rhvPre = docentes["rhv_pre"].([]map[string]interface{})
	}
	if docentes["rhv_pos"] != "{}" {
		rhvPos = docentes["rhv_pos"].([]map[string]interface{})
	}

	totalDocentes := make(map[string]interface{})

	totales := TotalesDocentes{}

	sueldoBasico := 0
	primaServicios := 0
	primaNavidad := 0
	primaVacaciones := 0
	bonificacion := 0
	interesesCesantias := 0
	cesantiasPublicas := 0
	cesantiasPrivadas := 0
	salud := 0
	pensionesPublicas := 0
	pensionesPrivadas := 0
	arl := 0
	caja := 0
	icbf := 0

	for i := 0; i < len(rhf); i++ {
		aux := rhf[i]

		if aux["sueldoBasico"] != nil {
			strSueldoBasico := strings.TrimLeft(aux["sueldoBasico"].(string), "$")
			strSueldoBasico = strings.ReplaceAll(strSueldoBasico, ",", "")
			arrSueldoBasico := strings.Split(strSueldoBasico, ".")
			auxSueldoBasico, err := strconv.Atoi(arrSueldoBasico[0])
			if err == nil {
				sueldoBasico += auxSueldoBasico * int(aux["cantidad"].(float64))
				totales.Rhf.SalarioBasico += auxSueldoBasico * int(aux["cantidad"].(float64))
			}
		}

		if aux["primaServicios"] != nil {
			strPrimaServicios := strings.TrimLeft(aux["primaServicios"].(string), "$")
			strPrimaServicios = strings.ReplaceAll(strPrimaServicios, ",", "")
			arrPrimaServicios := strings.Split(strPrimaServicios, ".")
			auxPrimaServicios, err := strconv.Atoi(arrPrimaServicios[0])
			if err == nil {
				primaServicios += auxPrimaServicios
				totales.Rhf.PrimaServicios += auxPrimaServicios
			}
		}

		if aux["primaNavidad"] != nil {
			strPrimaNavidad := strings.TrimLeft(aux["primaNavidad"].(string), "$")
			strPrimaNavidad = strings.ReplaceAll(strPrimaNavidad, ",", "")
			arrPrimaNavidad := strings.Split(strPrimaNavidad, ".")
			auxPrimaNavidad, err := strconv.Atoi(arrPrimaNavidad[0])
			if err == nil {
				primaNavidad += auxPrimaNavidad
				totales.Rhf.PrimaNavidad += auxPrimaNavidad
			}
		}

		if aux["primaVacaciones"] != nil {
			strPrimaVacaciones := strings.TrimLeft(aux["primaVacaciones"].(string), "$")
			strPrimaVacaciones = strings.ReplaceAll(strPrimaVacaciones, ",", "")
			arrPrimaVacaiones := strings.Split(strPrimaVacaciones, ".")
			auxPrimaVacaciones, err := strconv.Atoi(arrPrimaVacaiones[0])
			if err == nil {
				primaVacaciones += auxPrimaVacaciones
				totales.Rhf.PrimaVacaciones += auxPrimaVacaciones
			}
		}

		if aux["bonificacion"] != nil || aux["bonificacion"] != "N/A" {
			strBonificacion := strings.TrimLeft(aux["bonificacion"].(string), "$")
			strBonificacion = strings.ReplaceAll(strBonificacion, ",", "")
			arrBonificacion := strings.Split(strBonificacion, ".")
			auxBonificacion, err := strconv.Atoi(arrBonificacion[0])
			if err == nil {
				bonificacion += auxBonificacion
				totales.Rhf.Bonificacion += auxBonificacion
			}
		}

		if aux["interesesCesantias"] != nil || aux["interesesCesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["interesesCesantias"].(string), "$")
			strInteresesCesantias = strings.ReplaceAll(strInteresesCesantias, ",", "")
			arrInteresesCesantias := strings.Split(strInteresesCesantias, ".")
			auxInteresesCesantias, err := strconv.Atoi(arrInteresesCesantias[0])
			if err == nil {
				interesesCesantias += auxInteresesCesantias
				totales.Rhf.InteresesCesantias += auxInteresesCesantias
			}
		}

		if aux["cesantiasPublico"] != nil {
			strCesantiasPublico := strings.TrimLeft(aux["cesantiasPublico"].(string), "$")
			strCesantiasPublico = strings.ReplaceAll(strCesantiasPublico, ",", "")
			arrCesantiasPublico := strings.Split(strCesantiasPublico, ".")
			auxCesantiasPublico, err := strconv.Atoi(arrCesantiasPublico[0])
			if err == nil {
				cesantiasPublicas += auxCesantiasPublico
				totales.Rhf.CesantiasPublicas += auxCesantiasPublico
			}
		}

		if aux["cesantiasPrivado"] != nil {
			strCesantiasPrivado := strings.TrimLeft(aux["cesantiasPrivado"].(string), "$")
			strCesantiasPrivado = strings.ReplaceAll(strCesantiasPrivado, ",", "")
			arrCesantiasPrivado := strings.Split(strCesantiasPrivado, ".")
			auxCesantiasPrivado, err := strconv.Atoi(arrCesantiasPrivado[0])
			if err == nil {
				cesantiasPrivadas += auxCesantiasPrivado
				totales.Rhf.CesantiasPrivadas += auxCesantiasPrivado
			}
		}

		if aux["totalSalud"] != nil {
			strSalud := strings.TrimLeft(aux["totalSalud"].(string), "$")
			strSalud = strings.ReplaceAll(strSalud, ",", "")
			arrSalud := strings.Split(strSalud, ".")
			auxSalud, err := strconv.Atoi(arrSalud[0])
			if err == nil {
				salud += auxSalud
				totales.Rhf.Salud += auxSalud
			}
		}

		if aux["pensionesPublico"] != nil {
			strPensionesPublicas := strings.TrimLeft(aux["pensionesPublico"].(string), "$")
			strPensionesPublicas = strings.ReplaceAll(strPensionesPublicas, ",", "")
			arrPensionesPublicas := strings.Split(strPensionesPublicas, ".")
			auxPensionesPublicas, err := strconv.Atoi(arrPensionesPublicas[0])
			if err == nil {
				pensionesPublicas += auxPensionesPublicas
				totales.Rhf.PensionesPublicas += auxPensionesPublicas
			}
		}

		if aux["pensionesPrivado"] != nil {
			strPensionesPrivadas := strings.TrimLeft(aux["pensionesPrivado"].(string), "$")
			strPensionesPrivadas = strings.ReplaceAll(strPensionesPrivadas, ",", "")
			arrPensionesPrivadas := strings.Split(strPensionesPrivadas, ".")
			auxPensionesPrivadas, err := strconv.Atoi(arrPensionesPrivadas[0])
			if err == nil {
				pensionesPrivadas += auxPensionesPrivadas
				totales.Rhf.PensionesPrivadas += auxPensionesPrivadas
			}
		}

		if aux["caja"] != nil {
			strCaja := strings.TrimLeft(aux["caja"].(string), "$")
			strCaja = strings.ReplaceAll(strCaja, ",", "")
			arrCaja := strings.Split(strCaja, ".")
			auxCaja, err := strconv.Atoi(arrCaja[0])
			if err == nil {
				caja += auxCaja
				totales.Rhf.Caja += auxCaja
			}

		}

		if aux["totalArl"] != nil {
			strArl := strings.TrimLeft(aux["totalArl"].(string), "$")
			strArl = strings.ReplaceAll(strArl, ",", "")
			arrArl := strings.Split(strArl, ".")
			auxArl, err := strconv.Atoi(arrArl[0])
			if err == nil {
				arl += auxArl
				totales.Rhf.Arl += auxArl
			}
		}

		if aux["icbf"] != nil {
			strIcbf := strings.TrimLeft(aux["icbf"].(string), "$")
			strIcbf = strings.ReplaceAll(strIcbf, ",", "")
			arrIcbf := strings.Split(strIcbf, ".")
			auxIcbf, err := strconv.Atoi(arrIcbf[0])
			if err == nil {
				icbf += auxIcbf
				totales.Rhf.Icbf += auxIcbf
			}
		}
	}

	for i := 0; i < len(rhvPre); i++ {
		aux := rhvPre[i]

		if aux["sueldoBasico"] != nil {
			strSueldoBasico := strings.TrimLeft(aux["sueldoBasico"].(string), "$")
			strSueldoBasico = strings.ReplaceAll(strSueldoBasico, ",", "")
			arrSueldoBasico := strings.Split(strSueldoBasico, ".")
			auxSueldoBasico, err := strconv.Atoi(arrSueldoBasico[0])
			if err == nil {
				sueldoBasico += auxSueldoBasico * int(aux["cantidad"].(float64))
				totales.Rhv_pre.SalarioBasico += auxSueldoBasico * int(aux["cantidad"].(float64))
			}
		}

		if aux["primaServicios"] != nil {
			strPrimaServicios := strings.TrimLeft(aux["primaServicios"].(string), "$")
			strPrimaServicios = strings.ReplaceAll(strPrimaServicios, ",", "")
			arrPrimaServicios := strings.Split(strPrimaServicios, ".")
			auxPrimaServicios, err := strconv.Atoi(arrPrimaServicios[0])
			if err == nil {
				primaServicios += auxPrimaServicios
				totales.Rhv_pre.PrimaServicios += auxPrimaServicios
			}
		}

		if aux["primaNavidad"] != nil {
			strPrimaNavidad := strings.TrimLeft(aux["primaNavidad"].(string), "$")
			strPrimaNavidad = strings.ReplaceAll(strPrimaNavidad, ",", "")
			arrPrimaNavidad := strings.Split(strPrimaNavidad, ".")
			auxPrimaNavidad, err := strconv.Atoi(arrPrimaNavidad[0])
			if err == nil {
				primaNavidad += auxPrimaNavidad
				totales.Rhv_pre.PrimaNavidad += auxPrimaNavidad
			}
		}

		if aux["primaVacaciones"] != nil {
			strPrimaVacaciones := strings.TrimLeft(aux["primaVacaciones"].(string), "$")
			strPrimaVacaciones = strings.ReplaceAll(strPrimaVacaciones, ",", "")
			arrPrimaVacaiones := strings.Split(strPrimaVacaciones, ".")
			auxPrimaVacaciones, err := strconv.Atoi(arrPrimaVacaiones[0])
			if err == nil {
				primaVacaciones += auxPrimaVacaciones
				totales.Rhv_pre.PrimaVacaciones += auxPrimaVacaciones
			}
		}

		if aux["bonificacion"] != nil || aux["bonificacion"] != "N/A" {
			strBonificacion := strings.TrimLeft(aux["bonificacion"].(string), "$")
			strBonificacion = strings.ReplaceAll(strBonificacion, ",", "")
			arrBonificacion := strings.Split(strBonificacion, ".")
			auxBonificacion, err := strconv.Atoi(arrBonificacion[0])
			if err == nil {
				bonificacion += auxBonificacion
				totales.Rhv_pre.Bonificacion += auxBonificacion
			}
		}

		if aux["interesesCesantias"] != nil || aux["interesesCesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["interesesCesantias"].(string), "$")
			strInteresesCesantias = strings.ReplaceAll(strInteresesCesantias, ",", "")
			arrInteresesCesantias := strings.Split(strInteresesCesantias, ".")
			auxInteresesCesantias, err := strconv.Atoi(arrInteresesCesantias[0])
			if err == nil {
				interesesCesantias += auxInteresesCesantias
				totales.Rhv_pre.InteresesCesantias += auxInteresesCesantias
			}
		}

		if aux["cesantiasPublico"] != nil {
			strCesantiasPublico := strings.TrimLeft(aux["cesantiasPublico"].(string), "$")
			strCesantiasPublico = strings.ReplaceAll(strCesantiasPublico, ",", "")
			arrCesantiasPublico := strings.Split(strCesantiasPublico, ".")
			auxCesantiasPublico, err := strconv.Atoi(arrCesantiasPublico[0])
			if err == nil {
				cesantiasPublicas += auxCesantiasPublico
				totales.Rhv_pre.CesantiasPublicas += auxCesantiasPublico
			}
		}

		if aux["cesantiasPrivado"] != nil {
			strCesantiasPrivado := strings.TrimLeft(aux["cesantiasPrivado"].(string), "$")
			strCesantiasPrivado = strings.ReplaceAll(strCesantiasPrivado, ",", "")
			arrCesantiasPrivado := strings.Split(strCesantiasPrivado, ".")
			auxCesantiasPrivado, err := strconv.Atoi(arrCesantiasPrivado[0])
			if err == nil {
				cesantiasPrivadas += auxCesantiasPrivado
				totales.Rhv_pre.CesantiasPrivadas += auxCesantiasPrivado
			}
		}

		if aux["totalSalud"] != nil {
			strSalud := strings.TrimLeft(aux["totalSalud"].(string), "$")
			strSalud = strings.ReplaceAll(strSalud, ",", "")
			arrSalud := strings.Split(strSalud, ".")
			auxSalud, err := strconv.Atoi(arrSalud[0])
			if err == nil {
				salud += auxSalud
				totales.Rhv_pre.Salud += auxSalud
			}
		}

		if aux["pensionesPublico"] != nil {
			strPensionesPublicas := strings.TrimLeft(aux["pensionesPublico"].(string), "$")
			strPensionesPublicas = strings.ReplaceAll(strPensionesPublicas, ",", "")
			arrPensionesPublicas := strings.Split(strPensionesPublicas, ".")
			auxPensionesPublicas, err := strconv.Atoi(arrPensionesPublicas[0])
			if err == nil {
				pensionesPublicas += auxPensionesPublicas
				totales.Rhv_pre.PensionesPublicas += auxPensionesPublicas
			}
		}

		if aux["pensionesPrivado"] != nil {
			strPensionesPrivadas := strings.TrimLeft(aux["pensionesPrivado"].(string), "$")
			strPensionesPrivadas = strings.ReplaceAll(strPensionesPrivadas, ",", "")
			arrPensionesPrivadas := strings.Split(strPensionesPrivadas, ".")
			auxPensionesPrivadas, err := strconv.Atoi(arrPensionesPrivadas[0])
			if err == nil {
				pensionesPrivadas += auxPensionesPrivadas
				totales.Rhv_pre.PensionesPrivadas += auxPensionesPrivadas
			}
		}

		if aux["caja"] != nil {
			strCaja := strings.TrimLeft(aux["caja"].(string), "$")
			strCaja = strings.ReplaceAll(strCaja, ",", "")
			arrCaja := strings.Split(strCaja, ".")
			auxCaja, err := strconv.Atoi(arrCaja[0])
			if err == nil {
				caja += auxCaja
				totales.Rhv_pre.Caja += auxCaja
			}

		}

		if aux["totalArl"] != nil {
			strArl := strings.TrimLeft(aux["totalArl"].(string), "$")
			strArl = strings.ReplaceAll(strArl, ",", "")
			arrArl := strings.Split(strArl, ".")
			auxArl, err := strconv.Atoi(arrArl[0])
			if err == nil {
				arl += auxArl
				totales.Rhv_pre.Arl += auxArl
			}
		}

		if aux["icbf"] != nil {
			strIcbf := strings.TrimLeft(aux["icbf"].(string), "$")
			strIcbf = strings.ReplaceAll(strIcbf, ",", "")
			arrIcbf := strings.Split(strIcbf, ".")
			auxIcbf, err := strconv.Atoi(arrIcbf[0])
			if err == nil {
				icbf += auxIcbf
				totales.Rhv_pre.Icbf += auxIcbf
			}
		}
	}

	for i := 0; i < len(rhvPos); i++ {
		aux := rhvPos[i]

		if aux["sueldoBasico"] != nil {
			strSueldoBasico := strings.TrimLeft(aux["sueldoBasico"].(string), "$")
			strSueldoBasico = strings.ReplaceAll(strSueldoBasico, ",", "")
			arrSueldoBasico := strings.Split(strSueldoBasico, ".")
			auxSueldoBasico, err := strconv.Atoi(arrSueldoBasico[0])
			if err == nil {
				sueldoBasico += auxSueldoBasico * int(aux["cantidad"].(float64))
				totales.Rhv_pos.SalarioBasico += auxSueldoBasico * int(aux["cantidad"].(float64))
			}
		}

		if aux["primaServicios"] != nil {
			strPrimaServicios := strings.TrimLeft(aux["primaServicios"].(string), "$")
			strPrimaServicios = strings.ReplaceAll(strPrimaServicios, ",", "")
			arrPrimaServicios := strings.Split(strPrimaServicios, ".")
			auxPrimaServicios, err := strconv.Atoi(arrPrimaServicios[0])
			if err == nil {
				primaServicios += auxPrimaServicios
				totales.Rhv_pos.PrimaServicios += auxPrimaServicios
			}
		}

		if aux["primaNavidad"] != nil {
			strPrimaNavidad := strings.TrimLeft(aux["primaNavidad"].(string), "$")
			strPrimaNavidad = strings.ReplaceAll(strPrimaNavidad, ",", "")
			arrPrimaNavidad := strings.Split(strPrimaNavidad, ".")
			auxPrimaNavidad, err := strconv.Atoi(arrPrimaNavidad[0])
			if err == nil {
				primaNavidad += auxPrimaNavidad
				totales.Rhv_pos.PrimaNavidad += auxPrimaNavidad
			}
		}

		if aux["primaVacaciones"] != nil {
			strPrimaVacaciones := strings.TrimLeft(aux["primaVacaciones"].(string), "$")
			strPrimaVacaciones = strings.ReplaceAll(strPrimaVacaciones, ",", "")
			arrPrimaVacaiones := strings.Split(strPrimaVacaciones, ".")
			auxPrimaVacaciones, err := strconv.Atoi(arrPrimaVacaiones[0])
			if err == nil {
				primaVacaciones += auxPrimaVacaciones
				totales.Rhv_pos.PrimaVacaciones += auxPrimaVacaciones
			}
		}

		if aux["bonificacion"] != nil || aux["bonificacion"] != "N/A" {
			strBonificacion := strings.TrimLeft(aux["bonificacion"].(string), "$")
			strBonificacion = strings.ReplaceAll(strBonificacion, ",", "")
			arrBonificacion := strings.Split(strBonificacion, ".")
			auxBonificacion, err := strconv.Atoi(arrBonificacion[0])
			if err == nil {
				bonificacion += auxBonificacion
				totales.Rhv_pos.Bonificacion += auxBonificacion
			}
		}

		if aux["interesesCesantias"] != nil || aux["interesesCesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["interesesCesantias"].(string), "$")
			strInteresesCesantias = strings.ReplaceAll(strInteresesCesantias, ",", "")
			arrInteresesCesantias := strings.Split(strInteresesCesantias, ".")
			auxInteresesCesantias, err := strconv.Atoi(arrInteresesCesantias[0])
			if err == nil {
				interesesCesantias += auxInteresesCesantias
				totales.Rhv_pos.InteresesCesantias += auxInteresesCesantias
			}
		}

		if aux["cesantiasPublico"] != nil {
			strCesantiasPublico := strings.TrimLeft(aux["cesantiasPublico"].(string), "$")
			strCesantiasPublico = strings.ReplaceAll(strCesantiasPublico, ",", "")
			arrCesantiasPublico := strings.Split(strCesantiasPublico, ".")
			auxCesantiasPublico, err := strconv.Atoi(arrCesantiasPublico[0])
			if err == nil {
				cesantiasPublicas += auxCesantiasPublico
				totales.Rhv_pos.CesantiasPublicas += auxCesantiasPublico
			}
		}

		if aux["cesantiasPrivado"] != nil {
			strCesantiasPrivado := strings.TrimLeft(aux["cesantiasPrivado"].(string), "$")
			strCesantiasPrivado = strings.ReplaceAll(strCesantiasPrivado, ",", "")
			arrCesantiasPrivado := strings.Split(strCesantiasPrivado, ".")
			auxCesantiasPrivado, err := strconv.Atoi(arrCesantiasPrivado[0])
			if err == nil {
				cesantiasPrivadas += auxCesantiasPrivado
				totales.Rhv_pos.CesantiasPrivadas += auxCesantiasPrivado
			}
		}

		if aux["totalSalud"] != nil {
			strSalud := strings.TrimLeft(aux["totalSalud"].(string), "$")
			strSalud = strings.ReplaceAll(strSalud, ",", "")
			arrSalud := strings.Split(strSalud, ".")
			auxSalud, err := strconv.Atoi(arrSalud[0])
			if err == nil {
				salud += auxSalud
				totales.Rhv_pos.Salud += auxSalud
			}
		}

		if aux["pensionesPublico"] != nil {
			strPensionesPublicas := strings.TrimLeft(aux["pensionesPublico"].(string), "$")
			strPensionesPublicas = strings.ReplaceAll(strPensionesPublicas, ",", "")
			arrPensionesPublicas := strings.Split(strPensionesPublicas, ".")
			auxPensionesPublicas, err := strconv.Atoi(arrPensionesPublicas[0])
			if err == nil {
				pensionesPublicas += auxPensionesPublicas
				totales.Rhv_pos.PensionesPublicas += auxPensionesPublicas
			}
		}

		if aux["pensionesPrivado"] != nil {
			strPensionesPrivadas := strings.TrimLeft(aux["pensionesPrivado"].(string), "$")
			strPensionesPrivadas = strings.ReplaceAll(strPensionesPrivadas, ",", "")
			arrPensionesPrivadas := strings.Split(strPensionesPrivadas, ".")
			auxPensionesPrivadas, err := strconv.Atoi(arrPensionesPrivadas[0])
			if err == nil {
				pensionesPrivadas += auxPensionesPrivadas
				totales.Rhv_pos.PensionesPrivadas += auxPensionesPrivadas
			}
		}

		if aux["caja"] != nil {
			strCaja := strings.TrimLeft(aux["caja"].(string), "$")
			strCaja = strings.ReplaceAll(strCaja, ",", "")
			arrCaja := strings.Split(strCaja, ".")
			auxCaja, err := strconv.Atoi(arrCaja[0])
			if err == nil {
				caja += auxCaja
				totales.Rhv_pos.Caja += auxCaja
			}

		}

		if aux["totalArl"] != nil {
			strArl := strings.TrimLeft(aux["totalArl"].(string), "$")
			strArl = strings.ReplaceAll(strArl, ",", "")
			arrArl := strings.Split(strArl, ".")
			auxArl, err := strconv.Atoi(arrArl[0])
			if err == nil {
				arl += auxArl
				totales.Rhv_pos.Arl += auxArl
			}
		}

		if aux["icbf"] != nil {
			strIcbf := strings.TrimLeft(aux["icbf"].(string), "$")
			strIcbf = strings.ReplaceAll(strIcbf, ",", "")
			arrIcbf := strings.Split(strIcbf, ".")
			auxIcbf, err := strconv.Atoi(arrIcbf[0])
			if err == nil {
				icbf += auxIcbf
				totales.Rhv_pos.Icbf += auxIcbf
			}
		}
	}

	totalDocentes["sueldoBasico"] = sueldoBasico
	totalDocentes["primaServicios"] = primaServicios
	totalDocentes["primaNavidad"] = primaNavidad
	totalDocentes["primaVacaciones"] = primaVacaciones
	totalDocentes["bonificacion"] = bonificacion
	totalDocentes["interesesCesantias"] = interesesCesantias
	totalDocentes["cesantiasPublicas"] = cesantiasPublicas
	totalDocentes["cesantiasPrivadas"] = cesantiasPrivadas
	totalDocentes["salud"] = salud
	totalDocentes["pensionesPublicas"] = pensionesPublicas
	totalDocentes["pensionesPrivadas"] = pensionesPrivadas
	totalDocentes["arl"] = arl
	totalDocentes["caja"] = caja
	totalDocentes["icbf"] = icbf
	totalDocentes["TotalesPorTipo"] = totales

	return totalDocentes
}

func GetDataDocentes(docentes map[string]interface{}, dependencia_id string) map[string]interface{} {

	var respuestaDependencia map[string]interface{}
	dataDocentes := make(map[string]interface{})

	dataDocentes["tco"] = 0
	dataDocentes["mto"] = 0
	dataDocentes["hch"] = 0
	dataDocentes["hcp"] = 0
	dataDocentes["hchPos"] = 0
	dataDocentes["hcpPos"] = 0
	dataDocentes["valorPre"] = 0
	dataDocentes["valorPos"] = 0

	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+dependencia_id, &respuestaDependencia); err == nil {
		dataDocentes["nombreFacultad"] = respuestaDependencia["Nombre"]
	}

	if docentes["rhf"] != nil && docentes["rhf"] != "{}" {
		rhf := docentes["rhf"].([]map[string]interface{})
		for i := 0; i < len(rhf); i++ {
			if rhf[i]["tipo"] == "Tiempo Completo" {
				dataDocentes["tco"] = dataDocentes["tco"].(int) + 1
			}
			if rhf[i]["tipo"] == "Medio Tiempo" {
				dataDocentes["mto"] = dataDocentes["mto"].(int) + 1
			}
			strTotal := strings.TrimLeft(rhf[i]["total"].(string), "$")
			strTotal = strings.ReplaceAll(strTotal, ",", "")
			arrTotal := strings.Split(strTotal, ".")
			auxTotal, err := strconv.Atoi(arrTotal[0])
			if err == nil {
				dataDocentes["valorPre"] = dataDocentes["valorPre"].(int) + auxTotal
			}

		}
	}

	if docentes["rhv_pre"] != nil && docentes["rhv_pre"] != "{}" {
		rhvPre := docentes["rhv_pre"].([]map[string]interface{})
		for i := 0; i < len(rhvPre); i++ {
			if rhvPre[i]["tipo"] == "H. Catedra Honorarios" {
				dataDocentes["hch"] = dataDocentes["hch"].(int) + 1
			}
			if rhvPre[i]["tipo"] == "H. Catedra Prestacional" {
				dataDocentes["hcp"] = dataDocentes["hcp"].(int) + 1
			}
			strTotal := strings.TrimLeft(rhvPre[i]["total"].(string), "$")
			strTotal = strings.ReplaceAll(strTotal, ",", "")
			arrTotal := strings.Split(strTotal, ".")
			auxTotal, err := strconv.Atoi(arrTotal[0])
			if err == nil {
				dataDocentes["valorPre"] = dataDocentes["valorPre"].(int) + auxTotal
			}
		}
	}

	if docentes["rhv_pos"] != nil && docentes["rhv_pos"] != "{}" {
		rhvPos := docentes["rhv_pos"].([]map[string]interface{})
		for i := 0; i < len(rhvPos); i++ {
			if rhvPos[i]["tipo"] == "H. Catedra Honorarios" {
				dataDocentes["hchPos"] = dataDocentes["hchPos"].(int) + 1
			}
			if rhvPos[i]["tipo"] == "H. Catedra Prestacional" {
				dataDocentes["hcpPos"] = dataDocentes["hcpPos"].(int) + 1
			}
			strTotal := strings.TrimLeft(rhvPos[i]["total"].(string), "$")
			strTotal = strings.ReplaceAll(strTotal, ",", "")
			arrTotal := strings.Split(strTotal, ".")
			auxTotal, err := strconv.Atoi(arrTotal[0])
			if err == nil {
				dataDocentes["valorPos"] = dataDocentes["valorPos"].(int) + auxTotal
			}
		}

	}

	return dataDocentes

}

func SombrearCeldas(excel *excelize.File, idActividad int, sheetName string, hCell string, vCell string, style int, styleSombreado int) {
	if idActividad%2 == 0 {
		excel.SetCellStyle(sheetName, hCell, vCell, style)
	} else {
		excel.SetCellStyle(sheetName, hCell, vCell, styleSombreado)
	}
}

func Convert2Num(value interface{}) interface{} {
	switch value.(type) {
	case float64:
		return value.(float64)
	case string:
		num, _ := strconv.ParseFloat(value.(string), 64)
		return num
	default:
		return "-"
	}
}

func ConstruirExcelPlanAccionUnidad(esReporteAntiguo bool, datosReporte map[string]interface{}) (*excelize.File, error) {
	planes := datosReporte["planes"].(int)
	plan_id := datosReporte["plan_id"].(string)
	periodo := datosReporte["periodo"].([]map[string]interface{})
	planesFilter := datosReporte["planesFilter"].([]map[string]interface{})
	arregloPlanAnual := datosReporte["arregloPlanAnual"].([]map[string]interface{})

	unidadNombre := arregloPlanAnual[0]["nombreUnidad"].(string)
	sheetName := "Actividades del plan"
	consolidadoExcelPlanAnual := excelize.NewFile()
	indexPlan, _ := consolidadoExcelPlanAnual.NewSheet(sheetName)

	if len(planesFilter) <= 0 {
		return nil, errors.New("no hay planes")
	}

	if planes == 0 {
		consolidadoExcelPlanAnual.DeleteSheet("Sheet1")

		disable := false
		if err := consolidadoExcelPlanAnual.SetSheetView(sheetName, -1, &excelize.ViewOptions{
			ShowGridLines: &disable,
		}); err != nil {
			fmt.Println(err)
		}
	}

	//? Estilos
	styletitle, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"},
		Font:      &excelize.Font{Bold: true, Size: 18, Color: "000000"},
		Border: []excelize.Border{
			{Type: "right", Color: "ffffff", Style: 1},
			{Type: "left", Color: "ffffff", Style: 1},
			{Type: "top", Color: "ffffff", Style: 1},
			{Type: "bottom", Color: "ffffff", Style: 1},
		},
	})

	stylehead, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styletitles, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontent, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentC, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCL, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 4},
		},
	})
	stylecontentCLD, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCLS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 4},
		},
	})
	stylecontentCLDS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleLineamiento, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:   "center",
			Vertical:     "center",
			WrapText:     true,
			TextRotation: 90,
		},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleLineamientoSombra, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:   "center",
			Vertical:     "center",
			WrapText:     true,
			TextRotation: 90,
		},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if esReporteAntiguo {
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B1", "D1")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "E1", "G1")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "H1", "H2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "I1", "I2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "J1", "J2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "K1", "K2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "L1", "L2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "Q1", "Q2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "M1", "P1")
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 18)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "C", 11)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "D", "D", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "E", 16)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "F", "G", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "H", 6)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "J", 12)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "K", "K", 30)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "L", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "M", "N", 52)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "O", "O", 30)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "P", "P", 10)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "Q", "Q", 30)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B1", "Q1", stylehead)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B2", "Q2", styletitles)

		// encabezado excel
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B1", "Armonización PED")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B2", "Lineamiento")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C2", "Meta")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D2", "Estrategias")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E1", "Armonización Plan Indicativo")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E2", "Ejes transformadores")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "F2", "Lineamientos de acción")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "G2", "Estrategias")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "H2", "N°.")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "I2", "Ponderación de la actividad")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "J2", "Periodo de ejecución")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "K2", "Actividad")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Actividades específicas")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "M1", "Indicador")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "M2", "Nombre")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "N2", "Fórmula")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "O2", "Criterio del indicador")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "P2", "Meta")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "Q2", "Producto esperado")

		rowPos := 3

		for excelPlan := 0; excelPlan < len(arregloPlanAnual); excelPlan++ {
			datosExcelPlan := arregloPlanAnual[excelPlan]
			armoPED := datosExcelPlan["datosArmonizacion"].([]map[string]interface{})
			armoPI := datosExcelPlan["datosArmonizacionPI"].([]map[string]interface{})
			datosComplementarios := datosExcelPlan["datosComplementarios"].(map[string]interface{})
			indicadores := datosComplementarios["indicadores"].(map[string]interface{})

			MaxRowsXActivity := MinComMul_Armonization(armoPED, armoPI, len(indicadores))

			y_lin := rowPos
			h_lin := MaxRowsXActivity / len(armoPED)

			for _, lin := range armoPED {
				consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1), styleLineamiento, styleLineamientoSombra)
				y_met := y_lin
				h_met := h_lin / len(lin["meta"].([]map[string]interface{}))
				for _, met := range lin["meta"].([]map[string]interface{}) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(y_met), met["nombreMeta"])
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1), stylecontentC, stylecontentCS)
					y_est := y_met
					h_est := h_met / len(met["estrategias"].([]map[string]interface{}))
					for _, est := range met["estrategias"].([]map[string]interface{}) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(y_est), est["descripcionEstrategia"])
						if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
						} else {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
						}
						y_est += h_est
					}
					y_met += h_met
				}
				y_lin += h_lin
			}

			y_eje := rowPos
			h_eje := MaxRowsXActivity / len(armoPI)
			for _, eje := range armoPI {
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(y_eje), eje["nombreFactor"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1), stylecontentC, stylecontentCS)
				y_lin := y_eje
				h_lin := h_eje / len(eje["lineamientos"].([]map[string]interface{}))
				for _, lin := range eje["lineamientos"].([]map[string]interface{}) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1), stylecontentC, stylecontentCS)
					y_est := y_lin
					h_est := h_lin / len(lin["estrategias"].([]map[string]interface{}))
					for _, est := range lin["estrategias"].([]map[string]interface{}) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(y_est), est["descripcionEstrategia"])
						if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
						} else {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
						}
						y_est += h_est
					}
					y_lin += h_lin
				}
				y_eje += h_eje
			}

			consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(rowPos), "H"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(rowPos), "I"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(rowPos), "K"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(rowPos), excelPlan+1)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(rowPos), datosComplementarios["Ponderación de la actividad"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(rowPos), datosComplementarios["Periodo de ejecución"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(rowPos), datosComplementarios["Actividad general"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(rowPos), datosComplementarios["Tareas"])
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontent, stylecontentS)

			y_ind := rowPos
			h_ind := MaxRowsXActivity / len(indicadores)
			idx := int(0)
			var indicadoresVacios []map[string]interface{}
			var indicadoresNoVacios []map[string]interface{}
			var indicadoresOrdenados []map[string]interface{}

			for _, indicador := range indicadores {
				auxIndicador := indicador
				var nombreIndicador interface{}
				var formula interface{}
				var criterio interface{}
				var meta interface{}
				for key, element := range auxIndicador.(map[string]interface{}) {
					if strings.Contains(strings.ToLower(key), "nombre") {
						nombreIndicador = element
					}
					if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
						formula = element
					}
					if strings.Contains(strings.ToLower(key), "criterio") {
						criterio = element
					}
					if strings.Contains(strings.ToLower(key), "meta") {
						meta = element
					}
				}
				if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
					indicadoresVacios = append(indicadoresVacios, auxIndicador.(map[string]interface{}))
				} else {
					indicadoresNoVacios = append(indicadoresNoVacios, auxIndicador.(map[string]interface{}))
				}
			}

			indicadoresOrdenados = append(indicadoresNoVacios, indicadoresVacios...)
			for _, indicador := range indicadoresOrdenados {
				auxIndicador := indicador
				var nombreIndicador interface{}
				var formula interface{}
				var criterio interface{}
				var meta interface{}
				for key, element := range auxIndicador {
					if strings.Contains(strings.ToLower(key), "nombre") {
						nombreIndicador = element
					}
					if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
						formula = element
					}
					if strings.Contains(strings.ToLower(key), "criterio") {
						criterio = element
					}
					if strings.Contains(strings.ToLower(key), "meta") {
						meta = element
					}
				}
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(y_ind), nombreIndicador)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(y_ind), formula)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(y_ind), criterio)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(y_ind), meta)
				if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind-1), "M"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind-1), "N"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind-1), "O"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(y_ind-1), "P"+fmt.Sprint(y_ind+h_ind-1))
				} else {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind), "M"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind), "N"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(y_ind), "P"+fmt.Sprint(y_ind+h_ind-1))
				}

				idx++
				if idx < len(indicadores) {
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(y_ind), "P"+fmt.Sprint(y_ind+h_ind-1), stylecontentCL, stylecontentCLS)
				} else {
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(y_ind), "P"+fmt.Sprint(y_ind+h_ind-1), stylecontentCLD, stylecontentCLDS)
				}
				y_ind += h_ind
			}

			consolidadoExcelPlanAnual.MergeCell(sheetName, "Q"+fmt.Sprint(rowPos), "Q"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "Q"+fmt.Sprint(rowPos), datosComplementarios["Producto esperado"])
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "Q"+fmt.Sprint(rowPos), "Q"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)

			rowPos += MaxRowsXActivity

			consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
		}
	} else {
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B1", "D1")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "E1", "G1")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "H1", "H2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "I1", "I2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "J1", "J2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "K1", "K2")
		// consolidadoExcelPlanAnual.MergeCell(sheetName, "L1", "L2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "P1", "P2")
		consolidadoExcelPlanAnual.MergeCell(sheetName, "L1", "O1")
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 18)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "C", 11)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "D", "D", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "E", 16)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "F", "G", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "H", 6)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "J", 12)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "K", "K", 30)
		// consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "L", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "M", 52)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "N", "N", 30)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "O", "O", 10)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "P", "P", 30)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B1", "P1", stylehead)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B2", "P2", styletitles)

		// encabezado excel
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B1", "Armonización PED")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B2", "Lineamiento")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C2", "Meta")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D2", "Estrategias")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E1", "Armonización Plan Indicativo")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E2", "Ejes transformadores")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "F2", "Lineamientos de acción")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "G2", "Estrategias")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "H2", "N°.")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "I2", "Ponderación de la actividad")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "J2", "Periodo de ejecución")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "K2", "Actividad")
		// consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Actividades específicas")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "L1", "Indicador")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Nombre")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "M2", "Fórmula")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "N2", "Criterio del indicador")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "O2", "Meta")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "P2", "Producto esperado")

		rowPos := 3

		for excelPlan := 0; excelPlan < len(arregloPlanAnual); excelPlan++ {

			datosExcelPlan := arregloPlanAnual[excelPlan]
			armoPED := datosExcelPlan["datosArmonizacion"].([]map[string]interface{})
			armoPI := datosExcelPlan["datosArmonizacionPI"].([]map[string]interface{})
			datosComplementarios := datosExcelPlan["datosComplementarios"].(map[string]interface{})
			indicadores := datosComplementarios["indicadores"].(map[string]interface{})

			MaxRowsXActivity := MinComMul_Armonization(armoPED, armoPI, len(indicadores))

			y_lin := rowPos
			h_lin := MaxRowsXActivity / len(armoPED)

			for _, lin := range armoPED {
				consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1), styleLineamiento, styleLineamientoSombra)
				y_met := y_lin
				h_met := h_lin / len(lin["meta"].([]map[string]interface{}))
				for _, met := range lin["meta"].([]map[string]interface{}) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(y_met), met["nombreMeta"])
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1), stylecontentC, stylecontentCS)
					y_est := y_met
					h_est := h_met / len(met["estrategias"].([]map[string]interface{}))
					for _, est := range met["estrategias"].([]map[string]interface{}) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(y_est), est["descripcionEstrategia"])
						if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
						} else {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
						}
						y_est += h_est
					}
					y_met += h_met
				}
				y_lin += h_lin
			}

			y_eje := rowPos
			h_eje := MaxRowsXActivity / len(armoPI)
			for _, eje := range armoPI {
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(y_eje), eje["nombreFactor"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1), stylecontentC, stylecontentCS)
				y_lin := y_eje
				h_lin := h_eje / len(eje["lineamientos"].([]map[string]interface{}))
				for _, lin := range eje["lineamientos"].([]map[string]interface{}) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1), stylecontentC, stylecontentCS)
					y_est := y_lin
					h_est := h_lin / len(lin["estrategias"].([]map[string]interface{}))
					for _, est := range lin["estrategias"].([]map[string]interface{}) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(y_est), est["descripcionEstrategia"])
						if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
						} else {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
						}
						y_est += h_est
					}
					y_lin += h_lin
				}
				y_eje += h_eje
			}

			consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(rowPos), "H"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(rowPos), "I"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(rowPos), "K"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(rowPos), excelPlan+1)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(rowPos), datosComplementarios["Ponderación de la actividad"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(rowPos), datosComplementarios["Periodo de ejecución"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(rowPos), datosExcelPlan["nombreActividad"].(string))
			// consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(rowPos), datosComplementarios["Tareas"])
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontent, stylecontentS)

			y_ind := rowPos
			h_ind := MaxRowsXActivity / len(indicadores)
			idx := int(0)
			var indicadoresVacios []map[string]interface{}
			var indicadoresNoVacios []map[string]interface{}
			var indicadoresOrdenados []map[string]interface{}

			for _, indicador := range indicadores {
				auxIndicador := indicador
				var nombreIndicador interface{}
				var formula interface{}
				var criterio interface{}
				var meta interface{}
				for key, element := range auxIndicador.(map[string]interface{}) {
					if strings.Contains(strings.ToLower(key), "nombre") {
						nombreIndicador = element
					}
					if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
						formula = element
					}
					if strings.Contains(strings.ToLower(key), "criterio") {
						criterio = element
					}
					if strings.Contains(strings.ToLower(key), "meta") {
						meta = element
					}
				}
				if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
					indicadoresVacios = append(indicadoresVacios, auxIndicador.(map[string]interface{}))
				} else {
					indicadoresNoVacios = append(indicadoresNoVacios, auxIndicador.(map[string]interface{}))
				}
			}

			indicadoresOrdenados = append(indicadoresNoVacios, indicadoresVacios...)
			for _, indicador := range indicadoresOrdenados {
				auxIndicador := indicador
				var nombreIndicador interface{}
				var formula interface{}
				var criterio interface{}
				var meta interface{}
				for key, element := range auxIndicador {
					if strings.Contains(strings.ToLower(key), "nombre") {
						nombreIndicador = element
					}
					if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
						formula = element
					}
					if strings.Contains(strings.ToLower(key), "criterio") {
						criterio = element
					}
					if strings.Contains(strings.ToLower(key), "meta") {
						meta = element
					}
				}
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(y_ind), nombreIndicador)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(y_ind), formula)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(y_ind), criterio)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(y_ind), meta)
				if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(y_ind-1), "L"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind-1), "M"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind-1), "N"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind-1), "O"+fmt.Sprint(y_ind+h_ind-1))
				} else {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(y_ind), "L"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind), "M"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind), "N"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1))
				}

				idx++
				if idx < len(indicadores) {
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "L"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1), stylecontentCL, stylecontentCLS)
				} else {
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "L"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1), stylecontentCLD, stylecontentCLDS)
				}
				y_ind += h_ind
			}

			consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(rowPos), "P"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(rowPos), datosComplementarios["Producto esperado"])
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(rowPos), "P"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)

			rowPos += MaxRowsXActivity

			consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
		}
	}

	consolidadoExcelPlanAnual = TablaIdentificaciones(consolidadoExcelPlanAnual, plan_id)

	consolidadoExcelPlanAnual.InsertRows("Actividades del plan", 1, 7)
	consolidadoExcelPlanAnual.MergeCell("Actividades del plan", "C2", "P6")
	consolidadoExcelPlanAnual.SetCellStyle("Actividades del plan", "C2", "P6", styletitle)
	consolidadoExcelPlanAnual.SetCellStyle("Identificaciones", "C2", "G6", styletitle)

	if periodo[0] != nil {
		consolidadoExcelPlanAnual.SetCellValue("Actividades del plan", "C2", "Plan de Acción "+periodo[0]["Nombre"].(string)+"\n"+unidadNombre)
		consolidadoExcelPlanAnual.SetCellValue("Identificaciones", "C2", "Proyección de necesidades "+periodo[0]["Nombre"].(string)+"\n"+unidadNombre)
	} else {
		consolidadoExcelPlanAnual.SetCellValue("Actividades del plan", "C2", "Plan de Acción")
		consolidadoExcelPlanAnual.SetCellValue("Identificaciones", "C2", "Proyección de necesidades")
	}

	if err := consolidadoExcelPlanAnual.AddPicture("Actividades del plan", "B1", "static/img/UDEscudo2.png",
		&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "oneCell", OffsetX: 10}); err != nil {
		fmt.Println(err)
	}
	if err := consolidadoExcelPlanAnual.AddPicture("Identificaciones", "B1", "static/img/UDEscudo2.png",
		&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "absolute", OffsetX: 10}); err != nil {
		fmt.Println(err)
	}

	consolidadoExcelPlanAnual.SetColWidth("Actividades del plan", "A", "A", 2)

	return consolidadoExcelPlanAnual, nil
}

func ConstruirExcelPlanAccionGeneral(esReporteAntiguo bool, datosReporte map[string]interface{}) (*excelize.File, error) {
	var resPeriodo map[string]interface{}
	var periodo []map[string]interface{}
	contadorGeneral := 4
	rowPos := contadorGeneral + 5
	vigencia_id := datosReporte["vigencia_id"].(string)
	planesFilter := datosReporte["planesFilter"].([]map[string]interface{})
	arregloPlanAnual := datosReporte["arregloPlanAnual"].([]map[string]interface{})

	unidadNombre := arregloPlanAnual[0]["nombreUnidad"].(string)
	consolidadoExcelPlanAnual := excelize.NewFile()
	sheetName := "REPORTE GENERAL"
	indexPlan, _ := consolidadoExcelPlanAnual.NewSheet(sheetName)

	if len(planesFilter) <= 0 {
		return nil, errors.New("no hay planes")
	}

	consolidadoExcelPlanAnual.DeleteSheet("Sheet1")
	consolidadoExcelPlanAnual.InsertCols("REPORTE GENERAL", "A", 1)
	disable := false
	if err := consolidadoExcelPlanAnual.SetSheetView(sheetName, -1, &excelize.ViewOptions{
		ShowGridLines: &disable,
	}); err != nil {
		fmt.Println(err)
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Id:`+vigencia_id, &resPeriodo); err == nil {
		helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
	}

	// if planes == 0 {
		
	// }

	//? Estilos
	styletitle, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"},
		Font:      &excelize.Font{Bold: true, Size: 18, Color: "000000"},
		Border: []excelize.Border{
			{Type: "right", Color: "ffffff", Style: 1},
			{Type: "left", Color: "ffffff", Style: 1},
			{Type: "top", Color: "ffffff", Style: 1},
			{Type: "bottom", Color: "ffffff", Style: 1},
		},
	})
	stylehead, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styletitles, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontent, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentC, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleLineamiento, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:   "center",
			Vertical:     "center",
			WrapText:     true,
			TextRotation: 90,
		},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleLineamientoSombra, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:   "center",
			Vertical:     "center",
			WrapText:     true,
			TextRotation: 90,
		},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCL, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 4},
		},
	})
	stylecontentCLD, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCLS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 4},
		},
	})
	stylecontentCLDS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	if esReporteAntiguo {
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "P"+fmt.Sprint(contadorGeneral+1))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "D"+fmt.Sprint(contadorGeneral+2))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorGeneral+2), "G"+fmt.Sprint(contadorGeneral+2))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorGeneral+2), "H"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorGeneral+2), "I"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorGeneral+2), "J"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorGeneral+2), "K"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(contadorGeneral+2), "L"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "Q"+fmt.Sprint(contadorGeneral+2), "Q"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(contadorGeneral+2), "P"+fmt.Sprint(contadorGeneral+2))
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+1, 20)
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+2, 20)
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 20)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 19)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "P", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "C", 13)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "E", 16)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "H", 6)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "J", 12)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "K", "K", 30)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "L", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "M", "N", 52)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "O", "O", 30)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "P", "P", 10)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "Q", "Q", 30)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "Q"+fmt.Sprint(contadorGeneral+1), stylehead)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "Q"+fmt.Sprint(contadorGeneral+2), stylehead)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "Q"+fmt.Sprint(contadorGeneral+3), styletitles)
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 30)

		var tituloExcel string
		if periodo[0] != nil {
			tituloExcel = "Plan de acción " + periodo[0]["Nombre"].(string) + " - " + unidadNombre
		} else {
			tituloExcel = "Plan de acción - " + unidadNombre
		}

		// encabezado excel
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+1), tituloExcel)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "Armonización PED")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "Lineamiento")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorGeneral+3), "Meta")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorGeneral+3), "Estrategias")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorGeneral+2), "Armonización Plan Indicativo")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorGeneral+3), "Ejes transformadores")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorGeneral+3), "Lineamientos de acción")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorGeneral+3), "Estrategias")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorGeneral+3), "N°.")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorGeneral+3), "Ponderación de la actividad")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorGeneral+3), "Periodo de ejecución")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorGeneral+3), "Actividad")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+3), "Actividades específicas")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorGeneral+2), "Indicador")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorGeneral+3), "Nombre")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorGeneral+3), "Fórmula")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(contadorGeneral+3), "Criterio del indicador")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(contadorGeneral+3), "Meta")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "Q"+fmt.Sprint(contadorGeneral+3), "Producto esperado")
		consolidadoExcelPlanAnual.InsertRows(sheetName, 1, 1)

		for excelPlan := 0; excelPlan < len(arregloPlanAnual); excelPlan++ {

			datosExcelPlan := arregloPlanAnual[excelPlan]
			armoPED := datosExcelPlan["datosArmonizacion"].([]map[string]interface{})
			armoPI := datosExcelPlan["datosArmonizacionPI"].([]map[string]interface{})
			datosComplementarios := datosExcelPlan["datosComplementarios"].(map[string]interface{})
			indicadores := datosComplementarios["indicadores"].(map[string]interface{})

			MaxRowsXActivity := MinComMul_Armonization(armoPED, armoPI, len(indicadores))

			y_lin := rowPos
			h_lin := MaxRowsXActivity / len(armoPED)
			for _, lin := range armoPED {
				consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1), styleLineamiento, styleLineamientoSombra)
				y_met := y_lin
				h_met := h_lin / len(lin["meta"].([]map[string]interface{}))
				for _, met := range lin["meta"].([]map[string]interface{}) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(y_met), met["nombreMeta"])
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1), stylecontentC, stylecontentCS)
					y_est := y_met
					h_est := h_met / len(met["estrategias"].([]map[string]interface{}))
					for _, est := range met["estrategias"].([]map[string]interface{}) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(y_est), est["descripcionEstrategia"])
						if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
						} else {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
						}
						y_est += h_est
					}
					y_met += h_met
				}
				y_lin += h_lin
			}

			y_eje := rowPos
			h_eje := MaxRowsXActivity / len(armoPI)
			for _, eje := range armoPI {
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(y_eje), eje["nombreFactor"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1), stylecontentC, stylecontentCS)
				y_lin := y_eje
				h_lin := h_eje / len(eje["lineamientos"].([]map[string]interface{}))
				for _, lin := range eje["lineamientos"].([]map[string]interface{}) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1), stylecontentC, stylecontentCS)
					y_est := y_lin
					h_est := h_lin / len(lin["estrategias"].([]map[string]interface{}))
					for _, est := range lin["estrategias"].([]map[string]interface{}) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(y_est), est["descripcionEstrategia"])
						if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
						} else {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
						}
						y_est += h_est
					}
					y_lin += h_lin
				}
				y_eje += h_eje
			}

			consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(rowPos), "H"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(rowPos), "I"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(rowPos), "K"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(rowPos), excelPlan+1)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(rowPos), datosComplementarios["Ponderación de la actividad"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(rowPos), datosComplementarios["Periodo de ejecución"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(rowPos), datosComplementarios["Actividad general"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(rowPos), datosComplementarios["Tareas"])
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontent, stylecontentS)

			y_ind := rowPos
			h_ind := MaxRowsXActivity / len(indicadores)
			idx := int(0)

			var indicadoresVacios []map[string]interface{}
			var indicadoresNoVacios []map[string]interface{}
			var indicadoresOrdenados []map[string]interface{}

			for _, indicador := range indicadores {
				auxIndicador := indicador
				var nombreIndicador interface{}
				var formula interface{}
				var criterio interface{}
				var meta interface{}
				for key, element := range auxIndicador.(map[string]interface{}) {
					if strings.Contains(strings.ToLower(key), "nombre") {
						nombreIndicador = element
					}
					if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
						formula = element
					}
					if strings.Contains(strings.ToLower(key), "criterio") {
						criterio = element
					}
					if strings.Contains(strings.ToLower(key), "meta") {
						meta = element
					}
				}
				if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
					indicadoresVacios = append(indicadoresVacios, auxIndicador.(map[string]interface{}))
				} else {
					indicadoresNoVacios = append(indicadoresNoVacios, auxIndicador.(map[string]interface{}))
				}
			}

			indicadoresOrdenados = append(indicadoresNoVacios, indicadoresVacios...)
			for _, indicador := range indicadoresOrdenados {
				auxIndicador := indicador
				var nombreIndicador interface{}
				var formula interface{}
				var criterio interface{}
				var meta interface{}
				for key, element := range auxIndicador {
					if strings.Contains(strings.ToLower(key), "nombre") {
						nombreIndicador = element
					}
					if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
						formula = element
					}
					if strings.Contains(strings.ToLower(key), "criterio") {
						criterio = element
					}
					if strings.Contains(strings.ToLower(key), "meta") {
						meta = element
					}
				}
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(y_ind), nombreIndicador)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(y_ind), formula)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(y_ind), criterio)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(y_ind), meta)
				if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind-1), "M"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind-1), "N"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind-1), "O"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(y_ind-1), "P"+fmt.Sprint(y_ind+h_ind-1))
				} else {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind), "M"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind), "N"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(y_ind), "P"+fmt.Sprint(y_ind+h_ind-1))
				}

				idx++
				if idx < len(indicadores) {
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(y_ind), "P"+fmt.Sprint(y_ind+h_ind-1), stylecontentCL, stylecontentCLS)
				} else {
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(y_ind), "P"+fmt.Sprint(y_ind+h_ind-1), stylecontentCLD, stylecontentCLDS)
				}
				y_ind += h_ind
			}

			consolidadoExcelPlanAnual.MergeCell(sheetName, "Q"+fmt.Sprint(rowPos), "Q"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "Q"+fmt.Sprint(rowPos), datosComplementarios["Producto esperado"])
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "Q"+fmt.Sprint(rowPos), "Q"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)

			rowPos += MaxRowsXActivity

			contadorGeneral = rowPos - 2

			consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
		}
		arregloPlanAnual = nil
		consolidadoExcelPlanAnual.RemoveRow(sheetName, 1)
	} else {
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "P"+fmt.Sprint(contadorGeneral+1))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "D"+fmt.Sprint(contadorGeneral+2))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorGeneral+2), "G"+fmt.Sprint(contadorGeneral+2))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorGeneral+2), "H"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorGeneral+2), "I"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorGeneral+2), "J"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorGeneral+2), "K"+fmt.Sprint(contadorGeneral+3))
		// consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(contadorGeneral+2), "L"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(contadorGeneral+2), "P"+fmt.Sprint(contadorGeneral+3))
		consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(contadorGeneral+2), "O"+fmt.Sprint(contadorGeneral+2))
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+1, 20)
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+2, 20)
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 20)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 19)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "P", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "C", 13)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "E", 16)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "H", 6)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "J", 12)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "K", "K", 30)
		// consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "L", 35)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "M", 52)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "N", "N", 30)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "O", "O", 10)
		consolidadoExcelPlanAnual.SetColWidth(sheetName, "P", "P", 30)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "P"+fmt.Sprint(contadorGeneral+1), stylehead)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "P"+fmt.Sprint(contadorGeneral+2), stylehead)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "P"+fmt.Sprint(contadorGeneral+3), styletitles)
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 30)

		var tituloExcel string
		if periodo[0] != nil {
			tituloExcel = "Plan de acción " + periodo[0]["Nombre"].(string) + " - " + unidadNombre
		} else {
			tituloExcel = "Plan de acción - " + unidadNombre
		}

		// encabezado excel
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+1), tituloExcel)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "Armonización PED")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "Lineamiento")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorGeneral+3), "Meta")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorGeneral+3), "Estrategias")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorGeneral+2), "Armonización Plan Indicativo")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorGeneral+3), "Ejes transformadores")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorGeneral+3), "Lineamientos de acción")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorGeneral+3), "Estrategias")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorGeneral+3), "N°.")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorGeneral+3), "Ponderación de la actividad")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorGeneral+3), "Periodo de ejecución")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorGeneral+3), "Actividad")
		// consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+3), "Actividades específicas")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+2), "Indicador")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+3), "Nombre")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorGeneral+3), "Fórmula")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorGeneral+3), "Criterio del indicador")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(contadorGeneral+3), "Meta")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(contadorGeneral+3), "Producto esperado")
		consolidadoExcelPlanAnual.InsertRows(sheetName, 1, 1)

		for excelPlan := 0; excelPlan < len(arregloPlanAnual); excelPlan++ {

			datosExcelPlan := arregloPlanAnual[excelPlan]
			armoPED := datosExcelPlan["datosArmonizacion"].([]map[string]interface{})
			armoPI := datosExcelPlan["datosArmonizacionPI"].([]map[string]interface{})
			datosComplementarios := datosExcelPlan["datosComplementarios"].(map[string]interface{})
			indicadores := datosComplementarios["indicadores"].(map[string]interface{})

			MaxRowsXActivity := MinComMul_Armonization(armoPED, armoPI, len(indicadores))

			y_lin := rowPos
			h_lin := MaxRowsXActivity / len(armoPED)
			for _, lin := range armoPED {
				consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1), styleLineamiento, styleLineamientoSombra)
				y_met := y_lin
				h_met := h_lin / len(lin["meta"].([]map[string]interface{}))
				for _, met := range lin["meta"].([]map[string]interface{}) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(y_met), met["nombreMeta"])
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1), stylecontentC, stylecontentCS)
					y_est := y_met
					h_est := h_met / len(met["estrategias"].([]map[string]interface{}))
					for _, est := range met["estrategias"].([]map[string]interface{}) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(y_est), est["descripcionEstrategia"])
						if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
						} else {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
						}
						y_est += h_est
					}
					y_met += h_met
				}
				y_lin += h_lin
			}

			y_eje := rowPos
			h_eje := MaxRowsXActivity / len(armoPI)
			for _, eje := range armoPI {
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(y_eje), eje["nombreFactor"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1), stylecontentC, stylecontentCS)
				y_lin := y_eje
				h_lin := h_eje / len(eje["lineamientos"].([]map[string]interface{}))
				for _, lin := range eje["lineamientos"].([]map[string]interface{}) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1), stylecontentC, stylecontentCS)
					y_est := y_lin
					h_est := h_lin / len(lin["estrategias"].([]map[string]interface{}))
					for _, est := range lin["estrategias"].([]map[string]interface{}) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(y_est), est["descripcionEstrategia"])
						if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
						} else {
							SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
						}
						y_est += h_est
					}
					y_lin += h_lin
				}
				y_eje += h_eje
			}

			consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(rowPos), "H"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(rowPos), "I"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(rowPos), "K"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			// consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(rowPos), excelPlan+1)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(rowPos), datosComplementarios["Ponderación de la actividad"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(rowPos), datosComplementarios["Periodo de ejecución"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(rowPos), datosExcelPlan["nombreActividad"])
			// consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(rowPos), datosComplementarios["Tareas"])
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(rowPos), "K"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontent, stylecontentS)

			y_ind := rowPos
			h_ind := MaxRowsXActivity / len(indicadores)
			idx := int(0)

			var indicadoresVacios []map[string]interface{}
			var indicadoresNoVacios []map[string]interface{}
			var indicadoresOrdenados []map[string]interface{}

			for _, indicador := range indicadores {
				auxIndicador := indicador
				var nombreIndicador interface{}
				var formula interface{}
				var criterio interface{}
				var meta interface{}
				for key, element := range auxIndicador.(map[string]interface{}) {
					if strings.Contains(strings.ToLower(key), "nombre") {
						nombreIndicador = element
					}
					if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
						formula = element
					}
					if strings.Contains(strings.ToLower(key), "criterio") {
						criterio = element
					}
					if strings.Contains(strings.ToLower(key), "meta") {
						meta = element
					}
				}
				if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
					indicadoresVacios = append(indicadoresVacios, auxIndicador.(map[string]interface{}))
				} else {
					indicadoresNoVacios = append(indicadoresNoVacios, auxIndicador.(map[string]interface{}))
				}
			}

			indicadoresOrdenados = append(indicadoresNoVacios, indicadoresVacios...)
			for _, indicador := range indicadoresOrdenados {
				auxIndicador := indicador
				var nombreIndicador interface{}
				var formula interface{}
				var criterio interface{}
				var meta interface{}
				for key, element := range auxIndicador {
					if strings.Contains(strings.ToLower(key), "nombre") {
						nombreIndicador = element
					}
					if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
						formula = element
					}
					if strings.Contains(strings.ToLower(key), "criterio") {
						criterio = element
					}
					if strings.Contains(strings.ToLower(key), "meta") {
						meta = element
					}
				}
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(y_ind), nombreIndicador)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(y_ind), formula)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(y_ind), criterio)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(y_ind), meta)
				if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(y_ind-1), "L"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind-1), "M"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind-1), "N"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind-1), "O"+fmt.Sprint(y_ind+h_ind-1))
				} else {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(y_ind), "L"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind), "M"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind), "N"+fmt.Sprint(y_ind+h_ind-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1))
				}

				idx++
				if idx < len(indicadores) {
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "L"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1), stylecontentCL, stylecontentCLS)
				} else {
					SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "L"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1), stylecontentCLD, stylecontentCLDS)
				}
				y_ind += h_ind
			}

			consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(rowPos), "P"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(rowPos), datosComplementarios["Producto esperado"])
			SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(rowPos), "P"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)

			rowPos += MaxRowsXActivity

			contadorGeneral = rowPos - 2

			consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
		}
		arregloPlanAnual = nil
		consolidadoExcelPlanAnual.RemoveRow(sheetName, 1)
	}
	consolidadoExcelPlanAnual.InsertRows("REPORTE GENERAL", 1, 3)
	consolidadoExcelPlanAnual.MergeCell("REPORTE GENERAL", "C2", "P6")
	consolidadoExcelPlanAnual.SetCellStyle("REPORTE GENERAL", "C2", "P6", styletitle)
	if periodo[0] != nil {
		consolidadoExcelPlanAnual.SetCellValue("REPORTE GENERAL", "C2", "Plan de Acción Anual "+periodo[0]["Nombre"].(string)+"\nUniversidad Distrital Franciso José de Caldas")
	} else {
		consolidadoExcelPlanAnual.SetCellValue("REPORTE GENERAL", "C2", "Plan de Acción Anual \nUniversidad Distrital Franciso José de Caldas")
	}

	if err := consolidadoExcelPlanAnual.AddPicture("REPORTE GENERAL", "B1", "static/img/UDEscudo2.png",
		&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "oneCell", OffsetX: 10}); err != nil {
		fmt.Println(err)
	}

	consolidadoExcelPlanAnual.SetColWidth("REPORTE GENERAL", "A", "A", 2)

	return consolidadoExcelPlanAnual, nil
}
