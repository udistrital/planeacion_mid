package iniciativaInversionHelper

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

// crudInstrumentos responde como planes_crud los planes de un tipo y su árbol, y cuenta las llamadas.
type crudInstrumentos struct {
	mu       sync.Mutex
	tipos    map[string]string     // código del tipo de plan -> id
	planes   map[string][]Plan     // id del tipo -> planes activos con formato: true
	hijos    map[string][]Subgrupo // id del padre -> hijos en orden
	filtros  []map[string]string
	llamadas int
	falloEn  string
}

func (c *crudInstrumentos) BuscarPorCodigo(recurso, codigo string) (*Estado, error) {
	c.contar()
	if recurso != RecursoTipoPlan || c.tipos[codigo] == "" {
		return nil, errorNoEncontrado("No existe %s %s.", recurso, codigo)
	}
	return &Estado{ID: c.tipos[codigo], Codigo: codigo}, nil
}

func (c *crudInstrumentos) BuscarPlanes(filtros map[string]string) ([]Plan, error) {
	c.contar()
	c.mu.Lock()
	c.filtros = append(c.filtros, filtros)
	c.mu.Unlock()
	return c.planes[filtros["tipo_plan_id"]], nil
}

func (c *crudInstrumentos) HijosDe(id string) ([]Subgrupo, error) {
	c.contar()
	if id == c.falloEn {
		return nil, errors.New("crud caído")
	}
	return c.hijos[id], nil
}

func (c *crudInstrumentos) contar() {
	c.mu.Lock()
	c.llamadas++
	c.mu.Unlock()
}

func nodoPlan(id, nombre, descripcion string, hijos ...string) Subgrupo {
	return Subgrupo{ID: id, Nombre: nombre, Descripcion: descripcion, Activo: true, Hijos: hijos}
}

// crudConPED tiene un PED con dos lineamientos y sus metas, otro PED sin nodos, un PI y un PDD.
func crudConPED() *crudInstrumentos {
	inactivo := nodoPlan("l-x", "Lineamiento retirado", "")
	inactivo.Activo = false
	return &crudInstrumentos{
		tipos: map[string]string{CodigoTipoPlanEstrategico: "tipo-ped", CodigoTipoPlanIndicativo: "tipo-pi", CodigoTipoPlanDistrital: "tipo-pdd"},
		planes: map[string][]Plan{
			"tipo-ped": {
				{ID: "ped-1", Nombre: " Plan Estratégico  de Desarrollo", VigenciaAplica: `[{"Id":3,"Nombre":"2021"},{"Id":25,"Nombre":"2022"},{"Id":35,"Nombre":"2023"}]`},
				{ID: "ped-vacio", Nombre: "Plan sin nodos"},
			},
			"tipo-pi":  {{ID: "pi-1", Nombre: "Plan Indicativo", VigenciaAplica: `[{"Id":58,"Nombre":"2025"}]`}},
			"tipo-pdd": {{ID: "pdd-1", Nombre: "Plan Distrital de Desarrollo"}},
		},
		hijos: map[string][]Subgrupo{
			"ped-1": {
				nodoPlan("l-2", "Lineamiento 2. Investigar", "Fortalecer la investigación [[orden:2]]", "m-3"),
				inactivo,
				nodoPlan("l-1", "Lineamiento 1. Formar", "Formar ciudadanos [[orden:1]]", "m-1", "m-2"),
			},
			"l-1": {
				nodoPlan("m-1", "Meta 1", "Mantener la acreditación", "e-1"),
				nodoPlan("m-2", "Meta 2", "Meta 2"),
			},
			"l-2":   {nodoPlan("m-3", "Meta 3", "Aumentar los grupos de investigación")},
			"pi-1":  {nodoPlan("eje-1", "Eje transformador 1", "", "lin-1")},
			"lin-1": {nodoPlan("est-1", "Estrategia 1.1.1", "")},
			"eje-1": {nodoPlan("lin-1", "Lineamiento 1.1", "", "est-1")},
			"pdd-1": {nodoPlan("obj-1", "Objetivo estratégico 1", "", "edd-1")},
			"obj-1": {nodoPlan("edd-1", "Estrategia 1", "", "prog-1")},
			"edd-1": {nodoPlan("prog-1", "Programa intersectorial 1", "")},
		},
	}
}

