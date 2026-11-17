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

func TestCadaObjetivoDeLasAlternativasExigeSuDescripcion(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	delete(filaConID(t, filasEn(t, iniciativa["alternativas"], "actividades_por_objetivo"), "ao-2"), "descripcion_alternativa")

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["alternativas"]
	exigirError(t, errores, "actividades_por_objetivo[ao-2].descripcion_alternativa", "es obligatorio")
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

// ---------------------------------------------------------------------------
// Casillas y calculados numéricos (Alternativas de solución y Necesidades)
// ---------------------------------------------------------------------------

func campoDePrueba(clave, tipo string, extra Dato, hijos ...*Nodo) *Nodo {
	extra.Clave, extra.Tipo = clave, tipo
	return &Nodo{Nombre: clave, Dato: extra, Hijos: hijos}
}

func submoduloAlternativas() *Nodo {
	return &Nodo{Nombre: "Alternativas", Dato: Dato{Clave: "alternativas"}, Hijos: []*Nodo{
		campoDePrueba("lista", TipoRepetible, Dato{Requerido: true},
			campoDePrueba("descripcion", TipoTexto, Dato{Requerido: true}),
			campoDePrueba("pasa", TipoCasilla, Dato{Requerido: true, Unica: true})),
	}}
}

func alternativas(marcas ...interface{}) map[string]Datos {
	filas := []interface{}{}
	for i, marca := range marcas {
		fila := map[string]interface{}{"id": "a" + string(rune('1'+i)), "descripcion": "Alternativa"}
		if marca != nil {
			fila["pasa"] = marca
		}
		filas = append(filas, fila)
	}
	return map[string]Datos{"alternativas": {"lista": filas}}
}

func TestUnaCasillaUnicaSeMarcaEnUnaSolaFila(t *testing.T) {
	validador := NuevoValidador(catalogosDePrueba())
	submodulo := submoduloAlternativas()

	exigirError(t, validador.ValidarSubmodulo(submodulo, alternativas(nil, false)), "lista", "marque la casilla «pasa»")
	exigirError(t, validador.ValidarSubmodulo(submodulo, alternativas(true, false, true)), "lista", "está marcada en 2")
	exigirError(t, validador.ValidarSubmodulo(submodulo, alternativas("Sí", true)), "lista[a1].pasa", "marcada o sin marcar")
	exigirSinErrores(t, "alternativas", validador.ValidarSubmodulo(submodulo, alternativas(nil, true, false)))
}

func TestUnaCasillaUnicaOpcionalPuedeQuedarSinMarcar(t *testing.T) {
	submodulo := submoduloAlternativas()
	submodulo.Hijos[0].Hijos[1].Dato.Requerido = false
	validador := NuevoValidador(catalogosDePrueba())

	exigirSinErrores(t, "alternativas", validador.ValidarSubmodulo(submodulo, alternativas(nil, false)))
	exigirError(t, validador.ValidarSubmodulo(submodulo, alternativas(true, true)), "lista", "está marcada en 2")
}

// Análisis técnico trae la alternativa marcada en Alternativas de solución (referencia con marcada_en).
func TestLaReferenciaMarcadaTomaLaFilaMarcadaEnElOrigen(t *testing.T) {
	analisis := &Nodo{Nombre: "Análisis", Dato: Dato{Clave: "analisis"}, Hijos: []*Nodo{
		campoDePrueba("alternativa", TipoReferencia, Dato{Requerido: true, OpcionesDesde: "alternativas.lista", Mostrar: "descripcion", MarcadaEn: "pasa"}),
	}}
	validador := NuevoValidador(catalogosDePrueba())
	validar := func(datos map[string]Datos, alternativa interface{}) []ErrorCampo {
		datos["analisis"] = Datos{}
		if alternativa != nil {
			datos["analisis"]["alternativa"] = alternativa
		}
		return validador.ValidarSubmodulo(analisis, datos)
	}

	exigirError(t, validar(alternativas(false, false), nil), "alternativa", "aún no se ha marcado")
	exigirError(t, validar(alternativas(false, true), nil), "alternativa", "es obligatorio")
	exigirError(t, validar(alternativas(false, true), "a1"), "alternativa", "se marcó otra")
	exigirSinErrores(t, "analisis", validar(alternativas(false, true), "a2"))
	if !dependenciasDe(analisis)["alternativas"] {
		t.Fatal("el análisis debe depender de las alternativas para recibirlas en relacionadas")
	}
}

// Localización: sedes marcadas de un catálogo (multiple) y «¿Cuál?» solo si se marcó «Otra»
// (solo_si). Con el catálogo de procesos de prueba: «Adecuación» hace de «Otra».
func submoduloSedes() *Nodo {
	return &Nodo{Nombre: "Localización", Dato: Dato{Clave: "localizacion"}, Hijos: []*Nodo{
		{Nombre: "Sedes", Dato: Dato{Clave: "sedes"}, Hijos: []*Nodo{
			campoDePrueba("seleccionadas", TipoCatalogo, Dato{Requerido: true, Catalogo: "procesos", Multiple: true}),
			campoDePrueba("otra", TipoTexto, Dato{Requerido: true, SoloSi: &Condicion{Campo: "seleccionadas", Valor: "Adecuación"}}),
		}},
	}}
}

func validarSedes(t *testing.T, sedes Datos) []ErrorCampo {
	t.Helper()
	submodulo := submoduloSedes()
	catalogos := catalogosDePrueba()
	datos := NuevoNormalizador(catalogos).Normalizar(submodulo, Datos{"sedes": sedes})
	return NuevoValidador(catalogos).ValidarSubmodulo(submodulo, map[string]Datos{"localizacion": datos})
}

func TestLasSedesSonVariasOpcionesDelCatalogo(t *testing.T) {
	exigirSinErrores(t, "localizacion", validarSedes(t, Datos{"seleccionadas": []interface{}{"Administración", "Ampliación"}}))
	exigirError(t, validarSedes(t, Datos{"seleccionadas": []interface{}{}}), "sedes.seleccionadas", "es obligatorio")
	exigirError(t, validarSedes(t, Datos{"seleccionadas": []interface{}{"Administración", "Administración"}}), "sedes.seleccionadas", "dos veces")
	exigirError(t, validarSedes(t, Datos{"seleccionadas": []interface{}{"Sede inventada"}}), "sedes.seleccionadas", "ya no está en la lista")
	exigirError(t, validarSedes(t, Datos{"seleccionadas": "Administración"}), "sedes.seleccionadas", "lista de opciones")
}

func TestCualOtraSedeSoloSeExigeYSeGuardaSiSeMarcoOtra(t *testing.T) {
	exigirError(t, validarSedes(t, Datos{"seleccionadas": []interface{}{"Adecuación"}}), "sedes.otra", "es obligatorio")
	exigirSinErrores(t, "localizacion", validarSedes(t, Datos{"seleccionadas": []interface{}{"Adecuación"}, "otra": "Sede Calle 40"}))

	normalizado := NuevoNormalizador(catalogosDePrueba()).Normalizar(submoduloSedes(), Datos{"sedes": Datos{
		"seleccionadas": []interface{}{"Administración"}, "otra": "texto de cuando estaba marcada",
	}})
	if _, existe := normalizado["sedes"].(Datos)["otra"]; existe {
		t.Fatal("sin «Otra» marcada, el texto de ¿cuál? no se debe guardar")
	}
}

func submoduloNecesidades(operacion string) *Nodo {
	return &Nodo{Nombre: "Necesidades", Dato: Dato{Clave: "necesidades"}, Hijos: []*Nodo{
		campoDePrueba("demanda", TipoNumerico, Dato{Requerido: true}),
		campoDePrueba("oferta", TipoNumerico, Dato{Requerido: true}),
		campoDePrueba("deficit", TipoCalculado, Dato{Operacion: operacion, Operandos: []string{"oferta", "demanda"}}),
	}}
}

func TestElDeficitEsLaOfertaMenosLaDemandaYSeIgnoraElDelCliente(t *testing.T) {
	normalizador := NuevoNormalizador(catalogosDePrueba())
	submodulo := submoduloNecesidades(OperacionResta)

	normalizado := normalizador.Normalizar(submodulo, Datos{"demanda": 0.1, "oferta": 0.3, "deficit": 999.0})
	if normalizado["deficit"] != 0.2 {
		t.Fatalf("deficit = %v", normalizado["deficit"])
	}
	normalizado = normalizador.Normalizar(submodulo, Datos{"demanda": 80.0, "oferta": 30.0})
	if normalizado["deficit"] != -50.0 {
		t.Fatalf("deficit = %v", normalizado["deficit"])
	}
	exigirSinErrores(t, "necesidades", NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(submodulo, map[string]Datos{"necesidades": normalizado}))
}

func TestSinOperandosElDeficitQuedaVacioYSoloSeReportaElOperando(t *testing.T) {
	submodulo := submoduloNecesidades(OperacionResta)
	normalizado := NuevoNormalizador(catalogosDePrueba()).Normalizar(submodulo, Datos{"demanda": 10.0, "deficit": 5.0})
	if _, existe := normalizado["deficit"]; existe {
		t.Fatalf("sin oferta no debe haber déficit: %v", normalizado["deficit"])
	}

	errores := NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(submodulo, map[string]Datos{"necesidades": normalizado})
	if len(errores) != 1 {
		t.Fatalf("se esperaba solo el error de la oferta: %+v", errores)
	}
	exigirError(t, errores, "oferta", "es obligatorio")
}

// Cadena de valor: el costo total de cada actividad es la suma de sus periodos (vacío = 0).
func TestElCostoTotalEsLaSumaDeLosPeriodos(t *testing.T) {
	periodos := []string{"periodo_0", "periodo_1", "periodo_2"}
	hijos := []*Nodo{campoDePrueba("costo_total", TipoCalculado, Dato{Requerido: true, Operacion: OperacionSuma, Operandos: periodos})}
	for _, periodo := range periodos {
		hijos = append(hijos, campoDePrueba(periodo, TipoNumerico, Dato{}))
	}
	submodulo := &Nodo{Nombre: "Cadena", Dato: Dato{Clave: "cadena_valor"}, Hijos: hijos}
	normalizador, validador := NuevoNormalizador(catalogosDePrueba()), NuevoValidador(catalogosDePrueba())

	normalizado := normalizador.Normalizar(submodulo, Datos{"periodo_0": 1000.5, "periodo_2": 0.25, "costo_total": 1.0})
	if normalizado["costo_total"] != 1000.75 {
		t.Fatalf("costo_total = %v", normalizado["costo_total"])
	}
	exigirSinErrores(t, "cadena_valor", validador.ValidarSubmodulo(submodulo, map[string]Datos{"cadena_valor": normalizado}))

	vacio := normalizador.Normalizar(submodulo, Datos{})
	if _, existe := vacio["costo_total"]; existe {
		t.Fatalf("sin periodos no debe haber costo total: %v", vacio["costo_total"])
	}
	exigirError(t, validador.ValidarSubmodulo(submodulo, map[string]Datos{"cadena_valor": vacio}), "costo_total", "es obligatorio")
}

// Ingresos y beneficios: el valor total de cada año es la cantidad por el valor unitario; si
// falta alguno, queda sin valor (no es obligatorio: lo obligatorio son sus operandos).
func TestElValorTotalEsLaCantidadPorElValorUnitario(t *testing.T) {
	dos := 2
	submodulo := &Nodo{Nombre: "Ingresos", Dato: Dato{Clave: "ingresos_beneficios"}, Hijos: []*Nodo{
		campoDePrueba("cantidad", TipoNumerico, Dato{Requerido: true}),
		campoDePrueba("valor_unitario", TipoNumerico, Dato{Requerido: true}),
		campoDePrueba("valor_total", TipoCalculado, Dato{Operacion: OperacionMultiplicacion, Operandos: []string{"cantidad", "valor_unitario"}, Decimales: &dos}),
	}}
	normalizador := NuevoNormalizador(catalogosDePrueba())

	normalizado := normalizador.Normalizar(submodulo, Datos{"cantidad": 1.1, "valor_unitario": 3.0, "valor_total": 1.0})
	if normalizado["valor_total"] != 3.3 { // sin el ruido de 1,1 × 3 = 3,3000000000000003
		t.Fatalf("valor_total = %v", normalizado["valor_total"])
	}
	exigirSinErrores(t, "ingresos_beneficios", NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(submodulo, map[string]Datos{"ingresos_beneficios": normalizado}))

	if cero := normalizador.Normalizar(submodulo, Datos{"cantidad": 0.0, "valor_unitario": 500.0}); cero["valor_total"] != 0.0 {
		t.Fatalf("con cantidad 0 el valor total es 0: %v", cero["valor_total"])
	}
	if incompleto := normalizador.Normalizar(submodulo, Datos{"cantidad": 3.0}); incompleto["valor_total"] != nil {
		t.Fatalf("sin valor unitario no debe haber valor total: %v", incompleto["valor_total"])
	}
}

// Cadena de valor muestra arriba el objetivo general de Objetivos: no se guarda ni se valida,
// pero el submódulo lo recibe en relacionadas.
func TestUnCalculadoDesdeOtroSubmoduloNoSeGuardaYEsDependencia(t *testing.T) {
	submodulo := &Nodo{Nombre: "Cadena", Dato: Dato{Clave: "cadena_valor"}, Hijos: []*Nodo{
		campoDePrueba("objetivo_general", TipoCalculado, Dato{Requerido: true, Desde: "objetivos.objetivo_central"}),
	}}

	normalizado := NuevoNormalizador(catalogosDePrueba()).Normalizar(submodulo, Datos{"objetivo_general": "lo que mande el cliente"})
	if _, existe := normalizado["objetivo_general"]; existe {
		t.Fatalf("no se debe guardar: %v", normalizado["objetivo_general"])
	}
	exigirSinErrores(t, "cadena_valor", NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(submodulo, map[string]Datos{"cadena_valor": normalizado}))
	if !dependenciasDe(submodulo)["objetivos"] {
		t.Fatal("cadena_valor debe recibir objetivos en relacionadas")
	}
}

func TestUnaOperacionDesconocidaEsUnErrorDeConfiguracion(t *testing.T) {
	submodulo := submoduloNecesidades("division")
	errores := NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(submodulo, map[string]Datos{"necesidades": {"demanda": 1.0, "oferta": 2.0}})

	exigirError(t, errores, "deficit", "mal configurado")
}

// Un espejo puede salir de un grupo anidado ("submodulo.grupo.subgrupo"): una fila por
// cada fila del subgrupo en todas las filas del grupo (p. ej. Plurianual por actividad).
func TestUnEspejoDeUnGrupoAnidadoExigeUnaFilaPorCadaFilaInterna(t *testing.T) {
	campo := func(clave, tipo string, extra Dato) *Nodo {
		extra.Clave, extra.Tipo, extra.Requerido = clave, tipo, true
		return &Nodo{Nombre: clave, Dato: extra}
	}
	destino := &Nodo{Nombre: "Destino", Dato: Dato{Clave: "submodulo_destino"}, Hijos: []*Nodo{{
		Nombre: "Por hijo",
		Dato:   Dato{Clave: "por_hijo", Tipo: TipoRepetible, Requerido: true, EspejoDe: "submodulo_origen.padres.hijos", Mostrar: "descripcion"},
		Hijos:  []*Nodo{campo("valor", TipoTexto, Dato{})},
	}}}
	datos := datosConReferencias()
	datos["submodulo_origen"]["padres"] = append(datos["submodulo_origen"]["padres"].([]interface{}), map[string]interface{}{
		"id": "p2", "descripcion": "Padre 2",
		"hijos": []interface{}{map[string]interface{}{"id": "h3", "descripcion": "Hijo tres"}},
	})
	fila := func(id, origen string) map[string]interface{} {
		return map[string]interface{}{"id": id, "origen": origen, "valor": "x"}
	}
	validador := NuevoValidador(catalogosDePrueba())

	datos["submodulo_destino"] = Datos{"por_hijo": []interface{}{fila("x1", "h1"), fila("x3", "h3")}}
	exigirError(t, validador.ValidarSubmodulo(destino, datos), "por_hijo", "Hijo dos")

	datos["submodulo_destino"] = Datos{"por_hijo": []interface{}{fila("x1", "h1"), fila("x2", "h2"), fila("x3", "h3"), fila("x9", "p1")}}
	exigirError(t, validador.ValidarSubmodulo(destino, datos), "por_hijo[x9]", "ya no existe")

	datos["submodulo_destino"] = Datos{"por_hijo": []interface{}{fila("x1", "h1"), fila("x2", "h2"), fila("x3", "h3")}}
	exigirSinErrores(t, "submodulo_destino", validador.ValidarSubmodulo(destino, datos))
}
