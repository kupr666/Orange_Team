# Telegram App — Backend Architecture & API

## 📌 Назначение

Этот документ описывает примерную архитектуру backend для Telegram Mini App со следующими фичами:

1. Аутентификация через Telegram
2. Стрики привычек
3. Тренировки
4. Лидерборд
5. Случайные задания:
   - Daily Challenge
   - Weekly Challenge

Архитектура ориентирована на **модульный монолит** на Go. Для MVP это проще в разработке и эксплуатации, чем микросервисы, но при этом код разделяется по доменным модулям так, чтобы отдельные части можно было вынести в сервисы позже.

---

## 🛠 Технологический стек

| Категория | Технология |
|---|---|
| 🔧 Контроль версий | Git |
| 🐹 Язык разработки | Go 1.23+ |
| 🌐 API | REST + WebSocket |
| 🗄 База данных | PostgreSQL 18 |
| 🔄 Миграции | golang-migrate |
| 🌍 Конфигурация | Переменные окружения (`.env`) |
| 📝 Мониторинг | Структурированные логи |
| ⚙️ Сборка | Makefile |
| 🐳 Контейнеризация | Docker + Docker Compose |
| 🔐 Аутентификация | Telegram Mini App `initData` + JWT |
| 🔑 Хеширование паролей | bcrypt — только если появятся локальные credentials |

> Для чистого Telegram Mini App пароль пользователю не нужен. Telegram передает `initData`, backend проверяет его подпись и после успешной проверки выдает свой JWT. `bcrypt` пригодится для админки, email/password-входа или других локальных credentials.

---

# 1. Общая архитектура

```text
┌──────────────────────────── TELEGRAM ────────────────────────────┐
│                                                                  │
│  Telegram Client                                                 │
│       │                                                          │
│       │ открывает Mini App                                       │
│       ▼                                                          │
│  ┌─────────────────────────────┐                                 │
│  │ Telegram Mini App Frontend  │                                 │
│  │ React / Vue / другое SPA    │                                 │
│  │                             │                                 │
│  │ Telegram.WebApp.initData    │                                 │
│  └──────────────┬──────────────┘                                 │
│                 │ HTTPS / WSS                                    │
└─────────────────│─────────────────────────────────────────────────┘
                  ▼
┌──────────────────────────────────────────────────────────────────┐
│                         GO BACKEND                               │
│                                                                  │
│  Auth                                                            │
│  Users                                                           │
│  Habits / Streaks                                                │
│  Workouts                                                        │
│  Challenges                                                      │
│  Points                                                          │
│  Leaderboard                                                     │
│  Telegram                                                        │
│  WebSocket                                                       │
│                                                                  │
└─────────────────────┬────────────────────────────────────────────┘
                      │
             ┌────────┴────────┐
             ▼                 ▼
      ┌─────────────┐     ┌──────────────┐
      │ PostgreSQL  │     │ Telegram API │
      │     18      │     │   Bot API    │
      └─────────────┘     └──────────────┘
```

---

# 2. Основные архитектурные принципы

## 2.1. Модульный монолит

Все фичи работают в одном Go-приложении, но разделены по доменам:

```text
Auth
Users
Habits
Workouts
Challenges
Points
Leaderboard
Telegram
WebSocket
```

На первом этапе не нужны Kafka, RabbitMQ и отдельные микросервисы.

---

## 2.2. Domain Events

Связь между фичами лучше строить через внутренние события.

Например:

```text
HabitCompleted
      │
      ├── обновить streak
      ├── начислить XP
      ├── обновить challenge progress
      └── обновить leaderboard
```

И:

```text
WorkoutCompleted
      │
      ├── сохранить workout stats
      ├── начислить XP
      ├── обновить challenge progress
      └── обновить leaderboard
```

На MVP event bus может быть обычным in-memory механизмом внутри Go-процесса.

---

# 3. Структура проекта

Пример структуры:

