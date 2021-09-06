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

}
