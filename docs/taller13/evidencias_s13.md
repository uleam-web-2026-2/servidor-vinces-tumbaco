# Evidencias de ejecución — Taller 13

- **Rama de trabajo:** `taller13`
- **Paquetes verificados (`go test ./...`):**
  - `internal/avisos`: `TestEnviarSoloAlUsuarioSuscrito`, `TestQuitarConexion`, `TestConectarRechazaSinTokenOInvalido`, `TestConectarRechazaOrigenNoAutorizado`, `TestConectarRecibeAvisoYQuitaAlCerrar` (`PASS`).
  - `internal/services`: `TestCambiarEstadoEnviaAvisoSoloAParticipantes`, `TestCambiarEstadoNoEnviaAvisoSiFallaPersistencia`, `TestCambiarEstadoConParticipantesNulosNoNotifica` (`PASS`).
  - `internal/intercambios` y `internal/libros`: (`PASS`).
- **Compilación:** `go build ./...` ejecutado sobre el módulo.
- **Capturas / pruebas manuales:** adjuntar capturas de conexión WebSocket, cambio de estado y creación con rol administrador en esta carpeta (`docs/taller13/`) si el docente solicita evidencia visual adicional.
## Verificacion final en vivo y go test (2026-11-26)
- PATCH /intercambios/1 -> HTTP/1.1 200 OK
- SESION 1 (majo_user - ID 1) RECIBIO: {"id":1,"estado":"entregado","evento":"estado_cambiado"}
- SESION 2 (carlos_user - ID 2) RECIBIO: {"id":1,"estado":"entregado","evento":"estado_cambiado"}
- go test ./... -> ok en internal/avisos, internal/intercambios, internal/libros, internal/services