```text
.
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── telegram.go
│   │
│   ├── users/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── habits/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── streak.go
│   │   └── repository.go
│   │
│   ├── workouts/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── challenges/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── assignment.go
│   │   └── repository.go
│   │
│   ├── points/
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── leaderboard/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── telegram/
│   │   ├── client.go
│   │   └── webhook.go
│   │
│   ├── websocket/
│   │   ├── hub.go
│   │   └── client.go
│   │
│   ├── events/
│   │   ├── bus.go
│   │   └── events.go
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── logging.go
│   │   └── recovery.go
│   │
│   └── config/
│       └── config.go
│
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   └── ...
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
└── go.mod
```

---

# 4. API conventions

Базовый URL:

```text
/api/v1
```

Авторизованные запросы:

```http
Authorization: Bearer <JWT>
```

Content-Type:

```http
Content-Type: application/json
```

Стандартная ошибка:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "invalid request",
    "details": {
      "field": "name"
    }
  }
}
```

Пример кодов ошибок:

```text
UNAUTHORIZED
FORBIDDEN
NOT_FOUND
VALIDATION_ERROR
CONFLICT
INTERNAL_ERROR
TELEGRAM_AUTH_FAILED
HABIT_ALREADY_COMPLETED
WORKOUT_ALREADY_COMPLETED
CHALLENGE_EXPIRED
```

---

# 5. Фича: Authentication

## Цель

Авторизация пользователя Telegram Mini App без логина и пароля.

## Flow

```text
Telegram Client
      │
      ▼
Mini App получает initData
      │
      ▼
POST /api/v1/auth/telegram
      │
      ▼
Backend проверяет подпись Telegram
      │
      ▼
find/create user
      │
      ▼
