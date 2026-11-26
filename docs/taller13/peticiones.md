# Peticiones y pruebas — Taller 13 (WebSockets y avisos en tiempo real)

## 1. Conexión WebSocket (`/ws`)
- **Sin token:** `GET /ws` -> `401 Unauthorized` (`token requerido`).
- **Token inválido o expirado:** `GET /ws?token=token-invalido` -> `401 Unauthorized` (`token invalido o vencido`).
- **Origen no autorizado:** `GET /ws?token=<jwt>` con `Origin: http://evil.example.com` -> conexión rechazada por `OriginPatterns` (`http://localhost:5173`).
- **Token HS256 válido + Origin permitido (`http://localhost:5173`):** `GET /ws?token=<jwt>` -> `101 Switching Protocols` y suscripción activa por usuario y rol hasta el cierre de la conexión.

## 2. Eventos de dominio y notificaciones
- **Cambio de estado (`IntercambiosService.CambiarEstado`):**
  1. Valida estado permitido y transición de estado (`pendiente` -> `aceptado` / `cancelado`, `aceptado` -> `entregado` / `cancelado`).
  2. Persiste el cambio en repositorio (`s.Repositorio.CambiarEstado(intercambio)`).
  3. Solo tras guardar sin error, envía `Aviso{ID, Estado, Evento: "estado_cambiado"}` a los participantes no nulos (`UsuarioOfreceID`, `UsuarioRecibeID`) sin duplicar ni notificar a terceros.
- **Creación de intercambio (`IntercambiosService.Crear` — Reto):**
  1. Valida ejemplares, disponibilidad y propiedad del ejemplar ofrecido.
  2. Persiste el intercambio y su detalle en repositorio (`s.Repositorio.Crear(intercambio, detalle)`).
  3. Solo tras persistir sin error, envía `Aviso{ID, Estado, Evento: "creado"}` únicamente a conexiones con rol `administrador` mediante `EnviarPorRol`.