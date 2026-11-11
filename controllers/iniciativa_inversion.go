package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/planeacion_mid/helpers/iniciativaInversionHelper"
)

// IniciativaInversionController expone la formulación de las iniciativas de
// inversión con plantillas de Construcción. Solo traduce HTTP; la lógica está en
// helpers/iniciativaInversionHelper.
type IniciativaInversionController struct {
	beego.Controller
}

var (
	servicioIniciativa       *iniciativaInversionHelper.Servicio
	servicioIniciativaUnaVez sync.Once
	unidadesUsuario          *iniciativaInversionHelper.Unidades
	unidadesUsuarioUnaVez    sync.Once
)

func servicioIniciativas() *iniciativaInversionHelper.Servicio {
	servicioIniciativaUnaVez.Do(func() {
		servicioIniciativa = iniciativaInversionHelper.NuevoServicioPorDefecto(beego.AppConfig.String("PlanesService"))
	})
	return servicioIniciativa
}

func unidadesDeUsuarios() *iniciativaInversionHelper.Unidades {
	unidadesUsuarioUnaVez.Do(func() {
		unidadesUsuario = iniciativaInversionHelper.NuevasUnidadesPorDefecto(
			beego.AppConfig.String("TercerosService"), beego.AppConfig.String("OikosService"))
	})
	return unidadesUsuario
}

// URLMapping ...
func (c *IniciativaInversionController) URLMapping() {
	c.Mapping("ListarPlantillas", c.ListarPlantillas)
	c.Mapping("ListarIniciativas", c.ListarIniciativas)
	c.Mapping("CrearIniciativa", c.CrearIniciativa)
	c.Mapping("ObtenerIniciativa", c.ObtenerIniciativa)
	c.Mapping("ObtenerDefinicion", c.ObtenerDefinicion)
	c.Mapping("ObtenerSubmodulo", c.ObtenerSubmodulo)
	c.Mapping("GuardarSubmodulo", c.GuardarSubmodulo)
	c.Mapping("AgregarObservacion", c.AgregarObservacion)
	c.Mapping("RevisarSubmodulo", c.RevisarSubmodulo)
	c.Mapping("EjecutarAccion", c.EjecutarAccion)
	c.Mapping("ObtenerCatalogo", c.ObtenerCatalogo)
	c.Mapping("ObtenerUnidades", c.ObtenerUnidades)
}

// ObtenerUnidades ...
// @Title ObtenerUnidades
// @Description Unidades (dependencias de Oikos) del usuario según sus vinculaciones activas en Terceros, con los mismos cargos que Formulación. Lista vacía si no tiene ninguna
// @Param	documento	path	string	true	"Documento del usuario"
// @Success 200 {object} []iniciativaInversionHelper.Unidad
// @Failure 502 Error consultando Terceros u Oikos
// @router /unidades/:documento [get]
func (c *IniciativaInversionController) ObtenerUnidades() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	unidades, err := unidadesDeUsuarios().De(c.Ctx.Input.Param(":documento"))
	c.responder(http.StatusOK, unidades, err)
}

// ListarPlantillas ...
// @Title ListarPlantillas
// @Description Plantillas (construcciones) activas de iniciativas de inversión
// @Success 200 {object} []iniciativaInversionHelper.PlantillaResumen
// @Failure 502 Error consultando planes_crud
// @router /plantillas [get]
func (c *IniciativaInversionController) ListarPlantillas() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	plantillas, err := servicioIniciativas().Plantillas()
	c.responder(http.StatusOK, plantillas, err)
}

// ListarIniciativas ...
// @Title ListarIniciativas
// @Description Iniciativas activas con su número y su estado general, de la más reciente a la más antigua. Opcionalmente solo las de unas dependencias
// @Param	dependencias	query	string	false	"Ids de dependencia (Oikos) separados por coma"
// @Success 200 {object} []iniciativaInversionHelper.IniciativaListado
// @router /iniciativas [get]
func (c *IniciativaInversionController) ListarIniciativas() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	iniciativas, err := servicioIniciativas().ListarIniciativas(listaDeIds(c.GetString("dependencias")))
	c.responder(http.StatusOK, iniciativas, err)
}

