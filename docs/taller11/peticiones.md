# Pruebas de peticiones con Docker y Docker Compose

**Proyecto:** Plataforma de Intercambio de Libros Usados  
**Rama:** taller11  
**URL base:** `http://localhost:8080`

## Resultados de las seis peticiones

| N.º | Petición | Método | URL | En el contenedor | Con Compose |
|---|---|---|---|---|---|
| 1 | Listar intercambios | GET | `/intercambios` | 200 OK | 200 OK |
| 2 | Consultar un intercambio | GET | `/intercambios/1` | 200 OK | 200 OK |
| 3 | Crear un intercambio | POST | `/intercambios` | 200 OK | 201 Created |
| 4 | Actualizar un intercambio | PUT | `/intercambios/8` | 200 OK | 200 OK |
| 5 | Cambiar estado | PATCH | `/intercambios/8` | 200 OK | 200 OK |
| 6 | Eliminar un intercambio | DELETE | `/intercambios/8` | 200 OK | 200 OK |

## Observaciones

- La conexión con PostgreSQL desde el contenedor Docker se realizó correctamente.
- Las tablas se prepararon correctamente al iniciar el servidor.
- La ruta de salud `/salud` respondió con `estado: ok`.
- Las seis peticiones de Bruno finalizaron correctamente en el contenedor.
- Con Docker Compose, las seis peticiones también finalizaron correctamente.
- Para las peticiones protegidas se utilizó autenticación mediante Bearer Token.
- Se creó el intercambio con ID 8 para realizar las pruebas de actualización, cambio de estado y eliminación.
- Docker Compose permitió ejecutar la API y PostgreSQL como servicios independientes conectados mediante la red `kema-net`.
- PostgreSQL utiliza el volumen persistente `postgres_data` para conservar los datos.