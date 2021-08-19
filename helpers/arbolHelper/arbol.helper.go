package arbolHelper

import (
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
		forkData["activo"] = hijos[i].Activo

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

		err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/"+child, &res)

		if err != nil {
			return
		}

		helpers.LimpiezaRespuestaRefactor(res, &nodo)
		helpers.LimpiezaRespuestaRefactor(res, &nodoId)
		forkData["id"] = nodoId["_id"]
		forkData["nombre"] = nodo.Nombre
		forkData["descripcion"] = nodo.Descripcion
		forkData["activo"] = nodo.Activo

		if len(nodo.Hijos) > 0 {
			forkData["children"] = getChildren(nodo.Hijos)
		}
		childrenTree = append(childrenTree, forkData)

	}
	return
}
