# Decisiones técnicas del proyecto

## Plataforma de Intercambio de Libros Usados

### 1. Objetivo

Implementar una API REST para gestionar el intercambio de libros usados mediante el lenguaje Go y PostgreSQL, utilizando GORM para la persistencia de información. El sistema permite registrar usuarios, administrar libros y ejemplares, realizar solicitudes de intercambio y controlar sus estados.

Además, incorpora autenticación mediante JWT, permisos según roles y validaciones para mantener la integridad de las operaciones.

### 2. Tecnologías utilizadas

- **Go:** lenguaje utilizado para desarrollar la API REST.
- **PostgreSQL:** sistema de gestión de bases de datos relacional.
- **GORM:** ORM utilizado para realizar consultas y operaciones sobre PostgreSQL.
- **JWT:** mecanismo de autenticación y autorización mediante tokens.
- **OpenAPI 3.1:** especificación utilizada para documentar los endpoints de la API.
- **Redocly:** herramienta utilizada para validar y visualizar la documentación de la API.

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

Se definieron dos roles principales: **administrador** y **usuario**.

El administrador tiene permisos para consultar, crear, actualizar y eliminar intercambios, además de cambiar sus estados.

El usuario tiene permisos para consultar los intercambios, pero no puede realizar operaciones de modificación.

La autenticación se realiza mediante JWT. Cuando una solicitud protegida no incluye un token válido, el servidor responde con HTTP 401. Cuando un usuario autenticado intenta realizar una operación para la cual no tiene permisos, el servidor responde con HTTP 403.

### D3. Transiciones de estado y códigos HTTP

Los intercambios cuentan con cuatro estados:

- **Pendiente:** indica que el intercambio ha sido registrado y está esperando una respuesta.
- **Aceptado:** indica que la solicitud de intercambio ha sido aprobada.
- **Entregado:** indica que el intercambio ha finalizado.
- **Cancelado:** indica que el intercambio ha sido cancelado.

Las transiciones permitidas son:

| Estado actual | Nuevo estado permitido |
|---|---|
| Pendiente | Aceptado o cancelado |
| Aceptado | Entregado o cancelado |
| Entregado | Ninguno |
| Cancelado | Ninguno |

**Ejemplo de transición prohibida:**

Intentar cambiar directamente de `pendiente` a `entregado`.

El servidor responde con HTTP 409 (Conflict), porque la transición solicitada no cumple las reglas del negocio.

El cambio de estado se realiza mediante el endpoint:

`PATCH /intercambios/{id}`

### 3. Entidades del sistema

El proyecto contempla doce entidades principales que permiten organizar la información de los usuarios, libros, ejemplares y operaciones de intercambio.

#### 3.1. Usuario

Representa a las personas que utilizan la plataforma.

**Atributos:**
- `ID`: identificador único del usuario.
- `Nombre`: nombres del usuario.
- `Apellido`: apellidos del usuario.
- `Correo`: correo electrónico único.
- `Contrasena`: contraseña almacenada en la base de datos.
- `Telefono`: número de teléfono.
- `Direccion`: dirección del usuario.
- `Ciudad`: ciudad de residencia.
- `FechaRegistro`: fecha de registro.
- `Estado`: estado de la cuenta.
- `Rol`: tipo de usuario dentro del sistema.

**Relaciones:** puede estar relacionado con libros, ejemplares, solicitudes e intercambios.

#### 3.2. Libro

Representa la información bibliográfica de un libro publicado en la plataforma.

**Atributos:**
- `ID`: identificador único del libro.
- `Titulo`: título del libro.
- `Autor`: autor registrado como texto.
- `Genero`: género literario.
- `Estado`: condición general registrada para el libro.
- `ISBN`: código ISBN único.
- `Descripcion`: descripción del libro.
- `AnioPublicacion`: año de publicación.
- `Editorial`: editorial del libro.
- `Idioma`: idioma del libro.
- `Paginas`: cantidad de páginas.
- `Disponible`: indicador de disponibilidad.
- `UsuarioID`: identificador del usuario propietario, cuando corresponde.

