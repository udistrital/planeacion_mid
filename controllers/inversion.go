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
	c.Mapping("GetPlan", c.GetPlan)
	c.Mapping("GuardarMeta", c.GuardarMeta)
	c.Mapping("ArmonizarInversion", c.ArmonizarInversion)
	c.Mapping("ActualizarMetaPlan", c.ActualizarMetaPlan)
	c.Mapping("AllMetasPlan", c.AllMetasPlan)

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

		helpers.LimpiezaRespuestaRefactor(res, &infoProyect)
		//getProyect["nombre_proyecto"] = infoProyect["nombre"]
		//getProyect["codigo_proyecto"] = infoProyect["descripcion"]
		//getProyect["fecha_creacion"] = infoProyect["fecha_creacion"]
		//fmt.Println(infoProyect, "respuesta")
		//padreId := infoProyect["_id"].(string)
		//padreId := id
		var infoSubgrupos []map[string]interface{}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+id, &subgruposData); err == nil {
			fmt.Println(subgruposData, "respuesta")
			helpers.LimpiezaRespuestaRefactor(subgruposData, &infoSubgrupos)
			for i := range infoSubgrupos {
				var subgrupoDetalle map[string]interface{}
				var detalleSubgrupos []map[string]interface{}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "metas") {
					getProyect["subgrupo_id_metas"] = infoSubgrupos[i]["_id"]
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
						var dato_metas []map[string]interface{}
						getProyect["id_detalle_meta"] = detalleSubgrupos[0]["_id"]
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
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": getProyect}
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
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "InversionController" + "/" + (localError["funcion"]).(string))
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

// GuardarMeta ...
// @Title GuardarMeta
// @Description put Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /guardar_meta/:id [put]
func (c *InversionController) GuardarMeta() {
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
	//var dataProyectIn bool = false
	var respuestaGuardado map[string]interface{}

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	entrada = body["entrada"].(map[string]interface{})
	idSubDetalleProI := body["idSubDetalle"]
	indexMetaSubProI := body["indexMetaSubPro"]
	maxIndex := formulacionhelper.GetIndexActividad(entrada)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &resPlan); err != nil {
		panic(map[string]interface{}{"funcion": "GuardarMeta", "err": "Error get Plan \"id\"", "status": "400", "log": err})
	}
	helpers.LimpiezaRespuestaRefactor(resPlan, &plan)
	if plan["estado_plan_id"] != "614d3ad301c7a200482fabfd" {
		var res map[string]interface{}
		plan["estado_plan_id"] = "614d3ad301c7a200482fabfd"
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, "PUT", &res, plan); err != nil {
			panic(map[string]interface{}{"funcion": "GuardarMeta", "err": "Error actualizacion estado \"id\"", "status": "400", "log": err})
		}
	}

	for key, element := range entrada {

		var respuesta map[string]interface{}
		var respuestaLimpia []map[string]interface{}
		var subgrupo_detalle map[string]interface{}
		dato_plan := make(map[string]interface{})

		if element != "" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+key, &respuesta); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarMeta", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
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
				fmt.Sprintln(subgrupo_detalle["dato_plan"], "dato plan")
				//if !dataProyectIn {
				armonizacion_dato := make(map[string]interface{})
				aux := make(map[string]interface{})
				aux["idSubDetalleProI"] = idSubDetalleProI
				aux["indexMetaSubProI"] = indexMetaSubProI
				aux["indexMetaPlan"] = 1
				armonizacion_dato[i] = aux
				c, _ := json.Marshal(armonizacion_dato)
				strArmonizacion := string(c)
				subgrupo_detalle["armonizacion_dato"] = strArmonizacion
				fmt.Println(subgrupo_detalle["armonizacion_dato"], "armonización dato")
				//dataProyectIn = true
				//}
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
				fmt.Println(subgrupo_detalle["dato_plan"], "dato_plan 2")

				// //if !dataProyectIn {
				// armonizacion_dato := make(map[string]interface{})
				// if subgrupo_detalle["armonizacion_dato"] != nil {
				// 	armonizacion_dato_str := subgrupo_detalle["armonizacion_dato"].(string)
				// 	json.Unmarshal([]byte(armonizacion_dato_str), &armonizacion_dato)
				// 	aux := make(map[string]interface{})
				// 	aux["idSubDetalleProI"] = idSubDetalleProI
				// 	aux["indexMetaSubProI"] = indexMetaSubProI
				// 	aux["indexMetaPlan"] = i
				// 	armonizacion_dato[i] = aux
				// 	c, _ := json.Marshal(armonizacion_dato)
				// 	strArmonizacion := string(c)
				// 	subgrupo_detalle["armonizacion_dato"] = strArmonizacion
				// 	fmt.Println(subgrupo_detalle["armonizacion_dato"], "armonización dato 2")
				// }

				//dataProyectIn = true

				//}
			}
			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarMeta", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
			}
			helpers.LimpiezaRespuestaRefactor(res, &respuestaGuardado)
			//fmt.Println(res, "actividad Guardada")
			fmt.Println(respuestaGuardado, "actividad Guardada")
		}
	}
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuestaGuardado}
	c.ServeJSON()

}

