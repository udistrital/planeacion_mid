package iniciativaInversionHelper

import (
	"strconv"
	"strings"
)

// Utilidades para leer valores JSON decodificados (interface{}) sin pánicos.

func comoObjeto(valor interface{}) (Datos, bool) {
	objeto, ok := valor.(map[string]interface{})
	return objeto, ok
}

func comoLista(valor interface{}) ([]interface{}, bool) {
	if valor == nil {
		return nil, true
	}
	lista, ok := valor.([]interface{})
	return lista, ok
}

func comoTexto(valor interface{}) (string, bool) {
	texto, ok := valor.(string)
	return texto, ok
}

func comoNumero(valor interface{}) (float64, bool) {
	numero, ok := valor.(float64)
	return numero, ok
}

// filasDe devuelve solo los elementos que son objetos (filas válidas).
func filasDe(valor interface{}) []Datos {
	lista, _ := comoLista(valor)
	filas := make([]Datos, 0, len(lista))
	for _, elemento := range lista {
		if fila, ok := comoObjeto(elemento); ok {
			filas = append(filas, fila)
		}
	}
	return filas
}

// idDe devuelve el identificador de una fila.
func idDe(fila Datos) string {
	id, _ := comoTexto(fila["id"])
	return id
}

// textoPlano convierte un valor simple en texto legible.
func textoPlano(valor interface{}) string {
	switch v := valor.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

// resumir acorta un texto para usarlo dentro de un mensaje.
func resumir(texto string, maximo int) string {
	texto = strings.Join(strings.Fields(texto), " ")
	runas := []rune(texto)
	if len(runas) <= maximo {
		return texto
	}
	return string(runas[:maximo]) + "…"
}

// rutaCampo arma la ruta de un campo: "a.b".
func rutaCampo(base, clave string) string {
	if base == "" {
		return clave
	}
	return base + "." + clave
}

// rutaFila arma la ruta de una fila: "causas_directas[cd-1]".
func rutaFila(rutaGrupo, id string) string {
	return rutaGrupo + "[" + id + "]"
}
