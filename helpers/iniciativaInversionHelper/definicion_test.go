package iniciativaInversionHelper

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// arbolDePrueba sirve un árbol pequeño:
//
//	plan -> [B (módulo, orden 1), A (módulo, orden 0), C inactivo]
//	A    -> [a1 (campo de texto), a2 (grupo con dato "{}")]
type arbolDePrueba struct {
	planes        map[string]*Plan
	datos         map[string]string
	consultasHijo int32
}

func nuevoArbolDePrueba(formato bool) *arbolDePrueba {
	return &arbolDePrueba{
		planes: map[string]*Plan{"plantilla": {ID: "plantilla", Formato: formato, TipoPlanID: "t"}},
		datos: map[string]string{
			"A":  `{"clave":"primera","modulo":true}`,
			"B":  `{"clave":"segunda","modulo":true}`,
			"a1": `{"clave":"campo","type":"input","required":"true"}`,
			"a2": "{}",
		},
	}
}

func (a *arbolDePrueba) ObtenerPlan(id string) (*Plan, error) {
	if plan, existe := a.planes[id]; existe {
		return plan, nil
	}
	return nil, errorNoEncontrado("no existe %s", id)
}
func (a *arbolDePrueba) BuscarPlanes(map[string]string) ([]Plan, error) { return nil, nil }
func (a *arbolDePrueba) CrearPlan(Plan) (*Plan, error)                  { return nil, nil }
func (a *arbolDePrueba) ActualizarPlan(Plan) error                      { return nil }
func (a *arbolDePrueba) DesactivarPlan(string) error                    { return nil }

func (a *arbolDePrueba) HijosDe(id string) ([]Subgrupo, error) {
	atomic.AddInt32(&a.consultasHijo, 1)
	switch id {
	case "plantilla":
		return []Subgrupo{
			{ID: "B", Nombre: "Segundo", Descripcion: "Va segundo [[orden:0001]]", Activo: true},
			{ID: "A", Nombre: "Primero", Descripcion: "Va primero [[orden:0000]]", Activo: true, Hijos: []string{"a1", "a2"}},
			{ID: "C", Nombre: "Inactivo", Activo: false},
		}, nil
	case "A":
		return []Subgrupo{{ID: "a1", Nombre: "Campo", Activo: true}, {ID: "a2", Nombre: "Grupo", Activo: true}}, nil
	}
	return nil, nil
}

func (a *arbolDePrueba) DetallesDe(id string) ([]DetalleSubgrupo, error) {
	return []DetalleSubgrupo{{Dato: a.datos[id], Activo: true}}, nil
}

func TestCargadorRespetaElOrdenYOmiteInactivos(t *testing.T) {
	arbol := nuevoArbolDePrueba(true)
	definicion, err := NuevoCargadorDefiniciones(arbol, arbol, 2).Definicion("plantilla")
	if err != nil {
		t.Fatal(err)
	}
	if len(definicion.Modulos) != 2 || definicion.Modulos[0].ID != "A" || definicion.Modulos[1].ID != "B" {
		t.Fatalf("orden o filtro incorrecto: %+v", definicion.Modulos)
	}
	if definicion.Modulos[0].Descripcion != "Va primero" {
		t.Fatalf("no se quitó la marca de orden: %q", definicion.Modulos[0].Descripcion)
	}
	campo := definicion.Submodulo("campo")
	if campo == nil || campo.Dato.Tipo != TipoTexto || !bool(campo.Dato.Requerido) {
		t.Fatalf("no se interpretó el dato del campo: %+v", campo)
	}
	if definicion.Modulos[0].Hijos[1].EsCampo() {
		t.Fatal(`un dato "{}" es un grupo, no un campo`)
	}
	if len(definicion.Submodulos()) != 2 {
		t.Fatalf("los submódulos son los hijos de los módulos: %d", len(definicion.Submodulos()))
	}
}

func TestCargadorRechazaPlanesQueNoSonPlantilla(t *testing.T) {
	arbol := nuevoArbolDePrueba(false)
	if _, err := NuevoCargadorDefiniciones(arbol, arbol, 2).Definicion("plantilla"); err == nil {
		t.Fatal("se esperaba un error")
	}
}

func TestCargadorExigeModulosEnElPrimerNivel(t *testing.T) {
	arbol := nuevoArbolDePrueba(true)
	arbol.datos["B"] = `{"clave":"segunda"}` // sin "modulo": true
	_, err := NuevoCargadorDefiniciones(arbol, arbol, 2).Definicion("plantilla")
	if err == nil || !strings.Contains(err.Error(), "módulos") {
		t.Fatalf("se esperaba el error de módulos y llegó %v", err)
	}
}

func TestLasClavesDeSubmoduloNoSeRepitenEntreModulos(t *testing.T) {
	submodulo := func(clave string) *Nodo { return &Nodo{Dato: Dato{Clave: clave}} }
	definicion := &Definicion{Modulos: []*Nodo{
		{Nombre: "Uno", Dato: Dato{Clave: "uno", Modulo: true}, Hijos: []*Nodo{submodulo("nombre")}},
		{Nombre: "Dos", Dato: Dato{Clave: "dos", Modulo: true}, Hijos: []*Nodo{submodulo("nombre")}},
	}}
	if err := validarModulos(definicion); err == nil {
		t.Fatal("se esperaba un error por clave repetida")
	}
}

