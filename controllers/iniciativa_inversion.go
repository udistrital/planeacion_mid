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
		servicioIniciativa = iniciativaInversionHelper.NuevoServicioPorDefecto(
			beego.AppConfig.String("PlanesService"), beego.AppConfig.String("TercerosService"),
			beego.AppConfig.String("OikosService"), beego.AppConfig.String("GestorDocumental"))
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
	c.Mapping("SubirDocumento", c.SubirDocumento)
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
// @Description Iniciativas activas que consulta el usuario (la OAP, todas; los demás, las de sus unidades), con su número, su estado general y lo que el usuario hace con cada una, de la más reciente a la más antigua
// @Param	usuario	query	string	true	"Documento del usuario"
// @Param	roles	query	string	true	"Roles del usuario separados por coma"
// @Success 200 {object} []iniciativaInversionHelper.IniciativaListado
// @Failure 400 Falta el usuario o sus roles
// @router /iniciativas [get]
func (c *IniciativaInversionController) ListarIniciativas() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	iniciativas, err := servicioIniciativas().ListarIniciativas(c.actorDeConsulta())
	c.responder(http.StatusOK, iniciativas, err)
}

// CrearIniciativa ...
// @Title CrearIniciativa
// @Description Crea una iniciativa con todos sus submódulos sin formular, a nombre de quien la crea. Si no se indica la plantilla, usa la más reciente
// @Param	body	body	iniciativaInversionHelper.SolicitudCrearIniciativa	true	"usuario, roles, dependencia_id, plantilla_id (opcional), nombre (opcional), vigencias (opcional, lista JSON)"
// @Success 201 {object} iniciativaInversionHelper.ResumenIniciativa
// @Failure 400 Solicitud inválida
// @Failure 403 El usuario no formula en esa unidad
// @Failure 409 Fuera de las fechas de formulación de la unidad
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
// @Description Vista general: estado de la iniciativa y de cada submódulo, y lo que el usuario puede hacer con ella (permisos y acciones)
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	usuario	query	string	true	"Documento del usuario"
// @Param	roles	query	string	true	"Roles del usuario separados por coma"
// @Success 200 {object} iniciativaInversionHelper.ResumenIniciativa
// @Failure 403 La iniciativa es de una unidad que el usuario no consulta
// @Failure 404 No existe
// @router /iniciativas/:id [get]
func (c *IniciativaInversionController) ObtenerIniciativa() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	resumen, err := servicioIniciativas().Resumen(c.Ctx.Input.Param(":id"), c.actorDeConsulta())
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
// @Param	body	body	iniciativaInversionHelper.SolicitudGuardarSubmodulo	true	"usuario, roles, version (la consultada) y datos"
// @Success 200 {object} iniciativaInversionHelper.SubmoduloVista
// @Failure 403 El usuario no formula la iniciativa
// @Failure 409 La iniciativa no se puede editar, el submódulo está aprobado o cambió desde que se consultó, o está fuera de fechas
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

// SubirDocumento ...
// @Title SubirDocumento
// @Description Sube el archivo de un campo documento (Anexos) al gestor documental y responde su referencia. No cambia el submódulo: la referencia queda en la iniciativa al guardarlo (PUT .../submodulos/:clave)
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	clave	path	string	true	"Clave del submódulo (p. ej. anexos)"
// @Param	body	body	iniciativaInversionHelper.SolicitudDocumento	true	"usuario, roles, clave (ruta del campo sin ids de fila, p. ej. otros_anexos.archivo), nombre y archivo (base64)"
// @Success 201 {object} iniciativaInversionHelper.ReferenciaDocumento
// @Failure 400 Formato no permitido, archivo vacío o más grande de lo permitido
// @Failure 403 El usuario no formula la iniciativa
// @Failure 409 La iniciativa o el submódulo no se pueden editar, o está fuera de fechas
// @Failure 502 Error en el gestor documental
// @router /iniciativas/:id/submodulos/:clave/documentos [post]
func (c *IniciativaInversionController) SubirDocumento() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	var solicitud iniciativaInversionHelper.SolicitudDocumento
	if !c.leerCuerpo(&solicitud) {
		return
	}
	referencia, err := servicioIniciativas().SubirDocumento(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":clave"), solicitud)
	c.responder(http.StatusCreated, referencia, err)
}

