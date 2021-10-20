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
	fmt.Println(id)
	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada)
	for key, element := range entrada {
		var respuesta map[string]interface{}
		var respuestaLimpia []map[string]interface{}
		var subgrupo_detalle map[string]interface{}
		dato_plan := make(map[string]interface{})
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
	id := c.Ctx.Input.Param(":id")
	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada)
	arregloId = entrada["Data"]
	fmt.Println(id)
	for i := 0; i < len(arregloId); i++ {
		armonizacion = append(armonizacion, formulacionhelper.GetArmonizacion(arregloId[i]))
	}
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": armonizacion}
	c.ServeJSON()
}
