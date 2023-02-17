package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"

	"github.com/udistrital/planeacion_mid/helpers"
	formulacionhelper "github.com/udistrital/planeacion_mid/helpers/formulacionHelper"
	inversionhelper "github.com/udistrital/planeacion_mid/helpers/inversionHelper"
	"github.com/udistrital/utils_oas/request"
)

// InversionController operations for Inversion
type InversionController struct {
	beego.Controller
}

// URLMapping ...
func (c *InversionController) URLMapping() {
	c.Mapping("AddProyecto", c.AddProyecto)
	c.Mapping("EditProyecto", c.EditProyecto)
	c.Mapping("GuardarDocumentos", c.GuardarDocumentos)
	c.Mapping("GetProyectoId", c.GetProyectoId)
	c.Mapping("GetMetasProyect", c.GetMetasProyect)
	c.Mapping("GetAllProyectos", c.GetAllProyectos)
	c.Mapping("ActualizarSubgrupoDetalle", c.ActualizarSubgrupoDetalle)
	c.Mapping("ActualizarProyectoGeneral", c.ActualizarProyectoGeneral)
	c.Mapping("CrearPlan", c.CrearPlan)
	c.Mapping("GetPlanId", c.GetPlanId)
	c.Mapping("GuardarActividad", c.GuardarActividad)
}

// AddProyecto ...
// @Title AddProyecto
// @Description post AddProyecto
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @router /proyecto [post]
func (c *InversionController) AddProyecto() {
	var registroProyecto map[string]interface{}
	var idProyecto string
	var resPlan map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &registroProyecto); err == nil {
		plan := map[string]interface{}{
			"activo":        true,
			"nombre":        registroProyecto["nombre_proyecto"],
			"descripcion":   registroProyecto["codigo_proyecto"],
			"tipo_plan_id":  "63ca86f1b6c0e5725a977dae",
			"aplicativo_id": " ",
		}
		var respuesta map[string]interface{}

		err1 := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan", "POST", &respuesta, plan)
		if err1 == nil {
			resPlan = respuesta["Data"].(map[string]interface{})
			idProyecto = resPlan["_id"].(string)

			soportes := map[string]interface{}{"codigo_proyecto": registroProyecto["codigo_proyecto"], "data": registroProyecto["soportes"]}
			errSoporte := inversionhelper.ResgistrarInfoComplementaria(idProyecto, soportes, "soportes")

			fuentes := map[string]interface{}{"codigo_proyecto": registroProyecto["codigo_proyecto"], "data": registroProyecto["fuentes"]}
			errFuentes := inversionhelper.ResgistrarInfoComplementaria(idProyecto, fuentes, "fuentes apropiacion")
			inversionhelper.ActualizarPresupuestoDisponible(registroProyecto["fuentes"].([]interface{}))

			metas := map[string]interface{}{"codigo_proyecto": registroProyecto["codigo_proyecto"], "data": registroProyecto["metas"]}
			errMetas := inversionhelper.ResgistrarInfoComplementaria(idProyecto, metas, "metas asociadas al proyecto de inversion")

			if errSoporte != nil || errFuentes != nil || errMetas != nil {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "error", "Data": errSoporte}
				c.Abort("400")
			}
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": resPlan}
			c.ServeJSON()
		}

	} else {
		panic(map[string]interface{}{"funcion": "AddProyecto", "err": "Error Registrando Proyecto", "status": "400", "log": err})
	}
}

