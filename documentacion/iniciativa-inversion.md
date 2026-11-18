# Iniciativas de inversión (`/v1/iniciativa-inversion`)

Formulación de las iniciativas para proyectos de inversión a partir de plantillas armadas en Construcción. Por ahora cubre el módulo Identificación.

**Código:**
- `controllers/iniciativa_inversion.go`: solo traduce HTTP.
- `helpers/iniciativaInversionHelper/`: la lógica.

**Aislamiento:** no modifica ni usa los flujos de Formulación ni los de MACRO. Lo único compartido es el registro de la ruta en `routers/router.go` y su `commentsRouter` generado.

**Tipo de plan propio:** las plantillas y las iniciativas usan `IPI_SP` ("Iniciativa de inversión"). Los proyectos MACRO usan `PRI_SP`, así que no se mezclan: un plan de otro tipo no aparece en el listado y su consulta responde 404.

## Modelo

| Qué | Dónde (planes_crud) |
|---|---|
| Plantilla | `plan` con `formato: true` y `tipo_plan_id` = tipo `IPI_SP`; sus nodos en `subgrupo` / `subgrupo_detalle` (se carga con `database local/construccion-proyecto-inversion/cargar-construccion.js`). **Primer nivel = módulos** (`dato.modulo: true`: Identificación, Preparación y Anexos; el mid no depende de cuántos haya). **Segundo nivel = submódulos (hojas)**, con clave única en toda la plantilla |
| Iniciativa | `plan` de tipo `IPI_SP` con `formato: false`, `formato_id` = plantilla, `dependencia_id`, `estado_plan_id` (estado general) y, si se enviaron, `vigencia_aplica` (JSON de vigencias) |
| Datos de cada submódulo | `submodulo_iniciativa`: un documento por iniciativa y submódulo, con `dato` (JSON legible armado con las claves de la plantilla), `observaciones`, `estado_submodulo_id` y `version` |

**La iniciativa no copia el árbol de la plantilla:**
- La definición se lee de la plantilla y queda en una caché de 10 minutos.
- Una plantilla con iniciativas no se edita. Si cambia el instrumento, se carga una plantilla nueva.

**Cada nodo de la plantilla tiene en `dato`:**
- Una `clave` técnica, única entre hermanos.
- Las propiedades de su tipo (ver «Reglas» más abajo; `database local/construccion-proyecto-inversion/cargar-construccion.js` las valida al cargar la plantilla).

Ejemplo de `dato` del submódulo `problematica` (en el API del mid viaja como `datos`):

```json
{
  "problema_central": "…",
  "causas_directas": [
    { "id": "cd-1", "descripcion": "…", "linea_base": "…",
      "causas_indirectas": [ { "id": "ci-1", "descripcion": "…", "linea_base": "…" } ] }
  ]
}
```

- **Filas:** cada fila de un grupo repetible lleva un `id` estable.
- **Espejos:** las filas de un grupo espejo llevan además `origen`, que es el `id` de la fila de la que salen.
- **Referencias:** guardan el `id` de la fila referenciada.

## Endpoints

Todos responden `{ Success, Status, Message, Data }` con el **código HTTP real**.

`Message` es lo que ve el usuario: el cliente lo muestra tal cual. Por eso no lleva nombres de servicios, rutas, colecciones, códigos ni scripts. La causa técnica (`Detalle` de `ErrorNegocio`) solo va al log del mid. Si faltan los datos base del módulo (tipo de plan o estados) o la plantilla está mal armada, responde 500 con un mensaje general, y el log dice qué falta y con qué script se carga. Los mensajes técnicos del crud tampoco llegan al usuario; solo pasa el del 409 (el submódulo cambió desde que se consultó).

