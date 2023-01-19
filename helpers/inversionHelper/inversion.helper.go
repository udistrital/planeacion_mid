package inversionhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func RegistrarSoportes(idProyect string, registroProyecto map[string]interface{}) (map[string]interface{}, error) {
	var resPostSoportes map[string]interface{}
	planSoportes := make(map[string]interface{})
	clienteHttp := &http.Client{}
	planSoportes["activo"] = true
	planSoportes["padre"] = idProyect
	planSoportes["nombre"] = "soportes"
	planSoportes["descripcion"] = registroProyecto["codigo_proyecto"]
	fmt.Println(planSoportes, "soportes")
	aux, err := json.Marshal(planSoportes)
	if err != nil {
		return nil, err
	}
	peticion, err := http.NewRequest("POST", "http://"+beego.AppConfig.String("PlanesService")+"/subgrupo", bytes.NewBuffer(aux))
	if err != nil {
		return nil, err
	}
	peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
	respuesta, err := clienteHttp.Do(peticion)
	if err != nil {
		return nil, err
	}

	defer respuesta.Body.Close()

	cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cuerpoRespuesta, &resPostSoportes)
	return resPostSoportes, err
}

func RegistrarSoporteDetalle(idSoporte string, registroProyecto map[string]interface{}) (map[string]interface{}, error) {
	var resPostSoporteDetalle map[string]interface{}
	planSoportesDetalle := make(map[string]interface{})
	clienteHttp := &http.Client{}
	soportes, err := json.Marshal(registroProyecto["soportes"])
	if err != nil {
		return nil, err
	}
	planSoportesDetalle["activo"] = true
	planSoportesDetalle["subgrupo_id"] = idSoporte
	planSoportesDetalle["nombre"] = "soportes"
	planSoportesDetalle["descripcion"] = registroProyecto["codigo_proyecto"]
	planSoportesDetalle["dato"] = string(soportes)
	fmt.Println(planSoportesDetalle, "soportes")
	aux, err := json.Marshal(planSoportesDetalle)
	if err != nil {
		return nil, err
	}
	peticion, err := http.NewRequest("POST", "http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle", bytes.NewBuffer(aux))
	if err != nil {
		return nil, err
	}
	peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
	respuesta, err := clienteHttp.Do(peticion)
	if err != nil {
		return nil, err
	}

	defer respuesta.Body.Close()

	cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cuerpoRespuesta, &resPostSoporteDetalle)
	return resPostSoporteDetalle, err
}

func RegistrarFuentesApropiacion(idProyect string, registroProyecto map[string]interface{}) (map[string]interface{}, error) {
	var resPostFuentes map[string]interface{}
	planFuentes := make(map[string]interface{})
	clienteHttp := &http.Client{}
	planFuentes["activo"] = true
	planFuentes["padre"] = idProyect
	planFuentes["nombre"] = "fuentes apropiacion"
	planFuentes["descripcion"] = registroProyecto["codigo_proyecto"]
	fmt.Println(planFuentes, "fuentes")
	aux, err := json.Marshal(planFuentes)
	if err != nil {
		return nil, err
	}
	peticion, err := http.NewRequest("POST", "http://"+beego.AppConfig.String("PlanesService")+"/subgrupo", bytes.NewBuffer(aux))
	if err != nil {
		return nil, err
	}
	peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
	respuesta, err := clienteHttp.Do(peticion)
	if err != nil {
		return nil, err
	}

	defer respuesta.Body.Close()

	cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cuerpoRespuesta, &resPostFuentes)
	return resPostFuentes, err
}

func RegistrarFuentesDetalle(idFuentes string, registroProyecto map[string]interface{}) (map[string]interface{}, error) {
	var resPostFuentesDetalle map[string]interface{}
	planFuentesDetalle := make(map[string]interface{})
	fuentes, err := json.Marshal(registroProyecto["fuentes"])
	if err != nil {
		return nil, err
	}
	clienteHttp := &http.Client{}
	planFuentesDetalle["activo"] = true
	planFuentesDetalle["subgrupo_id"] = idFuentes
	planFuentesDetalle["nombre"] = "Fuentes"
	planFuentesDetalle["descripcion"] = registroProyecto["codigo_proyecto"]
	planFuentesDetalle["dato"] = string(fuentes)
	aux, err := json.Marshal(planFuentesDetalle)
	if err != nil {
		return nil, err
	}
	peticion, err := http.NewRequest("POST", "http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle", bytes.NewBuffer(aux))
	if err != nil {
		return nil, err
	}
	peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
	respuesta, err := clienteHttp.Do(peticion)
	if err != nil {
		return nil, err
	}

	defer respuesta.Body.Close()

	cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cuerpoRespuesta, &resPostFuentesDetalle)
	return resPostFuentesDetalle, err
}