func TestLosCatalogosDelPEDSalenDeSuArbolPorNivel(t *testing.T) {
	crud := crudConPED()
	repo := NuevoRepositorioConInstrumentos(&repoCatalogosContado{}, crud, time.Minute, 4)

	planes, err := repo.ElementosDeCatalogo("planes_ped")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(planes, []ItemCatalogo{{Codigo: "ped-1", Nombre: "Plan Estratégico de Desarrollo (2021–2023)"}}) {
		t.Fatalf("los planes sin nodos no se ofrecen y el nombre lleva las vigencias: %v", planes)
	}
	lineamientos, _ := repo.ElementosDeCatalogo("lineamientos_ped")
	esperados := []ItemCatalogo{
		{Codigo: "l-1", Nombre: "Lineamiento 1. Formar", Padre: "ped-1", Descripcion: "Formar ciudadanos"},
		{Codigo: "l-2", Nombre: "Lineamiento 2. Investigar", Padre: "ped-1", Descripcion: "Fortalecer la investigación"},
	}
	if !reflect.DeepEqual(lineamientos, esperados) {
		t.Fatalf("lineamientos activos, en el orden de la marca y sin ella:\n esperados %v\n llegaron  %v", esperados, lineamientos)
	}
	metas, _ := repo.ElementosDeCatalogo("metas_ped")
	esperadas := []ItemCatalogo{
		{Codigo: "m-1", Nombre: "Meta 1", Padre: "l-1", Descripcion: "Mantener la acreditación"},
		{Codigo: "m-2", Nombre: "Meta 2", Padre: "l-1"},
		{Codigo: "m-3", Nombre: "Meta 3", Padre: "l-2", Descripcion: "Aumentar los grupos de investigación"},
	}
	if !reflect.DeepEqual(metas, esperadas) {
		t.Fatalf("las metas van con su lineamiento como padre y sin repetir el nombre en la descripción:\n esperadas %v\n llegaron  %v", esperadas, metas)
	}
	filtros := crud.filtros[0]
	if filtros["tipo_plan_id"] != "tipo-ped" || filtros["formato"] != "true" || filtros["activo"] != "true" {
		t.Fatalf("solo se leen los planes activos armados en Construcción (formato: true): %v", filtros)
	}
}

func TestElPIllegaHastaLasEstrategias(t *testing.T) {
	repo := NuevoRepositorioConInstrumentos(&repoCatalogosContado{}, crudConPED(), time.Minute, 4)

	esperados := map[string][]ItemCatalogo{
		"planes_pi":       {{Codigo: "pi-1", Nombre: "Plan Indicativo (2025)"}},
		"ejes_pi":         {{Codigo: "eje-1", Nombre: "Eje transformador 1", Padre: "pi-1"}},
		"lineamientos_pi": {{Codigo: "lin-1", Nombre: "Lineamiento 1.1", Padre: "eje-1"}},
		"estrategias_pi":  {{Codigo: "est-1", Nombre: "Estrategia 1.1.1", Padre: "lin-1"}},
	}
	for catalogo, items := range esperados {
		llegaron, err := repo.ElementosDeCatalogo(catalogo)
		if err != nil || !reflect.DeepEqual(llegaron, items) {
			t.Fatalf("%s: se esperaba %v y llegó %v (err %v)", catalogo, items, llegaron, err)
		}
	}
}

