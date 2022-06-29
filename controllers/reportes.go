package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/leekchan/accounting"
	"github.com/udistrital/planeacion_mid/helpers"
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
	c.Mapping("PlanAccionAnual", c.PlanAccionAnual)
	c.Mapping("PlanAccionAnualGeneral", c.PlanAccionAnualGeneral)
	c.Mapping("Necesidades", c.Necesidades)
}

func CreateExcel(f *excelize.File, dir string) {
	if err := f.Save(); err != nil {
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
	var nombreUnidadVer string

	// excel
	var consolidadoExcel *excelize.File
	consolidadoExcel = excelize.NewFile()
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
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
				c.Abort("400")
			}
		}
		contadorDesagregado := 3
		stylehead, _ := consolidadoExcel.NewStyle(`{
			"alignment":{"horizontal":"center","vertical":"center","wrap_text":true}, 
			"font":{"bold":true,"color":"#FFFFFF"},
			"fill":{"type":"pattern","pattern":1,"color":["#CC0000"]},
			"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
		}`)
		styletitles, _ := consolidadoExcel.NewStyle(`{
			"alignment":{"horizontal":"center","vertical":"center","wrap_text":true}, 
			"font":{"bold":true},
			"fill":{"type":"pattern","pattern":1,"color":["#F2F2F2"]},
			"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
		}`)
		stylecontent, _ := consolidadoExcel.NewStyle(`{
			"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
			"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
		}`)
		for h := 0; h < len(data_identi); h++ {

			datosArreglo := data_identi[h]
			nombreUnidadVerIn := datosArreglo["unidad"].(string)
			if h == 0 {
				nombreUnidadVer = nombreUnidadVerIn
			}

			if nombreUnidadVerIn == nombreUnidadVer {
				nombreUnidadVer = datosArreglo["unidad"].(string)
				nombreHoja := fmt.Sprint(nombreUnidadVer)
				sheetName := nombreHoja
				index := consolidadoExcel.NewSheet(sheetName)
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
				nombreHoja := fmt.Sprint(nombreUnidadVer)
				sheetName := nombreHoja
				index := consolidadoExcel.NewSheet(sheetName)
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
		dataSend := make(map[string]interface{})

		buf, _ := consolidadoExcel.WriteToBuffer()
		strings.NewReader(buf.String())

		encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))

		dataSend["generalData"] = data_identi
		dataSend["excelB64"] = encoded

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}

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
	// var respuestaGeneral map[string]interface{}
	// var planesFilterGeneral []map[string]interface{}
	var respuesta map[string]interface{}
	var planesFilter []map[string]interface{}
	var res map[string]interface{}
	// var resPresupuesto map[string]interface{}
	var resArmo map[string]interface{}
	// var resEstrategia map[string]interface{}
	// var resMeta map[string]interface{}
	// var resLineamiento map[string]interface{}
	// var resPlan map[string]interface{}
	var respuestaUnidad []map[string]interface{}
	var hijosArmo []map[string]interface{}
	// var estrategiaData []map[string]interface{}
	// var metaData []map[string]interface{}
	// var LineamientoData []map[string]interface{}
	// var planData []map[string]interface{}
	var subgrupos []map[string]interface{}
	var plan_id string
	var actividadName string
	var arregloPlanAnual []map[string]interface{}
	// var arregloEstrategia []map[string]interface{}
	// var arregloMetaEst []map[string]interface{}
	// var identificacion map[string]interface{}
	// var datoPresupuesto map[string]interface{}
	// var identificacionres []map[string]interface{}
	// var data_identi []map[string]interface{}
	// var unidadId string
	// var nombrePlanDesarrollo string
	var nombreUnidad string
	// var unidadNombre string

	consolidadoExcelPlanAnual := excelize.NewFile()

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if body["unidad_id"].(string) == "" {

	} else if body["unidad_id"].(string) != "" {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",dependencia_id:"+body["unidad_id"].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
			for planes := 0; planes < len(planesFilter); planes++ {
				planesFilterData := planesFilter[planes]
				plan_id = planesFilterData["_id"].(string)

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id, &res); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
					for i := 0; i < len(subgrupos); i++ {
						if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
							actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
							var arregloLineamieto []map[string]interface{}
							var arregloLineamietoPI []map[string]interface{}

							for j := 0; j < len(actividades); j++ {
								arregloLineamieto = nil
								arregloLineamietoPI = nil
								actividad := actividades[j]
								actividadName = actividad["dato"].(string)
								index := fmt.Sprint(actividad["index"])
								datosArmonizacion := make(map[string]interface{})
								titulosArmonizacion := make(map[string]interface{})

								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &resArmo); err == nil {
									helpers.LimpiezaRespuestaRefactor(resArmo, &hijosArmo)
									reporteshelper.Limpia()
									tree := reporteshelper.BuildTreeFa(hijosArmo, index)

									treeDatos := tree[0]
									treeDatas := tree[1]
									treeArmo := tree[2]
									armonizacionTercer := treeArmo[0]

									armonizacionTercerNivel := armonizacionTercer["armo"].(map[string]interface{})["armonizacionPED"]
									armonizacionTercerNivelPI := armonizacionTercer["armo"].(map[string]interface{})["armonizacionPI"]

									for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
										treeDato := treeDatos[datoGeneral]
										treeData := treeDatas[0]
										if treeDato["sub"] == "" {
											if strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "ponderación") || strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "ponderacion") && strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "actividad") {
												datosArmonizacion["Ponderación de la actividad"] = treeData[fmt.Sprint(treeDato["id"])]
											} else if strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "período") || strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "periodo") && strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "ejecucion") || strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "ejecución") {
												datosArmonizacion["Periodo de ejecución"] = treeData[fmt.Sprint(treeDato["id"])]
											} else if strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "general") {
												datosArmonizacion["Actividad general"] = treeData[fmt.Sprint(treeDato["id"])]
											} else if strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "tarea") {
												datosArmonizacion["Tareas"] = treeData[fmt.Sprint(treeDato["id"])]
											} else {
												datosArmonizacion[treeDato["nombre"].(string)] = treeData[fmt.Sprint(treeDato["id"])]
											}
										}
									}
									treeIndicador := treeDatos[4]

									subIndicador := treeIndicador["sub"].([]map[string]interface{})
									for ind := 0; ind < len(subIndicador); ind++ {
										subIndicadorRes := subIndicador[ind]
										treeData := treeDatas[0]
										dataIndicador := make(map[string]interface{})
										auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
										for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
											dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[fmt.Sprint(auxSubIndicador[subInd]["id"])]
										}
										titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
									}
									datosArmonizacion["indicadores"] = titulosArmonizacion

									arregloLineamieto = reporteshelper.ArbolArmonizacion(armonizacionTercerNivel.(string))
									arregloLineamietoPI = reporteshelper.ArbolArmonizacionPI(armonizacionTercerNivelPI)
								} else {
									c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
									c.Abort("400")
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
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					c.Abort("400")
				}

				contadorLineamiento := 4
				// contadorMeta := 4
				// contadorEstrategia := 4

				contadorFactor := 4
				//contadorLineamientoPI := 4
				//contadorEstrategiaPI := 4
				contadorDataGeneral := 4
				unidadNombre := arregloPlanAnual[0]["nombreUnidad"]
				nombreHoja := fmt.Sprint(nombreUnidad)
				sheetName := nombreHoja
				indexPlan := consolidadoExcelPlanAnual.NewSheet(sheetName)

				stylehead, _ := consolidadoExcelPlanAnual.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"font":{"bold":true,"color":"#FFFFFF"},
					"fill":{"type":"pattern","pattern":1,"color":["#CC0000"]},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)
				styletitles, _ := consolidadoExcelPlanAnual.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"font":{"bold":true},
					"fill":{"type":"pattern","pattern":1,"color":["#F2F2F2"]},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)
				stylecontent, _ := consolidadoExcelPlanAnual.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)

				consolidadoExcelPlanAnual.MergeCell(sheetName, "A1", "N1")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "A2", "C2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "D2", "F2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "G2", "G3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "H2", "H3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "I2", "I3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "J2", "J3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "K2", "K3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "L2", "N2")
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 1, 20)
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 2, 20)
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 3, 20)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "A", "C", 70)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "D", "F", 70)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "N", 50)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "I", 20)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "J", "K", 80)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A1", "K1", stylehead)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A2", "N2", styletitles)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A3", "N3", styletitles)

				tituloExcel := fmt.Sprint("Plan de acción 2022 ", unidadNombre)

				// encabezado excel
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "A1", tituloExcel)
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "A2", "Armonización PED")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "A3", "Lineamiento")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "B3", "Meta")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "C3", "Estrategias")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "D2", "Armonización Plan Indicativo")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "D3", "Ejes transformadores")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "E3", "Lineaminetos de acción")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "F3", "Estrategias")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "G3", "N°.")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "H3", "Ponderación de la actividad")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "I3", "Periodo de ejecución")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "J3", "Actividad")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "K3", "Tareas")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Indicador")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "L3", "Nombre")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M3", "Fórmula")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "N3", "Meta")

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
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "A"+fmt.Sprint(contadorLineamiento), auxLineamiento)
						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorLineamiento), "N"+fmt.Sprint(contadorLineamiento), stylecontent)
						consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorLineamiento, 70)

						metas := datosArmo["meta"]

						contadorMetas := contadorLineamiento

						for j := 0; j < len(metas.([]map[string]interface{})); j++ {
							auxMeta := metas.([]map[string]interface{})[j]
							contadorMetaGeneralIn = contadorLineamiento

							consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorMetas), auxMeta["nombreMeta"])
							consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorMetas), "B"+fmt.Sprint(contadorMetas), stylecontent)

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

								consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorEstrategias), auxEstrategia["descripcionEstrategia"])
								consolidadoExcelPlanAnual.SetCellStyle(sheetName, "C"+fmt.Sprint(contadorEstrategias), "C"+fmt.Sprint(contadorEstrategias), stylecontent)

								if k == len(estrategias)-1 {
									contadorEstrategiaPEDOut = contadorMetas
								} else {
									contadorEstrategias = contadorEstrategias + 1
								}
							}

							contadorEstrategias = contadorMetas

							consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorMetaGeneralIn), "B"+fmt.Sprint(contadorMetaGeneralOut))

						}
						contadorMetas = contadorLineamiento
						contadorLineamientoGeneralOut = contadorMetaGeneralOut

						consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorLineamientoGeneralIn), "A"+fmt.Sprint(contadorLineamientoGeneralOut))

						contadorLineamiento = contadorLineamientoGeneralOut + 1

					}

					for i := 0; i < len(armoPI); i++ {
						datosArmo := armoPI[i]
						auxFactor := datosArmo["nombreFactor"]
						contadorFactorGeneralIn = contadorFactor

						// cuerpo del excel
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorFactor), auxFactor)
						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "D"+fmt.Sprint(contadorFactor), "N"+fmt.Sprint(contadorFactor), stylecontent)
						consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorFactor, 70)

						lineamientos := datosArmo["lineamientos"]

						contadorLineamientos := contadorFactor

						for j := 0; j < len(lineamientos.([]map[string]interface{})); j++ {
							auxLineamiento := lineamientos.([]map[string]interface{})[j]
							contadorLineamientoPIIn = contadorFactor

							consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorLineamientos), auxLineamiento["nombreLineamiento"])
							consolidadoExcelPlanAnual.SetCellStyle(sheetName, "E"+fmt.Sprint(contadorLineamientos), "E"+fmt.Sprint(contadorLineamientos), stylecontent)

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

								consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorEstrategias), auxEstrategia["descripcionEstrategia"])
								consolidadoExcelPlanAnual.SetCellStyle(sheetName, "F"+fmt.Sprint(contadorEstrategias), "F"+fmt.Sprint(contadorEstrategias), stylecontent)

								if k == len(estrategiasPI)-1 {
									contadorEstrategiaPIOut = contadorLineamientos
								} else {
									contadorEstrategias = contadorEstrategias + 1
								}
							}

							contadorEstrategias = contadorLineamientos

							consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorLineamientoPIIn), "E"+fmt.Sprint(contadorLineamientoPIOut))

						}
						contadorLineamientos = contadorFactor
						contadorFactorGeneralOut = contadorLineamientoPIOut

						consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorFactorGeneralIn), "D"+fmt.Sprint(contadorFactorGeneralOut))

						contadorFactor = contadorFactorGeneralOut + 1

					}

					consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorDataGeneral), datosExcelPlan["numeroActividad"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Ponderación de la actividad"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Periodo de ejecución"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Actividad general"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Tareas"])

					if contadorLineamientoGeneralOut > contadorFactorGeneralOut {
						contadorFactorGeneralOut = contadorLineamientoGeneralOut
						contadorFactor = contadorFactorGeneralOut + 1
					} else if contadorLineamientoGeneralOut < contadorFactorGeneralOut {
						contadorLineamientoGeneralOut = contadorFactorGeneralOut
						contadorLineamiento = contadorLineamientoGeneralOut + 1
					}

					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "G"+fmt.Sprint(contadorLineamiento), "N"+fmt.Sprint(contadorLineamiento), stylecontent)

					consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorLineamientoGeneralIn), "A"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorMetaGeneralIn), "B"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorEstrategiaPEDOut), "C"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorFactorGeneralIn), "D"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorLineamientoPIOut), "E"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorEstrategiaPIOut), "F"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorDataGeneral), "G"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorDataGeneral), "H"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorDataGeneral), "I"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorLineamientoGeneralOut))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorDataGeneral), "K"+fmt.Sprint(contadorLineamientoGeneralOut))

					indicadores := datosComplementarios["indicadores"].(map[string]interface{})
					contadorIndicadores := contadorDataGeneral
					for id, indicador := range indicadores {
						_ = id
						auxIndicador := indicador
						var nombreIndicador string
						var formula string
						var meta string

						for key, element := range auxIndicador.(map[string]interface{}) {
							if strings.Contains(strings.ToLower(key), "nombre") {
								nombreIndicador = element.(string)
							}
							if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
								formula = element.(string)
							}
							if strings.Contains(strings.ToLower(key), "meta") {
								meta = element.(string)
							}

						}

						consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorIndicadores), nombreIndicador)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorIndicadores), formula)
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorIndicadores), meta)

						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "L"+fmt.Sprint(contadorIndicadores), "L"+fmt.Sprint(contadorIndicadores), stylecontent)
						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "M"+fmt.Sprint(contadorIndicadores), "M"+fmt.Sprint(contadorIndicadores), stylecontent)
						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "N"+fmt.Sprint(contadorIndicadores), "N"+fmt.Sprint(contadorIndicadores), stylecontent)

						contadorIndicadores = contadorIndicadores + 1

					}

					contadorIndicadores--

					consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorLineamientoGeneralOut), "A"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorMetaGeneralOut), "B"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorEstrategiaPEDOut), "C"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorFactorGeneralOut), "D"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorLineamientoPIOut), "E"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorEstrategiaPIOut), "F"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorDataGeneral), "G"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorDataGeneral), "H"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorDataGeneral), "I"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorIndicadores))
					consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorDataGeneral), "K"+fmt.Sprint(contadorIndicadores))
					contadorDataGeneral = contadorIndicadores + 1
					contadorLineamiento = contadorIndicadores + 1
					contadorFactor = contadorIndicadores + 1

					consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)

				}

				consolidadoExcelPlanAnual = reporteshelper.TablaIdentificaciones(consolidadoExcelPlanAnual, plan_id)

			}

			//consolidadoExcelPlanAnual.SaveAs("plan_anual.xlsx")

			buf, _ := consolidadoExcelPlanAnual.WriteToBuffer()
			strings.NewReader(buf.String())

			encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))

			dataSend := make(map[string]interface{})

			dataSend["generalData"] = arregloPlanAnual
			dataSend["excelB64"] = encoded

			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}

		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
	}

	c.ServeJSON()
}

