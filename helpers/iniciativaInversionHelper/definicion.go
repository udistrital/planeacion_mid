package iniciativaInversionHelper

import (
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// FuenteDefiniciones entrega la definición (árbol de campos) de una plantilla.
type FuenteDefiniciones interface {
	Definicion(plantillaID string) (*Definicion, error)
}

// patronOrden es la marca que el cliente pone al final de la descripción de los
// nodos de primer nivel para fijar su orden (plan/utils/orden-descripcion.ts).
var patronOrden = regexp.MustCompile(`\s*\[\[orden:(\d+)\]\]\s*$`)

// cargadorDefiniciones arma la definición leyendo el árbol desde planes_crud.
// Pide los hijos y los detalles en paralelo, con un máximo de llamadas
// simultáneas para no saturar el crud. La definición es de solo lectura.
type cargadorDefiniciones struct {
	planes RepositorioPlanes
	arbol  RepositorioArbol
	cupos  chan struct{}
}

// NuevoCargadorDefiniciones crea un cargador con a lo sumo maxSimultaneas llamadas a la vez.
func NuevoCargadorDefiniciones(planes RepositorioPlanes, arbol RepositorioArbol, maxSimultaneas int) FuenteDefiniciones {
	if maxSimultaneas < 1 {
		maxSimultaneas = 1
	}
	return &cargadorDefiniciones{planes: planes, arbol: arbol, cupos: make(chan struct{}, maxSimultaneas)}
}

func (c *cargadorDefiniciones) Definicion(plantillaID string) (*Definicion, error) {
	plantilla, err := c.planes.ObtenerPlan(plantillaID)
	if err != nil {
		return nil, err
	}
	if !plantilla.Formato {
		return nil, errorSolicitud("La plantilla indicada no es válida.").conDetalle("el plan %s no es una plantilla (formato: false)", plantillaID)
	}
	modulos, err := c.cargarHijos(plantillaID)
	if err != nil {
		return nil, err
	}
	definicion := &Definicion{
		PlantillaID: plantilla.ID,
		Nombre:      plantilla.Nombre,
		TipoPlanID:  plantilla.TipoPlanID,
		Modulos:     modulos,
	}
	if err := validarModulos(definicion); err != nil {
		return nil, err
	}
	return definicion, nil
}

// validarModulos exige que el primer nivel sean módulos (dato.modulo) y que las
// claves de los submódulos no se repitan entre módulos.
func validarModulos(definicion *Definicion) error {
	claves := map[string]bool{}
	for _, modulo := range definicion.Modulos {
		if !modulo.Dato.Modulo {
			return errorConfiguracion("la plantilla %s no está organizada por módulos: «%s» debe tener dato.modulo = true",
				definicion.PlantillaID, modulo.Nombre)
		}
		for _, submodulo := range modulo.Hijos {
			if claves[submodulo.Clave()] {
				return errorConfiguracion("la plantilla %s repite la clave de submódulo %q", definicion.PlantillaID, submodulo.Clave())
			}
			claves[submodulo.Clave()] = true
		}
	}
	return nil
}

func (c *cargadorDefiniciones) cargarHijos(padreID string) ([]*Nodo, error) {
	var subgrupos []Subgrupo
	err := c.conCupo(func() (err error) {
		subgrupos, err = c.arbol.HijosDe(padreID)
		return err
	})
	if err != nil {
		return nil, err
	}

	activos := make([]Subgrupo, 0, len(subgrupos))
	for _, subgrupo := range subgrupos {
		if subgrupo.Activo {
			activos = append(activos, subgrupo)
		}
	}
	ordenarPorMarca(activos)

	nodos := make([]*Nodo, len(activos))
	fallas := make([]error, len(activos))
	var espera sync.WaitGroup
	for i := range activos {
		espera.Add(1)
		go func(i int) {
			defer espera.Done()
			nodos[i], fallas[i] = c.cargarNodo(activos[i])
		}(i)
	}
	espera.Wait()
	if err := errors.Join(fallas...); err != nil {
		return nil, err
	}
	return nodos, nil
}

func (c *cargadorDefiniciones) cargarNodo(subgrupo Subgrupo) (*Nodo, error) {
	dato, err := c.cargarDato(subgrupo.ID)
	if err != nil {
		return nil, err
	}
	nodo := &Nodo{
		ID:          subgrupo.ID,
		Nombre:      subgrupo.Nombre,
		Descripcion: patronOrden.ReplaceAllString(subgrupo.Descripcion, ""),
		Dato:        dato,
	}
	if len(subgrupo.Hijos) > 0 {
		if nodo.Hijos, err = c.cargarHijos(subgrupo.ID); err != nil {
			return nil, err
		}
	}
	return nodo, nil
}

func (c *cargadorDefiniciones) cargarDato(subgrupoID string) (Dato, error) {
	var detalles []DetalleSubgrupo
	err := c.conCupo(func() (err error) {
		detalles, err = c.arbol.DetallesDe(subgrupoID)
		return err
	})
	if err != nil {
		return Dato{}, err
	}
	for _, detalle := range detalles {
		if detalle.Activo {
			return interpretarDato(subgrupoID, detalle.Dato)
		}
	}
	return Dato{}, nil
}

// interpretarDato convierte el texto JSON de subgrupo_detalle.dato. Un texto en blanco
// o "{}" da un Dato sin tipo: el nodo es un grupo, no un campo.
func interpretarDato(subgrupoID, texto string) (Dato, error) {
	var dato Dato
	if strings.TrimSpace(texto) == "" {
		return dato, nil
	}
	if err := json.Unmarshal([]byte(texto), &dato); err != nil {
		return Dato{}, errorConfiguracion("la definición del nodo %s (subgrupo_detalle.dato) no es un JSON válido: %v", subgrupoID, err)
	}
	return dato, nil
}

func (c *cargadorDefiniciones) conCupo(llamada func() error) error {
	c.cupos <- struct{}{}
	defer func() { <-c.cupos }()
	return llamada()
}

// ordenarPorMarca respeta [[orden:N]] cuando existe; si no, conserva el orden recibido.
func ordenarPorMarca(subgrupos []Subgrupo) {
	sort.SliceStable(subgrupos, func(i, j int) bool {
		return ordenDeMarca(subgrupos[i].Descripcion) < ordenDeMarca(subgrupos[j].Descripcion)
	})
}

func ordenDeMarca(descripcion string) int {
	coincidencia := patronOrden.FindStringSubmatch(descripcion)
	if coincidencia == nil {
		return int(^uint(0) >> 1) // sin marca: al final, en el orden recibido
	}
	orden, _ := strconv.Atoi(coincidencia[1])
	return orden
}

// ---------------------------------------------------------------------------
// Caché
// ---------------------------------------------------------------------------

// definicionesEnCache guarda cada definición hasta que el mid se reinicia. Una
// plantilla con iniciativas no se edita (si cambia el instrumento se carga una
// nueva, con otro id), así que el árbol se recorre una sola vez por plantilla.
// Si varias peticiones piden a la vez una plantilla que no está en la caché,
// solo una la lee y las demás esperan ese resultado. Un error no se guarda: la
// siguiente petición vuelve a leer.
type definicionesEnCache struct {
	fuente FuenteDefiniciones

	mu           sync.Mutex
	definiciones map[string]*Definicion
	enCurso      map[string]*lecturaDefinicion
}

// lecturaDefinicion es la lectura de una plantilla que está en curso; listo se
// cierra cuando termina.
type lecturaDefinicion struct {
	listo      chan struct{}
	definicion *Definicion
	err        error
}

// NuevaCacheDefiniciones envuelve una fuente con caché en memoria.
func NuevaCacheDefiniciones(fuente FuenteDefiniciones) FuenteDefiniciones {
	return &definicionesEnCache{
		fuente:       fuente,
		definiciones: make(map[string]*Definicion),
		enCurso:      make(map[string]*lecturaDefinicion),
	}
}

func (c *definicionesEnCache) Definicion(plantillaID string) (*Definicion, error) {
	c.mu.Lock()
	if definicion, existe := c.definiciones[plantillaID]; existe {
		c.mu.Unlock()
		return definicion, nil
	}
	if lectura, existe := c.enCurso[plantillaID]; existe {
		c.mu.Unlock()
		<-lectura.listo
		return lectura.definicion, lectura.err
	}
	// El error inicial es el que reciben quienes esperan si la lectura se
	// interrumpe con un pánico antes de devolver un resultado.
	lectura := &lecturaDefinicion{
		listo: make(chan struct{}),
		err:   errors.New("se interrumpió la lectura de la plantilla " + plantillaID),
	}
	c.enCurso[plantillaID] = lectura
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.enCurso, plantillaID)
		if lectura.err == nil {
			c.definiciones[plantillaID] = lectura.definicion
		}
		c.mu.Unlock()
		close(lectura.listo)
	}()
	lectura.definicion, lectura.err = c.fuente.Definicion(plantillaID)
	return lectura.definicion, lectura.err
}
