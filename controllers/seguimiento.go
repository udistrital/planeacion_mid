package controllers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

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
	c.Mapping("GetActividadesGenerales", c.GetActividadesGenerales)
	c.Mapping("GetDataActividad", c.GetDataActividad)
	c.Mapping("GuardarSeguimiento", c.GuardarSeguimiento)
	c.Mapping("GetSeguimiento", c.GetSeguimiento)
	c.Mapping("GetIndicadores", c.GetIndicadores)
	c.Mapping("GetAvanceIndicador", c.GetAvanceIndicador)
	c.Mapping("GetEstadoTrimestre", c.GetEstadoTrimestre)
	c.Mapping("GuardarDocumentos", c.GuardarDocumentos)
	c.Mapping("GuardarCualitativo", c.GuardarCualitativo)
	c.Mapping("GuardarCuantitativo", c.GuardarCuantitativo)
	c.Mapping("ReportarActividad", c.ReportarActividad)
	c.Mapping("ReportarSeguimiento", c.ReportarSeguimiento)
	c.Mapping("RevisarActividad", c.RevisarActividad)
	c.Mapping("RevisarSeguimiento", c.RevisarSeguimiento)
}

// HabilitarReportes ...
// @Title HabilitarReportes
// @Description get Seguimiento
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403
// @router /habilitar_reportes [put]
func (c *SeguimientoController) HabilitarReportes() {

	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var res map[string]interface{}
	var resPut map[string]interface{}
	var reportes []map[string]interface{}
	var entrada map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento?query=periodo_id:"+entrada["periodo_id"].(string), &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &reportes)
		if len(reportes) > 0 {
			var element = reportes[0]
			element["activo"] = true
			element["fecha_inicio"] = entrada["fecha_inicio"]
			element["fecha_fin"] = entrada["fecha_fin"]

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+element["_id"].(string), "PUT", &resPut, element); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
			}

		} else {
			element := map[string]interface{}{
				"activo":       true,
				"fecha_inicio": entrada["fecha_inicio"],
				"fecha_fin":    entrada["fecha_fin"],
				"periodo_id":   entrada["periodo_id"],
			}

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento", "POST", &resPut, element); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle", "status": "400", "log": err})
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento?query=periodo_id:"+entrada["periodo_id"].(string), &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &reportes)
			}
		}
		c.Data["json"] = reportes
	} else {
		panic(err)
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
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	plan_id := c.Ctx.Input.Param(":plan")
	tipo := c.Ctx.Input.Param(":tipo")
	var res map[string]interface{}
	var resTrimestres map[string]interface{}
	var plan map[string]interface{}
	var respuestaPost map[string]interface{}
	var arrReportes []map[string]interface{}
	reporte := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan_id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &plan)
		trimestres := seguimientohelper.GetTrimestres(plan["vigencia"].(string))

		for i := 0; i < len(trimestres); i++ {
			periodo := int(trimestres[i]["Id"].(float64))
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento?query=periodo_id:"+strconv.Itoa(periodo), &resTrimestres); err == nil {
				reporte["nombre"] = "Seguimiento para el " + plan["nombre"].(string)
				reporte["descripcion"] = "Seguimiento para el " + plan["nombre"].(string) + " UNIVERSIDAD DISTRITAL FRANCISCO JOSE DE CALDAS"
				reporte["activo"] = false
				reporte["plan_id"] = plan_id
				reporte["estado_seguimiento_id"] = "61f237df25e40c57a60840d5"
				reporte["periodo_seguimiento_id"] = resTrimestres["Data"].([]interface{})[0].(map[string]interface{})["_id"]
				reporte["tipo_seguimiento_id"] = tipo
				reporte["dato"] = "{}"
				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaPost, reporte); err != nil {
					panic(map[string]interface{}{"funcion": "CrearReportes", "err": "Error creando reporte", "status": "400", "log": err})
				}

				arrReportes = append(arrReportes, respuestaPost["Data"].(map[string]interface{}))
				respuestaPost = nil
			} else {
				panic(err)
			}
		}

	} else {
		panic(err)
	}

	c.Data["json"] = arrReportes
	c.ServeJSON()
}

