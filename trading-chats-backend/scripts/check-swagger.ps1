$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

swag init -g cmd/api/main.go -o docs | Out-Host

git diff --exit-code -- docs/swagger.yaml docs/swagger.json docs/docs.go
