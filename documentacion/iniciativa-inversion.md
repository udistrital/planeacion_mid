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
| Iniciativa | `plan` de tipo `IPI_SP` con `formato: false`, `formato_id` = plantilla, `dependencia_id`, `estado_plan_id` (estado general), `responsable` (documento de quien la creó) y `responsable_nombre` (su nombre en Terceros al crearla) y, si se enviaron, `vigencia_aplica` (JSON de vigencias) |
| Datos de cada submódulo | `submodulo_iniciativa`: un documento por iniciativa y submódulo, con `dato` (JSON legible armado con las claves de la plantilla), `observaciones`, `estado_submodulo_id` y `version` |

**La iniciativa no copia el árbol de la plantilla:**
- La definición se lee de la plantilla una sola vez y queda en memoria hasta que se reinicia el mid. Si varias peticiones la piden a la vez, una la lee y las demás esperan ese resultado; un error no se guarda.
- Una plantilla con iniciativas no se edita. Si cambia el instrumento, se carga una plantilla nueva (con otro id, así que se lee la primera vez que se pide). Una plantilla corregida directamente en la base, en desarrollo, se ve al reiniciar el mid.

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

**Quién hace la solicitud (actor):** cada solicitud lleva el documento del usuario y todos sus roles. Las consultas los llevan en la URL (`?usuario=<documento>&roles=A,B`) y las escrituras, en el cuerpo (`usuario`, `roles`), junto a sus campos. Sin ellos responde 400. Con ellos el mid decide qué puede hacer el usuario (ver «Permisos»). Las consultas de la definición, de un submódulo, de los catálogos y de las unidades no los piden.