выдает JWT
```

## Endpoint

### POST `/api/v1/auth/telegram`

Request:

```json
{
  "initData": "query_id=...&user=...&auth_date=...&hash=..."
}
```

Response:

```json
{
  "accessToken": "jwt-token",
  "expiresIn": 3600,
  "user": {
    "id": "d193388d-ec75-45fe-b013-f81ec11bb431",
    "telegramId": 123456789,
    "username": "alex",
    "firstName": "Alex",
    "lastName": null,
    "avatarUrl": "https://...",
    "timezone": "Europe/London"
  }
}
```

Backend должен:

1. распарсить `initData`;
2. проверить Telegram signature/hash;
3. проверить `auth_date`;
4. найти пользователя по `telegram_id`;
5. создать пользователя, если он новый;
6. выдать JWT.

---

## GET `/api/v1/users/me`

Response:

```json
{
  "id": "d193388d-ec75-45fe-b013-f81ec11bb431",
  "telegramId": 123456789,
  "username": "alex",
  "firstName": "Alex",
  "avatarUrl": "https://...",
  "timezone": "Europe/London",
  "totalPoints": 1340,
  "createdAt": "2026-08-20T10:00:00Z"
}
```

---

## PATCH `/api/v1/users/me`

Можно менять настройки приложения.

Request:

```json
{
  "timezone": "Europe/London"
}
```

---

# 6. Фича: Habit Streak

## Цель

Пользователь создает привычки, отмечает выполнение и получает streak.

## Основные сущности

```text
Habit
HabitCompletion
```

## Таблица `habits`

```text
id              uuid PK
user_id         uuid FK users.id
name            varchar
description     text nullable
frequency       varchar
is_active       boolean
created_at      timestamptz
updated_at      timestamptz
```

На MVP:

```text
frequency = daily
```

Позже можно добавить:

```text
daily
weekdays
custom
```

---

## Таблица `habit_completions`

```text
id                uuid PK
habit_id          uuid FK habits.id
user_id           uuid FK users.id
completion_date   date
created_at        timestamptz
```

Ограничение:

```text
UNIQUE(habit_id, completion_date)
```

Именно `habit_completions` являются source of truth для streak.

---

## GET `/api/v1/habits`

Response:

```json
{
  "items": [
    {
      "id": "habit-uuid",
      "name": "Тренировка",
      "frequency": "daily",
      "currentStreak": 8,
      "longestStreak": 21,
      "completedToday": true
    }
  ]
}
```

---

## POST `/api/v1/habits`

Request:

```json
{
  "name": "Читать 20 минут",
  "description": "Минимум 20 минут каждый день",
  "frequency": "daily"
}
```

Response:

```json
{
  "id": "habit-uuid",
  "name": "Читать 20 минут",
  "frequency": "daily",
  "currentStreak": 0,
  "completedToday": false
}
```

---

## GET `/api/v1/habits/{habitId}`

Response:

```json
{
  "id": "habit-uuid",
  "name": "Читать 20 минут",
  "description": "Минимум 20 минут каждый день",
  "frequency": "daily",
  "currentStreak": 8,
  "longestStreak": 21,
  "completedToday": true
}
```

---

## PATCH `/api/v1/habits/{habitId}`

Request:

```json
{
  "name": "Читать 30 минут",
  "isActive": true
}
```

---

## DELETE `/api/v1/habits/{habitId}`

Рекомендуется soft-delete или `is_active = false`.

Response:

```http
204 No Content
```

---

## POST `/api/v1/habits/{habitId}/completions`

Отметить привычку выполненной.

Request:

```json
{
  "date": "2026-08-20"
}
```

Если разрешено отмечать только текущий день — поле `date` можно убрать и определять дату на backend по timezone пользователя.

Response:

```json
{
  "habitId": "habit-uuid",
  "completed": true,
  "currentStreak": 9,
  "longestStreak": 21,
  "pointsEarned": 10
}
```

Внутреннее событие:

```text
HabitCompleted
```

---

## DELETE `/api/v1/habits/{habitId}/completions/{date}`

Отменить отметку.

Response:

```json
{
  "habitId": "habit-uuid",
  "completed": false,
  "currentStreak": 8
}
```

---

## GET `/api/v1/habits/{habitId}/completions`

Query:

```text
?from=2026-08-01&to=2026-08-31
```

Response:

```json
{
  "items": [
    {
      "date": "2026-08-18"
    },
    {
      "date": "2026-08-19"
    },
    {
      "date": "2026-08-20"
    }
  ]
}
```

---

## GET `/api/v1/habits/{habitId}/stats`

Response:

```json
{
  "currentStreak": 8,
  "longestStreak": 21,
  "completedDays": 44,
  "completionRate7d": 0.86,
  "completionRate30d": 0.73
}
```

---

# 7. Фича: Workouts

## Основные сущности

```text
Exercise
Workout
WorkoutExercise
WorkoutSession
WorkoutSessionExercise
```

---

## Таблица `exercises`

```text
id              uuid PK
name            varchar
description     text
type            varchar
difficulty      varchar
video_url       text nullable
created_at      timestamptz
```

---

## Таблица `workouts`

```text
id                uuid PK
title             varchar
description       text
difficulty        varchar
duration_minutes  integer
reward_points     integer
is_active         boolean
created_at        timestamptz
updated_at        timestamptz
```

---

## Таблица `workout_exercises`

```text
workout_id        uuid FK
exercise_id       uuid FK
position          integer
sets              integer
reps              integer nullable
duration_seconds  integer nullable
rest_seconds      integer
```

---

## Таблица `workout_sessions`

```text
id                uuid PK
user_id           uuid FK
workout_id        uuid FK
status            varchar
started_at        timestamptz
completed_at      timestamptz nullable
duration_seconds  integer nullable
points_earned     integer default 0
```

Status:

```text
in_progress
completed
cancelled
```

---

## GET `/api/v1/workouts`

Query example:

```text
?difficulty=beginner&limit=20&offset=0
```

Response:

```json
{
  "items": [
    {
      "id": "workout-uuid",
      "title": "Full Body 20",
      "description": "Тренировка на всё тело",
      "difficulty": "beginner",
      "durationMinutes": 20,
      "exerciseCount": 5,
      "rewardPoints": 50
    }
  ]
}
```

---

## GET `/api/v1/workouts/{workoutId}`

Response:

```json
{
  "id": "workout-uuid",
  "title": "Full Body 20",
  "difficulty": "beginner",
  "durationMinutes": 20,
  "rewardPoints": 50,
  "exercises": [
    {
      "id": "exercise-1",
      "name": "Squats",
      "position": 1,
      "sets": 3,
      "reps": 15,
      "restSeconds": 60
    },
    {
      "id": "exercise-2",
      "name": "Push-ups",
      "position": 2,
      "sets": 3,
      "reps": 10,
      "restSeconds": 60
    }
  ]
}
```

---

## POST `/api/v1/workouts/{workoutId}/start`

Response:

```json
{
  "sessionId": "session-uuid",
  "workoutId": "workout-uuid",
  "status": "in_progress",
  "startedAt": "2026-08-20T12:30:00Z"
}
```

---

## GET `/api/v1/workout-sessions/{sessionId}`

Response:

```json
{
  "id": "session-uuid",
  "workoutId": "workout-uuid",
  "status": "in_progress",
  "startedAt": "2026-08-20T12:30:00Z"
}
```

---

## POST `/api/v1/workout-sessions/{sessionId}/complete`

Request:

```json
{
  "exercises": [
    {
      "exerciseId": "exercise-1",
      "setsCompleted": 3,
      "repsCompleted": 45
    },
    {
      "exerciseId": "exercise-2",
      "setsCompleted": 3,
      "repsCompleted": 30
    }
  ]
}
```

Response:

```json
{
  "sessionId": "session-uuid",
  "status": "completed",
  "durationSeconds": 1180,
  "pointsEarned": 50,
  "completedAt": "2026-08-20T12:50:00Z"
}
```

Внутреннее событие:

```text
WorkoutCompleted
```

---

## POST `/api/v1/workout-sessions/{sessionId}/cancel`

Response:

```json
{
  "sessionId": "session-uuid",
  "status": "cancelled"
}
```

---

## GET `/api/v1/users/me/workout-sessions`

Query:

```text
?limit=20&offset=0
```

Response:

```json
{
  "items": [
    {
      "id": "session-uuid",
      "workout": {
        "id": "workout-uuid",
        "title": "Full Body 20"
      },
      "completedAt": "2026-08-20T12:50:00Z",
      "durationSeconds": 1180,
      "pointsEarned": 50
    }
  ]
}
```

---

## GET `/api/v1/users/me/workout-stats`

Response:

```json
{
  "totalWorkouts": 42,
  "totalMinutes": 980,
  "workoutsThisWeek": 4,
  "currentWorkoutStreak": 5
}
```

---

# 8. Фича: Points / XP

Leaderboard и игровые механики лучше строить через единый ledger очков.

## Таблица `point_transactions`

```text
id            uuid PK
user_id       uuid FK users.id
source_type   varchar
source_id     uuid nullable
amount        integer
created_at    timestamptz
```

Пример `source_type`:

```text
habit_completion
workout_completion
daily_challenge
weekly_challenge
streak_bonus
admin_bonus
```

---

## GET `/api/v1/users/me/points`

Response:

```json
{
  "total": 1340,
  "items": [
    {
      "id": "points-1",
      "amount": 50,
      "sourceType": "workout_completion",
      "sourceId": "session-uuid",
      "createdAt": "2026-08-20T12:50:00Z"
    },
    {
      "id": "points-2",
      "amount": 10,
      "sourceType": "habit_completion",
      "sourceId": "completion-uuid",
      "createdAt": "2026-08-20T08:10:00Z"
    }
  ]
}
```

---

# 9. Фича: Leaderboard

## Цель

Показывать рейтинг пользователей по XP.

MVP может считать leaderboard напрямую из PostgreSQL.

Пример:

```sql
SELECT
    user_id,
    SUM(amount) AS points
