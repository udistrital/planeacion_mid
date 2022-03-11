package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	// "reflect"
	// "strconv"
	// "encoding/base64"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/xuri/excelize/v2"
	reporteshelper "github.com/udistrital/planeacion_mid/helpers/reportesHelper"
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

func CreateExcel(f *excelize.File, dir string){
	if err := f.SaveAs(dir); err != nil{
		fmt.Println(err)
	}
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
	var arreglo []map[string]interface{}
	// excel
	var consolidadoExcel *excelize.File
	consolidadoExcel = excelize.NewFile()
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
		fmt.Println(planesFilter)
		for i := 0; i < len(planesFilter); i++ {
			plan = planesFilter[i]
			planId := plan["_id"].(string)
			dependencia := plan["dependencia_id"].(string)
			if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia?query=Id:"+dependencia, &respuestaOikos); err == nil {
				nombreDep = respuestaOikos[0]
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=activo:true,plan_id:"+planId+",tipo_identificacion_id:"+"617b6630f6fc97b776279afa", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &identificacionres)
				identificacion = identificacionres[0]
				if identificacion["dato"] != "{}" {
					dato_str := identificacion["dato"].(string)
					json.Unmarshal([]byte(dato_str), &dato)
					for key := range dato {
						element := dato[key].(map[string]interface{})
						if element["activo"] == true {
							delete(element, "actividades")
							delete(element, "activo")
							delete(element, "index")
							element["unidad"] = nombreDep["Nombre"]
							data_identi = append(data_identi, element)
						}

					}
					arreglo = append(arreglo, data_identi...)

					for h := 0 ; h < len(arreglo); h++ {

						datosArreglo := arreglo[h]
						// fmt.Println(datosArreglo["Nombre"])
						nombreHoja := fmt.Sprint(datosArreglo["unidad"])
						sheetName := nombreHoja
						index := consolidadoExcel.NewSheet(sheetName)
						consolidadoExcel.SetCellValue(sheetName, "A1", "Dependencia Responsable")
							consolidadoExcel.SetCellValue(sheetName, "A2", "codigo del rubro")
							consolidadoExcel.SetCellValue(sheetName, "A3", datosArreglo["codigo"])

							consolidadoExcel.SetCellValue(sheetName, "B1", datosArreglo["unidad"])

							consolidadoExcel.SetCellValue(sheetName, "B2", "Nombre del rubro")
							consolidadoExcel.SetCellValue(sheetName, "B3", datosArreglo["Nombre"])
							consolidadoExcel.SetCellValue(sheetName, "C2", "valor")
							consolidadoExcel.SetCellValue(sheetName, "C3", datosArreglo["valor"])
							consolidadoExcel.SetCellValue(sheetName, "D2", "Descripcion del bien y/o servicio")
							consolidadoExcel.SetCellValue(sheetName, "D3", datosArreglo["descripcion"])
							consolidadoExcel.SetActiveSheet(index)
					}

				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": ""}
				}

			} else {
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
				c.Abort("400")
			}
		}
		CreateExcel(consolidadoExcel, "Consolidado Presupuestal.xls")
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": arreglo}


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
	var res map[string]interface{}
	var resSub map[string]interface{}
	var subgrupos []map[string]interface{}
	var hijos []map[string]interface{}
	var plan_id string
	var actividadName string
	// var datad map[string]interface{}
	// var resSubDetalle map[string]interface{}
	// var subgruposDetalle []map[string]interface{}
	// var data map[string]interface{}
	// var data_identi []map[string]interface{}
	var arregloPlanAnual []map[string]interface{}

	var generalData = make(map[string]interface{})
	var indicadoresData = make(map[string]interface{})

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if body["unidad_id"].(string) == ""{
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
			// fmt.Println(respuesta)
			for i := 0; i < len(planesFilter); i++ {
				planesFilterData := planesFilter[i]
				plan_id = planesFilterData["_id"].(string)
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &res); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
					index := body["index"].(string)
					datad := reporteshelper.GetDataSubgrupos(subgrupos, index)
					for i := range datad{
						test := datad[i]
						// el problema está aquí porque no sirve con el .(map[string]interface{})
						// data_identi = append(data_identi, test)
						fmt.Println(test)
					}
					// if datad != nil {
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": "not found", "Data": datad}
					// 	// dato_str := identificacion["dato"].(string)
					// 	// json.Unmarshal([]byte(data), &dato)
					// 	// for key := range datad {
					// 	// 	element := datad[key].(map[string]interface{})
					// 	// 	data_identi = append(data_identi, element)

					// 	// 	// if element["activo"] == true {
					// 	// 	// 	// delete(element, "actividades")
					// 	// 	// 	// delete(element, "activo")
					// 	// 	// 	// delete(element, "index")
					// 	// 	// 	// aqui dentro por medio del element puedo acceder tanto al nombre y a todos los demas atributos del objeto para hacer el excel sería acá
					// 	// 	// 	element["unidad"] = nombreDep["Nombre"]
					// 	// 	// 	data_identi = append(data_identi, element)
					// 	// 	// 	fmt.Println(dato);
					// 	// 	// }

					// 	// }
					// 	// // c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": data_identi}
					// 	arreglo = append(arreglo, datad...)

					// } else {
					// 	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": "not found", "Data": ""}
					// }

					// arreglo = append(arreglo, datad...)

					// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": data}
				}

			}
			// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": data}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	}else if body["unidad_id"].(string) != ""{
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",dependencia_id:"+body["unidad_id"].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
			// fmt.Println(respuesta)
			planesFilterData := planesFilter[0]
			plan_id = planesFilterData["_id"].(string)
			// fmt.Println(plan_id)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id, &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
		
				for i := 0; i < len(subgrupos); i++ {
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
						actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
						fmt.Println("este de abajo es actividaes")
						fmt.Println(actividades)
						for j := 0; j < len(actividades); j++ {
							actividad := actividades[j]
							actividadName = actividad["dato"].(string)
							// fmt.Println(actividadName)
							index := fmt.Sprint(actividad["index"])

							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &resSub); err == nil {
								helpers.LimpiezaRespuestaRefactor(resSub, &hijos)
								data := reporteshelper.GetDataSubgrupos(hijos, index)
								generalData := make(map[string]interface{})
								generalData["hijos"] = data
								generalData["nombreActividad"] = actividadName
								generalData["numeroActividad"] = index

								arregloPlanAnual = append(arregloPlanAnual, generalData)
								// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
							} else {
								c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
								c.Abort("400")
							}


						}
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": arregloPlanAnual}
		
						break
					}
				}
			} else {
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
				c.Abort("400")
			}



			// if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=activo:true,padre:"+plan_id, &res); err == nil {
			// if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &res); err == nil {
			// 	helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
			// 	index := body["index"].(string)
			// 	data := seguimientohelper.GetDataSubgrupos(subgrupos, index)
			// 	// primeros espacios de la tabla
			// 	lineamiento := data["lineamiento"]
			// 	fmt.Println(lineamiento)
			// 	metaEstrategica := data["meta_estrategica"]
			// 	fmt.Println(metaEstrategica)
			// 	estrategia := data["estrategia"]
			// 	fmt.Println(estrategia)

			// 	// siguientes espacios de la tabla
			// 	tarea := data["tarea"]
			// 	fmt.Println(tarea)

			// 	// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": data}

			// 	// if data = nil {
			// 	// 	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": "not found", "Data": ""}
			// 	// }else {
			// 	// 	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": data}
			// 	// }

			// 	// fmt.Println(subgrupos)
			// 	for i := 0; i < len(subgrupos); i++ {
			// 		subgrupo := subgrupos[i]
			// 		nombreIndicador := subgrupo["nombre"]
			// 		descripcionIndicador := subgrupo["descripcion"]
			// 		// subgrupo_id := subgrupo["_id"].(string)
			// 		if subgrupo["hijos"] != "[]"{
			// 			dato_str := subgrupo["hijos"]
			// 			fmt.Println(dato_str)
			// 			// json.Unmarshal([]byte(dato_str), &dato)
			// 			// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dato_str}
			// 		}else{
			// 			// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": "not found", "Data": ""}
			// 		}
			// 		// if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=activo:true,subgrupo_id:"+subgrupo_id, &resSubDetalle); err == nil {
			// 		// 	helpers.LimpiezaRespuestaRefactor(resSubDetalle, &subgruposDetalle)
			// 		// 	for j := 0; j < len(subgruposDetalle); j++ {
			// 		// 		subgrupoDetalle := subgruposDetalle[j]
			// 		// 		if subgrupoDetalle["dato"] != "{}"{
			// 		// 			dato_str := subgrupoDetalle["dato"].(string)
			// 		// 			json.Unmarshal([]byte(dato_str), &dato)
			// 		// 			// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": subgrupoDetalle}
			// 		// 		}else{
			// 		// 			// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "404", "Message": "not found", "Data": ""}
			// 		// 		}
			// 		// 	}
			// 		// 	// fmt.Println(subgruposDetalle)
			// 		// 	// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": subgruposDetalle}

			// 		// } else {
			// 		// 	c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			// 		// 	c.Abort("400")
			// 		// }

			// 		indicadoresData["nombreIndicador"] = nombreIndicador
			// 		indicadoresData["descripcionIndicador"] = descripcionIndicador
			// 		// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": data}
			// 	}

			// } else {
			// 	c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			// 	c.Abort("400")
			// }

			nombrePlan := planesFilterData["nombre"]
			descripcionPlan := planesFilterData["descripcion"]
			generalData["indicadores"] = indicadoresData
			generalData["nombrePlan"] = nombrePlan
			generalData["descripcionPlan"] = descripcionPlan

			// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": generalData}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
		// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": "{}"}
	}

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
