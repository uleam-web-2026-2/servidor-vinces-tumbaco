# Decisiones tÃƒÆ’Ã‚Â©cnicas de Docker

## D1. CompilaciÃƒÆ’Ã‚Â³n estÃƒÆ’Ã‚Â¡tica

Se utilizÃƒÆ’Ã‚Â³ `CGO_ENABLED=0` durante la compilaciÃƒÆ’Ã‚Â³n de Go para generar un ejecutable estÃƒÆ’Ã‚Â¡tico, sin dependencias de CGO. Esto permite ejecutar la aplicaciÃƒÆ’Ã‚Â³n en una imagen mÃƒÆ’Ã‚Â­nima.

## D2. Imagen base Distroless

Se seleccionÃƒÆ’Ã‚Â³ `gcr.io/distroless/static-debian13:nonroot` como imagen final porque contiene ÃƒÆ’Ã‚Âºnicamente los componentes necesarios para ejecutar la aplicaciÃƒÆ’Ã‚Â³n, reduciendo el tamaÃƒÆ’Ã‚Â±o de la imagen y la superficie de ataque.

## D3. EjecuciÃƒÆ’Ã‚Â³n sin privilegios

Se utiliza el usuario no root proporcionado por la imagen Distroless, con UID 65532. Esta decisiÃƒÆ’Ã‚Â³n evita ejecutar el servidor con privilegios de administrador dentro del contenedor y mejora la seguridad.

## D4. ConexiÃƒÆ’Ã‚Â³n con PostgreSQL desde Docker

Se utilizÃƒÆ’Ã‚Â³ `host.docker.internal` en la variable `DATABASE_URL` del archivo `.env.docker` para permitir que el contenedor de la API se comunique con PostgreSQL ejecutÃƒÆ’Ã‚Â¡ndose en el entorno del host.

La conexiÃƒÆ’Ã‚Â³n se verificÃƒÆ’Ã‚Â³ correctamente utilizando el puerto 5434 y la base de datos `intercambio_prueba`.

## D5. OrquestaciÃƒÆ’Ã‚Â³n con Docker Compose

Se utilizÃƒÆ’Ã‚Â³ el archivo `compose.yaml` para ejecutar la API y PostgreSQL mediante dos servicios: `api` y `db`. Esto facilita iniciar y administrar ambos contenedores con un solo comando.

## D6. Persistencia de datos

Se configurÃƒÆ’Ã‚Â³ el volumen `postgres_data` para almacenar la informaciÃƒÆ’Ã‚Â³n de PostgreSQL. Esto permite conservar los datos aunque los contenedores se detengan o se vuelvan a crear.

## D7. Red interna personalizada

Se creÃƒÆ’Ã‚Â³ la red `kema-net` para permitir la comunicaciÃƒÆ’Ã‚Â³n entre la API y PostgreSQL. La base de datos no publica su puerto al equipo anfitriÃƒÆ’Ã‚Â³n, por lo que ÃƒÆ’Ã‚Âºnicamente la API expone el puerto 8080.

## D8. Orden de inicio de los servicios

Se utilizÃƒÆ’Ã‚Â³ `depends_on` con la condiciÃƒÆ’Ã‚Â³n `service_healthy` para que la API espere a que PostgreSQL estÃƒÆ’Ã‚Â© preparado antes de iniciar. AdemÃƒÆ’Ã‚Â¡s, se configurÃƒÆ’Ã‚Â³ un `healthcheck` para comprobar la disponibilidad de la base de datos.

## D9. ComparaciÃƒÆ’Ã‚Â³n de las formas de ejecuciÃƒÆ’Ã‚Â³n

Se compararon tres formas de ejecutar el proyecto:

- **go run:** permite ejecutar la aplicaciÃƒÆ’Ã‚Â³n directamente durante el desarrollo, pero requiere tener Go y las dependencias configuradas.
- **docker run:** ejecuta la API dentro de un contenedor, pero la conexiÃƒÆ’Ã‚Â³n con PostgreSQL debe configurarse por separado.
- **docker compose:** permite administrar la API y PostgreSQL conjuntamente, utilizando una red y un volumen persistente.

Docker Compose facilita la administraciÃƒÆ’Ã‚Â³n del proyecto porque permite iniciar ambos servicios mediante `docker compose up --build`.

## Semana 12 Ãƒâ€šÃ‚Â· la seguridad de la API

### D1. Lo que encontramos y cerramos
Se revisaron las coincidencias de la bÃƒÆ’Ã‚Âºsqueda SQL (`patrones.txt`, `auditoria.txt` y `auditoria.md`): las consultas no concatenan texto del cliente directamente en el SQL sino que trabajan con parÃƒÆ’Ã‚Â¡metros o filtros parametrizados, y se reforzÃƒÆ’Ã‚Â³ la validaciÃƒÆ’Ã‚Â³n de identificadores numÃƒÆ’Ã‚Â©ricos para rechazar entradas invÃƒÆ’Ã‚Â¡lidas como `1 OR 1=1` con `400 Bad Request` en `inyeccion.txt`.

