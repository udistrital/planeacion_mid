package reporteshelper

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"sort"
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
	validDataT = []string{}
	ids = [][]string{}
	hijos_data = nil
	hijos_key = nil
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
	var actividades []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+subgrupo_id+"&fields=dato_plan", &res); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(res, &aux)
		subgrupoDetalle := aux[0]
		if subgrupoDetalle["dato_plan"] != nil {
			var datoPlan map[string]interface{}
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
		if !bandera {
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
			if !bandera {
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

func TablaIdentificaciones(consolidadoExcelPlanAnual *excelize.File, planId string, esReporteAntiguo bool) *excelize.File {
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
	return construirTablas(consolidadoExcelPlanAnual, recursos, contratistas, docentes, rubro, nombreRubro, esReporteAntiguo)
}

func construirTablas(consolidadoExcelPlanAnual *excelize.File, recursos []map[string]interface{}, contratistas []map[string]interface{}, docentes map[string]interface{}, rubro string, nombreRubro string, esReporteAntiguo bool) *excelize.File {
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
	stylecontentRubro, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Font:      &excelize.Font{Bold: true},
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
	if esReporteAntiguo {
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
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contador), "Valor Total Incremento")
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

		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentTotal)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "E"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentTotalCant)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "F"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentTotalM)

		contador++
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Rubro")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), rubro)
		consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contador), "G"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), nombreRubro)

		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador), stylecontentRubro)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "D"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentC)

		contador++
		contador++
	} else {
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "I"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Identificación de contratistas")
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), styletitles)

		contador++
		contador++

		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Descripción de la necesidad (objeto contractual)")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), "Equipo")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), "Requisitos mínimos (formación académica y experiencia)")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contador), "Perfil")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contador), "Cantidad")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contador), "Valor Total")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contador), "Valor Total Incremento")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contador), "Actividades")
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "J"+fmt.Sprint(contador), stylehead)

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
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), aux["equipoResponsable"])
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), aux["requisitos"])
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro/"+fmt.Sprint(aux["perfil"]), &respuestaParametro); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuestaParametro, &perfil)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contador), perfil["Nombre"])
			}
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contador), aux["cantidad"])
			aux2 := fmt.Sprintf("%v", aux["valorTotal"])
			strValor := strings.TrimLeft(aux2, "$")
			strValor = strings.ReplaceAll(strValor, ",", "")
			arrValor := strings.Split(strValor, ".")
			auxValor, err := strconv.Atoi(arrValor[0])
			if err == nil {
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contador), auxValor)
			}
			aux3 := fmt.Sprintf("%v", aux["valorTotalInc"])
			strValorInc := strings.TrimLeft(aux3, "$")
			strValorInc = strings.ReplaceAll(strValorInc, ",", "")
			arrValorInc := strings.Split(strValorInc, ".")
			auxValorInc, err := strconv.Atoi(arrValorInc[0])
			if err == nil {
				valorTotalInc = valorTotalInc + auxValorInc
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contador), auxValorInc)
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
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contador), strActividades)
			SombrearCeldas(consolidadoExcelPlanAnual, i, sheetName, "B"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), stylecontent, stylecontentS)
			SombrearCeldas(consolidadoExcelPlanAnual, i, sheetName, "G"+fmt.Sprint(contador), "H"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
			SombrearCeldas(consolidadoExcelPlanAnual, i, sheetName, "F"+fmt.Sprint(contador), "F"+fmt.Sprint(contador), stylecontentC, stylecontentCS)
			contador++
		}
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Total")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contador), total)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contador), valorTotal)
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contador), valorTotalInc)

		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontentTotal)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "G"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentTotalCant)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "H"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), stylecontentTotalM)

		contador++
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Rubro")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), rubro)
		consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contador), "H"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contador), nombreRubro)

		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontentRubro)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "E"+fmt.Sprint(contador), "H"+fmt.Sprint(contador), stylecontentC)

		contador++
		contador++
	}
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
		content, err := os.ReadFile("static/json/rubros.json")
		if err != nil {
			beego.Error("error leyendo archivo rubros.json:", err)
			return nil
		}

		rubrosJson := []map[string]interface{}{}
		if err := json.Unmarshal(content, &rubrosJson); err != nil {
			beego.Error("error decodificando rubros.json:", err)
			return nil
		}

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

	_ = consolidadoExcelPlanAnual.InsertRows(sheetName, 1, 7)
	_ = consolidadoExcelPlanAnual.MergeCell(sheetName, "C2", "G6")

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

		if aux["cesantias"] != nil || aux["cesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["cesantias"].(string), "$")
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

		if aux["cesantias"] != nil || aux["cesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["cesantias"].(string), "$")
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

		if aux["cesantias"] != nil || aux["cesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["cesantias"].(string), "$")
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
	totalDocentes["cesantias"] = interesesCesantias
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

func ConstruirExcelPlanAccionUnidad(planesFilter []map[string]interface{}, body map[string]interface{}) (*excelize.File, []map[string]interface{}, error) {
	var res map[string]interface{}
	var respuestaUnidad []map[string]interface{}
	var subgrupos []map[string]interface{}
	var plan_id string
	var actividadName string
	var arregloPlanAnual []map[string]interface{}
	var nombreUnidad string
	var resPeriodo map[string]interface{}
	var periodo []map[string]interface{}
	var unidadNombre string
	consolidadoExcelPlanAnual := excelize.NewFile()

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

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Id:`+body["vigencia"].(string), &resPeriodo); err == nil {
		helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
	}

	anioPlan := periodo[0]["Year"].(float64)
	var esReporteAntiguo bool = true

	if anioPlan > 2024 { //? Vigencias 2025 en adelante son reportes con nueva estructura
		esReporteAntiguo = false
	}
	if len(planesFilter) <= 0 {
		return nil, nil, errors.New("no hay planes")
	}

	for planes := 0; planes < len(planesFilter); planes++ {
		planesFilterData := planesFilter[planes]
		plan_id = planesFilterData["_id"].(string)
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id+"&fields=nombre,_id,hijos,activo", &res); err != nil {
			panic(err)
		}
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
		if esReporteAntiguo {
			for i := 0; i < len(subgrupos); i++ {
				if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
					actividades := GetActividades(subgrupos[i]["_id"].(string))
					var arregloLineamieto []map[string]interface{}
					var arregloLineamietoPI []map[string]interface{}
					if len(actividades) == 1 {
						for index := range actividades {
							if val, ok := actividades[index]["index"].(float64); ok {
								actividades[index]["index"] = fmt.Sprintf("%v", int(val))
							}
						}
					} else {
						sort.SliceStable(actividades, func(i int, j int) bool {
							if _, ok := actividades[i]["index"].(float64); ok {
								actividades[i]["index"] = fmt.Sprintf("%v", int(actividades[i]["index"].(float64)))
							}
							if _, ok := actividades[j]["index"].(float64); ok {
								actividades[j]["index"] = fmt.Sprintf("%v", int(actividades[j]["index"].(float64)))
							}
							aux, _ := strconv.Atoi((actividades[i]["index"]).(string))
							aux1, _ := strconv.Atoi((actividades[j]["index"]).(string))
							return aux < aux1
						})
					}

					LimpiarDetalles()
					for j := 0; j < len(actividades); j++ {
						arregloLineamieto = nil
						arregloLineamietoPI = nil
						actividad := actividades[j]
						actividadName = actividad["dato"].(string)
						index := actividad["index"].(string)
						datosArmonizacion := make(map[string]interface{})
						titulosArmonizacion := make(map[string]interface{})

						Limpia()
						tree := BuildTreeFa(subgrupos, index)
						treeDatos := tree[0]
						treeDatas := tree[1]
						treeArmo := tree[2]
						armonizacionTercer := treeArmo[0]
						var armonizacionTercerNivel interface{}
						var armonizacionTercerNivelPI interface{}

						if armonizacionTercer["armo"] != nil {
							armonizacionTercerNivel = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPED"]
							armonizacionTercerNivelPI = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPI"]
						}

						for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
							treeDato := treeDatos[datoGeneral]
							treeData := treeDatas[0]
							if treeDato["sub"] == "" {
								nombre := strings.ToLower(treeDato["nombre"].(string))
								if strings.Contains(nombre, "ponderación") || strings.Contains(nombre, "ponderacion") && strings.Contains(nombre, "actividad") {
									datosArmonizacion["Ponderación de la actividad"] = treeData[treeDato["id"].((string))]
								} else if strings.Contains(nombre, "período") || strings.Contains(nombre, "periodo") && strings.Contains(nombre, "ejecucion") || strings.Contains(nombre, "ejecución") {
									datosArmonizacion["Periodo de ejecución"] = treeData[treeDato["id"].(string)]
								} else if strings.Contains(nombre, "actividad") && strings.Contains(nombre, "general") {
									datosArmonizacion["Actividad general"] = treeData[treeDato["id"].(string)]
								} else if strings.Contains(nombre, "tarea") || strings.Contains(nombre, "actividades específicas") {
									datosArmonizacion["Tareas"] = treeData[treeDato["id"].(string)]
								} else {
									datosArmonizacion[treeDato["nombre"].(string)] = treeData[treeDato["id"].(string)]
								}
							}
						}
						var treeIndicador map[string]interface{}
						auxTree := tree[0]
						for i := 0; i < len(auxTree); i++ {
							subgrupo := auxTree[i]
							if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") {
								treeIndicador = auxTree[i]
							}
						}

						subIndicador := treeIndicador["sub"].([]map[string]interface{})
						for ind := 0; ind < len(subIndicador); ind++ {
							subIndicadorRes := subIndicador[ind]
							treeData := treeDatas[0]
							dataIndicador := make(map[string]interface{})
							auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
							for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
								if treeData[auxSubIndicador[subInd]["id"].(string)] == nil {
									treeData[auxSubIndicador[subInd]["id"].(string)] = ""
								}
								dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[auxSubIndicador[subInd]["id"].(string)]
							}
							titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
						}

						datosArmonizacion["indicadores"] = titulosArmonizacion
						arregloLineamieto = ArbolArmonizacionV2(armonizacionTercerNivel.(string))
						arregloLineamietoPI = ArbolArmonizacionPIV2(armonizacionTercerNivelPI.(string))
						generalData := make(map[string]interface{})
						if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+body["unidad_id"].(string), &respuestaUnidad); err == nil {
							aux := respuestaUnidad[0]
							dependenciaNombre := aux["DependenciaId"].(map[string]interface{})
							nombreUnidad = dependenciaNombre["Nombre"].(string)
						} else {
							panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
						}

						generalData["nombreUnidad"] = nombreUnidad
						generalData["nombreActividad"] = actividadName
						generalData["numeroActividad"] = index
						generalData["datosArmonizacion"] = arregloLineamieto
						generalData["datosArmonizacionPI"] = arregloLineamietoPI
						generalData["datosComplementarios"] = datosArmonizacion

						arregloPlanAnual = append(arregloPlanAnual, generalData)
					}
					break
				}
			}
			unidadNombre = arregloPlanAnual[0]["nombreUnidad"].(string)
			sheetName := "Actividades del plan"
			indexPlan, _ := consolidadoExcelPlanAnual.NewSheet(sheetName)

			if planes == 0 {
				_ = consolidadoExcelPlanAnual.DeleteSheet("Sheet1")

				disable := false
				if err := consolidadoExcelPlanAnual.SetSheetView(sheetName, -1, &excelize.ViewOptions{
					ShowGridLines: &disable,
				}); err != nil {
					fmt.Println(err)
				}
			}
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
		} else { //? Es reporte nuevo
			for i := 0; i < len(subgrupos); i++ {
				if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") {
					actividades := GetActividades(subgrupos[i]["_id"].(string))
					var arregloLineamieto []map[string]interface{}
					var arregloLineamietoPI []map[string]interface{}
					if len(actividades) == 1 {
						for index := range actividades {
							if val, ok := actividades[index]["index"].(float64); ok {
								actividades[index]["index"] = fmt.Sprintf("%v", int(val))
							}
						}
					} else {
						sort.SliceStable(actividades, func(i int, j int) bool {
							if _, ok := actividades[i]["index"].(float64); ok {
								actividades[i]["index"] = fmt.Sprintf("%v", int(actividades[i]["index"].(float64)))
							}
							if _, ok := actividades[j]["index"].(float64); ok {
								actividades[j]["index"] = fmt.Sprintf("%v", int(actividades[j]["index"].(float64)))
							}
							aux, _ := strconv.Atoi((actividades[i]["index"]).(string))
							aux1, _ := strconv.Atoi((actividades[j]["index"]).(string))
							return aux < aux1
						})
					}

					LimpiarDetalles()
					for j := 0; j < len(actividades); j++ {
						arregloLineamieto = nil
						arregloLineamietoPI = nil
						actividad := actividades[j]
						actividadName = actividad["dato"].(string)
						index := actividad["index"].(string)
						datosArmonizacion := make(map[string]interface{})
						titulosArmonizacion := make(map[string]interface{})

						Limpia()
						tree := BuildTreeFa(subgrupos, index)
						treeDatos := tree[0]
						treeDatas := tree[1]
						treeArmo := tree[2]
						armonizacionTercer := treeArmo[0]
						var armonizacionTercerNivel interface{}
						var armonizacionTercerNivelPI interface{}

						if armonizacionTercer["armo"] != nil {
							armonizacionTercerNivel = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPED"]
							armonizacionTercerNivelPI = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPI"]
						}

						for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
							treeDato := treeDatos[datoGeneral]
							treeData := treeDatas[0]
							if treeDato["sub"] == "" {
								nombre := strings.ToLower(treeDato["nombre"].(string))
								if strings.Contains(nombre, "ponderación") || strings.Contains(nombre, "ponderacion") && strings.Contains(nombre, "actividad") {
									datosArmonizacion["Ponderación de la actividad"] = treeData[treeDato["id"].((string))]
								} else if strings.Contains(nombre, "período") || strings.Contains(nombre, "periodo") && strings.Contains(nombre, "ejecucion") || strings.Contains(nombre, "ejecución") {
									datosArmonizacion["Periodo de ejecución"] = treeData[treeDato["id"].(string)]
								} else if strings.Contains(nombre, "actividad") {
									datosArmonizacion["Actividad general"] = treeData[treeDato["id"].(string)]
								} else if strings.Contains(nombre, "unidad") || strings.Contains(nombre, "grupo") {
									datosArmonizacion["Responsable"] = treeData[treeDato["id"].(string)]
								} else {
									datosArmonizacion[treeDato["nombre"].(string)] = treeData[treeDato["id"].(string)]
								}
							}
						}
						var treeIndicador map[string]interface{}
						auxTree := tree[0]
						for i := 0; i < len(auxTree); i++ {
							subgrupo := auxTree[i]
							if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") {
								treeIndicador = auxTree[i]
							}
						}

						subIndicador := treeIndicador["sub"].([]map[string]interface{})
						// for ind := 0; ind < len(subIndicador); ind++ {
						// 	subIndicadorRes := subIndicador[ind]
						// 	treeData := treeDatas[0]
						// 	dataIndicador := make(map[string]interface{})
						// 	auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
						// 	for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
						// 		if treeData[auxSubIndicador[subInd]["id"].(string)] == nil {
						// 			treeData[auxSubIndicador[subInd]["id"].(string)] = ""
						// 		}
						// 		dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[auxSubIndicador[subInd]["id"].(string)]
						// 	}
						// 	titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
						// }
						for ind := 0; ind < len(subIndicador); ind++ {
							subIndicadorRes := subIndicador[ind]
							treeData := treeDatas[0]
							dataIndicador := make(map[string]interface{})
							auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
							for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
								dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[auxSubIndicador[subInd]["id"].(string)]
							}
							titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
						}

						datosArmonizacion["indicadores"] = titulosArmonizacion
						arregloLineamieto = ArbolArmonizacionV2(armonizacionTercerNivel.(string))
						arregloLineamietoPI = ArbolArmonizacionPIV2(armonizacionTercerNivelPI.(string))

						generalData := make(map[string]interface{})
						if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+body["unidad_id"].(string), &respuestaUnidad); err == nil {
							aux := respuestaUnidad[0]
							dependenciaNombre := aux["DependenciaId"].(map[string]interface{})
							nombreUnidad = dependenciaNombre["Nombre"].(string)
						} else {
							panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
						}

						generalData["nombreUnidad"] = nombreUnidad
						generalData["nombreActividad"] = actividadName
						generalData["numeroActividad"] = index
						generalData["datosArmonizacion"] = arregloLineamieto
						generalData["datosArmonizacionPI"] = arregloLineamietoPI
						generalData["datosComplementarios"] = datosArmonizacion

						arregloPlanAnual = append(arregloPlanAnual, generalData)
					}
					break
				}
			}
			unidadNombre = arregloPlanAnual[0]["nombreUnidad"].(string)
			sheetName := "Actividades del plan"
			indexPlan, _ := consolidadoExcelPlanAnual.NewSheet(sheetName)

			if planes == 0 {
				_ = consolidadoExcelPlanAnual.DeleteSheet("Sheet1")

				disable := false
				if err := consolidadoExcelPlanAnual.SetSheetView(sheetName, -1, &excelize.ViewOptions{
					ShowGridLines: &disable,
				}); err != nil {
					fmt.Println(err)
				}
			}
			consolidadoExcelPlanAnual.MergeCell(sheetName, "B1", "D1")
			consolidadoExcelPlanAnual.MergeCell(sheetName, "E1", "G1")
			consolidadoExcelPlanAnual.MergeCell(sheetName, "H1", "H2")
			consolidadoExcelPlanAnual.MergeCell(sheetName, "I1", "I2")
			consolidadoExcelPlanAnual.MergeCell(sheetName, "J1", "J2")
			consolidadoExcelPlanAnual.MergeCell(sheetName, "K1", "K2")
			consolidadoExcelPlanAnual.MergeCell(sheetName, "P1", "P2")
			consolidadoExcelPlanAnual.MergeCell(sheetName, "Q1", "Q2")
			consolidadoExcelPlanAnual.MergeCell(sheetName, "L1", "O1")
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 18)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "C", 11)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "D", "D", 35)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "E", 16)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "F", "G", 35)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "H", 6)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "J", 12)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "K", "K", 30)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "M", 52)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "N", "N", 30)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "O", "O", 10)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "P", "P", 25)
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
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I2", "Peso (%)")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J2", "Periodo de ejecución")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K2", "Actividad")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L1", "Indicador")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Nombre")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "M2", "Fórmula")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "N2", "Criterio del indicador")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "O2", "Meta")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "P2", "Producto esperado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "Q2", "Unidad o grupo responsable")

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
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(rowPos), excelPlan+1)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(rowPos), datosComplementarios["Ponderación de la actividad"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(rowPos), datosComplementarios["Periodo de ejecución"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(rowPos), datosExcelPlan["nombreActividad"].(string))
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
					if (nombreIndicador == "" && formula == "" && criterio == "" && meta == "") || (nombreIndicador == nil && formula == nil && criterio == nil && meta == nil) {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(y_ind-1), "L"+fmt.Sprint(y_ind+h_ind-1))
						consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind-1), "M"+fmt.Sprint(y_ind+h_ind-1))
						consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind-1), "N"+fmt.Sprint(y_ind+h_ind-1))
						consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind-1), "O"+fmt.Sprint(y_ind+h_ind-1))
					} else {
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(y_ind), nombreIndicador)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(y_ind), formula)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(y_ind), criterio)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(y_ind), meta)
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
				consolidadoExcelPlanAnual.MergeCell(sheetName, "Q"+fmt.Sprint(rowPos), "Q"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "Q"+fmt.Sprint(rowPos), datosComplementarios["Responsable"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(rowPos), "P"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "Q"+fmt.Sprint(rowPos), "Q"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)

				rowPos += MaxRowsXActivity

				consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
			}
		}
	}
	consolidadoExcelPlanAnual = TablaIdentificaciones(consolidadoExcelPlanAnual, plan_id, esReporteAntiguo)
	_ = consolidadoExcelPlanAnual.InsertRows("Actividades del plan", 1, 7)
	_ = consolidadoExcelPlanAnual.MergeCell("Actividades del plan", "C2", "P6")
	_ = consolidadoExcelPlanAnual.SetCellStyle("Actividades del plan", "C2", "P6", styletitle)
	_ = consolidadoExcelPlanAnual.SetCellStyle("Identificaciones", "C2", "G6", styletitle)

	if periodo[0] != nil {
		_ = consolidadoExcelPlanAnual.SetCellValue("Actividades del plan", "C2", "Plan de Acción "+periodo[0]["Nombre"].(string)+"\n"+unidadNombre)
		_ = consolidadoExcelPlanAnual.SetCellValue("Identificaciones", "C2", "Proyección de necesidades "+periodo[0]["Nombre"].(string)+"\n"+unidadNombre)
	} else {
		_ = consolidadoExcelPlanAnual.SetCellValue("Actividades del plan", "C2", "Plan de Acción")
		_ = consolidadoExcelPlanAnual.SetCellValue("Identificaciones", "C2", "Proyección de necesidades")
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
	return consolidadoExcelPlanAnual, arregloPlanAnual, nil
}

func ConstruirExcelPlanAccionGeneral(planesFilter []map[string]interface{}, body map[string]interface{}) (*excelize.File, []map[string]interface{}, error) {
	var res map[string]interface{}
	var respuestaUnidad []map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaTipoPlan map[string]interface{}
	var estado map[string]interface{}
	var tipoPlan map[string]interface{}
	var subgrupos []map[string]interface{}
	var plan_id string
	var actividadName string
	var arregloPlanAnual []map[string]interface{}
	var arregloInfoReportes []map[string]interface{}
	var nombreUnidad string
	var idUnidad string
	var resPeriodo map[string]interface{}
	var periodo []map[string]interface{}
	var esReporteAntiguo bool = true
	contadorGeneral := 4
	consolidadoExcelPlanAnual := excelize.NewFile()

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

	for _, planes := range planesFilter {
		if idUnidad != planes["dependencia_id"].(string) {
			if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+planes["dependencia_id"].(string), &respuestaUnidad); err == nil {
				planes["nombreUnidad"] = respuestaUnidad[0]["DependenciaId"].(map[string]interface{})["Nombre"].(string)
				fmt.Sprintf("http://" + beego.AppConfig.String("OikosService") + "/dependencia_tipo_dependencia?query=DependenciaId:" + planes["dependencia_id"].(string))
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}
		}
	}

	for _, planes := range planesFilter {
		if idUnidad != planes["dependencia_id"].(string) {
			if err := request.GetJson(
				"http://"+beego.AppConfig.String("OikosService")+
					"/dependencia_tipo_dependencia?query=DependenciaId:"+planes["dependencia_id"].(string),
				&respuestaUnidad,
			); err == nil {
				planes["nombreUnidad"] = respuestaUnidad[0]["DependenciaId"].(map[string]interface{})["Nombre"].(string)
			} else {
				panic(map[string]interface{}{
					"funcion": "GetUnidades",
					"err":     "Error al consultar Oikos",
					"status":  "400",
					"log":     err,
				})
			}
		}
	}

	sort.SliceStable(planesFilter, func(i, j int) bool {
		a := (planesFilter)[i]["nombreUnidad"].(string)
		b := (planesFilter)[j]["nombreUnidad"].(string)
		return a < b
	})

	consolidadoExcelPlanAnual.InsertCols("REPORTE GENERAL", "A", 1)
	disable := false
	sheetName := "REPORTE GENERAL"
	indexPlan, _ := consolidadoExcelPlanAnual.NewSheet(sheetName)
	if err := consolidadoExcelPlanAnual.SetSheetView(sheetName, -1, &excelize.ViewOptions{
		ShowGridLines: &disable,
	}); err != nil {
		fmt.Println(err)
	}
	_ = consolidadoExcelPlanAnual.DeleteSheet("Sheet1")

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Id:`+body["vigencia"].(string), &resPeriodo); err == nil {
		helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
	}

	anioPlan := periodo[0]["Year"].(float64)

	if anioPlan > 2024 { //? Vigencias 2025 en adelante son reportes con nueva estructura
		esReporteAntiguo = false
	}

	for planes := 0; planes < len(planesFilter); planes++ {
		Limp()
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-plan/"+planesFilter[planes]["estado_plan_id"].(string), &respuestaEstado); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaEstado, &estado)
		} else {
			panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
		}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-plan/"+planesFilter[planes]["tipo_plan_id"].(string), &respuestaTipoPlan); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaTipoPlan, &tipoPlan)
		} else {
			panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
		}
		if esReporteAntiguo {
			planesFilterData := planesFilter[planes]
			plan_id = planesFilterData["_id"].(string)
			infoReporte := make(map[string]interface{})
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id+"&fields=nombre,_id,hijos,activo", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
				for i := 0; i < len(subgrupos); i++ {
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
						actividades := GetActividades(subgrupos[i]["_id"].(string))
						var arregloLineamieto []map[string]interface{}
						var arregloLineamietoPI []map[string]interface{}
						sort.SliceStable(actividades, func(i int, j int) bool {
							if _, ok := actividades[i]["index"].(float64); ok {
								actividades[i]["index"] = fmt.Sprintf("%v", int(actividades[i]["index"].(float64)))
							}
							if _, ok := actividades[j]["index"].(float64); ok {
								actividades[j]["index"] = fmt.Sprintf("%v", int(actividades[j]["index"].(float64)))
							}
							aux, _ := strconv.Atoi((actividades[i]["index"]).(string))
							aux1, _ := strconv.Atoi((actividades[j]["index"]).(string))
							return aux < aux1
						})
						LimpiarDetalles()
						for j := 0; j < len(actividades); j++ {
							arregloLineamieto = nil
							arregloLineamietoPI = nil
							actividad := actividades[j]
							actividadName = actividad["dato"].(string)
							index := actividad["index"].(string)
							datosArmonizacion := make(map[string]interface{})
							titulosArmonizacion := make(map[string]interface{})

							tree := BuildTreeFa(subgrupos, index)
							treeDatos := tree[0]
							treeDatas := tree[1]
							treeArmo := tree[2]
							armonizacionTercer := treeArmo[0]
							var armonizacionTercerNivel interface{}
							var armonizacionTercerNivelPI interface{}
							if armonizacionTercer["armo"] != nil {
								armonizacionTercerNivel = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPED"]
								armonizacionTercerNivelPI = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPI"]
							}

							for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
								treeDato := treeDatos[datoGeneral]
								treeData := treeDatas[0]
								if treeDato["sub"] == "" {
									nombreMinuscula := strings.ToLower(treeDato["nombre"].(string))
									if strings.Contains(nombreMinuscula, "ponderación") || strings.Contains(nombreMinuscula, "ponderacion") && strings.Contains(nombreMinuscula, "actividad") {
										datosArmonizacion["Ponderación de la actividad"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "período") || strings.Contains(nombreMinuscula, "periodo") && strings.Contains(nombreMinuscula, "ejecucion") || strings.Contains(nombreMinuscula, "ejecución") {
										datosArmonizacion["Periodo de ejecución"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "actividad") && strings.Contains(nombreMinuscula, "general") {
										datosArmonizacion["Actividad general"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "tarea") || strings.Contains(nombreMinuscula, "actividades específicas") {
										datosArmonizacion["Tareas"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "producto") {
										datosArmonizacion["Producto esperado"] = treeData[treeDato["id"].(string)]
									} else {
										datosArmonizacion[treeDato["nombre"].(string)] = treeData[treeDato["id"].(string)]
									}
								}
							}
							var treeIndicador map[string]interface{}
							auxTree := tree[0]
							for i := 0; i < len(auxTree); i++ {
								subgrupo := auxTree[i]
								if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") {
									treeIndicador = auxTree[i]
								}
							}

							subIndicador := treeIndicador["sub"].([]map[string]interface{})
							for ind := 0; ind < len(subIndicador); ind++ {
								subIndicadorRes := subIndicador[ind]
								treeData := treeDatas[0]
								dataIndicador := make(map[string]interface{})
								auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
								for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
									dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[auxSubIndicador[subInd]["id"].(string)]
								}
								titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
							}

							datosArmonizacion["indicadores"] = titulosArmonizacion
							arregloLineamieto = ArbolArmonizacionV2(armonizacionTercerNivel.(string))
							arregloLineamietoPI = ArbolArmonizacionPIV2(armonizacionTercerNivelPI.(string))

							generalData := make(map[string]interface{})
							nombreUnidad = planesFilterData["nombreUnidad"].(string)

							generalData["nombreUnidad"] = nombreUnidad
							generalData["nombreActividad"] = actividadName
							generalData["numeroActividad"] = index
							generalData["datosArmonizacion"] = arregloLineamieto
							generalData["datosArmonizacionPI"] = arregloLineamietoPI
							generalData["datosComplementarios"] = datosArmonizacion
							arregloPlanAnual = append(arregloPlanAnual, generalData)
						}
						break
					}
				}
			} else {
				panic(err)
			}

			infoReporte["tipo_plan"] = tipoPlan["nombre"]
			infoReporte["vigencia"] = body["vigencia"]
			infoReporte["estado_plan"] = estado["nombre"]
			infoReporte["nombre_unidad"] = nombreUnidad

			arregloInfoReportes = append(arregloInfoReportes, infoReporte)

			rowPos := contadorGeneral + 5

			unidadNombre := arregloPlanAnual[0]["nombreUnidad"]

			consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "Q"+fmt.Sprint(contadorGeneral+1))
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
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "Q", "Q", 25)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "Q"+fmt.Sprint(contadorGeneral+1), stylehead)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "Q"+fmt.Sprint(contadorGeneral+2), stylehead)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "Q"+fmt.Sprint(contadorGeneral+3), styletitles)
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 30)

			var tituloExcel string
			if periodo[0] != nil {
				tituloExcel = "Plan de acción " + periodo[0]["Nombre"].(string) + " - " + unidadNombre.(string)
			} else {
				tituloExcel = "Plan de acción - " + unidadNombre.(string)
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
			_ = consolidadoExcelPlanAnual.InsertRows(sheetName, 1, 1)

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
			_ = consolidadoExcelPlanAnual.RemoveRow(sheetName, 1)
		} else {
			planesFilterData := planesFilter[planes]
			plan_id = planesFilterData["_id"].(string)
			infoReporte := make(map[string]interface{})
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id+"&fields=nombre,_id,hijos,activo", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
				for i := 0; i < len(subgrupos); i++ {
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") {
						actividades := GetActividades(subgrupos[i]["_id"].(string))
						var arregloLineamieto []map[string]interface{}
						var arregloLineamietoPI []map[string]interface{}
						sort.SliceStable(actividades, func(i int, j int) bool {
							if _, ok := actividades[i]["index"].(float64); ok {
								actividades[i]["index"] = fmt.Sprintf("%v", int(actividades[i]["index"].(float64)))
							}
							if _, ok := actividades[j]["index"].(float64); ok {
								actividades[j]["index"] = fmt.Sprintf("%v", int(actividades[j]["index"].(float64)))
							}
							aux, _ := strconv.Atoi((actividades[i]["index"]).(string))
							aux1, _ := strconv.Atoi((actividades[j]["index"]).(string))
							return aux < aux1
						})
						LimpiarDetalles()
						for j := 0; j < len(actividades); j++ {
							arregloLineamieto = nil
							arregloLineamietoPI = nil
							actividad := actividades[j]
							actividadName = actividad["dato"].(string)
							index := actividad["index"].(string)
							datosArmonizacion := make(map[string]interface{})
							titulosArmonizacion := make(map[string]interface{})

							tree := BuildTreeFa(subgrupos, index)
							treeDatos := tree[0]
							treeDatas := tree[1]
							treeArmo := tree[2]
							armonizacionTercer := treeArmo[0]
							var armonizacionTercerNivel interface{}
							var armonizacionTercerNivelPI interface{}
							if armonizacionTercer["armo"] != nil {
								armonizacionTercerNivel = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPED"]
								armonizacionTercerNivelPI = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPI"]
							}

							for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
								treeDato := treeDatos[datoGeneral]
								treeData := treeDatas[0]
								if treeDato["sub"] == "" {
									nombreMinuscula := strings.ToLower(treeDato["nombre"].(string))
									if strings.Contains(nombreMinuscula, "ponderación") || strings.Contains(nombreMinuscula, "ponderacion") && strings.Contains(nombreMinuscula, "actividad") {
										datosArmonizacion["Ponderación de la actividad"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "período") || strings.Contains(nombreMinuscula, "periodo") && strings.Contains(nombreMinuscula, "ejecucion") || strings.Contains(nombreMinuscula, "ejecución") {
										datosArmonizacion["Periodo de ejecución"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "actividad") {
										datosArmonizacion["Actividad general"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "producto") {
										datosArmonizacion["Producto esperado"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "unidad") || strings.Contains(nombreMinuscula, "grupo") {
										datosArmonizacion["Responsable"] = treeData[treeDato["id"].(string)]
									} else {
										datosArmonizacion[treeDato["nombre"].(string)] = treeData[treeDato["id"].(string)]
									}
								}
							}
							var treeIndicador map[string]interface{}
							auxTree := tree[0]
							for i := 0; i < len(auxTree); i++ {
								subgrupo := auxTree[i]
								if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") {
									treeIndicador = auxTree[i]
								}
							}

							subIndicador := treeIndicador["sub"].([]map[string]interface{})
							for ind := 0; ind < len(subIndicador); ind++ {
								subIndicadorRes := subIndicador[ind]
								treeData := treeDatas[0]
								dataIndicador := make(map[string]interface{})
								auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
								for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
									dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[auxSubIndicador[subInd]["id"].(string)]
								}
								titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
							}

							datosArmonizacion["indicadores"] = titulosArmonizacion
							arregloLineamieto = ArbolArmonizacionV2(armonizacionTercerNivel.(string))
							arregloLineamietoPI = ArbolArmonizacionPIV2(armonizacionTercerNivelPI.(string))

							generalData := make(map[string]interface{})
							nombreUnidad = planesFilterData["nombreUnidad"].(string)

							generalData["nombreUnidad"] = nombreUnidad
							generalData["nombreActividad"] = actividadName
							generalData["numeroActividad"] = index
							generalData["datosArmonizacion"] = arregloLineamieto
							generalData["datosArmonizacionPI"] = arregloLineamietoPI
							generalData["datosComplementarios"] = datosArmonizacion
							arregloPlanAnual = append(arregloPlanAnual, generalData)
						}
						break
					}
				}
			} else {
				panic(err)
			}

			infoReporte["tipo_plan"] = tipoPlan["nombre"]
			infoReporte["vigencia"] = body["vigencia"]
			infoReporte["estado_plan"] = estado["nombre"]
			infoReporte["nombre_unidad"] = nombreUnidad

			arregloInfoReportes = append(arregloInfoReportes, infoReporte)

			rowPos := contadorGeneral + 5

			unidadNombre := arregloPlanAnual[0]["nombreUnidad"]

			consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "Q"+fmt.Sprint(contadorGeneral+1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "D"+fmt.Sprint(contadorGeneral+2))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorGeneral+2), "G"+fmt.Sprint(contadorGeneral+2))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorGeneral+2), "H"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorGeneral+2), "I"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorGeneral+2), "J"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorGeneral+2), "K"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(contadorGeneral+2), "P"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "Q"+fmt.Sprint(contadorGeneral+2), "Q"+fmt.Sprint(contadorGeneral+3))
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
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "M", 52)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "N", "N", 30)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "O", "O", 10)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "P", "P", 25)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "Q", "Q", 30)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "Q"+fmt.Sprint(contadorGeneral+1), stylehead)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "Q"+fmt.Sprint(contadorGeneral+2), stylehead)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "Q"+fmt.Sprint(contadorGeneral+3), styletitles)
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 30)

			var tituloExcel string
			if periodo[0] != nil {
				tituloExcel = "Plan de acción " + periodo[0]["Nombre"].(string) + " - " + unidadNombre.(string)
			} else {
				tituloExcel = "Plan de acción - " + unidadNombre.(string)
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
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorGeneral+3), "Peso (%)")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorGeneral+3), "Periodo de ejecución")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorGeneral+3), "Actividad")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+2), "Indicador")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+3), "Nombre")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorGeneral+3), "Fórmula")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorGeneral+3), "Criterio del indicador")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(contadorGeneral+3), "Meta")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(contadorGeneral+3), "Producto esperado")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "Q"+fmt.Sprint(contadorGeneral+3), "Unidad o grupo responsable")
			_ = consolidadoExcelPlanAnual.InsertRows(sheetName, 1, 1)

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
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(rowPos), excelPlan+1)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(rowPos), datosComplementarios["Ponderación de la actividad"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(rowPos), datosComplementarios["Periodo de ejecución"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(rowPos), datosExcelPlan["nombreActividad"])
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

						consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(y_ind), nombreIndicador)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(y_ind), formula)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(y_ind), criterio)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(y_ind), meta)
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
				consolidadoExcelPlanAnual.MergeCell(sheetName, "Q"+fmt.Sprint(rowPos), "Q"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "Q"+fmt.Sprint(rowPos), datosComplementarios["Responsable"])
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(rowPos), "P"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)
				SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "Q"+fmt.Sprint(rowPos), "Q"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)

				rowPos += MaxRowsXActivity

				contadorGeneral = rowPos - 2

				consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
			}
			arregloPlanAnual = nil
			_ = consolidadoExcelPlanAnual.RemoveRow(sheetName, 1)
		}
	}
	_ = consolidadoExcelPlanAnual.InsertRows("REPORTE GENERAL", 1, 3)
	_ = consolidadoExcelPlanAnual.MergeCell("REPORTE GENERAL", "C2", "P6")
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

	return consolidadoExcelPlanAnual, arregloInfoReportes, nil
}