// EditProyecto ...
// @Title EditProyecto
// @Description post EditProyecto
// @Param	body		body 	{}	true		"body for Plan content"
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /proyecto/:id [put]
func (c *InversionController) EditProyecto() {
	id := c.Ctx.Input.Param(":id")
	var registroProyecto map[string]interface{}
	var res map[string]interface{}
	var infoProyect map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &registroProyecto); err == nil {

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &infoProyect)
			infoProyect["nombre"] = registroProyecto["nombre_proyecto"]
			infoProyect["descripcion"] = registroProyecto["codigo_proyecto"]

			var respuesta map[string]interface{}
			err1 := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, "PUT", &respuesta, infoProyect)
			if err1 == nil {
				errSoporte := inversionhelper.ActualizarInfoComplDetalle(registroProyecto["id_detalle_soportes"].(string), registroProyecto["soportes"].([]interface{}))

				errFuentes := inversionhelper.ActualizarInfoComplDetalle(registroProyecto["id_detalle_fuentes"].(string), registroProyecto["fuentes"].([]interface{}))
				inversionhelper.ActualizarPresupuestoDisponible(registroProyecto["fuentes"].([]interface{}))

				errMetas := inversionhelper.ActualizarInfoComplDetalle(registroProyecto["id_detalle_metas"].(string), registroProyecto["metas"].([]interface{}))

				if errSoporte != nil || errFuentes != nil || errMetas != nil {
					c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "error", "Data": errSoporte}
					c.Abort("400")
				}
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": infoProyect}
				c.ServeJSON()
			}
		} else {
			panic(map[string]interface{}{"funcion": "GetProyectoId", "err": "Error obteniendo información plan", "status": "400", "log": err})
		}

	} else {
		panic(map[string]interface{}{"funcion": "AddProyecto", "err": "Error Registrando Proyecto", "status": "400", "log": err})
	}
}

// GuardarDocumentos ...
// @Title GuardarDocumentos
// @Description post AddProyecto
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /guardar_documentos [post]
func (c *InversionController) GuardarDocumentos() {
	var body map[string]interface{}
	var evidencias []map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err == nil {
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
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": evidencias}
			c.ServeJSON()
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	}
}

// GetProyectoId ...
// @Title GetProyectoId
// @Description get GetProyectoId
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /proyecto/:id [get]
func (c *InversionController) GetProyectoId() {
	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	getProyect := make(map[string]interface{})
	var infoProyect map[string]interface{}
	var subgruposData map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &infoProyect)
		getProyect["nombre_proyecto"] = infoProyect["nombre"]
		getProyect["codigo_proyecto"] = infoProyect["descripcion"]
		getProyect["fecha_creacion"] = infoProyect["fecha_creacion"]
		padreId := infoProyect["_id"].(string)

		var infoSubgrupos []map[string]interface{}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+padreId, &subgruposData); err == nil {
			helpers.LimpiezaRespuestaRefactor(subgruposData, &infoSubgrupos)
			for i := range infoSubgrupos {
				var subgrupoDetalle map[string]interface{}
				var detalleSubgrupos []map[string]interface{}

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
					helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)

					armonizacion_dato_str := detalleSubgrupos[0]["dato"].(string)
					var subgrupo_dato []map[string]interface{}
					json.Unmarshal([]byte(armonizacion_dato_str), &subgrupo_dato)

					if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "soporte") {
						getProyect["soportes"] = subgrupo_dato
						getProyect["id_detalle_soportes"] = detalleSubgrupos[0]["_id"]
					}
					if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "metas") {
						getProyect["metas"] = subgrupo_dato
						getProyect["id_detalle_metas"] = detalleSubgrupos[0]["_id"]
					}
					if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "fuentes") {
						getProyect["fuentes"] = subgrupo_dato
						getProyect["id_detalle_fuentes"] = detalleSubgrupos[0]["_id"]
					}
				}
			}
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": getProyect}
	} else {
		panic(map[string]interface{}{"funcion": "GetProyectoId", "err": "Error obteniendo información plan", "status": "400", "log": err})
	}
	c.ServeJSON()
}

