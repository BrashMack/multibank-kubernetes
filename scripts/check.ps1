$ErrorActionPreference = 'Stop'

function Assert-CommandSucceeded {
    param([string]$Message)

    if ($LASTEXITCODE -ne 0) {
        throw $Message
    }
}

$pods = kubectl get pods -n multibank -o json | ConvertFrom-Json
Assert-CommandSucceeded 'Failed to get pods.'
$failedPods = $pods.items | Where-Object {
    -not $_.metadata.deletionTimestamp -and $_.status.phase -notin @('Running', 'Succeeded')
} | ForEach-Object { $_.metadata.name }
if ($failedPods) {
    throw "Unhealthy pods found: $failedPods"
}

$services = kubectl get services -n multibank -o custom-columns=NAME:.metadata.name,TYPE:.spec.type --no-headers
Assert-CommandSucceeded 'Failed to get services.'
$invalidServices = $services -split "`n" | Where-Object { $_ -and $_ -notmatch '\s+ClusterIP$' }
if ($invalidServices) {
    throw "A non-ClusterIP service was found: $invalidServices"
}

$pvcStatus = kubectl get pvc postgres-data -n multibank -o jsonpath='{.status.phase}'
Assert-CommandSucceeded 'Failed to get PVC status.'
if ($pvcStatus -ne 'Bound') {
    throw "PVC postgres-data has status $pvcStatus"
}

kubectl wait --for=condition=Available deployment/backend -n multibank --timeout=60s
Assert-CommandSucceeded 'Backend is not available.'
kubectl wait --for=condition=Available deployment/frontend -n multibank --timeout=60s
Assert-CommandSucceeded 'Frontend is not available.'
kubectl wait --for=condition=Available deployment/adminer -n multibank --timeout=60s
Assert-CommandSucceeded 'Adminer is not available.'
kubectl wait --for=condition=Complete job/migrations -n multibank --timeout=60s
Assert-CommandSucceeded 'Migrations are not complete.'

curl.exe --fail --silent --show-error --resolve multibank.local:18080:127.0.0.1 http://multibank.local:18080/ --output NUL
Assert-CommandSucceeded 'Frontend HTTP check failed.'
curl.exe --fail --silent --show-error --resolve multibank.local:18080:127.0.0.1 http://multibank.local:18080/healthz --output NUL
Assert-CommandSucceeded 'Backend health check failed.'
curl.exe --fail --silent --show-error --resolve multibank.local:18080:127.0.0.1 http://multibank.local:18080/swagger/index.html --output NUL
Assert-CommandSucceeded 'Swagger HTTP check failed.'
curl.exe --fail --silent --show-error --resolve db.multibank.local:18080:127.0.0.1 http://db.multibank.local:18080/ --output NUL
Assert-CommandSucceeded 'Adminer HTTP check failed.'
curl.exe --fail --silent --show-error http://localhost:18080/db/ --output NUL
Assert-CommandSucceeded 'Adminer localhost route check failed.'

$loginRequest = Join-Path $PSScriptRoot 'login.json'
curl.exe --fail --silent --show-error --resolve multibank.local:18080:127.0.0.1 -H 'Content-Type: application/json' --data-binary "@$loginRequest" http://multibank.local:18080/api/v1/login --output NUL
Assert-CommandSucceeded 'Backend login check failed.'

Write-Host 'All checks passed.'
