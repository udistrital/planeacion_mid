package iniciativaInversionHelper

import (
	"regexp"
	"strings"
)

// Normalizador deja los datos de un submódulo con la forma exacta de la
// plantilla: descarta claves desconocidas, recorta textos y resuelve los campos
// calculados. No valida; eso lo hace el Validador.
type Normalizador struct {
	catalogos ProveedorCatalogos
}

// NuevoNormalizador crea un normalizador que usa los catálogos para mostrar
// nombres en los campos calculados.
func NuevoNormalizador(catalogos ProveedorCatalogos) *Normalizador {
	return &Normalizador{catalogos: catalogos}
}

// Normalizar devuelve una copia limpia de los datos del submódulo.
func (n *Normalizador) Normalizar(submodulo *Nodo, datos Datos) Datos {
	return n.objeto(submodulo.Hijos, datos)
}

func (n *Normalizador) objeto(nodos []*Nodo, datos Datos) Datos {
	salida := Datos{}
	for _, nodo := range nodos {
		clave := nodo.Clave()
		if clave == "" || nodo.Dato.Tipo == TipoCalculado {
			continue
		}
		valor := datos[clave]
		switch {
		case !nodo.EsCampo():
			subdatos, _ := comoObjeto(valor)
			salida[clave] = n.objeto(nodo.Hijos, subdatos)
		case nodo.Dato.Tipo == TipoRepetible:
			salida[clave] = n.filas(nodo, valor)
		case nodo.Dato.Tipo == TipoMatriz:
			salida[clave] = n.matriz(nodo, valor)
		case valor != nil:
			salida[clave] = limpiarEscalar(valor)
		}
	}
	// Los calculados van al final porque usan los valores ya limpios de sus hermanos.
	for _, nodo := range nodos {
		if nodo.Dato.Tipo == TipoCalculado && nodo.Clave() != "" {
			salida[nodo.Clave()] = n.valorCalculado(nodo, nodos, salida, datos[nodo.Clave()])
		}
	}
	return salida
}

// valorCalculado respeta lo que escribió el usuario si el campo es editable;
// si no es editable o lo dejó vacío, usa el valor calculado con los hermanos.
func (n *Normalizador) valorCalculado(nodo *Nodo, hermanos []*Nodo, valores Datos, enviado interface{}) string {
	if texto, esTexto := comoTexto(enviado); nodo.Dato.Editable && esTexto && strings.TrimSpace(texto) != "" {
		return transformar(nodo, texto)
	}
	return n.calcular(nodo, hermanos, valores)
}

// filas conserva el id (y el origen en los espejos) y limpia cada fila. Los
// elementos que no son objetos se conservan para que el validador los reporte.
func (n *Normalizador) filas(nodo *Nodo, valor interface{}) []interface{} {
	lista, ok := comoLista(valor)
	if !ok {
		return []interface{}{}
	}
	filas := make([]interface{}, 0, len(lista))
	for _, elemento := range lista {
		fila, esObjeto := comoObjeto(elemento)
		if !esObjeto {
			filas = append(filas, elemento)
			continue
		}
		limpia := n.objeto(nodo.Hijos, fila)
		limpia["id"] = limpiarEscalar(fila["id"])
		if nodo.EsEspejo() {
			limpia["origen"] = limpiarEscalar(fila["origen"])
		}
		filas = append(filas, limpia)
	}
	return filas
}

// matriz conserva solo las filas y columnas definidas.
func (n *Normalizador) matriz(nodo *Nodo, valor interface{}) Datos {
	entrada, _ := comoObjeto(valor)
	salida := Datos{}
	for _, fila := range nodo.Dato.Filas {
		celdas, _ := comoObjeto(entrada[fila.Clave])
		limpias := Datos{}
		for _, grupo := range nodo.Dato.GruposColumnas {
			for _, columna := range grupo.Columnas {
				if celda, existe := celdas[columna.Clave]; existe && celda != nil {
					limpias[columna.Clave] = celda
				}
			}
		}
		salida[fila.Clave] = limpias
	}
	return salida
}

var patronMarcador = regexp.MustCompile(`\{([a-z0-9_]+)\}`)

// calcular resuelve la plantilla del campo con los valores visibles de sus hermanos.
func (n *Normalizador) calcular(nodo *Nodo, hermanos []*Nodo, valores Datos) string {
	texto := patronMarcador.ReplaceAllStringFunc(nodo.Dato.Plantilla, func(marcador string) string {
		clave := marcador[1 : len(marcador)-1]
		return n.textoVisible(buscarNodo(hermanos, clave), valores[clave])
	})
	return transformar(nodo, texto)
}

// transformar quita espacios sobrantes y aplica la transformación del campo (p. ej. mayúsculas).
func transformar(nodo *Nodo, texto string) string {
	texto = strings.Join(strings.Fields(texto), " ")
	if nodo.Dato.Transformar == "mayusculas" {
		texto = strings.ToUpper(texto)
	}
	return texto
}

// textoVisible muestra el nombre del elemento de catálogo en lugar de su código.
func (n *Normalizador) textoVisible(nodo *Nodo, valor interface{}) string {
	codigo := textoPlano(valor)
	if nodo == nil || nodo.Dato.Tipo != TipoCatalogo || codigo == "" {
		return codigo
	}
	catalogo, err := n.catalogos.Catalogo(nodo.Dato.Catalogo)
	if err != nil {
		return codigo
	}
	if item, existe := catalogo.Buscar(codigo); existe {
		return item.Nombre
	}
	return codigo
}

func limpiarEscalar(valor interface{}) interface{} {
	if texto, ok := valor.(string); ok {
		return strings.TrimSpace(texto)
	}
	return valor
}
