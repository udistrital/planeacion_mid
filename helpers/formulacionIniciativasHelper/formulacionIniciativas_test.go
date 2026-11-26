package formulacionIniciativasHelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

const (
	rectoria   = `{"Id":7,"Nombre":"RECTORIA"}`
	ingenieria = `{"Id":14,"Nombre":"FACULTAD DE INGENIERIA"}`
	ciencias   = `{"Id":21,"Nombre":"FACULTAD DE CIENCIAS Y EDUCACION"}`
	i3         = `{"Id":209,"Nombre":"INSTITUTO DE INVESTIGACIÓN E INNOVACIÓN EN INGENERÍA -I3+"}`
)

// periodoDePrueba arma un registro como lo devuelve el crud: fechas a medianoche UTC y
// unidades_interes como texto JSON.
func periodoDePrueba(id, inicio, fin string, unidades ...string) map[string]interface{} {
	texto := "["
	for i, unidad := range unidades {
		if i > 0 {
			texto += ","
		}
		texto += unidad
	}
	texto += "]"
	return map[string]interface{}{
		"_id":                 id,
		"tipo_seguimiento_id": "6650c0a1b2c3d4e5f6a7b8c9",
		"fecha_inicio":        inicio + "T00:00:00.000Z",
		"fecha_fin":           fin + "T00:00:00.000Z",
		"unidades_interes":    texto,
		"planes_interes":      "[]",
		"activo":              true,
		"fecha_creacion":      "2026-11-01T15:00:00.000Z",
	}
}

// enBogota devuelve la hora indicada del día en Bogotá.
func enBogota(dia string, hora, minuto int) time.Time {
	fecha, err := time.ParseInLocation("2006-01-02", dia, zonaBogota)
	if err != nil {
		panic(err)
	}
	return fecha.Add(time.Duration(hora)*time.Hour + time.Duration(minuto)*time.Minute)
}

// conVigencia agrega periodo_id al registro, como los que se guardaron con una vigencia.
func conVigencia(registro map[string]interface{}, periodoID string) map[string]interface{} {
	registro["periodo_id"] = periodoID
	return registro
}

