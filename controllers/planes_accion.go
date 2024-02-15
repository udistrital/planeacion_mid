package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	formulacioPlanAccionnHelper "github.com/udistrital/planeacion_mid/helpers/planesAccionHelper"
)

type PlanesAccionController struct {
	beego.Controller
}

func (c *PlanesAccionController) URLMapping() {
	c.Mapping("ObtenerPlanesDeAccion", c.PlanesDeAccion)
}

// Get Planes De Accion ...
// @Title GetPlanesDeAccion
// @Description get Planes de Acción
// @Success 200 {object} models.Formulacion
// @Failure 400 bad response
// @router / [get]
func (c *PlanesAccionController) PlanesDeAccion() {
	defer helpers.ErrorController(c.Controller, "PlanesAccionController")
	if datos, err := formulacioPlanAccionnHelper.ObtenerPlanesAccion(); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": datos}
	} else {
		panic(map[string]interface{}{"funcion": "PlanesDeAccion", "err": err, "status": "400", "message": "Error obteniendo los datos"})
	}
	c.ServeJSON()
}
