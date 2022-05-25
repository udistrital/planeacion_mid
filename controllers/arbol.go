package controllers

import (
	//"fmt"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/helpers/arbolHelper"
	"github.com/udistrital/planeacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// ArbolController operations for Arbol
type ArbolController struct {
	beego.Controller
}

// URLMapping ...
func (c *ArbolController) URLMapping() {
	c.Mapping("GetArbol", c.GetArbol)
	c.Mapping("DeletePlan", c.DeletePlan)
	c.Mapping("DeleteNodo", c.DeleteNodo)
	c.Mapping("ActivarNodo", c.ActivarNodo)
	c.Mapping("ActivarPlan", c.ActivarPlan)

}

// GetArbol ...
// @Title GetArbol
// @Description get Arbol by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Arbol
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ArbolController) GetArbol() {

	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	var hijos []models.Nodo
	var hijosID []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		helpers.LimpiezaRespuestaRefactor(res, &hijosID)
		tree := arbolHelper.BuildTree(hijos, hijosID)
		fmt.Println(tree)
		if len(tree) != 0 {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": tree}
		} else {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()

}

// DeletePlan ...
// @Title DeletePlan
// @Description delete the Plan Arbol
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /desactivar_plan/:id [delete]
func (c *ArbolController) DeletePlan() {
	fmt.Println("entra melomano")
	id := c.Ctx.Input.Param(":id")
	var plan map[string]interface{}
	var res map[string]interface{}
	var resPut map[string]interface{}
	var resHijos map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res); err == nil {

		helpers.LimpiezaRespuestaRefactor(res, &plan)
		fmt.Println(plan)
		plan["activo"] = false
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan["_id"].(string), "PUT", &resPut, plan); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteHijos", "err": "Error actualizacion activo \"id\"", "status": "400", "log": err})
		}
		fmt.Println("entra aca primeros hijos")
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan["_id"].(string), &resHijos); err == nil {
			fmt.Println("consulta hijos")
			fmt.Println(resHijos)
			helpers.LimpiezaRespuestaRefactor(resHijos, &hijos)
			arbolHelper.DeleteHijos(hijos)
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": plan}

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": res, "Type": "error subgrupo"}
		c.Abort("400")
	}

	c.ServeJSON()

}

// DeleteNodo ...
// @Title DeleteNodo
// @Description delete the Nodo Arbol
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /desactivar_nodo/:id [delete]
func (c *ArbolController) DeleteNodo() {
	fmt.Println("entra melomano")
	id := c.Ctx.Input.Param(":id")
	var subgrupo map[string]interface{}
	var res map[string]interface{}
	var resPut map[string]interface{}
	var resHijos map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+id, &res); err == nil {

		helpers.LimpiezaRespuestaRefactor(res, &subgrupo)
		subgrupo["activo"] = false
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+subgrupo["_id"].(string), "PUT", &resPut, subgrupo); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteHijos", "err": "Error actualizacion activo \"id\"", "status": "400", "log": err})
		}
		fmt.Println("entra aca primeros hijos")
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+subgrupo["_id"].(string), &resHijos); err == nil {
			fmt.Println("consulta hijos")
			fmt.Println(resHijos)
			helpers.LimpiezaRespuestaRefactor(resHijos, &hijos)
			arbolHelper.DeleteHijos(hijos)
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": subgrupo}

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": res, "Type": "error subgrupo"}
		c.Abort("400")
	}

	c.ServeJSON()

}

// ActivarPlan ...
// @Title ActivarPlan
// @Description activar the Plan Arbol
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /activar_plan/:id [put]
func (c *ArbolController) ActivarPlan() {
	id := c.Ctx.Input.Param(":id")
	var plan map[string]interface{}
	var res map[string]interface{}
	var resPut map[string]interface{}
	var resHijos map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res); err == nil {

		helpers.LimpiezaRespuestaRefactor(res, &plan)
		plan["activo"] = true
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan["_id"].(string), "PUT", &resPut, plan); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteHijos", "err": "Error actualizacion activo \"id\"", "status": "400", "log": err})
		}
		fmt.Println("entra aca primeros hijos")
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan["_id"].(string), &resHijos); err == nil {
			helpers.LimpiezaRespuestaRefactor(resHijos, &hijos)
			arbolHelper.ActivarHijos(hijos)
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": plan}

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": res, "Type": "error subgrupo"}
		c.Abort("400")
	}

	c.ServeJSON()

}

// ActivarNodo ...
// @Title ActivarNodo
// @Description put the Nodo Arbol
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /activar_nodo/:id [put]
func (c *ArbolController) ActivarNodo() {
	id := c.Ctx.Input.Param(":id")
	var subgrupo map[string]interface{}
	var res map[string]interface{}
	var resPut map[string]interface{}
	var resHijos map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+id, &res); err == nil {

		helpers.LimpiezaRespuestaRefactor(res, &subgrupo)
		subgrupo["activo"] = true
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+subgrupo["_id"].(string), "PUT", &resPut, subgrupo); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteHijos", "err": "Error actualizacion activo \"id\"", "status": "400", "log": err})
		}
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+subgrupo["_id"].(string), &resHijos); err == nil {
			helpers.LimpiezaRespuestaRefactor(resHijos, &hijos)
			arbolHelper.ActivarHijos(hijos)
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": subgrupo}

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": res, "Type": "error subgrupo"}
		c.Abort("400")
	}

	c.ServeJSON()

}
