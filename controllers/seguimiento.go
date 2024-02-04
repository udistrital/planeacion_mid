package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
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
	c.Mapping("RetornarActividad", c.RetornarActividad)
	c.Mapping("MigrarInformacion", c.MigrarInformacion)
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

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+entrada["periodo_id"].(string), &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &reportes)
		if len(reportes) > 0 {
			var element = reportes[0]
			element["activo"] = true
			element["fecha_inicio"] = entrada["fecha_inicio"]
			element["fecha_fin"] = entrada["fecha_fin"]
			element["unidades_interes"] = "[]"
			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+element["_id"].(string), "PUT", &resPut, element); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
			}
		} else {
			element := map[string]interface{}{
				"tipo_seguimiento_id": "61f236f525e40c582a0840d0",
				"activo":              true,
				"fecha_inicio":        entrada["fecha_inicio"],
				"fecha_fin":           entrada["fecha_fin"],
				"periodo_id":          entrada["periodo_id"],
				"unidades_interes":    "[]",
			}

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento", "POST", &resPut, element); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle", "status": "400", "log": err})
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+entrada["periodo_id"].(string), &res); err == nil {
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
	var resPadres map[string]interface{}
	var resDependencia map[string]interface{}
	var resTrimestres map[string]interface{}
	var plan map[string]interface{}
	var planesPadre []map[string]interface{}
	var respuestaPost map[string]interface{}
	var arrReportes []map[string]interface{}
	reporte := make(map[string]interface{})
	nuevo := true

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan_id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &plan)
		trimestres := seguimientohelper.GetTrimestres(plan["vigencia"].(string))

		// Caso especial para el plan de acción, retomar avances de seguimiento de versiones anteriores
		if tipo == "61f236f525e40c582a0840d0" && plan["padre_plan_id"] != nil {
			nuevo = false

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=dependencia_id:"+plan["dependencia_id"].(string)+",vigencia:"+plan["vigencia"].(string)+",formato:false,nombre:"+url.QueryEscape(plan["nombre"].(string)), &resPadres); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPadres, &planesPadre)

				var seguimientosLlenos []map[string]interface{}
				var seguimientosVacios []map[string]interface{}

				for _, padre := range planesPadre {
					var resSeguimientos map[string]interface{}
					var seguimientos []map[string]interface{}
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+padre["_id"].(string), &resSeguimientos); err == nil {
						helpers.LimpiezaRespuestaRefactor(resSeguimientos, &seguimientos)
						for _, seguimiento := range seguimientos {
							if (len(seguimientosLlenos) + len(seguimientosVacios)) <= 4 {
								if fmt.Sprintf("%v", seguimiento["dato"]) != "{}" {
									seguimientosLlenos = append(seguimientosLlenos, seguimiento)
								} else {
									seguimientosVacios = append(seguimientosVacios, seguimiento)
								}
							} else {
								break
							}
						}
					}
				}

				var resActualizacion map[string]interface{}
				var resCreacion map[string]interface{}
				var resSeguimientoDetalle map[string]interface{}
				detalle := make(map[string]interface{})
				dato := make(map[string]interface{})
				var resEstado map[string]interface{}
				estado := map[string]interface{}{}

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}

				for _, seguimiento := range seguimientosVacios {
					// ? Inactivar el actual
					seguimiento["activo"] = false
					helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &resActualizacion, seguimiento)
					arrReportes = append(arrReportes, resActualizacion["Data"].(map[string]interface{}))
					// ? Crear el nuevo
					seguimiento["activo"] = true
					seguimiento["plan_id"] = plan_id
					delete(seguimiento, "_id")
					helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &resCreacion, seguimiento)
					arrReportes = append(arrReportes, resCreacion["Data"].(map[string]interface{}))
				}

				for _, seguimiento := range seguimientosLlenos {

					dato = map[string]interface{}{}
					datoStr := seguimiento["dato"].(string)
					json.Unmarshal([]byte(datoStr), &dato)

					listAct := make([]string, 0, len(dato))
					for k := range dato {
						listAct = append(listAct, k)
					}
					for _, idxAct := range listAct {
						id, existe := dato[idxAct].(map[string]interface{})["id"].(string)
						if existe && id != "" {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id, &resSeguimientoDetalle); err == nil {
								helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
								detalle = seguimientohelper.ConvertirStringJson(detalle)
								// ? Inactivar el actual
								detalle["activo"] = false
								seguimientohelper.GuardarDetalleSegimiento(detalle, true) // true => PUT
								// ? crear el nuevo
								detalle["activo"] = true
								detalle["estado"] = estado
								delete(detalle, "_id")
								delete(detalle, "cuantitativo")
								newDetalleId := seguimientohelper.GuardarDetalleSegimiento(detalle, false) // false => POST
								dato[idxAct].(map[string]interface{})["id"] = newDetalleId
							}
						}
					}

					// ? Inactiva el actual
					seguimiento["activo"] = false
					helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &resActualizacion, seguimiento)
					arrReportes = append(arrReportes, resActualizacion["Data"].(map[string]interface{}))

					// ? crear el nuevo
					seguimiento["activo"] = true
					seguimiento["plan_id"] = plan_id
					seguimiento["estado_seguimiento_id"] = "635c11e1e092c5fa5f099971" // En reporte
					valor, _ := json.Marshal(dato)
					str := string(valor)
					seguimiento["dato"] = str
					delete(seguimiento, "_id")
					helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &resCreacion, seguimiento)
					arrReportes = append(arrReportes, resCreacion["Data"].(map[string]interface{}))
				}
			}
		}

		if nuevo {
			for i := 0; i < len(trimestres); i++ {
				periodo := int(trimestres[i]["Id"].(float64))
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?query=tipo_seguimiento_id:`+tipo+`,periodo_id:`+strconv.Itoa(periodo), &resTrimestres); err == nil {
					reporte["nombre"] = "Seguimiento para el " + plan["nombre"].(string)
					reporte["descripcion"] = "Seguimiento " + plan["nombre"].(string)
					if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"dependencia/"+plan["dependencia_id"].(string), &resDependencia); err == nil {
						if resDependencia["Nombre"] != nil {
							reporte["descripcion"] = reporte["descripcion"].(string) + " dependencia " + resDependencia["Nombre"].(string)
						}
					}
					reporte["activo"] = false
					reporte["plan_id"] = plan_id
					reporte["estado_seguimiento_id"] = "61f237df25e40c57a60840d5"
					reporte["periodo_seguimiento_id"] = resTrimestres["Data"].([]interface{})[0].(map[string]interface{})["_id"]
					reporte["fecha_inicio"] = resTrimestres["Data"].([]interface{})[0].(map[string]interface{})["fecha_fin"]
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
	var resSeguimientoDetalle map[string]interface{}
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var seguimiento []map[string]interface{}
	var seguimientoDetalle []map[string]interface{}
	var datoPlan map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,_id:"+seguimiento_id, &resSeguimiento); err == nil {
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
									hayRegistro := false
									if reflect.TypeOf(actividad["index"]).String() == "string" {
										hayRegistro = indexActividad == actividad["index"]
									} else {
										hayRegistro = indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)
									}

									if hayRegistro {
										_, datosUnidos := element.(map[string]interface{})["estado"]
										if datosUnidos {
											actividad["estado"] = element.(map[string]interface{})["estado"]
											break
										} else {
											if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle?query=activo:true,_id:"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
												helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &seguimientoDetalle)
												dato := make(map[string]interface{})
												json.Unmarshal([]byte(seguimientoDetalle[0]["estado"].(string)), &dato)
												actividad["estado"] = dato
												break
											}
										}
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

// GuardarSeguimiento ...
// @Title GuardarSeguimiento
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /guardar_seguimiento/:plan_id/:index/:trimestre [put]
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
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
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
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			}

			if dato[indexActividad] == nil {
				body["estado"] = estado
				delete(body, "_id")
				dato[indexActividad] = map[string]interface{}{"id": seguimientohelper.GuardarDetalleSegimiento(body, false)}

				valor, _ := json.Marshal(dato)
				str := string(valor)
				seguimiento["dato"] = str
			} else {
				id, actualizar := dato[indexActividad].(map[string]interface{})["id"].(string)

				if actualizar && id != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
						detalle = seguimientohelper.ConvertirStringJson(detalle)
						detalle["estado"] = estado
						seguimientohelper.GuardarDetalleSegimiento(detalle, true)
					}
				} else {
					dato[indexActividad].(map[string]interface{})["estado"] = estado

					valor, _ := json.Marshal(dato)
					str := string(valor)
					seguimiento["dato"] = str
				}
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
						priodoId_rest, err := strconv.ParseFloat(test1, 32)
						if err != nil {
							fmt.Println(err)
						}
						periodId = priodoId_rest - 1
					} else {
						test1 = body["periodo_seguimiento_id"].(string)
						priodoId_rest, err := strconv.ParseFloat(test1, 32)
						if err != nil {
							fmt.Println(err)
						}
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
		aPe, err := strconv.ParseFloat(avancePeriodo, 32)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(aPe, err, reflect.TypeOf(avanceAcumulado))
		aAc, err := strconv.ParseFloat(avanceAcumulado, 32)
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

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": generalData}
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

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId, &resSeguimiento); err == nil {
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
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_documentos/:plan_id/:index/:trimestre [put]
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
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	comentario := false

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {

			for _, evidencia := range body["evidencia"].([]interface{}) {
				if evidencia.(map[string]interface{})["Enlace"] != nil {
					evidencias = append(evidencias, evidencia.(map[string]interface{}))
					if evidencia.(map[string]interface{})["Observacion"] != nil && evidencia.(map[string]interface{})["Observacion"] != "Sin observación" && evidencia.(map[string]interface{})["Observacion"] != "" {
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

			if body["unidad"].(bool) {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			} else if comentario {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:CO", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			} else {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			}

			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]

			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)
			if dato[indexActividad] == nil {
				detalle["evidencia"] = evidencias
				detalle["estado"] = estado
				delete(detalle, "_id")
				dato[indexActividad] = map[string]interface{}{"id": seguimientohelper.GuardarDetalleSegimiento(detalle, false)}
			} else {
				id, segregado := dato[indexActividad].(map[string]interface{})["id"]

				if segregado && id != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
						detalle = seguimientohelper.ConvertirStringJson(detalle)
						detalle["evidencia"] = evidencias
						detalle["estado"] = estado
						seguimientohelper.GuardarDetalleSegimiento(detalle, true)
					}
				} else {
					dato[indexActividad].(map[string]interface{})["evidencia"] = evidencias
					dato[indexActividad].(map[string]interface{})["estado"] = estado
				}
			}

			valor, _ := json.Marshal(dato)
			str := string(valor)
			seguimiento["dato"] = str

			estadoSeguimiento = seguimientohelper.GetEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarDocumentos", "err": "Error guardado documentos del seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"seguimiento": detalle["evidencia"], "estadoActividad": estado}}
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
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_cualitativo/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarCualitativo() {
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var resEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var cualitativo map[string]interface{}
	var informacion map[string]interface{}
	var estadoSeguimiento string
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	observacion := false
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]

			cualitativo = body["cualitativo"].(map[string]interface{})
			informacion = body["informacion"].(map[string]interface{})

			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if dato[indexActividad] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}

				detalle = map[string]interface{}{"estado": estado, "cualitativo": cualitativo, "informacion": informacion}
				dato[indexActividad] = map[string]interface{}{"id": seguimientohelper.GuardarDetalleSegimiento(detalle, false)}
			} else {
				id, segregado := dato[indexActividad].(map[string]interface{})["id"]

				if segregado && id != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
						detalle = seguimientohelper.ConvertirStringJson(detalle)
						estado = detalle["estado"].(map[string]interface{})
					}
				} else {
					estado = dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})
				}

				if estado["nombre"] == "Con observaciones" && body["dependencia"].(bool) {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				} else if estado["nombre"] == "Actividad reportada" || estado["nombre"] == "Con observaciones" {
					var codigo_abreviacion string

					observacion = seguimientohelper.ActividadConObservaciones(body)
					if observacion {
						codigo_abreviacion = "CO"
					} else {
						codigo_abreviacion = "AR"
					}

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				}

				if segregado && id != "" {
					detalle["cualitativo"] = cualitativo
					detalle["informacion"] = informacion
					detalle["estado"] = estado
					seguimientohelper.GuardarDetalleSegimiento(detalle, true)
				} else {
					dato[indexActividad].(map[string]interface{})["informacion"] = informacion
					dato[indexActividad].(map[string]interface{})["cualitativo"] = cualitativo
					dato[indexActividad].(map[string]interface{})["estado"] = estado
				}
			}
			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str

			estadoSeguimiento = seguimientohelper.GetEstadoSeguimiento(seguimiento)
			seguimiento["estado_seguimiento_id"] = estadoSeguimiento

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarCualitativo", "err": "Error actualizando componente cualitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
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
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @router /guardar_cuantitativo/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) GuardarCuantitativo() {
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var resEstado map[string]interface{}
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var cuantitativo map[string]interface{}
	var informacion map[string]interface{}
	var estadoSeguimiento string
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	observacion := false
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux[0]

			cuantitativo = body["cuantitativo"].(map[string]interface{})
			informacion = body["informacion"].(map[string]interface{})

			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			if dato[indexActividad] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}

				detalle = map[string]interface{}{"estado": estado, "cuantitativo": cuantitativo, "informacion": informacion}
				dato[indexActividad] = map[string]interface{}{"id": seguimientohelper.GuardarDetalleSegimiento(detalle, false)}
			} else {
				id, segregado := dato[indexActividad].(map[string]interface{})["id"]

				if segregado && id != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
						detalle = seguimientohelper.ConvertirStringJson(detalle)
						estado = detalle["estado"].(map[string]interface{})
					}
				} else {
					estado = dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})
				}

				if estado["nombre"] == "Con observaciones" && body["dependencia"].(bool) {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AER", &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				} else if estado["nombre"] == "Actividad reportada" || estado["nombre"] == "Con observaciones" {
					var codigo_abreviacion string

					observacion = seguimientohelper.ActividadConObservaciones(body)
					if observacion {
						codigo_abreviacion = "CO"
					} else {
						codigo_abreviacion = "AR"
					}

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
				}

				if segregado && id != "" {
					detalle["estado"] = estado
					detalle["cuantitativo"] = cuantitativo
					detalle["informacion"] = informacion
					seguimientohelper.GuardarDetalleSegimiento(detalle, true)
				} else {
					dato[indexActividad].(map[string]interface{})["informacion"] = informacion
					dato[indexActividad].(map[string]interface{})["cuantitativo"] = cuantitativo
					dato[indexActividad].(map[string]interface{})["estado"] = estado
				}
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
// @Description put Seguimiento by id
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403
// @router /reportar_actividad/:index [put]
func (c *SeguimientoController) ReportarActividad() {
	indexActividad := c.Ctx.Input.Param(":index")

	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	var body map[string]interface{}
	var estado map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := map[string]interface{}{}
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
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}

				id, segregable := dato[indexActividad].(map[string]interface{})["id"].(string)
				if segregable && id != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
						detalle = seguimientohelper.ConvertirStringJson(detalle)
						detalle["estado"] = estado

						seguimientohelper.GuardarDetalleSegimiento(detalle, true)
					}
				} else {
					dato[indexActividad].(map[string]interface{})["estado"] = estado

					b, _ := json.Marshal(dato)
					str := string(b)
					seguimiento["dato"] = str

					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
						panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
					}
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
// @Description put Seguimiento by id
// @Param	id			path 	string	true	"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403
// @router /reportar_seguimiento/:id [put]
func (c *SeguimientoController) ReportarSeguimiento() {
	idSeguimiento := c.Ctx.Input.Param(":id")

	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+idSeguimiento, &respuesta); err == nil {
		aux := make(map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux
		reportable, mensaje := seguimientohelper.SeguimientoReportable(seguimiento)
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

	c.ServeJSON()

}

// RevisarActividad ...
// @Title RevisarActividad
// @Description put Seguimiento by id
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /revision_actividad/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) RevisarActividad() {
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	var resDetalle map[string]interface{}
	detalle := map[string]interface{}{}
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
			if body["cualitativo"].(map[string]interface{})["observaciones"] != "" && body["cualitativo"].(map[string]interface{})["observaciones"] != "Sin observación" && body["cualitativo"].(map[string]interface{})["observaciones"] != nil {
				comentario = true
			}

			// Cuantitativo
			for _, indicador := range body["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{}) {
				if indicador.(map[string]interface{})["observaciones"] != "" && indicador.(map[string]interface{})["observaciones"] != "Sin observación" && indicador.(map[string]interface{})["observaciones"] != nil {
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
				}
			} else {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AAV", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			}

			id, segregado := body["id"].(string)
			if segregado && id != "" {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id, &resDetalle); err == nil {
					helpers.LimpiezaRespuestaRefactor(resDetalle, &detalle)
					detalle = seguimientohelper.ConvertirStringJson(detalle)
					detalle["evidencia"] = body["evidencia"]
					detalle["cualitativo"] = body["cualitativo"]
					detalle["cuantitativo"] = body["cuantitativo"]
					detalle["estado"] = estado
					seguimientohelper.GuardarDetalleSegimiento(detalle, true)
				}
			} else {
				dato[indexActividad].(map[string]interface{})["estado"] = estado

				b, _ := json.Marshal(dato)
				str := string(b)
				seguimiento["dato"] = str
			}

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
// @Description put Seguimiento by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :id is empty
// @router /revision_seguimiento/:id [put]
func (c *SeguimientoController) RevisarSeguimiento() {
	seguimientoId := c.Ctx.Input.Param(":id")

	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,_id:"+seguimientoId, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(respuesta, &aux)

		seguimiento = aux[0]
		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &dato)

		avalado, observacion, mensaje := seguimientohelper.SeguimientoAvalable(seguimiento)

		if avalado || observacion {
			var codigo_abreviacion string

			if observacion {
				codigo_abreviacion = "CO"
			} else {
				codigo_abreviacion = "AV"
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:"+codigo_abreviacion, &resEstado); err == nil {
				estado := map[string]interface{}{
					"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
				seguimiento["estado_seguimiento_id"] = estado["id"]
			}

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			}

			data := respuesta["Data"].(map[string]interface{})
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
		} else {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error", "Data": mensaje}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// RetornarActividad ...
// @Title RetornarActividad
// @Description Retorna la actividad de Avalado a en Revision
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403 :plan_id is empty
// @router /retornar_actividad/:plan_id/:index/:trimestre [put]
func (c *SeguimientoController) RetornarActividad() {
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	var resDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
			aux := make([]map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)

			seguimiento = aux[0]
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			id, segregado := body["id"].(string)
			if segregado && id != "" {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id, &resDetalle); err == nil {
					helpers.LimpiezaRespuestaRefactor(resDetalle, &detalle)
					detalle = seguimientohelper.ConvertirStringJson(detalle)
				}

				if detalle["estado"].(map[string]interface{})["id"] == "63793207242b813898e9856b" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:OAPC", &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
					seguimiento["estado_seguimiento_id"] = estado["id"]

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
					detalle["estado"] = estado
					seguimientohelper.GuardarDetalleSegimiento(detalle, true)

					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					}

					data := respuesta["Data"].(map[string]interface{})

					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": nil}
				}
			} else {
				dato[indexActividad] = body

				if dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})["id"] == "63793207242b813898e9856b" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:OAPC", &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
						}
					}
					seguimiento["estado_seguimiento_id"] = estado["id"]

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AR", &resEstado); err == nil {
						estado = map[string]interface{}{
							"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
							"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
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

					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": nil}
				}
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}

// MigrarInformacion ...
// @Title MigrarInformacion
// @Description post Segrar la informacion de los seguimientos
// @Param	plan_id		path 	string	true		"The key for staticblock"
// @Param	trimestre	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403
// @router /migrar_seguimiento/:plan_id/:trimestre [post]
func (c *SeguimientoController) MigrarInformacion() {
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
	trimestre := c.Ctx.Input.Param(":trimestre")

	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	resMigrado := []map[string]interface{}{}
	resNoMigrado := []map[string]interface{}{}
	dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_seguimiento_id:"+trimestre, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(respuesta, &aux)

		seguimiento = aux[0]
		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &dato)

		for indexActividad, actividad := range dato {

			id, segregado := actividad.(map[string]interface{})["id"].(string)
			if !segregado || id == "" {
				delete(dato[indexActividad].(map[string]interface{}), "_id")
				dato[indexActividad] = map[string]interface{}{"id": seguimientohelper.GuardarDetalleSegimiento(dato[indexActividad].(map[string]interface{}), false)}
				resMigrado = append(resMigrado, map[string]interface{}{"id": indexActividad})

			} else {
				resNoMigrado = append(resNoMigrado, map[string]interface{}{"id": indexActividad})
			}
			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
			}
		}

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}

	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": map[string]interface{}{"Actividades migradas:": resMigrado, "Actividades no migradas: ": resNoMigrado}}
	c.ServeJSON()
}

// VerificarActividad ...
// @Title VerificarActividad
// @Description put Seguimiento by id
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Seguimiento
// @Failure 403
// @router /verificar_actividad/:index [put]
func (c *SeguimientoController) VerificarActividad() {
	indexActividad := c.Ctx.Input.Param(":index")

	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var resEstado map[string]interface{}
	var body map[string]interface{}
	var estado map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	dato := make(map[string]interface{})
	observacion := false

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+body["SeguimientoId"].(string), &respuesta); err == nil {
			aux := make(map[string]interface{}, 1)
			helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
			seguimiento = aux

			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			observacion = seguimientohelper.ActividadConObservaciones(body)
			fmt.Println(observacion)
			if !observacion {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:AVV", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}

				id, segregable := dato[indexActividad].(map[string]interface{})["id"].(string)
				if segregable && id != "" {
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
						detalle = seguimientohelper.ConvertirStringJson(detalle)
						detalle["estado"] = estado

						seguimientohelper.GuardarDetalleSegimiento(detalle, true)
					}
				} else {
					dato[indexActividad].(map[string]interface{})["estado"] = estado

					b, _ := json.Marshal(dato)
					str := string(b)
					seguimiento["dato"] = str

					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
						panic(map[string]interface{}{"funcion": "GuardarCuantitativo", "err": "Error actualizando componente cuantitativo de seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
					}
				}

				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": seguimiento}
			} else {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error", "Data": "No se puede verificar una actividad con observaciones"}
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	}
	c.ServeJSON()

}
