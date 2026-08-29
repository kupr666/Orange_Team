# SoloLeveling

Go HTTP API с PostgreSQL и отдельным React/Vite-фронтендом.

## Требования

- Go 1.25.5 или новее;
- Docker с Docker Compose;
- Node.js 18+ и npm.

## 1. Настройка окружения

В корне проекта создай `.env` из примера:

```bash
cp .env.example .env
```

Заполни обязательные JWT-параметры:

```env
JWT_SECRET=<случайная строка длиной минимум 32 байта>
JWT_ISSUER=orange-team-api
JWT_AUDIENCE=orange-team-client
JWT_TTL=24h
```

Безопасный секрет можно сгенерировать командой:

```bash
openssl rand -hex 32
```

Остальные значения из `.env.example` подходят для локального запуска.

## 2. Запуск PostgreSQL и миграций

Из корня проекта выполни:

```bash
make env-up
make env-port-forward
make migrate-up
```

`env-up` запускает PostgreSQL, а `env-port-forward` открывает для локального API адрес `127.0.0.1:5432`. Если миграция запущена сразу после создания контейнера и база ещё не успела инициализироваться, подожди несколько секунд и повтори `make migrate-up`.

## 3. Запуск API

В первом терминале, из корня проекта:

```bash
make app-run
```

API будет доступен по адресу:

```text
http://127.0.0.1:5050/api/v1
```

## 4. Запуск фронтенда

Во втором терминале:

```bash
cd frontend
cp .env.example .env
npm ci
npm run dev
```

Открой в браузере:

```text
http://localhost:5173
```

В режиме разработки браузер отправляет запросы на `/api/v1`, а Vite перенаправляет их на Go API по адресу `http://127.0.0.1:5050`.

## Полная последовательность запуска

После первоначальной настройки `.env`:

```bash
# Терминал 1
make env-up
make env-port-forward
make migrate-up
make app-run
```

```bash
# Терминал 2
cd frontend
npm ci
npm run dev
```

## Частые ошибки

### `connect: connection refused` для `127.0.0.1:5432`

Не запущен контейнер, публикующий порт PostgreSQL:

```bash
make env-port-forward
```

### `required key ISSUER missing value`

В корневом `.env` не заполнены `JWT_ISSUER`, `JWT_AUDIENCE` или `JWT_TTL`. Проверь раздел «Настройка окружения» выше.

### `JWT secret must contain at least 32 bytes`

Сгенерируй более длинный `JWT_SECRET`:

```bash
openssl rand -hex 32
```

## Проверки

Тесты бэкенда:

```bash
go test ./...
```

Проверка и production-сборка фронтенда:

```bash
cd frontend
npm run build
```
