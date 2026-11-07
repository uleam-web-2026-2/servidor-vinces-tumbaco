# Peticiones - Taller 10

Tabla comparativa de las 6 peticiones solicitadas con sus códigos de estado antes y después:

| # | Descripción de la Petición | Estado (Antes) | Estado (Después) |
|---|---|---|---|
| 1 | Sin token (Acceso no autorizado) | 401 | 401 |
| 2 | Login (Obtención de credenciales) | 200 | 200 |
| 3 | Con token válido | 403 | 200 |
| 4 | Rol sin permiso adecuado | 200 | 403 |
| 5 | Cambio de estado permitido | 400 | 200 |
| 6 | Transición de estado prohibida | 200 | 422 |