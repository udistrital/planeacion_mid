package iniciativaInversionHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Los catálogos del Plan Distrital de Desarrollo y de los planes de la Universidad (Plan
// Estratégico de Desarrollo y Plan Indicativo) en Articulación instrumentos de planificación no
// están en elemento_catalogo: salen de los planes armados en Construcción de planes, los mismos
// con los que se armoniza en Formulación y en MACRO.
// Cada catálogo es un nivel del árbol: el primero, los planes; los siguientes, sus nodos, con
// código = id del subgrupo y padre = id del nodo de arriba, para que depende_de encadene
// plan → nivel 1 → nivel 2… Para el validador y el cliente son catálogos como los demás
// (/catalogos/<nombre>).

// Tipos de plan de los instrumentos de planificación.
const (
	CodigoTipoPlanDistrital   = "PDD_SP"
	CodigoTipoPlanEstrategico = "PD_SP"
	CodigoTipoPlanIndicativo  = "PLI_SP"
)

// instrumentoPlanificacion es un tipo de plan y sus catálogos, uno por nivel: [0] son los planes
// y [n], los nodos del nivel n.
type instrumentoPlanificacion struct {
	tipoPlan  string
	catalogos []string
}

// instrumentosPlanificacion son los planes con los que se articula una iniciativa.
var instrumentosPlanificacion = []instrumentoPlanificacion{
	// Plan Distrital de Desarrollo: objetivo estratégico → estrategia → programa intersectorial.
	{tipoPlan: CodigoTipoPlanDistrital, catalogos: []string{"planes_pdd", "objetivos_pdd", "estrategias_pdd", "programas_pdd"}},
	// Plan Estratégico de Desarrollo: lineamiento estratégico → meta estratégica.
	{tipoPlan: CodigoTipoPlanEstrategico, catalogos: []string{"planes_ped", "lineamientos_ped", "metas_ped"}},
	// Plan Indicativo: eje transformador → lineamiento de acción → estrategia.
	{tipoPlan: CodigoTipoPlanIndicativo, catalogos: []string{"planes_pi", "ejes_pi", "lineamientos_pi", "estrategias_pi"}},
}

// RepositorioInstrumentos lee de planes_crud los planes de un tipo y su árbol.
type RepositorioInstrumentos interface {
	BuscarPorCodigo(recurso, codigo string) (*Estado, error)
	BuscarPlanes(filtros map[string]string) ([]Plan, error)
	HijosDe(id string) ([]Subgrupo, error)
}

// RepositorioConInstrumentos lee los catálogos de los instrumentos de planificación de sus
// planes y los demás del repositorio base. Todos los catálogos de un instrumento salen de una
// sola lectura del árbol, que se guarda un rato: el formulario los pide juntos.
type RepositorioConInstrumentos struct {
	base     RepositorioElementosCatalogo
	crud     RepositorioInstrumentos
	duracion time.Duration
	ahora    func() time.Time
	cupos    chan struct{}

	porCatalogo map[string]nivelInstrumento
	arboles     map[string]*arbolInstrumento // por tipo de plan; el mapa no cambia después de crearse
}

// nivelInstrumento ubica un catálogo: de qué instrumento y qué nivel del árbol es.
type nivelInstrumento struct {
	instrumento instrumentoPlanificacion
	nivel       int
}

// arbolInstrumento son los niveles leídos de un instrumento, hasta que vencen.
type arbolInstrumento struct {
	mu      sync.Mutex
	niveles [][]ItemCatalogo
	vence   time.Time
}

// nodoInstrumento es un elemento de un nivel y si tiene hijos que leer.
type nodoInstrumento struct {
	item     ItemCatalogo
	conHijos bool
}

// NuevoRepositorioConInstrumentos arma el repositorio. duracion es cuánto se guarda cada árbol
// leído y maxSimultaneas, cuántas llamadas al crud se hacen a la vez al leerlo.
func NuevoRepositorioConInstrumentos(base RepositorioElementosCatalogo, crud RepositorioInstrumentos, duracion time.Duration, maxSimultaneas int) *RepositorioConInstrumentos {
	if maxSimultaneas < 1 {
		maxSimultaneas = 1
	}
	r := &RepositorioConInstrumentos{
		base:        base,
		crud:        crud,
		duracion:    duracion,
		ahora:       time.Now,
		cupos:       make(chan struct{}, maxSimultaneas),
		porCatalogo: map[string]nivelInstrumento{},
		arboles:     map[string]*arbolInstrumento{},
	}
	for _, instrumento := range instrumentosPlanificacion {
		r.arboles[instrumento.tipoPlan] = &arbolInstrumento{}
		for nivel, catalogo := range instrumento.catalogos {
			r.porCatalogo[catalogo] = nivelInstrumento{instrumento: instrumento, nivel: nivel}
		}
	}
	return r
}

// ElementosDeCatalogo implementa RepositorioElementosCatalogo.
func (r *RepositorioConInstrumentos) ElementosDeCatalogo(nombre string) ([]ItemCatalogo, error) {
	ubicacion, esDeInstrumento := r.porCatalogo[nombre]
	if !esDeInstrumento {
		return r.base.ElementosDeCatalogo(nombre)
	}
	niveles, err := r.niveles(ubicacion.instrumento)
	if err != nil {
		return nil, err
	}
	return niveles[ubicacion.nivel], nil
}

// niveles devuelve el árbol del instrumento, del que se leyó hace poco o leyéndolo de nuevo.
// Las consultas que llegan juntas esperan una sola lectura.
func (r *RepositorioConInstrumentos) niveles(instrumento instrumentoPlanificacion) ([][]ItemCatalogo, error) {
	arbol := r.arboles[instrumento.tipoPlan]
	arbol.mu.Lock()
	defer arbol.mu.Unlock()
	if arbol.niveles != nil && r.ahora().Before(arbol.vence) {
		return arbol.niveles, nil
	}
	niveles, err := r.leer(instrumento)
	if err != nil {
		return nil, err
	}
	arbol.niveles, arbol.vence = niveles, r.ahora().Add(r.duracion)
	return niveles, nil
}

// leer recorre los planes activos del tipo armados en Construcción (formato: true; las copias
// con formato: false son formulaciones de una dependencia) hasta el último nivel que se usa.
// Un plan sin nodos no sirve para articular, así que no se ofrece.
func (r *RepositorioConInstrumentos) leer(instrumento instrumentoPlanificacion) ([][]ItemCatalogo, error) {
	tipo, err := r.crud.BuscarPorCodigo(RecursoTipoPlan, instrumento.tipoPlan)
	if err != nil {
		return nil, err
	}
	planes, err := r.crud.BuscarPlanes(map[string]string{"tipo_plan_id": tipo.ID, "formato": "true", "activo": "true"})
	if err != nil {
		return nil, err
	}
	candidatos := make([]nodoInstrumento, 0, len(planes))
	for _, plan := range planes {
		candidatos = append(candidatos, nodoInstrumento{item: ItemCatalogo{Codigo: plan.ID, Nombre: nombreDePlan(plan)}, conHijos: true})
	}
	hijos, err := r.hijosDe(candidatos)
	if err != nil {
		return nil, err
	}

	niveles := make([][]ItemCatalogo, len(instrumento.catalogos))
	niveles[0] = []ItemCatalogo{}
	actual := []nodoInstrumento{}
	for i, plan := range candidatos {
		if len(hijos[i]) > 0 {
			niveles[0] = append(niveles[0], plan.item)
			actual = append(actual, hijos[i]...)
		}
	}
	if len(niveles[0]) == 0 {
		return nil, errorNoEncontrado("No se encontraron las opciones de este campo. Comuníquese con el administrador del sistema.").
			conDetalle("no hay planes %s activos con formato: true y con nodos. Se arman en Construcción de planes.", instrumento.tipoPlan)
	}
	niveles[1] = itemsDe(actual)
	for nivel := 2; nivel < len(niveles); nivel++ {
		if hijos, err = r.hijosDe(actual); err != nil {
			return nil, err
		}
		actual = []nodoInstrumento{}
		for _, grupo := range hijos {
			actual = append(actual, grupo...)
		}
		niveles[nivel] = itemsDe(actual)
	}
	return niveles, nil
}

// hijosDe lee en paralelo los nodos activos debajo de cada padre, en el orden del árbol.
func (r *RepositorioConInstrumentos) hijosDe(padres []nodoInstrumento) ([][]nodoInstrumento, error) {
	hijos := make([][]nodoInstrumento, len(padres))
	fallas := make([]error, len(padres))
	var espera sync.WaitGroup
	for i, padre := range padres {
		if !padre.conHijos {
			continue
		}
		espera.Add(1)
		go func(i int, padre string) {
			defer espera.Done()
			r.cupos <- struct{}{}
			subgrupos, err := r.crud.HijosDe(padre)
			<-r.cupos
			if err != nil {
				fallas[i] = err
				return
			}
			hijos[i] = nodosActivos(subgrupos, padre)
		}(i, padre.item.Codigo)
	}
	espera.Wait()
	if err := errors.Join(fallas...); err != nil {
		return nil, err
	}
	return hijos, nil
}

// nodosActivos convierte los subgrupos activos en elementos del catálogo. Respeta la marca de
// orden de Construcción ([[orden:N]]) y la quita de la descripción; si la descripción repite el
// nombre, no la lleva.
func nodosActivos(subgrupos []Subgrupo, padre string) []nodoInstrumento {
	activos := make([]Subgrupo, 0, len(subgrupos))
	for _, subgrupo := range subgrupos {
		if subgrupo.Activo {
			activos = append(activos, subgrupo)
		}
	}
	ordenarPorMarca(activos)
	nodos := make([]nodoInstrumento, 0, len(activos))
	for _, subgrupo := range activos {
		nombre := sinEspaciosDeMas(subgrupo.Nombre)
		if nombre == "" {
			continue
		}
		descripcion := sinEspaciosDeMas(patronOrden.ReplaceAllString(subgrupo.Descripcion, ""))
		if descripcion == nombre {
			descripcion = ""
		}
		nodos = append(nodos, nodoInstrumento{
			item:     ItemCatalogo{Codigo: subgrupo.ID, Nombre: nombre, Padre: padre, Descripcion: descripcion},
			conHijos: len(subgrupo.Hijos) > 0,
		})
	}
	return nodos
}

func itemsDe(nodos []nodoInstrumento) []ItemCatalogo {
	items := make([]ItemCatalogo, 0, len(nodos))
	for _, nodo := range nodos {
		items = append(items, nodo.item)
	}
	return items
}

func sinEspaciosDeMas(texto string) string {
	return strings.Join(strings.Fields(texto), " ")
}

// nombreDePlan es el nombre del plan con sus vigencias (p. ej. «Plan Estratégico de Desarrollo
// (2021–2025)»), que distinguen planes con el mismo nombre.
func nombreDePlan(plan Plan) string {
	nombre := sinEspaciosDeMas(plan.Nombre)
	if vigencias := textoDeVigencias(plan.VigenciaAplica); vigencias != "" {
		return nombre + " (" + vigencias + ")"
	}
	return nombre
}

// textoDeVigencias resume vigencia_aplica ([{"Id":3,"Nombre":"2021"},…]): un rango si son años
// seguidos («2021–2025») y, si no, la lista («2021, 2023»). Vacío si no tiene o no se entiende.
func textoDeVigencias(vigenciaAplica string) string {
	var vigencias []struct {
		Nombre string `json:"Nombre"`
	}
	if strings.TrimSpace(vigenciaAplica) == "" || json.Unmarshal([]byte(vigenciaAplica), &vigencias) != nil {
		return ""
	}
	nombres := make([]string, 0, len(vigencias))
	anios := make([]int, 0, len(vigencias))
	vistos := map[string]bool{}
	for _, vigencia := range vigencias {
		nombre := strings.TrimSpace(vigencia.Nombre)
		if nombre == "" || vistos[nombre] {
			continue
		}
		vistos[nombre] = true
		nombres = append(nombres, nombre)
		if anio, err := strconv.Atoi(nombre); err == nil {
			anios = append(anios, anio)
		}
	}
	if len(nombres) == 0 {
		return ""
	}
	if len(anios) != len(nombres) {
		return strings.Join(nombres, ", ")
	}
	sort.Ints(anios)
	for i := 1; i < len(anios); i++ {
		if anios[i] != anios[i-1]+1 {
			textos := make([]string, len(anios))
			for j, anio := range anios {
				textos[j] = strconv.Itoa(anio)
			}
			return strings.Join(textos, ", ")
		}
	}
	if len(anios) == 1 {
		return strconv.Itoa(anios[0])
	}
	return fmt.Sprintf("%d–%d", anios[0], anios[len(anios)-1])
}