func RegistrarMetas(idProyect string, registroProyecto map[string]interface{}) (map[string]interface{}, error) {
	var resPostMetas map[string]interface{}
	planMetas := make(map[string]interface{})
	clienteHttp := &http.Client{}
	planMetas["activo"] = true
	planMetas["padre"] = idProyect
	planMetas["nombre"] = "metas asociadas al proyecto de inversion"
	planMetas["descripcion"] = registroProyecto["codigo_proyecto"]
	aux, err := json.Marshal(planMetas)
	if err != nil {
		return nil, err
	}
	peticion, err := http.NewRequest("POST", "http://"+beego.AppConfig.String("PlanesService")+"/subgrupo", bytes.NewBuffer(aux))
	if err != nil {
		return nil, err
	}
	peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
	respuesta, err := clienteHttp.Do(peticion)
	if err != nil {
		return nil, err
	}

	defer respuesta.Body.Close()

	cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cuerpoRespuesta, &resPostMetas)
	return resPostMetas, err
}

func RegistrarMetasDetalle(idMetas string, registroProyecto map[string]interface{}) (map[string]interface{}, error) {
	var resPostMetasDetalle map[string]interface{}
	planMetasDetalle := make(map[string]interface{})
	metas, err := json.Marshal(registroProyecto["metas"])
	if err != nil {
		return nil, err
	}
	//url := "http://localhost:8070"
	clienteHttp := &http.Client{}
	planMetasDetalle["activo"] = true
	planMetasDetalle["subgrupo_id"] = idMetas
	planMetasDetalle["nombre"] = "Metas proyectos de inversion"
	planMetasDetalle["descripcion"] = registroProyecto["codigo_proyecto"]
	planMetasDetalle["dato"] = string(metas) //registroProyecto["metas"]
	aux, err := json.Marshal(planMetasDetalle)
	if err != nil {
		return nil, err
	}
	peticion, err := http.NewRequest("POST", "http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle", bytes.NewBuffer(aux))
	if err != nil {
		return nil, err
	}
	peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
	respuesta, err := clienteHttp.Do(peticion)
	if err != nil {
		return nil, err
	}

	defer respuesta.Body.Close()

	cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(cuerpoRespuesta, &resPostMetasDetalle)

	return resPostMetasDetalle, err
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
func GetDataProyects(id string) map[string]interface{} {
	var res map[string]interface{}
	getProyect := make(map[string]interface{})
	var infoProyect map[string]interface{}
	var subgruposData map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &infoProyect)
		getProyect["nombre_proyecto"] = infoProyect["nombre"]
		getProyect["codigo_proyecto"] = infoProyect["descripcion"]
		getProyect["fecha_creacion"] = infoProyect["fecha_creacion"]
		getProyect["id"] = infoProyect["_id"]
		padreId := infoProyect["_id"].(string)

		var infoSubgrupos []map[string]interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+padreId, &subgruposData); err == nil {
			helpers.LimpiezaRespuestaRefactor(subgruposData, &infoSubgrupos)
			//getProyect["fecha_creacion"]
			for i := range infoSubgrupos {
				var subgrupoDetalle map[string]interface{}
				var detalleSubgrupos []map[string]interface{}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "soporte") {
					getProyect["subgrupo_id_soportes"] = infoSubgrupos[i]["_id"]
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
						getProyect["soportes"] = detalleSubgrupos[0]["dato"]
						getProyect["id_detalle_soportes"] = detalleSubgrupos[0]["_id"]
					}
				}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "metas") {
					getProyect["subgrupo_id_metas"] = infoSubgrupos[i]["_id"]
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
						getProyect["metas"] = detalleSubgrupos[0]["dato"]
						getProyect["id_detalle_metas"] = detalleSubgrupos[0]["_id"]

					}
				}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "fuentes") {
					getProyect["subgrupo_id_fuentes"] = infoSubgrupos[i]["_id"]
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
						getProyect["fuentes"] = detalleSubgrupos[0]["dato"]
						getProyect["id_detalle_fuentes"] = detalleSubgrupos[0]["_id"]

					}
				}
			}
		}

	}

	return getProyect
}