// GetMetasProyect ...
// @Title GetMetasProyect
// @Description get GetMetasProyect
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /metaspro/:id [get]
func (c *InversionController) GetMetasProyect() {

	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	getProyect := make(map[string]interface{})
	var infoProyect map[string]interface{}
	var subgruposData map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res); err == nil {
		//fmt.Println(res, "respuesta")
		helpers.LimpiezaRespuestaRefactor(res, &infoProyect)
		//getProyect["nombre_proyecto"] = infoProyect["nombre"]
		//getProyect["codigo_proyecto"] = infoProyect["descripcion"]
		//getProyect["fecha_creacion"] = infoProyect["fecha_creacion"]
		//fmt.Println(infoProyect, "respuesta")
		padreId := infoProyect["_id"].(string)

		var infoSubgrupos []map[string]interface{}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+padreId, &subgruposData); err == nil {
			helpers.LimpiezaRespuestaRefactor(subgruposData, &infoSubgrupos)
			for i := range infoSubgrupos {
				var subgrupoDetalle map[string]interface{}
				var detalleSubgrupos []map[string]interface{}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "metas") {
					getProyect["subgrupo_id_metas"] = infoSubgrupos[i]["_id"]
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
						var dato_metas []map[string]interface{}
						datoMeta_str := detalleSubgrupos[0]["dato"].(string)
						fmt.Println(datoMeta_str, "datoMetas_str")
						json.Unmarshal([]byte(datoMeta_str), &dato_metas)
						fmt.Println(dato_metas, "datoMetas")
						getProyect["metas"] = dato_metas
						//getProyect["id_detalle_metas"] = detalleSubgrupos[0]["_id"]

					}
				}
			}
			//fmt.Println(getProyect)
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": getProyect["metas"]}
	} else {
		panic(map[string]interface{}{"funcion": "GetMetasProyect", "err": "Error obteniendo información Metas Proyecto Inversión", "status": "400", "log": err})
	}
	c.ServeJSON()
}

// GetAllProyectos ...
// @Title GetAllProyectos
// @Description get GetAllProyectos
// @Param	aplicativo_id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :aplicativo_id is empty
// @router /proyectos/:aplicativo_id [get]
func (c *InversionController) GetAllProyectos() {
	tipo_plan_id := c.Ctx.Input.Param(":aplicativo_id")

	var res map[string]interface{}
	var getProyect []map[string]interface{}
	var proyecto map[string]interface{}
	var dataProyects []map[string]interface{}

	proyect := make(map[string]interface{})
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+tipo_plan_id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &dataProyects)
		for i := range dataProyects {
			if dataProyects[i]["activo"] == true {
				proyect["id"] = dataProyects[i]["_id"]
				proyecto = inversionhelper.GetDataProyects(dataProyects[i])
			}
			getProyect = append(getProyect, proyecto)
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": getProyect}
	} else {
		panic(map[string]interface{}{"funcion": "GetProyectoId", "err": "Error obteniendo información plan", "status": "400", "log": err})
	}
	c.ServeJSON()
}

// ActualizarSubgrupoDetalle ...
// @Title ActualizarSubgrupoDetalle
// @Description put Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /actualiza_sub_detalle/:id [put]
func (c *InversionController) ActualizarSubgrupoDetalle() {
	var subDetalle map[string]interface{}
	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &subDetalle)
	if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+id, "PUT", &res, subDetalle); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Update Subgrupo Detalle Successful", "Data": res}
		c.ServeJSON()
	} else {
		panic(map[string]interface{}{"funcion": "ActualizarSubgrupoDetalle", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
	}

}

// ActualizarProyectoGeneral ...
// @Title ActualizarProyectoGeneral
// @Description put Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /actualizar_proyecto/:id [put]
func (c *InversionController) ActualizarProyectoGeneral() {
	var infoProyecto map[string]interface{}
	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &infoProyecto)
	if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, "PUT", &res, infoProyecto); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Update Plan Successful", "Data": res}
		c.ServeJSON()
	} else {
		panic(map[string]interface{}{"funcion": "ActualizarProyectoGeneral", "err": "Error actualizando plan \"id\"", "status": "400", "log": err})
	}

}

