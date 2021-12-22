package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"

	formulacionhelper "github.com/udistrital/planeacion_mid/helpers/formulacionHelper"
)

// FormulacionController operations for Formulacion
type FormulacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *FormulacionController) URLMapping() {
	c.Mapping("ClonarFormato", c.ClonarFormato)
	c.Mapping("GuardarActividad", c.GuardarActividad)
	c.Mapping("GetPlan", c.GetPlan)
	c.Mapping("ActualizarActividad", c.ActualizarActividad)
	c.Mapping("DeleteActividad", c.DeleteActividad)
	c.Mapping("GetAllActividades", c.GetAllActividades)
	c.Mapping("GuardarIdentificacion", c.GuardarIdentificacion)
	c.Mapping("GetAllIdentificacion", c.GetAllIdentificacion)
	c.Mapping("DeleteIdentificacion", c.DeleteIdentificacion)
	c.Mapping("VersionarPlan", c.VersionarPlan)
	c.Mapping("GetPlanVersiones", c.GetPlanVersiones)
	c.Mapping("PonderacionActividades", c.PonderacionActividades)
	c.Mapping("GetRubros", c.GetRubros)
	c.Mapping("GetUnidades", c.GetUnidades)

}

// ClonarFormato ...
// @Title ClonarFormato
// @Description post Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /clonar_formato/:id [post]
func (c *FormulacionController) ClonarFormato() {

	id := c.Ctx.Input.Param(":id")

	var respuesta map[string]interface{}
	var respuestaHijos map[string]interface{}
	var hijos []map[string]interface{}
	var planFormato map[string]interface{}
	var parametros map[string]interface{}

	plan := make(map[string]interface{})
	clienteHttp := &http.Client{}
	url := beego.AppConfig.String("PlanesService") + "/plan/"

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/plan/"+id, &respuesta); err == nil {

		helpers.LimpiezaRespuestaRefactor(respuesta, &planFormato)
		json.Unmarshal(c.Ctx.Input.RequestBody, &parametros)

		plan["nombre"] = "" + planFormato["nombre"].(string)
		plan["descripcion"] = planFormato["descripcion"].(string)
		plan["tipo_plan_id"] = planFormato["tipo_plan_id"].(string)
		plan["aplicativo_id"] = planFormato["aplicativo_id"].(string)
		plan["activo"] = planFormato["activo"]
		plan["formato"] = false
		plan["vigencia"] = parametros["vigencia"].(string)
		plan["dependencia_id"] = parametros["dependencia_id"].(string)
		plan["estado_plan_id"] = "614d3ad301c7a200482fabfd"

		var resPost map[string]interface{}
		var resLimpia map[string]interface{}

		aux, err := json.Marshal(plan)
		if err != nil {
			log.Fatalf("Error codificado: %v", err)
		}

		peticion, err := http.NewRequest("POST", url, bytes.NewBuffer(aux))
		if err != nil {
			log.Fatalf("Error creando peticion: %v", err)
		}
		peticion.Header.Set("Content-Type", "application/json; charset=UTF-8")
		respuesta, err := clienteHttp.Do(peticion)
		if err != nil {
			log.Fatalf("Error haciendo peticion: %v", err)
		}

		defer respuesta.Body.Close()

		cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
		if err != nil {
			log.Fatalf("Error leyendo peticion: %v", err)
		}

		json.Unmarshal(cuerpoRespuesta, &resPost)
		resLimpia = resPost["Data"].(map[string]interface{})
		padre := resLimpia["_id"].(string)
		c.Data["json"] = resPost

		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &respuestaHijos); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijos, &hijos)
			formulacionhelper.ClonarHijos(hijos, padre)
		}

	}
	c.ServeJSON()

}

