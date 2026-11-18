// Package iniciativaInversionHelper implementa la formulación de las iniciativas
// de inversión a partir de plantillas armadas en Construcción.
//
// Una iniciativa es un plan de tipo IPI_SP (formato: false) que apunta a su
// plantilla (formato_id). La plantilla define los submódulos y sus campos en
// subgrupo/subgrupo_detalle; los datos de cada submódulo se guardan en planes_crud
// como un JSON legible, en la colección submodulo_iniciativa (un documento por
// iniciativa y submódulo).
//
// Organización:
//   - modelo.go         tipos de la definición, los submódulos y las respuestas
//   - crud.go           cliente HTTP de planes_crud (revisa códigos de respuesta)
//   - definicion.go     lectura del árbol de la plantilla, con caché
//   - catalogos.go      catálogos (colección general elemento_catalogo), con caché
//   - sedes.go          catálogo de sedes, que sale de Oikos
//   - normalizacion.go  limpieza de datos y cálculo de campos calculados
//   - validacion.go     reglas por tipo de campo, espejos y referencias
//   - estados.go        tipo de plan, estados de la iniciativa y de los submódulos, y sus transiciones
//   - servicio.go       casos de uso que consume el controlador
//   - observaciones.go  historial de observaciones de cada submódulo
//   - unidades.go       unidades del usuario (Terceros y Oikos)
//   - errores.go        errores con el mensaje para el usuario y el detalle para el log
//   - valores.go        lectura segura de valores JSON y armado de rutas
//
// Este paquete no modifica ni depende de los flujos de Formulación.
package iniciativaInversionHelper
