package iniciativaInversionHelper

import (
	"errors"
	"testing"
	"time"
)

func TestLaIniciativaDeEjemploNoTieneErrores(t *testing.T) {
	for submodulo, errores := range validarIniciativa(t, definicionDePrueba(t), iniciativaDeEjemplo(t)) {
		exigirSinErrores(t, submodulo, errores)
	}
}

// Objetivos tiene un objetivo específico por cada causa directa de la Problemática (espejo).
func TestLosObjetivosEspecificosExigenUnaFilaPorCadaCausaDirecta(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	iniciativa["problematica"] = problematicaConOtraCausa(t)

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["objetivos"], "objetivos_especificos", "Nueva causa directa.")
}

func TestCadaCausaDirectaNecesitaAlMenosUnaCausaIndirecta(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	filaConID(t, filasEn(t, iniciativa["problematica"], "causas_directas"), "cd-2")["causas_indirectas"] = []interface{}{}

	exigirError(t, validarIniciativa(t, definicionDePrueba(t), iniciativa)["problematica"], "causas_directas[cd-2].causas_indirectas", "mínimo 1")
}

// La Cadena de valor tiene una sección por cada objetivo específico (espejo de Objetivos).
func TestLaCadenaDeValorExigeUnaSeccionPorCadaObjetivoEspecifico(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	secciones := filasEn(t, iniciativa["cadena_valor"], "productos_por_objetivo")
	iniciativa["cadena_valor"]["productos_por_objetivo"] = secciones[:2] // sin la de oe-3

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["cadena_valor"]
	exigirError(t, errores, "productos_por_objetivo", "Fortalecer los componentes tecnológicos según")
}

func TestLaCadenaDeValorRechazaSeccionesDeObjetivosInexistentes(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	filaConID(t, filasEn(t, iniciativa["cadena_valor"], "productos_por_objetivo"), "po-3")["origen"] = "oe-9"

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["cadena_valor"]
	exigirError(t, errores, "productos_por_objetivo[po-3]", "ya no existe")
	exigirError(t, errores, "productos_por_objetivo", "Fortalecer los componentes tecnológicos según")
}

func TestCadaAlternativaExigeSuDescripcion(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	delete(filaConID(t, filasEn(t, iniciativa["alternativas"], "alternativas_solucion"), "as-2"), "descripcion_alternativa")

	errores := validarIniciativa(t, definicionDePrueba(t), iniciativa)["alternativas"]
	exigirError(t, errores, "alternativas_solucion[as-2].descripcion_alternativa", "es obligatorio")
}

func TestElProgramaDebePertenecerAlSectorElegido(t *testing.T) {
	iniciativa := iniciativaDeEjemplo(t)
	iniciativa["articulacion_pnd"]["programa"] = "4599" // programa del sector 45, no del 22

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

// En la plantilla la población es opcional; una matriz obligatoria reporta las celdas que faltan.
func TestLaPoblacionEsOpcionalYUnaMatrizObligatoriaReportaLasCeldasFaltantes(t *testing.T) {
	poblacion := definicionDePrueba(t).Submodulo("poblacion")
	validador := NuevoValidador(catalogosDePrueba())
	vacia := map[string]Datos{"poblacion": {"afectada": map[string]interface{}{}}}
	exigirSinErrores(t, "poblacion", validador.ValidarSubmodulo(poblacion, vacia))

	obligatoria := *poblacion
	afectada := *poblacion.Hijo("afectada")
	afectada.Dato.Requerido = true
	obligatoria.Hijos = []*Nodo{&afectada}
	exigirError(t, validador.ValidarSubmodulo(&obligatoria, vacia), "afectada", "faltan 170 de 170")
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

const nombreSugerido = "FORTALECIMIENTO DE LA INFRAESTRUCTURA TECNOLÓGICA Y DE LOS COMPONENTES DE TRANSFORMACIÓN DIGITAL DE LA UNIVERSIDAD DISTRITAL FRANCISCO JOSÉ DE CALDAS EN BOGOTÁ D.C."

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

	problematica := NuevoNormalizador(catalogos).Normalizar(definicion.Submodulo("problematica"), ejemplo["problematica"])
	if len(filasEn(t, filaConID(t, filasEn(t, problematica, "causas_directas"), "cd-3"), "causas_indirectas")) != 3 {
		t.Fatal("se perdieron las causas indirectas anidadas")
	}
	objetivos := NuevoNormalizador(catalogos).Normalizar(definicion.Submodulo("objetivos"), ejemplo["objetivos"])
	if fila := filaConID(t, filasEn(t, objetivos, "objetivos_especificos"), "oe-2"); fila["origen"] != "cd-2" {
		t.Fatalf("origen = %v", fila["origen"])
	}
}

// ---------------------------------------------------------------------------
// Referencias "origen.*": se prueban con una definición pequeña.
// ---------------------------------------------------------------------------

func definicionConReferencias() *Definicion {
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
				campo("hijo", TipoReferencia, Dato{OpcionesDesde: "origen.hijos", Mostrar: "descripcion"}))),
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
	definicion := definicionConReferencias()
	errores := NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(definicion.Submodulo("submodulo_destino"), datosConReferencias("h1", "h9"))

	exigirError(t, errores, "espejos[e1].acciones[a2].hijo", "no es válida")
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
// cada fila del subgrupo en todas las filas del grupo (p. ej. el Cronograma, una fila por
// cada actividad de los productos de la Cadena de valor).
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