// GetPlan ...
// @Title GetPlan
// @Description get Plan by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /get_infoPlan/:id [get]
func (c *InversionController) GetPlan() {
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
	//var body map[string]interface{}
	var subgrupo []map[string]interface{}
	var res map[string]interface{}
	//getProyect := make(map[string]interface{})
	var id_subgrupoDetalle string
	var respuesta map[string]interface{}
	var respuestaLimpia []map[string]interface{}
	var subgrupo_detalle map[string]interface{}
	armo_dato := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=descripcion:Armonizar,activo:true,padre:"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupo)
		fmt.Println(res, "subgrupo")
		// for i := range subgrupo {
		// 	var subgrupoDetalle map[string]interface{}
		// 	var detalleSubgrupos []map[string]interface{}
		// 	if strings.Contains(strings.ToLower(subgrupo[i]["nombre"].(string)), "plan") {
		// 		//getProyect["subgrupo_id_metas"] = subgrupo[i]["_id"]
		// 		fmt.Println()
		// 		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+subgrupo[i]["_id"].(string), &subgrupoDetalle); err == nil {
		// 			helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
		// 			fmt.Println(subgrupoDetalle, "detalleSub")
		// 			var dato_metas []map[string]interface{}
		// 			datoMeta_str := detalleSubgrupos[0]["armonizacion_dato"].(string)
		// 			fmt.Println(datoMeta_str, "datoMetas_str")
		// 			json.Unmarshal([]byte(datoMeta_str), &dato_metas)
		// 			fmt.Println(dato_metas, "datoMetas")
		// 			getProyect["armonizacion"] = dato_metas
		// 			//getProyect["id_detalle_metas"] = detalleSubgrupos[0]["_id"]

		// 		} else {
		// 			panic(map[string]interface{}{"funcion": "GetPlan", "err": "Error get subgrupo-detalle", "status": "400", "log": err})
		// 		}
		// 	}
		// }

		id_subgrupoDetalle = subgrupo[0]["_id"].(string)
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=activo:true,subgrupo_id:"+id_subgrupoDetalle, &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
			subgrupo_detalle = respuestaLimpia[0]
			fmt.Println(subgrupo_detalle, "subgrupo Detalle")
			armonizacion_dato_str := subgrupo_detalle["armonizacion_dato"].(string)
			json.Unmarshal([]byte(armonizacion_dato_str), &armo_dato)

		} else {
			panic(map[string]interface{}{"funcion": "GetPlan", "err": "Error get subgrupo-detalle", "status": "400", "log": err})
		}

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": armo_dato}
	} else {
		panic(map[string]interface{}{"funcion": "GetPlan", "err": "Error get subgrupo", "status": "400", "log": err})

	}

	c.ServeJSON()
}

