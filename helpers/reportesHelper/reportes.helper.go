package reporteshelper

import (
	"encoding/json"
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
			var dato map[string]interface{}
			var data_identi []map[string]interface{}
			if identificacion["dato"] != nil && identificacion["dato"] != "{}" {
				result := make(map[string]interface{})
				dato_str := identificacion["dato"].(string)
				json.Unmarshal([]byte(dato_str), &dato)

				var identi map[string]interface{}
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
					data_identi = nil
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
		strValor := strings.TrimLeft(aux["valor"].(string), "$")
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
		infoDocentes := TotalDocentes(docentes)
		rubros := docentes["rubros"].([]map[string]interface{})
		consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Identificación docente")
		consolidadoExcelPlanAnual.SetRowHeight(sheetName, contador+1, 7)
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), styletitles)

		contador++
		contador++

		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), "Código del rubro")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), "Nombre del rubro")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), "Valor")
		consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylehead)

		contador++

		//Cuerpo Tabla
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), codigoRubrosDocentes(rubros, "Prima de Servicios"))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), "Prima de servicios - pregrado")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), infoDocentes["primaServicios"])
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "D"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), codigoRubrosDocentes(rubros, "Prima de navidad"))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), "Prima de navidad - posgrado")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), infoDocentes["primaNavidad"])
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "D"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
		contador++

		consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contador), codigoRubrosDocentes(rubros, "Prima de vacaciones"))
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contador), "Prima de vacaciones")
		consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contador), infoDocentes["primaVacaciones"])
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "B"+fmt.Sprint(contador), "C"+fmt.Sprint(contador), stylecontent, stylecontentS)
		SombrearCeldas(consolidadoExcelPlanAnual, contador, sheetName, "D"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontentMR, stylecontentMRS)
	}

	consolidadoExcelPlanAnual.InsertRows(sheetName, 1, 7)
	consolidadoExcelPlanAnual.MergeCell(sheetName, "C2", "G6")

	return consolidadoExcelPlanAnual
}

