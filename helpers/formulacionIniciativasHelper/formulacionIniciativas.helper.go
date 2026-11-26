package formulacionIniciativasHelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

// zonaBogota es la zona con la que se comparan por día las fechas de formulación.
var zonaBogota = time.FixedZone("COT", -5*60*60)

var formatoTipoSeguimientoID = regexp.MustCompile(`^[0-9a-fA-F]{24}$`)

// DefinirFechasIniciativasFormulacion habilita el rango de fechas para las unidades elegidas. Las
// fechas no son de una vigencia: cada unidad queda en un solo registro activo del tipo de
// seguimiento, y lo último que se guarda es lo que vale. Se les quita a los demás registros y se
// agrega al registro que ya tenga ese rango, o a uno nuevo. Devuelve los registros escritos.
func DefinirFechasIniciativasFormulacion(body map[string]interface{}) []interface{} {
	// 1. Validar la entrada
	tipoSeguimientoID := campoTexto(body, "tipo_seguimiento_id")
	if !formatoTipoSeguimientoID.MatchString(tipoSeguimientoID) {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Se requiere un tipo_seguimiento_id válido", "status": "400"})
	}
	inicio, errInicio := parseFechaHabilitacion(campoTexto(body, "fecha_inicio"))
	fin, errFin := parseFechaHabilitacion(campoTexto(body, "fecha_fin"))
	if errInicio != nil || errFin != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "fecha_inicio y fecha_fin deben ser fechas válidas", "status": "400"})
	}
	if diaDe(inicio).After(diaDe(fin)) {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "fecha_inicio no puede ser posterior a fecha_fin", "status": "400"})
	}
	seleccion, _ := ObtenerArrayUnidadesInteres(body)
	if len(seleccion) == 0 {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Debe elegir al menos una unidad", "status": "400"})
	}

	// 2. Traer los registros activos del tipo de seguimiento
	urlRegistros := beego.AppConfig.String("PlanesService") + "/periodo-seguimiento-iniciativas"
	consulta := "tipo_seguimiento_id:" + tipoSeguimientoID + ",activo:true"
	var resRegistros map[string]interface{}
	var registros []map[string]interface{}
	if err := helpers.SendJson(urlRegistros+"?query="+url.QueryEscape(consulta)+"&sortby=fecha_modificacion&order=desc&limit=0", "GET", &resRegistros, nil); err != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error al obtener los registros de fechas", "status": "400", "log": err})
	}
	exigirExito(resRegistros, "DefinirFechasIniciativasFormulacion", "Error al obtener los registros de fechas")
	helpers.LimpiezaRespuestaRefactor(resRegistros, &registros)

	// 3. Repartir las unidades
	actualizar, crear := repartirUnidades(registros, seleccion, inicio, fin)
	usuario := campoTexto(body, "usuario_modificacion")

	// 4. Guardar: primero el registro que queda con las unidades y después los demás, para que
	// una falla a mitad de camino no deje unidades sin fechas.
	escritos := []interface{}{}
	if crear != nil {
		crear["tipo_seguimiento_id"] = tipoSeguimientoID
		crear["usuario_modificacion"] = usuario
		escritos = append(escritos, guardarRegistro(urlRegistros, "POST", crear))
	}
	for _, registro := range actualizar {
		id, ok := registro["_id"].(string)
		if !ok || id == "" {
			panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Hay un registro de fechas sin _id", "status": "400"})
		}
		if usuario != "" {
			registro["usuario_modificacion"] = usuario
		}
		escritos = append(escritos, guardarRegistro(urlRegistros+"/"+url.PathEscape(id), "PUT", registro))
	}
	return escritos
}

// guardarRegistro envía el registro al crud y devuelve el registro guardado.
func guardarRegistro(urlRegistro, metodo string, registro map[string]interface{}) interface{} {
	var res map[string]interface{}
	if err := helpers.SendJson(urlRegistro, metodo, &res, registro); err != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error al guardar el registro de fechas", "status": "400", "log": err})
	}
	exigirExito(res, "DefinirFechasIniciativasFormulacion", "Error al guardar el registro de fechas")
	return res["Data"]
}

// exigirExito responde 400 con el mensaje del crud si este respondió Success: false. Hace falta
// porque la mayoría de los recursos del crud responden HTTP 200 aunque fallen.
func exigirExito(respuesta map[string]interface{}, funcion, mensaje string) {
	if exito, _ := respuesta["Success"].(bool); !exito {
		panic(map[string]interface{}{"funcion": funcion, "err": fmt.Sprint(mensaje, ": ", respuesta["Message"]), "status": "400"})
	}
}