// ArmonizarInversion ...
// @Title ArmonizarInversion
// @Description get Plan by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /armonizar/:id [put]
func (c *InversionController) ArmonizarInversion() {
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
	var body map[string]interface{}
	var subgrupo []map[string]interface{}
	var id_subgrupoDetalle string
	var respuesta map[string]interface{}
	var respuestaLimpia []map[string]interface{}
	var armonizacionUpdate []map[string]interface{}
	var subgrupo_detalle map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	armonizacion_data, _ := json.Marshal(body)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=descripcion:Armonizar,activo:true,padre:"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupo)
		//armo_dato := make(map[string]interface{})
		subgrupoPost := make(map[string]interface{})
		subDetallePost := make(map[string]interface{})

		if len(subgrupo) != 0 {
			id_subgrupoDetalle = subgrupo[0]["_id"].(string)
			fmt.Println(id_subgrupoDetalle, "id_subgrupoDetalle")
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=activo:true,subgrupo_id:"+id_subgrupoDetalle, &respuesta); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
				if len(respuestaLimpia) != 0 {
					fmt.Println(respuestaLimpia, "subgrupoDetalleb PUT")
					subgrupo_detalle = respuestaLimpia[0]
					// if subgrupo_detalle["armonizacion_dato"] != nil {
					// 	actividad := make(map[string]interface{})
					//armonizacion_dato_str := subgrupo_detalle["armonizacion_dato"].(string)
					//json.Unmarshal([]byte(armonizacion_dato_str), &armo_dato)

					// 	b, _ := json.Marshal(dato_plan)
					// 	str := string(b)
					// 	subgrupo_detalle["dato_plan"] = str
					// }
					subDetallePost["subgrupo_id"] = id_subgrupoDetalle
					subDetallePost["fecha_creacion"] = subgrupo_detalle["fecha_creacion"]
					subDetallePost["nombre"] = "Detalle Información Armonización"
					subDetallePost["descripcion"] = "Armonizar"
					subDetallePost["dato"] = " "
					subDetallePost["activo"] = true
					subDetallePost["armonizacion_dato"] = string(armonizacion_data)
					fmt.Println(subDetallePost, "dataJSON")
					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subDetallePost); err == nil {
						helpers.LimpiezaRespuestaRefactor(res, &armonizacionUpdate)
						fmt.Println(armonizacionUpdate, "update911")
					} else {
						panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
					}
				} else {
					subDetallePost["subgrupo_id"] = id_subgrupoDetalle
					subDetallePost["nombre"] = "Detalle Información Armonización"
					subDetallePost["descripcion"] = "Armonizar"
					subDetallePost["dato"] = " "
					subDetallePost["activo"] = true
					subDetallePost["armonizacion_dato"] = string(armonizacion_data)
					fmt.Println(subDetallePost["armonizacion_dato"], "dataJSON")

					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/", "POST", &res, subDetallePost); err == nil {
						helpers.LimpiezaRespuestaRefactor(res, &armonizacionUpdate)
						fmt.Println(res, "update926")
					} else {
						panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error registrando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
					}
				}

			} else {
				panic(map[string]interface{}{"funcion": "ArmonizarInversion", "err": "Error get subgrupo-detalle", "status": "400", "log": err})
			}
		} else {
			subgrupoPost["nombre"] = "Armonización Plan Inversión"
			subgrupoPost["descripcion"] = "Armonizar"
			subgrupoPost["padre"] = id
			subgrupoPost["activo"] = true
			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/", "POST", &respuesta, subgrupoPost); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
				fmt.Println(respuestaLimpia, "respuesta subgrupo POST")
				subgrupo_detalle = respuestaLimpia[0]
				subDetallePost["subgrupo_id"] = subgrupo_detalle["_id"]
				subDetallePost["nombre"] = "Detalle Información Armonización"
				subDetallePost["descripcion"] = "Armonizar"
				subDetallePost["dato"] = " "
				subDetallePost["activo"] = true
				subDetallePost["armonizacion_dato"] = string(armonizacion_data)
				fmt.Println(subDetallePost["armonizacion_dato"], "dataJSON")

				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/", "POST", &res, subDetallePost); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &armonizacionUpdate)
					fmt.Println(armonizacionUpdate, "update954")
				} else {
					panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error registrando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
				}

			} else {
				panic(map[string]interface{}{"funcion": "ArmonizarInversion", "err": "Error registrando subgrupo", "status": "400", "log": err})
			}
		}
		//formulacionhelper.Limpia()
		//tree := formulacionhelper.BuildTreeFa(hijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": res}
	} else {
		panic(err)

	}

	c.ServeJSON()
}

