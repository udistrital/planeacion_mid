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
// Quién publica y qué pasa lo decide el estado de la iniciativa (tareasPorEstado):
//   - 01, 02 y 06: comenta el formulador; el estado del submódulo no cambia.
//   - 03: observa el Jefe y 05: observa la OAP, sobre submódulos ya formulados (también un
//     aprobado, que así se desaprueba); el submódulo queda Con observaciones. En 05 es el
//     veredicto «Con observaciones», y la iniciativa puede cambiar de estado sola (cerrarRevision).
//   - En los demás estados no recibe observaciones.
//
// La observación queda firmada con el rol del carril con el que se publica. El nombre de quien
// publica se busca en Terceros con su documento; si no está o Terceros no responde, la
// observación se publica sin nombre.
func (s *Servicio) AgregarObservacion(iniciativaID, clave string, solicitud SolicitudObservacion) (*SubmoduloVista, error) {
	texto, err := textoObservacionValido(solicitud.Texto)
	if err != nil {
		return nil, err
	}
	q, err := s.solicitante(solicitud.Actor)
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

	nueva := NuevaObservacion{Texto: texto, Usuario: q.Usuario}
	tareasDelEstado := tareasPorEstado[ctx.estado.Codigo]
	switch {
	case tareasDelEstado.editar != nil:
		if nueva.Rol, err = s.autorizar(q, ctx, tareasDelEstado.editar); err != nil {
			return nil, err
		}
	case tareasDelEstado.observar != nil:
		if nueva.Rol, err = s.autorizar(q, ctx, tareasDelEstado.observar); err != nil {
			return nil, err
		}
		if err := s.exigirRevisable(nodo, documento); err != nil {
			return nil, err
		}
		conObservaciones, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloConObservaciones)
		if err != nil {
			return nil, err
		}
		nueva.EstadoSubmoduloID = conObservaciones.ID
	default:
		return nil, errorConflicto("La iniciativa está «%s»; ya no recibe observaciones.", ctx.estado.Nombre)
	}

	nombre, err := s.nombreTercero(nueva.Usuario)
	if err != nil {
		logs.Warn("iniciativaInversion: la observación del submódulo %s de la iniciativa %s se publica sin nombre: %v", clave, iniciativaID, err)
	}
	nueva.UsuarioNombre = nombre

	guardada, err := s.repo.AgregarObservacion(documento.ID, nueva)
	if err != nil {
		return nil, err
	}
	vista, err := s.vista(*guardada, s.validador.ValidarSubmodulo(nodo, ctx.datosIniciativa()))
	if err != nil {
		return nil, err
	}
	if ctx.estado.Codigo == EstadoIniciativaEnRevision {
		vista.EstadoIniciativa = s.cerrarRevision(iniciativaID)
	}
	return vista, nil
}

// exigirRevisable comprueba que el submódulo ya esté formulado para observarlo o aprobarlo.
func (s *Servicio) exigirRevisable(nodo *Nodo, documento *SubmoduloIniciativa) error {
	actual, err := s.sistema.EstadoSubmoduloPorID(documento.EstadoSubmoduloID)
	if err != nil {
		return err
	}
	if !SubmoduloRevisable(actual.Codigo) {
		return errorConflicto("El submódulo %s no está formulado; no se puede revisar.", nodo.Nombre)
	}
	return nil
}

// textoObservacionValido limpia el texto y exige que no esté vacío ni pase del máximo.
func textoObservacionValido(texto string) (string, error) {
	texto = strings.TrimSpace(texto)
	switch {
	case texto == "":
		return texto, errorSolicitud("Escriba la observación.")
	case utf8.RuneCountInString(texto) > MaximoTextoObservacion:
		return texto, errorSolicitud("La observación no puede superar %d caracteres.", MaximoTextoObservacion)
	}
	return texto, nil
}
