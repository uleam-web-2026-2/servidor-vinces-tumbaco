# Evidencias de ejecuciÃ³n â€” Taller 13

- **Rama de trabajo:** `taller13`
- **Paquetes verificados (`go test ./...`):**
  - `internal/avisos`: `TestEnviarSoloAlUsuarioSuscrito`, `TestQuitarConexion`, `TestConectarRechazaSinTokenOInvalido`, `TestConectarRechazaOrigenNoAutorizado`, `TestConectarRecibeAvisoYQuitaAlCerrar` (`PASS`).
  - `internal/services`: `TestCambiarEstadoEnviaAvisoSoloAParticipantes`, `TestCambiarEstadoNoEnviaAvisoSiFallaPersistencia`, `TestCambiarEstadoConParticipantesNulosNoNotifica` (`PASS`).
  - `internal/intercambios` y `internal/libros`: (`PASS`).
- **CompilaciÃ³n:** `go build ./...` ejecutado sobre el mÃ³dulo.
- **Capturas / pruebas manuales:** adjuntar capturas de conexiÃ³n WebSocket, cambio de estado y creaciÃ³n con rol administrador en esta carpeta (`docs/taller13/`) si el docente solicita evidencia visual adicional.
## Verificacion final en vivo y go test (2026-11-26)
- PATCH /intercambios/1 -> HTTP/1.1 200 OK
- SESION 1 (majo_user - ID 1) RECIBIO: {"id":1,"estado":"entregado","evento":"estado_cambiado"}
- SESION 2 (carlos_user - ID 2) RECIBIO: {"id":1,"estado":"entregado","evento":"estado_cambiado"}
- go test ./... -> ok en internal/avisos, internal/intercambios, internal/libros, internal/services

## Prueba adicional POST /intercambios y consulta de intercambios (2026-11-26)
- POST /intercambios con majo_user -> HTTP/1.1 403 Forbidden ({"error":"no tienes permiso para esta operacion"})
- SESION 1 (majo_user - ID 1) -> NO RECIBIO
- SESION 2 (carlos_user - ID 2) -> NO RECIBIO
- GET /intercambios -> [{"id_intercambio":1,"estado":"entregado","usuario_ofrece_id":1,"usuario_recibe_id":2,"fecha_solicitud":"2026-11-26T22:36:25.030845Z","fecha_confirmacion":null,"fecha_finalizacion":null,"observaciones":"Intercambio de libros de demostracion"}]