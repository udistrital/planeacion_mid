package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"],
        beego.ControllerComments{
            Method: "GetArbol",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"],
        beego.ControllerComments{
            Method: "DeletePlan",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"],
        beego.ControllerComments{
            Method: "GetFormato",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormatoController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "ActualizarActividad",
            Router: "/actualizar_actividad/:id/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "ClonarFormato",
            Router: "/clonar_formato/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "DeleteActividad",
            Router: "/delete_actividad/:id/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "DeleteIdentificacion",
            Router: "/delete_identificacion/:id/:idTipo/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GetAllActividades",
            Router: "/get_all_actividades/:id/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GetAllIdentificacion",
            Router: "/get_all_identificacion/:id/:idTipo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GetArbolArmonizacion",
            Router: "/get_arbol_armonizacion/:id/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GetPlan",
            Router: "/get_plan/:id/:index",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GetPlanVersiones",
            Router: "/get_plan_versiones/:unidad/:vigencia/:nombre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GetRubros",
            Router: "/get_rubros",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GetUnidades",
            Router: "/get_unidades",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GuardarActividad",
            Router: "/guardar_actividad/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "GuardarIdentificacion",
            Router: "/guardar_identificacion/:id/:idTipo",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "PonderacionActividades",
            Router: "/ponderacion_actividades/:plan",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "VersionarPlan",
            Router: "/versionar_plan/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "CrearReportes",
            Router: "/crear_reportes/:plan/:tipo",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GetActividadesGenerales",
            Router: "/get_actividades/:plan_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GetAvanceIndicador",
            Router: "/get_avance",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GetDataActividad",
            Router: "/get_data/:plan_id/:index",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GetIndicadores",
            Router: "/get_indicadores/:plan_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GetPeriodos",
            Router: "/get_periodos/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GetSeguimiento",
            Router: "/get_seguimiento/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarSeguimiento",
            Router: "/guardar_seguimiento/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "HabilitarReportes",
            Router: "/habilitar_reportes/:periodo",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