// CrearIniciativa ...
// @Title CrearIniciativa
// @Description Crea una iniciativa con todos sus submódulos sin formular. Si no se indica la plantilla, usa la más reciente
// @Param	body	body	iniciativaInversionHelper.SolicitudCrearIniciativa	true	"dependencia_id, plantilla_id (opcional), nombre (opcional), vigencias (opcional, lista JSON)"
// @Success 201 {object} iniciativaInversionHelper.ResumenIniciativa
// @Failure 400 Solicitud inválida
// @router /iniciativas [post]
func (c *IniciativaInversionController) CrearIniciativa() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	var solicitud iniciativaInversionHelper.SolicitudCrearIniciativa
	if !c.leerCuerpo(&solicitud) {
		return
	}
	resumen, err := servicioIniciativas().CrearIniciativa(solicitud)
	c.responder(http.StatusCreated, resumen, err)
}

// ObtenerIniciativa ...
// @Title ObtenerIniciativa
// @Description Vista general: estado de la iniciativa, estado y semáforo de cada submódulo y acciones posibles
// @Param	id	path	string	true	"Id de la iniciativa"
// @Success 200 {object} iniciativaInversionHelper.ResumenIniciativa
// @Failure 404 No existe
// @router /iniciativas/:id [get]
func (c *IniciativaInversionController) ObtenerIniciativa() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	resumen, err := servicioIniciativas().Resumen(c.Ctx.Input.Param(":id"))
	c.responder(http.StatusOK, resumen, err)
}

// ObtenerDefinicion ...
// @Title ObtenerDefinicion
// @Description Construcción (submódulos y campos) de la plantilla de la iniciativa
// @Param	id	path	string	true	"Id de la iniciativa"
// @Success 200 {object} iniciativaInversionHelper.Definicion
// @router /iniciativas/:id/definicion [get]
func (c *IniciativaInversionController) ObtenerDefinicion() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	definicion, err := servicioIniciativas().Definicion(c.Ctx.Input.Param(":id"))
	c.responder(http.StatusOK, definicion, err)
}

// ObtenerSubmodulo ...
// @Title ObtenerSubmodulo
// @Description Datos de un submódulo, sus errores vigentes y los datos de los submódulos de los que depende
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	clave	path	string	true	"Clave del submódulo (p. ej. problematica)"
// @Success 200 {object} iniciativaInversionHelper.SubmoduloConRelacionadas
// @router /iniciativas/:id/submodulos/:clave [get]
func (c *IniciativaInversionController) ObtenerSubmodulo() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	submodulo, err := servicioIniciativas().ObtenerSubmodulo(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":clave"))
	c.responder(http.StatusOK, submodulo, err)
}

// GuardarSubmodulo ...
// @Title GuardarSubmodulo
// @Description Guarda los datos de un submódulo. Responde los errores de validación; el submódulo queda Formulado solo si no hay errores
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	clave	path	string	true	"Clave del submódulo"
// @Param	body	body	iniciativaInversionHelper.SolicitudGuardarSubmodulo	true	"version (la consultada) y datos"
// @Success 200 {object} iniciativaInversionHelper.SubmoduloVista
// @Failure 409 La iniciativa no se puede editar o el submódulo cambió desde que se consultó
// @router /iniciativas/:id/submodulos/:clave [put]
func (c *IniciativaInversionController) GuardarSubmodulo() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	var solicitud iniciativaInversionHelper.SolicitudGuardarSubmodulo
	if !c.leerCuerpo(&solicitud) {
		return
	}
	submodulo, err := servicioIniciativas().GuardarSubmodulo(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":clave"), solicitud)
	c.responder(http.StatusOK, submodulo, err)
}

// AgregarObservacion ...
// @Title AgregarObservacion
// @Description Publica una observación en el historial del submódulo sin cambiar su versión. En formulación o por corregir comenta la dependencia; en revisión comenta Planeación y el submódulo queda Con observaciones
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	clave	path	string	true	"Clave del submódulo"
// @Param	body	body	iniciativaInversionHelper.SolicitudObservacion	true	"texto, rol y usuario (documento)"
// @Success 201 {object} iniciativaInversionHelper.SubmoduloVista
// @Failure 400 Solicitud inválida
// @Failure 409 El estado de la iniciativa o del submódulo no permite comentar
// @router /iniciativas/:id/submodulos/:clave/observaciones [post]
func (c *IniciativaInversionController) AgregarObservacion() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	var solicitud iniciativaInversionHelper.SolicitudObservacion
	if !c.leerCuerpo(&solicitud) {
		return
	}
	submodulo, err := servicioIniciativas().AgregarObservacion(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":clave"), solicitud)
	c.responder(http.StatusCreated, submodulo, err)
}

