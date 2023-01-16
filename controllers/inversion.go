package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego"

	"github.com/udistrital/planeacion_mid/helpers"
	inversionhelper "github.com/udistrital/planeacion_mid/helpers/inversionHelper"
	"github.com/udistrital/utils_oas/request"
)

type InversionController struct {
	beego.Controller
}

// URLMapping ...
func (c *InversionController) URLMapping() {
	c.Mapping("AddProyecto", c.AddProyecto)
	c.Mapping("GetProyectoId", c.GetProyectoId)
	c.Mapping("ActualizarSubgrupoDetalle", c.ActualizarSubgrupoDetalle)
	c.Mapping("ActualizarProyectoGeneral", c.ActualizarProyectoGeneral)
	// c.Mapping("GetActividadesGenerales", c.GetActividadesGenerales)
	// c.Mapping("GetDataActividad", c.GetDataActividad)
	// c.Mapping("GuardarSeguimiento", c.GuardarSeguimiento)
	// c.Mapping("GetSeguimiento", c.GetSeguimiento)
	// c.Mapping("GetIndicadores", c.GetIndicadores)
	// c.Mapping("GetAvanceIndicador", c.GetAvanceIndicador)
}

// AddProyecto ...
// @Title AddProyecto
// @Description post AddProyecto
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200
// @router /addProyecto [post]
func (c *InversionController) AddProyecto() {
	var registroProyecto map[string]interface{}
	plan := make(map[string]interface{})
	//var dataProyect map[string]interface{}
	var resPost map[string]interface{}
	var idProyect string
	var resSoportes map[string]interface{}
	clienteHttp := &http.Client{}
	var idPlan map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &registroProyecto); err == nil {
		fmt.Println(registroProyecto)
		plan["activo"] = true
		plan["nombre"] = registroProyecto["nombre_proyecto"]
		plan["descripcion"] = registroProyecto["codigo_proyecto"]
		plan["tipo_plan_id"] = " "
		plan["aplicativo_id"] = "idPlaneacion"
		fmt.Println(plan)
		aux, err := json.Marshal(plan)
		if err != nil {
			panic(err)
		}

		peticion, err := http.NewRequest("POST", "http://"+beego.AppConfig.String("PlanesService")+"/plan", bytes.NewBuffer(aux))
		if err != nil {
			panic(err)
		}
		peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
		respuesta, err := clienteHttp.Do(peticion)
		if err != nil {
			panic(err)
		}

		defer respuesta.Body.Close()

		cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
		if err != nil {
			panic(err)
		}

		json.Unmarshal(cuerpoRespuesta, &resPost)
		fmt.Println(resPost)
		idPlan = resPost["Data"].(map[string]interface{})
		idProyect = idPlan["_id"].(string)
		fmt.Println(idProyect, "prueba")
		resSoportes, err = inversionhelper.RegistrarSoportes(idProyect, registroProyecto)
		fmt.Println(err)
		soporte := resSoportes["Data"].(map[string]interface{})
		idSoporte := soporte["_id"].(string)
		resSoporteDetalle, e := inversionhelper.RegistrarSoporteDetalle(idSoporte, registroProyecto)
		fmt.Println(e)
		fmt.Println(resSoporteDetalle)
		resFuentesApropiacion, e := inversionhelper.RegistrarFuentesApropiacion(idProyect, registroProyecto)
		fmt.Println(e)
		fuentes := resFuentesApropiacion["Data"].(map[string]interface{})
		idFuentes := fuentes["_id"].(string)
		resFuentesDetalle, e := inversionhelper.RegistrarFuentesDetalle(idFuentes, registroProyecto)
		fmt.Println(resFuentesDetalle, e)
		resMetas, e := inversionhelper.RegistrarMetas(idProyect, registroProyecto)
		fmt.Println(e)
		metas := resMetas["Data"].(map[string]interface{})
		idMetas := metas["_id"].(string)
		resMetasDetalle, e := inversionhelper.RegistrarMetasDetalle(idMetas, registroProyecto)
		fmt.Println(resMetasDetalle)
		fmt.Println(e)
		if e != nil {
			c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "error", "Data": e}
			c.Abort("400")
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": resPost}
		c.ServeJSON()
	} else {
		panic(map[string]interface{}{"funcion": "AddProyecto", "err": "Error Registrando Proyecto", "status": "400", "log": err})
	}

}

// GetProyectoId ...
// @Title GetProyectoId
// @Description get GetProyectoId
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 403 :id is empty
// @router /getproyectoid/:id [get]
func (c *InversionController) GetProyectoId() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		localError := err.(map[string]interface{})
	// 		c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "InversionController" + "/" + (localError["funcion"]).(string))
	// 		c.Data["data"] = (localError["err"])
	// 		if status, ok := localError["status"]; ok {
	// 			c.Abort(status.(string))
	// 		} else {
	// 			c.Abort("404")
	// 		}
	// 	}
	// }()
	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	getProyect := make(map[string]interface{})
	var infoProyect map[string]interface{}
	var subgruposData map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &infoProyect)
		fmt.Println(infoProyect)
		getProyect["nombre_proyecto"] = infoProyect["nombre"]
		getProyect["codigo_proyecto"] = infoProyect["descripcion"]
		getProyect["fecha_creacion"] = infoProyect["fecha_creacion"]
		padreId := infoProyect["_id"].(string)

		var infoSubgrupos []map[string]interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+padreId, &subgruposData); err == nil {
			fmt.Println(subgruposData)
			helpers.LimpiezaRespuestaRefactor(subgruposData, &infoSubgrupos)
			//getProyect["fecha_creacion"]
			fmt.Println(infoSubgrupos)
			for i := range infoSubgrupos {
				var subgrupoDetalle map[string]interface{}
				var detalleSubgrupos []map[string]interface{}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "soporte") {
					getProyect["subgrupo_id_soportes"] = infoSubgrupos[i]["_id"]
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
						getProyect["soportes"] = detalleSubgrupos[0]["dato"]
						getProyect["id_detalle_soportes"] = detalleSubgrupos[0]["_id"]
					}
				}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "metas") {
					getProyect["subgrupo_id_metas"] = infoSubgrupos[i]["_id"]
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
						getProyect["metas"] = detalleSubgrupos[0]["dato"]
						getProyect["id_detalle_metas"] = detalleSubgrupos[0]["_id"]

					}
				}
				if strings.Contains(strings.ToLower(infoSubgrupos[i]["nombre"].(string)), "fuentes") {
					getProyect["subgrupo_id_fuentes"] = infoSubgrupos[i]["_id"]
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+infoSubgrupos[i]["_id"].(string), &subgrupoDetalle); err == nil {
						helpers.LimpiezaRespuestaRefactor(subgrupoDetalle, &detalleSubgrupos)
						getProyect["fuentes"] = detalleSubgrupos[0]["dato"]
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

// ActualizarSubgrupoDetalle ...
// @Title ActualizarSubgrupoDetalle
// @Description put Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object}
// @Failure 403 :id is empty
// @router /actualiza-sub-detalle/:id [put]
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
// @Success 200 {object}
// @Failure 403 :id is empty
// @router /actualizar-proyecto/:id [put]
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