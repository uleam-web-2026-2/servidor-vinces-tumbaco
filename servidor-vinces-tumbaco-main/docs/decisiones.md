# Decisiones técnicas del proyecto

## Plataforma de Intercambio de Libros Usados

### 1. Objetivo

Implementar una API REST para gestionar intercambios de libros usados mediante Go, PostgreSQL y GORM, incorporando autenticación JWT, permisos por roles y validaciones de estados.

### 2. Tecnologías utilizadas

- **Go:** desarrollo de la API REST.
- **PostgreSQL:** almacenamiento de información.
- **GORM:** conexión y operaciones con la base de datos.
- **JWT:** autenticación y autorización.
- **OpenAPI 3.1:** documentación de endpoints.
- **Redocly:** validación y generación de documentación HTML.

### D1. Matriz de rutas y permisos

| Método | Ruta | Administrador | Usuario | Sin token |
|---|---|---|---|---|
| POST | /login | Público | Público | Permitido |
| GET | /salud | Público | Público | Permitido |
| GET | /intercambios | Permitido | Permitido | 401 |
| GET | /intercambios/{id} | Permitido | Permitido | 401 |
| POST | /intercambios | Permitido | 403 | 401 |
| PUT | /intercambios/{id} | Permitido | 403 | 401 |
| PATCH | /intercambios/{id} | Permitido | 403 | 401 |
| DELETE | /intercambios/{id} | Permitido | 403 | 401 |

### D2. Justificación de los roles

Se definieron dos roles: **administrador** y **usuario**.

El administrador puede consultar, crear, actualizar y eliminar intercambios, además de cambiar sus estados.

El usuario puede consultar los intercambios, pero no modificarlos. Esta decisión protege la integridad de los datos y evita cambios no autorizados.

La autenticación utiliza JWT. Cuando falta un token válido, el servidor responde **HTTP 401**. Cuando un usuario autenticado intenta ejecutar una operación sin permisos, responde **HTTP 403**.

### D3. Transiciones de estado y HTTP 422

Los intercambios tienen cuatro estados:

- **Pendiente:** intercambio registrado.
- **Aceptado:** intercambio aprobado.
- **Entregado:** intercambio finalizado.
- **Cancelado:** intercambio cancelado.

Las transiciones permitidas son:

| Estado actual | Nuevo estado permitido |
|---|---|
| Pendiente | Aceptado o cancelado |
| Aceptado | Entregado o cancelado |
| Entregado | Ninguno |
| Cancelado | Ninguno |

**Ejemplo de transición prohibida:**

Intentar cambiar directamente de `pendiente` a `entregado`.

El servidor devuelve **HTTP 422 (Unprocessable Content)** porque la solicitud puede tener un formato JSON correcto, pero el cambio solicitado incumple las reglas del negocio.

El cambio de estado se realiza mediante `PATCH /intercambios/{id}`.

### 3. Pruebas realizadas

Se implementaron pruebas automatizadas para:

- Validación de estados.
- Transiciones permitidas y prohibidas.
- Inicio de sesión JWT.
- Autenticación y permisos.
- Verificación de las seis rutas protegidas sin token.

Las pruebas ejecutadas hasta el momento finalizaron correctamente.

### 4. Documentación

Se documentaron las rutas de la API con OpenAPI 3.1, incluyendo autenticación JWT y respuestas HTTP 401, 403 y 422.

La validación de Redocly finalizó correctamente y se generó el archivo `docs/api.html`.

### 5. Uso de inteligencia artificial

Se utilizó ChatGPT como herramienta de apoyo para comprender los requisitos del taller, revisar el código, elaborar pruebas y mejorar la documentación. Los cambios se comprobaron mediante comandos y pruebas ejecutadas durante el desarrollo.

### 6. Conclusión

La Plataforma de Intercambio de Libros Usados incorpora autenticación JWT, control de permisos por roles y reglas para los cambios de estado. Estas decisiones contribuyen a la seguridad y consistencia de las operaciones de la API.