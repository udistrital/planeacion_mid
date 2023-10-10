package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	evaluacionhelper "github.com/udistrital/planeacion_mid/helpers/evaluacionHelper"
	reporteshelper "github.com/udistrital/planeacion_mid/helpers/reportesHelper"
	"github.com/udistrital/utils_oas/request"
	"github.com/xuri/excelize/v2"
)

// ReportesController operations for Reportes
type ReportesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReportesController) URLMapping() {
	c.Mapping("Desagregado", c.Desagregado)
	c.Mapping("ValidarReporte", c.ValidarReporte)
	c.Mapping("PlanAccionAnual", c.PlanAccionAnual)
	c.Mapping("PlanAccionAnualGeneral", c.PlanAccionAnualGeneral)
	c.Mapping("Necesidades", c.Necesidades)
	c.Mapping("PlanAccionEvaluacion", c.PlanAccionEvaluacion)
}

func CreateExcel(f *excelize.File, dir string) {
	if err := f.Save(); err != nil {
		fmt.Println(err)
	}

}

// ValidarReporte ...
// @Title ValidarReporte
// @Description post ValidarReporte
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Reportes
// @router /validar_reporte [post]
func (c *ReportesController) ValidarReporte() {
	var res1 map[string]interface{}
	var resFilter []map[string]interface{}
	var body map[string]interface{}
	res := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if body["categoria"].(string) == "Evaluación" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",dependencia_id:"+body["unidad_id"].(string), &res1); err == nil {
			helpers.LimpiezaRespuestaRefactor(res1, &resFilter)

			if len(resFilter) == 0 {
				res["mensaje"] = "No existen planes para la unidad seleccionada"
				res["reporte"] = false
			} else {
				noPlan := true
				noVigencia := true
				noEstado := true
				for i := 0; i < len(resFilter); i++ {
					if resFilter[i]["nombre"] == body["nombre"].(string) {
						noPlan = false
						if resFilter[i]["vigencia"] == body["vigencia"].(string) {
							noVigencia = false
							if resFilter[i]["estado_plan_id"] == "6153355601c7a2365b2fb2a1" {
								noEstado = false
								res["mensaje"] = ""
								res["reporte"] = true
								break
							}
						}
					}
				}

				if noPlan {
					res["mensaje"] = "La unidad no tiene registros con el plan seleccionado"
					res["reporte"] = false
				} else if noVigencia {
					res["mensaje"] = "La unidad no cuenta con registros para la vigencia y el plan selecionados"
					res["reporte"] = false
				} else if noEstado {
					res["mensaje"] = "La unidad no cuenta con plan avalado"
					res["reporte"] = false
				}
			}
		} else {
			res["mensaje"] = "Ocurrio un error"
		}
	} else if body["categoria"].(string) == "Necesidades" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string), &res1); err == nil {
			helpers.LimpiezaRespuestaRefactor(res1, &resFilter)
			if len(resFilter) == 0 {
				res["mensaje"] = "No existen planes para la vigencia seleccionada"
				res["reporte"] = false
			} else {
				noPlan := true
				noEstado := true
				for i := 0; i < len(resFilter); i++ {
					if resFilter[i]["nombre"] == body["nombre"].(string) {
						noPlan = false
						if resFilter[i]["estado_plan_id"] == body["estado_plan_id"].(string) {
							noEstado = false
							res["mensaje"] = ""
							res["reporte"] = true
							break
						}
					}
				}

				if noPlan {
					res["mensaje"] = "No existen registros con el plan seleccionado"
					res["reporte"] = false
				} else if noEstado {
					res["mensaje"] = "No existen registros con el estado y plan seleccionado"
					res["reporte"] = false
				}
			}
		} else {
			res["mensaje"] = "Ocurrio un error"
		}
	} else if body["categoria"].(string) == "Plan de acción unidad" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",dependencia_id:"+body["unidad_id"].(string), &res1); err == nil {
			helpers.LimpiezaRespuestaRefactor(res1, &resFilter)
			if len(resFilter) == 0 {
				res["mensaje"] = "No existen planes para la unidad seleccionada"
				res["reporte"] = false
			} else {
				noPlan := true
				noVigencia := true
				noEstado := true
				for i := 0; i < len(resFilter); i++ {
					if resFilter[i]["nombre"] == body["nombre"].(string) {
						noPlan = false
						if resFilter[i]["vigencia"] == body["vigencia"].(string) {
							noVigencia = false
							if resFilter[i]["estado_plan_id"] == body["estado_plan_id"].(string) {
								noEstado = false
								res["mensaje"] = ""
								res["reporte"] = true
								break
							}
						}
					}
				}

				if noPlan {
					res["mensaje"] = "La unidad no tiene registros con el plan seleccionado"
					res["reporte"] = false
				} else if noVigencia {
					res["mensaje"] = "La unidad no cuenta con registros para la vigencia y el plan selecionados"
					res["reporte"] = false
				} else if noEstado {
					res["mensaje"] = "La unidad no cuenta con plan en el estado solicitado"
					res["reporte"] = false
				}
			}
		} else {
			res["mensaje"] = "Ocurrio un error"
		}
	} else if body["categoria"].(string) == "Plan de acción general" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string), &res1); err == nil {
			helpers.LimpiezaRespuestaRefactor(res1, &resFilter)
			if len(resFilter) == 0 {
				res["mensaje"] = "No existen planes para la vigencia seleccionada"
				res["reporte"] = false
			} else {
				noPlan := true
				noEstado := true
				for i := 0; i < len(resFilter); i++ {
					if resFilter[i]["nombre"] == body["nombre"].(string) {
						noPlan = false
						if resFilter[i]["estado_plan_id"] == body["estado_plan_id"].(string) {
							noEstado = false
							res["mensaje"] = ""
							res["reporte"] = true
							break
						}
					}
				}

				if noPlan {
					res["mensaje"] = "No existen registros con el plan seleccionado"
					res["reporte"] = false
				} else if noEstado {
					res["mensaje"] = "No existen registros con el estado y plan seleccionado"
					res["reporte"] = false
				}
			}
		} else {
			res["mensaje"] = "Ocurrio un error"
		}
	} else {
		res["mensaje"] = "Categoria incorrecta"
		res["reporte"] = false
	}

	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": res}
	c.ServeJSON()
}

// Desagregado ...
// @Title Desagregado
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Reportes
// @Failure 403 :plan_id is empty
// @router /desagregado [post]
func (c *ReportesController) Desagregado() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ReportesController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

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
	var nombreUnidadVer string

	// excel
	var consolidadoExcel = excelize.NewFile()
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)

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

				} else {
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": ""}
				}

			} else {
				panic(err)
			}
		}
		contadorDesagregado := 3
		stylehead, _ := consolidadoExcel.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
			Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
			Border: []excelize.Border{
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})
		styletitles, _ := consolidadoExcel.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Font:      &excelize.Font{Bold: true},
			Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
			Border: []excelize.Border{
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})
		stylecontent, _ := consolidadoExcel.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Border: []excelize.Border{
				{Type: "right", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})

		for h := 0; h < len(data_identi); h++ {

			datosArreglo := data_identi[h]
			nombreUnidadVerIn := datosArreglo["unidad"].(string)
			if h == 0 {
				nombreUnidadVer = nombreUnidadVerIn
			}

			if nombreUnidadVerIn == nombreUnidadVer {
				nombreUnidadVer = datosArreglo["unidad"].(string)
				nombreHoja := nombreUnidadVer
				sheetName := nombreHoja
				index, _ := consolidadoExcel.NewSheet(sheetName)
				consolidadoExcel.MergeCell(sheetName, "B1", "D1")

				consolidadoExcel.SetRowHeight(sheetName, 1, 20)
				consolidadoExcel.SetRowHeight(sheetName, 2, 20)
				consolidadoExcel.SetRowHeight(sheetName, contadorDesagregado, 50)

				consolidadoExcel.SetColWidth(sheetName, "A", "A", 30)
				consolidadoExcel.SetColWidth(sheetName, "B", "B", 50)
				consolidadoExcel.SetColWidth(sheetName, "C", "C", 30)
				consolidadoExcel.SetColWidth(sheetName, "D", "D", 60)

				consolidadoExcel.SetCellValue(sheetName, "A1", "Dependencia Responsable")
				consolidadoExcel.SetCellValue(sheetName, "B1", nombreUnidadVer)

				consolidadoExcel.SetCellValue(sheetName, "A2", "Código del rubro")
				consolidadoExcel.SetCellValue(sheetName, "B2", "Nombre del rubro")
				consolidadoExcel.SetCellValue(sheetName, "C2", "valor")
				consolidadoExcel.SetCellValue(sheetName, "D2", "Descripción del bien y/o servicio")

				consolidadoExcel.SetCellValue(sheetName, "A"+fmt.Sprint(contadorDesagregado), datosArreglo["codigo"])
				consolidadoExcel.SetCellValue(sheetName, "B"+fmt.Sprint(contadorDesagregado), datosArreglo["Nombre"])
				consolidadoExcel.SetCellValue(sheetName, "C"+fmt.Sprint(contadorDesagregado), datosArreglo["valor"])
				consolidadoExcel.SetCellValue(sheetName, "D"+fmt.Sprint(contadorDesagregado), datosArreglo["descripcion"])
				consolidadoExcel.SetCellStyle(sheetName, "A1", "A1", stylehead)
				consolidadoExcel.SetCellStyle(sheetName, "B1", "D1", stylehead)

				consolidadoExcel.SetCellStyle(sheetName, "A2", "D2", styletitles)
				consolidadoExcel.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorDesagregado), "D"+fmt.Sprint(contadorDesagregado), stylecontent)
				consolidadoExcel.SetActiveSheet(index)

				contadorDesagregado = contadorDesagregado + 1
			} else {
				contadorDesagregado = 3
				nombreUnidadVer = datosArreglo["unidad"].(string)
				nombreHoja := nombreUnidadVer
				sheetName := nombreHoja
				index, _ := consolidadoExcel.NewSheet(sheetName)
				consolidadoExcel.MergeCell(sheetName, "B1", "D1")

				consolidadoExcel.SetRowHeight(sheetName, 1, 20)
				consolidadoExcel.SetRowHeight(sheetName, 2, 20)
				consolidadoExcel.SetRowHeight(sheetName, contadorDesagregado, 50)

				consolidadoExcel.SetColWidth(sheetName, "A", "A", 30)
				consolidadoExcel.SetColWidth(sheetName, "B", "B", 50)
				consolidadoExcel.SetColWidth(sheetName, "C", "C", 30)
				consolidadoExcel.SetColWidth(sheetName, "D", "D", 60)

				consolidadoExcel.SetCellValue(sheetName, "A1", "Dependencia Responsable")
				consolidadoExcel.SetCellValue(sheetName, "B1", nombreUnidadVer)

				consolidadoExcel.SetCellValue(sheetName, "A2", "Código del rubro")
				consolidadoExcel.SetCellValue(sheetName, "B2", "Nombre del rubro")
				consolidadoExcel.SetCellValue(sheetName, "C2", "valor")
				consolidadoExcel.SetCellValue(sheetName, "D2", "Descripción del bien y/o servicio")

				consolidadoExcel.SetCellValue(sheetName, "A"+fmt.Sprint(contadorDesagregado), datosArreglo["codigo"])
				consolidadoExcel.SetCellValue(sheetName, "B"+fmt.Sprint(contadorDesagregado), datosArreglo["Nombre"])
				consolidadoExcel.SetCellValue(sheetName, "C"+fmt.Sprint(contadorDesagregado), datosArreglo["valor"])
				consolidadoExcel.SetCellValue(sheetName, "D"+fmt.Sprint(contadorDesagregado), datosArreglo["descripcion"])
				consolidadoExcel.SetCellStyle(sheetName, "A1", "A1", stylehead)
				consolidadoExcel.SetCellStyle(sheetName, "B1", "D1", stylehead)

				consolidadoExcel.SetCellStyle(sheetName, "A2", "D2", styletitles)
				consolidadoExcel.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorDesagregado), "D"+fmt.Sprint(contadorDesagregado), stylecontent)
				consolidadoExcel.SetActiveSheet(index)

			}

		}

		if len(consolidadoExcel.GetSheetList()) > 1 {
			consolidadoExcel.DeleteSheet("Sheet1")
		}

		dataSend := make(map[string]interface{})

		buf, _ := consolidadoExcel.WriteToBuffer()
		strings.NewReader(buf.String())

		encoded := base64.StdEncoding.EncodeToString([]byte(buf.Bytes()))

		dataSend["generalData"] = data_identi
		dataSend["excelB64"] = encoded

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}

	} else {
		panic(err)
	}
	c.ServeJSON()
}