// Cronograma: una fila por actividad de la Cadena de valor. La programación física de un año se
// exige solo si la actividad tiene costo (distinto de 0) en ese periodo (obligatorio_si), y la
// financiera solo muestra el costo de la fila de origen (desde "origen.<campo>"): no se guarda.
func submoduloCronograma() *Nodo {
	return &Nodo{Nombre: "Cronograma", Dato: Dato{Clave: "cronograma"}, Hijos: []*Nodo{
		campoDePrueba("actividades", TipoRepetible, Dato{Requerido: true, EspejoDe: "cadena_valor.objetivos.actividades", Mostrar: "actividad"},
			&Nodo{Nombre: "Física", Dato: Dato{Clave: "fisica"}, Hijos: []*Nodo{
				campoDePrueba("anio_0", TipoTexto, Dato{ObligatorioSi: "origen.periodo_0"}),
				campoDePrueba("anio_1", TipoTexto, Dato{ObligatorioSi: "origen.periodo_1"}),
				campoDePrueba("anio_2", TipoTexto, Dato{ObligatorioSi: "origen.periodo_2"}),
			}},
			&Nodo{Nombre: "Financiera", Dato: Dato{Clave: "financiera"}, Hijos: []*Nodo{
				campoDePrueba("anio_0", TipoCalculado, Dato{Desde: "origen.periodo_0"}),
			}}),
	}}
}

func validarCronograma(t *testing.T, fisica Datos) ([]ErrorCampo, Datos) {
	t.Helper()
	submodulo := submoduloCronograma()
	catalogos := catalogosDePrueba()
	cronograma := NuevoNormalizador(catalogos).Normalizar(submodulo, Datos{"actividades": []interface{}{
		map[string]interface{}{"id": "a1", "origen": "act-1", "fisica": fisica, "financiera": Datos{"anio_0": 999.0}},
	}})
	iniciativa := map[string]Datos{
		"cadena_valor": {"objetivos": []interface{}{map[string]interface{}{"id": "o1", "actividades": []interface{}{
			map[string]interface{}{"id": "act-1", "actividad": "Adecuar aulas", "periodo_0": 500.0, "periodo_1": 0.0},
		}}}},
		"cronograma": cronograma,
	}
	return NuevoValidador(catalogos).ValidarSubmodulo(submodulo, iniciativa), cronograma
}