func TestElPDDllegaHastaLosProgramasIntersectoriales(t *testing.T) {
	repo := NuevoRepositorioConInstrumentos(&repoCatalogosContado{}, crudConPED(), time.Minute, 4)

	esperados := map[string][]ItemCatalogo{
		"planes_pdd":      {{Codigo: "pdd-1", Nombre: "Plan Distrital de Desarrollo"}},
		"objetivos_pdd":   {{Codigo: "obj-1", Nombre: "Objetivo estratégico 1", Padre: "pdd-1"}},
		"estrategias_pdd": {{Codigo: "edd-1", Nombre: "Estrategia 1", Padre: "obj-1"}},
		"programas_pdd":   {{Codigo: "prog-1", Nombre: "Programa intersectorial 1", Padre: "edd-1"}},
	}
	for catalogo, items := range esperados {
		llegaron, err := repo.ElementosDeCatalogo(catalogo)
		if err != nil || !reflect.DeepEqual(llegaron, items) {
			t.Fatalf("%s: se esperaba %v y llegó %v (err %v)", catalogo, items, llegaron, err)
		}
	}
}

func TestLosCatalogosDeUnInstrumentoSalenDeUnaSolaLectura(t *testing.T) {
	crud := crudConPED()
	repo := NuevoRepositorioConInstrumentos(&repoCatalogosContado{}, crud, time.Minute, 4)
	ahora := time.Date(2026, 11, 19, 8, 0, 0, 0, time.UTC)
	repo.ahora = func() time.Time { return ahora }

	var espera sync.WaitGroup
	for _, catalogo := range []string{"planes_ped", "lineamientos_ped", "metas_ped", "metas_ped"} {
		espera.Add(1)
		go func(catalogo string) {
			defer espera.Done()
			if _, err := repo.ElementosDeCatalogo(catalogo); err != nil {
				t.Error(err)
			}
		}(catalogo)
	}
	espera.Wait()
	lectura := crud.llamadas
	if len(crud.filtros) != 1 {
		t.Fatalf("las consultas juntas debían esperar una sola lectura; se leyeron los planes %d veces", len(crud.filtros))
	}

	ahora = ahora.Add(2 * time.Minute)
	if _, err := repo.ElementosDeCatalogo("planes_ped"); err != nil {
		t.Fatal(err)
	}
	if crud.llamadas != 2*lectura {
		t.Fatalf("vencido el tiempo, el árbol se lee de nuevo: %d llamadas tras una lectura de %d", crud.llamadas, lectura)
	}
}

func TestLosDemasCatalogosNoLeenPlanes(t *testing.T) {
	crud := crudConPED()
	base := &repoCatalogosContado{items: []ItemCatalogo{{Codigo: "2202001", Nombre: "Documentos de planeación"}}}
	repo := NuevoRepositorioConInstrumentos(base, crud, time.Minute, 4)

	items, err := repo.ElementosDeCatalogo("productos")
	if err != nil || len(items) != 1 || base.consultas != 1 || crud.llamadas != 0 {
		t.Fatalf("debía consultar solo el repositorio base: items %v, err %v, consultas %d, llamadas a planes %d", items, err, base.consultas, crud.llamadas)
	}
}

func TestSinPlanesConNodosNiFallasSeInventanOpciones(t *testing.T) {
	sinNodos := crudConPED()
	sinNodos.hijos = map[string][]Subgrupo{}
	if items, err := NuevoRepositorioConInstrumentos(&repoCatalogosContado{}, sinNodos, time.Minute, 4).ElementosDeCatalogo("planes_ped"); err == nil {
		t.Fatalf("sin planes con nodos se esperaba un error y llegaron %v", items)
	}

	conFalla := crudConPED()
	conFalla.falloEn = "l-2"
	if items, err := NuevoRepositorioConInstrumentos(&repoCatalogosContado{}, conFalla, time.Minute, 4).ElementosDeCatalogo("metas_ped"); err == nil {
		t.Fatalf("si el crud falla en un nodo se esperaba un error y llegaron %v", items)
	}

	sinTipo := crudConPED()
	delete(sinTipo.tipos, CodigoTipoPlanIndicativo)
	if items, err := NuevoRepositorioConInstrumentos(&repoCatalogosContado{}, sinTipo, time.Minute, 4).ElementosDeCatalogo("ejes_pi"); err == nil {
		t.Fatalf("sin el tipo de plan se esperaba un error y llegaron %v", items)
	}
}

