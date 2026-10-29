# Мультибанк в Kubernetes

Развёртывание мультибанковского приложения в локальном Kubernetes-кластере. В состав приложения входят Frontend, Backend, PostgreSQL, два экземпляра Redis, мигратор базы данных и Adminer.

## Архитектура

```text
Ingress NGINX
├── /                  → frontend:80
├── /api               → backend:8080
├── /swagger           → backend:8080
└── /db                → adminer:8080

backend
├── postgres:5432
├── redis-approved:6379
└── redis-awaiting:6379
```

Все сервисы приложения имеют тип `ClusterIP`. PostgreSQL и Redis доступны только внутри кластера. Внешний HTTP-доступ организован через Ingress.

## Состав репозитория

```text
backend/               исходный код и Dockerfile Backend
frontend/              исходный код и Dockerfile Frontend
kind/cluster.yaml      конфигурация локального кластера kind
k8s/                   Kubernetes-манифесты
scripts/deploy.ps1     автоматическое развёртывание
scripts/check.ps1      проверка состояния приложения
```

Манифесты находятся в каталоге `k8s` и пронумерованы в порядке применения:

1. `00-namespace.yaml` — пространство имён `multibank`.
2. `01-config.yaml` — конфигурация Backend и SQL-миграции.
3. `02-secret.yaml` — демонстрационные пароли и секреты.
4. `03-postgres.yaml` — PostgreSQL, постоянный том и ClusterIP Service.
5. `04-redis.yaml` — два Redis и их ClusterIP Service.
6. `05-migrations.yaml` — Job для миграций PostgreSQL.
7. `06-backend.yaml` — Backend и ClusterIP Service.
8. `07-frontend.yaml` — Frontend и ClusterIP Service.
9. `08-adminer.yaml` — Adminer и ClusterIP Service.
10. `09-ingress.yaml` — маршруты Frontend, Backend и Adminer.

## Требования

- Windows 10 или Windows 11;
- Docker Desktop с Linux-контейнерами;
- `kubectl`;
- `kind` версии 0.33.0 или новее;
- PowerShell 5.1 или новее.

`kubectl` входит в поставку Docker Desktop. `kind` можно установить отдельно или положить файл `kind.exe` в каталог `.tools` в корне проекта:

```powershell
New-Item -ItemType Directory -Force .tools
Invoke-WebRequest `
  -Uri "https://kind.sigs.k8s.io/dl/v0.33.0/kind-windows-amd64" `
  -OutFile ".tools\kind.exe"
```

Перед запуском необходимо убедиться, что Docker Desktop работает:

```powershell
docker info
kubectl version --client
```

## Быстрый запуск

Из корня репозитория выполнить:

```powershell
powershell.exe -ExecutionPolicy Bypass -File .\scripts\deploy.ps1
```

Скрипт выполняет полный цикл:

1. собирает Docker-образы Frontend и Backend;
2. создаёт кластер `multibank`;
3. устанавливает Ingress NGINX;
4. загружает локальные образы в kind;
5. применяет манифесты в необходимом порядке;
6. ожидает готовности PostgreSQL и Redis;
7. запускает миграции;
8. запускает Backend, Frontend и Adminer;
9. выводит итоговое состояние ресурсов.

Первый запуск занимает несколько минут, так как Docker загружает базовые образы.

## Ручное развёртывание

### 1. Сборка образов

```powershell
docker build -t multibank-backend:1.0.0 .\backend
docker build -t multibank-frontend:1.0.0 .\frontend
```

### 2. Создание кластера

Если `kind` установлен глобально:

```powershell
kind create cluster --name multibank --config .\kind\cluster.yaml
```

Если `kind.exe` находится в `.tools`:

```powershell
.\.tools\kind.exe create cluster --name multibank --config .\kind\cluster.yaml
```

Конфигурация публикует Ingress на портах `18080` и `18443`, чтобы не конфликтовать с локальными веб-серверами.

### 3. Установка Ingress Controller

```powershell
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.15.1/deploy/static/provider/kind/deploy.yaml
kubectl rollout status deployment/ingress-nginx-controller -n ingress-nginx --timeout=180s
```

### 4. Загрузка локальных образов в kind

```powershell
kind load docker-image multibank-backend:1.0.0 --name multibank
kind load docker-image multibank-frontend:1.0.0 --name multibank
```

Для локального файла `.tools\kind.exe` команды имеют тот же вид:

```powershell
.\.tools\kind.exe load docker-image multibank-backend:1.0.0 --name multibank
.\.tools\kind.exe load docker-image multibank-frontend:1.0.0 --name multibank
```

### 5. Namespace, конфигурация и секреты

```powershell
kubectl apply -f .\k8s\00-namespace.yaml
kubectl apply -f .\k8s\01-config.yaml
kubectl apply -f .\k8s\02-secret.yaml
```

В `02-secret.yaml` находятся только демонстрационные значения. Для реального окружения их необходимо заменить и хранить вне репозитория.

### 6. PostgreSQL и Redis

```powershell
kubectl apply -f .\k8s\03-postgres.yaml
kubectl apply -f .\k8s\04-redis.yaml

kubectl rollout status deployment/postgres -n multibank --timeout=180s
kubectl rollout status deployment/redis-approved -n multibank --timeout=180s
kubectl rollout status deployment/redis-awaiting -n multibank --timeout=180s
```

