package iniciativaInversionHelper

import (
	"regexp"
	"sync"
	"time"
)

// ItemCatalogo es un elemento de un catálogo de las iniciativas. Padre es el código del
// elemento del que depende (p. ej. el sector de un programa).
type ItemCatalogo struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
	Padre  string `json:"padre,omitempty"`
}

// Catalogo es una lista de elementos con búsqueda por código.
type Catalogo struct {
	Nombre    string
	Items     []ItemCatalogo
	porCodigo map[string]ItemCatalogo
}

// NuevoCatalogo arma el catálogo con su índice por código.
func NuevoCatalogo(nombre string, items []ItemCatalogo) *Catalogo {
	catalogo := &Catalogo{Nombre: nombre, Items: items, porCodigo: make(map[string]ItemCatalogo, len(items))}
	for _, item := range items {
		catalogo.porCodigo[item.Codigo] = item
	}
	return catalogo
}

// Buscar devuelve el elemento con ese código.
func (c *Catalogo) Buscar(codigo string) (ItemCatalogo, bool) {
	item, existe := c.porCodigo[codigo]
	return item, existe
}

// Filtrar devuelve los elementos cuyo padre es el dado; sin padre, todos.
func (c *Catalogo) Filtrar(padre string) []ItemCatalogo {
	if padre == "" {
		return c.Items
	}
	filtrados := make([]ItemCatalogo, 0)
	for _, item := range c.Items {
		if item.Padre == padre {
			filtrados = append(filtrados, item)
		}
	}
	return filtrados
}

// ProveedorCatalogos entrega catálogos por nombre.
type ProveedorCatalogos interface {
	Catalogo(nombre string) (*Catalogo, error)
}

// RepositorioElementosCatalogo lee los elementos activos de un catálogo, en su
// orden (colección general elemento_catalogo de planes_crud).
type RepositorioElementosCatalogo interface {
	ElementosDeCatalogo(nombre string) ([]ItemCatalogo, error)
}

var patronNombreCatalogo = regexp.MustCompile(`^[a-z0-9_]+$`)

// CatalogosConCache entrega los catálogos de un repositorio y guarda cada uno en
// memoria durante un tiempo: el validador consulta los catálogos en cada guardado
// y cambian muy poco. Un catálogo vacío o un error no se guardan.
type CatalogosConCache struct {
	repo     RepositorioElementosCatalogo
	vigencia time.Duration
	ahora    func() time.Time

	mu       sync.Mutex
	cargados map[string]catalogoEnCache
}

type catalogoEnCache struct {
	catalogo *Catalogo
	vence    time.Time
}

// NuevosCatalogosConCache crea el proveedor; vigencia es cuánto dura cada catálogo en memoria.
func NuevosCatalogosConCache(repo RepositorioElementosCatalogo, vigencia time.Duration) *CatalogosConCache {
	return &CatalogosConCache{repo: repo, vigencia: vigencia, ahora: time.Now, cargados: make(map[string]catalogoEnCache)}
}

// Catalogo devuelve el catálogo, de la memoria si no ha vencido.
func (p *CatalogosConCache) Catalogo(nombre string) (*Catalogo, error) {
	if !patronNombreCatalogo.MatchString(nombre) {
		return nil, errorFormatoSolicitud("nombre de catálogo inválido: %q", nombre)
	}
	p.mu.Lock()
	enCache, existe := p.cargados[nombre]
	p.mu.Unlock()
	if existe && p.ahora().Before(enCache.vence) {
		return enCache.catalogo, nil
	}

	items, err := p.repo.ElementosDeCatalogo(nombre)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errorNoEncontrado("No se encontraron las opciones de este campo. Comuníquese con el administrador del sistema.").
			conDetalle("el catálogo %s no existe en elemento_catalogo o no tiene elementos activos. Se carga con database local/construccion-proyecto-inversion/cargar-catalogos.js.", nombre)
	}
	catalogo := NuevoCatalogo(nombre, items)
	p.mu.Lock()
	p.cargados[nombre] = catalogoEnCache{catalogo: catalogo, vence: p.ahora().Add(p.vigencia)}
	p.mu.Unlock()
	return catalogo, nil
}