// Como en la armonización de Formulación, en cada nivel se seleccionan varios elementos (multiple),
// cada uno debajo de alguno de los seleccionados en el nivel de arriba.
func TestLaValidacionEncadenaPlanLineamientosYMetas(t *testing.T) {
	repo := NuevoRepositorioConInstrumentos(&repoCatalogosContado{}, crudConPED(), time.Minute, 4)
	validador := NuevoValidador(NuevosCatalogosConCache(repo, time.Hour))
	submodulo := &Nodo{Nombre: "Articulación", Dato: Dato{Clave: "articulacion"}, Hijos: []*Nodo{
		{Nombre: "Plan Estratégico de Desarrollo", Dato: Dato{Clave: "ped"}, Hijos: []*Nodo{
			campoDePrueba("plan", TipoCatalogo, Dato{Requerido: true, Catalogo: "planes_ped", SinCodigo: true}),
			campoDePrueba("lineamiento_estrategico", TipoCatalogo,
				Dato{Requerido: true, Catalogo: "lineamientos_ped", DependeDe: "plan", SinCodigo: true, Multiple: true}),
			campoDePrueba("meta_estrategica", TipoCatalogo,
				Dato{Requerido: true, Catalogo: "metas_ped", DependeDe: "lineamiento_estrategico", SinCodigo: true, Multiple: true}),
		}},
	}}
	lista := func(codigos ...string) []interface{} {
		valores := make([]interface{}, len(codigos))
		for i, codigo := range codigos {
			valores[i] = codigo
		}
		return valores
	}
	validar := func(plan string, lineamientos, metas []interface{}) []ErrorCampo {
		return validador.ValidarSubmodulo(submodulo, map[string]Datos{"articulacion": {"ped": Datos{
			"plan": plan, "lineamiento_estrategico": lineamientos, "meta_estrategica": metas,
		}}})
	}

	exigirSinErrores(t, "articulacion", validar("ped-1", lista("l-1"), lista("m-2")))
	exigirSinErrores(t, "articulacion", validar("ped-1", lista("l-1", "l-2"), lista("m-1", "m-2", "m-3")))
	exigirError(t, validar("ped-1", lista("l-1"), lista("m-1", "m-3")), "ped.meta_estrategica", "no pertenece")
	exigirError(t, validar("ped-1", lista(), lista("m-1")), "ped.lineamiento_estrategico", "es obligatorio")
	exigirError(t, validar("ped-1", lista(), lista("m-1")), "ped.meta_estrategica", "primero seleccione")
	exigirError(t, validar("ped-1", lista("l-1"), lista()), "ped.meta_estrategica", "es obligatorio")
	exigirError(t, validar("ped-vacio", lista("l-1"), lista("m-1")), "ped.plan", "no existe en la lista")
}

func TestLasVigenciasSeResumen(t *testing.T) {
	casos := map[string]string{
		`[{"Id":3,"Nombre":"2021"},{"Id":25,"Nombre":"2022"},{"Id":35,"Nombre":"2023"}]`: "2021–2023",
		`[{"Id":35,"Nombre":"2023"},{"Id":3,"Nombre":"2021"}]`:                           "2021, 2023",
		`[{"Id":58,"Nombre":"2025"},{"Id":58,"Nombre":"2025"}]`:                          "2025",
		`[{"Id":1,"Nombre":"Primer semestre"}]`:                                          "Primer semestre",
		``:                                                                               "",
		`no es JSON`:                                                                     "",
	}
	for vigencias, esperado := range casos {
		if texto := textoDeVigencias(vigencias); texto != esperado {
			t.Errorf("%q: se esperaba %q y llegó %q", vigencias, esperado, texto)
		}
	}
}