// GetPeriodos ...
// @Title GetPeriodos
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /get_periodos/:vigencia [get]
func (c *SeguimientoController) GetPeriodos() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	vigencia := c.Ctx.Input.Param(":vigencia")
	if len(vigencia) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
	}
	trimestres := seguimientohelper.GetTrimestres(vigencia)
	fmt.Println(trimestres)
	if len(trimestres) == 0 || trimestres[0]["Id"] == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": nil}

	} else {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": trimestres}
	}
	c.ServeJSON()
}

// GetActividadesGenerales ...
// @Title GetActividadeGenerales
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /get_actividades/:seguimiento_id [get]
func (c *SeguimientoController) GetActividadesGenerales() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	seguimiento_id := c.Ctx.Input.Param(":seguimiento_id")
	var resSeguimiento map[string]interface{}
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var seguimiento []map[string]interface{}
	var datoPlan map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=_id:"+seguimiento_id, &resSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &seguimiento)
		if fmt.Sprintf("%v", seguimiento) != "[]" {
			planId := seguimiento[0]["plan_id"].(string)
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

				for i := 0; i < len(subgrupos); i++ {
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {

						actividades := seguimientohelper.GetActividades(subgrupos[i]["_id"].(string))

						if seguimiento[0]["dato"] == "{}" {
							for _, actividad := range actividades {
								actividad["estado"] = map[string]interface{}{"nombre": "Sin reporte"}
							}
						} else {
							dato_plan_str := seguimiento[0]["dato"].(string)
							json.Unmarshal([]byte(dato_plan_str), &datoPlan)

							for indexActividad, element := range datoPlan {
								for _, actividad := range actividades {
									if indexActividad == actividad["index"] {
										actividad["estado"] = element.(map[string]interface{})["estado"]
									}
								}
							}
							for _, actividad := range actividades {
								if actividad["estado"] == nil {
									actividad["estado"] = map[string]interface{}{"nombre": "Sin reporte"}
								}
							}
						}
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": actividades}
						break
					}
				}
			}
		} else {
			panic(err)
		}
	}
	c.ServeJSON()
}

// GetDataActividad ...
// @Title GetDataActividad
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /get_data/:plan_id/:index [get]
func (c *SeguimientoController) GetDataActividad() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	plan_id := c.Ctx.Input.Param(":plan_id")
	index := c.Ctx.Input.Param(":index")
	var res map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		data := seguimientohelper.GetDataSubgrupos(hijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// GuardarSeguimiento ...
// @Title GuardarSeguimiento
// @Description post Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /guardar_seguimiento/:plan_id/:index/:trimestre [post]
func (c *SeguimientoController) GuardarSeguimiento() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var evidencias []map[string]interface{}
	var resEstado map[string]interface{}
	var estadoSeguimiento string
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {

			for _, evidencia := range body["evidencia"].([]interface{}) {
				if evidencia.(map[string]interface{})["Enlace"] != nil {
					evidencias = append(evidencias, evidencia.(map[string]interface{}))
				}
			}
			body["evidencia"] = evidencias

			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]

			if seguimiento["dato"] == "{}" {
				dato[indexActividad] = body

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
				dato[indexActividad].(map[string]interface{})["estado"] = estado

				b, _ := json.Marshal(dato)
				str := string(b)
				seguimiento["dato"] = str
			} else {
				datoStr := seguimiento["dato"].(string)
				json.Unmarshal([]byte(datoStr), &dato)

				dato[indexActividad] = body
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
				dato[indexActividad].(map[string]interface{})["estado"] = estado

				b, _ := json.Marshal(dato)
				str := string(b)
				seguimiento["dato"] = str
			}
			estadoSeguimiento = seguimientohelper.GetEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarSeguimiento", "err": "Error actualizando seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta["Data"]}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	c.ServeJSON()

}

