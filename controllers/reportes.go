package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
	// "reflect"
	"strconv"
	// "encoding/base64"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"github.com/xuri/excelize/v2"
	reporteshelper "github.com/udistrital/planeacion_mid/helpers/reportesHelper"
	// formulacionhelper "github.com/udistrital/planeacion_mid/helpers/formulacionHelper"
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
		CreateExcel(consolidadoExcel, "Consolidado Presupuestal.xlsx")
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
	var respuestaGeneral map[string]interface{}
	var planesFilterGeneral []map[string]interface{}
	var respuesta map[string]interface{}
	var planesFilter []map[string]interface{}
	var res map[string]interface{}
	var resName map[string]interface{}
	var parametro_periodo_name []map[string]interface{}
	var resPeriodo map[string]interface{}
	var resPeriodos []map[string]interface{}
	var resSub map[string]interface{}
	var resArmo map[string]interface{}
	var resEstrategia map[string]interface{}
	var resMeta map[string]interface{}
	var resLineamiento map[string]interface{}
	var resPlan map[string]interface{}
	var hijosArmo []map[string]interface{}
	var estrategiaData []map[string]interface{}
	var metaData []map[string]interface{}
	var LineamientoData []map[string]interface{}
	var planData []map[string]interface{}
	var subgrupos []map[string]interface{}
	var hijos []map[string]interface{}
	var plan_id string
	var actividadName string
	var subgrupoDetalle map[string]interface{}
	var arregloPlanAnual []map[string]interface{}
	var arregloEstrategia []map[string]interface{}
	var arregloMetaEst []map[string]interface{}
	var arregloLineamieto []map[string]interface{}
	var planAnual []map[string]interface{}
	var respuestaDetalle map[string]interface{}
	var respuestaLimpiaDetalle []map[string]interface{}
	var seguimientoPeriodo map[string]interface{}
	var dato map[string]interface{}
	// var datoArmo map[string]interface{}
	var nombrePeriodo string
	var ponderacionActividades string
	var unidadId string
	var nombrePlanDesarrollo string
	var planAnualGeneral []map[string]interface{}

	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	if body["unidad_id"].(string) == ""{
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string), &respuestaGeneral); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuestaGeneral, &planesFilterGeneral)

			for datoUnidad := 0; datoUnidad < len(planesFilterGeneral); datoUnidad ++ {
				planesUnidad := planesFilterGeneral[datoUnidad]
				unidadId = planesUnidad["dependencia_id"].(string)

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",dependencia_id:"+unidadId, &respuesta); err == nil {
					helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
					for planes := 0; planes < len(planesFilter); planes ++ {
						planesFilterData := planesFilter[planes]
						fmt.Println(planesFilter)
						plan_id = planesFilterData["_id"].(string)
		
						if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id, &res); err == nil {
							helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
					
							for i := 0; i < len(subgrupos); i++ {
								if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
									actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
									for j := 0; j < len(actividades); j++ {
										actividad := actividades[j]
										actividadName = actividad["dato"].(string)
										index := fmt.Sprint(actividad["index"])
		
										if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+plan_id, &resPeriodo); err == nil {
											helpers.LimpiezaRespuestaRefactor(resPeriodo, &resPeriodos)
											for p:=0 ; p < len(resPeriodos); p++ {
												seguimientoPeriodo = resPeriodos[p]
											}
											seguimientoPeriodoStr := seguimientoPeriodo["dato"].(string)
											json.Unmarshal([]byte(seguimientoPeriodoStr), &dato)
											periodId := seguimientoPeriodo["periodo_id"].(string)
											if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodId, &resName); err == nil {
												helpers.LimpiezaRespuestaRefactor(resName, &parametro_periodo_name)
												paramIdlenName := parametro_periodo_name[0]
												paramIdName := paramIdlenName["ParametroId"].(map[string]interface{})
												nombrePeriodo = paramIdName["CodigoAbreviacion"].(string)
											}else{
												c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
												c.Abort("400")
											}
											
										}else{
											c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
											c.Abort("400")
										}
		
										if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &resSub); err == nil {
											helpers.LimpiezaRespuestaRefactor(resSub, &hijos)
											data := reporteshelper.GetDataSubgrupos(hijos, index)
											generalData := make(map[string]interface{})
		
											for h := 0; h < len(hijos); h++ {
												if strings.Contains(strings.ToUpper(hijos[h]["nombre"].(string)), "PONDERACIÓN") && strings.Contains(strings.ToUpper(hijos[h]["nombre"].(string)), "ACTIVIDAD") || strings.Contains(strings.ToUpper(hijos[h]["nombre"].(string)), "PONDERACIÓN") {
													if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijos[h]["_id"].(string), &respuestaDetalle); err == nil {
									
														helpers.LimpiezaRespuestaRefactor(respuestaDetalle, &respuestaLimpiaDetalle)
														subgrupoDetalle = respuestaLimpiaDetalle[0]
									
														if subgrupoDetalle["dato_plan"] != nil {
															datoPlan := make(map[string]map[string]interface{})
															json.Unmarshal([]byte(subgrupoDetalle["dato_plan"].(string)), &datoPlan)
		
															if datoPlan[strconv.Itoa(j+1)]["activo"] != false {
																ponderacionActividades = fmt.Sprint(datoPlan[strconv.Itoa(j+1)]["dato"])
															}
														}
									
													} else {
														panic(map[string]interface{}{"funcion": "PonderacionActividades", "err": "Error subgrupo_detalle plan \"plan\"", "status": "400", "log": err})
													}
												}
											}
		
											generalData["nombreActividad"] = actividadName
											generalData["numeroActividad"] = index
											generalData["ponderacion"] = ponderacionActividades
											generalData["periodo_ejecucion"] = nombrePeriodo
											generalData["hijos"] = data
											generalData["unidad_id"] = unidadId

											arregloPlanAnual = append(arregloPlanAnual, generalData)
											// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
										} else {
											c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
											c.Abort("400")
										}
		
		
									}
									// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": arregloPlanAnual}
					
									break
								}
							}
						} else {
							c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
							c.Abort("400")
						}
		
						planAnual = append(planAnual, arregloPlanAnual...)
					}
					// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": planAnual}
				} else {
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					c.Abort("400")
				}

				planAnualGeneral = append(planAnualGeneral, planAnual...)
			}
			
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": planAnualGeneral}

		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}

	}else if body["unidad_id"].(string) != ""{
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan?query=activo:true,tipo_plan_id:"+body["tipo_plan_id"].(string)+",vigencia:"+body["vigencia"].(string)+",estado_plan_id:"+body["estado_plan_id"].(string)+",dependencia_id:"+body["unidad_id"].(string), &respuesta); err == nil {
			helpers.LimpiezaRespuestaRefactor(respuesta, &planesFilter)
			for planes := 0; planes < len(planesFilter); planes ++ {
				planesFilterData := planesFilter[planes]
				fmt.Println(planesFilter)
				plan_id = planesFilterData["_id"].(string)

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+plan_id, &res); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &subgrupos)
					
					for i := 0; i < len(subgrupos); i++ {
						if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {
							actividades := reporteshelper.GetActividades(subgrupos[i]["_id"].(string))
							for j := 0; j < len(actividades); j++ {
								actividad := actividades[j]
								actividadName = actividad["dato"].(string)
								index := fmt.Sprint(actividad["index"])
								// datosArmonizacion := make(map[string]interface{})


								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &resArmo); err == nil {
									helpers.LimpiezaRespuestaRefactor(resArmo, &hijosArmo)
									reporteshelper.Limpia()
									tree := reporteshelper.BuildTreeFa(hijosArmo, index)
									treeArmo := tree[2]
									armonizacionTercer := treeArmo[0]
									armonizacionTercerNivel := armonizacionTercer["armo"]
									if armonizacionTercerNivel != nil {
										// fmt.Println(armonizacionTercerNivel, reflect.TypeOf(armonizacionTercerNivel))
										delimitador := ","
										output := strings.Split(armonizacionTercerNivel.(string), delimitador)
										fmt.Println(len(output))
										estrategiaDesc := make(map[string]interface{})
										metaEstrategica := make(map[string]interface{})
										lineamientoDesc := make(map[string]interface{})
										for estrategiaArmo := 0; estrategiaArmo < len(output); estrategiaArmo++{
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

														if metaDataStr != "[]"{
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
	
																// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": lineamientoEst}
	
			
															} else {
																c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
																c.Abort("400")
															}
														}else{
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

												}else{
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
											// estrategiaDesc["descripcion"] = estrategiaData
										}
										arregloEstrategia = append(arregloEstrategia, estrategiaDesc)
										metaEstrategica["estrategias"] = arregloEstrategia
										arregloMetaEst = append(arregloMetaEst, metaEstrategica)
										lineamientoDesc["meta"] = arregloMetaEst
										arregloLineamieto = append(arregloLineamieto, lineamientoDesc)

										// datosArmonizacion["armonizacion"] = arregloLineamieto

										// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": arregloLineamieto}
									}

								} else {
									c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
									c.Abort("400")
								}


								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+plan_id, &resPeriodo); err == nil {
									helpers.LimpiezaRespuestaRefactor(resPeriodo, &resPeriodos)
									for p:=0 ; p < len(resPeriodos); p++ {
										seguimientoPeriodo = resPeriodos[p]
									}
									seguimientoPeriodoStr := seguimientoPeriodo["dato"].(string)
									json.Unmarshal([]byte(seguimientoPeriodoStr), &dato)
									periodId := seguimientoPeriodo["periodo_id"].(string)
									if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodId, &resName); err == nil {
										helpers.LimpiezaRespuestaRefactor(resName, &parametro_periodo_name)
										paramIdlenName := parametro_periodo_name[0]
										paramIdName := paramIdlenName["ParametroId"].(map[string]interface{})
										nombrePeriodo = paramIdName["CodigoAbreviacion"].(string)
									}else{
										c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
										c.Abort("400")
									}
									
								}else{
									c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
									c.Abort("400")
								}

								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+plan_id, &resSub); err == nil {
									helpers.LimpiezaRespuestaRefactor(resSub, &hijos)
									data := reporteshelper.GetDataSubgrupos(hijos, index)
									generalData := make(map[string]interface{})

									for h := 0; h < len(hijos); h++ {
										if strings.Contains(strings.ToUpper(hijos[h]["nombre"].(string)), "PONDERACIÓN") && strings.Contains(strings.ToUpper(hijos[h]["nombre"].(string)), "ACTIVIDAD") || strings.Contains(strings.ToUpper(hijos[h]["nombre"].(string)), "PONDERACIÓN") {
											if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijos[h]["_id"].(string), &respuestaDetalle); err == nil {
							
												helpers.LimpiezaRespuestaRefactor(respuestaDetalle, &respuestaLimpiaDetalle)
												subgrupoDetalle = respuestaLimpiaDetalle[0]
							
												if subgrupoDetalle["dato_plan"] != nil {
													datoPlan := make(map[string]map[string]interface{})
													json.Unmarshal([]byte(subgrupoDetalle["dato_plan"].(string)), &datoPlan)

													if datoPlan[strconv.Itoa(j+1)]["activo"] != false {
														ponderacionActividades = fmt.Sprint(datoPlan[strconv.Itoa(j+1)]["dato"])
													}
												}
							
											} else {
												panic(map[string]interface{}{"funcion": "PonderacionActividades", "err": "Error subgrupo_detalle plan \"plan\"", "status": "400", "log": err})
											}
										}
									}

									generalData["nombreActividad"] = actividadName
									generalData["numeroActividad"] = index
									generalData["ponderacion"] = ponderacionActividades
									generalData["periodo_ejecucion"] = nombrePeriodo
									generalData["hijos"] = data
									generalData["datosArmonizacion"] = arregloLineamieto

									// arregloPlanAnual = append(arregloPlanAnual, generalData)
									arregloPlanAnual = append(arregloPlanAnual, generalData)
									// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
								} else {
									c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
									c.Abort("400")
								}


							}
							// c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": arregloPlanAnual}
			
							break
						}
					}
				} else {
					c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
					c.Abort("400")
				}

				planAnual = append(planAnual, arregloPlanAnual...)
			}
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Successful", "Data": planAnual}

		} else {
			c.Data["json"] = map[string]interface{}{"Code": "400", "Body": err, "Type": "error"}
			c.Abort("400")
		}
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
