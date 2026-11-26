# Peticiones y pruebas — Taller 13 (WebSockets y avisos en tiempo real)

## 1. Conexión WebSocket (`/ws`)
- **Sin token:** `GET /ws` -> `401 Unauthorized` (`token requerido`).
- **Token inválido o expirado:** `GET /ws?token=token-invalido` -> `401 Unauthorized` (`token invalido o vencido`).
- **Token HS256 válido + Origin permitido (`http://localhost:5173`):** `GET /ws?token=<jwt>` -> `101 Switching Protocols` y suscripción activa por clave de usuario hasta el cierre de la conexión.

## 2. Notificación al cambiar estado de intercambio
- Flujo en `IntercambiosService.CambiarEstado`:
  1. Valida estado permitido y transición de estado (`pendiente` -> `aceptado` / `cancelado`, `aceptado` -> `entregado` / `cancelado`).
  2. Persiste el cambio en repositorio (`s.Repositorio.CambiarEstado(intercambio)`).
  3. Solo tras guardar sin error, envía `Aviso{ID, Estado}` a los participantes no nulos (`UsuarioOfreceID`, `UsuarioRecibeID`) sin duplicar ni notificar a terceros.