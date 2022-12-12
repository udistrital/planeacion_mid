package seguimientohelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"log"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

func GetTrimestres(vigencia string) []map[string]interface{} {

	var res map[string]interface{}
	var trimestre []map[string]interface{}
	var trimestres []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T1", &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &trimestre)
		trimestres = append(trimestres, trimestre...)

		trimestre = nil
		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T2", &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &trimestre)
			trimestres = append(trimestres, trimestre...)

			trimestre = nil
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T3", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &trimestre)
				trimestres = append(trimestres, trimestre...)

				trimestre = nil
				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T4", &res); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &trimestre)
					trimestres = append(trimestres, trimestre...)
				} else {
					panic(map[string]interface{}{"funcion": "GetTrimestres", "err": "Error ", "status": "400", "log": err})
				}
			} else {
				panic(map[string]interface{}{"funcion": "GetTrimestres", "err": "Error ", "status": "400", "log": err})
			}
		} else {
			panic(map[string]interface{}{"funcion": "GetTrimestres", "err": "Error ", "status": "400", "log": err})
		}
	} else {
		panic(map[string]interface{}{"funcion": "GetTrimestres", "err": "Error ", "status": "400", "log": err})
	}

	return trimestres
}

func GetActividades(subgrupo_id string) []map[string]interface{} {
	var res map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var datoPlan map[string]interface{}
	var actividades []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+subgrupo_id, &res); err == nil {
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
			sort.SliceStable(actividades, func(i, j int) bool {
				var a int
				var b int
				if reflect.TypeOf(actividades[j]["index"]).String() == "string" {
					b, _ = strconv.Atoi(actividades[j]["index"].(string))
				} else {
					b = int(actividades[j]["index"].(float64))
				}

				if reflect.TypeOf(actividades[i]["index"]).String() == "string" {
					a, _ = strconv.Atoi(actividades[i]["index"].(string))
				} else {
					a = int(actividades[i]["index"].(float64))
				}
				return a < b
			})
		}
	} else {
		panic(map[string]interface{}{"Code": "400", "Body": err, "Type": "error"})

	}
	return actividades
}

func GetActividad(seguimiento map[string]interface{}, index string, trimestre string) map[string]interface{} {
	var data map[string]interface{}
	var resEstado map[string]interface{}
	var informacion map[string]interface{}
	var cuantitativo map[string]interface{}
	cualitativo := map[string]interface{}{}
	evidencia := []interface{}{}
	estado := map[string]interface{}{}

	dato := make(map[string]interface{})
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if dato[index] != nil {
		if dato[index].(map[string]interface{})["informacion"] == nil {
			informacion = GetInformacionPlan(seguimiento, index)
		} else {
			informacion = dato[index].(map[string]interface{})["informacion"].(map[string]interface{})
		}

		if dato[index].(map[string]interface{})["cuantitativo"] == nil {
			cuantitativo = GetCuantitativoPlan(seguimiento, index, trimestre)
		} else {
			cuantitativo = dato[index].(map[string]interface{})["cuantitativo"].(map[string]interface{})
		}

		if dato[index].(map[string]interface{})["evidencia"] != nil {
			evidencia = dato[index].(map[string]interface{})["evidencia"].([]interface{})
		}
		if dato[index].(map[string]interface{})["estado"] == nil {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &resEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			}
		} else {
			estado = dato[index].(map[string]interface{})["estado"].(map[string]interface{})
		}

		if dato[index].(map[string]interface{})["cualitativo"] == nil {
			cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
		} else {
			cualitativo = dato[index].(map[string]interface{})["cualitativo"].(map[string]interface{})
		}
	} else {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &resEstado); err == nil {
			estado = map[string]interface{}{
				"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
				"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
			}
		}
		informacion = GetInformacionPlan(seguimiento, index)
		cuantitativo = GetCuantitativoPlan(seguimiento, index, trimestre)
		cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
	}

	data = map[string]interface{}{
		"informacion":  informacion,
		"cualitativo":  cualitativo,
		"cuantitativo": cuantitativo,
		"estado":       estado,
		"evidencia":    evidencia,
	}

	return data
}