| Método y ruta | Qué hace | Errores |
|---|---|---|
| `GET /plantillas` | Plantillas activas con las que se pueden crear iniciativas, de la más reciente a la más antigua | 502 si falla el crud |
| `GET /iniciativas?usuario=&roles=` | Iniciativas activas que consulta el usuario (con un rol de la OAP, todas; si no, las de sus unidades), de la más reciente a la más antigua. Cada fila trae `numero` (consecutivo por orden de creación, contado sobre todas las iniciativas, activas o no), `estado` (general), `accion` (lo que el usuario hace con ella: `formular`, `revisar` o `consultar`), `fuera_de_fechas` (ver «Fechas de formulación») y `responsable` (nombre de quien la creó, el que se guardó al crearla en `responsable_nombre`; vacío si no se obtuvo de Terceros). El listado no consulta Terceros para los nombres; las unidades del usuario, una vez | 400 sin actor; 502 si falla el crud o Terceros |
| `POST /iniciativas` | Crea la iniciativa y todos sus submódulos en "Sin formular" (la iniciativa también), a nombre de quien la crea. Cuerpo: `{ usuario, roles, dependencia_id, plantilla_id?, nombre?, vigencias? }`. La crea el carril del formulador para una de sus unidades. El plan guarda el documento en `responsable` y su nombre en Terceros (`TercerosService`) en `responsable_nombre`; si no se obtiene el nombre, la crea igual, sin él. Sin `plantilla_id` usa la plantilla más reciente; sin `nombre`, "Iniciativa de inversión sin nombre" | 400 si falta `dependencia_id` o el actor; 403 si el usuario no formula en esa unidad; 409 si no hay plantillas o si la unidad está fuera de sus fechas de formulación |
| `GET /iniciativas/:id?usuario=&roles=` | Vista general: `estado`, `motivo_rechazo` (solo si está Rechazada), `dependencia_id`, `fuera_de_fechas`, `permisos` (ver «Permisos») y `modulos: [{ clave, nombre, submodulos: [{ clave, nombre, estado (con color), accion, observaciones, dependientes }] }]`. `accion` es lo que el usuario hace con el submódulo; `observaciones`, cuántas tiene; `dependientes` son las claves de los submódulos que leen datos de él | 400 sin actor; 403 si es de una unidad que el usuario no consulta; 404 |
| `GET /iniciativas/:id/definicion` | Construcción de la plantilla: `modulos`, y en cada uno sus submódulos y campos | 404 |
| `GET /iniciativas/:id/submodulos/:clave` | Datos, errores vigentes, historial de observaciones y `relacionadas` (datos de los submódulos de los que depende para espejos y referencias) | 404 |
| `PUT /iniciativas/:id/submodulos/:clave` | Guarda `{ usuario, roles, version, datos }`, aunque tenga errores (borrador). Responde el submódulo con sus `errores`; queda **Formulado** si no hay errores, **Formulando** si tiene información pero errores y **Sin formular** si está vacío. Con la primera información, la iniciativa pasa de Sin formular a **En formulación** | 403 si el usuario no la formula; 409 si el estado no deja editarla, si el submódulo está **Aprobado**, si está fuera de las fechas de su unidad o si `version` no es la vigente |
| `POST /iniciativas/:id/submodulos/:clave/documentos` | Sube el archivo de un campo `documento` al gestor documental (`GestorDocumental` de `conf/app.conf`, `POST /document/uploadAnyFormat`, tipo de documento 66) y responde 201 con su referencia `{ id, enlace, nombre, usuario, usuario_nombre, fecha }`. Como en las observaciones, `usuario_nombre` es el nombre en Terceros de quien lo sube (`TercerosService`); si no se obtiene, la referencia va sin él. Cuerpo: `{ usuario, roles, clave, nombre, archivo }`, con `clave` = ruta del campo en el submódulo sin ids de fila (`carta_presentacion`, `otros_anexos.archivo`) y `archivo` en base64. Se permite como guardar el submódulo. **No cambia el submódulo:** la referencia queda en la iniciativa con el `PUT` del submódulo, como cualquier valor. El cliente hace ese `PUT` solo, sin aviso, apenas responde la carga, y también al quitar un archivo o eliminar un adicional, así que esos cambios quedan de una vez; los nombres de los adicionales se guardan con Guardar o con el siguiente cambio de archivos. El archivo sigue en el gestor | 400 si la extensión no está en `formatos`, el archivo está vacío, tiene menos de 748 bytes (el gestor no recibe un base64 de menos de 1000 caracteres: responde «invalid file») o pasa de `tamano_maximo_mb`, o `clave` no es un campo `documento`; 403 y 409 como al guardar; 502 si falla el gestor |
| `POST /iniciativas/:id/submodulos/:clave/observaciones` | Publica una observación en el historial: `{ usuario, roles, texto }`. El mid la firma con el rol del carril con que se publica y agrega `usuario_nombre`, su nombre en Terceros (`TercerosService`); si no lo obtiene, la publica sin nombre. No cambia la `version` del submódulo. Del Jefe en 03 o de la OAP en 05, el submódulo queda **Con observaciones**; en 05, si era el último veredicto, la respuesta trae `estado_iniciativa` (ver «Estados») | 400 si falta el texto o el actor; 403 si el usuario no es del carril de ese estado; 409 si el estado no recibe observaciones, si el submódulo no está formulado (Jefe y OAP) o si está fuera de fechas |
| `PUT /iniciativas/:id/submodulos/:clave/revision` | Veredicto **Aprobado** de la OAP en 05: `{ usuario, roles, version }`. El historial se conserva. Si era el último veredicto, la respuesta trae `estado_iniciativa` | 403 si el usuario no es de la OAP; 409 si la iniciativa no está En revisión OAP, el submódulo no está formulado o está fuera de fechas |
| `PUT /iniciativas/:id/acciones/:accion` | `enviar`, `aprobar`, `devolver`, `iniciar-revision`, `rechazar` o `retractar-rechazo`. Cuerpo: `{ usuario, roles, motivo? }`; `motivo` es obligatorio al rechazar (máx. 4000): queda en `motivo_rechazo` del plan y el cliente lo manda también en la notificación. Responde la vista general | 400 si la acción no existe o falta el motivo; 403 si el usuario no es del carril de la acción; 409 si no se permite en el estado (o sus submódulos no lo cumplen) o está fuera de fechas |
| `GET /unidades/:documento` | Unidades del usuario (`[{ Id, Nombre }]`): tercero por documento y vinculaciones activas en Terceros, con los mismos cargos que Formulación (`/v1/formulacion/vinculacion_tercero`: 319, 312, 320 y 300), y nombre en Oikos. Lista vacía si no tiene. El documento solo admite letras y números (máx. 20). Usa `TercerosService` y `OikosService` de `conf/app.conf` | 400 si el documento no es válido; 502 si falla Terceros u Oikos |
| `GET /catalogos/:nombre?padre=` | Cualquier catálogo de la colección general `elemento_catalogo` del crud: `procesos`, `productos` (solo los del programa 2202 del Plan Nacional de Desarrollo), `municipios` (padre = departamento), `tipos_riesgo`… | 404 si no existe o no está cargado |

## Reglas

### Al guardar un submódulo

