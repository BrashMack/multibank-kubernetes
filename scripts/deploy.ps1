$ErrorActionPreference = 'Stop'

function Invoke-Checked {
    param(
        [scriptblock]$Command,
        [string]$Message
    )

    & $Command
    if ($LASTEXITCODE -ne 0) {
        throw $Message
    }
}

$projectRoot = Split-Path -Parent $PSScriptRoot
$clusterName = 'multibank'
$kind = Get-Command kind -ErrorAction SilentlyContinue

if (-not $kind) {
    $localKind = Join-Path $projectRoot '.tools\kind.exe'
    if (Test-Path -LiteralPath $localKind) {
        $kind = $localKind
    } else {
        throw 'kind was not found. Install kind or place kind.exe in the .tools directory.'
    }
}

Invoke-Checked { docker info *> $null } 'Docker is not available.'

Invoke-Checked { docker build -t multibank-backend:1.0.0 (Join-Path $projectRoot 'backend') } 'Backend image build failed.'
Invoke-Checked { docker build -t multibank-frontend:1.0.0 (Join-Path $projectRoot 'frontend') } 'Frontend image build failed.'

$clusters = & $kind get clusters
if ($LASTEXITCODE -ne 0) {
    throw 'Failed to get the kind cluster list.'
}
if ($clusters -notcontains $clusterName) {
    Invoke-Checked { & $kind create cluster --name $clusterName --config (Join-Path $projectRoot 'kind\cluster.yaml') } 'Failed to create the kind cluster.'
} else {
    Invoke-Checked { & $kind export kubeconfig --name $clusterName } 'Failed to export the cluster configuration.'
}

Invoke-Checked { kubectl apply -f 'https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.15.1/deploy/static/provider/kind/deploy.yaml' } 'Failed to install ingress-nginx.'
Invoke-Checked { kubectl rollout status deployment/ingress-nginx-controller -n ingress-nginx --timeout=180s } 'ingress-nginx did not become ready.'

Invoke-Checked { & $kind load docker-image multibank-backend:1.0.0 --name $clusterName } 'Failed to load the Backend image.'
Invoke-Checked { & $kind load docker-image multibank-frontend:1.0.0 --name $clusterName } 'Failed to load the Frontend image.'

$manifests = Join-Path $projectRoot 'k8s'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '00-namespace.yaml') } 'Failed to create the namespace.'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '01-config.yaml') } 'Failed to create ConfigMaps.'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '02-secret.yaml') } 'Failed to create the Secret.'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '03-postgres.yaml') } 'Failed to create PostgreSQL resources.'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '04-redis.yaml') } 'Failed to create Redis resources.'

Invoke-Checked { kubectl rollout status deployment/postgres -n multibank --timeout=180s } 'PostgreSQL did not become ready.'
Invoke-Checked { kubectl rollout status deployment/redis-approved -n multibank --timeout=180s } 'redis-approved did not become ready.'
Invoke-Checked { kubectl rollout status deployment/redis-awaiting -n multibank --timeout=180s } 'redis-awaiting did not become ready.'

Invoke-Checked { kubectl delete job migrations -n multibank --ignore-not-found=true } 'Failed to remove the previous migration Job.'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '05-migrations.yaml') } 'Failed to create the migration Job.'
Invoke-Checked { kubectl wait --for=condition=Complete job/migrations -n multibank --timeout=180s } 'Database migrations failed.'

Invoke-Checked { kubectl apply -f (Join-Path $manifests '06-backend.yaml') } 'Failed to create Backend resources.'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '07-frontend.yaml') } 'Failed to create Frontend resources.'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '08-adminer.yaml') } 'Failed to create Adminer resources.'
Invoke-Checked { kubectl apply -f (Join-Path $manifests '09-ingress.yaml') } 'Failed to create Ingress.'

Invoke-Checked { kubectl rollout status deployment/backend -n multibank --timeout=180s } 'Backend did not become ready.'
Invoke-Checked { kubectl rollout status deployment/frontend -n multibank --timeout=180s } 'Frontend did not become ready.'
Invoke-Checked { kubectl rollout status deployment/adminer -n multibank --timeout=180s } 'Adminer did not become ready.'

Invoke-Checked { kubectl get pods,services,ingress -n multibank } 'Failed to get the final resource status.'
