# RUS [GODanilich/hitalentGO](https://github.com/GODanilich/hitalentGO)

## Что это?

Это REST API сервис для управления организационной структурой:

- Департаменты с древовидной иерархией
- Сотрудники привязаны к департаментам
- Имеются операции CRUD

>**Данное приложение является тестовым заданием на стажировку**.

## Что использовано в разработке?

Приложение написано на **Go 1.26**.

Основные компоненты:

- **PostgreSQL 17**
- **GORM** для работы с БД
- **goose** для миграций
- **Docker** и **Docker Compose** для контейнеризации и оркестрации

## Запуск с помощью Docker-Compose

HTTP-запросы можно отправлять любым удобным инструментом: `curl`, Postman и т.д.

### Docker Compose

Запуск приложения вместе с БД и миграциями:

```bash
docker compose up --build
```

После старта сервис доступен на `http://localhost:8080`.

Сервисы в `docker-compose.yml`:

- `db` - PostgreSQL
- `migrate` - применение миграций (`goose up`)
- `app` - HTTP API

Проверка состояния:

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```json
{
  "status": "ok"
}
```

### Нативный запуск

1. Поднимите PostgreSQL и создайте БД.
2. В корне проекта необходимо создать файл `.env`:

```env
DB_URL=postgres://postgres:1@127.0.0.1:5432/hitalent?sslmode=disable
HTTP_ADDR=:8080
GORM_LOG_LEVEL=warn
```

`HTTP_ADDR` и `GORM_LOG_LEVEL` можно не задавать, есть значения по умолчанию.

3. Установите `goose` (опционально, но рекомендуется):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

4. Примените миграции:

```bash
goose -dir ./migrations postgres "postgres://postgres:1@127.0.0.1:5432/hitalent?sslmode=disable" up
```

5. Запустите API:

```bash
go run ./cmd/api
```

## Методы API

Базовый URL: `http://localhost:8080`

### 1) POST `/departments`

Создает департамент.

Тело запроса:

```json
{
  "name": "IT",
  "parent_id": null
}
```

Успешный ответ (`201 Created`):

```json
{
  "id": 1,
  "name": "IT",
  "parent_id": null,
  "created_at": "2026-03-01T12:00:00Z"
}
```

---

### 2) POST `/departments/{id}/employees`

Создает сотрудника в департаменте `{id}`.

Тело запроса:

```json
{
  "full_name": "Ivan Petrov",
  "position": "Backend Engineer",
  "hired_at": "2024-05-01"
}
```

`hired_at` опционален, формат даты строго `YYYY-MM-DD`.

Успешный ответ (`201 Created`):

```json
{
  "id": 10,
  "department_id": 1,
  "full_name": "Ivan Petrov",
  "position": "Backend Engineer",
  "hired_at": "2024-05-01T00:00:00Z",
  "created_at": "2026-03-01T12:01:00Z"
}
```

---

### 3) GET `/departments/{id}`

Возвращает древовидную структуру департамента.

Query-параметры:

- `depth` - глубина дерева, от `1` до `5` (по умолчанию `1`)
- `include_employees` - включать сотрудников (`true/false`, по умолчанию `true`)

Пример:

```bash
curl "http://localhost:8080/departments/1?depth=3&include_employees=true"
```

Ответ (`200 OK`):

```json
{
  "department": {
    "id": 1,
    "name": "Company",
    "parent_id": null,
    "created_at": "2026-03-01T12:00:00Z"
  },
  "employees": [],
  "children": [
    {
      "department": {
        "id": 2,
        "name": "IT",
        "parent_id": 1,
        "created_at": "2026-03-01T12:01:00Z"
      },
      "employees": [],
      "children": []
    }
  ]
}
```

---

### 4) PATCH `/departments/{id}`

Частично обновляет департамент.

Тело запроса (любое из полей):

```json
{
  "name": "Platform",
  "parent_id": null
}
```

Особенности:

- `name` - от 1 до 200 символов
- `parent_id` может быть `null` (сделать корневым)
- запрещено делать департамент родителем самого себя
- запрещены циклы в дереве

Успешный ответ (`200 OK`) - обновленный департамент.

---

### 5) DELETE `/departments/{id}`

Удаляет департамент.

Query-параметры:

- `mode=cascade` (по умолчанию) - удалить департамент и его поддерево
- `mode=reassign` - перед удалением перевести сотрудников удаляемого поддерева в другой департамент
- `reassign_to_department_id` - обязателен при `mode=reassign`

Примеры:

```bash
curl -X DELETE "http://localhost:8080/departments/2?mode=cascade"
```

```bash
curl -X DELETE "http://localhost:8080/departments/2?mode=reassign&reassign_to_department_id=3"
```

Успешный ответ: `204 No Content`.

Важно: `reassign_to_department_id` не может указывать на департамент внутри удаляемого поддерева.

---

### 6) GET `/health`

Проверка состояния сервиса и доступности БД.

Успешный ответ (`200 OK`):

```json
{
  "status": "ok"
}
```

Если БД недоступна (`503 Service Unavailable`):

```json
{
  "status": "fail",
  "db": "unavailable"
}
```

## Формат ошибок