// GuardarActividad ...
// @Title GuardarActividad
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /guardar_actividad/:id [put]
func (c *FormulacionController) GuardarActividad() {
	id := c.Ctx.Input.Param(":id")
	var body map[string]interface{}
	var res map[string]interface{}
	var entrada map[string]interface{}
	var resPlan map[string]interface{}
	var plan map[string]interface{}
	var armonizacionExecuted bool = false
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	entrada = body["entrada"].(map[string]interface{})
	armonizacion := body["armo"]

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/plan/"+id, &resPlan); err != nil {
		panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get Plan \"id\"", "status": "400", "log": err})
	}
	helpers.LimpiezaRespuestaRefactor(resPlan, &plan)
	if plan["estado_plan_id"] != "614d3ad301c7a200482fabfd" {
		var res map[string]interface{}
		plan["estado_plan_id"] = "614d3ad301c7a200482fabfd"
		if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/plan/"+id, "PUT", &res, plan); err != nil {
			panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizacion estado \"id\"", "status": "400", "log": err})
		}
	}

	for key, element := range entrada {

		var respuesta map[string]interface{}
		var respuestaLimpia []map[string]interface{}
		var subgrupo_detalle map[string]interface{}
		dato_plan := make(map[string]interface{})

		if element != "" {
			if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+key, &respuesta); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
			}
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
					armonizacion_dato[i] = armonizacion
					c, _ := json.Marshal(armonizacion_dato)
					strArmonizacion := string(c)
					subgrupo_detalle["armonizacion_dato"] = strArmonizacion
					armonizacionExecuted = true
				}
			} else {
				dato_plan_str := subgrupo_detalle["dato_plan"].(string)
				json.Unmarshal([]byte(dato_plan_str), &dato_plan)
				maxIndex := 0

				for key2 := range dato_plan {
					index, err := strconv.Atoi(key2)
					if err != nil {
						log.Panic(err)
					}
					if index > maxIndex {
						maxIndex = index
					}

				}

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

					armonizacion_dato[i] = armonizacion
					c, _ := json.Marshal(armonizacion_dato)
					strArmonizacion := string(c)
					subgrupo_detalle["armonizacion_dato"] = strArmonizacion
					armonizacionExecuted = true

				}
			}
			if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
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
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_plan/:id/:index [get]
func (c *FormulacionController) GetPlan() {
	id := c.Ctx.Input.Param(":id")
	index := c.Ctx.Input.Param(":index")
	var res map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		formulacionhelper.Limpia()
		tree := formulacionhelper.BuildTreeFa(hijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": tree}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// ActualizarActividad ...
// @Title ActualizarActividad
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /actualizar_actividad/:id/:index [put]
func (c *FormulacionController) ActualizarActividad() {
	id := c.Ctx.Input.Param(":id")
	index := c.Ctx.Input.Param(":index")

	var res map[string]interface{}
	var entrada map[string]interface{}
	var body map[string]interface{}

	fmt.Println(id)
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	entrada = body["entrada"].(map[string]interface{})
	armonizacion := body["armo"]

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
				if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id_subgrupoDetalle, &respuesta); err != nil {
					panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
				}
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)

				subgrupo_detalle = respuestaLimpia[0]
				if subgrupo_detalle["dato_plan"] != nil {
					actividad := make(map[string]interface{})
					dato_plan_str := subgrupo_detalle["dato_plan"].(string)
					json.Unmarshal([]byte(dato_plan_str), &dato_plan)
					for index_actividad := range dato_plan {
						if index_actividad == index {
							aux_actividad := dato_plan[index_actividad].(map[string]interface{})
							actividad["index"] = index_actividad
							actividad["dato"] = aux_actividad["dato"]
							actividad["activo"] = aux_actividad["activo"]
							actividad["observacion"] = element

							dato_plan[index_actividad] = actividad
						}
					}
					b, _ := json.Marshal(dato_plan)
					str := string(b)
					subgrupo_detalle["dato_plan"] = str
				}

				if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
					panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
				}

			}
		} else {
			id_subgrupoDetalle = key
			if element != "" {
				if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id_subgrupoDetalle, &respuesta); err != nil {
					panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
				}
				helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)

				subgrupo_detalle = respuestaLimpia[0]
				if subgrupo_detalle["armonizacion_dato"] != nil {
					dato_armonizacion_str := subgrupo_detalle["armonizacion_dato"].(string)
					json.Unmarshal([]byte(dato_armonizacion_str), &armonizacion_dato)
					if armonizacion_dato[index] != nil {
						armonizacion_dato[index] = armonizacion
					}
					c, _ := json.Marshal(armonizacion_dato)
					strArmonizacion := string(c)
					subgrupo_detalle["armonizacion_dato"] = strArmonizacion

				}
				if subgrupo_detalle["dato_plan"] != nil {
					actividad := make(map[string]interface{})
					dato_plan_str := subgrupo_detalle["dato_plan"].(string)
					json.Unmarshal([]byte(dato_plan_str), &dato_plan)
					for index_actividad := range dato_plan {
						if index_actividad == index {
							aux_actividad := dato_plan[index_actividad].(map[string]interface{})
							actividad["index"] = index_actividad
							actividad["dato"] = element
							actividad["activo"] = aux_actividad["activo"]
							if aux_actividad["observacion"] != nil {
								actividad["observacion"] = aux_actividad["observacion"]
							}
							dato_plan[index_actividad] = actividad
						}
					}
					b, _ := json.Marshal(dato_plan)
					str := string(b)
					subgrupo_detalle["dato_plan"] = str
				}
				if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
					panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
				}

			}
		}

	}
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": entrada}
	c.ServeJSON()

}

