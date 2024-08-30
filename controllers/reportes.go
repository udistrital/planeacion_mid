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
	seguimientohelper "github.com/udistrital/planeacion_mid/helpers/seguimientoHelper"
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
	if body["categoria"].(string) == "Evaluacion" {
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
							if resFilter[i]["estado_plan_id"] == "6153355601c7a2365b2fb2a1" { //Estado Avalado
								noEstado = false
								res["mensaje"] = ""
								res["reporte"] = true
								break
							}
						}
					}
				}

				if noPlan {
					res["mensaje"] = "La unidad no tiene registros con el plan seleccionado o el tipo-plan no coincide con un Plan de Acción de Funcionamiento"
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
					res["mensaje"] = "No existen registros con el plan seleccionado o el tipo-plan no coincide con un Plan de Acción de Funcionamiento"
					res["reporte"] = false
				} else if noEstado {
					res["mensaje"] = "No existen registros con el estado y plan seleccionado"
					res["reporte"] = false
				}
			}
		} else {
			res["mensaje"] = "Ocurrio un error"
		}
	} else if body["categoria"].(string) == "Plan_Accion_Unidad" {
		// TODO: hacer validación de tipo de plan de acción para mostrar que los demás planes aún no están soportados por el sistema.
		// "Por favor verificar el tipo de plan de acción. Actualmente NO soportado por el módulo de reportes."
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
					res["mensaje"] = "La unidad no tiene registros con el plan seleccionado o el tipo-plan no coincide con un Plan de Acción de Funcionamiento"
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
			res["mensaje"] = "Ocurrió un error"
		}
	} else if body["categoria"].(string) == "Plan_Accion_General" {
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
	var datosReporte map[string]interface{}
	nombre := c.Ctx.Input.Param(":nombre")
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
						if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") {
							actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
							var arregloLineamieto []map[string]interface{}
							var arregloLineamietoPI []map[string]interface{}
							if len(actividades) == 1 {
								for index := range actividades {
									if val, ok := actividades[index]["index"].(float64); ok {
										actividades[index]["index"] = fmt.Sprintf("%v", int(val))
									}
								}
							} else {
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
							}

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
								arregloLineamieto = reporteshelper.ArbolArmonizacionV2(armonizacionTercerNivel.(string))
								arregloLineamietoPI = reporteshelper.ArbolArmonizacionPIV2(armonizacionTercerNivelPI.(string))

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

				datosReporte = map[string]interface{}{
					"plan_id":          plan_id,
					"planes":           planes,
					"planesFilter":     planesFilter,
					"periodo":          periodo,
					"arregloPlanAnual": arregloPlanAnual,
				}
			}

			//TODO: falta discrimininar por vigencia

			reporteGenerado, errorReporte := reporteshelper.ConstruirExcelPlanAccionUnidad(false, datosReporte)

			if errorReporte != nil {
				panic(map[string]interface{}{"funcion": "PlanAccionAnual", "err": "Error en la generación del excel", "status": "400", "log": err})
			}

			buf, _ := reporteGenerado.WriteToBuffer()
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
	// var resPeriodo map[string]interface{}
	// var periodo []map[string]interface{}
	var datosReporte map[string]interface{}
	// contadorGeneral := 4

	// consolidadoExcelPlanAnual := excelize.NewFile()
	nombre := c.Ctx.Input.Param(":nombre")
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",nombre:"+nombre+"&fields=_id,dependencia_id,estado_plan_id,tipo_plan_id", &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
		for _, planes := range planesFilter {
			if idUnidad != planes["dependencia_id"].(string) {
				if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+planes["dependencia_id"].(string), &respuestaUnidad); err == nil {
					planes["nombreUnidad"] = respuestaUnidad[0]["DependenciaId"].(map[string]interface{})["Nombre"].(string)
					fmt.Sprintf("http://" + beego.AppConfig.String("OikosService") + "/dependencia_tipo_dependencia?query=DependenciaId:" + planes["dependencia_id"].(string))
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
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") {
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
							arregloLineamieto = reporteshelper.ArbolArmonizacionV2(armonizacionTercerNivel.(string))
							arregloLineamietoPI = reporteshelper.ArbolArmonizacionPIV2(armonizacionTercerNivelPI.(string))

							generalData := make(map[string]interface{})

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

			datosReporte = map[string]interface{} {
				"vigencia_id": body["vigencia"].(string),
				"planesFilter": planesFilter,
				"arregloPlanAnual": arregloPlanAnual,
			}
		}

		//TODO: falta discrimininar por vigencia
		reporteGenerado, errorReporte := reporteshelper.ConstruirExcelPlanAccionGeneral(true, datosReporte)

		if errorReporte != nil {
			panic(map[string]interface{}{"funcion": "PlanAccionAnualGeneral", "err": "Error en la generación del excel", "status": "400", "log": err})
		}

		buf, _ := reporteGenerado.WriteToBuffer()
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
								strValor2 := strings.TrimLeft(fmt.Sprintf("%v", recursos[i]["valor"]), "$")
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
						if len(unidades) > 1 {
							if unidades[0] == unidades[1] && (len(unidades)-1) == len(valores) {
								unidades = unidades[1:]
							}
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
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	excelArmonizacion := make([]map[string]interface{}, 0)
	nombre := c.Ctx.Input.Param(":nombre")
	consolidadoExcelEvaluacion := excelize.NewFile()
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if body["unidad_id"].(string) != "" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:6153355601c7a2365b2fb2a1,dependencia_id:"+body["unidad_id"].(string)+",nombre:"+nombre, &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planes)

			if len(planes) <= 0 {
				c.Abort("404")
			}

			trimestres := evaluacionhelper.GetPeriodosPlan(body["vigencia"].(string), planes[0]["_id"].(string))

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

			trimestreDelAnio := "-1"
			trimestresConContenido := make([]map[string]interface{}, 0)

			for indexTrimestre := 0; indexTrimestre <= len(trimestres)-1; indexTrimestre++ {
				auxEvaluacion := evaluacionhelper.GetEvaluacion(planes[0]["_id"].(string), trimestres, indexTrimestre)
				if fmt.Sprintf("%v", auxEvaluacion) != "[]" {
					evaluacion = auxEvaluacion
					trimestreDelAnio = trimestres[indexTrimestre]["codigo_trimestre"].(string)
					trimestresConContenido = append(trimestresConContenido, trimestres[indexTrimestre])
				}
			}

			trimestreVacio := map[string]interface{}{"actividad": 0.0, "acumulado": 0.0, "denominador": 0.0, "meta": 0.0, "numerador": 0.0, "periodo": 0.0, "numeradorAcumulado": 0.0, "denominadorAcumulado": 0.0, "brecha": 0.0, "cualitativo": map[string]interface{}{"reporte": "", "dificultades": ""}}
			for _, actividad := range evaluacion {
				for auxTrim := len(trimestresConContenido) - 1; auxTrim >= 0; auxTrim-- {
					seguimiento, err := seguimientohelper.GetSeguimiento(planes[0]["_id"].(string), actividad["numero"].(string), trimestresConContenido[auxTrim]["_id"].(string))
					if err == nil {
						switch auxTrim {
						case 0:
							actividad["trimestre1"].(map[string]interface{})["cualitativo"] = seguimiento["cualitativo"]
						case 1:
							actividad["trimestre2"].(map[string]interface{})["cualitativo"] = seguimiento["cualitativo"]
						case 2:
							actividad["trimestre3"].(map[string]interface{})["cualitativo"] = seguimiento["cualitativo"]
						case 3:
							actividad["trimestre4"].(map[string]interface{})["cualitativo"] = seguimiento["cualitativo"]
						}

					}
				}
			}

			switch trimestreDelAnio {
			case "4":
				for _, actividad := range evaluacion {
					if len(actividad["trimestre4"].(map[string]interface{})) == 0 {
						actividad["trimestre4"] = trimestreVacio
					}
					if len(actividad["trimestre3"].(map[string]interface{})) == 0 {
						actividad["trimestre3"] = trimestreVacio
					}
					if len(actividad["trimestre2"].(map[string]interface{})) == 0 {
						actividad["trimestre2"] = trimestreVacio
					}
					if len(actividad["trimestre1"].(map[string]interface{})) == 0 {
						actividad["trimestre1"] = trimestreVacio
					}
				}
			case "3":
				for _, actividad := range evaluacion {
					actividad["trimestre4"] = trimestreVacio
					if len(actividad["trimestre3"].(map[string]interface{})) == 0 {
						actividad["trimestre3"] = trimestreVacio
					}
					if len(actividad["trimestre2"].(map[string]interface{})) == 0 {
						actividad["trimestre2"] = trimestreVacio
					}
					if len(actividad["trimestre1"].(map[string]interface{})) == 0 {
						actividad["trimestre1"] = trimestreVacio
					}
				}
			case "2":
				for _, actividad := range evaluacion {
					actividad["trimestre4"] = trimestreVacio
					actividad["trimestre3"] = trimestreVacio
					if len(actividad["trimestre2"].(map[string]interface{})) == 0 {
						actividad["trimestre2"] = trimestreVacio
					}
					if len(actividad["trimestre1"].(map[string]interface{})) == 0 {
						actividad["trimestre1"] = trimestreVacio
					}
				}
			case "1":
				for _, actividad := range evaluacion {
					actividad["trimestre4"] = trimestreVacio
					actividad["trimestre3"] = trimestreVacio
					actividad["trimestre2"] = trimestreVacio
					if len(actividad["trimestre1"].(map[string]interface{})) == 0 {
						actividad["trimestre1"] = trimestreVacio
					}
				}
			case "-1":
				c.Abort("404")
			}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planes[0]["_id"].(string)+"&fields=nombre,_id,hijos,activo", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

				for i := 0; i < len(subgrupos); i++ {
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") {
						actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
						var arregloLineamietoPED []map[string]interface{}
						var arregloLineamietoPI []map[string]interface{}
						if len(actividades) == 1 {
							for index := range actividades {
								if val, ok := actividades[index]["index"].(float64); ok {
									actividades[index]["index"] = fmt.Sprintf("%v", int(val))
								}
							}
						} else {
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
						}

						reporteshelper.LimpiarDetalles()
						for j := 0; j < len(actividades); j++ {
							arregloLineamietoPED = nil
							arregloLineamietoPI = nil
							actividad := actividades[j]
							actividadName := actividad["dato"].(string)
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
							arregloLineamietoPED = reporteshelper.ArbolArmonizacionV2(armonizacionTercerNivel.(string))
							arregloLineamietoPI = reporteshelper.ArbolArmonizacionPIV2(armonizacionTercerNivelPI.(string))
							generalData := make(map[string]interface{})
							generalData["nombreActividad"] = actividadName
							generalData["numeroActividad"] = index
							generalData["datosArmonizacionPED"] = arregloLineamietoPED
							generalData["datosArmonizacionPI"] = arregloLineamietoPI
							generalData["datosComplementarios"] = datosArmonizacion
							excelArmonizacion = append(excelArmonizacion, generalData)

							b, _ := json.Marshal(excelArmonizacion)
							beego.Info(string(b))
						}
						break
					}
				}
			} else {
				panic(err)
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
			styleSubTitles, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
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
			styleContenidoCS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Font:      &excelize.Font{Bold: true},
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
			stylecontent, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			stylecontentC, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			styleLineamiento, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
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
			styleLineamientoSombra, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
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
			stylecontentCS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			stylecontentS, _ := consolidadoExcelEvaluacion.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{Horizontal: "justify", Vertical: "center", WrapText: true},
				Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})

			// Size
			consolidadoExcelEvaluacion.SetRowHeight(sheetName, 1, 12)
			consolidadoExcelEvaluacion.SetRowHeight(sheetName, 2, 27)
			consolidadoExcelEvaluacion.SetRowHeight(sheetName, 19, 31)
			consolidadoExcelEvaluacion.SetRowHeight(sheetName, 22, 27)

			consolidadoExcelEvaluacion.SetColWidth(sheetName, "A", "A", 3)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "B", "B", 19)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "C", "C", 13)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "D", "G", 35)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "E", "E", 16)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "H", "H", 4)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "I", "I", 8)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "J", "J", 13)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "K", "K", 42)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "L", "L", 16)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "M", "M", 21)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "N", "BG", 14)
			consolidadoExcelEvaluacion.SetColWidth(sheetName, "BJ", "BN", 3)
			// Merge
			consolidadoExcelEvaluacion.MergeCell(sheetName, "B4", "D4")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "B19", "E19")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "B21", "D21")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "E21", "G21")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "H21", "H22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "I21", "I22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "J21", "J22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "K21", "K22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "L21", "L22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "M21", "M22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "N21", "N22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "O21", "O22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "BL19", "BM19")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "BL21", "BL22")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "BM21", "BM22")

			consolidadoExcelEvaluacion.MergeCell(sheetName, "P21", "Z21")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "AA21", "AK21")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "AL21", "AV21")
			consolidadoExcelEvaluacion.MergeCell(sheetName, "AW21", "BG21")
			// Style
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B2", "T2", styleUnidad)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B4", "D4", styleTituloSB)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "E4", "E4", styleSombreadoSB)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B19", "B19", styleNegrilla)
			//  Estilos títulos
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B21", "BG21", styleTitulo)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "J22", "BG22", styleContenidoC)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B22", "G22", styleSubTitles)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "P22", "BG22", styleSubTitles)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z22", "Z22", styleContenidoCS)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AK22", "AK22", styleContenidoCS)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AV22", "AV22", styleContenidoCS)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BG22", "BG22", styleContenidoCS)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H22", "O22", styleTitulo)

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
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "B21", "Armonización PED")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "E21", "Armonización Plan Indicativo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "H21", "No.")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "I21", "Pond.")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "J21", "Periodo de ejecución")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "K21", "Actividad")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "L21", "Indicador asociado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "M21", "Fórmula del Indicador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "N21", "Meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "O21", "Tipo de Unidad")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "P21", "Trimestre I")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA21", "Trimestre II")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL21", "Trimestre III")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AW21", "Trimestre IV")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "B22", "Lineamiento")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "C22", "Meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "D22", "Estrategias")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "E22", "Ejes transformadores")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "F22", "Lineamientos de acción")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "G22", "Estrategias")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "P22", "Reporte de avance")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q22", "Dificultades")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "R22", "Numerador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "S22", "Denominador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "T22", "Indicador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "U22", "Numerador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "V22", "Denominador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "W22", "Indicador Acumulado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "X22", "Cumplimiento por meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y22", "Brecha")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z22", "Cumplimiento por actividad")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA22", "Reporte de avance")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB22", "Dificultades")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC22", "Numerador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD22", "Denominador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE22", "Indicador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF22", "Numerador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG22", "Denominador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH22", "Indicador Acumulado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI22", "Cumplimiento por meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ22", "Brecha")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK22", "Cumplimiento por actividad")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL22", "Reporte de avance")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM22", "Dificultades")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN22", "Numerador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO22", "Denominador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP22", "Indicador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ22", "Numerador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR22", "Denominador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS22", "Indicador Acumulado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AT22", "Cumplimiento por meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AU22", "Brecha")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AV22", "Cumplimiento por actividad")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AW22", "Reporte de avance")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX22", "Dificultades")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY22", "Numerador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AZ22", "Denominador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BA22", "Indicador del Periodo")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BB22", "Numerador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BC22", "Denominador Acumulador")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BD22", "Indicador Acumulado")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BE22", "Cumplimiento por meta")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BF22", "Brecha")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BG22", "Cumplimiento por actividad")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BL19", "Gráfica")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BL21", "No.")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BM21", "Cumplimiento")

			indice := 23
			indiceGraficos := 23
			// Agregar armonización
			for posArmonizacion := 0; posArmonizacion < len(excelArmonizacion); posArmonizacion++ {
				datosExcelArmonizacion := excelArmonizacion[posArmonizacion]
				armoPED := datosExcelArmonizacion["datosArmonizacionPED"].([]map[string]interface{})
				armoPI := datosExcelArmonizacion["datosArmonizacionPI"].([]map[string]interface{})
				datosComplementarios := datosExcelArmonizacion["datosComplementarios"].(map[string]interface{})
				indicadores := datosComplementarios["indicadores"].(map[string]interface{})
				numeroActividad := datosExcelArmonizacion["numeroActividad"]

				MaxRowsXActivity := reporteshelper.MinComMul_Armonization(armoPED, armoPI, len(indicadores))

				y_lin := indice
				h_lin := MaxRowsXActivity / len(armoPED)
				consolidadoExcelEvaluacion.SetRowHeight(sheetName, y_lin, 27*5)
				for _, lin := range armoPED {
					consolidadoExcelEvaluacion.MergeCell(sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1))
					consolidadoExcelEvaluacion.SetCellValue(sheetName, "B"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
					reporteshelper.SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "B"+fmt.Sprint(y_lin), "B"+fmt.Sprint(y_lin+h_lin-1), styleLineamiento, styleLineamientoSombra)
					y_met := y_lin
					h_met := h_lin / len(lin["meta"].([]map[string]interface{}))
					for _, met := range lin["meta"].([]map[string]interface{}) {
						consolidadoExcelEvaluacion.MergeCell(sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "C"+fmt.Sprint(y_met), met["nombreMeta"])
						reporteshelper.SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "C"+fmt.Sprint(y_met), "C"+fmt.Sprint(y_met+h_met-1), stylecontentC, stylecontentCS)
						y_est := y_met
						h_est := h_met / len(met["estrategias"].([]map[string]interface{}))
						for _, est := range met["estrategias"].([]map[string]interface{}) {
							consolidadoExcelEvaluacion.MergeCell(sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1))
							consolidadoExcelEvaluacion.SetCellValue(sheetName, "D"+fmt.Sprint(y_est), est["descripcionEstrategia"])
							if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
								reporteshelper.SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
							} else {
								reporteshelper.SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "D"+fmt.Sprint(y_est), "D"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
							}
							y_est += h_est
						}
						y_met += h_met
					}
					y_lin += h_lin
				}

				y_eje := indice
				h_eje := MaxRowsXActivity / len(armoPI)
				for _, eje := range armoPI {
					consolidadoExcelEvaluacion.MergeCell(sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1))
					consolidadoExcelEvaluacion.SetCellValue(sheetName, "E"+fmt.Sprint(y_eje), eje["nombreFactor"])
					reporteshelper.SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "E"+fmt.Sprint(y_eje), "E"+fmt.Sprint(y_eje+h_eje-1), stylecontentC, stylecontentCS)
					y_lin := y_eje
					h_lin := h_eje / len(eje["lineamientos"].([]map[string]interface{}))
					for _, lin := range eje["lineamientos"].([]map[string]interface{}) {
						consolidadoExcelEvaluacion.MergeCell(sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "F"+fmt.Sprint(y_lin), lin["nombreLineamiento"])
						reporteshelper.SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "F"+fmt.Sprint(y_lin), "F"+fmt.Sprint(y_lin+h_lin-1), stylecontentC, stylecontentCS)
						y_est := y_lin
						h_est := h_lin / len(lin["estrategias"].([]map[string]interface{}))
						for _, est := range lin["estrategias"].([]map[string]interface{}) {
							consolidadoExcelEvaluacion.MergeCell(sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1))
							consolidadoExcelEvaluacion.SetCellValue(sheetName, "G"+fmt.Sprint(y_est), est["descripcionEstrategia"])
							if (est["nombreEstrategia"].(string) == "No seleccionado") || strings.Contains(strings.ToLower(est["nombreEstrategia"].(string)), "no aplica") {
								reporteshelper.SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontentC, stylecontentCS)
							} else {
								reporteshelper.SombrearCeldas(consolidadoExcelEvaluacion, posArmonizacion, sheetName, "G"+fmt.Sprint(y_est), "G"+fmt.Sprint(y_est+h_est-1), stylecontent, stylecontentS)
							}
							y_est += h_est
						}
						y_lin += h_lin
					}
					y_eje += h_eje
				}

				for i, actividad := range evaluacion {
					if numeroActividad == actividad["numero"] {
						// Union de celdas
						consolidadoExcelEvaluacion.MergeCell(sheetName, "H"+fmt.Sprint(indice), "H"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "I"+fmt.Sprint(indice), "I"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "J"+fmt.Sprint(indice), "J"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "K"+fmt.Sprint(indice), "K"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "L"+fmt.Sprint(indice), "L"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "M"+fmt.Sprint(indice), "M"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "N"+fmt.Sprint(indice), "N"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "O"+fmt.Sprint(indice), "O"+fmt.Sprint(indice+MaxRowsXActivity-1))

						// Datos
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "H"+fmt.Sprint(indice), posArmonizacion+1)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "I"+fmt.Sprint(indice), actividad["ponderado"].(float64)/100)
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "J"+fmt.Sprint(indice), actividad["periodo"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "K"+fmt.Sprint(indice), datosExcelArmonizacion["nombreActividad"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "L"+fmt.Sprint(indice), actividad["indicador"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "M"+fmt.Sprint(indice), actividad["formula"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "N"+fmt.Sprint(indice), actividad["meta"].(float64))
						if actividad["unidad"] == "Porcentaje" {
							consolidadoExcelEvaluacion.SetCellValue(sheetName, "N"+fmt.Sprint(indice), actividad["meta"].(float64)/100)
						}
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "O"+fmt.Sprint(indice), actividad["unidad"])

						consolidadoExcelEvaluacion.MergeCell(sheetName, "P"+fmt.Sprint(indice), "P"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "Q"+fmt.Sprint(indice), "Q"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "R"+fmt.Sprint(indice), "R"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "S"+fmt.Sprint(indice), "S"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "T"+fmt.Sprint(indice), "T"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "U"+fmt.Sprint(indice), "U"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "V"+fmt.Sprint(indice), "V"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "W"+fmt.Sprint(indice), "W"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "X"+fmt.Sprint(indice), "X"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "Y"+fmt.Sprint(indice), "Y"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice+MaxRowsXActivity-1))

						consolidadoExcelEvaluacion.MergeCell(sheetName, "AA"+fmt.Sprint(indice), "AA"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AB"+fmt.Sprint(indice), "AB"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AC"+fmt.Sprint(indice), "AC"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AD"+fmt.Sprint(indice), "AD"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AE"+fmt.Sprint(indice), "AE"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AF"+fmt.Sprint(indice), "AF"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AG"+fmt.Sprint(indice), "AG"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AH"+fmt.Sprint(indice), "AH"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AI"+fmt.Sprint(indice), "AI"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AK"+fmt.Sprint(indice), "AK"+fmt.Sprint(indice+MaxRowsXActivity-1))

						consolidadoExcelEvaluacion.MergeCell(sheetName, "AL"+fmt.Sprint(indice), "AL"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AM"+fmt.Sprint(indice), "AM"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AN"+fmt.Sprint(indice), "AN"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AO"+fmt.Sprint(indice), "AO"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AQ"+fmt.Sprint(indice), "AQ"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AR"+fmt.Sprint(indice), "AR"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AT"+fmt.Sprint(indice), "AT"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AU"+fmt.Sprint(indice), "AU"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AV"+fmt.Sprint(indice), "AV"+fmt.Sprint(indice+MaxRowsXActivity-1))

						consolidadoExcelEvaluacion.MergeCell(sheetName, "AW"+fmt.Sprint(indice), "AW"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AX"+fmt.Sprint(indice), "AX"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AY"+fmt.Sprint(indice), "AY"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "AZ"+fmt.Sprint(indice), "AZ"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "BA"+fmt.Sprint(indice), "BA"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "BB"+fmt.Sprint(indice), "BB"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "BC"+fmt.Sprint(indice), "BC"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "BD"+fmt.Sprint(indice), "BD"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "BE"+fmt.Sprint(indice), "BE"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "BF"+fmt.Sprint(indice), "BF"+fmt.Sprint(indice+MaxRowsXActivity-1))
						consolidadoExcelEvaluacion.MergeCell(sheetName, "BG"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice+MaxRowsXActivity-1))

						// Trimestres
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "P"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["cualitativo"].(map[string]interface{})["reporte"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["cualitativo"].(map[string]interface{})["dificultades"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "R"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre1"].(map[string]interface{})["numerador"]))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "S"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre1"].(map[string]interface{})["denominador"]))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "T"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["periodo"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "U"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["numeradorAcumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "V"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["denominadorAcumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "W"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["acumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "X"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["meta"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["brecha"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z"+fmt.Sprint(indice), actividad["trimestre1"].(map[string]interface{})["actividad"].(float64))

						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["cualitativo"].(map[string]interface{})["reporte"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["cualitativo"].(map[string]interface{})["dificultades"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre2"].(map[string]interface{})["numerador"]))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre2"].(map[string]interface{})["denominador"]))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["periodo"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["numeradorAcumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["denominadorAcumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["acumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["meta"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["brecha"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK"+fmt.Sprint(indice), actividad["trimestre2"].(map[string]interface{})["actividad"].(float64))

						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["cualitativo"].(map[string]interface{})["reporte"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["cualitativo"].(map[string]interface{})["dificultades"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre3"].(map[string]interface{})["numerador"]))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre3"].(map[string]interface{})["denominador"]))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["periodo"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["numeradorAcumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["denominadorAcumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["acumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AT"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["meta"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AU"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["brecha"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AV"+fmt.Sprint(indice), actividad["trimestre3"].(map[string]interface{})["actividad"].(float64))

						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AW"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["cualitativo"].(map[string]interface{})["reporte"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["cualitativo"].(map[string]interface{})["dificultades"].(string))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre4"].(map[string]interface{})["numerador"]))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "AZ"+fmt.Sprint(indice), reporteshelper.Convert2Num(actividad["trimestre4"].(map[string]interface{})["denominador"]))
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BA"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["periodo"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BB"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["numeradorAcumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BC"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["denominadorAcumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BD"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["acumulado"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BE"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["meta"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BF"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["brecha"])
						consolidadoExcelEvaluacion.SetCellValue(sheetName, "BG"+fmt.Sprint(indice), actividad["trimestre4"].(map[string]interface{})["actividad"].(float64))

						// Estilos
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H"+fmt.Sprint(indice), "AY"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoC)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "I"+fmt.Sprint(indice), "I"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "K"+fmt.Sprint(indice), "K"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenido)

						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "P"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)

						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "R"+fmt.Sprint(indice), "S"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "U"+fmt.Sprint(indice), "V"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)

						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AC"+fmt.Sprint(indice), "AD"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AF"+fmt.Sprint(indice), "AG"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)

						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AN"+fmt.Sprint(indice), "AO"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AQ"+fmt.Sprint(indice), "AR"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)

						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AY"+fmt.Sprint(indice), "AZ"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BB"+fmt.Sprint(indice), "BC"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)

						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCPS)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AK"+fmt.Sprint(indice), "AK"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCPS)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AV"+fmt.Sprint(indice), "AV"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCPS)
						consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BG"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCPS)

						if actividad["unidad"] == "Porcentaje" {
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "N"+fmt.Sprint(indice), "N"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "T"+fmt.Sprint(indice), "T"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "W"+fmt.Sprint(indice), "W"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Y"+fmt.Sprint(indice), "Y"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AE"+fmt.Sprint(indice), "AE"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AH"+fmt.Sprint(indice), "AH"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AU"+fmt.Sprint(indice), "AU"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BA"+fmt.Sprint(indice), "BA"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BD"+fmt.Sprint(indice), "BD"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BF"+fmt.Sprint(indice), "BF"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCP)
						} else if actividad["unidad"] == "Tasa" {
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "N"+fmt.Sprint(indice), "N"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "T"+fmt.Sprint(indice), "T"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "W"+fmt.Sprint(indice), "W"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Y"+fmt.Sprint(indice), "Y"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AE"+fmt.Sprint(indice), "AE"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AH"+fmt.Sprint(indice), "AH"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AU"+fmt.Sprint(indice), "AU"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BA"+fmt.Sprint(indice), "BA"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BD"+fmt.Sprint(indice), "BD"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BF"+fmt.Sprint(indice), "BF"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCD)
						} else if actividad["unidad"] == "Unidad" {
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "N"+fmt.Sprint(indice), "N"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "T"+fmt.Sprint(indice), "T"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "W"+fmt.Sprint(indice), "W"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Y"+fmt.Sprint(indice), "Y"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AE"+fmt.Sprint(indice), "AE"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AH"+fmt.Sprint(indice), "AH"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AJ"+fmt.Sprint(indice), "AJ"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AP"+fmt.Sprint(indice), "AP"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AS"+fmt.Sprint(indice), "AS"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AU"+fmt.Sprint(indice), "AU"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BA"+fmt.Sprint(indice), "BA"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BD"+fmt.Sprint(indice), "BD"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
							consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BF"+fmt.Sprint(indice), "BF"+fmt.Sprint(indice+MaxRowsXActivity-1), styleContenidoCE)
						}

						// Unión de celdas por indicador
						if i > 0 {
							if actividad["numero"] == evaluacion[i-1]["numero"] {
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "H"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "I"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "J"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "K"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "Z"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "AK"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "AV"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "BG"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "BL"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.SetCellValue(sheetName, "BM"+fmt.Sprint(indice), nil)
								consolidadoExcelEvaluacion.MergeCell(sheetName, "P"+fmt.Sprint(indice-1), "P"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "Q"+fmt.Sprint(indice-1), "Q"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "AA"+fmt.Sprint(indice-1), "AA"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "AB"+fmt.Sprint(indice-1), "AB"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "AL"+fmt.Sprint(indice-1), "AL"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "AM"+fmt.Sprint(indice-1), "AM"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "AW"+fmt.Sprint(indice-1), "AW"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "AX"+fmt.Sprint(indice-1), "AX"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "B"+fmt.Sprint(indice-1), "B"+fmt.Sprint(indice+MaxRowsXActivity-1))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "C"+fmt.Sprint(indice-1), "C"+fmt.Sprint(indice+MaxRowsXActivity-1))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "D"+fmt.Sprint(indice-1), "D"+fmt.Sprint(indice+MaxRowsXActivity-1))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "E"+fmt.Sprint(indice-1), "E"+fmt.Sprint(indice+MaxRowsXActivity-1))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "F"+fmt.Sprint(indice-1), "F"+fmt.Sprint(indice+MaxRowsXActivity-1))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "G"+fmt.Sprint(indice-1), "G"+fmt.Sprint(indice+MaxRowsXActivity-1))
								consolidadoExcelEvaluacion.SetCellStyle(sheetName, "B"+fmt.Sprint(indice-1), "B"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
								consolidadoExcelEvaluacion.SetCellStyle(sheetName, "C"+fmt.Sprint(indice-1), "C"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
								consolidadoExcelEvaluacion.SetCellStyle(sheetName, "D"+fmt.Sprint(indice-1), "D"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
								consolidadoExcelEvaluacion.SetCellStyle(sheetName, "E"+fmt.Sprint(indice-1), "E"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
								consolidadoExcelEvaluacion.SetCellStyle(sheetName, "F"+fmt.Sprint(indice-1), "F"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)
								consolidadoExcelEvaluacion.SetCellStyle(sheetName, "G"+fmt.Sprint(indice-1), "G"+fmt.Sprint(indice+MaxRowsXActivity-1), styleLineamiento)

								consolidadoExcelEvaluacion.MergeCell(sheetName, "H"+fmt.Sprint(indice-1), "H"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "I"+fmt.Sprint(indice-1), "I"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "J"+fmt.Sprint(indice-1), "J"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "K"+fmt.Sprint(indice-1), "K"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "Z"+fmt.Sprint(indice-1), "Z"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "AK"+fmt.Sprint(indice-1), "AK"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "AV"+fmt.Sprint(indice-1), "AV"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.MergeCell(sheetName, "BG"+fmt.Sprint(indice-1), "BG"+fmt.Sprint(indice))
							} else {
								// Gaficos
								consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BL"+fmt.Sprint(indiceGraficos), "=H"+fmt.Sprint(indice))
								consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BM"+fmt.Sprint(indiceGraficos), "=IF(E4=\"Trimestre I\",Z"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AK"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AV"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",BG"+fmt.Sprint(indice)+"))))")
								indiceGraficos++
							}
						} else if i == 0 {
							// Gaficos
							consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BL"+fmt.Sprint(indiceGraficos), "=H"+fmt.Sprint(indice))
							consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BM"+fmt.Sprint(indiceGraficos), "=IF(E4=\"Trimestre I\",Z"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AK"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AV"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",BG"+fmt.Sprint(indice)+"))))")
							indiceGraficos++
						}
						indice += MaxRowsXActivity
					}
				}
			}

			consolidadoExcelEvaluacion.MergeCell(sheetName, "H"+fmt.Sprint(indice), "O"+fmt.Sprint(indice))
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "H"+fmt.Sprint(indice), "O"+fmt.Sprint(indice), styleTituloSB)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "P"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice), styleContenidoC)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "Z"+fmt.Sprint(indice), "Z"+fmt.Sprint(indice), styleContenidoCPSR)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AK"+fmt.Sprint(indice), "AK"+fmt.Sprint(indice), styleContenidoCPSR)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "AV"+fmt.Sprint(indice), "AV"+fmt.Sprint(indice), styleContenidoCPSR)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BG"+fmt.Sprint(indice), "BG"+fmt.Sprint(indice), styleContenidoCPSR)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BL19", "BL"+fmt.Sprint(indice+1), styleContenidoCI)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BM19", "BN"+fmt.Sprint(indice+1), styleContenidoCIP)
			consolidadoExcelEvaluacion.SetCellStyle(sheetName, "BM21", "BM22", styleContenidoCI)
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "H"+fmt.Sprint(indice), "Avance General del Plan de Acción")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "E4", "Trimestre I")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "P"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Q"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "R"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "S"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "T"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "U"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "V"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "W"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "X"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "Y"+fmt.Sprint(indice), "-")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AA"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AB"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AC"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AD"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AE"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AF"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AG"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AH"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AI"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AJ"+fmt.Sprint(indice), "-")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AL"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AM"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AN"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AO"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AP"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AQ"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AR"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AS"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AT"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AU"+fmt.Sprint(indice), "-")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AW"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AX"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AY"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "AZ"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BA"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BB"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BC"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BD"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BE"+fmt.Sprint(indice), "-")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BF"+fmt.Sprint(indice), "-")

			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BL"+fmt.Sprint(indice), "General")

			filaAnt := fmt.Sprint(indice - 1)
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "Z"+fmt.Sprint(indice), "=SUMPRODUCT(I23:I"+filaAnt+",Z23:Z"+filaAnt+")")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AK"+fmt.Sprint(indice), "=SUMPRODUCT(I23:I"+filaAnt+",AK23:AK"+filaAnt+")")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "AV"+fmt.Sprint(indice), "=SUMPRODUCT(I23:I"+filaAnt+",AV23:AV"+filaAnt+")")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BG"+fmt.Sprint(indice), "=SUMPRODUCT(I23:I"+filaAnt+",BG23:BG"+filaAnt+")")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BM"+fmt.Sprint(indice), "=IF(E4=\"Trimestre I\",Z"+fmt.Sprint(indice)+",IF(E4=\"Trimestre II\",AK"+fmt.Sprint(indice)+",IF(E4=\"Trimestre III\",AV"+fmt.Sprint(indice)+",IF(E4=\"Trimestre IV\",BG"+fmt.Sprint(indice)+"))))")
			consolidadoExcelEvaluacion.SetCellFormula(sheetName, "BN"+fmt.Sprint(indice), "=100%-BM"+fmt.Sprint(indice))
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BM"+fmt.Sprint(indice+1), "Avance")
			consolidadoExcelEvaluacion.SetCellValue(sheetName, "BN"+fmt.Sprint(indice+1), "Restante")

			consolidadoExcelEvaluacion.AddChart(sheetName, "B5", &excelize.Chart{
				Type: excelize.Pie,
				Series: []excelize.ChartSeries{
					{
						Name:       "",
						Categories: sheetName + "!$BM$" + fmt.Sprint(indice+1) + ":$BN$" + fmt.Sprint(indice+1),
						Values:     sheetName + "!$BM$" + fmt.Sprint(indice) + ":$BN$" + fmt.Sprint(indice),
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
						Categories: sheetName + "!$BL$23:$BL$" + fmt.Sprint(indiceGraficos-1),
						Values:     sheetName + "!$BM$23:$BM$" + fmt.Sprint(indiceGraficos-1),
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
