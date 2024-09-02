package seguimientohelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	comunhelper "github.com/udistrital/planeacion_mid/helpers/comunHelper"
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

func ObtenerTrimestres(vigencia string) ([]map[string]interface{}, error) {

	var res map[string]interface{}
	var trimestre []map[string]interface{}
	var trimestres []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T1", &res); err != nil {
		return nil, errors.New("No se encontró el trimestre 1 para la vigencia seleccionada")
	}

	helpers.LimpiezaRespuestaRefactor(res, &trimestre)
	if len(trimestre[0]) <= 0 {
		return nil, errors.New("No se encontró el trimestre 1 para la vigencia seleccionada")
	}
	trimestres = append(trimestres, trimestre...)
	trimestre = nil

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T2", &res); err != nil {
		return nil, errors.New("No se encontró el trimestre 2 para la vigencia seleccionada")
	}

	helpers.LimpiezaRespuestaRefactor(res, &trimestre)
	if len(trimestre[0]) <= 0 {
		return nil, errors.New("No se encontró el trimestre 2 para la vigencia seleccionada")
	}
	trimestres = append(trimestres, trimestre...)
	trimestre = nil

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T3", &res); err != nil {
		return nil, errors.New("No se encontró el trimestre 3 para la vigencia seleccionada")
	}

	helpers.LimpiezaRespuestaRefactor(res, &trimestre)
	if len(trimestre[0]) <= 0 {
		return nil, errors.New("No se encontró el trimestre 3 para la vigencia seleccionada")
	}
	trimestres = append(trimestres, trimestre...)
	trimestre = nil

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T4", &res); err != nil {
		return nil, errors.New("No se encontró el trimestre 4 para la vigencia seleccionada")
	}

	helpers.LimpiezaRespuestaRefactor(res, &trimestre)
	if len(trimestre[0]) <= 0 {
		return nil, errors.New("No se encontró el trimestre 4 para la vigencia seleccionada")
	}
	trimestres = append(trimestres, trimestre...)

	return trimestres, nil

	// TODO: Analizar por qué no funciona el siguiente código
	// trimestres := make([]map[string]interface{}, 4)
	// errorTrimestres := make(chan error, 4)

	// var wg sync.WaitGroup
	// wg.Add(4)

	// for i := 0; i < 4; i++ {
	// 	go func(index int) {
	// 		defer wg.Done()
	// 		var numeroTrimestre = index + 1
	// 		var numeroTrimestreStr = strconv.Itoa(numeroTrimestre)
	// 		var res map[string]interface{}
	// 		beego.Info("http://" + beego.AppConfig.String("ParametrosService") + "/parametro_periodo?query=PeriodoId:" + vigencia + ",ParametroId__CodigoAbreviacion:T" + numeroTrimestreStr)
	// 		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T"+numeroTrimestreStr, &res); err != nil {
	// 			errorTrimestres <- err
	// 			return
	// 		}
	// 		// url := fmt.Sprintf("http://%s/parametro_periodo?query=PeriodoId:%s,ParametroId__CodigoAbreviacion:T%d", beego.AppConfig.String("ParametrosService"), vigencia, index+1)

	// 		var trimestre []map[string]interface{}
	// 		helpers.LimpiezaRespuestaRefactor(res, &trimestre)
	// 		beego.Info(res)
	// 		if len(trimestre) > 0 {
	// 			trimestres[index] = trimestre[0]
	// 		} else {
	// 			errorTrimestres <- errors.New("No se encontró el trimestre")
	// 		}
	// 	}(i)
	// }

	// wg.Wait()
	// close(errorTrimestres)

	// // Check for errors
	// for err := range errorTrimestres {
	// 	return nil, err
	// }

	// // Flatten trimestres slice
	// var result []map[string]interface{}
	// for _, trimestre := range trimestres {
	// 	if len(trimestre) == 0 {
	// 		return nil, errors.New("No se encontró el trimestre")
	// 	}
	// 	result = append(result, trimestre)
	// }

	// return result, nil
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
	var resDetalle map[string]interface{}
	var informacion map[string]interface{}
	var cuantitativo map[string]interface{}
	cualitativo := map[string]interface{}{}
	evidencia := []interface{}{}
	evidenciaSeg := []map[string]interface{}{}
	estado := map[string]interface{}{}
	detalle := map[string]interface{}{}
	id := ""
	dato := make(map[string]interface{})
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if dato[index] != nil {
		idS, segregado := dato[index].(map[string]interface{})["id"]

		if segregado && idS != "" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[index].(map[string]interface{})["id"].(string), &resDetalle); err == nil {
				if resDetalle["Data"] != "null" {
					helpers.LimpiezaRespuestaRefactor(resDetalle, &detalle)
					detalle = ConvertirStringJson(detalle)

					id = detalle["_id"].(string)

					if len(detalle["informacion"].(map[string]interface{})) == 0 {
						informacion = GetInformacionPlan(seguimiento, index)
					} else {
						informacion = detalle["informacion"].(map[string]interface{})
					}

					if len(detalle["cuantitativo"].(map[string]interface{})) == 0 {
						cuantitativo = GetCuantitativoPlan(seguimiento, index, trimestre)
					} else {
						cuantitativo = detalle["cuantitativo"].(map[string]interface{})
					}

					if len(detalle["cualitativo"].(map[string]interface{})) == 0 {
						cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
					} else {
						cualitativo = detalle["cualitativo"].(map[string]interface{})
					}

					if len(detalle["evidencia"].([]map[string]interface{})) != 0 {
						evidenciaSeg = detalle["evidencia"].([]map[string]interface{})
					}

					if len(detalle["estado"].(map[string]interface{})) == 0 {
						if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &resEstado); err == nil {
							estado = map[string]interface{}{
								"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
								"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
							}
						}
					} else {
						estado = detalle["estado"].(map[string]interface{})
					}
				}
			}
		} else {
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
		"id":           id,
		"informacion":  informacion,
		"cualitativo":  cualitativo,
		"cuantitativo": cuantitativo,
		"estado":       estado,
		"evidencia":    evidencia,
	}

	if id != "" {
		data["evidencia"] = evidenciaSeg
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
		"producto":    "",
		"nombre":      "",
		"descripcion": "",
		"index":       index,
		"unidad":      "",
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+seguimiento["plan_id"].(string), &resPlan); err == nil {
		informacion["nombre"] = resPlan["Data"].(map[string]interface{})["nombre"]
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
			nombreHijo := strings.ToLower(hijo["nombre"].(string))
			if hijo["activo"] == true {
				var res map[string]interface{}

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijo["_id"].(string), &res); err == nil {
					datoPlan := make(map[string]interface{})
					dato := make(map[string]interface{})

					nombreDetalle := strings.ToLower(res["Data"].([]interface{})[0].(map[string]interface{})["nombre"].(string))
					if strings.Contains(nombreDetalle, "indicadores") || strings.Contains(nombreDetalle, "indicador") {
						continue
					}

					json.Unmarshal([]byte(res["Data"].([]interface{})[0].(map[string]interface{})["dato"].(string)), &dato)
					if dato["required"] == false || dato["required"] == "false" {
						continue
					}

					json.Unmarshal([]byte(res["Data"].([]interface{})[0].(map[string]interface{})["dato_plan"].(string)), &datoPlan)

					if datoPlan[index] == nil {
						continue
					}

					switch {
					case strings.Contains(nombreHijo, "ponderación"):
						informacion["ponderacion"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "periodo") || strings.Contains(nombreHijo, "período"):
						informacion["periodo"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "tareas") || strings.Contains(nombreHijo, "actividades específicas"):
						informacion["tarea"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "producto"):
						informacion["producto"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "actividad"):
						informacion["descripcion"] = datoPlan[index].(map[string]interface{})["dato"]
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
	respuestas := make([]map[string]interface{}, 0)
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
							"detalleReporte": "",
						}

						respuesta := map[string]interface{}{
							"indicador":            0,
							"indicadorAcumulado":   0,
							"avanceAcumulado":      0,
							"brechaExistente":      0,
							"acumuladoNumerador":   0,
							"acumuladoDenominador": 0,
							"meta":                 0,
						}

						for _, hijoI := range hijosIndicadores {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijoI.(string), &resDetalle); err == nil {
								var subgrupo_detalle []map[string]interface{}
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
											respuesta["nombre"] = dato_plan[index].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "meta"):
											informacion["meta"] = dato_plan[index].(map[string]interface{})["dato"]
											if reflect.TypeOf(dato_plan[index].(map[string]interface{})["dato"]).String() == "string" {
												respuesta["meta"], _ = strconv.ParseFloat(dato_plan[index].(map[string]interface{})["dato"].(string), 64)
											} else {
												respuesta["meta"] = dato_plan[index].(map[string]interface{})["dato"].(float64)
											}
											continue
										case strings.Contains(nombreDetalle, "fórmula"):
											informacion["formula"] = dato_plan[index].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "criterio"):
											informacion["denominador"] = dato_plan[index].(map[string]interface{})["dato"]
											if informacion["denominador"] == "Denominador fijo" {
												// informacion["reporteDenominador"] = GetDenominadorFijo(seguimiento, len(indicadores), index)
											}
											continue
										case strings.Contains(nombreDetalle, "tendencia"):
											informacion["tendencia"] = strings.Trim(dato_plan[index].(map[string]interface{})["dato"].(string), " ")
											continue
										case strings.Contains(nombreDetalle, "unidad de medida"):
											informacion["unidad"] = strings.Trim(dato_plan[index].(map[string]interface{})["dato"].(string), " ")
											respuesta["unidad"] = strings.Trim(dato_plan[index].(map[string]interface{})["dato"].(string), " ")
											continue
										}
									}
								}
							}
						}

						if informacion["reporteDenominador"] == 1.0 {
							informacion["reporteDenominador"] = nil
						}

						if informacion["nombre"] != nil && informacion["nombre"] != "" {
							indicadores = append(indicadores, informacion)
							respuestas = append(respuestas, respuesta)
						}

						if informacion["denominador"] == nil {
							informacion["denominador"] = ""
						}

						respuestas = GetRespuestaAnterior(seguimiento, len(indicadores)-1, respuestas, index, trimestre)
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

							if fmt.Sprint(reflect.TypeOf(seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"])) == "int" || fmt.Sprint(reflect.TypeOf(seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"])) == "float64" {
								return seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"].(float64)
							} else {
								aux2, err := strconv.ParseFloat(seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"].(string), 64)
								if err == nil {
									return aux2
								}
							}
						}
						break
					}
				}
			}
		}
	}

	return 1
}