// repartirUnidades calcula qué escribir para que las unidades de la selección queden solo en el
// registro con el rango [inicio, fin]. El destino es el primer registro con las mismas fechas,
// comparando por día. A los demás se les quitan las unidades de la selección y, si quedan sin
// unidades, se inactivan. Las unidades se comparan solo por Id.
//
// Devuelve los registros que cambiaron, con el destino de primero, y, si no hay destino, el
// registro a crear. Ese registro trae las fechas, las unidades, planes_interes y activo; quien
// llama completa tipo_seguimiento_id y usuario_modificacion. No modifica los registros recibidos.
func repartirUnidades(registros []map[string]interface{}, seleccion []map[string]interface{}, inicio, fin time.Time) (actualizar []map[string]interface{}, crear map[string]interface{}) {
	seleccion = unidadesSinRepetir(seleccion)
	idsSeleccion := make(map[int]bool, len(seleccion))
	for _, unidad := range seleccion {
		idsSeleccion[idUnidad(unidad)] = true
	}

	indiceDestino := -1
	for i, registro := range registros {
		if mismoRango(registro, inicio, fin) {
			indiceDestino = i
			break
		}
	}

	if indiceDestino >= 0 {
		destino := copiarRegistro(registros[indiceDestino])
		actuales, _ := ObtenerArrayUnidadesInteres(destino)
		combinadas := make([]map[string]interface{}, 0, len(actuales)+len(seleccion))
		combinadas = append(combinadas, actuales...)
		combinadas = unidadesSinRepetir(append(combinadas, seleccion...))
		if !mismosIds(actuales, combinadas) {
			destino["unidades_interes"] = textoUnidades(combinadas)
			actualizar = append(actualizar, destino)
		}
	} else {
		crear = map[string]interface{}{
			"fecha_inicio":     inicio.Format(time.RFC3339Nano),
			"fecha_fin":        fin.Format(time.RFC3339Nano),
			"unidades_interes": textoUnidades(seleccion),
			"planes_interes":   "[]",
			"activo":           true,
		}
	}

	for i, registro := range registros {
		if i == indiceDestino {
			continue
		}
		actuales, _ := ObtenerArrayUnidadesInteres(registro)
		restantes := make([]map[string]interface{}, 0, len(actuales))
		for _, unidad := range actuales {
			if !idsSeleccion[idUnidad(unidad)] {
				restantes = append(restantes, unidad)
			}
		}
		if len(restantes) == len(actuales) {
			continue
		}
		copia := copiarRegistro(registro)
		copia["unidades_interes"] = textoUnidades(restantes)
		if len(restantes) == 0 {
			copia["activo"] = false
		}
		actualizar = append(actualizar, copia)
	}
	return actualizar, crear
}

// mismoRango indica si el registro tiene el mismo fecha_inicio y fecha_fin, comparando por día.
func mismoRango(registro map[string]interface{}, inicio, fin time.Time) bool {
	inicioRegistro, errInicio := parseFecha(fmt.Sprint(registro["fecha_inicio"]), time.UTC)
	finRegistro, errFin := parseFecha(fmt.Sprint(registro["fecha_fin"]), time.UTC)
	if errInicio != nil || errFin != nil {
		return false
	}
	return inicioRegistro.Equal(diaDe(inicio)) && finRegistro.Equal(diaDe(fin))
}

// diaDe devuelve la medianoche UTC del día de t en UTC, que es el día que se guarda en el crud.
func diaDe(t time.Time) time.Time {
	utc := t.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

// parseFechaHabilitacion lee una fecha RFC 3339 (como la envía el cliente) o AAAA-MM-DD.
func parseFechaHabilitacion(valor string) (time.Time, error) {
	if fecha, err := time.Parse(time.RFC3339, valor); err == nil {
		return fecha, nil
	}
	return time.Parse("2006-01-02", valor)
}

// unidadesSinRepetir devuelve las unidades sin Ids repetidos, en el orden en que aparecen.
func unidadesSinRepetir(unidades []map[string]interface{}) []map[string]interface{} {
	vistos := make(map[int]bool, len(unidades))
	resultado := make([]map[string]interface{}, 0, len(unidades))
	for _, unidad := range unidades {
		id := idUnidad(unidad)
		if vistos[id] {
			continue
		}
		vistos[id] = true
		resultado = append(resultado, unidad)
	}
	return resultado
}

func mismosIds(a, b []map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if idUnidad(a[i]) != idUnidad(b[i]) {
			return false
		}
	}
	return true
}

