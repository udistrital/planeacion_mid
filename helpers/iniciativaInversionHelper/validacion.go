package iniciativaInversionHelper

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/astaxie/beego/logs"
)

// Validador revisa los datos de un submódulo contra su definición. Cada tipo de
// campo tiene su regla en reglas; para soportar un tipo nuevo basta registrarla.
type Validador struct {
	catalogos ProveedorCatalogos
	reglas    map[string]reglaCampo
}

// reglaCampo valida un campo. existe indica si el valor vino y no es null.
type reglaCampo func(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, alcance alcance)

// alcance es el contexto de un nivel del formulario.
type alcance struct {
	hermanos   Datos   // valores del mismo nivel (para depende_de)
	campos     []*Nodo // campos del mismo nivel (para nombrarlos en los mensajes)
	filaOrigen Datos   // fila de origen del espejo más cercano (para "origen.*")
}

// campoHermano nombra para el usuario el campo del mismo nivel con esa clave
// (p. ej. «Sector»); los mensajes no muestran claves internas.
func (a alcance) campoHermano(clave string) string {
	for _, nodo := range a.campos {
		if nodo.Clave() == clave {
			return "«" + nodo.Nombre + "»"
		}
	}
	return "el campo del que depende"
}

// ejecucion acumula los errores de una validación.
type ejecucion struct {
	validador  *Validador
	iniciativa map[string]Datos // datos de todos los submódulos, por clave de submódulo
	errores    []ErrorCampo
}

// NuevoValidador crea el validador con las reglas de todos los tipos de campo.
func NuevoValidador(catalogos ProveedorCatalogos) *Validador {
	v := &Validador{catalogos: catalogos}
	v.reglas = map[string]reglaCampo{
		TipoTexto:      validarTexto,
		TipoNumerico:   validarNumerico,
		TipoSelect:     validarSelect,
		TipoCatalogo:   validarCatalogo,
		TipoCalculado:  validarTexto, // ya calculado (o escrito, si es editable): se valida como texto
		TipoRepetible:  validarRepetible,
		TipoReferencia: validarReferencia,
		TipoMatriz:     validarMatriz,
	}
	return v
}

// ValidarSubmodulo valida iniciativa[clave del submódulo]. iniciativa trae los datos
// de todos los submódulos porque los espejos y las referencias cruzan submódulos.
func (v *Validador) ValidarSubmodulo(submodulo *Nodo, iniciativa map[string]Datos) []ErrorCampo {
	e := &ejecucion{validador: v, iniciativa: iniciativa}
	e.objeto(submodulo.Hijos, iniciativa[submodulo.Clave()], "", alcance{})
	if e.errores == nil {
		return []ErrorCampo{}
	}
	return e.errores
}

func (e *ejecucion) agregar(ruta, formato string, args ...interface{}) {
	e.errores = append(e.errores, ErrorCampo{Ruta: ruta, Mensaje: fmt.Sprintf(formato, args...)})
}

// malConfigurado reporta un campo que la plantilla define mal. El usuario no puede
// corregirlo: ve un aviso general y el detalle va al log.
func (e *ejecucion) malConfigurado(nodo *Nodo, ruta, formato string, args ...interface{}) {
	logs.Warn("iniciativaInversion: el campo %s (%s) está mal configurado en la plantilla: %s", ruta, nodo.Nombre, fmt.Sprintf(formato, args...))
	e.agregar(ruta, "%s: el campo está mal configurado en la plantilla. Comuníquese con el administrador del sistema.", nodo.Nombre)
}

func (e *ejecucion) faltaObligatorio(nodo *Nodo, ruta string) {
	if nodo.Dato.Requerido {
		e.agregar(ruta, "%s: es obligatorio.", nodo.Nombre)
	}
}

// objeto valida un nivel: cada hijo según su tipo; los grupos sin tipo se recorren.
func (e *ejecucion) objeto(nodos []*Nodo, valores Datos, ruta string, padre alcance) {
	actual := alcance{hermanos: valores, campos: nodos, filaOrigen: padre.filaOrigen}
	for _, nodo := range nodos {
		clave := nodo.Clave()
		if clave == "" {
			continue
		}
		valor := valores[clave]
		rutaNodo := rutaCampo(ruta, clave)
		if !nodo.EsCampo() {
			subdatos, _ := comoObjeto(valor)
			e.objeto(nodo.Hijos, subdatos, rutaNodo, actual)
			continue
		}
		regla, soportado := e.validador.reglas[nodo.Dato.Tipo]
		if !soportado {
			e.malConfigurado(nodo, rutaNodo, "el tipo de campo %q no está soportado", nodo.Dato.Tipo)
			continue
		}
		regla(e, nodo, valor, valor != nil, rutaNodo, actual)
	}
}

