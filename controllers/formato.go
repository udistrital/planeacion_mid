package controllers

import (
	//"fmt"

	//"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/helpers/formatoHelper"
	"github.com/udistrital/planeacion_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// FormatoController operations for Formato
type FormatoController struct {
	beego.Controller
}

// URLMapping ...
func (c *FormatoController) URLMapping() {
	// c.Mapping("Post", c.Post)
	// c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetFormato", c.GetFormato)
	// c.Mapping("GetAll", c.GetAll)
	// c.Mapping("Put", c.Put)
	// c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Formato
// @Param	body		body 	models.Formato	true		"body for Formato content"
// @Success 201 {object} models.Formato
// @Failure 403 body is empty
// @router / [post]
func (c *FormatoController) Post() {

}

// GetOne ...
// @Title GetFormato
// @Description get Formato by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formato
// @Failure 403 :id is empty
// @router /:id [get]
func (c *FormatoController) GetFormato() {

	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	var hijos []models.Nodo
	var plan map[string]interface{}
	var hijosID []map[string]interface{}

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		helpers.LimpiezaRespuestaRefactor(res, &hijosID)
		err := request.GetJson(beego.AppConfig.String("PlanesService")+"/plan/"+id, &res)
		if err != nil {
			return
		}
		helpers.LimpiezaRespuestaRefactor(res, &plan)
		formatoHelper.Limpia(plan)
		tree := formatoHelper.BuildTreeFa(hijos, hijosID)
		c.Data["json"] = tree
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get Formato
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Formato
// @Failure 403
// @router / [get]
func (c *FormatoController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Formato
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Formato	true		"body for Formato content"
// @Success 200 {object} models.Formato
// @Failure 403 :id is not int
// @router /:id [put]
func (c *FormatoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Formato
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *FormatoController) Delete() {

}
