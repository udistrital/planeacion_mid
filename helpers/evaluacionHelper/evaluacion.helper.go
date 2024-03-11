package evaluacionhelper

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	seguimientohelper "github.com/udistrital/planeacion_mid/helpers/seguimientoHelper"
	"github.com/udistrital/utils_oas/request"
)

const (
	ABREVIACION_AVALADO_PARA_SEGUIMIENTO string = "AV"
	ABREVIACION_SEGUIMIENTO_PLAN_ACCION  string = "S_SP"
)

func GetEvaluacionTrimestre(planId string, periodoId string, actividadId string) []map[string]interface{} {
	var resSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var evaluacion []map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	actividades := make(map[string]interface{})
	detalle := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=estado_seguimiento_id:622ba49216511e93a95c326d,plan_id:`+planId+`,periodo_seguimiento_id:`+periodoId, &resSeguimiento); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &aux)
		if fmt.Sprintf("%v", aux) == "[]" {
			return nil
		}

		seguimiento = aux[0]

		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &actividades)

		if actividades[actividadId] == nil {
			return nil
		}

		var indicadores []interface{}
		var resultados []interface{}
		id, segregado := actividades[actividadId].(map[string]interface{})["id"]

		if segregado && id != "" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id.(string), &resSeguimientoDetalle); err == nil {
				helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
				detalle = seguimientohelper.ConvertirStringJson(detalle)
				if fmt.Sprintf("%v", detalle["cuantitativo"]) != "map[]" {
					indicadores = detalle["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})
					resultados = detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})
				} else {
					indicadores = []interface{}{}
					resultados = []interface{}{}
				}
			}
		} else {
			indicadores = actividades[actividadId].(map[string]interface{})["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})
			resultados = actividades[actividadId].(map[string]interface{})["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})
		}

		for i := 0; i < len(indicadores); i++ {
			var metaA float64
			if indicadores[i].(map[string]interface{})["meta"] == nil {
				metaA = 0
			} else {
				if reflect.TypeOf(indicadores[i].(map[string]interface{})["meta"]).String() == "string" {
					metaA, _ = strconv.ParseFloat(indicadores[i].(map[string]interface{})["meta"].(string), 64)
				} else {
					metaA = indicadores[i].(map[string]interface{})["meta"].(float64)
				}
			}

			evaluacion = append(evaluacion, map[string]interface{}{
				"indicador":            indicadores[i].(map[string]interface{})["nombre"],
				"formula":              indicadores[i].(map[string]interface{})["formula"],
				"metaA":                metaA,
				"unidad":               indicadores[i].(map[string]interface{})["unidad"],
				"numerador":            indicadores[i].(map[string]interface{})["reporteNumerador"],
				"denominador":          indicadores[i].(map[string]interface{})["reporteDenominador"],
				"periodo":              resultados[i].(map[string]interface{})["indicador"],
				"acumulado":            resultados[i].(map[string]interface{})["indicadorAcumulado"],
				"meta":                 resultados[i].(map[string]interface{})["avanceAcumulado"],
				"numeradorAcumulado":   resultados[i].(map[string]interface{})["acumuladoNumerador"],
				"denominadorAcumulado": resultados[i].(map[string]interface{})["acumuladoDenominador"],
				"brecha":               resultados[i].(map[string]interface{})["brechaExistente"],
				"actividad":            0})
		}
		return evaluacion
	}
	return nil
}

func GetEvaluacion(planId string, periodos []map[string]interface{}, trimestre int) []map[string]interface{} {
	var resSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var evaluacion []map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	actividades := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=estado_seguimiento_id:622ba49216511e93a95c326d,plan_id:`+planId+`,periodo_seguimiento_id:`+periodos[trimestre]["_id"].(string), &resSeguimiento); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &aux)
		if fmt.Sprintf("%v", aux) == "[]" {
			return nil
		}

		seguimiento = aux[0]
		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &actividades)

		for actividadId, act := range actividades {
			id, segregado := actividades[actividadId].(map[string]interface{})["id"].(string)
			var actividad map[string]interface{}

			if segregado && id != "" {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id, &resSeguimientoDetalle); err == nil {
					helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
					actividad = seguimientohelper.ConvertirStringJson(detalle)
				}
			} else {
				actividad = act.(map[string]interface{})
			}
			for indexPeriodo, periodo := range periodos {
				if indexPeriodo > trimestre {
					break
				}
				resIndicadores := GetEvaluacionTrimestre(planId, periodo["_id"].(string), actividadId)
				for _, resIndicador := range resIndicadores {

					indice := -1
					for index, eval := range evaluacion {
						if eval["numero"] == actividad["informacion"].(map[string]interface{})["index"] && eval["indicador"] == resIndicador["indicador"] {
							indice = index
							break
						}
					}

					var trimestreNom string
					if indexPeriodo == 0 {
						trimestreNom = "trimestre1"
					} else if indexPeriodo == 1 {
						trimestreNom = "trimestre2"
					} else if indexPeriodo == 2 {
						trimestreNom = "trimestre3"
					} else if indexPeriodo == 3 {
						trimestreNom = "trimestre4"
					}

					if indice == -1 {
						evaluacionAct := map[string]interface{}{
							"actividad":  actividad["informacion"].(map[string]interface{})["descripcion"],
							"numero":     actividad["informacion"].(map[string]interface{})["index"],
							"periodo":    actividad["informacion"].(map[string]interface{})["periodo"],
							"ponderado":  actividad["informacion"].(map[string]interface{})["ponderacion"],
							"trimestre1": make(map[string]interface{}),
							"trimestre2": make(map[string]interface{}),
							"trimestre3": make(map[string]interface{}),
							"trimestre4": make(map[string]interface{}),
						}
						evaluacionAct["indicador"] = resIndicador["indicador"]
						evaluacionAct["unidad"] = resIndicador["unidad"]
						evaluacionAct["formula"] = resIndicador["formula"]
						evaluacionAct["meta"] = resIndicador["metaA"].(float64)
						evaluacionAct[trimestreNom] = map[string]interface{}{
							"acumulado":            resIndicador["acumulado"],
							"denominador":          resIndicador["denominador"],
							"meta":                 resIndicador["meta"],
							"numerador":            resIndicador["numerador"],
							"periodo":              resIndicador["periodo"],
							"numeradorAcumulado":   resIndicador["numeradorAcumulado"],
							"denominadorAcumulado": resIndicador["denominadorAcumulado"],
							"brecha":               resIndicador["brecha"],
						}

						evaluacion = append(evaluacion, evaluacionAct)
					} else {
						evaluacion[indice][trimestreNom] = map[string]interface{}{
							"acumulado":            resIndicador["acumulado"],
							"denominador":          resIndicador["denominador"],
							"meta":                 resIndicador["meta"],
							"numerador":            resIndicador["numerador"],
							"periodo":              resIndicador["periodo"],
							"numeradorAcumulado":   resIndicador["numeradorAcumulado"],
							"denominadorAcumulado": resIndicador["denominadorAcumulado"],
							"brecha":               resIndicador["brecha"],
						}
					}
				}
			}
		}

		helpers.SortSlice(&evaluacion, "numero")
		agrupacion_actividades := make(map[string][]int)
		for i, eval := range evaluacion {
			if _, ok := agrupacion_actividades[eval["numero"].(string)]; !ok {
				agrupacion_actividades[eval["numero"].(string)] = []int{}
			}
			agrupacion_actividades[eval["numero"].(string)] = append(agrupacion_actividades[eval["numero"].(string)], i)
		}

		for _, idxs := range agrupacion_actividades {
			sum1 := 0.0
			sum2 := 0.0
			sum3 := 0.0
			sum4 := 0.0

			for _, i := range idxs {
				if fmt.Sprintf("%v", evaluacion[i]["trimestre1"]) != "map[]" {
					if evaluacion[i]["trimestre1"].(map[string]interface{})["meta"].(float64) > 1 {
						sum1 = sum1 + 1.0
					} else {
						sum1 = sum1 + evaluacion[i]["trimestre1"].(map[string]interface{})["meta"].(float64)
					}
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre2"]) != "map[]" {
					if evaluacion[i]["trimestre2"].(map[string]interface{})["meta"].(float64) > 1 {
						sum2 = sum2 + 1.0
					} else {
						sum2 = sum2 + evaluacion[i]["trimestre2"].(map[string]interface{})["meta"].(float64)
					}
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre3"]) != "map[]" {
					if evaluacion[i]["trimestre3"].(map[string]interface{})["meta"].(float64) > 1 {
						sum3 = sum3 + 1.0
					} else {
						sum3 = sum3 + evaluacion[i]["trimestre3"].(map[string]interface{})["meta"].(float64)
					}
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre4"]) != "map[]" {
					if evaluacion[i]["trimestre4"].(map[string]interface{})["meta"].(float64) > 1 {
						sum4 = sum4 + 1.0
					} else {
						sum4 = sum4 + evaluacion[i]["trimestre4"].(map[string]interface{})["meta"].(float64)
					}
				}
			}

			cont := len(idxs)
			cumplActividad1 := math.Floor((sum1/float64(cont))*1000) / 1000
			cumplActividad2 := math.Floor((sum2/float64(cont))*1000) / 1000
			cumplActividad3 := math.Floor((sum3/float64(cont))*1000) / 1000
			cumplActividad4 := math.Floor((sum4/float64(cont))*1000) / 1000

			if cumplActividad1 > 1 {
				cumplActividad1 = 1
			}

			if cumplActividad2 > 1 {
				cumplActividad2 = 1
			}

			if cumplActividad3 > 1 {
				cumplActividad3 = 1
			}

			if cumplActividad4 > 1 {
				cumplActividad4 = 1
			}

			for _, i := range idxs {
				if fmt.Sprintf("%v", evaluacion[i]["trimestre1"]) != "map[]" {
					evaluacion[i]["trimestre1"].(map[string]interface{})["actividad"] = cumplActividad1
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre2"]) != "map[]" {
					evaluacion[i]["trimestre2"].(map[string]interface{})["actividad"] = cumplActividad2
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre3"]) != "map[]" {
					evaluacion[i]["trimestre3"].(map[string]interface{})["actividad"] = cumplActividad3
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre4"]) != "map[]" {
					evaluacion[i]["trimestre4"].(map[string]interface{})["actividad"] = cumplActividad4
				}
			}
		}

		return evaluacion
	}
	return nil
}

