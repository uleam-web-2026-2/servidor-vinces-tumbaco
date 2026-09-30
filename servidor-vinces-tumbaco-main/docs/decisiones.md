\# Decisión de diseño D2: eliminación de libros con préstamos



\## Decisión



No se permite eliminar un libro cuando tiene préstamos registrados.



\## Justificación



Esta decisión permite conservar la información de los préstamos relacionados con cada libro y evita eliminar un libro que todavía tiene información asociada.



Cuando se intenta eliminar un libro que tiene préstamos, el sistema responde con el código HTTP 422.



\## Alternativa



La alternativa sería permitir la eliminación del libro y eliminar también sus préstamos relacionados. Esta opción podría provocar pérdida de información de los préstamos, por lo que no se utiliza en este proyecto.

## Semana 4 · D1 · Variables de configuración

**Decisión:** `DATABASE_URL` no tiene valor por defecto. `PUERTO` tiene como valor por defecto `8080` y `TIEMPO_ESPERA_SEGUNDOS` tiene como valor por defecto `5`.

**Por qué:** `DATABASE_URL` es necesaria para conectarse a nuestra base de datos y no sería seguro inventar una conexión. El puerto y el tiempo de espera tienen valores comunes para poder iniciar el servidor sin configurarlos cada vez.

**Qué descartamos:** No usar un valor por defecto para `DATABASE_URL`, porque podría intentar conectarse a una base de datos incorrecta.

## Semana 4 · D2 · Regla propia de negocio

**Decisión:** No se permite crear un libro cuando el título y el autor son iguales.

**Por qué:** Esta regla evita registrar datos incorrectos o poco útiles en la información de los libros. Cuando se incumple, el sistema responde con HTTP 422.

**Qué descartamos:** Permitir que el título y el autor sean iguales, porque podría generar registros confusos.

## Semana 4 · D3 · Pruebas

**Prueba que no escribimos:** No escribimos una prueba que consulte directamente la base de datos, porque las pruebas de esta semana se diseñaron para validar primero las respuestas del servidor sin depender de PostgreSQL.

**Dónde empieza la base:** En el manejador `CrearLibro`, el acceso a la base de datos comienza en la línea 65, con `m.DB.Create(&libro)`. Antes de esa línea se validan los datos, por eso las pruebas de error no necesitan acceder a la base de datos.

**Qué haría falta para probarla:** Sería necesario levantar una base de datos de prueba y preparar datos para comprobar las consultas directamente.
