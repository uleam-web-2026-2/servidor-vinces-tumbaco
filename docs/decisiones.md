# Decisiones tÃ©cnicas de Docker

## D1. CompilaciÃ³n estÃ¡tica

Se utilizÃ³ `CGO_ENABLED=0` durante la compilaciÃ³n de Go para generar un ejecutable estÃ¡tico, sin dependencias de CGO. Esto permite ejecutar la aplicaciÃ³n en una imagen mÃ­nima.

## D2. Imagen base Distroless

Se seleccionÃ³ `gcr.io/distroless/static-debian13:nonroot` como imagen final porque contiene Ãºnicamente los componentes necesarios para ejecutar la aplicaciÃ³n, reduciendo el tamaÃ±o de la imagen y la superficie de ataque.

## D3. EjecuciÃ³n sin privilegios

Se utiliza el usuario no root proporcionado por la imagen Distroless, con UID 65532. Esta decisiÃ³n evita ejecutar el servidor con privilegios de administrador dentro del contenedor y mejora la seguridad.

## D4. ConexiÃ³n con PostgreSQL desde Docker

Se utilizÃ³ `host.docker.internal` en la variable `DATABASE_URL` del archivo `.env.docker` para permitir que el contenedor de la API se comunique con PostgreSQL ejecutÃ¡ndose en el entorno del host.

La conexiÃ³n se verificÃ³ correctamente utilizando el puerto 5434 y la base de datos `intercambio_prueba`.

## D5. OrquestaciÃ³n con Docker Compose

Se utilizÃ³ el archivo `compose.yaml` para ejecutar la API y PostgreSQL mediante dos servicios: `api` y `db`. Esto facilita iniciar y administrar ambos contenedores con un solo comando.

## D6. Persistencia de datos

Se configurÃ³ el volumen `postgres_data` para almacenar la informaciÃ³n de PostgreSQL. Esto permite conservar los datos aunque los contenedores se detengan o se vuelvan a crear.

## D7. Red interna personalizada

Se creÃ³ la red `kema-net` para permitir la comunicaciÃ³n entre la API y PostgreSQL. La base de datos no publica su puerto al equipo anfitriÃ³n, por lo que Ãºnicamente la API expone el puerto 8080.

## D8. Orden de inicio de los servicios

Se utilizÃ³ `depends_on` con la condiciÃ³n `service_healthy` para que la API espere a que PostgreSQL estÃ© preparado antes de iniciar. AdemÃ¡s, se configurÃ³ un `healthcheck` para comprobar la disponibilidad de la base de datos.

## D9. ComparaciÃ³n de las formas de ejecuciÃ³n

Se compararon tres formas de ejecutar el proyecto:

- **go run:** permite ejecutar la aplicaciÃ³n directamente durante el desarrollo, pero requiere tener Go y las dependencias configuradas.
- **docker run:** ejecuta la API dentro de un contenedor, pero la conexiÃ³n con PostgreSQL debe configurarse por separado.
- **docker compose:** permite administrar la API y PostgreSQL conjuntamente, utilizando una red y un volumen persistente.

Docker Compose facilita la administraciÃ³n del proyecto porque permite iniciar ambos servicios mediante `docker compose up --build`.

## Semana 12 Â· la seguridad de la API

### D1. Lo que encontramos y cerramos
Se revisaron las coincidencias de la bÃºsqueda SQL (`patrones.txt`, `auditoria.txt` y `auditoria.md`): las consultas no concatenan texto del cliente directamente en el SQL sino que trabajan con parÃ¡metros o filtros parametrizados, y se reforzÃ³ la validaciÃ³n de identificadores numÃ©ricos para rechazar entradas invÃ¡lidas como `1 OR 1=1` con `400 Bad Request` en `inyeccion.txt`.