func GetPeriodos(vigencia string) []map[string]interface{} {
	var periodos []map[string]interface{}
	var resPeriodo map[string]interface{}
	var wg sync.WaitGroup
	trimestres := seguimientohelper.GetTrimestres(vigencia)
	periodosMutex := sync.Mutex{}

	for _, trimestre := range trimestres {
		wg.Add(1)
		if fmt.Sprintf("%v", trimestre) == "map[]" {
			wg.Done()
			continue
		}
		go func(trimestreId int, wg *sync.WaitGroup, periodos *[]map[string]interface{}) {
			periodosMutex.Lock()
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?fields=_id,periodo_id&query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+strconv.Itoa(trimestreId), &resPeriodo); err == nil {
				var periodo []map[string]interface{}
				helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
				(*periodos) = append((*periodos), periodo...)
			}
			periodosMutex.Unlock()
			wg.Done()
		}(int(trimestre["Id"].(float64)), &wg, &periodos)
	}

	wg.Wait()

	helpers.SortSlice(&periodos, "periodo_id")
	return periodos
}

func GetPlanesParaEvaluar() (planes []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			outputError = map[string]interface{}{
				"funcion": "GetPlanesParaEvaluar",
				"err":     localError["err"],
				"status":  localError["status"],
			}
			panic(outputError)
		}
	}()

	var respuestaEstado map[string]interface{}
	var respuestaTipoSeguimiento map[string]interface{}
	var respuestaSeguimiento map[string]interface{}

	var estadoSeguimiento []map[string]interface{}
	var tipoSeguimiento []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_AVALADO_PARA_SEGUIMIENTO, &respuestaEstado); err != nil {
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaEstado, &estadoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_SEGUIMIENTO_PLAN_ACCION, &respuestaTipoSeguimiento); err != nil {
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaTipoSeguimiento, &tipoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:`+tipoSeguimiento[0]["_id"].(string)+`,estado_seguimiento_id:`+estadoSeguimiento[0]["_id"].(string), &respuestaSeguimiento); err == nil {
		var seguimientos []map[string]interface{}
		helpers.LimpiezaRespuestaRefactor(respuestaSeguimiento, &seguimientos)
		for _, seguimiento := range seguimientos {
			// Esta en los planes que ya se trajeron?
			var respuestaPlan map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan/`+seguimiento["plan_id"].(string), &respuestaPlan); err == nil {
				var plan map[string]interface{}
				existeNombrePlan := false
				helpers.LimpiezaRespuestaRefactor(respuestaPlan, &plan)
				for _, nombre := range planes {
					if nombre == plan["nombre"].(string) {
						existeNombrePlan = true
					}
				}
				if !existeNombrePlan {
					planes = append(planes, plan["nombre"].(string))
				}
			} else {
				panic(err)
			}
		}
	} else {
		panic(err)
	}
	return planes, outputError
}