func TestLaProgramacionFisicaSoloSeExigeEnLosAniosConCosto(t *testing.T) {
	errores, _ := validarCronograma(t, Datos{})
	if len(errores) != 1 {
		t.Fatalf("solo el año 0 tiene costo; se esperaba un error: %+v", errores)
	}
	exigirError(t, errores, "actividades[a1].fisica.anio_0", "es obligatorio")

	errores, cronograma := validarCronograma(t, Datos{"anio_0": "Enero - Diciembre", "anio_2": "Marzo"})
	exigirSinErrores(t, "cronograma", errores)
	fila := filasEn(t, cronograma, "actividades")[0].(map[string]interface{})
	if financiera := fila["financiera"].(Datos); len(financiera) != 0 {
		t.Fatalf("la programación financiera se trae de la fila de origen y no se guarda: %v", financiera)
	}
	if !dependenciasDe(submoduloCronograma())["cadena_valor"] {
		t.Fatal("el cronograma debe recibir la cadena de valor en relacionadas")
	}
}

// Presupuesto inicial: un espejo de las actividades en el que todo se trae (desde). Solo guarda id y
// origen de cada fila, queda sin errores con una fila por actividad y avisa si falta alguna.
func TestUnResumenTraidoSoloGuardaLasFilasYExigeUnaPorActividad(t *testing.T) {
	submodulo := &Nodo{Nombre: "Presupuesto", Dato: Dato{Clave: "presupuesto_inicial"}, Hijos: []*Nodo{
		campoDePrueba("actividades", TipoRepetible, Dato{Requerido: true, EspejoDe: "cadena_valor.objetivos.actividades", Mostrar: "actividad"},
			campoDePrueba("producto", TipoCalculado, Dato{Desde: "cadena_valor.objetivos.producto", Catalogo: "productos"}),
			campoDePrueba("anio_0", TipoCalculado, Dato{Desde: "origen.periodo_0"})),
	}}
	catalogos := catalogosDePrueba()
	normalizar := func(filas ...interface{}) Datos {
		return NuevoNormalizador(catalogos).Normalizar(submodulo, Datos{"actividades": filas})
	}
	iniciativa := map[string]Datos{
		"cadena_valor": {"objetivos": []interface{}{map[string]interface{}{"id": "o1", "producto": "2201001", "actividades": []interface{}{
			map[string]interface{}{"id": "act-1", "actividad": "Adecuar aulas", "periodo_0": 500.0},
			map[string]interface{}{"id": "act-2", "actividad": "Comprar equipos", "periodo_0": 250.0},
		}}}},
	}
	validador := NuevoValidador(catalogos)

	iniciativa["presupuesto_inicial"] = normalizar(
		map[string]interface{}{"id": "pi-1", "origen": "act-1", "producto": "lo que mande el cliente", "anio_0": 1.0},
		map[string]interface{}{"id": "pi-2", "origen": "act-2"},
	)
	exigirSinErrores(t, "presupuesto_inicial", validador.ValidarSubmodulo(submodulo, iniciativa))
	if fila := filasEn(t, iniciativa["presupuesto_inicial"], "actividades")[0].(map[string]interface{}); len(fila) != 2 {
		t.Fatalf("de cada fila solo se guardan id y origen: %v", fila)
	}

	iniciativa["presupuesto_inicial"] = normalizar(map[string]interface{}{"id": "pi-1", "origen": "act-1"})
	exigirError(t, validador.ValidarSubmodulo(submodulo, iniciativa), "actividades", "Comprar equipos")
}

// Trazadores presupuestales: espejo de un catálogo en el que se llenan solo los que aplican.
// Grupos étnicos se desglosa en sus comunidades (espejo de otro catálogo con solo_si) y los
// demás llevan discriminado y beneficiarios propios (solo_si con distinto). Se pueden agregar
// hasta dos trazadores que no están en el catálogo («Otro», con su «¿Cuál?»).
type repoCatalogosPorNombre map[string][]ItemCatalogo

func (r repoCatalogosPorNombre) ElementosDeCatalogo(nombre string) ([]ItemCatalogo, error) {
	return r[nombre], nil
}