| Método y ruta | Qué hace | Errores |
|---|---|---|
| `GET /plantillas` | Plantillas activas con las que se pueden crear iniciativas, de la más reciente a la más antigua | 502 si falla el crud |
| `GET /iniciativas?dependencias=1,2` | Iniciativas activas de esas dependencias (sin el filtro, todas), de la más reciente a la más antigua. Cada fila trae `numero` (consecutivo por orden de creación, contado sobre todas las iniciativas, activas o no), `estado` (general), `permite_editar` y `responsable` (nombre en Terceros de quien la creó, consultado una vez por documento; vacío si la iniciativa no tiene responsable o no se obtiene el nombre) | |
| `POST /iniciativas` | Crea la iniciativa y todos sus submódulos en "Sin formular" (la iniciativa también). Cuerpo: `{ dependencia_id, responsable, plantilla_id?, nombre?, vigencias? }`, con `responsable` = documento de quien la crea (el mid no lo exige). Sin `plantilla_id` usa la plantilla más reciente; sin `nombre`, "Iniciativa de inversión sin nombre" | 400 si falta `dependencia_id`; 409 si no hay plantillas |
| `GET /iniciativas/:id` | Vista general: estado general y `modulos: [{ clave, nombre, submodulos: [{ clave, nombre, estado (con color), observaciones }] }]`, más `permite_editar`, `permite_revisar` y `acciones_posibles` | 404 |
| `GET /iniciativas/:id/definicion` | Construcción de la plantilla: `modulos`, y en cada uno sus submódulos y campos | 404 |
| `GET /iniciativas/:id/submodulos/:clave` | Datos, errores vigentes, historial de observaciones y `relacionadas` (datos de los submódulos de los que depende para espejos y referencias) | 404 |
| `PUT /iniciativas/:id/submodulos/:clave` | Guarda `{ version, datos }`, aunque tenga errores (borrador). Responde el submódulo con sus `errores`; queda **Formulado** si no hay errores, **Formulando** si tiene información pero errores y **Sin formular** si está vacío. Con la primera información, la iniciativa pasa de Sin formular a **Formulación iniciada** | 409 si la iniciativa no se puede editar o si `version` no es la vigente |
| `POST /iniciativas/:id/submodulos/:clave/observaciones` | Publica una observación en el historial: `{ texto, rol, usuario }`, con `usuario` = documento de quien comenta. El mid agrega `usuario_nombre`, su nombre en Terceros (`TercerosService`); si no lo obtiene, la publica sin nombre. No cambia la `version` del submódulo. Con la iniciativa en revisión, el submódulo queda **Con observaciones** | 400 si falta un dato o `usuario` no es un documento (letras y números, máx. 20); 409 si la iniciativa está avalada, o si está en revisión y el submódulo no está formulado |
| `PUT /iniciativas/:id/submodulos/:clave/revision` | Planeación aprueba el submódulo: `{ version }` → **Aprobado**. El historial se conserva | 400 si se envían observaciones (van por el `POST`); 409 si la iniciativa no está en revisión o el submódulo no está formulado |
| `PUT /iniciativas/:id/acciones/:accion` | `enviar`, `devolver` o `avalar` | 400 si la acción no existe, 409 si no se permite |
| `GET /unidades/:documento` | Unidades del usuario (`[{ Id, Nombre }]`): tercero por documento y vinculaciones activas en Terceros, con los mismos cargos que Formulación (`/v1/formulacion/vinculacion_tercero`: 319, 312, 320 y 300), y nombre en Oikos. Lista vacía si no tiene. El documento solo admite letras y números (máx. 20). Usa `TercerosService` y `OikosService` de `conf/app.conf` | 400 si el documento no es válido; 502 si falla Terceros u Oikos |
| `GET /catalogos/:nombre?padre=` | Cualquier catálogo de la colección general `elemento_catalogo` del crud: `procesos`, `sectores_pnd`, `programas_pnd` (padre = sector), `productos` (padre = programa), `municipios` (padre = departamento), `tipos_riesgo`… | 404 si no existe o no está cargado |

## Reglas

### Al guardar un submódulo

