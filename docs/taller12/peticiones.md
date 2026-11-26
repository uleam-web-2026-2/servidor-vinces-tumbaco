# Semana 12 · las seis peticiones del hito 2

| Petición | Antes (taller 11) | Con las defensas | Con el usuario de la API (reto) |
| --- | --- | --- | --- |
| 1 sin token | 401 | 401 | 401 |
| 2 login | 200 | 200 | 200 |
| 3 con token | 200 | 200 | 200 |
| 4 rol sin permiso | 403 | 403 | 403 |
| 5 cambio permitido | 200 | 200 | 200 |
| 6 transición prohibida | 200 | 409 | 409 |

## Explicación del `409 Conflict` en la sexta petición

La sexta petición responde **`409 Conflict`** (en lugar de `422`) porque el formato de la petición es válido, pero la transición de estado solicitada entra en conflicto con el estado actual del intercambio en la base de datos según las reglas de negocio.
