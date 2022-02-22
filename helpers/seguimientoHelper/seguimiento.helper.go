package seguimientohelper

import (
	"encoding/json"

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

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId:641", &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &trimestre)
		trimestres = append(trimestres, trimestre...)

		trimestre = nil
		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId:642", &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &trimestre)
			trimestres = append(trimestres, trimestre...)

			trimestre = nil
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId:643", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &trimestre)
				trimestres = append(trimestres, trimestre...)

				trimestre = nil
				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId:644", &res); err == nil {
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

		}
	} else {
		panic(map[string]interface{}{"Code": "400", "Body": err, "Type": "error"})

	}
	return actividades
}

func GetDataSubgrupos(subgrupos []map[string]interface{}, index string) map[string]interface{} {
	var data map[string]interface{}
	auxSubgrupo := make(map[string]interface{})

	for i := range subgrupos {
		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "lineamiento") {
			aux := getSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
			auxSubgrupo["lineamiento"] = aux["dato"]

		}

		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "meta") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "estratégica") {
			aux := getSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
			auxSubgrupo["meta_estrategica"] = aux["dato"]

		}

		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "estrategia") {
			aux := getSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
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
					aux := getSubgrupoDetalle(hijos[j]["_id"].(string), index)
					auxSubgrupo["indicador"] = aux["dato"]
					indicadores = append(indicadores, aux)
				}
				if strings.Contains(strings.ToLower(hijos[j]["nombre"].(string)), "meta") {
					aux := getSubgrupoDetalle(hijos[j]["_id"].(string), index)
					auxSubgrupo["meta"] = aux["dato"]
					metas = append(metas, aux)

				}
			}
			auxSubgrupo["indicador"] = indicadores
			auxSubgrupo["meta"] = metas

		}

		if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "tarea") {
			aux := getSubgrupoDetalle(subgrupos[i]["_id"].(string), index)
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

func getSubgrupoDetalle(subgrupo_id string, index string) map[string]interface{} {
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
