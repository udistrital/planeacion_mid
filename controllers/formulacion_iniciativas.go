package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"

	formulacionIniciativasHelper "github.com/udistrital/planeacion_mid/helpers/formulacionIniciativasHelper"
)

// FormulacionIniciativasController operations for Formulacion Iniciativas
type FormulacionIniciativasController struct {
	beego.Controller
}

func (c *FormulacionIniciativasController) URLMapping() {
	c.Mapping("DefinirFechasIniciativas", c.DefinirFechasIniciativas)
	c.Mapping("ConsultarFechasIniciativas", c.ConsultarFechasIniciativas)

}

// DefinirFechasIniciativas ...
// @Title DefinirFechasIniciativas
// @Description Peticion POST para definir fechas en iniciativas de proyectos de inversión para los procesos de formulación
// @Param	body		body 	{}	true		"body for Project content"
// @Success 200 {object} models.Formulacion
// @Failure 400 bad request
// @router /habilitar_fechas_iniciativa [post]
func (c *FormulacionIniciativasController) DefinirFechasIniciativas() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "FormulacionIniciativasController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var res interface{}

	// Decodificar JSON desde el cuerpo de la solicitud
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error al decodificar JSON desde el cuerpo de la solicitud", "status": "400", "log": err})
	}
	res = formulacionIniciativasHelper.DefinirFechasIniciativasFormulacion(body)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": res}
	c.ServeJSON()
}

// ConsultarFechasIniciativas ...
// @Title ConsultarFechasIniciativas
// @Description Peticion POST para consultar las fechas de formulación de iniciativas de proyectos de inversión de una vigencia y una unidad
// @Param	body		body	{}	true		"body for Project content"
// @Success 200 {object}
// @Failure 400 bad request
// @router /consultar_fechas_iniciativa [post]
func (c *FormulacionIniciativasController) ConsultarFechasIniciativas() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "FormulacionIniciativasController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var res interface{}

	// Decodificar JSON desde el cuerpo de la solicitud
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err != nil {
		panic(map[string]interface{}{"funcion": "ConsultarFechasIniciativasFormulacion", "err": "Error al decodificar JSON desde el cuerpo de la solicitud", "status": "400", "log": err})
	}

	res = formulacionIniciativasHelper.ConsultarFechasIniciativasFormulacion(body)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": res}
	c.ServeJSON()
}
