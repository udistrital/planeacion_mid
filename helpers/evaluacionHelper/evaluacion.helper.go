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

func GetEvaluacion(planId string, periodos []map[string]interface{}, trimestre int) []map[string]interface{} {
	var resSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var evaluacion []map[string]interface{}
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
			actividad := act.(map[string]interface{})
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
							"numero":     actividad["informacion"].(map[string]interface{})["index"],
							"ponderado":  actividad["informacion"].(map[string]interface{})["ponderacion"],
							"periodo":    actividad["informacion"].(map[string]interface{})["periodo"],
							"actividad":  actividad["informacion"].(map[string]interface{})["descripcion"],
							"trimestre1": make(map[string]interface{}),
							"trimestre2": make(map[string]interface{}),
							"trimestre3": make(map[string]interface{}),
							"trimestre4": make(map[string]interface{}),
						}
						evaluacionAct["indicador"] = resIndicador["indicador"]
						evaluacionAct["formula"] = resIndicador["formula"]
						evaluacionAct["meta"] = resIndicador["metaA"]
						evaluacionAct[trimestreNom] = map[string]interface{}{
							"numerador":   resIndicador["numerador"],
							"denominador": resIndicador["denominador"],
							"periodo":     resIndicador["periodo"],
							"acumulado":   resIndicador["acumulado"],
							"meta":        resIndicador["meta"],
						}

						evaluacion = append(evaluacion, evaluacionAct)
					} else {
						evaluacion[indice][trimestreNom] = map[string]interface{}{
							"numerador":   resIndicador["numerador"],
							"denominador": resIndicador["denominador"],
							"periodo":     resIndicador["periodo"],
							"acumulado":   resIndicador["acumulado"],
							"meta":        resIndicador["meta"],
						}
					}
				}
			}
		}

		helpers.SortSlice(&evaluacion, "numero")
		cont := 0
		ant := 0
		size := len(evaluacion) - 1
		for index, eval := range evaluacion {
			if index == 0 {
				cont++
				continue
			}

			if eval["numero"] == evaluacion[index-1]["numero"] && size != index {
				cont++
			} else {
				if size == index {
					cont++
				}
				sum1 := 0.0
				sum2 := 0.0
				sum3 := 0.0
				sum4 := 0.0

				for i := ant; i < cont+ant; i++ {
					if fmt.Sprintf("%v", evaluacion[i]["trimestre1"]) != "map[]" {
						sum1 = sum1 + evaluacion[i]["trimestre1"].(map[string]interface{})["meta"].(float64)
					}
					if fmt.Sprintf("%v", evaluacion[i]["trimestre2"]) != "map[]" {
						sum2 = sum2 + evaluacion[i]["trimestre2"].(map[string]interface{})["meta"].(float64)
					}
					if fmt.Sprintf("%v", evaluacion[i]["trimestre3"]) != "map[]" {
						sum3 = sum3 + evaluacion[i]["trimestre3"].(map[string]interface{})["meta"].(float64)
					}
					if fmt.Sprintf("%v", evaluacion[i]["trimestre4"]) != "map[]" {
						sum4 = sum4 + evaluacion[i]["trimestre4"].(map[string]interface{})["meta"].(float64)
					}
				}

				cumplActividad1 := math.Floor((sum1/float64(cont))*100) / 100
				cumplActividad2 := math.Floor((sum2/float64(cont))*100) / 100
				cumplActividad3 := math.Floor((sum3/float64(cont))*100) / 100
				cumplActividad4 := math.Floor((sum4/float64(cont))*100) / 100

				if cumplActividad1 > 100 {
					cumplActividad1 = 100
				}

				if cumplActividad2 > 100 {
					cumplActividad2 = 100
				}

				if cumplActividad3 > 100 {
					cumplActividad3 = 100
				}

				if cumplActividad4 > 100 {
					cumplActividad4 = 100
				}

				for i := ant; i < cont+ant; i++ {
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

				cont = 1
				ant = index
			}
		}

		return evaluacion
	}
	return nil
}

func GetEvaluacionTrimestre(planId string, periodoId string, actividadId string) []map[string]interface{} {
	var resSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var evaluacion []map[string]interface{}
	actividades := make(map[string]interface{})

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
		indicadores := actividades[actividadId].(map[string]interface{})["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})
		resultados := actividades[actividadId].(map[string]interface{})["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})
		for i := 0; i < len(indicadores); i++ {

			var meta float64
			if reflect.TypeOf(indicadores[i].(map[string]interface{})["meta"]).String() == "string" {
				meta, _ = strconv.ParseFloat(indicadores[i].(map[string]interface{})["meta"].(string), 64)
			} else {
				meta = indicadores[i].(map[string]interface{})["meta"].(float64)
			}

			var brecha float64
			if reflect.TypeOf(resultados[i].(map[string]interface{})["brechaExistente"]).String() == "string" {
				brecha, _ = strconv.ParseFloat(resultados[i].(map[string]interface{})["brechaExistente"].(string), 64)
			} else {
				brecha = resultados[i].(map[string]interface{})["brechaExistente"].(float64)
			}

			evaluacion = append(evaluacion, map[string]interface{}{
				"indicador":   indicadores[i].(map[string]interface{})["nombre"],
				"formula":     indicadores[i].(map[string]interface{})["formula"],
				"metaA":       indicadores[i].(map[string]interface{})["meta"],
				"numerador":   indicadores[i].(map[string]interface{})["reporteNumerador"],
				"denominador": indicadores[i].(map[string]interface{})["reporteDenominador"],
				"periodo":     resultados[i].(map[string]interface{})["indicador"],
				"acumulado":   resultados[i].(map[string]interface{})["indicadorAcumulado"],
				"meta":        math.Floor((meta-brecha)*100) / 100,
				"actividad":   0})
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
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+strconv.Itoa(trimestreId), &resPeriodo); err == nil {
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