// idUnidad devuelve el Id de una unidad leída con ObtenerArrayUnidadesInteres.
func idUnidad(unidad map[string]interface{}) int {
	id, _ := unidad["Id"].(int)
	return id
}

// textoUnidades convierte las unidades al texto JSON de unidades_interes, solo con Id y Nombre.
func textoUnidades(unidades []map[string]interface{}) string {
	limpias := make([]map[string]interface{}, 0, len(unidades))
	for _, unidad := range unidades {
		limpias = append(limpias, map[string]interface{}{"Id": unidad["Id"], "Nombre": unidad["Nombre"]})
	}
	texto, err := json.Marshal(limpias)
	if err != nil {
		panic(map[string]interface{}{"funcion": "DefinirFechasIniciativasFormulacion", "err": "Error al convertir las unidades a JSON", "status": "400", "log": err})
	}
	return string(texto)
}

func copiarRegistro(registro map[string]interface{}) map[string]interface{} {
	copia := make(map[string]interface{}, len(registro))
	for clave, valor := range registro {
		copia[clave] = valor
	}
	return copia
}

// campoTexto devuelve el campo del body como texto sin espacios, o "" si no viene.
func campoTexto(body map[string]interface{}, campo string) string {
	valor, ok := body[campo]
	if !ok || valor == nil {
		return ""
	}
	return strings.TrimSpace(textoID(valor))
}

// textoID convierte un Id a texto sin notación científica: 14, 14.0 y "14" dan "14".
func textoID(valor interface{}) string {
	switch v := valor.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case json.Number:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}

// UnidadHabilitada indica si el periodo que rige a la dependencia cubre el día de hoy en Bogotá,
// inclusivo en los dos extremos. Los periodos llegan del más reciente (fecha_modificacion) al
// más antiguo, y rige el primero que incluye la dependencia: si una
// unidad quedó en dos periodos, por dos guardados al mismo tiempo, vale lo último que se guardó.
// El Id se compara como texto, así que 14 y "14" son iguales. Responde error si un periodo tiene
// unidades_interes o fechas inválidas.
func UnidadHabilitada(periodos []map[string]interface{}, dependenciaID string, hoy time.Time) (bool, error) {
	ahora := hoy.In(zonaBogota)
	dia := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, zonaBogota)

	for _, periodo := range periodos {
		incluye, err := periodoIncluyeUnidad(periodo, dependenciaID)
		if err != nil {
			return false, err
		}
		if !incluye {
			continue
		}

		inicio, err := parseFecha(fmt.Sprint(periodo["fecha_inicio"]), zonaBogota)
		if err != nil {
			return false, fmt.Errorf("fecha_inicio inválida en el periodo %v: %w", periodo["_id"], err)
		}
		fin, err := parseFecha(fmt.Sprint(periodo["fecha_fin"]), zonaBogota)
		if err != nil {
			return false, fmt.Errorf("fecha_fin inválida en el periodo %v: %w", periodo["_id"], err)
		}
		return !dia.Before(inicio) && !dia.After(fin), nil
	}
	return false, nil
}

// periodoIncluyeUnidad indica si unidades_interes del periodo tiene la dependencia. Compara el
// Id como texto, así que 14 y "14" son iguales.
func periodoIncluyeUnidad(periodo map[string]interface{}, dependenciaID string) (bool, error) {
	texto, ok := periodo["unidades_interes"].(string)
	if !ok {
		return false, fmt.Errorf("el periodo %v no tiene unidades_interes como texto", periodo["_id"])
	}
	var unidades []map[string]interface{}
	decodificador := json.NewDecoder(bytes.NewReader([]byte(texto)))
	decodificador.UseNumber()
	if err := decodificador.Decode(&unidades); err != nil {
		return false, fmt.Errorf("unidades_interes inválido en el periodo %v: %w", periodo["_id"], err)
	}
	for _, unidad := range unidades {
		if textoID(unidad["Id"]) == dependenciaID {
			return true, nil
		}
	}
	return false, nil
}

// Toma solo la parte YYYY-MM-DD para comparar por día y evitar desfases de zona horaria
func parseFecha(valor string, loc *time.Location) (time.Time, error) {
	if len(valor) < 10 {
		return time.Time{}, fmt.Errorf("fecha inválida: %q", valor)
	}
	return time.ParseInLocation("2006-01-02", valor[:10], loc)
}