### D2. El origen que puede leer nuestras respuestas
Se configurÃƒÆ’Ã‚Â³ `AllowedOrigins` ÃƒÆ’Ã‚Âºnicamente con `http://localhost:5173` y no con `*` para restringir el acceso desde el navegador solo al frontend esperado. CORS va antes de las rutas/manejadores para atender las peticiones previas (`OPTIONS`) sin exigir token; `curl` desde `http://malo.example` sÃƒÆ’Ã‚Â­ pasa porque `curl` no aplica polÃƒÆ’Ã‚Â­ticas CORS del navegador (`cors_curl_ejecuta.txt`), y al publicar el frontend en otro origen habrÃƒÆ’Ã‚Â¡ que actualizar `AllowedOrigins` al origen exacto de producciÃƒÆ’Ã‚Â³n.

### D3. El lÃƒÆ’Ã‚Â­mite del login
Se aplicÃƒÆ’Ã‚Â³ un lÃƒÆ’Ã‚Â­mite de 5 intentos por minuto por IP (`httprate.LimitBy(5, time.Minute, porIP)`) ÃƒÆ’Ã‚Âºnicamente en `POST /login` para frenar intentos repetitivos de credenciales sin bloquear el resto de endpoints ni las pruebas normales del flujo. A partir del sexto intento dentro de la ventana se recibe `429 Too Many Requests` con `Retry-After: 60` (`rafaga.txt`).

### D4. DetrÃƒÆ’Ã‚Â¡s de un proxy
`ClientIPFromRemoteAddr` toma la IP desde la direcciÃƒÆ’Ã‚Â³n remota de la conexiÃƒÆ’Ã‚Â³n (`RemoteAddr`) y la normaliza con `CanonicalizeIP`. Si la API queda detrÃƒÆ’Ã‚Â¡s de un proxy inverso sin configurar adecuadamente la extracciÃƒÆ’Ã‚Â³n de IP confiable, todas las peticiones podrÃƒÆ’Ã‚Â­an verse con la IP del proxy o, si se confÃƒÆ’Ã‚Â­a ciegamente en `X-Forwarded-For`, un cliente podrÃƒÆ’Ã‚Â­a falsificar esa cabecera para evadir el lÃƒÆ’Ã‚Â­mite por IP.

## Semana 12 Ãƒâ€šÃ‚Â· reto: mÃƒÆ’Ã‚Â­nimo privilegio en la base

### D5. Con quÃƒÆ’Ã‚Â© usuario se conecta la API
La API se conecta con el usuario restringido `api_app` (sin atributos de superusuario, creaciÃƒÆ’Ã‚Â³n de bases ni creaciÃƒÆ’Ã‚Â³n de roles), con solo los permisos necesarios sobre los datos/secuencias de `intercambio_prueba` y sin ser dueÃƒÆ’Ã‚Â±o de las tablas. Al intentar `DROP TABLE intercambios` con `api_app` (`reto_drop.txt`), PostgreSQL rechaza la operaciÃƒÆ’Ã‚Â³n (`must be owner of table intercambios`).

### D6. QuiÃƒÆ’Ã‚Â©n crea y cambia las tablas
La creaciÃƒÆ’Ã‚Â³n y los cambios de esquema/migraciones corresponden a un usuario administrativo o de migraciones (como `postgres` o el rol propietario de esquema/migraciones en el momento de migrar/desplegar), no al usuario de tiempo de ejecuciÃƒÆ’Ã‚Â³n de la API. Si se cambia un modelo o esquema y se arranca la API con `api_app`, no debe depender de que la conexiÃƒÆ’Ã‚Â³n de runtime altere o cree tablas DDL en caliente.

### D7. Lo que esto no frena
Aunque el usuario `api_app` limita el daÃƒÆ’Ã‚Â±o en la base (por ejemplo, no puede borrar tablas ni administrar roles/bases como superusuario), no reemplaza las validaciones de aplicaciÃƒÆ’Ã‚Â³n, autorizaciÃƒÆ’Ã‚Â³n/propiedad de registros ni consultas parametrizadas: si existiera una vulnerabilidad en una consulta a la que `api_app` tiene acceso, aÃƒÆ’Ã‚Âºn podrÃƒÆ’Ã‚Â­a leer o modificar datos dentro de los permisos que sÃƒÆ’Ã‚Â­ tiene otorgados.

## Taller 13 - WebSockets y avisos de intercambio

- **D1 - Identidad de usuario:** Se unifica la clave de suscripcion en Avisos con el identificador de usuario extraido del JWT (id o sub) y con los IDs de participantes (UsuarioOfreceID, UsuarioRecibeID) de Intercambio, evitando suposiciones no definidas entre cuentas demo y filas de base de datos.
- **D2 - Orden de persistencia y notificacion:** CambiarEstado primero guarda el estado en el repositorio y unicamente si la persistencia tiene exito emite Aviso{ID, Estado} a los participantes validos sin duplicados.
- **D3 - Seguridad en WebSocket:** /ws valida JWT (HS256 y JWT_SECRET de al menos 32 caracteres) antes de aceptar la conexion, restringe OriginPatterns a http://localhost:5173 y limpia la suscripcion al cerrar la conexion.
- **D4 - Aviso al crear para rol administrador (Reto):** IntercambiosService.Crear (en internal/services/intercambios_service.go) persiste Intercambio y DetalleIntercambio en PostgreSQL y luego invoca la interfaz NotificadorCreacion sin conocer WebSocket ni el nombre del rol (dministrador). En main.go, 
otificadorCreacionRol resuelve el rol y llama a Avisos.EnviarPorRol en internal/avisos/avisos.go excluyendo el creadorID (UsuarioSolicitanteID) y omitiendo otros roles; verificado con TestRetoAvisaALosQueAdministran en docs/taller13/prueba_reto.txt.