// DeleteActividad ...
// @Title DeleteActividad
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /delete_actividad/:id/:index [put]
func (c *FormulacionController) DeleteActividad() {
	id := c.Ctx.Input.Param(":id")
	index := c.Ctx.Input.Param(":index")

	var res map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		formulacionhelper.RecorrerHijos(hijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Actividades Inactivas"}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// GetAllActividades ...
// @Title GetAllActividades
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_all_actividades/:id/ [get]
func (c *FormulacionController) GetAllActividades() {
	id := c.Ctx.Input.Param(":id")
	var res map[string]interface{}
	var hijos []map[string]interface{}
	var tabla map[string]interface{}
	formulacionhelper.LimpiaTabla()
	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		for i := 0; i < len(hijos); i++ {
			if len(hijos[i]["hijos"].([]interface{})) != 0 {
				tabla = formulacionhelper.GetTabla(hijos[i]["hijos"].([]interface{}))
			}
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": tabla}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}
	c.ServeJSON()
}

// GetArbolArmonizacion ...
// @Title GetArbolArmonizacion
// @Description post Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_arbol_armonizacion/:id/ [post]
func (c *FormulacionController) GetArbolArmonizacion() {

	var entrada map[string][]string
	var arregloId []string
	var armonizacion []map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada)
	arregloId = entrada["Data"]
	for i := 0; i < len(arregloId); i++ {
		armonizacion = append(armonizacion, formulacionhelper.GetArmonizacion(arregloId[i]))
	}
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": armonizacion}
	c.ServeJSON()
}

// GuardarIdentificacion ...
// @Title GuardarIdentificacion
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	idTipo		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /guardar_identificacion/:id/:idTipo [put]
func (c *FormulacionController) GuardarIdentificacion() {
	id := c.Ctx.Input.Param(":id")
	tipoIdenti := c.Ctx.Input.Param(":idTipo")
	var entrada map[string]interface{}
	var res map[string]interface{}
	var resJ map[string]interface{}
	var respuesta []map[string]interface{}
	var idStr string
	var identificacion map[string]interface{}

	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada)

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id+",tipo_identificacion_id:"+tipoIdenti, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &respuesta)
		jsonString, _ := json.Marshal(respuesta[0]["_id"])
		json.Unmarshal(jsonString, &idStr)
		identificacion = respuesta[0]
		b, _ := json.Marshal(entrada)
		str := string(b)
		identificacion["dato"] = str
		if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/identificacion/"+idStr, "PUT", &resJ, identificacion); err != nil {
			panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error actualizando identificacion \"idStr\"", "status": "400", "log": err})
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Registro de identificación"}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()

}

