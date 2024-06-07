package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/models"
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
	c.Mapping("VinculacionTercero", c.VinculacionTercero)
	c.Mapping("Planes", c.Planes)
	c.Mapping("VerificarIdentificaciones", c.VerificarIdentificaciones)
	c.Mapping("PlanesEnFormulacion", c.PlanesEnFormulacion)
	c.Mapping("GetPlanesUnidadesComun", c.GetPlanesUnidadesComun)
	c.Mapping("DefinirFechasFuncionamiento", c.DefinirFechasFuncionamiento)
	c.Mapping("CalculosDocentes", c.CalculosDocentes)
	c.Mapping("EstructuraPlanes", c.EstructuraPlanes)
	c.Mapping("VinculacionTerceroByEmail", c.VinculacionTerceroByEmail)
	c.Mapping("CambioCargoIdVinculacionTercero", c.CambioCargoIdVinculacionTercero)
	c.Mapping("VinculacionTerceroByIdentificacion", c.VinculacionTerceroByIdentificacion)
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

	var respuesta map[string]interface{}
	var respuestaHijos map[string]interface{}
	var hijos []map[string]interface{}
	var planFormato map[string]interface{}
	var parametros map[string]interface{}

	plan := make(map[string]interface{})
	clienteHttp := &http.Client{}
	url := "http://" + beego.AppConfig.String("PlanesService") + "/plan/"

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &respuesta); err == nil {

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
		plan["formato_id"] = id
		plan["nueva_estructura"] = true

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

		cuerpoRespuesta, err := io.ReadAll(respuesta.Body)
		if err != nil {
			log.Fatalf("Error leyendo peticion: %v", err)
		}

		json.Unmarshal(cuerpoRespuesta, &resPost)
		resLimpia = resPost["Data"].(map[string]interface{})
		padre := resLimpia["_id"].(string)
		c.Data["json"] = resPost

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &respuestaHijos); err == nil {
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
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_plan/:id/:index [get]
func (c *FormulacionController) GetPlan() {
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
	var entrada map[string]interface{}
	var body map[string]interface{}

	_ = id
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	entrada = body["entrada"].(map[string]interface{})
	armonizacion := body["armo"]
	armonizacionPI := body["armoPI"]
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

				if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
					panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
				}

			}
			continue
		}
		id_subgrupoDetalle = key
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+id_subgrupoDetalle, &respuesta); err != nil {
			panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
		}
		helpers.LimpiezaRespuestaRefactor(respuesta, &respuestaLimpia)

		subgrupo_detalle = respuestaLimpia[0]
		if subgrupo_detalle["armonizacion_dato"] != nil {
			dato_armonizacion_str := subgrupo_detalle["armonizacion_dato"].(string)
			json.Unmarshal([]byte(dato_armonizacion_str), &armonizacion_dato)
			if armonizacion_dato[index] != nil {
				aux := make(map[string]interface{})
				aux["armonizacionPED"] = armonizacion
				aux["armonizacionPI"] = armonizacionPI
				armonizacion_dato[index] = aux
			}
			c, _ := json.Marshal(armonizacion_dato)
			strArmonizacion := string(c)
			subgrupo_detalle["armonizacion_dato"] = strArmonizacion

		}

		nuevoDato := true
		actividad := make(map[string]interface{})

		if subgrupo_detalle["dato_plan"] != nil {
			dato_plan_str := subgrupo_detalle["dato_plan"].(string)
			json.Unmarshal([]byte(dato_plan_str), &dato_plan)

			for index_actividad := range dato_plan {
				if index_actividad == index {
					nuevoDato = false
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
		}

		if nuevoDato {
			actividad["index"] = index
			actividad["dato"] = element
			actividad["activo"] = true
			dato_plan[index] = actividad
		}

		b, _ := json.Marshal(dato_plan)
		str := string(b)
		subgrupo_detalle["dato_plan"] = str

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/"+subgrupo_detalle["_id"].(string), "PUT", &res, subgrupo_detalle); err != nil {
			panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error actualizando subgrupo-detalle \"subgrupo_detalle[\"_id\"].(string)\"", "status": "400", "log": err})
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
		formulacionhelper.RecorrerHijos(hijos, index)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Actividades Inactivas"}
	} else {
		panic(err)
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
	var res map[string]interface{}
	var hijos []map[string]interface{}
	var tabla map[string]interface{}
	var auxHijos []interface{}
	formulacionhelper.LimpiaTabla()
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &hijos)
		for i := 0; i < len(hijos); i++ {
			auxHijos = append(auxHijos, hijos[i]["_id"])
		}
		tabla = formulacionhelper.GetTabla(auxHijos)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": tabla}
	} else {
		panic(err)
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
	tipoIdenti := c.Ctx.Input.Param(":idTipo")
	var entrada map[string]interface{}
	var res map[string]interface{}
	var resJ map[string]interface{}
	var respuesta []map[string]interface{}
	var idStr string
	var identificacion map[string]interface{}

	json.Unmarshal(c.Ctx.Input.RequestBody, &entrada)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id+",tipo_identificacion_id:"+tipoIdenti, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &respuesta)

		if tipoIdenti == "61897518f6fc97091727c3c3" { // ? Recurso docente unicamente
			if len(respuesta) > 0 {
				if strings.Contains(respuesta[0]["dato"].(string), "ids_detalle") {
					identificacion = respuesta[0]
					dato_json := map[string]interface{}{}
					json.Unmarshal([]byte(identificacion["dato"].(string)), &dato_json)

					iddetail := ""
					identificacionDetalle := map[string]interface{}{}
					errIdentificacionDetalle := error(nil)

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rhf"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						identificacionDetallePut := identificacionDetalle["Data"].(map[string]interface{})
						identificacionDetallePut["dato"] = entrada["rhf"]
						identificacionDetalle = map[string]interface{}{}
						errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, "PUT", &identificacionDetalle, identificacionDetallePut)
						if errIdentificacionDetalle != nil || identificacionDetalle["Status"] != "200" {
							panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error actualizando detalle identificacion \"rhf\"", "status": "400", "log": err})
						}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error consultando detalle identificacion \"rhf\"", "status": "400", "log": err})
					}

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rhv_pre"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						identificacionDetallePut := identificacionDetalle["Data"].(map[string]interface{})
						identificacionDetallePut["dato"] = entrada["rhv_pre"]
						identificacionDetalle = map[string]interface{}{}
						errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, "PUT", &identificacionDetalle, identificacionDetallePut)
						if errIdentificacionDetalle != nil || identificacionDetalle["Status"] != "200" {
							panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error actualizando detalle identificacion \"rhv_pre\"", "status": "400", "log": err})
						}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error consultando detalle identificacion \"rhv_pre\"", "status": "400", "log": err})
					}

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rhv_pos"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						identificacionDetallePut := identificacionDetalle["Data"].(map[string]interface{})
						identificacionDetallePut["dato"] = entrada["rhv_pos"]
						identificacionDetalle = map[string]interface{}{}
						errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, "PUT", &identificacionDetalle, identificacionDetallePut)
						if errIdentificacionDetalle != nil || identificacionDetalle["Status"] != "200" {
							panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error actualizando detalle identificacion \"rhv_pos\"", "status": "400", "log": err})
						}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error consultando detalle identificacion \"rhv_pos\"", "status": "400", "log": err})
					}

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rubros"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						identificacionDetallePut := identificacionDetalle["Data"].(map[string]interface{})
						identificacionDetallePut["dato"] = entrada["rubros"]
						identificacionDetalle = map[string]interface{}{}
						errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, "PUT", &identificacionDetalle, identificacionDetallePut)
						if errIdentificacionDetalle != nil || identificacionDetalle["Status"] != "200" {
							panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error actualizando detalle identificacion \"rubros\"", "status": "400", "log": err})
						}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error consultando detalle identificacion \"rubros\"", "status": "400", "log": err})
					}

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rubros_pos"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						identificacionDetallePut := identificacionDetalle["Data"].(map[string]interface{})
						identificacionDetallePut["dato"] = entrada["rubros_pos"]
						identificacionDetalle = map[string]interface{}{}
						errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, "PUT", &identificacionDetalle, identificacionDetallePut)
						if errIdentificacionDetalle != nil || identificacionDetalle["Status"] != "200" {
							panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error actualizando detalle identificacion \"rubros_pos\"", "status": "400", "log": err})
						}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error consultando detalle identificacion \"rubros_pos\"", "status": "400", "log": err})
					}

					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Registro de identificación"}

				} else {
					// ? -- Transicional mientras migran datos --
					identificacion = respuesta[0]
					detalles := map[string]interface{}{
						"rhf":        "",
						"rhv_pre":    "",
						"rhv_pos":    "",
						"rubros":     "",
						"rubros_pos": "",
					}
					identificacionDetalle := map[string]interface{}{}
					errIdentificacionDetalle := error(nil)

					data := map[string]interface{}{
						"dato": entrada["rhf"],
					}
					errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/", "POST", &identificacionDetalle, data)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "201" {
						detalles["rhf"] = identificacionDetalle["Data"].(map[string]interface{})["_id"].(string)
						identificacionDetalle = map[string]interface{}{}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error guardando detalle identificacion \"rhf\"", "status": "400", "log": err})
					}

					data = map[string]interface{}{
						"dato": entrada["rhv_pre"],
					}
					errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/", "POST", &identificacionDetalle, data)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "201" {
						detalles["rhv_pre"] = identificacionDetalle["Data"].(map[string]interface{})["_id"].(string)
						identificacionDetalle = map[string]interface{}{}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error guardando detalle identificacion \"rhv_pre\"", "status": "400", "log": err})
					}

					data = map[string]interface{}{
						"dato": entrada["rhv_pos"],
					}
					errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/", "POST", &identificacionDetalle, data)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "201" {
						detalles["rhv_pos"] = identificacionDetalle["Data"].(map[string]interface{})["_id"].(string)
						identificacionDetalle = map[string]interface{}{}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error guardando detalle identificacion \"rhv_pos\"", "status": "400", "log": err})
					}

					data = map[string]interface{}{
						"dato": entrada["rubros"],
					}
					errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/", "POST", &identificacionDetalle, data)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "201" {
						detalles["rubros"] = identificacionDetalle["Data"].(map[string]interface{})["_id"].(string)
						identificacionDetalle = map[string]interface{}{}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error guardando detalle identificacion \"rubros\"", "status": "400", "log": err})
					}

					data = map[string]interface{}{
						"dato": entrada["rubros_pos"],
					}
					errIdentificacionDetalle = helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/", "POST", &identificacionDetalle, data)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "201" {
						detalles["rubros_pos"] = identificacionDetalle["Data"].(map[string]interface{})["_id"].(string)
						identificacionDetalle = map[string]interface{}{}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error guardando detalle identificacion \"rubros_pos\"", "status": "400", "log": err})
					}

					bt, _ := json.Marshal(map[string]interface{}{"ids_detalle": detalles})
					identificacion["dato"] = string(bt)

					identificacionAns := map[string]interface{}{}
					errIdentificacion := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion/"+identificacion["_id"].(string), "PUT", &identificacionAns, identificacion)
					if errIdentificacion == nil && identificacionAns["Status"] == "200" {
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Registro de identificación"}
					} else {
						panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error actualizando identificacion", "status": "400", "log": err})
					}

				}
			} else {
				panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error sin dato identificacion", "status": "400", "log": err})
			}
		} else {
			jsonString, _ := json.Marshal(respuesta[0]["_id"])
			json.Unmarshal(jsonString, &idStr)

			identificacion = respuesta[0]
			b, _ := json.Marshal(entrada)
			str := string(b)

			identificacion["dato"] = str
			if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion/"+idStr, "PUT", &resJ, identificacion); err != nil {
				panic(map[string]interface{}{"funcion": "GuardarIdentificacion", "err": "Error actualizando identificacion \"idStr\"", "status": "400", "log": err})
			}
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Registro de identificación"}
		}

	} else {
		panic(err)
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
	tipoIdenti := c.Ctx.Input.Param(":idTipo")
	var respuesta []map[string]interface{}
	var res map[string]interface{}
	var identificacion map[string]interface{}
	var dato map[string]interface{}
	var data_identi []map[string]interface{}
	if tipoIdenti == "61897518f6fc97091727c3c3" { // ? Recurso docente unicamente
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id+",tipo_identificacion_id:"+tipoIdenti, &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &respuesta)
			if len(respuesta) > 0 {
				if strings.Contains(respuesta[0]["dato"].(string), "ids_detalle") { // ? Nuevo método fraccionado
					identificacion = respuesta[0]
					dato_json := map[string]interface{}{}
					json.Unmarshal([]byte(identificacion["dato"].(string)), &dato_json)
					result := dato_json["ids_detalle"].(map[string]interface{})

					iddetail := ""
					identificacionDetalle := map[string]interface{}{}
					errIdentificacionDetalle := error(nil)

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rhf"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_str := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						dato := map[string]interface{}{}
						datos := []map[string]interface{}{}
						json.Unmarshal([]byte(dato_str), &dato)
						for key := range dato {
							element := dato[key].(map[string]interface{})
							if element["activo"] == true {
								datos = append(datos, element)
							}
						}
						if len(datos) > 0 {
							result["rhf"] = datos
						} else {
							result["rhf"] = "{}"
						}
					} else {
						result["rhf"] = "{}"
					}

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rhv_pre"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_str := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						dato := map[string]interface{}{}
						datos := []map[string]interface{}{}
						json.Unmarshal([]byte(dato_str), &dato)
						for key := range dato {
							element := dato[key].(map[string]interface{})
							if element["activo"] == true {
								datos = append(datos, element)
							}
						}
						if len(datos) > 0 {
							result["rhv_pre"] = datos
						} else {
							result["rhv_pre"] = "{}"
						}
					} else {
						result["rhv_pre"] = "{}"
					}

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rhv_pos"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_str := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						dato := map[string]interface{}{}
						datos := []map[string]interface{}{}
						json.Unmarshal([]byte(dato_str), &dato)
						for key := range dato {
							element := dato[key].(map[string]interface{})
							if element["activo"] == true {
								datos = append(datos, element)
							}
						}
						if len(datos) > 0 {
							result["rhv_pos"] = datos
						} else {
							result["rhv_pos"] = "{}"
						}
					} else {
						result["rhv_pos"] = "{}"
					}

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rubros"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_str := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						dato := map[string]interface{}{}
						datos := []map[string]interface{}{}
						json.Unmarshal([]byte(dato_str), &dato)
						for key := range dato {
							element := dato[key].(map[string]interface{})
							if element["activo"] == true {
								datos = append(datos, element)
							}
						}
						if len(datos) > 0 {
							result["rubros"] = datos
						} else {
							result["rubros"] = "{}"
						}
					} else {
						result["rubros"] = "{}"
					}

					iddetail = dato_json["ids_detalle"].(map[string]interface{})["rubros_pos"].(string)
					errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
					if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
						dato_str := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
						dato := map[string]interface{}{}
						datos := []map[string]interface{}{}
						json.Unmarshal([]byte(dato_str), &dato)
						for key := range dato {
							element := dato[key].(map[string]interface{})
							if element["activo"] == true {
								datos = append(datos, element)
							}
						}
						if len(datos) > 0 {
							result["rubros_pos"] = datos
						} else {
							result["rubros_pos"] = "{}"
						}
					} else {
						result["rubros_pos"] = "{}"
					}

					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": result}

				} else {
					identificacion = respuesta[0]
					if identificacion["dato"] != nil && identificacion["dato"] != "{}" { // ? Antiguo método unificado
						result := make(map[string]interface{})
						dato_str := identificacion["dato"].(string)
						json.Unmarshal([]byte(dato_str), &dato)

						var identi map[string]interface{} = nil
						dato_aux := ""

						dato_aux = dato["rhf"].(string)
						if dato_aux == "{}" {
							result["rhf"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rhf"] = data_identi
						}

						identi = nil
						data_identi = nil

						dato_aux = dato["rhv_pre"].(string)
						if dato_aux == "{}" {
							result["rhv_pre"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rhv_pre"] = data_identi
						}

						identi = nil
						data_identi = nil

						dato_aux = dato["rhv_pos"].(string)
						if dato_aux == "{}" {
							result["rhv_pos"] = "{}"
						} else {
							json.Unmarshal([]byte(dato_aux), &identi)
							for key := range identi {
								element := identi[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							result["rhv_pos"] = data_identi
						}

						identi = nil
						data_identi = nil

						if dato["rubros"] != nil {
							dato_aux = dato["rubros"].(string)
							if dato_aux == "{}" {
								result["rubros"] = "{}"
							} else {
								json.Unmarshal([]byte(dato_aux), &identi)
								for key := range identi {
									element := identi[key].(map[string]interface{})
									if element["activo"] == true {
										data_identi = append(data_identi, element)
									}
								}
								result["rubros"] = data_identi
							}
						}

						identi = nil
						data_identi = nil

						if dato["rubros_pos"] != nil {
							dato_aux = dato["rubros_pos"].(string)
							if dato_aux == "{}" {
								result["rubros_pos"] = "{}"
							} else {
								json.Unmarshal([]byte(dato_aux), &identi)
								for key := range identi {
									element := identi[key].(map[string]interface{})
									if element["activo"] == true {
										data_identi = append(data_identi, element)
									}
								}
								result["rubros_pos"] = data_identi
							}
						}

						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": result}

					} else {
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}
					}
				}
			} else {
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}
			}
		} else {
			panic(err)
		}
	} else {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id+",tipo_identificacion_id:"+tipoIdenti, &res); err == nil {
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

				sort.SliceStable(data_identi, func(i, j int) bool {
					a, _ := strconv.Atoi(data_identi[i]["index"].(string))
					b, _ := strconv.Atoi(data_identi[j]["index"].(string))
					return a < b
				})

				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data_identi}

			} else {
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}
			}

		} else {
			panic(err)
		}

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
	idTipo := c.Ctx.Input.Param(":idTipo")
	var idStr string
	var res map[string]interface{}
	var respuesta []map[string]interface{}
	var identificacion map[string]interface{}
	var dato map[string]interface{}
	var resJ map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id+",tipo_identificacion_id:"+idTipo, &res); err == nil {
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
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion/"+idStr, "PUT", &resJ, identificacion); err != nil {
			panic(map[string]interface{}{"funcion": "DeleteIdentificacion", "err": "Error eliminando identificacion \"idStr\"", "status": "400", "log": err})
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "Identificación Inactiva"}
	} else {
		panic(err)
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

	var respuesta map[string]interface{}
	var respuestaHijos map[string]interface{}
	var respuestaIdentificacion map[string]interface{}
	var hijos []map[string]interface{}
	var identificaciones []map[string]interface{}
	var planPadre map[string]interface{}
	var respuestaPost map[string]interface{}
	var planVersionado map[string]interface{}
	plan := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &respuesta); err == nil {

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
		if value, ok := planPadre["formato_id"].(string); ok {
			plan["formato_id"] = value
		}

		if _, ok := planPadre["nueva_estructura"]; ok {
			plan["nueva_estructura"] = true
		}

		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan", "POST", &respuestaPost, plan); err != nil {
			panic(map[string]interface{}{"funcion": "VersionarPlan", "err": "Error versionando plan \"plan[\"_id\"].(string)\"", "status": "400", "log": err})
		}
		planVersionado = respuestaPost["Data"].(map[string]interface{})
		c.Data["json"] = respuestaPost

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id, &respuestaIdentificacion); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaIdentificacion, &identificaciones)
			if len(identificaciones) != 0 {
				formulacionhelper.VersionarIdentificaciones(identificaciones, planVersionado["_id"].(string))
			}
		}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+id, &respuestaHijos); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaHijos, &hijos)
			formulacionhelper.VersionarHijos(hijos, planVersionado["_id"].(string))
		}

		var resPadres map[string]interface{}
		var planesPadre []map[string]interface{}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=dependencia_id:"+plan["dependencia_id"].(string)+",vigencia:"+plan["vigencia"].(string)+",formato:false,nombre:"+url.QueryEscape(plan["nombre"].(string)), &resPadres); err == nil {
			helpers.LimpiezaRespuestaRefactor(resPadres, &planesPadre)
			for _, padre := range planesPadre {
				var resActualizacion map[string]interface{}
				if padre["_id"].(string) != planVersionado["_id"].(string) {
					padre["activo"] = false
					helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+padre["_id"].(string), "PUT", &resActualizacion, padre)
				}
			}
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
	unidad := c.Ctx.Input.Param(":unidad")
	vigencia := c.Ctx.Input.Param(":vigencia")
	nombre := c.Ctx.Input.Param(":nombre")

	var respuesta map[string]interface{}
	var versiones []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=dependencia_id:"+unidad+",vigencia:"+vigencia+",formato:false,nombre:"+nombre, &respuesta); err == nil {
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
	plan := c.Ctx.Input.Param(":plan")
	var respuesta map[string]interface{}
	var respuestaDetalle map[string]interface{}
	var respuestaLimpiaDetalle []map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var hijos []map[string]interface{}
	var hijosFiltrado []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &hijos)
		for i := 0; i < len(hijos); i++ {
			if hijos[i]["activo"] == true {
				hijosFiltrado = append(hijosFiltrado, hijos[i])
			}
		}

		for i := 0; i < len(hijosFiltrado); i++ {
			if strings.Contains(strings.ToUpper(hijosFiltrado[i]["nombre"].(string)), "PONDERACIÓN") && strings.Contains(strings.ToUpper(hijosFiltrado[i]["nombre"].(string)), "ACTIVIDAD") || strings.Contains(strings.ToUpper(hijosFiltrado[i]["nombre"].(string)), "PONDERACIÓN") {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijosFiltrado[i]["_id"].(string), &respuestaDetalle); err == nil {
					helpers.LimpiezaRespuestaRefactor(respuestaDetalle, &respuestaLimpiaDetalle)
					subgrupoDetalle = respuestaLimpiaDetalle[0]

					if subgrupoDetalle["dato_plan"] != nil {
						var suma float64 = 0
						datoPlan := make(map[string]map[string]interface{})
						json.Unmarshal([]byte(subgrupoDetalle["dato_plan"].(string)), &datoPlan)

						ponderacionActividades := make(map[string]interface{})

						for j, dato := range datoPlan {
							if dato["activo"] != false && len(dato) != 0 {
								ponderacionActividades["Actividad "+(j)] = dato["dato"]
								suma += dato["dato"].(float64)
								suma = math.Round(suma*100) / 100
							}
						}

						ponderacionActividades["Total"] = suma
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ponderacionActividades}
					}

				} else {
					panic(map[string]interface{}{"funcion": "PonderacionActividades", "err": "Error subgrupo_detalle plan \"plan\"", "status": "400", "log": err})
				}
			} else {
				c.Data["json"] = map[string]interface{}{"Success": false, "Status": "400", "Message": "Error-tipo1", "Data": `El formato usado para definir la actividad, debe de tener como primer item las palabras "Ponderación" y "Actividad" o solo "Ponderación"`}
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

	var respuesta map[string]interface{}
	var rubros []interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanCuentasService")+"/arbol_rubro", &respuesta); err != nil {
		panic(map[string]interface{}{"funcion": "GetRubros", "err": "Error arbol_rubros", "status": "400", "log": err})
	} else {
		rubros = respuesta["Body"].([]interface{})
		for i := 0; i < len(rubros); i++ {
			if strings.ToUpper(rubros[i].(map[string]interface{})["Nombre"].(string)) == "GASTOS" {
				aux := rubros[i].(map[string]interface{})
				hojas := formulacionhelper.GetHijosRubro(aux["Hijos"].([]interface{}))
				if len(hojas) != 0 {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": hojas}
				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}
				}
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

	var respuesta []map[string]interface{}
	var unidades []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:2&limit=0", &respuesta); err == nil {
		for i := 0; i < len(respuesta); i++ {
			aux := respuesta[i]["DependenciaId"].(map[string]interface{})
			aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
			unidades = append(unidades, aux)
		}
		respuesta = nil

		if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:3&limit=0", &respuesta); err == nil {
			for i := 0; i < len(respuesta); i++ {
				aux := respuesta[i]["DependenciaId"].(map[string]interface{})
				aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
				unidades = append(unidades, aux)
			}
			respuesta = nil

			if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:4&limit=0", &respuesta); err == nil {
				for i := 0; i < len(respuesta); i++ {
					aux := respuesta[i]["DependenciaId"].(map[string]interface{})
					aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
					unidades = append(unidades, aux)
				}
				respuesta = nil

				if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:5&limit=0", &respuesta); err == nil {
					for i := 0; i < len(respuesta); i++ {
						aux := respuesta[i]["DependenciaId"].(map[string]interface{})
						aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
						unidades = append(unidades, aux)
					}
					respuesta = nil
					if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:6&limit=0", &respuesta); err == nil {
						for i := 0; i < len(respuesta); i++ {
							aux := respuesta[i]["DependenciaId"].(map[string]interface{})
							aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
							unidades = append(unidades, aux)
						}
						respuesta = nil
						if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:7&limit=0", &respuesta); err == nil {
							for i := 0; i < len(respuesta); i++ {
								aux := respuesta[i]["DependenciaId"].(map[string]interface{})
								aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
								unidades = append(unidades, aux)
							}
							respuesta = nil

							if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:8&limit=0", &respuesta); err == nil {
								for i := 0; i < len(respuesta); i++ {
									aux := respuesta[i]["DependenciaId"].(map[string]interface{})
									aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
									unidades = append(unidades, aux)
								}
								respuesta = nil

								if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:13&limit=0", &respuesta); err == nil {
									for i := 0; i < len(respuesta); i++ {
										aux := respuesta[i]["DependenciaId"].(map[string]interface{})
										aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
										unidades = append(unidades, aux)
									}
									respuesta = nil

									if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:15&limit=0", &respuesta); err == nil {
										for i := 0; i < len(respuesta); i++ {
											aux := respuesta[i]["DependenciaId"]
											if strings.Contains(aux.(map[string]interface{})["Nombre"].(string), "DOCTORADO") {
												aux := respuesta[i]["DependenciaId"].(map[string]interface{})
												aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
												unidades = append(unidades, aux)
											}
										}
										respuesta = nil

										if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:11&limit=0", &respuesta); err == nil {
											for i := 0; i < len(respuesta); i++ {
												aux := respuesta[i]["DependenciaId"].(map[string]interface{})
												aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
												unidades = append(unidades, aux)
											}
											respuesta = nil

											if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:28&limit=0", &respuesta); err == nil {
												for i := 0; i < len(respuesta); i++ {
													aux := respuesta[i]["DependenciaId"].(map[string]interface{})
													aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
													unidades = append(unidades, aux)
												}
												respuesta = nil

												if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:33&limit=0", &respuesta); err == nil {
													for i := 0; i < len(respuesta); i++ {
														aux := respuesta[i]["DependenciaId"].(map[string]interface{})
														aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
														unidades = append(unidades, aux)
													}
													respuesta = nil

													if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:14,DependenciaId:171&limit=0", &respuesta); err == nil {
														for i := 0; i < len(respuesta); i++ {
															aux := respuesta[i]["DependenciaId"].(map[string]interface{})
															aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
															unidades = append(unidades, aux)
														}
														respuesta = nil

														if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:10,DependenciaId:96&limit=0", &respuesta); err == nil {
															for i := 0; i < len(respuesta); i++ {
																aux := respuesta[i]["DependenciaId"].(map[string]interface{})
																aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
																unidades = append(unidades, aux)
															}
															respuesta = nil
															c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": unidades}

														} else {
															panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
														}

														if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:33,DependenciaId:222&limit=0", &respuesta); err == nil {
															for i := 0; i < len(respuesta); i++ {
																aux := respuesta[i]["DependenciaId"].(map[string]interface{})
																aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
																unidades = append(unidades, aux)
															}
															respuesta = nil
															c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": unidades}

														} else {
															panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
														}

														if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:10,DependenciaId:97&limit=0", &respuesta); err == nil {
															for i := 0; i < len(respuesta); i++ {
																aux := respuesta[i]["DependenciaId"].(map[string]interface{})
																aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
																unidades = append(unidades, aux)
															}
															respuesta = nil
															c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": unidades}

														} else {
															panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
														}

														if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:10,DependenciaId:209&limit=0", &respuesta); err == nil {
															for i := 0; i < len(respuesta); i++ {
																aux := respuesta[i]["DependenciaId"].(map[string]interface{})
																aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
																unidades = append(unidades, aux)
															}
															respuesta = nil
															c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": unidades}

														} else {
															panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
														}

														if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:10,DependenciaId:92&limit=0", &respuesta); err == nil {
															for i := 0; i < len(respuesta); i++ {
																aux := respuesta[i]["DependenciaId"].(map[string]interface{})
																aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
																unidades = append(unidades, aux)
															}
															respuesta = nil
															c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": unidades}

														} else {
															panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
														}

														if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=TipoDependenciaId:14,DependenciaId:42&limit=0", &respuesta); err == nil {
															for i := 0; i < len(respuesta); i++ {
																aux := respuesta[i]["DependenciaId"].(map[string]interface{})
																aux["TipoDependencia"] = respuesta[i]["TipoDependenciaId"]
																unidades = append(unidades, aux)
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
										panic(map[string]interface{}{"fuplann": "GetUnidades", "err": "Error ", "status": "400", "log": err})
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
}