1. **Normalizar:**
   - Se quitan las claves que no están en la plantilla y se recortan los textos.
   - Se recalculan los campos `calculado` y se ignora el valor que envíe el cliente. Excepción: si el campo es `editable` y el cliente envía un texto no vacío, se guarda ese texto (con la misma transformación, p. ej. mayúsculas). Así funciona el nombre de la iniciativa: el sistema lo propone y el usuario puede cambiarlo. Una `plantilla` sin `{clave}` es un texto fijo que se guarda igual en todas las iniciativas: así van el sector (`22 - EDUCACIÓN`) y el programa (`2202 - Calidad y fomento de la educación superior`) del Plan Nacional de Desarrollo en Articulación, que no se seleccionan. Con la `plantilla` vacía (`""`), el campo se muestra sin valor y se guarda vacío: así van el indicador y la unidad de medida del producto en Cadena de valor (`indicador_producto`, `unidad_medida_producto`) mientras no esté la relación producto → indicador.
   - Un `calculado` con `operacion` es numérico: `"operacion": "resta", "operandos": ["oferta", "demanda"]` guarda el primer campo hermano menos los demás (el déficit de Necesidades); si falta un operando, queda sin valor. `"operacion": "suma"` guarda la suma, con los vacíos en 0 (el costo total de cada actividad de Cadena de valor, de `periodo_0` a `periodo_4`); si todos están vacíos, queda sin valor. `"operacion": "multiplicacion"` guarda el producto (el valor total de cada año de Ingresos y beneficios: `cantidad` × `valor_unitario`); si falta un operando, queda sin valor.
   - Un `calculado` con `desde: "<submodulo>.<campo>"` solo muestra ese valor de otro submódulo (el objetivo general en Cadena de valor): no se guarda ni se valida, y el submódulo de origen va en `relacionadas`. Con `desde: "origen.<campo>"`, en un campo de un espejo, muestra ese campo de la fila de origen (en el Cronograma, el costo de cada periodo de la actividad en Cadena de valor y su costo total); tampoco se guarda. Con `desde: "<ruta de un grupo>.<campo>"`, en un espejo cuyo `espejo_de` pasa por ese grupo, muestra el campo de la fila del grupo que contiene la de origen (en el Presupuesto inicial, `cadena_valor.productos_por_objetivo.productos.producto` es el producto de cada actividad). Con `catalogo`, el cliente muestra el nombre del elemento en vez del código.
   - Un espejo en el que todos los campos se traen (`desde`) es un resumen de solo lectura: solo guarda `id` y `origen` de cada fila, y queda Formulado al guardarlo con una fila por elemento de origen (el Presupuesto inicial de la alternativa seleccionada). Con `totales: { "columnas": true }`, el cliente agrega una fila con la suma de cada columna numérica.