func GetInformacionPlan(seguimiento map[string]interface{}, index string) map[string]interface{} {
	var resPlan map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var resInformacion map[string]interface{}
	var respuestaDependencia []map[string]interface{}
	var hijos []map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var periodo []map[string]interface{}

	informacion := map[string]interface{}{
		"ponderacion": "",
		"periodo":     "",
		"tarea":       "",
		"indicador":   "",
		"producto":    "",
		"nombre":      "",
		"descripcion": "",
		"index":       index,
		"unidad":      "",
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+seguimiento["plan_id"].(string), &resPlan); err == nil {
		informacion["nombre"] = resPlan["Data"].(map[string]interface{})["nombre"]
		informacion["descripcion"] = resPlan["Data"].(map[string]interface{})["descripcion"]
		informacion["unidad"] = resPlan["Data"].(map[string]interface{})["dependencia_id"]
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
			helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
			informacion["trimestre"] = periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"]
		}
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+seguimiento["plan_id"].(string), &resInformacion); err == nil {
		helpers.LimpiezaRespuestaRefactor(resInformacion, &hijos)
		for _, hijo := range hijos {
			if hijo["activo"] == true {
				var res map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijo["_id"].(string), &res); err == nil {
					dato := make(map[string]interface{})
					json.Unmarshal([]byte(res["Data"].([]interface{})[0].(map[string]interface{})["dato_plan"].(string)), &dato)
					nombreDetalle := strings.ToLower(res["Data"].([]interface{})[0].(map[string]interface{})["nombre"].(string))

					if dato[index] == nil {
						break
					}

					switch {
					case strings.Contains(nombreDetalle, "ponderación"):
						informacion["ponderacion"] = dato[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreDetalle, "periodo"):
						informacion["periodo"] = dato[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreDetalle, "tareas"):
						informacion["tarea"] = dato[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreDetalle, "indicadores"):
						informacion["indicador"] = dato[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreDetalle, "producto"):
						informacion["producto"] = dato[index].(map[string]interface{})["dato"]
						continue
					}
				}
			}
		}
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+informacion["unidad"].(string), &respuestaDependencia); err == nil {
		informacion["unidad"] = respuestaDependencia[0]["DependenciaId"].(map[string]interface{})["Nombre"]
	} else {
		informacion["unidad"] = nil
	}

	return informacion
}

func GetCuantitativoPlan(seguimiento map[string]interface{}, index string, trimestre string) map[string]interface{} {
	var resInformacion map[string]interface{}
	var resDetalle map[string]interface{}
	var hijos []interface{}
	var subgrupos []map[string]interface{}
	var indicadores []map[string]interface{}
	var respuestas []map[string]interface{}
	var subgrupo_detalle []map[string]interface{}
	response := map[string]interface{}{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+seguimiento["plan_id"].(string), &resInformacion); err == nil {
		helpers.LimpiezaRespuestaRefactor(resInformacion, &subgrupos)
		for _, subgrupo := range subgrupos {
			if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") && subgrupo["activo"] == true {

				hijos = subgrupo["hijos"].([]interface{})
				hijos = append(hijos, subgrupo["_id"])

				for _, hijo := range hijos {
					var res map[string]interface{}
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijo.(string), &res); err == nil {
						hijosIndicadores := res["Data"].(map[string]interface{})["hijos"].([]interface{})

						var dato_plan map[string]interface{}

						informacion := map[string]interface{}{
							"reporteNumerador":   0,
							"reporteDenominador": 1,
							"detalleReporte":     "",
						}

						respuesta := map[string]interface{}{
							"indicador":            0,
							"indicadorAcumulado":   0,
							"avanceAcumulado":      0,
							"brechaExistente":      0,
							"acumuladoNumerador":   0,
							"acumuladoDenominador": 0,
						}

						for _, hijoI := range hijosIndicadores {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijoI.(string), &resDetalle); err == nil {
								helpers.LimpiezaRespuestaRefactor(resDetalle, &subgrupo_detalle)

								if len(subgrupo_detalle) > 0 {
									if subgrupo_detalle[0]["dato_plan"] != nil {
										dato_plan_str := subgrupo_detalle[0]["dato_plan"].(string)
										json.Unmarshal([]byte(dato_plan_str), &dato_plan)
										nombreDetalle := strings.ToLower(subgrupo_detalle[0]["nombre"].(string))

										if dato_plan[index] == nil || dato_plan[index].(map[string]interface{})["dato"] == "" {
											break
										}

										switch {
										case strings.Contains(nombreDetalle, "nombre"):
											informacion["nombre"] = dato_plan[index].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "meta"):
											informacion["meta"] = dato_plan[index].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "fórmula"):
											informacion["formula"] = dato_plan[index].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "criterio"):
											informacion["denominador"] = dato_plan[index].(map[string]interface{})["dato"]
											if informacion["denominador"] == "Denominador fijo" {
												informacion["reporteDenominador"] = GetDenominadorFijo(seguimiento, len(indicadores), index)
											}
											continue
										case strings.Contains(nombreDetalle, "tendencia"):
											informacion["tendencia"] = dato_plan[index].(map[string]interface{})["dato"]
											continue
										}
									}
								}
							}
						}

						if informacion["nombre"] != nil && informacion["nombre"] != "" {
							indicadores = append(indicadores, informacion)
							respuestas = append(respuestas, respuesta)
						}

						if informacion["denominador"] == nil {
							informacion["denominador"] = ""
						}

						respuestas = GetRespuestaAcumulado(seguimiento, len(indicadores)-1, respuestas, index, trimestre, informacion["denominador"].(string))
					}
				}
				break
			}
		}
	}

	response["indicadores"] = indicadores
	response["resultados"] = respuestas
	return response
}