func validarTrazadores(t *testing.T, filas ...interface{}) ([]ErrorCampo, Datos) {
	t.Helper()
	uno := 1.0
	noEtnicos := &Condicion{Campo: "origen", Valor: "etn", Distinto: true}
	submodulo := &Nodo{Nombre: "Trazadores", Dato: Dato{Clave: "trazadores_presupuestales"}, Hijos: []*Nodo{
		campoDePrueba("trazadores", TipoRepetible, Dato{EspejoDe: "catalogo.trazadores", Mostrar: "nombre", FilasOpcionales: true,
			FilasAdicionales: &FilasAdicionales{Origen: "Otro", Max: 2, Nombre: "Otro: ¿Cuál?", Etiqueta: "Otro trazador"}},
			campoDePrueba("cual", TipoTexto, Dato{Requerido: true, SoloSi: &Condicion{Campo: "origen", Valor: "Otro"}}),
			campoDePrueba("descripcion", TipoTexto, Dato{Requerido: true}),
			campoDePrueba("discriminado", TipoTexto, Dato{SoloSi: noEtnicos}),
			campoDePrueba("beneficiarios", TipoNumerico, Dato{Requerido: true, SoloSi: noEtnicos}),
			campoDePrueba("grupos", TipoRepetible, Dato{EspejoDe: "catalogo.grupos", Mostrar: "nombre", FilasOpcionales: true, Min: &uno,
				SoloSi: &Condicion{Campo: "origen", Valor: "etn"}},
				campoDePrueba("beneficiarios", TipoNumerico, Dato{Requerido: true}))),
	}}
	catalogos := NuevosCatalogosConCache(repoCatalogosPorNombre{
		"trazadores": {{Codigo: "cul", Nombre: "Cultura ciudadana", Descripcion: "Busca consolidar un modelo…"}, {Codigo: "etn", Nombre: "Grupos étnicos"}},
		"grupos":     {{Codigo: "ind", Nombre: "Comunidades indígenas"}, {Codigo: "rrom", Nombre: "Comunidades Rrom"}},
	}, time.Hour)
	datos := NuevoNormalizador(catalogos).Normalizar(submodulo, Datos{"trazadores": filas})
	return NuevoValidador(catalogos).ValidarSubmodulo(submodulo, map[string]Datos{"trazadores_presupuestales": datos}), datos
}

func TestLosTrazadoresSeLlenanSiAplicanYGruposEtnicosSeDesglosaPorComunidad(t *testing.T) {
	sinNinguno, _ := validarTrazadores(t)
	exigirSinErrores(t, "trazadores_presupuestales", sinNinguno)

	cultura, _ := validarTrazadores(t, Datos{"id": "t1", "origen": "cul", "descripcion": "Talleres de convivencia"})
	exigirError(t, cultura, "trazadores[t1].beneficiarios", "es obligatorio")

	etnicos := func(grupos ...interface{}) Datos {
		return Datos{"id": "t2", "origen": "etn", "descripcion": "Becas", "discriminado": "no aplica", "beneficiarios": 9.0, "grupos": grupos}
	}
	sinComunidades, datos := validarTrazadores(t, etnicos(Datos{"id": "g1", "origen": "ind"}, Datos{"id": "g2", "origen": "rrom"}))
	exigirError(t, sinComunidades, "trazadores[t2].grupos", "al menos 1")
	fila := filasEn(t, datos, "trazadores")[0].(map[string]interface{})
	if _, existe := fila["discriminado"]; existe {
		t.Fatal("en grupos étnicos no se guarda el discriminado: lo reemplazan las comunidades")
	}
	if _, existe := fila["beneficiarios"]; existe {
		t.Fatal("en grupos étnicos los beneficiarios van por comunidad")
	}

	conIndigenas, _ := validarTrazadores(t, etnicos(Datos{"id": "g1", "origen": "ind", "beneficiarios": 40.0}, Datos{"id": "g2", "origen": "rrom"}))
	exigirSinErrores(t, "trazadores_presupuestales", conIndigenas)
}