// fuenteContada cuenta las lecturas y responde con leer.
type fuenteContada struct {
	lecturas int32
	leer     func(lectura int32) (*Definicion, error)
}

func (f *fuenteContada) Definicion(plantillaID string) (*Definicion, error) {
	return f.leer(atomic.AddInt32(&f.lecturas, 1))
}

func TestCacheLeeCadaPlantillaUnaSolaVez(t *testing.T) {
	arbol := nuevoArbolDePrueba(true)
	cache := NuevaCacheDefiniciones(NuevoCargadorDefiniciones(arbol, arbol, 2))

	for i := 0; i < 3; i++ {
		if _, err := cache.Definicion("plantilla"); err != nil {
			t.Fatal(err)
		}
	}
	// Dos consultas de hijos: las de la plantilla y las del módulo A.
	if consultas := atomic.LoadInt32(&arbol.consultasHijo); consultas != 2 {
		t.Fatalf("se esperaba leer el árbol una vez (2 consultas de hijos) y hubo %d", consultas)
	}
}

func TestCacheComparteLaLecturaEnCurso(t *testing.T) {
	entro := make(chan struct{})
	liberar := make(chan struct{})
	fuente := &fuenteContada{leer: func(lectura int32) (*Definicion, error) {
		if lectura == 1 {
			close(entro)
		}
		<-liberar
		return &Definicion{PlantillaID: "plantilla"}, nil
	}}
	cache := NuevaCacheDefiniciones(fuente)

	const peticiones = 5
	resultados := make([]*Definicion, peticiones)
	fallas := make([]error, peticiones)
	var iniciadas, terminadas sync.WaitGroup
	iniciadas.Add(peticiones)
	terminadas.Add(peticiones)
	for i := 0; i < peticiones; i++ {
		go func(i int) {
			defer terminadas.Done()
			iniciadas.Done()
			resultados[i], fallas[i] = cache.Definicion("plantilla")
		}(i)
	}
	<-entro
	iniciadas.Wait()
	close(liberar)
	terminadas.Wait()

	if lecturas := atomic.LoadInt32(&fuente.lecturas); lecturas != 1 {
		t.Fatalf("se esperaba una sola lectura y hubo %d", lecturas)
	}
	for i := 0; i < peticiones; i++ {
		if fallas[i] != nil || resultados[i] != resultados[0] {
			t.Fatalf("la petición %d no recibió la definición compartida: %v", i, fallas[i])
		}
	}
}

func TestCacheNoGuardaErrores(t *testing.T) {
	fuente := &fuenteContada{leer: func(lectura int32) (*Definicion, error) {
		if lectura == 1 {
			return nil, errorServicio("planes_crud no responde")
		}
		return &Definicion{PlantillaID: "plantilla"}, nil
	}}
	cache := NuevaCacheDefiniciones(fuente)

	if _, err := cache.Definicion("plantilla"); err == nil {
		t.Fatal("se esperaba el error de la primera lectura")
	}
	for i := 0; i < 2; i++ {
		if definicion, err := cache.Definicion("plantilla"); err != nil || definicion == nil {
			t.Fatalf("después del error se esperaba leer de nuevo: %v", err)
		}
	}
	if lecturas := atomic.LoadInt32(&fuente.lecturas); lecturas != 2 {
		t.Fatalf("se esperaban 2 lecturas (la fallida y la buena) y hubo %d", lecturas)
	}
}

func TestCacheSeRecuperaDeUnaLecturaInterrumpida(t *testing.T) {
	entro := make(chan struct{})
	liberar := make(chan struct{})
	fuente := &fuenteContada{leer: func(lectura int32) (*Definicion, error) {
		if lectura == 1 {
			close(entro)
			<-liberar
			panic("falla inesperada")
		}
		return &Definicion{PlantillaID: "plantilla"}, nil
	}}
	cache := NuevaCacheDefiniciones(fuente)

	interrumpida := make(chan interface{})
	go func() {
		defer func() { interrumpida <- recover() }()
		cache.Definicion("plantilla")
	}()
	<-entro
	espera := make(chan error)
	go func() {
		_, err := cache.Definicion("plantilla")
		espera <- err
	}()
	close(liberar)

	if recuperado := <-interrumpida; recuperado == nil {
		t.Fatal("el pánico de la fuente debe llegar a quien hizo la lectura")
	}
	// Quien esperaba recibe un error o, si llegó después de la falla, una lectura nueva.
	if err := <-espera; err != nil && !strings.Contains(err.Error(), "se interrumpió") {
		t.Fatalf("error inesperado para quien esperaba: %v", err)
	}
	if definicion, err := cache.Definicion("plantilla"); err != nil || definicion == nil {
		t.Fatalf("después de la falla se esperaba leer de nuevo: %v", err)
	}
}