// ---------------------------------------------------------------------------
// Campos simples
// ---------------------------------------------------------------------------

func validarTexto(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, _ alcance) {
	texto, esTexto := comoTexto(valor)
	switch {
	case !existe:
		e.faltaObligatorio(nodo, ruta)
	case !esTexto:
		e.agregar(ruta, "%s: debe ser texto.", nodo.Nombre)
	case strings.TrimSpace(texto) == "":
		e.faltaObligatorio(nodo, ruta)
	case nodo.Dato.MaxLongitud != nil && utf8.RuneCountInString(texto) > *nodo.Dato.MaxLongitud:
		e.agregar(ruta, "%s: admite máximo %d caracteres.", nodo.Nombre, *nodo.Dato.MaxLongitud)
	}
}

func validarNumerico(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, _ alcance) {
	if !existe {
		e.faltaObligatorio(nodo, ruta)
		return
	}
	numero, esNumero := comoNumero(valor)
	if !esNumero {
		e.agregar(ruta, "%s: debe ser un número.", nodo.Nombre)
		return
	}
	if mensaje := revisarNumero(numero, nodo.Dato); mensaje != "" {
		e.agregar(ruta, "%s: %s", nodo.Nombre, mensaje)
	}
}

// revisarNumero aplica min, max y decimales; devuelve "" si el número es válido.
func revisarNumero(numero float64, dato Dato) string {
	if dato.Min != nil && numero < *dato.Min {
		return fmt.Sprintf("debe ser mayor o igual a %v.", *dato.Min)
	}
	if dato.Max != nil && numero > *dato.Max {
		return fmt.Sprintf("debe ser menor o igual a %v.", *dato.Max)
	}
	if dato.Decimales != nil {
		escala := math.Pow(10, float64(*dato.Decimales))
		if math.Abs(numero*escala-math.Round(numero*escala)) > 1e-9 {
			if *dato.Decimales == 0 {
				return "debe ser un número entero."
			}
			return fmt.Sprintf("admite máximo %d decimales.", *dato.Decimales)
		}
	}
	return ""
}

func validarSelect(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, _ alcance) {
	texto, esTexto := comoTexto(valor)
	if !existe || (esTexto && texto == "") {
		e.faltaObligatorio(nodo, ruta)
		return
	}
	for _, opcion := range nodo.Dato.Opciones {
		if esTexto && opcion.Valor == texto {
			return
		}
	}
	e.agregar(ruta, "%s: el valor seleccionado no es una opción válida.", nodo.Nombre)
}

func validarCatalogo(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, alcance alcance) {
	codigo, esTexto := comoTexto(valor)
	if !existe || (esTexto && codigo == "") {
		e.faltaObligatorio(nodo, ruta)
		return
	}
	if !esTexto {
		e.agregar(ruta, "%s: seleccione una opción de la lista.", nodo.Nombre)
		return
	}
	catalogo, err := e.validador.catalogos.Catalogo(nodo.Dato.Catalogo)
	if err != nil {
		logs.Warn("iniciativaInversion: no se pudo consultar el catálogo %s para %s: %v", nodo.Dato.Catalogo, ruta, err)
		e.agregar(ruta, "%s: no se pudieron consultar las opciones del campo. Intente de nuevo más tarde.", nodo.Nombre)
		return
	}
	item, existeItem := catalogo.Buscar(codigo)
	if !existeItem {
		e.agregar(ruta, "%s: la opción seleccionada no existe en la lista; seleccione otra.", nodo.Nombre)
		return
	}
	if nodo.Dato.DependeDe == "" {
		return
	}
	padre, _ := comoTexto(alcance.hermanos[nodo.Dato.DependeDe])
	if padre == "" {
		e.agregar(ruta, "%s: primero seleccione %s.", nodo.Nombre, alcance.campoHermano(nodo.Dato.DependeDe))
	} else if item.Padre != padre {
		e.agregar(ruta, "%s: «%s» no pertenece a la opción seleccionada en %s.", nodo.Nombre, item.Nombre, alcance.campoHermano(nodo.Dato.DependeDe))
	}
}