// VinculacionTercero ...
// @Title VinculacionTercero
// @Description get VinculacionTercero
// @Param	tercero_id	path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /vinculacion_tercero/:tercero_id [get]
func (c *FormulacionController) VinculacionTercero() {

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

	terceroId := c.Ctx.Input.Param(":tercero_id")
	var vinculaciones []models.Vinculacion
	if err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/vinculacion?query=Activo:true,TerceroPrincipalId:"+terceroId, &vinculaciones); err != nil {
		panic(map[string]interface{}{"funcion": "VinculacionTercero", "err": "Error get vinculacion", "status": "400", "log": err})
	} else {
		for i := 0; i < len(vinculaciones); i++ {
			if vinculaciones[i].CargoId == 319 || vinculaciones[i].CargoId == 312 || vinculaciones[i].CargoId == 320 {
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": vinculaciones[i]}
				break
			} else {
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}
			}
		}
	}
	c.ServeJSON()
}

// VinculacionTerceroByEmail ...
// @Title VinculacionTerceroByEmail
// @Description get VinculacionTerceroByEmail
// @Param	tercero_email	path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /vinculacion_tercero_email/:tercero_email [get]
func (c *FormulacionController) VinculacionTerceroByEmail() {

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

	terceroEmail := c.Ctx.Input.Param(":tercero_email")
	var vinculaciones []models.Vinculacion

	idNoRegistra, idJefeOficina, idAsistenteDependencia, err := formulacionhelper.ObtenerIdParametros()
	if err != nil {
		panic(map[string]interface{}{"funcion": "VinculacionTercero", "err": "Error get parametros", "status": "400", "log": err})
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/vinculacion?query=Activo:true,TerceroPrincipalId__UsuarioWSO2:"+terceroEmail, &vinculaciones); err != nil {
		panic(map[string]interface{}{"funcion": "VinculacionTercero", "err": "Error get vinculacion", "status": "400", "log": err})
	} else {
		var vinculacionesResponse []models.Vinculacion
		for i := 0; i < len(vinculaciones); i++ {
			if vinculaciones[i].CargoId == int(idJefeOficina) || vinculaciones[i].CargoId == int(idAsistenteDependencia) || vinculaciones[i].CargoId == int(idNoRegistra) {
				vinculacionesResponse = append(vinculacionesResponse, vinculaciones[i])
			}
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": vinculacionesResponse}
	}
	c.ServeJSON()
}

// VinculacionTerceroByIdentificacion ...
// @Title VinculacionTerceroByIdentificacion
// @Description get VinculacionTerceroByIdentificacion
// @Param	identificacion	path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /vinculacion_tercero_identificacion/:identificacion [get]
func (c *FormulacionController) VinculacionTerceroByIdentificacion() {

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

	identificacionTercero := c.Ctx.Input.Param(":identificacion")
	var vinculaciones []models.Vinculacion
	var tercero []models.DatosIdentificacion

	idNoRegistra, idJefeOficina, idAsistenteDependencia, err := formulacionhelper.ObtenerIdParametros()
	if err != nil {
		panic(map[string]interface{}{"funcion": "VinculacionTerceroByIdentificacion", "err": "Error get parametros", "status": "400", "log": err})
	}

	s := "Numero:" + identificacionTercero + ",Activo:true"
	if err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/datos_identificacion?query="+url.QueryEscape(s), &tercero); err != nil || tercero[0].TerceroID.ID == 0 {
		panic(map[string]interface{}{"funcion": "VinculacionTerceroByIdentificacion", "err": "Error get tercero", "status": "400", "log": err})
	}

	TerceroIdStr := fmt.Sprintf("%d", tercero[0].TerceroID.ID)
	if err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/vinculacion?query=Activo:true,TerceroPrincipalId:"+TerceroIdStr, &vinculaciones); err != nil {
		panic(map[string]interface{}{"funcion": "VinculacionTerceroByIdentificacion", "err": "Error get vinculacion", "status": "400", "log": err})
	} else {
		var vinculacionesResponse []models.Vinculacion
		for i := 0; i < len(vinculaciones); i++ {
			if vinculaciones[i].CargoId == int(idJefeOficina) || vinculaciones[i].CargoId == int(idAsistenteDependencia) || vinculaciones[i].CargoId == int(idNoRegistra) {
				vinculacionesResponse = append(vinculacionesResponse, vinculaciones[i])
			}
		}
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": vinculacionesResponse}
	}
	c.ServeJSON()
}