### D2. El origen que puede leer nuestras respuestas
Se configurÃ³ `AllowedOrigins` Ãºnicamente con `http://localhost:5173` y no con `*` para restringir el acceso desde el navegador solo al frontend esperado. CORS va antes de las rutas/manejadores para atender las peticiones previas (`OPTIONS`) sin exigir token; `curl` desde `http://malo.example` sÃ­ pasa porque `curl` no aplica polÃ­ticas CORS del navegador (`cors_curl_ejecuta.txt`), y al publicar el frontend en otro origen habrÃ¡ que actualizar `AllowedOrigins` al origen exacto de producciÃ³n.

### D3. El lÃ­mite del login
Se aplicÃ³ un lÃ­mite de 5 intentos por minuto por IP (`httprate.LimitBy(5, time.Minute, porIP)`) Ãºnicamente en `POST /login` para frenar intentos repetitivos de credenciales sin bloquear el resto de endpoints ni las pruebas normales del flujo. A partir del sexto intento dentro de la ventana se recibe `429 Too Many Requests` con `Retry-After: 60` (`rafaga.txt`).

### D4. DetrÃ¡s de un proxy
`ClientIPFromRemoteAddr` toma la IP desde la direcciÃ³n remota de la conexiÃ³n (`RemoteAddr`) y la normaliza con `CanonicalizeIP`. Si la API queda detrÃ¡s de un proxy inverso sin configurar adecuadamente la extracciÃ³n de IP confiable, todas las peticiones podrÃ­an verse con la IP del proxy o, si se confÃ­a ciegamente en `X-Forwarded-For`, un cliente podrÃ­a falsificar esa cabecera para evadir el lÃ­mite por IP.

## Semana 12 Â· reto: mÃ­nimo privilegio en la base

### D5. Con quÃ© usuario se conecta la API
La API se conecta con el usuario restringido `api_app` (sin atributos de superusuario, creaciÃ³n de bases ni creaciÃ³n de roles), con solo los permisos necesarios sobre los datos/secuencias de `intercambio_prueba` y sin ser dueÃ±o de las tablas. Al intentar `DROP TABLE intercambios` con `api_app` (`reto_drop.txt`), PostgreSQL rechaza la operaciÃ³n (`must be owner of table intercambios`).

### D6. QuiÃ©n crea y cambia las tablas
La creaciÃ³n y los cambios de esquema/migraciones corresponden a un usuario administrativo o de migraciones (como `postgres` o el rol propietario de esquema/migraciones en el momento de migrar/desplegar), no al usuario de tiempo de ejecuciÃ³n de la API. Si se cambia un modelo o esquema y se arranca la API con `api_app`, no debe depender de que la conexiÃ³n de runtime altere o cree tablas DDL en caliente.

### D7. Lo que esto no frena
Aunque el usuario `api_app` limita el daÃ±o en la base (por ejemplo, no puede borrar tablas ni administrar roles/bases como superusuario), no reemplaza las validaciones de aplicaciÃ³n, autorizaciÃ³n/propiedad de registros ni consultas parametrizadas: si existiera una vulnerabilidad en una consulta a la que `api_app` tiene acceso, aÃºn podrÃ­a leer o modificar datos dentro de los permisos que sÃ­ tiene otorgados.

## Taller 13 — WebSockets y avisos de intercambio

- **D1 — Identidad de usuario:** Se unifica la clave de suscripción en Avisos con el identificador de usuario extraído del JWT (id o sub) y con los IDs de participantes (UsuarioOfreceID, UsuarioRecibeID) de Intercambio, evitando suposiciones no definidas entre cuentas demo y filas de base de datos.
- **D2 — Orden de persistencia y notificación:** CambiarEstado primero guarda el estado en el repositorio y únicamente si la persistencia tiene éxito emite Aviso{ID, Estado} a los participantes válidos sin duplicados.
- **D3 — Seguridad en WebSocket:** /ws valida JWT (HS256 y JWT_SECRET de al menos 32 caracteres) antes de aceptar la conexión, restringe OriginPatterns a http://localhost:5173 y limpia la suscripción al cerrar la conexión.