func GetRespuestaAnterior(dataSeg map[string]interface{}, index int, respuestas []map[string]interface{}, indexActividad string, trimestre string) []map[string]interface{} {
	plan_id := dataSeg["plan_id"].(string)
	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var seguimientos []map[string]interface{}
	var periodo []map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+plan_id, &resSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)

		acumuladoNumerador := 0.0
		acumuladoDenominador := 0.0
		indicadorAcumulado := 0.0
		avanceAcumulado := 0.0
		brechaExistente := 0.0
		divisionCero := false
		for _, seguimiento := range seguimientos {

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
					helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
					tri, _ := strconv.Atoi(string(trimestre[1]))
					segTrimestre, _ := strconv.Atoi(string(periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)[1]))

					if (tri - 1) == segTrimestre {
						if seguimiento["dato"] != "{}" {
							dato := make(map[string]interface{})
							datoStr := seguimiento["dato"].(string)
							json.Unmarshal([]byte(datoStr), &dato)
							if dato[indexActividad] == nil {
								respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
								respuestas[index]["avanceAcumulado"] = avanceAcumulado
								respuestas[index]["brechaExistente"] = brechaExistente
								respuestas[index]["divisionCero"] = divisionCero
								continue
							}

							id, segregado := dato[indexActividad].(map[string]interface{})["id"]
							if segregado && id != "" {
								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
									helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
									detalle = ConvertirStringJson(detalle)

									if fmt.Sprintf("%v", detalle["cuantitativo"]) == "map[]" {
										respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
										respuestas[index]["avanceAcumulado"] = avanceAcumulado
										respuestas[index]["brechaExistente"] = brechaExistente
										respuestas[index]["divisionCero"] = divisionCero
										continue
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["indicadorAcumulado"] != nil {
										indicadorAcumulado += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["indicadorAcumulado"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["avanceAcumulado"] != nil {
										avanceAcumulado += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["avanceAcumulado"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["brechaExistente"] != nil {
										brechaExistente += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["brechaExistente"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["divisionCero"] != nil {
										divisionCero = detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["divisionCero"].(bool)
									} else {
										divisionCero = false
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoDenominador"] != nil {
										acumuladoDenominador += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoDenominador"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoNumerador"] != nil {
										acumuladoNumerador += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoNumerador"].(float64)
									}
								}
							} else {
								seguimientoActividad := dato[indexActividad].(map[string]interface{})
								if seguimientoActividad["cuantitativo"] == nil {
									respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
									respuestas[index]["avanceAcumulado"] = avanceAcumulado
									respuestas[index]["brechaExistente"] = brechaExistente
									respuestas[index]["divisionCero"] = divisionCero
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

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["divisionCero"] != nil {
									divisionCero = seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["divisionCero"].(bool)
								} else {
									divisionCero = false
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoDenominador"] != nil {
									acumuladoDenominador += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoDenominador"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoNumerador"] != nil {
									acumuladoNumerador += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoNumerador"].(float64)
								}
							}
						}

						respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
						respuestas[index]["avanceAcumulado"] = avanceAcumulado
						respuestas[index]["brechaExistente"] = brechaExistente
						respuestas[index]["acumuladoNumerador"] = acumuladoNumerador
						respuestas[index]["acumuladoDenominador"] = acumuladoDenominador
						respuestas[index]["divisionCero"] = divisionCero
						break
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
			_, datosUnidos := actividad.(map[string]interface{})["estado"]
			if datosUnidos {
				if actividad.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad en reporte" {
					enReporte = false
				}
			} else {
				var resSeguimientoDetalle map[string]interface{}
				var seguimientoDetalle map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+actividad.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
					helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &seguimientoDetalle)
					dato := make(map[string]interface{})
					json.Unmarshal([]byte(seguimientoDetalle["estado"].(string)), &dato)
					if dato["nombre"] != "Actividad en reporte" {
						enReporte = false
					}
				}
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
	var resSeguimientoDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	var cuantitativo interface{}
	var cualitativo interface{}
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if dato[indexActividad] == nil {
		return false, map[string]interface{}{"error": 1, "motivo": "Actividad sin seguimiento"}
	} else {
		_, datosUnidos := dato[indexActividad].(map[string]interface{})["estado"]
		if datosUnidos {
			estado = dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})
			cuantitativo = dato[indexActividad].(map[string]interface{})["cuantitativo"]
			cualitativo = dato[indexActividad].(map[string]interface{})["cualitativo"]
		} else {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
				helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
				detalle = ConvertirStringJson(detalle)
				estado = detalle["estado"].(map[string]interface{})
				cualitativo = detalle["cualitativo"]
				cuantitativo = detalle["cuantitativo"]
			}

			if fmt.Sprintf("%v", cuantitativo) == "map[]" {
				cuantitativo = nil
			}

			if fmt.Sprintf("%v", cualitativo) == "map[]" {
				cualitativo = nil
			}
		}

		if estado["nombre"] != "Actividad en reporte" {
			return false, map[string]interface{}{"error": 2, "motivo": "El estado de la actividad no es el adecuado"}
		}

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
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})

	planId := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") {

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
						id, segregado := element.(map[string]interface{})["id"]
						if segregado && id != "" {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id.(string), &resSeguimientoDetalle); err == nil {
								helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
								detalle = ConvertirStringJson(detalle)
							}

							for _, actividad := range actividades {
								if reflect.TypeOf(actividad["index"]).String() == "string" {
									if indexActividad == actividad["index"] {
										actividad["estado"] = detalle["estado"]
									}
								} else {
									if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
										actividad["estado"] = detalle["estado"]
									}
								}
							}
						} else {
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
					}

					for _, actividad := range actividades {
						if actividad["estado"] == nil {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								dato[actividad["index"].(string)] = actividad["dato"]
							} else {
								dato[strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)] = actividad["dato"]
							}
						} else if actividad["estado"].(map[string]interface{})["nombre"] != "Actividad reportada" && actividad["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && actividad["estado"].(map[string]interface{})["nombre"] != "Actividad Verificada" {
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
			if (indicador.(map[string]interface{})["observaciones_dependencia"] != "" && indicador.(map[string]interface{})["observaciones_dependencia"] != "Sin observación" && indicador.(map[string]interface{})["observaciones_dependencia"] != nil) ||
				(indicador.(map[string]interface{})["observaciones_planeacion"] != "" && indicador.(map[string]interface{})["observaciones_planeacion"] != "Sin observación" && indicador.(map[string]interface{})["observaciones_planeacion"] != nil) {
				return true
			}
		}
	}

	if seguimiento["cualitativo"] != nil {
		cualitativo = seguimiento["cualitativo"].(map[string]interface{})
		if (cualitativo["observaciones_dependencia"] != "" && cualitativo["observaciones_dependencia"] != "Sin observación" && cualitativo["observaciones_dependencia"] != nil) ||
			(cualitativo["observaciones_planeacion"] != "" && cualitativo["observaciones_planeacion"] != "Sin observación" && cualitativo["observaciones_planeacion"] != nil) {
			return true
		}
	}

	if seguimiento["evidencia"] != nil {
		for _, evidencia := range seguimiento["evidencia"].([]map[string]interface{}) {
			if evidencia["Observacion"] != "" && evidencia["Observacion"] != "Sin observación" && evidencia["Observacion"] != nil {
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
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	observaciones := false
	avaladas := false
	estado := map[string]interface{}{}

	planId := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") {

				actividades := GetActividades(subgrupos[i]["_id"].(string))

				dato_plan_str := seguimiento["dato"].(string)
				json.Unmarshal([]byte(dato_plan_str), &datoPlan)

				for indexActividad, element := range datoPlan {
					id, segregado := element.(map[string]interface{})["id"]

					for _, actividad := range actividades {
						if reflect.TypeOf(actividad["index"]).String() == "string" {
							if indexActividad == actividad["index"] {
								if segregado && id != "" {
									if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
										helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
										detalle = ConvertirStringJson(detalle)
										estado = detalle["estado"].(map[string]interface{})
										if estado["nombre"] != "Actividad avalada" && estado["nombre"] != "Con observaciones" {
											dato[indexActividad] = actividad["dato"]
										}
									}
								} else {
									if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
										dato[indexActividad] = actividad["dato"]
									}
								}
							}
						} else if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
							if segregado && id != "" {
								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
									helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
									detalle = ConvertirStringJson(detalle)
									estado = detalle["estado"].(map[string]interface{})
									if estado["nombre"] != "Actividad avalada" && estado["nombre"] != "Con observaciones" {
										dato[indexActividad] = actividad["dato"]
									}
								}
							} else {
								if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
									dato[indexActividad] = actividad["dato"]
								}
							}
						}
					}

					if segregado && id != "" {
						if estado["nombre"] == "Con observaciones" {
							observaciones = true
						}

						if estado["nombre"] == "Actividad avalada" {
							avaladas = true
						}
					} else {
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
	}

	if fmt.Sprintf("%v", dato) != "map[]" {
		return false, false, map[string]interface{}{"error": 1, "motivo": "Hay actividades sin revisar", "actividades": dato}
	}

	return avaladas, observaciones, nil
}

