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
}

// HabilitarReportes ...
// @Title HabilitarReportes
// @Description get Seguimiento
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403
// @router /habilitar_reportes [put]
func (c *SeguimientoController) HabilitarReportes() {
	var res map[string]interface{}
	var resPut map[string]interface{}
	var reportes []map[string]interface{}
	var entrada map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=periodo_id:"+entrada["periodo_id"].(string), &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &reportes)
		for key, element := range reportes {
			_ = key
			element["activo"] = true
			element["fecha_inicio"] = entrada["fecha_inicio"]
			element["fecha_fin"] = entrada["fecha_fin"]

			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+element["_id"].(string), "PUT", &resPut, element); err != nil {
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

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+plan_id, &res); err == nil {
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

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento", "POST", &respuestaPost, reporte); err != nil {
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
	if trimestres[0]["Id"] == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}

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
// @router /get_actividades/:plan_id [get]
func (c *SeguimientoController) GetActividadesGenerales() {
	plan_id := c.Ctx.Input.Param(":plan_id")
	var res map[string]interface{}
	var subgrupos []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
				actividades := seguimientohelper.GetActividades(subgrupos[i]["_id"].(string))
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": actividades}

				break
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
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
	plan_id := c.Ctx.Input.Param(":plan_id")
	index := c.Ctx.Input.Param(":index")
	var res map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		data := seguimientohelper.GetDataSubgrupos(hijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
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
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	dato := make(map[string]interface{})

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_id:"+trimestre, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux[0]
		if seguimiento["dato"] == "{}" {
			dato[indexActividad] = body
			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str
		} else {
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)

			dato[indexActividad] = body
			b, _ := json.Marshal(dato)
			str := string(b)
			seguimiento["dato"] = str
		}
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento/"+seguimiento["_id"].(string), "PUT", &respuesta, seguimiento); err != nil {
			panic(map[string]interface{}{"funcion": "GuardarSeguimiento", "err": "Error actualizando seguimiento \"seguimiento[\"_id\"].(string)\"", "status": "400", "log": err})
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta["Data"]}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
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
	planId := c.Ctx.Input.Param(":plan_id")
	indexActividad := c.Ctx.Input.Param(":index")
	trimestre := c.Ctx.Input.Param(":trimestre")
	var respuesta map[string]interface{}
	var seguimiento map[string]interface{}
	var seguimientoActividad map[string]interface{}
	dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+planId+",periodo_id:"+trimestre, &respuesta); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
		seguimiento = aux[0]

		if seguimiento["dato"] == "{}" {
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}
		} else {
			datoStr := seguimiento["dato"].(string)
			json.Unmarshal([]byte(datoStr), &dato)
			seguimientoActividad = dato[indexActividad].(map[string]interface{})
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": seguimientoActividad}

		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
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
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					c.Abort("400")
				}
				break
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
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

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_id:"+body["periodo_id"].(string), &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &avancedata)
		fmt.Println(res)
		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+body["periodo_id"].(string), &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &parametro_periodo)
			paramIdlen := parametro_periodo[0]
			fmt.Println()
			paramId := paramIdlen["ParametroId"].(map[string]interface{})
			if paramId["CodigoAbreviacion"] != "T1" {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_id:"+body["periodo_id"].(string), &res1); err == nil {
					helpers.LimpiezaRespuestaRefactor(res1, &avancedata1)
					seguimiento1 = avancedata1[0]
					datoStrUltimoTrimestre := seguimiento1["dato"].(string)
					if datoStrUltimoTrimestre == "{}" {
						test1 = body["periodo_id"].(string)
						priodoId_rest, err := strconv.ParseFloat(test1, 8)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(test1)
						periodId = priodoId_rest - 1
					} else {
						test1 = body["periodo_id"].(string)
						priodoId_rest, err := strconv.ParseFloat(test1, 8)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(test1)
						periodId = priodoId_rest
					}

				} else {
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					c.Abort("400")
				}
				periodIdString = fmt.Sprint(periodId)
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_id:"+periodIdString, &res2); err == nil {
					helpers.LimpiezaRespuestaRefactor(res2, &avancedata2)
					seguimiento = avancedata2[0]
					datoStr := seguimiento["dato"].(string)
					json.Unmarshal([]byte(datoStr), &dato)
					indicador1 := dato[body["index"].(string)].(map[string]interface{})
					fmt.Println(len(indicador1))
					avanceIndicador1 := indicador1[body["Nombre_del_indicador"].(string)].(map[string]interface{})
					avanceAcumulado = avanceIndicador1["avanceAcumulado"].(string)
					testavancePeriodo = avanceIndicador1["avancePeriodo"].(string)
				} else {
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					c.Abort("400")
				}
				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+body["periodo_id"].(string), &resName); err == nil {
					helpers.LimpiezaRespuestaRefactor(resName, &parametro_periodo_name)
					paramIdlenName := parametro_periodo_name[0]
					fmt.Println()
					paramIdName := paramIdlenName["ParametroId"].(map[string]interface{})
					nombrePeriodo = paramIdName["CodigoAbreviacion"].(string)
				} else {
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					c.Abort("400")
				}
			} else {
				fmt.Println("")
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
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
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}
