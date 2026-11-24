package iniciativaInversionHelper

import (
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/udistrital/planeacion_mid/helpers/formulacionIniciativasHelper"
)

// Las iniciativas se formulan solo en las fechas que Planeación habilita para cada unidad: los
// registros activos del tipo de seguimiento FIPI_SP en periodo_seguimiento (ver FECHAS.md).
// Fuera de ellas no se crean iniciativas, no se guardan submódulos, no se suben anexos y no se
// envía a revisión. Consultar, comentar y la revisión de Planeación no dependen de las fechas.
//
// Cuenta la unidad de la iniciativa (dependencia_id), no la de quien la abre. Como la iniciativa
// no guarda la vigencia, sirve un registro de cualquier vigencia que cubra hoy a esa unidad.

// mensajeFueraDeFechas es lo que ve el usuario cuando intenta formular fuera de las fechas.
const mensajeFueraDeFechas = "La formulación de iniciativas no está habilitada para su unidad en estas fechas."

// periodosFormulacion trae los registros activos de fechas de formulación de iniciativas.
func (s *Servicio) periodosFormulacion() ([]map[string]interface{}, error) {
	tipo, err := s.sistema.TipoSeguimientoFormulacion()
	if err != nil {
		return nil, err
	}
	return s.repo.PeriodosActivos(tipo.ID)
}

// habilitadaHoy indica si alguno de los registros cubre hoy a la dependencia.
func habilitadaHoy(periodos []map[string]interface{}, dependenciaID string) (bool, error) {
	habilitada, err := formulacionIniciativasHelper.UnidadHabilitada(periodos, strings.TrimSpace(dependenciaID), time.Now())
	if err != nil {
		return false, errorConfiguracion("un registro de fechas de formulación de iniciativas (%s) está mal formado: %v", CodigoTipoSeguimientoFormulacion, err)
	}
	return habilitada, nil
}

// exigirFechasHabilitadas responde 409 si hoy no se formula en la dependencia, como cuando el
// estado de la iniciativa no permite editarla: el cliente recarga y la muestra en solo lectura.
func (s *Servicio) exigirFechasHabilitadas(dependenciaID string) error {
	periodos, err := s.periodosFormulacion()
	if err != nil {
		return err
	}
	habilitada, err := habilitadaHoy(periodos, dependenciaID)
	if err != nil {
		return err
	}
	if !habilitada {
		return errorConflicto("%s", mensajeFueraDeFechas)
	}
	return nil
}

// fechasDeConsulta indica, en la vista general y en el listado, qué iniciativas están fuera de las
// fechas de su unidad. Lee los registros la primera vez que los necesita y una sola vez por
// consulta. Si no se pueden revisar, la iniciativa se da por fuera de fechas (así queda cerrada,
// igual que al guardar) y el motivo va al log: la consulta sale igual.
type fechasDeConsulta struct {
	servicio *Servicio
	leidas   bool
	periodos []map[string]interface{}
	err      error
}

// fueraDeFechas indica si la iniciativa se puede editar por su estado, pero hoy no hay fechas
// habilitadas para su unidad. En los estados que no permiten editar no consulta las fechas.
func (f *fechasDeConsulta) fueraDeFechas(plan *Plan, estado string) bool {
	if !PermiteEditar(estado) {
		return false
	}
	if !f.leidas {
		f.leidas = true
		if f.periodos, f.err = f.servicio.periodosFormulacion(); f.err != nil {
			logs.Warn("iniciativaInversion: no se pudieron leer las fechas de formulación; las iniciativas editables se muestran fuera de fechas: %v", f.err)
		}
	}
	if f.err != nil {
		return true
	}
	habilitada, err := habilitadaHoy(f.periodos, plan.DependenciaID)
	if err != nil {
		logs.Warn("iniciativaInversion: la iniciativa %s se muestra fuera de fechas porque no se pudieron revisar: %v", plan.ID, err)
		return true
	}
	return !habilitada
}
