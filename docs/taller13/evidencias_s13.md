# Evidencias de ejecución — Taller 13

- **Rama de trabajo:** `taller13`
- **Paquetes verificados (`go test ./...`):**
  - `internal/avisos`: `TestEnviarSoloAlUsuarioSuscrito`, `TestQuitarConexion`, `TestConectarRechazaSinTokenOInvalido`, `TestConectarRechazaOrigenNoAutorizado`, `TestConectarRecibeAvisoYQuitaAlCerrar` (`PASS`).
  - `internal/services`: `TestCambiarEstadoEnviaAvisoSoloAParticipantes`, `TestCambiarEstadoNoEnviaAvisoSiFallaPersistencia`, `TestCambiarEstadoConParticipantesNulosNoNotifica` (`PASS`).
  - `internal/intercambios` y `internal/libros`: (`PASS`).
- **Compilación:** `go build ./...` ejecutado sobre el módulo.
- **Capturas / pruebas manuales:** adjuntar capturas de conexión WebSocket y cambio de estado en esta carpeta (`docs/taller13/`) si el docente solicita evidencia visual adicional.