package iniciativaInversionHelper

import (
	"testing"
)

func TestElEjemploDelExcelNoTieneErrores(t *testing.T) {
	for submodulo, errores := range validarIniciativa(t, definicionDePrueba(t), iniciativaDeEjemplo(t)) {
		exigirSinErrores(t, submodulo, errores)
	}
}

func TestLosObjetivosNoDependenDeLaProblematica(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	problematica := iniciativa["problematica"]
	problematica["causas_directas"] = append(filasEn(t, problematica, "causas_directas"), map[string]interface{}{
		"id": "cd-4", "descripcion": "Nueva causa directa.", "linea_base": "Sin medición.",
		"causas_indirectas": []interface{}{map[string]interface{}{"id": "ci-9", "descripcion": "Nueva causa indirecta."}},
	})

	exigirSinErrores(t, "objetivos", validarIniciativa(t, definicionDePrueba(t), iniciativa)["objetivos"])
}

func TestCadaObjetivoEspecificoNecesitaAlMenosUnMedio(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	filaConID(t, filasEn(t, iniciativa["objetivos"], "objetivos_especificos"), "oe-2")["medios"] = []interface{}{}

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["objetivos"], "objetivos_especificos[oe-2].medios", "mínimo 1")
}

func TestLasAlternativasExigenUnaFilaPorCadaObjetivoEspecifico(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	actividades := filasEn(t, iniciativa["alternativas"], "actividades_por_objetivo")
	iniciativa["alternativas"]["actividades_por_objetivo"] = actividades[:2] // se quita la de oe-3

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["alternativas"]
	exigirError(t, errores, "actividades_por_objetivo", "Fortalecer los componentes tecnológicos según")
}

func TestLasAlternativasRechazanFilasDeObjetivosQueYaNoExisten(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	filaConID(t, filasEn(t, iniciativa["alternativas"], "actividades_por_objetivo"), "ao-3")["origen"] = "oe-9"

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["alternativas"]
	exigirError(t, errores, "actividades_por_objetivo[ao-3]", "ya no existe")
	exigirError(t, errores, "actividades_por_objetivo", "Fortalecer los componentes tecnológicos según")
}

func TestElProgramaDebePertenecerAlSectorElegido(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	iniciativa["articulacion_pnd"]["programa"] = "2201" // programa del sector 22, no del 45

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["articulacion_pnd"]
	exigirError(t, errores, "programa", "no pertenece")
}

func TestCodigoInexistenteEnCatalogo(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	iniciativa["nombre"]["producto_principal"] = "9999999"

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["nombre"], "producto_principal", "no existe")
}

func TestParticipantesExigeMinimoDos(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	iniciativa["participantes"]["actores"] = filasEn(t, iniciativa["participantes"], "actores")[:1]

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["participantes"], "actores", "mínimo 2")
}

func TestFilasSinIdentificadorOConIdentificadorRepetido(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	actores := filasEn(t, iniciativa["participantes"], "actores")
	actores[1].(map[string]interface{})["id"] = "pa-1"
	delete(actores[2].(map[string]interface{}), "id")

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["participantes"]
	exigirError(t, errores, "actores", "repetido")
	exigirError(t, errores, "actores", "no tiene identificador")
}

func TestSelectRechazaValoresFueraDeLasOpciones(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	filaConID(t, filasEn(t, iniciativa["participantes"], "actores"), "pa-1")["posicion"] = "Aliado"

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["participantes"], "actores[pa-1].posicion", "no es una opción válida")
}

func TestMatrizIncompletaReportaLasCeldasFaltantes(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	iniciativa["poblacion"]["afectada"] = map[string]interface{}{}

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["poblacion"], "afectada", "faltan 170 de 170")
}

