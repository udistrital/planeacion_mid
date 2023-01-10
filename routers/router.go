// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/planeacion_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/arbol",
			beego.NSInclude(
				&controllers.ArbolController{},
			),
		),
		beego.NSNamespace("/formato",
			beego.NSInclude(
				&controllers.FormatoController{},
			),
		),
		beego.NSNamespace("/formulacion",
			beego.NSInclude(
				&controllers.FormulacionController{},
			),
		),
		beego.NSNamespace("/seguimiento",
			beego.NSInclude(
				&controllers.SeguimientoController{},
			),
		),
		beego.NSNamespace("/reportes",
			beego.NSInclude(
				&controllers.ReportesController{},
			),
		),

		beego.NSNamespace("/inversion",
			beego.NSInclude(
				&controllers.InversionController{},
			),
		),

		beego.NSNamespace("/evaluacion",
		beego.NSInclude(
			&controllers.EvaluacionController{},
		),
	),

	)
	beego.AddNamespace(ns)
}
