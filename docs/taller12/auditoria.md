# Auditoría de posibles inyecciones SQL

## Revisión de coincidencias

| Archivo y línea (auditoría original) | ¿Arma SQL? | ¿Llega texto del cliente? | ¿Va aparte como parámetro? | Análisis y acción |
|---|---|---|---|---|
| internal/almacen/filas/consultas.sql.go:19 | Sí | Sí, el identificador recibido por el repositorio | Sí, mediante SQLC y `$1` | Consulta parametrizada. No requiere cambios. |
| internal/almacen/filas/consultas.sql.go:37 | Sí | Sí, ID y estado recibidos por el repositorio | Sí, mediante SQLC y `$1`, `$2` | Actualización parametrizada. No requiere cambios. |
| internal/libros/manejadores.go:20 | Sí | No | No aplica; `Order("id ASC")` es una expresión fija | No recibe texto del cliente. Se mantiene la expresión fija. |
| internal/libros/manejadores.go:30 | Sí | Sí, parámetro `estado` de la URL | Sí, mediante `?` | Filtro parametrizado. No requiere cambios. |
| internal/libros/manejadores.go:33 | Sí | El filtro anterior puede contener `estado` | Sí, conserva la parametrización de GORM | Ejecuta la consulta construida con el filtro parametrizado. |
| internal/libros/manejadores.go:84 | Sí | Sí, ID de la ruta | No aplica como parámetro explícito; se entrega como identificador a GORM | `ObtenerLibro` valida el ID con `strconv.ParseUint` antes de consultar. |
| internal/libros/manejadores.go:103 | Sí | Sí, ID de la ruta | No aplica como parámetro explícito; se entrega como identificador a GORM | Se agregó validación numérica con `strconv.ParseUint` antes de consultar en `ActualizarLibro`. |
| internal/libros/manejadores.go:148 | Sí | Sí, ID de la ruta | No aplica como parámetro explícito; se entrega como identificador a GORM | Se agregó validación numérica con `strconv.ParseUint` antes de consultar en `EliminarLibro`. |
| internal/libros/manejadores.go:160 | Sí | Sí, ID de la ruta | Sí, mediante `?` | Filtro parametrizado. No requiere cambios. |
| internal/repository/intercambios_repository.go:34 | Sí | No se observa un dato del cliente en esta operación | No aplica; consulta de listado mediante GORM | No hay filtro dinámico proporcionado por el cliente. |
| internal/repository/intercambios_repository.go:65 | Sí | Sí, identificador recibido por el repositorio | Sí, mediante `?` | Consulta parametrizada. No requiere cambios. |
| internal/repository/intercambios_repository.go:117 | Sí | Sí, identificador recibido por el repositorio | Sí, mediante `?` | Filtro parametrizado. No requiere cambios. |

## Coincidencias que no construyen SQL

| Archivo y líneas (auditoría original) | Motivo |
|---|---|
| internal/almacen/filas/db.go:13, 15 y 16 | Son declaraciones de interfaces para ejecutar consultas; no construyen SQL por sí mismas. |
| internal/libros/manejadores.go:22 | Obtiene el parámetro `estado` de la URL; no construye SQL por sí solo. |
| main.go:151 | Configura la dirección del servidor HTTP; no construye SQL. |

## Conclusión

En las coincidencias revisadas no se identificó concatenación de datos del cliente dentro de consultas SQL. Las consultas de SQLC utilizan parámetros y los filtros de GORM revisados emplean parámetros enlazados.

Se reforzó la validación de identificadores en las funciones `ActualizarLibro` y `EliminarLibro`, utilizando `strconv.ParseUint`, igual que en `ObtenerLibro`. Se agregaron pruebas para verificar que las tres rutas rechacen identificadores no numéricos con HTTP 400.

Las pruebas de `internal/libros` finalizaron correctamente.

La auditoría estática no sustituye la prueba de inyección con peticiones HTTP, que se documentará en `inyeccion.txt`.