1. **Normalizar:**
   - Se quitan las claves que no están en la plantilla y se recortan los textos.
   - Se recalculan los campos `calculado` y se ignora el valor que envíe el cliente. Excepción: si el campo es `editable` y el cliente envía un texto no vacío, se guarda ese texto (con la misma transformación, p. ej. mayúsculas). Así funciona el nombre de la iniciativa: el sistema lo propone y el usuario puede cambiarlo.
   - Un `calculado` con `operacion` es numérico: `"operacion": "resta", "operandos": ["oferta", "demanda"]` guarda el primer campo hermano menos los demás (el déficit de Necesidades); si falta un operando, queda sin valor. `"operacion": "suma"` guarda la suma, con los vacíos en 0 (el costo total de cada actividad de Cadena de valor, de `periodo_0` a `periodo_4`); si todos están vacíos, queda sin valor. `"operacion": "multiplicacion"` guarda el producto (el valor total de cada año de Ingresos y beneficios: `cantidad` × `valor_unitario`); si falta un operando, queda sin valor.
   - Un `calculado` con `desde: "<submodulo>.<campo>"` solo muestra ese valor de otro submódulo (el objetivo general en Cadena de valor): no se guarda ni se valida, y el submódulo de origen va en `relacionadas`. Con `desde: "origen.<campo>"`, en un campo de un espejo, muestra ese campo de la fila de origen (en el Cronograma, el costo de cada periodo de la actividad en Cadena de valor y su costo total); tampoco se guarda. Con `desde: "<ruta de un grupo>.<campo>"`, en un espejo cuyo `espejo_de` pasa por ese grupo, muestra el campo de la fila del grupo que contiene la de origen (en el Presupuesto inicial, `cadena_valor.productos_por_objetivo.productos.producto` es el producto de cada actividad). Con `catalogo`, el cliente muestra el nombre del elemento en vez del código.
   - Un espejo en el que todos los campos se traen (`desde`) es un resumen de solo lectura: solo guarda `id` y `origen` de cada fila, y queda Formulado al guardarlo con una fila por elemento de origen (el Presupuesto inicial de la alternativa seleccionada). Con `totales: { "columnas": true }`, el cliente agrega una fila con la suma de cada columna numérica.
2. **Validar** según el tipo de campo:
   - Obligatorios y longitud máxima.
   - Mínimo y máximo, y decimales.
   - Códigos de catálogo y su `depende_de`. Toda lista de opciones es un catálogo: la plantilla no lleva opciones escritas.
   - Mínimo y máximo de filas, e `id` únicos.
   - Espejos: exactamente una fila por cada fila de origen. Con `espejo_de: "catalogo.<nombre>"` las filas de origen son los elementos del catálogo (`{ id: codigo, nombre }`, con `mostrar: "nombre"`); el mid se los manda al cliente en `relacionadas.catalogo.<nombre>` (los resultados esperados de Objetivos). Con `filas_opcionales: true`, la fila vacía no se exige ni se valida (no aplica), la que tiene algún dato se valida completa y `min` es cuántas filas se llenan como mínimo. `filas_obligatorias` (lista de ids de origen, p. ej. códigos del catálogo) son las filas que sí se exigen: en Resultados esperados, estudiantes, docentes y administrativos beneficiados. Un campo de un espejo puede tener `solo_si: { campo: "origen", valor: "<id o código>" }`: solo aplica a esa fila (el «¿Cuál?» de la fila «Otro»); el cliente lo escribe en la celda del elemento de origen.
   - Referencias: solo filas válidas en su alcance. Con `marcada_en` (clave de una casilla única del grupo de origen), la única opción es la fila marcada: el cliente la llena sola y, si después se marca otra, el submódulo pide guardarse de nuevo (Análisis técnico con la alternativa que pasa a Preparación).
   - Matrices completas.
   - Catálogos con `multiple`: lista de códigos del catálogo, sin repetir; el cliente los muestra como casillas (las sedes y los factores de Localización). El catálogo `sedes` no está en `elemento_catalogo`: `RepositorioConSedes` (`sedes.go`) lo lee de Oikos (`OikosService`), de los espacios físicos de tipo SEDE (38) activos, con el `Id` como código (el `Codigo` de Oikos se repite) y «Otra» al final.
   - `solo_si: { campo, valor }`: el campo solo se exige y se guarda si el campo hermano vale `valor` o, si es una lista, lo incluye («¿Cuál otra sede?» solo si se marcó «Otra»). Con `"distinto": true`, al revés: aplica si no lo vale. Si no aplica, el normalizador lo descarta. En Trazadores presupuestales, Discriminado y Número de beneficiarios aplican a todos los trazadores menos a Grupos Étnicos, y la lista de comunidades (un espejo de `catalogo.grupos_etnicos_trazador`, con `filas_opcionales` y `min: 1`) solo a Grupos Étnicos.
   - Los elementos de un catálogo pueden traer `descripcion` (campo opcional de `elemento_catalogo`). En un espejo de catálogo, la fila de origen la incluye y el cliente la muestra como ayuda bajo el nombre (la definición de cada trazador presupuestal).
   - `obligatorio_si: "origen.<campo>"` (con `required: "false"`), en un campo de un espejo: el campo se exige solo si ese campo de la fila de origen tiene valor; un número cuenta si es distinto de 0. En el Cronograma, la programación física de un año es obligatoria solo si la actividad tiene costo en ese periodo de la Cadena de valor. Fuera de un espejo es un error de configuración.
   - `numerar_por: "<grupo>"`, en una lista, solo lo usa el cliente: numera las filas n.m, donde n es la fila de `<grupo>` (de la ruta de la lista o, si es espejo, de su `espejo_de`) que las contiene y m su consecutivo dentro de ella. Las actividades de Cadena de valor se numeran por objetivo específico (1.1, 1.2… son las del objetivo 1, aunque estén en productos distintos) y el Cronograma muestra el número de cada actividad. El mid solo lo devuelve en la definición.
   - Casillas (`casilla`): `true` o `false` (ausente = sin marcar). Con `unica`, se marca en una sola fila de su lista y, si es obligatoria, en exactamente una (la alternativa que pasa al módulo de Preparación).
