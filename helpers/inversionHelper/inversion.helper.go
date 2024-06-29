package inversionhelper

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

var validDataT = []string{}

func Limpia() {
	data_source = nil
	displayed_columns = nil
	validDataT = []string{}
}

func RegistrarProyecto(registroProyecto map[string]interface{}) map[string]interface{} {
	var respuestaProyecto map[string]interface{}
	plan := make(map[string]interface{})
	plan["activo"] = true
	plan["nombre"] = registroProyecto["nombre_proyecto"]
	plan["descripcion"] = registroProyecto["codigo_proyecto"]
	if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan", "POST", &respuestaProyecto, plan); err != nil {
		panic(map[string]interface{}{"funcion": "AddProyecto", "err": "Error versionando plan \"plan[\"_id\"].(string)\"", "status": "400", "log": err})
	}
	return respuestaProyecto
}

func ResgistrarInfoComplementaria(idProyecto string, infoProyecto map[string]interface{}, nombreCoplementaria string) error {
	var resSubgrupo map[string]interface{}
	infoSubgrupo := map[string]interface{}{
		"activo":      true,
		"padre":       idProyecto,
		"nombre":      nombreCoplementaria,
		"descripcion": infoProyecto["codigo_proyecto"],
	}
	errSubgrupo := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo", "POST", &resSubgrupo, infoSubgrupo)
	if errSubgrupo == nil {
		idSubgrupo := resSubgrupo["Data"].(map[string]interface{})["_id"].(string)
		detalle, _ := json.Marshal(infoProyecto["data"])
		var resDetalle map[string]interface{}

		subgrupoDetalle := map[string]interface{}{
			"activo":      true,
			"subgrupo_id": idSubgrupo,
			"nombre":      nombreCoplementaria,
			"descripcion": infoProyecto["codigo_proyecto"],
			"dato":        string(detalle),
		}

		errDetalle := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle", "POST", &resDetalle, subgrupoDetalle)
		return errDetalle
	}

	return errSubgrupo
}

func ActualizarInfoComplDetalle(idSubgrupo string, detalleData []interface{}) error {
	var resSubgrupo map[string]interface{}
	var subgrupo map[string]interface{}

	errGet := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+idSubgrupo, &resSubgrupo)
	if errGet == nil {
		helpers.LimpiezaRespuestaRefactor(resSubgrupo, &subgrupo)
		detalle, _ := json.Marshal(detalleData)
		subgrupo["dato"] = string(detalle)

		var resDetalle map[string]interface{}
		errDetalle := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+idSubgrupo, "PUT", &resDetalle, subgrupo)
		return errDetalle
	}
	return errGet
}

func ActualizarPresupuestoDisponible(infoFuente []interface{}) {
	for _, fuente := range infoFuente {
		var dataFuente map[string]interface{}
		err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/fuentes-apropiacion/"+fuente.(map[string]interface{})["id"].(string), &dataFuente)
		if err == nil {
			resFuente := dataFuente["Data"].(map[string]interface{})
			var dataFuente map[string]interface{}
			resFuente["presupuestoDisponible"] = fuente.(map[string]interface{})["presupuestoDisponible"]
			helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/fuentes-apropiacion/"+fuente.(map[string]interface{})["id"].(string), "PUT", &dataFuente, resFuente)
		}
	}
}

// func GetIdSbugrupoDetalle(padreId string) map[string]interface{} {

// 	var res []map[string]interface{}
// 	var infoSubgrupos map[string]interface{}
// 	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+padreId, &res); err == nil {
// 		if res[0]["Data"] != nil {

// 		}
// 		// for i := 0; i < len(res["Data"]); i++ {
// 		// 	if res["Data"][i]["nombre"] == "soportes" {
// 		// 		idSubgrupoSoportes = res["Data"][i]["_id"].(string)
// 		// 	}
// 		// }
// 		// s, err := json.Marshal(res["Data"])
// 		// fmt.Println(s[0], "primera posicion")
// 		// if err != nil {
// 		// 	panic(err)
// 		// }
// 		//fmt.Println(s)
// 		//json.Unmarshal(s, &infoSubgrupos)
// 		fmt.Println(infoSubgrupos)
// 		helpers.LimpiezaRespuestaRefactor(res[0], &infoSubgrupos)
// 		//fmt.Println(res, "respuesta subgrupos")
// 	}

