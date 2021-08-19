package controllers

import (
	//"fmt"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/planes_mid/helpers"
	"github.com/udistrital/planes_mid/helpers/arbolHelper"
	"github.com/udistrital/planes_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// ArbolController operations for Arbol
type ArbolController struct {
	beego.Controller
}

// URLMapping ...
func (c *ArbolController) URLMapping() {
	//c.Mapping("Post", c.Post)
	c.Mapping("GetArbol", c.GetArbol)
	//c.Mapping("BuildTree", c.BuildTree)
	//c.Mapping("Put", c.Put)
	//c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Arbol
// @Param	body		body 	models.Arbol	true		"body for Arbol content"
// @Success 201 {object} models.Arbol
// @Failure 403 body is empty
// @router / [post]
func (c *ArbolController) Post() {

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
	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		helpers.LimpiezaRespuestaRefactor(res, &hijosID)
		tree := arbolHelper.BuildTree(hijos, hijosID)
		c.Data["json"] = tree
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()

}

// GetAll ...
// @Title GetAll
// @Description get Arbol
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Arbol
// @Failure 403
// @router / [get]
func (c *ArbolController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Arbol
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Arbol	true		"body for Arbol content"
// @Success 200 {object} models.Arbol
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ArbolController) Put() {

}

// Delete ...
// @Title DeletePlan
// @Description delete the Plan Arbol
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ArbolController) DeletePlan() {
	id := c.Ctx.Input.Param(":id")
	var plan map[string]interface{}

	if err := request.SendJson(beego.AppConfig.String("PlanesService")+"/plan/delete_plan/"+id, "PUT", &plan, plan); err == nil {

		var res map[string]interface{}
		var hijos []models.Nodo
		var hijosId []map[string]interface{}
		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {

			helpers.LimpiezaRespuestaRefactor(res, &hijos)
			helpers.LimpiezaRespuestaRefactor(res, &hijosId)

			for i := 0; i < len(hijos); i++ {
				var subgrupo map[string]interface{}
				var idSubgrupo string

				jsonString, _ := json.Marshal(hijosId[i]["_id"])
				json.Unmarshal(jsonString, &idSubgrupo)
				request.SendJson(beego.AppConfig.String("PlanesService")+"/subgrupo/delete_nodo/"+idSubgrupo, "PUT", &res, subgrupo)
			}

		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": res, "Type": "error subgrupo"}
			c.Abort("400")
		}

		c.Data["json"] = plan

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error plan"}
		c.Abort("400")
	}

	c.ServeJSON()

}
