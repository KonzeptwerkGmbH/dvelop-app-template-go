#!/bin/bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Azure App Service Deployment Script
# ---------------------------------------------------------------------------
# Baut die App als selbstenthaltendes Linux-Binary (Assets eingebettet)
# und deployed es auf den konfigurierten Azure App Service.
#
# Voraussetzungen:
#   - Go >= 1.21 installiert
#   - Azure CLI installiert und eingeloggt (az login)
#
# Verwendung:
#   ./azure-deploy.sh
# ---------------------------------------------------------------------------

RESOURCE_GROUP="D-VELOP_DEV"
APP_NAME="konzeptwerk-saschajacobi"
DIST_DIR="$(cd "$(dirname "$0")" && pwd)/dist"
BINARY="$DIST_DIR/app"
ZIP_FILE="$DIST_DIR/deploy.zip"

# ---------------------------------------------------------------------------
echo "==> Schritt 1: Code generieren (Assets einbetten)"
go generate -tags release ./...

echo "==> Schritt 2: Linux-Binary bauen (amd64, release)"
GOOS=linux GOARCH=amd64 go build -tags release \
    -o "$BINARY" \
    ./cmd/app/

echo "    Binary: $(ls -lh "$BINARY" | awk '{print $5, $9}')"

echo "==> Schritt 3: Deployment-ZIP erstellen"
rm -f "$ZIP_FILE"
zip -j "$ZIP_FILE" "$BINARY"
echo "    ZIP:    $(ls -lh "$ZIP_FILE" | awk '{print $5, $9}')"

echo "==> Schritt 4: Azure-Login prüfen"
az account show --query "{subscription:name, tenant:tenantDisplayName}" -o table

echo "==> Schritt 5: Startup-Kommando setzen"
az webapp config set \
    --name "$APP_NAME" \
    --resource-group "$RESOURCE_GROUP" \
    --startup-file "/home/site/wwwroot/app" \
    --query "{startupCommand:appCommandLine}" -o table

echo "==> Schritt 6: ZIP deployen"
az webapp deploy \
    --name "$APP_NAME" \
    --resource-group "$RESOURCE_GROUP" \
    --type zip \
    --src-path "$ZIP_FILE"

echo ""
echo "==> Deployment abgeschlossen!"
APP_URL=$(az webapp show \
    --name "$APP_NAME" \
    --resource-group "$RESOURCE_GROUP" \
    --query "defaultHostName" -o tsv)
echo "    URL: https://$APP_URL/$APP_NAME/"