FROM point_transactions
WHERE created_at >= $1
  AND created_at < $2
GROUP BY user_id
ORDER BY points DESC;
```

---

## GET `/api/v1/leaderboard`

Query:

```text
?period=weekly&limit=100
```

Допустимые периоды:

```text
daily
weekly
monthly
all_time
```

Response:

```json
{
  "period": "weekly",
  "top": [
    {
      "rank": 1,
      "user": {
        "id": "user-1",
        "username": "john",
        "firstName": "John",
        "avatarUrl": "https://..."
      },
      "points": 980
    },
    {
      "rank": 2,
      "user": {
        "id": "user-2",
        "username": "max",
        "firstName": "Max"
      },
      "points": 870
    }
  ],
  "me": {
    "rank": 37,
    "points": 430
  }
}
```

---

## GET `/api/v1/leaderboard/me`

Query:

```text
?period=weekly
```

Response:

```json
{
  "rank": 37,
  "points": 430,
  "usersAbove": [
    {
      "rank": 36,
      "username": "alex2",
      "points": 440
    }
  ],
  "usersBelow": [
    {
      "rank": 38,
      "username": "bob",
      "points": 420
    }
  ]
}
```

---

# 10. Leaderboard WebSocket

WebSocket имеет смысл использовать для live-обновлений leaderboard.

Endpoint:

```text
GET /api/v1/ws
```

После подключения клиент может подписаться:

```json
{
  "type": "subscribe",
  "channel": "leaderboard:weekly"
}
```

Server event:

```json
{
  "type": "leaderboard.updated",
  "data": {
    "period": "weekly",
    "userId": "user-uuid",
    "rank": 14,
    "points": 620
  }
}
```

Также можно отправлять:

```text
challenge.updated
points.earned
workout.completed
habit.completed
```

WebSocket не должен быть обязательным для основного функционала. REST остается source of truth для клиента.

---

# 11. Фича: Daily / Weekly Challenge

## Основные сущности

```text
Challenge
UserChallenge
```

---

## Таблица `challenges`

```text
id                uuid PK
title             varchar
description       text
type              varchar
difficulty        varchar
reward_points     integer
target_type       varchar
target_value      integer nullable
is_active         boolean
created_at        timestamptz
updated_at        timestamptz
```

`type`:

```text
daily
weekly
```

Примеры `target_type`:

```text
manual
complete_workouts
complete_habits
earn_points
workout_minutes
```

Пример:

```text
title        = "4 тренировки за неделю"
type         = weekly
target_type  = complete_workouts
target_value = 4
reward       = 150
```

---

## Таблица `user_challenges`

```text
id               uuid PK
user_id          uuid FK
challenge_id     uuid FK
period_key       varchar
status           varchar
progress         integer
assigned_at      timestamptz
completed_at     timestamptz nullable
expires_at       timestamptz
```

Status:

```text
active
completed
expired
```

`period_key` позволяет не назначить пользователю два одинаковых daily/weekly задания в одном периоде.

Например:

```text
daily:2026-08-20
weekly:2026-W34
```

---

# 12. Ленивая выдача случайного задания

Для MVP не обязательно запускать job в 00:00 для каждого пользователя.

Можно выдавать challenge лениво:

```text
GET /challenges/daily
        │
        ▼