**Relaciones:** se relaciona con ejemplares, autores y categorías.

#### 3.3. Autor

Representa a los autores de los libros registrados.

**Atributos:**
- `ID`: identificador único del autor.
- `Nombre`: nombre del autor.
- `Apellido`: apellido del autor.
- `Nacionalidad`: nacionalidad del autor.
- `FechaNacimiento`: fecha de nacimiento, cuando está registrada.

**Relaciones:** puede estar relacionado con varios libros mediante la entidad intermedia LibroAutor.

#### 3.4. Categoria

Representa las categorías utilizadas para clasificar los libros.

**Atributos:**
- `ID`: identificador único de la categoría.
- `Nombre`: nombre de la categoría.
- `Descripcion`: descripción de la categoría.

**Relaciones:** puede estar relacionada con varios libros mediante LibroCategoria.

#### 3.5. LibroAutor

Es una entidad intermedia que permite relacionar los libros con sus autores.

**Atributos:**
- `LibroID`: identificador del libro.
- `AutorID`: identificador del autor.
- `Rol`: función o participación del autor en el libro.

**Relaciones:** establece una relación entre Libro y Autor.

#### 3.6. LibroCategoria

Es una entidad intermedia que permite relacionar los libros con sus categorías.

**Atributos:**
- `LibroID`: identificador del libro.
- `CategoriaID`: identificador de la categoría.

**Relaciones:** establece una relación entre Libro y Categoria.

#### 3.7. EjemplarLibro

Representa una copia física de un libro que pertenece a un usuario y puede participar en un intercambio.

**Atributos:**
- `ID`: identificador único del ejemplar.
- `LibroID`: identificador del libro al que pertenece.
- `UsuarioID`: identificador del propietario.
- `EstadoFisico`: condición física del ejemplar.
- `Disponibilidad`: situación de disponibilidad.
- `FechaRegistro`: fecha de registro.
- `Observaciones`: información adicional sobre el ejemplar.

**Relaciones:** pertenece a un Libro y a un Usuario. También puede intervenir en solicitudes y detalles de intercambios.

#### 3.8. SolicitudIntercambio

Representa una solicitud realizada por un usuario para intercambiar dos ejemplares.

**Atributos:**
- `ID`: identificador único de la solicitud.
- `UsuarioSolicitanteID`: usuario que realiza la solicitud.
- `EjemplarOfrecidoID`: ejemplar que ofrece.
- `EjemplarDeseadoID`: ejemplar que desea recibir.
- `FechaSolicitud`: fecha de creación de la solicitud.
- `Estado`: estado de la solicitud.
- `Mensaje`: mensaje adicional de la solicitud.

**Relaciones:** se relaciona con Usuario y con los dos ejemplares involucrados.

#### 3.9. Intercambio

Representa la operación principal de intercambio entre usuarios.

**Atributos:**
- `ID`: identificador único del intercambio.
- `LibroID`: identificador del libro, cuando corresponde.
- `Persona`: información de la persona asociada al registro.
- `Estado`: situación actual del intercambio.
- `UsuarioOfreceID`: usuario que ofrece el ejemplar.
- `UsuarioRecibeID`: usuario que recibe el ejemplar.
- `FechaSolicitud`: fecha en que se solicita el intercambio.
- `FechaConfirmacion`: fecha de aceptación.
- `FechaFinalizacion`: fecha de finalización.
- `Observaciones`: información adicional.
- `SolicitanteID`: identificador del solicitante, cuando corresponde.
- `Creado`: fecha de creación del registro, cuando está registrada.

**Relaciones:** puede estar relacionado con Libro, Usuario y las entidades DetalleIntercambio, Entrega y Resena.

#### 3.10. DetalleIntercambio

Registra los ejemplares físicos que participan en un intercambio.