Ошибки API возвращаются в едином формате:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "depth must be between 1 and 5"
  }
}
```

Возможные варианты `code`:

- `VALIDATION_ERROR` (`400`)
- `NOT_FOUND` (`404`)
- `CONFLICT` (`409`)
- `INTERNAL` (`500`)







# ENG [GODanilich/hitalentGO](https://github.com/GODanilich/hitalentGO)

## What is it?

A REST API service for managing an organizational structure:

- Departments with a tree-like hierarchy
- Employees are assigned to departments
- CRUD operations are available

> **This application was created as a test assignment for an internship.**

## Tech Stack

The application is written in **Go 1.26**.

Main components:

- **PostgreSQL 17**
- **GORM** for working with the database
- **goose** for database migrations
- **Docker** and **Docker Compose** for containerization and orchestration

## Run with Docker Compose

You can send HTTP requests with any tool you like: `curl`, Postman, etc.

### Docker Compose

Start the application together with the database and migrations:

```bash
docker compose up --build
```

After startup, the service is available at `http://localhost:8080`.

Services in `docker-compose.yml`:

- `db` — PostgreSQL
- `migrate` — applies migrations (`goose up`)
- `app` — HTTP API

Health check:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

### Run natively

1. Start PostgreSQL and create the database.
2. In the project root, create a `.env` file:

```env
DB_URL=postgres://postgres:1@127.0.0.1:5432/hitalent?sslmode=disable
HTTP_ADDR=:8080
GORM_LOG_LEVEL=warn
```

`HTTP_ADDR` and `GORM_LOG_LEVEL` are optional — defaults are provided.

3. Install `goose` (optional, but recommended):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

4. Apply migrations:

```bash
goose -dir ./migrations postgres "postgres://postgres:1@127.0.0.1:5432/hitalent?sslmode=disable" up
```

5. Run the API:

```bash
go run ./cmd/api
```

## API Methods

Base URL: `http://localhost:8080`

### 1) POST `/departments`

Creates a department.

Request body:

```json
{
  "name": "IT",
  "parent_id": null
}
```

Successful response (`201 Created`):

```json
{
  "id": 1,
  "name": "IT",
  "parent_id": null,
  "created_at": "2026-03-01T12:00:00Z"
}
```

---

### 2) POST `/departments/{id}/employees`

Creates an employee in department `{id}`.

Request body:

```json
{
  "full_name": "Ivan Petrov",
  "position": "Backend Engineer",
  "hired_at": "2024-05-01"
}
```

`hired_at` is optional. Date format must be strictly `YYYY-MM-DD`.

Successful response (`201 Created`):

```json
{
  "id": 10,
  "department_id": 1,
  "full_name": "Ivan Petrov",
  "position": "Backend Engineer",
  "hired_at": "2024-05-01T00:00:00Z",
  "created_at": "2026-03-01T12:01:00Z"
}
```

---

### 3) GET `/departments/{id}`

Returns the department tree structure.

Query parameters:

- `depth` — tree depth, from `1` to `5` (default: `1`)
- `include_employees` — include employees (`true/false`, default: `true`)

Example:

```bash
curl "http://localhost:8080/departments/1?depth=3&include_employees=true"
```

Response (`200 OK`):

```json
{
  "department": {
    "id": 1,
    "name": "Company",
    "parent_id": null,
    "created_at": "2026-03-01T12:00:00Z"
  },
  "employees": [],
  "children": [
    {
      "department": {
        "id": 2,
        "name": "IT",
        "parent_id": 1,
        "created_at": "2026-03-01T12:01:00Z"
      },
      "employees": [],
      "children": []
    }
  ]
}
```

---

### 4) PATCH `/departments/{id}`

Partially updates a department.

Request body (any of the fields):

```json
{
  "name": "Platform",
  "parent_id": null
}
```

Notes:

- `name` — 1 to 200 characters
- `parent_id` can be `null` (make it a root department)
- a department cannot be its own parent
- cycles in the tree are not allowed

Successful response (`200 OK`) — the updated department.

---

### 5) DELETE `/departments/{id}`

Deletes a department.

Query parameters:

- `mode=cascade` (default) — delete the department and its subtree
- `mode=reassign` — before deletion, reassign employees from the deleted subtree to another department
- `reassign_to_department_id` — required when `mode=reassign`

Examples:

```bash
curl -X DELETE "http://localhost:8080/departments/2?mode=cascade"
```

```bash
curl -X DELETE "http://localhost:8080/departments/2?mode=reassign&reassign_to_department_id=3"
```

Successful response: `204 No Content`.

Important: `reassign_to_department_id` cannot point to a department inside the subtree being deleted.

---

### 6) GET `/health`

Checks service status and database availability.

Successful response (`200 OK`):

```json
{
  "status": "ok"
}
```

If the database is unavailable (`503 Service Unavailable`):

```json
{
  "status": "fail",
  "db": "unavailable"
}
```

## Error Format

API errors are returned in a unified format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "depth must be between 1 and 5"
  }
}
```

Possible `code` values:

- `VALIDATION_ERROR` (`400`)
- `NOT_FOUND` (`404`)
- `CONFLICT` (`409`)
- `INTERNAL` (`500`)