Есть user_challenge на сегодня?
        │
   ┌────┴────┐
   │         │
  Да        Нет
   │         │
   │         ▼
   │    выбрать случайный
   │    challenge
   │         │
   │         ▼
   │    INSERT user_challenge
   │         │
   └────┬────┘
        ▼
   вернуть challenge
```

То же самое для weekly.

---

## GET `/api/v1/challenges/daily`

Response:

```json
{
  "id": "user-challenge-uuid",
  "type": "daily",
  "status": "active",
  "challenge": {
    "id": "challenge-uuid",
    "title": "50 приседаний",
    "description": "Сделай 50 приседаний в течение дня",
    "difficulty": "medium",
    "rewardPoints": 25,
    "targetType": "manual",
    "targetValue": 1
  },
  "progress": {
    "current": 0,
    "target": 1
  },
  "expiresAt": "2026-08-20T23:59:59Z"
}
```

---

## GET `/api/v1/challenges/weekly`

Response:

```json
{
  "id": "user-challenge-uuid",
  "type": "weekly",
  "status": "active",
  "challenge": {
    "id": "challenge-uuid",
    "title": "4 тренировки за неделю",
    "description": "Заверши любые 4 тренировки",
    "rewardPoints": 150,
    "targetType": "complete_workouts",
    "targetValue": 4
  },
  "progress": {
    "current": 2,
    "target": 4
  },
  "expiresAt": "2026-08-23T23:59:59Z"
}
```

---

## POST `/api/v1/challenges/{userChallengeId}/complete`

Использовать только для заданий, которые пользователь подтверждает вручную.

Response:

```json
{
  "id": "user-challenge-uuid",
  "status": "completed",
  "pointsEarned": 25,
  "completedAt": "2026-08-20T13:05:00Z"
}
```

---

## Автоматический challenge progress

Если задание:

```text
"Сделай 4 тренировки за неделю"
```

frontend не должен отправлять `complete`.

Backend обновляет progress автоматически:

```text
WorkoutCompleted
      │
      ▼