2. **Validar** según el tipo de campo:
   - Obligatorios y longitud máxima.
   - Mínimo y máximo, y decimales.
   - Códigos de catálogo y su `depende_de`. Toda lista de opciones es un catálogo: la plantilla no lleva opciones escritas.
   - Mínimo y máximo de filas, e `id` únicos.
   - Espejos: exactamente una fila por cada fila de origen. Con `espejo_de: "catalogo.<nombre>"` las filas de origen son los elementos del catálogo (`{ id: codigo, nombre }`, con `mostrar: "nombre"`); el mid se los manda al cliente en `relacionadas.catalogo.<nombre>` (los resultados esperados de Objetivos). Con `filas_opcionales: true`, la fila vacía no se exige ni se valida (no aplica), la que tiene algún dato se valida completa y `min` es cuántas filas se llenan como mínimo. `filas_obligatorias` (lista de ids de origen, p. ej. códigos del catálogo) son las filas que sí se exigen: en Resultados esperados, estudiantes, docentes y administrativos beneficiados. Un campo de un espejo puede tener `solo_si: { campo: "origen", valor: "<id o código>" }`: solo aplica a esa fila (el «¿Cuál?» de la fila «Otro»); el cliente lo escribe en la celda del elemento de origen. `filas_adicionales: { origen, max, nombre, etiqueta }` deja que el usuario agregue filas que no salen del origen, hasta `max` si lo tiene (en Trazadores presupuestales, hasta dos «Otro: ¿Cuál?»; en Fuentes de financiación, sin límite, «Otras, ¿cuáles?»): su `origen` vale `origen`, que no puede ser un elemento del origen; se validan como las demás (con `filas_opcionales`, la vacía no se valida) y sus campos propios llevan `solo_si` sobre ese valor (el `cual`). `nombre` es el texto de su celda de origen (solo lo usa el cliente) y `etiqueta`, el de los botones de agregar y eliminar y el del error del mid cuando hay más de `max`. `columna_codigo: { nombre, campo }`, en un espejo de catálogo, también es solo del cliente: muestra el código de cada elemento en una columna antes de su nombre (Fuentes de financiación), y las filas agregadas escriben el suyo en el campo `campo` (texto con `solo_si` sobre su origen; el mid lo valida como cualquier campo). En una matriz, `etiqueta_fila` es el encabezado de la columna de las filas («Grupo etáreo y ciclo de vida» en Población).
   - Referencias: solo filas válidas en su alcance. Con `marcada_en` (clave de una casilla única del grupo de origen), la única opción es la fila marcada: el cliente la llena sola y, si después se marca otra, el submódulo pide guardarse de nuevo (Análisis técnico con la alternativa que pasa a Preparación).
   - Matrices completas.
   - Catálogos con `multiple`: lista de códigos del catálogo, sin repetir; el cliente los muestra como casillas (las sedes y los factores de Localización). Con `depende_de`, cada código debe colgar de lo seleccionado en ese campo, que puede ser un código o una lista: así se guardan los niveles del PDD, el PED y el PI en Articulación (varios por nivel, como en la armonización de Formulación), que el cliente elige en un árbol. El catálogo `sedes` no está en `elemento_catalogo`: `RepositorioConSedes` (`sedes.go`) lo lee de Oikos (`OikosService`), de los espacios físicos de tipo SEDE (38) activos, con el `Id` como código (el `Codigo` de Oikos se repite) y «Otra» al final.
   - Los catálogos del Plan Distrital de Desarrollo, del Plan Estratégico de Desarrollo y del Plan Indicativo (Articulación) tampoco están en `elemento_catalogo`: `RepositorioConInstrumentos` (`instrumentos.go`) los lee de los planes activos de tipo `PDD_SP`, `PD_SP` y `PLI_SP` armados en Construcción (`formato: true`; las copias con `formato: false` son formulaciones de una dependencia), los mismos con los que se armoniza en Formulación y en MACRO. Cada catálogo es un nivel del árbol: `planes_pdd` (los planes, con sus vigencias en el nombre; un plan sin nodos no se ofrece), `objetivos_pdd`, `estrategias_pdd` y `programas_pdd`; `planes_ped`, `lineamientos_ped` y `metas_ped`; `planes_pi`, `ejes_pi`, `lineamientos_pi` y `estrategias_pi`. Si un nivel queda vacío, su catálogo responde 404 como cualquier catálogo sin elementos. El código es el id del subgrupo y el padre, el del nodo de arriba, así que `depende_de` encadena plan → lineamiento → meta. La descripción del nodo va en la `descripcion` del elemento (el enunciado de cada meta). Todos los catálogos de un plan salen de una sola lectura del árbol, que se guarda un minuto; después, como todo catálogo, cada uno queda 10 minutos en `CatalogosConCache`, así que un cambio en esos planes tarda hasta 10 minutos en verse (o reiniciar el mid).
   - `sin_codigo: true`, en un catálogo, solo lo usa el cliente: la lista muestra el nombre sin el código, porque el código es un id interno (los de Articulación). Si el elemento trae `descripcion`, el cliente la muestra en gris bajo su nombre en la lista y bajo el campo al elegirlo.
   - `solo_si: { campo, valor }`: el campo solo se exige y se guarda si el campo hermano vale `valor` o, si es una lista, lo incluye («¿Cuál otra sede?» solo si se marcó «Otra»). Con `"distinto": true`, al revés: aplica si no lo vale. Si no aplica, el normalizador lo descarta. En Trazadores presupuestales, Discriminado y Número de beneficiarios aplican a todos los trazadores menos a Grupos Étnicos, y la lista de comunidades (un espejo de `catalogo.grupos_etnicos_trazador`, con `filas_opcionales` y `min: 1`) solo a Grupos Étnicos.
   - Los elementos de un catálogo pueden traer `descripcion` (campo opcional de `elemento_catalogo`). En un espejo de catálogo, la fila de origen la incluye y el cliente la muestra como ayuda bajo el nombre (la definición de cada trazador presupuestal).
   - `obligatorio_si: "origen.<campo>"` (con `required: "false"`), en un campo de un espejo: el campo se exige solo si ese campo de la fila de origen tiene valor; un número cuenta si es distinto de 0. En el Cronograma, la programación física de un año es obligatoria solo si la actividad tiene costo en ese periodo de la Cadena de valor. Fuera de un espejo es un error de configuración.
   - `numerar_por: "<grupo>"`, en una lista, solo lo usa el cliente: numera las filas n.m, donde n es la fila de `<grupo>` (de la ruta de la lista o, si es espejo, de su `espejo_de`) que las contiene y m su consecutivo dentro de ella. Las actividades de Cadena de valor se numeran por objetivo específico (1.1, 1.2… son las del objetivo 1, aunque estén en productos distintos) y el Cronograma muestra el número de cada actividad. El mid solo lo devuelve en la definición.
   - Casillas (`casilla`): `true` o `false` (ausente = sin marcar). Con `unica`, se marca en una sola fila de su lista y, si es obligatoria, en exactamente una (la alternativa que pasa al módulo de Preparación).
   - Documentos (`documento`, los de Anexos): el valor es la referencia que respondió `POST .../documentos`. Debe traer `enlace` y `nombre`, y la extensión del nombre debe estar en `formatos`; el normalizador guarda solo `id`, `enlace`, `nombre`, `usuario`, `usuario_nombre` y `fecha`. `formatos` y `tamano_maximo_mb` también los usa el cliente para revisar el archivo antes de subirlo; el cliente muestra «Subido el … por `usuario_nombre`» y, sin nombre, solo la fecha. Ejemplo: `"carta_presentacion": { "id": 123, "enlace": "uuid", "nombre": "carta.pdf", "usuario": "1234567", "usuario_nombre": "LAURA GOMEZ RIOS", "fecha": "2026-11-20T15:00:00Z" }`; los adicionales son un repetible con `nombre` y `archivo` (`"otros_anexos": [ { "id": "oa-1", "nombre": "Cotización", "archivo": { … } } ]`).