// CambioCargoIdVinculacionTercero ...
// @Title CambioCargoIdVinculacionTercero
// @Description put Cambio de cargo_id en VinculacionTercero by idVinculacion
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} {}
// @Failure 403 :id is empty
// @router /cargo_vinculacion/:id [put]
func (c *FormulacionController) CambioCargoIdVinculacionTercero() {
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
	var body map[string]interface{}
	id := c.Ctx.Input.Param(":id")
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	res, err := formulacionhelper.CambioCargoIdVinculacionTercero(id, body)

	if err != nil {
		panic(map[string]interface{}{"funcion": "CambioCargoIdVinculacionTercero", "err": "Error actualizando cargo_id en vinculacion", "status": "400", "log": err})
	}

	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": res}
	c.ServeJSON()
}

// Planes ...
// @Title Planes
// @Description get Rubros
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /planes [get]
func (c *FormulacionController) Planes() {

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

	var respuesta map[string]interface{}
	var res map[string]interface{}
	var planes []map[string]interface{}
	var planesPED []map[string]interface{}
	var planesPI []map[string]interface{}
	var tipoPlanes []map[string]interface{}
	var plan map[string]interface{}
	var arregloPlanes []map[string]interface{}
	var auxArregloPlanes []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=formato:true", &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planes)

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=tipo_plan_id:6239117116511e20405d408b", &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planesPI)
		} else {
			panic(err)
		}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=tipo_plan_id:616513b91634adfaffed52bf", &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planesPED)
		} else {
			panic(err)
		}

		auxArregloPlanes = append(auxArregloPlanes, planes...)
		auxArregloPlanes = append(auxArregloPlanes, planesPI...)
		auxArregloPlanes = append(auxArregloPlanes, planesPED...)

		for i := 0; i < len(auxArregloPlanes); i++ {
			plan = auxArregloPlanes[i]
			tipoPlanId := plan["tipo_plan_id"].(string)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-plan?query=_id:"+tipoPlanId, &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &tipoPlanes)
				tipoPlan := tipoPlanes[0]
				nombreTipoPlan := tipoPlan["nombre"]
				planesTipo := make(map[string]interface{})
				planesTipo["_id"] = plan["_id"]
				planesTipo["nombre"] = plan["nombre"]
				planesTipo["descripcion"] = plan["descripcion"]
				planesTipo["tipo_plan_id"] = tipoPlanId
				planesTipo["formato"] = plan["formato"]
				planesTipo["vigencia"] = plan["vigencia"]
				planesTipo["dependencia_id"] = plan["dependencia_id"]
				planesTipo["aplicativo_id"] = plan["aplicativo_id"]
				planesTipo["activo"] = plan["activo"]
				planesTipo["nombre_tipo_plan"] = nombreTipoPlan

				arregloPlanes = append(arregloPlanes, planesTipo)

				if arregloPlanes != nil {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": arregloPlanes}
				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": ""}
				}

			} else {
				panic(err)
			}
		}

	} else {
		panic(err)
	}
	c.ServeJSON()
}