ChallengeProgressHandler
      │
      ├── progress 2/4
      ├── progress 3/4
      └── progress 4/4
               │
               ▼
        ChallengeCompleted
               │
               ▼
         начислить XP
```

Аналогично можно сделать:

```text
HabitCompleted
PointsEarned
WorkoutCompleted
```

---

# 13. API Summary

## Auth

```text
POST   /api/v1/auth/telegram
```

## User

```text
GET    /api/v1/users/me
PATCH  /api/v1/users/me
GET    /api/v1/users/me/points
```

## Habits

```text
GET    /api/v1/habits
POST   /api/v1/habits

GET    /api/v1/habits/{habitId}
PATCH  /api/v1/habits/{habitId}
DELETE /api/v1/habits/{habitId}

POST   /api/v1/habits/{habitId}/completions
DELETE /api/v1/habits/{habitId}/completions/{date}

GET    /api/v1/habits/{habitId}/completions
GET    /api/v1/habits/{habitId}/stats
```

## Workouts

```text
GET    /api/v1/workouts
GET    /api/v1/workouts/{workoutId}

POST   /api/v1/workouts/{workoutId}/start

GET    /api/v1/workout-sessions/{sessionId}
POST   /api/v1/workout-sessions/{sessionId}/complete
POST   /api/v1/workout-sessions/{sessionId}/cancel

GET    /api/v1/users/me/workout-sessions
GET    /api/v1/users/me/workout-stats
```

## Challenges

```text
GET    /api/v1/challenges/daily
GET    /api/v1/challenges/weekly

POST   /api/v1/challenges/{userChallengeId}/complete
```

## Leaderboard

```text
GET    /api/v1/leaderboard
GET    /api/v1/leaderboard/me
```

## WebSocket

```text
GET    /api/v1/ws
```

---

# 14. Схема связей между фичами

```text
                        ┌──────────────┐
                        │     User     │
                        └──────┬───────┘
                               │
             ┌─────────────────┼─────────────────┐
             │                 │                 │
             ▼                 ▼                 ▼
        ┌─────────┐      ┌──────────┐      ┌───────────┐
        │ Habits  │      │ Workouts │      │ Challenges│
        └────┬────┘      └────┬─────┘      └─────┬─────┘
             │                │                  │
             ▼                ▼                  │
      HabitCompleted    WorkoutCompleted         │
             │                │                  │
             └────────┬───────┴──────────┬───────┘
                      ▼                  ▼
                ┌──────────┐       ┌────────────┐
                │  Points  │       │ Challenge  │
                │   / XP   │       │  Progress  │
                └────┬─────┘       └────────────┘
                     │
                     ▼
              ┌─────────────┐
              │ Leaderboard │
              └─────────────┘
```

---

# 15. Базовые XP rewards

Для MVP можно начать с простой системы:

| Действие | XP |
|---|---:|
| Выполнение привычки | +10 |
| Завершение тренировки | +50 |
| Daily challenge | +25 |
| Weekly challenge | +150 |
| 7-day streak bonus | +50 |
| 30-day streak bonus | +250 |

Значения лучше держать в конфигурации или в БД, а не хардкодить по всему проекту.

---

# 16. PostgreSQL — основные таблицы

Минимальный набор:

```text
users

habits
habit_completions

exercises
workouts
workout_exercises
workout_sessions
workout_session_exercises

challenges
user_challenges

point_transactions
```

Связи:

```text
users
 │
 ├── habits
 │      └── habit_completions
 │
 ├── workout_sessions
 │      └── workout_session_exercises
 │
 ├── user_challenges
 │
 └── point_transactions
```

---

# 17. Индексы

Рекомендуемые индексы:

```sql
CREATE UNIQUE INDEX users_telegram_id_idx
ON users (telegram_id);

CREATE INDEX habits_user_id_idx
ON habits (user_id);

CREATE UNIQUE INDEX habit_completions_unique_idx
ON habit_completions (habit_id, completion_date);