// PlanAccionAnualGeneral ...
// @Title PlanAccionAnualGeneral
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Reportes
// @Failure 403 :plan_id is empty
// @router /plan_anual_general [post]
func (c *ReportesController) PlanAccionAnualGeneral() {
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var planesFilter []map[string]interface{}
	var res map[string]interface{}
	var resArmo map[string]interface{}
	var respuestaUnidad []map[string]interface{}
	var respuestaEstado map[string]interface{}
	var respuestaTipoPlan map[string]interface{}
	var estado map[string]interface{}
	var tipoPlan map[string]interface{}
	var hijosArmo []map[string]interface{}
	var subgrupos []map[string]interface{}
	var plan_id string
	var actividadName string
	var arregloPlanAnual []map[string]interface{}
	var arregloInfoReportes []map[string]interface{}

	var nombreUnidad string
	contadorGeneral := 0

	consolidadoExcelPlanAnual := excelize.NewFile()

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)

		for planes := 0; planes < len(planesFilter); planes++ {
			planesFilterData := planesFilter[planes]
			plan_id = planesFilterData["_id"].(string)
			infoReporte := make(map[string]interface{})

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id, &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
				for i := 0; i < len(subgrupos); i++ {
					if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
						actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
						var arregloLineamieto []map[string]interface{}
						var arregloLineamietoPI []map[string]interface{}

						for j := 0; j < len(actividades); j++ {
							arregloLineamieto = nil
							arregloLineamietoPI = nil
							actividad := actividades[j]
							actividadName = actividad["dato"].(string)
							index := fmt.Sprint(actividad["index"])
							datosArmonizacion := make(map[string]interface{})
							titulosArmonizacion := make(map[string]interface{})

							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &resArmo); err == nil {
								helpers.LimpiezaRespuestaRefactor(resArmo, &hijosArmo)
								reporteshelper.Limpia()
								tree := reporteshelper.BuildTreeFa(hijosArmo, index)

								treeDatos := tree[0]
								treeDatas := tree[1]
								treeArmo := tree[2]
								armonizacionTercer := treeArmo[0]
								armonizacionTercerNivel := armonizacionTercer["armo"].(map[string]interface{})["armonizacionPED"]
								armonizacionTercerNivelPI := armonizacionTercer["armo"].(map[string]interface{})["armonizacionPI"]

								for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
									treeDato := treeDatos[datoGeneral]
									treeData := treeDatas[0]
									if treeDato["sub"] == "" {
										if strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "ponderación") || strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "ponderacion") && strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "actividad") {
											datosArmonizacion["Ponderación de la actividad"] = treeData[fmt.Sprint(treeDato["id"])]
										} else if strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "período") || strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "periodo") && strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "ejecucion") || strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "ejecución") {
											datosArmonizacion["Periodo de ejecución"] = treeData[fmt.Sprint(treeDato["id"])]
										} else if strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "general") {
											datosArmonizacion["Actividad general"] = treeData[fmt.Sprint(treeDato["id"])]
										} else if strings.Contains(strings.ToLower(treeDato["nombre"].(string)), "tarea") {
											datosArmonizacion["Tareas"] = treeData[fmt.Sprint(treeDato["id"])]
										} else {
											datosArmonizacion[treeDato["nombre"].(string)] = treeData[fmt.Sprint(treeDato["id"])]
										}
									}
								}
								treeIndicador := treeDatos[4]

								subIndicador := treeIndicador["sub"].([]map[string]interface{})
								for ind := 0; ind < len(subIndicador); ind++ {
									subIndicadorRes := subIndicador[ind]
									treeData := treeDatas[0]
									dataIndicador := make(map[string]interface{})
									auxSubIndicador := subIndicadorRes["sub"].([]map[string]interface{})
									for subInd := 0; subInd < len(auxSubIndicador); subInd++ {
										dataIndicador[auxSubIndicador[subInd]["nombre"].(string)] = treeData[fmt.Sprint(auxSubIndicador[subInd]["id"])]
									}
									titulosArmonizacion[subIndicadorRes["nombre"].(string)] = dataIndicador
								}
								datosArmonizacion["indicadores"] = titulosArmonizacion

								arregloLineamieto = reporteshelper.ArbolArmonizacion(armonizacionTercerNivel.(string))
								arregloLineamietoPI = reporteshelper.ArbolArmonizacionPI(armonizacionTercerNivelPI)
							} else {
								c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
								c.Abort("400")
							}

							generalData := make(map[string]interface{})

							if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+planesFilter[planes]["dependencia_id"].(string), &respuestaUnidad); err == nil {
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
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
				c.Abort("400")
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

			contadorLineamiento := contadorGeneral + 4
			contadorFactor := contadorGeneral + 4
			contadorDataGeneral := contadorGeneral + 4

			unidadNombre := arregloPlanAnual[0]["nombreUnidad"]
			sheetName := "REPORTE GENERAL"
			indexPlan := consolidadoExcelPlanAnual.NewSheet(sheetName)

			stylehead, _ := consolidadoExcelPlanAnual.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"font":{"bold":true,"color":"#FFFFFF"},
					"fill":{"type":"pattern","pattern":1,"color":["#CC0000"]},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)
			styletitles, _ := consolidadoExcelPlanAnual.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"font":{"bold":true},
					"fill":{"type":"pattern","pattern":1,"color":["#F2F2F2"]},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)
			stylecontent, _ := consolidadoExcelPlanAnual.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)

			consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorGeneral+1), "N"+fmt.Sprint(contadorGeneral+1))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorGeneral+2), "C"+fmt.Sprint(contadorGeneral+2))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorGeneral+2), "F"+fmt.Sprint(contadorGeneral+2))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorGeneral+2), "G"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorGeneral+2), "H"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorGeneral+2), "I"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorGeneral+2), "J"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorGeneral+2), "K"+fmt.Sprint(contadorGeneral+3))
			consolidadoExcelPlanAnual.MergeCell(sheetName, "L"+fmt.Sprint(contadorGeneral+2), "N"+fmt.Sprint(contadorGeneral+2))
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+1, 20)
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+2, 20)
			consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorGeneral+3, 20)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "A", "C", 70)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "D", "F", 70)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "N", 50)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "H", "I", 20)
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "J", "K", 80)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorGeneral+1), "K"+fmt.Sprint(contadorGeneral+1), stylehead)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorGeneral+2), "N"+fmt.Sprint(contadorGeneral+2), styletitles)
			consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorGeneral+3), "N"+fmt.Sprint(contadorGeneral+3), styletitles)

			tituloExcel := fmt.Sprint("Plan de acción 2022 ", unidadNombre)

			// encabezado excel
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "A"+fmt.Sprint(contadorGeneral+1), tituloExcel)
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "A"+fmt.Sprint(contadorGeneral+2), "Armonización PED")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "A"+fmt.Sprint(contadorGeneral+3), "Lineamiento")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorGeneral+3), "Meta")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorGeneral+3), "Estrategias")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorGeneral+2), "Armonización Plan Indicativo")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorGeneral+3), "Factores")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorGeneral+3), "Lineaminetos de acción")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorGeneral+3), "Estrategias")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorGeneral+3), "N°.")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorGeneral+3), "Ponderación de la actividad")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorGeneral+3), "Periodo de ejecución")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorGeneral+3), "Actividad")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorGeneral+3), "Tareas")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+2), "Indicador")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+3), "Nombre")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorGeneral+3), "Fórmula")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorGeneral+3), "Meta")

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
				//fmt.Println(contadorLineamientoGeneralIn + contadorFactorGeneralOut + contadorMetaGeneralIn + contadorMetaGeneralOut + contadorEstrategiaPEDIn + contadorFactorGeneralIn + contadorEstrategiaPEDOut + contadorLineamientoPIIn + contadorLineamientoPIOut + contadorEstrategiaPIIn + contadorEstrategiaPIOut)

				for i := 0; i < len(armoPED); i++ {
					datosArmo := armoPED[i]
					auxLineamiento := datosArmo["nombreLineamiento"]
					contadorLineamientoGeneralIn = contadorLineamiento

					// cuerpo del excel
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A"+fmt.Sprint(contadorLineamiento), auxLineamiento)
					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorLineamiento), "N"+fmt.Sprint(contadorLineamiento), stylecontent)
					consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorLineamiento, 70)

					metas := datosArmo["meta"]

					contadorMetas := contadorLineamiento

					for j := 0; j < len(metas.([]map[string]interface{})); j++ {
						auxMeta := metas.([]map[string]interface{})[j]
						contadorMetaGeneralIn = contadorLineamiento

						consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorMetas), auxMeta["nombreMeta"])
						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "B"+fmt.Sprint(contadorMetas), "B"+fmt.Sprint(contadorMetas), stylecontent)

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

							consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorEstrategias), auxEstrategia["descripcionEstrategia"])
							consolidadoExcelPlanAnual.SetCellStyle(sheetName, "C"+fmt.Sprint(contadorEstrategias), "C"+fmt.Sprint(contadorEstrategias), stylecontent)

							if k == len(estrategias)-1 {
								contadorEstrategiaPEDOut = contadorMetas
							} else {
								contadorEstrategias = contadorEstrategias + 1
							}
						}

						contadorEstrategias = contadorMetas

						consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorMetaGeneralIn), "B"+fmt.Sprint(contadorMetaGeneralOut))

					}
					contadorMetas = contadorLineamiento
					contadorLineamientoGeneralOut = contadorMetaGeneralOut

					consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorLineamientoGeneralIn), "A"+fmt.Sprint(contadorLineamientoGeneralOut))

					contadorLineamiento = contadorLineamientoGeneralOut + 1

				}

				for i := 0; i < len(armoPI); i++ {
					datosArmo := armoPI[i]
					auxFactor := datosArmo["nombreFactor"]
					contadorFactorGeneralIn = contadorFactor

					// cuerpo del excel
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorFactor), auxFactor)
					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "D"+fmt.Sprint(contadorFactor), "N"+fmt.Sprint(contadorFactor), stylecontent)
					consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorFactor, 70)

					lineamientos := datosArmo["lineamientos"]

					contadorLineamientos := contadorFactor

					for j := 0; j < len(lineamientos.([]map[string]interface{})); j++ {
						auxLineamiento := lineamientos.([]map[string]interface{})[j]
						contadorLineamientoPIIn = contadorFactor

						consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorLineamientos), auxLineamiento["nombreLineamiento"])
						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "E"+fmt.Sprint(contadorLineamientos), "E"+fmt.Sprint(contadorLineamientos), stylecontent)

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

							consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorEstrategias), auxEstrategia["descripcionEstrategia"])
							consolidadoExcelPlanAnual.SetCellStyle(sheetName, "F"+fmt.Sprint(contadorEstrategias), "F"+fmt.Sprint(contadorEstrategias), stylecontent)

							if k == len(estrategiasPI)-1 {
								contadorEstrategiaPIOut = contadorLineamientos
							} else {
								contadorEstrategias = contadorEstrategias + 1
							}
						}

						contadorEstrategias = contadorLineamientos

						consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorLineamientoPIIn), "E"+fmt.Sprint(contadorLineamientoPIOut))

					}
					contadorLineamientos = contadorFactor
					contadorFactorGeneralOut = contadorLineamientoPIOut

					consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorFactorGeneralIn), "D"+fmt.Sprint(contadorFactorGeneralOut))

					contadorFactor = contadorFactorGeneralOut + 1

				}

				consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorDataGeneral), datosExcelPlan["numeroActividad"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Ponderación de la actividad"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Periodo de ejecución"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Actividad general"])
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Tareas"])

				if contadorLineamientoGeneralOut > contadorFactorGeneralOut {
					contadorFactorGeneralOut = contadorLineamientoGeneralOut
					contadorFactor = contadorFactorGeneralOut + 1
				} else if contadorLineamientoGeneralOut < contadorFactorGeneralOut {
					contadorLineamientoGeneralOut = contadorFactorGeneralOut
					contadorLineamiento = contadorLineamientoGeneralOut + 1
				}

				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "G"+fmt.Sprint(contadorLineamiento), "N"+fmt.Sprint(contadorLineamiento), stylecontent)

				consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorLineamientoGeneralIn), "A"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorMetaGeneralIn), "B"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorEstrategiaPEDOut), "C"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorFactorGeneralIn), "D"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorLineamientoPIOut), "E"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorEstrategiaPIOut), "F"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorDataGeneral), "G"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorDataGeneral), "H"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorDataGeneral), "I"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorLineamientoGeneralOut))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorDataGeneral), "K"+fmt.Sprint(contadorLineamientoGeneralOut))

				indicadores := datosComplementarios["indicadores"].(map[string]interface{})
				contadorIndicadores := contadorDataGeneral
				for id, indicador := range indicadores {
					_ = id
					auxIndicador := indicador
					var nombreIndicador string
					var formula string
					var meta string

					for key, element := range auxIndicador.(map[string]interface{}) {
						if strings.Contains(strings.ToLower(key), "nombre") {
							nombreIndicador = element.(string)
						}
						if strings.Contains(strings.ToLower(key), "formula") || strings.Contains(strings.ToLower(key), "fórmula") {
							formula = element.(string)
						}
						if strings.Contains(strings.ToLower(key), "meta") {
							meta = element.(string)
						}

					}

					consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorIndicadores), nombreIndicador)
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorIndicadores), formula)
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorIndicadores), meta)

					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "L"+fmt.Sprint(contadorIndicadores), "L"+fmt.Sprint(contadorIndicadores), stylecontent)
					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "M"+fmt.Sprint(contadorIndicadores), "M"+fmt.Sprint(contadorIndicadores), stylecontent)
					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "N"+fmt.Sprint(contadorIndicadores), "N"+fmt.Sprint(contadorIndicadores), stylecontent)

					contadorIndicadores = contadorIndicadores + 1

				}

				contadorIndicadores--

				consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorLineamientoGeneralOut), "A"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorMetaGeneralOut), "B"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorEstrategiaPEDOut), "C"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorFactorGeneralOut), "D"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorLineamientoPIOut), "E"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorEstrategiaPIOut), "F"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorDataGeneral), "G"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorDataGeneral), "H"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "I"+fmt.Sprint(contadorDataGeneral), "I"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "J"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorIndicadores))
				consolidadoExcelPlanAnual.MergeCell(sheetName, "K"+fmt.Sprint(contadorDataGeneral), "K"+fmt.Sprint(contadorIndicadores))
				contadorDataGeneral = contadorIndicadores + 1
				contadorLineamiento = contadorIndicadores + 1
				contadorFactor = contadorIndicadores + 1

				consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)

			}

			contadorGeneral = contadorDataGeneral + 2
			arregloPlanAnual = nil
		}

		//consolidadoExcelPlanAnual.SaveAs("plan_anual.xlsx")

		buf, _ := consolidadoExcelPlanAnual.WriteToBuffer()
		strings.NewReader(buf.String())

		encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))

		dataSend := make(map[string]interface{})

		dataSend["generalData"] = arregloInfoReportes
		dataSend["excelB64"] = encoded

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}