func SeguimientoVerificable(seguimiento map[string]interface{}) (bool, bool, map[string]interface{}) {
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var datoPlan map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	observaciones := false
	avaladas := false
	estado := map[string]interface{}{}

	planId := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") {

				actividades := GetActividades(subgrupos[i]["_id"].(string))

				dato_plan_str := seguimiento["dato"].(string)
				json.Unmarshal([]byte(dato_plan_str), &datoPlan)

				for indexActividad, element := range datoPlan {
					id, segregado := element.(map[string]interface{})["id"]

					for _, actividad := range actividades {
						if reflect.TypeOf(actividad["index"]).String() == "string" {
							if indexActividad == actividad["index"] {
								if segregado && id != "" {
									if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
										helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
										detalle = ConvertirStringJson(detalle)
										estado = detalle["estado"].(map[string]interface{})
										if estado["nombre"] != "Actividad Verificada" && estado["nombre"] != "Con observaciones" && estado["nombre"] != "Actividad avalada" {
											dato[indexActividad] = actividad["dato"]
										}
									}
								} else {
									if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad Verificada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
										dato[indexActividad] = actividad["dato"]
									}
								}
							}
						} else if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
							if segregado && id != "" {
								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
									helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
									detalle = ConvertirStringJson(detalle)
									estado = detalle["estado"].(map[string]interface{})
									if estado["nombre"] != "Actividad Verificada" && estado["nombre"] != "Con observaciones" && estado["nombre"] != "Actividad avalada" {
										dato[indexActividad] = actividad["dato"]
									}
								}
							} else {
								if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad Verificada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" {
									dato[indexActividad] = actividad["dato"]
								}
							}
						}
					}

					if segregado && id != "" {
						if estado["nombre"] == "Con observaciones" {
							observaciones = true
						}

						if estado["nombre"] == "Actividad Verificada" {
							avaladas = true
						}
					} else {
						if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Con observaciones" {
							observaciones = true
						}

						if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Actividad Verificada" {
							avaladas = true
						}
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

