package inversionhelper

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

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