func TestMatrizRechazaNegativosYDecimales(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	infancia := iniciativa["poblacion"]["objetivo"].(map[string]interface{})["infancia"].(map[string]interface{})
	infancia["mujeres"] = -1.0
	infancia["hombres"] = 2.5

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["poblacion"]
	exigirError(t, errores, "objetivo.infancia.mujeres", "mayor o igual a 0")
	exigirError(t, errores, "objetivo.infancia.hombres", "entero")
}

func TestTextoObligatorioVacio(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	iniciativa["problematica"]["problema_central"] = "   "

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["problematica"], "problema_central", "es obligatorio")
}

// ---------------------------------------------------------------------------
// Nombre del proyecto: calculado, pero editable
// ---------------------------------------------------------------------------

const nombreSugerido = "IMPLEMENTACIÓN DE LA POLÍTICA DE ESPACIO PÚBLICO PARA LA GENERACIÓN DE MÁS Y MEJORES ÁREAS PARA ENCUENTRO, CUIDADO Y DISFRUTE EN BOGOTÁ D.C."

func normalizarNombre(t *testing.T, nombre interface{}) Datos {
	t.Helper()
	datos := iniciativaDeEjemplo(t)["nombre"]
	if nombre != nil {
		datos["nombre_proyecto"] = nombre
	}
	datos["clave_que_no_existe"] = "x"
	return NuevoNormalizador(catalogosDePrueba()).Normalizar(definicionDePrueba(t).Submodulo("nombre"), datos)
}

func TestSinNombreEscritoSeProponeElCalculado(t *testing.T) {
	for _, vacio := range []interface{}{nil, "", "   "} {
		normalizado := normalizarNombre(t, vacio)
		if normalizado["nombre_proyecto"] != nombreSugerido {
			t.Fatalf("con %q: nombre_proyecto = %q", vacio, normalizado["nombre_proyecto"])
		}
		if _, existe := normalizado["clave_que_no_existe"]; existe {
			t.Fatal("la clave desconocida debió descartarse")
		}
	}
}

func TestElUsuarioPuedeEscribirSuPropioNombre(t *testing.T) {
	normalizado := normalizarNombre(t, "  Modernización   de la infraestructura TIC  ")
	if normalizado["nombre_proyecto"] != "MODERNIZACIÓN DE LA INFRAESTRUCTURA TIC" {
		t.Fatalf("nombre_proyecto = %q", normalizado["nombre_proyecto"])
	}
}

func TestUnCalculadoNoEditableIgnoraLoQueEnvieElCliente(t *testing.T) {
	definicion := definicionDePrueba(t)
	submodulo := definicion.Submodulo("nombre")
	nodo := *submodulo.Hijo("nombre_proyecto")
	nodo.Dato.Editable = false
	copia := *submodulo
	copia.Hijos = append([]*Nodo{}, submodulo.Hijos...)
	for i, hijo := range copia.Hijos {
		if hijo.Clave() == "nombre_proyecto" {
			copia.Hijos[i] = &nodo
		}
	}
	datos := iniciativaDeEjemplo(t)["nombre"]
	datos["nombre_proyecto"] = "otro nombre"

	normalizado := NuevoNormalizador(catalogosDePrueba()).Normalizar(&copia, datos)
	if normalizado["nombre_proyecto"] != nombreSugerido {
		t.Fatalf("nombre_proyecto = %q", normalizado["nombre_proyecto"])
	}
}

func TestElNombreEsObligatorio(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	for _, clave := range []string{"proceso", "objeto", "localizacion"} {
		delete(iniciativa["nombre"], clave)
	}

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["nombre"], "nombre_proyecto", "es obligatorio")
}