// Necesidades ...
// @Title Necesidades
// @Description post Reportes by id
// @Param	body		body 	{}	true		"body for Plan content"
// @Success 201 {object} models.Reportes
// @Failure 403 :plan_id is empty
// @router /necesidades [post]
func (c *ReportesController) Necesidades() {
	var body map[string]interface{}
	var respuesta map[string]interface{}
	var respuestaIdentificaciones map[string]interface{}
	var identificaciones []map[string]interface{}
	var planes []map[string]interface{}
	var recursos []map[string]interface{}
	var docentes map[string]interface{}
	var recursosGeneral []map[string]interface{}
	var docentesGeneral map[string]interface{}
	docentesPregrado := make(map[string]interface{})
	docentesPosgrado := make(map[string]interface{})
	var arrDataDocentes []map[string]interface{}

	docentesPregrado["tco"] = 0
	docentesPregrado["mto"] = 0
	docentesPregrado["hch"] = 0
	docentesPregrado["hcp"] = 0
	docentesPregrado["valor"] = 0

	docentesPosgrado["hch"] = 0
	docentesPosgrado["hcp"] = 0
	docentesPosgrado["valor"] = 0

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

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

	necesidadesExcel := excelize.NewFile()
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	stylecontent, _ := necesidadesExcel.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)
	styletitles, _ := necesidadesExcel.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"font":{"bold":true,"family":"Arial", "size":26,"color":"#000000"},
					"fill":{"type":"pattern","pattern":1,"color":["#F2F2F2"]},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)

	stylesubtitles, _ := necesidadesExcel.NewStyle(`{
					"alignment":{"horizontal":"left","vertical":"center","wrap_text":true},
					"font":{"bold":true,"family":"Arial", "size":20,"color":"#000000"},
					"fill":{"type":"pattern","pattern":1,"color":["#F2F2F2"]},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)
	stylehead, _ := necesidadesExcel.NewStyle(`{
					"alignment":{"horizontal":"center","vertical":"center","wrap_text":true},
					"font":{"bold":true,"color":"#FFFFFF"},
					"fill":{"type":"pattern","pattern":1,"color":["#CC0000"]},
					"border":[{"type":"right","color":"#000000","style":1},{"type":"left","color":"#000000","style":1},{"type":"top","color":"#000000","style":1},{"type":"bottom","color":"#000000","style":1}]
				}`)

	necesidadesExcel.NewSheet("Necesidades")

	necesidadesExcel.MergeCell("Necesidades", "A1", "F1")
	necesidadesExcel.MergeCell("Necesidades", "A1", "A2")
	necesidadesExcel.MergeCell("Necesidades", "A3", "F3")
	necesidadesExcel.MergeCell("Necesidades", "A3", "A4")

	necesidadesExcel.SetColWidth("Necesidades", "A", "F", 30)

	necesidadesExcel.SetCellValue("Necesidades", "A1", "Necesidades Presupuestales")
	necesidadesExcel.SetCellStyle("Necesidades", "A1", "F1", styletitles)

	necesidadesExcel.SetCellValue("Necesidades", "A3", "Identificación de recursos:")
	necesidadesExcel.SetCellStyle("Necesidades", "A3", "F3", stylesubtitles)

	necesidadesExcel.SetCellStyle("Necesidades", "A200", "F200", stylecontent)

	necesidadesExcel.SetCellValue("Necesidades", "A5", "Código del rubro")
	necesidadesExcel.SetCellValue("Necesidades", "B5", "Nombre del rubro")
	necesidadesExcel.SetCellValue("Necesidades", "C5", "Valor")
	necesidadesExcel.SetCellStyle("Necesidades", "A5", "C5", stylehead)
	necesidadesExcel.SetRowHeight("Necesidades", 5, 35)
	contador := 6

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planes)

		for i := 0; i < len(planes); i++ {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=plan_id:"+planes[i]["_id"].(string), &respuestaIdentificaciones); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuestaIdentificaciones, &identificaciones)

				for i := 0; i < len(identificaciones); i++ {
					identificacion := identificaciones[i]
					if strings.Contains(strings.ToLower(identificacion["nombre"].(string)), "recurso") {
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
					} else if strings.Contains(strings.ToLower(identificacion["nombre"].(string)), "docente") {
						var dato map[string]interface{}
						var data_identi []map[string]interface{}
						if identificacion["dato"] != nil && identificacion["dato"] != "{}" {
							result := make(map[string]interface{})
							dato_str := identificacion["dato"].(string)
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

							docentes = result
						}
					}
				}

				for i := 0; i < len(recursos); i++ {
					var aux bool
					if len(recursos) == 0 {
						recursosGeneral = append(recursosGeneral, recursos[i])
					} else {
						for j := 0; j < len(recursosGeneral); j++ {
							if recursosGeneral[j]["codigo"] == recursos[i]["codigo"] {
								aux = true
								break
							} else {
								aux = false
							}
						}

						if !aux {
							recursosGeneral = append(recursosGeneral, recursos[i])
						}
					}
				}

				if docentes["rubros"] != nil {
					rubros := docentes["rubros"].([]map[string]interface{})
					for i := 0; i < len(rubros); i++ {
						if rubros[i]["rubro"] != "" {
							var respuestaRubro map[string]interface{}
							rubro := make(map[string]interface{})
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanCuentasService")+"/arbol_rubro/"+rubros[i]["rubro"].(string), &respuestaRubro); err == nil {
								aux := respuestaRubro["Body"].(map[string]interface{})
								rubro["codigo"] = aux["Codigo"]
								rubro["nombre"] = aux["Nombre"]
								rubro["categoria"] = rubros[i]["categoria"]
								recursosGeneral = append(recursosGeneral, rubro)
							}
						}
					}
				}

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
		}

		for i := 0; i < len(recursosGeneral); i++ {
			if recursosGeneral[i]["categoria"] != nil {
				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "servicio") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + primaServicios
						}
					} else {
						recursosGeneral[i]["valor"] = primaServicios
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "navidad") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + primaNavidad
						}
					} else {
						recursosGeneral[i]["valor"] = primaNavidad
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "vacaciones") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + primaVacaciones
						}
					} else {
						recursosGeneral[i]["valor"] = primaVacaciones
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "bonificacion") || strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "bonificación") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + bonificacion
						}
					} else {
						recursosGeneral[i]["valor"] = bonificacion
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "interes") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "cesantía") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + interesesCesantias
						}
					} else {
						recursosGeneral[i]["valor"] = interesesCesantias
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "cesantía") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "público") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + cesantiasPublicas
						}
					} else {
						recursosGeneral[i]["valor"] = cesantiasPublicas
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "cesantía") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "privado") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + cesantiasPrivadas
						}
					} else {
						recursosGeneral[i]["valor"] = cesantiasPrivadas
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "salud") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + salud
						}
					} else {
						recursosGeneral[i]["valor"] = salud
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "pension") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "público") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + pensionesPublicas
						}
					} else {
						recursosGeneral[i]["valor"] = pensionesPublicas
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "pension") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "privado") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + pensionesPrivadas
						}
					} else {
						recursosGeneral[i]["valor"] = pensionesPrivadas
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "arl") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + arl
						}
					} else {
						recursosGeneral[i]["valor"] = arl
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "ccf") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + caja
						}
					} else {
						recursosGeneral[i]["valor"] = caja
					}
				}

				if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "icbf") {
					if recursosGeneral[i]["valor"] != nil {
						strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						arrValor := strings.Split(strValor, ".")
						auxValor, err := strconv.Atoi(arrValor[0])
						if err == nil {
							recursosGeneral[i]["valor"] = auxValor + icbf
						}
					} else {
						recursosGeneral[i]["valor"] = icbf
					}
				}
			}

		}

		//Completado de tablas
		for i := 0; i < len(recursosGeneral); i++ {
			necesidadesExcel.SetCellValue("Necesidades", "A"+fmt.Sprint(contador), recursosGeneral[i]["codigo"])
			if recursosGeneral[i]["Nombre"] != nil {
				necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), recursosGeneral[i]["Nombre"])
				necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), recursosGeneral[i]["valor"])

			} else {
				necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), recursosGeneral[i]["nombre"])
				necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), ac.FormatMoney(recursosGeneral[i]["valor"]))

			}
			necesidadesExcel.SetCellStyle("Necesidades", "A"+fmt.Sprint(contador), "C"+fmt.Sprint(contador), stylecontent)
			contador++
		}

		contador++
		contador++
		necesidadesExcel.MergeCell("Necesidades", "A"+fmt.Sprint(contador), "F"+fmt.Sprint(contador))
		necesidadesExcel.MergeCell("Necesidades", "A"+fmt.Sprint(contador), "A"+fmt.Sprint(contador+1))

		necesidadesExcel.SetColWidth("Necesidades", "A", "A", 30)
		necesidadesExcel.SetColWidth("Necesidades", "B", "E", 15)
		necesidadesExcel.SetColWidth("Necesidades", "F", "F", 30)
		necesidadesExcel.SetColWidth("Necesidades", "G", "H", 15)
		necesidadesExcel.SetColWidth("Necesidades", "I", "I", 15)

		necesidadesExcel.SetCellValue("Necesidades", "A"+fmt.Sprint(contador), "Docentes por tipo de vinculación:")
		necesidadesExcel.SetCellStyle("Necesidades", "A"+fmt.Sprint(contador), "F"+fmt.Sprint(contador), stylesubtitles)

		necesidadesExcel.SetCellStyle("Necesidades", "A200", "F200", stylecontent)
		contador++
		contador++
		necesidadesExcel.SetCellValue("Necesidades", "A"+fmt.Sprint(contador), "Facultad")
		necesidadesExcel.MergeCell("Necesidades", "A"+fmt.Sprint(contador), "A"+fmt.Sprint(contador+1))
		necesidadesExcel.MergeCell("Necesidades", "B"+fmt.Sprint(contador), "F"+fmt.Sprint(contador))

		necesidadesExcel.MergeCell("Necesidades", "G"+fmt.Sprint(contador), "I"+fmt.Sprint(contador))

		necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), "Pregrado")
		necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "Posgrado")

		necesidadesExcel.SetCellStyle("Necesidades", "A"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), stylehead)

		contador++
		necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), "TCO")
		necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), "MTO")
		necesidadesExcel.SetCellValue("Necesidades", "D"+fmt.Sprint(contador), "HCH")
		necesidadesExcel.SetCellValue("Necesidades", "E"+fmt.Sprint(contador), "HCP")
		necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), "Valor")
		necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), "HCH")
		necesidadesExcel.SetCellValue("Necesidades", "H"+fmt.Sprint(contador), "HCP")
		necesidadesExcel.SetCellValue("Necesidades", "I"+fmt.Sprint(contador), "Valor")
		necesidadesExcel.SetCellStyle("Necesidades", "B"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), stylecontent)

		contador++

		for i := 0; i < len(arrDataDocentes); i++ {
			necesidadesExcel.SetCellValue("Necesidades", "A"+fmt.Sprint(contador), arrDataDocentes[i]["nombreFacultad"])
			necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), arrDataDocentes[i]["tco"])
			necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), arrDataDocentes[i]["mto"])
			necesidadesExcel.SetCellValue("Necesidades", "D"+fmt.Sprint(contador), arrDataDocentes[i]["hch"])
			necesidadesExcel.SetCellValue("Necesidades", "E"+fmt.Sprint(contador), arrDataDocentes[i]["hcp"])
			necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), ac.FormatMoney(arrDataDocentes[i]["valor"]))
			necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), arrDataDocentes[i]["hchPos"])
			necesidadesExcel.SetCellValue("Necesidades", "H"+fmt.Sprint(contador), arrDataDocentes[i]["hcpPos"])
			necesidadesExcel.SetCellValue("Necesidades", "I"+fmt.Sprint(contador), ac.FormatMoney(arrDataDocentes[i]["valor"]))
			necesidadesExcel.SetCellStyle("Necesidades", "A"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), stylecontent)
			contador++
		}

		buf, _ := necesidadesExcel.WriteToBuffer()
		strings.NewReader(buf.String())

		encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": encoded}

		necesidadesExcel.SaveAs("necesidades.xlsx")

	} else {
		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
		c.Abort("400")
	}

	c.ServeJSON()
}
