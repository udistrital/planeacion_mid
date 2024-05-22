package services

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	evaluacionhelper "github.com/udistrital/planeacion_mid/helpers/evaluacionHelper"
	"github.com/udistrital/utils_oas/request"
)

const (
	ABREVIACION_AVALADO_PARA_SEGUIMIENTO string = "AV"
	ABREVIACION_SEGUIMIENTO_PLAN_ACCION  string = "S_SP"
)

func GetPlanesPeriodo(vigencia string, unidad string) (interface{}, error) {

	if len(vigencia) == 0 || len(unidad) == 0 {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: ")
	}
	if respuesta, err := evaluacionhelper.GetPlanesPeriodo(unidad, vigencia); err == nil {
		return respuesta, nil
	} else {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: ")
	}
}

func GetEvaluacion(vigencia string, plan string, periodoId string) (interface{}, error) {

	var evaluacion []map[string]interface{}

	if len(vigencia) == 0 || len(plan) == 0 || len(periodoId) == 0 {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	if len(plan) == 0 {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	if len(periodoId) == 0 {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}

	trimestres := evaluacionhelper.GetPeriodos(vigencia)
	if len(trimestres) == 0 {
		return nil, nil
	} else {
		i := 0
		for index, periodo := range trimestres {
			if periodo["_id"] == periodoId {
				i = index
				break
			}
		}
		evaluacion = evaluacionhelper.GetEvaluacion(plan, trimestres, i)
		return evaluacion, nil
	}
}

func Unidades(plan string, vigencia string) (interface{}, error) {

	if nombrePlan, err := url.QueryUnescape(plan); err == nil {
		if data, err := evaluacionhelper.GetUnidadesPorPlanYVigencia(nombrePlan, vigencia); err == nil {
			return data, nil
		} else {
			return nil, errors.New("Error obteniendo las unidades del plan y la vigencia dados ")
		}
	} else {
		return nil, errors.New("Error obteniendo las unidades del plan y la vigencia dados " + err.Error())

	}
}

func Avances(plan string, vigencia string, unidad string) (interface{}, error) {

	if nombrePlan, err1 := url.QueryUnescape(plan); err1 == nil {
		if data, err2 := evaluacionhelper.GetAvances(nombrePlan, vigencia, unidad); err2 == nil {
			return data, nil
		} else {
			return nil, errors.New("Error obteniendo los avances ")
		}
	} else {
		return nil, errors.New("Error obteniendo los avances " + err1.Error())
	}
}

func PlanesAEvaluar() (planes []string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}
	}()
	var respuestaEstado map[string]interface{}
	var respuestaTipoSeguimiento map[string]interface{}
	var respuestaSeguimiento map[string]interface{}

	var estadoSeguimiento []map[string]interface{}
	var tipoSeguimiento []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_AVALADO_PARA_SEGUIMIENTO, &respuestaEstado); err != nil {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	request.LimpiezaRespuestaRefactor(respuestaEstado, &estadoSeguimiento)
	fmt.Println("Entered PlanesAevaluar", planes)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_SEGUIMIENTO_PLAN_ACCION, &respuestaTipoSeguimiento); err != nil {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	request.LimpiezaRespuestaRefactor(respuestaTipoSeguimiento, &tipoSeguimiento)
	fmt.Println("Entered PlanesAevaluar", planes)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:`+tipoSeguimiento[0]["_id"].(string)+`,estado_seguimiento_id:`+estadoSeguimiento[0]["_id"].(string)+`,activo:true`, &respuestaSeguimiento); err == nil {
		var seguimientos []map[string]interface{}
		request.LimpiezaRespuestaRefactor(respuestaSeguimiento, &seguimientos)
		for _, seguimiento := range seguimientos {
			// Esta en los planes que ya se trajeron?
			var respuestaPlan map[string]interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan/`+seguimiento["plan_id"].(string), &respuestaPlan); err == nil {
				var plan map[string]interface{}
				existeNombrePlan := false
				request.LimpiezaRespuestaRefactor(respuestaPlan, &plan)
				fmt.Println("Entered PlanesAevaluar REQUEST LIMPIUEZA RESPUESTAPLAN", seguimiento["plan_id"].(string))
				for _, nombre := range planes {
					if nombre == plan["nombre"].(string) {
						existeNombrePlan = true
					}
				}
				if !existeNombrePlan {
					planes = append(planes, plan["nombre"].(string))
				}

			} else {
				outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
			}
		}
	} else {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")

	}
	fmt.Println("Entered PlanesAevaluar", planes)
	fmt.Println("Entered PlanesAevaluar", outputError)
	return planes, outputError
}