// GetAllIdentificacion ...
// @Title GetAllIdentificacion
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	idTipo		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_all_identificacion/:id/:idTipo [get]
func (c *FormulacionController) GetAllIdentificacion() {
	id := c.Ctx.Input.Param(":id")
	tipoIdenti := c.Ctx.Input.Param(":idTipo")
	var respuesta []map[string]interface{}
	var res map[string]interface{}
	var identificacion map[string]interface{}
	var dato map[string]interface{}
	var data_identi []map[string]interface{}
	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id+",tipo_identificacion_id:"+tipoIdenti, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &respuesta)
		identificacion = respuesta[0]

		if identificacion["dato"] != nil {
			dato_str := identificacion["dato"].(string)
			json.Unmarshal([]byte(dato_str), &dato)
			for key := range dato {
				element := dato[key].(map[string]interface{})
				if element["activo"] == true {
					data_identi = append(data_identi, element)
				}
			}
		}

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data_identi}

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// DeleteIdentificacion ...
// @Title DeleteIdentificacion
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	idTipo		path 	string	true		"The key for staticblock"
// @Param	index		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /delete_identificacion/:id/:idTipo/:index [put]
func (c *FormulacionController) DeleteIdentificacion() {
	id := c.Ctx.Input.Param(":id")
	index := c.Ctx.Input.Param(":index")
	idTipo := c.Ctx.Input.Param(":idTipo")
	var idStr string
	var res map[string]interface{}
	var respuesta []map[string]interface{}
	var identificacion map[string]interface{}
	var dato map[string]interface{}
	var resJ map[string]interface{}

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id+",tipo_identificacion_id:"+idTipo, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &respuesta)
		identificacion = respuesta[0]

		jsonString, _ := json.Marshal(respuesta[0]["_id"])
		json.Unmarshal(jsonString, &idStr)

		if identificacion["dato"] != nil {
			dato_str := identificacion["dato"].(string)
			json.Unmarshal([]byte(dato_str), &dato)
			for key := range dato {
				intVar, _ := strconv.Atoi(key)
				intVar = intVar + 1
				strr := strconv.Itoa(intVar)
				if strr == index {
					element := dato[key].(map[string]interface{})
					element["activo"] = false
					dato[key] = element
				}
			}
			b, _ := json.Marshal(dato)
			str := string(b)
			identificacion["dato"] = str
		}
		if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/identificacion/"+idStr, "PUT", &resJ, identificacion); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteIdentificacion", "err": "Error eliminando identificacion \"idStr\"", "status": "400", "log": err})
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Identificación Inactiva"}
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}
}

// VersionarPlan ...
// @Title VersionarPlan
// @Description post Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /versionar_plan/:id [post]
func (c *FormulacionController) VersionarPlan() {

	id := c.Ctx.Input.Param(":id")

	var respuesta map[string]interface{}
	var respuestaHijos map[string]interface{}
	var respuestaIdentificacion map[string]interface{}
	var hijos []map[string]interface{}
	var identificaciones []map[string]interface{}
	var planPadre map[string]interface{}
	var respuestaPost map[string]interface{}
	var planVersionado map[string]interface{}
	plan := make(map[string]interface{})

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/plan/"+id, &respuesta); err == nil {

		helpers.LimpiezaRespuestaRefactor(respuesta, &planPadre)

		plan["nombre"] = planPadre["nombre"].(string)
		plan["descripcion"] = planPadre["descripcion"].(string)
		plan["tipo_plan_id"] = planPadre["tipo_plan_id"].(string)
		plan["aplicativo_id"] = planPadre["aplicativo_id"].(string)
		plan["activo"] = planPadre["activo"]
		plan["formato"] = false
		plan["vigencia"] = planPadre["vigencia"].(string)
		plan["dependencia_id"] = planPadre["dependencia_id"].(string)
		plan["estado_plan_id"] = "614d3ad301c7a200482fabfd"
		plan["padre_plan_id"] = id

		if err := helpers.SendJson(beego.AppConfig.String("PlanesService")+"/plan", "POST", &respuestaPost, plan); err != nil {
			panic(map[string]interface{}{"funcion": "VersionarPlan", "err": "Error versionando plan \"plan[\"_id\"].(string)\"", "status": "400", "log": err})
		}
		planVersionado = respuestaPost["Data"].(map[string]interface{})
		c.Data["json"] = respuestaPost

		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id, &respuestaIdentificacion); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaIdentificacion, &identificaciones)
			if len(identificaciones) != 0 {
				formulacionhelper.VersionarIdentificaciones(identificaciones, planVersionado["_id"].(string))
			}
		}

		if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &respuestaHijos); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijos, &hijos)
			formulacionhelper.VersionarHijos(hijos, planVersionado["_id"].(string))
		}

	}
	c.ServeJSON()
}