func TestUnidadHabilitada(t *testing.T) {
	casos := []struct {
		nombre string
		// periodos va del más reciente al más antiguo, como los entrega el crud.
		periodos         []map[string]interface{}
		dependenciaID    string
		hoy              time.Time
		quiereHabilitada bool
	}{
		{
			nombre: "la unidad está en el segundo registro con fechas vigentes (F1)",
			periodos: []map[string]interface{}{
				periodoDePrueba("A", "2026-11-10", "2026-11-30", rectoria),
				periodoDePrueba("B", "2026-11-20", "2026-12-10", ingenieria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereHabilitada: true,
		},
		{
			nombre: "un registro vencido más antiguo con la unidad no la bloquea (F2)",
			periodos: []map[string]interface{}{
				periodoDePrueba("B", "2026-11-20", "2026-12-10", ingenieria),
				periodoDePrueba("viejo", "2025-11-01", "2025-11-30", rectoria, ingenieria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereHabilitada: true,
		},
		{
			nombre: "rige el registro más reciente aunque uno más antiguo cubra hoy",
			periodos: []map[string]interface{}{
				periodoDePrueba("cierre", "2026-11-10", "2026-11-14", rectoria, ingenieria),
				periodoDePrueba("apertura", "2026-11-17", "2026-11-30", rectoria, ingenieria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-24", 10, 0),
			quiereHabilitada: false,
		},
		{
			nombre: "un registro vencido que no tiene la unidad no la bloquea (F2)",
			periodos: []map[string]interface{}{
				periodoDePrueba("viejo", "2025-11-01", "2025-11-30", rectoria),
				periodoDePrueba("B", "2026-11-20", "2026-12-10", ingenieria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereHabilitada: true,
		},
		{
			nombre: "la unidad solo está en un registro vencido",
			periodos: []map[string]interface{}{
				periodoDePrueba("A", "2026-11-01", "2026-11-15", ingenieria),
				periodoDePrueba("B", "2026-11-20", "2026-12-10", rectoria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereHabilitada: false,
		},
		{
			nombre: "la unidad no está en ningún registro",
			periodos: []map[string]interface{}{
				periodoDePrueba("A", "2026-11-10", "2026-11-30", rectoria, ciencias),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereHabilitada: false,
		},
		{
			nombre:        "hoy es el día de inicio",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", ingenieria)},
			dependenciaID: "14", hoy: enBogota("2026-11-10", 0, 0),
			quiereHabilitada: true,
		},
		{
			nombre:        "hoy es el día de fin, de noche en Bogotá (ya es el día siguiente en UTC)",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", ingenieria)},
			dependenciaID: "14", hoy: enBogota("2026-11-30", 23, 30),
			quiereHabilitada: true,
		},
		{
			nombre:        "el día anterior al inicio",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", ingenieria)},
			dependenciaID: "14", hoy: enBogota("2026-11-09", 23, 59),
			quiereHabilitada: false,
		},
		{
			nombre:        "el día siguiente al fin",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", ingenieria)},
			dependenciaID: "14", hoy: enBogota("2026-12-01", 0, 0),
			quiereHabilitada: false,
		},
		{
			nombre: "sirve un registro guardado con vigencia que cubra hoy",
			periodos: []map[string]interface{}{
				conVigencia(periodoDePrueba("con vigencia", "2026-11-10", "2026-11-30", ingenieria), "11"),
				periodoDePrueba("sin vigencia", "2026-11-10", "2026-11-30", rectoria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereHabilitada: true,
		},
		{
			nombre:        "la unidad 209, con caracteres especiales en el nombre",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", rectoria, i3)},
			dependenciaID: "209", hoy: enBogota("2026-11-25", 10, 0),
			quiereHabilitada: true,
		},
		{
			nombre:        "un Id guardado como texto se compara igual",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", `{"Id":"14","Nombre":"FACULTAD DE INGENIERIA"}`)},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereHabilitada: true,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			habilitada, err := UnidadHabilitada(caso.periodos, caso.dependenciaID, caso.hoy)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if habilitada != caso.quiereHabilitada {
				t.Fatalf("se esperaba habilitada=%v y llegó %v", caso.quiereHabilitada, habilitada)
			}
		})
	}
}

func TestUnidadHabilitadaFallaConDatosInvalidos(t *testing.T) {
	casos := []struct {
		nombre  string
		periodo map[string]interface{}
	}{
		{"unidades_interes no es JSON", map[string]interface{}{"_id": "A", "unidades_interes": "no es json", "fecha_inicio": "2026-11-10T00:00:00.000Z", "fecha_fin": "2026-11-30T00:00:00.000Z"}},
		{"unidades_interes no es texto", map[string]interface{}{"_id": "A", "fecha_inicio": "2026-11-10T00:00:00.000Z", "fecha_fin": "2026-11-30T00:00:00.000Z"}},
		{"fecha_fin inválida", map[string]interface{}{"_id": "A", "unidades_interes": "[" + ingenieria + "]", "fecha_inicio": "2026-11-10T00:00:00.000Z", "fecha_fin": "x"}},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			if _, err := UnidadHabilitada([]map[string]interface{}{caso.periodo}, "14", enBogota("2026-11-25", 10, 0)); err == nil {
				t.Fatal("se esperaba un error")
			}
		})
	}
}

// seleccionDePrueba arma la selección como la lee DefinirFechasIniciativasFormulacion del body.
func seleccionDePrueba(unidades ...string) []map[string]interface{} {
	texto := "["
	for i, unidad := range unidades {
		if i > 0 {
			texto += ","
		}
		texto += unidad
	}
	seleccion, _ := ObtenerArrayUnidadesInteres(map[string]interface{}{"unidades_interes": texto + "]"})
	return seleccion
}

func fecha(valor string) time.Time {
	t, err := time.Parse(time.RFC3339, valor)
	if err != nil {
		panic(err)
	}
	return t
}

// idsDe decodifica unidades_interes y devuelve sus Ids, y falla si hay algo más que Id y Nombre.
func idsDe(t *testing.T, registro map[string]interface{}) []int {
	t.Helper()
	var unidades []map[string]interface{}
	if err := json.Unmarshal([]byte(registro["unidades_interes"].(string)), &unidades); err != nil {
		t.Fatalf("unidades_interes inválido: %v", err)
	}
	ids := []int{}
	for _, unidad := range unidades {
		if len(unidad) != 2 || unidad["Nombre"] == nil {
			t.Fatalf("cada unidad debe tener solo Id y Nombre: %v", unidad)
		}
		ids = append(ids, int(unidad["Id"].(float64)))
	}
	return ids
}

func TestRepartirUnidades(t *testing.T) {
	// esperado describe un registro escrito: su _id (vacío para el que se crea), sus unidades y activo.
	type esperado struct {
		id     string
		ids    []int
		activo bool
	}
	casos := []struct {
		nombre         string
		registros      []map[string]interface{}
		seleccion      []map[string]interface{}
		inicio, fin    string
		quiereActuales []esperado
		quiereCrear    *esperado
	}{
		{
			nombre:    "cerrar con un rango que ya pasó saca a las unidades del registro que las abría",
			registros: []map[string]interface{}{periodoDePrueba("A", "2026-11-17", "2026-11-30", rectoria, ingenieria)},
			seleccion: seleccionDePrueba(rectoria, ingenieria),
			inicio:    "2026-11-10T00:00:00.000Z", fin: "2026-11-14T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{}, false}},
			quiereCrear:    &esperado{"", []int{7, 14}, true},
		},
		{
			nombre:    "un registro guardado con vigencia y el mismo rango es el destino",
			registros: []map[string]interface{}{conVigencia(periodoDePrueba("A", "2026-11-20", "2026-12-10", rectoria), "11")},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7, 14}, true}},
		},
		{
			nombre:    "ejemplo A y B: Ingeniería cambia de rango",
			registros: []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", rectoria, ingenieria)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7}, true}},
			quiereCrear:    &esperado{"", []int{14}, true},
		},
		{
			nombre: "la selección está repartida en dos registros y se quita de los dos (F3)",
			registros: []map[string]interface{}{
				periodoDePrueba("A", "2026-11-10", "2026-11-30", rectoria, ingenieria),
				periodoDePrueba("C", "2026-12-01", "2026-12-15", ciencias, i3),
			},
			seleccion: seleccionDePrueba(ingenieria, i3),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7}, true}, {"C", []int{21}, true}},
			quiereCrear:    &esperado{"", []int{14, 209}, true},
		},
		{
			nombre:    "el registro que queda vacío se inactiva",
			registros: []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", ingenieria)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{}, false}},
			quiereCrear:    &esperado{"", []int{14}, true},
		},
		{
			nombre: "ya existe un registro con las mismas fechas: se agregan ahí sin repetir",
			registros: []map[string]interface{}{
				periodoDePrueba("A", "2026-11-10", "2026-11-30", rectoria, ingenieria),
				periodoDePrueba("B", "2026-11-20", "2026-12-10", ciencias),
			},
			seleccion: seleccionDePrueba(ingenieria, ciencias),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"B", []int{21, 14}, true}, {"A", []int{7}, true}},
		},
		{
			nombre:    "el destino se compara por día aunque el cliente envíe otra hora",
			registros: []map[string]interface{}{periodoDePrueba("B", "2026-11-20", "2026-12-10", ciencias)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T05:00:00.000Z", fin: "2026-12-10T05:00:00.000Z",
			quiereActuales: []esperado{{"B", []int{21, 14}, true}},
		},
		{
			nombre:    "si la selección ya está solo en su destino no se escribe nada",
			registros: []map[string]interface{}{periodoDePrueba("B", "2026-11-20", "2026-12-10", ingenieria, ciencias)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
		},
		{
			nombre:    "no hay registros: se crea uno",
			seleccion: seleccionDePrueba(rectoria, ingenieria),
			inicio:    "2026-11-10T00:00:00.000Z", fin: "2026-11-30T00:00:00.000Z",
			quiereCrear: &esperado{"", []int{7, 14}, true},
		},
		{
			nombre:    "la unidad 209 se encuentra aunque su nombre tenga + (F4)",
			registros: []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", i3, rectoria)},
			seleccion: seleccionDePrueba(i3),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7}, true}},
			quiereCrear:    &esperado{"", []int{209}, true},
		},
		{
			nombre:    "la unidad 209 no se duplica en su destino (F4)",
			registros: []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", i3, rectoria)},
			seleccion: seleccionDePrueba(i3),
			inicio:    "2026-11-10T00:00:00.000Z", fin: "2026-11-30T00:00:00.000Z",
		},
		{
			nombre:    "Ids repetidos en un registro se quitan todos (F6)",
			registros: []map[string]interface{}{periodoDePrueba("A", "2026-11-10", "2026-11-30", ingenieria, ingenieria, rectoria)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7}, true}},
			quiereCrear:    &esperado{"", []int{14}, true},
		},
		{
			nombre:    "Ids repetidos en el destino quedan una sola vez (F6)",
			registros: []map[string]interface{}{periodoDePrueba("B", "2026-11-20", "2026-12-10", ingenieria, ingenieria)},
			seleccion: seleccionDePrueba(ciencias),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"B", []int{14, 21}, true}},
		},
		{
			nombre:    "Ids repetidos en la selección se guardan una sola vez",
			seleccion: seleccionDePrueba(ingenieria, ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereCrear: &esperado{"", []int{14}, true},
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			inicio, fin := fecha(caso.inicio), fecha(caso.fin)
			actualizar, crear := repartirUnidades(caso.registros, caso.seleccion, inicio, fin)

			if len(actualizar) != len(caso.quiereActuales) {
				t.Fatalf("se esperaban %d registros para actualizar y llegaron %d: %v", len(caso.quiereActuales), len(actualizar), actualizar)
			}
			for i, quiere := range caso.quiereActuales {
				registro := actualizar[i]
				if registro["_id"] != quiere.id {
					t.Fatalf("actualizar[%d]: se esperaba el registro %s y llegó %v", i, quiere.id, registro["_id"])
				}
				if ids := idsDe(t, registro); !reflect.DeepEqual(ids, quiere.ids) {
					t.Fatalf("registro %s: se esperaban las unidades %v y llegaron %v", quiere.id, quiere.ids, ids)
				}
				if registro["activo"] != quiere.activo {
					t.Fatalf("registro %s: se esperaba activo=%v y llegó %v", quiere.id, quiere.activo, registro["activo"])
				}
				if registro["fecha_creacion"] == nil || registro["tipo_seguimiento_id"] == nil {
					t.Fatalf("registro %s: el PUT debe llevar el registro completo: %v", quiere.id, registro)
				}
			}

			if caso.quiereCrear == nil {
				if crear != nil {
					t.Fatalf("no se esperaba crear un registro y llegó %v", crear)
				}
				return
			}
			if crear == nil {
				t.Fatal("se esperaba crear un registro")
			}
			if ids := idsDe(t, crear); !reflect.DeepEqual(ids, caso.quiereCrear.ids) {
				t.Fatalf("registro nuevo: se esperaban las unidades %v y llegaron %v", caso.quiereCrear.ids, ids)
			}
			if crear["activo"] != true || crear["planes_interes"] != "[]" {
				t.Fatalf("registro nuevo: se esperaba activo=true y planes_interes=[]: %v", crear)
			}
			if !fecha(crear["fecha_inicio"].(string)).Equal(inicio) || !fecha(crear["fecha_fin"].(string)).Equal(fin) {
				t.Fatalf("registro nuevo: las fechas deben ser las enviadas: %v", crear)
			}
			if _, tieneVigencia := crear["periodo_id"]; tieneVigencia {
				t.Fatalf("registro nuevo: no lleva vigencia: %v", crear)
			}
		})
	}
}