// GetSeguimiento ...
// @Title GetSeguimiento
// @Description get Seguimiento
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /get_seguimiento/:plan_id/:index/:trimestre [get]
func (c *SeguimientoController) GetSeguimiento() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestreId := c.Ctx.Input.Param(":trimestre")
	var respuesta map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var resEstado map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var seguimientoActividad map[string]interface{}
	var periodo []map[string]interface{}
	var trimestre string
	dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestreId, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux[0]

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
			helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
				trimestre = periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)
			}
		}

		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &dato)

		actividad, _ := json.Marshal(seguimientohelper.GetActividad(seguimiento, indexActividad, trimestre))
		json.Unmarshal([]byte(string(actividad)), &seguimientoActividad)
		seguimientoActividad["_id"] = seguimiento["_id"].(string)
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+seguimiento["estado_seguimiento_id"].(string), &resEstado); err == nil {
			seguimientoActividad["estadoSeguimiento"] = resEstado["Data"].(map[string]interface{})["nombre"].(string)
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": seguimientoActividad}
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// GetIndicadores ...
// @Title Indicadores
// @Description get Seguimiento
// @Param	plan_id 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /get_indicadores/:plan_id [get]
func (c *SeguimientoController) GetIndicadores() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	plan_id := c.Ctx.Input.Param(":plan_id")
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var hijos []map[string]interface{}
	var indicadores []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "indicador") {

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+subgrupos[i]["_id"].(string), &res); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &hijos)
					for j := range hijos {
						if strings.Contains(strings.ToLower(hijos[j]["nombre"].(string)), "indicador") {
							aux := hijos[j]
							indicadores = append(indicadores, aux)
						}
					}
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": indicadores}
				} else {
					panic(err)
				}
				break
			}
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GetAvanceIndicador ...
// @Title GetAvanceIndicador
// @Description post Seguimiento by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /get_avance [post]
func (c *SeguimientoController) GetAvanceIndicador() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var res map[string]interface{}
	var avancedata []map[string]interface{}
	var res1 map[string]interface{}
	var avancedata1 []map[string]interface{}
	var res2 map[string]interface{}
	var resName map[string]interface{}
	var parametro_periodo_name []map[string]interface{}
	var avancedata2 []map[string]interface{}
	var parametro_periodo []map[string]interface{}
	var dato map[string]interface{}
	var seguimiento map[string]interface{}
	var seguimiento1 map[string]interface{}
	var test1 string
	var periodIdString string
	var periodId float64
	var avanceAcumulado string
	var testavancePeriodo string
	var nombrePeriodo string
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_seguimiento_id:"+body["periodo_seguimiento_id"].(string), &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &avancedata)

		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+body["periodo_seguimiento_id"].(string), &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &parametro_periodo)
			paramIdlen := parametro_periodo[0]

			paramId := paramIdlen["ParametroId"].(map[string]interface{})
			if paramId["CodigoAbreviacion"] != "T1" {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_seguimiento_id:"+body["periodo_seguimiento_id"].(string), &res1); err == nil {
					helpers.LimpiezaRespuestaRefactor(res1, &avancedata1)
					seguimiento1 = avancedata1[0]
					datoStrUltimoTrimestre := seguimiento1["dato"].(string)
					if datoStrUltimoTrimestre == "{}" {
						test1 = body["periodo_seguimiento_id"].(string)
						priodoId_rest, err := strconv.ParseFloat(test1, 8)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(test1)
						periodId = priodoId_rest - 1
					} else {
						test1 = body["periodo_seguimiento_id"].(string)
						priodoId_rest, err := strconv.ParseFloat(test1, 8)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(test1)
						periodId = priodoId_rest
					}

				} else {
					panic(err)
				}
				periodIdString = fmt.Sprint(periodId)
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_seguimiento_id:"+periodIdString, &res2); err == nil {
					helpers.LimpiezaRespuestaRefactor(res2, &avancedata2)
					seguimiento = avancedata2[0]
					datoStr := seguimiento["dato"].(string)
					json.Unmarshal([]byte(datoStr), &dato)
					indicador1 := dato[body["index"].(string)].(map[string]interface{})
					avanceIndicador1 := indicador1[body["Nombre_del_indicador"].(string)].(map[string]interface{})
					avanceAcumulado = avanceIndicador1["avanceAcumulado"].(string)
					testavancePeriodo = avanceIndicador1["avancePeriodo"].(string)
				} else {
					panic(err)
				}
				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+body["periodo_seguimiento_id"].(string), &resName); err == nil {
					helpers.LimpiezaRespuestaRefactor(resName, &parametro_periodo_name)
					paramIdlenName := parametro_periodo_name[0]

					paramIdName := paramIdlenName["ParametroId"].(map[string]interface{})
					nombrePeriodo = paramIdName["CodigoAbreviacion"].(string)
				} else {
					panic(err)
				}
			} else {
				fmt.Println("")
			}
		} else {
			panic(err)
		}
		avancePeriodo := body["avancePeriodo"].(string)
		aPe, err := strconv.ParseFloat(avancePeriodo, 8)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(aPe, err, reflect.TypeOf(avanceAcumulado))
		aAc, err := strconv.ParseFloat(avanceAcumulado, 8)
		if err != nil {
			fmt.Println(err)
		}
		totalAcumulado := fmt.Sprint(aPe + aAc)
		generalData := make(map[string]interface{})
		generalData["avancePeriodo"] = avancePeriodo
		generalData["periodIdString"] = periodIdString
		generalData["avanceAcumulado"] = totalAcumulado
		generalData["avancePeriodoPrev"] = testavancePeriodo
		generalData["avanceAcumuladoPrev"] = avanceAcumulado
		generalData["nombrePeriodo"] = nombrePeriodo

		fmt.Println(avanceAcumulado, reflect.TypeOf(avanceAcumulado))
		fmt.Println(avancePeriodo, reflect.TypeOf(avancePeriodo))
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": generalData}
		fmt.Println(dato)
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// GetEstadoTrimestre ...
// @Title GetEstadoTrimestre
// @Description get Seguimiento del trimestre correspondiente
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @Failure 404 not found resource
// @router /get_estado_trimestre/:plan_id/:trimestre [get]
func (c *SeguimientoController) GetEstadoTrimestre() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "SeguimientoController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var planes []map[string]interface{}
	var periodoSeguimiento []map[string]interface{}

	planId := c.Ctx.Input.Param(":plan_id")
	trimestre := c.Ctx.Input.Param(":trimestre")

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=plan_id:"+planId, &resSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &planes)

		for _, plan := range planes {
			var periodo []map[string]interface{}
			periodoSeguimientoId := plan["periodo_seguimiento_id"].(string)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento?query=_id:"+periodoSeguimientoId, &resPeriodoSeguimiento); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)
				if fmt.Sprintf("%v", periodoSeguimiento[0]) != "map[]" {

					if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento[0]["periodo_id"].(string)+",ParametroId__CodigoAbreviacion:"+trimestre, &resPeriodo); err == nil {
						helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
						plan["periodo_seguimiento_id"] = periodoSeguimiento[0]

						if fmt.Sprintf("%v", periodo[0]) != "map[]" {
							var resEstado map[string]interface{}

							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+plan["estado_seguimiento_id"].(string), &resEstado); err == nil {
								plan["estado_seguimiento_id"] = resEstado["Data"]

								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan["plan_id"].(string), &resEstado); err == nil {
									plan["plan_id"] = resEstado["Data"]

									c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": plan}
									break
								}
							}
						}
					}
				}
			}
		}

	} else {
		panic(err)
	}
	c.ServeJSON()
}