// VerificarIdentificaciones ...
// @Title VerificarIdentificaciones
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /verificar_identificaciones/:id [get]
func (c *FormulacionController) VerificarIdentificaciones() {

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
	var respuesta map[string]interface{}
	var respuestaPlan map[string]interface{}
	var respuestaDependencia []map[string]interface{}
	var dependencia map[string]interface{}
	var plan map[string]interface{}
	var identificaciones []map[string]interface{}
	var bandera bool

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+id, &respuestaPlan); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuestaPlan, &plan)

		if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+plan["dependencia_id"].(string), &respuestaDependencia); err == nil {
			dependencia = respuestaDependencia[0]

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+id, &respuesta); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &identificaciones)

				tipoDependencia := dependencia["TipoDependenciaId"].(map[string]interface{})
				id := dependencia["DependenciaId"].(map[string]interface{})["Id"]
				if (tipoDependencia["Id"] == 2.00 || id == 67.00) && id != 8.0 {
					bandera = formulacionhelper.VerificarDataIdentificaciones(identificaciones, "facultad")
				} else {
					bandera = formulacionhelper.VerificarDataIdentificaciones(identificaciones, "unidad")
				}

			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}

	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": bandera}
	c.ServeJSON()
}

// Get Planes En Formulacion ...
// @Title GetPlanesEnFormulacion
// @Description get Planes en formulacion
// @Success 200 {object} models.Formulacion
// @Failure 400 bad response
// @router /planes_formulacion [get]
func (c *FormulacionController) PlanesEnFormulacion() {
	defer helpers.ErrorController(c.Controller, "PlanesFormulacionController")

	if resumenPlanesActivos, err := formulacionhelper.ObtenerPlanesFormulacion(); err != nil {
		panic(map[string]interface{}{"funcion": "PlanesEnFormulacion", "err": err, "status": "400"})
	} else {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": resumenPlanesActivos}
	}

	c.ServeJSON()
}

