# Decisiones técnicas del proyecto

## Plataforma de Intercambio de Libros Usados

### 1. Objetivo

Implementar una API REST que permita gestionar intercambios de libros usados, aplicando validaciones, autenticación JWT y control de permisos.

### 2. Tecnologías utilizadas

- **Go:** desarrollo del servidor y los endpoints REST.
- **PostgreSQL:** almacenamiento de los datos.
- **GORM:** manejo de modelos y operaciones con la base de datos.
- **JWT:** autenticación mediante tokens.
- **OpenAPI 3.1:** documentación de los endpoints.
- **Redocly:** validación y generación de documentación HTML.

### 3. Decisiones de implementación

**Validación de estados:** Se definieron cuatro estados: pendiente, aceptado, entregado y cancelado.

**Transiciones permitidas:** Un intercambio puede pasar de pendiente a aceptado o cancelado, y de aceptado a entregado o cancelado. Las transiciones no permitidas devuelven HTTP 422.

**Autenticación JWT:** Se implementó el endpoint POST /login con dos usuarios de demostración. El token tiene una duración de 30 minutos.

**Control de roles:** Los administradores pueden consultar, crear, actualizar y eliminar intercambios. Los usuarios normales solamente pueden consultar.

**Respuestas de seguridad:** La API devuelve HTTP 401 cuando falta un token válido y HTTP 403 cuando el usuario no tiene permisos suficientes.

### 4. Pruebas realizadas

Se ejecutaron 28 casos de prueba automáticos mediante `go test -v ./...`:

- 7 pruebas de validación de estados.
- 12 pruebas de transiciones.
- 5 pruebas de inicio de sesión JWT.
- 4 pruebas de autenticación y autorización.

Todos los casos finalizaron correctamente.

### 5. Documentación

Se utilizó OpenAPI 3.1 para describir los endpoints, las estructuras JSON, la autenticación JWT y las respuestas HTTP.

La documentación se generó en `docs/api.html` mediante Redocly.

### 6. Uso de inteligencia artificial

Se utilizó ChatGPT como herramienta de apoyo para comprender los requisitos, desarrollar ejemplos de código, preparar pruebas automáticas y revisar la documentación. Las soluciones fueron ejecutadas y verificadas durante el desarrollo del proyecto.

### 7. Conclusión

La API permite gestionar intercambios de libros usados con validaciones de estado, autenticación JWT y permisos por roles. Las pruebas automáticas realizadas finalizaron correctamente y la documentación HTML fue generada.