//		return infoSubgrupos
//	}

func GetDataProyects(infoProyect map[string]interface{}) map[string]interface{} {
	getProyect := make(map[string]interface{})
	var subgruposData map[string]interface{}
	var infoSubgrupos []map[string]interface{}

	getProyect["nombre_proyecto"] = infoProyect["nombre"]
	getProyect["codigo_proyecto"] = infoProyect["descripcion"]
	getProyect["fecha_creacion"] = infoProyect["fecha_creacion"]
	getProyect["id"] = infoProyect["_id"]

	padreId := infoProyect["_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+padreId, &subgruposData); err == nil {
		helpers.LimpiezaRespuestaRefactor(subgruposData, &infoSubgrupos)
		for i := range infoSubgrupos {
			var subgrupoDetalle map[string]interface{}
			var detalleSubgrupos []map[string]interface{}
			if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "fuentes") {

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
					helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)

					armonizacion_dato_str := detalleSubgrupos[0]["dato"].(string)
					var subgrupo_dato []map[string]interface{}
					json.Unmarshal([]byte(armonizacion_dato_str), &subgrupo_dato)

					getProyect["subgrupo_id_fuentes"] = infoSubgrupos[i]["_id"]
					getProyect["fuentes"] = subgrupo_dato
					getProyect["id_detalle_fuentes"] = detalleSubgrupos[0]["_id"]
				}
			}
		}
	}

	return getProyect
}