3. **Revalidar los submódulos dependientes.** Si el submódulo cambiado hace inválido otro que depende de él, ese otro vuelve a "Formulando" (o a "Sin formular", si está vacío). Por ejemplo, Objetivos al agregar una causa directa en Problemática, o Análisis técnico al marcar otra alternativa.
4. **Actualizar el nombre del plan** con el campo marcado con `nombre_plan` (el nombre de la iniciativa, propuesto o escrito por el usuario).

Los errores llegan como `[{ ruta, mensaje }]`. La ruta usa el formato `causas_directas[cd-1].causas_indirectas[ci-2].descripcion`.

### Estados

Son los de la «Relación de estados» del mockup (`docs/Refinamiento mockup para Banco de proyectos OAP.pdf`): cada estado general va con uno de submódulo. Se buscan por `codigo_abreviacion`; los crea y actualiza `database local/construccion-proyecto-inversion/cargar-estados.js`.

| Iniciativa (`estado-plan`) | Submódulo (`estado-submodulo`) | Cuándo |
|---|---|---|
| `SF_IPI` Sin formular | `SF_SUB` Sin formular | Nada diligenciado todavía |
| `EF_IPI` Formulación iniciada | `FN_SUB` Formulando | Con información y sin enviar; el submódulo, incompleto o con errores |
| `ER_IPI` Enviado a revisión | `FOR_SUB` Formulado | Todo formulado y enviado a Planeación; el submódulo, completo y sin errores |
| `AV_IPI` Avalado | `SO_SUB` Aprobado | Todo revisado y aprobado por Planeación |
| `PC_IPI` Por corregir | `CO_SUB` Con observaciones | Al menos un submódulo con observaciones por atender |

Son códigos propios para que las iniciativas no aparezcan en las bandejas de Formulación, que usan `EF_SP`, `ER_SP`, etc. El tipo de plan `IPI_SP` también lo crea `cargar-estados.js`.

| Acción | Desde | Hacia | Requisito |
|---|---|---|---|
| `enviar` | Formulación iniciada, Por corregir | Enviado a revisión | Todos los submódulos Formulados o Aprobados, y sin pendientes al volver a validarlos |
| `devolver` | Enviado a revisión | Por corregir | Al menos un submódulo Con observaciones |
| `avalar` | Enviado a revisión | Avalado | Todos los submódulos Aprobados |

- Al enviar se validan de nuevo todos los submódulos, sin llamadas extra al crud. Así no pasa un submódulo que quedó Formulado con datos que ya no son válidos (p. ej. si al guardar otro falló devolverlo a Formulando). Si alguno tiene pendientes, responde 409 con sus nombres.
- Solo se guardan submódulos con la iniciativa Sin formular, Formulación iniciada o Por corregir.
- Solo se revisan submódulos con la iniciativa Enviada a revisión.
- El control por rol (quién formula y quién revisa) lo hace el cliente, como en el resto de la aplicación.

### Observaciones

Cada submódulo guarda un **historial** de observaciones en `submodulo_iniciativa.observaciones`, en orden de registro. Cada observación es sobre el submódulo completo (no hay observaciones por campo):

```json
{ "_id": "…", "texto": "Precisar la fuente de la línea base.",
  "rol": "PLANEACION", "usuario": "1012345678", "usuario_nombre": "CARLOS RUIZ PENA",
  "fecha": "2026-11-04T10:15:00Z" }
```

- **Quién comenta lo decide el estado de la iniciativa:**
  - **Sin formular**, **Formulación iniciada** o **Por corregir**: comenta la dependencia. El estado del submódulo no cambia.
  - **Enviada a revisión**: comenta Planeación, solo sobre submódulos ya formulados. El submódulo queda **Con observaciones**.
  - **Avalada**: no se comenta (409).