// GuardarDocumentos ...
// @Title GuardarDocumentos
// @Description post Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_documentos/:plan_id/:index/:trimestre [post]
func (c *SeguimientoController) GuardarDocumentos() {
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var resEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var evidencias []map[string]interface{}
	var estadoSeguimiento string
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	comentario := false

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {

			for _, evidencia := range body["evidencia"].([]interface{}) {
				if evidencia.(map[string]interface{})["Enlace"] != nil {
					evidencias = append(evidencias, evidencia.(map[string]interface{}))
					if evidencia.(map[string]interface{})["Observacion"] != nil && evidencia.(map[string]interface{})["Observacion"] != "Sin observación" {
						comentario = true
					}
				}
			}

			if body["documento"] != nil {
				resDocs := helpers.GuardarDocumento(body["documento"].([]interface{}))

				for _, doc := range resDocs {

					evidencias = append(evidencias, map[string]interface{}{
						"Id":     doc.(map[string]interface{})["Id"],
						"Enlace": doc.(map[string]interface{})["Enlace"],
						"nombre": doc.(map[string]interface{})["Nombre"],
						"TipoDocumento": map[string]interface{}{
							"id":                doc.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["Id"],
							"codigoAbreviacion": doc.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["CodigoAbreviacion"],
						},
						"Observacion": "",
						"Activo":      true,
					})
				}
			}

			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]

			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)
			if dato[indexActividad] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
				dato[indexActividad] = map[string]interface{}{"estado": estado, "evidencia": evidencias}
			} else {
				estado = dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})
				if estado["nombre"] != "Actividad en reporte" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				}
				if comentario {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:CO", &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				}
				dato[indexActividad].(map[string]interface{})["evidencia"] = evidencias
				dato[indexActividad].(map[string]interface{})["estado"] = estado
			}

			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str

			estadoSeguimiento = seguimientohelper.GetEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarDocumentos", "err": "Error guardado documentos del seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta["Data"]}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}
	c.ServeJSON()

}

