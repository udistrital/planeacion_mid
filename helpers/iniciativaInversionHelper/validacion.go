package iniciativaInversionHelper

import (
	"fmt"
	"math"
	"regexp"
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
		TipoCatalogo:   validarCatalogo,
		TipoCalculado:  validarCalculado,
		TipoRepetible:  validarRepetible,
		TipoReferencia: validarReferencia,
		TipoMatriz:     validarMatriz,
		TipoCasilla:    validarCasilla,
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
		if clave == "" || !nodo.visible(valores) {
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
		campo, configurado := e.segunOrigen(nodo, actual, rutaNodo)
		if !configurado {
			continue
		}
		regla(e, campo, valor, valor != nil, rutaNodo, actual)
	}
}

// segunOrigen resuelve obligatorio_si "origen.<campo>": el campo se exige solo si ese campo de la
// fila de origen del espejo tiene valor (p. ej. la programación física de un año en que la actividad
// tiene costo en la Cadena de valor). Devuelve una copia del nodo con required resuelto, o false si
// está mal configurado.
func (e *ejecucion) segunOrigen(nodo *Nodo, alcance alcance, ruta string) (*Nodo, bool) {
	if nodo.Dato.ObligatorioSi == "" {
		return nodo, true
	}
	campo, deOrigen := strings.CutPrefix(nodo.Dato.ObligatorioSi, prefijoOrigen)
	if !deOrigen || campo == "" || alcance.filaOrigen == nil {
		e.malConfigurado(nodo, ruta, "obligatorio_si %q: se espera origen.<campo> en un campo de un grupo espejo", nodo.Dato.ObligatorioSi)
		return nil, false
	}
	copia := *nodo
	copia.Dato.Requerido = Booleano(conValor(alcance.filaOrigen[campo]))
	return &copia, true
}

