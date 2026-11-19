# Decisiones técnicas de Docker

## D1. Compilación estática

Se utilizó `CGO_ENABLED=0` durante la compilación de Go para generar un ejecutable estático, sin dependencias de CGO. Esto permite ejecutar la aplicación en una imagen mínima.

## D2. Imagen base Distroless

Se seleccionó `gcr.io/distroless/static-debian13:nonroot` como imagen final porque contiene únicamente los componentes necesarios para ejecutar la aplicación, reduciendo el tamaño de la imagen y la superficie de ataque.

## D3. Ejecución sin privilegios

Se utiliza el usuario no root proporcionado por la imagen Distroless, con UID 65532. Esta decisión evita ejecutar el servidor con privilegios de administrador dentro del contenedor y mejora la seguridad.

## D4. Conexión con PostgreSQL desde Docker

Se utilizó `host.docker.internal` en la variable `DATABASE_URL` del archivo `.env.docker` para permitir que el contenedor de la API se comunique con PostgreSQL ejecutándose en el entorno del host.

La conexión se verificó correctamente utilizando el puerto 5434 y la base de datos `intercambio_prueba`.

## D5. Orquestación con Docker Compose

Se utilizó el archivo `compose.yaml` para ejecutar la API y PostgreSQL mediante dos servicios: `api` y `db`. Esto facilita iniciar y administrar ambos contenedores con un solo comando.

## D6. Persistencia de datos

Se configuró el volumen `postgres_data` para almacenar la información de PostgreSQL. Esto permite conservar los datos aunque los contenedores se detengan o se vuelvan a crear.

## D7. Red interna personalizada

Se creó la red `kema-net` para permitir la comunicación entre la API y PostgreSQL. La base de datos no publica su puerto al equipo anfitrión, por lo que únicamente la API expone el puerto 8080.

## D8. Orden de inicio de los servicios

Se utilizó `depends_on` con la condición `service_healthy` para que la API espere a que PostgreSQL esté preparado antes de iniciar. Además, se configuró un `healthcheck` para comprobar la disponibilidad de la base de datos.

## D9. Comparación de las formas de ejecución

Se compararon tres formas de ejecutar el proyecto:

- **go run:** permite ejecutar la aplicación directamente durante el desarrollo, pero requiere tener Go y las dependencias configuradas.
- **docker run:** ejecuta la API dentro de un contenedor, pero la conexión con PostgreSQL debe configurarse por separado.
- **docker compose:** permite administrar la API y PostgreSQL conjuntamente, utilizando una red y un volumen persistente.

Docker Compose facilita la administración del proyecto porque permite iniciar ambos servicios mediante `docker compose up --build`.

## Semana 12 · la seguridad de la API

### D1. Lo que encontramos y cerramos
Se revisaron las coincidencias de la búsqueda SQL (`patrones.txt`, `auditoria.txt` y `auditoria.md`): las consultas no concatenan texto del cliente directamente en el SQL sino que trabajan con parámetros o filtros parametrizados, y se reforzó la validación de identificadores numéricos para rechazar entradas inválidas como `1 OR 1=1` con `400 Bad Request` en `inyeccion.txt`.

### D2. El origen que puede leer nuestras respuestas
Se configuró `AllowedOrigins` únicamente con `http://localhost:5173` y no con `*` para restringir el acceso desde el navegador solo al frontend esperado. CORS va antes de las rutas/manejadores para atender las peticiones previas (`OPTIONS`) sin exigir token; `curl` desde `http://malo.example` sí pasa porque `curl` no aplica políticas CORS del navegador (`cors_curl_ejecuta.txt`), y al publicar el frontend en otro origen habrá que actualizar `AllowedOrigins` al origen exacto de producción.

### D3. El límite del login
Se aplicó un límite de 5 intentos por minuto por IP (`httprate.LimitBy(5, time.Minute, porIP)`) únicamente en `POST /login` para frenar intentos repetitivos de credenciales sin bloquear el resto de endpoints ni las pruebas normales del flujo. A partir del sexto intento dentro de la ventana se recibe `429 Too Many Requests` con `Retry-After: 60` (`rafaga.txt`).

### D4. Detrás de un proxy
`ClientIPFromRemoteAddr` toma la IP desde la dirección remota de la conexión (`RemoteAddr`) y la normaliza con `CanonicalizeIP`. Si la API queda detrás de un proxy inverso sin configurar adecuadamente la extracción de IP confiable, todas las peticiones podrían verse con la IP del proxy o, si se confía ciegamente en `X-Forwarded-For`, un cliente podría falsificar esa cabecera para evadir el límite por IP.

## Semana 12 · reto: mínimo privilegio en la base

### D5. Con qué usuario se conecta la API
La API se conecta con el usuario restringido `api_app` (sin atributos de superusuario, creación de bases ni creación de roles), con solo los permisos necesarios sobre los datos/secuencias de `intercambio_prueba` y sin ser dueño de las tablas. Al intentar `DROP TABLE intercambios` con `api_app` (`reto_drop.txt`), PostgreSQL rechaza la operación (`must be owner of table intercambios`).

### D6. Quién crea y cambia las tablas
La creación y los cambios de esquema/migraciones corresponden a un usuario administrativo o de migraciones (como `postgres` o el rol propietario de esquema/migraciones en el momento de migrar/desplegar), no al usuario de tiempo de ejecución de la API. Si se cambia un modelo o esquema y se arranca la API con `api_app`, no debe depender de que la conexión de runtime altere o cree tablas DDL en caliente.

### D7. Lo que esto no frena
Aunque el usuario `api_app` limita el daño en la base (por ejemplo, no puede borrar tablas ni administrar roles/bases como superusuario), no reemplaza las validaciones de aplicación, autorización/propiedad de registros ni consultas parametrizadas: si existiera una vulnerabilidad en una consulta a la que `api_app` tiene acceso, aún podría leer o modificar datos dentro de los permisos que sí tiene otorgados.