func TestRepartirUnidadesNoModificaLosRegistrosRecibidos(t *testing.T) {
	registros := []map[string]interface{}{
		periodoDePrueba("A", "2026-11-10", "2026-11-30", rectoria, ingenieria),
		periodoDePrueba("B", "2026-11-20", "2026-12-10", ciencias),
	}
	antes := []map[string]interface{}{copiarRegistro(registros[0]), copiarRegistro(registros[1])}

	repartirUnidades(registros, seleccionDePrueba(ingenieria), fecha("2026-11-20T00:00:00.000Z"), fecha("2026-12-10T00:00:00.000Z"))

	if !reflect.DeepEqual(registros, antes) {
		t.Fatalf("los registros recibidos cambiaron:\nantes:   %v\ndespués: %v", antes, registros)
	}
}

func TestEjemploAyBQuedaCadaUnidadEnUnSoloRegistro(t *testing.T) {
	// 1. Rectoría e Ingeniería del 10 al 30 de noviembre: se crea A.
	_, crearA := repartirUnidades(nil, seleccionDePrueba(rectoria, ingenieria), fecha("2026-11-10T00:00:00.000Z"), fecha("2026-11-30T00:00:00.000Z"))
	crearA["_id"] = "A"

	// 2. Ingeniería pasa al 20 de noviembre-10 de diciembre: A se queda con Rectoría y se crea B.
	actualizar, crearB := repartirUnidades([]map[string]interface{}{crearA}, seleccionDePrueba(ingenieria), fecha("2026-11-20T00:00:00.000Z"), fecha("2026-12-10T00:00:00.000Z"))
	if len(actualizar) != 1 || crearB == nil {
		t.Fatalf("se esperaba actualizar A y crear B: %v %v", actualizar, crearB)
	}
	crearB["_id"] = "B"
	periodos := []map[string]interface{}{actualizar[0], crearB}

	// 3. El 25 de noviembre las dos unidades están habilitadas, cada una por su registro.
	for _, dependencia := range []string{"7", "14"} {
		habilitada, err := UnidadHabilitada(periodos, dependencia, enBogota("2026-11-25", 10, 0))
		if err != nil || !habilitada {
			t.Fatalf("la unidad %s debería estar habilitada: habilitada=%v err=%v", dependencia, habilitada, err)
		}
	}
	// El 5 de diciembre solo Ingeniería sigue habilitada.
	if habilitada, _ := UnidadHabilitada(periodos, "7", enBogota("2026-12-05", 10, 0)); habilitada {
		t.Fatal("Rectoría no debería estar habilitada el 5 de diciembre")
	}
	if habilitada, _ := UnidadHabilitada(periodos, "14", enBogota("2026-12-05", 10, 0)); !habilitada {
		t.Fatal("Ingeniería debería estar habilitada el 5 de diciembre")
	}
}