3. **Revalidar los submódulos dependientes.** Si el submódulo cambiado hace inválido otro que depende de él, ese otro vuelve a "Formulando" (o a "Sin formular", si está vacío). Por ejemplo, Objetivos al agregar una causa directa en Problemática, o Análisis técnico al marcar otra alternativa. Los **Aprobados** no cambian: nadie los edita, y quien revisa decide si los vuelve a observar.
4. **Actualizar el nombre del plan** con el campo marcado con `nombre_plan` (el nombre de la iniciativa, propuesto o escrito por el usuario).

Los errores llegan como `[{ ruta, mensaje }]`. La ruta usa el formato `causas_directas[cd-1].causas_indirectas[ci-2].descripcion`.

### Estados

El flujo completo, con quién puede qué, está en `docs/Máquina de estados - Iniciativas de inversión.md` (y su diagrama `.drawio`). Los estados se buscan por `codigo_abreviacion`; los crea y actualiza `database local/construccion-proyecto-inversion/cargar-estados.js`, que también renombra en el mismo registro los códigos viejos (`AV_IPI` → `APR_IPI`, `SO_SUB` → `APR_SUB`). Son códigos propios para que las iniciativas no aparezcan en las bandejas de Formulación, que usan `EF_SP`, `ER_SP`, etc.

| N.° | Iniciativa (`estado-plan`) | Qué significa |
|:---:|---|---|
| 01 | `SF_IPI` Sin formular | Recién creada |
| 02 | `EF_IPI` En formulación | La unidad la está diligenciando |
| 03 | `PA_IPI` Por aprobar (Jefe de unidad) | Enviada al Jefe, que la revisa |
| 04 | `PR_IPI` Por revisar OAP | Radicada ante la OAP |
| 05 | `ER_IPI` En revisión OAP | La OAP da un veredicto por submódulo |
| 06 | `PC_IPI` Por corregir | Devuelta con observaciones |
| 07 | `APR_IPI` Aprobado | Todos los submódulos aprobados. Fin |
| 08 | `RCH_IPI` Rechazada | Rechazada por la OAP; vuelve a 06 si la OAP se retracta |

Submódulo (`estado-submodulo`): `SF_SUB` Sin formular, `FN_SUB` Formulando y `FOR_SUB` Formulado los pone el guardado; `CO_SUB` Con observaciones, el Jefe en 03 o la OAP en 05 al publicar una observación (también sobre un aprobado, que así se desaprueba); `APR_SUB` Aprobado, la OAP en 05. Nadie edita un aprobado.

Las transiciones son datos (`transiciones`, en `estados.go`):

