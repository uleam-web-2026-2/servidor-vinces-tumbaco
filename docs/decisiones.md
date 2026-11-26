# Decisiones tÃƒÂ©cnicas de Docker

## D1. CompilaciÃƒÂ³n estÃƒÂ¡tica

Se utilizÃƒÂ³ `CGO_ENABLED=0` durante la compilaciÃƒÂ³n de Go para generar un ejecutable estÃƒÂ¡tico, sin dependencias de CGO. Esto permite ejecutar la aplicaciÃƒÂ³n en una imagen mÃƒÂ­nima.

## D2. Imagen base Distroless

Se seleccionÃƒÂ³ `gcr.io/distroless/static-debian13:nonroot` como imagen final porque contiene ÃƒÂºnicamente los componentes necesarios para ejecutar la aplicaciÃƒÂ³n, reduciendo el tamaÃƒÂ±o de la imagen y la superficie de ataque.

## D3. EjecuciÃƒÂ³n sin privilegios

Se utiliza el usuario no root proporcionado por la imagen Distroless, con UID 65532. Esta decisiÃƒÂ³n evita ejecutar el servidor con privilegios de administrador dentro del contenedor y mejora la seguridad.

## D4. ConexiÃƒÂ³n con PostgreSQL desde Docker

Se utilizÃƒÂ³ `host.docker.internal` en la variable `DATABASE_URL` del archivo `.env.docker` para permitir que el contenedor de la API se comunique con PostgreSQL ejecutÃƒÂ¡ndose en el entorno del host.

La conexiÃƒÂ³n se verificÃƒÂ³ correctamente utilizando el puerto 5434 y la base de datos `intercambio_prueba`.

## D5. OrquestaciÃƒÂ³n con Docker Compose

Se utilizÃƒÂ³ el archivo `compose.yaml` para ejecutar la API y PostgreSQL mediante dos servicios: `api` y `db`. Esto facilita iniciar y administrar ambos contenedores con un solo comando.

## D6. Persistencia de datos

Se configurÃƒÂ³ el volumen `postgres_data` para almacenar la informaciÃƒÂ³n de PostgreSQL. Esto permite conservar los datos aunque los contenedores se detengan o se vuelvan a crear.

## D7. Red interna personalizada

Se creÃƒÂ³ la red `kema-net` para permitir la comunicaciÃƒÂ³n entre la API y PostgreSQL. La base de datos no publica su puerto al equipo anfitriÃƒÂ³n, por lo que ÃƒÂºnicamente la API expone el puerto 8080.

## D8. Orden de inicio de los servicios

Se utilizÃƒÂ³ `depends_on` con la condiciÃƒÂ³n `service_healthy` para que la API espere a que PostgreSQL estÃƒÂ© preparado antes de iniciar. AdemÃƒÂ¡s, se configurÃƒÂ³ un `healthcheck` para comprobar la disponibilidad de la base de datos.

## D9. ComparaciÃƒÂ³n de las formas de ejecuciÃƒÂ³n

Se compararon tres formas de ejecutar el proyecto:

- **go run:** permite ejecutar la aplicaciÃƒÂ³n directamente durante el desarrollo, pero requiere tener Go y las dependencias configuradas.
- **docker run:** ejecuta la API dentro de un contenedor, pero la conexiÃƒÂ³n con PostgreSQL debe configurarse por separado.
- **docker compose:** permite administrar la API y PostgreSQL conjuntamente, utilizando una red y un volumen persistente.

Docker Compose facilita la administraciÃƒÂ³n del proyecto porque permite iniciar ambos servicios mediante `docker compose up --build`.

## Semana 12 Ã‚Â· la seguridad de la API

### D1. Lo que encontramos y cerramos
Se revisaron las coincidencias de la bÃƒÂºsqueda SQL (`patrones.txt`, `auditoria.txt` y `auditoria.md`): las consultas no concatenan texto del cliente directamente en el SQL sino que trabajan con parÃƒÂ¡metros o filtros parametrizados, y se reforzÃƒÂ³ la validaciÃƒÂ³n de identificadores numÃƒÂ©ricos para rechazar entradas invÃƒÂ¡lidas como `1 OR 1=1` con `400 Bad Request` en `inyeccion.txt`.

