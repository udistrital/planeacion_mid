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

			anioPlan := periodo[0]["Year"].(float64)
			var esReporteAntiguo bool = true

			if anioPlan > 2024 { //? Vigencias 2025 en adelante son reportes con nueva estructura
				esReporteAntiguo = false
			}

			reporteGenerado, errorReporte := reporteshelper.ConstruirExcelPlanAccionUnidad(esReporteAntiguo, datosReporte)

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
	var resPeriodo map[string]interface{}
	var periodo []map[string]interface{}
	var datosReporte map[string]interface{}

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

		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+`/periodo?query=Id:`+body["vigencia"].(string), &resPeriodo); err == nil {
			helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
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

			datosReporte = map[string]interface{}{
				"periodo":          periodo,
				"planesFilter":     planesFilter,
				"arregloPlanAnual": arregloPlanAnual,
			}
		}

		anioPlan := periodo[0]["Year"].(float64)
		var esReporteAntiguo bool = true

		if anioPlan > 2024 { //? Vigencias 2025 en adelante son reportes con nueva estructura
			esReporteAntiguo = false
		}

		reporteGenerado, errorReporte := reporteshelper.ConstruirExcelPlanAccionGeneral(esReporteAntiguo, datosReporte)

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
	var datosReporte map[string]interface{}
	excelArmonizacion := make([]map[string]interface{}, 0)
	nombre := c.Ctx.Input.Param(":nombre")
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
						}
						break
					}
				}
			} else {
				panic(err)
			}

			datosReporte = map[string]interface{}{
				"unidadNombre":      unidadNombre,
				"periodo":           periodo,
				"evaluacion":        evaluacion,
				"excelArmonizacion": excelArmonizacion,
			}

			anioPlan := periodo[0]["Year"].(float64)
			var esReporteAntiguo bool = true

			if anioPlan > 2024 { //? Vigencias 2025 en adelante son reportes con nueva estructura
				esReporteAntiguo = false
			}

			reporteGenerado, errorReporte := reporteshelper.ConstruirExcelPlanAccionEvaluacion(esReporteAntiguo, datosReporte)

			if errorReporte != nil {
				panic(map[string]interface{}{"funcion": "PlanAccionEvaluacion", "err": "Error en la generación del excel", "status": "400", "log": err})
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