// GetPlanVersiones ...
// @Title GetPlanVersiones
// @Description get Formulacion by id
// @Param	unidad		path 	string	true		"The key for staticblock"
// @Param	vigencia		path 	string	true		"The key for staticblock"
// @Param	nombre		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_plan_versiones/:unidad/:vigencia/:nombre [get]
func (c *FormulacionController) GetPlanVersiones() {
	unidad := c.Ctx.Input.Param(":unidad")
	vigencia := c.Ctx.Input.Param(":vigencia")
	nombre := c.Ctx.Input.Param(":nombre")

	var respuesta map[string]interface{}
	var versiones []map[string]interface{}

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/plan?query=dependencia_id:"+unidad+",vigencia:"+vigencia+",formato:false,nombre:"+nombre, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &versiones)
		versionesOrdenadas := formulacionhelper.OrdenarVersiones(versiones)
		c.Data["json"] = versionesOrdenadas

	}
	c.ServeJSON()
}

// GetPonderacionActividades ...
// @Title GetPonderacionActividades
// @Description get Formulacion by id
// @Param	plan		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /ponderacion_actividades/:plan [get]
func (c *FormulacionController) PonderacionActividades() {
	plan := c.Ctx.Input.Param(":plan")
	var respuesta map[string]interface{}
	var respuestaDetalle map[string]interface{}
	var respuestaLimpiaDetalle []map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var hijos []map[string]interface{}

	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &hijos)

		for i := 0; i < len(hijos); i++ {
			if strings.Contains(strings.ToUpper(hijos[i]["nombre"].(string)), "PONDERACIÓN") && strings.Contains(strings.ToUpper(hijos[i]["nombre"].(string)), "ACTIVIDAD") || strings.Contains(strings.ToUpper(hijos[i]["nombre"].(string)), "PONDERACIÓN") {
				if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijos[i]["_id"].(string), &respuestaDetalle); err == nil {

					helpers.LimpiezaRespuestaRefactor(respuestaDetalle, &respuestaLimpiaDetalle)
					subgrupoDetalle = respuestaLimpiaDetalle[0]

					if subgrupoDetalle["dato_plan"] != nil {
						var suma float64 = 0
						datoPlan := make(map[string]map[string]interface{})
						json.Unmarshal([]byte(subgrupoDetalle["dato_plan"].(string)), &datoPlan)

						ponderacionActividades := make(map[string]interface{})

						for j := 1; j <= len(datoPlan); j++ {
							if datoPlan[strconv.Itoa(j)]["activo"] != false {
								ponderacionActividades["Actividad "+strconv.Itoa(j)] = datoPlan[strconv.Itoa(j)]["dato"]
								suma += datoPlan[strconv.Itoa(j)]["dato"].(float64)
							}

						}

						ponderacionActividades["Total"] = suma
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ponderacionActividades}
					}

				} else {
					panic(map[string]interface{}{"funcion": "PonderacionActividades", "err": "Error subgrupo_detalle plan \"plan\"", "status": "400", "log": err})
				}
			}
		}
	} else {
		panic(map[string]interface{}{"funcion": "PonderacionActividades", "err": "Error subgrupo_hijos plan \"plan\"", "status": "400", "log": err})
	}

	c.ServeJSON()
}