| Acción | Desde | Hacia | Carril | Requisito |
|---|---|---|---|---|
| `enviar` | 02, 06 | 03 | Formulador | Todos los submódulos Formulados o Aprobados, y los no aprobados sin pendientes al volver a validarlos |
| `aprobar` | 03 | 04 | Jefe | Ningún submódulo Con observaciones |
| `devolver` | 03 | 06 | Jefe | Al menos un submódulo Con observaciones |
| `iniciar-revision` | 04 | 05 | OAP | — |
| `rechazar` | 05 | 08 | OAP | `motivo` |
| `retractar-rechazo` | 08 | 06 | OAP | — |
| (automático) | 05 | 07 o 06 | — | Tras cada veredicto (`EstadoTrasRevision`): todos Aprobados → 07; todos con veredicto y alguno Con observaciones → 06 |

- Al enviar se validan de nuevo los submódulos no aprobados, sin llamadas extra al crud. Así no pasa un submódulo que quedó Formulado con datos que ya no son válidos (p. ej. si al guardar otro falló devolverlo a Formulando). Si alguno tiene pendientes, responde 409 con sus nombres.
- **Cambio automático:** después de cada veredicto (`cerrarRevision`), el mid relee la iniciativa y sus submódulos y, si sigue En revisión OAP y ya todos tienen veredicto, la cambia. Si dos revisores dan a la vez los dos últimos veredictos, el segundo la encuentra fuera de 05 y no hace nada. Si el cambio falla, el veredicto ya quedó: va al log y se aplica con el siguiente veredicto. La respuesta del veredicto trae `estado_iniciativa` cuando la iniciativa cambió, para que el cliente lo notifique.
- **Rechazada** es solo un estado: el plan sigue con `activo: true`, aparece en los listados y nadie lo edita.

### Permisos

Lo resuelve `permisos.go`. Cada carril del diagrama es un grupo de roles:

| Carril | Roles | Sobre qué iniciativas |
|---|---|---|
| Formulador | `JEFE_DEPENDENCIA`, `ASISTENTE_DEPENDENCIA`, `ASISTENTE_PLANEACION` | Las de sus unidades |
| Jefe | `JEFE_DEPENDENCIA` | Las de sus unidades |
| OAP | `PLANEACION`, `ASISTENTE_PLANEACION` | Todas |

- **Qué hace cada carril en cada estado** (`tareasPorEstado`): el formulador edita y sube anexos en 01, 02 y 06; el Jefe observa en 03; la OAP observa y aprueba submódulos en 05. Las acciones van con su carril en `transiciones`. En cada estado solo actúa su carril; los demás consultan. Nada se asigna a una persona.
- **Varios roles:** se suman los carriles. Cuando más de un rol del usuario sirve, firma con el de mayor jerarquía: `PLANEACION`, `ASISTENTE_PLANEACION`, `JEFE_DEPENDENCIA`, `ASISTENTE_DEPENDENCIA`.
- **Unidades del usuario:** son las de `GET /unidades/:documento` (Terceros, mismos cargos que Formulación). Las consulta el mid solo cuando un carril las necesita, una vez por solicitud, y las guarda 5 minutos por usuario (`CacheDependencias`, en `unidades.go`); las fallas no se guardan.
- **Consultar:** con un rol de la OAP, cualquier iniciativa; si no, solo las de sus unidades (403 en la vista general).
- **En las consultas,** `permisos` de la vista general dice lo que el usuario puede hacer hoy: `editar`, `observar`, `aprobar_submodulos` y `acciones: [{ accion, habilitada, motivo }]`, las de su carril en ese estado, con el motivo si no está habilitada (faltan submódulos por formular, hay observaciones, fuera de fechas). El cliente solo pinta lo que dice el mid.
- **El rol y el documento los envía el cliente,** que los toma del token, como en el resto de SISGPLAN: el mid no valida el token. Sí valida el carril, la unidad y las fechas.

### Fechas de formulación