- **Se agregan, nunca se reemplazan.** Forman un hilo, como un blog, con las rondas de revisión. Aprobar (`PUT /revision`) cambia el estado, pero no borra el historial.
- **No cambian la `version` del submódulo.** Comentar no provoca un 409 a quien tiene el formulario abierto.
- **`rol` y `usuario` los envía el cliente.** El cliente toma el documento del token y el rol que corresponde al momento: uno de Planeación al revisar, uno de la dependencia al formular. El mid no valida el token, igual que el resto de SISGPLAN.
- **`usuario_nombre` lo pone el mid** al publicar: busca el documento en Terceros (`nombreTercero`, en `servicio.go`). Si el documento no está o Terceros no responde, la observación se publica sin nombre y el motivo va al log.

## Organización del paquete

| Archivo | Responsabilidad |
|---|---|
| `modelo.go` | Tipos de la definición, documentos del crud, solicitudes y respuestas |
| `crud.go` | Interfaces pequeñas por recurso (`RepositorioPlanes`, `RepositorioArbol`, `RepositorioCatalogos`, `RepositorioSubmodulos`, `RepositorioElementosCatalogo`) y el cliente HTTP. A diferencia de `helpers.SendJson`, **trata 4xx/5xx y `Success: false` como error** |
| `definicion.go` | Lee el árbol de la plantilla. Una llamada para los hijos de cada nodo y otra para su detalle, en paralelo con cupo; marca `[[orden:N]]`; valida que el primer nivel sean módulos y que no se repitan claves de submódulo; caché |
| `catalogos.go` | Catálogos. `CatalogosConCache` los lee de planes_crud (`/elemento-catalogo`, a través de `RepositorioElementosCatalogo`) y guarda cada uno 10 minutos en memoria; un catálogo vacío o un error no se guardan. Para pasarlos a Parámetros basta otra implementación de `RepositorioElementosCatalogo` |
| `normalizacion.go` | Limpieza y campos calculados |
| `validacion.go` | Una regla por tipo de campo en un mapa: para soportar un tipo nuevo se registra una regla |
| `estados.go` | Tipo de plan, estados, transiciones (como datos) y resolución de códigos a ids |
| `servicio.go` | Casos de uso; recibe sus dependencias por constructor |
| `observaciones.go` | Caso de uso de agregar observaciones: validación y quién comenta en cada estado |
| `unidades.go` | Unidades del usuario (Terceros + Oikos), con su propio cliente HTTP |
| `errores.go` | `ErrorNegocio`: código HTTP, `Mensaje` para el usuario y `Detalle` técnico, que solo va al log |

## Pruebas

```bash
go test ./helpers/iniciativaInversionHelper/
```

Cubren:
- Validación de la iniciativa de ejemplo (`testdata/iniciativa_ejemplo.json`, copia de `database local/construccion-proyecto-inversion/iniciativa.ejemplo.json`) con la plantilla vigente sin Anexos (`testdata/iniciativa.plantilla.json`) y copias de los catálogos (`testdata/catalogos/`).
- Caché de catálogos: vigencia, catálogo vacío, error del crud y nombre inválido.
- Estados y transiciones.
- Servicio completo con un repositorio en memoria (crear con la plantilla indicada o la más reciente, listado con número y filtro por unidad, guardar, 409 por versión, revalidación de dependientes, flujo de revisión hasta el aval, y que un plan de otro tipo no sea iniciativa).
- Observaciones: quién comenta en cada estado, que no cambien la versión, que el historial no se borre y la validación de la solicitud.
- Carga del árbol y caché.

**Contra un crud real:** carga el tipo de plan y los estados (`cargar-estados.js`), los catálogos (`cargar-catalogos.js`) y la plantilla (`cargar-construccion.js`), todos en `database local/construccion-proyecto-inversion/`.

**Rutas generadas:** si se cambian las anotaciones `@router`, hay que regenerar `routers/commentsRouter_controllers.go`. Fuera del GOPATH, Beego no lo regenera al arrancar: `bee generate routers` lo genera en `routers/commentsRouter.go` con los imports de Beego v2; se pasan a `github.com/astaxie/beego`, las rutas entre comillas en vez de acentos graves, se guarda como `commentsRouter_controllers.go` (borrando el otro) y se le pasa `gofmt`.
