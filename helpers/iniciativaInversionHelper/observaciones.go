package iniciativaInversionHelper

import (
	"strings"
	"unicode/utf8"

	"github.com/astaxie/beego/logs"
)

// MaximoTextoObservacion es el largo máximo de una observación (igual que en planes_crud).
const MaximoTextoObservacion = 4000

// AgregarObservacion publica un comentario en el historial del submódulo. No cambia
// la versión del submódulo, así que no choca con quien está editando los datos.
// Quién comenta y qué pasa lo decide el estado de la iniciativa:
//   - Sin formular, Formulación iniciada o Por corregir: comenta la dependencia; el
//     estado del submódulo no cambia.
//   - Enviado a revisión: comenta Planeación, solo sobre submódulos formulados; el
//     submódulo queda "Con observaciones".
//   - Avalado: no recibe observaciones.
//
// El nombre de quien comenta se busca en Terceros con su documento. Si no está o
// Terceros no responde, la observación se publica sin nombre.
func (s *Servicio) AgregarObservacion(iniciativaID, clave string, solicitud SolicitudObservacion) (*SubmoduloVista, error) {
	nueva, err := observacionValida(solicitud)
	if err != nil {
		return nil, err
	}
	ctx, err := s.cargarIniciativa(iniciativaID)
	if err != nil {
		return nil, err
	}
	nodo, documento, err := ctx.submodulo(clave)
	if err != nil {
		return nil, err
	}

	switch {
	case PermiteRevisar(ctx.estado.Codigo):
		if err := s.exigirRevisable(ctx, nodo, documento); err != nil {
			return nil, err
		}
		conObservaciones, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloConObservaciones)
		if err != nil {
			return nil, err
		}
		nueva.EstadoSubmoduloID = conObservaciones.ID
	case PermiteEditar(ctx.estado.Codigo):
		// La dependencia comenta mientras formula o corrige.
	default:
		return nil, errorConflicto("La iniciativa está «%s»; ya no recibe observaciones.", ctx.estado.Nombre)
	}

	nombre, err := nombreTercero(nueva.Usuario)
	if err != nil {
		logs.Warn("iniciativaInversion: la observación del submódulo %s de la iniciativa %s se publica sin nombre: %v", clave, iniciativaID, err)
	}
	nueva.UsuarioNombre = nombre

	guardada, err := s.repo.AgregarObservacion(documento.ID, nueva)
	if err != nil {
		return nil, err
	}
	return s.vista(*guardada, s.validador.ValidarSubmodulo(nodo, ctx.datosIniciativa()))
}

// exigirRevisable comprueba que Planeación pueda revisar el submódulo: la iniciativa
// debe estar en revisión y el submódulo, formulado.
func (s *Servicio) exigirRevisable(ctx *contextoIniciativa, nodo *Nodo, documento *SubmoduloIniciativa) error {
	if !PermiteRevisar(ctx.estado.Codigo) {
		return errorConflicto("Solo se revisan iniciativas enviadas a revisión; esta está «%s».", ctx.estado.Nombre)
	}
	actual, err := s.sistema.EstadoSubmoduloPorID(documento.EstadoSubmoduloID)
	if err != nil {
		return err
	}
	if !SubmoduloRevisable(actual.Codigo) {
		return errorConflicto("El submódulo %s no está formulado; no se puede revisar.", nodo.Nombre)
	}
	return nil
}

// observacionValida limpia la solicitud y exige texto, rol y usuario; el usuario es
// el documento de quien comenta.
func observacionValida(solicitud SolicitudObservacion) (NuevaObservacion, error) {
	nueva := NuevaObservacion{
		Texto:   strings.TrimSpace(solicitud.Texto),
		Rol:     strings.TrimSpace(solicitud.Rol),
		Usuario: strings.TrimSpace(solicitud.Usuario),
	}
	switch {
	case nueva.Texto == "":
		return nueva, errorSolicitud("Escriba la observación.")
	case utf8.RuneCountInString(nueva.Texto) > MaximoTextoObservacion:
		return nueva, errorSolicitud("La observación no puede superar %d caracteres.", MaximoTextoObservacion)
	case nueva.Rol == "" || nueva.Usuario == "":
		return nueva, errorSolicitud("No se pudo identificar quién publica la observación. Vuelva a iniciar sesión e intente de nuevo.").
			conDetalle("rol y usuario son obligatorios (rol vacío: %t, usuario vacío: %t)", nueva.Rol == "", nueva.Usuario == "")
	case !patronDocumento.MatchString(nueva.Usuario):
		return nueva, errorSolicitud("No se pudo identificar quién publica la observación. Vuelva a iniciar sesión e intente de nuevo.").
			conDetalle("el usuario no cumple el formato de documento: solo letras y números, máximo 20")
	}
	return nueva, nil
}

// tieneTexto indica si alguna observación trae texto.
func tieneTexto(observaciones []Observacion) bool {
	for _, observacion := range observaciones {
		if strings.TrimSpace(observacion.Texto) != "" {
			return true
		}
	}
	return false
}