// CrearPlan ...
// @Title CrearPlan
// @Description post CrearPlan
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403
// @router /crearplan [post]
func (c *InversionController) CrearPlan() {

	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "FormulacionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	//id := c.Ctx.Input.Param(":id")

	var respuesta map[string]interface{}
	//var respuestaHijos map[string]interface{}
	//var hijos []map[string]interface{}
	var planFormato map[string]interface{}
	var parametros map[string]interface{}
	var respuestaPost map[string]interface{}
	var planSubgrupo map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &parametros)
	plan := make(map[string]interface{})
	var respuestaHijos map[string]interface{}
	var hijos []map[string]interface{}
	//subgrupoMetas := make(map[string]interface{})
	//subgrupoArmonizacion := make(map[string]interface{})
	//subgrupoActividades := make(map[string]interface{})
	//var infoSubgrupos []map[string]interface{}
	//subDetalleMetas := make(map[string]interface{})
	//subDetalleArmonizacionPDD := make(map[string]interface{})
	//subDetalleArmonizacionPED := make(map[string]interface{})
	//subDetalleArmonizacion := make(map[string]interface{})
	//subDetalleActividades := make(map[string]interface{})
	//var infoSubDetalles []map[string]interface{}
	//var idSubMetas string
	//var idSubArmonizacion string
	//var idSubActividades string
	//clienteHttp := &http.Client{}
	//url := "http://" + beego.AppConfig.String("PlanesService") + "/plan/"
	id := parametros["id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &respuesta); err == nil {

		helpers.LimpiezaRespuestaRefactor(respuesta, &planFormato)

		plan["nombre"] = "" + planFormato["nombre"].(string)
		plan["descripcion"] = planFormato["descripcion"].(string)
		plan["tipo_plan_id"] = planFormato["tipo_plan_id"].(string)
		plan["aplicativo_id"] = planFormato["aplicativo_id"].(string)
		plan["activo"] = planFormato["activo"]
		plan["formato"] = false
		plan["vigencia"] = parametros["vigencia"].(string)
		plan["dependencia_id"] = parametros["dependencia_id"].(string)
		plan["estado_plan_id"] = "614d3ad301c7a200482fabfd"

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/", "POST", &respuestaPost, plan); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaPost, &planSubgrupo)
			padre := planSubgrupo["_id"].(string)
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &respuestaHijos); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuestaHijos, &hijos)
				formulacionhelper.ClonarHijos(hijos, padre)
			}
			// subgrupoMetas["padre"] = planSubgrupo["_id"]
			// subgrupoMetas["activo"] = true
			// subgrupoMetas["nombre"] = "Metas Planes Inversión"
			// subgrupoMetas["descripcion"] = "Plan de Acción de Inversión"
			// infoSubgrupos = append(infoSubgrupos, subgrupoMetas)
			// subgrupoArmonizacion["padre"] = planSubgrupo["_id"]
			// subgrupoArmonizacion["activo"] = true
			// subgrupoArmonizacion["nombre"] = "Armonización Planes Inversión"
			// subgrupoArmonizacion["descripcion"] = "Plan de Acción de Inversión"
			// infoSubgrupos = append(infoSubgrupos, subgrupoArmonizacion)
			// subgrupoActividades["padre"] = planSubgrupo["_id"]
			// subgrupoActividades["activo"] = true
			// subgrupoActividades["nombre"] = "Actividades Planes Inversión"
			// subgrupoActividades["descripcion"] = "Plan de Acción de Inversión"
			// infoSubgrupos = append(infoSubgrupos, subgrupoActividades)

			// fmt.Println(infoSubgrupos, "subgrupos")
			// for _, subgrupo := range infoSubgrupos {
			// 	var resSubgrupo map[string]interface{}
			// 	var subgrupos map[string]interface{}

			// 	if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo", "POST", &resSubgrupo, subgrupo); err == nil {
			// 		helpers.LimpiezaRespuestaRefactor(resSubgrupo, &subgrupos)

			// 		if strings.Contains(strings.ToLower(subgrupos["nombre"].(string)), "metas") {
			// 			idSubMetas = subgrupos["_id"].(string)
			// 		}
			// 		if strings.Contains(strings.ToLower(subgrupos["nombre"].(string)), "armonización") {
			// 			idSubArmonizacion = subgrupos["_id"].(string)
			// 		}
			// 		if strings.Contains(strings.ToLower(subgrupos["nombre"].(string)), "actividades") {
			// 			idSubActividades = subgrupos["_id"].(string)
			// 		}

			// 		fmt.Println(subgrupos, "respuesta sub")
			// 	} else {
			// 		panic(map[string]interface{}{"funcion": "CrearPlan", "err": "Error Registrando Subgrupo", "status": "400", "log": err})
			// 	}
			// }
			// subDetalleMetas["nombre"] = "Metas Planes Inversión"
			// subDetalleMetas["descripcion"] = "Plan de Acción de Inversión"
			// subDetalleMetas["subgrupo_id"] = idSubMetas
			// subDetalleMetas["dato"] = "{}"
			// subDetalleMetas["dato_plan"] = ""
			// subDetalleMetas["armonizacion_dato"] = ""
			// infoSubDetalles = append(infoSubDetalles, subDetalleMetas)
			// subDetalleArmonizacion["nombre"] = "Armonización PDD"
			// subDetalleArmonizacion["descripcion"] = "Plan de Acción de Inversión"
			// subDetalleArmonizacion["subgrupo_id"] = idSubArmonizacion
			// subDetalleArmonizacion["dato"] = "{}"
			// subDetalleArmonizacion["dato_plan"] = ""
			// subDetalleArmonizacion["armonizacion_dato"] = ""
			// infoSubDetalles = append(infoSubDetalles, subDetalleArmonizacion)
			// // subDetalleArmonizacionPED["nombre"] = "Armonización PED"
			// // subDetalleArmonizacionPED["descripcion"] = "Plan de Acción de Inversión"
			// // subDetalleArmonizacionPED["subgrupo_id"] = idSubArmonizacion
			// // subDetalleArmonizacionPED["dato"] = "{}"
			// // subDetalleArmonizacionPED["dato_plan"] = ""
			// // subDetalleArmonizacionPED["armonizacion_dato"] = ""
			// // infoSubDetalles = append(infoSubDetalles, subDetalleArmonizacionPED)
			// // subDetalleArmonizacionPI["nombre"] = "Armonización PI"
			// // subDetalleArmonizacionPI["descripcion"] = "Plan de Acción de Inversión"
			// // subDetalleArmonizacionPI["subgrupo_id"] = idSubArmonizacion
			// // subDetalleArmonizacionPI["dato"] = "{}"
			// // subDetalleArmonizacionPI["dato_plan"] = ""
			// // subDetalleArmonizacionPI["armonizacion_dato"] = ""
			// //infoSubDetalles = append(infoSubDetalles, subDetalleArmonizacionPI)
			// subDetalleActividades["nombre"] = "Actividades Plan de Inversión"
			// subDetalleActividades["descripcion"] = "Plan de Acción de Inversión"
			// subDetalleActividades["subgrupo_id"] = idSubActividades
			// subDetalleActividades["dato"] = "{}"
			// subDetalleActividades["dato_plan"] = ""
			// subDetalleActividades["armonizacion_dato"] = ""
			// infoSubDetalles = append(infoSubDetalles, subDetalleActividades)

			// for _, detalle := range infoSubDetalles {
			// 	var resSubDetalle map[string]interface{}
			// 	var subDetalle map[string]interface{}
			// 	if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/", "POST", &resSubDetalle, detalle); err == nil {
			// 		helpers.LimpiezaRespuestaRefactor(resSubDetalle, &subDetalle)
			// 		fmt.Println(subDetalle, "resDetalle")
			// 	} else {
			// 		panic(map[string]interface{}{"funcion": "CrearPlan", "err": "Error registrando subgrupo detalle", "status": "400", "log": err})

			// 	}
			// }

		} else {
			panic(map[string]interface{}{"funcion": "CrearPlan", "err": "Error creando plan", "status": "400", "log": err})
		}

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Update Plan Successful", "Data": planSubgrupo}
		c.ServeJSON()
	} else {
		panic(map[string]interface{}{"funcion": "CrearPlan", "err": "Error consultando datos Plan Formato", "status": "400", "log": err})
	}

	//var resPost map[string]interface{}
	//var resLimpia map[string]interface{}

	// aux, err := json.Marshal(plan)
	// if err != nil {
	// 	log.Fatalf("Error codificado: %v", err)
	// }

	// peticion, err := http.NewRequest("POST", url, bytes.NewBuffer(aux))
	// if err != nil {
	// 	log.Fatalf("Error creando peticion: %v", err)
	// }
	// peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
	// respuesta, err := clienteHttp.Do(peticion)
	// if err != nil {
	// 	log.Fatalf("Error haciendo peticion: %v", err)
	// }

	// defer respuesta.Body.Close()

	// cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
	// if err != nil {
	// 	log.Fatalf("Error leyendo peticion: %v", err)
	// }

	// json.Unmarshal(cuerpoRespuesta, &resPost)
	// resLimpia = resPost["Data"].(map[string]interface{})
	// padre := resLimpia["_id"].(string)
	// c.Data["json"] = resPost

	// if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &respuestaHijos); err == nil {
	// 	helpers.LimpiezaRespuestaRefactor(respuestaHijos, &hijos)
	// 	formulacionhelper.ClonarHijos(hijos, padre)
	// }

}