Las iniciativas se trabajan solo en las fechas que Planeación habilita para cada unidad: los registros activos del tipo de seguimiento `FIPI_SP` en `periodo_seguimiento`. Lo resuelve `fechas.go`:
- **Qué se cierra:** todo lo que escribe, de cualquier carril: crear, guardar, subir un anexo, publicar observaciones, dar veredictos y ejecutar acciones responden 409 «La formulación de iniciativas no está habilitada para la unidad de esta iniciativa en estas fechas.», como cuando el estado no lo permite. La iniciativa queda congelada en su estado y solo se consulta; para seguir, Planeación amplía las fechas.
- **Qué unidad cuenta:** la `dependencia_id` de la iniciativa (al crearla, la de la solicitud), no la de quien la abre. Como la iniciativa no guarda la vigencia, sirve un registro de cualquier vigencia que cubra hoy esa unidad, por día en hora de Bogotá e inclusivo en los dos extremos (`UnidadHabilitada`, de `formulacionIniciativasHelper`).
- **En las consultas:** el listado y la vista general traen `fuera_de_fechas`: el usuario tendría algo que hacer en ese estado, pero hoy no hay fechas para la unidad. Entonces sus `permisos` no habilitan nada y cada acción trae el mensaje de fechas como motivo. Los registros se leen una vez por consulta, y solo si el usuario tiene algo que hacer en alguna iniciativa. Si no se pueden leer (falla el crud o un registro está mal formado), la iniciativa se muestra fuera de fechas y el motivo va al log; al escribir, ese error se responde (502 o 500).
- **El tipo de seguimiento** se busca por código (lo crea `cargar-estados.js`) y queda en memoria cuando existe. Sin él, las iniciativas se consultan, pero no se crean ni se formulan (500, con el motivo en el log).

### Observaciones

Cada submódulo guarda un **historial** de observaciones en `submodulo_iniciativa.observaciones`, en orden de registro. Cada observación es sobre el submódulo completo (no hay observaciones por campo):

```json
{ "_id": "…", "texto": "Precisar la fuente de la línea base.",
  "rol": "PLANEACION", "usuario": "1012345678", "usuario_nombre": "CARLOS RUIZ PENA",
  "fecha": "2026-11-04T10:15:00Z" }
```

- **Quién publica lo decide el estado de la iniciativa** (`tareasPorEstado`). Las observaciones son para que la unidad corrija:
  - **03**: observa el Jefe; **05**: observa la OAP. Solo sobre submódulos ya formulados, también un aprobado. El submódulo queda **Con observaciones**; en 05 es el veredicto «Con observaciones».
  - **01, 02, 04, 06, 07 y 08**: nadie publica (409). En 01, 02 y 06 la unidad solo formula y lee las observaciones.
- **Se agregan, nunca se reemplazan.** Forman un hilo, como un blog, con las rondas de revisión. Aprobar (`PUT /revision`) cambia el estado, pero no borra el historial.
- **No cambian la `version` del submódulo.** Quien revisa sigue dando veredictos con la versión que tiene abierta, sin 409.
- **`rol` lo pone el mid:** el del carril con que se publica, según la jerarquía de «Permisos». `usuario` es el documento del actor.
- **`usuario_nombre` lo pone el mid** al publicar: busca el documento en Terceros con el cliente que recibe el servicio (`FuenteNombreTercero`, que implementa `ClienteUnidades` en `unidades.go`). Si el documento no está o Terceros no responde, la observación se publica sin nombre y el motivo va al log.

## Organización del paquete

| Archivo | Responsabilidad |
|---|---|
| `modelo.go` | Tipos de la definición, documentos del crud, solicitudes y respuestas |
| `crud.go` | Interfaces pequeñas por recurso (`RepositorioPlanes`, `RepositorioArbol`, `RepositorioCatalogos`, `RepositorioSubmodulos`, `RepositorioPeriodos`, `RepositorioElementosCatalogo`) y el cliente HTTP. A diferencia de `helpers.SendJson`, **trata 4xx/5xx y `Success: false` como error** |
| `definicion.go` | Lee el árbol de la plantilla. Una llamada para los hijos de cada nodo y otra para su detalle, en paralelo con cupo; marca `[[orden:N]]`; valida que el primer nivel sean módulos y que no se repitan claves de submódulo; caché sin vencimiento con una sola lectura por plantilla |
| `catalogos.go` | Catálogos. `CatalogosConCache` los lee de planes_crud (`/elemento-catalogo`, a través de `RepositorioElementosCatalogo`) y guarda cada uno 10 minutos en memoria; un catálogo vacío o un error no se guardan. Para pasarlos a Parámetros basta otra implementación de `RepositorioElementosCatalogo` |
| `normalizacion.go` | Limpieza y campos calculados |
| `validacion.go` | Una regla por tipo de campo en un mapa: para soportar un tipo nuevo se registra una regla |
| `estados.go` | Tipo de plan, estados, transiciones (como datos, con su carril), cambio automático tras la revisión y resolución de códigos a ids |
| `permisos.go` | Roles, carriles, qué hace cada carril en cada estado (`tareasPorEstado`), el actor de la solicitud, `autorizar` (carril, unidad y fechas) y los permisos que traen las consultas |
| `servicio.go` | Casos de uso; recibe sus dependencias por constructor |
| `fechas.go` | Fechas de formulación: cierra todo lo que escribe fuera de ellas y calcula `fuera_de_fechas` para las consultas |
| `observaciones.go` | Caso de uso de agregar observaciones: validación y quién publica en cada estado |
| `documentos.go` | Campos `documento`: caso de uso de subir el archivo (permiso, campo, formato y tamaño), referencia que se guarda y cliente HTTP del gestor documental (`GestorDocumental`, inyectado como los demás) |
| `unidades.go` | Unidades del usuario (Terceros + Oikos), con su propio cliente HTTP, que el servicio también recibe para el nombre de las personas (Terceros) y el de las dependencias (Oikos), y la caché de unidades que usan los permisos |
| `errores.go` | `ErrorNegocio`: código HTTP, `Mensaje` para el usuario y `Detalle` técnico, que solo va al log |

