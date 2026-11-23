package formulacionIniciativasHelper

import (
	"encoding/json"
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
func periodoDePrueba(id, periodoID, inicio, fin string, unidades ...string) map[string]interface{} {
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
		"periodo_id":          periodoID,
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

func TestUnidadHabilitada(t *testing.T) {
	casos := []struct {
		nombre           string
		periodos         []map[string]interface{}
		dependenciaID    string
		periodoID        string
		hoy              time.Time
		quiereEncontrada bool
		quiereHabilitada bool
	}{
		{
			nombre: "la unidad está en el segundo registro con fechas vigentes (F1)",
			periodos: []map[string]interface{}{
				periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", rectoria),
				periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", ingenieria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: true, quiereHabilitada: true,
		},
		{
			nombre: "el primer registro está vencido y el segundo vigente (F2)",
			periodos: []map[string]interface{}{
				periodoDePrueba("viejo", "11", "2025-11-01", "2025-11-30", rectoria, ingenieria),
				periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", ingenieria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: true, quiereHabilitada: true,
		},
		{
			nombre: "un registro vencido que no tiene la unidad no la bloquea (F2)",
			periodos: []map[string]interface{}{
				periodoDePrueba("viejo", "11", "2025-11-01", "2025-11-30", rectoria),
				periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", ingenieria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: true, quiereHabilitada: true,
		},
		{
			nombre: "la unidad solo está en un registro vencido",
			periodos: []map[string]interface{}{
				periodoDePrueba("A", "12", "2026-11-01", "2026-11-15", ingenieria),
				periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", rectoria),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: true, quiereHabilitada: false,
		},
		{
			nombre: "la unidad no está en ningún registro",
			periodos: []map[string]interface{}{
				periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", rectoria, ciencias),
			},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: false, quiereHabilitada: false,
		},
		{
			nombre:        "hoy es el día de inicio",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", ingenieria)},
			dependenciaID: "14", hoy: enBogota("2026-11-10", 0, 0),
			quiereEncontrada: true, quiereHabilitada: true,
		},
		{
			nombre:        "hoy es el día de fin, de noche en Bogotá (ya es el día siguiente en UTC)",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", ingenieria)},
			dependenciaID: "14", hoy: enBogota("2026-11-30", 23, 30),
			quiereEncontrada: true, quiereHabilitada: true,
		},
		{
			nombre:        "el día anterior al inicio",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", ingenieria)},
			dependenciaID: "14", hoy: enBogota("2026-11-09", 23, 59),
			quiereEncontrada: true, quiereHabilitada: false,
		},
		{
			nombre:        "el día siguiente al fin",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", ingenieria)},
			dependenciaID: "14", hoy: enBogota("2026-12-01", 0, 0),
			quiereEncontrada: true, quiereHabilitada: false,
		},
		{
			nombre: "con periodoID solo cuenta esa vigencia",
			periodos: []map[string]interface{}{
				periodoDePrueba("otra", "11", "2026-11-10", "2026-11-30", ingenieria),
				periodoDePrueba("esta", "12", "2026-11-10", "2026-11-30", rectoria),
			},
			dependenciaID: "14", periodoID: "12", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: false, quiereHabilitada: false,
		},
		{
			nombre: "con periodoID encuentra la unidad en esa vigencia",
			periodos: []map[string]interface{}{
				periodoDePrueba("otra", "11", "2026-11-01", "2026-11-05", rectoria),
				periodoDePrueba("esta", "12", "2026-11-10", "2026-11-30", ingenieria),
			},
			dependenciaID: "14", periodoID: "12", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: true, quiereHabilitada: true,
		},
		{
			nombre:        "la unidad 209, con caracteres especiales en el nombre",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", rectoria, i3)},
			dependenciaID: "209", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: true, quiereHabilitada: true,
		},
		{
			nombre:        "un Id guardado como texto se compara igual",
			periodos:      []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", `{"Id":"14","Nombre":"FACULTAD DE INGENIERIA"}`)},
			dependenciaID: "14", hoy: enBogota("2026-11-25", 10, 0),
			quiereEncontrada: true, quiereHabilitada: true,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			encontrada, habilitada, err := unidadHabilitada(caso.periodos, caso.dependenciaID, caso.periodoID, caso.hoy)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if encontrada != caso.quiereEncontrada || habilitada != caso.quiereHabilitada {
				t.Fatalf("se esperaba encontrada=%v habilitada=%v y llegó encontrada=%v habilitada=%v",
					caso.quiereEncontrada, caso.quiereHabilitada, encontrada, habilitada)
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
			if _, _, err := unidadHabilitada([]map[string]interface{}{caso.periodo}, "14", "", enBogota("2026-11-25", 10, 0)); err == nil {
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
			nombre:    "ejemplo A y B: Ingeniería cambia de rango",
			registros: []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", rectoria, ingenieria)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7}, true}},
			quiereCrear:    &esperado{"", []int{14}, true},
		},
		{
			nombre: "la selección está repartida en dos registros y se quita de los dos (F3)",
			registros: []map[string]interface{}{
				periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", rectoria, ingenieria),
				periodoDePrueba("C", "12", "2026-12-01", "2026-12-15", ciencias, i3),
			},
			seleccion: seleccionDePrueba(ingenieria, i3),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7}, true}, {"C", []int{21}, true}},
			quiereCrear:    &esperado{"", []int{14, 209}, true},
		},
		{
			nombre:    "el registro que queda vacío se inactiva",
			registros: []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", ingenieria)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{}, false}},
			quiereCrear:    &esperado{"", []int{14}, true},
		},
		{
			nombre: "ya existe un registro con las mismas fechas: se agregan ahí sin repetir",
			registros: []map[string]interface{}{
				periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", rectoria, ingenieria),
				periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", ciencias),
			},
			seleccion: seleccionDePrueba(ingenieria, ciencias),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"B", []int{21, 14}, true}, {"A", []int{7}, true}},
		},
		{
			nombre:    "el destino se compara por día aunque el cliente envíe otra hora",
			registros: []map[string]interface{}{periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", ciencias)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T05:00:00.000Z", fin: "2026-12-10T05:00:00.000Z",
			quiereActuales: []esperado{{"B", []int{21, 14}, true}},
		},
		{
			nombre:    "si la selección ya está solo en su destino no se escribe nada",
			registros: []map[string]interface{}{periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", ingenieria, ciencias)},
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
			registros: []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", i3, rectoria)},
			seleccion: seleccionDePrueba(i3),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7}, true}},
			quiereCrear:    &esperado{"", []int{209}, true},
		},
		{
			nombre:    "la unidad 209 no se duplica en su destino (F4)",
			registros: []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", i3, rectoria)},
			seleccion: seleccionDePrueba(i3),
			inicio:    "2026-11-10T00:00:00.000Z", fin: "2026-11-30T00:00:00.000Z",
		},
		{
			nombre:    "Ids repetidos en un registro se quitan todos (F6)",
			registros: []map[string]interface{}{periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", ingenieria, ingenieria, rectoria)},
			seleccion: seleccionDePrueba(ingenieria),
			inicio:    "2026-11-20T00:00:00.000Z", fin: "2026-12-10T00:00:00.000Z",
			quiereActuales: []esperado{{"A", []int{7}, true}},
			quiereCrear:    &esperado{"", []int{14}, true},
		},
		{
			nombre:    "Ids repetidos en el destino quedan una sola vez (F6)",
			registros: []map[string]interface{}{periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", ingenieria, ingenieria)},
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
		})
	}
}

func TestRepartirUnidadesNoModificaLosRegistrosRecibidos(t *testing.T) {
	registros := []map[string]interface{}{
		periodoDePrueba("A", "12", "2026-11-10", "2026-11-30", rectoria, ingenieria),
		periodoDePrueba("B", "12", "2026-11-20", "2026-12-10", ciencias),
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
	crearA["periodo_id"] = "12"

	// 2. Ingeniería pasa al 20 de noviembre-10 de diciembre: A se queda con Rectoría y se crea B.
	actualizar, crearB := repartirUnidades([]map[string]interface{}{crearA}, seleccionDePrueba(ingenieria), fecha("2026-11-20T00:00:00.000Z"), fecha("2026-12-10T00:00:00.000Z"))
	if len(actualizar) != 1 || crearB == nil {
		t.Fatalf("se esperaba actualizar A y crear B: %v %v", actualizar, crearB)
	}
	crearB["_id"] = "B"
	crearB["periodo_id"] = "12"
	periodos := []map[string]interface{}{actualizar[0], crearB}

	// 3. El 25 de noviembre las dos unidades están habilitadas, cada una por su registro.
	for _, dependencia := range []string{"7", "14"} {
		encontrada, habilitada, err := unidadHabilitada(periodos, dependencia, "12", enBogota("2026-11-25", 10, 0))
		if err != nil || !encontrada || !habilitada {
			t.Fatalf("la unidad %s debería estar habilitada: encontrada=%v habilitada=%v err=%v", dependencia, encontrada, habilitada, err)
		}
	}
	// El 5 de diciembre solo Ingeniería sigue habilitada.
	if _, habilitada, _ := unidadHabilitada(periodos, "7", "12", enBogota("2026-12-05", 10, 0)); habilitada {
		t.Fatal("Rectoría no debería estar habilitada el 5 de diciembre")
	}
	if _, habilitada, _ := unidadHabilitada(periodos, "14", "12", enBogota("2026-12-05", 10, 0)); !habilitada {
		t.Fatal("Ingeniería debería estar habilitada el 5 de diciembre")
	}
}