CREATE INDEX workout_sessions_user_completed_idx
ON workout_sessions (user_id, completed_at DESC);

CREATE INDEX user_challenges_user_period_idx
ON user_challenges (user_id, period_key);

CREATE INDEX point_transactions_user_created_idx
ON point_transactions (user_id, created_at DESC);

CREATE INDEX point_transactions_created_idx
ON point_transactions (created_at);
```

---

# 18. Переменные окружения

`.env.example`:

```env
APP_ENV=development
APP_PORT=8080

DATABASE_URL=postgres://postgres:postgres@postgres:5432/telegram_app?sslmode=disable

JWT_SECRET=change-me
JWT_TTL=1h

TELEGRAM_BOT_TOKEN=change-me

LOG_LEVEL=debug
```

> `TELEGRAM_BOT_TOKEN` нельзя передавать во frontend.

---

# 19. Docker Compose

Минимальная инфраструктура:

```text
docker-compose
│
├── api
└── postgres
```

Redis на MVP не обязателен.

Позже его можно добавить для:

```text
leaderboard cache
rate limiting
distributed locks
WebSocket pub/sub
background jobs
```

---

# 20. Makefile

Пример команд:

```text
make build
make run
make test

make docker-up
make docker-down

make migrate-up
make migrate-down
make migrate-create name=create_users

make lint
```

---

# 21. Логирование

Рекомендуется JSON structured logging.

Пример:

```json
{
  "level": "info",
  "time": "2026-08-20T12:50:00Z",
  "request_id": "req-123",
  "user_id": "user-uuid",
  "method": "POST",
  "path": "/api/v1/workout-sessions/session-uuid/complete",
  "status": 200,
  "duration_ms": 35
}
```

Не логировать:

```text
JWT
Telegram Bot Token
полный Telegram initData
password hashes
секреты
```

---

# 22. Healthcheck

## GET `/health`

Response:

```json
{
  "status": "ok"
}
```

## GET `/ready`

Response:

```json
{
  "status": "ok",
  "database": "ok"
}
```

---

# 23. MVP порядок реализации

Рекомендуемый порядок:

```text
1. Bootstrap проекта
2. PostgreSQL + migrations
3. Telegram Auth + JWT
4. Users
5. Habits + completions + streak
6. Workouts + workout sessions
7. Points / XP
8. Leaderboard
9. Daily challenge
10. Weekly challenge
11. WebSocket live updates
12. Telegram notifications
```

---

# 24. Главный backend flow

```text
                    Telegram Mini App
                            │
                            ▼
                    REST / WebSocket
                            │
                            ▼
                     ┌─────────────┐
                     │   Go API    │
                     └──────┬──────┘
                            │
             ┌──────────────┼──────────────┐
             ▼              ▼              ▼
          Habits         Workouts      Challenges
             │              │              │
             ▼              ▼              ▼
      HabitCompleted  WorkoutCompleted ChallengeCompleted
             │              │              │
             └──────────────┬┴──────────────┘
                            ▼
                         Points
                            │
                            ▼
                       Leaderboard
                            │
                            ▼
                       PostgreSQL
```

---

# 25. Что можно добавить после MVP

Возможные следующие фичи:

```text
Achievements / badges
Friends
Private leaderboards
Teams
Push/Telegram notifications
Workout plans
Habit reminders
Premium subscription
Admin panel
Challenge moderation
Anti-cheat
Redis leaderboard
Background jobs
Analytics
```

---

# 26. Итог

Для первой версии достаточно следующей архитектуры:

```text
Telegram Mini App
       │
       ▼
   Go REST API
       │
       ├── Auth
       ├── Users
       ├── Habits
       ├── Workouts
       ├── Challenges
       ├── Points
       └── Leaderboard
       │
       ▼
   PostgreSQL 18
```

WebSocket добавляется для live-обновлений, но не является обязательным для основного CRUD/API.

Главная связь приложения:

```text
User Action
    │
    ▼
Domain Event
    │
    ├── Progress
    ├── Streak
    ├── Challenge
    └── XP
         │
         ▼
    Leaderboard
```

Такой подход позволяет быстро сделать MVP и позже масштабировать отдельные модули без полной переделки backend.