func GetDataProyect(proyect map[string]interface{}) map[string]interface{} {
	getProyect := make(map[string]interface{})
	var subgruposData map[string]interface{}
	var infoSubgrupos []map[string]interface{}

	getProyect["nombre_proyecto"] = proyect["nombre"]
	getProyect["codigo_proyecto"] = proyect["descripcion"]
	getProyect["fecha_creacion"] = proyect["fecha_creacion"]
	getProyect["id"] = proyect["_id"]

	padreId := proyect["_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+padreId, &subgruposData); err == nil {
		helpers.LimpiezaRespuestaRefactor(subgruposData, &infoSubgrupos)
		for i := range infoSubgrupos {
			var subgrupoDetalle map[string]interface{}
			var detalleSubgrupos []map[string]interface{}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
				helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)

				armonizacion_dato_str := detalleSubgrupos[0]["dato"].(string)
				var subgrupo_dato []map[string]interface{}
				json.Unmarshal([]byte(armonizacion_dato_str), &subgrupo_dato)

				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "soporte") {
					getProyect["subgrupo_id_soportes"] = infoSubgrupos[i]["_id"]
					getProyect["soportes"] = subgrupo_dato
					getProyect["id_detalle_soportes"] = detalleSubgrupos[0]["_id"]
				}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "metas") {
					getProyect["subgrupo_id_metas"] = infoSubgrupos[i]["_id"]
					getProyect["metas"] = subgrupo_dato
					getProyect["id_detalle_metas"] = detalleSubgrupos[0]["_id"]
				}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "fuentes") {
					getProyect["subgrupo_id_fuentes"] = infoSubgrupos[i]["_id"]
					getProyect["fuentes"] = subgrupo_dato
					getProyect["id_detalle_fuentes"] = detalleSubgrupos[0]["_id"]
				}
			}
		}
	}

	return getProyect
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
					fmt.Println(dato_armonizacion, "armonizacion_dato")
					aux := dato_armonizacion[index]
					if aux != nil {
						armonizacion["idSubDetalleProI"] = aux.(map[string]interface{})["idSubDetalleProI"]
						armonizacion["indexMetaSubProI"] = aux.(map[string]interface{})["indexMetaSubProI"]
						armonizacion["indexMetaPlan"] = index

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
							if actividad["observacion"] != nil || actividad["observacion"] != ""{
								keyObservacion := v + "_o"
								forkData[keyObservacion] = getObservacion(actividad)
							} else {
								keyObservacion := v + "_o"
								forkData[keyObservacion] = ""
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

var data_source []map[string]interface{}
var displayed_columns []string

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
	var subgrupo_proyect map[string]interface{}
	var dataProyect map[string]interface{}
	var dato_plan map[string]interface{}
	var totalPresupuestoActividad int
	var armonizacion_dato map[string]interface{}
	var dato []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo["_id"].(string), &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
		subgrupo_detalle = respuestaLimpia[0]
		if data_source == nil {
			if subgrupo_detalle["dato_plan"] != nil {
				dato_plan_str := subgrupo_detalle["dato_plan"].(string)
				json.Unmarshal([]byte(dato_plan_str), &dato_plan)
				//fmt.Println(subgrupo_detalle, "subgrupoDetalle")

				armonizacion_dato_str := subgrupo_detalle["armonizacion_dato"].(string)
				json.Unmarshal([]byte(armonizacion_dato_str), &armonizacion_dato)

				fmt.Println(armonizacion_dato, "JSON")
				for key := range dato_plan {
					actividad := make(map[string]interface{})
					element := dato_plan[key].(map[string]interface{})
					actividad["index"] = key
					//i, _ := strconv.Atoi(key)

					actividad[subgrupo["nombre"].(string)] = element["dato"]
					actividad["activo"] = element["activo"]
					if armonizacion_dato[key] != nil {
						//fmt.Println(armonizacion_dato[i], "entra a armonizacion dato")
						dataProyect = armonizacion_dato[key].(map[string]interface{})
						actividad["presupuesto"] = dataProyect["presupuesto_programado"]
						actividad["presupuesto_programado"] = dataProyect["presupuesto_programado"]
						fmt.Println(actividad["presupuesto"], "presupuesto")
						if dataProyect["idSubDetalleProI"] != nil {

							idSubDetalleProI := dataProyect["idSubDetalleProI"].(string) //para las actividades no está guardado éste valor
							indexMetaSubProI := dataProyect["indexMetaSubProI"]
							//totalPresupuestoActividad = dataProyect["totalPresupuestoActividad"].(int)
							fmt.Println(totalPresupuestoActividad, "total")
							//fmt.Println(idSubDetalleProI, "idSubDetalleProI")
							var respuestaLimpia2 map[string]interface{}
							var res map[string]interface{}
							//var posicion
							//var indexPro int
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+idSubDetalleProI, &res); err != nil {
								panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
							}
							helpers.LimpiezaRespuestaRefactor(res, &respuestaLimpia2)
							// b, err := json.Marshal(res["Data"])
							// if err != nil {
							// 	panic(err)
							// }
							// json.Unmarshal(b, &respuestaLimpia2)
							//fmt.Println(respuestaLimpia2, "metas pro inversion")
							subgrupo_proyect = respuestaLimpia2
							dato_str := subgrupo_proyect["dato"].(string)
							json.Unmarshal([]byte(dato_str), &dato)
							//fmt.Println(dato, "dato")
							for key2 := range dato {
								//fmt.Println("entra a dato")
								//metaProyect := dato[key2]
								// j, err := strconv.Atoi(indexMetaSubProI.(string))
								// if err != nil {
								// 	panic(err)
								// }
								//indexPro = j
								//posicion = dato[key2]["posicion"]
								//j = dato[key2]["posicion"].(string)
								//fmt.Println(dato[key2]["posicion"] == j, "j")
								if indexMetaSubProI == dato[key2]["posicion"] {
									//fmt.Println("entra a presupuesto")
									actividad["meta"] = dato[key2]["descripcion"]
									//actividad["presupuesto_programado"] = dato[key2]["presupuestoT"]
									actividad["posicion"] = dato[key2]["posicion"]
									actividad["indexMetaSubProI"] = indexMetaSubProI

								}

							}
						}
						actividad["indexMeta"] = dataProyect["indexMetaSubProI"]
						//actividad["presupuesto"] = dataProyect["totalPresupuestoActividad"]
					}
					data_source = append(data_source, actividad)
					//fmt.Println(actividad, "actividad")
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

func GetSons(hijos []interface{}, index string) {
	//tabla := make(map[string]interface{})
	var respuesta map[string]interface{}
	var subgrupo map[string]interface{}

	for i := 0; i < len(hijos); i++ {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijos[i].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &subgrupo)
			if subgrupo["bandera_tabla"] == true {
				DeleteMetas(subgrupo, index)
			}

			if len(subgrupo["hijos"].([]interface{})) != 0 {
				var respuestaHijos map[string]interface{}
				var subgrupoHijo map[string]interface{}
				for j := 0; j < len(subgrupo["hijos"].([]interface{})); j++ {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+subgrupo["hijos"].([]interface{})[j].(string), &respuestaHijos); err == nil {
						helpers.LimpiezaRespuestaRefactor(respuestaHijos, &subgrupoHijo)
						if subgrupoHijo["bandera_tabla"] == true {
							displayed_columns = append(displayed_columns, subgrupoHijo["nombre"].(string))
							DeleteMetas(subgrupoHijo, index)
						}
					}
				}
			}
		}
	}
	// tabla["displayed_columns"] = displayed_columns
	// tabla["data_source"] = data_source
	// return tabla
}

func DeleteMetas(subgrupo map[string]interface{}, index string) {
	var respuesta map[string]interface{}
	var res map[string]interface{}
	var respuestaLimpia []map[string]interface{}
	var subgrupo_detalle map[string]interface{}
	var dato_plan map[string]interface{}
	actividad := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo["_id"].(string), &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
		subgrupo_detalle = respuestaLimpia[0]

		if subgrupo_detalle["dato_plan"] != nil {
			dato_plan_str := subgrupo_detalle["dato_plan"].(string)
			json.Unmarshal([]byte(dato_plan_str), &dato_plan)

			for key := range dato_plan {
				if key == index {
					element := dato_plan[key].(map[string]interface{})
					actividad["index"] = element["index"]
					actividad["dato"] = element["dato"]
					actividad["activo"] = false
					if element["observacion"] != nil {
						actividad["observacion"] = element["observacion"]
					}
					dato_plan[key] = actividad

				}
			}
			b, _ := json.Marshal(dato_plan)
			str := string(b)
			subgrupo_detalle["dato_plan"] = str
		}
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteActividad", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
		}
		fmt.Println(res, "respuesta delete")
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
						desactivarMeta(subgrupo["_id"].(string), index)
					}
				} else {
					panic(map[string]interface{}{"funcion": "InactivarMeta", "err": "Error obteniendo subgrupo \"subgrupo[\"_id\"].(string)\"", "status": "400", "log": err})
				}
			}
		} else {
			desactivarMeta(hijos[i]["_id"].(string), index)
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
				desactivarMeta(subgrupo["_id"].(string), index)
			}
		} else {
			panic(map[string]interface{}{"funcion": "InactivarMeta", "err": "Error obteniendo subgrupo \"subgrupo[\"_id\"].(string)\"", "status": "400", "log": err})
		}
	}
}

