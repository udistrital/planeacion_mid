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

		encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))

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

				contadorLineamiento := 3
				contadorFactor := 3
				contadorDataGeneral := 3
				unidadNombre = arregloPlanAnual[0]["nombreUnidad"].(string)
				sheetName := "Actividades del plan"
				indexPlan, _ := consolidadoExcelPlanAnual.NewSheet(sheetName)

				if planes == 0 {

					styledefault, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
						Border: []excelize.Border{
							{Type: "right", Color: "ffffff", Style: 1},
							{Type: "left", Color: "ffffff", Style: 1},
							{Type: "top", Color: "ffffff", Style: 1},
							{Type: "bottom", Color: "ffffff", Style: 1},
						},
					})
					consolidadoExcelPlanAnual.InsertCols("Actividades del plan", "A", 1)
					consolidadoExcelPlanAnual.SetColStyle(sheetName, "A:Q", styledefault)
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
				consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 33)
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
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "F2", "Lineaminetos de acción")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "G2", "Estrategias")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "H2", "N°.")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "I2", "Ponderación de la actividad")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "J2", "Periodo de ejecución")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "K2", "Actividad")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "L2", "Tareas")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M1", "Indicador")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "M2", "Nombre")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "N2", "Fórmula")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "O2", "Meta")
				consolidadoExcelPlanAnual.SetCellValue(sheetName, "P2", "Producto esperado")

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
								reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategias), "D"+fmt.Sprint(contadorEstrategias), stylecontent, stylecontentS)

								if k == len(estrategias)-1 {
									contadorEstrategiaPEDOut = contadorMetas
								} else {
									contadorEstrategias = contadorEstrategias + 1
								}
							}

							contadorEstrategias = contadorMetas
							consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorMetaGeneralIn), "C"+fmt.Sprint(contadorMetaGeneralOut))
						}
						contadorMetas = contadorLineamiento
						contadorLineamientoGeneralOut = contadorMetaGeneralOut

						consolidadoExcelPlanAnual.MergeCell(sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralIn), "B"+fmt.Sprint(contadorLineamientoGeneralOut))
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "B"+fmt.Sprint(contadorLineamientoGeneralIn), "B"+fmt.Sprint(contadorLineamientoGeneralOut), styleLineamiento, styleLineamientoSombra)

						contadorLineamiento = contadorLineamientoGeneralOut + 1
					}
					for i := 0; i < len(armoPI); i++ {
						datosArmo := armoPI[i]
						auxFactor := datosArmo["nombreFactor"]
						contadorFactorGeneralIn = contadorFactor
						// cuerpo del excel
						consolidadoExcelPlanAnual.SetCellValue(sheetName, "E"+fmt.Sprint(contadorFactor), auxFactor)
						consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorFactor, 70)

						lineamientos := datosArmo["lineamientos"]
						contadorLineamientos := contadorFactor

						for j := 0; j < len(lineamientos.([]map[string]interface{})); j++ {
							auxLineamiento := lineamientos.([]map[string]interface{})[j]
							contadorLineamientoPIIn = contadorFactor
							consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorLineamientos), auxLineamiento["nombreLineamiento"])

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
								reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategias), "G"+fmt.Sprint(contadorEstrategias), stylecontent, stylecontentS)

								if k == len(estrategiasPI)-1 {
									contadorEstrategiaPIOut = contadorLineamientos
								} else {
									contadorEstrategias = contadorEstrategias + 1
								}
							}

							contadorEstrategias = contadorLineamientos
							consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorLineamientoPIIn), "F"+fmt.Sprint(contadorLineamientoPIOut))
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientoPIIn), "F"+fmt.Sprint(contadorLineamientoPIOut), stylecontent, stylecontentS)

						}
						contadorLineamientos = contadorFactor
						contadorFactorGeneralOut = contadorLineamientoPIOut

						consolidadoExcelPlanAnual.MergeCell(sheetName, "E"+fmt.Sprint(contadorFactorGeneralIn), "E"+fmt.Sprint(contadorFactorGeneralOut))
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactorGeneralIn), "E"+fmt.Sprint(contadorFactorGeneralOut), stylecontent, stylecontentS)

						contadorFactor = contadorFactorGeneralOut + 1

					}
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorDataGeneral), datosExcelPlan["numeroActividad"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Ponderación de la actividad"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Periodo de ejecución"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Actividad general"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Tareas"])
					consolidadoExcelPlanAnual.SetCellValue(sheetName, "P"+fmt.Sprint(contadorDataGeneral), datosComplementarios["Producto esperado "])

					if contadorLineamientoGeneralOut > contadorFactorGeneralOut {
						contadorFactorGeneralOut = contadorLineamientoGeneralOut
						contadorFactor = contadorFactorGeneralOut + 1
					} else if contadorLineamientoGeneralOut < contadorFactorGeneralOut {
						contadorLineamientoGeneralOut = contadorFactorGeneralOut
						contadorLineamiento = contadorLineamientoGeneralOut + 1
					}

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
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentS)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)

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
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(contadorIndicadores), "O"+fmt.Sprint(contadorIndicadores), stylecontentCL, stylecontentCLS)
						contadorIndicadores = contadorIndicadores + 1
					}
					
					if excelPlan+1 != len(arregloPlanAnual) {
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(contadorIndicadores), "O"+fmt.Sprint(contadorIndicadores), stylecontentCL, stylecontentCLS)
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
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactorGeneralOut), "E"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(contadorEstrategiaPIOut), "L"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorIndicadores), stylecontentC, stylecontentCS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(contadorIndicadores), "O"+fmt.Sprint(contadorIndicadores), stylecontentCLD, stylecontentCLDS)
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
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactorGeneralOut), "E"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "K"+fmt.Sprint(contadorEstrategiaPIOut), "L"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorDataGeneral), "J"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(contadorIndicadores), "O"+fmt.Sprint(contadorIndicadores), stylecontentCLD, stylecontentCLDS)
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontentC, stylecontentCS)
					}

					contadorDataGeneral = contadorIndicadores + 1
					contadorLineamiento = contadorIndicadores + 1
					contadorFactor = contadorIndicadores + 1
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
				&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "oneCell", OffsetX: 50}); err != nil {
				fmt.Println(err)
			}
			if err := consolidadoExcelPlanAnual.AddPicture("Identificaciones", "B1", "static/img/UDEscudo2.png",
				&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "absolute", OffsetX: 10}); err != nil {
				fmt.Println(err)
			}

			if len(consolidadoExcelPlanAnual.GetSheetList()) > 1 {
				consolidadoExcelPlanAnual.DeleteSheet("Sheet1")
			}

			consolidadoExcelPlanAnual.SetColWidth("Actividades del plan", "A", "A", 2)
			buf, _ := consolidadoExcelPlanAnual.WriteToBuffer()
			strings.NewReader(buf.String())
			encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))

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

							if idUnidad != planesFilter[planes]["dependencia_id"].(string) {
								if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+planesFilter[planes]["dependencia_id"].(string), &respuestaUnidad); err == nil {
									aux := respuestaUnidad[0]
									dependenciaNombre := aux["DependenciaId"].(map[string]interface{})
									nombreUnidad = dependenciaNombre["Nombre"].(string)
									idUnidad = planesFilter[planes]["dependencia_id"].(string)

								} else {
									panic(map[string]interface{}{"funcion": "GetUnidades", "err": "Error ", "status": "400", "log": err})
								}
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
				styledefault, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
					Border: []excelize.Border{
						{Type: "right", Color: "ffffff", Style: 1},
						{Type: "left", Color: "ffffff", Style: 1},
						{Type: "top", Color: "ffffff", Style: 1},
						{Type: "bottom", Color: "ffffff", Style: 1},
					},
				})

				consolidadoExcelPlanAnual.InsertCols("REPORTE GENERAL", "A", 1)
				consolidadoExcelPlanAnual.SetColStyle(sheetName, "A:Q", styledefault)

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
				Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
				Border: []excelize.Border{
					{Type: "right", Color: "000000", Style: 1},
					{Type: "left", Color: "000000", Style: 1},
					{Type: "top", Color: "000000", Style: 1},
					{Type: "bottom", Color: "000000", Style: 1},
				},
			})
			stylecontentSombra, _ := consolidadoExcelPlanAnual.NewStyle(&excelize.Style{
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
			consolidadoExcelPlanAnual.SetColWidth(sheetName, "B", "B", 33)
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
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorGeneral+3), "Lineaminetos de acción")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "G"+fmt.Sprint(contadorGeneral+3), "Estrategias")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "H"+fmt.Sprint(contadorGeneral+3), "N°.")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "I"+fmt.Sprint(contadorGeneral+3), "Ponderación de la actividad")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "J"+fmt.Sprint(contadorGeneral+3), "Periodo de ejecución")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "K"+fmt.Sprint(contadorGeneral+3), "Actividad")
			consolidadoExcelPlanAnual.SetCellValue(sheetName, "L"+fmt.Sprint(contadorGeneral+3), "Tareas")
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
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorLineamiento), "P"+fmt.Sprint(contadorLineamiento), stylecontent, stylecontentSombra)
					consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorLineamiento, 70)

					metas := datosArmo["meta"]

					contadorMetas := contadorLineamiento

					for j := 0; j < len(metas.([]map[string]interface{})); j++ {
						auxMeta := metas.([]map[string]interface{})[j]
						contadorMetaGeneralIn = contadorLineamiento

						consolidadoExcelPlanAnual.SetCellValue(sheetName, "C"+fmt.Sprint(contadorMetas), auxMeta["nombreMeta"])
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorMetas), "C"+fmt.Sprint(contadorMetas), stylecontent, stylecontentSombra)

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
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategias), "D"+fmt.Sprint(contadorEstrategias), stylecontent, stylecontentSombra)

							if k == len(estrategias)-1 {
								contadorEstrategiaPEDOut = contadorMetas
							} else {
								contadorEstrategias = contadorEstrategias + 1
							}
						}

						contadorEstrategias = contadorMetas

						consolidadoExcelPlanAnual.MergeCell(sheetName, "C"+fmt.Sprint(contadorMetaGeneralIn), "C"+fmt.Sprint(contadorMetaGeneralOut))

					}
					contadorMetas = contadorLineamiento
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
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactor), "P"+fmt.Sprint(contadorFactor), stylecontent, stylecontentSombra)
					consolidadoExcelPlanAnual.SetRowHeight(sheetName, contadorFactor, 70)

					lineamientos := datosArmo["lineamientos"]

					contadorLineamientos := contadorFactor

					for j := 0; j < len(lineamientos.([]map[string]interface{})); j++ {
						auxLineamiento := lineamientos.([]map[string]interface{})[j]
						contadorLineamientoPIIn = contadorFactor

						consolidadoExcelPlanAnual.SetCellValue(sheetName, "F"+fmt.Sprint(contadorLineamientos), auxLineamiento["nombreLineamiento"])
						reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientos), "F"+fmt.Sprint(contadorLineamientos), stylecontent, stylecontentSombra)

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
							reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategias), "G"+fmt.Sprint(contadorEstrategias), stylecontent, stylecontentSombra)

							if k == len(estrategiasPI)-1 {
								contadorEstrategiaPIOut = contadorLineamientos
							} else {
								contadorEstrategias = contadorEstrategias + 1
							}
						}

						contadorEstrategias = contadorLineamientos

						consolidadoExcelPlanAnual.MergeCell(sheetName, "F"+fmt.Sprint(contadorLineamientoPIIn), "F"+fmt.Sprint(contadorLineamientoPIOut))

					}
					contadorLineamientos = contadorFactor
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
					contadorFactor = contadorFactorGeneralOut + 1
				} else if contadorLineamientoGeneralOut < contadorFactorGeneralOut {
					contadorLineamientoGeneralOut = contadorFactorGeneralOut
					contadorLineamiento = contadorLineamientoGeneralOut + 1
				}

				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorLineamiento), "P"+fmt.Sprint(contadorLineamiento), stylecontent, stylecontentSombra)

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
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorMetaGeneralIn), "C"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentSombra)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentSombra)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactorGeneralIn), "E"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentSombra)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentSombra)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentSombra)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentSombra)
				reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorLineamientoGeneralOut), stylecontent, stylecontentSombra)

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

					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "M"+fmt.Sprint(contadorIndicadores), "O"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentSombra)
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
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "C"+fmt.Sprint(contadorMetaGeneralOut), "C"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "D"+fmt.Sprint(contadorEstrategiaPEDOut), "D"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "E"+fmt.Sprint(contadorFactorGeneralOut), "E"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "F"+fmt.Sprint(contadorLineamientoPIOut), "F"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "G"+fmt.Sprint(contadorEstrategiaPIOut), "G"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "H"+fmt.Sprint(contadorDataGeneral), "L"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentSombra)
					reporteshelper.SombrearCeldas(consolidadoExcelPlanAnual, excelPlan, sheetName, "P"+fmt.Sprint(contadorDataGeneral), "P"+fmt.Sprint(contadorIndicadores), stylecontent, stylecontentSombra)
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
				}
				contadorDataGeneral = contadorIndicadores + 1
				contadorLineamiento = contadorIndicadores + 1
				contadorFactor = contadorIndicadores + 1
				consolidadoExcelPlanAnual.SetActiveSheet(indexPlan)
			}
			contadorGeneral = contadorDataGeneral - 1
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
			&excelize.GraphicOptions{ScaleX: 0.1, ScaleY: 0.1, Positioning: "oneCell", OffsetX: 50}); err != nil {
			fmt.Println(err)
		}

		if len(consolidadoExcelPlanAnual.GetSheetList()) > 1 {
			consolidadoExcelPlanAnual.DeleteSheet("Sheet1")
		}

		consolidadoExcelPlanAnual.SetColWidth("REPORTE GENERAL", "A", "A", 2)
		buf, _ := consolidadoExcelPlanAnual.WriteToBuffer()
		strings.NewReader(buf.String())
		encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))

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
	// var docentesGeneral map[string]interface{}
	docentesPregrado := make(map[string]interface{})
	docentesPosgrado := make(map[string]interface{})
	// var arrDataDocentes []map[string]interface{}
	nombre := c.Ctx.Input.Param(":nombre")

	docentesPregrado["tco"] = 0
	docentesPregrado["mto"] = 0
	docentesPregrado["hch"] = 0
	docentesPregrado["hcp"] = 0
	docentesPregrado["valor"] = 0

	docentesPosgrado["hch"] = 0
	docentesPosgrado["hcp"] = 0
	docentesPosgrado["valor"] = 0

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	// Comentariado temporalmente por no uso de docentes
	// primaServicios := 0
	// primaNavidad := 0
	// primaVacaciones := 0
	// bonificacion := 0
	// interesesCesantias := 0
	// cesantiasPublicas := 0
	// cesantiasPrivadas := 0
	// salud := 0
	// pensionesPublicas := 0
	// pensionesPrivadas := 0
	// arl := 0
	// caja := 0
	// icbf := 0

	necesidadesExcel := excelize.NewFile()
	stylecontent, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	styletitles, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Family: "Arial", Size: 26, Color: "000000"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"F2F2F2"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	stylesubtitles, _ := necesidadesExcel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Bold: true, Family: "Arial", Size: 20, Color: "000000"},
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

	necesidadesExcel.NewSheet("Necesidades")
	necesidadesExcel.DeleteSheet("Sheet1")

	necesidadesExcel.MergeCell("Necesidades", "A1", "F1")
	necesidadesExcel.MergeCell("Necesidades", "A1", "A2")
	necesidadesExcel.MergeCell("Necesidades", "A3", "F3")
	necesidadesExcel.MergeCell("Necesidades", "A3", "A4")

	necesidadesExcel.SetColWidth("Necesidades", "A", "c", 30)
	necesidadesExcel.SetColWidth("Necesidades", "D", "D", 50)
	necesidadesExcel.SetColWidth("Necesidades", "E", "F", 20)

	necesidadesExcel.SetCellValue("Necesidades", "A1", "Necesidades Presupuestales")
	necesidadesExcel.SetCellStyle("Necesidades", "A1", "F1", styletitles)

	necesidadesExcel.SetCellValue("Necesidades", "A3", "Identificación de recursos:")
	necesidadesExcel.SetCellStyle("Necesidades", "A3", "F3", stylesubtitles)

	necesidadesExcel.SetCellStyle("Necesidades", "A200", "F200", stylecontent)

	necesidadesExcel.SetCellValue("Necesidades", "A5", "Código del rubro")
	necesidadesExcel.SetCellValue("Necesidades", "B5", "Nombre del rubro")
	necesidadesExcel.SetCellValue("Necesidades", "C5", "Valor")
	necesidadesExcel.SetCellValue("Necesidades", "D5", "Dependencias")
	necesidadesExcel.SetCellStyle("Necesidades", "A5", "D5", stylehead)
	necesidadesExcel.SetRowHeight("Necesidades", 5, 35)
	contador := 6
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",nombre:"+nombre, &respuesta); err == nil {
		helpers.LimpiezaRespuestaRefactor(respuesta, &planes)
		for i := 0; i < len(planes); i++ {
			flag := true
			// var docentes map[string]interface{}
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
					}
					if strings.Contains(strings.ToLower(identificacion["nombre"].(string)), "contratista") && flag {
						if identificacion["dato"] != nil {
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
					}
					// comentariado temporalmente por no uso de docentes
					/*else if strings.Contains(strings.ToLower(identificacion["nombre"].(string)), "docente") {
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
					}*/
				}
				for i := 0; i < len(recursos); i++ {
					var aux bool
					var aux1 []string
					if len(recursosGeneral) == 0 {
						recursosGeneral = append(recursosGeneral, recursos[i])
						aux1 = append(aux1, dependencia_nombre)
						recursosGeneral[len(recursosGeneral)-1]["unidades"] = aux1
						unidades_total = append(unidades_total, dependencia_nombre)
					} else {
						for j := 0; j < len(recursosGeneral); j++ {
							if recursosGeneral[j]["codigo"] == recursos[i]["codigo"] {
								flag := false
								for k := 0; k < len(recursosGeneral[j]["unidades"].([]string)); k++ {
									aux2 := recursosGeneral[j]["unidades"].([]string)
									if aux2[k] == dependencia_nombre {
										flag = true
									}
								}
								if !flag {
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
									var auxValor int
									var auxValor2 int
									if fmt.Sprint(reflect.TypeOf(recursosGeneral[j]["valor"])) == "int" {
										auxValor = recursosGeneral[j]["valor"].(int)
									} else {
										strValor := strings.TrimLeft(recursosGeneral[j]["valor"].(string), "$")
										strValor = strings.ReplaceAll(strValor, ",", "")
										arrValor := strings.Split(strValor, ".")
										aux1, err := strconv.Atoi(arrValor[0])
										if err == nil {
											auxValor = aux1
										}
									}

									if fmt.Sprint(reflect.TypeOf(recursos[i]["valor"])) == "int" {
										auxValor = recursos[i]["valor"].(int)
									} else {
										strValor2 := strings.TrimLeft(recursos[i]["valor"].(string), "$")
										strValor2 = strings.ReplaceAll(strValor2, ",", "")
										arrValor2 := strings.Split(strValor2, ".")
										aux2, err := strconv.Atoi(arrValor2[0])
										if err == nil {
											auxValor2 = aux2
										}
									}
									recursosGeneral[j]["valor"] = auxValor + auxValor2
								} else {
									recursosGeneral[j]["valor"] = recursos[i]["valor"]
								}
								aux = true
								break
							} else {
								aux = false
							}
						}
						if !aux {
							flag := false
							recursosGeneral = append(recursosGeneral, recursos[i])
							aux1 = append(aux1, dependencia_nombre)
							recursosGeneral[len(recursosGeneral)-1]["unidades"] = aux1
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
						rubrosGeneral = append(rubrosGeneral, rubros[i])
						aux1 = append(aux1, dependencia_nombre)
						rubrosGeneral[len(rubrosGeneral)-1]["unidades"] = aux1
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
									rubrosGeneral[j]["totalInc"] = auxValor + auxValor2
								} else {
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
							recursosGeneral = append(recursosGeneral, recursos[i])
							aux1 = append(aux1, dependencia_nombre)
							recursosGeneral[len(recursosGeneral)-1]["unidades"] = aux1
							for k := 0; k < len(unidades_total); k++ {
								if unidades_total[k] == dependencia_nombre {
									flag = true
								}
							}
							if !flag {
								unidades_total = append(unidades_total, dependencia_nombre)
							}
						}
						if !aux {
							flag := false
							rubrosGeneral = append(rubrosGeneral, rubros[i])
							aux1 = append(aux1, dependencia_nombre)
							rubrosGeneral[len(rubrosGeneral)-1]["unidades"] = aux1
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
				// comentariado temporalmente por no uso de docentes
				/*if len(docentes) > 0 {
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
				}*/

				// comentariado temporalmente por no uso de docentes
				/*if docentes["rubros"] != nil {
					var aux bool
					var respuestaRubro map[string]interface{}
					rubros := docentes["rubros"].([]map[string]interface{})
					for i := 0; i < len(rubros); i++ {
						if rubros[i]["rubro"] != "" {
							for j := 0; j < len(recursosGeneral); j++ {
								if recursosGeneral[j]["codigo"] == rubros[i]["rubro"] {
									aux = true
									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "servicio") {
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
									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "navidad") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "vacaciones") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "bonificacion") || strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "bonificación") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "interes") && strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "cesantía") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "cesantía") && strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "público") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "cesantía") && strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "privado") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "salud") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "pension") && strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "público") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "pension") && strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "privado") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "arl") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "ccf") {
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

									if strings.Contains(strings.ToLower(rubros[i]["categoria"].(string)), "icbf") {
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
							if !aux {
								rubro := make(map[string]interface{})
								if err := request.GetJson("http://"+beego.AppConfig.String("PlanCuentasService")+"/arbol_rubro/"+rubros[i]["rubro"].(string), &respuestaRubro); err == nil {
									aux := respuestaRubro["Body"].(map[string]interface{})
									rubro["codigo"] = aux["Codigo"]
									rubro["nombre"] = aux["Nombre"]
									rubro["categoria"] = rubros[i]["categoria"]

									if rubro["categoria"] != nil {
										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(rubro["categoria"].(string)), "servicio") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(rubro["categoria"].(string)), "navidad") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(rubro["categoria"].(string)), "vacaciones") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "bonificacion") || strings.Contains(strings.ToLower(rubro["categoria"].(string)), "bonificación") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "interes") && strings.Contains(strings.ToLower(rubro["categoria"].(string)), "cesantía") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "cesantía") && strings.Contains(strings.ToLower(rubro["categoria"].(string)), "público") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "cesantía") && strings.Contains(strings.ToLower(rubro["categoria"].(string)), "privado") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "salud") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "pension") && strings.Contains(strings.ToLower(rubro["categoria"].(string)), "público") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "pension") && strings.Contains(strings.ToLower(rubro["categoria"].(string)), "privado") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "arl") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "ccf") {
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

										if strings.Contains(strings.ToLower(rubro["categoria"].(string)), "icbf") {
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
				}*/
			}
		}
		/*
			// for i := 0; i < len(recursosGeneral); i++ {
			// 	if recursosGeneral[i]["categoria"] != nil {
			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "servicio") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + primaServicios
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = primaServicios
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "navidad") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + primaNavidad
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = primaNavidad
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "prima") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "vacaciones") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + primaVacaciones
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = primaVacaciones
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "bonificacion") || strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "bonificación") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + bonificacion
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = bonificacion
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "interes") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "cesantía") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + interesesCesantias
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = interesesCesantias
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "cesantía") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "público") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + cesantiasPublicas
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = cesantiasPublicas
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "cesantía") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "privado") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + cesantiasPrivadas
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = cesantiasPrivadas
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "salud") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + salud
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = salud
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "pension") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "público") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + pensionesPublicas
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = pensionesPublicas
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "pension") && strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "privado") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + pensionesPrivadas
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = pensionesPrivadas
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "arl") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + arl
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = arl
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "ccf") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + caja
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = caja
			// 			}
			// 		}

			// 		if strings.Contains(strings.ToLower(recursosGeneral[i]["categoria"].(string)), "icbf") {
			// 			if recursosGeneral[i]["valor"] != nil {
			// 				strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
			// 				strValor = strings.ReplaceAll(strValor, ",", "")
			// 				arrValor := strings.Split(strValor, ".")
			// 				auxValor, err := strconv.Atoi(arrValor[0])
			// 				if err == nil {
			// 					recursosGeneral[i]["valor"] = auxValor + icbf
			// 				}
			// 			} else {
			// 				recursosGeneral[i]["valor"] = icbf
			// 			}
			// 		}
			// 	}

			// }
		*/
		//Completado de tablas

		for i := 0; i < len(recursosGeneral); i++ {
			unidades := ""
			necesidadesExcel.SetCellValue("Necesidades", "A"+fmt.Sprint(contador), recursosGeneral[i]["codigo"])
			if recursosGeneral[i]["Nombre"] != nil {
				necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), recursosGeneral[i]["Nombre"])
				if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
					necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), recursosGeneral[i]["valor"])
				} else {
					strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
					strValor = strings.ReplaceAll(strValor, ",", "")
					arrValor := strings.Split(strValor, ".")
					auxValor, err := strconv.Atoi(arrValor[0])
					if err == nil {
						necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), auxValor)
					}
				}
				if recursosGeneral[i]["unidades"] != nil {
					aux2 := recursosGeneral[i]["unidades"].([]string)
					for j := 0; j < len(aux2); j++ {
						unidades = unidades + aux2[j] + ", "
					}
					unidades = strings.TrimRight(unidades, ", ")
					necesidadesExcel.SetCellValue("Necesidades", "D"+fmt.Sprint(contador), unidades)
				}

			} else {
				necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), recursosGeneral[i]["nombre"])
				if fmt.Sprint(reflect.TypeOf(recursosGeneral[i]["valor"])) == "int" {
					necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), recursosGeneral[i]["valor"])
				} else {
					strValor := strings.TrimLeft(recursosGeneral[i]["valor"].(string), "$")
					strValor = strings.ReplaceAll(strValor, ",", "")
					arrValor := strings.Split(strValor, ".")
					auxValor, err := strconv.Atoi(arrValor[0])
					if err == nil {
						necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), auxValor)
					}
				}
			}
			necesidadesExcel.SetCellStyle("Necesidades", "A"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent)
			contador++
		}

		for i := 0; i < len(rubrosGeneral); i++ {
			unidades := ""
			necesidadesExcel.SetCellValue("Necesidades", "A"+fmt.Sprint(contador), rubrosGeneral[i]["rubro"])
			if rubrosGeneral[i]["rubroNombre"] != nil {
				necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), rubrosGeneral[i]["rubroNombre"])
				if fmt.Sprint(reflect.TypeOf(rubrosGeneral[i]["totalInc"])) == "float64" {
					necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), rubrosGeneral[i]["totalInc"])
				} else {
					strValor := strings.TrimLeft(rubrosGeneral[i]["totalInc"].(string), "$")
					strValor = strings.ReplaceAll(strValor, ",", "")
					auxValor, err := strconv.ParseFloat(strValor, 64)
					if err == nil {
						necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), auxValor)
					}
				}
				if rubrosGeneral[i]["unidades"] != nil {
					aux2 := rubrosGeneral[i]["unidades"].([]string)
					for j := 0; j < len(aux2); j++ {
						unidades = unidades + aux2[j] + ", "
					}
					unidades = strings.TrimRight(unidades, ", ")
					necesidadesExcel.SetCellValue("Necesidades", "D"+fmt.Sprint(contador), unidades)
				}

			} else {
				necesidadesExcel.SetCellValue("Necesidades", "B"+fmt.Sprint(contador), rubrosGeneral[i]["rubroNombre"])
				if fmt.Sprint(reflect.TypeOf(rubrosGeneral[i]["totalInc"])) == "float64" {
					necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), rubrosGeneral[i]["totalInc"])
				} else {
					if rubrosGeneral[i]["totalInc"] != nil {
						strValor := strings.TrimLeft(rubrosGeneral[i]["totalInc"].(string), "$")
						strValor = strings.ReplaceAll(strValor, ",", "")
						auxValor, err := strconv.ParseFloat(strValor, 64)
						if err == nil {
							necesidadesExcel.SetCellValue("Necesidades", "C"+fmt.Sprint(contador), auxValor)
						}
					} else {
						contador--
					}
				}
			}
			necesidadesExcel.SetCellStyle("Necesidades", "A"+fmt.Sprint(contador), "D"+fmt.Sprint(contador), stylecontent)
			contador++
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

		contador = 1

		necesidadesExcel.NewSheet("Total Unidades")

		necesidadesExcel.MergeCell("Total Unidades", "A", "B")
		necesidadesExcel.MergeCell("Total Unidades", "A1", "A2")

		necesidadesExcel.SetColWidth("Total Unidades", "A", "B", 30)

		necesidadesExcel.MergeCell("Total Unidades", "A"+fmt.Sprint(contador), "B"+fmt.Sprint(contador))
		necesidadesExcel.MergeCell("Total Unidades", "A"+fmt.Sprint(contador), "A"+fmt.Sprint(contador+1))
		necesidadesExcel.SetCellValue("Total Unidades", "A"+fmt.Sprint(contador), "Total de unidades generadas:")
		necesidadesExcel.SetCellStyle("Total Unidades", "A"+fmt.Sprint(contador), "B"+fmt.Sprint(contador), stylesubtitles)
		necesidadesExcel.SetCellStyle("Total Unidades", "A"+fmt.Sprint(contador+1), "B"+fmt.Sprint(contador+1), stylesubtitles)
		necesidadesExcel.SetCellStyle("Total Unidades", "A200", "F200", stylecontent)
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
		// comentariado temporalmente por no uso de docentes
		/*necesidadesExcel.MergeCell("Necesidades", "A"+fmt.Sprint(contador), "F"+fmt.Sprint(contador))
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
			necesidadesExcel.SetCellValue("Necesidades", "F"+fmt.Sprint(contador), arrDataDocentes[i]["valorPre"])
			necesidadesExcel.SetCellValue("Necesidades", "G"+fmt.Sprint(contador), arrDataDocentes[i]["hchPos"])
			necesidadesExcel.SetCellValue("Necesidades", "H"+fmt.Sprint(contador), arrDataDocentes[i]["hcpPos"])
			necesidadesExcel.SetCellValue("Necesidades", "I"+fmt.Sprint(contador), arrDataDocentes[i]["valorPos"])
			necesidadesExcel.SetCellStyle("Necesidades", "A"+fmt.Sprint(contador), "I"+fmt.Sprint(contador), stylecontent)
			contador++
		}*/

		buf, _ := necesidadesExcel.WriteToBuffer()
		strings.NewReader(buf.String())
		encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))
		dataSend := make(map[string]interface{})
		dataSend["generalData"] = arregloInfoReportes
		dataSend["excelB64"] = encoded

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": dataSend}

	} else {
		panic(err)
	}

	c.ServeJSON()
}