### 7. Миграции базы данных

```powershell
kubectl apply -f .\k8s\05-migrations.yaml
kubectl wait --for=condition=Complete job/migrations -n multibank --timeout=180s
kubectl logs job/migrations -n multibank
```

При повторном ручном запуске Job сначала следует удалить:

```powershell
kubectl delete job migrations -n multibank --ignore-not-found=true
kubectl apply -f .\k8s\05-migrations.yaml
```

### 8. Backend, Frontend и Adminer

```powershell
kubectl apply -f .\k8s\06-backend.yaml
kubectl apply -f .\k8s\07-frontend.yaml
kubectl apply -f .\k8s\08-adminer.yaml

kubectl rollout status deployment/backend -n multibank --timeout=180s
kubectl rollout status deployment/frontend -n multibank --timeout=180s
kubectl rollout status deployment/adminer -n multibank --timeout=180s
```

### 9. Ingress

```powershell
kubectl apply -f .\k8s\09-ingress.yaml
kubectl get ingress -n multibank
```

## Адреса приложения

После развёртывания доступны:

- Frontend: [http://localhost:18080](http://localhost:18080)
- Backend Swagger: [http://localhost:18080/swagger/index.html](http://localhost:18080/swagger/index.html)
- Adminer: [http://localhost:18080/db/](http://localhost:18080/db/)
- Backend health check: [http://localhost:18080/healthz](http://localhost:18080/healthz)

Демонстрационный пользователь Frontend:

```text
Логин: demo
Пароль: multibank
```

Данные для входа в Adminer:

```text
Движок: PostgreSQL
Сервер: postgres
Пользователь: multibank
Пароль: multibank-demo-password
База данных: multibank
```

В Adminer после входа должны отображаться таблицы `clients` и `schema_migrations`.

### Именованные адреса

Ingress также настроен для доменов `multibank.local` и `db.multibank.local`. Для их использования необходимо запустить PowerShell от имени администратора и добавить в файл `C:\Windows\System32\drivers\etc\hosts` строку:

```text
127.0.0.1 multibank.local db.multibank.local
```

После этого будут доступны:

- `http://multibank.local:18080`
- `http://multibank.local:18080/swagger/index.html`
- `http://db.multibank.local:18080`

## Проверка

Полная автоматическая проверка:

```powershell
powershell.exe -ExecutionPolicy Bypass -File .\scripts\check.ps1
```

Успешный результат заканчивается строкой:

```text
All checks passed.
```

Скрипт проверяет:

- отсутствие неработающих Pod-ов;
- тип всех сервисов `ClusterIP`;
- состояние `Bound` у тома PostgreSQL;
- готовность Backend, Frontend и Adminer;
- завершение Job миграций;
- доступность Frontend, Backend, Swagger и Adminer через Ingress;
- авторизацию демонстрационного пользователя через Backend и PostgreSQL.

Состояние ресурсов можно проверить вручную:

```powershell
kubectl get pods -n multibank
kubectl get services -n multibank
kubectl get endpoints -n multibank
kubectl get ingress -n multibank
kubectl get pvc -n multibank
```

Ожидаемый результат:

- Pod-ы `postgres`, `redis-approved`, `redis-awaiting`, `backend`, `frontend` и `adminer` находятся в состоянии `Running` и `Ready`;
- Pod миграций находится в состоянии `Completed`;
- все сервисы имеют тип `ClusterIP`;
- PVC `postgres-data` находится в состоянии `Bound`;
- Ingress содержит адреса `localhost`, `multibank.local` и `db.multibank.local`.

Проверка данных PostgreSQL:

```powershell
kubectl exec deployment/postgres -n multibank -- `
  psql -U multibank -d multibank -c "SELECT login FROM clients ORDER BY login;"
```

Проверка сохранения данных после перезапуска PostgreSQL:

```powershell
kubectl rollout restart deployment/postgres -n multibank
kubectl rollout status deployment/postgres -n multibank --timeout=180s
kubectl exec deployment/postgres -n multibank -- `
  psql -U multibank -d multibank -c "SELECT COUNT(*) FROM clients;"
```

После перезапуска таблица должна сохраниться, а запрос должен вернуть шесть пользователей.

## Особенности демонстрационного окружения

Внешние банковские API отключены параметром `Providers.OpenAPI.Disabled`. Для их работы нужны закрытые командные ключи, которые не входят в учебный репозиторий. Авторизация, PostgreSQL, Redis, Swagger и маршрутизация Ingress работают независимо от внешних API.

Функции камеры зависят от разрешений браузера и не требуются для проверки Kubernetes-конфигурации.

## Обновление приложения

После изменения исходного кода необходимо пересобрать и повторно загрузить образ:

```powershell
docker build -t multibank-backend:1.0.0 .\backend
kind load docker-image multibank-backend:1.0.0 --name multibank
kubectl rollout restart deployment/backend -n multibank
```

Для Frontend используются имя образа `multibank-frontend:1.0.0` и Deployment `frontend`.

## Удаление

Удалить всё приложение вместе с локальным кластером:

```powershell
kind delete cluster --name multibank
```

Если `kind.exe` находится в каталоге `.tools`:

```powershell
.\.tools\kind.exe delete cluster --name multibank
```
