# Coinflip Game

Оригинальная игра Coinflip на Go + PostgreSQL с настраиваемым RTP и Provably Fair.

## Требования

- Go 1.22+
- PostgreSQL 16+
- Windows / Linux / macOS

## Установка

### 1. Клонировать репозиторий
```
git clone https://github.com/USERNAME/coinflip.git
cd coinflip
```

### 2. Установить и настроить PostgreSQL

Windows: https://www.postgresql.org/download/windows/
Linux: `sudo apt install postgresql`

Создать базу и пользователя:
```sql
CREATE USER coinflip WITH PASSWORD 'coinflip_pass';
CREATE DATABASE coinflip OWNER coinflip;
GRANT ALL PRIVILEGES ON DATABASE coinflip TO coinflip;
```

### 3. Применить миграции
```
psql -U coinflip -d coinflip -h localhost -f internal/db/migrations.sql
```

### 4. Создать тестового пользователя
```
psql -U coinflip -d coinflip -h localhost -c "INSERT INTO users (username, balance) VALUES ('demo', 1000);"
```

### 5. Настроить `.env`
```
cp .env.example .env
```
Отредактировать под свои параметры.

### 6. Установить зависимости и запустить
```
go mod tidy
go run ./cmd/server
```

Открыть: http://localhost:8080

## API

- `POST /api/play` — сделать ставку
- `GET /api/balance?user_id=1` — баланс
- `GET /api/history?user_id=1&limit=20` — история
- `GET /api/admin/rtp` — текущий RTP (нужен `X-Admin-Token`)
- `POST /api/admin/rtp` — сменить RTP

## Смена RTP на лету
```
curl -X POST -H "X-Admin-Token: your_token" \
  -H "Content-Type: application/json" \
  -d '{"value": 95.5}' \
  http://localhost:8080/api/admin/rtp
```

## Структура проекта
```
coinflip/
├── cmd/server/         # точка входа
├── internal/
│   ├── config/         # загрузка конфига
│   ├── db/             # подключение к PostgreSQL и миграции
│   ├── game/           # логика игры, RTP, Provably Fair
│   └── handlers/       # HTTP-хендлеры
└── web/client/         # фронтенд
```

## Лицензия
MIT