// baseDePrueba simula la colección periodo_seguimiento: guarda lo que calcula repartirUnidades
// y entrega los registros activos del más reciente al más antiguo, como los pide el mid.
type baseDePrueba struct {
	registros []map[string]interface{} // del más reciente al más antiguo
	creados   int
}

func (b *baseDePrueba) activos() []map[string]interface{} {
	activos := []map[string]interface{}{}
	for _, registro := range b.registros {
		if registro["activo"] == true {
			activos = append(activos, registro)
		}
	}
	return activos
}

// calcular lee los registros activos y reparte la selección como al habilitar, sin guardar.
func (b *baseDePrueba) calcular(inicio, fin string, unidades ...string) ([]map[string]interface{}, map[string]interface{}) {
	actualizar, crear := repartirUnidades(b.activos(), seleccionDePrueba(unidades...), fecha(inicio+"T00:00:00.000Z"), fecha(fin+"T00:00:00.000Z"))
	if crear != nil {
		b.creados++
		crear["_id"] = fmt.Sprintf("nuevo-%d", b.creados)
	}
	return actualizar, crear
}

// escribir guarda en el orden del mid, primero el registro nuevo y después los demás; cada
// registro que se escribe pasa a ser el más reciente.
func (b *baseDePrueba) escribir(actualizar []map[string]interface{}, crear map[string]interface{}) {
	escritos := actualizar
	if crear != nil {
		escritos = append([]map[string]interface{}{crear}, actualizar...)
	}
	for _, escrito := range escritos {
		registros := []map[string]interface{}{escrito}
		for _, registro := range b.registros {
			if registro["_id"] != escrito["_id"] {
				registros = append(registros, registro)
			}
		}
		b.registros = registros
	}
}

