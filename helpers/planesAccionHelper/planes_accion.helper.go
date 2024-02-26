package planesaccionhelper

import (
	"sort"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"

	"github.com/udistrital/planeacion_mid/helpers"
	formulacionhelper "github.com/udistrital/planeacion_mid/helpers/formulacionHelper"
)

func obtenerEstadosSeguimiento() map[string]string {
	defer func() {
		if err := recover(); err != nil {
			panic(map[string]interface{}{"funcion": "obtenerEstadosSeguimiento", "err": "Error obteniendo los estados", "status": "400", "log": err})
		}
	}()
	estados := make(map[string]string)
	var respuestaEstados map[string]interface{}
	var estadoFormulacion []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento", &respuestaEstados); err != nil {
		panic(err)
	}
	helpers.LimpiezaRespuestaRefactor(respuestaEstados, &estadoFormulacion)
	for _, estado := range estadoFormulacion {
		estados[estado["_id"].(string)] = estado["nombre"].(string)
	}
	return estados
}

/*
Retornar planes con:
 - Unidad Académica y / o Administrativa (Nombre)
 - Vigencia
 - Plan de acción (Nombre)
 - Estado del plan de acción (Formulación / Seguimiento)
 - Versión en el caso de Formulación (N/A para los demás casos)
@return Un arreglo con los id's y datos más relevantes
*/
func ObtenerPlanesAccion() (resumenPlanes []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			outputError = map[string]interface{}{
				"funcion": "ObtenerPlanesAccion/" + localError["funcion"].(string),
				"err":     localError["err"],
				"status":  localError["status"],
			}
			panic(outputError)
		}
	}()
	var planesSeguimiento []map[string]interface{}
	planesAvalados := make(map[string]map[string]interface{})
	var resSeguimiento map[string]interface{}

	var planesFormulacion, err = formulacionhelper.ObtenerPlanesFormulacion()
	if err != nil {
		panic(err)
	}

	for _, plan := range planesFormulacion {
		plan["fase"] = "Formulación"
		if plan["estado"] == "Aval" {
			planesAvalados[plan["id"].(string)] = plan
		}
		resumenPlanes = append(resumenPlanes, plan)
	}

	estadosSeguimiento := obtenerEstadosSeguimiento()
	// Obtiene los planes que hay en seguimiento
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento`, &resSeguimiento); err != nil {
		panic(err)
	}

	helpers.LimpiezaRespuestaRefactor(resSeguimiento, &planesSeguimiento)
	for _, planSeguimiento := range planesSeguimiento {
		if planSeguimiento["plan_id"] != nil {
			if plan, existePlan := planesAvalados[planSeguimiento["plan_id"].(string)]; existePlan {
				planNuevo := make(map[string]interface{})
				planNuevo["id"] = plan["id"]
				planNuevo["dependencia_id"] = plan["dependencia_id"]
				planNuevo["dependencia_nombre"] = plan["dependencia_nombre"]
				planNuevo["vigencia_id"] = plan["vigencia_id"]
				planNuevo["vigencia"] = plan["vigencia"]
				planNuevo["nombre"] = plan["nombre"]
				planNuevo["estado_id"] = planSeguimiento["estado_seguimiento_id"]
				planNuevo["estado"] = estadosSeguimiento[planSeguimiento["estado_seguimiento_id"].(string)]
				planNuevo["ultima_modificacion"] = planSeguimiento["fecha_modificacion"]
				planNuevo["fase"] = "Seguimiento"
				// La versión solo se devolverá en caso de que exista, esto solo aplicla en los planes en formulación, por tanto, se debé manejar esto en el cliente
				resumenPlanes = append(resumenPlanes, planNuevo)
			}
		}
	}
	sort.Slice(resumenPlanes, func(i, j int) bool {
		return resumenPlanes[i]["ultima_modificacion"].(string) > resumenPlanes[j]["ultima_modificacion"].(string)
	})
	return resumenPlanes, outputError
}

func ObtenerPlanesDeAccionPorUnidad(unidadID string) (planes []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			localError := err.(map[string]interface{})
			outputError = map[string]interface{}{
				"funcion": "ObtenerPlanesAccion/" + localError["funcion"].(string),
				"err":     localError["err"],
				"status":  localError["status"],
			}
			panic(outputError)
		}
	}()

	if planesAccion, err := ObtenerPlanesAccion(); err != nil {
		panic(err)
	} else {
		for _, plan := range planesAccion {
			if plan["dependencia_id"] == unidadID {
				planes = append(planes, plan)
			}
		}
	}
	return planes, outputError
}