func ConstruirExcelPlanAccionEvaluacion(esReporteAntiguo bool, datosReporte map[string]interface{}) (*excelize.File, error) {
	unidadNombre := datosReporte["unidadNombre"].(string)
	periodo := datosReporte["periodo"].([]map[string]interface{})
	evaluacion := datosReporte["evaluacion"].([]map[string]interface{})
	excelArmonizacion := datosReporte["excelArmonizacion"].([]map[string]interface{})

	consolidadoExcelEvaluacion := excelize.NewFile()
	sheetName := "Evaluación"
	consolidadoExcelEvaluacion.NewSheet(sheetName)
	_ = consolidadoExcelEvaluacion.DeleteSheet("Sheet1")

	disable := false
	if err := consolidadoExcelEvaluacion.SetSheetView(sheetName, -1, &excelize.ViewOptions{
		ShowGridLines: &disable,
	}); err != nil {
		fmt.Println(err)
	}

	styleUnidad, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Font:      &excelize.Font{Bold: true, Color: "000000", Family: "Bahnschrift SemiBold SemiConden", Size: 20},
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})
	styleTituloSB, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Family: "Calibri", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
	})
	styleSombreadoSB, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Color: "000000", Family: "Calibri", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
	})
	styleNegrilla, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Color: "000000", Family: "Calibri", Size: 12, Bold: true},
	})
	styleTitulo, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Family: "Calibri", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleContenido, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleContenidoC, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleContenidoCI, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Color: "FFFFFF"},
	})
	styleContenidoCIP, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		NumFmt:    10,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Color: "FFFFFF"},
	})
	styleContenidoCE, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		NumFmt:    1,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleContenidoCD, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		NumFmt:    4,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleSubTitles, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
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
	styleContenidoCS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FCE4D6"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleContenidoCP, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		NumFmt:    10,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleContenidoCPSR, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		NumFmt:    10,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Family: "Calibri", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleContenidoCPS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		NumFmt:    10,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FCE4D6"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontent, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentC, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styleLineamiento, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
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
	styleLineamientoSombra, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
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
	stylecontentCS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// Size
	consolidadoExcelEvaluacion.SetRowHeight(sheetName, 1, 12)
	consolidadoExcelEvaluacion.SetRowHeight(sheetName, 2, 27)
	consolidadoExcelEvaluacion.SetRowHeight(sheetName, 19, 31)
	consolidadoExcelEvaluacion.SetRowHeight(sheetName, 22, 27)

	consolidadoExcelEvaluacion.SetColWidth(sheetName, "A", "A", 3)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "B", "B", 19)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "C", "C", 13)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "D", "G", 35)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "E", "E", 16)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "H", "H", 4)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "I", "I", 8)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "J", "J", 13)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "K", "K", 42)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "L", "L", 16)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "M", "M", 21)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "N", "BG", 14)
	consolidadoExcelEvaluacion.SetColWidth(sheetName, "BJ", "BN", 3)
	// Merge
	consolidadoExcelEvaluacion.MergeCell(sheetName, "B4", "D4")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "B19", "E19")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "B21", "D21")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "E21", "G21")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "H21", "H22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "I21", "I22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "J21", "J22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "K21", "K22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "L21", "L22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "M21", "M22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "N21", "N22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "O21", "O22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "BL19", "BM19")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "BL21", "BL22")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "BM21", "BM22")

	consolidadoExcelEvaluacion.MergeCell(sheetName, "P21", "Z21")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "AA21", "AK21")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "AL21", "AV21")
	consolidadoExcelEvaluacion.MergeCell(sheetName, "AW21", "BG21")
	// Style
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B2", "T2", styleUnidad)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B4", "D4", styleTituloSB)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "E4", "E4", styleSombreadoSB)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B19", "B19", styleNegrilla)
	//  Estilos títulos
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B21", "BG21", styleTitulo)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "J22", "BG22", styleContenidoC)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B22", "G22", styleSubTitles)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "P22", "BG22", styleSubTitles)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z22", "Z22", styleContenidoCS)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AK22", "AK22", styleContenidoCS)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AV22", "AV22", styleContenidoCS)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BG22", "BG22", styleContenidoCS)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H22", "O22", styleTitulo)

	if periodo[0] != nil {
		consolidadoExcelEvaluacion.SetCellValue(sheetName, "B2", "Evaluación Plan de Acción "+periodo[0]["Nombre"].(string)+" - "+unidadNombre)
		consolidadoExcelEvaluacion.SetCellValue(sheetName, "B19", "Cumplimiento General Plan de Acción "+periodo[0]["Nombre"].(string)+" - "+unidadNombre)
	} else {
		consolidadoExcelEvaluacion.SetCellValue(sheetName, "B2", "Evaluación Plan de Acción - "+unidadNombre)
		consolidadoExcelEvaluacion.SetCellValue(sheetName, "B19", "Cumplimiento General Plan de Acción - "+unidadNombre)
	}

	ddRango1 := excelize.NewDataValidation(true)
	ddRango1.Sqref = "E4:E4"
	ddRango1.SetDropList([]string{"Trimestre I", "Trimestre II", "Trimestre III", "Trimestre IV"})

	if err := consolidadoExcelEvaluacion.AddDataValidation(sheetName, ddRango1); err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Titles
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "B4", "Seleccione el periodo:")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "B21", "Armonización PED")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "E21", "Armonización Plan Indicativo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "H21", "No.")
	if esReporteAntiguo {
		consolidadoExcelEvaluacion.SetCellValue(sheetName, "I21", "Pond.")
	} else {
		consolidadoExcelEvaluacion.SetCellValue(sheetName, "I21", "Peso (%)")
	}
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "J21", "Periodo de ejecución")
	if esReporteAntiguo {
		consolidadoExcelEvaluacion.SetCellValue(sheetName, "K21", "Actividad General")
	} else {
		consolidadoExcelEvaluacion.SetCellValue(sheetName, "K21", "Actividad")
	}
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "L21", "Indicador asociado")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "M21", "Fórmula del Indicador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "N21", "Meta")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "O21", "Tipo de Unidad")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "P21", "Trimestre I")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA21", "Trimestre II")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL21", "Trimestre III")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AW21", "Trimestre IV")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "B22", "Lineamiento")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "C22", "Meta")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "D22", "Estrategias")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "E22", "Ejes transformadores")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "F22", "Lineamientos de acción")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "G22", "Estrategias")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "P22", "Reporte de avance")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q22", "Dificultades")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "R22", "Numerador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "S22", "Denominador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "T22", "Indicador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "U22", "Numerador Acumulador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "V22", "Denominador Acumulador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "W22", "Indicador Acumulado")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "X22", "Cumplimiento por meta")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y22", "Brecha")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z22", "Cumplimiento por actividad")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA22", "Reporte de avance")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB22", "Dificultades")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC22", "Numerador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD22", "Denominador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE22", "Indicador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF22", "Numerador Acumulador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG22", "Denominador Acumulador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH22", "Indicador Acumulado")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI22", "Cumplimiento por meta")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ22", "Brecha")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK22", "Cumplimiento por actividad")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL22", "Reporte de avance")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM22", "Dificultades")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN22", "Numerador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO22", "Denominador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP22", "Indicador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ22", "Numerador Acumulador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR22", "Denominador Acumulador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS22", "Indicador Acumulado")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AT22", "Cumplimiento por meta")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AU22", "Brecha")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AV22", "Cumplimiento por actividad")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AW22", "Reporte de avance")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX22", "Dificultades")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY22", "Numerador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AZ22", "Denominador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BA22", "Indicador del Periodo")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BB22", "Numerador Acumulador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BC22", "Denominador Acumulador")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BD22", "Indicador Acumulado")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BE22", "Cumplimiento por meta")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BF22", "Brecha")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BG22", "Cumplimiento por actividad")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BL19", "Gráfica")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BL21", "No.")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BM21", "Cumplimiento")

	indice := 23
	indiceGraficos := 23
	// Agregar armonización
	for posArmonizacion := 0; posArmonizacion < len(excelArmonizacion); posArmonizacion++ {
		datosExcelArmonizacion := excelArmonizacion[posArmonizacion]
		armoPED := datosExcelArmonizacion["datosArmonizacionPED"].([]map[string]interface{})
		armoPI := datosExcelArmonizacion["datosArmonizacionPI"].([]map[string]interface{})
		datosComplementarios := datosExcelArmonizacion["datosComplementarios"].(map[string]interface{})
		indicadores := datosComplementarios["indicadores"].(map[string]interface{})
		numeroActividad := datosExcelArmonizacion["numeroActividad"]

		MaxRowsXActivity := MinComMul_Armonization(armoPED, armoPI, len(indicadores))

		y_lin := indice
		h_lin := MaxRowsXActivity / len(armoPED)
		consolidadoExcelEvaluacion.SetRowHeight(sheetName, y_lin, 27*5)
		for _, lin := range armoPED {
			consolidadoExcelEvaluacion.MergeCell(sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1))
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "B"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
			SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1), styleLineamiento, styleLineamientoSombra)
			y_met := y_lin
			h_met := h_lin / len(lin["meta"].([]map[string]interface{}))
			for _, met := range lin["meta"].([]map[string]interface{}) {
				consolidadoExcelEvaluacion.MergeCell(sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "C"+fmt.Sprint(y_met), met["nombreMeta"])
				SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1), stylecontentC, stylecontentCS)
				y_est := y_met
				h_est := h_met / len(met["estrategias"].([]map[string]interface{}))
				for _, est := range met["estrategias"].([]map[string]interface{}) {
					consolidadoExcelEvaluacion.MergeCell(sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1))
					consolidadoExcelEvaluacion.SetCellValue(sheetName, "D"+fmt.Sprint(y_est), est["descripcionEstrategia"])
					if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
						SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
					} else {
						SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
					}
					y_est += h_est
				}
				y_met += h_met
			}
			y_lin += h_lin
		}

		y_eje := indice
		h_eje := MaxRowsXActivity / len(armoPI)
		for _, eje := range armoPI {
			consolidadoExcelEvaluacion.MergeCell(sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1))
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "E"+fmt.Sprint(y_eje), eje["nombreFactor"])
			SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1), stylecontentC, stylecontentCS)
			y_lin := y_eje
			h_lin := h_eje / len(eje["lineamientos"].([]map[string]interface{}))
			for _, lin := range eje["lineamientos"].([]map[string]interface{}) {
				consolidadoExcelEvaluacion.MergeCell(sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "F"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
				SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1), stylecontentC, stylecontentCS)
				y_est := y_lin
				h_est := h_lin / len(lin["estrategias"].([]map[string]interface{}))
				for _, est := range lin["estrategias"].([]map[string]interface{}) {
					consolidadoExcelEvaluacion.MergeCell(sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1))
					consolidadoExcelEvaluacion.SetCellValue(sheetName, "G"+fmt.Sprint(y_est), est["descripcionEstrategia"])
					if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
						SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
					} else {
						SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
					}
					y_est += h_est
				}
				y_lin += h_lin
			}
			y_eje += h_eje
		}

		for i, actividad := range evaluacion {
			if numeroActividad == actividad["numero"] {
				// Union de celdas
				consolidadoExcelEvaluacion.MergeCell(sheetName, "H"+fmt.Sprint(indice), "H"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "I"+fmt.Sprint(indice), "I"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "J"+fmt.Sprint(indice), "J"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "K"+fmt.Sprint(indice), "K"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "L"+fmt.Sprint(indice), "L"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "M"+fmt.Sprint(indice), "M"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "N"+fmt.Sprint(indice), "N"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "O"+fmt.Sprint(indice), "O"+fmt.Sprint(indice+MaxRowsXActivity-1))

				// Datos
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "H"+fmt.Sprint(indice), posArmonizacion+1)
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "I"+fmt.Sprint(indice), actividad["ponderado"].(float64)/100)
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "J"+fmt.Sprint(indice), actividad["periodo"])
				if esReporteAntiguo {
					consolidadoExcelEvaluacion.SetCellValue(sheetName, "K"+fmt.Sprint(indice), actividad["actividad"])
				} else {
					consolidadoExcelEvaluacion.SetCellValue(sheetName, "K"+fmt.Sprint(indice), datosExcelArmonizacion["nombreActividad"].(string))
				}
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "L"+fmt.Sprint(indice), actividad["indicador"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "M"+fmt.Sprint(indice), actividad["formula"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "N"+fmt.Sprint(indice), actividad["meta"].(float64))
				if actividad["unidad"] == "Porcentaje" {
					consolidadoExcelEvaluacion.SetCellValue(sheetName, "N"+fmt.Sprint(indice), actividad["meta"].(float64)/100)
				}
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "O"+fmt.Sprint(indice), actividad["unidad"])

				consolidadoExcelEvaluacion.MergeCell(sheetName, "P"+fmt.Sprint(indice), "P"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "Q"+fmt.Sprint(indice), "Q"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "R"+fmt.Sprint(indice), "R"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "S"+fmt.Sprint(indice), "S"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "T"+fmt.Sprint(indice), "T"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "U"+fmt.Sprint(indice), "U"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "V"+fmt.Sprint(indice), "V"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "W"+fmt.Sprint(indice), "W"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "X"+fmt.Sprint(indice), "X"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "Y"+fmt.Sprint(indice), "Y"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice+MaxRowsXActivity-1))

				consolidadoExcelEvaluacion.MergeCell(sheetName, "AA"+fmt.Sprint(indice), "AA"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AB"+fmt.Sprint(indice), "AB"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AC"+fmt.Sprint(indice), "AC"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AD"+fmt.Sprint(indice), "AD"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AE"+fmt.Sprint(indice), "AE"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AF"+fmt.Sprint(indice), "AF"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AG"+fmt.Sprint(indice), "AG"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AH"+fmt.Sprint(indice), "AH"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AI"+fmt.Sprint(indice), "AI"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AK"+fmt.Sprint(indice), "AK"+fmt.Sprint(indice+MaxRowsXActivity-1))

				consolidadoExcelEvaluacion.MergeCell(sheetName, "AL"+fmt.Sprint(indice), "AL"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AM"+fmt.Sprint(indice), "AM"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AN"+fmt.Sprint(indice), "AN"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AO"+fmt.Sprint(indice), "AO"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AQ"+fmt.Sprint(indice), "AQ"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AR"+fmt.Sprint(indice), "AR"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AT"+fmt.Sprint(indice), "AT"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AU"+fmt.Sprint(indice), "AU"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AV"+fmt.Sprint(indice), "AV"+fmt.Sprint(indice+MaxRowsXActivity-1))

				consolidadoExcelEvaluacion.MergeCell(sheetName, "AW"+fmt.Sprint(indice), "AW"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AX"+fmt.Sprint(indice), "AX"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AY"+fmt.Sprint(indice), "AY"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "AZ"+fmt.Sprint(indice), "AZ"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "BA"+fmt.Sprint(indice), "BA"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "BB"+fmt.Sprint(indice), "BB"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "BC"+fmt.Sprint(indice), "BC"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "BD"+fmt.Sprint(indice), "BD"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "BE"+fmt.Sprint(indice), "BE"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "BF"+fmt.Sprint(indice), "BF"+fmt.Sprint(indice+MaxRowsXActivity-1))
				consolidadoExcelEvaluacion.MergeCell(sheetName, "BG"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice+MaxRowsXActivity-1))

				// Trimestres
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "P"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["cualitativo"].(map[string]interface{})["reporte"].(string))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["cualitativo"].(map[string]interface{})["dificultades"].(string))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "R"+fmt.Sprint(indice), Convert2Num(actividad["trimestre1"].(map[string]interface{})["numerador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "S"+fmt.Sprint(indice), Convert2Num(actividad["trimestre1"].(map[string]interface{})["denominador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "T"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "U"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["numeradorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "V"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["denominadorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "W"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["acumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "X"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["meta"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["brecha"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["actividad"].(float64))

				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["cualitativo"].(map[string]interface{})["reporte"].(string))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["cualitativo"].(map[string]interface{})["dificultades"].(string))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC"+fmt.Sprint(indice), Convert2Num(actividad["trimestre2"].(map[string]interface{})["numerador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD"+fmt.Sprint(indice), Convert2Num(actividad["trimestre2"].(map[string]interface{})["denominador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["numeradorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["denominadorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["acumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["meta"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["brecha"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["actividad"].(float64))

				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["cualitativo"].(map[string]interface{})["reporte"].(string))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["cualitativo"].(map[string]interface{})["dificultades"].(string))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN"+fmt.Sprint(indice), Convert2Num(actividad["trimestre3"].(map[string]interface{})["numerador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO"+fmt.Sprint(indice), Convert2Num(actividad["trimestre3"].(map[string]interface{})["denominador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["numeradorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["denominadorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["acumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AT"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["meta"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AU"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["brecha"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AV"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["actividad"].(float64))

				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AW"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["cualitativo"].(map[string]interface{})["reporte"].(string))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["cualitativo"].(map[string]interface{})["dificultades"].(string))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY"+fmt.Sprint(indice), Convert2Num(actividad["trimestre4"].(map[string]interface{})["numerador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AZ"+fmt.Sprint(indice), Convert2Num(actividad["trimestre4"].(map[string]interface{})["denominador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "BA"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "BB"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["numeradorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "BC"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["denominadorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "BD"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["acumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "BE"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["meta"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "BF"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["brecha"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "BG"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["actividad"].(float64))

				// Estilos
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H"+fmt.Sprint(indice), "AY"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoC)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "I"+fmt.Sprint(indice), "I"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "K"+fmt.Sprint(indice), "K"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenido)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "P"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "R"+fmt.Sprint(indice), "S"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "U"+fmt.Sprint(indice), "V"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AC"+fmt.Sprint(indice), "AD"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AF"+fmt.Sprint(indice), "AG"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AN"+fmt.Sprint(indice), "AO"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AQ"+fmt.Sprint(indice), "AR"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AY"+fmt.Sprint(indice), "AZ"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BB"+fmt.Sprint(indice), "BC"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCPS)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AK"+fmt.Sprint(indice), "AK"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCPS)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AV"+fmt.Sprint(indice), "AV"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCPS)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BG"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCPS)

				if actividad["unidad"] == "Porcentaje" {
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "N"+fmt.Sprint(indice), "N"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "T"+fmt.Sprint(indice), "T"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "W"+fmt.Sprint(indice), "W"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Y"+fmt.Sprint(indice), "Y"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AE"+fmt.Sprint(indice), "AE"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AH"+fmt.Sprint(indice), "AH"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AU"+fmt.Sprint(indice), "AU"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BA"+fmt.Sprint(indice), "BA"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BD"+fmt.Sprint(indice), "BD"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BF"+fmt.Sprint(indice), "BF"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
				} else if actividad["unidad"] == "Tasa" {
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "N"+fmt.Sprint(indice), "N"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "T"+fmt.Sprint(indice), "T"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "W"+fmt.Sprint(indice), "W"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Y"+fmt.Sprint(indice), "Y"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AE"+fmt.Sprint(indice), "AE"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AH"+fmt.Sprint(indice), "AH"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AU"+fmt.Sprint(indice), "AU"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BA"+fmt.Sprint(indice), "BA"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BD"+fmt.Sprint(indice), "BD"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BF"+fmt.Sprint(indice), "BF"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
				} else if actividad["unidad"] == "Unidad" {
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "N"+fmt.Sprint(indice), "N"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "T"+fmt.Sprint(indice), "T"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "W"+fmt.Sprint(indice), "W"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Y"+fmt.Sprint(indice), "Y"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AE"+fmt.Sprint(indice), "AE"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AH"+fmt.Sprint(indice), "AH"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AU"+fmt.Sprint(indice), "AU"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BA"+fmt.Sprint(indice), "BA"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BD"+fmt.Sprint(indice), "BD"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BF"+fmt.Sprint(indice), "BF"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
				}

				// Unión de celdas por indicador
				if i > 0 {
					if actividad["numero"] == evaluacion[i-1]["numero"] {
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "H"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "I"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "J"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "K"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AV"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BG"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BL"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BM"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.MergeCell(sheetName, "P"+fmt.Sprint(indice-1), "P"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "Q"+fmt.Sprint(indice-1), "Q"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AA"+fmt.Sprint(indice-1), "AA"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AB"+fmt.Sprint(indice-1), "AB"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AL"+fmt.Sprint(indice-1), "AL"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AM"+fmt.Sprint(indice-1), "AM"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AW"+fmt.Sprint(indice-1), "AW"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AX"+fmt.Sprint(indice-1), "AX"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "B"+fmt.Sprint(indice-1), "B"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "C"+fmt.Sprint(indice-1), "C"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "D"+fmt.Sprint(indice-1), "D"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "E"+fmt.Sprint(indice-1), "E"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "F"+fmt.Sprint(indice-1), "F"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "G"+fmt.Sprint(indice-1), "G"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B"+fmt.Sprint(indice-1), "B"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "C"+fmt.Sprint(indice-1), "C"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "D"+fmt.Sprint(indice-1), "D"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "E"+fmt.Sprint(indice-1), "E"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "F"+fmt.Sprint(indice-1), "F"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "G"+fmt.Sprint(indice-1), "G"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)

						consolidadoExcelEvaluacion.MergeCell(sheetName, "H"+fmt.Sprint(indice-1), "H"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "I"+fmt.Sprint(indice-1), "I"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "J"+fmt.Sprint(indice-1), "J"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "K"+fmt.Sprint(indice-1), "K"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "Z"+fmt.Sprint(indice-1), "Z"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AK"+fmt.Sprint(indice-1), "AK"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AV"+fmt.Sprint(indice-1), "AV"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "BG"+fmt.Sprint(indice-1), "BG"+fmt.Sprint(indice))
					} else {
						// Gaficos
						consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BL"+fmt.Sprint(indiceGraficos), "=H"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BM"+fmt.Sprint(indiceGraficos), "=IF(E4=\"Trimestre I\",Z"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AK"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AV"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",BG"+fmt.Sprint(indice)+"))))")
						indiceGraficos++
					}
				} else if i == 0 {
					// Gaficos
					consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BL"+fmt.Sprint(indiceGraficos), "=H"+fmt.Sprint(indice))
					consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BM"+fmt.Sprint(indiceGraficos), "=IF(E4=\"Trimestre I\",Z"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AK"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AV"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",BG"+fmt.Sprint(indice)+"))))")
					indiceGraficos++
				}
				indice += MaxRowsXActivity
			}
		}
	}

	consolidadoExcelEvaluacion.MergeCell(sheetName, "H"+fmt.Sprint(indice), "O"+fmt.Sprint(indice))
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H"+fmt.Sprint(indice), "O"+fmt.Sprint(indice), styleTituloSB)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "P"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice), styleContenidoC)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice), styleContenidoCPSR)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AK"+fmt.Sprint(indice), "AK"+fmt.Sprint(indice), styleContenidoCPSR)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AV"+fmt.Sprint(indice), "AV"+fmt.Sprint(indice), styleContenidoCPSR)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BG"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice), styleContenidoCPSR)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BL19", "BL"+fmt.Sprint(indice+1), styleContenidoCI)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BM19", "BN"+fmt.Sprint(indice+1), styleContenidoCIP)
	consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BM21", "BM22", styleContenidoCI)
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "H"+fmt.Sprint(indice), "Avance General del Plan de Acción")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "E4", "Trimestre I")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "P"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "R"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "S"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "T"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "U"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "V"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "W"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "X"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y"+fmt.Sprint(indice), "-")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ"+fmt.Sprint(indice), "-")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AT"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AU"+fmt.Sprint(indice), "-")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AW"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "AZ"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BA"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BB"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BC"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BD"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BE"+fmt.Sprint(indice), "-")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BF"+fmt.Sprint(indice), "-")

	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BL"+fmt.Sprint(indice), "General")

	filaAnt := fmt.Sprint(indice - 1)
	consolidadoExcelEvaluacion.SetCellFormula(sheetName, "Z"+fmt.Sprint(indice), "=SUMPRODUCT(I23:I"+filaAnt+",Z23:Z"+filaAnt+")")
	consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AK"+fmt.Sprint(indice), "=SUMPRODUCT(I23:I"+filaAnt+",AK23:AK"+filaAnt+")")
	consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AV"+fmt.Sprint(indice), "=SUMPRODUCT(I23:I"+filaAnt+",AV23:AV"+filaAnt+")")
	consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BG"+fmt.Sprint(indice), "=SUMPRODUCT(I23:I"+filaAnt+",BG23:BG"+filaAnt+")")
	consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BM"+fmt.Sprint(indice), "=IF(E4=\"Trimestre I\",Z"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AK"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AV"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",BG"+fmt.Sprint(indice)+"))))")
	consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BN"+fmt.Sprint(indice), "=100%-BM"+fmt.Sprint(indice))
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BM"+fmt.Sprint(indice+1), "Avance")
	consolidadoExcelEvaluacion.SetCellValue(sheetName, "BN"+fmt.Sprint(indice+1), "Restante")

	consolidadoExcelEvaluacion.AddChart(sheetName, "B5", &excelize.Chart{
		Type: excelize.Pie,
		Series: []excelize.ChartSeries{
			{
				Name:       "",
				Categories: sheetName + "!$BM$" + fmt.Sprint(indice+1) + ":$BN$" + fmt.Sprint(indice+1),
				Values:     sheetName + "!$BM$" + fmt.Sprint(indice) + ":$BN$" + fmt.Sprint(indice),
			},
		},
		Format: excelize.GraphicOptions{
			ScaleX:          1.0,
			ScaleY:          1.0,
			OffsetX:         15,
			OffsetY:         10,
			LockAspectRatio: false,
			Locked:          &disable,
		},
		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowSerName:     false,
			ShowVal:         false,
		},
		ShowBlanksAs: "zero",
		Dimension: excelize.ChartDimension{
			Height: 265,
			Width:  454,
		},
		XAxis: excelize.ChartAxis{
			None: true,
		},
		YAxis: excelize.ChartAxis{
			None: true,
		},
	})

	consolidadoExcelEvaluacion.AddChart(sheetName, "F4", &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       "",
				Categories: sheetName + "!$BL$23:$BL$" + fmt.Sprint(indiceGraficos-1),
				Values:     sheetName + "!$BM$23:$BM$" + fmt.Sprint(indiceGraficos-1),
			},
		},
		Format: excelize.GraphicOptions{
			OffsetX:         15,
			LockAspectRatio: false,
			Locked:          &disable,
		},
		Dimension: excelize.ChartDimension{
			Height: 344,
			Width:  1605,
		},
		PlotArea: excelize.ChartPlotArea{
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     false,
			ShowSerName:     false,
			ShowVal:         true,
		},
		YAxis: excelize.ChartAxis{
			MajorGridLines: true,
			Font:           excelize.Font{Family: "Calibri", Size: 9, Color: "000000"},
		},
		XAxis: excelize.ChartAxis{
			Font: excelize.Font{Family: "Calibri", Size: 9, Color: "000000"},
		},
		VaryColors:   &disable,
		ShowBlanksAs: "span",
	})

	return consolidadoExcelEvaluacion, nil
}