// GetRubros ...
// @Title GetRubros
// @Description get Rubros
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_rubros [get]
func (c *FormulacionController) GetRubros() {

	var respuesta map[string]interface{}
	var rubros []interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanCuentasService")+"/arbol_rubro", &respuesta); err != nil {
		panic(map[string]interface{}{"funcion": "GetRubros", "err": "Error arbol_rubros", "status": "400", "log": err})
	} else {
		rubros = respuesta["Body"].([]interface{})
		for i := 0; i < len(rubros); i++ {
			if rubros[i].(map[string]interface{})["Nombre"] == "GASTOS" {
				aux := rubros[i].(map[string]interface{})
				hojas := formulacionhelper.GetHijosRubro(aux["Hijos"].([]interface{}))

				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": hojas}
				break
			}
		}
	}
	c.ServeJSON()
}

// GetUnidades ...
// @Title GetUnidades
// @Description get Unidades
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_unidades [get]
func (c *FormulacionController) GetUnidades() {

	var respuesta []map[string]interface{}
	var unidades []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:2&limit=0", &respuesta); err == nil {
		for i := 0; i < len(respuesta); i++ {
			unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
		}
		respuesta = nil

		if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:3&limit=0", &respuesta); err == nil {
			for i := 0; i < len(respuesta); i++ {
				unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
			}
			respuesta = nil

			if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:4&limit=0", &respuesta); err == nil {
				for i := 0; i < len(respuesta); i++ {
					unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
				}
				respuesta = nil

				if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:5&limit=0", &respuesta); err == nil {
					for i := 0; i < len(respuesta); i++ {
						unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
					}
					respuesta = nil
					if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:6&limit=0", &respuesta); err == nil {
						for i := 0; i < len(respuesta); i++ {
							unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
						}
						respuesta = nil
						if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:7&limit=0", &respuesta); err == nil {
							for i := 0; i < len(respuesta); i++ {
								unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
							}
							respuesta = nil

							if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:8&limit=0", &respuesta); err == nil {
								for i := 0; i < len(respuesta); i++ {
									unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
								}
								respuesta = nil

								if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:13&limit=0", &respuesta); err == nil {
									for i := 0; i < len(respuesta); i++ {
										unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
									}
									respuesta = nil

									if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:15&limit=0", &respuesta); err == nil {
										for i := 0; i < len(respuesta); i++ {
											aux := respuesta[i]["DependenciaId"]
											if strings.Contains(aux.(map[string]interface{})["Nombre"].(string), "DOCTORADO") {
												unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
											}
										}
										respuesta = nil

										if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:11&limit=0", &respuesta); err == nil {
											for i := 0; i < len(respuesta); i++ {
												unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
											}
											respuesta = nil

											if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:28&limit=0", &respuesta); err == nil {
												for i := 0; i < len(respuesta); i++ {
													unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
												}
												respuesta = nil

												if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:33&limit=0", &respuesta); err == nil {
													for i := 0; i < len(respuesta); i++ {
														unidades = append(unidades, respuesta[i]["DependenciaId"].(map[string]interface{}))
													}
													respuesta = nil
													c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": unidades}

												} else {
													panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
												}

											} else {
												panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
											}
										} else {
											panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
										}
									} else {
										panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
									}
								} else {
									panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
								}
							} else {
								panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
							}
						} else {
							panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
						}
					} else {
						panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
					}
				} else {
					panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
				}
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}

		} else {
			panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
		}

	} else {
		panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
	}

	c.ServeJSON()
}