// PlanAccionAnual ...
// @Title PlanAccionAnual
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Param	nombre		path 	string	true		"The key for staticblock"
// @Success 201 {object} models.Reportes
// @Failure 403 :plan_id is empty
// @router /plan_anual/:nombre [post]
func (c *ReportesController) PlanAccionAnual() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ReportesController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var planesFilter []map[string]interface{}
	var res map[string]interface{}
	var respuestaUnidad []map[string]interface{}
	var subgrupos []map[string]interface{}
	var plan_id string
	var actividadName string
	var arregloPlanAnual []map[string]interface{}
	var nombreUnidad string
	var resPeriodo map[string]interface{}
	var periodo []map[string]interface{}
	var unidadNombre string
	nombre := c.Ctx.Input.Param(":nombre")
	consolidadoExcelPlanAnual := excelize.NewFile()
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if body["unidad_id"].(string) != "" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",dependencia_id:"+body["unidad_id"].(string)+",nombre:"+nombre, &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)

			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Id:`+body["vigencia"].(string), &resPeriodo); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
			}

			for planes := 0; planes < len(planesFilter); planes++ {
				planesFilterData := planesFilter[planes]
				plan_id = planesFilterData["_id"].(string)

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id+"&fields=nombre,_id,hijos,activo", &res); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
					for i := 0; i < len(subgrupos); i++ {
						if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
							actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
							var arregloLineamieto []map[string]interface{}
							var arregloLineamietoPI []map[string]interface{}
							sort.SliceStable(actividades, func(i int, j int) bool {
								if _, ok := actividades[i]["index"].(float64); ok {
									actividades[i]["index"] = fmt.Sprintf("%v", int(actividades[i]["index"].(float64)))
								}
								if _, ok := actividades[j]["index"].(float64); ok {
									actividades[j]["index"] = fmt.Sprintf("%v", int(actividades[j]["index"].(float64)))
								}
								aux, _ := strconv.Atoi((actividades[i]["index"]).(string))
								aux1, _ := strconv.Atoi((actividades[j]["index"]).(string))
								return aux < aux1
							})
							reporteshelper.LimpiarDetalles()
							for j := 0; j < len(actividades); j++ {
								arregloLineamieto = nil
								arregloLineamietoPI = nil
								actividad := actividades[j]
								actividadName = actividad["dato"].(string)
								index := actividad["index"].(string)
								datosArmonizacion := make(map[string]interface{})
								titulosArmonizacion := make(map[string]interface{})

								reporteshelper.Limpia()
								tree := reporteshelper.BuildTreeFa(subgrupos, index)
								treeDatos := tree[0]
								treeDatas := tree[1]
								treeArmo := tree[2]
								armonizacionTercer := treeArmo[0]
								var armonizacionTercerNivel interface{}
								var armonizacionTercerNivelPI interface{}

								if armonizacionTercer["armo"] != nil {
									armonizacionTercerNivel = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPED"]
									armonizacionTercerNivelPI = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPI"]
								}

								for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
									treeDato := treeDatos[datoGeneral]
									treeData := treeDatas[0]
									if treeDato["sub"] == "" {
										nombre := strings.ToLower(treeDato["nombre"].(string))
										if strings.Contains(nombre, "ponderación") || strings.Contains(nombre, "ponderacion") && strings.Contains(nombre, "actividad") {
											datosArmonizacion["Ponderación de la actividad"] = treeData[treeDato["id"].((string))]
										} else if strings.Contains(nombre, "período") || strings.Contains(nombre, "periodo") && strings.Contains(nombre, "ejecucion") || strings.Contains(nombre, "ejecución") {
											datosArmonizacion["Periodo de ejecución"] = treeData[treeDato["id"].(string)]
										} else if strings.Contains(nombre, "actividad") && strings.Contains(nombre, "general") {
											datosArmonizacion["Actividad general"] = treeData[treeDato["id"].(string)]
										} else if strings.Contains(nombre, "tarea") || strings.Contains(nombre, "actividades específicas") {
											datosArmonizacion["Tareas"] = treeData[treeDato["id"].(string)]
										} else {
											datosArmonizacion[treeDato["nombre"].(string)] = treeData[treeDato["id"].(string)]
										}
									}
								}
								var treeIndicador map[string]interface{}
								auxTree := tree[0]
								for i := 0; i < len(auxTree); i++ {
									subgrupo := auxTree[i]
									if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") {
										treeIndicador = auxTree[i]
									}
								}

								subIndicador := treeIndicador["sub"].([]map[string]interface{})
								for ind := 0; ind < len(subIndicador); ind++ {
									subIndicadorRes := subIndicador[ind]
									treeData := treeDatas[0]
									dataIndicador := make(map[string]interface{})
									auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
									for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
										if treeData[auxSubIndicador[subInd]["id"].(string)] == nil {
											treeData[auxSubIndicador[subInd]["id"].(string)] = ""
										}
										dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[auxSubIndicador[subInd]["id"].(string)]
									}
									titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
								}

								datosArmonizacion["indicadores"] = titulosArmonizacion
								if armonizacionTercerNivel != nil {
									//arregloLineamieto = reporteshelper.ArbolArmonizacion(armonizacionTercerNivel.(string))
									arregloLineamieto = reporteshelper.ArbolArmonizacionV2(armonizacionTercerNivel.(string))
								} else {
									arregloLineamieto = []map[string]interface{}{}
								}
								if armonizacionTercerNivelPI != nil {
									//arregloLineamietoPI = reporteshelper.ArbolArmonizacionPI(armonizacionTercerNivelPI)
									arregloLineamietoPI = reporteshelper.ArbolArmonizacionPIV2(armonizacionTercerNivelPI.(string))
								} else {
									arregloLineamietoPI = []map[string]interface{}{}
								}

								generalData := make(map[string]interface{})
								if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+body["unidad_id"].(string), &respuestaUnidad); err == nil {
									aux := respuestaUnidad[0]
									dependenciaNombre := aux["DependenciaId"].(map[string]interface{})
									nombreUnidad = dependenciaNombre["Nombre"].(string)
								} else {
									panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
								}

								generalData["nombreUnidad"] = nombreUnidad
								generalData["nombreActividad"] = actividadName
								generalData["numeroActividad"] = index
								generalData["datosArmonizacion"] = arregloLineamieto
								generalData["datosArmonizacionPI"] = arregloLineamietoPI
								generalData["datosComplementarios"] = datosArmonizacion

								arregloPlanAnual = append(arregloPlanAnual, generalData)

							}
							break
						}
					}
				} else {
					panic(err)
				}

				unidadNombre = arregloPlanAnual[0]["nombreUnidad"].(string)
				sheetName := "Actividades del plan"
				indexPlan, _ := consolidadoExcelPlanAnual.NewSheet(sheetName)

				if planes == 0 {
					consolidadoExcelPlanAnual.DeleteSheet("Sheet1")

					disable := false
					if err := consolidadoExcelPlanAnual.SetSheetView(sheetName, -1, &excelize.ViewOptions{
						ShowGridLines: &disable,
					}); err != nil {
						fmt.Println(err)
					}
				}
				stylehead, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
					Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "top", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				styletitles, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Font:      &excelize.Font{Bold: true},
					Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "top", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				stylecontent, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "top", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				stylecontentS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
					Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "top", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				stylecontentC, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "top", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				stylecontentCL, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 4},
					},
				})
				stylecontentCLD, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				stylecontentCS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "top", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				stylecontentCLS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 4},
					},
				})
				stylecontentCLDS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
					Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				styleLineamiento, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{
						Horizontal:   "center",
						Vertical:     "center",
						WrapText:     true,
						TextRotation: 90,
					},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "top", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})
				styleLineamientoSombra, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Alignment: &excelize.Alignment{
						Horizontal:   "center",
						Vertical:     "center",
						WrapText:     true,
						TextRotation: 90,
					},
					Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
					Border: []excelize.Border{
						{Type: "right", Color: "000000", Style: 1},
						{Type: "left", Color: "000000", Style: 1},
						{Type: "top", Color: "000000", Style: 1},
						{Type: "bottom", Color: "000000", Style: 1},
					},
				})

				consolidadoExcelPlanAnual.MergeCell(sheetName, "B1", "D1")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E1", "G1")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "H1", "H2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "I1", "I2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "J1", "J2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "K1", "K2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "L1", "L2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "P1", "P2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "M1", "O1")
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 18)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "P", 35)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "C", 11)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "E", 16)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "H", 6)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "J", 12)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "K", "K", 30)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "L", 35)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "M", "N", 52)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "O", "O", 10)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "P", "P", 30)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B1", "P1", stylehead)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B2", "P2", styletitles)

				// encabezado excel
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "B1", "Armonización PED")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "B2", "Lineamiento")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "C2", "Meta")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "D2", "Estrategias")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "E1", "Armonización Plan Indicativo")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "E2", "Ejes transformadores")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "F2", "Lineamientos de acción")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "G2", "Estrategias")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "H2", "N°.")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "I2", "Ponderación de la actividad")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "J2", "Periodo de ejecución")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "K2", "Actividad")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Actividades específicas")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M1", "Indicador")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M2", "Nombre")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "N2", "Fórmula")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "O2", "Meta")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "P2", "Producto esperado")

				rowPos := 3

				for excelPlan := 0; excelPlan < len(arregloPlanAnual); excelPlan++ {

					datosExcelPlan := arregloPlanAnual[excelPlan]
					armoPED := datosExcelPlan["datosArmonizacion"].([]map[string]interface{})
					armoPI := datosExcelPlan["datosArmonizacionPI"].([]map[string]interface{})
					datosComplementarios := datosExcelPlan["datosComplementarios"].(map[string]interface{})
					indicadores := datosComplementarios["indicadores"].(map[string]interface{})

					MaxRowsXActivity := reporteshelper.MinComMul_Armonization(armoPED, armoPI, len(indicadores))

					y_lin := rowPos
					h_lin := MaxRowsXActivity / len(armoPED)
					for _, lin := range armoPED {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1), styleLineamiento, styleLineamientoSombra)
						y_met := y_lin
						h_met := h_lin / len(lin["meta"].([]map[string]interface{}))
						for _, met := range lin["meta"].([]map[string]interface{}) {
							consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1))
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(y_met), met["nombreMeta"])
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1), stylecontentC, stylecontentCS)
							y_est := y_met
							h_est := h_met / len(met["estrategias"].([]map[string]interface{}))
							for _, est := range met["estrategias"].([]map[string]interface{}) {
								consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1))
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(y_est), est["descripcionEstrategia"])
								reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
								y_est += h_est
							}
							y_met += h_met
						}
						y_lin += h_lin
					}

					y_eje := rowPos
					h_eje := MaxRowsXActivity / len(armoPI)
					for _, eje := range armoPI {
						consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(y_eje), eje["nombreFactor"])
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1), stylecontent, stylecontentS)
						y_lin := y_eje
						h_lin := h_eje / len(eje["lineamientos"].([]map[string]interface{}))
						for _, lin := range eje["lineamientos"].([]map[string]interface{}) {
							consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1))
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1), stylecontent, stylecontentS)
							y_est := y_lin
							h_est := h_lin / len(lin["estrategias"].([]map[string]interface{}))
							for _, est := range lin["estrategias"].([]map[string]interface{}) {
								consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1))
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(y_est), est["descripcionEstrategia"])
								reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
								y_est += h_est
							}
							y_lin += h_lin
						}
						y_eje += h_eje
					}

					consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(rowPos), "H"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(rowPos), "I"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(rowPos), "K"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(rowPos), datosExcelPlan["numeroActividad"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(rowPos), datosComplementarios["Ponderación de la actividad"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(rowPos), datosComplementarios["Periodo de ejecución"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(rowPos), datosComplementarios["Actividad general"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(rowPos), datosComplementarios["Tareas"])
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(rowPos), "J"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(rowPos), "L"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontent, stylecontentS)

					y_ind := rowPos
					h_ind := MaxRowsXActivity / len(indicadores)
					idx := int(0)
					for _, indicador := range indicadores {
						auxIndicador := indicador
						var nombreIndicador interface{}
						var formula interface{}
						var meta interface{}
						for key, element := range auxIndicador.(map[string]interface{}) {
							if strings.Contains(strings.ToLower(key), "nombre") {
								nombreIndicador = element
							}
							if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
								formula = element
							}
							if strings.Contains(strings.ToLower(key), "meta") {
								meta = element
							}
						}
						consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(y_ind), "M"+fmt.Sprint(y_ind+h_ind-1))
						consolidadoExcelPlanAnual.MergeCell(sheetName, "N"+fmt.Sprint(y_ind), "N"+fmt.Sprint(y_ind+h_ind-1))
						consolidadoExcelPlanAnual.MergeCell(sheetName, "O"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1))
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(y_ind), nombreIndicador)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(y_ind), formula)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(y_ind), meta)
						idx++
						if idx < len(indicadores) {
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1), stylecontentCL, stylecontentCLS)
						} else {
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(y_ind), "O"+fmt.Sprint(y_ind+h_ind-1), stylecontentCLD, stylecontentCLDS)
						}
						y_ind += h_ind
					}

					consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(rowPos), "P"+fmt.Sprint(rowPos+MaxRowsXActivity-1))
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(rowPos), datosComplementarios["Producto esperado"])
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(rowPos), "P"+fmt.Sprint(rowPos+MaxRowsXActivity-1), stylecontentC, stylecontentCS)

					rowPos += MaxRowsXActivity

					consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
				}
				consolidadoExcelPlanAnual = reporteshelper.TablaIdentificaciones(consolidadoExcelPlanAnual, plan_id)

			}

			if len(planesFilter) <= 0 {
				c.Abort("404")
			}

			styletitle, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"},
				Font:      &excelize.Font{Bold: true, Size: 18, Color: "000000"},
				Border: []excelize.Border{
					{Type: "right", Color: "ffffff", Style: 1},
					{Type: "left", Color: "ffffff", Style: 1},
					{Type: "top", Color: "ffffff", Style: 1},
					{Type: "bottom", Color: "ffffff", Style: 1},
				},
			})

			consolidadoExcelPlanAnual.InsertRows("Actividades del plan", 1, 7)
			consolidadoExcelPlanAnual.MergeCell("Actividades del plan", "C2", "P6")
			consolidadoExcelPlanAnual.SetCellStyle("Actividades del plan", "C2", "P6", styletitle)
			consolidadoExcelPlanAnual.SetCellStyle("Identificaciones", "C2", "G6", styletitle)

			if periodo[0] != nil {
				consolidadoExcelPlanAnual.SetCellValue("Actividades del plan", "C2", "Plan de Acción "+periodo[0]["Nombre"].(string)+"\n"+unidadNombre)
				consolidadoExcelPlanAnual.SetCellValue("Identificaciones", "C2", "Proyección de necesidades "+periodo[0]["Nombre"].(string)+"\n"+unidadNombre)
			} else {
				consolidadoExcelPlanAnual.SetCellValue("Actividades del plan", "C2", "Plan de Acción")
				consolidadoExcelPlanAnual.SetCellValue("Identificaciones", "C2", "Proyección de necesidades")
			}

			if err := consolidadoExcelPlanAnual.AddPicture("Actividades del plan", "B1", "static/img/UDEscudo2.png",
				&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "oneCell", OffsetX: 10}); err != nil {
				fmt.Println(err)
			}
			if err := consolidadoExcelPlanAnual.AddPicture("Identificaciones", "B1", "static/img/UDEscudo2.png",
				&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "absolute", OffsetX: 10}); err != nil {
				fmt.Println(err)
			}

			consolidadoExcelPlanAnual.SetColWidth("Actividades del plan", "A", "A", 2)
			buf, _ := consolidadoExcelPlanAnual.WriteToBuffer()
			strings.NewReader(buf.String())
			encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

			dataSend := make(map[string]interface{})
			dataSend["generalData"] = arregloPlanAnual
			dataSend["excelB64"] = encoded

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}
		} else {
			panic(err)
		}
	}

	c.ServeJSON()
}

// PlanAccionAnualGeneral ...
// @Title PlanAccionAnualGeneral
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Param	nombre		path 	string	true		"The key for staticblock"
// @Success 201 {object} models.Reportes
// @Failure 403 :plan_id is empty
// @router /plan_anual_general/:nombre [post]
func (c *ReportesController) PlanAccionAnualGeneral() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ReportesController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var planesFilter []map[string]interface{}
	var res map[string]interface{}
	var respuestaUnidad []map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaTipoPlan map[string]interface{}
	var estado map[string]interface{}
	var tipoPlan map[string]interface{}
	var subgrupos []map[string]interface{}
	var plan_id string
	var actividadName string
	var arregloPlanAnual []map[string]interface{}
	var arregloInfoReportes []map[string]interface{}
	var nombreUnidad string
	var idUnidad string
	var resPeriodo map[string]interface{}
	var periodo []map[string]interface{}
	contadorGeneral := 4

	consolidadoExcelPlanAnual := excelize.NewFile()
	nombre := c.Ctx.Input.Param(":nombre")
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",nombre:"+nombre+"&fields=_id,dependencia_id,estado_plan_id,tipo_plan_id", &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
		for _, planes := range planesFilter {
			if idUnidad != planes["dependencia_id"].(string) {
				if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+planes["dependencia_id"].(string), &respuestaUnidad); err == nil {
					planes["nombreUnidad"] = respuestaUnidad[0]["DependenciaId"].(map[string]interface{})["Nombre"].(string)
				} else {
					panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
				}
			}
		}

		sort.SliceStable(planesFilter, func(i, j int) bool {
			a := (planesFilter)[i]["nombreUnidad"].(string)
			b := (planesFilter)[j]["nombreUnidad"].(string)
			return a < b
		})

		for planes := 0; planes < len(planesFilter); planes++ {
			reporteshelper.Limp()
			planesFilterData := planesFilter[planes]
			plan_id = planesFilterData["_id"].(string)
			infoReporte := make(map[string]interface{})
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id+"&fields=nombre,_id,hijos,activo", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
				for i := 0; i < len(subgrupos); i++ {
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
						actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
						var arregloLineamieto []map[string]interface{}
						var arregloLineamietoPI []map[string]interface{}
						sort.SliceStable(actividades, func(i int, j int) bool {
							if _, ok := actividades[i]["index"].(float64); ok {
								actividades[i]["index"] = fmt.Sprintf("%v", int(actividades[i]["index"].(float64)))
							}
							if _, ok := actividades[j]["index"].(float64); ok {
								actividades[j]["index"] = fmt.Sprintf("%v", int(actividades[j]["index"].(float64)))
							}
							aux, _ := strconv.Atoi((actividades[i]["index"]).(string))
							aux1, _ := strconv.Atoi((actividades[j]["index"]).(string))
							return aux < aux1
						})
						reporteshelper.LimpiarDetalles()
						for j := 0; j < len(actividades); j++ {
							arregloLineamieto = nil
							arregloLineamietoPI = nil
							actividad := actividades[j]
							actividadName = actividad["dato"].(string)
							index := actividad["index"].(string)
							datosArmonizacion := make(map[string]interface{})
							titulosArmonizacion := make(map[string]interface{})

							//reporteshelper.Limpia()
							tree := reporteshelper.BuildTreeFa(subgrupos, index)
							treeDatos := tree[0]
							treeDatas := tree[1]
							treeArmo := tree[2]
							armonizacionTercer := treeArmo[0]
							var armonizacionTercerNivel interface{}
							var armonizacionTercerNivelPI interface{}
							if armonizacionTercer["armo"] != nil {
								armonizacionTercerNivel = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPED"]
								armonizacionTercerNivelPI = armonizacionTercer["armo"].(map[string]interface{})["armonizacionPI"]
							}

							for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
								treeDato := treeDatos[datoGeneral]
								treeData := treeDatas[0]
								if treeDato["sub"] == "" {
									nombreMinuscula := strings.ToLower(treeDato["nombre"].(string))
									if strings.Contains(nombreMinuscula, "ponderación") || strings.Contains(nombreMinuscula, "ponderacion") && strings.Contains(nombreMinuscula, "actividad") {
										datosArmonizacion["Ponderación de la actividad"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "período") || strings.Contains(nombreMinuscula, "periodo") && strings.Contains(nombreMinuscula, "ejecucion") || strings.Contains(nombreMinuscula, "ejecución") {
										datosArmonizacion["Periodo de ejecución"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "actividad") && strings.Contains(nombreMinuscula, "general") {
										datosArmonizacion["Actividad general"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "tarea") || strings.Contains(nombreMinuscula, "actividades específicas") {
										datosArmonizacion["Tareas"] = treeData[treeDato["id"].(string)]
									} else if strings.Contains(nombreMinuscula, "producto") {
										datosArmonizacion["Producto esperado"] = treeData[treeDato["id"].(string)]
									} else {
										datosArmonizacion[treeDato["nombre"].(string)] = treeData[treeDato["id"].(string)]
									}
								}
							}
							var treeIndicador map[string]interface{}
							auxTree := tree[0]
							for i := 0; i < len(auxTree); i++ {
								subgrupo := auxTree[i]
								if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") {
									treeIndicador = auxTree[i]
								}
							}

							subIndicador := treeIndicador["sub"].([]map[string]interface{})
							for ind := 0; ind < len(subIndicador); ind++ {
								subIndicadorRes := subIndicador[ind]
								treeData := treeDatas[0]
								dataIndicador := make(map[string]interface{})
								auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
								for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
									dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[auxSubIndicador[subInd]["id"].(string)]
								}
								titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
							}

							datosArmonizacion["indicadores"] = titulosArmonizacion
							if armonizacionTercerNivel != nil {
								arregloLineamieto = reporteshelper.ArbolArmonizacion(armonizacionTercerNivel.(string))
							} else {
								arregloLineamieto = []map[string]interface{}{}
							}
							if armonizacionTercerNivelPI != nil {
								arregloLineamietoPI = reporteshelper.ArbolArmonizacionPI(armonizacionTercerNivelPI)
							} else {
								arregloLineamietoPI = []map[string]interface{}{}
							}

							generalData := make(map[string]interface{})

							// if idUnidad != planesFilter[planes]["dependencia_id"].(string) {
							// 	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+planesFilter[planes]["dependencia_id"].(string), &respuestaUnidad); err == nil {
							// 		aux := respuestaUnidad[0]
							// 		dependenciaNombre := aux["DependenciaId"].(map[string]interface{})
							// 		nombreUnidad = dependenciaNombre["Nombre"].(string)
							// 		idUnidad = planesFilter[planes]["dependencia_id"].(string)

							// 	} else {
							// 		panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
							// 	}
							// }

							nombreUnidad = planesFilterData["nombreUnidad"].(string)

							generalData["nombreUnidad"] = nombreUnidad
							generalData["nombreActividad"] = actividadName
							generalData["numeroActividad"] = index
							generalData["datosArmonizacion"] = arregloLineamieto
							generalData["datosArmonizacionPI"] = arregloLineamietoPI
							generalData["datosComplementarios"] = datosArmonizacion
							arregloPlanAnual = append(arregloPlanAnual, generalData)

						}
						break
					}

				}
			} else {
				panic(err)
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-plan/"+planesFilter[planes]["estado_plan_id"].(string), &respuestaEstado); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuestaEstado, &estado)
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-plan/"+planesFilter[planes]["tipo_plan_id"].(string), &respuestaTipoPlan); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuestaTipoPlan, &tipoPlan)
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}

			infoReporte["tipo_plan"] = tipoPlan["nombre"]
			infoReporte["vigencia"] = body["vigencia"]
			infoReporte["estado_plan"] = estado["nombre"]
			infoReporte["nombre_unidad"] = nombreUnidad

			arregloInfoReportes = append(arregloInfoReportes, infoReporte)

			contadorLineamiento := contadorGeneral + 5
			contadorFactor := contadorGeneral + 5
			contadorDataGeneral := contadorGeneral + 5

			unidadNombre := arregloPlanAnual[0]["nombreUnidad"]
			sheetName := "REPORTE GENERAL"
			indexPlan, _ := consolidadoExcelPlanAnual.NewSheet(sheetName)

			if planes == 0 {
				consolidadoExcelPlanAnual.DeleteSheet("Sheet1")
				consolidadoExcelPlanAnual.InsertCols("REPORTE GENERAL", "A", 1)
				disable := false
				if err := consolidadoExcelPlanAnual.SetSheetView(sheetName, -1, &excelize.ViewOptions{
					ShowGridLines: &disable,
				}); err != nil {
					fmt.Println(err)
				}

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Id:`+body["vigencia"].(string), &resPeriodo); err == nil {
					helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
				}
			}

			stylehead, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styletitles, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: true},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			stylecontent, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			stylecontentS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			stylecontentC, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			stylecontentCS, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleLineamiento, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal:   "center",
					Vertical:     "center",
					WrapText:     true,
					TextRotation: 90,
				},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleLineamientoSombra, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal:   "center",
					Vertical:     "center",
					WrapText:     true,
					TextRotation: 90,
				},
				Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})

			consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "P"+fmt.Sprint(contadorGeneral+1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "D"+fmt.Sprint(contadorGeneral+2))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorGeneral+2), "G"+fmt.Sprint(contadorGeneral+2))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorGeneral+2), "H"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorGeneral+2), "I"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorGeneral+2), "J"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorGeneral+2), "K"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(contadorGeneral+2), "L"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(contadorGeneral+2), "P"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "M"+fmt.Sprint(contadorGeneral+2), "O"+fmt.Sprint(contadorGeneral+2))
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+1, 20)
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+2, 20)
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 20)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 19)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "P", 35)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "C", "C", 13)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "E", 16)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "H", 6)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "J", 12)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "K", "K", 30)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "L", 35)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "M", "N", 52)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "O", "O", 10)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "P", "P", 30)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+1), "P"+fmt.Sprint(contadorGeneral+1), stylehead)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "P"+fmt.Sprint(contadorGeneral+2), stylehead)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "P"+fmt.Sprint(contadorGeneral+3), styletitles)
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 30)

			var tituloExcel string
			if periodo[0] != nil {
				tituloExcel = "Plan de acción " + periodo[0]["Nombre"].(string) + " - " + unidadNombre.(string)
			} else {
				tituloExcel = "Plan de acción - " + unidadNombre.(string)
			}

			// encabezado excel
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+1), tituloExcel)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+2), "Armonización PED")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "Lineamiento")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorGeneral+3), "Meta")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorGeneral+3), "Estrategias")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorGeneral+2), "Armonización Plan Indicativo")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorGeneral+3), "Ejes transformadores")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorGeneral+3), "Lineamientos de acción")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorGeneral+3), "Estrategias")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorGeneral+3), "N°.")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorGeneral+3), "Ponderación de la actividad")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorGeneral+3), "Periodo de ejecución")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorGeneral+3), "Actividad")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+3), "Actividades específicas")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorGeneral+2), "Indicador")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorGeneral+3), "Nombre")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorGeneral+3), "Fórmula")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(contadorGeneral+3), "Meta")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(contadorGeneral+3), "Producto esperado")
			consolidadoExcelPlanAnual.InsertRows(sheetName, 1, 1)

			for excelPlan := 0; excelPlan < len(arregloPlanAnual); excelPlan++ {
				datosExcelPlan := arregloPlanAnual[excelPlan]
				armoPED := datosExcelPlan["datosArmonizacion"].([]map[string]interface{})
				armoPI := datosExcelPlan["datosArmonizacionPI"].([]map[string]interface{})
				datosComplementarios := datosExcelPlan["datosComplementarios"].(map[string]interface{})

				var contadorLineamientoGeneralIn int
				var contadorLineamientoGeneralOut int

				var contadorMetaGeneralIn int
				var contadorMetaGeneralOut int

				var contadorEstrategiaPEDIn int
				var contadorEstrategiaPEDOut int

				var contadorFactorGeneralIn int
				var contadorFactorGeneralOut int

				var contadorLineamientoPIIn int
				var contadorLineamientoPIOut int

				var contadorEstrategiaPIIn int
				var contadorEstrategiaPIOut int

				_ = contadorEstrategiaPEDIn
				_ = contadorEstrategiaPIIn

				for i := 0; i < len(armoPED); i++ {
					datosArmo := armoPED[i]
					auxLineamiento := datosArmo["nombreLineamiento"]
					contadorLineamientoGeneralIn = contadorLineamiento

					// cuerpo del excel
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorLineamiento), auxLineamiento)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(contadorLineamiento), "B"+fmt.Sprint(contadorLineamiento), styleLineamiento, styleLineamientoSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorLineamiento), "P"+fmt.Sprint(contadorLineamiento), stylecontentC, stylecontentCS)
					consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorLineamiento, 70)

					metas := datosArmo["meta"]

					contadorMetas := contadorLineamiento

					for j := 0; j < len(metas.([]map[string]interface{})); j++ {
						auxMeta := metas.([]map[string]interface{})[j]
						contadorMetaGeneralIn = contadorLineamiento

						consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorMetas), auxMeta["nombreMeta"])
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorMetas), "C"+fmt.Sprint(contadorMetas), stylecontentC, stylecontentCS)

						if j == len(metas.([]map[string]interface{}))-1 {
							contadorMetaGeneralOut = contadorMetas
						} else {
							contadorMetas = contadorMetas + 1
						}

						estrategias := auxMeta["estrategias"].([]map[string]interface{})
						contadorEstrategias := contadorMetas
						for k := 0; k < len(estrategias); k++ {
							auxEstrategia := estrategias[k]
							contadorEstrategiaPEDIn = contadorMetas

							consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorEstrategias), auxEstrategia["descripcionEstrategia"])
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategias), "D"+fmt.Sprint(contadorEstrategias), stylecontentC, stylecontentCS)

							if k == len(estrategias)-1 {
								contadorEstrategiaPEDOut = contadorMetas
							} else {
								contadorEstrategias = contadorEstrategias + 1
							}
						}

						contadorEstrategias = contadorMetas

						consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorMetaGeneralIn), "C"+fmt.Sprint(contadorMetaGeneralOut))

					}
					// contadorMetas = contadorLineamiento
					contadorLineamientoGeneralOut = contadorMetaGeneralOut

					consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralIn), "B"+fmt.Sprint(contadorLineamientoGeneralOut))

					contadorLineamiento = contadorLineamientoGeneralOut + 1
				}

				for i := 0; i < len(armoPI); i++ {
					datosArmo := armoPI[i]
					auxFactor := datosArmo["nombreFactor"]
					contadorFactorGeneralIn = contadorFactor

					// cuerpo del excel
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorFactor), auxFactor)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactor), "P"+fmt.Sprint(contadorFactor), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(contadorFactor), "L"+fmt.Sprint(contadorFactor), stylecontent, stylecontentS)
					consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorFactor, 70)

					lineamientos := datosArmo["lineamientos"]

					contadorLineamientos := contadorFactor

					for j := 0; j < len(lineamientos.([]map[string]interface{})); j++ {
						auxLineamiento := lineamientos.([]map[string]interface{})[j]
						contadorLineamientoPIIn = contadorFactor

						consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorLineamientos), auxLineamiento["nombreLineamiento"])
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientos), "F"+fmt.Sprint(contadorLineamientos), stylecontentC, stylecontentCS)

						if j == len(lineamientos.([]map[string]interface{}))-1 {
							contadorLineamientoPIOut = contadorLineamientos
						} else {
							contadorLineamientos = contadorLineamientos + 1
						}

						estrategiasPI := auxLineamiento["estrategias"].([]map[string]interface{})
						contadorEstrategias := contadorLineamientos
						for k := 0; k < len(estrategiasPI); k++ {
							auxEstrategia := estrategiasPI[k]
							contadorEstrategiaPEDIn = contadorLineamientos

							consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorEstrategias), auxEstrategia["descripcionEstrategia"])
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategias), "G"+fmt.Sprint(contadorEstrategias), stylecontentC, stylecontentCS)

							if k == len(estrategiasPI)-1 {
								contadorEstrategiaPIOut = contadorLineamientos
							} else {
								contadorEstrategias = contadorEstrategias + 1
							}
						}

						contadorEstrategias = contadorLineamientos

						consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorLineamientoPIIn), "F"+fmt.Sprint(contadorLineamientoPIOut))

					}
					// contadorLineamientos = contadorFactor
					contadorFactorGeneralOut = contadorLineamientoPIOut

					consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorFactorGeneralIn), "E"+fmt.Sprint(contadorFactorGeneralOut))

					contadorFactor = contadorFactorGeneralOut + 1
				}

				consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorDataGeneral), datosExcelPlan["numeroActividad"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Ponderación de la actividad"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Periodo de ejecución"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Actividad general"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Tareas"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Producto esperado"])

				if contadorLineamientoGeneralOut > contadorFactorGeneralOut {
					contadorFactorGeneralOut = contadorLineamientoGeneralOut
					// contadorFactor = contadorFactorGeneralOut + 1
				} else if contadorLineamientoGeneralOut < contadorFactorGeneralOut {
					contadorLineamientoGeneralOut = contadorFactorGeneralOut
					contadorLineamiento = contadorLineamientoGeneralOut + 1
				}

				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorLineamiento), "P"+fmt.Sprint(contadorLineamiento), stylecontentC, stylecontentCS)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(contadorLineamiento), "L"+fmt.Sprint(contadorLineamiento), stylecontent, stylecontentS)

				consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralIn), "B"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorMetaGeneralIn), "C"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorFactorGeneralIn), "E"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorDataGeneral), "H"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorDataGeneral), "I"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorDataGeneral), "K"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorLineamientoGeneralOut))

				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralIn), "B"+fmt.Sprint(contadorLineamientoGeneralOut), styleLineamiento, styleLineamientoSombra)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorMetaGeneralIn), "C"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactorGeneralIn), "E"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentS)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)

				indicadores := datosComplementarios["indicadores"].(map[string]interface{})
				contadorIndicadores := contadorDataGeneral
				for id, indicador := range indicadores {
					_ = id
					auxIndicador := indicador
					var nombreIndicador interface{}
					var formula interface{}
					var meta interface{}

					for key, element := range auxIndicador.(map[string]interface{}) {
						if strings.Contains(strings.ToLower(key), "nombre") {
							nombreIndicador = element
						}
						if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
							formula = element
						}
						if strings.Contains(strings.ToLower(key), "meta") {
							meta = element
						}

					}

					consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorIndicadores), nombreIndicador)
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorIndicadores), formula)
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "O"+fmt.Sprint(contadorIndicadores), meta)

					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(contadorIndicadores), "O"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
					contadorIndicadores = contadorIndicadores + 1
				}

				contadorIndicadores--
				if contadorLineamientoGeneralOut < contadorIndicadores {
					consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralOut), "B"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorMetaGeneralOut), "C"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorFactorGeneralOut), "E"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorDataGeneral), "H"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorDataGeneral), "I"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorDataGeneral), "K"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorIndicadores))

					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralOut), "B"+fmt.Sprint(contadorIndicadores), styleLineamiento, styleLineamientoSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorMetaGeneralOut), "C"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactorGeneralOut), "E"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
				} else {
					contadorIndicadores = contadorLineamientoGeneralOut
					consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralOut), "B"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorMetaGeneralOut), "C"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorFactorGeneralOut), "E"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorDataGeneral), "H"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorDataGeneral), "I"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorDataGeneral), "K"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorLineamientoGeneralOut))

					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralOut), "B"+fmt.Sprint(contadorLineamientoGeneralOut), styleLineamiento, styleLineamientoSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorMetaGeneralOut), "C"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactorGeneralOut), "E"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
				}
				contadorDataGeneral = contadorIndicadores + 1
				contadorLineamiento = contadorIndicadores + 1
				contadorFactor = contadorIndicadores + 1
				consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
			}
			contadorGeneral = contadorDataGeneral - 2
			arregloPlanAnual = nil
			consolidadoExcelPlanAnual.RemoveRow(sheetName, 1)
		}

		if len(planesFilter) <= 0 {
			c.Abort("404")
		}

		styletitle, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"},
			Font:      &excelize.Font{Bold: true, Size: 18, Color: "000000"},
			Border: []excelize.Border{
				{Type: "right", Color: "ffffff", Style: 1},
				{Type: "left", Color: "ffffff", Style: 1},
				{Type: "top", Color: "ffffff", Style: 1},
				{Type: "bottom", Color: "ffffff", Style: 1},
			},
		})

		consolidadoExcelPlanAnual.InsertRows("REPORTE GENERAL", 1, 3)
		consolidadoExcelPlanAnual.MergeCell("REPORTE GENERAL", "C2", "P6")
		consolidadoExcelPlanAnual.SetCellStyle("REPORTE GENERAL", "C2", "P6", styletitle)
		if periodo[0] != nil {
			consolidadoExcelPlanAnual.SetCellValue("REPORTE GENERAL", "C2", "Plan de Acción Anual "+periodo[0]["Nombre"].(string)+"\nUniversidad Distrital Franciso José de Caldas")
		} else {
			consolidadoExcelPlanAnual.SetCellValue("REPORTE GENERAL", "C2", "Plan de Acción Anual \nUniversidad Distrital Franciso José de Caldas")
		}

		if err := consolidadoExcelPlanAnual.AddPicture("REPORTE GENERAL", "B1", "static/img/UDEscudo2.png",
			&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "oneCell", OffsetX: 10}); err != nil {
			fmt.Println(err)
		}

		consolidadoExcelPlanAnual.SetColWidth("REPORTE GENERAL", "A", "A", 2)
		buf, _ := consolidadoExcelPlanAnual.WriteToBuffer()
		strings.NewReader(buf.String())
		encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

		dataSend := make(map[string]interface{})
		dataSend["generalData"] = arregloInfoReportes
		dataSend["excelB64"] = encoded

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// Necesidades ...
// @Title Necesidades
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Param	nombre		path 	string	true		"The key for staticblock"
// @Success 201 {object} models.Reportes
// @Failure 403 :plan_id is empty
// @router /necesidades/:nombre [post]
func (c *ReportesController) Necesidades() {

	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ReportesController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var respuestaIdentificaciones map[string]interface{}
	var identificaciones []map[string]interface{}
	var planes []map[string]interface{}
	var recursos []map[string]interface{}
	var recursosGeneral []map[string]interface{}
	var rubros []map[string]interface{}
	var rubrosGeneral []map[string]interface{}
	var unidades_total []string
	var unidades_rubros_total []string
	var respuestaEstado map[string]interface{}
	var estado map[string]interface{}
	var respuestaTipo map[string]interface{}
	var tipo map[string]interface{}
	var arregloInfoReportes []map[string]interface{}
	docentesPregrado := make(map[string]interface{})
	docentesPosgrado := make(map[string]interface{})
	var docentesGeneral map[string]interface{}
	var arrDataDocentes []map[string]interface{}
	nombre := c.Ctx.Input.Param(":nombre")

	primaServicios := 0
	primaNavidad := 0
	primaVacaciones := 0
	bonificacion := 0
	interesesCesantias := 0
	cesantiasPublicas := 0
	cesantiasPrivadas := 0
	salud := 0
	pensionesPublicas := 0
	pensionesPrivadas := 0
	arl := 0
	caja := 0
	icbf := 0
	docentesPregrado["tco"] = 0
	docentesPregrado["mto"] = 0
	docentesPregrado["hch"] = 0
	docentesPregrado["hcp"] = 0
	docentesPregrado["valor"] = 0

	docentesPosgrado["hch"] = 0
	docentesPosgrado["hcp"] = 0
	docentesPosgrado["valor"] = 0

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	necesidadesExcel := excelize.NewFile()
	stylecontent, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentS, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentM, _ := necesidadesExcel.NewStyle(&excelize.Style{
		NumFmt:    183,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentMS, _ := necesidadesExcel.NewStyle(&excelize.Style{
		NumFmt:    183,
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styletitles, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylehead, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCL, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 4},
		},
	})
	stylecontentCML, _ := necesidadesExcel.NewStyle(&excelize.Style{
		NumFmt:    183,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 4},
		},
	})
	stylecontentCMD, _ := necesidadesExcel.NewStyle(&excelize.Style{
		NumFmt:    183,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"C2C2C2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCM, _ := necesidadesExcel.NewStyle(&excelize.Style{
		NumFmt:    183,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"C2C2C2"}},
		Border: []excelize.Border{
			{Type: "top", Color: "ffffff", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentC, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"C2C2C2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "ffffff", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCD, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"C2C2C2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylecontentCLS, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 4},
		},
	})
	stylecontentCMLS, _ := necesidadesExcel.NewStyle(&excelize.Style{
		NumFmt:    183,
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 4},
		},
	})

	necesidadesExcel.NewSheet("Necesidades")
	necesidadesExcel.DeleteSheet("Sheet1")
	disable := false
	if err := necesidadesExcel.SetSheetView("Necesidades", -1, &excelize.ViewOptions{
		ShowGridLines: &disable,
	}); err != nil {
		fmt.Println(err)
	}

	necesidadesExcel.MergeCell("Necesidades", "C1", "E1")

	necesidadesExcel.SetColWidth("Necesidades", "A", "A", 4)
	necesidadesExcel.SetColWidth("Necesidades", "B", "B", 26)
	necesidadesExcel.SetColWidth("Necesidades", "C", "C", 15)
	necesidadesExcel.SetColWidth("Necesidades", "C", "E", 15)
	necesidadesExcel.SetColWidth("Necesidades", "F", "F", 20)
	necesidadesExcel.SetColWidth("Necesidades", "G", "G", 35)
	necesidadesExcel.SetColWidth("Necesidades", "H", "I", 12)
	necesidadesExcel.SetColWidth("Necesidades", "J", "J", 35)

	necesidadesExcel.SetCellValue("Necesidades", "B1", "Código del rubro")
	necesidadesExcel.SetCellValue("Necesidades", "C1", "Nombre del rubro")
	necesidadesExcel.SetCellValue("Necesidades", "F1", "Valor")
	necesidadesExcel.SetCellValue("Necesidades", "G1", "Dependencias")
	necesidadesExcel.SetCellStyle("Necesidades", "B1", "G1", stylehead)

	contador := 2
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",nombre:"+nombre, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planes)
		for i := 0; i < len(planes); i++ {
			flag := true
			var docentes map[string]interface{}
			var aux map[string]interface{}
			var dependencia_nombre string
			dependencia := planes[i]["dependencia_id"]
			if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia/"+dependencia.(string), &aux); err == nil {
				dependencia_nombre = aux["Nombre"].(string)
			}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+planes[i]["_id"].(string)+",activo:true", &respuestaIdentificaciones); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuestaIdentificaciones, &identificaciones)
				for i := 0; i < len(identificaciones); i++ {
					identificacion := identificaciones[i]
					nombre := strings.ToLower(identificacion["nombre"].(string))
					if strings.Contains(nombre, "recurso") {
						if identificacion["dato"] != nil {
							var dato map[string]interface{}
							var data_identi []map[string]interface{}
							dato_str := identificacion["dato"].(string)
							json.Unmarshal([]byte(dato_str), &dato)
							for key := range dato {
								element := dato[key].(map[string]interface{})
								if element["activo"] == true {
									data_identi = append(data_identi, element)
								}
							}
							recursos = data_identi

						}
					}
					if strings.Contains(nombre, "contratista") && flag {
						if identificacion["dato"] != nil && identificacion["dato"].(string) != "{}" {
							var dato map[string]interface{}
							var dato_contratistas []map[string]interface{}
							dato_str := identificacion["dato"].(string)
							json.Unmarshal([]byte(dato_str), &dato)
							element := dato["0"].(map[string]interface{})
							if element["activo"] == true {
								dato_contratistas = append(dato_contratistas, element)
								flag = false
							}
							rubros = dato_contratistas

						}
					} else if strings.Contains(nombre, "docente") {
						dato := map[string]interface{}{}
						var data_identi []map[string]interface{}
						if identificacion["dato"] != nil && identificacion["dato"] != "{}" {
							result := make(map[string]interface{})
							dato_str := identificacion["dato"].(string)

							// ? Se tiene en cuenta la nueva estructura la info ahora está en identificacion-detalle, pero tambien tiene en cuenta la estructura de indentificaciones viejas (else)
							if strings.Contains(dato_str, "ids_detalle") {
								json.Unmarshal([]byte(dato_str), &dato)

								var identi map[string]interface{}
								iddetail := ""
								identificacionDetalle := map[string]interface{}{}
								errIdentificacionDetalle := error(nil)

								iddetail = dato["ids_detalle"].(map[string]interface{})["rhf"].(string)
								errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
								if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
									dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
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
									data_identi = nil
									identificacionDetalle = map[string]interface{}{}
								} else {
									result["rhf"] = "{}"
								}

								iddetail = dato["ids_detalle"].(map[string]interface{})["rhv_pre"].(string)
								errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
								if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
									dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
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
									data_identi = nil
									identificacionDetalle = map[string]interface{}{}
								} else {
									result["rhv_pre"] = "{}"
								}

								iddetail = dato["ids_detalle"].(map[string]interface{})["rhv_pos"].(string)
								errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
								if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
									dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
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
									data_identi = nil
									identificacionDetalle = map[string]interface{}{}
								} else {
									result["rhv_pos"] = "{}"
								}

								iddetail = dato["ids_detalle"].(map[string]interface{})["rubros"].(string)
								errIdentificacionDetalle = request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion-detalle/"+iddetail, &identificacionDetalle)
								if errIdentificacionDetalle == nil && identificacionDetalle["Status"] == "200" && identificacionDetalle["Data"] != nil {
									dato_aux := identificacionDetalle["Data"].(map[string]interface{})["dato"].(string)
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
									data_identi = nil
									identificacionDetalle = map[string]interface{}{}
								} else {
									result["rubros"] = "{}"
								}

							} else {
								json.Unmarshal([]byte(dato_str), &dato)

								var identi map[string]interface{}
								dato_aux := dato["rhf"].(string)
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

								data_identi = nil
							}

							docentes = result
						}
					}
				}
				for i := 0; i < len(recursos); i++ {
					var aux bool
					var aux1 []string
					if len(recursosGeneral) == 0 {
						recursosGeneral = append(recursosGeneral, recursos[i])

						var valorU []float64
						var auxValor float64
						if fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "int" {
							auxValor = recursos[i]["valor"].(float64)
						} else {
							strValor2 := strings.TrimLeft(recursos[i]["valor"].(string), "$")
							strValor2 = strings.ReplaceAll(strValor2, ",", "")
							arrValor2 := strings.Split(strValor2, ".")
							aux2, err := strconv.ParseFloat(arrValor2[0], 64)
							if err == nil {
								auxValor = aux2
							}
						}
						index := len(recursosGeneral) - 1
						valorU = append(valorU, auxValor)
						recursosGeneral[index]["valorU"] = valorU
						aux1 = append(aux1, dependencia_nombre)
						recursosGeneral[index]["unidades"] = aux1
						unidades_total = append(unidades_total, dependencia_nombre)
					} else {
						for j := 0; j < len(recursosGeneral); j++ {
							if recursosGeneral[j]["codigo"] == recursos[i]["codigo"] {
								if recursosGeneral[j]["unidades"] != nil {
									recursosGeneral[j]["unidades"] = append(recursosGeneral[j]["unidades"].([]string), dependencia_nombre)
								}
								flag1 := false
								for k := 0; k < len(unidades_total); k++ {
									if unidades_total[k] == dependencia_nombre {
										flag1 = true
									}
								}
								if !flag1 {
									unidades_total = append(unidades_total, dependencia_nombre)
								}

								if recursosGeneral[j]["valor"] != nil {
									var auxValor float64
									var auxValor2 float64
									if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" || fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "float64" {
										auxValor = recursosGeneral[j]["valor"].(float64)
									} else {
										strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
										strValor = strings.ReplaceAll(strValor, ",", "")
										arrValor := strings.Split(strValor, ".")
										aux1, err := strconv.ParseFloat(arrValor[0], 64)
										if err == nil {
											auxValor = aux1
										}
									}

									if fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "int" || fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "float64" {
										auxValor2 = recursos[i]["valor"].(float64)
									} else {
										strValor2 := strings.TrimLeft(recursos[i]["valor"].(string), "$")
										strValor2 = strings.ReplaceAll(strValor2, ",", "")
										arrValor2 := strings.Split(strValor2, ".")
										aux2, err := strconv.ParseFloat(arrValor2[0], 64)
										if err == nil {
											auxValor2 = aux2
										}
									}

									recursosGeneral[j]["valor"] = auxValor + auxValor2
									if recursosGeneral[j]["valorU"] == nil {
										var valorU []float64
										valorU = append(valorU, auxValor2)
										recursosGeneral[j]["valorU"] = valorU
									} else {
										valorU := recursosGeneral[j]["valorU"].([]float64)
										recursosGeneral[j]["valorU"] = append(valorU, auxValor2)
									}
								} else {
									var valorU []float64
									var auxValor float64
									if fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "int" || fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "float64" {
										auxValor := recursos[i]["valor"].(float64)
										valorU = append(valorU, auxValor)
									} else {
										strValor2 := strings.TrimLeft(recursos[i]["valor"].(string), "$")
										strValor2 = strings.ReplaceAll(strValor2, ",", "")
										arrValor2 := strings.Split(strValor2, ".")
										auxValor, err := strconv.ParseFloat(arrValor2[0], 64)
										if err == nil {
											valorU = append(valorU, auxValor)
										}
									}
									recursosGeneral[j]["valor"] = auxValor
									recursosGeneral[j]["valorU"] = valorU
								}
								aux = true
								break
							} else {
								aux = false
							}
						}
						if !aux {
							flag := false
							var valorU []float64
							var auxValor float64
							if fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "int" || fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "flaot64" {
								auxValor = recursos[i]["valor"].(float64)
							} else {
								strValor2 := strings.TrimLeft(recursos[i]["valor"].(string), "$")
								strValor2 = strings.ReplaceAll(strValor2, ",", "")
								arrValor2 := strings.Split(strValor2, ".")
								aux2, err := strconv.ParseFloat(arrValor2[0], 64)
								if err == nil {
									auxValor = aux2
								}
							}
							recursosGeneral = append(recursosGeneral, recursos[i])
							index := len(recursosGeneral) - 1
							valorU = append(valorU, auxValor)
							recursosGeneral[index]["valorU"] = valorU
							aux1 = append(aux1, dependencia_nombre)
							recursosGeneral[index]["unidades"] = aux1
							for k := 0; k < len(unidades_total); k++ {
								if unidades_total[k] == dependencia_nombre {
									flag = true
								}
							}
							if !flag {
								unidades_total = append(unidades_total, dependencia_nombre)
							}
						}
					}
				}

				for i := 0; i < len(rubros); i++ {
					var aux bool
					var aux1 []string
					if len(rubrosGeneral) == 0 {
						var auxValor2 float64
						var valorU []float64
						if _, ok := rubros[i]["totalInc"].(float64); ok {
							rubros[i]["totalInc"] = fmt.Sprintf("%f", rubros[i]["totalInc"])
						}
						if rubros[i]["totalInc"] != nil {
							auxValor2, _ = strconv.ParseFloat(rubros[i]["totalInc"].(string), 64)
						} else {
							auxValor2 = 0.0
						}

						rubrosGeneral = append(rubrosGeneral, rubros[i])
						index := len(rubrosGeneral) - 1
						valorU = append(valorU, auxValor2)
						rubrosGeneral[index]["valorU"] = valorU
						aux1 = append(aux1, dependencia_nombre)
						rubrosGeneral[index]["unidades"] = aux1
						unidades_rubros_total = append(unidades_rubros_total, dependencia_nombre)
					} else {
						for j := 0; j < len(rubrosGeneral); j++ {
							if rubrosGeneral[j]["rubro"] == rubros[i]["rubro"] {
								flag := false
								for k := 0; k < len(rubrosGeneral[j]["unidades"].([]string)); k++ {
									aux2 := rubrosGeneral[j]["unidades"].([]string)
									if aux2[k] == dependencia_nombre {
										flag = true
									}
								}
								if !flag {
									rubrosGeneral[j]["unidades"] = append(rubrosGeneral[j]["unidades"].([]string), dependencia_nombre)
								}
								flag2 := false
								for k := 0; k < len(unidades_total); k++ {
									if unidades_total[k] == dependencia_nombre {
										flag2 = true
									}
								}
								if !flag2 {
									unidades_total = append(unidades_total, dependencia_nombre)
								}
								flag1 := false
								for k := 0; k < len(unidades_rubros_total); k++ {
									if unidades_rubros_total[k] == dependencia_nombre {
										flag1 = true
									}
								}
								if !flag1 {
									unidades_rubros_total = append(unidades_rubros_total, dependencia_nombre)
								}
								if rubrosGeneral[j]["totalInc"] != nil {
									var auxValor float64
									var auxValor2 float64
									if _, ok := rubrosGeneral[j]["totalInc"].(float64); ok {
										rubrosGeneral[j]["totalInc"] = fmt.Sprintf("%f", rubrosGeneral[j]["totalInc"])
									}
									auxValor, _ = strconv.ParseFloat(rubrosGeneral[j]["totalInc"].(string), 64)
									auxValor2, _ = strconv.ParseFloat(rubros[i]["totalInc"].(string), 64)

									if rubrosGeneral[j]["valorU"] == nil {
										var valorU []float64
										valorU = append(valorU, auxValor2)
										rubrosGeneral[j]["valorU"] = valorU
									} else {
										rubrosGeneral[j]["valorU"] = append(rubrosGeneral[j]["valorU"].([]float64), auxValor2)
									}

									rubrosGeneral[j]["totalInc"] = auxValor + auxValor2
								} else {
									rubrosGeneral[j]["valorU"] = append(rubrosGeneral[j]["valorU"].([]float64), 0.0)
									rubrosGeneral[j]["totalInc"] = rubros[i]["totalInc"]
								}
								aux = true
								break
							} else {
								aux = false
							}
						}
						if !aux {
							flag := false
							var valorU []float64
							var auxValor float64
							// ? puede haber recursos[] sin datos
							if len(recursos) > 0 {
								if fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "int" || fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "float64" {
									auxValor = recursos[i]["valor"].(float64)
								} else {
									strValor2 := strings.TrimLeft(recursos[i]["valor"].(string), "$")
									strValor2 = strings.ReplaceAll(strValor2, ",", "")
									arrValor2 := strings.Split(strValor2, ".")
									aux2, err := strconv.ParseFloat(arrValor2[0], 64)
									if err == nil {
										auxValor = aux2
									}
								}
								rubrosGeneral = append(rubrosGeneral, recursos[i])

								index := len(rubrosGeneral) - 1
								valorU = append(valorU, auxValor)
								rubrosGeneral[index]["valorU"] = valorU
								aux1 = append(aux1, dependencia_nombre)
								rubrosGeneral[index]["unidades"] = aux1
								for k := 0; k < len(unidades_total); k++ {
									if unidades_total[k] == dependencia_nombre {
										flag = true
									}
								}
							}
							if !flag {
								unidades_total = append(unidades_total, dependencia_nombre)
							}
						}
						if !aux {
							flag := false
							var valorU []float64
							var auxValor2 float64
							if rubros[i]["totalInc"] != nil {
								if _, ok := rubros[i]["totalInc"].(float64); ok {
									rubros[i]["totalInc"] = fmt.Sprintf("%f", rubros[i]["totalInc"])
								}
								auxValor2, _ = strconv.ParseFloat(rubros[i]["totalInc"].(string), 64)
								valorU = append(valorU, auxValor2)
							} else {
								valorU = append(valorU, 0.0)
							}

							rubrosGeneral = append(rubrosGeneral, rubros[i])
							index := len(rubrosGeneral) - 1
							rubrosGeneral[index]["valorU"] = valorU
							aux1 = append(aux1, dependencia_nombre)
							rubrosGeneral[index]["unidades"] = aux1
							for k := 0; k < len(unidades_rubros_total); k++ {
								if unidades_rubros_total[k] == dependencia_nombre {
									flag = true
								}
							}
							if !flag {
								unidades_rubros_total = append(unidades_rubros_total, dependencia_nombre)
							}
						}
					}
				}

				if len(docentes) > 0 {
					docentesGeneral = reporteshelper.TotalDocentes(docentes)
					primaServicios = primaServicios + docentesGeneral["primaServicios"].(int)
					primaNavidad = primaNavidad + docentesGeneral["primaNavidad"].(int)
					primaVacaciones = primaVacaciones + docentesGeneral["primaVacaciones"].(int)
					bonificacion = bonificacion + docentesGeneral["bonificacion"].(int)
					interesesCesantias = interesesCesantias + docentesGeneral["interesesCesantias"].(int)
					cesantiasPublicas = cesantiasPublicas + docentesGeneral["cesantiasPublicas"].(int)
					cesantiasPrivadas = cesantiasPrivadas + docentesGeneral["cesantiasPrivadas"].(int)
					salud = salud + docentesGeneral["salud"].(int)
					pensionesPublicas = pensionesPublicas + docentesGeneral["pensionesPublicas"].(int)
					pensionesPrivadas = pensionesPrivadas + docentesGeneral["pensionesPrivadas"].(int)
					arl = arl + docentesGeneral["arl"].(int)
					caja = caja + docentesGeneral["caja"].(int)
					icbf = icbf + docentesGeneral["icbf"].(int)
					arrDataDocentes = append(arrDataDocentes, reporteshelper.GetDataDocentes(docentes, planes[i]["dependencia_id"].(string)))
				}

				if docentes["rubros"] != nil {
					var aux bool
					var respuestaRubro map[string]interface{}
					rubros := docentes["rubros"].([]map[string]interface{})
					for i := 0; i < len(rubros); i++ {
						if rubros[i]["rubro"] != "" {
							for j := 0; j < len(recursosGeneral); j++ {
								if recursosGeneral[j]["codigo"] == rubros[i]["rubro"] {
									aux = true
									categoria := strings.ToLower(rubros[i]["categoria"].(string))
									if strings.Contains(categoria, "prima") && strings.Contains(categoria, "servicio") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + primaServicios
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + primaServicios
												}
											}
										} else {
											recursosGeneral[j]["valor"] = primaServicios
										}
									}

									if strings.Contains(categoria, "prima") && strings.Contains(categoria, "navidad") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + primaNavidad
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + primaNavidad
												}

											}

										} else {
											recursosGeneral[j]["valor"] = primaNavidad
										}
									}

									if strings.Contains(categoria, "prima") && strings.Contains(categoria, "vacaciones") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + primaVacaciones
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + primaVacaciones
												}
											}
										} else {
											recursosGeneral[j]["valor"] = primaVacaciones
										}
									}

									if strings.Contains(categoria, "bonificacion") || strings.Contains(categoria, "bonificación") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + bonificacion
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + bonificacion
												}
											}

										} else {
											recursosGeneral[j]["valor"] = bonificacion
										}
									}

									if strings.Contains(categoria, "interes") && strings.Contains(categoria, "cesantía") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + interesesCesantias
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + interesesCesantias
												}
											}

										} else {
											recursosGeneral[j]["valor"] = interesesCesantias
										}
									}

									if strings.Contains(categoria, "cesantía") && strings.Contains(categoria, "público") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + cesantiasPublicas
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + cesantiasPublicas
												}
											}
										} else {
											recursosGeneral[j]["valor"] = cesantiasPublicas
										}
									}

									if strings.Contains(categoria, "cesantía") && strings.Contains(categoria, "privado") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + cesantiasPrivadas
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + cesantiasPrivadas
												}
											}
										} else {
											recursosGeneral[j]["valor"] = cesantiasPrivadas
										}
									}

									if strings.Contains(categoria, "salud") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + salud
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + salud
												}
											}
										} else {
											recursosGeneral[j]["valor"] = salud
										}
									}

									if strings.Contains(categoria, "pension") && strings.Contains(categoria, "público") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + pensionesPublicas
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + pensionesPublicas
												}
											}

										} else {
											recursosGeneral[j]["valor"] = pensionesPublicas
										}
									}

									if strings.Contains(categoria, "pension") && strings.Contains(categoria, "privado") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + pensionesPrivadas
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + pensionesPrivadas
												}
											}

										} else {
											recursosGeneral[j]["valor"] = pensionesPrivadas
										}
									}

									if strings.Contains(categoria, "arl") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + arl
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + arl
												}
											}

										} else {
											recursosGeneral[j]["valor"] = arl
										}
									}

									if strings.Contains(categoria, "ccf") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + caja
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + caja
												}
											}

										} else {
											recursosGeneral[j]["valor"] = caja
										}
									}

									if strings.Contains(categoria, "icbf") {
										if recursosGeneral[j]["valor"] != nil {
											if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
												recursosGeneral[j]["valor"] = recursosGeneral[j]["valor"].(int) + icbf
											} else {
												strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													recursosGeneral[j]["valor"] = auxValor + icbf
												}
											}

										} else {
											recursosGeneral[j]["valor"] = icbf
										}
									}
									break
								} else {
									aux = false
								}
							}
							if !aux && rubros[i]["rubro"] != nil {
								rubro := make(map[string]interface{})
								if err := request.GetJson("http://"+beego.AppConfig.String("PlanCuentasService")+"/arbol_rubro/"+rubros[i]["rubro"].(string), &respuestaRubro); err == nil {
									if respuestaRubro["Body"] == nil {
										continue
									}
									aux := respuestaRubro["Body"].(map[string]interface{})
									rubro["codigo"] = aux["Codigo"]
									rubro["nombre"] = aux["Nombre"]
									rubro["categoria"] = rubros[i]["categoria"]

									if rubro["categoria"] != nil {
										categoria := strings.ToLower(rubro["categoria"].(string))

										if strings.Contains(categoria, "prima") && strings.Contains(categoria, "servicio") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + primaServicios
												}
											} else {
												rubro["valor"] = primaServicios
											}
										}

										if strings.Contains(categoria, "prima") && strings.Contains(categoria, "navidad") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + primaNavidad
												}
											} else {
												rubro["valor"] = primaNavidad
											}
										}

										if strings.Contains(categoria, "prima") && strings.Contains(categoria, "vacaciones") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + primaVacaciones
												}
											} else {
												rubro["valor"] = primaVacaciones
											}
										}

										if strings.Contains(categoria, "bonificacion") || strings.Contains(categoria, "bonificación") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + bonificacion
												}
											} else {
												rubro["valor"] = bonificacion
											}
										}

										if strings.Contains(categoria, "interes") && strings.Contains(categoria, "cesantía") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + interesesCesantias
												}
											} else {
												rubro["valor"] = interesesCesantias
											}
										}

										if strings.Contains(categoria, "cesantía") && strings.Contains(categoria, "público") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + cesantiasPublicas
												}
											} else {
												rubro["valor"] = cesantiasPublicas
											}
										}

										if strings.Contains(categoria, "cesantía") && strings.Contains(categoria, "privado") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + cesantiasPrivadas
												}
											} else {
												rubro["valor"] = cesantiasPrivadas
											}
										}

										if strings.Contains(categoria, "salud") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + salud
												}
											} else {
												rubro["valor"] = salud
											}
										}

										if strings.Contains(categoria, "pension") && strings.Contains(categoria, "público") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + pensionesPublicas
												}
											} else {
												rubro["valor"] = pensionesPublicas
											}
										}

										if strings.Contains(categoria, "pension") && strings.Contains(categoria, "privado") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + pensionesPrivadas
												}
											} else {
												rubro["valor"] = pensionesPrivadas
											}
										}

										if strings.Contains(categoria, "arl") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + arl
												}
											} else {
												rubro["valor"] = arl
											}
										}

										if strings.Contains(categoria, "ccf") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + caja
												}
											} else {
												rubro["valor"] = caja
											}
										}

										if strings.Contains(categoria, "icbf") {
											if rubro["valor"] != nil {
												strValor := strings.TrimLeft(rubro["valor"].(string), "$")
												strValor = strings.ReplaceAll(strValor, ",", "")
												arrValor := strings.Split(strValor, ".")
												auxValor, err := strconv.Atoi(arrValor[0])
												if err == nil {
													rubro["valor"] = auxValor + icbf
												}
											} else {
												rubro["valor"] = icbf
											}
										}
									}
									recursosGeneral = append(recursosGeneral, rubro)
								}
							}

						}
					}
				}
			}
		}

		for i := 0; i < len(recursosGeneral); i++ {
			if recursosGeneral[i]["categoria"] != nil {
				categoria := strings.ToLower(recursosGeneral[i]["categoria"].(string))
				if strings.Contains(categoria, "prima") && strings.Contains(categoria, "servicio") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + primaServicios
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + primaServicios
							}
						}
					} else {
						recursosGeneral[i]["valor"] = primaServicios
					}
				}

				if strings.Contains(categoria, "prima") && strings.Contains(categoria, "navidad") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + primaNavidad
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + primaNavidad
							}
						}
					} else {
						recursosGeneral[i]["valor"] = primaNavidad
					}
				}

				if strings.Contains(categoria, "prima") && strings.Contains(categoria, "vacaciones") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + primaVacaciones
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + primaVacaciones
							}
						}
					} else {
						recursosGeneral[i]["valor"] = primaVacaciones
					}
				}

				if strings.Contains(categoria, "bonificacion") || strings.Contains(categoria, "bonificación") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + bonificacion
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + bonificacion
							}
						}
					} else {
						recursosGeneral[i]["valor"] = bonificacion
					}
				}

				if strings.Contains(categoria, "interes") && strings.Contains(categoria, "cesantía") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + interesesCesantias
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + interesesCesantias
							}
						}
					} else {
						recursosGeneral[i]["valor"] = interesesCesantias
					}
				}

				if strings.Contains(categoria, "cesantía") && strings.Contains(categoria, "público") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + cesantiasPublicas
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + cesantiasPublicas
							}
						}
					} else {
						recursosGeneral[i]["valor"] = cesantiasPublicas
					}
				}

				if strings.Contains(categoria, "cesantía") && strings.Contains(categoria, "privado") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + cesantiasPrivadas
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + cesantiasPrivadas
							}
						}
					} else {
						recursosGeneral[i]["valor"] = cesantiasPrivadas
					}
				}

				if strings.Contains(categoria, "salud") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + salud
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + salud
							}
						}
					} else {
						recursosGeneral[i]["valor"] = salud
					}
				}

				if strings.Contains(categoria, "pension") && strings.Contains(categoria, "público") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + pensionesPublicas
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + pensionesPublicas
							}
						}
					} else {
						recursosGeneral[i]["valor"] = pensionesPublicas
					}
				}

				if strings.Contains(categoria, "pension") && strings.Contains(categoria, "privado") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + pensionesPrivadas
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + pensionesPrivadas
							}
						}
					} else {
						recursosGeneral[i]["valor"] = pensionesPrivadas
					}
				}

				if strings.Contains(categoria, "arl") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + arl
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + arl
							}
						}
					} else {
						recursosGeneral[i]["valor"] = arl
					}
				}

				if strings.Contains(categoria, "ccf") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + caja
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + caja
							}
						}
					} else {
						recursosGeneral[i]["valor"] = caja
					}
				}

				if strings.Contains(categoria, "icbf") {
					if recursosGeneral[i]["valor"] != nil {
						if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
							recursosGeneral[i]["valor"] = recursosGeneral[i]["valor"].(int) + icbf
						} else {
							strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
							strValor = strings.ReplaceAll(strValor, ",", "")
							arrValor := strings.Split(strValor, ".")
							auxValor, err := strconv.Atoi(arrValor[0])
							if err == nil {
								recursosGeneral[i]["valor"] = auxValor + icbf
							}
						}
					} else {
						recursosGeneral[i]["valor"] = icbf
					}
				}
			}
		}

		if len(arrDataDocentes) < 7 {
			ingenieria := true
			ciencias := true
			ambiente := true
			tecnologica := true
			ASAB := true
			ILUD := true
			matematicas := true
			for _, data := range arrDataDocentes {
				facultad := strings.ToLower(data["nombreFacultad"].(string))
				if strings.Contains(facultad, "ingenieria") {
					ingenieria = false
				}
				if strings.Contains(facultad, "ciencias") {
					ciencias = false
				}
				if strings.Contains(facultad, "medio ambiente") {
					ambiente = false
				}
				if strings.Contains(facultad, "tecnologica") {
					tecnologica = false
				}
				if strings.Contains(facultad, "asab") {
					ASAB = false
				}
				if strings.Contains(facultad, "ilud") {
					ILUD = false
				}
				if strings.Contains(facultad, "matematicas") {
					matematicas = false
				}
			}
			if ingenieria {
				vacio := map[string]interface{}{"hch": 0, "hchPos": 0, "hcp": 0, "hcpPos": 0, "mto": 0, "nombreFacultad": "FACULTAD DE INGENIERIA", "tco": 0, "valorPos": 0, "valorPre": 0}
				arrDataDocentes = append(arrDataDocentes, vacio)
			}
			if ciencias {
				vacio := map[string]interface{}{"hch": 0, "hchPos": 0, "hcp": 0, "hcpPos": 0, "mto": 0, "nombreFacultad": "FACULTAD DE CIENCIAS Y EDUCACION", "tco": 0, "valorPos": 0, "valorPre": 0}
				arrDataDocentes = append(arrDataDocentes, vacio)
			}
			if ambiente {
				vacio := map[string]interface{}{"hch": 0, "hchPos": 0, "hcp": 0, "hcpPos": 0, "mto": 0, "nombreFacultad": "FACULTAD DE MEDIO AMBIENTE", "tco": 0, "valorPos": 0, "valorPre": 0}
				arrDataDocentes = append(arrDataDocentes, vacio)
			}
			if tecnologica {
				vacio := map[string]interface{}{"hch": 0, "hchPos": 0, "hcp": 0, "hcpPos": 0, "mto": 0, "nombreFacultad": "FACULTAD TECNOLOGICA", "tco": 0, "valorPos": 0, "valorPre": 0}
				arrDataDocentes = append(arrDataDocentes, vacio)
			}
			if ASAB {
				vacio := map[string]interface{}{"hch": 0, "hchPos": 0, "hcp": 0, "hcpPos": 0, "mto": 0, "nombreFacultad": "FACULTAD DE ARTES - ASAB", "tco": 0, "valorPos": 0, "valorPre": 0}
				arrDataDocentes = append(arrDataDocentes, vacio)
			}
			if ILUD {
				vacio := map[string]interface{}{"hch": 0, "hchPos": 0, "hcp": 0, "hcpPos": 0, "mto": 0, "nombreFacultad": "INSTITUTO DE LENGUAS - ILUD", "tco": 0, "valorPos": 0, "valorPre": 0}
				arrDataDocentes = append(arrDataDocentes, vacio)
			}
			if matematicas {
				vacio := map[string]interface{}{"hch": 0, "hchPos": 0, "hcp": 0, "hcpPos": 0, "mto": 0, "nombreFacultad": "FACULTAD DE CIENCIAS MATEMATICAS Y NATURALES", "tco": 0, "valorPos": 0, "valorPre": 0}
				arrDataDocentes = append(arrDataDocentes, vacio)
			}
		}
		//Completado de tablas
		idActividad := 0
		for i := 0; i < len(recursosGeneral); i++ {
			necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), recursosGeneral[i]["codigo"])
			if recursosGeneral[i]["Nombre"] != nil {
				necesidadesExcel.MergeCell("Necesidades", "C"+fmt.Sprint(contador), "E"+fmt.Sprint(contador))
				necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), recursosGeneral[i]["Nombre"])

				if recursosGeneral[i]["unidades"] != nil {
					unidades := recursosGeneral[i]["unidades"].([]string)
					valores := recursosGeneral[i]["valorU"].([]float64)
					necesidadesExcel.MergeCell("Necesidades", "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+len(unidades)))
					necesidadesExcel.MergeCell("Necesidades", "C"+fmt.Sprint(contador), "C"+fmt.Sprint(contador+len(unidades)))
					reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador+len(unidades)), stylecontent, stylecontentS)
					for j := 0; j < len(unidades); j++ {
						necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador+j), unidades[j])
						necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador+j), valores[j])
						reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "F"+fmt.Sprint(contador+j), "F"+fmt.Sprint(contador+j), stylecontentCML, stylecontentCMLS)
						reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "G"+fmt.Sprint(contador+j), "G"+fmt.Sprint(contador+j), stylecontentCL, stylecontentCLS)
					}
					contador = contador + len(unidades)
				}
			} else {
				necesidadesExcel.MergeCell("Necesidades", "C"+fmt.Sprint(contador), "E"+fmt.Sprint(contador))
				necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), recursosGeneral[i]["nombre"])
				if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
					necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), recursosGeneral[i]["valor"])
					necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "Total")
				} else {
					strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
					strValor = strings.ReplaceAll(strValor, ",", "")
					arrValor := strings.Split(strValor, ".")
					auxValor, err := strconv.Atoi(arrValor[0])
					if err == nil {
						necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), auxValor)
						necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "Total")
					}
				}
			}
			reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontent, stylecontentS)
			reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "F"+fmt.Sprint(contador), "F"+fmt.Sprint(contador), stylecontentCMD, stylecontentCMD)
			reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "G"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentCD, stylecontentCD)

			if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" || fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "float64" {
				necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), recursosGeneral[i]["valor"])
				necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "Total")
			} else {
				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
				strValor = strings.ReplaceAll(strValor, ",", "")
				arrValor := strings.Split(strValor, ".")
				auxValor, err := strconv.Atoi(arrValor[0])
				if err == nil {
					necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), auxValor)
					necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "Total")
				}
			}
			idActividad = i
			contador++
		}

		for i := 0; i < len(rubrosGeneral); i++ {
			if rubrosGeneral[i]["rubro"] != nil {
				idActividad++
				necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), rubrosGeneral[i]["rubro"])
				if rubrosGeneral[i]["rubroNombre"] != nil {
					necesidadesExcel.MergeCell("Necesidades", "C"+fmt.Sprint(contador), "E"+fmt.Sprint(contador))
					necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), rubrosGeneral[i]["rubroNombre"])
					if rubrosGeneral[i]["unidades"] != nil {

						unidades := rubrosGeneral[i]["unidades"].([]string)
						valores := rubrosGeneral[i]["valorU"].([]float64)

						// TODO: Revisar este machete, hay menos valores que unidades, se repite la primera unidad, el problema se presenta más arriba.
						if unidades[0] == unidades[1] && (len(unidades)-1) == len(valores) {
							unidades = unidades[1:]
						}
						// --- end of machete

						necesidadesExcel.MergeCell("Necesidades", "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+len(unidades)))
						necesidadesExcel.MergeCell("Necesidades", "C"+fmt.Sprint(contador), "C"+fmt.Sprint(contador+len(unidades)))
						reporteshelper.SombrearCeldas(necesidadesExcel, idActividad, "Necesidades", "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador+len(unidades)), stylecontent, stylecontentS)
						for j := 0; j < len(unidades); j++ {
							necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador+j), unidades[j])
							if j < len(valores) {
								necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador+j), valores[j])
							} else {
								necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador+j), 0.0)
							}
							reporteshelper.SombrearCeldas(necesidadesExcel, idActividad, "Necesidades", "F"+fmt.Sprint(contador+j), "F"+fmt.Sprint(contador+j), stylecontentCML, stylecontentCMLS)
							reporteshelper.SombrearCeldas(necesidadesExcel, idActividad, "Necesidades", "G"+fmt.Sprint(contador+j), "G"+fmt.Sprint(contador+j), stylecontentCL, stylecontentCLS)
						}

						if len(rubrosGeneral[i]["unidades"].([]string)) == 1 {
							necesidadesExcel.MergeCell("Necesidades", "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+len(unidades)+1))
							necesidadesExcel.MergeCell("Necesidades", "C"+fmt.Sprint(contador), "C"+fmt.Sprint(contador+len(unidades)+1))
							reporteshelper.SombrearCeldas(necesidadesExcel, idActividad, "Necesidades", "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador+len(unidades)+1), stylecontent, stylecontentS)
							contador = contador + len(unidades) + 1
						} else {
							contador = contador + len(unidades)
						}
					}
				} else {
					necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), rubrosGeneral[i]["rubroNombre"])
					necesidadesExcel.MergeCell("Necesidades", "C"+fmt.Sprint(contador), "E"+fmt.Sprint(contador))
				}

				if fmt.Sprint(reflect.TypeOf(rubrosGeneral[i]["totalInc"])) == "float64" {
					necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), rubrosGeneral[i]["totalInc"])
					necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "Total")
				} else {
					if rubrosGeneral[i]["totalInc"] != nil {
						strValor := strings.TrimLeft(rubrosGeneral[i]["totalInc"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						auxValor, err := strconv.ParseFloat(strValor, 64)
						if err == nil {
							necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), auxValor)
							necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "Total")
						}
					}
				}

				reporteshelper.SombrearCeldas(necesidadesExcel, idActividad, "Necesidades", "B"+fmt.Sprint(contador), "E"+fmt.Sprint(contador), stylecontent, stylecontentS)
				reporteshelper.SombrearCeldas(necesidadesExcel, idActividad, "Necesidades", "F"+fmt.Sprint(contador), "F"+fmt.Sprint(contador), stylecontentCMD, stylecontentCMD)
				reporteshelper.SombrearCeldas(necesidadesExcel, idActividad, "Necesidades", "G"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentCD, stylecontentCD)

				contador++
			}
		}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-plan/"+planes[0]["estado_plan_id"].(string), &respuestaEstado); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaEstado, &estado)
		} else {
			panic(map[string]interface{}{"funcion": "getNecesidades", "err": "Error ", "status": "400", "log": err})
		}

		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-plan/"+planes[0]["tipo_plan_id"].(string), &respuestaTipo); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaTipo, &tipo)
		} else {
			panic(map[string]interface{}{"funcion": "getNecesidades", "err": "Error ", "status": "400", "log": err})
		}

		contador++
		necesidadesExcel.MergeCell("Necesidades", "B"+fmt.Sprint(contador), "J"+fmt.Sprint(contador))

		necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), "Docentes por tipo de vinculación:")
		necesidadesExcel.SetCellStyle("Necesidades", "B"+fmt.Sprint(contador), "J"+fmt.Sprint(contador), styletitles)

		contador++
		necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), "Facultad")
		necesidadesExcel.MergeCell("Necesidades", "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1))
		necesidadesExcel.MergeCell("Necesidades", "C"+fmt.Sprint(contador), "G"+fmt.Sprint(contador))

		necesidadesExcel.MergeCell("Necesidades", "H"+fmt.Sprint(contador), "J"+fmt.Sprint(contador))

		necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), "Pregrado")
		necesidadesExcel.SetCellValue("Necesidades", "H"+fmt.Sprint(contador), "Posgrado")

		necesidadesExcel.SetCellStyle("Necesidades", "B"+fmt.Sprint(contador), "J"+fmt.Sprint(contador), stylehead)
		necesidadesExcel.SetCellStyle("Necesidades", "B"+fmt.Sprint(contador), "B"+fmt.Sprint(contador+1), stylehead)

		contador++
		necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), "TCO")
		necesidadesExcel.SetCellValue("Necesidades", "D"+fmt.Sprint(contador), "MTO")
		necesidadesExcel.SetCellValue("Necesidades", "E"+fmt.Sprint(contador), "HCH")
		necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), "HCP")
		necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "Valor")
		necesidadesExcel.SetCellValue("Necesidades", "H"+fmt.Sprint(contador), "HCH")
		necesidadesExcel.SetCellValue("Necesidades", "I"+fmt.Sprint(contador), "HCP")
		necesidadesExcel.SetCellValue("Necesidades", "J"+fmt.Sprint(contador), "Valor")
		necesidadesExcel.SetCellStyle("Necesidades", "C"+fmt.Sprint(contador), "J"+fmt.Sprint(contador), styletitles)

		contador++

		tco := 0
		mto := 0
		hch := 0
		hcp := 0
		valorPre := 0.0
		hchPos := 0
		hcpPos := 0
		valorPos := 0.0

		for i := 0; i < len(arrDataDocentes); i++ {
			necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), arrDataDocentes[i]["nombreFacultad"])
			necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), arrDataDocentes[i]["tco"])
			if fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["tco"])) == "int" {
				tco += arrDataDocentes[i]["tco"].(int)
			} else {
				aux2, _ := strconv.Atoi(arrDataDocentes[i]["tco"].(string))
				tco += aux2
			}

			necesidadesExcel.SetCellValue("Necesidades", "D"+fmt.Sprint(contador), arrDataDocentes[i]["mto"])
			if fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["mto"])) == "int" {
				mto += arrDataDocentes[i]["mto"].(int)
			} else {
				aux2, _ := strconv.Atoi(arrDataDocentes[i]["mto"].(string))
				mto += aux2
			}

			necesidadesExcel.SetCellValue("Necesidades", "E"+fmt.Sprint(contador), arrDataDocentes[i]["hch"])
			if fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["hch"])) == "int" {
				hch += arrDataDocentes[i]["hch"].(int)
			} else {
				aux2, _ := strconv.Atoi(arrDataDocentes[i]["hch"].(string))
				hch += aux2
			}

			necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), arrDataDocentes[i]["hcp"])
			if fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["hcp"])) == "int" {
				hcp += arrDataDocentes[i]["hcp"].(int)
			} else {
				aux2, _ := strconv.Atoi(arrDataDocentes[i]["hcp"].(string))
				hcp += aux2
			}

			necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), arrDataDocentes[i]["valorPre"])
			if fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["valorPre"])) == "int" || fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["valorPre"])) == "float64" {
				valorPre += float64(arrDataDocentes[i]["valorPre"].(int))
			} else {
				aux2, _ := strconv.ParseFloat(arrDataDocentes[i]["valorPre"].(string), 64)
				valorPre += aux2
			}

			necesidadesExcel.SetCellValue("Necesidades", "H"+fmt.Sprint(contador), arrDataDocentes[i]["hchPos"])
			if fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["hchPos"])) == "int" {
				hchPos += arrDataDocentes[i]["hchPos"].(int)
			} else {
				aux2, _ := strconv.Atoi(arrDataDocentes[i]["hchPos"].(string))
				hchPos += aux2
			}

			necesidadesExcel.SetCellValue("Necesidades", "I"+fmt.Sprint(contador), arrDataDocentes[i]["hcpPos"])
			if fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["hcpPos"])) == "int" {
				hcpPos += arrDataDocentes[i]["hcpPos"].(int)
			} else {
				aux2, _ := strconv.Atoi(arrDataDocentes[i]["hcpPos"].(string))
				hcpPos += aux2
			}

			necesidadesExcel.SetCellValue("Necesidades", "J"+fmt.Sprint(contador), arrDataDocentes[i]["valorPos"])
			if fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["valorPos"])) == "int" || fmt.Sprint(reflect.TypeOf(arrDataDocentes[i]["valorPos"])) == "float64" {
				valorPos += float64(arrDataDocentes[i]["valorPos"].(int))
			} else {
				aux2, _ := strconv.ParseFloat(arrDataDocentes[i]["valorPos"].(string), 64)
				valorPos += aux2
			}

			reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "B"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), stylecontent, stylecontentS)
			reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "G"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentM, stylecontentMS)
			reporteshelper.SombrearCeldas(necesidadesExcel, i, "Necesidades", "J"+fmt.Sprint(contador), "J"+fmt.Sprint(contador), stylecontentM, stylecontentMS)
			contador++
		}

		necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), "Total")
		necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), tco)
		necesidadesExcel.SetCellValue("Necesidades", "D"+fmt.Sprint(contador), mto)
		necesidadesExcel.SetCellValue("Necesidades", "E"+fmt.Sprint(contador), hch)
		necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), hcp)
		necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), valorPre)
		necesidadesExcel.SetCellValue("Necesidades", "H"+fmt.Sprint(contador), hchPos)
		necesidadesExcel.SetCellValue("Necesidades", "I"+fmt.Sprint(contador), hcpPos)
		necesidadesExcel.SetCellValue("Necesidades", "J"+fmt.Sprint(contador), valorPos)

		reporteshelper.SombrearCeldas(necesidadesExcel, 0, "Necesidades", "B"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), stylecontentC, stylecontentC)
		reporteshelper.SombrearCeldas(necesidadesExcel, 0, "Necesidades", "G"+fmt.Sprint(contador), "G"+fmt.Sprint(contador), stylecontentCM, stylecontentCM)
		reporteshelper.SombrearCeldas(necesidadesExcel, 0, "Necesidades", "J"+fmt.Sprint(contador), "J"+fmt.Sprint(contador), stylecontentCM, stylecontentCM)

		styletitle, _ := necesidadesExcel.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{WrapText: true, Vertical: "center"},
			Font:      &excelize.Font{Bold: true, Size: 18, Color: "000000"},
			Border: []excelize.Border{
				{Type: "right", Color: "ffffff", Style: 1},
				{Type: "left", Color: "ffffff", Style: 1},
				{Type: "top", Color: "ffffff", Style: 1},
				{Type: "bottom", Color: "ffffff", Style: 1},
			},
		})
		necesidadesExcel.InsertRows("Necesidades", 1, 7)
		necesidadesExcel.MergeCell("Necesidades", "C2", "G6")
		necesidadesExcel.SetCellStyle("Necesidades", "C2", "G6", styletitle)
		var resPeriodo map[string]interface{}
		var periodo []map[string]interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Id:`+body["vigencia"].(string), &resPeriodo); err == nil {
			helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
		}

		if periodo[0] != nil {
			necesidadesExcel.SetCellValue("Necesidades", "C2", "Consolidado proyeccción de necesidades "+periodo[0]["Nombre"].(string))
		} else {
			necesidadesExcel.SetCellValue("Necesidades", "C2", "Consolidado proyeccción de necesidades")
		}

		if err := necesidadesExcel.AddPicture("Necesidades", "B1", "static/img/UDEscudo2.png",
			&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "oneCell", OffsetX: 25}); err != nil {
			fmt.Println(err)
		}

		contador = 1

		necesidadesExcel.NewSheet("Total Unidades")
		necesidadesExcel.MergeCell("Total Unidades", "A", "B")
		necesidadesExcel.MergeCell("Total Unidades", "A1", "A2")
		necesidadesExcel.SetColWidth("Total Unidades", "A", "B", 30)

		necesidadesExcel.MergeCell("Total Unidades", "A"+fmt.Sprint(contador), "B"+fmt.Sprint(contador))
		necesidadesExcel.MergeCell("Total Unidades", "A"+fmt.Sprint(contador), "A"+fmt.Sprint(contador+1))
		necesidadesExcel.SetCellValue("Total Unidades", "A"+fmt.Sprint(contador), "Total de unidades generadas:")
		necesidadesExcel.SetCellStyle("Total Unidades", "A"+fmt.Sprint(contador), "B"+fmt.Sprint(contador), styletitles)
		necesidadesExcel.SetCellStyle("Total Unidades", "A"+fmt.Sprint(contador+1), "B"+fmt.Sprint(contador+1), styletitles)
		contador++
		contador++
		necesidadesExcel.SetCellValue("Total Unidades", "A"+fmt.Sprint(contador), "Total de Unidades Generadas")
		necesidadesExcel.SetCellValue("Total Unidades", "B"+fmt.Sprint(contador), "Unidades Generadas")
		necesidadesExcel.SetCellStyle("Total Unidades", "A"+fmt.Sprint(contador), "B"+fmt.Sprint(contador), stylehead)
		contador++
		unid_total := ""
		for j := 0; j < len(unidades_total); j++ {
			infoReporte := make(map[string]interface{})
			infoReporte["vigencia"] = body["vigencia"].(string)
			infoReporte["estado_plan"] = estado["nombre"]
			infoReporte["tipo_plan"] = tipo["nombre"]
			infoReporte["nombre_unidad"] = unidades_total[j]
			unid_total = unid_total + unidades_total[j] + ", "
			arregloInfoReportes = append(arregloInfoReportes, infoReporte)
		}
		unid_total = strings.TrimRight(unid_total, ", ")
		necesidadesExcel.SetCellValue("Total Unidades", "A"+fmt.Sprint(contador), len(unidades_total))
		necesidadesExcel.SetCellValue("Total Unidades", "B"+fmt.Sprint(contador), unid_total)
		necesidadesExcel.SetCellStyle("Total Unidades", "A"+fmt.Sprint(contador), "B"+fmt.Sprint(contador), stylecontent)

		buf, _ := necesidadesExcel.WriteToBuffer()
		strings.NewReader(buf.String())
		encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
		dataSend := make(map[string]interface{})
		dataSend["generalData"] = arregloInfoReportes
		dataSend["excelB64"] = encoded

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}

	} else {
		panic(err)
	}

	c.ServeJSON()
}

// PlanAccionEvaluacion ...
// @Title PlanAccionEvaluacion
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Param	nombre		path 	string	true		"The key for staticblock"
// @Success 201 {object} models.Reportes
// @Failure 403 :nombre is empty
// @router /plan_anual_evaluacion/:nombre [post]
func (c *ReportesController) PlanAccionEvaluacion() {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "ReportesController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	var body map[string]interface{}
	var respuesta map[string]interface{}
	var planes []map[string]interface{}
	var arregloPlanAnual []map[string]interface{}
	var periodo []map[string]interface{}
	var unidadNombre string
	var evaluacion []map[string]interface{}
	var respuestaOikos []map[string]interface{}
	var resPeriodo map[string]interface{}
	nombre := c.Ctx.Input.Param(":nombre")
	consolidadoExcelEvaluacion := excelize.NewFile()
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if body["unidad_id"].(string) != "" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:6153355601c7a2365b2fb2a1,dependencia_id:"+body["unidad_id"].(string)+",nombre:"+nombre, &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planes)

			trimestres := evaluacionhelper.GetPeriodos(body["vigencia"].(string))

			if len(planes) <= 0 {
				c.Abort("404")
			}

			dependencia := body["unidad_id"].(string)
			if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia?query=Id:"+dependencia, &respuestaOikos); err == nil {
				unidadNombre = respuestaOikos[0]["Nombre"].(string)
				arregloPlanAnual = append(arregloPlanAnual, map[string]interface{}{"nombreUnidad": unidadNombre})
			} else {
				panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Id:`+body["vigencia"].(string), &resPeriodo); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
			}

			var index int
			for index = 3; index >= 0; index-- {
				evaluacion = evaluacionhelper.GetEvaluacion(planes[0]["_id"].(string), trimestres, index)
				if fmt.Sprintf("%v", evaluacion) != "[]" {
					break
				}
			}

			trimestreVacio := map[string]interface{}{"actividad": 0.0, "acumulado": 0.0, "denominador": 0.0, "meta": 0.0, "numerador": 0.0, "periodo": 0.0, "numeradorAcumulado": 0.0, "denominadorAcumulado": 0.0, "brecha": 0.0}

			switch index {
			case 2:
				for _, actividad := range evaluacion {
					actividad["trimestre4"] = trimestreVacio
				}
			case 1:
				for _, actividad := range evaluacion {
					actividad["trimestre4"] = trimestreVacio
					actividad["trimestre3"] = trimestreVacio
				}
			case 0:
				for _, actividad := range evaluacion {
					actividad["trimestre4"] = trimestreVacio
					actividad["trimestre3"] = trimestreVacio
					actividad["trimestre2"] = trimestreVacio
				}
			case -1:
				c.Abort("404")
			}

			sheetName := "Evaluación"
			consolidadoExcelEvaluacion.NewSheet(sheetName)
			consolidadoExcelEvaluacion.DeleteSheet("Sheet1")

			disable := false
			if err := consolidadoExcelEvaluacion.SetSheetView(sheetName, -1, &excelize.ViewOptions{
				ShowGridLines: &disable,
			}); err != nil {
				fmt.Println(err)
			}

			styleUnidad, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Vertical: "center"},
				Font:      &excelize.Font{Bold: true, Color: "000000", Family: "Bahnschrift SemiBold SemiConden", Size: 20},
				Border: []excelize.Border{
					{Type: "bottom", Color: "000000", Style: 2},
				},
			})
			styleTituloSB, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
				Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Family: "Calibri", Size: 11},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
			})
			styleSombreadoSB, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
				Font:      &excelize.Font{Color: "000000", Family: "Calibri", Size: 11},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
			})
			styleNegrilla, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Color: "000000", Family: "Calibri", Size: 12, Bold: true},
			})
			styleTitulo, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Family: "Calibri", Size: 11},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleContenido, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleContenidoC, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleContenidoCI, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
				Font:      &excelize.Font{Color: "FFFFFF"},
			})
			styleContenidoCIP, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				NumFmt:    10,
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
				Font:      &excelize.Font{Color: "FFFFFF"},
			})
			styleContenidoCE, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				NumFmt:    1,
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleContenidoCD, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				NumFmt:    4,
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleContenidoCS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FCE4D6"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleContenidoCP, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				NumFmt:    10,
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleContenidoCPSR, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				NumFmt:    10,
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Family: "Calibri", Size: 11},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"CC0000"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleContenidoCPS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				NumFmt:    10,
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"FCE4D6"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})

			// Size
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "A", "A", 3)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "B", "B", 4)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "B", "B", 4)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "C", "C", 8)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "D", "D", 13)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "D", "D", 13)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "E", "E", 42)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "F", "F", 16)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "G", "G", 21)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "H", "AS", 14)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "AV", "AY", 3)
			consolidadoExcelEvaluacion.SetRowHeight(sheetName, 1, 12)
			consolidadoExcelEvaluacion.SetRowHeight(sheetName, 2, 27)
			consolidadoExcelEvaluacion.SetRowHeight(sheetName, 19, 31)
			consolidadoExcelEvaluacion.SetRowHeight(sheetName, 22, 27)
			// Merge
			consolidadoExcelEvaluacion.MergeCell(sheetName, "B4", "D4")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "B19", "E19")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "AX19", "AY19")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "B21", "B22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "C21", "C22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "D21", "D22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "E21", "E22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "F21", "F22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "G21", "G22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "H21", "H22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "I21", "I22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "AX21", "AX22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "AY21", "AY22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "J21", "R21")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "S21", "AA21")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "AB21", "AJ21")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "AK21", "AS21")
			// Style
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B2", "T2", styleUnidad)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B4", "D4", styleTituloSB)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "E4", "E4", styleSombreadoSB)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B19", "B19", styleNegrilla)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B21", "AS22", styleTitulo)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "J22", "AS22", styleContenidoC)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "R22", "R22", styleContenidoCS)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AA22", "AA22", styleContenidoCS)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ22", "AJ22", styleContenidoCS)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS22", "AS22", styleContenidoCS)

			if periodo[0] != nil {
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "B2", "Evaluación Plan de Acción "+periodo[0]["Nombre"].(string)+" - "+unidadNombre)
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "B19", "Cumplimiento General Plan de Acción "+periodo[0]["Nombre"].(string)+" - "+unidadNombre)
			} else {
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "B2", "Evaluación Plan de Acción - "+unidadNombre)
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "B19", "Cumplimiento General Plan de Acción - "+unidadNombre)
			}

			ddRango1 := excelize.NewDataValidation(true)
			ddRango1.Sqref = "E4:E4"
			ddRango1.SetDropList([]string{"Trimestre I", "Trimestre II", "Trimestre III", "Trimestre IV"})

			if err = consolidadoExcelEvaluacion.AddDataValidation(sheetName, ddRango1); err != nil {
				fmt.Println(err)
				return
			}

			// Titles
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "B4", "Seleccione el periodo:")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "B21", "No.")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "C21", "Pond.")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "D21", "Periodo de ejecución")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "E21", "Actividad General")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "F21", "Indicador asociado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "G21", "Fórmula del Indicador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "H21", "Meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "I21", "Tipo de Unidad")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "J21", "Trimestre I")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "S21", "Trimestre II")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB21", "Trimestre III")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK21", "Trimestre IV")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "J22", "Numerador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "K22", "Denominador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "L22", "Indicador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "M22", "Numerador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "N22", "Denominador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "O22", "Indicador Acumulado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "P22", "Cumplimiento por meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q22", "Brecha")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "R22", "Cumplimiento por actividad")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "S22", "Numerador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "T22", "Denominador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "U22", "Indicador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "V22", "Numerador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "W22", "Denominador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "X22", "Indicador Acumulado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y22", "Cumplimiento por meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z22", "Brecha")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA22", "Cumplimiento por actividad")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB22", "Numerador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC22", "Denominador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD22", "Indicador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE22", "Numerador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF22", "Denominador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG22", "Indicador Acumulado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH22", "Cumplimiento por meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI22", "Brecha")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ22", "Cumplimiento por actividad")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK22", "Numerador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL22", "Denominador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM22", "Indicador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN22", "Numerador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO22", "Denominador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP22", "Indicador Acumulado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ22", "Cumplimiento por meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR22", "Brecha")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS22", "Cumplimiento por actividad")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX19", "Gráfica")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX21", "No.")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY21", "Cumplimiento")

			indice := 23
			indiceGraficos := 23
			for i, actividad := range evaluacion {
				// Datos
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "B"+fmt.Sprint(indice), actividad["numero"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "C"+fmt.Sprint(indice), actividad["ponderado"].(float64)/100)
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "D"+fmt.Sprint(indice), actividad["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "E"+fmt.Sprint(indice), actividad["actividad"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "F"+fmt.Sprint(indice), actividad["indicador"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "G"+fmt.Sprint(indice), actividad["formula"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "H"+fmt.Sprint(indice), actividad["meta"].(float64))
				if actividad["unidad"] == "Porcentaje" {
					consolidadoExcelEvaluacion.SetCellValue(sheetName, "H"+fmt.Sprint(indice), actividad["meta"].(float64)/100)
				}
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "I"+fmt.Sprint(indice), actividad["unidad"])

				// Trimestres
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "J"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre1"].(map[string]interface{})["numerador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "K"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre1"].(map[string]interface{})["denominador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "L"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "M"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["numeradorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "N"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["denominadorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "O"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["acumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "P"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["meta"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["brecha"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "R"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["actividad"].(float64))

				consolidadoExcelEvaluacion.SetCellValue(sheetName, "S"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre2"].(map[string]interface{})["numerador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "T"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre2"].(map[string]interface{})["denominador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "U"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "V"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["numeradorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "W"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["denominadorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "X"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["acumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["meta"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["brecha"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["actividad"].(float64))

				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre3"].(map[string]interface{})["numerador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre3"].(map[string]interface{})["denominador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["numeradorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["denominadorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["acumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["meta"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["brecha"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["actividad"].(float64))

				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre4"].(map[string]interface{})["numerador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre4"].(map[string]interface{})["denominador"]))
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["periodo"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["numeradorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["denominadorAcumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["acumulado"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["meta"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["brecha"])
				consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["actividad"].(float64))

				// Estilos
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice), styleContenidoC)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "C"+fmt.Sprint(indice), "C"+fmt.Sprint(indice), styleContenidoCP)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "E"+fmt.Sprint(indice), "E"+fmt.Sprint(indice), styleContenido)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "J"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice), styleContenidoCP)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "J"+fmt.Sprint(indice), "K"+fmt.Sprint(indice), styleContenidoCD)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "M"+fmt.Sprint(indice), "N"+fmt.Sprint(indice), styleContenidoCD)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "S"+fmt.Sprint(indice), "T"+fmt.Sprint(indice), styleContenidoCD)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "V"+fmt.Sprint(indice), "W"+fmt.Sprint(indice), styleContenidoCD)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AB"+fmt.Sprint(indice), "AC"+fmt.Sprint(indice), styleContenidoCD)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AE"+fmt.Sprint(indice), "AF"+fmt.Sprint(indice), styleContenidoCD)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AK"+fmt.Sprint(indice), "AL"+fmt.Sprint(indice), styleContenidoCD)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AN"+fmt.Sprint(indice), "AO"+fmt.Sprint(indice), styleContenidoCD)

				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "R"+fmt.Sprint(indice), "R"+fmt.Sprint(indice), styleContenidoCPS)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AA"+fmt.Sprint(indice), "AA"+fmt.Sprint(indice), styleContenidoCPS)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice), styleContenidoCPS)
				consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice), styleContenidoCPS)

				if actividad["unidad"] == "Porcentaje" {
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H"+fmt.Sprint(indice), "H"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "L"+fmt.Sprint(indice), "L"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "O"+fmt.Sprint(indice), "O"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Q"+fmt.Sprint(indice), "Q"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "U"+fmt.Sprint(indice), "U"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "X"+fmt.Sprint(indice), "X"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AD"+fmt.Sprint(indice), "AD"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AG"+fmt.Sprint(indice), "AG"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AI"+fmt.Sprint(indice), "AI"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AM"+fmt.Sprint(indice), "AM"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice), styleContenidoCP)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AR"+fmt.Sprint(indice), "AR"+fmt.Sprint(indice), styleContenidoCP)
				}

				if actividad["unidad"] == "Tasa" {
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H"+fmt.Sprint(indice), "H"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "L"+fmt.Sprint(indice), "L"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "O"+fmt.Sprint(indice), "O"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Q"+fmt.Sprint(indice), "Q"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "U"+fmt.Sprint(indice), "U"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "X"+fmt.Sprint(indice), "X"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AD"+fmt.Sprint(indice), "AD"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AG"+fmt.Sprint(indice), "AG"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AI"+fmt.Sprint(indice), "AI"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AM"+fmt.Sprint(indice), "AM"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice), styleContenidoCD)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AR"+fmt.Sprint(indice), "AR"+fmt.Sprint(indice), styleContenidoCD)
				}

				if actividad["unidad"] == "Unidad" {
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H"+fmt.Sprint(indice), "H"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "L"+fmt.Sprint(indice), "L"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "O"+fmt.Sprint(indice), "O"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Q"+fmt.Sprint(indice), "Q"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "U"+fmt.Sprint(indice), "U"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "X"+fmt.Sprint(indice), "X"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AD"+fmt.Sprint(indice), "AD"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AG"+fmt.Sprint(indice), "AG"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AI"+fmt.Sprint(indice), "AI"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AM"+fmt.Sprint(indice), "AM"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice), styleContenidoCE)
					consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AR"+fmt.Sprint(indice), "AR"+fmt.Sprint(indice), styleContenidoCE)
				}

				// Unión de celdas por indicador
				if i > 0 {
					if actividad["numero"] == evaluacion[i-1]["numero"] {
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "B"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "C"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "D"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "E"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "R"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX"+fmt.Sprint(indice), nil)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY"+fmt.Sprint(indice), nil)

						consolidadoExcelEvaluacion.MergeCell(sheetName, "B"+fmt.Sprint(indice-1), "B"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "C"+fmt.Sprint(indice-1), "C"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "D"+fmt.Sprint(indice-1), "D"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "E"+fmt.Sprint(indice-1), "E"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "R"+fmt.Sprint(indice-1), "R"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AA"+fmt.Sprint(indice-1), "AA"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AJ"+fmt.Sprint(indice-1), "AJ"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AS"+fmt.Sprint(indice-1), "AS"+fmt.Sprint(indice))
					} else {
						// Gaficos
						consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AX"+fmt.Sprint(indiceGraficos), "=B"+fmt.Sprint(indice))
						consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AY"+fmt.Sprint(indiceGraficos), "=IF(E4=\"Trimestre I\",R"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AA"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AJ"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",AS"+fmt.Sprint(indice)+"))))")
						indiceGraficos++
					}
				} else if i == 0 {
					// Gaficos
					consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AX"+fmt.Sprint(indiceGraficos), "=B"+fmt.Sprint(indice))
					consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AY"+fmt.Sprint(indiceGraficos), "=IF(E4=\"Trimestre I\",R"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AA"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AJ"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",AS"+fmt.Sprint(indice)+"))))")
					indiceGraficos++
				}
				indice++
			}

			consolidadoExcelEvaluacion.MergeCell(sheetName, "B"+fmt.Sprint(indice), "I"+fmt.Sprint(indice))
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B"+fmt.Sprint(indice), "I"+fmt.Sprint(indice), styleTituloSB)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "J"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice), styleContenidoC)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "R"+fmt.Sprint(indice), "R"+fmt.Sprint(indice), styleContenidoCPSR)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AA"+fmt.Sprint(indice), "AA"+fmt.Sprint(indice), styleContenidoCPSR)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice), styleContenidoCPSR)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice), styleContenidoCPSR)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AX19", "AX"+fmt.Sprint(indice+1), styleContenidoCI)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AY19", "AZ"+fmt.Sprint(indice+1), styleContenidoCIP)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AY21", "AY22", styleContenidoCI)

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "B"+fmt.Sprint(indice), "Avance General del Plan de Acción")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "E4", "Trimestre I") //
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "J"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "K"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "L"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "M"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "N"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "O"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "P"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q"+fmt.Sprint(indice), "-")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "S"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "T"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "U"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "V"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "W"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "X"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z"+fmt.Sprint(indice), "-")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI"+fmt.Sprint(indice), "-")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR"+fmt.Sprint(indice), "-")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX"+fmt.Sprint(indice), "General")

			filaAnt := fmt.Sprint(indice - 1)
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "R"+fmt.Sprint(indice), "=SUMPRODUCT(C23:C"+filaAnt+",R23:R"+filaAnt+")")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AA"+fmt.Sprint(indice), "=SUMPRODUCT(C23:C"+filaAnt+",AA23:AA"+filaAnt+")")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AJ"+fmt.Sprint(indice), "=SUMPRODUCT(C23:C"+filaAnt+",AJ23:AJ"+filaAnt+")")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AS"+fmt.Sprint(indice), "=SUMPRODUCT(C23:C"+filaAnt+",AS23:AS"+filaAnt+")")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AY"+fmt.Sprint(indice), "=IF(E4=\"Trimestre I\",R"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AA"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AJ"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",AS"+fmt.Sprint(indice)+"))))")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AZ"+fmt.Sprint(indice), "=100%-AY"+fmt.Sprint(indice))
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY"+fmt.Sprint(indice+1), "Avance")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AZ"+fmt.Sprint(indice+1), "Restante")

			consolidadoExcelEvaluacion.AddChart(sheetName, "B5", &excelize.Chart{
				Type: excelize.Pie,
				Series: []excelize.ChartSeries{
					{
						Name:       "",
						Categories: sheetName + "!$AY$" + fmt.Sprint(indice+1) + ":$AZ$" + fmt.Sprint(indice+1),
						Values:     sheetName + "!$AY$" + fmt.Sprint(indice) + ":$AZ$" + fmt.Sprint(indice),
					},
				},
				Format: excelize.GraphicOptions{
					ScaleX:          1.0,
					ScaleY:          1.0,
					OffsetX:         15,
					OffsetY:         10,
					LockAspectRatio: false,
					Locked:          &disable,
				},
				PlotArea: excelize.ChartPlotArea{
					ShowCatName:     false,
					ShowLeaderLines: false,
					ShowPercent:     true,
					ShowSerName:     false,
					ShowVal:         false,
				},
				ShowBlanksAs: "zero",
				Dimension: excelize.ChartDimension{
					Height: 265,
					Width:  454,
				},
				XAxis: excelize.ChartAxis{
					None: true,
				},
				YAxis: excelize.ChartAxis{
					None: true,
				},
			})

			consolidadoExcelEvaluacion.AddChart(sheetName, "F4", &excelize.Chart{
				Type: excelize.Col,
				Series: []excelize.ChartSeries{
					{
						Name:       "",
						Categories: sheetName + "!$AX$23:$AX$" + fmt.Sprint(indiceGraficos-1),
						Values:     sheetName + "!$AY$23:$AY$" + fmt.Sprint(indiceGraficos-1),
					},
				},
				Format: excelize.GraphicOptions{
					OffsetX:         15,
					LockAspectRatio: false,
					Locked:          &disable,
				},
				Dimension: excelize.ChartDimension{
					Height: 344,
					Width:  1605,
				},
				PlotArea: excelize.ChartPlotArea{
					ShowCatName:     false,
					ShowLeaderLines: false,
					ShowPercent:     false,
					ShowSerName:     false,
					ShowVal:         true,
				},
				YAxis: excelize.ChartAxis{
					MajorGridLines: true,
					Font:           excelize.Font{Family: "Calibri", Size: 9, Color: "000000"},
				},
				XAxis: excelize.ChartAxis{
					Font: excelize.Font{Family: "Calibri", Size: 9, Color: "000000"},
				},
				VaryColors:   &disable,
				ShowBlanksAs: "span",
			})

			buf, _ := consolidadoExcelEvaluacion.WriteToBuffer()
			strings.NewReader(buf.String())
			encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

			dataSend := make(map[string]interface{})
			dataSend["generalData"] = arregloPlanAnual
			dataSend["excelB64"] = encoded

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}
		} else {
			panic(err)
		}
	}

	c.ServeJSON()
}
