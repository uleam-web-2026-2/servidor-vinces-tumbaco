Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Preparando entorno para el servidor Go  " -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# Limpiar y descargar dependencias de Go
Write-Host "`n[1/2] Sincronizando dependencias de Go..." -ForegroundColor Yellow
go mod tidy

# Verificar que exista el archivo .env
Write-Host "[2/2] Verificando archivo de configuracion..." -ForegroundColor Yellow
if (-not (Test-Path ".env")) {
    if (Test-Path ".env.example") {
        Copy-Item ".env.example" ".env"
        Write-Host " Archivo .env creado a partir de .env.example" -ForegroundColor Green
    } else {
        Write-Host " Advertencia: No se encontro .env ni .env.example" -ForegroundColor Red
    }
} else {
    Write-Host " Archivo .env detectado correctamente." -ForegroundColor Green
}

Write-Host "`n ¡Entorno preparado con exito!" -ForegroundColor Green