### D2. El origen que puede leer nuestras respuestas
Se configurÃƒÂ³ `AllowedOrigins` ÃƒÂºnicamente con `http://localhost:5173` y no con `*` para restringir el acceso desde el navegador solo al frontend esperado. CORS va antes de las rutas/manejadores para atender las peticiones previas (`OPTIONS`) sin exigir token; `curl` desde `http://malo.example` sÃƒÂ­ pasa porque `curl` no aplica polÃƒÂ­ticas CORS del navegador (`cors_curl_ejecuta.txt`), y al publicar el frontend en otro origen habrÃƒÂ¡ que actualizar `AllowedOrigins` al origen exacto de producciÃƒÂ³n.

### D3. El lÃƒÂ­mite del login
Se aplicÃƒÂ³ un lÃƒÂ­mite de 5 intentos por minuto por IP (`httprate.LimitBy(5, time.Minute, porIP)`) ÃƒÂºnicamente en `POST /login` para frenar intentos repetitivos de credenciales sin bloquear el resto de endpoints ni las pruebas normales del flujo. A partir del sexto intento dentro de la ventana se recibe `429 Too Many Requests` con `Retry-After: 60` (`rafaga.txt`).

### D4. DetrÃƒÂ¡s de un proxy
`ClientIPFromRemoteAddr` toma la IP desde la direcciÃƒÂ³n remota de la conexiÃƒÂ³n (`RemoteAddr`) y la normaliza con `CanonicalizeIP`. Si la API queda detrÃƒÂ¡s de un proxy inverso sin configurar adecuadamente la extracciÃƒÂ³n de IP confiable, todas las peticiones podrÃƒÂ­an verse con la IP del proxy o, si se confÃƒÂ­a ciegamente en `X-Forwarded-For`, un cliente podrÃƒÂ­a falsificar esa cabecera para evadir el lÃƒÂ­mite por IP.

## Semana 12 Ã‚Â· reto: mÃƒÂ­nimo privilegio en la base

### D5. Con quÃƒÂ© usuario se conecta la API
La API se conecta con el usuario restringido `api_app` (sin atributos de superusuario, creaciÃƒÂ³n de bases ni creaciÃƒÂ³n de roles), con solo los permisos necesarios sobre los datos/secuencias de `intercambio_prueba` y sin ser dueÃƒÂ±o de las tablas. Al intentar `DROP TABLE intercambios` con `api_app` (`reto_drop.txt`), PostgreSQL rechaza la operaciÃƒÂ³n (`must be owner of table intercambios`).

### D6. QuiÃƒÂ©n crea y cambia las tablas
La creaciÃƒÂ³n y los cambios de esquema/migraciones corresponden a un usuario administrativo o de migraciones (como `postgres` o el rol propietario de esquema/migraciones en el momento de migrar/desplegar), no al usuario de tiempo de ejecuciÃƒÂ³n de la API. Si se cambia un modelo o esquema y se arranca la API con `api_app`, no debe depender de que la conexiÃƒÂ³n de runtime altere o cree tablas DDL en caliente.

### D7. Lo que esto no frena
Aunque el usuario `api_app` limita el daÃƒÂ±o en la base (por ejemplo, no puede borrar tablas ni administrar roles/bases como superusuario), no reemplaza las validaciones de aplicaciÃƒÂ³n, autorizaciÃƒÂ³n/propiedad de registros ni consultas parametrizadas: si existiera una vulnerabilidad en una consulta a la que `api_app` tiene acceso, aÃƒÂºn podrÃƒÂ­a leer o modificar datos dentro de los permisos que sÃƒÂ­ tiene otorgados.

## Taller 13 - WebSockets y avisos de intercambio

- **D1 - Identidad de usuario:** Se unifica la clave de suscripcion en Avisos con el identificador de usuario extraido del JWT (id o sub) y con los IDs de participantes (UsuarioOfreceID, UsuarioRecibeID) de Intercambio, evitando suposiciones no definidas entre cuentas demo y filas de base de datos.
- **D2 - Orden de persistencia y notificacion:** CambiarEstado primero guarda el estado en el repositorio y unicamente si la persistencia tiene exito emite Aviso{ID, Estado} a los participantes validos sin duplicados.
- **D3 - Seguridad en WebSocket:** /ws valida JWT (HS256 y JWT_SECRET de al menos 32 caracteres) antes de aceptar la conexion, restringe OriginPatterns a http://localhost:5173 y limpia la suscripcion al cerrar la conexion.
- **D4 - Aviso al crear para rol administrador:** Tras persistir un nuevo intercambio, el servicio envia un aviso en tiempo real unicamente a los clientes conectados con rol administrador (sin notificar a otros roles), manteniendo la persistencia previa antes de emitir el evento.