// CalculosDocentes ...
// @Title CalculosDocentes
// @Description post Formulacion
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /calculos_docentes [post]
func (c *FormulacionController) CalculosDocentes() {

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

	//Obtener respuesta del body
	var body map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		panic(map[string]interface{}{"funcion": "CalculosDocentes", "err": "Error al decodificar el cuerpo de la solicitud", "status": "400", "log": err})
	}

	// Obtener Desagregado
	body["vigencia"] = body["vigencia"].(float64) - 1
	bodyResolucionesDocente := formulacionhelper.ConstruirCuerpoRD(body)
	respuestaPost, err := formulacionhelper.GetDesagregado(bodyResolucionesDocente)
	if err != nil {
		panic(map[string]interface{}{"funcion": "CalculosDocentes", "err": "Error al obtener desagregado", "status": "400", "log": err})
	}
	result := respuestaPost["Data"].([]interface{})

	//Peticion GET hacia Parametros Service
	vigenciaStr := strconv.FormatFloat(body["vigencia"].(float64), 'f', 0, 64)
	salarioMinimo, err := formulacionhelper.GetSalarioMinimo(vigenciaStr)
	if err != nil {
		panic(map[string]interface{}{"funcion": "CalculosDocentes", "err": "Error al obtener salario minimo", "status": "400", "log": err})
	}

	// Objeto para hacer los cálculos necesarios
	data := body
	data["resolucionDocente"] = result[0].(map[string]interface{})
	data["salarioMinimo"] = salarioMinimo["Valor"]
	delete(body, "vigencia")
	delete(body, "categoria")
	delete(body, "tipo")

	// Realizar los calculos
	dataFinal, err := formulacionhelper.GetCalculos(data)
	if err != nil {
		panic(map[string]interface{}{"funcion": "CalculosDocentes", "err": "Error al intentar realizar los calculos", "status": "400", "log": err})
	}

	c.Data["json"] = map[string]interface{}{
		"Success": true,
		"Status":  "200",
		"Message": "Successful",
		"Data":    dataFinal,
	}
	c.ServeJSON()
}