func (b *baseDePrueba) habilitar(inicio, fin string, unidades ...string) {
	b.escribir(b.calcular(inicio, fin, unidades...))
}

// Vale lo último que se guarda, sea para abrir o para cerrar.
func TestValeLoUltimoQueSeGuarda(t *testing.T) {
	hoy := enBogota("2026-11-24", 10, 0)
	base := &baseDePrueba{}
	exigir := func(paso string, quiere bool) {
		t.Helper()
		for _, dependencia := range []string{"7", "14"} {
			habilitada, err := UnidadHabilitada(base.activos(), dependencia, hoy)
			if err != nil || habilitada != quiere {
				t.Fatalf("%s: la unidad %s quedó habilitada=%v (err=%v) y se esperaba %v", paso, dependencia, habilitada, err, quiere)
			}
		}
	}

	// 1. Rectoría e Ingeniería se abren del 17 al 30 de noviembre.
	base.habilitar("2026-11-17", "2026-11-30", rectoria, ingenieria)
	exigir("abrir", true)

	// 2. Se cierran con un rango que ya pasó: salen del registro que las abría.
	base.habilitar("2026-11-10", "2026-11-14", rectoria, ingenieria)
	if activos := base.activos(); len(activos) != 1 || activos[0]["fecha_fin"] != fecha("2026-11-14T00:00:00.000Z").Format(time.RFC3339Nano) {
		t.Fatalf("solo debe quedar activo el registro que las cierra: %v", activos)
	}
	exigir("cerrar", false)

	// 3. Se vuelven a abrir.
	base.habilitar("2026-11-17", "2026-11-30", rectoria, ingenieria)
	exigir("volver a abrir", true)
}

