# Orden de los elementos del formato

## Decisión de diseño

No se agregó un campo `orden` al modelo MongoDB. Para un nodo cuyo padre es
otro subgrupo, la posición del identificador en `subgrupo.hijos` es la fuente
de verdad. Para los nodos cuyo padre es directamente un plan, `planes_crud`
ordena por `fecha_creacion` ascendente y usa `_id` ascendente como desempate
estable.

## Lectura

`GET /v1/arbol/{planId}` y `GET /v1/formato/{planId}` conservan el orden
recibido desde `planes_crud` en todos los niveles. El MID incluye además
`orden`, calculado desde cero dentro de cada arreglo `children` o `sub`; este
valor no se persiste.

Los nodos inactivos conservan su posición en `hijos`. Cuando una respuesta los
filtra, los nodos visibles reciben posiciones consecutivas sin modificar MongoDB.

## Clonación y propagación

Las operaciones de clonación recorren secuencialmente los hijos entregados por
`planes_crud`, por lo que el nuevo padre recibe sus IDs en el mismo orden.
Después de `formulacion/estructura_planes`, el MID llama a
`PUT /subgrupo/{idPadre}` una vez por cada subgrupo padre con el arreglo completo
de hijos en el orden de la plantilla. Los nodos exclusivos del plan se conservan
al final en su orden previo.

El primer nivel no se reordena durante la propagación: al no existir un arreglo
`hijos` en el plan, mantiene el criterio `fecha_creacion`, `_id`.

## Dependencia con planes_crud

El contrato de actualización, sus validaciones y el procedimiento para ordenar
datos históricos se documentan en `planes_crud/documentacion/orden-formato.md`.
