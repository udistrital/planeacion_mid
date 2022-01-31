package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	seguimientohelper "github.com/udistrital/planeacion_mid/helpers/seguimientoHelper"
	"github.com/udistrital/utils_oas/request"
)

// SeguimientoController operations for Seguimiento
type SeguimientoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SeguimientoController) URLMapping() {
	c.Mapping("HabilitarReportes", c.HabilitarReportes)
	c.Mapping("CrearReportes", c.CrearReportes)
	c.Mapping("GetPeriodos", c.GetPeriodos)
}

// HabilitarReportes ...
// @Title HabilitarReportes
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /habilitar_reportes/:periodo [put]
func (c *SeguimientoController) HabilitarReportes() {
	periodo := c.Ctx.Input.Param(":periodo")
	var res map[string]interface{}
	var resPut map[string]interface{}
	var reportes []map[string]interface{}

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/seguimiento?query=periodo_id:"+periodo, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &reportes)
		for key, element := range reportes {
			_ = key
			element["activo"] = true

			if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/seguimiento/"+element["_id"].(string), "PUT", &resPut, element); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
			}

		}
		c.Data["json"] = reportes
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// CrearReportes ...
// @Title CrearReportes
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /crear_reportes/:plan/:tipo [post]
func (c *SeguimientoController) CrearReportes() {
	plan_id := c.Ctx.Input.Param(":plan")
	tipo := c.Ctx.Input.Param(":tipo")
	var res map[string]interface{}
	var plan map[string]interface{}
	var respuestaPost map[string]interface{}
	var arrReportes []map[string]interface{}
	reporte := make(map[string]interface{})

	trimestres := seguimientohelper.GetTrimestres("25")

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/plan/"+plan_id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &plan)
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	for i := 0; i < len(trimestres); i++ {
		reporte["nombre"] = "Seguimiento para el " + plan["nombre"].(string)
		reporte["descripcion"] = "Seguimiento para el " + plan["nombre"].(string) + " UNIVERSIDAD DISTRITAL FRANCISCO JOSE DE CALDAS"
		reporte["activo"] = false
		reporte["plan_id"] = plan_id
		reporte["estado_seguimiento_id"] = "61f237df25e40c57a60840d5"
		reporte["periodo_id"] = trimestres[i]["Id"]
		reporte["tipo_seguimiento_id"] = tipo
		reporte["dato"] = "{}"

		if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaPost, reporte); err != nil {
			panic(map[string]interface{}{"funcion": "CrearReportes", "err": "Error creando reporte", "status": "400", "log": err})
		}

		arrReportes = append(arrReportes, respuestaPost["Data"].(map[string]interface{}))
		respuestaPost = nil
	}

	c.Data["json"] = arrReportes
	c.ServeJSON()
}

// GetPeriodos ...
// @Title GetPeriodos
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /get_periodos/:vigencia [get]
func (c *SeguimientoController) GetPeriodos() {
	vigencia := c.Ctx.Input.Param(":vigencia")
	trimestres := seguimientohelper.GetTrimestres(vigencia)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": trimestres}
	c.ServeJSON()
}