func validarMatriz(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, _ alcance) {
	celdasPorFila, esObjeto := comoObjeto(valor)
	if existe && !esObjeto {
		e.agregar(ruta, "%s: formato de matriz inválido.", nodo.Nombre)
		return
	}
	faltantes, total := 0, 0
	for _, fila := range nodo.Dato.Filas {
		celdas, _ := comoObjeto(celdasPorFila[fila.Clave])
		for _, grupo := range nodo.Dato.GruposColumnas {
			for _, columna := range grupo.Columnas {
				total++
				celda, hay := celdas[columna.Clave]
				if !hay || celda == nil {
					faltantes++
					continue
				}
				rutaCelda := rutaCampo(rutaCampo(ruta, fila.Clave), columna.Clave)
				numero, esNumero := comoNumero(celda)
				if !esNumero {
					e.agregar(rutaCelda, "%s (%s, %s): debe ser un número.", nodo.Nombre, fila.Nombre, columna.Nombre)
				} else if mensaje := revisarNumero(numero, nodo.Dato); mensaje != "" {
					e.agregar(rutaCelda, "%s (%s, %s): %s", nodo.Nombre, fila.Nombre, columna.Nombre, mensaje)
				}
			}
		}
	}
	if faltantes > 0 && nodo.Dato.Requerido {
		e.agregar(ruta, "%s: faltan %d de %d celdas por diligenciar (use 0 si no aplica).", nodo.Nombre, faltantes, total)
	}
}

// ---------------------------------------------------------------------------
// Grupos repetibles, espejos y referencias
// ---------------------------------------------------------------------------

func validarRepetible(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, padre alcance) {
	lista, esLista := comoLista(valor)
	if existe && !esLista {
		e.agregar(ruta, "%s: debe ser una lista.", nodo.Nombre)
		return
	}
	filas := e.filasValidas(nodo, lista, ruta)
	if nodo.EsEspejo() {
		e.espejo(nodo, filas, ruta)
		return
	}
	e.cantidadDeFilas(nodo, len(lista), ruta)
	for _, fila := range filas {
		e.objeto(nodo.Hijos, fila, rutaFila(ruta, idDe(fila)), padre)
	}
	e.cobertura(nodo, filas, ruta, padre)
}

// filasValidas descarta (y reporta) los elementos que no son objetos o que no
// tienen un id único.
func (e *ejecucion) filasValidas(nodo *Nodo, lista []interface{}, ruta string) []Datos {
	vistos := map[string]bool{}
	filas := make([]Datos, 0, len(lista))
	for posicion, elemento := range lista {
		fila, esObjeto := comoObjeto(elemento)
		id := idDe(fila)
		switch {
		case !esObjeto:
			e.agregar(ruta, "%s: el elemento %d no tiene el formato de una fila.", nodo.Nombre, posicion+1)
		case id == "":
			e.agregar(ruta, "%s: la fila %d no tiene identificador.", nodo.Nombre, posicion+1)
		case vistos[id]:
			e.agregar(ruta, "%s: el identificador de fila %q está repetido.", nodo.Nombre, id)
		default:
			vistos[id] = true
			filas = append(filas, fila)
		}
	}
	return filas
}

func (e *ejecucion) cantidadDeFilas(nodo *Nodo, cantidad int, ruta string) {
	minimo := 0
	if nodo.Dato.Min != nil {
		minimo = int(*nodo.Dato.Min)
	} else if nodo.Dato.Requerido {
		minimo = 1
	}
	etiqueta := strings.ToLower(nodo.Dato.EtiquetaFila)
	if etiqueta == "" {
		etiqueta = "fila"
	}
	if cantidad < minimo {
		e.agregar(ruta, "%s: se requiere mínimo %d (%s), hay %d.", nodo.Nombre, minimo, etiqueta, cantidad)
	}
	if nodo.Dato.Max != nil && cantidad > int(*nodo.Dato.Max) {
		e.agregar(ruta, "%s: se permiten máximo %d (%s), hay %d.", nodo.Nombre, int(*nodo.Dato.Max), etiqueta, cantidad)
	}
}

// espejo exige exactamente una fila por cada fila del grupo de origen.
func (e *ejecucion) espejo(nodo *Nodo, filas []Datos, ruta string) {
	origen, err := e.filasDeRuta(nodo.Dato.EspejoDe)
	if err != "" {
		e.malConfigurado(nodo, ruta, "espejo_de: %s", err)
		return
	}
	idsOrigen := map[string]bool{}
	for _, filaOrigen := range origen {
		idsOrigen[idDe(filaOrigen)] = true
	}

	porOrigen := map[string]Datos{}
	for _, fila := range filas {
		idOrigen, _ := comoTexto(fila["origen"])
		rutaDeFila := rutaFila(ruta, idDe(fila))
		switch {
		case idOrigen == "":
			e.agregar(rutaDeFila, "%s: la fila no indica de qué elemento sale.", nodo.Nombre)
		case !idsOrigen[idOrigen]:
			e.agregar(rutaDeFila, "%s: la fila corresponde a un elemento que ya no existe; elimínela.", nodo.Nombre)
		case porOrigen[idOrigen] != nil:
			e.agregar(rutaDeFila, "%s: hay más de una fila para el mismo elemento.", nodo.Nombre)
		default:
			porOrigen[idOrigen] = fila
		}
	}

	for _, filaOrigen := range origen {
		fila := porOrigen[idDe(filaOrigen)]
		if fila == nil {
			e.agregar(ruta, "%s: falta diligenciar el elemento que corresponde a «%s».",
				nodo.Nombre, resumir(textoPlano(filaOrigen[nodo.Dato.Mostrar]), 80))
			continue
		}
		e.objeto(nodo.Hijos, fila, rutaFila(ruta, idDe(fila)), alcance{filaOrigen: filaOrigen})
	}
}

