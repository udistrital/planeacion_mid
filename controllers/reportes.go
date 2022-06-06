package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
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
}

func CreateExcel(f *excelize.File, dir string) {
	fmt.Println(dir)
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
	fmt.Println(reflect.TypeOf(consolidadoExcelPlanAnual))

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
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "D3", "Factores")
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

					fmt.Println(contadorLineamientoGeneralIn + contadorFactorGeneralOut + contadorMetaGeneralIn + contadorMetaGeneralOut + contadorEstrategiaPEDIn + contadorFactorGeneralIn + contadorEstrategiaPEDOut + contadorLineamientoPIIn + contadorLineamientoPIOut + contadorEstrategiaPIIn + contadorEstrategiaPIOut)

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

						fmt.Println("Metas---------------")
						fmt.Println(metas)
						for j := 0; j < len(metas.([]map[string]interface{})); j++ {
							auxMeta := metas.([]map[string]interface{})[j]
							contadorMetaGeneralIn = contadorLineamiento
							fmt.Println("auxmeta........")
							fmt.Println(auxMeta)
							fmt.Println(contadorMetas)
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

						fmt.Println("MERGE ABER")
						fmt.Println(contadorLineamiento)
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

						fmt.Println("Lineamientos PI---------------")
						fmt.Println(lineamientos)
						for j := 0; j < len(lineamientos.([]map[string]interface{})); j++ {
							auxLineamiento := lineamientos.([]map[string]interface{})[j]
							contadorLineamientoPIIn = contadorFactor
							fmt.Println("aux LinemientoPI........")
							fmt.Println(auxLineamiento)
							fmt.Println(contadorLineamientos)
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

						fmt.Println("MERGE ABER")
						fmt.Println(contadorFactor)
						consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorFactorGeneralIn), "D"+fmt.Sprint(contadorFactorGeneralOut))

						contadorFactor = contadorFactorGeneralOut + 1

					}

					consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorDataGeneral), datosExcelPlan["numeroActividad"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Ponderación de la actividad"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Periodo de ejecución"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Actividad general"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Tareas"])

					fmt.Println("contadores finales lineamiento facor ....")
					fmt.Println(contadorLineamientoGeneralOut)
					fmt.Println(contadorFactorGeneralOut)
					if contadorLineamientoGeneralOut > contadorFactorGeneralOut {
						contadorFactorGeneralOut = contadorLineamientoGeneralOut
						contadorFactor = contadorFactorGeneralOut + 1
					} else if contadorLineamientoGeneralOut < contadorFactorGeneralOut {
						contadorLineamientoGeneralOut = contadorFactorGeneralOut
						contadorLineamiento = contadorLineamientoGeneralOut + 1
					}

					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "G"+fmt.Sprint(contadorLineamiento), "N"+fmt.Sprint(contadorLineamiento), stylecontent)

					fmt.Println("contadores data general lineamientos despues de indicadores")
					fmt.Println(contadorDataGeneral)
					fmt.Println(contadorLineamientoGeneralOut)
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

						fmt.Println("Adentro del arreglo de indicadores")
						fmt.Println(contadorIndicadores)
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
					fmt.Println("contador indicadores")
					fmt.Println(contadorIndicadores)
					fmt.Println(contadorDataGeneral)
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
					fmt.Println("Contadores inciales ........")
					fmt.Println(contadorLineamientoGeneralIn)
					fmt.Println(contadorMetaGeneralIn)
					fmt.Println(contadorEstrategiaPEDIn)
					fmt.Println(contadorFactorGeneralIn)
					fmt.Println(contadorLineamientoPIIn)
					fmt.Println(contadorEstrategiaPIOut)

					fmt.Println("Contadores importantes ........")
					fmt.Println(contadorLineamientoGeneralOut)
					fmt.Println(contadorMetaGeneralOut)
					fmt.Println(contadorEstrategiaPEDOut)
					fmt.Println(contadorFactorGeneralOut)
					fmt.Println(contadorLineamientoPIOut)
					fmt.Println(contadorEstrategiaPIOut)
					fmt.Println(contadorDataGeneral)

					consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)

				}
			}

			consolidadoExcelPlanAnual.SaveAs("plan_anual.xlsx")

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
	var respuestaGeneral map[string]interface{}
	var planesFilterGeneral []map[string]interface{}
	var respuesta map[string]interface{}
	var planesFilter []map[string]interface{}
	var res map[string]interface{}
	var resPresupuesto map[string]interface{}
	var resArmo map[string]interface{}
	var resEstrategia map[string]interface{}
	var resMeta map[string]interface{}
	var resLineamiento map[string]interface{}
	var resPlan map[string]interface{}
	var respuestaUnidad []map[string]interface{}
	var hijosArmo []map[string]interface{}
	var estrategiaData []map[string]interface{}
	var metaData []map[string]interface{}
	var LineamientoData []map[string]interface{}
	var planData []map[string]interface{}
	var subgrupos []map[string]interface{}
	var plan_id string
	var actividadName string
	var arregloPlanAnual []map[string]interface{}
	var arregloEstrategia []map[string]interface{}
	var arregloMetaEst []map[string]interface{}
	var identificacion map[string]interface{}
	var datoPresupuesto map[string]interface{}
	var identificacionres []map[string]interface{}
	var data_identi []map[string]interface{}
	var unidadId string
	var nombrePlanDesarrollo string
	var nombreUnidad string
	var unidadNombre string
	var arregloLineamieto []map[string]interface{}

	consolidadoExcelPlanAnual := excelize.NewFile()
	fmt.Println(reflect.TypeOf(consolidadoExcelPlanAnual))

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuestaGeneral); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuestaGeneral, &planesFilterGeneral)
		fmt.Println("----------planes--------------")
		fmt.Println(planesFilterGeneral)
		for datoUnidad := 0; datoUnidad < len(planesFilterGeneral); datoUnidad++ {
			planesUnidad := planesFilterGeneral[datoUnidad]
			unidadId = planesUnidad["dependencia_id"].(string)

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",dependencia_id:"+unidadId, &respuesta); err == nil {
				helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
				for planes := 0; planes < len(planesFilter); planes++ {
					planesFilterData := planesFilter[planes]
					plan_id = planesFilterData["_id"].(string)
					generalData := make(map[string]interface{})

					data_identi = nil
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/identificacion?query=activo:true,plan_id:"+plan_id+",tipo_identificacion_id:"+"617b6630f6fc97b776279afa", &resPresupuesto); err == nil {
						helpers.LimpiezaRespuestaRefactor(resPresupuesto, &identificacionres)
						// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": identificacionres}
						identificacion = identificacionres[0]
						if identificacion["dato"] != "{}" {
							dato_str := identificacion["dato"].(string)
							json.Unmarshal([]byte(dato_str), &datoPresupuesto)
							// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dato}
							for key := range datoPresupuesto {
								element := datoPresupuesto[key].(map[string]interface{})
								if element["activo"] == true {
									delete(element, "actividades")
									delete(element, "activo")
									delete(element, "index")
									data_identi = append(data_identi, element)
								}

							}
							// arreglo = append(arreglo, data_identi...)

						} else {
							c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": ""}
						}

					} else {
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
						c.Abort("400")
					}

					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id, &res); err == nil {
						helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

						for i := 0; i < len(subgrupos); i++ {
							if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
								actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
								for j := 0; j < len(actividades); j++ {
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
										armonizacionTercerNivel := armonizacionTercer["armo"]
										for datoGeneral := 0; datoGeneral < len(treeDatos); datoGeneral++ {
											treeDato := treeDatos[datoGeneral]
											treeData := treeDatas[0]
											if treeDato["sub"] == "" {
												datosArmonizacion[treeDato["nombre"].(string)] = treeData[fmt.Sprint(treeDato["id"])]
											}
										}
										fmt.Println(treeDatos)
										treeIndicador := treeDatos[4]
										subIndicador := treeIndicador["sub"].([]map[string]interface{})

										for ind := 0; ind < len(subIndicador); ind++ {
											subIndicadorRes := subIndicador[ind]
											treeData := treeDatas[0]
											titulosArmonizacion[subIndicadorRes["nombre"].(string)] = treeData[fmt.Sprint(subIndicadorRes["id"])]
										}
										datosArmonizacion["indicadores"] = titulosArmonizacion

										if armonizacionTercerNivel != nil {
											delimitador := ","
											output := strings.Split(armonizacionTercerNivel.(string), delimitador)
											estrategiaDesc := make(map[string]interface{})
											metaEstrategica := make(map[string]interface{})
											lineamientoDesc := make(map[string]interface{})
											for estrategiaArmo := 0; estrategiaArmo < len(output); estrategiaArmo++ {
												estrategiaId := output[estrategiaArmo]
												if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=_id:"+estrategiaId, &resEstrategia); err == nil {
													helpers.LimpiezaRespuestaRefactor(resEstrategia, &estrategiaData)
													estrategiaEst := estrategiaData[0]
													idPadre := estrategiaEst["padre"].(string)
													estrategiaDataStr := fmt.Sprint(estrategiaData)
													if estrategiaDataStr != "[]" {
														estrategiaDesc["descripcionEstrategia"] = estrategiaEst["descripcion"]
														estrategiaDesc["nombreEstrategia"] = estrategiaEst["nombre"]

														if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=_id:"+idPadre, &resMeta); err == nil {
															helpers.LimpiezaRespuestaRefactor(resMeta, &metaData)
															metaEst := metaData[0]
															padreMeta := metaEst["padre"].(string)
															metaDataStr := fmt.Sprint(metaData)

															if metaDataStr != "[]" {
																metaEstrategica["nombreMeta"] = metaEst["nombre"]

																if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=_id:"+padreMeta, &resLineamiento); err == nil {
																	helpers.LimpiezaRespuestaRefactor(resLineamiento, &LineamientoData)
																	lineamientoEst := LineamientoData[0]
																	lineamientoDesc["nombreLineamiento"] = lineamientoEst["nombre"]
																	padreLineamiento := lineamientoEst["padre"].(string)

																	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=_id:"+padreLineamiento, &resPlan); err == nil {
																		helpers.LimpiezaRespuestaRefactor(resPlan, &planData)
																		planDesarrollo := planData[0]
																		nombrePlanDesarrollo = planDesarrollo["nombre"].(string)
																		lineamientoDesc["nombrePlanDesarrollo"] = nombrePlanDesarrollo

																	} else {
																		c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
																		c.Abort("400")
																	}

																} else {
																	c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
																	c.Abort("400")
																}
															} else {
																lineamientoDesc["nombreLineamiento"] = metaEst["nombre"]
																if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=_id:"+padreMeta, &resPlan); err == nil {
																	helpers.LimpiezaRespuestaRefactor(resPlan, &planData)
																	planDesarrollo := planData[0]
																	nombrePlanDesarrollo = planDesarrollo["nombre"].(string)
																	lineamientoDesc["nombrePlanDesarrollo"] = nombrePlanDesarrollo

																} else {
																	c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
																	c.Abort("400")
																}
															}

														} else {
															c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
															c.Abort("400")
														}

													} else {
														lineamientoDesc["nombreLineamiento"] = estrategiaEst["nombre"]
														if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=_id:"+idPadre, &resPlan); err == nil {
															helpers.LimpiezaRespuestaRefactor(resPlan, &planData)
															planDesarrollo := planData[0]
															nombrePlanDesarrollo = planDesarrollo["nombre"].(string)
															lineamientoDesc["nombrePlanDesarrollo"] = nombrePlanDesarrollo

														} else {
															c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
															c.Abort("400")
														}

													}

												} else {
													c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
													c.Abort("400")
												}
											}
											arregloEstrategia = append(arregloEstrategia, estrategiaDesc)
											metaEstrategica["estrategias"] = arregloEstrategia
											arregloMetaEst = append(arregloMetaEst, metaEstrategica)
											lineamientoDesc["meta"] = arregloMetaEst
											arregloLineamieto = append(arregloLineamieto, lineamientoDesc)

										}

									} else {
										c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
										c.Abort("400")
									}

									// generalData := make(map[string]interface{})

									if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+unidadId, &respuestaUnidad); err == nil {
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
									generalData["datosComplementarios"] = datosArmonizacion
									// generalData["presupuesto"] = data_identi

									// arregloPlanAnual = append(arregloPlanAnual, generalData)

								}
								break
							}
						}
					} else {
						c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
						c.Abort("400")
					}

					generalData["presupuesto"] = data_identi

					arregloPlanAnual = append(arregloPlanAnual, generalData)

				}

			} else {
				c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
				c.Abort("400")
			}

		}

		contadorLineamiento := 4
		contadorMeta := 4
		contadorEstrategia := 4
		// definicion de estilos del excel
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

		// creacion del excel
		for excelPlan := 0; excelPlan < len(arregloPlanAnual); excelPlan++ {
			datosExcelPlan := arregloPlanAnual[excelPlan]
			armo := datosExcelPlan["datosArmonizacion"].([]map[string]interface{})

			unidadNombreIn := datosExcelPlan["nombreUnidad"].(string)
			if excelPlan == 0 {
				unidadNombre = unidadNombreIn
			}
			if unidadNombreIn == unidadNombre {
				unidadNombre = datosExcelPlan["nombreUnidad"].(string)
				nombreHoja := fmt.Sprint(unidadNombre)
				sheetName := nombreHoja
				indexPlan := consolidadoExcelPlanAnual.NewSheet(sheetName)
				consolidadoExcelPlanAnual.MergeCell(sheetName, "A1", "N1")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "A2", "C2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "D2", "D3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E2", "E3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "F2", "F3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "G2", "G3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "H2", "H3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "I2", "K2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "L2", "N2")
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 1, 20)
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 2, 20)
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 3, 20)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "A", "C", 70)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "K", 50)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "N", 50)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "F", 20)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "G", "H", 80)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A1", "K1", stylehead)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A2", "N2", styletitles)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A3", "N3", styletitles)

				for dataExcelArmo := 0; dataExcelArmo < len(armo); dataExcelArmo++ {
					contadorLineamientoIn := contadorLineamiento
					datosArmo := armo[dataExcelArmo]
					lineamiento := datosArmo["nombreLineamiento"]
					planDesarrolloName := datosArmo["nombrePlanDesarrollo"]
					tituloExcel := fmt.Sprint("Plan de acción 2022 ", unidadNombre)
					// encabezado excel
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A1", tituloExcel)
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A2", planDesarrolloName)
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A3", "Lineamiento")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "B3", "Meta")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "C3", "Estrategia")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "D3", "N°.")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "E3", "Ponderación de la actividad")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "F3", "Periodo de ejecución")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "G3", "Actividad")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "H3", "Tareas")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "I2", "Indicador")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "I3", "Nombre")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "J3", "Fórmula")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "K3", "Meta")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Presupuesto")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "L3", "Código del rubro")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "M3", "Nombre del rubro")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "N3", "Valor")

					// cuerpo del excel
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A"+fmt.Sprint(contadorLineamiento), lineamiento)
					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorLineamiento), "K"+fmt.Sprint(contadorLineamiento), stylecontent)
					consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorLineamiento, 70)
					metaEstr := datosArmo["meta"].([]map[string]interface{})
					for dataExcelMetaEstr := 0; dataExcelMetaEstr < len(metaEstr); dataExcelMetaEstr++ {
						contadorMetaIn := contadorMeta
						datosMeta := metaEstr[dataExcelMetaEstr]
						nombreMetaEstr := datosMeta["nombreMeta"]
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorMeta), nombreMetaEstr)
						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorMeta), "K"+fmt.Sprint(contadorMeta), stylecontent)
						consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorMeta, 70)

						estrategia := datosMeta["estrategias"].([]map[string]interface{})
						for dataExcelEstrategia := 0; dataExcelEstrategia < len(estrategia); dataExcelEstrategia++ {
							contadorEstrategiaIn := contadorEstrategia
							contadorEstrategiaOut := contadorEstrategia
							datosEstrategia := estrategia[dataExcelEstrategia]
							descEstrategia := datosEstrategia["descripcionEstrategia"]
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorEstrategia), descEstrategia)
							consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorEstrategia), "N"+fmt.Sprint(contadorEstrategia), stylecontent)
							consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorEstrategia, 70)

							datosComplementarios := datosExcelPlan["datosComplementarios"].(map[string]interface{})
							numeroActividadExcel := datosExcelPlan["numeroActividad"]
							ponderacionActividad := datosComplementarios["Ponderación de la actividad"]
							periodoEjecucion := datosComplementarios["Periodo de ejecución"]
							actividadGeneral := datosComplementarios["Actividad general"]
							tareas := datosComplementarios["Tareas"]
							indicador := datosComplementarios["indicadores"].(map[string]interface{})
							nombreIndicador1 := indicador["Nombre del indicador"]
							nombreIndicador2 := indicador["Nombre indicador 2"]
							formula1 := indicador["Fórmula"]
							formula2 := indicador["Fórmula 2"]
							meta1 := indicador["Meta"]
							meta2 := indicador["Meta 2"]

							consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorEstrategia), numeroActividadExcel)
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorEstrategia), ponderacionActividad)
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorEstrategia), periodoEjecucion)
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorEstrategia), actividadGeneral)
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorEstrategia), tareas)
							if len(indicador) > 3 {
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorEstrategia), nombreIndicador1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorEstrategia+1), nombreIndicador2)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorEstrategia), formula1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorEstrategia+1), formula2)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorEstrategia), meta1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorEstrategia+1), meta2)
								consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorEstrategia+1), "N"+fmt.Sprint(contadorEstrategia+1), stylecontent)
								consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorEstrategia+1, 70)
								contadorEstrategia = contadorEstrategia + 1
								contadorEstrategiaOut = contadorEstrategia
							} else {
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorEstrategia), nombreIndicador1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorEstrategia), formula1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorEstrategia), meta1)
							}
							consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorEstrategiaIn), "C"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorEstrategiaIn), "D"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorEstrategiaIn), "E"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorEstrategiaIn), "F"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorEstrategiaIn), "G"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorEstrategiaIn), "H"+fmt.Sprint(contadorEstrategiaOut))

							if contadorEstrategiaIn > contadorEstrategia {
								contadorMeta = contadorEstrategiaIn
							} else {
								contadorMeta = contadorEstrategia
							}
							contadorEstrategia = contadorEstrategia + 1
						}

						contadorMetaOut := contadorMeta
						consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorMetaIn), "B"+fmt.Sprint(contadorMetaOut))
						if contadorMeta >= contadorLineamiento {
							contadorLineamiento = contadorMeta
						}
						contadorMeta = contadorMeta + 1
					}

					contadorLineamientoOut := contadorLineamiento

					consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorLineamientoIn), "A"+fmt.Sprint(contadorLineamientoOut))

					contadorLineamiento = contadorLineamiento + 1
				}
				consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
			} else {
				contadorLineamiento = 4
				contadorMeta = 4
				contadorEstrategia = 4
				// contadorPresupuestoOut := 0
				unidadNombre = datosExcelPlan["nombreUnidad"].(string)
				nombreHoja := fmt.Sprint(unidadNombre)
				sheetName := nombreHoja
				indexPlan := consolidadoExcelPlanAnual.NewSheet(sheetName)
				consolidadoExcelPlanAnual.MergeCell(sheetName, "A1", "N1")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "A2", "C2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "D2", "D3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "E2", "E3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "F2", "F3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "G2", "G3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "H2", "H3")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "I2", "K2")
				consolidadoExcelPlanAnual.MergeCell(sheetName, "L2", "N2")
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 1, 20)
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 2, 20)
				consolidadoExcelPlanAnual.SetRowHeight(sheetName, 3, 20)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "A", "C", 70)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "I", "K", 50)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "L", "N", 50)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "E", "F", 20)
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "G", "H", 80)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A1", "K1", stylehead)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A2", "N2", styletitles)
				consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A3", "N3", styletitles)

				for dataExcelArmo := 0; dataExcelArmo < len(armo); dataExcelArmo++ {
					contadorLineamientoIn := contadorLineamiento
					datosArmo := armo[dataExcelArmo]
					lineamiento := datosArmo["nombreLineamiento"]
					planDesarrolloName := datosArmo["nombrePlanDesarrollo"]
					tituloExcel := fmt.Sprint("Plan de acción 2022 ", unidadNombre)
					// encabezado excel
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A1", tituloExcel)
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A2", planDesarrolloName)
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A3", "Lineamiento")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "B3", "Meta")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "C3", "Estrategia")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "D3", "N°.")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "E3", "Ponderación de la actividad")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "F3", "Periodo de ejecución")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "G3", "Actividad")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "H3", "Tareas")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "I2", "Indicador")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "I3", "Nombre")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "J3", "Fórmula")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "K3", "Meta")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Presupuesto")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "L3", "Código del rubro")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "M3", "Nombre del rubro")
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "N3", "Valor")

					// cuerpo del excel
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "A"+fmt.Sprint(contadorLineamiento), lineamiento)
					consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorLineamiento), "K"+fmt.Sprint(contadorLineamiento), stylecontent)
					consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorLineamiento, 70)
					contadorPresupuesto := 4

					presupuestoExc := datosExcelPlan["presupuesto"].([]map[string]interface{})
					if dataExcelArmo == 0 {
						if presupuestoExc != nil {
							for dataExcelPresupuesto := 0; dataExcelPresupuesto < len(presupuestoExc); dataExcelPresupuesto++ {
								datosPresupuesto := presupuestoExc[dataExcelPresupuesto]
								nombrePresupuesto := datosPresupuesto["Nombre"]
								codigoPresupuesto := datosPresupuesto["codigo"]
								valorPresupuesto := datosPresupuesto["valor"]
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorPresupuesto), codigoPresupuesto)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "M"+fmt.Sprint(contadorPresupuesto), nombrePresupuesto)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "N"+fmt.Sprint(contadorPresupuesto), valorPresupuesto)
								consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorPresupuesto, 70)

								contadorPresupuesto = contadorPresupuesto + 1

							}

						}
					}

					metaEstr := datosArmo["meta"].([]map[string]interface{})
					for dataExcelMetaEstr := 0; dataExcelMetaEstr < len(metaEstr); dataExcelMetaEstr++ {
						contadorMetaIn := contadorMeta
						datosMeta := metaEstr[dataExcelMetaEstr]
						nombreMetaEstr := datosMeta["nombreMeta"]
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "B"+fmt.Sprint(contadorMeta), nombreMetaEstr)
						consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorMeta), "K"+fmt.Sprint(contadorMeta), stylecontent)
						consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorMeta, 70)

						estrategia := datosMeta["estrategias"].([]map[string]interface{})
						for dataExcelEstrategia := 0; dataExcelEstrategia < len(estrategia); dataExcelEstrategia++ {
							contadorEstrategiaIn := contadorEstrategia
							contadorEstrategiaOut := contadorEstrategia
							datosEstrategia := estrategia[dataExcelEstrategia]
							descEstrategia := datosEstrategia["descripcionEstrategia"]
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorEstrategia), descEstrategia)
							consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorEstrategia), "N"+fmt.Sprint(contadorEstrategia), stylecontent)
							consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorEstrategia, 70)

							datosComplementarios := datosExcelPlan["datosComplementarios"].(map[string]interface{})
							numeroActividadExcel := datosExcelPlan["numeroActividad"]
							ponderacionActividad := datosComplementarios["Ponderación de la actividad"]
							periodoEjecucion := datosComplementarios["Periodo de ejecución"]
							actividadGeneral := datosComplementarios["Actividad general"]
							tareas := datosComplementarios["Tareas"]
							indicador := datosComplementarios["indicadores"].(map[string]interface{})
							nombreIndicador1 := indicador["Nombre del indicador"]
							nombreIndicador2 := indicador["Nombre indicador 2"]
							formula1 := indicador["Fórmula"]
							formula2 := indicador["Fórmula 2"]
							meta1 := indicador["Meta"]
							meta2 := indicador["Meta 2"]

							consolidadoExcelPlanAnual.SetCellValue(sheetName, "D"+fmt.Sprint(contadorEstrategia), numeroActividadExcel)
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorEstrategia), ponderacionActividad)
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorEstrategia), periodoEjecucion)
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorEstrategia), actividadGeneral)
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorEstrategia), tareas)
							if len(indicador) > 3 {
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorEstrategia), nombreIndicador1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorEstrategia+1), nombreIndicador2)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorEstrategia), formula1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorEstrategia+1), formula2)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorEstrategia), meta1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorEstrategia+1), meta2)
								consolidadoExcelPlanAnual.SetCellStyle(sheetName, "A"+fmt.Sprint(contadorEstrategia+1), "N"+fmt.Sprint(contadorEstrategia+1), stylecontent)
								consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorEstrategia+1, 70)
								contadorEstrategia = contadorEstrategia + 1
								contadorEstrategiaOut = contadorEstrategia
							} else {
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorEstrategia), nombreIndicador1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorEstrategia), formula1)
								consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorEstrategia), meta1)
							}
							consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorEstrategiaIn), "C"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "D"+fmt.Sprint(contadorEstrategiaIn), "D"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorEstrategiaIn), "E"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorEstrategiaIn), "F"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "G"+fmt.Sprint(contadorEstrategiaIn), "G"+fmt.Sprint(contadorEstrategiaOut))
							consolidadoExcelPlanAnual.MergeCell(sheetName, "H"+fmt.Sprint(contadorEstrategiaIn), "H"+fmt.Sprint(contadorEstrategiaOut))

							if contadorEstrategiaIn > contadorEstrategia {
								contadorMeta = contadorEstrategiaIn
							} else {
								contadorMeta = contadorEstrategia
							}
							contadorEstrategia = contadorEstrategia + 1
						}

						contadorMetaOut := contadorMeta
						consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorMetaIn), "B"+fmt.Sprint(contadorMetaOut))
						if contadorMeta >= contadorLineamiento {
							contadorLineamiento = contadorMeta
						}
						contadorMeta = contadorMeta + 1
					}

					contadorLineamientoOut := contadorLineamiento

					consolidadoExcelPlanAnual.MergeCell(sheetName, "A"+fmt.Sprint(contadorLineamientoIn), "A"+fmt.Sprint(contadorLineamientoOut))
					contadorLineamiento = contadorLineamiento + 1
				}
				consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
			}

		}
		consolidadoExcelPlanAnual.SaveAs("plan_anual_.xlsx")

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