// EstructuraPlanes ...
// @Title EstructuraPlanes
// @Description put Formulacion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /estructura_planes/:id [put]
func (c *FormulacionController) EstructuraPlanes() {
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
	//Obtener plantilla por id
	plantilla, err := formulacionhelper.GetPlantilla(id)
	if err != nil {
		panic(map[string]interface{}{"funcion": "EstructuraPlanes", "err": "Error al obtener plantilla", "status": "400", "log": err})
	}

	//Obtener los planes en estado "En formulacion" asociados a la plantilla
	planes, err := formulacionhelper.GetPlanesPorNombre(plantilla["nombre"].(string))
	if err != nil {
		panic(map[string]interface{}{"funcion": "EstructuraPlanes", "err": "Error al obtener planes asociados a plantilla", "status": "400", "log": err})
	}

	//Obtener el formato de la plantilla
	formatoPLantilla, err := formulacionhelper.GetFormato(id)
	if err != nil {
		panic(map[string]interface{}{"funcion": "EstructuraPlanes", "err": "Error al obtener formato de plantilla", "status": "400", "log": err})
	}

	//Obtener lista plana del formato
	listaPlantilla, err := formulacionhelper.ConvArbolAListaPlana(formatoPLantilla[0], id, true)
	if err != nil {
		panic(map[string]interface{}{"funcion": "EstructuraPlanes", "err": "Error al obtener el valor de referencia", "status": "400", "log": err})
	}

	//Obtener los formatos de los planes y comparar con el formato de la plantilla
	for _, plan := range planes {
		planId := plan["_id"].(string)
		formatoPlan, err := formulacionhelper.GetFormato(planId)
		if err != nil {
			panic(map[string]interface{}{"funcion": "EstructuraPlanes", "err": "Error al obtener formato de plan", "status": "400", "log": err})
		}
		listaPlan, err := formulacionhelper.ConvArbolAListaPlana(formatoPlan[0], planId, false)
		if err == nil {
			formulacionhelper.ActualizarEstructuraPlan(listaPlantilla, listaPlan, planId)
		}
	}

	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": "La estructura de los planes fue actualizada correctamente"}
	c.ServeJSON()
}