func validarReferencia(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, alcance alcance) {
	ids := idsReferenciados(valor)
	if !existe || len(ids) == 0 {
		e.faltaObligatorio(nodo, ruta)
		return
	}
	if !nodo.Dato.Multiple && len(ids) > 1 {
		e.agregar(ruta, "%s: solo se permite seleccionar una opción.", nodo.Nombre)
		return
	}
	opciones, err := e.opcionesDe(nodo, alcance)
	if err != "" {
		e.malConfigurado(nodo, ruta, "opciones_desde: %s", err)
		return
	}
	for _, id := range ids {
		if _, valida := opciones[id]; !valida {
			e.agregar(ruta, "%s: la opción seleccionada no es válida aquí (puede haber sido eliminada).", nodo.Nombre)
			return
		}
	}
}

// cobertura exige que cada opción de una referencia con cubrir_todas quede
// seleccionada en al menos una fila (p. ej. toda causa indirecta con su medio).
func (e *ejecucion) cobertura(nodo *Nodo, filas []Datos, ruta string, padre alcance) {
	for _, hijo := range nodo.Hijos {
		if hijo.Dato.Tipo != TipoReferencia || !hijo.Dato.CubrirTodas {
			continue
		}
		opciones, err := e.listaDeOpciones(hijo, padre)
		if err != "" {
			continue // el error ya se reporta en cada fila
		}
		usadas := map[string]bool{}
		for _, fila := range filas {
			for _, id := range idsReferenciados(fila[hijo.Clave()]) {
				usadas[id] = true
			}
		}
		etiqueta := strings.ToLower(nodo.Dato.EtiquetaFila)
		if etiqueta == "" {
			etiqueta = "elemento"
		}
		for _, opcion := range opciones {
			if !usadas[idDe(opcion)] {
				e.agregar(ruta, "%s: ningún %s atiende «%s».",
					nodo.Nombre, etiqueta, resumir(textoPlano(opcion[hijo.Dato.Mostrar]), 80))
			}
		}
	}
}

func idsReferenciados(valor interface{}) []string {
	if texto, ok := comoTexto(valor); ok {
		if texto == "" {
			return nil
		}
		return []string{texto}
	}
	lista, _ := comoLista(valor)
	ids := make([]string, 0, len(lista))
	for _, elemento := range lista {
		if texto, ok := comoTexto(elemento); ok && texto != "" {
			ids = append(ids, texto)
		}
	}
	return ids
}

// opcionesDe indexa por id las filas que puede referenciar el campo.
func (e *ejecucion) opcionesDe(nodo *Nodo, alcance alcance) (map[string]Datos, string) {
	lista, err := e.listaDeOpciones(nodo, alcance)
	if err != "" {
		return nil, err
	}
	opciones := make(map[string]Datos, len(lista))
	for _, fila := range lista {
		opciones[idDe(fila)] = fila
	}
	return opciones, ""
}

// listaDeOpciones resuelve opciones_desde: "origen.<grupo>" (dentro de la fila de
// origen del espejo) o "<submodulo>.<grupo>".
func (e *ejecucion) listaDeOpciones(nodo *Nodo, alcance alcance) ([]Datos, string) {
	desde := nodo.Dato.OpcionesDesde
	if strings.HasPrefix(desde, prefijoOrigen) {
		if alcance.filaOrigen == nil {
			return nil, "la referencia " + desde + " solo se puede usar dentro de un grupo espejo."
		}
		return filasDe(alcance.filaOrigen[strings.TrimPrefix(desde, prefijoOrigen)]), ""
	}
	return e.filasDeRuta(desde)
}

// filasDeRuta devuelve las filas de "<submodulo>.<grupo>" en los datos de la iniciativa.
func (e *ejecucion) filasDeRuta(ruta string) ([]Datos, string) {
	partes := strings.Split(ruta, ".")
	if len(partes) != 2 {
		return nil, fmt.Sprintf("la ruta %q no es válida (se espera submodulo.grupo).", ruta)
	}
	return filasDe(e.iniciativa[partes[0]][partes[1]]), ""
}
