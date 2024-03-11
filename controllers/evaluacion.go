package controllers

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	evaluacionhelper "github.com/udistrital/planeacion_mid/helpers/evaluacionHelper"
	seguimientohelper "github.com/udistrital/planeacion_mid/helpers/seguimientoHelper"
	"github.com/udistrital/utils_oas/request"
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

	var resPlan map[string]interface{}
	var resSeguimiento map[string]interface{}
	var respuesta []map[string]interface{}
	vigencia := c.Ctx.Input.Param(":vigencia")
	unidad := c.Ctx.Input.Param(":unidad")

	if len(vigencia) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}
	if len(unidad) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan?query=estado_plan_id:6153355601c7a2365b2fb2a1,dependencia_id:`+unidad+`,vigencia:`+vigencia, &resPlan); err == nil {
		planes := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(resPlan, &planes)
		if fmt.Sprintf("%v", planes) == "[]" {
			c.Abort("404")
		}

		periodos := evaluacionhelper.GetPeriodos(vigencia)
		trimestres := seguimientohelper.GetTrimestres(vigencia)
		for _, plan := range planes {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,estado_seguimiento_id:622ba49216511e93a95c326d,plan_id:`+plan["_id"].(string), &resSeguimiento); err == nil {
				seguimientos := make([]map[string]interface{}, 1)
				helpers.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)
				if fmt.Sprintf("%v", seguimientos) == "[]" {
					continue
				}

				var periodosSelecionados []map[string]interface{}
				for _, seguimiento := range seguimientos {
					for _, periodo := range periodos {
						if seguimiento["periodo_seguimiento_id"] == periodo["_id"] {
							for _, trimestre := range trimestres {
								var trimestreId float64
								if reflect.TypeOf(trimestre["Id"]).String() == "string" {
									trimestreId, _ = strconv.ParseFloat(trimestre["Id"].(string), 64)
								} else {
									trimestreId = trimestre["Id"].(float64)
								}
								var periodoId float64
								if reflect.TypeOf(periodo["periodo_id"]).String() == "string" {
									periodoId, _ = strconv.ParseFloat(periodo["periodo_id"].(string), 64)
								} else {
									periodoId = periodo["periodo_id"].(float64)
								}

								if trimestreId == periodoId {
									periodosSelecionados = append(periodosSelecionados, map[string]interface{}{"nombre": trimestre["ParametroId"].(map[string]interface{})["Nombre"].(string), "id": periodo["_id"]})
									break
								}
							}
							break
						}
					}
				}

				respuesta = append(respuesta, map[string]interface{}{"plan": plan["nombre"], "id": plan["_id"], "periodos": periodosSelecionados})
			}
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
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
// @Title GetPlanesAEvaluar
// @Description get Planes que se pueden evaluar
// @Success 200
// @Failure 404
// @router /planes/ [get]
func (c *EvaluacionController) PlanesAEvaluar() {
	defer helpers.ErrorController(c.Controller, "EvaluacionController")

	if datos, err := evaluacionhelper.GetPlanesParaEvaluar(); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": datos}
	} else {
		panic(map[string]interface{}{"funcion": "PlanesAEvaluar", "err": err, "status": "400", "message": "Error obteniendo los planes a evaluar"})
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
			panic(map[string]interface{}{"funcion": "Unidades", "err": err, "status": "400", "message": "Error obteniendo las unidades del plan y la vigencia dados"})
		}
	} else {
		panic(map[string]interface{}{"funcion": "Unidades", "err": err, "status": "400", "message": "Error obteniendo las unidades del plan y la vigencia dados"})
	}

	c.ServeJSON()
}