func ConvertirJsonString(diccionario map[string]interface{}) map[string]interface{} {
	dicStrings := map[string]interface{}{}

	for clave, valor := range diccionario {
		jsonBytes, _ := json.Marshal(valor)
		if string(jsonBytes) == "null" {
			dicStrings[clave] = "{}"
		} else if clave == "_id" || clave == "fecha_creacion" {
			dicStrings[clave] = valor
		} else {
			dicStrings[clave] = string(jsonBytes)
		}
	}

	return dicStrings
}

func ConvertirStringJson(diccionario map[string]interface{}) map[string]interface{} {
	dicStrings := map[string]interface{}{}
	for clave, valor := range diccionario {
		if clave == "informacion" || clave == "cualitativo" || clave == "cuantitativo" || clave == "estado" {
			datoJson := make(map[string]interface{})
			json.Unmarshal([]byte(valor.(string)), &datoJson)
			dicStrings[clave] = datoJson
		} else if clave == "evidencia" {
			var datoJson []map[string]interface{}
			json.Unmarshal([]byte(valor.(string)), &datoJson)
			dicStrings[clave] = datoJson
		} else {
			dicStrings[clave] = valor
		}
	}
	return dicStrings
}

func GuardarDetalleSegimiento(detalle map[string]interface{}, actualizar bool) string {
	var res map[string]interface{}
	var id string

	detalle = ConvertirJsonString(detalle)

	if _, existe := detalle["informacion"]; !existe {
		detalle["informacion"] = "{}"
	}
	if _, existe := detalle["cualitativo"]; !existe {
		detalle["cualitativo"] = "{}"
	}
	if _, existe := detalle["cuantitativo"]; !existe {
		detalle["cuantitativo"] = "{}"
	}
	if _, existe := detalle["evidencia"]; !existe {
		detalle["evidencia"] = "[]"
	}

	if actualizar {
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+detalle["_id"].(string), "PUT", &res, detalle); err == nil {
			aux := make(map[string]interface{})
			helpers.LimpiezaRespuestaRefactor(res, &aux)
			id = aux["_id"].(string)
		}
	} else {
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle", "POST", &res, detalle); err == nil {
			aux := make(map[string]interface{})
			helpers.LimpiezaRespuestaRefactor(res, &aux)
			id = aux["_id"].(string)
		}
	}

	return id
}