func TestSePuedenAgregarHastaDosTrazadoresQueNoEstanEnElCatalogo(t *testing.T) {
	otro := func(id string, campos Datos) Datos {
		fila := Datos{"id": id, "origen": "Otro"}
		for clave, valor := range campos {
			fila[clave] = valor
		}
		return fila
	}
	completo := Datos{"cual": "Trazador de Paz", "descripcion": "Cátedra de paz", "beneficiarios": 30.0}

	vacia, _ := validarTrazadores(t, otro("o1", nil))
	exigirSinErrores(t, "trazadores_presupuestales", vacia)

	sinCual, _ := validarTrazadores(t, otro("o1", Datos{"descripcion": "Cátedra de paz", "beneficiarios": 30.0}))
	exigirError(t, sinCual, "trazadores[o1].cual", "es obligatorio")

	dos, datos := validarTrazadores(t, otro("o1", completo), otro("o2", completo))
	exigirSinErrores(t, "trazadores_presupuestales", dos)
	if fila := filasEn(t, datos, "trazadores")[0].(map[string]interface{}); fila["origen"] != "Otro" || fila["cual"] != "Trazador de Paz" {
		t.Fatalf("la fila agregada conserva su origen y su «¿Cuál?»: %v", fila)
	}

	tres, _ := validarTrazadores(t, otro("o1", completo), otro("o2", completo), otro("o3", completo))
	exigirError(t, tres, "trazadores", "máximo 2 (otro trazador), hay 3")

	cultura, datos := validarTrazadores(t, Datos{"id": "t1", "origen": "cul", "cual": "No aplica", "descripcion": "Talleres", "beneficiarios": 5.0})
	exigirSinErrores(t, "trazadores_presupuestales", cultura)
	if _, existe := filasEn(t, datos, "trazadores")[0].(map[string]interface{})["cual"]; existe {
		t.Fatal("el «¿Cuál?» solo se guarda en las filas agregadas")
	}
}

// Fuentes de financiación: las otras fuentes se agregan sin límite (filas_adicionales sin max) y el
// espejo exige al menos una fuente con valores.
func TestSinMaxSePuedenAgregarTodasLasOtrasFuentesYSeExigeAlMenosUna(t *testing.T) {
	uno := 1.0
	submodulo := &Nodo{Nombre: "Fuentes", Dato: Dato{Clave: "fuentes_financiacion"}, Hijos: []*Nodo{
		campoDePrueba("fuentes", TipoRepetible, Dato{EspejoDe: "catalogo.fuentes", Mostrar: "nombre", FilasOpcionales: true, Min: &uno,
			FilasAdicionales: &FilasAdicionales{Origen: "Otra"}},
			campoDePrueba("cual", TipoTexto, Dato{Requerido: true, SoloSi: &Condicion{Campo: "origen", Valor: "Otra"}}),
			campoDePrueba("anio_0", TipoNumerico, Dato{})),
	}}
	catalogos := NuevosCatalogosConCache(repoCatalogosPorNombre{"fuentes": {{Codigo: "1-100-I013", Nombre: "VA-Estampilla Universidad Distrital"}}}, time.Hour)
	validar := func(filas ...interface{}) []ErrorCampo {
		return NuevoValidador(catalogos).ValidarSubmodulo(submodulo, map[string]Datos{"fuentes_financiacion": {"fuentes": filas}})
	}

	exigirError(t, validar(), "fuentes", "al menos 1")
	otras := []interface{}{}
	for _, id := range []string{"o1", "o2", "o3", "o4"} {
		otras = append(otras, Datos{"id": id, "origen": "Otra", "cual": "Recursos de " + id, "anio_0": 10.0})
	}
	exigirSinErrores(t, "fuentes_financiacion", validar(otras...))
}

func TestElOrigenDeUnEspejoDeCatalogoTraeLaDescripcionDelElemento(t *testing.T) {
	filas := filasDeCatalogo(NuevoCatalogo("trazadores", []ItemCatalogo{
		{Codigo: "cul", Nombre: "Cultura ciudadana", Descripcion: "Busca consolidar un modelo…"},
		{Codigo: "juv", Nombre: "Juventud"},
	}))
	if filas[0]["descripcion"] != "Busca consolidar un modelo…" {
		t.Fatalf("la fila de origen debe traer la descripción: %v", filas[0])
	}
	if _, existe := filas[1]["descripcion"]; existe {
		t.Fatalf("sin descripción no se agrega la clave: %v", filas[1])
	}
}