// Si dos personas guardan al mismo tiempo, las dos leen los mismos registros y la unidad queda
// en dos; rige el que se guardó de último.
func TestDosGuardadosAlMismoTiempoRigeElUltimo(t *testing.T) {
	for _, abrirDeUltimo := range []bool{false, true} {
		t.Run(fmt.Sprintf("abrir de último: %v", abrirDeUltimo), func(t *testing.T) {
			base := &baseDePrueba{}
			abrirActualizar, abrirCrear := base.calcular("2026-11-17", "2026-11-30", ingenieria)
			cerrarActualizar, cerrarCrear := base.calcular("2026-11-10", "2026-11-14", ingenieria)
			if abrirDeUltimo {
				base.escribir(cerrarActualizar, cerrarCrear)
				base.escribir(abrirActualizar, abrirCrear)
			} else {
				base.escribir(abrirActualizar, abrirCrear)
				base.escribir(cerrarActualizar, cerrarCrear)
			}
			if activos := base.activos(); len(activos) != 2 {
				t.Fatalf("la unidad debe haber quedado en dos registros: %v", activos)
			}

			habilitada, err := UnidadHabilitada(base.activos(), "14", enBogota("2026-11-24", 10, 0))
			if err != nil || habilitada != abrirDeUltimo {
				t.Fatalf("habilitada=%v (err=%v) y se esperaba %v", habilitada, err, abrirDeUltimo)
			}
		})
	}
}
