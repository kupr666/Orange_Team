# Контракт успешных ответов Workout API

Документ фиксирует JSON-ответы CRUD-эндпоинтов для сущности `workout`.

## Общие правила

- Базовый путь: `/app/v1/workouts`.
- Названия JSON-полей записываются в `snake_case`.
- Дата и время передаются строкой в формате RFC 3339, например
  `2026-08-23T10:00:00Z`.
- Если тренировка ещё не начата или не завершена, соответствующие поля времени
  содержат `null`.
- `intensity` содержит `null`, пока интенсивность не задана.
- Все ответы с телом имеют `Content-Type: application/json`.

## Объект `Workout`

Полный объект тренировки имеет следующий вид:

```json
{
  "id": 42,
  "version": 2,
  "user_id": 7,
  "status": "in_progress",
  "started_at": "2026-08-23T10:00:00Z",
  "completed_at": null,
  "created_at": "2026-08-23T09:55:00Z",
  "updated_at": "2026-08-23T10:00:00Z",
  "workout_score": 0,
  "intensity": null,
  "personal_score_coefficient": 3
}
```

| Поле | JSON-тип | Может быть `null` | Описание |
|---|---|---:|---|
| `id` | integer | нет | Идентификатор тренировки |
| `version` | integer | нет | Текущая версия записи |
| `user_id` | integer | нет | Идентификатор владельца тренировки |
| `status` | string | нет | Текущий статус тренировки |
| `started_at` | string | да | Время начала тренировки в RFC 3339 |
| `completed_at` | string | да | Время завершения тренировки в RFC 3339 |
| `created_at` | string | нет | Время создания записи в RFC 3339 |
| `updated_at` | string | нет | Время последнего изменения в RFC 3339 |
| `workout_score` | integer | нет | Количество баллов за тренировку |
| `intensity` | integer | да | Интенсивность тренировки |
| `personal_score_coefficient` | integer | нет | Персональный коэффициент начисления баллов |

Допустимые значения поля `status` фиксируются отдельным контрактом жизненного
цикла тренировки.

## Сводка ответов

| Эндпоинт | Статус | Тело ответа |
|---|---:|---|
| `POST /app/v1/workouts` | `201 Created` | `WorkoutCreated` |
| `GET /app/v1/workouts/{id}` | `200 OK` | `Workout` |
| `GET /app/v1/workouts` | `200 OK` | массив `Workout` |
| `PATCH /app/v1/workouts/{id}` | `200 OK` | обновлённый `Workout` |
| `DELETE /app/v1/workouts/{id}` | `204 No Content` | отсутствует |

## POST `/app/v1/workouts`

После создания сервер возвращает только идентификатор новой тренировки.

```http
HTTP/1.1 201 Created
Content-Type: application/json
```

```json
{
  "id": 42
}
```

### `WorkoutCreated`

| Поле | JSON-тип | Может быть `null` | Описание |
|---|---|---:|---|
| `id` | integer | нет | Идентификатор созданной тренировки |

## GET `/app/v1/workouts/{id}`

Возвращает полный объект тренировки с указанным `id`.

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

Тело ответа соответствует объекту `Workout`.

## GET `/app/v1/workouts`

Возвращает JSON-массив полных объектов `Workout`.

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
[
  {
    "id": 42,
    "version": 2,
    "user_id": 7,
    "status": "in_progress",
    "started_at": "2026-08-23T10:00:00Z",
    "completed_at": null,
    "created_at": "2026-08-23T09:55:00Z",
    "updated_at": "2026-08-23T10:00:00Z",
    "workout_score": 0,
    "intensity": null,
    "personal_score_coefficient": 3
  }
]
```

Если тренировок нет, сервер возвращает пустой массив:

```json
[]
```

## PATCH `/app/v1/workouts/{id}`

Идентификатор изменяемой тренировки передаётся в URL. После успешного изменения
сервер возвращает полный объект `Workout` с актуальными значениями полей.
Значение `version` в ответе содержит новую версию записи.

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

Тело ответа соответствует объекту `Workout`.

## DELETE `/app/v1/workouts/{id}`

После успешного удаления сервер не возвращает JSON.

```http
HTTP/1.1 204 No Content
```

У ответа отсутствует тело.

## Границы контракта

Этот документ описывает только успешные ответы, представленные на схеме.
Форматы тел запросов и ошибок должны быть зафиксированы отдельно.
