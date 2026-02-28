# API Organization

Небольшой HTTP API для управления организационной структурой: подразделения и сотрудники. Поддерживает иерархию подразделений, выбор глубины поддерева и массовые операции удаления с переносом сотрудников.

## Возможности
- Иерархия подразделений с вложенностью.
- Создание, просмотр, обновление и удаление подразделений.
- Создание сотрудников внутри подразделения.
- Получение подразделения с поддеревом заданной глубины.
- Опциональная загрузка сотрудников при просмотре подразделения.
- Миграции БД через `goose` при старте сервиса.

## Стек
- Go `1.25`.
- PostgreSQL.
- GORM.
- Goose (миграции).
- Logrus (логирование).

## Быстрый старт

### Локальный запуск
1. Поднять PostgreSQL и создать базу.
1. Установить переменные окружения (см. ниже).
1. Запустить приложение:

```bash
go run ./cmd/api
```

### Docker Compose

```bash
docker compose up --build
```

По умолчанию приложение будет доступно на `http://localhost:8080`.

## Конфигурация

Переменные окружения (значения по умолчанию в скобках):
- `DB_HOST` (`localhost`)
- `DB_PORT` (`5432`)
- `DB_USER` (`postgres`)
- `DB_PASSWORD` (`postgres`)
- `DB_NAME` (`organization`)
- `PORT` (`8080`)

В `docker-compose.yml` используется PostgreSQL на порту `5433` хоста и контейнерный порт `5432`.

## API

### Подразделения
- `POST /departments` — создать подразделение.
- `GET /departments/{id}` — получить подразделение.
- `PATCH /departments/{id}` — обновить подразделение.
- `DELETE /departments/{id}` — удалить подразделение.

Параметры `GET /departments/{id}`:
- `depth` — глубина поддерева от `1` до `5` (по умолчанию `1`).
- `include_employees` — включать сотрудников (`true` или `false`, по умолчанию `true`).

Параметры `DELETE /departments/{id}`:
- `mode` — `cascade` или `reassign` (обязательный).
- `reassign_to_department_id` — обязателен при `mode=reassign`.

### Сотрудники
- `POST /departments/{id}/employees` — создать сотрудника в подразделении.

## Примеры запросов

Создать корневое подразделение:

```bash
curl -X POST http://localhost:8080/departments \
  -H 'Content-Type: application/json' \
  -d '{"name":"Engineering"}'
```

Создать дочернее подразделение:

```bash
curl -X POST http://localhost:8080/departments \
  -H 'Content-Type: application/json' \
  -d '{"name":"Platform","parent_id":1}'
```

Получить подразделение с поддеревом глубины 3 без сотрудников:

```bash
curl 'http://localhost:8080/departments/1?depth=3&include_employees=false'
```

Обновить название подразделения:

```bash
curl -X PATCH http://localhost:8080/departments/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"R&D"}'
```

Удалить подразделение с каскадом:

```bash
curl -X DELETE 'http://localhost:8080/departments/1?mode=cascade'
```

Удалить подразделение с переносом сотрудников:

```bash
curl -X DELETE 'http://localhost:8080/departments/2?mode=reassign&reassign_to_department_id=1'
```

Создать сотрудника:

```bash
curl -X POST http://localhost:8080/departments/1/employees \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"Ivan Petrov","position":"Backend Engineer","hired_at":"2025-01-15T00:00:00Z"}'
```

## Примеры ответов

Создание подразделения (`POST /departments`):

```json
{
  "id": 1,
  "name": "Engineering",
  "parent_id": null,
  "created_at": "2025-02-28T12:00:00Z"
}
```

Получение подразделения с поддеревом и сотрудниками (`GET /departments/1?depth=2&include_employees=true`):

```json
{
  "id": 1,
  "name": "Engineering",
  "parent_id": null,
  "created_at": "2025-02-28T12:00:00Z",
  "children": [
    {
      "id": 2,
      "name": "Platform",
      "parent_id": 1,
      "created_at": "2025-02-28T12:10:00Z"
    }
  ],
  "employees": [
    {
      "id": 10,
      "department_id": 1,
      "full_name": "Ivan Petrov",
      "position": "Backend Engineer",
      "hired_at": "2025-01-15T00:00:00Z",
      "created_at": "2025-02-28T12:20:00Z"
    }
  ]
}
```

Создание сотрудника (`POST /departments/1/employees`):

```json
{
  "id": 10,
  "department_id": 1,
  "full_name": "Ivan Petrov",
  "position": "Backend Engineer",
  "hired_at": "2025-01-15T00:00:00Z",
  "created_at": "2025-02-28T12:20:00Z"
}
```

## Схема БД

`departments`:
- `id` `SERIAL` первичный ключ.
- `name` `VARCHAR(200)` не `NULL`.
- `parent_id` `INT` с `FK` на `departments(id)` и `ON DELETE CASCADE`.
- `created_at` `TIMESTAMP` с `DEFAULT NOW()`.
- Уникальность `name` в рамках одного `parent_id`.
- Уникальность `name` среди корневых подразделений (`parent_id IS NULL`).

`employees`:
- `id` `SERIAL` первичный ключ.
- `department_id` `INT` не `NULL`, `FK` на `departments(id)` с `ON DELETE CASCADE`.
- `full_name` `VARCHAR(200)` не `NULL`.
- `position` `VARCHAR(200)` не `NULL`.
- `hired_at` `DATE`, может быть `NULL`.
- `created_at` `TIMESTAMP` с `DEFAULT NOW()`.
- Индекс по `department_id`.

## Тесты

```bash
go test ./...
```

## Структура проекта
- `cmd/api` — точка входа приложения.
- `internal/app` — инициализация приложения и маршрутизация.
- `internal/handlers` — HTTP обработчики.
- `internal/service` — бизнес-логика.
- `internal/repository` — слой доступа к данным.
- `internal/models` — модели.
- `internal/db` — подключение к БД и миграции.
- `migrations` — SQL-миграции.