func ObtenerArrayUnidadesInteres(body map[string]interface{}) ([]map[string]interface{}, error) {
	// Acceder a la lista de unidades de interés
	unidadesInteresString, ok := body["unidades_interes"].(string)
	if !ok {
		panic(map[string]interface{}{"funcion": "ObtenerArrayUnidadesInteres", "err": "No se pudo obtener la lista de unidades de interés como cadena JSON", "status": "400"})
	}

	// Decodificar la cadena JSON a una lista de interfaces
	var unidadesInteres []interface{}
	err := json.Unmarshal([]byte(unidadesInteresString), &unidadesInteres)
	if err != nil {
		panic(map[string]interface{}{"funcion": "ObtenerArrayUnidadesInteres", "err": "Error al decodificar la lista de unidades de interés", "status": "400"})
	}

	// Almacenar las unidades de interés en un array
	var unidades []map[string]interface{}
	for _, unidad := range unidadesInteres {
		if unidadMap, ok := unidad.(map[string]interface{}); ok {
			// Verificar si la clave "Id" está presente y tiene el tipo correcto
			if id, idOk := unidadMap["Id"].(float64); idOk {
				// Convertir el valor de "Id" a entero
				idInt := int(id)
				unidades = append(unidades, map[string]interface{}{
					"Id":     idInt,
					"Nombre": unidadMap["Nombre"],
					// Otras propiedades de interés si es necesario
				})
			} else {
				panic(map[string]interface{}{"funcion": "ObtenerArrayUnidadesInteres", "err": "El valor de Id no es un número", "status": "400"})
			}
		}
	}
	return unidades, nil
}

// ConsultarFechasIniciativasFormulacion devuelve las fechas de formulación de una unidad
// (dependenciaId): las del registro que la rige, el mismo que usa UnidadHabilitada.
func ConsultarFechasIniciativasFormulacion(body map[string]interface{}) map[string]interface{} {
	var resTipoSeguimiento map[string]interface{}
	var resFechas map[string]interface{}
	var tipoSeguimiento []map[string]interface{}
	var fechas []map[string]interface{}
	respuesta := map[string]interface{}{"encontrado": false}

	if body["dependenciaId"] == nil {
		panic(map[string]interface{}{"funcion": "ConsultarFechasIniciativasFormulacion", "err": "Se requiere el campo dependenciaId", "status": "400"})
	}
	dependenciaId := textoID(body["dependenciaId"])

	// 1. Obtener el tipo de seguimiento de formulación de iniciativas
	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/tipo-seguimiento/?query=codigo_abreviacion:FIPI_SP", &resTipoSeguimiento); err != nil {
		panic(map[string]interface{}{"funcion": "ConsultarFechasIniciativasFormulacion", "err": "Error al obtener el tipo de seguimiento", "status": "400", "log": err})
	}
	exigirExito(resTipoSeguimiento, "ConsultarFechasIniciativasFormulacion", "Error al obtener el tipo de seguimiento")
	helpers.LimpiezaRespuestaRefactor(resTipoSeguimiento, &tipoSeguimiento)
	if len(tipoSeguimiento) == 0 {
		panic(map[string]interface{}{"funcion": "ConsultarFechasIniciativasFormulacion", "err": "No se encontró el tipo de seguimiento", "status": "404"})
	}

	// 2. Traer los registros activos, el más reciente primero
	consulta := "tipo_seguimiento_id:" + tipoSeguimiento[0]["_id"].(string) + ",activo:true"
	if err := request.GetJson(beego.AppConfig.String("PlanesService")+"/periodo-seguimiento-iniciativas?query="+consulta+"&sortby=fecha_modificacion&order=desc&limit=0", &resFechas); err != nil {
		panic(map[string]interface{}{"funcion": "ConsultarFechasIniciativasFormulacion", "err": "Error al obtener las fechas", "status": "400", "log": err})
	}
	exigirExito(resFechas, "ConsultarFechasIniciativasFormulacion", "Error al obtener las fechas")
	helpers.LimpiezaRespuestaRefactor(resFechas, &fechas)

	// 3. Buscar el primer registro (el más reciente) que incluya la unidad
	for _, registro := range fechas {
		if registroIncluyeUnidad(registro, dependenciaId) {
			respuesta["encontrado"] = true
			respuesta["fecha_inicio"] = registro["fecha_inicio"]
			respuesta["fecha_fin"] = registro["fecha_fin"]
			break
		}
	}
	return respuesta
}

// registroIncluyeUnidad indica si el registro tiene la unidad en unidades_interes.
// Compara el Id como texto, por lo que 14 y "14" se consideran iguales.
func registroIncluyeUnidad(registro map[string]interface{}, dependenciaId string) bool {
	incluye, err := periodoIncluyeUnidad(registro, dependenciaId)
	return err == nil && incluye
}
