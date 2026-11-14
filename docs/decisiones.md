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