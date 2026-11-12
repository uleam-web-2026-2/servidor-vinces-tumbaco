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