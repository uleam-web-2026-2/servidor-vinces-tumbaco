Write-Host "==========================================" -ForegroundColor Magenta
Write-Host "  Verificacion de Entregables del Hito 2 " -ForegroundColor Magenta
Write-Host "==========================================" -ForegroundColor Magenta

$listo = $true

# Check 1: Carpeta docs
if (Test-Path "docs") {
    Write-Host "[OK] Carpeta 'docs' presente." -ForegroundColor Green
} else {
    Write-Host "[ERROR] Falta la carpeta 'docs'." -ForegroundColor Red
    $listo = $false
}

# Check 2: Archivo main.go
if (Test-Path "main.go") {
    Write-Host "[OK] Archivo 'main.go' presente." -ForegroundColor Green
} else {
    Write-Host "[ERROR] Falta 'main.go' en la raiz." -ForegroundColor Red
    $listo = $false
}

# Check 3: Archivo compose.yaml
if (Test-Path "compose.yaml") {
    Write-Host "[OK] Archivo 'compose.yaml' presente." -ForegroundColor Green
} else {
    Write-Host "[ERROR] Falta 'compose.yaml'." -ForegroundColor Red
    $listo = $false
}

# Check 4: Pruebas unitarias en Go
Write-Host "`nEjecutando pruebas unitarias de Go..." -ForegroundColor Yellow
go test ./...
if ($LASTEXITCODE -eq 0) {
    Write-Host "[OK] Todas las pruebas unitarias pasaron." -ForegroundColor Green
} else {
    Write-Host "[ERROR] Algunas pruebas unitarias fallaron." -ForegroundColor Red
    $listo = $false
}

Write-Host "`n------------------------------------------"
if ($listo) {
    Write-Host " RESULTADO: ¡El proyecto cumple con los requisitos del Hito 2!" -ForegroundColor Green
} else {
    Write-Host " RESULTADO: Hay elementos pendientes por corregir." -ForegroundColor Red
}