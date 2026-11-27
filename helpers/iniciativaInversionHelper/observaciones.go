package iniciativaInversionHelper

import (
	"strings"
	"unicode/utf8"

	"github.com/astaxie/beego/logs"
)

// MaximoTextoObservacion es el largo máximo de una observación (igual que en planes_crud).
const MaximoTextoObservacion = 4000

// AgregarObservacion publica una observación en el historial del submódulo y lo deja Con
// observaciones, para que la unidad lo corrija. Pueden publicar:
//   - el Jefe, con la iniciativa Por aprobar;
//   - la OAP, con la iniciativa En revisión OAP. Es su veredicto «Con observaciones»; si es el
//     último, la iniciativa cambia de estado sola (cerrarRevision).
//
// En los demás estados nadie publica. Solo se observa un submódulo ya formulado (también uno
// aprobado, que así deja de estarlo). La versión del submódulo no cambia.
//
// La observación queda firmada con el rol con el que se publica y con el nombre de quien la
// publica, que se busca en Terceros; si no se encuentra, se publica sin nombre.
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

	carril := tareasPorEstado[ctx.estado.Codigo].observar
	if carril == nil {
		return nil, errorConflicto("La iniciativa está «%s»; en este estado no recibe observaciones.", ctx.estado.Nombre)
	}
	rol, err := s.autorizar(q, ctx, carril)
	if err != nil {
		return nil, err
	}
	if err := s.exigirRevisable(nodo, documento); err != nil {
		return nil, err
	}
	conObservaciones, err := s.sistema.EstadoSubmodulo(EstadoSubmoduloConObservaciones)
	if err != nil {
		return nil, err
	}
	nueva := NuevaObservacion{Texto: texto, Rol: rol, Usuario: q.Usuario, EstadoSubmoduloID: conObservaciones.ID}

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