// GuardarCualitativo ...
// @Title GuardarCualitativo
// @Description post Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_cualitativo/:plan_id/:index/:trimestre [post]
func (c *SeguimientoController) GuardarCualitativo() {
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var resEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var cualitativo map[string]interface{}
	var estadoSeguimiento string
	observacion := false
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]

			cualitativo = body["cualitativo"].(map[string]interface{})

			if cualitativo["observaciones"] != "" && cualitativo["observaciones"] != "Sin observación" {
				observacion = true
			}

			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if dato[indexActividad] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
				dato[indexActividad] = map[string]interface{}{"estado": estado, "cualitativo": cualitativo}
			} else {
				estado = dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})

				if estado["nombre"] == "Actividad reportada" {
					var codigo_abreviacion string

					if observacion {
						codigo_abreviacion = "CO"
					} else {
						codigo_abreviacion = "AAV"
					}

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				}

				dato[indexActividad].(map[string]interface{})["cualitativo"] = cualitativo
				dato[indexActividad].(map[string]interface{})["estado"] = estado
			}

			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str

			estadoSeguimiento = seguimientohelper.GetEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarCualitativo", "err": "Error actualizando componente cualitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			if respuesta["Status"] == "400" {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error", "Data": respuesta["Message"]}
			}
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta["Data"]}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}
	c.ServeJSON()

}

// GuardarCuantitativo ...
// @Title GuardarCuantitativo
// @Description post Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_cuantitativo/:plan_id/:index/:trimestre [post]
func (c *SeguimientoController) GuardarCuantitativo() {
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var resEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var cuantitativo map[string]interface{}
	var estadoSeguimiento string
	observacion := false
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]

			cuantitativo = body["cuantitativo"].(map[string]interface{})
			for _, indicador := range cuantitativo["indicadores"].([]interface{}) {
				if indicador.(map[string]interface{})["observaciones"] != "" && indicador.(map[string]interface{})["observaciones"] != "Sin observación" {
					observacion = true
				}
			}

			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if dato[indexActividad] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
				dato[indexActividad] = map[string]interface{}{"estado": estado, "cuantitativo": cuantitativo}
			} else {
				estado = dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})

				if estado["nombre"] == "Actividad reportada" {
					var codigo_abreviacion string

					if observacion {
						codigo_abreviacion = "CO"
					} else {
						codigo_abreviacion = "AAV"
					}

					fmt.Println(codigo_abreviacion)
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				}

				dato[indexActividad].(map[string]interface{})["cuantitativo"] = cuantitativo
				dato[indexActividad].(map[string]interface{})["estado"] = estado
			}

			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str

			estadoSeguimiento = seguimientohelper.GetEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta["Data"]}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}
	c.ServeJSON()

}