func codigoRubrosDocentes(rubros []map[string]interface{}, categoria string) string {
	var codigo string
	for i := 0; i < len(rubros); i++ {
		rubro := rubros[i]
		if rubro["categoria"] == categoria {
			codigo = rubro["rubro"].(string)
			break
		}
	}
	return codigo
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
		if aux["primaServicios"] != nil {
			strPrimaServicios := strings.TrimLeft(aux["primaServicios"].(string), "$")
			strPrimaServicios = strings.ReplaceAll(strPrimaServicios, ",", "")
			arrPrimaServicios := strings.Split(strPrimaServicios, ".")
			auxPrimaServicios, err := strconv.Atoi(arrPrimaServicios[0])
			if err == nil {
				primaServicios += auxPrimaServicios
			}
		}

		if aux["primaNavidad"] != nil {
			strPrimaNavidad := strings.TrimLeft(aux["primaNavidad"].(string), "$")
			strPrimaNavidad = strings.ReplaceAll(strPrimaNavidad, ",", "")
			arrPrimaNavidad := strings.Split(strPrimaNavidad, ".")
			auxPrimaNavidad, err := strconv.Atoi(arrPrimaNavidad[0])
			if err == nil {
				primaNavidad += auxPrimaNavidad
			}
		}

		if aux["primaVacaciones"] != nil {
			strPrimaVacaciones := strings.TrimLeft(aux["primaVacaciones"].(string), "$")
			strPrimaVacaciones = strings.ReplaceAll(strPrimaVacaciones, ",", "")
			arrPrimaVacaiones := strings.Split(strPrimaVacaciones, ".")
			auxPrimaVacaciones, err := strconv.Atoi(arrPrimaVacaiones[0])
			if err == nil {
				primaVacaciones += auxPrimaVacaciones
			}
		}

		if aux["bonificacion"] != nil || aux["bonificacion"] != "N/A" {
			strBonificacion := strings.TrimLeft(aux["bonificacion"].(string), "$")
			strBonificacion = strings.ReplaceAll(strBonificacion, ",", "")
			arrBonificacion := strings.Split(strBonificacion, ".")
			auxBonificacion, err := strconv.Atoi(arrBonificacion[0])
			if err == nil {
				bonificacion += auxBonificacion
			}
		}

		if aux["interesesCesantias"] != nil || aux["interesesCesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["interesesCesantias"].(string), "$")
			strInteresesCesantias = strings.ReplaceAll(strInteresesCesantias, ",", "")
			arrInteresesCesantias := strings.Split(strInteresesCesantias, ".")
			auxInteresesCesantias, err := strconv.Atoi(arrInteresesCesantias[0])
			if err == nil {
				interesesCesantias += auxInteresesCesantias
			}
		}

		if aux["cesantiasPublico"] != nil {
			strCesantiasPublico := strings.TrimLeft(aux["cesantiasPublico"].(string), "$")
			strCesantiasPublico = strings.ReplaceAll(strCesantiasPublico, ",", "")
			arrCesantiasPublico := strings.Split(strCesantiasPublico, ".")
			auxCesantiasPublico, err := strconv.Atoi(arrCesantiasPublico[0])
			if err == nil {
				cesantiasPublicas += auxCesantiasPublico
			}
		}

		if aux["cesantiasPrivado"] != nil {
			strCesantiasPrivado := strings.TrimLeft(aux["cesantiasPrivado"].(string), "$")
			strCesantiasPrivado = strings.ReplaceAll(strCesantiasPrivado, ",", "")
			arrCesantiasPrivado := strings.Split(strCesantiasPrivado, ".")
			auxCesantiasPrivado, err := strconv.Atoi(arrCesantiasPrivado[0])
			if err == nil {
				cesantiasPrivadas += auxCesantiasPrivado
			}
		}

		if aux["totalSalud"] != nil {
			strSalud := strings.TrimLeft(aux["totalSalud"].(string), "$")
			strSalud = strings.ReplaceAll(strSalud, ",", "")
			arrSalud := strings.Split(strSalud, ".")
			auxSalud, err := strconv.Atoi(arrSalud[0])
			if err == nil {
				salud += auxSalud
			}
		}

		if aux["pensionesPublico"] != nil {
			strPensionesPublicas := strings.TrimLeft(aux["pensionesPublico"].(string), "$")
			strPensionesPublicas = strings.ReplaceAll(strPensionesPublicas, ",", "")
			arrPensionesPublicas := strings.Split(strPensionesPublicas, ".")
			auxPensionesPublicas, err := strconv.Atoi(arrPensionesPublicas[0])
			if err == nil {
				pensionesPublicas += auxPensionesPublicas
			}
		}

		if aux["pensionesPrivado"] != nil {
			strPensionesPrivadas := strings.TrimLeft(aux["pensionesPrivado"].(string), "$")
			strPensionesPrivadas = strings.ReplaceAll(strPensionesPrivadas, ",", "")
			arrPensionesPrivadas := strings.Split(strPensionesPrivadas, ".")
			auxPensionesPrivadas, err := strconv.Atoi(arrPensionesPrivadas[0])
			if err == nil {
				pensionesPrivadas += auxPensionesPrivadas
			}
		}

		if aux["caja"] != nil {
			strCaja := strings.TrimLeft(aux["caja"].(string), "$")
			strCaja = strings.ReplaceAll(strCaja, ",", "")
			arrCaja := strings.Split(strCaja, ".")
			auxCaja, err := strconv.Atoi(arrCaja[0])
			if err == nil {
				caja += auxCaja
			}

		}

		if aux["totalArl"] != nil {
			strArl := strings.TrimLeft(aux["totalArl"].(string), "$")
			strArl = strings.ReplaceAll(strArl, ",", "")
			arrArl := strings.Split(strArl, ".")
			auxArl, err := strconv.Atoi(arrArl[0])
			if err == nil {
				arl += auxArl
			}
		}

		if aux["icbf"] != nil {
			strIcbf := strings.TrimLeft(aux["icbf"].(string), "$")
			strIcbf = strings.ReplaceAll(strIcbf, ",", "")
			arrIcbf := strings.Split(strIcbf, ".")
			auxIcbf, err := strconv.Atoi(arrIcbf[0])
			if err == nil {
				icbf += auxIcbf
			}
		}
	}

	for i := 0; i < len(rhvPre); i++ {
		aux := rhvPre[i]
		if aux["primaServicios"] != nil {
			strPrimaServicios := strings.TrimLeft(aux["primaServicios"].(string), "$")
			strPrimaServicios = strings.ReplaceAll(strPrimaServicios, ",", "")
			arrPrimaServicios := strings.Split(strPrimaServicios, ".")
			auxPrimaServicios, err := strconv.Atoi(arrPrimaServicios[0])
			if err == nil {
				primaServicios += auxPrimaServicios
			}
		}

		if aux["primaNavidad"] != nil {
			strPrimaNavidad := strings.TrimLeft(aux["primaNavidad"].(string), "$")
			strPrimaNavidad = strings.ReplaceAll(strPrimaNavidad, ",", "")
			arrPrimaNavidad := strings.Split(strPrimaNavidad, ".")
			auxPrimaNavidad, err := strconv.Atoi(arrPrimaNavidad[0])
			if err == nil {
				primaNavidad += auxPrimaNavidad
			}
		}

		if aux["primaVacaciones"] != nil {
			strPrimaVacaciones := strings.TrimLeft(aux["primaVacaciones"].(string), "$")
			strPrimaVacaciones = strings.ReplaceAll(strPrimaVacaciones, ",", "")
			arrPrimaVacaiones := strings.Split(strPrimaVacaciones, ".")
			auxPrimaVacaciones, err := strconv.Atoi(arrPrimaVacaiones[0])
			if err == nil {
				primaVacaciones += auxPrimaVacaciones
			}
		}

		if aux["bonificacion"] != nil || aux["bonificacion"] != "N/A" {
			strBonificacion := strings.TrimLeft(aux["bonificacion"].(string), "$")
			strBonificacion = strings.ReplaceAll(strBonificacion, ",", "")
			arrBonificacion := strings.Split(strBonificacion, ".")
			auxBonificacion, err := strconv.Atoi(arrBonificacion[0])
			if err == nil {
				bonificacion += auxBonificacion
			}
		}

		if aux["interesesCesantias"] != nil || aux["interesesCesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["interesesCesantias"].(string), "$")
			strInteresesCesantias = strings.ReplaceAll(strInteresesCesantias, ",", "")
			arrInteresesCesantias := strings.Split(strInteresesCesantias, ".")
			auxInteresesCesantias, err := strconv.Atoi(arrInteresesCesantias[0])
			if err == nil {
				interesesCesantias += auxInteresesCesantias
			}
		}

		if aux["cesantiasPublico"] != nil {
			strCesantiasPublico := strings.TrimLeft(aux["cesantiasPublico"].(string), "$")
			strCesantiasPublico = strings.ReplaceAll(strCesantiasPublico, ",", "")
			arrCesantiasPublico := strings.Split(strCesantiasPublico, ".")
			auxCesantiasPublico, err := strconv.Atoi(arrCesantiasPublico[0])
			if err == nil {
				cesantiasPublicas += auxCesantiasPublico
			}
		}

		if aux["cesantiasPrivado"] != nil {
			strCesantiasPrivado := strings.TrimLeft(aux["cesantiasPrivado"].(string), "$")
			strCesantiasPrivado = strings.ReplaceAll(strCesantiasPrivado, ",", "")
			arrCesantiasPrivado := strings.Split(strCesantiasPrivado, ".")
			auxCesantiasPrivado, err := strconv.Atoi(arrCesantiasPrivado[0])
			if err == nil {
				cesantiasPrivadas += auxCesantiasPrivado
			}
		}

		if aux["totalSalud"] != nil {
			strSalud := strings.TrimLeft(aux["totalSalud"].(string), "$")
			strSalud = strings.ReplaceAll(strSalud, ",", "")
			arrSalud := strings.Split(strSalud, ".")
			auxSalud, err := strconv.Atoi(arrSalud[0])
			if err == nil {
				salud += auxSalud
			}
		}

		if aux["pensionesPublico"] != nil {
			strPensionesPublicas := strings.TrimLeft(aux["pensionesPublico"].(string), "$")
			strPensionesPublicas = strings.ReplaceAll(strPensionesPublicas, ",", "")
			arrPensionesPublicas := strings.Split(strPensionesPublicas, ".")
			auxPensionesPublicas, err := strconv.Atoi(arrPensionesPublicas[0])
			if err == nil {
				pensionesPublicas += auxPensionesPublicas
			}
		}

		if aux["pensionesPrivado"] != nil {
			strPensionesPrivadas := strings.TrimLeft(aux["pensionesPrivado"].(string), "$")
			strPensionesPrivadas = strings.ReplaceAll(strPensionesPrivadas, ",", "")
			arrPensionesPrivadas := strings.Split(strPensionesPrivadas, ".")
			auxPensionesPrivadas, err := strconv.Atoi(arrPensionesPrivadas[0])
			if err == nil {
				pensionesPrivadas += auxPensionesPrivadas
			}
		}

		if aux["caja"] != nil {
			strCaja := strings.TrimLeft(aux["caja"].(string), "$")
			strCaja = strings.ReplaceAll(strCaja, ",", "")
			arrCaja := strings.Split(strCaja, ".")
			auxCaja, err := strconv.Atoi(arrCaja[0])
			if err == nil {
				caja += auxCaja
			}

		}

		if aux["totalArl"] != nil {
			strArl := strings.TrimLeft(aux["totalArl"].(string), "$")
			strArl = strings.ReplaceAll(strArl, ",", "")
			arrArl := strings.Split(strArl, ".")
			auxArl, err := strconv.Atoi(arrArl[0])
			if err == nil {
				arl += auxArl
			}
		}

		if aux["icbf"] != nil {
			strIcbf := strings.TrimLeft(aux["icbf"].(string), "$")
			strIcbf = strings.ReplaceAll(strIcbf, ",", "")
			arrIcbf := strings.Split(strIcbf, ".")
			auxIcbf, err := strconv.Atoi(arrIcbf[0])
			if err == nil {
				icbf += auxIcbf
			}
		}
	}

	for i := 0; i < len(rhvPos); i++ {
		aux := rhvPos[i]
		if aux["primaServicios"] != nil {
			strPrimaServicios := strings.TrimLeft(aux["primaServicios"].(string), "$")
			strPrimaServicios = strings.ReplaceAll(strPrimaServicios, ",", "")
			arrPrimaServicios := strings.Split(strPrimaServicios, ".")
			auxPrimaServicios, err := strconv.Atoi(arrPrimaServicios[0])
			if err == nil {
				primaServicios += auxPrimaServicios
			}
		}

		if aux["primaNavidad"] != nil {
			strPrimaNavidad := strings.TrimLeft(aux["primaNavidad"].(string), "$")
			strPrimaNavidad = strings.ReplaceAll(strPrimaNavidad, ",", "")
			arrPrimaNavidad := strings.Split(strPrimaNavidad, ".")
			auxPrimaNavidad, err := strconv.Atoi(arrPrimaNavidad[0])
			if err == nil {
				primaNavidad += auxPrimaNavidad
			}
		}

		if aux["primaVacaciones"] != nil {
			strPrimaVacaciones := strings.TrimLeft(aux["primaVacaciones"].(string), "$")
			strPrimaVacaciones = strings.ReplaceAll(strPrimaVacaciones, ",", "")
			arrPrimaVacaiones := strings.Split(strPrimaVacaciones, ".")
			auxPrimaVacaciones, err := strconv.Atoi(arrPrimaVacaiones[0])
			if err == nil {
				primaVacaciones += auxPrimaVacaciones
			}
		}

		if aux["bonificacion"] != nil || aux["bonificacion"] != "N/A" {
			strBonificacion := strings.TrimLeft(aux["bonificacion"].(string), "$")
			strBonificacion = strings.ReplaceAll(strBonificacion, ",", "")
			arrBonificacion := strings.Split(strBonificacion, ".")
			auxBonificacion, err := strconv.Atoi(arrBonificacion[0])
			if err == nil {
				bonificacion += auxBonificacion
			}
		}

		if aux["interesesCesantias"] != nil || aux["interesesCesantias"] != "N/A" {
			strInteresesCesantias := strings.TrimLeft(aux["interesesCesantias"].(string), "$")
			strInteresesCesantias = strings.ReplaceAll(strInteresesCesantias, ",", "")
			arrInteresesCesantias := strings.Split(strInteresesCesantias, ".")
			auxInteresesCesantias, err := strconv.Atoi(arrInteresesCesantias[0])
			if err == nil {
				interesesCesantias += auxInteresesCesantias
			}
		}

		if aux["cesantiasPublico"] != nil {
			strCesantiasPublico := strings.TrimLeft(aux["cesantiasPublico"].(string), "$")
			strCesantiasPublico = strings.ReplaceAll(strCesantiasPublico, ",", "")
			arrCesantiasPublico := strings.Split(strCesantiasPublico, ".")
			auxCesantiasPublico, err := strconv.Atoi(arrCesantiasPublico[0])
			if err == nil {
				cesantiasPublicas += auxCesantiasPublico
			}
		}

		if aux["cesantiasPrivado"] != nil {
			strCesantiasPrivado := strings.TrimLeft(aux["cesantiasPrivado"].(string), "$")
			strCesantiasPrivado = strings.ReplaceAll(strCesantiasPrivado, ",", "")
			arrCesantiasPrivado := strings.Split(strCesantiasPrivado, ".")
			auxCesantiasPrivado, err := strconv.Atoi(arrCesantiasPrivado[0])
			if err == nil {
				cesantiasPrivadas += auxCesantiasPrivado
			}
		}

		if aux["totalSalud"] != nil {
			strSalud := strings.TrimLeft(aux["totalSalud"].(string), "$")
			strSalud = strings.ReplaceAll(strSalud, ",", "")
			arrSalud := strings.Split(strSalud, ".")
			auxSalud, err := strconv.Atoi(arrSalud[0])
			if err == nil {
				salud += auxSalud
			}
		}

		if aux["pensionesPublico"] != nil {
			strPensionesPublicas := strings.TrimLeft(aux["pensionesPublico"].(string), "$")
			strPensionesPublicas = strings.ReplaceAll(strPensionesPublicas, ",", "")
			arrPensionesPublicas := strings.Split(strPensionesPublicas, ".")
			auxPensionesPublicas, err := strconv.Atoi(arrPensionesPublicas[0])
			if err == nil {
				pensionesPublicas += auxPensionesPublicas
			}
		}

		if aux["pensionesPrivado"] != nil {
			strPensionesPrivadas := strings.TrimLeft(aux["pensionesPrivado"].(string), "$")
			strPensionesPrivadas = strings.ReplaceAll(strPensionesPrivadas, ",", "")
			arrPensionesPrivadas := strings.Split(strPensionesPrivadas, ".")
			auxPensionesPrivadas, err := strconv.Atoi(arrPensionesPrivadas[0])
			if err == nil {
				pensionesPrivadas += auxPensionesPrivadas
			}
		}

		if aux["caja"] != nil {
			strCaja := strings.TrimLeft(aux["caja"].(string), "$")
			strCaja = strings.ReplaceAll(strCaja, ",", "")
			arrCaja := strings.Split(strCaja, ".")
			auxCaja, err := strconv.Atoi(arrCaja[0])
			if err == nil {
				caja += auxCaja
			}

		}

		if aux["totalArl"] != nil {
			strArl := strings.TrimLeft(aux["totalArl"].(string), "$")
			strArl = strings.ReplaceAll(strArl, ",", "")
			arrArl := strings.Split(strArl, ".")
			auxArl, err := strconv.Atoi(arrArl[0])
			if err == nil {
				arl += auxArl
			}
		}

		if aux["icbf"] != nil {
			strIcbf := strings.TrimLeft(aux["icbf"].(string), "$")
			strIcbf = strings.ReplaceAll(strIcbf, ",", "")
			arrIcbf := strings.Split(strIcbf, ".")
			auxIcbf, err := strconv.Atoi(arrIcbf[0])
			if err == nil {
				icbf += auxIcbf
			}
		}
	}

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