// RevisarSubmodulo ...
// @Title RevisarSubmodulo
// @Description Planeación da el submódulo por revisado sin observaciones. Las observaciones se registran con POST .../observaciones y su historial no se borra
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	clave	path	string	true	"Clave del submódulo"
// @Param	body	body	iniciativaInversionHelper.SolicitudRevisarSubmodulo	true	"version (la consultada)"
// @Success 200 {object} iniciativaInversionHelper.SubmoduloVista
// @Failure 400 Se enviaron observaciones (van por POST .../observaciones)
// @Failure 409 La iniciativa no está en revisión o el submódulo no está formulado
// @router /iniciativas/:id/submodulos/:clave/revision [put]
func (c *IniciativaInversionController) RevisarSubmodulo() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	var solicitud iniciativaInversionHelper.SolicitudRevisarSubmodulo
	if !c.leerCuerpo(&solicitud) {
		return
	}
	submodulo, err := servicioIniciativas().RevisarSubmodulo(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":clave"), solicitud)
	c.responder(http.StatusOK, submodulo, err)
}

// EjecutarAccion ...
// @Title EjecutarAccion
// @Description Cambia el estado general de la iniciativa: enviar, devolver, avalar o no-avalar
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	accion	path	string	true	"enviar | devolver | avalar | no-avalar"
// @Success 200 {object} iniciativaInversionHelper.ResumenIniciativa
// @Failure 409 La acción no se permite en el estado actual
// @router /iniciativas/:id/acciones/:accion [put]
func (c *IniciativaInversionController) EjecutarAccion() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	resumen, err := servicioIniciativas().EjecutarAccion(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":accion"))
	c.responder(http.StatusOK, resumen, err)
}

// ObtenerCatalogo ...
// @Title ObtenerCatalogo
// @Description Elementos de un catálogo de las iniciativas (procesos, sectores_pnd, programas_pnd, productos...)
// @Param	nombre	path	string	true	"Nombre del catálogo"
// @Param	padre	query	string	false	"Código del padre para filtrar (p. ej. sector de los programas)"
// @Success 200 {object} []iniciativaInversionHelper.ItemCatalogo
// @router /catalogos/:nombre [get]
func (c *IniciativaInversionController) ObtenerCatalogo() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	items, err := servicioIniciativas().Catalogo(c.Ctx.Input.Param(":nombre"), c.GetString("padre"))
	c.responder(http.StatusOK, items, err)
}

// listaDeIds separa "1, 2,3" en ["1", "2", "3"], sin vacíos.
func listaDeIds(texto string) []string {
	ids := []string{}
	for _, id := range strings.Split(texto, ",") {
		if id = strings.TrimSpace(id); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

// leerCuerpo decodifica el JSON del cuerpo; si es inválido responde 400.
func (c *IniciativaInversionController) leerCuerpo(destino interface{}) bool {
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, destino); err != nil {
		logs.Warn("IniciativaInversionController: %s %s: el cuerpo no es un JSON válido: %v", c.Ctx.Request.Method, c.Ctx.Request.URL.Path, err)
		c.escribir(http.StatusBadRequest, false, "No fue posible procesar la solicitud. Recargue la página e intente de nuevo.", nil)
		return false
	}
	return true
}

// responder escribe el sobre { Success, Status, Message, Data } con el código
// HTTP real. Los errores de validación van en Data. Message es lo que ve el
// usuario; la causa técnica del error (Detalle) solo va al log.
func (c *IniciativaInversionController) responder(statusExito int, data interface{}, err error) {
	if err != nil {
		negocio := iniciativaInversionHelper.ComoErrorNegocio(err)
		if negocio.Detalle != "" {
			registrar := logs.Warn
			if negocio.Status >= http.StatusInternalServerError {
				registrar = logs.Error
			}
			registrar("IniciativaInversionController: %s %s -> %d: %s", c.Ctx.Request.Method, c.Ctx.Request.URL.Path, negocio.Status, negocio.Detalle)
		}
		var detalle interface{}
		if len(negocio.Errores) > 0 {
			detalle = negocio.Errores
		}
		c.escribir(negocio.Status, false, negocio.Mensaje, detalle)
		return
	}
	c.escribir(statusExito, true, "Successful", data)
}

func (c *IniciativaInversionController) escribir(status int, exito bool, mensaje string, data interface{}) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{
		"Success": exito,
		"Status":  strconv.Itoa(status),
		"Message": mensaje,
		"Data":    data,
	}
	c.ServeJSON()
}
