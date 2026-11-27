# Peticiones y pruebas Ã¢â‚¬â€ Taller 13 (WebSockets y avisos en tiempo real)

## 1. ConexiÃƒÂ³n WebSocket (`/ws`)
- **Sin token:** `GET /ws` -> `401 Unauthorized` (`token requerido`).
- **Token invÃƒÂ¡lido o expirado:** `GET /ws?token=token-invalido` -> `401 Unauthorized` (`token invalido o vencido`).
- **Origen no autorizado:** `GET /ws?token=<jwt>` con `Origin: http://evil.example.com` -> conexiÃƒÂ³n rechazada por `OriginPatterns` (`http://localhost:5173`).
- **Token HS256 vÃƒÂ¡lido + Origin permitido (`http://localhost:5173`):** `GET /ws?token=<jwt>` -> `101 Switching Protocols` y suscripciÃƒÂ³n activa por usuario y rol hasta el cierre de la conexiÃƒÂ³n.

## 2. Eventos de dominio y notificaciones
- **Cambio de estado (`IntercambiosService.CambiarEstado`):**
  1. Valida estado permitido y transiciÃƒÂ³n de estado (`pendiente` -> `aceptado` / `cancelado`, `aceptado` -> `entregado` / `cancelado`).
  2. Persiste el cambio en repositorio (`s.Repositorio.CambiarEstado(intercambio)`).
  3. Solo tras guardar sin error, envÃƒÂ­a `Aviso{ID, Estado, Evento: "estado_cambiado"}` a los participantes no nulos (`UsuarioOfreceID`, `UsuarioRecibeID`) sin duplicar ni notificar a terceros.
- **CreaciÃƒÂ³n de intercambio (`IntercambiosService.Crear` Ã¢â‚¬â€ Reto):**
  1. Valida ejemplares, disponibilidad y propiedad del ejemplar ofrecido.
  2. Persiste el intercambio y su detalle en repositorio (`s.Repositorio.Crear(intercambio, detalle)`).
  3. Solo tras persistir sin error, envÃƒÂ­a `Aviso{ID, Estado, Evento: "creado"}` ÃƒÂºnicamente a conexiones con rol `administrador` mediante `EnviarPorRol`.
## 3. Linea de base â€” 6 peticiones de Bruno (Fase 0.4)

| # | Peticion | Antes (taller 12) | Con el reto | Observacion |
|---|---|---|---|---|
| 1 | Inicio de sesion valido (POST /login) | 200 OK (verificado en flujo taller 13) | 200 OK | Emite JWT HS256 |
| 2 | Inicio de sesion repetido (POST /login) | 429 Too Many Requests en rafaga corta (registrado en windows.md) | 429 Too Many Requests | Limitador de tasa activo en rafaga |
| 3 | Recurso protegido sin token | Sin registro previo en taller 12 | 401 Unauthorized | Requiere token Bearer |
| 4 | Recurso protegido con rol sin permiso (POST /intercambios con majo_user) | 403 Forbidden (no tienes permiso para esta operacion) | 403 Forbidden | Restringido por rol/permisos |
| 5 | Registro nuevo (POST registro) | Sin registro previo en taller 12 | 201 Created / 200 OK (segun registro previo) | Crea usuario nuevo cuando no existe |
| 6 | Consulta de recurso (GET /intercambios) | 200 OK (1 intercambio devuelto en prueba en vivo) | 200 OK | Lista intercambios existentes |
