package controllers

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	evaluacionhelper "github.com/udistrital/planeacion_mid/helpers/evaluacionHelper"
)

// EvaluacionController operations for Evaluacion
type EvaluacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *EvaluacionController) URLMapping() {
	c.Mapping("GetEvaluacion", c.GetEvaluacion)
	c.Mapping("GetPlanesPeriodo", c.GetPlanesPeriodo)
	c.Mapping("PlanesAEvaluar", c.PlanesAEvaluar)
	c.Mapping("Unidades", c.Unidades)
	c.Mapping("Avances", c.Avances)
}

// GetPlanesPeriodo ...
// @Title GetPlanesPeriodo
// @Description get Planes y vigencias para la unidad y vigencia dado
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Param	unidad 		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /planes_periodo/:vigencia/:unidad [get]
func (c *EvaluacionController) GetPlanesPeriodo() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "EvaluacionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	vigencia := c.Ctx.Input.Param(":vigencia")
	unidad := c.Ctx.Input.Param(":unidad")

	if len(vigencia) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}
	if len(unidad) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}

	if respuesta, err := evaluacionhelper.GetPlanesPeriodo(unidad, vigencia); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// GetEvaluacion ...
// @Title GetEvaluacion
// @Description get Evaluacion
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Param	plan 		path 	string	true		"The key for staticblock"
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:vigencia/:plan:/:periodo [get]
func (c *EvaluacionController) GetEvaluacion() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "EvaluacionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var evaluacion []map[string]interface{}
	vigencia := c.Ctx.Input.Param(":vigencia")
	plan := c.Ctx.Input.Param(":plan")
	periodoId := c.Ctx.Input.Param(":periodo")

	if len(vigencia) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}
	if len(plan) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}
	if len(periodoId) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}

	trimestres := evaluacionhelper.GetPeriodos(vigencia)
	if len(trimestres) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": nil}
	} else {
		i := 0
		for index, periodo := range trimestres {
			if periodo["_id"] == periodoId {
				i = index
				break
			}
		}

		evaluacion = evaluacionhelper.GetEvaluacion(plan, trimestres, i)

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": evaluacion}
	}

	c.ServeJSON()
}

// Get Planes A Evaluar ...
// @Title PlanesAEvaluar
// @Description get Planes que se pueden evaluar
// @Success 200
// @Failure 404
// @router /planes/ [get]
func (c *EvaluacionController) PlanesAEvaluar() {
	defer helpers.ErrorController(c.Controller, "EvaluacionController")

	if datos, err := evaluacionhelper.PlanesAEvaluar(); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": datos}
	} else {
		panic(map[string]interface{}{"funcion": "PlanesAEvaluar", "err": err, "status": "404", "message": "Error obteniendo los planes a evaluar"})
	}
	c.ServeJSON()
}

// Get Unidades ...
// @Title GetUnidades
// @Description get Unidades
// @Param	plan 		path 	string	true		"The key for staticblock"
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /unidades/:plan:/:vigencia [get]
func (c *EvaluacionController) Unidades() {
	defer helpers.ErrorController(c.Controller, "EvaluacionController")

	plan := c.Ctx.Input.Param(":plan")
	vigencia := c.Ctx.Input.Param(":vigencia")

	if nombrePlan, err := url.QueryUnescape(plan); err == nil {
		if data, err := evaluacionhelper.GetUnidadesPorPlanYVigencia(nombrePlan, vigencia); err == nil {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			panic(map[string]interface{}{"funcion": "Unidades", "err": err, "status": "404", "message": "Error obteniendo las unidades del plan y la vigencia dados"})
		}
	} else {
		panic(map[string]interface{}{"funcion": "Unidades", "err": err, "status": "404", "message": "Error obteniendo las unidades del plan y la vigencia dados"})
	}

	c.ServeJSON()
}

// Get Avance ...
// @Title GetAvance
// @Description get Avance de Unidad
// @Param	plan 		path 	string	true		"The key for staticblock"
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Param	unidad 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /avance/:plan:/:vigencia/:unidad [get]
func (c *EvaluacionController) Avances() {
	defer helpers.ErrorController(c.Controller, "EvaluacionController")

	plan := c.Ctx.Input.Param(":plan")
	vigencia := c.Ctx.Input.Param(":vigencia")
	unidad := c.Ctx.Input.Param(":unidad")

	if nombrePlan, err1 := url.QueryUnescape(plan); err1 == nil {
		if data, err2 := evaluacionhelper.GetAvances(nombrePlan, vigencia, unidad); err2 == nil {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			panic(map[string]interface{}{"funcion": "Avances", "err": err2, "status": "404", "message": "Error obteniendo los avances"})
		}
	} else {
		panic(map[string]interface{}{"funcion": "Avances", "err": err1, "status": "404", "message": "No se pudo obtener el nombre"})
	}

	c.ServeJSON()
}