func GetUnidadesPorPlanYVigencia(nombrePlan string, vigencia string) (unidades []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			outputError = map[string]interface{}{
				"funcion": "GetUnidadesPorPlanYVigencia",
				"err":     localError["err"],
				"status":  localError["status"],
			}
			panic(outputError)
		}
	}()

	var respuestaEstado map[string]interface{}
	var respuestaTipoSeguimiento map[string]interface{}
	var respuestaSeguimiento map[string]interface{}
	var respuestaTipoDependencia []map[string]interface{}

	var estadoSeguimiento []map[string]interface{}
	var tipoSeguimiento []map[string]interface{}
	idsUnidades := make([]string, 0)
	unidades = make([]map[string]interface{}, 0)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_AVALADO_PARA_SEGUIMIENTO, &respuestaEstado); err != nil {
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaEstado, &estadoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_SEGUIMIENTO_PLAN_ACCION, &respuestaTipoSeguimiento); err != nil {
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaTipoSeguimiento, &tipoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:`+tipoSeguimiento[0]["_id"].(string)+`,estado_seguimiento_id:`+estadoSeguimiento[0]["_id"].(string), &respuestaSeguimiento); err == nil {
		var seguimientos []map[string]interface{}
		helpers.LimpiezaRespuestaRefactor(respuestaSeguimiento, &seguimientos)
		for _, seguimiento := range seguimientos {
			// Esta en los planes que ya se trajeron?
			var respuestaPlan map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan/`+seguimiento["plan_id"].(string), &respuestaPlan); err == nil {
				var plan map[string]interface{}
				helpers.LimpiezaRespuestaRefactor(respuestaPlan, &plan)

				if plan["nombre"] == nombrePlan && vigencia == plan["vigencia"] {
					existeIdUnidad := false
					for _, idUnidad := range idsUnidades {
						fmt.Println(plan["dependencia_id"])
						if idUnidad == plan["dependencia_id"].(string) {
							existeIdUnidad = true
						}
					}
					if !existeIdUnidad {
						idsUnidades = append(idsUnidades, plan["dependencia_id"].(string))
					}
				}

			} else {
				panic(err)
			}
		}
		if len(idsUnidades) > 0 {
			valor := 0
			for _, idUnidad := range idsUnidades {
				valor++
				if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId__Id:"+idUnidad, &respuestaTipoDependencia); err == nil {
					aux := respuestaTipoDependencia[0]["DependenciaId"].(map[string]interface{})
					delete(aux, "DependenciaTipoDependencia")
					aux["TipoDependencia"] = respuestaTipoDependencia[0]["TipoDependenciaId"]
					unidades = append(unidades, aux)
					respuestaTipoDependencia = nil
				} else {
					panic(err)
				}
			}
		}
	} else {
		panic(err)
	}
	return unidades, outputError
}
