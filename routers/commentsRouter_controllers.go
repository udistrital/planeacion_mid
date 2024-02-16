package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
            Method: "ActivarNodo",
            Router: "/activar_nodo/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"],
        beego.ControllerComments{
            Method: "ActivarPlan",
            Router: "/activar_plan/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"],
        beego.ControllerComments{
            Method: "DeleteNodo",
            Router: "/desactivar_nodo/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ArbolController"],
        beego.ControllerComments{
            Method: "DeletePlan",
            Router: "/desactivar_plan/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:EvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:EvaluacionController"],
        beego.ControllerComments{
            Method: "GetEvaluacion",
            Router: "/:vigencia/:plan:/:periodo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:EvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:EvaluacionController"],
        beego.ControllerComments{
            Method: "GetPlanesPeriodo",
            Router: "/planes_periodo/:vigencia/:unidad",
            AllowHTTPMethods: []string{"get"},
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
            Method: "Planes",
            Router: "/planes",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "PlanesEnFormulacion",
            Router: "/planes_formulacion",
            AllowHTTPMethods: []string{"get"},
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
            Method: "VerificarIdentificaciones",
            Router: "/verificar_identificaciones/:id",
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

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:FormulacionController"],
        beego.ControllerComments{
            Method: "VinculacionTercero",
            Router: "/vinculacion_tercero/:tercero_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "ActualizarSubgrupoDetalle",
            Router: "/actualiza_sub_detalle/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "ActualizarActividad",
            Router: "/actualizar_actividad/:id/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "ActualizarMetaPlan",
            Router: "/actualizar_meta/:id/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "ActualizarPresupuestoMeta",
            Router: "/actualizar_presupuesto_meta/:id/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "ActualizarProyectoGeneral",
            Router: "/actualizar_proyecto/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "ActualizarTablaActividad",
            Router: "/actualizar_tabla_actividad/:id/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "AllMetasPlan",
            Router: "/all_metas/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "ArmonizarInversion",
            Router: "/armonizar/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "CrearGrupoMeta",
            Router: "/crear_grupo_meta",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "CrearPlan",
            Router: "/crearplan",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "GetPlan",
            Router: "/get_infoPlan/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "GetPlanId",
            Router: "/get_plan/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "GuardarDocumentos",
            Router: "/guardar_documentos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "GuardarMeta",
            Router: "/guardar_meta/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "InactivarMeta",
            Router: "/inactivar_meta/:id/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "ProgMagnitudesPlan",
            Router: "/magnitudes/:id/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "MagnitudesProgramadas",
            Router: "/magnitudes/:id/:indexMeta",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "GetMetasProyect",
            Router: "/metaspro/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "AddProyecto",
            Router: "/proyecto",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "GetProyectoId",
            Router: "/proyecto/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "EditProyecto",
            Router: "/proyecto/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "GetAllProyectos",
            Router: "/proyectos/:aplicativo_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "VerificarMagnitudesProgramadas",
            Router: "/verificar_magnitudes/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:InversionController"],
        beego.ControllerComments{
            Method: "VersionarPlan",
            Router: "/versionar_plan/:id",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:PlanesAccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:PlanesAccionController"],
        beego.ControllerComments{
            Method: "PlanesDeAccion",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:PlanesAccionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:PlanesAccionController"],
        beego.ControllerComments{
            Method: "PlanesDeAccionPorUnidad",
            Router: "/:unidad_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "Desagregado",
            Router: "/desagregado",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "Necesidades",
            Router: "/necesidades/:nombre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "PlanAccionAnual",
            Router: "/plan_anual/:nombre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "PlanAccionEvaluacion",
            Router: "/plan_anual_evaluacion/:nombre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "PlanAccionAnualGeneral",
            Router: "/plan_anual_general/:nombre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:ReportesController"],
        beego.ControllerComments{
            Method: "ValidarReporte",
            Router: "/validar_reporte",
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
            Router: "/get_actividades/:seguimiento_id",
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
            Method: "GetEstadoTrimestre",
            Router: "/get_estado_trimestre/:plan_id/:trimestre",
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
            Method: "GuardarCualitativo",
            Router: "/guardar_cualitativo/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarCuantitativo",
            Router: "/guardar_cuantitativo/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarDocumentos",
            Router: "/guardar_documentos/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "GuardarSeguimiento",
            Router: "/guardar_seguimiento/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "HabilitarReportes",
            Router: "/habilitar_reportes",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "MigrarInformacion",
            Router: "/migrar_seguimiento/:plan_id/:trimestre",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ReportarActividad",
            Router: "/reportar_actividad/:index",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "ReportarSeguimiento",
            Router: "/reportar_seguimiento/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "RetornarActividad",
            Router: "/retornar_actividad/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "RevisarActividad",
            Router: "/revision_actividad/:plan_id/:index/:trimestre",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_mid/controllers:SeguimientoController"],
        beego.ControllerComments{
            Method: "RevisarSeguimiento",
            Router: "/revision_seguimiento/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