// GetPlanId ...
// @Title GetPlanId
// @Description get GetPlanId
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /get_plan/:id [get]
func (c *InversionController) GetPlanId() {
	fmt.Println("llega a la función")
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "InversionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	getProyect := make(map[string]interface{})
	var infoProyect map[string]interface{}
	//var subgruposData map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &infoProyect)
		getProyect["nombre_proyecto"] = infoProyect["nombre"]
		getProyect["codigo_proyecto"] = infoProyect["descripcion"]
		getProyect["fecha_creacion"] = infoProyect["fecha_creacion"]
		// padreId := infoProyect["_id"].(string)

		// var infoSubgrupos []map[string]interface{}
		// if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+padreId, &subgruposData); err == nil {
		// 	helpers.LimpiezaRespuestaRefactor(subgruposData, &infoSubgrupos)
		// 	for i := range infoSubgrupos {
		// 		var subgrupoDetalle map[string]interface{}
		// 		var detalleSubgrupos []map[string]interface{}
		// 		if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "soporte") {
		// 			getProyect["subgrupo_id_soportes"] = infoSubgrupos[i]["_id"]
		// 			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
		// 				helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
		// 				getProyect["soportes"] = detalleSubgrupos[0]["dato"]
		// 				getProyect["id_detalle_soportes"] = detalleSubgrupos[0]["_id"]
		// 			}
		// 		}
		// 		if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "metas") {
		// 			getProyect["subgrupo_id_metas"] = infoSubgrupos[i]["_id"]
		// 			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
		// 				helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
		// 				getProyect["metas"] = detalleSubgrupos[0]["dato"]
		// 				getProyect["id_detalle_metas"] = detalleSubgrupos[0]["_id"]

		// 			}
		// 		}
		// 		if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "fuentes") {
		// 			getProyect["subgrupo_id_fuentes"] = infoSubgrupos[i]["_id"]
		// 			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
		// 				helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
		// 				getProyect["fuentes"] = detalleSubgrupos[0]["dato"]
		// 				getProyect["id_detalle_fuentes"] = detalleSubgrupos[0]["_id"]

		// 			}
		// 		}
		// 	}
		// }
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": getProyect}
	} else {
		panic(map[string]interface{}{"funcion": "GetProyectoId", "err": "Error obteniendo información plan", "status": "400", "log": err})
	}
	c.ServeJSON()
}