func CambiarEstadoPlan(plan map[string]interface{}, idEstado string) (map[string]interface{}, error) {
	idPlan := plan["_id"].(string)
	plan["estado_plan_id"] = idEstado
	var res map[string]interface{}
	if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+idPlan, "PUT", &res, plan); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

// encodeBase62 genera el hash a partir de un planId en formato hexadecimal
func EncodeBase62(actividadID string) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	num := new(big.Int)
	num.SetString(actividadID, 16)

	var encoded string

	zero := big.NewInt(0)
	base := big.NewInt(62)

	for num.Cmp(zero) > 0 {
		var remainder big.Int
		num.DivMod(num, base, &remainder)
		encoded = string(charset[remainder.Int64()]) + encoded
	}

	return encoded
}

func DecodeBase62(base62Str string) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	decodedNum := new(big.Int)

	for _, char := range base62Str {
		index := strings.IndexByte(charset, byte(char))
		if index == -1 {
			panic("Invalid character found in base62 string")
		}
		decodedNum.Mul(decodedNum, big.NewInt(62))
		decodedNum.Add(decodedNum, big.NewInt(int64(index)))
	}

	return fmt.Sprintf("%x", decodedNum)
}

func GetSeguimiento(planId string, indexActividad string, trimestreId string) (map[string]interface{}, error) {
	var respuesta map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var resEstado map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var seguimientoActividad map[string]interface{}
	var periodo []map[string]interface{}
	var trimestre string

	id_actividad := EncodeBase62(planId + "" + indexActividad)

	dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestreId, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux[0]

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
			helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
				trimestre = periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)
			}
		}

		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &dato)

		actividad, _ := json.Marshal(GetActividad(seguimiento, indexActividad, trimestre))
		json.Unmarshal([]byte(string(actividad)), &seguimientoActividad)
		seguimientoActividad["_id"] = seguimiento["_id"].(string)
		seguimientoActividad["id_actividad"] = id_actividad

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+seguimiento["estado_seguimiento_id"].(string), &resEstado); err == nil {
			seguimientoActividad["estadoSeguimiento"] = resEstado["Data"].(map[string]interface{})["nombre"].(string)
		}

		return seguimientoActividad, nil
	} else {
		return nil, err
	}
}