// ActualizarMetaPlan ...
// @Title ActualizarMetaPlan
// @Description put Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /actualizar_meta/:id/:index [put]
func (c *InversionController) ActualizarMetaPlan() {
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
	index := c.Ctx.Input.Param(":index")

	var res map[string]interface{}
	var entrada map[string]interface{}
	var body map[string]interface{}

	_ = id
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	entrada = body["entrada"].(map[string]interface{})
	idSubDetalleProI := body["idSubDetalle"]
	indexMetaSubProI := body["indexMetaSubPro"]
	for key, element := range entrada {
		var respuesta map[string]interface{}
		var respuestaLimpia []map[string]interface{}
		var subgrupo_detalle map[string]interface{}
		dato_plan := make(map[string]interface{})
		var armonizacion_dato map[string]interface{}
		var id_subgrupoDetalle string
		keyStr := strings.Split(key, "_")

		if len(keyStr) > 1 && keyStr[1] == "o" {
			id_subgrupoDetalle = keyStr[0]
			if element != "" {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id_subgrupoDetalle, &respuesta); err != nil {
					panic(map[string]interface{}{"funcion": "ActualizarMetaPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
				}
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)

				subgrupo_detalle = respuestaLimpia[0]
				if subgrupo_detalle["dato_plan"] != nil {
					meta := make(map[string]interface{})
					dato_plan_str := subgrupo_detalle["dato_plan"].(string)
					json.Unmarshal([]byte(dato_plan_str), &dato_plan)
					fmt.Println(dato_plan, "dato_plan")
					for index_actividad := range dato_plan {
						if index_actividad == index {
							aux_actividad := dato_plan[index_actividad].(map[string]interface{})
							meta["index"] = index_actividad
							meta["dato"] = aux_actividad["dato"]
							meta["activo"] = aux_actividad["activo"]
							meta["observacion"] = element
							dato_plan[index_actividad] = meta

							// aux := make(map[string]interface{})
							// aux["idSubDetalleProI"] = idSubDetalleProI
							// aux["indexMetaSubProI"] = indexMetaSubProI
							// armonizacion_dato[index] = aux
						}
					}
					b, _ := json.Marshal(dato_plan)
					str := string(b)
					subgrupo_detalle["dato_plan"] = str
				}

				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
					panic(map[string]interface{}{"funcion": "ActualizarMetaPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
				}
				fmt.Println(res, "res 1058")

			}
			continue
		}
		id_subgrupoDetalle = key
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id_subgrupoDetalle, &respuesta); err != nil {
			panic(map[string]interface{}{"funcion": "ActualizarMetaPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
		}
		helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)

		subgrupo_detalle = respuestaLimpia[0]
		if subgrupo_detalle["armonizacion_dato"] != nil {
			dato_armonizacion_str := subgrupo_detalle["armonizacion_dato"].(string)
			json.Unmarshal([]byte(dato_armonizacion_str), &armonizacion_dato)
			if armonizacion_dato[index] != nil {
				aux := make(map[string]interface{})
				aux["idSubDetalleProI"] = idSubDetalleProI
				aux["indexMetaSubProI"] = indexMetaSubProI
				armonizacion_dato[index] = aux
				fmt.Println(armonizacion_dato, "armonizacion_dato")
			}
			c, _ := json.Marshal(armonizacion_dato)
			strArmonizacion := string(c)
			subgrupo_detalle["armonizacion_dato"] = strArmonizacion

		}

		nuevoDato := true
		meta := make(map[string]interface{})

		if subgrupo_detalle["dato_plan"] != nil {
			dato_plan_str := subgrupo_detalle["dato_plan"].(string)
			json.Unmarshal([]byte(dato_plan_str), &dato_plan)

			for index_actividad := range dato_plan {
				if index_actividad == index {
					nuevoDato = false
					aux_actividad := dato_plan[index_actividad].(map[string]interface{})
					meta["index"] = index_actividad
					meta["dato"] = element
					meta["activo"] = aux_actividad["activo"]
					if aux_actividad["observacion"] != nil {
						meta["observacion"] = aux_actividad["observacion"]
					}
					dato_plan[index_actividad] = meta
				}
			}
		}

		if nuevoDato {
			meta["index"] = index
			meta["dato"] = element
			meta["activo"] = true
			dato_plan[index] = meta
		}

		b, _ := json.Marshal(dato_plan)
		str := string(b)
		subgrupo_detalle["dato_plan"] = str

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
			panic(map[string]interface{}{"funcion": "ActualizarMetaPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
		}
		fmt.Println(res, "res 1121")
	}
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": entrada}
	c.ServeJSON()

}

// AllMetasPlan ...
// @Title AllMetasPlan
// @Description put Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /all_metas/:id [get]
func (c *InversionController) AllMetasPlan() {
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
	var hijos []map[string]interface{}
	var tabla map[string]interface{}
	var metas []map[string]interface{}
	var auxHijos []interface{}
	var data_source []map[string]interface{}
	inversionhelper.Limpia()
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		for i := 0; i < len(hijos); i++ {
			auxHijos = append(auxHijos, hijos[i]["_id"])
		}
		fmt.Println(auxHijos, "auxhijos")
		tabla = inversionhelper.GetTabla(auxHijos)
		metas = tabla["data_source"].([]map[string]interface{})
		fmt.Println(tabla, "tabla")
		for indexMeta := range metas {
			if metas[indexMeta]["activo"] == true {
				data_source = append(data_source, metas[indexMeta])
			}
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data_source}
	} else {
		panic(map[string]interface{}{"funcion": "AllMetasPlan", "err": "Error al consultar metas del plan \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
	}
	c.ServeJSON()
	// id := c.Ctx.Input.Param(":id")
	// index := c.Ctx.Input.Param(":index")
	// var res map[string]interface{}
	// var hijos []map[string]interface{}

	// if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
	// 	helpers.LimpiezaRespuestaRefactor(res, &hijos)
	// 	inversionhelper.Limpia()
	// 	tree := inversionhelper.BuildTreeFa(hijos, index)
	// 	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": tree}
	// } else {
	// 	panic(err)

	// }

	c.ServeJSON()

}

// InactivarMeta ...
// @Title InactivarMeta
// @Description put Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /inactivar_meta/:id/:index [put]
func (c *InversionController) InactivarMeta() {
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
	index := c.Ctx.Input.Param(":index")

	// var res map[string]interface{}
	// var hijos []map[string]interface{}
	// inversionhelper.Limpia()

	// if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
	// 	helpers.LimpiezaRespuestaRefactor(res, &hijos)
	// 	inversionhelper.GetSons(hijos, index)
	// 	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Meta Inactivada"}
	// } else {
	// 	panic(err)
	// }

	// c.ServeJSON()
	var res map[string]interface{}
	var hijos []map[string]interface{}
	var tabla map[string]interface{}
	var auxHijos []interface{}
	inversionhelper.Limpia()
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		for i := 0; i < len(hijos); i++ {
			auxHijos = append(auxHijos, hijos[i]["_id"])
		}
		fmt.Println(auxHijos, "auxhijos")
		inversionhelper.GetSons(auxHijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": tabla}
	} else {
		panic(map[string]interface{}{"funcion": "AllMetasPlan", "err": "Error al consultar metas del plan \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
	}
	c.ServeJSON()
}

// ProgMagnitudesPlan ...
// @Title ProgMagnitudesPlan
// @Description get Plan by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @Failure 403 :id is empty
// @router /magnitudes/:id/:index [put]
func (c *InversionController) ProgMagnitudesPlan() {
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
	index := c.Ctx.Input.Param(":index")
	var res map[string]interface{}
	var body map[string]interface{}
	var subgrupo []map[string]interface{}
	var id_subgrupoDetalle string
	//var index string
	var respuesta map[string]interface{}
	var respuestaLimpia []map[string]interface{}
	var magnitudesUpdate []map[string]interface{}
	var subgrupo_detalle map[string]interface{}
	dato := make(map[string]interface{})

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	//index = body["indiceMetaProyecto"].(string)
	//magnitudes_data, _ := json.Marshal(body)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=descripcion:Magnitudes,activo:true,padre:"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupo)
		//armo_dato := make(map[string]interface{})
		subgrupoPost := make(map[string]interface{})
		subDetallePost := make(map[string]interface{})

		if len(subgrupo) != 0 {
			id_subgrupoDetalle = subgrupo[0]["_id"].(string)
			fmt.Println(id_subgrupoDetalle, "id_subgrupoDetalle")
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=activo:true,subgrupo_id:"+id_subgrupoDetalle, &respuesta); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
				fmt.Println(respuestaLimpia, "subgrupoDetalleb PUT")
				if len(respuestaLimpia) != 0 {
					fmt.Println("ingresa a PUT")
					subgrupo_detalle = respuestaLimpia[0]
					fmt.Println(subgrupo_detalle["dato"], "subgrupo detalle put")
					magnitud := make(map[string]interface{})

					if subgrupo_detalle["dato"] == nil {
						magnitud["index"] = index
						magnitud["dato"] = body
						magnitud["activo"] = true
						i := strconv.Itoa(magnitud["index"].(int))
						dato[i] = magnitud
						b, _ := json.Marshal(dato)
						str := string(b)
						subgrupo_detalle["dato"] = str
					} else {
						dato_str := subgrupo_detalle["dato"].(string)
						json.Unmarshal([]byte(dato_str), &dato)
						magnitud["index"] = index
						magnitud["dato"] = body
						magnitud["activo"] = true
						//i := strconv.Itoa(magnitud["index"].(int))
						dato[index] = magnitud
						b, _ := json.Marshal(dato)
						str := string(b)
						subgrupo_detalle["dato"] = str
					}
					subDetallePost["dato"] = subgrupo_detalle["dato"]
					subDetallePost["subgrupo_id"] = id_subgrupoDetalle
					subDetallePost["fecha_creacion"] = subgrupo_detalle["fecha_creacion"]
					subDetallePost["nombre"] = "Detalle Información Programación de Magnitudes y Presupuesto"
					subDetallePost["descripcion"] = "Magnitudes"
					subDetallePost["activo"] = true

					fmt.Println(subDetallePost, "dataJSON")

					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subDetallePost); err == nil {
						helpers.LimpiezaRespuestaRefactor(res, &magnitudesUpdate)
						fmt.Println(magnitudesUpdate, "update1290")
					} else {
						panic(map[string]interface{}{"funcion": "ProgMagnitudesPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
					}
				} else {
					magnitud := make(map[string]interface{})
					magnitud["index"] = index
					magnitud["dato"] = body
					magnitud["activo"] = true
					i := strconv.Itoa(magnitud["index"].(int))
					dato[i] = magnitud
					b, _ := json.Marshal(dato)
					str := string(b)
					subgrupo_detalle["dato"] = str

					subDetallePost["subgrupo_id"] = id_subgrupoDetalle
					subDetallePost["nombre"] = "Detalle Información Programación de Magnitudes y Presupuesto"
					subDetallePost["descripcion"] = "Magnitudes"
					subDetallePost["dato"] = subgrupo_detalle["dato"]
					subDetallePost["activo"] = true
					fmt.Println(subDetallePost, "dataJSON")

					if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/", "POST", &res, subDetallePost); err == nil {
						helpers.LimpiezaRespuestaRefactor(res, &magnitudesUpdate)
						fmt.Println(res, "update1304")
					} else {
						panic(map[string]interface{}{"funcion": "ProgMagnitudesPlan", "err": "Error registrando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
					}
				}

			} else {
				panic(map[string]interface{}{"funcion": "ProgMagnitudesPlan", "err": "Error get subgrupo-detalle", "status": "400", "log": err})
			}
		} else {
			subgrupoPost["nombre"] = "Programación Magnitudes y Prespuesto Plan de Inversión"
			subgrupoPost["descripcion"] = "Magnitudes"
			subgrupoPost["padre"] = id
			subgrupoPost["activo"] = true
			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/", "POST", &respuesta, subgrupoPost); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
				fmt.Println(respuestaLimpia, "respuesta subgrupo POST")
				subgrupo_detalle = respuestaLimpia[0]
				magnitud := make(map[string]interface{})
				magnitud["index"] = index
				magnitud["dato"] = body
				magnitud["activo"] = true
				i := strconv.Itoa(magnitud["index"].(int))
				dato[i] = magnitud
				b, _ := json.Marshal(dato)
				str := string(b)
				subgrupo_detalle["dato"] = str
				subDetallePost["subgrupo_id"] = subgrupo_detalle["_id"]
				subDetallePost["nombre"] = "Detalle Información Programación de Magnitudes y Presupuesto"
				subDetallePost["descripcion"] = "Magnitudes"
				subDetallePost["dato"] = subgrupo_detalle["dato"]
				subDetallePost["activo"] = true
				fmt.Println(subDetallePost, "dataJSON")

				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/", "POST", &res, subDetallePost); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &magnitudesUpdate)
					fmt.Println(magnitudesUpdate, "update1331")
				} else {
					panic(map[string]interface{}{"funcion": "ProgMagnitudesPlan", "err": "Error registrando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
				}

			} else {
				panic(map[string]interface{}{"funcion": "ProgMagnitudesPlan", "err": "Error registrando subgrupo", "status": "400", "log": err})
			}
		}
		//formulacionhelper.Limpia()
		//tree := formulacionhelper.BuildTreeFa(hijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": res}
	} else {
		panic(err)

	}

	c.ServeJSON()
}

// MagnitudesProgramadas ...
// @Title MagnitudesProgramadas
// @Description get MagnitudesProgramadas
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /magnitudes/:id/:indexMeta [get]
func (c *InversionController) MagnitudesProgramadas() {
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
	index := c.Ctx.Input.Param(":indexMeta")
	var res map[string]interface{}
	//var body map[string]interface{}
	var subgrupo map[string]interface{}
	//var id_subgrupoDetalle string
	var respuesta map[string]interface{}
	var respuestaLimpia []map[string]interface{}
	//var armonizacionUpdate []map[string]interface{}
	var subgrupo_detalle map[string]interface{}
	//json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	//armonizacion_data, _ := json.Marshal(body)
	dato := make(map[string]interface{})
	var magnitud map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=descripcion:Magnitudes,activo:true,padre:"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &respuestaLimpia)
		subgrupo = respuestaLimpia[0]
		fmt.Println(res, "subgrupo")
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=activo:true,subgrupo_id:"+subgrupo["_id"].(string), &respuesta); err != nil {
			panic(map[string]interface{}{"funcion": "MagnitudesProgramadas", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
		}
		helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)
		subgrupo_detalle = respuestaLimpia[0]

		if subgrupo_detalle["dato"] != nil {

			dato_str := subgrupo_detalle["dato"].(string)
			json.Unmarshal([]byte(dato_str), &dato)
			for index_actividad := range dato {
				if index_actividad == index {
					aux_actividad := dato[index_actividad].(map[string]interface{})
					magnitud = aux_actividad
				}
			}

		}
		fmt.Println(magnitud)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": magnitud}
	} else {
		panic(map[string]interface{}{"funcion": "MagnitudesProgramadas", "err": "Error consultando subgrupo", "status": "400", "log": err})
	}
	c.ServeJSON()
}
