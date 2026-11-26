# Peticiones y pruebas â€” Taller 13 (WebSockets y avisos en tiempo real)

## 1. ConexiÃ³n WebSocket (`/ws`)
- **Sin token:** `GET /ws` -> `401 Unauthorized` (`token requerido`).
- **Token invÃ¡lido o expirado:** `GET /ws?token=token-invalido` -> `401 Unauthorized` (`token invalido o vencido`).
- **Origen no autorizado:** `GET /ws?token=<jwt>` con `Origin: http://evil.example.com` -> conexiÃ³n rechazada por `OriginPatterns` (`http://localhost:5173`).
- **Token HS256 vÃ¡lido + Origin permitido (`http://localhost:5173`):** `GET /ws?token=<jwt>` -> `101 Switching Protocols` y suscripciÃ³n activa por usuario y rol hasta el cierre de la conexiÃ³n.

## 2. Eventos de dominio y notificaciones
- **Cambio de estado (`IntercambiosService.CambiarEstado`):**
  1. Valida estado permitido y transiciÃ³n de estado (`pendiente` -> `aceptado` / `cancelado`, `aceptado` -> `entregado` / `cancelado`).
  2. Persiste el cambio en repositorio (`s.Repositorio.CambiarEstado(intercambio)`).
  3. Solo tras guardar sin error, envÃ­a `Aviso{ID, Estado, Evento: "estado_cambiado"}` a los participantes no nulos (`UsuarioOfreceID`, `UsuarioRecibeID`) sin duplicar ni notificar a terceros.
- **CreaciÃ³n de intercambio (`IntercambiosService.Crear` â€” Reto):**
  1. Valida ejemplares, disponibilidad y propiedad del ejemplar ofrecido.
  2. Persiste el intercambio y su detalle en repositorio (`s.Repositorio.Crear(intercambio, detalle)`).
  3. Solo tras persistir sin error, envÃ­a `Aviso{ID, Estado, Evento: "creado"}` Ãºnicamente a conexiones con rol `administrador` mediante `EnviarPorRol`.
## 3. Linea de base — 6 peticiones de Bruno (Fase 0.4)

| # | Peticion | Antes (taller 12) | Observacion |
|---|---|---|---|
| 1 | Inicio de sesion valido (POST /login) | 200 OK (verificado en flujo taller 13) | Emite JWT HS256 |
| 2 | Inicio de sesion repetido (POST /login) | 429 Too Many Requests en rafaga corta (registrado en windows.md) | Limitador de tasa activo |
| 3 | Recurso protegido sin token | Sin registro previo en taller 12 | Requiere token Bearer |
| 4 | Recurso protegido con rol sin permiso (POST /intercambios con majo_user) | 403 Forbidden (
o tienes permiso para esta operacion) | Verificado en vivo |
| 5 | Registro nuevo (POST registro) | Sin registro previo en taller 12 | Gasta el registro al crear usuario |
| 6 | Consulta de recurso (GET /intercambios) | 200 OK (1 intercambio devuelto en prueba en vivo) | Lista intercambios existentes |