func GetDataSubgrupos(subgrupos []map[string]interface{}, index string) map[string]interface{} {
	var data map[string]interface{}
	auxSubgrupo := make(map[string]interface{})

	for i := range subgrupos {
		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
			aux := GetSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
			auxSubgrupo["actividad"] = aux["dato"]

		}

		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "lineamiento") {
			aux := GetSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
			auxSubgrupo["lineamiento"] = aux["dato"]

		}

		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "meta") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "estratégica") {
			aux := GetSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
			auxSubgrupo["meta_estrategica"] = aux["dato"]

		}

		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "estrategia") {
			aux := GetSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
			auxSubgrupo["estrategia"] = aux["dato"]

		}

		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "indicador") {
			var res map[string]interface{}
			var hijos []map[string]interface{}
			var indicadores []map[string]interface{}
			var metas []map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+subgrupos[i]["_id"].(string), &res); err != nil {
				panic(map[string]interface{}{"funcion": "GetDataSubgrupos", "err": "Error get indicador \"key\"", "status": "400", "log": err})
			}

			helpers.LimpiezaRespuestaRefactor(res, &hijos)
			for j := range hijos {
				if strings.Contains(strings.ToLower(hijos[j]["nombre"].(string)), "indicador") {
					aux := GetSubgrupoDetalle(hijos[j]["_id"].(string), index)
					auxSubgrupo["indicador"] = aux["dato"]
					indicadores = append(indicadores, aux)
				}
				if strings.Contains(strings.ToLower(hijos[j]["nombre"].(string)), "meta") {
					aux := GetSubgrupoDetalle(hijos[j]["_id"].(string), index)
					auxSubgrupo["meta"] = aux["dato"]
					metas = append(metas, aux)

				}
			}
			auxSubgrupo["indicador"] = indicadores
			auxSubgrupo["meta"] = metas

		}

		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "tarea") {
			aux := GetSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
			auxSubgrupo["tarea"] = aux["dato"]
		}

		data = auxSubgrupo
	}

	if data["lineamiento"] == nil {
		data["lineamiento"] = "No Aplica"
	}

	if data["meta_estrategica"] == nil {
		data["meta_estrategica"] = "No Aplica"
	}

	if data["estrategia"] == nil {
		data["estrategia"] = "No Aplica"
	}

	if data["indicador"] == nil {
		data["indicador"] = "No Aplica"
		data["meta"] = "No Aplica"
	}

	if data["tarea"] == nil {
		data["tarea"] = "No Aplica"
	}
	return data
}

func GetDenominadorFijo(dataSeg map[string]interface{}, index int, indexActividad string) float64 {
	plan_id := dataSeg["plan_id"].(string)
	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var seguimientos []map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var periodo []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+plan_id, &resSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)

		for _, seguimiento := range seguimientos {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
					helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
					trimestre := periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"]
					if trimestre == "T1" {
						if seguimiento["dato"] != "{}" {
							dato := make(map[string]interface{})
							datoStr := seguimiento["dato"].(string)
							json.Unmarshal([]byte(datoStr), &dato)
							if dato[indexActividad] == nil {
								break
							}

							seguimientoActividad := dato[indexActividad].(map[string]interface{})
							if seguimientoActividad["cuantitativo"] == nil {
								break
							}
							return seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"].(float64)
						}
						break
					}
				}
			}
		}
	}

	return 1
}