func GetEstadoTrimestre(planId string, trimestre string) (respuesta map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			outputError = map[string]interface{}{
				"funcion": "GetPlanesPeriodo",
				"err":     localError["err"],
				"status":  localError["status"],
			}
			panic(outputError)
		}
	}()

	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var planes []map[string]interface{}
	var periodoSeguimiento []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId, &resSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &planes)

		for _, plan := range planes {
			var periodo []map[string]interface{}
			periodoSeguimientoId := plan["periodo_seguimiento_id"].(string)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento?query=_id:"+periodoSeguimientoId, &resPeriodoSeguimiento); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)
				if fmt.Sprintf("%v", periodoSeguimiento[0]) != "map[]" {

					if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento[0]["periodo_id"].(string)+",ParametroId__CodigoAbreviacion:"+trimestre, &resPeriodo); err == nil {
						helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
						plan["periodo_seguimiento_id"] = periodoSeguimiento[0]

						if fmt.Sprintf("%v", periodo[0]) != "map[]" {
							var resEstado map[string]interface{}

							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+plan["estado_seguimiento_id"].(string), &resEstado); err == nil {
								plan["estado_seguimiento_id"] = resEstado["Data"]

								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan["plan_id"].(string), &resEstado); err == nil {
									plan["plan_id"] = resEstado["Data"]

									return plan, nil
								}
							}
						}
					}
				}
			}
		}

	} else {
		panic(map[string]interface{}{
			"err":    err,
			"status": "404",
		})
	}
	return nil, outputError
}

