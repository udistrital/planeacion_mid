package iniciativaInversionHelper

import (
	"math"
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
	// Lo que no aplica (solo_si) no se guarda: p. ej. el texto de «Otra» si se desmarcó.
	for _, nodo := range nodos {
		if !nodo.visible(salida) {
			delete(salida, nodo.Clave())
		}
	}
	// Los calculados van al final porque usan los valores ya limpios de sus hermanos. Los que
	// traen un valor de otro submódulo (desde) solo se muestran: no se guardan.
	for _, nodo := range nodos {
		if nodo.Dato.Tipo != TipoCalculado || nodo.Clave() == "" || nodo.Dato.Desde != "" {
			continue
		}
		if nodo.Dato.Operacion == "" {
			salida[nodo.Clave()] = n.valorCalculado(nodo, nodos, salida, datos[nodo.Clave()])
		} else if resultado, ok := operar(nodo.Dato, salida); ok {
			salida[nodo.Clave()] = resultado
		}
	}
	return salida
}

// operar resuelve un calculado numérico con los valores de sus hermanos. Sin resultado (el
// campo queda sin valor) si la operación no existe, si en una resta o una multiplicación
// falta algún operando o si en una suma están todos vacíos.
func operar(dato Dato, valores Datos) (float64, bool) {
	if !operacionValida(dato) {
		return 0, false
	}
	resultado, hayAlguno := 0.0, false
	for _, clave := range dato.Operandos {
		numero, esNumero := comoNumero(valores[clave])
		if !esNumero {
			if dato.Operacion == OperacionSuma {
				continue // en la suma, vacío = 0
			}
			return 0, false
		}
		switch {
		case !hayAlguno:
			resultado = numero
		case dato.Operacion == OperacionResta:
			resultado -= numero
		case dato.Operacion == OperacionMultiplicacion:
			resultado *= numero
		default:
			resultado += numero
		}
		hayAlguno = true
	}
	if !hayAlguno {
		return 0, false
	}
	return redondear(resultado, dato.Decimales), true
}

// redondear quita el ruido del punto flotante (0,3 - 0,1 = 0,19999…): a los decimales
// del campo o, si no los define, a nueve.
func redondear(numero float64, decimales *int) float64 {
	cifras := 9
	if decimales != nil {
		cifras = *decimales
	}
	escala := math.Pow(10, float64(cifras))
	return math.Round(numero*escala) / escala
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
