package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"

	"github.com/udistrital/planeacion_mid/helpers"
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
	c.Mapping("GetAllProyectos", c.GetAllProyectos)
	c.Mapping("ActualizarSubgrupoDetalle", c.ActualizarSubgrupoDetalle)
	c.Mapping("ActualizarProyectoGeneral", c.ActualizarProyectoGeneral)
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

// GetAllProyectos ...
// @Title GetAllProyectos
// @Description get GetAllProyectos
// @Param	aplicativo_id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :aplicativo_id is empty
// @router /getproyectos/:aplicativo_id [get]
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
