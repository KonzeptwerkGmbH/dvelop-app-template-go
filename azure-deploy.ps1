# ---------------------------------------------------------------------------
# Azure App Service Deployment Script (PowerShell)
# ---------------------------------------------------------------------------
# Baut die App als selbstenthaltendes Linux-Binary (Assets eingebettet)
# und deployed es auf den konfigurierten Azure App Service.
#
# Voraussetzungen:
#   - Go >= 1.21 installiert
#   - Azure CLI installiert und eingeloggt (az login)
#
# Verwendung:
#   .\azure-deploy.ps1
# ---------------------------------------------------------------------------

$ErrorActionPreference = "Stop"

$ResourceGroup = "D-VELOP_DEV"
$AppName       = "konzeptwerk-saschajacobi"
$DistDir       = Join-Path $PSScriptRoot "dist"
$Binary        = Join-Path $DistDir "app"
$ZipFile       = Join-Path $DistDir "deploy.zip"

# ---------------------------------------------------------------------------
Write-Host "==> Schritt 1: Code generieren (Assets einbetten)" -ForegroundColor Cyan
go generate -tags release ./...
if ($LASTEXITCODE -ne 0) { throw "go generate fehlgeschlagen" }

Write-Host "==> Schritt 2: Linux-Binary bauen (amd64, release)" -ForegroundColor Cyan
$env:GOOS   = "linux"
$env:GOARCH = "amd64"
go build -tags release -o $Binary ./cmd/app/
if ($LASTEXITCODE -ne 0) { throw "go build fehlgeschlagen" }
$env:GOOS   = ""
$env:GOARCH = ""

$size = (Get-Item $Binary).Length / 1MB
Write-Host ("    Binary: {0:N1} MB  -> {1}" -f $size, $Binary)

Write-Host "==> Schritt 3: Deployment-ZIP erstellen" -ForegroundColor Cyan
if (Test-Path $ZipFile) { Remove-Item $ZipFile }
Compress-Archive -Path $Binary -DestinationPath $ZipFile
$zipSize = (Get-Item $ZipFile).Length / 1MB
Write-Host ("    ZIP:    {0:N1} MB  -> {1}" -f $zipSize, $ZipFile)

Write-Host "==> Schritt 4: Azure-Login prüfen" -ForegroundColor Cyan
az account show --query "{subscription:name, tenant:tenantDisplayName}" -o table
if ($LASTEXITCODE -ne 0) { throw "Nicht bei Azure eingeloggt. Bitte 'az login' ausführen." }

Write-Host "==> Schritt 5: Startup-Kommando setzen" -ForegroundColor Cyan
az webapp config set `
    --name $AppName `
    --resource-group $ResourceGroup `
    --startup-file "/home/site/wwwroot/app" `
    --query "{startupCommand:appCommandLine}" -o table
if ($LASTEXITCODE -ne 0) { throw "Startup-Kommando konnte nicht gesetzt werden" }

Write-Host "==> Schritt 6: ZIP deployen" -ForegroundColor Cyan
az webapp deploy `
    --name $AppName `
    --resource-group $ResourceGroup `
    --type zip `
    --src-path $ZipFile
if ($LASTEXITCODE -ne 0) { throw "Deployment fehlgeschlagen" }

Write-Host ""
Write-Host "==> Deployment abgeschlossen!" -ForegroundColor Green
$AppUrl = az webapp show `
    --name $AppName `
    --resource-group $ResourceGroup `
    --query "defaultHostName" -o tsv
Write-Host ("    URL: https://{0}/{1}/" -f $AppUrl, $AppName) -ForegroundColor Green