**Atributos:**
- `ID`: identificador único del detalle.
- `IntercambioID`: identificador del intercambio.
- `EjemplarEntregadoID`: ejemplar que se entrega.
- `EjemplarRecibidoID`: ejemplar que se recibe.
- `Cantidad`: cantidad de ejemplares registrados.

**Relaciones:** se relaciona con Intercambio y con dos registros de EjemplarLibro.

#### 3.11. Entrega

Registra la información relacionada con la entrega de los ejemplares de un intercambio.

**Atributos:**
- `ID`: identificador único de la entrega.
- `IntercambioID`: identificador del intercambio.
- `LugarEntrega`: lugar establecido para la entrega.
- `FechaProgramada`: fecha prevista para realizarla.
- `FechaReal`: fecha en que se realizó.
- `Estado`: situación de la entrega.

**Relaciones:** cada registro de entrega está asociado a un Intercambio.

#### 3.12. Resena

Permite registrar la valoración de los usuarios después de un intercambio.

**Atributos:**
- `ID`: identificador único de la reseña.
- `IntercambioID`: identificador del intercambio.
- `UsuarioAutorID`: usuario que realiza la reseña.
- `UsuarioEvaluadoID`: usuario que recibe la valoración.
- `Puntuacion`: puntuación asignada.
- `Comentario`: comentario de la reseña.
- `Fecha`: fecha de creación de la reseña.

**Relaciones:** se relaciona con Intercambio y con dos registros de Usuario: quien escribe y quien recibe la valoración.

### 4. Pruebas realizadas

Se ejecutaron pruebas automatizadas para comprobar:

- Validación de los estados de los intercambios.
- Transiciones permitidas y prohibidas.
- Inicio de sesión mediante JWT.
- Autenticación y permisos.
- Rechazo de solicitudes sin token.

También se realizaron comprobaciones de las rutas principales, consulta de intercambios, creación, actualización de estado y eliminación de un registro de prueba.

El comando utilizado para ejecutar las pruebas fue:

```powershell
go test ./...
```

Las pruebas existentes finalizaron correctamente después de implementar la arquitectura por capas.

### 5. Documentación de la API

El proyecto dispone de documentación de los endpoints mediante OpenAPI 3.1. Esta documentación describe las rutas, los métodos HTTP, los esquemas de solicitudes, los mecanismos de autenticación y las respuestas de la API.

También se cuenta con el archivo `docs/api.html` para visualizar la documentación.

La documentación debe mantenerse sincronizada con los modelos, los campos JSON y los códigos HTTP utilizados en la implementación.

### 6. Arquitectura por capas

El código de la entidad principal se organizó mediante una arquitectura por capas, aplicando el principio de responsabilidad única.

La estructura implementada es:

```text
main.go
internal/
├── config/
├── intercambios/
│   ├── auth.go
│   ├── handlers.go
│   ├── modelos.go
│   ├── rutas.go
│   ├── semilla.go
│   └── pruebas
├── repository/
│   └── intercambios_repository.go
├── services/
│   └── intercambios_service.go
└── libros/
```

Las responsabilidades son:

- **main.go:** punto de entrada de la aplicación. Carga la configuración, valida la conexión a PostgreSQL, prepara las entidades y levanta el servidor HTTP. También conecta las capas.
- **internal/config/:** administra la configuración y las variables de entorno.
- **internal/intercambios/:** contiene los modelos de datos, los manejadores HTTP, las rutas, la autenticación y las pruebas de la entidad principal.
- **internal/repository/:** contiene las operaciones de acceso a PostgreSQL mediante GORM.
- **internal/services/:** contiene las reglas del negocio y las validaciones de los intercambios.
- **internal/libros/:** contiene los modelos, manejadores y operaciones relacionadas con la entidad de libros.

El flujo de una solicitud es:

```text
Cliente HTTP
     |
     v
Handler
     |
     v
Service
     |
     v
Repository
     |
     v
PostgreSQL
```