// AgregarObservacion ...
// @Title AgregarObservacion
// @Description Publica una observación en el historial del submódulo sin cambiar su versión. En 01, 02 y 06 comenta el formulador; en 03 observa el Jefe y en 05 la OAP, y el submódulo queda Con observaciones (en 05 la iniciativa puede cambiar de estado sola: va en estado_iniciativa)
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	clave	path	string	true	"Clave del submódulo"
// @Param	body	body	iniciativaInversionHelper.SolicitudObservacion	true	"usuario, roles y texto"
// @Success 201 {object} iniciativaInversionHelper.SubmoduloVista
// @Failure 400 Solicitud inválida
// @Failure 403 El usuario no actúa en el carril de este estado
// @Failure 409 El estado de la iniciativa o del submódulo no permite comentar, o está fuera de fechas
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
// @Description La OAP aprueba el submódulo (veredicto Aprobado) mientras revisa la iniciativa. El veredicto Con observaciones se da con POST .../observaciones. Con el último veredicto la iniciativa cambia de estado sola: va en estado_iniciativa
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	clave	path	string	true	"Clave del submódulo"
// @Param	body	body	iniciativaInversionHelper.SolicitudRevisarSubmodulo	true	"usuario, roles y version (la consultada)"
// @Success 200 {object} iniciativaInversionHelper.SubmoduloVista
// @Failure 403 El usuario no es de la OAP
// @Failure 409 La iniciativa no está En revisión OAP, el submódulo no está formulado o está fuera de fechas
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
// @Description Cambia el estado general de la iniciativa con una acción del carril del usuario
// @Param	id	path	string	true	"Id de la iniciativa"
// @Param	accion	path	string	true	"enviar | aprobar | devolver | iniciar-revision | rechazar | retractar-rechazo"
// @Param	body	body	iniciativaInversionHelper.SolicitudAccion	true	"usuario, roles y, al rechazar, motivo"
// @Success 200 {object} iniciativaInversionHelper.ResumenIniciativa
// @Failure 400 La acción no existe o falta el motivo del rechazo
// @Failure 403 El usuario no actúa en el carril de la acción
// @Failure 409 La acción no se permite en el estado actual o está fuera de fechas
// @router /iniciativas/:id/acciones/:accion [put]
func (c *IniciativaInversionController) EjecutarAccion() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	var solicitud iniciativaInversionHelper.SolicitudAccion
	if !c.leerCuerpo(&solicitud) {
		return
	}
	resumen, err := servicioIniciativas().EjecutarAccion(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":accion"), solicitud)
	c.responder(http.StatusOK, resumen, err)
}

// ObtenerCatalogo ...
// @Title ObtenerCatalogo
// @Description Elementos de un catálogo de las iniciativas (procesos, productos, municipios...)
// @Param	nombre	path	string	true	"Nombre del catálogo"
// @Param	padre	query	string	false	"Código del padre para filtrar (p. ej. departamento de los municipios)"
// @Success 200 {object} []iniciativaInversionHelper.ItemCatalogo
// @router /catalogos/:nombre [get]
func (c *IniciativaInversionController) ObtenerCatalogo() {
	defer helpers.ErrorController(c.Controller, "IniciativaInversionController")
	items, err := servicioIniciativas().Catalogo(c.Ctx.Input.Param(":nombre"), c.GetString("padre"))
	c.responder(http.StatusOK, items, err)
}

// actorDeConsulta toma de la consulta (?usuario=...&roles=A,B) quién hace una solicitud de lectura;
// las que escriben lo traen en el cuerpo. El servicio valida que no falten.
func (c *IniciativaInversionController) actorDeConsulta() iniciativaInversionHelper.Actor {
	return iniciativaInversionHelper.Actor{Usuario: c.GetString("usuario"), Roles: listaSeparadaPorComas(c.GetString("roles"))}
}

// listaSeparadaPorComas separa "A, B,C" en ["A", "B", "C"], sin vacíos.
func listaSeparadaPorComas(texto string) []string {
	elementos := []string{}
	for _, elemento := range strings.Split(texto, ",") {
		if elemento = strings.TrimSpace(elemento); elemento != "" {
			elementos = append(elementos, elemento)
		}
	}
	return elementos
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
