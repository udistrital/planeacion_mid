package controllers

import (
	"encoding/json"
	"fmt"
	// "reflect"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

// ReportesController operations for Reportes
type ReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReportesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Desagregado", c.Desagregado)
	c.Mapping("PlanAccionAnual", c.PlanAccionAnual)
}

// Desagregado ...
// @Title Desagregado
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Reportes
// @Failure 403 :plan_id is empty
// @router /desagregado [post]
func (c *ReportesController) Desagregado() {
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var planesFilter []map[string]interface{}
	var plan map[string]interface{}
	var respuestaOikos []map[string]interface{}
	var nombreDep map[string]interface{}
	var identificacionres []map[string]interface{}
	var res map[string]interface{}
	var identificacion map[string]interface{}
	var dato map[string]interface{}
	var data_identi []map[string]interface{}
	var identificacion_data map[string]interface{}

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
		fmt.Println(respuesta)

		for i := 0; i < len(planesFilter); i++ {
			plan = planesFilter[i]
			planId := plan["_id"].(string)
			fmt.Println(plan)	
			dependencia := plan["dependencia_id"].(string)
			fmt.Println(dependencia)

			dependenciaId, err := strconv.ParseFloat(dependencia, 8)
			if err != nil {
				fmt.Println(err)
			}
			// priodoId_rest, err := strconv.ParseFloat(test1, 8)
			// 	if err != nil {
			// 		fmt.Println(err)
			// 	}
			fmt.Println(dependenciaId+1)
			 
			if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia?query=Id:8", &respuestaOikos); err == nil {
				nombreDep = respuestaOikos[0]
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+planId+",tipo_identificacion_id:"+"617b6630f6fc97b776279afa", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &identificacionres)
				identificacion = identificacionres[0]
				fmt.Println(identificacion)		
				if identificacion["dato"] != nil {
					dato_str := identificacion["dato"].(string)
					json.Unmarshal([]byte(dato_str), &dato)
					for key := range dato {
						element := dato[key].(map[string]interface{})
						if element["activo"] == true {
							data_identi = append(data_identi, element)
							fmt.Println(data_identi);
						}
						// identificacion_data = data_identi[i]
						// identificacionData := make(map[string]interface{})
						// identificacionData["codigo"] = identificacion_data["codigo"]
						// identificacionData["concepto"] = identificacion_data["Nombre"]
						// identificacionData["valor"] = identificacion_data["valor"]
						// identificacionData["unidad"] = nombreDep["Nombre"]
						// identificacionData["descripcion"] = identificacion_data["descripcion"]
						// fmt.Println(identificacionData)
						// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": identificacionData}
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": data_identi}

					}
		
				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": ""}
				}
		
			} else {
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
				c.Abort("400")
			}
		}
		
		identificacion_data = data_identi[0]
		identificacionData := make(map[string]interface{})
		identificacionData["codigo"] = identificacion_data["codigo"]
		identificacionData["concepto"] = identificacion_data["Nombre"]
		identificacionData["valor"] = identificacion_data["valor"]
		identificacionData["unidad"] = nombreDep["Nombre"]
		identificacionData["descripcion"] = identificacion_data["descripcion"]
		fmt.Println(identificacionData)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": identificacionData}


	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}
	c.ServeJSON()
}

// PlanAccionAnual ...
// @Title PlanAccionAnual
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Reportes
// @Failure 403 :plan_id is empty
// @router /plan_anual [post]
func (c *ReportesController) PlanAccionAnual() {
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var planesFilter []map[string]interface{}


	// var parametro_periodo []map[string]interface{}
	// var dato map[string]interface{}
	// var seguimiento map[string]interface{}
	// var test1 string
	// var periodIdString string
	// var periodId float64
	// var avanceAcumulado string
	// var testavancePeriodo string
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if body["unidad_id"].(string) == ""{
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
			fmt.Println(respuesta)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": planesFilter}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	}else{
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": "{}"}
	}

	// if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_id:"+body["periodo_id"].(string), &res); err == nil {
	// 	helpers.LimpiezaRespuestaRefactor(res, &avancedata)
	// 	fmt.Println(res)
	// 	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+body["periodo_id"].(string), &res); err == nil {
	// 		helpers.LimpiezaRespuestaRefactor(res, &parametro_periodo)
	// 		paramIdlen := parametro_periodo[0]
	// 		fmt.Println()
	// 		paramId := paramIdlen["ParametroId"].(map[string]interface{})
	// 		if paramId["CodigoAbreviacion"] != "T1"{
	// 			test1 = body["periodo_id"].(string)
	// 			priodoId_rest, err := strconv.ParseFloat(test1, 8)
	// 			if err != nil {
	// 				fmt.Println(err)
	// 			}
	// 			fmt.Println(test1)
	// 			periodId = priodoId_rest - 1
	// 			periodIdString = fmt.Sprint(periodId)
	// 			if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+body["plan_id"].(string)+",periodo_id:"+periodIdString, &res); err == nil {
	// 				helpers.LimpiezaRespuestaRefactor(res, &avancedata)
	// 				seguimiento = avancedata[0]
	// 				datoStr := seguimiento["dato"].(string)
	// 				json.Unmarshal([]byte(datoStr), &dato)
	// 				indicador1 := dato[body["index"].(string)].(map[string]interface{})
	// 				fmt.Println(len(indicador1))
	// 				avanceIndicador1 := indicador1[body["Nombre_del_indicador"].(string)].(map[string]interface{})
	// 				avanceAcumulado = avanceIndicador1["avanceAcumulado"].(string)
	// 				testavancePeriodo = avanceIndicador1["avancePeriodo"].(string)
	// 			}else{
	// 				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	// 				c.Abort("400")
	// 			}
	// 		}else{
	// 			fmt.Println("hi")
	// 		}
	// 	}else {
	// 		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	// 		c.Abort("400")
	// 	}
	// 	avancePeriodo := body["avancePeriodo"].(string)
	// 	aPe, err := strconv.ParseFloat(avancePeriodo, 8)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(aPe, err, reflect.TypeOf(avanceAcumulado))
	// 	aAc, err := strconv.ParseFloat(avanceAcumulado, 8)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	totalAcumulado := fmt.Sprint(aPe + aAc)
	// 	generalData := make(map[string]interface{})
	// 	generalData["avancePeriodo"] = avancePeriodo
	// 	generalData["periodIdString"] = periodIdString
	// 	generalData["avanceAcumulado"] = totalAcumulado
	// 	generalData["avancePeriodoPrev"] = testavancePeriodo
	// 	generalData["avanceAcumuladoPrev"] = avanceAcumulado

	// 	fmt.Println(avanceAcumulado, reflect.TypeOf(avanceAcumulado))
	// 	fmt.Println(avancePeriodo, reflect.TypeOf(avancePeriodo))
	// 	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": generalData}
	// 	fmt.Println(dato)
	// } else {
	// 	c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
	// 	c.Abort("400")
	// }

	c.ServeJSON()
}

// Post ...
// @Title Create
// @Description create Reportes
// @Param	body		body 	models.Reportes	true		"body for Reportes content"
// @Success 201 {object} models.Reportes
// @Failure 403 body is empty
// @router / [post]
func (c *ReportesController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Reportes by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Reportes
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ReportesController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Reportes
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Reportes
// @Failure 403
// @router / [get]
func (c *ReportesController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Reportes
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Reportes	true		"body for Reportes content"
// @Success 200 {object} models.Reportes
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ReportesController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Reportes
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ReportesController) Delete() {

}