## Pruebas

```bash
go test ./helpers/iniciativaInversionHelper/
```

Cubren:
- Validación de la iniciativa de ejemplo (`testdata/iniciativa_ejemplo.json`, copia de `database local/construccion-proyecto-inversion/iniciativa.ejemplo.json`) con la plantilla vigente sin Anexos (`testdata/iniciativa.plantilla.json`) y copias de los catálogos (`testdata/catalogos/`).
- Caché de catálogos: vigencia, catálogo vacío, error del crud y nombre inválido.
- Estados y transiciones: cada acción permitida y rechazada por estado y por requisito, y el cambio automático tras la revisión.
- Permisos: el rol de cada carril según la jerarquía, la acción de cada submódulo, la caché de unidades y el error sin fuente de unidades.
- Servicio completo con un repositorio en memoria y Terceros, Oikos y las unidades falsos, con un usuario por rol (formulador, Jefe, OAP y uno con Jefe y OAP): crear solo como formulador y en sus unidades, con el usuario y sus roles obligatorios y el nombre del responsable guardado (o sin él si Terceros no responde); listado con número, lo que ve cada uno y la acción de cada fila; guardar, 409 por versión y por submódulo aprobado, revalidación de dependientes (salvo los aprobados); el flujo completo 01 → 07 con devoluciones del Jefe, el cambio automático a 06 y a 07 y un aprobado desaprobado por el Jefe; rechazo con motivo y retractación; qué hace cada carril en cada estado; el Jefe de otra unidad; la suma de carriles con varios roles; y que un plan de otro tipo no sea iniciativa.
- Observaciones: quién publica en cada estado y con qué rol, que no cambien la versión, que el historial no se borre y la validación de la solicitud.
- Documentos: subir responde la referencia sin cambiar el submódulo; se rechaza (sin llamar al gestor) un campo que no es documento, un formato no permitido, un archivo vacío o muy grande, una iniciativa que no se puede editar y un submódulo aprobado; la referencia al guardar y su validación; el cliente del gestor contra un servidor de prueba. Usan un Anexos armado en la prueba (`definicionConAnexos`).
- Fechas de formulación: fuera de ellas no se crea, no se guarda, no se suben anexos ni se envía, el Jefe no observa y la OAP no da veredictos ni rechaza; quien no tiene nada que hacer en el estado no las consulta; cuenta la unidad de la iniciativa y el listado las lee una vez; si no se pueden revisar o falta `FIPI_SP`, solo se consulta; el cliente del crud pide los registros activos y trata `Success: false` como error.
- Carga del árbol y caché de definiciones: una lectura por plantilla, compartida por las peticiones simultáneas; un error o una lectura interrumpida no se guardan.

**Contra un crud real:** carga el tipo de plan y los estados (`cargar-estados.js`), los catálogos (`cargar-catalogos.js`) y la plantilla (`cargar-construccion.js`), todos en `database local/construccion-proyecto-inversion/`.

**Rutas generadas:** si se cambian las anotaciones `@router`, hay que regenerar `routers/commentsRouter_controllers.go`. Fuera del GOPATH, Beego no lo regenera al arrancar: `bee generate routers` lo genera en `routers/commentsRouter.go` con los imports de Beego v2; se pasan a `github.com/astaxie/beego`, las rutas entre comillas en vez de acentos graves, se guarda como `commentsRouter_controllers.go` (borrando el otro) y se le pasa `gofmt`.
