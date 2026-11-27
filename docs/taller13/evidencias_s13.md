# Evidencias de ejecución — Taller 13

- **Rama de trabajo:** `taller13`
- **Arquitectura de notificaciones:** El gestor de avisos en memoria funciona en una sola instancia del servidor, compartiendo una única instancia entre el endpoint WebSocket (`/ws`) y `IntercambiosService`.
- **Paquetes verificados (`go test ./...`):**
  - `internal/avisos`: `TestEnviarSoloAlUsuarioSuscrito`, `TestQuitarConexion`, `TestRetoAvisaALosQueAdministran`, `TestConectarRechazaSinTokenOInvalido`, `TestConectarRecibeAvisoYQuitaAlCerrar`, `TestConectarRechazaOrigenNoAutorizado` (`PASS`).
  - `internal/services`: `TestCambiarEstadoEnviaAvisoSoloAParticipantes`, `TestCambiarEstadoNoEnviaAvisoSiFallaPersistencia`, `TestCambiarEstadoConParticipantesNulosNoNotifica` (`PASS`).
  - `internal/intercambios` y `internal/libros`: (`PASS`).
- **Compilación:** `go build ./...` ejecutado sin errores sobre el módulo.
- **Evidencias visuales y de red en `docs/taller13/`:**
  - `aviso_navegador.png`: captura de conexión WebSocket y recepción en tiempo real de los eventos `creado` y `estado_cambiado`.
  - `ws_permitido.txt`, `ws_sin_token.txt`, `ws_otro_origen.txt`: pruebas de conexión válida, rechazo sin token y rechazo por origen no autorizado.

## Verificación final en vivo y go test (2026-11-26)
- `PATCH /intercambios/1` -> `HTTP/1.1 200 OK`
- `SESION 1 (majo_user - ID 1)` RECIBIÓ: `{"id":1,"estado":"entregado","evento":"estado_cambiado"}`
- `SESION 2 (carlos_user - ID 2)` RECIBIÓ: `{"id":1,"estado":"entregado","evento":"estado_cambiado"}`
- `go test ./...` -> `ok` en `internal/avisos`, `internal/intercambios`, `internal/libros`, `internal/services`

## Prueba adicional POST /intercambios y consulta de intercambios (2026-11-26)
- `POST /intercambios` con `majo_user` -> `HTTP/1.1 403 Forbidden` (`{"error":"no tienes permiso para esta operacion"}`)
- `SESION 1 (majo_user - ID 1)` -> NO RECIBIÓ
- `SESION 2 (carlos_user - ID 2)` -> NO RECIBIÓ
- `GET /intercambios` -> `[{"id_intercambio":1,"estado":"entregado","usuario_ofrece_id":1,"usuario_recibe_id":2,"fecha_solicitud":"2026-11-26T22:36:25.030845Z","fecha_confirmacion":null,"fecha_finalizacion":null,"observaciones":"Intercambio de libros de demostracion"}]`