El Handler recibe la petición y devuelve la respuesta HTTP. El Service valida las reglas del negocio y el Repository ejecuta las operaciones de persistencia. De esta manera, los manejadores no realizan consultas directas a la base de datos.

### 7. Uso de inteligencia artificial

Se utilizó ChatGPT como herramienta de apoyo para comprender los requisitos del taller, revisar el código, elaborar pruebas y mejorar la documentación. Los cambios se comprobaron mediante comandos y pruebas ejecutadas durante el desarrollo.

### 8. Conclusión

La Plataforma de Intercambio de Libros Usados contempla doce entidades que permiten organizar la información de usuarios, libros, autores, categorías, ejemplares, solicitudes, intercambios, entregas y reseñas.

Además, incorpora autenticación JWT, control de permisos por roles, validaciones de estados y una arquitectura organizada por capas. Esta separación facilita el mantenimiento del código, organiza las responsabilidades y mejora la estructura de la API.
## Anexo A · Semana 10 · Transacciones

### D1 · El caso
El caso de uso analizado es una operación de confirmación y cobro (o el equivalente en KEMA al aceptar y registrar un intercambio de libros). Este flujo escribe de manera consecutiva en dos tablas distintas de la base de datos (por ejemplo, actualizar el estado de la solicitud/libro y registrar el movimiento de transacción o detalle). Si la segunda escritura falla a mitad de camino por una restricción de la base de datos (como una validación `CHECK`), la primera escritura permanece guardada, dejando el sistema en un estado inconsistente.

### D2 · Opciones que consideramos
1. **Opción 1: Transacción gestionada directamente en la capa de servicio.** 
   - *Funcionamiento:* El servicio importa `gorm` o `database/sql` y abre explícitamente el bloque de transacción (`BEGIN`, `COMMIT`, `ROLLBACK`).
   - *Restricción que rompe:* Rompe el principio de Arquitectura Limpia al acoplar la lógica de negocio a detalles específicos de persistencia de bases de datos.
2. **Opción 2: Transacción encapsulada en el adaptador del repositorio.**
   - *Funcionamiento:* El servicio invoca un único método del repositorio, y es la implementación del adaptador la que envuelve las operaciones dentro de una transacción nativa.
   - *Restricción que cumple:* Mantiene el servicio totalmente agnóstico de la base de datos y permite que las pruebas unitarias sigan usando un repositorio falso (mock) sin base de datos.

### D3 · La que elegimos y por qué
Elegimos la **Opción 2 (transacción en el adaptador)** porque en la plataforma **KEMA** es fundamental proteger el desacoplamiento entre capas. Queremos que el servicio de negocio mantenga su portabilidad y que las pruebas unitarias continúen corriendo de forma rápida en memoria sin requerir una conexión real a PostgreSQL.

### D4 · La que descartamos y por qué
Descartamos la **Opción 1 (transacción en el servicio)** porque nos habría traído el problema concreto de contaminar las reglas del negocio con código acoplado a GORM y SQL, imposibilitando la ejecución de pruebas unitarias limpias con repositorios simulados.

### D5 · Con sqlc
Para el adaptador con `sqlc`, se requiere utilizar el soporte de transacciones mediante objetos de tipo `Queries` vinculados a la transacción (`WithTx`), abriendo la transacción a nivel de la conexión o pool (`db.BeginTx`). En nuestro diseño actual basado en GORM, la estructura de capas e interfaces se mantiene intacta; la lógica transaccional se aísla por completo dentro de cada adaptador correspondiente sin alterar los contratos del servicio.

### D6 · Archivos
- **Archivos que tocamos:** El archivo `docs/decisiones.md` (incorporando este Anexo A con el análisis de diseño).
- **Archivos que quedaron igual:** Las capas de servicios (`internal/services/`), las interfaces de los repositorios, las rutas, los manejadores (`handlers.go`) y las pruebas unitarias existentes con el repositorio falso, los cuales no requirieron modificaciones gracias a una correcta separación de responsabilidades.