// DefinirFechasFuncionamiento ...
// @Title DefinirFechasFuncionamiento
// @Description Peticion POST para definir fechas en planes de acción de funcionamiento e inversión para los procesos de formulación y seguimiento
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 400 bad request
// @router /habilitar_fechas_funcionamiento [post]
func (c *FormulacionController) DefinirFechasFuncionamiento() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/" + "FormulacionController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var res interface{}

	// Decodificar JSON desde el cuerpo de la solicitud
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err != nil {
		// Manejar el error, imprimirlo o devolver una respuesta de error al cliente
		c.Abort("400") // Bad Request
		return
	}
	res = formulacionhelper.DefinirFechasFormulacionSeguimiento(body)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": res}
	c.ServeJSON()
}

// GetPlanesUnidadesComun ...
// @Title GetPlanesUnidadesComun
// @Description post Get planes en comun con unidades by id periodo-seguimiento
// @Param	id		path 	string	true		"The key for staticblock"
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 200 {object} models.Formulacion
// @Failure 403 :id is empty
// @router /get_planes_unidades_comun/:id [post]
func (c *FormulacionController) GetPlanesUnidadesComun() {
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
	var body_unidades map[string]interface{}
	var periodo_seguimiento map[string]interface{}
	var res map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &body_unidades)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetPlanesUnidadesComun", "err": "Error al decodificar el cuerpo de la solicitud", "status": "400", "log": err})
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+id, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &periodo_seguimiento)
		unidadesValidadas := formulacionhelper.ValidarUnidadesPlanes(periodo_seguimiento, body_unidades)

		// Verificar el resultado
		if len(unidadesValidadas) > 0 {
			fmt.Println("Unidades de intersección:", unidadesValidadas)
			planesInteres := periodo_seguimiento["planes_interes"]
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": planesInteres}
		} else {
			// fmt.Println("No hay unidades en la intersección")
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": "Successful", "Data": "Not found"}
		}
		// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": planesInteres}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