func GetRespuestaAcumulado(dataSeg map[string]interface{}, index int, respuestas []map[string]interface{}, indexActividad string, trimestre string, denominador string) []map[string]interface{} {
	plan_id := dataSeg["plan_id"].(string)
	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var seguimientos []map[string]interface{}
	var periodo []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+plan_id, &resSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)

		acumuladoNumerador := 0.0
		acumuladoDenominador := 0.0
		indicadorAcumulado := 0.0
		avanceAcumulado := 0.0
		brechaExistente := 0.0
		for _, seguimiento := range seguimientos {

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
					helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
					tri, _ := strconv.Atoi(string(trimestre[1]))
					segTrimestre, _ := strconv.Atoi(string(periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)[1]))

					if tri >= segTrimestre {
						if seguimiento["dato"] != "{}" {
							dato := make(map[string]interface{})
							datoStr := seguimiento["dato"].(string)
							json.Unmarshal([]byte(datoStr), &dato)
							if dato[indexActividad] == nil {
								respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
								respuestas[index]["avanceAcumulado"] = avanceAcumulado
								respuestas[index]["brechaExistente"] = brechaExistente
								continue
							}

							seguimientoActividad := dato[indexActividad].(map[string]interface{})
							if seguimientoActividad["cuantitativo"] == nil {
								respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
								respuestas[index]["avanceAcumulado"] = avanceAcumulado
								respuestas[index]["brechaExistente"] = brechaExistente
								continue
							}

							if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["indicadorAcumulado"] != nil {
								indicadorAcumulado += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["indicadorAcumulado"].(float64)
							}

							if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["avanceAcumulado"] != nil {
								avanceAcumulado += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["avanceAcumulado"].(float64)
							}

							if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["brechaExistente"] != nil {
								brechaExistente += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["brechaExistente"].(float64)
							}

							if denominador == "Denominador fijo" {
								acumuladoDenominador = seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"].(float64)
							} else {
								acumuladoDenominador += seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"].(float64)
							}

							if seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteNumerador"] == nil {
								acumuladoNumerador += seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteNumerador"].(float64)
							}
						}
						respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
						respuestas[index]["avanceAcumulado"] = avanceAcumulado
						respuestas[index]["brechaExistente"] = brechaExistente
						respuestas[index]["acumuladoNumerador"] = acumuladoNumerador
						respuestas[index]["acumuladoDenominador"] = acumuladoDenominador
					}
				}
			}

		}
	}

	return respuestas
}

func GetSubgrupoDetalle(subgrupo_id string, index string) map[string]interface{} {
	var respuesta map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var datoPlan map[string]interface{}
	var data map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo_id, &respuesta); err != nil {
		panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
	}
	aux := make([]map[string]interface{}, 1)
	helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
	subgrupoDetalle = aux[0]
	if subgrupoDetalle["dato_plan"] != nil {
		dato_plan_str := subgrupoDetalle["dato_plan"].(string)
		json.Unmarshal([]byte(dato_plan_str), &datoPlan)

		data = datoPlan[index].(map[string]interface{})
	}
	return data
}

func GetEstadoSeguimiento(seguimiento map[string]interface{}) string {
	var resEstado map[string]interface{}
	enReporte := true
	estado := map[string]interface{}{}
	dato := make(map[string]interface{})

	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+seguimiento["estado_seguimiento_id"].(string), &resEstado); err == nil {
		estado = map[string]interface{}{
			"nombre": resEstado["Data"].(map[string]interface{})["nombre"],
			"id":     resEstado["Data"].(map[string]interface{})["_id"],
		}

		for _, actividad := range dato {
			if actividad.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad en reporte" {
				enReporte = false
			}
		}

		if enReporte {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:ER", &resEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			}
		}
	}

	return estado["id"].(string)
}

func ActividadReportable(seguimiento map[string]interface{}, indexActividad string) (bool, map[string]interface{}) {
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if dato[indexActividad] == nil {
		return false, map[string]interface{}{"error": 1, "motivo": "Actividad sin seguimiento"}
	} else {
		estado = dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})
		if estado["nombre"] != "Actividad en reporte" {
			return false, map[string]interface{}{"error": 2, "motivo": "El estado de la actividad no es el adecuado"}
		}

		cuantitativo := dato[indexActividad].(map[string]interface{})["cuantitativo"]
		cualitativo := dato[indexActividad].(map[string]interface{})["cualitativo"]

		if cuantitativo == nil {
			return false, map[string]interface{}{"error": 3, "motivo": "Componenten cuantitativo sin guardar"}
		}

		if cualitativo == nil {
			return false, map[string]interface{}{"error": 4, "motivo": "Componenten cualitativo sin guardar"}
		} else {
			cualitativo := cualitativo.(map[string]interface{})
			if cualitativo["dificultades"] == "" || cualitativo["productos"] == "" || cualitativo["reporte"] == "" {
				return false, map[string]interface{}{"error": 5, "motivo": "Campos vacios en el componenten cualitativo"}
			}
		}
	}

	return true, nil
}

