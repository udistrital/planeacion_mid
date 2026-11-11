package iniciativaInversionHelper

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"
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

func TestCacheEvitaRecargarHastaQueVence(t *testing.T) {
	arbol := nuevoArbolDePrueba(true)
	cache := NuevaCacheDefiniciones(NuevoCargadorDefiniciones(arbol, arbol, 2), time.Minute).(*definicionesEnCache)
	ahora := time.Now()
	cache.ahora = func() time.Time { return ahora }

	for i := 0; i < 3; i++ {
		if _, err := cache.Definicion("plantilla"); err != nil {
			t.Fatal(err)
		}
	}
	consultasIniciales := atomic.LoadInt32(&arbol.consultasHijo)

	ahora = ahora.Add(2 * time.Minute)
	if _, err := cache.Definicion("plantilla"); err != nil {
		t.Fatal(err)
	}
	if consultasIniciales != 2 || atomic.LoadInt32(&arbol.consultasHijo) != 4 {
		t.Fatalf("consultas: iniciales=%d finales=%d", consultasIniciales, atomic.LoadInt32(&arbol.consultasHijo))
	}
}
