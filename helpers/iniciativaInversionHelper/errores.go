package iniciativaInversionHelper

import (
	"errors"
	"fmt"
	"net/http"
)

// Mensajes para cuando el usuario no puede corregir la causa. La causa técnica
// (servicio, ruta, código, script de carga) va en el Detalle del error.
const (
	mensajeErrorInterno      = "Ocurrió un error inesperado. Intente de nuevo más tarde."
	mensajeServicio          = "No fue posible consultar la información en este momento. Intente de nuevo más tarde."
	mensajeConfiguracion     = "El módulo de iniciativas de inversión no está configurado por completo. Comuníquese con el administrador del sistema."
	mensajeSolicitudInvalida = "No fue posible procesar la solicitud. Recargue la página e intente de nuevo."
)

// ErrorNegocio es un error con el código HTTP y el mensaje que ve el usuario y,
// si es de validación, el detalle por campo. Mensaje llega a la pantalla, así que
// no lleva nombres de servicios, rutas, colecciones, códigos ni scripts: eso va
// en Detalle, que el controlador escribe en el log y nunca envía al cliente.
type ErrorNegocio struct {
	Status  int
	Mensaje string
	Detalle string
	Errores []ErrorCampo
}

func (e *ErrorNegocio) Error() string {
	if e.Detalle == "" {
		return e.Mensaje
	}
	return e.Mensaje + " (" + e.Detalle + ")"
}

// conDetalle agrega la causa técnica, que solo va al log.
func (e *ErrorNegocio) conDetalle(formato string, args ...interface{}) *ErrorNegocio {
	e.Detalle = fmt.Sprintf(formato, args...)
	return e
}

func nuevoError(status int, formato string, args ...interface{}) *ErrorNegocio {
	return &ErrorNegocio{Status: status, Mensaje: fmt.Sprintf(formato, args...)}
}

// errorSolicitud: el usuario puede corregir la solicitud (p. ej. escribir la observación).
func errorSolicitud(formato string, args ...interface{}) *ErrorNegocio {
	return nuevoError(http.StatusBadRequest, formato, args...)
}

func errorNoEncontrado(formato string, args ...interface{}) *ErrorNegocio {
	return nuevoError(http.StatusNotFound, formato, args...)
}

// errorConflicto se usa cuando el estado actual no permite la operación.
func errorConflicto(formato string, args ...interface{}) *ErrorNegocio {
	return nuevoError(http.StatusConflict, formato, args...)
}

// errorFormatoSolicitud se usa cuando la solicitud no cumple el contrato del API
// (falta un campo técnico o tiene otro formato). Es una falla del cliente, no algo
// que el usuario pueda corregir: ve un mensaje general y el detalle va al log.
func errorFormatoSolicitud(formato string, args ...interface{}) *ErrorNegocio {
	return nuevoError(http.StatusBadRequest, "%s", mensajeSolicitudInvalida).conDetalle(formato, args...)
}

// errorServicio se usa cuando falla un servicio externo (planes_crud, Terceros u
// Oikos). formato describe la falla para el log.
func errorServicio(formato string, args ...interface{}) *ErrorNegocio {
	return nuevoError(http.StatusBadGateway, "%s", mensajeServicio).conDetalle(formato, args...)
}

// errorConfiguracion se usa cuando faltan los datos base del módulo (tipo de
// plan, estados, catálogos) o la plantilla está mal armada. Lo resuelve un
// administrador; formato dice en el log qué falta y cómo cargarlo.
func errorConfiguracion(formato string, args ...interface{}) *ErrorNegocio {
	return nuevoError(http.StatusInternalServerError, "%s", mensajeConfiguracion).conDetalle(formato, args...)
}

// ComoErrorNegocio extrae el ErrorNegocio de err; si no lo es, devuelve un 500
// con un mensaje general y err como detalle.
func ComoErrorNegocio(err error) *ErrorNegocio {
	var negocio *ErrorNegocio
	if errors.As(err, &negocio) {
		return negocio
	}
	return nuevoError(http.StatusInternalServerError, "%s", mensajeErrorInterno).conDetalle("%v", err)
}
