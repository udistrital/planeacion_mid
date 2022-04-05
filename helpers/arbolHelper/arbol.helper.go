package arbolHelper

import (
	//"fmt"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func BuildTree(hijos []models.Nodo, hijosID []map[string]interface{}) []map[string]interface{} {

	var tree []map[string]interface{}

	for i := 0; i < len(hijos); i++ {

		forkData := make(map[string]interface{})
		forkData["id"] = hijosID[i]["_id"]
		forkData["nombre"] = hijos[i].Nombre
		forkData["descripcion"] = hijos[i].Descripcion
		if hijos[i].Activo {
			forkData["activo"] = "activo"
		} else {
			forkData["activo"] = "inactivo"
		}

		if len(hijos[i].Hijos) > 0 {
			forkData["children"] = make([]map[string]interface{}, len(getChildren(hijos[i].Hijos)))
			forkData["children"] = getChildren(hijos[i].Hijos)
		}
		tree = append(tree, forkData)
	}

	return tree

}

func getChildren(children []string) (childrenTree []map[string]interface{}) {
	var res map[string]interface{}
	var nodo models.Nodo
	var nodoId map[string]interface{}
	for _, child := range children {
		forkData := make(map[string]interface{})

		err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+child, &res)
		if err != nil {
			return
		}

		helpers.LimpiezaRespuestaRefactor(res, &nodo)
		helpers.LimpiezaRespuestaRefactor(res, &nodoId)
		forkData["id"] = nodoId["_id"]
		forkData["nombre"] = nodo.Nombre
		forkData["descripcion"] = nodo.Descripcion
		if nodo.Activo == true {
			forkData["activo"] = "activo"
		} else {
			forkData["activo"] = "inactivo"
		}

		if len(nodo.Hijos) > 0 {
			forkData["children"] = getChildren(nodo.Hijos)
		}
		childrenTree = append(childrenTree, forkData)

	}
	return
}

func DeleteHijos(children []map[string]interface{}) {
	var res map[string]interface{}
	var resHijos map[string]interface{}
	var hijos []map[string]interface{}
	for i := 0; i < len(children); i++ {
		aux := children[i]
		aux["activo"] = false
		fmt.Println(aux)
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+aux["_id"].(string), "PUT", &res, aux); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteHijos", "err": "Error actualizacion activo \"id\"", "status": "400", "log": err})
		}

		if len(aux["hijos"].([]interface{})) > 0 {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+aux["_id"].(string), &resHijos); err == nil {
				helpers.LimpiezaRespuestaRefactor(resHijos, &hijos)
				fmt.Println(hijos)
				DeleteHijos(hijos)
			}
		}
	}

}

func ActivarHijos(children []map[string]interface{}) {
	var res map[string]interface{}
	var resHijos map[string]interface{}
	var hijos []map[string]interface{}
	for i := 0; i < len(children); i++ {
		aux := children[i]
		aux["activo"] = true
		fmt.Println(aux)
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+aux["_id"].(string), "PUT", &res, aux); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteHijos", "err": "Error actualizacion activo \"id\"", "status": "400", "log": err})
		}

		if len(aux["hijos"].([]interface{})) > 0 {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+aux["_id"].(string), &resHijos); err == nil {
				helpers.LimpiezaRespuestaRefactor(resHijos, &hijos)
				fmt.Println(hijos)
				ActivarHijos(hijos)
			}
		}
	}

}
