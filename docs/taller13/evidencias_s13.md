# Evidencias de ejecución — Taller 13

- **Rama de trabajo:** `taller13`
- **Paquetes verificados:**
  - `internal/avisos`: `TestEnviarSoloAlUsuarioSuscrito`, `TestQuitarConexion`, `TestConectarRechazaSinTokenOInvalido`, `TestConectarRecibeAvisoYQuitaAlCerrar` (`PASS`).
  - `internal/services`: `TestCambiarEstadoEnviaAvisoSoloAParticipantes`, `TestCambiarEstadoNoEnviaAvisoSiFallaPersistencia`, `TestCambiarEstadoConParticipantesNulosNoNotifica` (`PASS`).
  - `internal/intercambios`: `TestLogin`, `TestProteger`, `TestEstadoValido`, `TestTransicionPermitida`, `TestSinTokenLasRutasProtegidasResponden401` (`PASS`).
- **Compilación:** `go build ./...` sin errores.
- **Capturas / pruebas manuales:** adjuntar capturas de conexión WebSocket y cambio de estado en esta carpeta (`docs/taller13/`) si el docente solicita evidencia visual adicional.