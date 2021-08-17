package arbolHelper

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/planes_mid/helpers"
	"github.com/udistrital/planes_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func BuildTree(hijos []models.Nodo) []map[string]interface{} {

	var tree []map[string]interface{}

	for _, hijo := range hijos {
		forkData := make(map[string]interface{})
		//forkData["data"] = hijo
		forkData["id"] = hijo.Id
		forkData["nombre"] = hijo.Nombre
		forkData["descripcion"] = hijo.Descripcion
		forkData["activo"] = hijo.Activo
		fmt.Println(getChildren(hijo.Hijos))
		if len(hijo.Hijos) > 0 {
			forkData["children"] = make([]map[string]interface{}, len(getChildren(hijo.Hijos)))
			forkData["children"] = getChildren(hijo.Hijos)
		}
		//fmt.Println(forkData["children"])
		tree = append(tree, forkData)
	}

	return tree

}

func getChildren(children []string) (childrenTree []map[string]interface{}) {
	var res map[string]interface{}
	var nodo models.Nodo
	for _, child := range children {
		forkData := make(map[string]interface{})

		err := request.GetJson(beego.AppConfig.String("UrlCrud")+"/subgrupo/"+child, &res)
		if err != nil {
			return
		}
		helpers.LimpiezaRespuestaRefactor(res, &nodo)
		//forkData["data"] = nodo
		fmt.Println(res)
		forkData["id"] = nodo.Id
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