func TestObligatorioSiFueraDeUnEspejoEsUnErrorDeConfiguracion(t *testing.T) {
	submodulo := &Nodo{Nombre: "Suelto", Dato: Dato{Clave: "suelto"}, Hijos: []*Nodo{
		campoDePrueba("anio_0", TipoTexto, Dato{ObligatorioSi: "origen.periodo_0"}),
	}}
	errores := NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(submodulo, map[string]Datos{"suelto": {}})

	exigirError(t, errores, "anio_0", "mal configurado")
}

// ---------------------------------------------------------------------------
// Espejo de un catálogo con filas opcionales (Resultados esperados de Objetivos)
// ---------------------------------------------------------------------------

func submoduloResultados() *Nodo {
	return &Nodo{Nombre: "Objetivos", Dato: Dato{Clave: "objetivos"}, Hijos: []*Nodo{
		campoDePrueba("resultados", TipoRepetible, Dato{Requerido: true, EspejoDe: "catalogo.resultados", Mostrar: "nombre",
			FilasOpcionales: true, FilasObligatorias: []string{"est"}},
			campoDePrueba("cual", TipoTexto, Dato{Requerido: true, SoloSi: &Condicion{Campo: "origen", Valor: "Otro"}}),
			campoDePrueba("cantidad", TipoNumerico, Dato{Requerido: true}),
			campoDePrueba("sede", TipoTexto, Dato{Requerido: true})),
	}}
}

func validarResultados(repo RepositorioElementosCatalogo, filas ...interface{}) []ErrorCampo {
	return NuevoValidador(NuevosCatalogosConCache(repo, time.Hour)).ValidarSubmodulo(submoduloResultados(),
		map[string]Datos{"objetivos": {"resultados": filas}})
}

func resultadosDePrueba() *repoCatalogosContado {
	return &repoCatalogosContado{items: []ItemCatalogo{
		{Codigo: "est", Nombre: "Número de estudiantes beneficiados"},
		{Codigo: "doc", Nombre: "Número de docentes beneficiados"},
		{Codigo: "aul", Nombre: "Número de aulas"},
		{Codigo: "Otro", Nombre: "Otro: ¿Cuál?"},
	}}
}

func filaResultado(origen string, campos Datos) Datos {
	fila := Datos{"id": "r-" + origen, "origen": origen}
	for clave, valor := range campos {
		fila[clave] = valor
	}
	return fila
}

func estudiantes() Datos {
	return filaResultado("est", Datos{"cantidad": 26257.0, "sede": "Todas las sedes"})
}

func TestUnEspejoDeCatalogoSoloExigeLasFilasQueSeLlenanYLasObligatorias(t *testing.T) {
	exigirSinErrores(t, "objetivos", validarResultados(resultadosDePrueba(),
		estudiantes(),
		filaResultado("doc", Datos{"cantidad": 120.0, "sede": "Macarena"}),
		filaResultado("aul", Datos{"sede": "  "})))
	// El cliente puede no mandar las filas opcionales vacías.
	exigirSinErrores(t, "objetivos", validarResultados(resultadosDePrueba(), estudiantes()))
}

func TestUnaFilaObligatoriaDelCatalogoNoSePuedeDejarVacia(t *testing.T) {
	errores := validarResultados(resultadosDePrueba(), filaResultado("est", nil))
	exigirError(t, errores, "resultados[r-est].cantidad", "es obligatorio")
	exigirError(t, errores, "resultados[r-est].sede", "es obligatorio")
	exigirError(t, validarResultados(resultadosDePrueba()), "resultados", "«Número de estudiantes beneficiados»")
}

func TestUnaFilaDelCatalogoQueSeEmpiezaALlenarSeExigeCompleta(t *testing.T) {
	errores := validarResultados(resultadosDePrueba(), estudiantes(), filaResultado("doc", Datos{"cantidad": 120.0}))
	exigirError(t, errores, "resultados[r-doc].sede", "es obligatorio")
}