func ObtenerPromedioBrechayEstado(requestBody []byte) (respuesta []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			outputError = map[string]interface{}{
				"funcion": "GetPlanesPeriodo",
				"err":     localError["err"],
				"status":  localError["status"],
			}
			panic(outputError)
		}
	}()

	var body map[string]interface{}
	var trimestres []map[string]interface{}

	if err := json.Unmarshal(requestBody, &body); err == nil {
		nombrePlan := body["nombre"].(string)
		id := body["id"].(string)
		vigencia := body["vigencia"].(string)
		dependencia := body["dependencia"].(string)

		periodos := GetTrimestres(vigencia)
		if len(periodos) == 0 {
			panic(map[string]interface{}{
				"err":    "El plan no tiene definido Trimestres",
				"status": "404",
			})
		}

		for _, periodo := range periodos {
			trimestre := map[string]interface{}{
				"codigo": periodo["ParametroId"].(map[string]interface{})["CodigoAbreviacion"],
				"nombre": periodo["ParametroId"].(map[string]interface{})["Nombre"],
			}
			trimestres = append(trimestres, trimestre)
		}

		for _, tr := range trimestres {
			estado, err := GetEstadoTrimestre(id, tr["codigo"].(string))
			if err != nil {
				panic(map[string]interface{}{
					"err":    err,
					"status": "404",
				})
			}
			tr["estado"] = estado["estado_seguimiento_id"].(map[string]interface{})["nombre"]
		}

		unidades, errUnd := comunhelper.GetUnidadesPorPlanYVigencia(nombrePlan, vigencia)
		if errUnd != nil {
			panic(map[string]interface{}{
				"err":    errUnd,
				"status": "404",
			})
		}
		if len(unidades) != 0 {
			planesPeriodo, errPer := comunhelper.GetPlanesPeriodo(dependencia, vigencia)
			if errPer != nil {
				panic(map[string]interface{}{
					"err":    errPer,
					"status": "404",
				})
			}

			var planEspecifico map[string]interface{}
			for _, item := range planesPeriodo {
				if item["plan"] == nombrePlan {
					planEspecifico = item
				}
			}
			if planEspecifico != nil {
				pers := planEspecifico["periodos"].([]map[string]interface{})
				ultimoPeriodo := pers[len(pers)-1]
				ultimoPeriodoID := ultimoPeriodo["id"].(string)

				periodosPlan := comunhelper.GetPeriodosPlan(vigencia, id)
				if len(periodosPlan) == 0 {
					panic(map[string]interface{}{
						"err":    "El plan no posee trimestres",
						"status": "404",
					})
				} else {
					var evaluacion []map[string]interface{}
					for posicionTri, tri := range periodosPlan {
						if tri["_id"] == ultimoPeriodoID {
							evaluacion = comunhelper.GetEvaluacion(planEspecifico["id"].(string), periodosPlan, posicionTri)
							break
						}
					}
					if evaluacion == nil {
						panic(map[string]interface{}{
							"err":    "El plan no posee evaluacion",
							"status": "404",
						})
					}
					var brechasT1 []float64
					var brechasT2 []float64
					var brechasT3 []float64
					var brechasT4 []float64
					for _, eval := range evaluacion {
						if eval["trimestre1"] != nil && eval["trimestre1"].(map[string]interface{})["brecha"] != "" && len(eval["trimestre1"].(map[string]interface{})) > 0 {
							brechasT1 = append(brechasT1, eval["trimestre1"].(map[string]interface{})["brecha"].(float64))
						}
						if eval["trimestre2"] != nil && eval["trimestre2"].(map[string]interface{})["brecha"] != "" && len(eval["trimestre2"].(map[string]interface{})) > 0 {
							brechasT2 = append(brechasT2, eval["trimestre2"].(map[string]interface{})["brecha"].(float64))
						}
						if eval["trimestre3"] != nil && eval["trimestre3"].(map[string]interface{})["brecha"] != "" && len(eval["trimestre3"].(map[string]interface{})) > 0 {
							brechasT3 = append(brechasT3, eval["trimestre3"].(map[string]interface{})["brecha"].(float64))
						}
						if eval["trimestre4"] != nil && eval["trimestre4"].(map[string]interface{})["brecha"] != "" && len(eval["trimestre4"].(map[string]interface{})) > 0 {
							brechasT4 = append(brechasT4, eval["trimestre4"].(map[string]interface{})["brecha"].(float64))
						}
					}

					for _, tr := range trimestres {
						var suma float64 = 0
						var prod float64 = 0
						if tr["codigo"] == "T1" && len(brechasT1) != 0 {
							for _, numero := range brechasT1 {
								suma += numero
							}
							if len(brechasT1) > 0 {
								prod = ((suma / float64(len(brechasT1))))
							}

							prodFormatted := fmt.Sprintf("%.2f", prod*100)
							tr["promedioBrechas"] = prodFormatted
						} else if tr["codigo"] == "T2" && len(brechasT2) != 0 {
							for _, numero := range brechasT2 {
								suma += numero
							}
							if len(brechasT2) > 0 {
								prod = ((suma / float64(len(brechasT2))))
							}

							prodFormatted := fmt.Sprintf("%.2f", prod*100)
							tr["promedioBrechas"] = prodFormatted
						} else if tr["codigo"] == "T3" && len(brechasT3) != 0 {
							for _, numero := range brechasT3 {
								suma += numero
							}
							if len(brechasT3) > 0 {
								prod = ((suma / float64(len(brechasT3))))
							}

							prodFormatted := fmt.Sprintf("%.2f", prod*100)
							tr["promedioBrechas"] = prodFormatted
						} else if tr["codigo"] == "T4" && len(brechasT4) != 0 {
							for _, numero := range brechasT4 {
								suma += numero
							}
							if len(brechasT4) > 0 {
								prod = ((suma / float64(len(brechasT4))))
							}

							prodFormatted := fmt.Sprintf("%.2f", prod*100)
							tr["promedioBrechas"] = prodFormatted
						} else {
							tr["promedioBrechas"] = 0
						}
					}
				}

			} else {
				for _, tr := range trimestres {
					tr["promedioBrechas"] = 0
				}
			}
		} else {
			for _, tr := range trimestres {
				tr["promedioBrechas"] = 0
			}
		}
	} else {
		panic(map[string]interface{}{
			"err":    err,
			"status": "404",
		})
	}

	return trimestres, outputError
}