// conValor indica si un valor cuenta como diligenciado para obligatorio_si: un número distinto de 0
// o cualquier otro contenido (un costo de 0 es como no tener costo).
func conValor(valor interface{}) bool {
	if numero, esNumero := comoNumero(valor); esNumero {
		return numero != 0
	}
	return tieneContenido(valor)
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

// validarCalculado: el valor ya lo resolvió el normalizador. El de texto (o el escrito, si es
// editable) se valida como texto; el numérico solo exige que la operación esté bien definida.
func validarCalculado(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, alcance alcance) {
	if nodo.Dato.Desde != "" {
		return // solo muestra un valor de otro submódulo; no se guarda
	}
	if nodo.Dato.Operacion == "" {
		validarTexto(e, nodo, valor, existe, ruta, alcance)
		return
	}
	if !operacionValida(nodo.Dato) {
		e.malConfigurado(nodo, ruta, "la operación %q con %d operandos no está soportada", nodo.Dato.Operacion, len(nodo.Dato.Operandos))
		return
	}
	if !existe {
		e.faltaObligatorio(nodo, ruta)
	}
}

// validarCasilla: marcada (true) o sin marcar (false o ausente). Si es obligatoria debe quedar
// marcada, salvo con unica: entonces lo obligatorio es marcarla en una fila de la lista (casillasUnicas).
func validarCasilla(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, _ alcance) {
	marcada, esBooleano := valor.(bool)
	switch {
	case existe && !esBooleano:
		e.agregar(ruta, "%s: debe quedar marcada o sin marcar.", nodo.Nombre)
	case !marcada && !nodo.Dato.Unica:
		e.faltaObligatorio(nodo, ruta)
	}
}

func validarCatalogo(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, alcance alcance) {
	if nodo.Dato.Multiple {
		e.catalogoMultiple(nodo, valor, existe, ruta)
		return
	}
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

// catalogoMultiple valida un catálogo con multiple: una lista de códigos sin repetir, todos del
// catálogo (p. ej. las sedes marcadas). Si es obligatorio, al menos uno.
func (e *ejecucion) catalogoMultiple(nodo *Nodo, valor interface{}, existe bool, ruta string) {
	lista, esLista := comoLista(valor)
	switch {
	case existe && !esLista:
		e.agregar(ruta, "%s: debe ser una lista de opciones.", nodo.Nombre)
		return
	case nodo.Dato.DependeDe != "":
		e.malConfigurado(nodo, ruta, "un catálogo multiple no admite depende_de")
		return
	case len(lista) == 0:
		e.faltaObligatorio(nodo, ruta)
		return
	}
	catalogo, err := e.validador.catalogos.Catalogo(nodo.Dato.Catalogo)
	if err != nil {
		logs.Warn("iniciativaInversion: no se pudo consultar el catálogo %s para %s: %v", nodo.Dato.Catalogo, ruta, err)
		e.agregar(ruta, "%s: no se pudieron consultar las opciones del campo. Intente de nuevo más tarde.", nodo.Nombre)
		return
	}
	vistos := map[string]bool{}
	for _, elemento := range lista {
		codigo, esTexto := comoTexto(elemento)
		switch {
		case !esTexto || codigo == "":
			e.agregar(ruta, "%s: hay una opción vacía o inválida.", nodo.Nombre)
		case vistos[codigo]:
			e.agregar(ruta, "%s: «%s» está marcada dos veces.", nodo.Nombre, codigo)
		default:
			if _, existeItem := catalogo.Buscar(codigo); !existeItem {
				e.agregar(ruta, "%s: «%s» ya no está en la lista; desmárquela.", nodo.Nombre, codigo)
			}
		}
		vistos[codigo] = true
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
	e.casillasUnicas(nodo, filas, ruta)
	if nodo.EsEspejo() {
		e.espejo(nodo, filas, ruta)
		return
	}
	e.cantidadDeFilas(nodo, len(lista), ruta)
	for _, fila := range filas {
		e.objeto(nodo.Hijos, fila, rutaFila(ruta, idDe(fila)), padre)
	}
}

// casillasUnicas exige que cada casilla con unica quede marcada en una sola fila de la lista
// (p. ej. la alternativa que pasa al módulo de Preparación); si es obligatoria, en exactamente una.
func (e *ejecucion) casillasUnicas(nodo *Nodo, filas []Datos, ruta string) {
	for _, hijo := range nodo.Hijos {
		if hijo.Dato.Tipo != TipoCasilla || !hijo.Dato.Unica {
			continue
		}
		marcadas := 0
		for _, fila := range filas {
			if marcada, _ := fila[hijo.Clave()].(bool); marcada {
				marcadas++
			}
		}
		switch {
		case marcadas > 1:
			e.agregar(ruta, "%s: la casilla «%s» solo se puede marcar en una fila; está marcada en %d.", nodo.Nombre, hijo.Nombre, marcadas)
		case marcadas == 0 && len(filas) > 0 && bool(hijo.Dato.Requerido):
			e.agregar(ruta, "%s: marque la casilla «%s» en una de las filas.", nodo.Nombre, hijo.Nombre)
		}
	}
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

// etiquetaDeFila es la etiqueta de las filas para los mensajes, en minúscula y sin el «No.» final
// del título (p. ej. «Bien o servicio a entregar o demanda a satisfacer No.»).
func etiquetaDeFila(nodo *Nodo) string {
	return strings.ToLower(patronNumeroFinal.ReplaceAllString(nodo.Dato.EtiquetaFila, ""))
}

var patronNumeroFinal = regexp.MustCompile(`(?i)\s+no\.?$`)

func (e *ejecucion) cantidadDeFilas(nodo *Nodo, cantidad int, ruta string) {
	minimo := 0
	if nodo.Dato.Min != nil {
		minimo = int(*nodo.Dato.Min)
	} else if nodo.Dato.Requerido {
		minimo = 1
	}
	etiqueta := etiquetaDeFila(nodo)
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

// espejo exige exactamente una fila por cada fila del grupo de origen (o por cada elemento del
// catálogo). Con filas_opcionales, la fila vacía no se exige ni se valida (no aplica), salvo las
// de filas_obligatorias, y min es cuántas filas se llenan como mínimo. Con filas_adicionales, acepta
// además las filas agregadas por el usuario (hasta max, si lo tiene), que se validan igual que las demás.
func (e *ejecucion) espejo(nodo *Nodo, filas []Datos, ruta string) {
	origen, ok := e.origenDeEspejo(nodo, ruta)
	if !ok {
		return
	}
	idsOrigen := map[string]bool{}
	for _, filaOrigen := range origen {
		idsOrigen[idDe(filaOrigen)] = true
	}
	adicionales := nodo.Dato.FilasAdicionales
	if adicionales != nil && (adicionales.Origen == "" || adicionales.Max < 0 || idsOrigen[adicionales.Origen]) {
		e.malConfigurado(nodo, ruta, "filas_adicionales: origen debe ser un valor que no esté en el origen del espejo y max no puede ser negativo")
		return
	}

	porOrigen := map[string]Datos{}
	agregadas := []Datos{}
	for _, fila := range filas {
		idOrigen, _ := comoTexto(fila["origen"])
		rutaDeFila := rutaFila(ruta, idDe(fila))
		switch {
		case idOrigen == "":
			e.agregar(rutaDeFila, "%s: la fila no indica de qué elemento sale.", nodo.Nombre)
		case adicionales != nil && idOrigen == adicionales.Origen:
			agregadas = append(agregadas, fila)
		case !idsOrigen[idOrigen]:
			e.agregar(rutaDeFila, "%s: la fila corresponde a un elemento que ya no existe; elimínela.", nodo.Nombre)
		case porOrigen[idOrigen] != nil:
			e.agregar(rutaDeFila, "%s: hay más de una fila para el mismo elemento.", nodo.Nombre)
		default:
			porOrigen[idOrigen] = fila
		}
	}

	diligenciadas := 0
	for _, filaOrigen := range origen {
		fila := porOrigen[idDe(filaOrigen)]
		if nodo.Dato.FilasOpcionales && !filaDiligenciada(nodo.Hijos, fila) && !filaObligatoria(nodo, idDe(filaOrigen)) {
			continue
		}
		if fila == nil {
			e.agregar(ruta, "%s: falta diligenciar el elemento que corresponde a «%s».",
				nodo.Nombre, resumir(textoPlano(filaOrigen[nodo.Dato.Mostrar]), 80))
			continue
		}
		diligenciadas++
		e.objeto(nodo.Hijos, fila, rutaFila(ruta, idDe(fila)), alcance{filaOrigen: filaOrigen})
	}
	if adicionales != nil {
		diligenciadas += e.filasAgregadas(nodo, agregadas, ruta)
	}
	if nodo.Dato.FilasOpcionales && nodo.Dato.Min != nil && diligenciadas < int(*nodo.Dato.Min) {
		e.agregar(ruta, "%s: diligencie al menos %d (hay %d).", nodo.Nombre, int(*nodo.Dato.Min), diligenciadas)
	}
}

// filasAgregadas valida las filas que el usuario agregó a un espejo con filas_adicionales (p. ej. los
// otros trazadores presupuestales): como mucho max (si lo tiene) y, con filas_opcionales, la vacía no
// se valida. No tienen fila de origen: la de su alcance solo lleva el id (Origen), así que "origen.<campo>"
// no trae nada. Devuelve cuántas se validaron.
func (e *ejecucion) filasAgregadas(nodo *Nodo, filas []Datos, ruta string) int {
	adicionales := nodo.Dato.FilasAdicionales
	if adicionales.Max > 0 && len(filas) > adicionales.Max {
		etiqueta := strings.ToLower(adicionales.Etiqueta)
		if etiqueta == "" {
			etiqueta = "filas agregadas"
		}
		e.agregar(ruta, "%s: se permiten máximo %d (%s), hay %d.", nodo.Nombre, adicionales.Max, etiqueta, len(filas))
	}
	validadas := 0
	for _, fila := range filas {
		if nodo.Dato.FilasOpcionales && !filaDiligenciada(nodo.Hijos, fila) {
			continue
		}
		validadas++
		e.objeto(nodo.Hijos, fila, rutaFila(ruta, idDe(fila)), alcance{filaOrigen: Datos{"id": adicionales.Origen}})
	}
	return validadas
}

// filaObligatoria indica si la fila de ese elemento de origen se exige aunque el espejo tenga
// filas_opcionales (p. ej. los resultados de estudiantes, docentes y administrativos beneficiados).
func filaObligatoria(nodo *Nodo, idOrigen string) bool {
	for _, id := range nodo.Dato.FilasObligatorias {
		if id == idOrigen {
			return true
		}
	}
	return false
}

// origenDeEspejo devuelve las filas de origen de un espejo: las de otro grupo o, con
// "catalogo.<nombre>", los elementos del catálogo. Si no se pueden obtener, reporta el error.
func (e *ejecucion) origenDeEspejo(nodo *Nodo, ruta string) ([]Datos, bool) {
	if nombre, esCatalogo := catalogoDeRuta(nodo.Dato.EspejoDe); esCatalogo {
		catalogo, err := e.validador.catalogos.Catalogo(nombre)
		if err != nil {
			logs.Warn("iniciativaInversion: no se pudo consultar el catálogo %s para %s: %v", nombre, ruta, err)
			e.agregar(ruta, "%s: no se pudieron consultar los elementos de la lista. Intente de nuevo más tarde.", nodo.Nombre)
			return nil, false
		}
		return filasDeCatalogo(catalogo), true
	}
	origen, err := e.filasDeRuta(nodo.Dato.EspejoDe)
	if err != "" {
		e.malConfigurado(nodo, ruta, "espejo_de: %s", err)
		return nil, false
	}
	return origen, true
}

// filaDiligenciada indica si la fila tiene algún valor en sus campos.
func filaDiligenciada(nodos []*Nodo, fila Datos) bool {
	return fila != nil && tieneInformacion(nodos, fila)
}

// tieneInformacion indica si hay algo diligenciado en esos campos: un texto, número, opción,
// casilla marcada, celda o documento, también dentro de grupos y filas. No cuentan los calculados
// (salen de los demás) ni el id y el origen de las filas, que pone el sistema (p. ej. las filas de
// un espejo que nadie ha llenado).
func tieneInformacion(nodos []*Nodo, datos Datos) bool {
	for _, nodo := range nodos {
		clave := nodo.Clave()
		if clave == "" || nodo.Dato.Tipo == TipoCalculado {
			continue
		}
		valor := datos[clave]
		switch {
		case !nodo.EsCampo():
			if subdatos, esObjeto := comoObjeto(valor); esObjeto && tieneInformacion(nodo.Hijos, subdatos) {
				return true
			}
		case nodo.Dato.Tipo == TipoRepetible:
			for _, fila := range filasDe(valor) {
				if tieneInformacion(nodo.Hijos, fila) {
					return true
				}
			}
		case tieneContenido(valor):
			return true
		}
	}
	return false
}

func tieneContenido(valor interface{}) bool {
	switch v := valor.(type) {
	case nil:
		return false
	case string:
		return strings.TrimSpace(v) != ""
	case bool:
		return v
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		for _, interno := range v {
			if tieneContenido(interno) {
				return true
			}
		}
		return false
	}
	return true
}

func validarReferencia(e *ejecucion, nodo *Nodo, valor interface{}, existe bool, ruta string, alcance alcance) {
	ids := idsReferenciados(valor)
	if nodo.Dato.MarcadaEn != "" {
		e.referenciaMarcada(nodo, ids, ruta, alcance)
		return
	}
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

// referenciaMarcada valida una referencia con marcada_en: su única opción es la fila marcada en
// esa casilla del grupo de origen (p. ej. la alternativa que pasa al módulo de Preparación). El
// cliente la llena sola; si después se marca otra fila, el submódulo se debe revisar y guardar.
func (e *ejecucion) referenciaMarcada(nodo *Nodo, ids []string, ruta string, alcance alcance) {
	opciones, err := e.opcionesDe(nodo, alcance)
	if err != "" {
		e.malConfigurado(nodo, ruta, "opciones_desde: %s", err)
		return
	}
	switch {
	case len(opciones) == 0:
		if nodo.Dato.Requerido {
			e.agregar(ruta, "%s: aún no se ha marcado cuál pasa; márquela primero en el submódulo de donde se trae.", nodo.Nombre)
		}
	case len(ids) == 0:
		e.faltaObligatorio(nodo, ruta)
	case len(ids) > 1 || opciones[ids[0]] == nil:
		e.agregar(ruta, "%s: se marcó otra en el submódulo de donde se trae; revise este submódulo y guárdelo de nuevo.", nodo.Nombre)
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
// origen del espejo) o "<submodulo>.<grupo>". Con marcada_en, solo las filas que
// tienen marcada esa casilla.
func (e *ejecucion) listaDeOpciones(nodo *Nodo, alcance alcance) ([]Datos, string) {
	desde := nodo.Dato.OpcionesDesde
	var filas []Datos
	if strings.HasPrefix(desde, prefijoOrigen) {
		if alcance.filaOrigen == nil {
			return nil, "la referencia " + desde + " solo se puede usar dentro de un grupo espejo."
		}
		filas = filasDe(alcance.filaOrigen[strings.TrimPrefix(desde, prefijoOrigen)])
	} else {
		var err string
		if filas, err = e.filasDeRuta(desde); err != "" {
			return nil, err
		}
	}
	if nodo.Dato.MarcadaEn == "" {
		return filas, ""
	}
	marcadas := []Datos{}
	for _, fila := range filas {
		if marcada, _ := fila[nodo.Dato.MarcadaEn].(bool); marcada {
			marcadas = append(marcadas, fila)
		}
	}
	return marcadas, ""
}

// filasDeRuta devuelve las filas de "<submodulo>.<grupo>" en los datos de la iniciativa.
// La ruta puede seguir por grupos anidados ("<submodulo>.<grupo>.<subgrupo>"): entonces
// devuelve, en orden, las filas de <subgrupo> de todas las filas de <grupo>.
func (e *ejecucion) filasDeRuta(ruta string) ([]Datos, string) {
	partes := strings.Split(ruta, ".")
	if len(partes) < 2 || strings.Contains(ruta, "..") || strings.HasSuffix(ruta, ".") {
		return nil, fmt.Sprintf("la ruta %q no es válida (se espera submodulo.grupo).", ruta)
	}
	filas := filasDe(e.iniciativa[partes[0]][partes[1]])
	for _, clave := range partes[2:] {
		anidadas := []Datos{}
		for _, fila := range filas {
			anidadas = append(anidadas, filasDe(fila[clave])...)
		}
		filas = anidadas
	}
	return filas, ""
}
