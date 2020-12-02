#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Get component data and set necessary variables
$component = Get-Content -Path "component.json" | ConvertFrom-Json
$rcImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-rc"

Write-Host $rcImage
# Build docker image
docker build -f docker/Dockerfile -t $rcImage .

if ($LastExitCode -ne 0 -and $env:RETRY -eq $true) {
    # if package failed and retries enabled run package script again
    Write-Host "Package failed, but retries enabled, so restarting package script again..."
    ./package.ps1
}

# Set environment variables
$env:IMAGE = $rcImage

try {
    # Workaround to remove dangling images
    docker-compose -f ./docker/docker-compose.yml down

    docker-compose -f ./docker/docker-compose.yml up -d

    Start-Sleep -Seconds 15
    Invoke-WebRequest -Uri http://localhost:8080/heartbeat

    Write-Host "The container was successfully built."
    
    # Save the result to avoid overwriting it with the "down" command below
    $exitCode = $LastExitCode 
} finally {
    docker-compose -f ./docker/docker-compose.yml down
}

# Return the exit code of the "docker-compose.yml up" command
exit $exitCode 
