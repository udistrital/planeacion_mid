package formatoHelper

import (
	"encoding/json"
	"fmt"

	//"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

var validDataT = []string{}

func Limpia() {
	validDataT = []string{}
}

func BuildTreeFa(hijos []models.Nodo, hijosID []map[string]interface{}) [][]map[string]interface{} {

	var tree []map[string]interface{}
	var requeridos []map[string]interface{}
	var nodo []models.NodoDetalle
	var res map[string]interface{}
	var result [][]map[string]interface{}

	for i := 0; i < len(hijos); i++ {
		if hijos[i].Activo {
			forkData := make(map[string]interface{})
			var id string
			forkData["id"] = hijosID[i]["_id"]
			forkData["nombre"] = hijos[i].Nombre
			jsonString, _ := json.Marshal(hijosID[i]["_id"])
			json.Unmarshal(jsonString, &id)
			if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &nodo)
				if len(nodo) == 0 {
					// forkData["type"] = ""
					// forkData["required"] = ""
				} else {
					var deta models.Dato
					json.Unmarshal([]byte(nodo[0].Dato), &deta)
					forkData["type"] = deta.Type
					forkData["required"] = deta.Required
				}
			}
			if len(hijos[i].Hijos) > 0 {
				if len(getChildren(hijos[i].Hijos)) == 0 {
					forkData["sub"] = ""
				} else {
					forkData["sub"] = make([]map[string]interface{}, len(getChildren(hijos[i].Hijos)))
					forkData["sub"] = getChildren(hijos[i].Hijos)
				}
			}
			tree = append(tree, forkData)
			add(id)
		} else {
		}
	}
	requeridos = convert(validDataT)
	result = append(result, tree)
	result = append(result, requeridos)
	return result
}

func getChildren(children []string) (childrenTree []map[string]interface{}) {
	var res map[string]interface{}
	var resp map[string]interface{}
	var nodo models.Nodo
	var nodoId map[string]interface{}
	var detalle []models.NodoDetalle
	for _, child := range children {
		forkData := make(map[string]interface{})
		var id string
		err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+child, &res)
		if err != nil {
			return
		}
		helpers.LimpiezaRespuestaRefactor(res, &nodo)
		helpers.LimpiezaRespuestaRefactor(res, &nodoId)
		if nodo.Activo == true {
			forkData["id"] = nodoId["_id"]
			forkData["nombre"] = nodo.Nombre
			jsonString, _ := json.Marshal(nodoId["_id"])
			json.Unmarshal(jsonString, &id)
			if err_ := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &resp); err_ == nil {
				helpers.LimpiezaRespuestaRefactor(resp, &detalle)
				if len(detalle) == 0 {
					// forkData["type"] = ""
					// forkData["required"] = ""
				} else {
					var deta models.Dato
					json.Unmarshal([]byte(detalle[0].Dato), &deta)
					forkData["type"] = deta.Type
					forkData["required"] = deta.Required
				}
			}
			if len(nodo.Hijos) > 0 {
				if len(getChildren(nodo.Hijos)) == 0 {
					forkData["sub"] = ""
				} else {
					forkData["sub"] = getChildren(nodo.Hijos)
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
	if !contains(validDataT, id) {
		validDataT = append(validDataT, id)
	}
}

func convert(valid []string) []map[string]interface{} {
	fmt.Println(valid)

	var validadores []map[string]interface{}
	forkData := make(map[string]interface{})
	for _, v := range valid {
		forkData[v] = ""
	}

	validadores = append(validadores, forkData)
	fmt.Println(validadores)
	return validadores
}