func SeguimientoReportable(seguimiento map[string]interface{}) (bool, map[string]interface{}) {
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var datoPlan map[string]interface{}

	dato := make(map[string]interface{})

	planId := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {

				actividades := GetActividades(subgrupos[i]["_id"].(string))
				if seguimiento["dato"] == "{}" {
					for _, actividad := range actividades {
						dato[actividad["index"].(string)] = actividad["dato"]
					}
					return false, map[string]interface{}{"error": 1, "motivo": "No hay actividades resportadas", "actividades": dato}
				} else {
					dato_plan_str := seguimiento["dato"].(string)
					json.Unmarshal([]byte(dato_plan_str), &datoPlan)

					for indexActividad, element := range datoPlan {
						for _, actividad := range actividades {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								if indexActividad == actividad["index"] {
									actividad["estado"] = element.(map[string]interface{})["estado"]
								}
							} else {
								if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
									actividad["estado"] = element.(map[string]interface{})["estado"]
								}
							}
						}
					}

					for _, actividad := range actividades {
						if actividad["estado"] == nil {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								dato[actividad["index"].(string)] = actividad["dato"]
							} else {
								dato[strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)] = actividad["dato"]
							}
						} else if actividad["estado"].(map[string]interface{})["nombre"] != "Actividad reportada" && actividad["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								dato[actividad["index"].(string)] = actividad["dato"]
							} else {
								dato[strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)] = actividad["dato"]
							}
						}
					}

					if fmt.Sprintf("%v", dato) != "map[]" {
						return false, map[string]interface{}{"error": 2, "motivo": "Hay actividades sin resportar", "actividades": dato}
					} else {
						return true, nil
					}
				}
			}
		}
	}
	return true, nil
}

func ActividadConObservaciones(seguimiento map[string]interface{}) bool {
	var cuantitativo map[string]interface{}
	var cualitativo map[string]interface{}

	if seguimiento["cuantitativo"] != nil {
		cuantitativo = seguimiento["cuantitativo"].(map[string]interface{})
		for _, indicador := range cuantitativo["indicadores"].([]interface{}) {
			if indicador.(map[string]interface{})["observaciones"] != "" && indicador.(map[string]interface{})["observaciones"] != "Sin observación" && indicador.(map[string]interface{})["observaciones"] != nil {
				return true
			}
		}
	}

	if seguimiento["cualitativo"] != nil {
		cualitativo = seguimiento["cualitativo"].(map[string]interface{})
		if cualitativo["observaciones"] != "" && cualitativo["observaciones"] != "Sin observación" && cualitativo["observaciones"] != nil {
			return true
		}
	}

	if seguimiento["evidencias"] != nil {
		for _, evidencia := range seguimiento["evidencias"].([]interface{}) {
			if evidencia.(map[string]interface{})["Observacion"] != "" && evidencia.(map[string]interface{})["Observacion"] != "Sin observación" && evidencia.(map[string]interface{})["Observacion"] != nil {
				return true
			}
		}
	}

	return false
}

func SeguimientoAvalable(seguimiento map[string]interface{}) (bool, bool, map[string]interface{}) {
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var datoPlan map[string]interface{}
	dato := make(map[string]interface{})
	observaciones := false
	avaladas := false

	planId := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {

				actividades := GetActividades(subgrupos[i]["_id"].(string))

				dato_plan_str := seguimiento["dato"].(string)
				json.Unmarshal([]byte(dato_plan_str), &datoPlan)

				for indexActividad, element := range datoPlan {
					for _, actividad := range actividades {

						if reflect.TypeOf(actividad["index"]).String() == "string" {
							if indexActividad == actividad["index"] {
								if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
									dato[indexActividad] = actividad["dato"]
								}
							}
						} else {
							if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
								if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
									dato[indexActividad] = actividad["dato"]
								}
							}
						}
					}

					if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Con observaciones" {
						observaciones = true
					}

					if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Actividad avalada" {
						avaladas = true
					}
				}
			}
		}
	}

	if fmt.Sprintf("%v", dato) != "map[]" {
		return false, false, map[string]interface{}{"error": 1, "motivo": "Hay actividades sin revisar", "actividades": dato}
	}

	return avaladas, observaciones, nil
}