// GuardarActividad ...
// @Title GuardarActividad
// @Description put Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /guardar_actividad/:id [put]
func (c *InversionController) GuardarActividad() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "InversionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()
	id := c.Ctx.Input.Param(":id")
	var body map[string]interface{}
	var res map[string]interface{}
	var entrada map[string]interface{}
	var resPlan map[string]interface{}
	var plan map[string]interface{}
	var armonizacionExecuted bool = false

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	entrada = body
	armonizacion := body["armo"]
	armonizacionPI := body["armoPI"]
	maxIndex := formulacionhelper.GetIndexActividad(entrada)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &resPlan); err != nil {
		panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get Plan \"id\"", "status": "400", "log": err})
	}
	helpers.LimpiezaRespuestaRefactor(resPlan, &plan)
	if plan["estado_plan_id"] != "614d3ad301c7a200482fabfd" {
		var res map[string]interface{}
		plan["estado_plan_id"] = "614d3ad301c7a200482fabfd"
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, "PUT", &res, plan); err != nil {
			panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizacion estado \"id\"", "status": "400", "log": err})
		}
	}

	for key, element := range entrada {

		var respuesta map[string]interface{}
		var respuestaLimpia []map[string]interface{}
		var subgrupo_detalle map[string]interface{}
		dato_plan := make(map[string]interface{})

		if element != "" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+key, &respuesta); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
			}
			fmt.Println(key, "key")
			helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
			subgrupo_detalle = respuestaLimpia[0]
			actividad := make(map[string]interface{})

			if subgrupo_detalle["dato_plan"] == nil {
				actividad["index"] = 1
				actividad["dato"] = element
				actividad["activo"] = true
				i := strconv.Itoa(actividad["index"].(int))
				dato_plan[i] = actividad

				b, _ := json.Marshal(dato_plan)
				str := string(b)
				subgrupo_detalle["dato_plan"] = str
				if !armonizacionExecuted {
					armonizacion_dato := make(map[string]interface{})
					aux := make(map[string]interface{})
					aux["armonizacionPED"] = armonizacion
					aux["armonizacionPI"] = armonizacionPI
					armonizacion_dato[i] = aux
					c, _ := json.Marshal(armonizacion_dato)
					strArmonizacion := string(c)
					subgrupo_detalle["armonizacion_dato"] = strArmonizacion
					armonizacionExecuted = true
				}
			} else {
				dato_plan_str := subgrupo_detalle["dato_plan"].(string)
				json.Unmarshal([]byte(dato_plan_str), &dato_plan)

				actividad["index"] = maxIndex + 1
				actividad["dato"] = element
				actividad["activo"] = true
				i := strconv.Itoa(actividad["index"].(int))
				dato_plan[i] = actividad
				b, _ := json.Marshal(dato_plan)
				str := string(b)
				subgrupo_detalle["dato_plan"] = str

				if !armonizacionExecuted {
					armonizacion_dato := make(map[string]interface{})
					if subgrupo_detalle["armonizacion_dato"] != nil {
						armonizacion_dato_str := subgrupo_detalle["armonizacion_dato"].(string)
						json.Unmarshal([]byte(armonizacion_dato_str), &armonizacion_dato)
					}
					aux := make(map[string]interface{})
					aux["armonizacionPED"] = armonizacion
					aux["armonizacionPI"] = armonizacionPI
					armonizacion_dato[i] = aux
					c, _ := json.Marshal(armonizacion_dato)
					strArmonizacion := string(c)
					subgrupo_detalle["armonizacion_dato"] = strArmonizacion
					armonizacionExecuted = true

				}
			}
			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
			}
		}
	}
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": entrada}
	c.ServeJSON()

}

// GetPlan ...
// @Title GetPlan
// @Description get Plan by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /get_plan/:id/:index [get]
func (c *InversionController) GetPlan() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "FormulacionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	id := c.Ctx.Input.Param(":id")
	index := c.Ctx.Input.Param(":index")
	var res map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		formulacionhelper.Limpia()
		tree := formulacionhelper.BuildTreeFa(hijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": tree}
	} else {
		panic(err)

	}

	c.ServeJSON()
}