// ReportarActividad ...
// @Title ReportarActividad
// @Description post Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 403
// @router /reportar_actividad/:index [put]
func (c *SeguimientoController) ReportarActividad() {
	indexActividad := c.Ctx.Input.Param(":index")

	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	var body map[string]interface{}
	dato := make(map[string]interface{})

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+body["SeguimientoId"].(string), &respuesta); err == nil {
			aux := make(map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux

			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			reportable, mensaje := seguimientohelper.ActividadReportable(seguimiento, indexActividad)
			if reportable {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &resEstado); err == nil {
					estado := map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
					dato[indexActividad].(map[string]interface{})["estado"] = estado
				}

				b, _ := json.Marshal(dato)
				str := string(b)
				seguimiento["dato"] = str

				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
					panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
				}

				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": seguimiento}
			} else {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error", "Data": mensaje}
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()

}

// ReportarSeguimiento ...
// @Title ReportarSeguimiento
// @Description post Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 403
// @router /reportar_seguimiento [put]
func (c *SeguimientoController) ReportarSeguimiento() {
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	var body map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+body["SeguimientoId"].(string), &respuesta); err == nil {
			aux := make(map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux
			reportable, mensaje := seguimientohelper.ActividadSeguimiento(seguimiento)
			if reportable {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:EAR", &resEstado); err == nil {
					seguimiento["estado_seguimiento_id"] = resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"]
				}

				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
					panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
				}

				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": seguimiento}
			} else {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error", "Data": mensaje}
			}

		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()

}

// RevisarActividad ...
// @Title RevisarActividad
// @Description post Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /revision_actividad/:plan_id/:index/:trimestre [post]
func (c *SeguimientoController) RevisarActividad() {
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	comentario := false

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)

			seguimiento = aux[0]
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			dato[indexActividad] = body

			// Cualitativo
			if body["cualitativo"].(map[string]interface{})["observaciones"] != "" && body["cualitativo"].(map[string]interface{})["observaciones"] != "Sin observación" {
				comentario = true
			}

			// Cuantitativo
			for _, indicador := range body["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{}) {
				if indicador.(map[string]interface{})["observaciones"] != "" && indicador.(map[string]interface{})["observaciones"] != "Sin observación" {
					comentario = true
					break
				}
			}

			// Evidencia
			for _, evidencia := range body["evidencia"].([]interface{}) {
				if evidencia.(map[string]interface{})["Observacion"] != "" && evidencia.(map[string]interface{})["Observacion"] != "Sin observación" {
					comentario = true
					break
				}
			}

			if comentario {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:CO", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
					seguimiento["estado_seguimiento_id"] = estado["id"]
				}
			} else {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AAV", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			}

			dato[indexActividad].(map[string]interface{})["estado"] = estado

			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			}
			data := respuesta["Data"].(map[string]interface{})
			data["Observación"] = comentario
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()

}

// RevisarSeguimiento ...
// @Title RevisarSeguimiento
// @Description post Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /revision_seguimiento/:plan_id/:trimestre [post]
func (c *SeguimientoController) RevisarSeguimiento() {
	planId := c.Ctx.Input.Param(":plan_id")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	avalado := true

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)

			seguimiento = aux[0]
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			for _, actividad := range dato {
				if actividad.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad valada" {
					avalado = false
					break
				}
			}

			if avalado {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AV", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
					seguimiento["estado_seguimiento_id"] = estado["id"]
				}
				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
				}
				data := respuesta["Data"].(map[string]interface{})
				// data["Observación"] = comentario
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
			} else {
			}

		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()

}