func desactivarMeta(subgrupo_id string, index string) {
	var respuesta map[string]interface{}
	var res map[string]interface{}
	var respuestaLimpia []map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var dato_plan map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo_id, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
		subgrupoDetalle = respuestaLimpia[0]
		if subgrupoDetalle["dato_plan"] != nil {
			meta := make(map[string]interface{})
			dato_plan_str := subgrupoDetalle["dato_plan"].(string)
			json.Unmarshal([]byte(dato_plan_str), &dato_plan)
			for index_meta := range dato_plan {
				if index_meta == index {
					aux_meta := dato_plan[index_meta].(map[string]interface{})
					meta["activo"] = false
					meta["index"] = index_meta
					meta["dato"] = aux_meta["dato"]
					if aux_meta["observacion"] != nil {
						meta["observacion"] = aux_meta["observacion"]
					}
					dato_plan[index_meta] = meta
					fmt.Println(dato_plan[index_meta], "Dato Plan")
				}
			}
			b, _ := json.Marshal(dato_plan)
			str := string(b)
			subgrupoDetalle["dato_plan"] = str
		}

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupoDetalle["_id"].(string), "PUT", &res, subgrupoDetalle); err == nil {
			fmt.Println(res, "res 622")
		} else {
			panic(map[string]interface{}{"funcion": "InactivarMeta", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
		}

	} else {
		panic(map[string]interface{}{"funcion": "InactivarMeta", "err": "Error obteniendo subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
	}
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

func OrdenarVersiones(versiones []map[string]interface{}) []map[string]interface{} {
	var versionesOrdenadas []map[string]interface{}

	for i := range versiones {
		if versiones[i]["padre_plan_id"] == nil {
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