func TestUnEspejoDeCatalogoExigeElMinimoDeFilasLlenas(t *testing.T) {
	submodulo := submoduloResultados()
	minimo := 2.0
	submodulo.Hijos[0].Dato.Min = &minimo
	validador := NuevoValidador(NuevosCatalogosConCache(resultadosDePrueba(), time.Hour))

	errores := validador.ValidarSubmodulo(submodulo, map[string]Datos{"objetivos": {"resultados": []interface{}{estudiantes()}}})
	exigirError(t, errores, "resultados", "diligencie al menos 2 (hay 1)")
}

func TestUnEspejoDeCatalogoRechazaElementosQueYaNoEstanEnLaLista(t *testing.T) {
	errores := validarResultados(resultadosDePrueba(), estudiantes(),
		filaResultado("viejo", Datos{"cantidad": 1.0, "sede": "Macarena"}))
	exigirError(t, errores, "resultados[r-viejo]", "ya no existe")
}

func TestSiNoSePuedeConsultarElCatalogoDelEspejoSeAvisa(t *testing.T) {
	errores := validarResultados(&repoCatalogosContado{err: errors.New("crud caído")}, estudiantes())
	exigirError(t, errores, "resultados", "no se pudieron consultar")
}

func TestUnEspejoDeCatalogoNoEsDependenciaDeOtroSubmodulo(t *testing.T) {
	submodulo := submoduloResultados()
	if len(dependenciasDe(submodulo)) != 0 {
		t.Fatalf("el catálogo no es un submódulo: %v", dependenciasDe(submodulo))
	}
	if !catalogosEspejadosDe(submodulo)["resultados"] {
		t.Fatal("el submódulo debe recibir los elementos del catálogo resultados en relacionadas")
	}
}

func TestLaFilaOtroDelCatalogoEsOpcionalYSiSeLlenaPideCual(t *testing.T) {
	exigirSinErrores(t, "objetivos", validarResultados(resultadosDePrueba(), estudiantes(), filaResultado("Otro", nil)))

	errores := validarResultados(resultadosDePrueba(), estudiantes(), filaResultado("Otro", Datos{"cantidad": 3.0, "sede": "Macarena"}))
	exigirError(t, errores, "resultados[r-Otro].cual", "es obligatorio")

	exigirSinErrores(t, "objetivos", validarResultados(resultadosDePrueba(), estudiantes(),
		filaResultado("Otro", Datos{"cual": "Número de huertas", "cantidad": 3.0, "sede": "Macarena"})))
}

func TestElCualSoloSeGuardaEnLaFilaOtro(t *testing.T) {
	normalizado := NuevoNormalizador(catalogosDePrueba()).Normalizar(submoduloResultados(), Datos{"resultados": []interface{}{
		filaResultado("doc", Datos{"cual": "no aplica aquí", "cantidad": 1.0, "sede": "Macarena"}),
		filaResultado("Otro", Datos{"cual": "  Número de huertas ", "cantidad": 3.0, "sede": "Macarena"}),
	}})
	filas := normalizado["resultados"].([]interface{})
	if _, existe := filas[0].(Datos)["cual"]; existe {
		t.Fatal("el ¿cuál? no se debe guardar en una fila que no es «Otro»")
	}
	otro := filas[1].(Datos)
	if otro["cual"] != "Número de huertas" || otro["origen"] != "Otro" || otro["id"] != "r-Otro" {
		t.Fatalf("la fila «Otro» debe conservar su ¿cuál?, su origen y su id: %v", otro)
	}
}

func TestLosMensajesNombranLaFilaSinElNoDelTitulo(t *testing.T) {
	bienes := campoDePrueba("bienes", TipoRepetible, Dato{Requerido: true, EtiquetaFila: "Bien o servicio a entregar o demanda a satisfacer No."},
		campoDePrueba("descripcion", TipoTexto, Dato{Requerido: true}))
	submodulo := &Nodo{Nombre: "Necesidades", Dato: Dato{Clave: "necesidades"}, Hijos: []*Nodo{bienes}}

	errores := NuevoValidador(catalogosDePrueba()).ValidarSubmodulo(submodulo, map[string]Datos{"necesidades": {}})
	exigirError(t, errores, "bienes", "(bien o servicio a entregar o demanda a satisfacer)")
}
