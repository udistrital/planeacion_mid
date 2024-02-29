package formatoHelper

import (
	"encoding/json"

	//"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

var validDataT = []string{}
var estadoPlan string

func Limpia(plan map[string]interface{}) {
	validDataT = []string{}
	jsonString, _ := json.Marshal(plan["estado_plan_id"])
	json.Unmarshal(jsonString, &estadoPlan)
}

func BuildTreeFaActEst(hijos []models.Nodo, hijosID []map[string]interface{}) [][]map[string]interface{} {
	return ConstruirArbolFormato(hijos, hijosID, true)
}

func BuildTreeFa(hijos []models.Nodo, hijosID []map[string]interface{}) [][]map[string]interface{} {
	return ConstruirArbolFormato(hijos, hijosID, false)
}

func ConstruirArbolFormato(hijos []models.Nodo, hijosID []map[string]interface{}, activos bool) [][]map[string]interface{} {
	var tree []map[string]interface{}
	var requeridos []map[string]interface{}
	var nodo []models.NodoDetalle
	var res map[string]interface{}
	var result [][]map[string]interface{}

	for i := 0; i < len(hijos); i++ {
		if activos || hijos[i].Activo {
			forkData := make(map[string]interface{})
			var id string
			forkData["id"] = hijosID[i]["_id"]
			forkData["nombre"] = hijos[i].Nombre
			forkData["ref"] = hijosID[i]["ref"]
			jsonString, _ := json.Marshal(hijosID[i]["_id"])
			json.Unmarshal(jsonString, &id)
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &nodo)
				if len(nodo) == 0 {
					// forkData["type"] = ""
					// forkData["required"] = ""
				} else {
					var deta map[string]interface{}
					json.Unmarshal([]byte(nodo[0].Dato), &deta)
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
			if len(hijos[i].Hijos) > 0 {
				if len(getChildren(hijos[i].Hijos, activos)) == 0 {
					forkData["sub"] = ""
				} else {
					forkData["sub"] = make([]map[string]interface{}, len(getChildren(hijos[i].Hijos, activos)))
					forkData["sub"] = getChildren(hijos[i].Hijos, activos)
				}
			}
			tree = append(tree, forkData)
			add(id)
		}
	}
	requeridos = convert(validDataT)
	result = append(result, tree)
	result = append(result, requeridos)
	return result
}

func getChildren(children []string, activos bool) (childrenTree []map[string]interface{}) {
	var res map[string]interface{}
	var resp map[string]interface{}
	var nodo models.Nodo
	var nodoId map[string]interface{}
	var detalle []models.NodoDetalle
	for _, child := range children {
		forkData := make(map[string]interface{})
		var id string
		err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+child, &res)
		if err != nil {
			return
		}
		helpers.LimpiezaRespuestaRefactor(res, &nodo)
		helpers.LimpiezaRespuestaRefactor(res, &nodoId)
		if activos || nodo.Activo {
			forkData["id"] = nodoId["_id"]
			forkData["nombre"] = nodo.Nombre
			forkData["ref"] = nodoId["ref"]
			jsonString, _ := json.Marshal(nodoId["_id"])
			json.Unmarshal(jsonString, &id)
			if err_ := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id, &resp); err_ == nil {
				helpers.LimpiezaRespuestaRefactor(resp, &detalle)
				if len(detalle) == 0 {
					// forkData["type"] = ""
					// forkData["required"] = ""
				} else {
					var deta map[string]interface{}
					json.Unmarshal([]byte(detalle[0].Dato), &deta)
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
			if len(nodo.Hijos) > 0 {
				if len(getChildren(nodo.Hijos, activos)) == 0 {
					forkData["sub"] = ""
				} else {
					forkData["sub"] = getChildren(nodo.Hijos, activos)
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

	var validadores []map[string]interface{}
	forkData := make(map[string]interface{})
	for _, v := range valid {
		if v == "" {

		} else {
			forkData[v] = ""
			forkData[v+"_o"] = ""
		}
	}

	validadores = append(validadores, forkData)
	return validadores
}
