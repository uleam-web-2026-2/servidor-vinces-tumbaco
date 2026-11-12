# Pruebas de peticiones en contenedor Docker

**Proyecto:** Plataforma de Intercambio de Libros Usados  
**Rama:** taller11  
**Entorno:** Contenedor Docker (`kema`)  
**URL base:** `http://localhost:8080`

## Resultados de las seis peticiones

| N.º | Petición | Método | URL | En el contenedor |
|---|---|---|---|---|
| 1 | Listar intercambios | GET | `/intercambios` | 200 OK |
| 2 | Consultar un intercambio | GET | `/intercambios/1` | 200 OK |
| 3 | Crear un intercambio | POST | `/intercambios` | 200 OK |
| 4 | Actualizar un intercambio | PUT | `/intercambios/8` | 200 OK |
| 5 | Cambiar estado | PATCH | `/intercambios/8` | 200 OK |
| 6 | Eliminar un intercambio | DELETE | `/intercambios/8` | 200 OK |

## Observaciones

- La conexión con PostgreSQL desde el contenedor Docker se realizó correctamente.
- Las tablas se prepararon correctamente al iniciar el servidor.
- La ruta de salud `/salud` respondió con `estado: ok`.
- Las seis peticiones de Bruno finalizaron con respuesta HTTP 200 OK.
- Para las peticiones protegidas se utilizó autenticación mediante Bearer Token.
- Se creó el intercambio con ID 8 para realizar las pruebas de actualización, cambio de estado y eliminación.