func TestNormalizarConservaIdYOrigenDeLasFilas(t *testing.T) {
	definicion := definicionDePrueba(t)
	catalogos := catalogosDePrueba()
	ejemplo := iniciativaDeEjemplo(t)

	objetivos := NuevoNormalizador(catalogos).Normalizar(definicion.Submodulo("objetivos"), ejemplo["objetivos"])
	if len(filasEn(t, filaConID(t, filasEn(t, objetivos, "objetivos_especificos"), "oe-3"), "medios")) != 3 {
		t.Fatal("se perdieron los medios anidados")
	}
	alternativas := NuevoNormalizador(catalogos).Normalizar(definicion.Submodulo("alternativas"), ejemplo["alternativas"])
	if fila := filaConID(t, filasEn(t, alternativas, "actividades_por_objetivo"), "ao-2"); fila["origen"] != "oe-2" {
		t.Fatalf("origen = %v", fila["origen"])
	}
}

// ---------------------------------------------------------------------------
// Referencias y cobertura: el motor las soporta aunque la plantilla de Identificación hoy no
// las use. Se prueban con una definición pequeña.
// ---------------------------------------------------------------------------

func definicionConReferencias(cubrirTodas bool) *Definicion {
	campo := func(clave, tipo string, extra Dato, hijos ...*Nodo) *Nodo {
		extra.Clave, extra.Tipo, extra.Requerido = clave, tipo, true
		return &Nodo{Nombre: clave, Dato: extra, Hijos: hijos}
	}
	origen := &Nodo{Nombre: "Origen", Dato: Dato{Clave: "submodulo_origen"}, Hijos: []*Nodo{
		campo("padres", TipoRepetible, Dato{},
			campo("descripcion", TipoTexto, Dato{}),
			campo("hijos", TipoRepetible, Dato{}, campo("descripcion", TipoTexto, Dato{}))),
	}}
	destino := &Nodo{Nombre: "Destino", Dato: Dato{Clave: "submodulo_destino"}, Hijos: []*Nodo{
		campo("espejos", TipoRepetible, Dato{EspejoDe: "submodulo_origen.padres", Mostrar: "descripcion"},
			campo("acciones", TipoRepetible, Dato{EtiquetaFila: "Acción"},
				campo("hijo", TipoReferencia, Dato{OpcionesDesde: "origen.hijos", Mostrar: "descripcion", CubrirTodas: cubrirTodas}))),
	}}
	return &Definicion{Modulos: []*Nodo{{Dato: Dato{Clave: "m", Modulo: true}, Hijos: []*Nodo{origen, destino}}}}
}

func datosConReferencias(referencias ...string) map[string]Datos {
	acciones := []interface{}{}
	for i, referencia := range referencias {
		acciones = append(acciones, map[string]interface{}{"id": "a" + string(rune('1'+i)), "hijo": referencia})
	}
	return map[string]Datos{
		"submodulo_origen": {"padres": []interface{}{map[string]interface{}{
			"id": "p1", "descripcion": "Padre 1",
			"hijos": []interface{}{
				map[string]interface{}{"id": "h1", "descripcion": "Hijo uno"},
				map[string]interface{}{"id": "h2", "descripcion": "Hijo dos"},
			},
		}}},
		"submodulo_destino": {"espejos": []interface{}{map[string]interface{}{"id": "e1", "origen": "p1", "acciones": acciones}}},
	}
}

func TestUnaReferenciaSoloAceptaFilasDeSuFilaDeOrigen(t *testing.T) {
	definicion := definicionConReferencias(false)
	errores := NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(definicion.Submodulo("submodulo_destino"), datosConReferencias("h1", "h9"))

	exigirError(t, errores, "espejos[e1].acciones[a2].hijo", "no es válida")
}

func TestCubrirTodasExigeReferenciarCadaOpcion(t *testing.T) {
	definicion := definicionConReferencias(true)
	validador := NuevoValidador(catalogosDePrueba())

	exigirError(t, validador.ValidarSubmodulo(definicion.Submodulo("submodulo_destino"), datosConReferencias("h1")), "espejos[e1].acciones", "Hijo dos")
	exigirSinErrores(t, "submodulo_destino", validador.ValidarSubmodulo(definicion.Submodulo("submodulo_destino"), datosConReferencias("h1", "h2")))
}
